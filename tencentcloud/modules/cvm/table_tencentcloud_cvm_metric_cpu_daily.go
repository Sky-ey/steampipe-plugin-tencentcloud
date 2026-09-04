package cvm

import (
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// Table Definition

func TableTencentcloudCvmMetricCpuDaily() *plugin.Table {
	return cvmMetricTable(
		"tencentcloud_cvm_metric_cpu_daily",
		"Daily CPU utilization metric of Tencent Cloud CVM instances (last 30 days by default; a timestamp predicate narrows the queried range).",
		"CPUUsage",
		"The CPU utilization percentage of the instance at the timestamp.",
		86400,
		30*24*time.Hour,
	)
}
