package cvm

import (
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// Table Definition

func TableTencentcloudCvmMetricWanOutHourly() *plugin.Table {
	return cvmMetricTable(
		"tencentcloud_cvm_metric_wan_out_hourly",
		"Hourly public network outbound bandwidth metric of Tencent Cloud CVM instances (last 24 hours by default; a timestamp predicate narrows the queried range).",
		"WanOuttraffic",
		"The average public network outbound bandwidth of the instance in Mbps at the timestamp.",
		3600,
		24*time.Hour,
	)
}
