package okta

import (
	"context"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableOktaGroupRule() *plugin.Table {
	return &plugin.Table{
		Name:        "okta_group_rule",
		Description: "Retrieve group rules for Okta. Group rules define conditions and actions for automating group membership.",
		Get: &plugin.GetConfig{
			Hydrate:           getOktaGroupRule,
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"Not found"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listOktaGroupRules,
		},
		Columns: []*plugin.Column{
			// Basic columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique identifier of the group rule."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the group rule."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Status of the group rule (e.g., ACTIVE, INACTIVE)."},
			{Name: "created", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the group rule was created."},
			{Name: "last_updated", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the group rule was last updated."},

			// JSON columns
			{Name: "conditions", Type: proto.ColumnType_JSON, Description: "Conditions that trigger this group rule."},
			{Name: "actions", Type: proto.ColumnType_JSON, Description: "Actions performed when the rule conditions are met."},

			// Steampipe-specific
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "Title of the group rule."},
		},
	}
}

//// LIST FUNCTION

func listOktaGroupRules(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listOktaGroupRules")

	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("listOktaGroupRules", "connection_error", err)
		return nil, err
	}

	// Fetch group rules
	groupRules, _, err := client.Group.ListGroupRules(ctx, nil)
	if err != nil {
		logger.Error("listOktaGroupRules", "list_error", err)
		return nil, err
	}

	for _, rule := range groupRules {
		d.StreamListItem(ctx, rule)

		// Exit if the context is canceled
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// GET FUNCTION

func getOktaGroupRule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getOktaGroupRule")

	// Retrieve the rule ID from the query or hydrate data
	var ruleId string
	if h.Item != nil {
		ruleId = h.Item.(*okta.GroupRule).Id
	} else {
		ruleId = d.EqualsQuals["id"].GetStringValue()
	}

	if ruleId == "" {
		return nil, nil
	}

	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("getOktaGroupRule", "connection_error", err)
		return nil, err
	}

	// Fetch the group rule by ID
	groupRule, _, err := client.Group.GetGroupRule(ctx, ruleId, nil)
	if err != nil {
		logger.Error("getOktaGroupRule", "get_error", err)
		return nil, err
	}

	return groupRule, nil
}