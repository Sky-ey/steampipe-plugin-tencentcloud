---
title: "Steampipe Table: tencentcloud_vpc_eni - Query Tencent Cloud VPC elastic network interfaces using SQL"
description: "Allows users to query elastic network interfaces (ENIs) and their private IPv4/IPv6 addresses."
folder: "VPC"
---

# Table: tencentcloud_vpc_eni - Query Tencent Cloud VPC elastic network interfaces using SQL

An Elastic Network Interface (ENI) is a virtual network card that can be attached to a CVM instance. An ENI carries one or more private IP addresses, a MAC address, a list of bound security groups, and (optionally) an attached CVM instance. Secondary ENIs enable multi-NIC scenarios and IP-forwarding workloads.

## Table Usage Guide

The `tencentcloud_vpc_eni` table provides you with detailed information about ENIs within your Tencent Cloud account. Use this table to inventory ENIs per VPC or subnet, audit ENI state and attachment, find secondary ENIs, and review the security groups and IP addresses attached to each ENI.

## Examples

### Basic info

List every ENI with its key attributes.

```sql+postgres
select
  network_interface_id,
  network_interface_name,
  vpc_id,
  subnet_id,
  instance_id,
  state,
  "primary",
  region
from
  tencentcloud_vpc_eni;
```

```sql+sqlite
select
  network_interface_id,
  network_interface_name,
  vpc_id,
  subnet_id,
  instance_id,
  state,
  "primary",
  region
from
  tencentcloud_vpc_eni;
```

### Find ENIs attached to a specific CVM instance

Filter ENIs by their attached instance ID.

```sql+postgres
select
  network_interface_id,
  network_interface_name,
  "primary",
  mac_address,
  attachment
from
  tencentcloud_vpc_eni
where
  instance_id = 'ins-xxxxxxxx';
```

```sql+sqlite
select
  network_interface_id,
  network_interface_name,
  "primary",
  mac_address,
  attachment
from
  tencentcloud_vpc_eni
where
  instance_id = 'ins-xxxxxxxx';
```

### List unattached ENIs

Identify ENIs that are not bound to any CVM instance.

```sql+postgres
select
  network_interface_id,
  network_interface_name,
  vpc_id,
  subnet_id,
  state
from
  tencentcloud_vpc_eni
where
  instance_id = '';
```

```sql+sqlite
select
  network_interface_id,
  network_interface_name,
  vpc_id,
  subnet_id,
  state
from
  tencentcloud_vpc_eni
where
  instance_id = '';
```

### Find ENIs in a specific subnet

Filter ENIs by their parent subnet ID.

```sql+postgres
select
  network_interface_id,
  network_interface_name,
  instance_id,
  private_ip_address_set
from
  tencentcloud_vpc_eni
where
  subnet_id = 'subnet-xxxxxxxx';
```

```sql+sqlite
select
  network_interface_id,
  network_interface_name,
  instance_id,
  private_ip_address_set
from
  tencentcloud_vpc_eni
where
  subnet_id = 'subnet-xxxxxxxx';
```

### Count ENIs by state

Summarize how many ENIs exist in each state.

```sql+postgres
select
  state,
  count(*) as count
from
  tencentcloud_vpc_eni
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
  tencentcloud_vpc_eni
group by
  state
order by
  count desc;
```
