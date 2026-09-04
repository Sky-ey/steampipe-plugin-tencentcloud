---
title: "Steampipe Table: tencentcloud_vpc_subnet - Query Tencent Cloud VPC subnets using SQL"
description: "Allows users to query VPC subnets to inspect IP ranges, zones, and route table associations."
folder: "VPC"
---

# Table: tencentcloud_vpc_subnet - Query Tencent Cloud VPC subnets using SQL

A subnet is a sub-network within a VPC. Each subnet belongs to a single availability zone and is associated with one route table. Subnets hold the IP ranges that CVM, database, and other Tencent Cloud resources draw private IPs from.

## Table Usage Guide

The `tencentcloud_vpc_subnet` table provides you with detailed information about subnets within your Tencent Cloud account. Use this table to inventory subnets per VPC or zone, audit CIDR exhaustion (`available_ip_address_count`), inspect route table associations, and identify default or CDC subnets.

## Examples

### Basic info

List every subnet with its key attributes.

```sql+postgres
select
  subnet_id,
  subnet_name,
  vpc_id,
  cidr_block,
  zone,
  available_ip_address_count,
  region
from
  tencentcloud_vpc_subnet;
```

```sql+sqlite
select
  subnet_id,
  subnet_name,
  vpc_id,
  cidr_block,
  zone,
  available_ip_address_count,
  region
from
  tencentcloud_vpc_subnet;
```

### List subnets of a specific VPC

Filter subnets by their parent VPC ID.

```sql+postgres
select
  subnet_id,
  subnet_name,
  cidr_block,
  zone,
  route_table_id
from
  tencentcloud_vpc_subnet
where
  vpc_id = 'vpc-xxxxxxxx';
```

```sql+sqlite
select
  subnet_id,
  subnet_name,
  cidr_block,
  zone,
  route_table_id
from
  tencentcloud_vpc_subnet
where
  vpc_id = 'vpc-xxxxxxxx';
```

### Find subnets running low on IP addresses

Identify subnets where fewer than 16 IP addresses remain.

```sql+postgres
select
  subnet_id,
  subnet_name,
  zone,
  available_ip_address_count,
  total_ip_address_count
from
  tencentcloud_vpc_subnet
where
  available_ip_address_count < 16
order by
  available_ip_address_count asc;
```

```sql+sqlite
select
  subnet_id,
  subnet_name,
  zone,
  available_ip_address_count,
  total_ip_address_count
from
  tencentcloud_vpc_subnet
where
  available_ip_address_count < 16
order by
  available_ip_address_count asc;
```

### Count subnets by zone

Summarize the subnet distribution across availability zones.

```sql+postgres
select
  zone,
  count(*) as count
from
  tencentcloud_vpc_subnet
group by
  zone
order by
  count desc;
```

```sql+sqlite
select
  zone,
  count(*) as count
from
  tencentcloud_vpc_subnet
group by
  zone
order by
  count desc;
```

### List subnets with a specific tag

Filter subnets by a specific tag key-value pair.

```sql+postgres
select
  subnet_id,
  subnet_name,
  tags
from
  tencentcloud_vpc_subnet
where
  tags ->> 'Environment' = 'Production';
```

```sql+sqlite
select
  subnet_id,
  subnet_name,
  tags
from
  tencentcloud_vpc_subnet
where
  json_extract(tags, '$.Environment') = 'Production';
```
