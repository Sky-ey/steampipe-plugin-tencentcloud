---
title: "Steampipe Table: tencentcloud_cbs_disk_storage_pool - Query Tencent Cloud CBS Dedicated Cluster Disk Storage Pools using SQL"
description: "Query Tencent Cloud CBS dedicated cluster disk storage pools, including state, zone, capacity, and expiry."
folder: "CBS"
---

# Table: tencentcloud_cbs_disk_storage_pool - Query Tencent Cloud CBS Dedicated Cluster Disk Storage Pools using SQL

The `tencentcloud_cbs_disk_storage_pool` table uses the `DescribeDiskStoragePool` API to query the dedicated cloud disk clusters (CDC storage pools) under the current account. Dedicated clusters provide exclusive, physically isolated CBS storage capacity for workloads with strict compliance or performance isolation requirements.

## Table Usage Guide

Each row represents a dedicated cloud disk cluster (storage pool). Use this table to inspect the state and zone of your dedicated clusters, check how much capacity is used, and identify clusters that are about to expire or have become unavailable for creating new disks.

## Examples

### List disk storage pools

```sql+postgres
select
  cdc_id,
  cdc_name,
  cdc_state,
  zone,
  disk_type,
  region
from
  tencentcloud_cbs_disk_storage_pool;
```

```sql+sqlite
select
  cdc_id,
  cdc_name,
  cdc_state,
  zone,
  disk_type,
  region
from
  tencentcloud_cbs_disk_storage_pool;
```

### Find storage pools that are not in a normal state

```sql+postgres
select
  cdc_id,
  cdc_name,
  cdc_state,
  zone,
  region
from
  tencentcloud_cbs_disk_storage_pool
where
  cdc_state <> 'NORMAL';
```

```sql+sqlite
select
  cdc_id,
  cdc_name,
  cdc_state,
  zone,
  region
from
  tencentcloud_cbs_disk_storage_pool
where
  cdc_state <> 'NORMAL';
```

### Check capacity usage of each storage pool

```sql+postgres
select
  cdc_id,
  cdc_name,
  disk_total,
  disk_available,
  disk_total - disk_available as disk_used,
  disk_number,
  region
from
  tencentcloud_cbs_disk_storage_pool;
```

```sql+sqlite
select
  cdc_id,
  cdc_name,
  disk_total,
  disk_available,
  disk_total - disk_available as disk_used,
  disk_number,
  region
from
  tencentcloud_cbs_disk_storage_pool;
```

### Find storage pools expiring within the next 30 days

```sql+postgres
select
  cdc_id,
  cdc_name,
  expired_time,
  disk_available,
  region
from
  tencentcloud_cbs_disk_storage_pool
where
  expired_time is not null
  and expired_time < now() + interval '30 days';
```

```sql+sqlite
select
  cdc_id,
  cdc_name,
  expired_time,
  disk_available,
  region
from
  tencentcloud_cbs_disk_storage_pool
where
  expired_time is not null
  and expired_time < datetime('now', '+30 days');
```

### Count storage pools by disk type

```sql+postgres
select
  disk_type,
  count(*) as pool_count,
  sum(disk_available) as total_available_gib
from
  tencentcloud_cbs_disk_storage_pool
group by
  disk_type;
```

```sql+sqlite
select
  disk_type,
  count(*) as pool_count,
  sum(disk_available) as total_available_gib
from
  tencentcloud_cbs_disk_storage_pool
group by
  disk_type;
```

### Filter storage pools in a specific availability zone

```sql+postgres
select
  cdc_id,
  cdc_name,
  zone,
  disk_type,
  disk_total,
  disk_available,
  region
from
  tencentcloud_cbs_disk_storage_pool
where
  zone = 'ap-singapore-1';
```

```sql+sqlite
select
  cdc_id,
  cdc_name,
  zone,
  disk_type,
  disk_total,
  disk_available,
  region
from
  tencentcloud_cbs_disk_storage_pool
where
  zone = 'ap-singapore-1';
```
