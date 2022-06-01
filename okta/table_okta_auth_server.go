package okta

import (
	"context"
	"strings"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

//// TABLE DEFINITION

func tableOktaAuthServer() *plugin.Table {
	return &plugin.Table{
		Name:        "okta_auth_server",
		Description: "Represents an Okta Authorization Server.",
		Get: &plugin.GetConfig{
			Hydrate:           getOktaAuthServer,
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"Not found"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listOktaAuthServers,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "name", Require: plugin.Optional},
			},
		},

		Columns: []*plugin.Column{
			// Top columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name for the authorization server."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique key for the authorization server."},

			// Other columns
			{Name: "created", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the authorization server was created."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "A human-readable description of the authorization server."},
			{Name: "issuer", Type: proto.ColumnType_STRING, Description: "The issuer URI of the authorization server."},
			{Name: "issuer_mode", Type: proto.ColumnType_STRING, Description: "The issuer mode of the authorization server."},
			{Name: "last_updated", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the authorization server was last updated."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of the authorization server."},

			// JSON Columns
			{Name: "audiences", Type: proto.ColumnType_JSON, Description: "The audiences of the authorization server."},
			{Name: "credentials", Type: proto.ColumnType_JSON, Description: "The authorization server credentials."},

			// Steampipe Columns
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: titleDescription},
		},
	}
}

//// LIST FUNCTION

func listOktaAuthServers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("listOktaAuthServers", "connect_error", err)
		return nil, err
	}

	// Default maximum limit set as per documentation
	// https://developer.okta.com/docs/reference/api/authorization-servers/#list-authorization-servers
	input := query.Params{
		Limit: 200,
	}

	if d.KeyColumnQualString("name") != "" {
		input.Q = d.KeyColumnQualString("name")
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < input.Limit {
			input.Limit = *limit
		}
	}

	servers, resp, err := client.AuthorizationServer.ListAuthorizationServers(ctx, &input)
	if err != nil {
		logger.Error("listOktaAuthServers", "list_auth_servers_error", err)
		if strings.Contains(err.Error(), "Not found") {
			return nil, nil
		}
		return nil, err
	}

	for _, server := range servers {
		d.StreamListItem(ctx, server)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	// paging
	for resp.HasNextPage() {
		var nextAuthorizationServerSet []*okta.AuthorizationServer
		resp, err = resp.Next(ctx, &nextAuthorizationServerSet)
		if err != nil {
			logger.Error("listOktaAuthServers", "list_auth_servers_paging_error", err)
			return nil, err
		}
		for _, server := range nextAuthorizationServerSet {
			d.StreamListItem(ctx, server)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getOktaAuthServer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getOktaAuthServer")

	authServerId := d.KeyColumnQuals["id"].GetStringValue()

	if authServerId == "" {
		return nil, nil
	}

	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("getOktaAuthServer", "connect_error", err)
		return nil, err
	}

	server, _, err := client.AuthorizationServer.GetAuthorizationServer(ctx, authServerId)
	if err != nil {
		logger.Error("getOktaAuthServer", "get_auth_servers_error", err)
		return nil, err
	}

	return server, nil
}
