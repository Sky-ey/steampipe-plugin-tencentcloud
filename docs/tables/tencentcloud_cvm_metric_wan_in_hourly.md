---
title: "Steampipe Table: tencentcloud_cvm_metric_wan_in_hourly - Query CVM instance hourly public network inbound bandwidth using SQL"
description: "Query the hourly public network inbound bandwidth metric of Tencent Cloud CVM instances over the last 24 hours."
folder: "CVM"
---

# Table: tencentcloud_cvm_metric_wan_in_hourly - Query CVM instance hourly public network inbound bandwidth using SQL

Tencent Cloud Monitor (Cloud Monitor) collects per-instance public network inbound bandwidth for every CVM instance. The `tencentcloud_cvm_metric_wan_in_hourly` table uses the `GetMonitorData` API (namespace `QCE/CVM`, metric `WanIntraffic`, period 3600 seconds) and returns one row per instance per hour over the last 24 hours, across all regions or a single region when the `region` column is filtered.

Each row represents a single hourly data point of one instance: the instance ID, the timestamp of the data point, and the public network inbound bandwidth value (in Mbps). The values are the inbound traffic of the instance together with its associated EIP and CLB.

## Table Usage Guide

Use this table for recent, finer-grained troubleshooting — spotting public network inbound bandwidth spikes within the last day, correlating load with deployments, or verifying autoscaling behaviour. Hourly granularity is retained for 93 days, but this table only reads the last 24 hours to keep the row count manageable. For month-scale trend analysis use `tencentcloud_cvm_metric_wan_in_daily`.

Filtering by `instance_id` queries only that instance; omitting it fans out across every instance in the region. Only instances with a public IP address are included in the fan-out; instances without one are skipped, because the `GetMonitorData` API rejects an entire batch with `InvalidParameterValue` when any instance in it has no public IP.

When `instance_id` is omitted, the fan-out makes one paginated `DescribeInstances` pass per scanned region to enumerate the instance IDs, then calls `GetMonitorData` once per batch of up to 10 instances — so N instances in a region cost roughly N/10 `GetMonitorData` calls. This repeats in every scanned region: a single region when the `region` column is filtered (or when the connection has no `regions` list configured), or every region matched by the configured `regions` patterns (e.g. `["*"]`) otherwise.

Keep call volume in mind on large accounts: the `GetMonitorData` API includes a free quota of 1 million calls per month, and calls beyond the quota are billed. Filtering by `instance_id`, or narrowing the scanned regions and the `timestamp` window, keeps usage within the free quota.

By default the table reads the last 24 hours. A `timestamp` predicate in the `WHERE` clause narrows the time range and is pushed down to the `GetMonitorData` API (`StartTime`/`EndTime`), so you can query any window within the 93-day retention period instead of the default lookback.

## Examples

### Basic info

```sql+postgres
select
  instance_id,
  timestamp,
  value,
  region
from
  tencentcloud_cvm_metric_wan_in_hourly
order by
  instance_id,
  timestamp;
```

```sql+sqlite
select
  instance_id,
  timestamp,
  value,
  region
from
  tencentcloud_cvm_metric_wan_in_hourly
order by
  instance_id,
  timestamp;
```

### Hourly public network inbound bandwidth over a custom time range

The `timestamp` range is pushed down to the `GetMonitorData` API, so only the requested window is read.

```sql+postgres
select
  instance_id,
  timestamp,
  value
from
  tencentcloud_cvm_metric_wan_in_hourly
where
  timestamp >= '2026-07-01T00:00:00Z'
  and timestamp < '2026-08-01T00:00:00Z'
order by
  instance_id,
  timestamp;
```

```sql+sqlite
select
  instance_id,
  timestamp,
  value
from
  tencentcloud_cvm_metric_wan_in_hourly
where
  timestamp >= '2026-07-01T00:00:00Z'
  and timestamp < '2026-08-01T00:00:00Z'
order by
  instance_id,
  timestamp;
```

### Hourly public network inbound bandwidth of a single instance

```sql+postgres
select
  instance_id,
  timestamp,
  value
from
  tencentcloud_cvm_metric_wan_in_hourly
where
  instance_id = 'ins-xxxxxxxx'
order by
  timestamp;
```

```sql+sqlite
select
  instance_id,
  timestamp,
  value
from
  tencentcloud_cvm_metric_wan_in_hourly
where
  instance_id = 'ins-xxxxxxxx'
order by
  timestamp;
```

### Average hourly public network inbound bandwidth per instance

```sql+postgres
select
  instance_id,
  round(avg(value)::numeric, 2) as avg_wan_in_mbps,
  round(max(value)::numeric, 2) as max_wan_in_mbps
from
  tencentcloud_cvm_metric_wan_in_hourly
group by
  instance_id
order by
  avg_wan_in_mbps desc;
```

```sql+sqlite
select
  instance_id,
  round(avg(value), 2) as avg_wan_in_mbps,
  round(max(value), 2) as max_wan_in_mbps
from
  tencentcloud_cvm_metric_wan_in_hourly
group by
  instance_id
order by
  avg_wan_in_mbps desc;
```

### Find instances with inbound bandwidth above 100 Mbps

```sql+postgres
select
  instance_id,
  timestamp,
  value,
  region
from
  tencentcloud_cvm_metric_wan_in_hourly
where
  value > 100
order by
  value desc;
```

```sql+sqlite
select
  instance_id,
  timestamp,
  value,
  region
from
  tencentcloud_cvm_metric_wan_in_hourly
where
  value > 100
order by
  value desc;
```
