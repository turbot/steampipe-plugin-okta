package okta

import (
	"context"
	"fmt"
	"strings"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableOktaGroup() *plugin.Table {
	return &plugin.Table{
		Name:        "okta_group",
		Description: "A Group is made up of users. Groups are useful for representing roles, relationships, and can even be used for subscription tiers.",
		Get: &plugin.GetConfig{
			Hydrate:           getOktaGroup,
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"Not found"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listOktaGroups,
			KeyColumns: plugin.KeyColumnSlice{
				// Key fields
				{Name: "id", Require: plugin.Optional},
				{Name: "type", Require: plugin.Optional},
				{Name: "filter", Require: plugin.Optional},
				{Name: "last_updated", Operators: []string{">", ">=", "=", "<", "<="}, Require: plugin.Optional},
				{Name: "last_membership_updated", Operators: []string{">", ">=", "=", "<", "<="}, Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Top Columns
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Profile.Name"), Description: "Name of the Group."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique key for Group."},
			{Name: "description", Type: proto.ColumnType_STRING, Transform: transform.FromField("Profile.Description"), Description: "Description of the Group."},
			{Name: "created", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when Group was created."},

			// Other Columns
			{Name: "filter", Type: proto.ColumnType_STRING, Transform: transform.FromQual("filter"), Description: "Filter string to [filter](https://developer.okta.com/docs/reference/api/users/#list-users-with-a-filter) users. Input filter query should not be encoded."},
			{Name: "last_membership_updated", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when Group's memberships were last updated."},
			{Name: "last_updated", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when Group's profile was last updated."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Determines how a Group's Profile and memberships are managed. Can be one of OKTA_GROUP, APP_GROUP or BUILT_IN."},

			// JSON Columns
			{Name: "profile", Type: proto.ColumnType_JSON, Description: "The Group's Profile properties."},
			{Name: "object_class", Type: proto.ColumnType_JSON, Description: "Determines the Group's profile."},
			{Name: "group_members", Type: proto.ColumnType_JSON, Hydrate: listGroupMembers, Transform: transform.From(transformGroupMembers), Description: "List of all users that are a member of this Group."},

			// Steampipe Columns
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: titleDescription},
		},
	}
}

//// LIST FUNCTION

func listOktaGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("listOktaGroups", "connect", err)
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
		plugin.Logger(ctx).Debug("Filter", "input.Filter", input.Filter)
	}

	groups, resp, err := client.Group.ListGroups(ctx, &input)
	if err != nil {
		logger.Error("listOktaGroups", "list groups", err)
		return nil, err
	}

	for _, user := range groups {
		d.StreamListItem(ctx, user)
	}

	// paging
	for resp.HasNextPage() {
		var nextGroupSet []*okta.Group
		resp, err = resp.Next(ctx, &nextGroupSet)
		if err != nil {
			logger.Error("listOktaGroups", "list group paging", err)
			return nil, err
		}
		for _, user := range nextGroupSet {
			d.StreamListItem(ctx, user)
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getOktaGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Debug("getOktaGroup")

	var groupId string
	if h.Item != nil {
		groupId = h.Item.(*okta.Group).Id
	} else {
		groupId = d.KeyColumnQuals["id"].GetStringValue()
	}

	if groupId == "" {
		return nil, nil
	}

	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("getOktaGroup", "connect", err)
		return nil, err
	}

	group, _, err := client.Group.GetGroup(ctx, groupId)
	if err != nil {
		logger.Error("getOktaGroup", "get group", err)
		return nil, err
	}

	return group, nil
}

func listGroupMembers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Debug("listGroupMembers")

	var groupId string
	if h.Item != nil {
		groupId = h.Item.(*okta.Group).Id
	} else {
		groupId = d.KeyColumnQuals["id"].GetStringValue()
	}

	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("listGroupMembers", "connect", err)
		return nil, err
	}

	groupMembers, resp, err := client.Group.ListGroupUsers(ctx, groupId, &query.Params{})
	if err != nil {
		logger.Error("listGroupMembers", "get group", err)
		return nil, err
	}

	// paging
	for resp.HasNextPage() {
		var nextgroupMembersSet []*okta.User
		resp, err = resp.Next(ctx, &groupMembers)
		if err != nil {
			logger.Error("listOktaGroups", "list group paging", err)
			return nil, err
		}
		groupMembers = append(groupMembers, nextgroupMembersSet...)
	}

	return groupMembers, nil
}

//// TRANSFORM FUNCTION

func transformGroupMembers(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	users := d.HydrateItem.([]*okta.User)
	var usersData = []map[string]string{}

	for _, user := range users {
		userProfile := *user.Profile
		usersData = append(usersData, map[string]string{
			"id":    user.Id,
			"email": userProfile["email"].(string),
			"login": userProfile["login"].(string),
		})
	}

	return usersData, nil
}
