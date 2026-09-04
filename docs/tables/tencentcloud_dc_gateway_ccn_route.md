---
title: "Steampipe Table: tencentcloud_dc_gateway_ccn_route - Query Tencent Cloud Direct Connect gateway CCN routes using SQL"
description: "Allows users to query the IDC IP ranges (CCN routes) learned by each Tencent Cloud Direct Connect gateway."
folder: "DC"
---

# Table: tencentcloud_dc_gateway_ccn_route - Query Tencent Cloud Direct Connect gateway CCN routes using SQL

A Tencent Cloud Direct Connect gateway that attaches to a CCN learns IDC IP ranges from the on-premises network and publishes them to the CCN as CCN routes. Each route entry records the destination CIDR block, the BGP AS-Path attribute, the route-learning type (BGP or STATIC), and the last update time.

## Table Usage Guide

The `tencentcloud_dc_gateway_ccn_route` table provides you with detailed information about the CCN routes learned by each Direct Connect gateway. This is a sub-resource table that fans out from a parent Direct Connect gateway: pass `direct_connect_gateway_id` to list routes for a specific gateway, otherwise the plugin enumerates every Direct Connect gateway in the region. Use this table as a network administrator to audit IDC IP ranges advertised to the CCN, inspect BGP AS-Path attributes, and verify route-learning types.

## Examples

### Basic info

List every Direct Connect gateway CCN route and its core attributes.

```sql+postgres
select
  route_id,
  direct_connect_gateway_id,
  destination_cidr_block,
  ccn_route_type,
  as_path,
  region
from
  tencentcloud_dc_gateway_ccn_route;
```

```sql+sqlite
select
  route_id,
  direct_connect_gateway_id,
  destination_cidr_block,
  ccn_route_type,
  as_path,
  region
from
  tencentcloud_dc_gateway_ccn_route;
```

### List CCN routes for a specific Direct Connect gateway

Filter routes by the parent Direct Connect gateway ID. This is the recommended fan-out pattern for sub-resource tables.

```sql+postgres
select
  route_id,
  destination_cidr_block,
  as_path,
  description,
  update_time
from
  tencentcloud_dc_gateway_ccn_route
where
  direct_connect_gateway_id = 'dcg-9o233uri';
```

```sql+sqlite
select
  route_id,
  destination_cidr_block,
  as_path,
  description,
  update_time
from
  tencentcloud_dc_gateway_ccn_route
where
  direct_connect_gateway_id = 'dcg-9o233uri';
```

### Filter routes by learning type

List CCN routes that were learned dynamically via BGP.

```sql+postgres
select
  direct_connect_gateway_id,
  route_id,
  destination_cidr_block,
  as_path,
  region
from
  tencentcloud_dc_gateway_ccn_route
where
  ccn_route_type = 'BGP';
```

```sql+sqlite
select
  direct_connect_gateway_id,
  route_id,
  destination_cidr_block,
  as_path,
  region
from
  tencentcloud_dc_gateway_ccn_route
where
  ccn_route_type = 'BGP';
```

### Count CCN routes per Direct Connect gateway

Summarize how many CCN routes each Direct Connect gateway has learned.

```sql+postgres
select
  direct_connect_gateway_id,
  count(*) as count
from
  tencentcloud_dc_gateway_ccn_route
group by
  direct_connect_gateway_id
order by
  count desc;
```

```sql+sqlite
select
  direct_connect_gateway_id,
  count(*) as count
from
  tencentcloud_dc_gateway_ccn_route
group by
  direct_connect_gateway_id
order by
  count desc;
```

### Inspect routes with a non-empty AS-Path

Find BGP-learned routes whose AS-Path attribute is populated.

```sql+postgres
select
  direct_connect_gateway_id,
  route_id,
  destination_cidr_block,
  as_path,
  update_time,
  region
from
  tencentcloud_dc_gateway_ccn_route
where
  as_path is not null
  and jsonb_array_length(as_path) > 0;
```

```sql+sqlite
select
  direct_connect_gateway_id,
  route_id,
  destination_cidr_block,
  as_path,
  update_time,
  region
from
  tencentcloud_dc_gateway_ccn_route
where
  as_path is not null
  and json_array_length(as_path) > 0;
```
