package okta

import (
	"context"
	"strings"

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
			ShouldIgnoreError: isNotFoundError([]string{"Not found"}),
		},
		List: &plugin.ListConfig{
			ParentHydrate: getOrListOktaApplications,
			Hydrate:       listApplicationAssignedGroups,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "app_id", Require: plugin.Optional},
			},
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
	AppId string
	okta.ApplicationGroupAssignment
}

//// LIST FUNCTION

func listApplicationAssignedGroups(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listApplicationAssignedGroups")
	appId := h.Item.(*okta.Application).Id

	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("listApplicationAssignedGroups", "connect_error", err)
		return nil, err
	}

	// Default maximum limit set as per documentation
	// https://developer.okta.com/docs/reference/api/apps/#list-groups-assigned-to-application
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
		logger.Error("listApplicationAssignedGroups", "list_app_groups_error", err)
		return nil, err
	}

	for _, group := range groups {
		d.StreamListItem(ctx, AppGroupInfo{appId, *group})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	// paging
	for resp.HasNextPage() {
		var nextGroupSet []*okta.ApplicationGroupAssignment
		resp, err = resp.Next(ctx, &nextGroupSet)
		if err != nil {
			logger.Error("listApplicationAssignedGroups", "list_app_groups_paging_error", err)
			return nil, err
		}
		for _, group := range nextGroupSet {
			d.StreamListItem(ctx, AppGroupInfo{appId, *group})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTION

func getApplicationAssignedGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getApplicationAssignedGroup")
	appId := d.KeyColumnQuals["app_id"].GetStringValue()
	groupId := d.KeyColumnQuals["id"].GetStringValue()

	if appId == "" || groupId == "" {
		return nil, nil
	}

	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("getApplicationAssignedGroup", "connect_error", err)
		return nil, err
	}

	group, _, err := client.Application.GetApplicationGroupAssignment(ctx, appId, groupId, &query.Params{})
	if err != nil {
		logger.Error("getApplicationAssignedGroup", "get_app_group_error", err)
		return nil, err
	}

	return AppGroupInfo{appId, *group}, nil
}

//// PARENT HYDRATE FUNCTION

func getOrListOktaApplications(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getOrListOktaApplications")
	appID := d.KeyColumnQuals["app_id"].GetStringValue()

	// List application API doesn't support filtering by app ID, so call the get
	// function to reduce API calls
	if appID != "" {
		// The okta_application table uses the "id" column instead
		d.KeyColumnQuals["id"] = d.KeyColumnQuals["app_id"]
		app, err := getOktaApplication(ctx, d, h)
		if err != nil && !strings.Contains(err.Error(), "Not found") {
			logger.Error("getOrListOktaApplications", "get_application_error", err)
			return nil, err
		}
		d.StreamListItem(ctx, app)
		return nil, nil
	}

	_, err := listOktaApplications(ctx, d, h)
	if err != nil {
		logger.Error("getOrListOktaApplications", "list_applications_error", err)
		return nil, err
	}

	return nil, nil
}
