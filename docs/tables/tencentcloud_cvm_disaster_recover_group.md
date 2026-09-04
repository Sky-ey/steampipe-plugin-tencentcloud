---
title: "Steampipe Table: tencentcloud_cvm_disaster_recover_group - Query Tencent Cloud Spread Placement Groups using SQL"
description: "Query Tencent Cloud CVM spread placement groups that distribute instances across physical hardware to reduce hardware-failure impact."
folder: "CVM"
---

# Table: tencentcloud_cvm_disaster_recover_group - Query Tencent Cloud Spread Placement Groups using SQL

Tencent Cloud spread placement groups (disaster recover groups) distribute CVM instances across different physical hardware (host, switch or rack) so that a single hardware failure affects as few instances as possible. Each group has a CVM quota, a current instance count and a type that determines the failure-domain granularity.

## Table Usage Guide

Use this table to inventory spread placement groups across regions, audit which groups are near their CVM quota, inspect the instance members of each group, and find groups by type (HOST / SW / RACK) for capacity planning.

## Examples

### Basic info

```sql+postgres
select
  disaster_recover_group_id,
  name,
  type,
  region
from tencentcloud_cvm_disaster_recover_group
order by region, name;
```
```sql+sqlite
select
  disaster_recover_group_id,
  name,
  type,
  region
from tencentcloud_cvm_disaster_recover_group
order by region, name;
```

### Quota usage and remaining capacity

```sql+postgres
select
  disaster_recover_group_id,
  name,
  type,
  cvm_quota_total,
  current_num,
  cvm_quota_total - current_num as remaining,
  region
from tencentcloud_cvm_disaster_recover_group
order by remaining asc;
```
```sql+sqlite
select
  disaster_recover_group_id,
  name,
  type,
  cvm_quota_total,
  current_num,
  cvm_quota_total - current_num as remaining,
  region
from tencentcloud_cvm_disaster_recover_group
order by remaining asc;
```

### Count groups by type

```sql+postgres
select
  type,
  count(*) as group_count,
  sum(cvm_quota_total) as total_quota,
  sum(current_num) as total_instances
from tencentcloud_cvm_disaster_recover_group
group by type
order by group_count desc;
```
```sql+sqlite
select
  type,
  count(*) as group_count,
  sum(cvm_quota_total) as total_quota,
  sum(current_num) as total_instances
from tencentcloud_cvm_disaster_recover_group
group by type
order by group_count desc;
```

### Groups that are near or at quota

```sql+postgres
select
  disaster_recover_group_id,
  name,
  type,
  cvm_quota_total,
  current_num,
  region
from tencentcloud_cvm_disaster_recover_group
where current_num >= cvm_quota_total;
```
```sql+sqlite
select
  disaster_recover_group_id,
  name,
  type,
  cvm_quota_total,
  current_num,
  region
from tencentcloud_cvm_disaster_recover_group
where current_num >= cvm_quota_total;
```

### Inspect instance members of a group

```sql+postgres
select
  disaster_recover_group_id,
  name,
  instance_ids
from tencentcloud_cvm_disaster_recover_group
where jsonb_array_length(instance_ids) > 0;
```
```sql+sqlite
select
  disaster_recover_group_id,
  name,
  instance_ids
from tencentcloud_cvm_disaster_recover_group
where json_array_length(instance_ids) > 0;
```
