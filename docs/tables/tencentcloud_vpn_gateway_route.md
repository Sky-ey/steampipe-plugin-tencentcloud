---
title: "Steampipe Table: tencentcloud_vpn_gateway_route - Query Tencent Cloud VPN gateway routes using SQL"
description: "Allows users to query VPN gateway destination routes that route IDC IP ranges to VPN tunnels or CCN instances."
folder: "VPN"
---

# Table: tencentcloud_vpn_gateway_route - Query Tencent Cloud VPN gateway routes using SQL

A Tencent Cloud VPN gateway maintains destination routes that direct IDC IP ranges to either a VPN tunnel (VPNCONN) or a CCN instance. Each route records the destination CIDR block, the next hop type and instance ID, the route priority (0 or 100), the enabled status, and the route type (VPC, CCN, Static, or BGP). The DescribeVpnGatewayRoutes API lives in the VPC SDK package, but the resource is managed under the VPN product line.

## Table Usage Guide

The `tencentcloud_vpn_gateway_route` table provides you with detailed information about the destination routes on a VPN gateway. This is a sub-resource table that fans out from a parent VPN gateway: pass `vpn_gateway_id` to list routes for a specific gateway, otherwise the plugin enumerates every VPN gateway in the region. Use this table as a network administrator to audit destination CIDRs, inspect next-hop tunnels or CCN instances, verify route priorities, and review enabled vs disabled routes.

## Examples

### Basic info

List every VPN gateway route and its core attributes.

```sql+postgres
select
  vpn_gateway_id,
  route_id,
  destination_cidr_block,
  instance_type,
  instance_id,
  priority,
  status,
  type,
  region
from
  tencentcloud_vpn_gateway_route;
```

```sql+sqlite
select
  vpn_gateway_id,
  route_id,
  destination_cidr_block,
  instance_type,
  instance_id,
  priority,
  status,
  type,
  region
from
  tencentcloud_vpn_gateway_route;
```

### List routes for a specific VPN gateway

Filter routes by the parent VPN gateway ID. This is the recommended fan-out pattern for sub-resource tables.

```sql+postgres
select
  route_id,
  destination_cidr_block,
  instance_type,
  instance_id,
  priority,
  status,
  type
from
  tencentcloud_vpn_gateway_route
where
  vpn_gateway_id = 'vpngw-f49l6u0z';
```

```sql+sqlite
select
  route_id,
  destination_cidr_block,
  instance_type,
  instance_id,
  priority,
  status,
  type
from
  tencentcloud_vpn_gateway_route
where
  vpn_gateway_id = 'vpngw-f49l6u0z';
```

### Filter disabled routes

Find VPN gateway routes that are currently disabled.

```sql+postgres
select
  vpn_gateway_id,
  route_id,
  destination_cidr_block,
  instance_id,
  type,
  region
from
  tencentcloud_vpn_gateway_route
where
  status = 'DISABLE';
```

```sql+sqlite
select
  vpn_gateway_id,
  route_id,
  destination_cidr_block,
  instance_id,
  type,
  region
from
  tencentcloud_vpn_gateway_route
where
  status = 'DISABLE';
```

### Count routes by next-hop type per VPN gateway

Summarize how many routes of each next-hop type exist per VPN gateway.

```sql+postgres
select
  vpn_gateway_id,
  instance_type,
  count(*) as count
from
  tencentcloud_vpn_gateway_route
group by
  vpn_gateway_id,
  instance_type
order by
  vpn_gateway_id,
  count desc;
```

```sql+sqlite
select
  vpn_gateway_id,
  instance_type,
  count(*) as count
from
  tencentcloud_vpn_gateway_route
group by
  vpn_gateway_id,
  instance_type
order by
  vpn_gateway_id,
  count desc;
```

### Look up a route by destination CIDR block

Find the route that matches a specific destination CIDR block.

```sql+postgres
select
  vpn_gateway_id,
  route_id,
  destination_cidr_block,
  instance_type,
  instance_id,
  priority,
  status,
  region
from
  tencentcloud_vpn_gateway_route
where
  destination_cidr_block = '10.0.0.0/16';
```

```sql+sqlite
select
  vpn_gateway_id,
  route_id,
  destination_cidr_block,
  instance_type,
  instance_id,
  priority,
  status,
  region
from
  tencentcloud_vpn_gateway_route
where
  destination_cidr_block = '10.0.0.0/16';
```
