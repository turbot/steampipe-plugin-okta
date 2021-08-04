package okta

import (
	"context"
	"fmt"
	"strings"

	"github.com/ettle/strcase"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableOktaUser() *plugin.Table {
	return &plugin.Table{
		Name:        "okta_user",
		Description: "Represents an Okta user account.",
		Get: &plugin.GetConfig{
			Hydrate:           getOktaUser,
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"Not found"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listOktaUsers,
			KeyColumns: plugin.KeyColumnSlice{
				// https://developer.okta.com/docs/reference/api/users/#list-users-with-a-filter
				// https://developer.okta.com/docs/reference/api-overview/#filter
				// Key fields
				{Name: "id", Require: plugin.Optional},
				{Name: "login", Require: plugin.Optional},
				{Name: "email", Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
				{Name: "filter", Require: plugin.Optional},
				{Name: "last_updated", Operators: []string{">", ">=", "=", "<", "<="}, Require: plugin.Optional},
			},
		},

		Columns: []*plugin.Column{
			// top columns
			{Name: "login", Type: proto.ColumnType_STRING, Transform: transform.From(userProfile), Description: "Unique identifier for the user (username)."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique key for user."},
			{Name: "email", Type: proto.ColumnType_STRING, Transform: transform.From(userProfile), Description: "Primary email address of user."},
			{Name: "created", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when user was created."},
			{Name: "filter", Type: proto.ColumnType_STRING, Transform: transform.FromQual("filter"), Description: "Filter string to [filter](https://developer.okta.com/docs/reference/api/users/#list-users-with-a-filter) users. Input filter query should not be encoded."},

			// other columns
			{Name: "activated", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when transition to ACTIVE status completed."},
			{Name: "last_login", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of last login."},
			{Name: "last_updated", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when user was last updated."},
			{Name: "password_changed", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when password last changed."},
			{Name: "self_link", Type: proto.ColumnType_STRING, Transform: transform.FromField("Links.self.href"), Description: "A self-referential link to this user."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Current status of user. Can be one of the STAGED, PROVISIONED, ACTIVE, RECOVERY, LOCKED_OUT, PASSWORD_EXPIRED, SUSPENDED, or DEPROVISIONED."},
			{Name: "status_changed", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when status last changed."},
			{Name: "transitioning_to_status", Type: proto.ColumnType_STRING, Description: "Target status of an in-progress asynchronous status transition."},

			// JSON columns
			{Name: "profile", Type: proto.ColumnType_JSON, Description: "User profile properties."},
			{Name: "type", Type: proto.ColumnType_JSON, Description: "User type that determines the schema for the user's profile."},
			{Name: "user_groups", Type: proto.ColumnType_JSON, Hydrate: listUserGroups, Transform: transform.From(transformUserGroups), Description: "List of groups of which the user is a member."},
		},
	}
}

//// LIST FUNCTION

func listOktaUsers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := Connect(ctx, d)
	if err != nil {
		return nil, err
	}

	input := query.Params{}
	equalQuals := d.KeyColumnQuals
	quals := d.Quals

	var queryFilter string
	filter := buildQueryFilter(equalQuals)

	// TODO - optimize or move it to a utility function
	// https://developer.okta.com/docs/reference/api-overview/#operators
	if quals["last_updated"] != nil {
		for _, q := range quals["last_updated"].Quals {
			timeString := q.Value.GetTimestampValue().AsTime().Format(filterTimeFormat)
			filter = append(filter, fmt.Sprintf("%s %s \"%s\"", "lastUpdated", operatorsMap[q.Operator], timeString))
		}
	}

	if equalQuals["filter"] != nil {
		queryFilter = equalQuals["filter"].GetStringValue()
	}

	if queryFilter != "" {
		input.Filter = queryFilter
	} else if len(filter) > 0 {
		input.Filter = strings.Join(filter, " and ")
	}

	if input.Filter != "" {
		plugin.Logger(ctx).Error("Filter", "input.Filter", input.Filter)
	}

	users, resp, err := client.User.ListUsers(ctx, &input)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		d.StreamListItem(ctx, user)
	}

	// paging
	for resp.HasNextPage() {
		var nextUserSet []*okta.User
		resp, err = resp.Next(ctx, &nextUserSet)
		if err != nil {
			return nil, err
		}
		for _, user := range nextUserSet {
			d.StreamListItem(ctx, user)
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getOktaUser(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Debug("getOktaUser")
	var userId string
	if h.Item != nil {
		userId = h.Item.(*okta.User).Id
	} else {
		userId = d.KeyColumnQuals["id"].GetStringValue()
	}

	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("getOktaUser", "connect", err)
		return nil, err
	}

	user, _, err := client.User.GetUser(ctx, userId)
	if err != nil {
		logger.Error("getOktaUser", "get user", err)
		return nil, err
	}

	return user, nil
}

func listUserGroups(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user := h.Item.(*okta.User)
	plugin.Logger(ctx).Debug("listUserGroups")
	client, err := Connect(ctx, d)
	if err != nil {
		return nil, err
	}

	groups, resp, err := client.User.ListUserGroups(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	for resp.HasNextPage() {
		var nextGroupSet []*okta.Group
		resp, err = resp.Next(ctx, &nextGroupSet)
		if err != nil {
			return nil, err
		}
		groups = append(groups, nextGroupSet...)
	}

	return groups, nil
}

//// TRANSFORM FUNCTION

func userProfile(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	user := d.HydrateItem.(*okta.User)
	if user.Profile == nil {
		return nil, nil
	}
	userProfile := *user.Profile

	return userProfile[strcase.ToCamel(d.ColumnName)], nil
}

func transformUserGroups(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	groups := d.HydrateItem.([]*okta.Group)
	var groupsData = []map[string]string{}

	for _, group := range groups {
		groupsData = append(groupsData, map[string]string{
			"id":   group.Id,
			"name": group.Profile.Name,
			"type": group.Type,
		})
	}

	return groupsData, nil
}

//// other useful functions

func buildQueryFilter(equalQuals plugin.KeyColumnEqualsQualMap) []string {
	filters := []string{}

	filterQuals := map[string]string{
		"id":     "id",
		"email":  "profile.email",
		"login":  "profile.login",
		"status": "status",
	}

	for qual, filterColumn := range filterQuals {
		if equalQuals[qual] != nil {
			filters = append(filters, fmt.Sprintf("%s eq \"%s\"", filterColumn, equalQuals[qual].GetStringValue()))
		}
	}

	return filters
}
