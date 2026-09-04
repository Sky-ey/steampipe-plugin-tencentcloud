package cvm

import (
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// Table Definition

func TableTencentcloudCvmMetricMemDaily() *plugin.Table {
	return cvmMetricTable(
		"tencentcloud_cvm_metric_mem_daily",
		"Daily memory utilization metric of Tencent Cloud CVM instances (last 30 days by default; a timestamp predicate narrows the queried range).",
		"MemUsage",
		"The memory utilization percentage of the instance at the timestamp (excluding buffers and cache).",
		86400,
		30*24*time.Hour,
	)
}
