package cvm

import (
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// Table Definition

func TableTencentcloudCvmMetricDiskDaily() *plugin.Table {
	return cvmMetricTable(
		"tencentcloud_cvm_metric_disk_daily",
		"Daily disk utilization metric of Tencent Cloud CVM instances (last 30 days by default; a timestamp predicate narrows the queried range).",
		"CvmDiskUsage",
		"The disk utilization percentage of the instance at the timestamp.",
		86400,
		30*24*time.Hour,
	)
}
