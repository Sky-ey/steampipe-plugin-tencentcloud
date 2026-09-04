---
title: "Steampipe Table: tencentcloud_cvm_auto_scaling_group - Query Tencent Cloud Auto Scaling Groups using SQL"
description: "Query Tencent Cloud Auto Scaling (AS) groups and their scaling policies, instance counts and load balancer associations."
folder: "CVM"
---

# Table: tencentcloud_cvm_auto_scaling_group - Query Tencent Cloud Auto Scaling Groups using SQL

Tencent Cloud Auto Scaling (AS) groups automatically add or remove CVM instances based on configured scaling policies, launch configurations and scheduled tasks. Each group maintains a desired/min/max instance count, health checks and termination policies, and can be associated with CLB load balancers and multi-AZ subnets.

## Table Usage Guide

Use this table to inventory auto scaling groups across regions, audit scaling policies and capacity bounds, find groups with abnormal status or disabled scaling, and inspect launch configuration and VPC associations.

## Examples

### Basic info

```sql+postgres
select
  auto_scaling_group_id,
  auto_scaling_group_name,
  auto_scaling_group_status,
  enabled_status,
  region
from tencentcloud_cvm_auto_scaling_group
order by region, auto_scaling_group_name;
```
```sql+sqlite
select
  auto_scaling_group_id,
  auto_scaling_group_name,
  auto_scaling_group_status,
  enabled_status,
  region
from tencentcloud_cvm_auto_scaling_group
order by region, auto_scaling_group_name;
```

### Find groups with abnormal status or disabled scaling

```sql+postgres
select
  auto_scaling_group_id,
  auto_scaling_group_name,
  auto_scaling_group_status,
  enabled_status,
  in_activity_status,
  region
from tencentcloud_cvm_auto_scaling_group
where auto_scaling_group_status <> 'NORMAL'
  or enabled_status = 'DISABLED';
```
```sql+sqlite
select
  auto_scaling_group_id,
  auto_scaling_group_name,
  auto_scaling_group_status,
  enabled_status,
  in_activity_status,
  region
from tencentcloud_cvm_auto_scaling_group
where auto_scaling_group_status <> 'NORMAL'
  or enabled_status = 'DISABLED';
```

### Capacity and instance count overview

```sql+postgres
select
  auto_scaling_group_id,
  auto_scaling_group_name,
  min_size,
  max_size,
  desired_capacity,
  instance_count,
  in_service_instance_count,
  region
from tencentcloud_cvm_auto_scaling_group
order by region, desired_capacity desc;
```
```sql+sqlite
select
  auto_scaling_group_id,
  auto_scaling_group_name,
  min_size,
  max_size,
  desired_capacity,
  instance_count,
  in_service_instance_count,
  region
from tencentcloud_cvm_auto_scaling_group
order by region, desired_capacity desc;
```

### Scaling policy and health check details

```sql+postgres
select
  auto_scaling_group_id,
  auto_scaling_group_name,
  launch_configuration_id,
  health_check_type,
  instance_allocation_policy,
  retry_policy,
  multi_zone_subnet_policy,
  termination_policy_set,
  region
from tencentcloud_cvm_auto_scaling_group;
```
```sql+sqlite
select
  auto_scaling_group_id,
  auto_scaling_group_name,
  launch_configuration_id,
  health_check_type,
  instance_allocation_policy,
  retry_policy,
  multi_zone_subnet_policy,
  termination_policy_set,
  region
from tencentcloud_cvm_auto_scaling_group;
```

### Count groups by launch configuration

```sql+postgres
select
  launch_configuration_id,
  count(*) as group_count
from tencentcloud_cvm_auto_scaling_group
group by launch_configuration_id
order by group_count desc;
```
```sql+sqlite
select
  launch_configuration_id,
  count(*) as group_count
from tencentcloud_cvm_auto_scaling_group
group by launch_configuration_id
order by group_count desc;
```

### Filter by tag

```sql+postgres
select
  auto_scaling_group_id,
  auto_scaling_group_name,
  tags
from tencentcloud_cvm_auto_scaling_group
where tags ->> 'env' = 'production';
```
```sql+sqlite
select
  auto_scaling_group_id,
  auto_scaling_group_name,
  tags
from tencentcloud_cvm_auto_scaling_group
where json_extract(tags, '$.env') = 'production';
```
