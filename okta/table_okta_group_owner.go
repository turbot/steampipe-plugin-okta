package okta

import (
	"context"
	"time"

	"github.com/okta/okta-sdk-golang/v2/okta"
	oktav4 "github.com/okta/okta-sdk-golang/v4/okta"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableOktaGroupOwner() *plugin.Table {
	return &plugin.Table{
		Name:        "okta_group_owner",
		Description: "An Okta Group owner is a designated individual responsible for managing and overseeing a specific group within the Okta identity and access management platform.",
		List: &plugin.ListConfig{
			Hydrate:       listOktaGroupOwners,
			ParentHydrate: listOktaGroups,
			KeyColumns: plugin.OptionalColumns([]string{"group_id"}),
		},
		Columns: commonColumns([]*plugin.Column{
			{Name: "group_id", Type: proto.ColumnType_STRING, Description: "Unique key for Group."},
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "The display name of the group owner."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The ID of the group owner."},
			{Name: "last_updated", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("LastUpdated").Transform(transform.NullIfZeroValue), Description: "Timestamp when the group owner was last updated."},
			{Name: "origin_id", Type: proto.ColumnType_STRING, Description: "The ID of the app instance if the originType is APPLICATION. This value is NULL if originType is OKTA_DIRECTORY."},
			{Name: "origin_type", Type: proto.ColumnType_STRING, Description: "The source where group ownership is managed."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The entity type of the owner."},
			{Name: "resolved", Type: proto.ColumnType_BOOL, Description: "If originType is APPLICATION, this parameter is set to FALSE until the owner’s originId is reconciled with an associated Okta ID."},

			// JSON Columns
			{Name: "additional_properties", Type: proto.ColumnType_JSON, Description: "The additional properties for the owner."},

			// Steampipe Columns
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("DisplayName"), Description: titleDescription},
		}),
	}
}

type GroupOwner struct {
	GroupId              *string
	DisplayName          *string
	Id                   *string
	LastUpdated          *time.Time
	OriginId             *string
	OriginType           *string
	Resolved             *bool
	Type                 *string
	AdditionalProperties map[string]interface{}
}

//// LIST FUNCTION

func listOktaGroupOwners(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	var groupId string
	if h.Item != nil {
		groupId = h.Item.(*okta.Group).Id
	} else {
		groupId = d.EqualsQuals["group_id"].GetStringValue()
	}

	// Restrict API call based on group_id query parameter.
	if d.EqualsQuals["group_id"] != nil && d.EqualsQuals["group_id"].GetStringValue() != groupId {
		return nil, nil
	}

	client, err := ConnectV4(ctx, d)
	if err != nil {
		logger.Error("okta_group_owner.listGroupOwners", "connect_error", err)
		return nil, err
	}

	groupOwnerReq := client.GroupAPI.ListGroupOwners(ctx, groupId)

	owners, resp, err := groupOwnerReq.Execute()
	if err != nil {
		logger.Error("okta_group_owner.listGroupOwners", "api_error", err)
		return nil, err
	}

	for _, owner := range owners {
		d.StreamListItem(ctx, GroupOwner{
				GroupId:              &groupId,
				DisplayName:          owner.DisplayName,
				Id:                   owner.Id,
				LastUpdated:          owner.LastUpdated,
				OriginId:             owner.OriginId,
				OriginType:           owner.OriginType,
				Resolved:             owner.Resolved,
				Type:                 owner.Type,
				AdditionalProperties: owner.AdditionalProperties,
			})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	// paging
	for resp.HasNextPage() {
		var nextGroupOwners []oktav4.GroupOwner
		resp, err = resp.Next(&nextGroupOwners)
		if err != nil {
			logger.Error("okta_group_owner.listGroupOwners", "api_paging_error", err)
			return nil, err
		}
		for _, owner := range nextGroupOwners {
			d.StreamListItem(ctx, GroupOwner{
				GroupId:              &groupId,
				DisplayName:          owner.DisplayName,
				Id:                   owner.Id,
				LastUpdated:          owner.LastUpdated,
				OriginId:             owner.OriginId,
				OriginType:           owner.OriginType,
				Resolved:             owner.Resolved,
				Type:                 owner.Type,
				AdditionalProperties: owner.AdditionalProperties,
			})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
