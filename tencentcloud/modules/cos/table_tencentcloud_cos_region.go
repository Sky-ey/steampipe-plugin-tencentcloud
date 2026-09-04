package cos

import (
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// TableTencentcloudCosRegion returns the region metadata table for COS,
// backed by the generic region.DescribeRegions API (Product="cos").
//
// COS uses an independent RESTful/S3V4 architecture for bucket and object
// operations, but the region metadata table relies on the generic region
// service (same as all other *_region tables), not the COS SDK.
func TableTencentcloudCosRegion() *plugin.Table {
	return utils.TableRegion("cos", "cos")
}
