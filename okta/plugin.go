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
			"okta_application": tableOktaApplication(),
			"okta_group":       tableOktaGroup(),
			"okta_policy":      tableOktaPolicy(),
			"okta_user":        tableOktaUser(),
		},
	}

	return p
}
