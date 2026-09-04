---
title: "Steampipe Table: tencentcloud_vpc_nat_gateway - Query Tencent Cloud VPC NAT gateways using SQL"
description: "Allows users to query VPC NAT gateways to inspect bandwidth, public IPs, and SNAT/DNAT rules."
folder: "VPC"
---

# Table: tencentcloud_vpc_nat_gateway - Query Tencent Cloud VPC NAT gateways using SQL

A NAT gateway provides SNAT (source network address translation) and DNAT (destination network address translation) for resources in a private subnet so they can reach the internet or accept inbound port forwarding. Each NAT gateway carries one or more elastic public IPs, a bandwidth cap, and optional sets of SNAT/DNAT rules.

## Table Usage Guide

The `tencentcloud_vpc_nat_gateway` table provides you with detailed information about NAT gateways within your Tencent Cloud account. Use this table to inventory NAT gateways per VPC, audit state and bandwidth, inspect bound public IPs, and review the SNAT/DNAT rule sets exposed as JSON columns.

## Examples

### Basic info

List every NAT gateway with its key attributes.

```sql+postgres
select
  nat_gateway_id,
  nat_gateway_name,
  vpc_id,
  state,
  network_state,
  internet_max_bandwidth_out,
  region
from
  tencentcloud_vpc_nat_gateway;
```

```sql+sqlite
select
  nat_gateway_id,
  nat_gateway_name,
  vpc_id,
  state,
  network_state,
  internet_max_bandwidth_out,
  region
from
  tencentcloud_vpc_nat_gateway;
```

### Find NAT gateways of a specific VPC

Filter NAT gateways by their parent VPC ID.

```sql+postgres
select
  nat_gateway_id,
  nat_gateway_name,
  zone,
  nat_product_version,
  public_ip_address_set
from
  tencentcloud_vpc_nat_gateway
where
  vpc_id = 'vpc-xxxxxxxx';
```

```sql+sqlite
select
  nat_gateway_id,
  nat_gateway_name,
  zone,
  nat_product_version,
  public_ip_address_set
from
  tencentcloud_vpc_nat_gateway
where
  vpc_id = 'vpc-xxxxxxxx';
```

### Inspect DNAT (port forwarding) rules

Expand the `destination_ip_port_translation_nat_rule_set` JSON column to view individual port forwarding rules.

```sql+postgres
select
  nat_gateway_id,
  rule ->> 'IpProtocol' as protocol,
  rule ->> 'PublicIpAddress' as public_ip,
  rule ->> 'PublicPort' as public_port,
  rule ->> 'PrivateIpAddress' as private_ip,
  rule ->> 'PrivatePort' as private_port,
  rule ->> 'Description' as description
from
  tencentcloud_vpc_nat_gateway,
  jsonb_array_elements(destination_ip_port_translation_nat_rule_set) as rule;
```

```sql+sqlite
select
  nat_gateway_id,
  json_extract(rule.value, '$.IpProtocol') as protocol,
  json_extract(rule.value, '$.PublicIpAddress') as public_ip,
  json_extract(rule.value, '$.PublicPort') as public_port,
  json_extract(rule.value, '$.PrivateIpAddress') as private_ip,
  json_extract(rule.value, '$.PrivatePort') as private_port,
  json_extract(rule.value, '$.Description') as description
from
  tencentcloud_vpc_nat_gateway,
  json_each(destination_ip_port_translation_nat_rule_set) as rule;
```

### Find dedicated (exclusive) NAT gateways

Identify NAT gateways deployed on a dedicated gateway cluster.

```sql+postgres
select
  nat_gateway_id,
  nat_gateway_name,
  vpc_id,
  exclusive_gateway_bandwidth
from
  tencentcloud_vpc_nat_gateway
where
  is_exclusive;
```

```sql+sqlite
select
  nat_gateway_id,
  nat_gateway_name,
  vpc_id,
  exclusive_gateway_bandwidth
from
  tencentcloud_vpc_nat_gateway
where
  is_exclusive = 1;
```

### Count NAT gateways by state

Summarize how many NAT gateways exist in each lifecycle state.

```sql+postgres
select
  state,
  count(*) as count
from
  tencentcloud_vpc_nat_gateway
group by
  state
order by
  count desc;
```

```sql+sqlite
select
  state,
  count(*) as count
from
  tencentcloud_vpc_nat_gateway
group by
  state
order by
  count desc;
```
