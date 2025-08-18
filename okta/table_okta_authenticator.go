package okta

import (
	"context"

	oktaV5 "github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

// tableOktaAuthenticator returns the table definition for okta_authenticator.
func tableOktaAuthenticator() *plugin.Table {
	return &plugin.Table{
		Name:        "okta_authenticator",
		Description: "Represents an Okta Authenticator configured in the organization.",
		Get: &plugin.GetConfig{
			Hydrate:           getOktaAuthenticator,
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"Not found"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listOktaAuthenticators,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Display name of the authenticator."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique identifier for the authenticator."},

			// Other columns
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Type of the authenticator (e.g., email, password, phone, security_key, security_question, etc.)."},
			{Name: "key", Type: proto.ColumnType_STRING, Description: "Key for the authenticator (e.g., okta_email, okta_password, phone_number)."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Lifecycle status of the authenticator (ACTIVE or INACTIVE)."},
			{Name: "created", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the authenticator was created."},
			{Name: "last_updated", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the authenticator was last updated."},

			// JSON columns
			{Name: "settings", Type: proto.ColumnType_JSON, Description: "Settings for the authenticator."},
			{Name: "provider", Type: proto.ColumnType_JSON, Description: "Provider configuration for app-type authenticators."},
			{Name: "links", Type: proto.ColumnType_JSON, Description: "HAL links related to the authenticator."},

			// Steampipe standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: titleDescription},
		}),
	}
}

//// LIST FUNCTION

func listOktaAuthenticators(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client, err := ConnectV5(ctx, d)
	if err != nil {
		logger.Error("okta_authenticator.listOktaAuthenticators", "connect_error", err)
		return nil, err
	}

	req := client.AuthenticatorAPI.ListAuthenticators(ctx)
	authenticators, resp, err := req.Execute()
	if err != nil {
		logger.Error("okta_authenticator.listOktaAuthenticators", "api_error", err)
		return nil, err
	}

	for _, item := range authenticators {
		if v := item.GetActualInstance(); v != nil {
			d.StreamListItem(ctx, v)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	// pagination
	for resp.HasNextPage() {
		var nextSet []oktaV5.ListAuthenticators200ResponseInner
		resp, err = resp.Next(&nextSet)
		if err != nil {
			logger.Error("okta_authenticator.listOktaAuthenticators", "paging_error", err)
			return nil, err
		}
		for _, item := range nextSet {
			if v := item.GetActualInstance(); v != nil {
				d.StreamListItem(ctx, v)

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}
	}

	return nil, nil
}

//// GET FUNCTION

func getOktaAuthenticator(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	id := d.EqualsQualString("id")
	if id == "" {
		return nil, nil
	}

	client, err := ConnectV5(ctx, d)
	if err != nil {
		logger.Error("okta_authenticator.getOktaAuthenticator", "connect_error", err)
		return nil, err
	}

	req := client.AuthenticatorAPI.GetAuthenticator(ctx, id)
	result, _, err := req.Execute()
	if err != nil {
		logger.Error("okta_authenticator.getOktaAuthenticator", "api_error", err)
		return nil, err
	}

	if result != nil {
		return *result, nil
	}
	return nil, nil
}
