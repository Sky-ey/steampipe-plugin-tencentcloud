package cvm

import (
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// Table Definition

func TableTencentcloudCvmMetricWanOutDaily() *plugin.Table {
	return cvmMetricTable(
		"tencentcloud_cvm_metric_wan_out_daily",
		"Daily public network outbound bandwidth metric of Tencent Cloud CVM instances (last 30 days by default; a timestamp predicate narrows the queried range).",
		"WanOuttraffic",
		"The average public network outbound bandwidth of the instance in Mbps at the timestamp.",
		86400,
		30*24*time.Hour,
	)
}
