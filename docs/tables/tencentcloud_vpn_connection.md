---
title: "Steampipe Table: tencentcloud_vpn_connection - Query Tencent Cloud VPN connections using SQL"
description: "Allows users to query VPN connections (IPsec tunnels) that bind a VPN gateway to a customer gateway."
folder: "VPN"
---

# Table: tencentcloud_vpn_connection - Query Tencent Cloud VPN connections using SQL

A Tencent Cloud VPN connection is an IPsec tunnel that binds a VPN gateway to a customer gateway. Each connection records the IKE and IPsec options, the SPD (Security Policy Database) set, the negotiation type, the pre-shared key, the route type, the lifecycle state, the connectivity (net) status, and optional health check and DPD (Dead Peer Detection) configuration. The DescribeVpnConnections API lives in the VPC SDK package, but the resource is managed under the VPN product line.

## Table Usage Guide

The `tencentcloud_vpn_connection` table provides you with detailed information about VPN tunnels. Use this table as a network administrator to inventory tunnels, audit the VPN gateway and customer gateway that each tunnel binds, inspect IKE/IPsec options, verify health check and DPD configuration, and review tag-based classifications. Pass `vpn_connection_id` for a direct lookup, or filter by `vpc_id`, `vpn_gateway_id`, `customer_gateway_id`, or `vpn_connection_name`.

## Examples

### Basic info

List every VPN connection and its core attributes.

```sql+postgres
select
  vpn_connection_id,
  vpn_connection_name,
  vpn_gateway_id,
  customer_gateway_id,
  state,
  net_status,
  route_type,
  region
from
  tencentcloud_vpn_connection;
```

```sql+sqlite
select
  vpn_connection_id,
  vpn_connection_name,
  vpn_gateway_id,
  customer_gateway_id,
  state,
  net_status,
  route_type,
  region
from
  tencentcloud_vpn_connection;
```

### Filter VPN connections by state

Identify VPN tunnels that are currently available.

```sql+postgres
select
  vpn_connection_id,
  vpn_connection_name,
  vpn_gateway_id,
  customer_gateway_id,
  net_status,
  region
from
  tencentcloud_vpn_connection
where
  state = 'AVAILABLE';
```

```sql+sqlite
select
  vpn_connection_id,
  vpn_connection_name,
  vpn_gateway_id,
  customer_gateway_id,
  net_status,
  region
from
  tencentcloud_vpn_connection
where
  state = 'AVAILABLE';
```

### List tunnels for a specific VPN gateway

Filter VPN connections by the parent VPN gateway ID.

```sql+postgres
select
  vpn_connection_id,
  vpn_connection_name,
  customer_gateway_id,
  route_type,
  net_status,
  health_check_status
from
  tencentcloud_vpn_connection
where
  vpn_gateway_id = 'vpngw-f49l6u0z';
```

```sql+sqlite
select
  vpn_connection_id,
  vpn_connection_name,
  customer_gateway_id,
  route_type,
  net_status,
  health_check_status
from
  tencentcloud_vpn_connection
where
  vpn_gateway_id = 'vpngw-f49l6u0z';
```

### Count VPN connections by state per region

Summarize how many VPN connections of each state exist per region.

```sql+postgres
select
  region,
  state,
  count(*) as count
from
  tencentcloud_vpn_connection
group by
  region,
  state
order by
  region,
  count desc;
```

```sql+sqlite
select
  region,
  state,
  count(*) as count
from
  tencentcloud_vpn_connection
group by
  region,
  state
order by
  region,
  count desc;
```

### List VPN connections by tag

Filter VPN connections by a specific tag key-value pair.

```sql+postgres
select
  vpn_connection_id,
  vpn_connection_name,
  vpn_gateway_id,
  tags
from
  tencentcloud_vpn_connection
where
  tags ->> 'Environment' = 'Production';
```

```sql+sqlite
select
  vpn_connection_id,
  vpn_connection_name,
  vpn_gateway_id,
  tags
from
  tencentcloud_vpn_connection
where
  json_extract(tags, '$.Environment') = 'Production';
```
