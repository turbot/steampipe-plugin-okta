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

func tableOktaSignonPolicy() *plugin.Table {
	return &plugin.Table{
		Name:        "okta_signon_policy",
		Description: "Okta Sign On Policy controls the manner in which a user is allowed to sign on to Okta, including whether they are challenged for multifactor authentication (MFA) and how long they are allowed to remain signed in before re-authenticating.",
		List: &plugin.ListConfig{
			Hydrate: listOktaSignonPolicies,
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
			{Name: "rules", Type: proto.ColumnType_JSON, Transform: transform.FromField("Embedded.rules"), Description: "Each Policy may contain one or more Rules. Rules, like Policies, contain conditions that must be satisfied for the Rule to be applied."},

			// Steampipe Columns
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: titleDescription},
		}),
	}
}

//// LIST FUNCTION

func listOktaSignonPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("listOktaSignonPolicies", "connection_error", err)
		return nil, err
	}

	config := GetConfig(d.Connection)
	if config.EngineType == nil {
		return nil, fmt.Errorf("'engine_type' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}

	if config.EngineType != nil && *config.EngineType == "identity" {
		return listOktaSignonPoliciesIdentityEngine(ctx, client, d)
	}
	return listOktaSignonPoliciesClassicEngine(ctx, client, d)
}

//// CLASSIC ENGINE LIST FUNCTION

func listOktaSignonPoliciesClassicEngine(ctx context.Context, client *okta.Client, d *plugin.QueryData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Define query parameters for listing policies
	input := &query.Params{
		Type:   "OKTA_SIGN_ON",
		Expand: "rules",
	}

	// Fetch policies
	policies, resp, err := client.Policy.ListPolicies(ctx, input)
	if err != nil {
		logger.Error("listOktaSignonPoliciesClassicEngine", "list_policies_error", err)
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
			logger.Error("listOktaSignonPoliciesClassicEngine", "list_policies_paging_error", err)
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

func listOktaSignonPoliciesIdentityEngine(ctx context.Context, client *okta.Client, d *plugin.QueryData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Define query parameters for listing policies
	input := &query.Params{
		Type: "OKTA_SIGN_ON",
	}

	// Fetch policies
	policies, resp, err := client.Policy.ListPolicies(ctx, input)
	if err != nil {
		logger.Error("listOktaSignonPoliciesIdentityEngine", "list_policies_error", err)
		return nil, err
	}

	// Stream each policy item with additional rule fetching
	for _, policy := range policies {
		// Fetch rules for Identity Engine
		rules, _, err := client.Policy.ListPolicyRules(ctx, policy.Id)
		if err != nil {
			logger.Error("listOktaSignonPoliciesIdentityEngine", "list_policy_rules_error", err)
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
			logger.Error("listOktaSignonPoliciesIdentityEngine", "list_policies_paging_error", err)
			return nil, err
		}
		for _, policy := range nextPolicySet {
			// Fetch rules for Identity Engine
			rules, _, err := client.Policy.ListPolicyRules(ctx, policy.Id)
			if err != nil {
				logger.Error("listOktaSignonPoliciesIdentityEngine", "list_policy_rules_error", err)
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
