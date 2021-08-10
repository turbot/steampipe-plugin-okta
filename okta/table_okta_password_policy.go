package okta

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableOktaPasswordPolicy() *plugin.Table {
	return &plugin.Table{
		Name:        "okta_password_policy",
		Description: "The Password Policy determines the requirements for a user's password length and complexity, as well as the frequency with which a password must be changed. This Policy also governs the recovery operations that may be performed by the User, including change password, reset (forgot) password, and self-service password unlock.",
		List: &plugin.ListConfig{
			Hydrate: listOktaPasswordPolicies,
		},
		Columns: policyColumns([]*plugin.Column{
			{Name: "settings", Type: proto.ColumnType_JSON, Description: "Settings of the Policy."},
		}),
	}
}

func listOktaPasswordPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client, err := Connect(ctx, d)

	input := &query.Params{}
	if err != nil {
		logger.Error("listOktaPolicies", "connect", err)
		return nil, err
	}

	if d.Table.Name == "okta_password_policy" {
		input.Type = "PASSWORD"
		input.Expand = "rules"
	}

	policyType := d.KeyColumnQuals["type"].GetStringValue()
	if input == nil && policyType == "" {
		return nil, nil
	} else {
		policyType = input.Type
	}

	if !helpers.StringSliceContains(policyTypes, policyType) {
		return nil, fmt.Errorf("%s is not a valid policy type. Valid policy types are: %s", policyType, strings.Join(policyTypes, ", "))
	}

	policies, resp, err := ListPasswordPolicies(ctx, *client, input)
	if err != nil {
		logger.Error("listOktaPolicies", "list policies", err)
		return nil, err
	}

	for _, policy := range policies {
		d.StreamListItem(ctx, policy)
	}

	// paging
	for resp.HasNextPage() {
		var nextPolicySet []*okta.Policy
		resp, err = resp.Next(ctx, &nextPolicySet)
		if err != nil {
			logger.Error("listOktaPolicies", "list policies paging", err)
			return nil, err
		}
		for _, policy := range nextPolicySet {
			d.StreamListItem(ctx, policy)
		}
	}

	return nil, err
}

// generic policy missing Settings field
type PasswordPolicy struct {
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

// Gets all password policies with the specified type.
func ListPasswordPolicies(ctx context.Context, client okta.Client, qp *query.Params) ([]*PasswordPolicy, *okta.Response, error) {
	url := "/api/v1/policies"
	if qp != nil {
		url = url + qp.String()
	}

	requestExecutor := client.GetRequestExecutor()
	req, err := requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var policies []*PasswordPolicy

	resp, err := requestExecutor.Do(ctx, req, &policies)
	if err != nil {
		return nil, resp, err
	}

	return policies, resp, nil
}
