package okta

import (
	"context"

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
		List: &plugin.ListConfig{
			Hydrate: listOktaUsers,
			// 	KeyColumns: plugin.KeyColumnSlice{
			// 		// Key fields
			// 		{Name: "id", Require: plugin.Optional},
			// 		{Name: "user_principal_name", Require: plugin.Optional},
			// 		{Name: "filter", Require: plugin.Optional},

			// 		// Other fields for filtering OData
			// 		{Name: "user_type", Require: plugin.Optional},
			// 		{Name: "account_enabled", Require: plugin.Optional, Operators: []string{"<>", "="}},
			// 		{Name: "display_name", Require: plugin.Optional},
			// 		{Name: "surname", Require: plugin.Optional},
			// 	},
		},

		Columns: []*plugin.Column{
			// top columns
			{Name: "login", Type: proto.ColumnType_STRING, Transform: transform.From(userProfile), Description: "Unique identifier for the user (username)."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique key for user."},
			{Name: "email", Type: proto.ColumnType_STRING, Transform: transform.From(userProfile), Description: "Primary email address of user."},
			{Name: "created", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when user was created."},

			// other columns
			{Name: "activated", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when transition to ACTIVE status completed."},
			{Name: "last_login", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of last login."},
			{Name: "last_updated", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when user was last updated."},
			{Name: "password_changed", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when password last changed."},
			{Name: "self_link", Type: proto.ColumnType_STRING, Transform: transform.FromField("Links.self.href"), Description: "A self-referential link to this user."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Current status of user. Can be one of the STAGED, PROVISIONED, ACTIVE, RECOVERY, LOCKED_OUT, PASSWORD_EXPIRED, SUSPENDED, or DEPROVISIONED."},
			{Name: "status_changed", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when status last changed."},
			{Name: "transitioning_to_status", Type: proto.ColumnType_STRING, Description: "Target status of an in-progress asynchronous status transition."},

			// JSOn columns
			{Name: "profile", Type: proto.ColumnType_JSON, Description: "User profile properties."},
			{Name: "type", Type: proto.ColumnType_JSON, Description: "User type that determines the schema for the user's profile."},
		},
	}
}

func userProfile(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	user := d.HydrateItem.(*okta.User)
	if user.Profile == nil {
		return nil, nil
	}
	userProfile := *user.Profile

	return userProfile[strcase.ToCamel(d.ColumnName)], nil
}

//// LIST FUNCTION

func listOktaUsers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := Connect(ctx, d)
	if err != nil {
		return nil, err
	}

	input := query.Params{}
	// if helpers.StringSliceContains(d.QueryContext.Columns, "member_of") {
	// 	input.Expand = odata.Expand{
	// 		Relationship: "memberOf",
	// 		Select:       []string{"id", "displayName"},
	// 	}
	// }

	// equalQuals := d.KeyColumnQuals
	// quals := d.Quals

	// var queryFilter string
	// filter := buildQueryFilter(equalQuals)
	// filter = append(filter, buildBoolNEFilter(quals)...)

	// if equalQuals["filter"] != nil {
	// 	queryFilter = equalQuals["filter"].GetStringValue()
	// }

	// if queryFilter != "" {
	// 	input.Filter = queryFilter
	// } else if len(filter) > 0 {
	// 	input.Filter = strings.Join(filter, " and ")
	// }

	// if input.Filter != "" {
	// 	plugin.Logger(ctx).Debug("Filter", "input.Filter", input.Filter)
	// }

	pagesLeft := true
	for pagesLeft {
		users, _, err := client.User.ListUsers(ctx, &input)
		if err != nil {
			// if isNotFoundError(err) {
			// 	return nil, nil
			// }
			return nil, err
		}

		for _, user := range users {
			d.StreamListItem(ctx, user)
		}
		pagesLeft = false
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

// func getTenantId(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 	plugin.Logger(ctx).Debug("getTenantId")

// 	session, err := GetNewSession(ctx, d)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return session.TenantID, nil
// }

// func buildQueryFilter(equalQuals plugin.KeyColumnEqualsQualMap) []string {
// 	filters := []string{}

// 	filterQuals := map[string]string{
// 		"display_name":             "string",
// 		"id":                       "string",
// 		"surname":                  "string",
// 		"user_principal_name":      "string",
// 		"user_type":                "string",
// 		"account_enabled":          "bool",
// 		"mail_enabled":             "bool",
// 		"security_enabled":         "bool",
// 		"on_premises_sync_enabled": "bool",
// 	}

// 	for qual, qualType := range filterQuals {
// 		switch qualType {
// 		case "string":
// 			if equalQuals[qual] != nil {
// 				filters = append(filters, fmt.Sprintf("%s eq '%s'", strcase.ToCamel(qual), equalQuals[qual].GetStringValue()))
// 			}
// 		case "bool":
// 			if equalQuals[qual] != nil {
// 				filters = append(filters, fmt.Sprintf("%s eq %t", strcase.ToCamel(qual), equalQuals[qual].GetBoolValue()))
// 			}
// 		}
// 	}

// 	return filters
// }

// func buildBoolNEFilter(quals plugin.KeyColumnQualMap) []string {
// 	filters := []string{}

// 	filterQuals := []string{
// 		"account_enabled",
// 		"mail_enabled",
// 		"on_premises_sync_enabled",
// 		"security_enabled",
// 	}

// 	for _, qual := range filterQuals {
// 		if quals[qual] != nil {
// 			for _, q := range quals[qual].Quals {
// 				value := q.Value.GetBoolValue()
// 				if q.Operator == "<>" {
// 					filters = append(filters, fmt.Sprintf("%s eq %t", strcase.ToCamel(qual), !value))
// 					break
// 				}
// 			}
// 		}
// 	}

// 	return filters
// }
