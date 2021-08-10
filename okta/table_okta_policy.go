package okta

import (
	"context"
	"fmt"
	"strings"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableOktaPolicy() *plugin.Table {
	return &plugin.Table{
		Name:        "okta_policy",
		Description: "An Okta organization (org) is a root object and a container for all other Okta objects. It contains resources such as users, groups, and applications, as well as policy and configurations for your Okta environment.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("type"),
			Hydrate:    listOktaPolicies,
		},
		Columns: []*plugin.Column{
			// Top Columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the Policy."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Identifier of the Policy."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description of the Policy."},
			{Name: "created", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the Policy was created."},
			{Name: "type", Type: proto.ColumnType_STRING, Transform: transform.FromQual("type"), Description: "Specifies the type of Policy. Valid values: OKTA_SIGN_ON, PASSWORD, MFA_ENROLL or IDP_DISCOVERY."},

			// Other Columns
			{Name: "last_updated", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the Policy was last modified."},
			{Name: "priority", Type: proto.ColumnType_INT, Description: "Priority of the Policy."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Status of the Policy: ACTIVE or INACTIVE."},
			{Name: "system", Type: proto.ColumnType_BOOL, Description: "This is set to true on system policies, which cannot be deleted."},

			// JSON Columns
			{Name: "conditions", Type: proto.ColumnType_JSON, Description: "Conditions for Policy."},
			{Name: "settings", Type: proto.ColumnType_JSON, Description: "Priority of the Policy."},
			{Name: "data", Type: proto.ColumnType_JSON, Transform: transform.FromValue(), Description: "Priority of the Policy."},
		},
	}
}

//// LIST FUNCTION

func listOktaPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("listOktaPolicies", "connect", err)
		return nil, err
	}

	policyType := d.KeyColumnQuals["type"].GetStringValue()
	if policyType == "" {
		return nil, nil
	}

	if !helpers.StringSliceContains(policyTypes, policyType) {
		return nil, fmt.Errorf("%s is not a valid policy type. Valid policy types are: %s", policyType, strings.Join(policyTypes, ", "))
	}

	policies, resp, err := client.Policy.ListPolicies(ctx, &query.Params{Type: policyType, Expand: "rules"})
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
