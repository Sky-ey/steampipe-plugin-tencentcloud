---
title: "Steampipe Table: tencentcloud_vpn_gateway - Query Tencent Cloud VPN gateways using SQL"
description: "Allows users to query VPN gateways that establish IPsec/SSL tunnels between a VPC and on-premises networks, or CCN-attached gateways."
folder: "VPN"
---

# Table: tencentcloud_vpn_gateway - Query Tencent Cloud VPN gateways using SQL

Tencent Cloud VPN gateways establish IPsec or SSL tunnels between a VPC and on-premises networks, or act as CCN-attached gateways. Each gateway records the public IP address, maximum outbound bandwidth, gateway type (IPSEC, SSL, or CCN), associated VPC or CCN instance, zone, billing type, and renewal flag. The DescribeVpnGateways API lives in the VPC SDK package, but the resource is managed under the VPN product line.

## Table Usage Guide

The `tencentcloud_vpn_gateway` table provides you with detailed information about VPN gateway instances. Use this table as a network administrator to inventory gateways across regions, audit gateway types, verify VPC/CCN associations, inspect public IP and bandwidth allocations, and review billing and renewal configuration. Pass `vpn_gateway_id` for a direct lookup, or filter by `vpc_id`, `vpn_gateway_name`, `type`, `zone`, `public_ip_address`, or `renew_flag`.

## Examples

### Basic info

List every VPN gateway and its core attributes.

```sql+postgres
select
  vpn_gateway_id,
  vpn_gateway_name,
  type,
  state,
  public_ip_address,
  internet_max_bandwidth_out,
  vpc_id,
  region
from
  tencentcloud_vpn_gateway;
```

```sql+sqlite
select
  vpn_gateway_id,
  vpn_gateway_name,
  type,
  state,
  public_ip_address,
  internet_max_bandwidth_out,
  vpc_id,
  region
from
  tencentcloud_vpn_gateway;
```

### Filter VPN gateways by state

Identify VPN gateways that are currently available.

```sql+postgres
select
  vpn_gateway_id,
  vpn_gateway_name,
  type,
  public_ip_address,
  vpc_id,
  region
from
  tencentcloud_vpn_gateway
where
  state = 'AVAILABLE';
```

```sql+sqlite
select
  vpn_gateway_id,
  vpn_gateway_name,
  type,
  public_ip_address,
  vpc_id,
  region
from
  tencentcloud_vpn_gateway
where
  state = 'AVAILABLE';
```

### Count VPN gateways by type per region

Summarize how many VPN gateways of each type exist per region.

```sql+postgres
select
  region,
  type,
  count(*) as count
from
  tencentcloud_vpn_gateway
group by
  region,
  type
order by
  region,
  count desc;
```

```sql+sqlite
select
  region,
  type,
  count(*) as count
from
  tencentcloud_vpn_gateway
group by
  region,
  type
order by
  region,
  count desc;
```

### Inspect CCN-type VPN gateways

Find VPN gateways that attach to a CCN.

```sql+postgres
select
  vpn_gateway_id,
  vpn_gateway_name,
  network_instance_id,
  public_ip_address,
  zone,
  region
from
  tencentcloud_vpn_gateway
where
  type = 'CCN';
```

```sql+sqlite
select
  vpn_gateway_id,
  vpn_gateway_name,
  network_instance_id,
  public_ip_address,
  zone,
  region
from
  tencentcloud_vpn_gateway
where
  type = 'CCN';
```

### List VPN gateways with blocked public IP addresses

Find VPN gateways whose public IP address is blocked.

```sql+postgres
select
  vpn_gateway_id,
  vpn_gateway_name,
  public_ip_address,
  restrict_state,
  region
from
  tencentcloud_vpn_gateway
where
  is_address_blocked;
```

```sql+sqlite
select
  vpn_gateway_id,
  vpn_gateway_name,
  public_ip_address,
  restrict_state,
  region
from
  tencentcloud_vpn_gateway
where
  is_address_blocked = 1;
```
