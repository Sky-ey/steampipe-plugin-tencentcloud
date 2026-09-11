---
title: "Steampipe Table: tencentcloud_clb_target_group_instance - Query Tencent Cloud CLB Target Group Instances using SQL"
description: "Query Tencent Cloud Load Balancer backend servers registered in target groups, including IPs, ports, and weights."
folder: "CLB"
---

# Table: tencentcloud_clb_target_group_instance - Query Tencent Cloud CLB Target Group Instances using SQL

Tencent Cloud Load Balancer (CLB) target group instances are the backend servers registered in target groups. The `tencentcloud_clb_target_group_instance` table uses the `DescribeTargetGroupInstances` API and returns the backend servers of target groups in each region, including their instance ID and name, service port, forwarding weight, IP addresses, and registration time.

## Table Usage Guide

The `tencentcloud_clb_target_group_instance` table in Steampipe provides information about the backend servers registered in CLB target groups within your Tencent Cloud account. This table allows you, as a DevOps engineer or cloud administrator, to query registered servers, including the target group, instance type and ID, port, weight, public and private IPs, and registration time. You can use this table to audit target group registrations, find misconfigured weights, and cross-reference with CVM instances.

## Examples

### Basic info

Get an overview of all backend servers registered in target groups.

```sql+postgres
select
  target_group_id,
  instance_id,
  instance_name,
  type,
  port,
  weight,
  region
from
  tencentcloud_clb_target_group_instance;
```

```sql+sqlite
select
  target_group_id,
  instance_id,
  instance_name,
  type,
  port,
  weight,
  region
from
  tencentcloud_clb_target_group_instance;
```

### List servers of a specific target group

The `target_group_id` qual is pushed down to the API as a filter.

```sql+postgres
select
  instance_id,
  instance_name,
  port,
  weight,
  private_ip_addresses,
  registered_time
from
  tencentcloud_clb_target_group_instance
where
  target_group_id = 'lbtg-12345678'
  and region = 'ap-guangzhou';
```

```sql+sqlite
select
  instance_id,
  instance_name,
  port,
  weight,
  private_ip_addresses,
  registered_time
from
  tencentcloud_clb_target_group_instance
where
  target_group_id = 'lbtg-12345678'
  and region = 'ap-guangzhou';
```

### Count servers per target group

```sql+postgres
select
  target_group_id,
  count(*) as server_count
from
  tencentcloud_clb_target_group_instance
group by
  target_group_id;
```

```sql+sqlite
select
  target_group_id,
  count(*) as server_count
from
  tencentcloud_clb_target_group_instance
group by
  target_group_id;
```

### Find zero-weight servers

Identify backend servers that receive no traffic due to a forwarding weight of zero.

```sql+postgres
select
  target_group_id,
  instance_id,
  instance_name,
  port,
  weight,
  region
from
  tencentcloud_clb_target_group_instance
where
  weight = 0;
```

```sql+sqlite
select
  target_group_id,
  instance_id,
  instance_name,
  port,
  weight,
  region
from
  tencentcloud_clb_target_group_instance
where
  weight = 0;
```

### Find servers of a specific instance

The `instance_id` qual is pushed down to the API as a filter.

```sql+postgres
select
  target_group_id,
  port,
  weight,
  public_ip_addresses,
  private_ip_addresses,
  region
from
  tencentcloud_clb_target_group_instance
where
  instance_id = 'ins-12345678';
```

```sql+sqlite
select
  target_group_id,
  port,
  weight,
  public_ip_addresses,
  private_ip_addresses,
  region
from
  tencentcloud_clb_target_group_instance
where
  instance_id = 'ins-12345678';
```
