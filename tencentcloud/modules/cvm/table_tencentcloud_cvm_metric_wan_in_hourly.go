package cvm

import (
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// Table Definition

func TableTencentcloudCvmMetricWanInHourly() *plugin.Table {
	return cvmMetricTable(
		"tencentcloud_cvm_metric_wan_in_hourly",
		"Hourly public network inbound bandwidth metric of Tencent Cloud CVM instances (last 24 hours by default; a timestamp predicate narrows the queried range).",
		"WanIntraffic",
		"The average public network inbound bandwidth of the instance in Mbps at the timestamp.",
		3600,
		24*time.Hour,
	)
}
