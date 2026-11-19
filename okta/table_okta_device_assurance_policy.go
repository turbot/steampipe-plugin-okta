package okta

import (
	"context"
	"encoding/json"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableOktaDeviceAssurancePolicy() *plugin.Table {
	return &plugin.Table{
		Name:        "okta_device_assurance_policy",
		Description: "Device assurance policies define the device posture requirements (OS version, disk encryption, screen lock, and so on) that a device must meet to access resources protected by an authentication policy.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getOktaDeviceAssurancePolicy,
		},
		List: &plugin.ListConfig{
			Hydrate: listOktaDeviceAssurancePolicies,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top Columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Display name of the device assurance policy."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Identifier of the device assurance policy."},
			{Name: "platform", Type: proto.ColumnType_STRING, Description: "Platform the policy applies to: ANDROID, CHROMEOS, IOS, MACOS, or WINDOWS."},

			// Other Columns
			{Name: "created_by", Type: proto.ColumnType_STRING, Description: "Identifier of the user who created the policy."},
			{Name: "created_date", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the policy was created."},
			{Name: "jailbreak", Type: proto.ColumnType_BOOL, Description: "Indicates whether jailbroken or rooted devices are allowed. False means such devices do not meet the policy. Applies to ANDROID and IOS."},
			{Name: "last_update", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the policy was last modified."},
			{Name: "last_updated_by", Type: proto.ColumnType_STRING, Description: "Identifier of the user who last modified the policy."},
			{Name: "secure_hardware_present", Type: proto.ColumnType_BOOL, Description: "Indicates whether a secure hardware module (for example a TPM or Secure Enclave) is required. Applies to ANDROID, MACOS, and WINDOWS."},

			// JSON Columns
			{Name: "disk_encryption_type", Type: proto.ColumnType_JSON, Description: "Disk encryption types the device must use, as an object with an `include` list. Applies to ANDROID, MACOS, and WINDOWS."},
			{Name: "os_version", Type: proto.ColumnType_JSON, Description: "Minimum operating system version required, as an object with a `minimum` value."},
			{Name: "os_version_constraints", Type: proto.ColumnType_JSON, Description: "Per-major-version operating system constraints. Applies to WINDOWS."},
			{Name: "screen_lock_type", Type: proto.ColumnType_JSON, Description: "Screen lock types the device must use, as an object with an `include` list (for example BIOMETRIC, PASSCODE)."},
			{Name: "third_party_signal_providers", Type: proto.ColumnType_JSON, Description: "Third-party signal provider requirements (for example endpoint security integrations)."},

			// Steampipe Columns
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: titleDescription},
		}),
	}
}

// deviceAssurancePolicy is a flat view over the platform-specific policy
// models returned by the API. Fields that do not apply to a platform are nil.
type deviceAssurancePolicy struct {
	Id                        *string       `json:"id"`
	Name                      *string       `json:"name"`
	Platform                  *string       `json:"platform"`
	CreatedBy                 *string       `json:"createdBy"`
	CreatedDate               *string       `json:"createdDate"`
	LastUpdate                *string       `json:"lastUpdate"`
	LastUpdatedBy             *string       `json:"lastUpdatedBy"`
	Jailbreak                 *bool         `json:"jailbreak"`
	SecureHardwarePresent     *bool         `json:"secureHardwarePresent"`
	DiskEncryptionType        interface{}   `json:"diskEncryptionType"`
	OsVersion                 interface{}   `json:"osVersion"`
	OsVersionConstraints      []interface{} `json:"osVersionConstraints"`
	ScreenLockType            interface{}   `json:"screenLockType"`
	ThirdPartySignalProviders interface{}   `json:"thirdPartySignalProviders"`
}

//// LIST FUNCTION

func listOktaDeviceAssurancePolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := ConnectV5(ctx, d)
	if err != nil {
		logger.Error("okta_device_assurance_policy.listOktaDeviceAssurancePolicies", "connect_error", err)
		return nil, err
	}

	// The list endpoint returns all policies in a single response (no paging).
	policies, _, err := client.DeviceAssuranceAPI.ListDeviceAssurancePolicies(ctx).Execute()
	if err != nil {
		logger.Error("okta_device_assurance_policy.listOktaDeviceAssurancePolicies", "api_error", err)
		return nil, err
	}

	for _, policy := range policies {
		row, err := flattenDeviceAssurancePolicy(policy.GetActualInstance())
		if err != nil {
			logger.Error("okta_device_assurance_policy.listOktaDeviceAssurancePolicies", "parse_error", err)
			return nil, err
		}
		if row == nil {
			continue
		}
		d.StreamListItem(ctx, row)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTION

func getOktaDeviceAssurancePolicy(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	id := d.EqualsQuals["id"].GetStringValue()
	if id == "" {
		return nil, nil
	}

	client, err := ConnectV5(ctx, d)
	if err != nil {
		logger.Error("okta_device_assurance_policy.getOktaDeviceAssurancePolicy", "connect_error", err)
		return nil, err
	}

	policy, _, err := client.DeviceAssuranceAPI.GetDeviceAssurancePolicy(ctx, id).Execute()
	if err != nil {
		logger.Error("okta_device_assurance_policy.getOktaDeviceAssurancePolicy", "api_error", err)
		return nil, err
	}
	if policy == nil {
		return nil, nil
	}

	return flattenDeviceAssurancePolicy(policy.GetActualInstance())
}

//// TRANSFORM FUNCTION

func flattenDeviceAssurancePolicy(instance interface{}) (*deviceAssurancePolicy, error) {
	if instance == nil {
		return nil, nil
	}

	raw, err := json.Marshal(instance)
	if err != nil {
		return nil, err
	}

	var row deviceAssurancePolicy
	if err := json.Unmarshal(raw, &row); err != nil {
		return nil, err
	}

	return &row, nil
}
