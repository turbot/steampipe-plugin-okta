package okta

import (
	"context"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableOktaDomain() *plugin.Table {
	return &plugin.Table{
		Name:        "okta_domain",
		Description: "Represents an Okta custom domain.",
		// Get: &plugin.GetConfig{
		// 	Hydrate:           getOktaAuthorizationServer,
		// 	KeyColumns:        plugin.SingleColumn("id"),
		// 	ShouldIgnoreError: isNotFoundError([]string{"Not found"}),
		// },
		List: &plugin.ListConfig{
			Hydrate: listDomians,
		},

		Columns: []*plugin.Column{
			// Top columns
			{Name: "domain", Type: proto.ColumnType_STRING, Description: "The domain name."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique key for the custom domain."},
			{Name: "certificate_source_type", Type: proto.ColumnType_STRING, Transform: transform.FromField("CertificateSourcetype"), Description: ""},
			{Name: "validation_status", Type: proto.ColumnType_STRING, Description: ""},

			{Name: "dns_records", Type: proto.ColumnType_JSON, Description: ""},
			{Name: "public_certificate", Type: proto.ColumnType_JSON, Description: ""},

			// Steampipe Columns
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Domain"), Description: titleDescription},
		},
	}
}

//// LIST FUNCTION

func listDomians(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("listOktaAuthorizationServers", "connect", err)
		return nil, err
	}

	domains, resp, err := client.Domain.ListDomains(ctx)
	if err != nil {
		logger.Error("listOktaAuthorizationServers", "list auth servers", err)
		return nil, err
	}

	for _, domain := range domains.Domains {
		d.StreamListItem(ctx, domain)
	}

	// paging
	for resp.HasNextPage() {
		var nextDomainListSet []*okta.DomainListResponse
		resp, err = resp.Next(ctx, &nextDomainListSet)
		if err != nil {
			logger.Error("listOktaAuthorizationServers", "list paging", err)
			return nil, err
		}
		for _, domainListResponse := range nextDomainListSet {
			for _, domain := range domainListResponse.Domains {
				d.StreamListItem(ctx, domain)
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

// func getOktaAuthorizationServer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 	logger := plugin.Logger(ctx)
// 	logger.Debug("getOktaAuthorizationServer")

// 	id := d.KeyColumnQuals["id"].GetStringValue()

// 	if id == "" {
// 		return nil, nil
// 	}

// 	client, err := Connect(ctx, d)
// 	if err != nil {
// 		logger.Error("getOktaAuthorizationServer", "connect.error", err)
// 		return nil, err
// 	}

// 	authorizationServer, _, err := client.AuthorizationServer.GetAuthorizationServer(ctx, id)
// 	if err != nil {
// 		logger.Error("getOktaAuthorizationServer", "get_authorization_server.error", err)
// 		return nil, err
// 	}

// 	return authorizationServer, nil
// }
