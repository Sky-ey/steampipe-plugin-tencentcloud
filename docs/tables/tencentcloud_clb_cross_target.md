---
title: "Steampipe Table: tencentcloud_clb_cross_target - Query Tencent Cloud CLB Cross Targets using SQL"
description: "Query Tencent Cloud Load Balancer cross-VPC backend targets available through CCN cross-domain 2.0."
folder: "CLB"
---

# Table: tencentcloud_clb_cross_target - Query Tencent Cloud CLB Cross Targets using SQL

Tencent Cloud Load Balancer (CLB) cross targets are the CVM and ENI instances in other VPCs or regions that can serve as backends through CCN cross-domain 2.0. The `tencentcloud_clb_cross_target` table uses the `DescribeCrossTargets` API and returns the candidate backend instances in each region, including their VPC, IP, instance ID, and region.

## Table Usage Guide

The `tencentcloud_clb_cross_target` table in Steampipe provides information about cross-VPC backend candidates within your Tencent Cloud account. This table allows you, as a DevOps engineer or cloud administrator, to query which CVMs and ENIs are available as backends for cross-domain load balancing, including their local and remote VPCs, IP addresses, and regions. You can use this table to plan CCN cross-domain 2.0 deployments and audit which instances are reachable across VPCs.

## Examples

### Basic cross target inventory

Get an overview of all cross-VPC backend candidates.

```sql+postgres
select
  local_vpc_id,
  vpc_id,
  vpc_name,
  ip,
  instance_id,
  instance_name,
  target_region
from
  tencentcloud_clb_cross_target;
```

```sql+sqlite
select
  local_vpc_id,
  vpc_id,
  vpc_name,
  ip,
  instance_id,
  instance_name,
  target_region
from
  tencentcloud_clb_cross_target;
```

### Find targets in a specific VPC

The `vpc_id` qual is pushed down to the API as a filter.

```sql+postgres
select
  local_vpc_id,
  vpc_id,
  vpc_name,
  ip,
  instance_id,
  instance_name,
  target_region
from
  tencentcloud_clb_cross_target
where
  vpc_id = 'vpc-12345678';
```

```sql+sqlite
select
  local_vpc_id,
  vpc_id,
  vpc_name,
  ip,
  instance_id,
  instance_name,
  target_region
from
  tencentcloud_clb_cross_target
where
  vpc_id = 'vpc-12345678';
```

### Count cross targets per remote VPC

```sql+postgres
select
  vpc_id,
  vpc_name,
  count(*) as target_count
from
  tencentcloud_clb_cross_target
group by
  vpc_id,
  vpc_name
order by
  target_count desc;
```

```sql+sqlite
select
  vpc_id,
  vpc_name,
  count(*) as target_count
from
  tencentcloud_clb_cross_target
group by
  vpc_id,
  vpc_name
order by
  target_count desc;
```

### Find targets by IP

The `ip` qual is pushed down to the API as a filter.

```sql+postgres
select
  local_vpc_id,
  vpc_id,
  instance_id,
  instance_name,
  eni_id,
  target_region
from
  tencentcloud_clb_cross_target
where
  ip = '192.168.0.1';
```

```sql+sqlite
select
  local_vpc_id,
  vpc_id,
  instance_id,
  instance_name,
  eni_id,
  target_region
from
  tencentcloud_clb_cross_target
where
  ip = '192.168.0.1';
```

### Find targets by listener or location

The `listener_id` and `location_id` quals are pushed down to the API as filters (DescribeCrossTargets `listener-id` / `location-id`). These columns are filter-only — they echo the qual value you supplied and are not part of the DescribeCrossTargets response payload.

```sql+postgres
select
  local_vpc_id,
  vpc_id,
  instance_id,
  instance_name,
  target_region,
  listener_id,
  location_id
from
  tencentcloud_clb_cross_target
where
  listener_id = 'lbl-xxxxxxxx';
```

```sql+sqlite
select
  local_vpc_id,
  vpc_id,
  instance_id,
  instance_name,
  target_region,
  listener_id,
  location_id
from
  tencentcloud_clb_cross_target
where
  listener_id = 'lbl-xxxxxxxx';
```

### Count cross targets per region

```sql+postgres
select
  target_region,
  count(*) as target_count
from
  tencentcloud_clb_cross_target
group by
  target_region;
```

```sql+sqlite
select
  target_region,
  count(*) as target_count
from
  tencentcloud_clb_cross_target
group by
  target_region;
```
