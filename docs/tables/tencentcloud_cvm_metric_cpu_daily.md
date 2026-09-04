---
title: "Steampipe Table: tencentcloud_cvm_metric_cpu_daily - Query CVM instance daily CPU utilization using SQL"
description: "Query the daily CPU utilization metric of Tencent Cloud CVM instances over the last 30 days."
folder: "CVM"
---

# Table: tencentcloud_cvm_metric_cpu_daily - Query CVM instance daily CPU utilization using SQL

Tencent Cloud Monitor (Cloud Monitor) collects per-instance CPU utilization for every CVM instance. The `tencentcloud_cvm_metric_cpu_daily` table uses the `GetMonitorData` API (namespace `QCE/CVM`, metric `CPUUsage`, period 86400 seconds) and returns one row per instance per day over the last 30 days, across all regions or a single region when the `region` column is filtered.

Each row represents a single daily data point of one instance: the instance ID, the timestamp of the data point, and the CPU utilization value (a percentage).

## Table Usage Guide

Use this table to spot long-running CPU trends, find consistently overloaded or idle instances, and feed capacity-planning reports. Because daily granularity is retained for 186 days, this table is suited to month-scale analysis. For finer-grained, recent troubleshooting use `tencentcloud_cvm_metric_cpu_hourly`.

Filtering by `instance_id` queries only that instance; omitting it fans out across every instance in the region.

By default the table reads the last 30 days. A `timestamp` predicate in the `WHERE` clause narrows the time range and is pushed down to the `GetMonitorData` API (`StartTime`/`EndTime`), so you can query any window within the 186-day retention period instead of the default lookback.

## Examples

### Basic daily CPU utilization

```sql+postgres
select
  instance_id,
  timestamp,
  value,
  region
from
  tencentcloud_cvm_metric_cpu_daily
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
  tencentcloud_cvm_metric_cpu_daily
order by
  instance_id,
  timestamp;
```

### Daily CPU utilization over a custom time range

The `timestamp` range is pushed down to the `GetMonitorData` API, so only the requested window is read.

```sql+postgres
select
  instance_id,
  timestamp,
  value
from
  tencentcloud_cvm_metric_cpu_daily
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
  tencentcloud_cvm_metric_cpu_daily
where
  timestamp >= '2026-07-01T00:00:00Z'
  and timestamp < '2026-08-01T00:00:00Z'
order by
  instance_id,
  timestamp;
```

### Daily CPU utilization of a single instance

```sql+postgres
select
  instance_id,
  timestamp,
  value
from
  tencentcloud_cvm_metric_cpu_daily
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
  tencentcloud_cvm_metric_cpu_daily
where
  instance_id = 'ins-xxxxxxxx'
order by
  timestamp;
```

### Average daily CPU utilization per instance

```sql+postgres
select
  instance_id,
  round(avg(value)::numeric, 2) as avg_cpu_utilization,
  round(max(value)::numeric, 2) as max_cpu_utilization
from
  tencentcloud_cvm_metric_cpu_daily
group by
  instance_id
order by
  avg_cpu_utilization desc;
```

```sql+sqlite
select
  instance_id,
  round(avg(value), 2) as avg_cpu_utilization,
  round(max(value), 2) as max_cpu_utilization
from
  tencentcloud_cvm_metric_cpu_daily
group by
  instance_id
order by
  avg_cpu_utilization desc;
```

### Find instances with high daily CPU utilization

```sql+postgres
select
  instance_id,
  timestamp,
  value,
  region
from
  tencentcloud_cvm_metric_cpu_daily
where
  value > 80
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
  tencentcloud_cvm_metric_cpu_daily
where
  value > 80
order by
  value desc;
```

### Count instances by average daily CPU band

```sql+postgres
select
  case
    when avg(value) >= 80 then 'high (>=80)'
    when avg(value) >= 50 then 'medium (50-80)'
    else 'low (<50)'
  end as cpu_band,
  count(distinct instance_id) as instance_count
from
  tencentcloud_cvm_metric_cpu_daily
group by
  instance_id;
```

```sql+sqlite
select
  case
    when avg(value) >= 80 then 'high (>=80)'
    when avg(value) >= 50 then 'medium (50-80)'
    else 'low (<50)'
  end as cpu_band,
  count(distinct instance_id) as instance_count
from
  tencentcloud_cvm_metric_cpu_daily
group by
  instance_id;
```
