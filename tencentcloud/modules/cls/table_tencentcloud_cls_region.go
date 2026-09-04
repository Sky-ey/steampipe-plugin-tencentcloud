package cls

import (
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// TableTencentcloudClsRegion returns the region metadata table for CLS,
func TableTencentcloudClsRegion() *plugin.Table {
	return utils.TableRegion("cls", "cls")
}
