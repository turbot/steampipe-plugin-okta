package main

import (
	"github.com/turbot/steampipe-plugin-okta/okta"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		PluginFunc: okta.Plugin})
}
