package okta

import (
	"context"

	"github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
		Columns: commonColumns([]*plugin.Column{
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
		}),
	}
}

//// LIST FUNCTION

func listOktaNetworkZones(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := ConnectV5(ctx, d)
	if err != nil {
		logger.Error("listOktaNetworkZones", "connect_error", err)
		return nil, err
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := int64(200)
	if d.QueryContext.Limit != nil {
		if *d.QueryContext.Limit < limit {
			limit = *d.QueryContext.Limit
		}
	}

	// Request
	zoneReq := client.NetworkZoneAPI.ListNetworkZones(ctx)

	zones, resp, err := zoneReq.Limit(int32(limit)).Execute()
	if err != nil {
		logger.Error("okta_factor.listOktaNetworkZones", "api_error", err)
		return nil, err
	}

	for _, zone := range zones {
		z := processNetworkZones(zone)
		if z != nil {
			d.StreamListItem(ctx, z)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	// paging
	for resp.HasNextPage() {
		var nextZoneSet []okta.ListNetworkZones200ResponseInner

		resp, err = resp.Next(&nextZoneSet)
		if err != nil {
			logger.Error("listOktaNetworkZones", "list_network_zones_paging_error", err)
			return nil, err
		}

		for _, zone := range nextZoneSet {
			z := processNetworkZones(zone)
			if z != nil {
				d.StreamListItem(ctx, z)

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}

	}

	return nil, err
}

//// HYDRATE FUNCTION

func getOktaNetworkZone(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client, err := ConnectV5(ctx, d)

	id := d.EqualsQuals["id"].GetStringValue()
	if id == "" {
		return nil, nil
	}
	if err != nil {
		logger.Error("getOktaNetworkZone", "connect_error", err)
		return nil, err
	}

	networkZone, _, err := client.NetworkZoneAPI.GetNetworkZone(ctx, id).Execute()
	if err != nil {
		logger.Error("getOktaNetworkZone", "get_network_zone_error", err)
		return nil, err
	}

	if networkZone != nil {
		return processNetworkZones(*networkZone), nil
	}

	return nil, nil
}

// Helper function to process and stream network zones
func processNetworkZones(networkZone okta.ListNetworkZones200ResponseInner) interface{} {
	if networkZone.DynamicNetworkZone != nil {
		return networkZone.DynamicNetworkZone
	}
	if networkZone.EnhancedDynamicNetworkZone != nil {
		return networkZone.EnhancedDynamicNetworkZone
	}
	if networkZone.IPNetworkZone != nil {
		return networkZone.IPNetworkZone
	}

	return nil
}
