---
title: "Steampipe Table: tencentcloud_vpc_region - Query Tencent Cloud VPC Regions using SQL"
description: "Allows users to query VPC regions in Tencent Cloud to retrieve region IDs, descriptions, and availability states."
folder: "VPC"
---

# Table: tencentcloud_vpc_region - Query Tencent Cloud VPC Regions using SQL

Tencent Cloud VPC regions are physically isolated geographic locations where cloud resources are hosted. Each region has a unique ID (e.g., ap-guangzhou), a description (e.g., South China (Guangzhou)), and an availability state indicating whether the region is currently available.

## Table Usage Guide

The `tencentcloud_vpc_region` table in Steampipe provides you with the list of regions available to VPC in your Tencent Cloud account. You can use this table, as a DevOps engineer or cloud administrator, to enumerate valid region IDs before provisioning resources, check region availability states, and build region-aware inventory queries.

## Examples

### Basic info

List all regions available to VPC, along with their descriptions and availability states.

```sql+postgres
select
  region,
  region_name,
  region_state
from
  tencentcloud_vpc_region;
```

```sql+sqlite
select
  region,
  region_name,
  region_state
from
  tencentcloud_vpc_region;
```

### Get a specific region

Retrieve details of a single region by its ID. The `region` qualifier is matched client-side after a single DescribeRegions call.

```sql+postgres
select
  region,
  region_name,
  region_state
from
  tencentcloud_vpc_region
where
  region = 'ap-guangzhou';
```

```sql+sqlite
select
  region,
  region_name,
  region_state
from
  tencentcloud_vpc_region
where
  region = 'ap-guangzhou';
```

### List available regions only

Find regions that are currently in the AVAILABLE state, useful for validating deployment targets.

```sql+postgres
select
  region,
  region_name
from
  tencentcloud_vpc_region
where
  region_state = 'AVAILABLE';
```

```sql+sqlite
select
  region,
  region_name
from
  tencentcloud_vpc_region
where
  region_state = 'AVAILABLE';
```

### List region IDs and akas

Show the region ID alongside its globally unique akas, convenient when joining with other Tencent Cloud resources that reference a region.

```sql+postgres
select
  region,
  title,
  akas
from
  tencentcloud_vpc_region
order by
  region;
```

```sql+sqlite
select
  region,
  title,
  akas
from
  tencentcloud_vpc_region
order by
  region;
```
