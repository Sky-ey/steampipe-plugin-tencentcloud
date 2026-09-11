---
title: "Steampipe Table: tencentcloud_cbs_disk_backup - Query Tencent Cloud CBS Disk Backups using SQL"
description: "Query Tencent Cloud CBS disk backup points, including source disk, state, progress, size, and encryption details."
folder: "CBS"
---

# Table: tencentcloud_cbs_disk_backup - Query Tencent Cloud CBS Disk Backups using SQL

The `tencentcloud_cbs_disk_backup` table uses the `DescribeDiskBackups` API to query CBS disk backup points.

## Table Usage Guide

Use this table to audit CBS disk backup points: verify backup coverage per disk, check backup sizes and states, and trace each backup point back to its source disk. The `disk_backup_id` and `disk_id` quals push down to the `DescribeDiskBackups` Filters for efficient filtered listing and single-backup lookups.

## Examples

### Basic info

```sql+postgres
select
  disk_backup_id,
  disk_backup_name,
  disk_id,
  disk_size,
  disk_backup_state,
  create_time,
  region
from
  tencentcloud_cbs_disk_backup;
```

```sql+sqlite
select
  disk_backup_id,
  disk_backup_name,
  disk_id,
  disk_size,
  disk_backup_state,
  create_time,
  region
from
  tencentcloud_cbs_disk_backup;
```

### Find backup points that are not ready

```sql+postgres
select
  disk_backup_id,
  disk_id,
  disk_backup_state,
  percent,
  region
from
  tencentcloud_cbs_disk_backup
where
  disk_backup_state <> 'NORMAL';
```

```sql+sqlite
select
  disk_backup_id,
  disk_id,
  disk_backup_state,
  percent,
  region
from
  tencentcloud_cbs_disk_backup
where
  disk_backup_state <> 'NORMAL';
```

### Find backups for a specific disk

```sql+postgres
select
  disk_backup_id,
  disk_backup_name,
  disk_backup_state,
  percent,
  create_time,
  region
from
  tencentcloud_cbs_disk_backup
where
  disk_id = 'disk-12345678';
```

```sql+sqlite
select
  disk_backup_id,
  disk_backup_name,
  disk_backup_state,
  percent,
  create_time,
  region
from
  tencentcloud_cbs_disk_backup
where
  disk_id = 'disk-12345678';
```

### Find encrypted source disks

```sql+postgres
select
  disk_backup_id,
  disk_backup_name,
  disk_id,
  encrypt,
  region
from
  tencentcloud_cbs_disk_backup
where
  encrypt;
```

```sql+sqlite
select
  disk_backup_id,
  disk_backup_name,
  disk_id,
  encrypt,
  region
from
  tencentcloud_cbs_disk_backup
where
  encrypt = true;
```
