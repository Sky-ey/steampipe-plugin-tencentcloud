---
title: "Steampipe Table: tencentcloud_cvm_zone - Query Tencent Cloud CVM Zones using SQL"
description: "Allows users to query CVM availability zones in Tencent Cloud to retrieve zone IDs, descriptions, and availability states."
folder: "CVM"
---

# Table: tencentcloud_cvm_zone - Query Tencent Cloud CVM Zones using SQL

Tencent Cloud CVM availability zones are physically isolated locations within a region where cloud resources are hosted. Each zone has a unique ID (e.g., ap-guangzhou-3), a description (e.g., 广州三区, Guangzhou Zone 3), a numeric zone ID (e.g., 100003), and an availability state indicating whether the zone is currently available.

## Table Usage Guide

The `tencentcloud_cvm_zone` table in Steampipe provides you with the list of availability zones available to CVM in your Tencent Cloud account. Since DescribeZones is region-scoped, this table discovers CVM regions through DescribeRegions, queries each region, and merges the results. You can use this table, as a DevOps engineer or cloud administrator, to enumerate valid zone IDs before provisioning instances, check zone availability states, and build zone-aware inventory queries.

## Examples

### Basic info

List all availability zones across all CVM regions, along with their descriptions and availability states.

```sql+postgres
select
  zone,
  zone_name,
  zone_id,
  zone_state,
  region
from
  tencentcloud_cvm_zone;
```

```sql+sqlite
select
  zone,
  zone_name,
  zone_id,
  zone_state,
  region
from
  tencentcloud_cvm_zone;
```

### Get a specific zone

Retrieve details of a single availability zone by its ID. The `zone` qualifier is matched client-side after a single DescribeZones call.

```sql+postgres
select
  zone,
  zone_name,
  zone_id,
  zone_state,
  region
from
  tencentcloud_cvm_zone
where
  zone = 'ap-guangzhou-3';
```

```sql+sqlite
select
  zone,
  zone_name,
  zone_id,
  zone_state,
  region
from
  tencentcloud_cvm_zone
where
  zone = 'ap-guangzhou-3';
```

### List available zones in a specific region

Find zones that are currently in the AVAILABLE state in a given region, useful for validating deployment targets.

```sql+postgres
select
  zone,
  zone_name
from
  tencentcloud_cvm_zone
where
  region = 'ap-guangzhou'
  and zone_state = 'AVAILABLE';
```

```sql+sqlite
select
  zone,
  zone_name
from
  tencentcloud_cvm_zone
where
  region = 'ap-guangzhou'
  and zone_state = 'AVAILABLE';
```

### Count zones by state

Count how many availability zones are in each state, useful for a quick overview of deployable capacity across regions.

```sql+postgres
select
  region,
  zone_state,
  count(*) as zone_count
from
  tencentcloud_cvm_zone
group by
  region,
  zone_state;
```

```sql+sqlite
select
  region,
  zone_state,
  count(*) as zone_count
from
  tencentcloud_cvm_zone
group by
  region,
  zone_state;
```

### List zones with numeric zone IDs

Show the zone ID and numeric zone ID together, convenient when provisioning resources that accept the numeric zone identifier.

```sql+postgres
select
  zone,
  zone_name,
  zone_id
from
  tencentcloud_cvm_zone
order by
  region,
  zone;
```

```sql+sqlite
select
  zone,
  zone_name,
  zone_id
from
  tencentcloud_cvm_zone
order by
  region,
  zone;
```
