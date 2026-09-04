package vpc

import (
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// TableTencentcloudVpcRegion returns the region metadata table for VPC,
func TableTencentcloudVpcRegion() *plugin.Table {
	return utils.TableRegion("vpc", "vpc")
}
