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

func tableOktaAuthenticationPolicy() *plugin.Table {
	return &plugin.Table{
		Name:        "okta_authentication_policy",
		Description: "Okta Authentication Policy controls the manner in which a user is authenticated, including MFA requirements.",
		List: &plugin.ListConfig{
			Hydrate: listOktaAuthenticationPolicies,
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
			{Name: "rules", Type: proto.ColumnType_JSON, Hydrate: getOktaPolicyRules, Transform: transform.FromValue(), Description: "Each Policy may contain one or more Rules. Rules, like Policies, contain conditions that must be satisfied for the Rule to be applied."},
			{Name: "resource_mapping", Type: proto.ColumnType_JSON, Hydrate: getOktaPolicyAssociatedResources, Transform: transform.FromValue(), Description: "The resources that are mapped to the Policy."},

			// Steampipe Columns
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: titleDescription},
		},
	}
}

//// LIST FUNCTION

func listOktaAuthenticationPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	input := &query.Params{
		Type: "ACCESS_POLICY",
	}

	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("listOktaAuthenticationPolicies", "connect_error", err)
		return nil, err
	}

	policies, resp, err := client.Policy.ListPolicies(ctx, input)
	if err != nil {
		logger.Error("listOktaAuthenticationPolicies", "list_policies_error", err)
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
		var nextPolicySet []*okta.Policy
		resp, err = resp.Next(ctx, &nextPolicySet)
		if err != nil {
			logger.Error("listOktaAuthenticationPolicies", "list_policies_paging_error", err)
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
