---
title: "Steampipe Table: tencentcloud_clb_target - Query Tencent Cloud CLB Targets using SQL"
description: "Query Tencent Cloud Load Balancer backend targets, including bound instances, ports, weights, and IP addresses."
folder: "CLB"
---

# Table: tencentcloud_clb_target - Query Tencent Cloud CLB Targets using SQL

Tencent Cloud Load Balancer (CLB) targets are the backend servers bound to listeners or forwarding rules. The `tencentcloud_clb_target` table uses the `DescribeTargets` API to query backend targets for each load balancer and listener. Because a target has no independent ID, this table expands targets from their parent load balancer and listener; it has no Get hydrate, and queries are filtered through the `load_balancer_id`, `listener_id`, `listener_protocol`, and `listener_port` quals.

## Table Usage Guide

The `tencentcloud_clb_target` table in Steampipe provides detailed information about CLB backend targets within your Tencent Cloud account. This table allows you, as a DevOps engineer or cloud administrator, to query backend target details, including the owning load balancer and listener, target type, instance ID, service port, forwarding weight, and public or private IP addresses. You can use this table to audit backend bindings, find low-weight targets, and identify ENI-based backends.

## Examples

### Basic info

Get an overview of all backend targets, including their bound load balancers, listeners, instance IDs, and ports.

```sql+postgres
select
  load_balancer_id,
  listener_id,
  listener_protocol,
  listener_port,
  instance_id,
  instance_name,
  port,
  weight,
  region
from
  tencentcloud_clb_target;
```

```sql+sqlite
select
  load_balancer_id,
  listener_id,
  listener_protocol,
  listener_port,
  instance_id,
  instance_name,
  port,
  weight,
  region
from
  tencentcloud_clb_target;
```

### Find targets for a specific listener

The `load_balancer_id` and `listener_id` quals are pushed down to the API.

```sql+postgres
select
  load_balancer_id,
  listener_id,
  instance_id,
  instance_name,
  port,
  weight,
  private_ip_addresses,
  public_ip_addresses,
  registered_time,
  region
from
  tencentcloud_clb_target
where
  load_balancer_id = 'lb-12345678'
  and listener_id = 'lbl-12345678'
  and region = 'ap-guangzhou';
```

```sql+sqlite
select
  load_balancer_id,
  listener_id,
  instance_id,
  instance_name,
  port,
  weight,
  private_ip_addresses,
  public_ip_addresses,
  registered_time,
  region
from
  tencentcloud_clb_target
where
  load_balancer_id = 'lb-12345678'
  and listener_id = 'lbl-12345678'
  and region = 'ap-guangzhou';
```

### Find ENI-based targets

List backend targets that are elastic network interfaces rather than CVM instances.

```sql+postgres
select
  load_balancer_id,
  listener_id,
  eni_id,
  port,
  weight,
  region
from
  tencentcloud_clb_target
where
  target_type = 'ENI';
```

```sql+sqlite
select
  load_balancer_id,
  listener_id,
  eni_id,
  port,
  weight,
  region
from
  tencentcloud_clb_target
where
  target_type = 'ENI';
```

### Find targets with zero weight

List backend targets that are bound but currently receive no traffic.

```sql+postgres
select
  load_balancer_id,
  listener_id,
  instance_id,
  instance_name,
  port,
  region
from
  tencentcloud_clb_target
where
  weight = 0;
```

```sql+sqlite
select
  load_balancer_id,
  listener_id,
  instance_id,
  instance_name,
  port,
  region
from
  tencentcloud_clb_target
where
  weight = 0;
```

### Find targets for a specific backend instance

```sql+postgres
select
  load_balancer_id,
  listener_id,
  listener_protocol,
  listener_port,
  port,
  weight,
  region
from
  tencentcloud_clb_target
where
  instance_id = 'ins-12345678';
```

```sql+sqlite
select
  load_balancer_id,
  listener_id,
  listener_protocol,
  listener_port,
  port,
  weight,
  region
from
  tencentcloud_clb_target
where
  instance_id = 'ins-12345678';
```
