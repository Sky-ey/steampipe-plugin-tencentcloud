---
title: "Steampipe Table: tencentcloud_vpn_customer_gateway - Query Tencent Cloud VPN customer gateways using SQL"
description: "Allows users to query VPN customer gateways that register the public IP address of the on-premises VPN peer."
folder: "VPN"
---

# Table: tencentcloud_vpn_customer_gateway - Query Tencent Cloud VPN customer gateways using SQL

A Tencent Cloud VPN customer gateway registers the public IP address of the on-premises VPN peer. It is a lightweight resource that pairs with a VPN gateway to form a VPN connection (IPsec tunnel). The customer gateway records only its ID, name, public IP address, and creation time. The DescribeCustomerGateways API lives in the VPC SDK package, but the resource is managed under the VPN product line.

## Table Usage Guide

The `tencentcloud_vpn_customer_gateway` table provides you with detailed information about customer gateway instances. Use this table as a network administrator to inventory customer gateways across regions, audit the registered public IP addresses of on-premises peers, and verify the creation time of each registration. Pass `customer_gateway_id` for a direct lookup, or filter by `customer_gateway_name` or `ip_address`.

## Examples

### Basic info

List every VPN customer gateway and its core attributes.

```sql+postgres
select
  customer_gateway_id,
  customer_gateway_name,
  ip_address,
  created_time,
  region
from
  tencentcloud_vpn_customer_gateway;
```

```sql+sqlite
select
  customer_gateway_id,
  customer_gateway_name,
  ip_address,
  created_time,
  region
from
  tencentcloud_vpn_customer_gateway;
```

### Look up a customer gateway by public IP address

Find the customer gateway that registers a specific on-premises peer IP.

```sql+postgres
select
  customer_gateway_id,
  customer_gateway_name,
  ip_address,
  created_time,
  region
from
  tencentcloud_vpn_customer_gateway
where
  ip_address = '203.0.113.10';
```

```sql+sqlite
select
  customer_gateway_id,
  customer_gateway_name,
  ip_address,
  created_time,
  region
from
  tencentcloud_vpn_customer_gateway
where
  ip_address = '203.0.113.10';
```

### Count customer gateways by region

Summarize how many customer gateways exist in each region.

```sql+postgres
select
  region,
  count(*) as count
from
  tencentcloud_vpn_customer_gateway
group by
  region
order by
  count desc;
```

```sql+sqlite
select
  region,
  count(*) as count
from
  tencentcloud_vpn_customer_gateway
group by
  region
order by
  count desc;
```

### Inspect recently created customer gateways

Find customer gateways created within the last 30 days.

```sql+postgres
select
  customer_gateway_id,
  customer_gateway_name,
  ip_address,
  created_time,
  region
from
  tencentcloud_vpn_customer_gateway
where
  created_time > now() - interval '30 days';
```

```sql+sqlite
select
  customer_gateway_id,
  customer_gateway_name,
  ip_address,
  created_time,
  region
from
  tencentcloud_vpn_customer_gateway
where
  created_time > datetime('now', '-30 days');
```

### Search customer gateways by name pattern

Find customer gateways whose name matches a pattern.

```sql+postgres
select
  customer_gateway_id,
  customer_gateway_name,
  ip_address,
  region
from
  tencentcloud_vpn_customer_gateway
where
  customer_gateway_name ilike '%beijing%';
```

```sql+sqlite
select
  customer_gateway_id,
  customer_gateway_name,
  ip_address,
  region
from
  tencentcloud_vpn_customer_gateway
where
  customer_gateway_name like '%beijing%';
```
