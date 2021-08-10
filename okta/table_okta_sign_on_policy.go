package okta

import (
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableOktaSignOnPolicy() *plugin.Table {
	return &plugin.Table{
		Name:        "okta_sign_on_policy",
		Description: "Okta Sign On Policy controls the manner in which a user is allowed to sign on to Okta.",
		List: &plugin.ListConfig{
			Hydrate: listOktaPolicies,
		},
		Columns: commonPolicyColumns(),
	}
}
