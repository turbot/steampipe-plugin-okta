package okta

import (
	"context"

	"github.com/okta/okta-sdk-golang/v4/okta"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableOktaDevice() *plugin.Table {
	return &plugin.Table{
		Name:        "okta_device",
		Description: "Okta’s device management is a crucial part of its broader suite of identity and access management solutions, helping organizations to secure their IT environments in an increasingly mobile and cloud-centric world.",
		Get: &plugin.GetConfig{
			Hydrate:           getOktaDevice,
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"Not found"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listOktaDevices,
			// TODO: Optional qualifier
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func:           listOktaDevices,
				MaxConcurrency: 10,
			},
		},
		Columns: []*plugin.Column{
			// Top Columns
			{Name: "display_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Profile.DisplayName"), Description: "Display name of the device."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique key for the device."},
			{Name: "created", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when device was created."},
			{Name: "last_updated", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the device record was last updated.", Transform: transform.FromField("LastUpdated").Transform(transform.NullIfZeroValue)},
			{Name: "resource_id", Type: proto.ColumnType_STRING, Description: "Alternate key for the Id."},
			{Name: "resource_type", Type: proto.ColumnType_STRING, Description: "The resource type."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The state object of the device."},

			// JSON Columns
			{Name: "profile", Type: proto.ColumnType_JSON, Description: "The Device's Profile properties."},
			{Name: "resource_display_name", Type: proto.ColumnType_JSON, Description: "Display name of the device."},
			{Name: "links", Type: proto.ColumnType_JSON, Description: "Specifies link relations (see Web Linking) available for the current status of an application using the JSON Hypertext Application Language specification."},
			{Name: "embedded", Type: proto.ColumnType_JSON, Description: "List of associated users for the device if the expand=user query parameter is specified in the request. Use expand=userSummary to get only a summary of each associated user for the device."},
			{Name: "additional_properties", Type: proto.ColumnType_JSON, Description: "additional properties of the device."},

			// Steampipe Columns
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Profile.DisplayName"), Description: titleDescription},
		},
	}
}

//// LIST FUNCTION

func listOktaDevices(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client, err := ConnectV4(ctx, d)
	if err != nil {
		logger.Error("okta_device.listOktaDevices", "connect_error", err)
		return nil, err
	}

	// Default maximum limit set as per documentation
	// https://developer.okta.com/docs/api/openapi/okta-management/management/tag/Device/#tag/Device/operation/listDevices!in=query&path=limit&t=request
	maxLimit := int64(20)
	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < maxLimit {
			maxLimit = *limit
		}
	}

	deviceReq := client.DeviceAPI.ListDevices(ctx).Expand("userSummary")
	deviceReq = deviceReq.Limit(int32(maxLimit))

	devices, resp, err := deviceReq.Execute()
	if err != nil {
		logger.Error("okta_device.listOktaDevices", "api_error", err)
		return nil, err
	}

	for _, device := range devices {
		d.StreamListItem(ctx, device)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	// paging
	for resp.HasNextPage() {
		var nextDeviceSet []okta.DeviceList
		resp, err = resp.Next(&nextDeviceSet)
		if err != nil {
			logger.Error("okta_device.listOktaDevices", "paging_error", err)
			return nil, err
		}
		for _, device := range nextDeviceSet {
			d.StreamListItem(ctx, device)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getOktaDevice(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	deviceId := d.EqualsQualString("id")

	if deviceId == "" {
		return nil, nil
	}

	client, err := ConnectV4(ctx, d)
	if err != nil {
		logger.Error("okta_device.getOktaDevice", "connect_error", err)
		return nil, err
	}

	deviceReq := client.DeviceAPI.GetDevice(ctx, deviceId)
	result, _, err := deviceReq.Execute()
	if err != nil {
		logger.Error("okta_device.getOktaDevice", "api_error", err)
		return nil, err
	}

	if result != nil {
		return *result, err
	}

	return nil, nil
}
