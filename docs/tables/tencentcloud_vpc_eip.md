---
title: "Steampipe Table: tencentcloud_vpc_eip - Query Tencent Cloud VPC elastic public IPs using SQL"
description: "Allows users to query elastic public IP (EIP) addresses and their binding status."
folder: "VPC"
---

# Table: tencentcloud_vpc_eip - Query Tencent Cloud VPC elastic public IPs using SQL

An Elastic Public IP (EIP) is a public IP address that you can purchase and bind to a CVM, NAT gateway, ENI, CLB, or other Tencent Cloud resources. An EIP carries a bandwidth value, an internet charge type (prepaid by month, pay-as-you-go by traffic, etc.), and an ISP type (BGP, CMCC, CTCC, CUCC).

## Table Usage Guide

The `tencentcloud_vpc_eip` table provides you with detailed information about EIPs within your Tencent Cloud account. Use this table to inventory EIPs by status, find unbound (idle) EIPs that incur cost without value, audit bound resource types, and review bandwidth and billing configuration.

## Examples

### Basic info

List every EIP with its key attributes.

```sql+postgres
select
  address_id,
  address_name,
  address_ip,
  address_status,
  address_type,
  bandwidth,
  internet_charge_type,
  region
from
  tencentcloud_vpc_eip;
```

```sql+sqlite
select
  address_id,
  address_name,
  address_ip,
  address_status,
  address_type,
  bandwidth,
  internet_charge_type,
  region
from
  tencentcloud_vpc_eip;
```

### Find unbound EIPs

Identify EIPs that are not bound to any resource — these continue to incur charges.

```sql+postgres
select
  address_id,
  address_ip,
  bandwidth,
  internet_charge_type
from
  tencentcloud_vpc_eip
where
  address_status = 'UNBIND';
```

```sql+sqlite
select
  address_id,
  address_ip,
  bandwidth,
  internet_charge_type
from
  tencentcloud_vpc_eip
where
  address_status = 'UNBIND';
```

### Find EIPs bound to a specific instance

Filter EIPs by their bound instance ID.

```sql+postgres
select
  address_id,
  address_ip,
  address_status,
  instance_type,
  network_interface_id,
  private_address_ip
from
  tencentcloud_vpc_eip
where
  instance_id = 'ins-xxxxxxxx';
```

```sql+sqlite
select
  address_id,
  address_ip,
  address_status,
  instance_type,
  network_interface_id,
  private_address_ip
from
  tencentcloud_vpc_eip
where
  instance_id = 'ins-xxxxxxxx';
```

### Find overdue (isolated) EIPs

Identify EIPs that have been isolated due to arrears.

```sql+postgres
select
  address_id,
  address_ip,
  address_status,
  deadline_date
from
  tencentcloud_vpc_eip
where
  is_arrears;
```

```sql+sqlite
select
  address_id,
  address_ip,
  address_status,
  deadline_date
from
  tencentcloud_vpc_eip
where
  is_arrears = 1;
```

### Count EIPs by status

Summarize how many EIPs exist in each status.

```sql+postgres
select
  address_status,
  count(*) as count
from
  tencentcloud_vpc_eip
group by
  address_status
order by
  count desc;
```

```sql+sqlite
select
  address_status,
  count(*) as count
from
  tencentcloud_vpc_eip
group by
  address_status
order by
  count desc;
```
