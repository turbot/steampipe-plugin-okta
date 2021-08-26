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
			{Name: "links", Type: proto.ColumnType_JSON, Description: "The authorization server link properties."},

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
		logger.Error("listOktaAuthServers", "connect", err)
		return nil, err
	}

	input := query.Params{}
	servers, resp, err := client.AuthorizationServer.ListAuthorizationServers(ctx, &input)
	if err != nil {
		logger.Error("listOktaAuthServers", "list auth servers", err)
		if strings.Contains(err.Error(), "Not found") {
			return nil, nil
		}
		return nil, err
	}

	for _, server := range servers {
		d.StreamListItem(ctx, server)
	}

	// paging
	for resp.HasNextPage() {
		var nextAuthorizationServerSet []*okta.AuthorizationServer
		resp, err = resp.Next(ctx, &nextAuthorizationServerSet)
		if err != nil {
			logger.Error("listOktaAuthServers", "list auth servers paging", err)
			return nil, err
		}
		for _, server := range nextAuthorizationServerSet {
			d.StreamListItem(ctx, server)
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getOktaAuthServer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Debug("getOktaAuthServer")

	authServerId := d.KeyColumnQuals["id"].GetStringValue()

	if authServerId == "" {
		return nil, nil
	}

	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("getOktaAuthServer", "connect", err)
		return nil, err
	}

	server, _, err := client.AuthorizationServer.GetAuthorizationServer(ctx, authServerId)
	if err != nil {
		logger.Error("getOktaAuthServer", "get auth server", err)
		return nil, err
	}

	return server, nil
}
