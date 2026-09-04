package dc

import (
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// TableTencentcloudDcRegion returns the region metadata table for DC,
func TableTencentcloudDcRegion() *plugin.Table {
	return utils.TableRegion("dc", "dc")
}
