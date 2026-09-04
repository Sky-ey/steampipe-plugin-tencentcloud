---
title: "Steampipe Table: tencentcloud_cos_region - Query Tencent Cloud COS Regions using SQL"
description: "Allows users to query COS regions in Tencent Cloud to retrieve region IDs, descriptions, and availability states."
folder: "COS"
---

# Table: tencentcloud_cos_region - Query Tencent Cloud COS Regions using SQL

Tencent Cloud COS (Cloud Object Storage) regions are physically isolated geographic locations where COS buckets are hosted. Each region has a unique ID (e.g., ap-guangzhou), a description (e.g., South China (Guangzhou)), and an availability state indicating whether the region is currently available. Although COS uses an independent RESTful/S3V4 architecture for bucket and object operations, the region metadata is queried via the generic region service.

## Table Usage Guide

The `tencentcloud_cos_region` table in Steampipe provides you with the list of regions available to COS in your Tencent Cloud account. You can use this table, as a DevOps engineer or cloud administrator, to enumerate valid region IDs before creating buckets, check region availability states, and build region-aware inventory queries.

## Examples

### Basic info

List all regions available to COS, along with their descriptions and availability states.

```sql+postgres
select
  region,
  region_name,
  region_state
from
  tencentcloud_cos_region;
```

```sql+sqlite
select
  region,
  region_name,
  region_state
from
  tencentcloud_cos_region;
```

### Get a specific region

Retrieve details of a single region by its ID. The `region` qualifier is matched client-side after a single DescribeRegions call.

```sql+postgres
select
  region,
  region_name,
  region_state
from
  tencentcloud_cos_region
where
  region = 'ap-guangzhou';
```

```sql+sqlite
select
  region,
  region_name,
  region_state
from
  tencentcloud_cos_region
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
  tencentcloud_cos_region
where
  region_state = 'AVAILABLE';
```

```sql+sqlite
select
  region,
  region_name
from
  tencentcloud_cos_region
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
  tencentcloud_cos_region
order by
  region;
```

```sql+sqlite
select
  region,
  title,
  akas
from
  tencentcloud_cos_region
order by
  region;
```
