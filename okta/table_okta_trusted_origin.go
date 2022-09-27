package okta

import (
	"context"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableOktaTrustedOrigin() *plugin.Table {
	return &plugin.Table{
		Name:        "okta_trusted_origin",
		Description: "Trusted Origin is a security-based concept that combines the URI scheme, hostname, and port number of a page.",
		Get: &plugin.GetConfig{
			Hydrate:           getOktaTrustedOrigin,
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"Not found"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listOktaTrustedOrigins,
		},
		Columns: []*plugin.Column{
			// Top Columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the trusted origin."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "A unique key for the trusted origin."},

			// Other Columns
			{Name: "created", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the trusted origin was created."},
			{Name: "created_by", Type: proto.ColumnType_STRING, Description: "The ID of the user who created the trusted origin."},
			{Name: "last_updated", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the trusted origin was last updated."},
			{Name: "last_updated_by", Type: proto.ColumnType_STRING, Description: "The ID of the user who last updated the trusted origin."},
			{Name: "origin", Type: proto.ColumnType_STRING, Description: "The origin of the trusted origin."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Current status of the trusted origin. Valid values are 'ACTIVE' or 'INACTIVE'."},

			// JSON Columns
			{Name: "scopes", Type: proto.ColumnType_JSON, Description: "The scopes for the trusted origin. Valid values are 'CORS' or 'REDIRECT'."},

			// Steampipe Columns
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: titleDescription},
		},
	}
}

//// LIST FUNCTION

func listOktaTrustedOrigins(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("listOktaTrustedOrigins", "connect_error", err)
		return nil, err
	}

	// Maximum limit isn't mentioned in the documentation
	// Default maximum limit is set as 200
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

	origins, resp, err := client.TrustedOrigin.ListOrigins(ctx, &input)
	if err != nil {
		logger.Error("listOktaTrustedOrigins", "list_origins_error", err)
		return nil, err
	}

	for _, origin := range origins {
		d.StreamListItem(ctx, origin)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	// paging
	for resp.HasNextPage() {
		var nextOriginSet []*okta.TrustedOrigin
		resp, err = resp.Next(ctx, &nextOriginSet)
		if err != nil {
			logger.Error("listOktaTrustedOrigins", "list_origins_paging_error", err)
			return nil, err
		}
		for _, origin := range nextOriginSet {
			d.StreamListItem(ctx, origin)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getOktaTrustedOrigin(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getOktaTrustedOrigin")

	trustedOriginId := d.KeyColumnQuals["id"].GetStringValue()

	if trustedOriginId == "" {
		return nil, nil
	}

	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("getOktaTrustedOrigin", "connect_error", err)
		return nil, err
	}

	app, _, err := client.TrustedOrigin.GetOrigin(ctx, trustedOriginId)
	if err != nil {
		logger.Error("getOktaTrustedOrigin", "get_origin_error", err)
		return nil, err
	}

	return app, nil
}
