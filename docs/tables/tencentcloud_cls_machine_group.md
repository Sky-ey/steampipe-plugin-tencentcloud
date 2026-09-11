---
title: "Steampipe Table: tencentcloud_cls_machine_group - Query Tencent Cloud CLS Machine Groups using SQL"
description: "The tencentcloud_cls_machine_group table queries Cloud Log Service (CLS) machine groups, collections of machines (by IP or label) running the LogListener agent that collect and ship logs to CLS topics."
folder: "CLS"
---

# Table: tencentcloud_cls_machine_group - Query Tencent Cloud CLS Machine Groups using SQL

A CLS machine group is a collection of machines (identified by IP address or by label) that run the LogListener agent. The group defines the collection behavior — including the LogListener auto-upgrade window, service logging, and offline-machine cleanup time — and is referenced by collection rules to ship logs to CLS topics.

## Table Usage Guide

Use this table to inventory CLS machine groups across regions, inspect group membership type (IP vs label), audit LogListener auto-upgrade windows, and verify service-logging and cleanup configuration. The `machine_group_id` and `machine_group_name` quals push down to the `DescribeMachineGroups` Filters for efficient single-group or filtered listing.

## Examples

### Basic info

```sql+postgres
select
  machine_group_id,
  machine_group_name,
  os_type,
  machine_group_type,
  auto_update,
  create_time,
  region
from tencentcloud_cls_machine_group
order by create_time desc;
```

```sql+sqlite
select
  machine_group_id,
  machine_group_name,
  os_type,
  machine_group_type,
  auto_update,
  create_time,
  region
from tencentcloud_cls_machine_group
order by create_time desc;
```

### Query a specific machine group by ID

```sql+postgres
select
  machine_group_id,
  machine_group_name,
  machine_group_type,
  machine_group_values,
  service_logging,
  delay_cleanup_time,
  tags
from tencentcloud_cls_machine_group
where machine_group_id = '4ae9b3c1-xxxx-xxxx-xxxx-xxxxxxxxxxxx';
```

```sql+sqlite
select
  machine_group_id,
  machine_group_name,
  machine_group_type,
  machine_group_values,
  service_logging,
  delay_cleanup_time,
  tags
from tencentcloud_cls_machine_group
where machine_group_id = '4ae9b3c1-xxxx-xxxx-xxxx-xxxxxxxxxxxx';
```

### Find machine groups with auto-update enabled

```sql+postgres
select
  machine_group_id,
  machine_group_name,
  auto_update,
  update_start_time,
  update_end_time,
  region
from tencentcloud_cls_machine_group
where auto_update = 'true'
order by region;
```

```sql+sqlite
select
  machine_group_id,
  machine_group_name,
  auto_update,
  update_start_time,
  update_end_time,
  region
from tencentcloud_cls_machine_group
where auto_update = 'true'
order by region;
```

### Count machine groups by operating system and region

```sql+postgres
select
  region,
  case os_type when 0 then 'Linux' when 1 then 'Windows' else 'Unknown' end as os,
  count(*) as group_count
from tencentcloud_cls_machine_group
group by region, os_type
order by group_count desc;
```

```sql+sqlite
select
  region,
  case os_type when 0 then 'Linux' when 1 then 'Windows' else 'Unknown' end as os,
  count(*) as group_count
from tencentcloud_cls_machine_group
group by region, os_type
order by group_count desc;
```

### List machine groups with service logging enabled

```sql+postgres
select
  machine_group_id,
  machine_group_name,
  service_logging,
  delay_cleanup_time,
  meta_tags,
  region
from tencentcloud_cls_machine_group
where service_logging = true;
```

```sql+sqlite
select
  machine_group_id,
  machine_group_name,
  service_logging,
  delay_cleanup_time,
  meta_tags,
  region
from tencentcloud_cls_machine_group
where service_logging = true;
```
