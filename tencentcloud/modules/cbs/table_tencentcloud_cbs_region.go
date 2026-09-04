package cbs

import (
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// TableTencentcloudCbsRegion returns the region metadata table for CBS,
func TableTencentcloudCbsRegion() *plugin.Table {
	return utils.TableRegion("cbs", "cbs")
}
