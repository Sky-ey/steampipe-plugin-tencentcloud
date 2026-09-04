package cvm

import (
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// Table Definition

func TableTencentcloudCvmMetricWanInDaily() *plugin.Table {
	return cvmMetricTable(
		"tencentcloud_cvm_metric_wan_in_daily",
		"Daily public network inbound bandwidth metric of Tencent Cloud CVM instances (last 30 days by default; a timestamp predicate narrows the queried range).",
		"WanIntraffic",
		"The average public network inbound bandwidth of the instance in Mbps at the timestamp.",
		86400,
		30*24*time.Hour,
	)
}
