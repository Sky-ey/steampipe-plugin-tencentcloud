package clb

import (
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// TableTencentcloudClbRegion returns the region metadata table for CLB,
func TableTencentcloudClbRegion() *plugin.Table {
	return utils.TableRegion("clb", "clb")
}
