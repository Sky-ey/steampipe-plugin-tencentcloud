---
title: "Steampipe Table: tencentcloud_cbs_auto_snapshot_policy - Query Tencent Cloud CBS Auto Snapshot Policies using SQL"
description: "Query Tencent Cloud CBS auto snapshot policies, including schedules, retention, bindings, replication, and tags."
folder: "CBS"
---

# Table: tencentcloud_cbs_auto_snapshot_policy - Query Tencent Cloud CBS Auto Snapshot Policies using SQL

The `tencentcloud_cbs_auto_snapshot_policy` table uses the `DescribeAutoSnapshotPolicies` API to query scheduled snapshot policies and their retention and resource bindings.

## Examples

### List auto snapshot policies

```sql+postgres
select
  auto_snapshot_policy_id,
  auto_snapshot_policy_name,
  auto_snapshot_policy_state,
  is_activated,
  next_trigger_time,
  retention_days,
  region
from
  tencentcloud_cbs_auto_snapshot_policy;
```

```sql+sqlite
select
  auto_snapshot_policy_id,
  auto_snapshot_policy_name,
  auto_snapshot_policy_state,
  is_activated,
  next_trigger_time,
  retention_days,
  region
from
  tencentcloud_cbs_auto_snapshot_policy;
```

### Find inactive policies

```sql+postgres
select
  auto_snapshot_policy_id,
  auto_snapshot_policy_name,
  disk_id_set,
  instance_id_set,
  region
from
  tencentcloud_cbs_auto_snapshot_policy
where
  not is_activated;
```

```sql+sqlite
select
  auto_snapshot_policy_id,
  auto_snapshot_policy_name,
  disk_id_set,
  instance_id_set,
  region
from
  tencentcloud_cbs_auto_snapshot_policy
where
  not is_activated;
```

### Review execution schedules

```sql+postgres
select
  auto_snapshot_policy_id,
  auto_snapshot_policy_name,
  policy,
  advanced_retention_policy,
  region
from
  tencentcloud_cbs_auto_snapshot_policy;
```

```sql+sqlite
select
  auto_snapshot_policy_id,
  auto_snapshot_policy_name,
  policy,
  advanced_retention_policy,
  region
from
  tencentcloud_cbs_auto_snapshot_policy;
```

### Count policies by state

```sql+postgres
select
  auto_snapshot_policy_state,
  count(*) as policy_count
from
  tencentcloud_cbs_auto_snapshot_policy
group by
  auto_snapshot_policy_state;
```

```sql+sqlite
select
  auto_snapshot_policy_state,
  count(*) as policy_count
from
  tencentcloud_cbs_auto_snapshot_policy
group by
  auto_snapshot_policy_state;
```
