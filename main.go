package main

import (
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		PluginFunc: tencentcloud.Plugin,
	})
}
