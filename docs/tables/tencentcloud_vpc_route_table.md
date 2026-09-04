---
title: "Steampipe Table: tencentcloud_vpc_route_table - Query Tencent Cloud VPC route tables using SQL"
description: "Allows users to query VPC route tables to inspect their routing policies and subnet associations."
folder: "VPC"
---

# Table: tencentcloud_vpc_route_table - Query Tencent Cloud VPC route tables using SQL

A route table controls how traffic is forwarded within a VPC. Each VPC has one default (main) route table and may have additional custom route tables that are explicitly associated with subnets. Route tables contain routing policies (`route_set`) that map destination CIDR blocks to next hops.

## Table Usage Guide

The `tencentcloud_vpc_route_table` table provides you with detailed information about route tables within your Tencent Cloud account. Use this table to inventory route tables per VPC, identify the main route table, inspect subnet associations (`association_set`), and audit routing policies (`route_set`) as JSON.

## Examples

### Basic info

List every route table with its key attributes.

```sql+postgres
select
  route_table_id,
  route_table_name,
  vpc_id,
  main,
  region
from
  tencentcloud_vpc_route_table;
```

```sql+sqlite
select
  route_table_id,
  route_table_name,
  vpc_id,
  main,
  region
from
  tencentcloud_vpc_route_table;
```

### Find the main route table of each VPC

Identify the default route table that Tencent Cloud provisions per VPC.

```sql+postgres
select
  route_table_id,
  route_table_name,
  vpc_id,
  region
from
  tencentcloud_vpc_route_table
where
  main;
```

```sql+sqlite
select
  route_table_id,
  route_table_name,
  vpc_id,
  region
from
  tencentcloud_vpc_route_table
where
  main = 1;
```

### List route tables of a specific VPC

Filter route tables by their parent VPC ID.

```sql+postgres
select
  route_table_id,
  route_table_name,
  association_set
from
  tencentcloud_vpc_route_table
where
  vpc_id = 'vpc-xxxxxxxx';
```

```sql+sqlite
select
  route_table_id,
  route_table_name,
  association_set
from
  tencentcloud_vpc_route_table
where
  vpc_id = 'vpc-xxxxxxxx';
```

### Inspect routing policies of a route table

Expand the `route_set` JSON column to view individual routing entries.

```sql+postgres
select
  route_table_id,
  route ->> 'DestinationCidrBlock' as destination,
  route ->> 'GatewayType' as gateway_type,
  route ->> 'GatewayId' as gateway_id,
  route ->> 'RouteType' as route_type,
  (route ->> 'Enabled')::bool as enabled
from
  tencentcloud_vpc_route_table,
  jsonb_array_elements(route_set) as route
where
  route_table_id = 'rtb-xxxxxxxx';
```

```sql+sqlite
select
  route_table_id,
  json_extract(route.value, '$.DestinationCidrBlock') as destination,
  json_extract(route.value, '$.GatewayType') as gateway_type,
  json_extract(route.value, '$.GatewayId') as gateway_id,
  json_extract(route.value, '$.RouteType') as route_type,
  json_extract(route.value, '$.Enabled') as enabled
from
  tencentcloud_vpc_route_table,
  json_each(route_set) as route
where
  route_table_id = 'rtb-xxxxxxxx';
```

### Count route tables by VPC

Summarize how many route tables each VPC has.

```sql+postgres
select
  vpc_id,
  count(*) as count
from
  tencentcloud_vpc_route_table
group by
  vpc_id
order by
  count desc;
```

```sql+sqlite
select
  vpc_id,
  count(*) as count
from
  tencentcloud_vpc_route_table
group by
  vpc_id
order by
  count desc;
```
