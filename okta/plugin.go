package okta

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

const pluginName = "steampipe-plugin-okta"

// Plugin creates this (okta) plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             pluginName,
		DefaultTransform: transform.FromCamel(),
		ConnectionKeyColumns: []plugin.ConnectionKeyColumn{
			{
				Name:    "domain",
				Hydrate: getOktaDomainName,
			},
		},
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
		},
		TableMap: map[string]*plugin.Table{
			"okta_app_assigned_group":    tableOktaApplicationAssignedGroup(),
			"okta_app_assigned_user":     tableOktaApplicationAssignedUser(),
			"okta_application":           tableOktaApplication(),
			"okta_auth_server":           tableOktaAuthServer(),
			"okta_authentication_policy": tableOktaAuthenticationPolicy(),
			"okta_device":                tableOktaDevice(),
			"okta_factor":                tableOktaFactor(),
			"okta_group":                 tableOktaGroup(),
			"okta_group_owner":           tableOktaGroupOwner(),
			"okta_group_rule":            tableOktaGroupRule(),
			"okta_idp_discovery_policy":  tableOktaIdpDiscoveryPolicy(),
			"okta_mfa_policy":            tableOktaMfaPolicy(),
			"okta_network_zone":          tableOktaNetworkZone(),
			"okta_password_policy":       tableOktaPasswordPolicy(),
			"okta_signon_policy":         tableOktaSignonPolicy(),
			"okta_trusted_origin":        tableOktaTrustedOrigin(),
			"okta_user":                  tableOktaUser(),
			"okta_user_type":             tableOktaUserType(),
		},
	}

	return p
}
