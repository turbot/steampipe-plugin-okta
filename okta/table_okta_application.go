package okta

import (
	"context"
	"strings"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableOktaApplication() *plugin.Table {
	return &plugin.Table{
		Name:        "okta_application",
		Description: "An Application holds information about the protocol in which it wants Okta to communicate, policies for accessing the application, and which users can use the application after identifying themselves.",
		Get: &plugin.GetConfig{
			Hydrate:           getOktaApplication,
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"Not found"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listOktaApplications,
			KeyColumns: plugin.KeyColumnSlice{
				// https://developer.okta.com/docs/reference/api/apps/#filters
				{Name: "name", Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
				{Name: "filter", Require: plugin.Optional},
			},
		},

		Columns: []*plugin.Column{
			// Top Columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Unique key for app definition."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique key for app."},
			{Name: "label", Type: proto.ColumnType_STRING, Description: "User-defined display name for app."},
			{Name: "created", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when user was created."},
			{Name: "filter", Type: proto.ColumnType_STRING, Transform: transform.FromQual("filter"), Description: "Filter string to [filter](https://developer.okta.com/docs/reference/api/users/#list-users-with-a-filter) users. Input filter query should not be encoded."},

			// Other Columns
			{Name: "last_updated", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when app was last updated."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Current status of app. Valid values are ACTIVE or INACTIVE."},
			{Name: "sign_on_mode", Type: proto.ColumnType_STRING, Description: "Authentication mode of app. Can be one of AUTO_LOGIN, BASIC_AUTH, BOOKMARK, BROWSER_PLUGIN, Custom, OPENID_CONNECT, SAML_1_1, SAML_2_0, SECURE_PASSWORD_STORE and WS_FEDERATION."},

			// JSON Columns
			{Name: "settings", Type: proto.ColumnType_JSON, Description: "Settings for app."},
			{Name: "visibility", Type: proto.ColumnType_JSON, Description: "Visibility settings for app."},
			{Name: "credentials", Type: proto.ColumnType_JSON, Description: "Credentials for the specified signOnMode."},
			{Name: "accessibility", Type: proto.ColumnType_JSON, Description: "Access settings for app."},

			// Steampipe Columns
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: titleDescription},
		},
	}
}

//// LIST FUNCTION

func listOktaApplications(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("listOktaApplications", "connect_error", err)
		return nil, err
	}

	// Default maximum limit set as per documentation
	// https://developer.okta.com/docs/reference/api/apps/#list-applications
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

	equalQuals := d.KeyColumnQuals
	filter := buildQueryFilter(equalQuals, []string{"name", "status"})
	var queryFilter string

	if equalQuals["filter"] != nil {
		queryFilter = equalQuals["filter"].GetStringValue()
	}

	if queryFilter != "" {
		input.Filter = queryFilter
	} else if len(filter) > 0 {
		input.Filter = strings.Join(filter, " and ")
	}

	applications, resp, err := client.Application.ListApplications(ctx, &input)
	if err != nil {
		logger.Error("listOktaApplications", "list_applications_error", err)
		if strings.Contains(err.Error(), "Not found") {
			return nil, nil
		}
		return nil, err
	}

	for _, app := range applications {
		d.StreamListItem(ctx, app)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	// paging
	for resp.HasNextPage() {
		var nextApplicationSet []*okta.Application
		resp, err = resp.Next(ctx, &nextApplicationSet)
		if err != nil {
			logger.Error("listOktaApplications", "list_applications_paging_error", err)
			return nil, err
		}
		for _, app := range nextApplicationSet {
			d.StreamListItem(ctx, app)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTION

func getOktaApplication(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getOktaApplication")
	appId := d.KeyColumnQuals["id"].GetStringValue()

	// Empty check for appId
	if appId == "" {
		return nil, nil
	}

	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("getOktaApplication", "connect_error", err)
		return nil, err
	}

	app, _, err := client.Application.GetApplication(ctx, appId, okta.NewApplication(), &query.Params{})
	if err != nil {
		logger.Error("getOktaApplication", "get_application_error", err)
		return nil, err
	}

	return app, nil
}
