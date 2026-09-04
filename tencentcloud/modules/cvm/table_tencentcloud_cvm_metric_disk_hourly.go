package cvm

import (
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// Table Definition

func TableTencentcloudCvmMetricDiskHourly() *plugin.Table {
	return cvmMetricTable(
		"tencentcloud_cvm_metric_disk_hourly",
		"Hourly disk utilization metric of Tencent Cloud CVM instances (last 24 hours by default; a timestamp predicate narrows the queried range).",
		"DiskUsage",
		"The disk utilization percentage of the instance at the timestamp.",
		3600,
		24*time.Hour,
	)
}
