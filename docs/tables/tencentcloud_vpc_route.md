---
title: "Steampipe Table: tencentcloud_vpc_route - Query Tencent Cloud VPC routing entries using SQL"
description: "Allows users to query individual VPC routing policies expanded from each route table."
folder: "VPC"
---

# Table: tencentcloud_vpc_route - Query Tencent Cloud VPC routing entries using SQL

Each VPC route table contains a list of routing policies (`RouteSet`). The `tencentcloud_vpc_route` table expands those entries into one row per routing policy, so you can filter, sort, and aggregate individual routes without parsing the JSON column of `tencentcloud_vpc_route_table`. The Tencent Cloud SDK does not expose a dedicated `DescribeRoutes` API, so this table is implemented as an expansion view of `DescribeRouteTables`.

## Table Usage Guide

The `tencentcloud_vpc_route` table provides you with detailed information about individual routing policies. Use this table to inventory routing entries across all route tables, find routes with a given next hop, audit gateway types, or filter routes by `route_table_id` and `vpc_id`.

## Examples

### Basic info

List every routing policy with its key attributes.

```sql+postgres
select
  route_table_id,
  vpc_id,
  destination_cidr_block,
  gateway_type,
  gateway_id,
  route_type,
  enabled
from
  tencentcloud_vpc_route;
```

```sql+sqlite
select
  route_table_id,
  vpc_id,
  destination_cidr_block,
  gateway_type,
  gateway_id,
  route_type,
  enabled
from
  tencentcloud_vpc_route;
```

### List routes of a specific route table

Filter routing policies by their parent route table ID.

```sql+postgres
select
  route_id,
  route_item_id,
  destination_cidr_block,
  gateway_type,
  gateway_id,
  enabled
from
  tencentcloud_vpc_route
where
  route_table_id = 'rtb-xxxxxxxx';
```

```sql+sqlite
select
  route_id,
  route_item_id,
  destination_cidr_block,
  gateway_type,
  gateway_id,
  enabled
from
  tencentcloud_vpc_route
where
  route_table_id = 'rtb-xxxxxxxx';
```

### Find routes that forward to a NAT gateway

Identify routing policies whose next hop is a NAT gateway.

```sql+postgres
select
  route_table_id,
  vpc_id,
  destination_cidr_block,
  gateway_id
from
  tencentcloud_vpc_route
where
  gateway_type = 'NAT';
```

```sql+sqlite
select
  route_table_id,
  vpc_id,
  destination_cidr_block,
  gateway_id
from
  tencentcloud_vpc_route
where
  gateway_type = 'NAT';
```

### Find disabled routes

Identify routing policies that have been disabled.

```sql+postgres
select
  route_table_id,
  vpc_id,
  destination_cidr_block,
  gateway_type,
  gateway_id
from
  tencentcloud_vpc_route
where
  not enabled;
```

```sql+sqlite
select
  route_table_id,
  vpc_id,
  destination_cidr_block,
  gateway_type,
  gateway_id
from
  tencentcloud_vpc_route
where
  enabled = 0;
```

### Count routes by gateway type

Summarize how many routing policies reference each next-hop type.

```sql+postgres
select
  gateway_type,
  count(*) as count
from
  tencentcloud_vpc_route
group by
  gateway_type
order by
  count desc;
```

```sql+sqlite
select
  gateway_type,
  count(*) as count
from
  tencentcloud_vpc_route
group by
  gateway_type
order by
  count desc;
```
