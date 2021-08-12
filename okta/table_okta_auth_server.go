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

func tableOktaAuthServer() *plugin.Table {
	return &plugin.Table{
		Name:        "okta_auth_server",
		Description: "Represents an Okta user account.",
		Get: &plugin.GetConfig{
			Hydrate:           getOktaAuthorizationServer,
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"Not found"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listOktaAuthorizationServers,
		},

		Columns: []*plugin.Column{
			// Top columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of authorization server."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique key for the authorization server."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Indicates whether a custom authorization server is ACTIVE or INACTIVE."},

			// Other columns
			{Name: "created", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the User Type was created."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of a custom authorization server."},
			{Name: "issuer", Type: proto.ColumnType_STRING, Description: "The complete URL for a custom authorization server. This becomes the iss claim in an access token. issuerMode is visible if you have the custom URL Domain feature enabled."},
			{Name: "issuer_mode", Type: proto.ColumnType_STRING, Description: "Indicates which value is specified in the issuer of the tokens that a custom authorization server returns: the original Okta org domain URL or a custom domain URL."},
			{Name: "last_updated", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when authorization server was last updated."},

			{Name: "audiences", Type: proto.ColumnType_JSON, Description: "The recipients that the tokens are intended for. This becomes the aud claim in an access token. Currently, Okta supports only one audience."},
			{Name: "credentials", Type: proto.ColumnType_JSON, Description: "Keys and settings used to sign tokens."},

			// Steampipe Columns
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: titleDescription},
		},
	}
}

//// LIST FUNCTION

func listOktaAuthorizationServers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("listOktaAuthorizationServers", "connect", err)
		return nil, err
	}

	authorizationServers, resp, err := client.AuthorizationServer.ListAuthorizationServers(ctx, &query.Params{})
	if err != nil {
		logger.Error("listOktaAuthorizationServers", "list auth servers", err)
		return nil, err
	}

	for _, authorizationServer := range authorizationServers {
		d.StreamListItem(ctx, authorizationServer)
	}

	// paging
	for resp.HasNextPage() {
		var nextAuthorizationServerSet []*okta.AuthorizationServer
		resp, err = resp.Next(ctx, &nextAuthorizationServerSet)
		if err != nil {
			logger.Error("listOktaAuthorizationServers", "list paging", err)
			return nil, err
		}
		for _, authorizationServer := range nextAuthorizationServerSet {
			d.StreamListItem(ctx, authorizationServer)
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getOktaAuthorizationServer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Debug("getOktaAuthorizationServer")

	id := d.KeyColumnQuals["id"].GetStringValue()

	if id == "" {
		return nil, nil
	}

	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("getOktaAuthorizationServer", "connect.error", err)
		return nil, err
	}

	authorizationServer, _, err := client.AuthorizationServer.GetAuthorizationServer(ctx, id)
	if err != nil {
		logger.Error("getOktaAuthorizationServer", "get_authorization_server.error", err)
		return nil, err
	}

	return authorizationServer, nil
}
