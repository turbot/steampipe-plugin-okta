package okta

import (
	"context"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableOktaIdpDiscoveryPolicy() *plugin.Table {
	return &plugin.Table{
		Name:        "okta_idp_discovery_policy",
		Description: "The IdP Discovery Policy determines where to route Users when they are attempting to sign in to your org. Users can be routed to a variety of Identity Providers (SAML2, IWA, AgentlessDSSO, X509, FACEBOOK, GOOGLE, LINKEDIN, MICROSOFT, OIDC) based on multiple conditions.",
		List: &plugin.ListConfig{
			Hydrate: listOktaIdpDiscoveryPolicies,
		},
		Columns: []*plugin.Column{
			// Top Columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the Policy."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Identifier of the Policy."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description of the Policy."},
			{Name: "created", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the Policy was created."},

			// Other Columns
			{Name: "last_updated", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the Policy was last modified."},
			{Name: "priority", Type: proto.ColumnType_INT, Description: "Priority of the Policy."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Status of the Policy: ACTIVE or INACTIVE."},
			{Name: "system", Type: proto.ColumnType_BOOL, Description: "This is set to true on system policies, which cannot be deleted."},

			// JSON Columns
			{Name: "conditions", Type: proto.ColumnType_JSON, Description: "Conditions for Policy."},
			{Name: "rules", Type: proto.ColumnType_JSON, Transform: transform.FromP(getpolicyRules, "rules"), Description: "Each Policy may contain one or more Rules. Rules, like Policies, contain conditions that must be satisfied for the Rule to be applied."},

			// Steampipe Columns
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: titleDescription},
		},
	}
}

func listOktaIdpDiscoveryPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client, err := Connect(ctx, d)

	input := &query.Params{
		Limit: 200,
	}

	if err != nil {
		logger.Error("listOktaIdpDiscoveryPolicies", "connect_error", err)
		return nil, err
	}
	if d.Table.Name == "okta_idp_discovery_policy" {
		input.Type = "IDP_DISCOVERY"
		input.Expand = "rules"
	}
	policies, resp, err := client.Policy.ListPolicies(ctx, input)
	if err != nil {
		logger.Error("listOktaIdpDiscoveryPolicies", "list_policies_error", err)
		return nil, err
	}
	for _, policy := range policies {
		d.StreamListItem(ctx, policy)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	// paging
	for resp.HasNextPage() {
		var nextPolicySet []*okta.AuthorizationServerPolicy
		resp, err = resp.Next(ctx, &nextPolicySet)
		if err != nil {
			logger.Error("listOktaIdpDiscoveryPolicies", "list_policies_paging_error", err)
			return nil, err
		}
		for _, policy := range nextPolicySet {
			d.StreamListItem(ctx, policy)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, err
}

//// TRANSFORM FUNCTION

func getpolicyRules(_ context.Context, d *transform.TransformData) (interface{}, error) {
	policy := d.HydrateItem.(*okta.AuthorizationServerPolicy)

	if policy.Embedded != nil {
		rules := policy.Embedded.(map[string]interface{})
		return rules["rules"], nil
	}
	return nil, nil
}
