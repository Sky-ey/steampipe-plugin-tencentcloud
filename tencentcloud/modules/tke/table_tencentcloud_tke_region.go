package tke

import (
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// TableTencentcloudTkeRegion returns the region metadata table for TKE,
func TableTencentcloudTkeRegion() *plugin.Table {
	return utils.TableRegion("tke", "tke")
}
