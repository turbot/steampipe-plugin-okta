package okta

import (
	"context"
	"fmt"
	"strings"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableOktaSystemLog() *plugin.Table {
	return &plugin.Table{
		Name:        "okta_system_log",
		Description: "Get all the system log events.",
		List: &plugin.ListConfig{
			Hydrate: listOktaSystemLogs,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "log_actor_id", Require: plugin.Optional},
				{Name: "log_ip_address", Require: plugin.Optional},
				{Name: "log_event_type", Require: plugin.Optional},
				{Name: "filter", Require: plugin.Optional},
				{Name: "log_event_time", Operators: []string{">", ">=", "=", "<", "<="}, Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "log_actor_id",
				Type:        proto.ColumnType_STRING,
				Description: "The Id of the log actor.",
				Transform:   transform.FromField("Actor.Id"),
			},
			{
				Name:        "log_actor_name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the log actor.",
				Transform:   transform.FromField("Actor.DisplayName"),
			},
			{
				Name:        "log_actor_email",
				Type:        proto.ColumnType_STRING,
				Description: "The email of the log actor.",
				Transform:   transform.From(actorTurbotData),
			},
			{
				Name:        "log_actor_username",
				Type:        proto.ColumnType_STRING,
				Description: "The username of the log actor.",
				Transform:   transform.From(actorTurbotData),
			},
			{
				Name:        "actor",
				Type:        proto.ColumnType_JSON,
				Description: "Represents who or what performed the action.",
			},
			{
				Name:        "authentication_context",
				Type:        proto.ColumnType_JSON,
				Description: "Provides context about the authentication.",
			},
			{
				Name:        "log_ip_address",
				Type:        proto.ColumnType_IPADDR,
				Description: "The log IP address.",
				Transform:   transform.FromField("Client.IpAddress"),
			},
			{
				Name:        "client",
				Type:        proto.ColumnType_JSON,
				Description: "Provides details about the client or user-agent making the request.",
			},
			{
				Name:        "debug_context",
				Type:        proto.ColumnType_JSON,
				Description: "Information useful for debugging.",
			},
			{
				Name:        "display_message",
				Type:        proto.ColumnType_STRING,
				Description: "A human-readable message for the event.",
			},
			{
				Name:        "log_event_type",
				Type:        proto.ColumnType_STRING,
				Description: "The type or nature of the event.",
				Transform:   transform.FromField("EventType"),
			},
			{
				Name:        "legacy_event_type",
				Type:        proto.ColumnType_STRING,
				Description: "The type of the event in older versions (if any).",
			},
			{
				Name:        "outcome",
				Type:        proto.ColumnType_JSON,
				Description: "Represents the result of the action (success, failure, etc.)",
			},
			{
				Name:        "log_event_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Timestamp indicating when the log event was generated or published.",
				Transform:   transform.FromField("Published"),
			},
			{
				Name:        "request",
				Type:        proto.ColumnType_JSON,
				Description: "Details about the incoming request.",
			},
			{
				Name:        "security_context",
				Type:        proto.ColumnType_JSON,
				Description: "Context related to the security aspects of the event.",
			},
			{
				Name:        "severity",
				Type:        proto.ColumnType_STRING,
				Description: "The seriousness or urgency of the event (like 'info', 'warning', 'error', etc.)",
			},
			{
				Name:        "target",
				Type:        proto.ColumnType_JSON,
				Description: "Represents the target or the object being acted upon.",
			},
			{
				Name:        "transaction",
				Type:        proto.ColumnType_JSON,
				Description: "Details about the transaction.",
			},
			{
				Name:        "log_event_id",
				Type:        proto.ColumnType_STRING,
				Description: "A unique identifier for the log event.",
				Transform:   transform.FromField("Uuid"),
			},
			{
				Name:        "version",
				Type:        proto.ColumnType_STRING,
				Description: "The version of the event.",
			},
			{
				Name:        "filter",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("filter"),
				Description: "Filter string to [filter](https://developer.okta.com/docs/reference/api/system-log/#bounded-requests) events. Input filter query should not be encoded.",
			},

			// Steampipe Columns
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Uuid"),
				Description: titleDescription,
			},
		},
	}
}

//// LIST FUNCTION

func listOktaSystemLogs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := Connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listOktaSystemLogs", "connect_error", err)
		return nil, err
	}

	// Default maximum limit set as per documentation
	// https://developer.okta.com/docs/reference/api/system-log/#request-parameters
	input := query.Params{
		Limit: 1000,
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < input.Limit {
			input.Limit = *limit
		}
	}

	equalQuals := d.EqualsQuals
	var queryFilter string
	filter := buildSystemLogQueryFilter(equalQuals, []string{"log_actor_id", "log_event_type", "log_ip_address"})

	// set the start and end time based on the provided log_event_time
	// https://developer.okta.com/docs/reference/api/system-log/#request-parameters
	if d.Quals["log_event_time"] != nil {
		for _, q := range d.Quals["log_event_time"].Quals {
			timestamp := q.Value.GetTimestampValue().AsTime().Format(filterTimeFormat)
			switch q.Operator {
			case "=":
				input.Since = timestamp
				input.Until = timestamp
			case ">=", ">":
				input.Since = timestamp
			case "<", "<=":
				input.Until = timestamp
			}
		}
	}

	if d.EqualsQualString("filter") != "" {
		queryFilter = d.EqualsQualString("filter")
	}

	if queryFilter != "" {
		input.Filter = queryFilter
	} else if len(filter) > 0 {
		input.Filter = strings.Join(filter, " and ")
	}
	events, resp, err := client.LogEvent.GetLogs(ctx, &input)
	if err != nil {
		plugin.Logger(ctx).Error("listOktaSystemLogs", "GetLogs", err)
		return nil, err
	}
	for _, event := range events {
		d.StreamListItem(ctx, event)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	// paging
	for {
		var nextLogSet []*okta.LogEvent
		resp, err = resp.Next(ctx, &nextLogSet)
		if err != nil {
			plugin.Logger(ctx).Error("listOktaSystemLogs", "list_logs_paging_error", err)
			return nil, err
		}
		if len(nextLogSet) == 0 {
			break
		}
		for _, event := range nextLogSet {
			d.StreamListItem(ctx, event)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func buildSystemLogQueryFilter(equalQuals plugin.KeyColumnEqualsQualMap, filterKeys []string) []string {
	filters := []string{}

	for k, v := range equalQuals {
		if v != nil && helpers.StringSliceContains(filterKeys, k) {
			if k == "log_actor_id" {
				filters = append(filters, fmt.Sprintf("actor.id eq \"%s\"", v.GetStringValue()))
			} else if k == "log_ip_address" {
				filters = append(filters, fmt.Sprintf("client.ipAddress eq \"%s\"", v.GetInetValue().GetAddr()))
			} else {
				filters = append(filters, fmt.Sprintf("eventType eq \"%s\"", v.GetStringValue()))
			}
		}
	}

	return filters
}

//// TRANSFORM FUNCTIONS

func actorTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	event := d.HydrateItem.(*okta.LogEvent)

	if strings.Contains(event.Actor.AlternateId, "@") {
		return event.Actor.AlternateId, nil
	}

	return nil, nil
}
