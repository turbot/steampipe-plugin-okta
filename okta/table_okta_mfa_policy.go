package okta

import (
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableOktaMfaPolicy() *plugin.Table {
	return &plugin.Table{
		Name:        "okta_mfa_policy",
		Description: "The Multifactor (MFA) Enrollment Policy controls which MFA methods are available for a User, as well as when a User may enroll in a particular Factor.",
		List: &plugin.ListConfig{
			Hydrate: listPolicies,
		},
		Columns: listPoliciesWithSettingsColumns(),
	}
}
