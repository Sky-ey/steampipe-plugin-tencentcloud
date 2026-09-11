---
title: "Steampipe Table: tencentcloud_clb_block_ip - Query Tencent Cloud CLB Blocked IPs using SQL"
description: "Query Tencent Cloud Load Balancer blocked IP (blacklist) entries, including block and expiration times."
folder: "CLB"
---

# Table: tencentcloud_clb_block_ip - Query Tencent Cloud CLB Blocked IPs using SQL

Tencent Cloud Load Balancer (CLB) IP blocklists let you blacklist client IPs at the load balancer level to mitigate abuse. The `tencentcloud_clb_block_ip` table uses the `DescribeBlockIPList` API and returns the blocked IP entries of each load balancer in every region, including the blocked IP, block time, and expiration time.

## Table Usage Guide

The `tencentcloud_clb_block_ip` table in Steampipe provides information about the blocked IP entries of CLB load balancers within your Tencent Cloud account. This table allows you, as a DevOps engineer or cloud administrator, to query blacklist entries, including the owning load balancer, the blocked IP, the real client IP field, and the block/expiration times. You can use this table to audit IP blocklists, find expired blocks, and review which load balancers have active blacklists.

## Examples

### Basic info

Get an overview of all blocked IP entries across load balancers.

```sql+postgres
select
  load_balancer_id,
  ip,
  create_time,
  expire_time,
  region
from
  tencentcloud_clb_block_ip;
```

```sql+sqlite
select
  load_balancer_id,
  ip,
  create_time,
  expire_time,
  region
from
  tencentcloud_clb_block_ip;
```

### List blocked IPs of a specific load balancer

The `load_balancer_id` qual is pushed down to the API.

```sql+postgres
select
  ip,
  client_ip_field,
  create_time,
  expire_time
from
  tencentcloud_clb_block_ip
where
  load_balancer_id = 'lb-12345678'
  and region = 'ap-guangzhou';
```

```sql+sqlite
select
  ip,
  client_ip_field,
  create_time,
  expire_time
from
  tencentcloud_clb_block_ip
where
  load_balancer_id = 'lb-12345678'
  and region = 'ap-guangzhou';
```

### Count blocked IPs per load balancer

```sql+postgres
select
  load_balancer_id,
  count(*) as blocked_count
from
  tencentcloud_clb_block_ip
group by
  load_balancer_id
order by
  blocked_count desc;
```

```sql+sqlite
select
  load_balancer_id,
  count(*) as blocked_count
from
  tencentcloud_clb_block_ip
group by
  load_balancer_id
order by
  blocked_count desc;
```

### Find expired blocks

Identify blocked IPs whose block has already expired and may no longer be enforced.

```sql+postgres
select
  load_balancer_id,
  ip,
  create_time,
  expire_time,
  region
from
  tencentcloud_clb_block_ip
where
  expire_time < now();
```

```sql+sqlite
select
  load_balancer_id,
  ip,
  create_time,
  expire_time,
  region
from
  tencentcloud_clb_block_ip
where
  expire_time < datetime('now');
```

### Find blocks for a specific IP

Search for a particular client IP across all load balancers.

```sql+postgres
select
  load_balancer_id,
  create_time,
  expire_time,
  region
from
  tencentcloud_clb_block_ip
where
  ip = '203.0.113.10';
```

```sql+sqlite
select
  load_balancer_id,
  create_time,
  expire_time,
  region
from
  tencentcloud_clb_block_ip
where
  ip = '203.0.113.10';
```
