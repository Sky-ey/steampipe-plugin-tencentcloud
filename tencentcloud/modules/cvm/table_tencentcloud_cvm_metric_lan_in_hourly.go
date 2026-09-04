package cvm

import (
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// Table Definition

func TableTencentcloudCvmMetricLanInHourly() *plugin.Table {
	return cvmMetricTable(
		"tencentcloud_cvm_metric_lan_in_hourly",
		"Hourly private network inbound bandwidth metric of Tencent Cloud CVM instances (last 24 hours by default; a timestamp predicate narrows the queried range).",
		"LanIntraffic",
		"The average private network inbound bandwidth of the instance in Mbps at the timestamp.",
		3600,
		24*time.Hour,
	)
}
