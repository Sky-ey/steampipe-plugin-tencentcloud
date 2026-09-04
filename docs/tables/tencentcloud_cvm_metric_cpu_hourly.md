---
title: "Steampipe Table: tencentcloud_cvm_metric_cpu_hourly - Query CVM instance hourly CPU utilization using SQL"
description: "Query the hourly CPU utilization metric of Tencent Cloud CVM instances over the last 24 hours."
folder: "CVM"
---

# Table: tencentcloud_cvm_metric_cpu_hourly - Query CVM instance hourly CPU utilization using SQL

Tencent Cloud Monitor (Cloud Monitor) collects per-instance CPU utilization for every CVM instance. The `tencentcloud_cvm_metric_cpu_hourly` table uses the `GetMonitorData` API (namespace `QCE/CVM`, metric `CPUUsage`, period 3600 seconds) and returns one row per instance per hour over the last 24 hours, across all regions or a single region when the `region` column is filtered.

Each row represents a single hourly data point of one instance: the instance ID, the timestamp of the data point, and the CPU utilization value (a percentage).

## Table Usage Guide

Use this table for recent, finer-grained troubleshooting — spotting CPU spikes within the last day, correlating load with deployments, or verifying autoscaling behaviour. Hourly granularity is retained for 93 days, but this table only reads the last 24 hours to keep the row count manageable. For month-scale trend analysis use `tencentcloud_cvm_metric_cpu_daily`.

Filtering by `instance_id` queries only that instance; omitting it fans out across every instance in the region.

By default the table reads the last 24 hours. A `timestamp` predicate in the `WHERE` clause narrows the time range and is pushed down to the `GetMonitorData` API (`StartTime`/`EndTime`), so you can query any window within the 93-day retention period instead of the default lookback.

## Examples

### Basic hourly CPU utilization

```sql+postgres
select
  instance_id,
  timestamp,
  value,
  region
from
  tencentcloud_cvm_metric_cpu_hourly
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
  tencentcloud_cvm_metric_cpu_hourly
order by
  instance_id,
  timestamp;
```

### Hourly CPU utilization over a custom time range

The `timestamp` range is pushed down to the `GetMonitorData` API, so only the requested window is read.

```sql+postgres
select
  instance_id,
  timestamp,
  value
from
  tencentcloud_cvm_metric_cpu_hourly
where
  timestamp >= '2026-08-12T00:00:00Z'
  and timestamp < '2026-08-13T00:00:00Z'
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
  tencentcloud_cvm_metric_cpu_hourly
where
  timestamp >= '2026-08-12T00:00:00Z'
  and timestamp < '2026-08-13T00:00:00Z'
order by
  instance_id,
  timestamp;
```

### Hourly CPU utilization of a single instance

```sql+postgres
select
  instance_id,
  timestamp,
  value
from
  tencentcloud_cvm_metric_cpu_hourly
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
  tencentcloud_cvm_metric_cpu_hourly
where
  instance_id = 'ins-xxxxxxxx'
order by
  timestamp;
```

### Peak hourly CPU utilization per instance

```sql+postgres
select
  instance_id,
  round(max(value)::numeric, 2) as peak_cpu_utilization,
  round(avg(value)::numeric, 2) as avg_cpu_utilization
from
  tencentcloud_cvm_metric_cpu_hourly
group by
  instance_id
order by
  peak_cpu_utilization desc;
```

```sql+sqlite
select
  instance_id,
  round(max(value), 2) as peak_cpu_utilization,
  round(avg(value), 2) as avg_cpu_utilization
from
  tencentcloud_cvm_metric_cpu_hourly
group by
  instance_id
order by
  peak_cpu_utilization desc;
```

### Find hourly CPU spikes above a threshold

```sql+postgres
select
  instance_id,
  timestamp,
  value,
  region
from
  tencentcloud_cvm_metric_cpu_hourly
where
  value > 90
order by
  timestamp desc;
```

```sql+sqlite
select
  instance_id,
  timestamp,
  value,
  region
from
  tencentcloud_cvm_metric_cpu_hourly
where
  value > 90
order by
  timestamp desc;
```

### Count hourly data points per instance

```sql+postgres
select
  instance_id,
  count(*) as data_points,
  min(timestamp) as first_point,
  max(timestamp) as last_point
from
  tencentcloud_cvm_metric_cpu_hourly
group by
  instance_id
order by
  data_points desc;
```

```sql+sqlite
select
  instance_id,
  count(*) as data_points,
  min(timestamp) as first_point,
  max(timestamp) as last_point
from
  tencentcloud_cvm_metric_cpu_hourly
group by
  instance_id
order by
  data_points desc;
```
