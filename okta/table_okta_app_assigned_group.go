package okta

import (
	"context"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableOktaApplicationAssignedGroup() *plugin.Table {
	return &plugin.Table{
		Name:        "okta_app_assigned_group",
		Description: "Represents an application group assignment.",
		Get: &plugin.GetConfig{
			Hydrate:           getApplicationAssignedGroup,
			KeyColumns:        plugin.AllColumns([]string{"id", "app_id"}),
		},
		List: &plugin.ListConfig{
			ParentHydrate: listOktaApplications,
			Hydrate: listApplicationAssignedGroups,
		},

		Columns: []*plugin.Column{
			// Top Columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique key for the group."},
			{Name: "app_id", Type: proto.ColumnType_STRING, Description: "Unique key for the application."},

			// Other Columns
			{Name: "last_updated", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when group was last updated."},
			{Name: "priority", Type: proto.ColumnType_INT, Description: "Priority of the group"},

			// JSON Columns
			{Name: "links", Type: proto.ColumnType_JSON, Description: "The link details of the group."},
			{Name: "profile", Type: proto.ColumnType_JSON, Description: "The profile details of the group."},

			// Steampipe Columns
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: titleDescription},
		},
	}
}

type AppGroupInfo struct {
	AppId   string
	okta.ApplicationGroupAssignment
}

//// LIST FUNCTION

func listApplicationAssignedGroups(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listApplicationAssignedGroups")
	var appId string

	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("listApplicationAssignedGroups", "connect", err)
		return nil, err
	}

	if h.Item != nil {
		appId = h.Item.(*okta.Application).Id
	} else {
		appId = d.KeyColumnQuals["app_id"].GetStringValue()
	}
	
	input := query.Params{
		Limit: 200,
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < input.Limit {
			input.Limit = *limit
		}
	}

	groups, resp, err := client.Application.ListApplicationGroupAssignments(ctx, appId, &input)
	
	if err != nil {
		logger.Error("listApplicationAssignedGroups", "error_ListApplicationGroupAssignments", err)
		return nil, err
	}

	for _, group := range groups {
		d.StreamListItem(ctx, AppGroupInfo{appId, *group})
	}

	// paging
	for resp.HasNextPage() {
		var nextGroupSet []*okta.ApplicationGroupAssignment
		resp, err = resp.Next(ctx, &nextGroupSet)
		if err != nil {
			logger.Error("listApplicationAssignedGroups", "error_ListApplicationGroupAssignments_paging", err)
			return nil, err
		}
		for _, group := range nextGroupSet {
			d.StreamListItem(ctx, AppGroupInfo{appId, *group})
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTION

func getApplicationAssignedGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Debug("getApplicationAssignedGroup")
	appId := d.KeyColumnQuals["app_id"].GetStringValue()
	groupId := d.KeyColumnQuals["id"].GetStringValue()

	if appId == "" || groupId == "" {
		return nil, nil
	}

	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("getApplicationAssignedGroup", "connect", err)
		return nil, err
	}

	group, _, err := client.Application.GetApplicationGroupAssignment(ctx, appId, groupId, &query.Params{})
	if err != nil {
		logger.Error("getApplicationAssignedGroup", "error_GetApplicationGroupAssignment", err)
		return nil, err
	}

	return AppGroupInfo{appId, *group}, nil
}
