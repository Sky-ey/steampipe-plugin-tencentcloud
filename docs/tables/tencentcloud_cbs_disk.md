---
title: "Steampipe Table: tencentcloud_cbs_disk - Query Tencent Cloud CBS Disks using SQL"
description: "Query Tencent Cloud Block Storage disks, including state, capacity, billing, attachment, encryption, snapshots, backups, and tags."
folder: "CBS"
---

# Table: tencentcloud_cbs_disk - Query Tencent Cloud CBS Disks using SQL

Tencent Cloud Block Storage (CBS) provides persistent block storage for cloud workloads. The `tencentcloud_cbs_disk` table uses the `DescribeDisks` API and returns disks across all regions, or a single region when the `region` column is filtered.

## Examples

### Basic disk inventory

```sql+postgres
select
  disk_id,
  disk_name,
  disk_type,
  disk_usage,
  disk_size,
  disk_state,
  instance_id,
  region
from
  tencentcloud_cbs_disk;
```

```sql+sqlite
select
  disk_id,
  disk_name,
  disk_type,
  disk_usage,
  disk_size,
  disk_state,
  instance_id,
  region
from
  tencentcloud_cbs_disk;
```

### Find unattached disks

```sql+postgres
select
  disk_id,
  disk_name,
  disk_size,
  disk_charge_type,
  create_time,
  region
from
  tencentcloud_cbs_disk
where
  disk_state = 'UNATTACHED';
```

```sql+sqlite
select
  disk_id,
  disk_name,
  disk_size,
  disk_charge_type,
  create_time,
  region
from
  tencentcloud_cbs_disk
where
  disk_state = 'UNATTACHED';
```

### Find unencrypted data disks

```sql+postgres
select
  disk_id,
  disk_name,
  disk_type,
  disk_size,
  region
from
  tencentcloud_cbs_disk
where
  disk_usage = 'DATA_DISK'
  and not encrypt;
```

```sql+sqlite
select
  disk_id,
  disk_name,
  disk_type,
  disk_size,
  region
from
  tencentcloud_cbs_disk
where
  disk_usage = 'DATA_DISK'
  and not encrypt;
```

### Count disks by state

```sql+postgres
select
  disk_state,
  count(*) as disk_count
from
  tencentcloud_cbs_disk
group by
  disk_state
order by
  disk_count desc;
```

```sql+sqlite
select
  disk_state,
  count(*) as disk_count
from
  tencentcloud_cbs_disk
group by
  disk_state
order by
  disk_count desc;
```

### Find disks by tag

```sql+postgres
select
  disk_id,
  disk_name,
  tags,
  region
from
  tencentcloud_cbs_disk,
  jsonb_each_text(tags) as kv
where
  kv.key = 'env'
  and kv.value = 'production';
```

```sql+sqlite
select
  disk_id,
  disk_name,
  tags,
  region
from
  tencentcloud_cbs_disk,
  json_each(tags) as kv
where
  kv.key = 'env'
  and kv.value = 'production';
```
