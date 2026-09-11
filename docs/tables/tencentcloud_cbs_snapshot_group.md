---
title: "Steampipe Table: tencentcloud_cbs_snapshot_group - Query Tencent Cloud CBS Snapshot Groups using SQL"
description: "Query Tencent Cloud CBS snapshot groups, including state, member snapshots, image associations, and retention."
folder: "CBS"
---

# Table: tencentcloud_cbs_snapshot_group - Query Tencent Cloud CBS Snapshot Groups using SQL

The `tencentcloud_cbs_snapshot_group` table uses the `DescribeSnapshotGroups` API to query groups of related CBS snapshots.

## Table Usage Guide

Use this table to work with snapshot groups, which tie together the snapshots of all disks attached to an instance at a point in time — useful for auditing consistent multi-disk backups and planning instance-level restores. The `snapshot_group_id` and `snapshot_id` quals push down to the `DescribeSnapshotGroups` Filters for efficient filtered listing and single-group lookups.

## Examples

### Basic info

```sql+postgres
select
  snapshot_group_id,
  snapshot_group_name,
  snapshot_group_type,
  snapshot_group_state,
  snapshot_id_set,
  create_time,
  region
from
  tencentcloud_cbs_snapshot_group;
```

```sql+sqlite
select
  snapshot_group_id,
  snapshot_group_name,
  snapshot_group_type,
  snapshot_group_state,
  snapshot_id_set,
  create_time,
  region
from
  tencentcloud_cbs_snapshot_group;
```

### Find snapshot groups that are not ready

```sql+postgres
select
  snapshot_group_id,
  snapshot_group_name,
  snapshot_group_state,
  percent,
  region
from
  tencentcloud_cbs_snapshot_group
where
  snapshot_group_state <> 'NORMAL';
```

```sql+sqlite
select
  snapshot_group_id,
  snapshot_group_name,
  snapshot_group_state,
  percent,
  region
from
  tencentcloud_cbs_snapshot_group
where
  snapshot_group_state <> 'NORMAL';
```

### Find expiring snapshot groups

```sql+postgres
select
  snapshot_group_id,
  snapshot_group_name,
  deadline_time,
  snapshot_id_set,
  region
from
  tencentcloud_cbs_snapshot_group
where
  not is_permanent
  and deadline_time is not null;
```

```sql+sqlite
select
  snapshot_group_id,
  snapshot_group_name,
  deadline_time,
  snapshot_id_set,
  region
from
  tencentcloud_cbs_snapshot_group
where
  not is_permanent
  and deadline_time is not null;
```

### Count snapshot groups by state

```sql+postgres
select
  snapshot_group_state,
  count(*) as group_count
from
  tencentcloud_cbs_snapshot_group
group by
  snapshot_group_state;
```

```sql+sqlite
select
  snapshot_group_state,
  count(*) as group_count
from
  tencentcloud_cbs_snapshot_group
group by
  snapshot_group_state;
```
