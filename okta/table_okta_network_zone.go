package okta

import (
	"context"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableOktaNetworkZone() *plugin.Table {
	return &plugin.Table{
		Name:        "okta_network_zone",
		Description: "The Okta Zones provides operations to manage Zones in your organization. There are two usage Zone types: Policy Network Zones and Block List Network Zones. Policy Network Zones are used to guide policy decisions. Block List Network Zones are used to deny access from certain IP addresses, locations, proxy types, or Autonomous System Numbers (ASNs) before policy evaluation.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getOktaNetworkZone,
		},
		List: &plugin.ListConfig{
			Hydrate: listOktaNetworkZones,
		},
		Columns: []*plugin.Column{
			// Top Columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Unique name for the zone."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Identifier of the network zone."},
			{Name: "created", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the network zone was created."},

			// Other Columns
			{Name: "last_updated", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the network zone was last modified."},
			{Name: "proxy_type", Type: proto.ColumnType_STRING, Description: "One of: '' or null (when not specified), Any (meaning any proxy), Tor, NotTorAnonymizer."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Status of the network zone: ACTIVE or INACTIVE."},
			{Name: "system", Type: proto.ColumnType_BOOL, Description: "Indicates if this is a system network zone."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The type of the network zone."},
			{Name: "usage", Type: proto.ColumnType_STRING, Description: "Usage of Zone: POLICY, BLOCKLIST."},

			// JSON Columns
			{Name: "asns", Type: proto.ColumnType_JSON, Description: "Format of each array value: a string representation of an ASN numeric value."},
			{Name: "gateways", Type: proto.ColumnType_JSON, Description: "IP addresses (range or CIDR form) of the zone."},
			{Name: "locations", Type: proto.ColumnType_JSON, Description: "The geolocations of the zone."},
			{Name: "proxies", Type: proto.ColumnType_JSON, Description: "IP addresses (range or CIDR form) that are allowed to forward a request from gateway addresses. These proxies are automatically trusted by Threat Insights. These proxies are used to identify the client IP of a request."},

			// Steampipe Columns
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: titleDescription},
		},
	}
}

//// LIST FUNCTION

func listOktaNetworkZones(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("listOktaNetworkZones", "connect_error", err)
		return nil, err
	}

	// Maximum limit isn't mentioned in the documentation
	// Default maximum limit is set as 1000
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

	networkZones, resp, err := client.NetworkZone.ListNetworkZones(ctx, &input)
	if err != nil {
		logger.Error("listOktaNetworkZones", "list_network_zones_error", err)
		return nil, err
	}

	for _, networkZone := range networkZones {
		d.StreamListItem(ctx, networkZone)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	// paging
	for resp.HasNextPage() {
		var nextZoneSet []*okta.NetworkZone
		resp, err = resp.Next(ctx, &nextZoneSet)
		if err != nil {
			logger.Error("listOktaNetworkZones", "list_network_zones_paging_error", err)
			return nil, err
		}
		for _, networkZone := range nextZoneSet {
			d.StreamListItem(ctx, networkZone)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTION

func getOktaNetworkZone(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client, err := Connect(ctx, d)

	id := d.KeyColumnQuals["id"].GetStringValue()
	if id == "" {
		return nil, nil
	}
	if err != nil {
		logger.Error("getOktaNetworkZone", "connect_error", err)
		return nil, err
	}

	networkZone, _, err := client.NetworkZone.GetNetworkZone(ctx, id)
	if err != nil {
		logger.Error("getOktaNetworkZone", "get_network_zone_error", err)
		return nil, err
	}

	return networkZone, err
}
