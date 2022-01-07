package okta

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

const pluginName = "steampipe-plugin-okta"

// Plugin creates this (okta) plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             pluginName,
		DefaultTransform: transform.FromCamel(),
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		TableMap: map[string]*plugin.Table{
			"okta_app_assigned_group":   tableOktaApplicationAssignedGroup(),
			"okta_application":          tableOktaApplication(),
			"okta_auth_server":          tableOktaAuthServer(),
			"okta_factor":               tableOktaFactor(),
			"okta_group":                tableOktaGroup(),
			"okta_idp_discovery_policy": tableOktaIdpDiscoveryPolicy(),
			"okta_mfa_policy":           tableOktaMfaPolicy(),
			"okta_network_zone":         tableOktaNetworkZone(),
			"okta_password_policy":      tableOktaPasswordPolicy(),
			"okta_signon_policy":        tableOktaSignonPolicy(),
			"okta_trusted_origin":       tableOktaTrustedOrigin(),
			"okta_user":                 tableOktaUser(),
			"okta_user_type":            tableOktaUserType(),
		},
	}

	return p
}
