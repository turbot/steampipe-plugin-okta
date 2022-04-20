package okta

import (
	"context"
	"strings"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

//// TABLE DEFINITION

func tableOktaUserType() *plugin.Table {
	return &plugin.Table{
		Name:        "okta_user_type",
		Description: "Represents an Okta user account.",
		Get: &plugin.GetConfig{
			Hydrate:           getOktaUserType,
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"Not found"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listOktaUserTypes,
		},

		Columns: []*plugin.Column{
			// Top columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name for the type."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique key for the User Type."},
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "The display name for the type."},

			// Other columns
			{Name: "created", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the User Type was created."},
			{Name: "created_by", Type: proto.ColumnType_STRING, Description: "The user ID of the creator of this type."},
			{Name: "default", Type: proto.ColumnType_BOOL, Description: "Boolean to indicate if this type is the default."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "A human-readable description of the type."},
			{Name: "last_updated", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the User Type was last updated."},
			{Name: "last_updated_by", Type: proto.ColumnType_STRING, Description: "The user ID of the last user to edit this type."},

			// Steampipe Columns
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: titleDescription},
		},
	}
}

//// LIST FUNCTION

func listOktaUserTypes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("listOktaUserTypes", "connect_error", err)
		return nil, err
	}

	userTypes, resp, err := client.UserType.ListUserTypes(ctx)
	if err != nil {
		logger.Error("listOktaUserTypes", "list_user_types_error", err)
		if strings.Contains(err.Error(), "Not found") {
			return nil, nil
		}
		return nil, err
	}

	for _, userType := range userTypes {
		d.StreamListItem(ctx, userType)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	// paging
	for resp.HasNextPage() {
		var nextUserTypeSet []*okta.UserType
		resp, err = resp.Next(ctx, &nextUserTypeSet)
		if err != nil {
			logger.Error("listOktaUserTypes", "list_user_types_paging_error", err)
			return nil, err
		}
		for _, user := range nextUserTypeSet {
			d.StreamListItem(ctx, user)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getOktaUserType(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getOktaUserType")

	userTypeId := d.KeyColumnQuals["id"].GetStringValue()

	if userTypeId == "" {
		return nil, nil
	}

	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("getOktaUserType", "connect_error", err)
		return nil, err
	}

	userType, _, err := client.UserType.GetUserType(ctx, userTypeId)
	if err != nil {
		logger.Error("getOktaUserType", "get_user_type_error", err)
		return nil, err
	}

	return userType, nil
}
