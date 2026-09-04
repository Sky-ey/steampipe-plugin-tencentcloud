package cvm

import (
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// Table Definition

func TableTencentcloudCvmMetricMemHourly() *plugin.Table {
	return cvmMetricTable(
		"tencentcloud_cvm_metric_mem_hourly",
		"Hourly memory utilization metric of Tencent Cloud CVM instances (last 24 hours by default; a timestamp predicate narrows the queried range).",
		"MemUsage",
		"The memory utilization percentage of the instance at the timestamp (excluding buffers and cache).",
		3600,
		24*time.Hour,
	)
}
