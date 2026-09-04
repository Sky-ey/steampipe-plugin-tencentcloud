---
title: "Steampipe Table: tencentcloud_dc_direct_connect_tunnel - Query Tencent Cloud Direct Connect tunnels using SQL"
description: "Allows users to query Direct Connect dedicated tunnels that connect a Direct Connect to a VPC, BMVPC, or CCN."
folder: "DC"
---

# Table: tencentcloud_dc_direct_connect_tunnel - Query Tencent Cloud Direct Connect tunnels using SQL

A Tencent Cloud Direct Connect dedicated tunnel is a logical connection that rides on top of a physical Direct Connect and terminates on a VPC, BMVPC, or CCN via a Direct Connect gateway. Each tunnel records the network type, routing type (BGP or static), VLAN, bandwidth, Tencent-side and customer-side IP addresses, BGP peer configuration, and route filter prefixes.

## Table Usage Guide

The `tencentcloud_dc_direct_connect_tunnel` table provides you with detailed information about dedicated tunnels. Use this table as a network administrator to inventory tunnels, audit the parent Direct Connect for each tunnel, inspect BGP vs static routing, verify VLAN and bandwidth assignments, and review tag-based classifications. Pass `direct_connect_id` to filter tunnels by parent physical connection.

## Examples

### Basic info

List every Direct Connect tunnel and its core attributes.

```sql+postgres
select
  direct_connect_tunnel_id,
  direct_connect_tunnel_name,
  state,
  direct_connect_id,
  network_type,
  route_type,
  vlan,
  bandwidth,
  region
from
  tencentcloud_dc_direct_connect_tunnel;
```

```sql+sqlite
select
  direct_connect_tunnel_id,
  direct_connect_tunnel_name,
  state,
  direct_connect_id,
  network_type,
  route_type,
  vlan,
  bandwidth,
  region
from
  tencentcloud_dc_direct_connect_tunnel;
```

### List tunnels for a specific Direct Connect

Filter tunnels by the parent Direct Connect ID.

```sql+postgres
select
  direct_connect_tunnel_id,
  direct_connect_tunnel_name,
  network_type,
  route_type,
  direct_connect_gateway_id,
  state
from
  tencentcloud_dc_direct_connect_tunnel
where
  direct_connect_id = 'dc-abcdefgh';
```

```sql+sqlite
select
  direct_connect_tunnel_id,
  direct_connect_tunnel_name,
  network_type,
  route_type,
  direct_connect_gateway_id,
  state
from
  tencentcloud_dc_direct_connect_tunnel
where
  direct_connect_id = 'dc-abcdefgh';
```

### Filter tunnels by routing type

List tunnels that use BGP routing.

```sql+postgres
select
  direct_connect_tunnel_id,
  direct_connect_tunnel_name,
  direct_connect_id,
  bgp_peer,
  region
from
  tencentcloud_dc_direct_connect_tunnel
where
  route_type = 'BGP';
```

```sql+sqlite
select
  direct_connect_tunnel_id,
  direct_connect_tunnel_name,
  direct_connect_id,
  bgp_peer,
  region
from
  tencentcloud_dc_direct_connect_tunnel
where
  route_type = 'BGP';
```

### Count tunnels by network type per region

Summarize how many tunnels of each network type exist per region.

```sql+postgres
select
  region,
  network_type,
  count(*) as count
from
  tencentcloud_dc_direct_connect_tunnel
group by
  region,
  network_type
order by
  region,
  count desc;
```

```sql+sqlite
select
  region,
  network_type,
  count(*) as count
from
  tencentcloud_dc_direct_connect_tunnel
group by
  region,
  network_type
order by
  region,
  count desc;
```

### List tunnels by tag

Filter Direct Connect tunnels by a specific tag key-value pair.

```sql+postgres
select
  direct_connect_tunnel_id,
  direct_connect_tunnel_name,
  tags
from
  tencentcloud_dc_direct_connect_tunnel
where
  tags ->> 'Environment' = 'Production';
```

```sql+sqlite
select
  direct_connect_tunnel_id,
  direct_connect_tunnel_name,
  tags
from
  tencentcloud_dc_direct_connect_tunnel
where
  json_extract(tags, '$.Environment') = 'Production';
```
