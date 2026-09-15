---
title: "Steampipe Table: tencentcloud_clb_target_group - Query Tencent Cloud CLB Target Groups using SQL"
description: "Query Tencent Cloud Load Balancer target groups, including protocol, port, health checks, and bound instances."
folder: "CLB"
---

# Table: tencentcloud_clb_target_group - Query Tencent Cloud CLB Target Groups using SQL

Tencent Cloud Load Balancer (CLB) target groups group backend servers so they can be shared across listeners and load balancers. The `tencentcloud_clb_target_group` table uses the `DescribeTargetGroups` API and returns target groups in each region, including their protocol, port, health check configuration, and associated rules. The `associated_rule` column carries the listener/rule bindings (load balancer, listener, location, domain, URL and port) so you can see where each target group is referenced without querying the CLB listener APIs.

## Table Usage Guide

The `tencentcloud_clb_target_group` table in Steampipe provides detailed information about CLB target groups within your Tencent Cloud account. This table allows you, as a DevOps engineer or cloud administrator, to query target group details, including the target group ID and name, backend protocol and port, VPC, scheduling algorithm, session persistence, health check configuration, and associated listener rules. You can use this table to audit target group configuration, find legacy v1 groups, review health check settings, and trace which load balancers and rules reference each target group.

## Examples

### Basic info

Get an overview of all target groups, including their names, backend protocols, ports, and VPCs.

```sql+postgres
select
  target_group_id,
  target_group_name,
  target_group_type,
  protocol,
  port,
  vpc_id,
  region
from
  tencentcloud_clb_target_group;
```

```sql+sqlite
select
  target_group_id,
  target_group_name,
  target_group_type,
  protocol,
  port,
  vpc_id,
  region
from
  tencentcloud_clb_target_group;
```

### Find v2 target groups

List current-generation (v2) target groups, which support advanced scheduling and health check features.

```sql+postgres
select
  target_group_id,
  target_group_name,
  protocol,
  port,
  schedule_algorithm,
  session_expire_time,
  region
from
  tencentcloud_clb_target_group
where
  target_group_type = 'v2';
```

```sql+sqlite
select
  target_group_id,
  target_group_name,
  protocol,
  port,
  schedule_algorithm,
  session_expire_time,
  region
from
  tencentcloud_clb_target_group
where
  target_group_type = 'v2';
```

### Count target groups by type

```sql+postgres
select
  target_group_type,
  count(*) as group_count
from
  tencentcloud_clb_target_group
group by
  target_group_type;
```

```sql+sqlite
select
  target_group_type,
  count(*) as group_count
from
  tencentcloud_clb_target_group
group by
  target_group_type;
```

### Find target groups in a specific VPC

The `vpc_id` qual is pushed down to the API as a filter.

```sql+postgres
select
  target_group_id,
  target_group_name,
  protocol,
  port,
  region
from
  tencentcloud_clb_target_group
where
  vpc_id = 'vpc-12345678'
  and region = 'ap-guangzhou';
```

```sql+sqlite
select
  target_group_id,
  target_group_name,
  protocol,
  port,
  region
from
  tencentcloud_clb_target_group
where
  vpc_id = 'vpc-12345678'
  and region = 'ap-guangzhou';
```

### Find target groups by tag

```sql+postgres
select
  target_group_id,
  target_group_name,
  tags,
  region
from
  tencentcloud_clb_target_group,
  jsonb_each_text(tags) as kv
where
  kv.key = 'env'
  and kv.value = 'production';
```

```sql+sqlite
select
  target_group_id,
  target_group_name,
  tags,
  region
from
  tencentcloud_clb_target_group,
  json_each(tags) as kv
where
  kv.key = 'env'
  and kv.value = 'production';
```
