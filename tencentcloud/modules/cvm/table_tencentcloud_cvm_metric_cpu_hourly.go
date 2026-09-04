package cvm

import (
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// Table Definition

func TableTencentcloudCvmMetricCpuHourly() *plugin.Table {
	return cvmMetricTable(
		"tencentcloud_cvm_metric_cpu_hourly",
		"Hourly CPU utilization metric of Tencent Cloud CVM instances (last 24 hours by default; a timestamp predicate narrows the queried range).",
		"CPUUsage",
		"The CPU utilization percentage of the instance at the timestamp.",
		3600,
		24*time.Hour,
	)
}
