package okta

import (
	"context"
	"fmt"

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
		Columns: commonColumns([]*plugin.Column{
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
		}),
	}
}

func listOktaIdpDiscoveryPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// logger := plugin.Logger(ctx)
	client, err := Connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listOktaIdpDiscoveryPolicies", "connect_error", err)
		return nil, err
	}

	config := GetConfig(d.Connection)
	if config.EngineType == nil {
		return nil, fmt.Errorf("'engine_type' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}

	if config.EngineType != nil && *config.EngineType == "identity" {
		return listOktaIdpDiscoveryPoliciesIdentityEngine(ctx, client, d)
	}
	return listOktaIdpDiscoveryPoliciesClassicEngine(ctx, client, d)
}

//// CLASSIC ENGINE LIST FUNCTION

func listOktaIdpDiscoveryPoliciesClassicEngine(ctx context.Context, client *okta.Client, d *plugin.QueryData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Define query parameters for listing policies
	input := &query.Params{
		Limit:  200,
		Type:   "IDP_DISCOVERY",
		Expand: "rules",
	}

	// Fetch policies
	policies, resp, err := client.Policy.ListPolicies(ctx, input)
	if err != nil {
		logger.Error("listOktaIdpDiscoveryPoliciesClassicEngine", "list_policies_error", err)
		return nil, err
	}

	// Stream each policy item
	for _, policy := range policies {
		d.StreamListItem(ctx, policy)

		// Check if the context is canceled or the row limit is reached
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	// Handle paging
	for resp.HasNextPage() {
		var nextPolicySet []*okta.Policy
		resp, err = resp.Next(ctx, &nextPolicySet)
		if err != nil {
			logger.Error("listOktaIdpDiscoveryPoliciesClassicEngine", "list_policies_paging_error", err)
			return nil, err
		}
		for _, policy := range nextPolicySet {
			d.StreamListItem(ctx, policy)

			// Check if the context is canceled or the row limit is reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// IDENTITY ENGINE LIST FUNCTION

func listOktaIdpDiscoveryPoliciesIdentityEngine(ctx context.Context, client *okta.Client, d *plugin.QueryData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Define query parameters for listing policies
	input := &query.Params{
		Limit: 200,
		Type:  "IDP_DISCOVERY",
	}

	// Fetch policies
	policies, resp, err := client.Policy.ListPolicies(ctx, input)
	if err != nil {
		logger.Error("listOktaIdpDiscoveryPoliciesIdentityEngine", "list_policies_error", err)
		return nil, err
	}

	// Stream each policy item with additional rule fetching
	for _, policy := range policies {
		// Fetch rules for Identity Engine
		rules, _, err := client.Policy.ListPolicyRules(ctx, policy.Id)
		if err != nil {
			logger.Error("listOktaIdpDiscoveryPoliciesIdentityEngine", "list_policy_rules_error", err)
			return nil, err
		}
		policy.Embedded = map[string]interface{}{
			"rules": rules,
		}
		d.StreamListItem(ctx, policy)

		// Check if the context is canceled or the row limit is reached
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	// Handle paging
	for resp.HasNextPage() {
		var nextPolicySet []*okta.Policy
		resp, err = resp.Next(ctx, &nextPolicySet)
		if err != nil {
			logger.Error("listOktaIdpDiscoveryPoliciesIdentityEngine", "list_policies_paging_error", err)
			return nil, err
		}
		for _, policy := range nextPolicySet {
			// Fetch rules for Identity Engine
			rules, _, err := client.Policy.ListPolicyRules(ctx, policy.Id)
			if err != nil {
				logger.Error("listOktaIdpDiscoveryPoliciesIdentityEngine", "list_policy_rules_error", err)
				return nil, err
			}
			policy.Embedded = map[string]interface{}{
				"rules": rules,
			}
			d.StreamListItem(ctx, policy)

			// Check if the context is canceled or the row limit is reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// TRANSFORM FUNCTION

func getpolicyRules(_ context.Context, d *transform.TransformData) (interface{}, error) {

	switch item := d.HydrateItem.(type) {
	case *okta.AuthorizationServerPolicy:
		if item.Embedded != nil {
			rules := item.Embedded.(map[string]interface{})
			return rules["rules"], nil
		}
	case *okta.Policy:
		if item.Embedded != nil {
			rules := item.Embedded.(map[string]interface{})
			return rules["rules"], nil
		}
	}

	return nil, nil
}
