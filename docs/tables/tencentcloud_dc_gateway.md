---
title: "Steampipe Table: tencentcloud_dc_gateway - Query Tencent Cloud Direct Connect gateways using SQL"
description: "Allows users to query Direct Connect gateways that bridge Direct Connect tunnels with VPC or CCN networks."
folder: "DC"
---

# Table: tencentcloud_dc_gateway - Query Tencent Cloud Direct Connect gateways using SQL

A Tencent Cloud Direct Connect gateway bridges Direct Connect tunnels with a VPC or a CCN. Each gateway records the gateway type (NORMAL or NAT), the network type (VPC or CCN), the associated VPC or CCN ID, the CCN route-learning type (BGP or STATIC), and optional NAT and BGP community settings. The DescribeDirectConnectGateways API lives in the VPC SDK package, but the resource is managed under the Direct Connect product line.

## Table Usage Guide

The `tencentcloud_dc_gateway` table provides you with detailed information about Direct Connect gateways. Use this table as a network administrator to inventory gateways, audit gateway and network types, verify VPC/CCN associations, inspect BGP configuration, and review zone and high-availability placement. Pass `direct_connect_gateway_id` for a direct lookup, or filter by `vpc_id`, `ccn_id`, `gateway_type`, or `network_type`.

## Examples

### Basic info

List every Direct Connect gateway and its core attributes.

```sql+postgres
select
  direct_connect_gateway_id,
  direct_connect_gateway_name,
  gateway_type,
  network_type,
  vpc_id,
  ccn_id,
  ccn_route_type,
  region
from
  tencentcloud_dc_gateway;
```

```sql+sqlite
select
  direct_connect_gateway_id,
  direct_connect_gateway_name,
  gateway_type,
  network_type,
  vpc_id,
  ccn_id,
  ccn_route_type,
  region
from
  tencentcloud_dc_gateway;
```

### Filter gateways by network type

List Direct Connect gateways that attach to a CCN.

```sql+postgres
select
  direct_connect_gateway_id,
  direct_connect_gateway_name,
  ccn_id,
  ccn_route_type,
  enable_bgp,
  region
from
  tencentcloud_dc_gateway
where
  network_type = 'CCN';
```

```sql+sqlite
select
  direct_connect_gateway_id,
  direct_connect_gateway_name,
  ccn_id,
  ccn_route_type,
  enable_bgp,
  region
from
  tencentcloud_dc_gateway
where
  network_type = 'CCN';
```

### Count gateways by type per region

Summarize how many gateways of each type exist per region.

```sql+postgres
select
  region,
  gateway_type,
  count(*) as count
from
  tencentcloud_dc_gateway
group by
  region,
  gateway_type
order by
  region,
  count desc;
```

```sql+sqlite
select
  region,
  gateway_type,
  count(*) as count
from
  tencentcloud_dc_gateway
group by
  region,
  gateway_type
order by
  region,
  count desc;
```

### Inspect NAT-type Direct Connect gateways

Find gateways that are configured as NAT gateways.

```sql+postgres
select
  direct_connect_gateway_id,
  direct_connect_gateway_name,
  nat_gateway_id,
  direct_connect_gateway_ip,
  zone,
  region
from
  tencentcloud_dc_gateway
where
  gateway_type = 'NAT';
```

```sql+sqlite
select
  direct_connect_gateway_id,
  direct_connect_gateway_name,
  nat_gateway_id,
  direct_connect_gateway_ip,
  zone,
  region
from
  tencentcloud_dc_gateway
where
  gateway_type = 'NAT';
```

### List Direct Connect gateways attached to a specific VPC

Filter gateways by the associated VPC ID.

```sql+postgres
select
  direct_connect_gateway_id,
  direct_connect_gateway_name,
  gateway_type,
  network_type,
  enable_bgp,
  region
from
  tencentcloud_dc_gateway
where
  vpc_id = 'vpc-xxxxxxxx';
```

```sql+sqlite
select
  direct_connect_gateway_id,
  direct_connect_gateway_name,
  gateway_type,
  network_type,
  enable_bgp,
  region
from
  tencentcloud_dc_gateway
where
  vpc_id = 'vpc-xxxxxxxx';
```
