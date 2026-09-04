---
title: "Steampipe Table: tencentcloud_ccn_route - Query Tencent Cloud CCN routes using SQL"
description: "Allows users to query routing policies that route traffic between associated instances of a Tencent Cloud CCN."
folder: "CCN"
---

# Table: tencentcloud_ccn_route - Query Tencent Cloud CCN routes using SQL

A Tencent Cloud CCN (Cloud Connect Network) maintains a routing table that directs traffic between its associated instances. Each route entry records the destination CIDR block, the next hop (an associated VPC or Direct Connect instance), whether the route is dynamic (BGP) or static, the route priority, and the enabled state.

## Table Usage Guide

The `tencentcloud_ccn_route` table provides you with detailed information about the routing policies of a CCN. This is a sub-resource table that fans out from a parent CCN: pass `ccn_id` to list routes for a specific CCN, otherwise the plugin enumerates every CCN in the account. Use this table as a network administrator to audit destination CIDRs, inspect next-hop instances, verify BGP vs static routes, and review route priorities.

## Examples

### Basic info

List every CCN route and its core attributes.

```sql+postgres
select
  ccn_id,
  route_id,
  destination_cidr_block,
  instance_type,
  instance_id,
  enabled,
  is_bgp,
  route_priority
from
  tencentcloud_ccn_route;
```

```sql+sqlite
select
  ccn_id,
  route_id,
  destination_cidr_block,
  instance_type,
  instance_id,
  enabled,
  is_bgp,
  route_priority
from
  tencentcloud_ccn_route;
```

### List routes for a specific CCN

Filter routes by the parent CCN ID. This is the recommended fan-out pattern for sub-resource tables.

```sql+postgres
select
  route_id,
  destination_cidr_block,
  instance_type,
  instance_id,
  instance_region,
  enabled,
  route_priority
from
  tencentcloud_ccn_route
where
  ccn_id = 'ccn-f49l6u0z';
```

```sql+sqlite
select
  route_id,
  destination_cidr_block,
  instance_type,
  instance_id,
  instance_region,
  enabled,
  route_priority
from
  tencentcloud_ccn_route
where
  ccn_id = 'ccn-f49l6u0z';
```

### Inspect disabled routes

Find CCN routes that are currently disabled.

```sql+postgres
select
  ccn_id,
  route_id,
  destination_cidr_block,
  instance_id,
  extra_state
from
  tencentcloud_ccn_route
where
  not enabled;
```

```sql+sqlite
select
  ccn_id,
  route_id,
  destination_cidr_block,
  instance_id,
  extra_state
from
  tencentcloud_ccn_route
where
  enabled = 0;
```

### List dynamic (BGP) CCN routes

Find CCN routes that were learned dynamically via BGP.

```sql+postgres
select
  ccn_id,
  route_id,
  destination_cidr_block,
  instance_id,
  route_priority
from
  tencentcloud_ccn_route
where
  is_bgp;
```

```sql+sqlite
select
  ccn_id,
  route_id,
  destination_cidr_block,
  instance_id,
  route_priority
from
  tencentcloud_ccn_route
where
  is_bgp = 1;
```

### Count routes per CCN by next-hop type

Summarize how many routes of each next-hop type exist per CCN.

```sql+postgres
select
  ccn_id,
  instance_type,
  count(*) as count
from
  tencentcloud_ccn_route
group by
  ccn_id,
  instance_type
order by
  ccn_id,
  count desc;
```

```sql+sqlite
select
  ccn_id,
  instance_type,
  count(*) as count
from
  tencentcloud_ccn_route
group by
  ccn_id,
  instance_type
order by
  ccn_id,
  count desc;
```
