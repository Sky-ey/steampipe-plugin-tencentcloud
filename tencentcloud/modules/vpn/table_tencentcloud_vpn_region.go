package vpn

import (
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// TableTencentcloudVpnRegion returns the region metadata table for VPN.
func TableTencentcloudVpnRegion() *plugin.Table {
	return utils.TableRegion("vpn", "vpc")
}
