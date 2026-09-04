---
title: "Steampipe Table: tencentcloud_cvm_region - Query Tencent Cloud CVM Regions using SQL"
description: "Allows users to query CVM regions in Tencent Cloud to retrieve region IDs, descriptions, and availability states."
folder: "CVM"
---

# Table: tencentcloud_cvm_region - Query Tencent Cloud CVM Regions using SQL

Tencent Cloud CVM regions are physically isolated geographic locations where cloud resources are hosted. Each region has a unique ID (e.g., ap-guangzhou), a description (e.g., 华南地区(广州)), and an availability state indicating whether the region is currently available.

## Table Usage Guide

The `tencentcloud_cvm_region` table in Steampipe provides you with the list of regions available to CVM in your Tencent Cloud account. You can use this table, as a DevOps engineer or cloud administrator, to enumerate valid region IDs before provisioning instances, check region availability states, and build region-aware inventory queries.

## Examples

### Basic info

List all regions available to CVM, along with their descriptions and availability states.

```sql+postgres
select
  region,
  region_name,
  region_state
from
  tencentcloud_cvm_region;
```

```sql+sqlite
select
  region,
  region_name,
  region_state
from
  tencentcloud_cvm_region;
```

### Get a specific region

Retrieve details of a single region by its ID. The `region` qualifier is matched client-side after a single DescribeRegions call.

```sql+postgres
select
  region,
  region_name,
  region_state
from
  tencentcloud_cvm_region
where
  region = 'ap-guangzhou';
```

```sql+sqlite
select
  region,
  region_name,
  region_state
from
  tencentcloud_cvm_region
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
  tencentcloud_cvm_region
where
  region_state = 'AVAILABLE';
```

```sql+sqlite
select
  region,
  region_name
from
  tencentcloud_cvm_region
where
  region_state = 'AVAILABLE';
```

### Count regions by state

Count how many regions are in each availability state, useful for a quick health overview of the account's deployable regions.

```sql+postgres
select
  region_state,
  count(*) as region_count
from
  tencentcloud_cvm_region
group by
  region_state;
```

```sql+sqlite
select
  region_state,
  count(*) as region_count
from
  tencentcloud_cvm_region
group by
  region_state;
```

### List region IDs and akas

Show the region ID alongside its globally unique akas, convenient when joining with other Tencent Cloud resources that reference a region.

```sql+postgres
select
  region,
  title,
  akas
from
  tencentcloud_cvm_region
order by
  region;
```

```sql+sqlite
select
  region,
  title,
  akas
from
  tencentcloud_cvm_region
order by
  region;
```
