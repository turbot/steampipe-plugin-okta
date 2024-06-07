package okta

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type oktaConfig struct {
	Domain         *string `hcl:"domain"`
	Token          *string `hcl:"token"`
	ClientID       *string `hcl:"client_id"`
	PrivateKey     *string `hcl:"private_key"`
	PrivateKeyID   *string `hcl:"private_key_id"`
	RequestTimeout *int64  `hcl:"request_timeout"`
	MaxRetries     *int32  `hcl:"max_retries"`
	MaxBackoff     *int64  `hcl:"max_backoff"`
}

func ConfigInstance() interface{} {
	return &oktaConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) oktaConfig {
	if connection == nil || connection.Config == nil {
		return oktaConfig{}
	}
	config, _ := connection.Config.(oktaConfig)
	return config
}
