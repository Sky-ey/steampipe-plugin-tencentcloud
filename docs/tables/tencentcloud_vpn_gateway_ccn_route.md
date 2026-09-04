---
title: "Steampipe Table: tencentcloud_vpn_gateway_ccn_route - Query Tencent Cloud VPN gateway CCN routes using SQL"
description: "Allows users to query CCN routes propagated by a CCN-type VPN gateway (IDC IP ranges advertised to the CCN)."
folder: "VPN"
---

# Table: tencentcloud_vpn_gateway_ccn_route - Query Tencent Cloud VPN gateway CCN routes using SQL

A Tencent Cloud CCN-type VPN gateway propagates IDC IP ranges to the CCN as CCN routes. Each route entry records the route ID, the enabled status (ENABLE or DISABLE), and the destination CIDR block. This is a lightweight sub-resource with a minimal schema. The DescribeVpnGatewayCcnRoutes API lives in the VPC SDK package, but the resource is managed under the VPN product line.

## Table Usage Guide

The `tencentcloud_vpn_gateway_ccn_route` table provides you with detailed information about the CCN routes propagated by a CCN-type VPN gateway. This is a sub-resource table that fans out from a parent VPN gateway: pass `vpn_gateway_id` to list routes for a specific gateway, otherwise the plugin enumerates every VPN gateway in the region. Use this table as a network administrator to audit the IDC IP ranges that each CCN-type VPN gateway advertises to the CCN, and to verify which routes are enabled.

## Examples

### Basic info

List every VPN gateway CCN route and its core attributes.

```sql+postgres
select
  vpn_gateway_id,
  route_id,
  destination_cidr_block,
  status,
  region
from
  tencentcloud_vpn_gateway_ccn_route;
```

```sql+sqlite
select
  vpn_gateway_id,
  route_id,
  destination_cidr_block,
  status,
  region
from
  tencentcloud_vpn_gateway_ccn_route;
```

### List CCN routes for a specific VPN gateway

Filter routes by the parent VPN gateway ID. This is the recommended fan-out pattern for sub-resource tables.

```sql+postgres
select
  route_id,
  destination_cidr_block,
  status
from
  tencentcloud_vpn_gateway_ccn_route
where
  vpn_gateway_id = 'vpngw-f49l6u0z';
```

```sql+sqlite
select
  route_id,
  destination_cidr_block,
  status
from
  tencentcloud_vpn_gateway_ccn_route
where
  vpn_gateway_id = 'vpngw-f49l6u0z';
```

### Filter disabled CCN routes

Find VPN gateway CCN routes that are currently disabled.

```sql+postgres
select
  vpn_gateway_id,
  route_id,
  destination_cidr_block,
  region
from
  tencentcloud_vpn_gateway_ccn_route
where
  status = 'DISABLE';
```

```sql+sqlite
select
  vpn_gateway_id,
  route_id,
  destination_cidr_block,
  region
from
  tencentcloud_vpn_gateway_ccn_route
where
  status = 'DISABLE';
```

### Count CCN routes per VPN gateway

Summarize how many CCN routes each VPN gateway has propagated.

```sql+postgres
select
  vpn_gateway_id,
  count(*) as count
from
  tencentcloud_vpn_gateway_ccn_route
group by
  vpn_gateway_id
order by
  count desc;
```

```sql+sqlite
select
  vpn_gateway_id,
  count(*) as count
from
  tencentcloud_vpn_gateway_ccn_route
group by
  vpn_gateway_id
order by
  count desc;
```

### Inspect enabled CCN routes by region

List enabled CCN routes and group them by region.

```sql+postgres
select
  region,
  vpn_gateway_id,
  route_id,
  destination_cidr_block
from
  tencentcloud_vpn_gateway_ccn_route
where
  status = 'ENABLE'
order by
  region,
  vpn_gateway_id;
```

```sql+sqlite
select
  region,
  vpn_gateway_id,
  route_id,
  destination_cidr_block
from
  tencentcloud_vpn_gateway_ccn_route
where
  status = 'ENABLE'
order by
  region,
  vpn_gateway_id;
```
