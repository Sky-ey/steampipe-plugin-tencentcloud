---
title: "Steampipe Table: tencentcloud_clb_instance - Query Tencent Cloud CLB Instances using SQL"
description: "Query Tencent Cloud Load Balancer instances, including network type, VIPs, billing, security groups, and tags."
folder: "CLB"
---

# Table: tencentcloud_clb_instance - Query Tencent Cloud CLB Instances using SQL

Tencent Cloud Load Balancer (CLB) instances distribute incoming traffic across multiple backend servers. The `tencentcloud_clb_instance` table uses the `DescribeLoadBalancers` API and returns load balancers across all configured regions, or a single region when the `region` column is filtered.

## Table Usage Guide

The `tencentcloud_clb_instance` table in Steampipe provides detailed information about CLB instances within your Tencent Cloud account. This table allows you, as a DevOps engineer or cloud administrator, to query load balancer details, including the instance ID, name, network type, VIP addresses, billing mode, security groups, backend target count, and tags. You can use this table to audit load balancer inventory, find isolated or blocked instances, and review tag-based classifications.

## Examples

### Basic instance inventory

Get an overview of all load balancers in your account, including their names, network types, VIP addresses, and running status.

```sql+postgres
select
  load_balancer_id,
  load_balancer_name,
  load_balancer_type,
  load_balancer_vips,
  status,
  region
from
  tencentcloud_clb_instance;
```

```sql+sqlite
select
  load_balancer_id,
  load_balancer_name,
  load_balancer_type,
  load_balancer_vips,
  status,
  region
from
  tencentcloud_clb_instance;
```

### Find public load balancers

List all public (internet-facing) load balancers, which use the `OPEN` network type.

```sql+postgres
select
  load_balancer_id,
  load_balancer_name,
  load_balancer_vips,
  charge_type,
  region
from
  tencentcloud_clb_instance
where
  load_balancer_type = 'OPEN';
```

```sql+sqlite
select
  load_balancer_id,
  load_balancer_name,
  load_balancer_vips,
  charge_type,
  region
from
  tencentcloud_clb_instance
where
  load_balancer_type = 'OPEN';
```

### Count load balancers by status

```sql+postgres
select
  status,
  count(*) as lb_count
from
  tencentcloud_clb_instance
group by
  status
order by
  lb_count desc;
```

```sql+sqlite
select
  status,
  count(*) as lb_count
from
  tencentcloud_clb_instance
group by
  status
order by
  lb_count desc;
```

### Find load balancers by tag

```sql+postgres
select
  load_balancer_id,
  load_balancer_name,
  tags,
  region
from
  tencentcloud_clb_instance,
  jsonb_each_text(tags) as kv
where
  kv.key = 'env'
  and kv.value = 'production';
```

```sql+sqlite
select
  load_balancer_id,
  load_balancer_name,
  tags,
  region
from
  tencentcloud_clb_instance,
  json_each(tags) as kv
where
  kv.key = 'env'
  and kv.value = 'production';
```

### Find a specific load balancer

Look up a single load balancer by its ID. The `load_balancer_id` qual is pushed down to the API as `LoadBalancerIds`.

```sql+postgres
select
  load_balancer_id,
  load_balancer_name,
  load_balancer_type,
  status,
  charge_type,
  vpc_id,
  subnet_id,
  target_count,
  region
from
  tencentcloud_clb_instance
where
  load_balancer_id = 'lb-12345678'
  and region = 'ap-guangzhou';
```

```sql+sqlite
select
  load_balancer_id,
  load_balancer_name,
  load_balancer_type,
  status,
  charge_type,
  vpc_id,
  subnet_id,
  target_count,
  region
from
  tencentcloud_clb_instance
where
  load_balancer_id = 'lb-12345678'
  and region = 'ap-guangzhou';
```
