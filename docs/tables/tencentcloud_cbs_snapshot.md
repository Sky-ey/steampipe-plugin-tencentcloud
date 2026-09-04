---
title: "Steampipe Table: tencentcloud_cbs_snapshot - Query Tencent Cloud CBS Snapshots using SQL"
description: "Query Tencent Cloud CBS snapshots, including lifecycle, source disk, encryption, sharing, image associations, and tags."
folder: "CBS"
---

# Table: tencentcloud_cbs_snapshot - Query Tencent Cloud CBS Snapshots using SQL

The `tencentcloud_cbs_snapshot` table uses the `DescribeSnapshots` API to query CBS snapshots and their lifecycle, sharing, encryption, and source disk information.

## Examples

### Basic snapshot inventory

```sql+postgres
select
  snapshot_id,
  snapshot_name,
  snapshot_state,
  disk_id,
  disk_size,
  create_time,
  region
from
  tencentcloud_cbs_snapshot;
```

```sql+sqlite
select
  snapshot_id,
  snapshot_name,
  snapshot_state,
  disk_id,
  disk_size,
  create_time,
  region
from
  tencentcloud_cbs_snapshot;
```

### Find non-permanent snapshots nearing expiration

```sql+postgres
select
  snapshot_id,
  snapshot_name,
  deadline_time,
  disk_id,
  region
from
  tencentcloud_cbs_snapshot
where
  not is_permanent
  and deadline_time is not null
order by
  deadline_time;
```

```sql+sqlite
select
  snapshot_id,
  snapshot_name,
  deadline_time,
  disk_id,
  region
from
  tencentcloud_cbs_snapshot
where
  not is_permanent
  and deadline_time is not null
order by
  deadline_time;
```

### Find shared snapshots

```sql+postgres
select
  snapshot_id,
  snapshot_name,
  share_reference,
  time_start_share,
  region
from
  tencentcloud_cbs_snapshot
where
  share_reference > 0;
```

```sql+sqlite
select
  snapshot_id,
  snapshot_name,
  share_reference,
  time_start_share,
  region
from
  tencentcloud_cbs_snapshot
where
  share_reference > 0;
```

### Find snapshots by tag

```sql+postgres
select
  snapshot_id,
  snapshot_name,
  tags,
  region
from
  tencentcloud_cbs_snapshot,
  jsonb_each_text(tags) as kv
where
  kv.key = 'env'
  and kv.value = 'production';
```

```sql+sqlite
select
  snapshot_id,
  snapshot_name,
  tags,
  region
from
  tencentcloud_cbs_snapshot,
  json_each(tags) as kv
where
  kv.key = 'env'
  and kv.value = 'production';
```
