---
title: "Steampipe Table: tencentcloud_vpc_bandwidth_package - Query Tencent Cloud VPC bandwidth packages using SQL"
description: "Allows users to query VPC bandwidth packages and their attached resources."
folder: "VPC"
---

# Table: tencentcloud_vpc_bandwidth_package - Query Tencent Cloud VPC bandwidth packages using SQL

A bandwidth package aggregates and shares bandwidth across multiple public IP resources (EIPs and load balancers) in the same region, so that all attached resources consume a single shared bandwidth pool. Bandwidth packages come in several network types (BGP, dedicated BGP, AIA BGP, static single-line) and billing modes (enhanced 95th percentile `ENHANCED_95_PERCENTILE`, monthly top 5 `TOP5_POSTPAID_BY_MONTH`, monthly fixed bandwidth `FIXED_PREPAID_BY_MONTH`, bandwidth-based `BANDWIDTH_POSTPAID_BY_DAY`, traffic-based `TRAFFIC_POSTPAID_BY_HOUR`, etc.).

## Table Usage Guide

The `tencentcloud_vpc_bandwidth_package` table provides detailed information about bandwidth packages within your Tencent Cloud account. Use this table to inventory bandwidth packages by type and billing mode, find packages whose bandwidth limit may need adjustment, audit which EIPs and load balancers are attached to each package, and review creation times.

## Examples

### Basic info

List every bandwidth package with its key attributes.

```sql+postgres
select
  bandwidth_package_id,
  bandwidth_package_name,
  network_type,
  charge_type,
  status,
  bandwidth,
  region
from
  tencentcloud_vpc_bandwidth_package;
```

```sql+sqlite
select
  bandwidth_package_id,
  bandwidth_package_name,
  network_type,
  charge_type,
  status,
  bandwidth,
  region
from
  tencentcloud_vpc_bandwidth_package;
```

### List packages by network type

Filter bandwidth packages by their network type, for example BGP or dedicated BGP.

```sql+postgres
select
  bandwidth_package_id,
  bandwidth_package_name,
  network_type,
  charge_type,
  bandwidth
from
  tencentcloud_vpc_bandwidth_package
where
  network_type = 'BGP';
```

```sql+sqlite
select
  bandwidth_package_id,
  bandwidth_package_name,
  network_type,
  charge_type,
  bandwidth
from
  tencentcloud_vpc_bandwidth_package
where
  network_type = 'BGP';
```

### List packages with no bandwidth limit

Find bandwidth packages configured with unlimited bandwidth (`-1`).

```sql+postgres
select
  bandwidth_package_id,
  bandwidth_package_name,
  network_type,
  bandwidth
from
  tencentcloud_vpc_bandwidth_package
where
  bandwidth = -1;
```

```sql+sqlite
select
  bandwidth_package_id,
  bandwidth_package_name,
  network_type,
  bandwidth
from
  tencentcloud_vpc_bandwidth_package
where
  bandwidth = -1;
```

### Get the resources attached to each package

Review the EIPs and load balancers sharing each bandwidth package by expanding the `resource_set` JSON column.

```sql+postgres
select
  bandwidth_package_id,
  bandwidth_package_name,
  r ->> 'ResourceType' as resource_type,
  r ->> 'ResourceId' as resource_id,
  r ->> 'AddressIp' as address_ip
from
  tencentcloud_vpc_bandwidth_package,
  jsonb_array_elements(resource_set) as r;
```

```sql+sqlite
select
  bandwidth_package_id,
  bandwidth_package_name,
  json_extract(r.value, '$.ResourceType') as resource_type,
  json_extract(r.value, '$.ResourceId') as resource_id,
  json_extract(r.value, '$.AddressIp') as address_ip
from
  tencentcloud_vpc_bandwidth_package,
  json_each(resource_set) as r;
```

### Count packages by status

Summarize how many bandwidth packages exist in each lifecycle status.

```sql+postgres
select
  status,
  count(*) as count
from
  tencentcloud_vpc_bandwidth_package
group by
  status
order by
  count desc;
```

```sql+sqlite
select
  status,
  count(*) as count
from
  tencentcloud_vpc_bandwidth_package
group by
  status
order by
  count desc;
```
