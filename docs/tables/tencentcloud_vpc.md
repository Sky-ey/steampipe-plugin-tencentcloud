---
title: "Steampipe Table: tencentcloud_vpc - Query Tencent Cloud VPCs using SQL"
description: "Allows users to query VPC instances in Tencent Cloud to inspect private network topology."
folder: "VPC"
---

# Table: tencentcloud_vpc - Query Tencent Cloud VPCs using SQL

Tencent Cloud Virtual Private Cloud (VPC) is a logically isolated network environment in the cloud. You can define the IPv4/IPv6 CIDR block, DHCP options, DNS servers, and auxiliary CIDRs for each VPC. VPCs are the foundation for subnets, route tables, security groups, and other network resources.

## Table Usage Guide

The `tencentcloud_vpc` table provides you with detailed information about VPC instances within your Tencent Cloud account. Use this table as a DevOps engineer or cloud administrator to inventory VPCs across regions, identify default VPCs, audit CIDR allocations, inspect DHCP/DNS configuration, and review tag-based classifications.

## Examples

### Basic info

List every VPC and its core attributes.

```sql+postgres
select
  vpc_id,
  vpc_name,
  cidr_block,
  is_default,
  region
from
  tencentcloud_vpc;
```

```sql+sqlite
select
  vpc_id,
  vpc_name,
  cidr_block,
  is_default,
  region
from
  tencentcloud_vpc;
```

### Find the default VPC per region

Identify the default VPC that Tencent Cloud provisions automatically.

```sql+postgres
select
  vpc_id,
  vpc_name,
  cidr_block,
  region
from
  tencentcloud_vpc
where
  is_default;
```

```sql+sqlite
select
  vpc_id,
  vpc_name,
  cidr_block,
  region
from
  tencentcloud_vpc
where
  is_default = 1;
```

### Inspect VPCs with IPv6 enabled

Find VPCs that have an IPv6 CIDR block allocated.

```sql+postgres
select
  vpc_id,
  vpc_name,
  ipv6_cidr_block,
  ipv6_cidr_block_set
from
  tencentcloud_vpc
where
  ipv6_cidr_block is not null
  and ipv6_cidr_block <> '';
```

```sql+sqlite
select
  vpc_id,
  vpc_name,
  ipv6_cidr_block,
  ipv6_cidr_block_set
from
  tencentcloud_vpc
where
  ipv6_cidr_block is not null
  and ipv6_cidr_block <> '';
```

### Count VPCs by region

Summarize how many VPCs exist in each region.

```sql+postgres
select
  region,
  count(*) as count
from
  tencentcloud_vpc
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
  tencentcloud_vpc
group by
  region
order by
  count desc;
```

### List VPCs with a specific tag

Filter VPCs by a specific tag key-value pair.

```sql+postgres
select
  vpc_id,
  vpc_name,
  tags
from
  tencentcloud_vpc
where
  tags ->> 'Environment' = 'Production';
```

```sql+sqlite
select
  vpc_id,
  vpc_name,
  tags
from
  tencentcloud_vpc
where
  json_extract(tags, '$.Environment') = 'Production';
```
