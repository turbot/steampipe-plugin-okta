package okta

import (
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/schema"
)

type azureADConfig struct {
	Domain *string `cty:"domain"`
	Token  *string `cty:"token"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"domain": {
		Type: schema.TypeString,
	},
	"token": {
		Type: schema.TypeString,
	},
}

func ConfigInstance() interface{} {
	return &azureADConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) azureADConfig {
	if connection == nil || connection.Config == nil {
		return azureADConfig{}
	}
	config, _ := connection.Config.(azureADConfig)
	return config
}
