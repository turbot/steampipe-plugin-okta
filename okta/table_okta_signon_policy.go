package okta

import (
	"context"
	"strings"

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
			{Name: "rules", Type: proto.ColumnType_JSON, Hydrate: getOktaPolicyRules, Transform: transform.FromValue(), Description: "Each Policy may contain one or more Rules. Rules, like Policies, contain conditions that must be satisfied for the Rule to be applied."},

			// Steampipe Columns
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: titleDescription},
		}),
	}
}

func listOktaSignonPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client, err := Connect(ctx, d)

	input := &query.Params{}
	if err != nil {
		logger.Error("listOktaSignonPolicies", "connect_error", err)
		return nil, err
	}

	if d.Table.Name == "okta_signon_policy" {
		input.Type = "OKTA_SIGN_ON"
	}

	policies, resp, err := client.Policy.ListPolicies(ctx, input)
	if err != nil {
		logger.Error("listOktaSignonPolicies", "list_policies_error", err)
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
		// var nextPolicySet []*okta.Policy
		nextToken := strings.Split(strings.Split(resp.NextPage, "after=")[0], "&")[0]
		input.After = nextToken
		policies, resp, err = client.Policy.ListPolicies(ctx, input)
		if err != nil {
			logger.Error("listOktaSignonPolicies", "list_policies_paging_error", err)
			return nil, err
		}
		plugin.Logger(ctx).Error("Next page: ", resp.NextPage)
		for _, policy := range policies {
			d.StreamListItem(ctx, policy)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTION

func getOktaPolicyRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	if h.Item == nil {
		return nil, nil
	}
	policyId := ""

	switch item := h.Item.(type) {
	case *PolicyStructure:
		policyId = item.Id
	case *okta.Policy:
		policyId = item.Id
	case *okta.AuthorizationServerPolicy:
		policyId = item.Id
	}

	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("getOktaPolicyRules", "connect_error", err)
		return nil, err
	}

	var rules []*okta.PolicyRule

	policyRules, resp, err := client.Policy.ListPolicyRules(ctx, policyId)
	if err != nil {
		logger.Error("getOktaPolicyRules", "list_policies_error", err)
		return nil, err
	}

	rules = append(rules, policyRules...)

	// paging
	for resp.HasNextPage() {
		var nextPolicyRules []*okta.PolicyRule
		resp, err = resp.Next(ctx, &nextPolicyRules)
		if err != nil {
			logger.Error("getOktaPolicyRules", "list_policies_paging_error", err)
			return nil, err
		}
		rules = append(rules, nextPolicyRules...)
	}

	return rules, nil
}
