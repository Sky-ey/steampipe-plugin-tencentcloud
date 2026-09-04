---
title: "Steampipe Table: tencentcloud_cvm_image - Query Tencent Cloud CVM Images using SQL"
description: "Allows users to query CVM images in Tencent Cloud to retrieve information about public, private, and shared machine images."
folder: "CVM"
---

# Table: tencentcloud_cvm_image - Query Tencent Cloud CVM Images using SQL

Tencent Cloud CVM images are templates used to launch instances, containing an operating system and optionally pre-installed software. Images come in three types: public images (officially provided by Tencent Cloud), private images (created by your account), and shared images (shared by other accounts). Images capture the full state of a system disk and can be used to quickly replicate environments, scale fleets, and standardize deployments across regions.

## Table Usage Guide

The `tencentcloud_cvm_image` table in Steampipe provides you with detailed information about CVM images within your Tencent Cloud account. This table allows you, as a DevOps engineer or cloud administrator, to query image-specific details, including the image ID, name, type, state, operating system platform, architecture, size, and associated snapshots and tags. You can utilize this table to gather insights on your images, such as auditing which images are available, finding deprecated images still in use, tracking custom image creation, and reviewing tag-based classifications.

## Examples

### Basic info

Get an overview of all images visible to your account, including their names, types, states, and platforms. This is useful for inventory management and understanding which base images are available for launching instances.

```sql+postgres
select
  image_id,
  image_name,
  image_type,
  image_state,
  platform,
  os_name
from
  tencentcloud_cvm_image;
```

```sql+sqlite
select
  image_id,
  image_name,
  image_type,
  image_state,
  platform,
  os_name
from
  tencentcloud_cvm_image;
```

### List private (custom) images

Identify all custom images created by your account. This helps track internally maintained images, manage storage costs, and audit image sprawl.

```sql+postgres
select
  image_id,
  image_name,
  image_state,
  image_size,
  created_time
from
  tencentcloud_cvm_image
where
  image_type = 'PRIVATE_IMAGE';
```

```sql+sqlite
select
  image_id,
  image_name,
  image_state,
  image_size,
  created_time
from
  tencentcloud_cvm_image
where
  image_type = 'PRIVATE_IMAGE';
```

### Find images in a failed state

Discover images whose creation or import has failed. This is critical for operational monitoring, allowing you to investigate build pipeline issues and clean up unusable images.

```sql+postgres
select
  image_id,
  image_name,
  image_state,
  image_source,
  created_time
from
  tencentcloud_cvm_image
where
  image_state in ('CREATEFAILED', 'IMPORTFAILED');
```

```sql+sqlite
select
  image_id,
  image_name,
  image_state,
  image_source,
  created_time
from
  tencentcloud_cvm_image
where
  image_state in ('CREATEFAILED', 'IMPORTFAILED');
```

### List deprecated images

Find images that have been marked as deprecated so you can plan migration away from them before they become unavailable. This is important for keeping instance launches on supported base images.

```sql+postgres
select
  image_id,
  image_name,
  platform,
  os_name,
  image_type
from
  tencentcloud_cvm_image
where
  image_deprecated;
```

```sql+sqlite
select
  image_id,
  image_name,
  platform,
  os_name,
  image_type
from
  tencentcloud_cvm_image
where
  image_deprecated = 1;
```

### Search images by platform

Filter images by platform, such as CentOS or Ubuntu, to find suitable base images for new instances. The platform filter supports fuzzy matching and is pushed down to the DescribeImages API.

```sql+postgres
select
  image_id,
  image_name,
  os_name,
  architecture,
  image_size
from
  tencentcloud_cvm_image
where
  platform = 'Ubuntu';
```

```sql+sqlite
select
  image_id,
  image_name,
  os_name,
  architecture,
  image_size
from
  tencentcloud_cvm_image
where
  platform = 'Ubuntu';
```

### Find images with specific tags

Query images that have been tagged with a specific key-value pair, such as an environment classification. Tag-based filtering is useful for organizing resources and applying policies based on metadata.

```sql+postgres
select
  image_id,
  image_name,
  tags
from
  tencentcloud_cvm_image
where
  tags ->> 'Environment' = 'Production';
```

```sql+sqlite
select
  image_id,
  image_name,
  tags
from
  tencentcloud_cvm_image
where
  json_extract(tags, '$.Environment') = 'Production';
```

### Inspect snapshots backing custom images

Examine the snapshots associated with each image, including the disk usage (system or data disk) and disk size. This aids in understanding image composition and estimating snapshot storage consumption.

```sql+postgres
select
  image_id,
  image_name,
  snapshot ->> 'SnapshotId' as snapshot_id,
  snapshot ->> 'DiskUsage' as disk_usage,
  snapshot ->> 'DiskSize' as disk_size
from
  tencentcloud_cvm_image,
  jsonb_array_elements(snapshot_set) as snapshot;
```

```sql+sqlite
select
  image_id,
  image_name,
  json_extract(snapshot.value, '$.SnapshotId') as snapshot_id,
  json_extract(snapshot.value, '$.DiskUsage') as disk_usage,
  json_extract(snapshot.value, '$.DiskSize') as disk_size
from
  tencentcloud_cvm_image,
  json_each(snapshot_set) as snapshot;
```
