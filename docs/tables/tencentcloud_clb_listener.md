---
title: "Steampipe Table: tencentcloud_clb_listener - Query Tencent Cloud CLB Listeners using SQL"
description: "Query Tencent Cloud Load Balancer listeners, including protocol, port, scheduling, health checks, and forwarding rules."
folder: "CLB"
---

# Table: tencentcloud_clb_listener - Query Tencent Cloud CLB Listeners using SQL

Tencent Cloud Load Balancer (CLB) listeners define how a load balancer accepts traffic, including the protocol, port, scheduling algorithm, health check, and forwarding rules. The `tencentcloud_clb_listener` table uses the `DescribeListeners` API to query listeners for load balancers in each region.

## Table Usage Guide

The `tencentcloud_clb_listener` table in Steampipe provides detailed information about CLB listeners within your Tencent Cloud account. This table allows you, as a DevOps engineer or cloud administrator, to query listener details, including the owning load balancer ID, listener protocol and port, scheduling algorithm, session persistence, health check configuration, and bound target groups. You can use this table to audit listener configuration, find HTTPS listeners, and review session or health check settings.

## Examples

### Basic info

Get an overview of all listeners, including their protocols, ports, and owning load balancers.

```sql+postgres
select
  load_balancer_id,
  listener_id,
  listener_name,
  protocol,
  port,
  region
from
  tencentcloud_clb_listener;
```

```sql+sqlite
select
  load_balancer_id,
  listener_id,
  listener_name,
  protocol,
  port,
  region
from
  tencentcloud_clb_listener;
```

### Find HTTPS listeners

List all HTTPS listeners that terminate TLS traffic.

```sql+postgres
select
  load_balancer_id,
  listener_id,
  listener_name,
  port,
  certificate,
  region
from
  tencentcloud_clb_listener
where
  protocol = 'HTTPS';
```

```sql+sqlite
select
  load_balancer_id,
  listener_id,
  listener_name,
  port,
  certificate,
  region
from
  tencentcloud_clb_listener
where
  protocol = 'HTTPS';
```

### Count listeners by protocol

```sql+postgres
select
  protocol,
  count(*) as listener_count
from
  tencentcloud_clb_listener
group by
  protocol
order by
  listener_count desc;
```

```sql+sqlite
select
  protocol,
  count(*) as listener_count
from
  tencentcloud_clb_listener
group by
  protocol
order by
  listener_count desc;
```

### Find listeners for a specific load balancer

The `load_balancer_id` qual is pushed down to the API.

```sql+postgres
select
  listener_id,
  listener_name,
  protocol,
  port,
  scheduler,
  session_expire_time,
  region
from
  tencentcloud_clb_listener
where
  load_balancer_id = 'lb-12345678'
  and region = 'ap-guangzhou';
```

```sql+sqlite
select
  listener_id,
  listener_name,
  protocol,
  port,
  scheduler,
  session_expire_time,
  region
from
  tencentcloud_clb_listener
where
  load_balancer_id = 'lb-12345678'
  and region = 'ap-guangzhou';
```

### Find listeners with health checks enabled

List listeners that have a health check configuration.

```sql+postgres
select
  load_balancer_id,
  listener_id,
  listener_name,
  protocol,
  port,
  health_check,
  region
from
  tencentcloud_clb_listener
where
  health_check is not null;
```

```sql+sqlite
select
  load_balancer_id,
  listener_id,
  listener_name,
  protocol,
  port,
  health_check,
  region
from
  tencentcloud_clb_listener
where
  health_check is not null;
```
