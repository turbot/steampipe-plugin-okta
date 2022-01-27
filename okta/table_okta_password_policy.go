package okta

import (
	"context"
	"time"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableOktaPasswordPolicy() *plugin.Table {
	return &plugin.Table{
		Name:        "okta_password_policy",
		Description: "The Password Policy determines the requirements for a user's password length and complexity, as well as the frequency with which a password must be changed. This Policy also governs the recovery operations that may be performed by the User, including change password, reset (forgot) password, and self-service password unlock.",
		List: &plugin.ListConfig{
			Hydrate: listPolicies,
		},
		Columns: listPoliciesWithSettingsColumns(),
	}
}

func listPoliciesWithSettingsColumns() []*plugin.Column {
	return []*plugin.Column{
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
		{Name: "settings", Type: proto.ColumnType_JSON, Description: "Settings of the password policy."},

		// Steampipe Columns
		{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: titleDescription},
	}
}

func listPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client, err := Connect(ctx, d)

	input := &query.Params{}
	if err != nil {
		logger.Error("listOktaPolicies", "connect_error", err)
		return nil, err
	}

	input.Expand = "rules"
	switch d.Table.Name {
	case "okta_password_policy":
		input.Type = "PASSWORD"
	case "okta_mfa_policy":
		input.Type = "MFA_ENROLL"
	}

	policies, resp, err := listPoliciesWithSettings(ctx, *client, input)
	if err != nil {
		logger.Error("listPolicies", "list_policies_with_settings_error", err)
		return nil, err
	}

	for _, policy := range policies {
		d.StreamListItem(ctx, policy)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	// paging
	for resp.HasNextPage() {
		var nextPolicySet []*okta.Policy
		resp, err = resp.Next(ctx, &nextPolicySet)
		if err != nil {
			logger.Error("listPolicies", "list_policies_with_settings_paging_error", err)
			return nil, err
		}
		for _, policy := range nextPolicySet {
			d.StreamListItem(ctx, policy)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

// Generic policy returned by
func listPoliciesWithSettings(ctx context.Context, client okta.Client, qp *query.Params) ([]*PolicyStructure, *okta.Response, error) {
	url := "/api/v1/policies"
	if qp != nil {
		url = url + qp.String()
	}

	requestExecutor := client.GetRequestExecutor()
	req, err := requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var policies []*PolicyStructure

	resp, err := requestExecutor.Do(ctx, req, &policies)
	if err != nil {
		return nil, resp, err
	}

	return policies, resp, nil
}

// generic policy missing Settings field
type PolicyStructure struct {
	Embedded    interface{}                `json:"_embedded,omitempty"`
	Links       interface{}                `json:"_links,omitempty"`
	Settings    interface{}                `json:"settings,omitempty"`
	Conditions  *okta.PolicyRuleConditions `json:"conditions,omitempty"`
	Created     *time.Time                 `json:"created,omitempty"`
	Description string                     `json:"description,omitempty"`
	Id          string                     `json:"id,omitempty"`
	LastUpdated *time.Time                 `json:"lastUpdated,omitempty"`
	Name        string                     `json:"name,omitempty"`
	Priority    int64                      `json:"priority,omitempty"`
	Status      string                     `json:"status,omitempty"`
	System      *bool                      `json:"system,omitempty"`
	Type        string                     `json:"type,omitempty"`
}
