package cvm

import (
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// TableTencentcloudCvmRegion returns the region metadata table for CVM,
func TableTencentcloudCvmRegion() *plugin.Table {
	return utils.TableRegion("cvm", "cvm")
}
