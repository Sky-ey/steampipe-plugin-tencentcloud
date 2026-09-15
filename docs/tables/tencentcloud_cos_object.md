---
title: "Steampipe Table: tencentcloud_cos_object - Query Tencent Cloud Object Storage objects using SQL"
description: "Allows users to query objects stored in Tencent Cloud COS buckets."
folder: "COS"
---

# Table: tencentcloud_cos_object - Query Tencent Cloud Object Storage objects using SQL

Objects are the files stored inside a COS bucket. Each object has a key (name), size, ETag, storage class, last-modified time, and owner. COS supports multiple storage classes (`STANDARD`, `STANDARD_IA`, `ARCHIVE`, `DEEP_ARCHIVE`, `INTELLIGENT_TIERING`) to balance access frequency and cost.

## Table Usage Guide

The `tencentcloud_cos_object` table provides information about objects across COS buckets. Pass `bucket_name` to list objects in a specific bucket; when it is omitted, the plugin enumerates every bucket in each scanned region and then lists each bucket's objects. Filtering by `region` narrows the scanned regions (it is registered as a key column). Use `prefix` to limit results to a key prefix (the equivalent of a folder path), and `delimiter` to group keys by common prefix. Use this table as a storage administrator to inventory large objects, find archived objects that need restoration, audit storage class distribution, and locate objects by prefix. Join with `tencentcloud_cos_bucket` on `bucket_name` and `region` to correlate object metadata with bucket-level configuration.

> **Caution**: COS buckets can contain millions of objects. Listing all objects across all buckets without a `bucket_name` or `prefix` filter may be slow and consume significant API calls. Always narrow the query with `bucket_name` and/or `prefix` when possible.

## Examples

### Basic info

List objects in a specific bucket.

```sql+postgres
select
  key,
  bucket_name,
  size,
  storage_class,
  last_modified
from
  tencentcloud_cos_object
where
  bucket_name = 'my-bucket-1250000000';
```

```sql+sqlite
select
  key,
  bucket_name,
  size,
  storage_class,
  last_modified
from
  tencentcloud_cos_object
where
  bucket_name = 'my-bucket-1250000000';
```

### List objects under a prefix

```sql+postgres
select
  key,
  size,
  storage_class,
  last_modified
from
  tencentcloud_cos_object
where
  bucket_name = 'my-bucket-1250000000'
  and prefix = 'logs/2024/';
```

```sql+sqlite
select
  key,
  size,
  storage_class,
  last_modified
from
  tencentcloud_cos_object
where
  bucket_name = 'my-bucket-1250000000'
  and prefix = 'logs/2024/';
```

### Find the largest objects

```sql+postgres
select
  key,
  bucket_name,
  size,
  storage_class
from
  tencentcloud_cos_object
where
  bucket_name = 'my-bucket-1250000000'
order by
  size desc
limit 20;
```

```sql+sqlite
select
  key,
  bucket_name,
  size,
  storage_class
from
  tencentcloud_cos_object
where
  bucket_name = 'my-bucket-1250000000'
order by
  size desc
limit 20;
```

### Find archived objects that may need restoration

```sql+postgres
select
  key,
  bucket_name,
  storage_class,
  storage_tier,
  restore_status,
  last_modified
from
  tencentcloud_cos_object
where
  storage_class in ('ARCHIVE', 'DEEP_ARCHIVE');
```

```sql+sqlite
select
  key,
  bucket_name,
  storage_class,
  storage_tier,
  restore_status,
  last_modified
from
  tencentcloud_cos_object
where
  storage_class in ('ARCHIVE', 'DEEP_ARCHIVE');
```

### Count objects and total size by storage class

```sql+postgres
select
  bucket_name,
  storage_class,
  count(*) as object_count,
  sum(size) as total_bytes
from
  tencentcloud_cos_object
where
  bucket_name = 'my-bucket-1250000000'
group by
  bucket_name,
  storage_class
order by
  total_bytes desc;
```

```sql+sqlite
select
  bucket_name,
  storage_class,
  count(*) as object_count,
  sum(size) as total_bytes
from
  tencentcloud_cos_object
where
  bucket_name = 'my-bucket-1250000000'
group by
  bucket_name,
  storage_class
order by
  total_bytes desc;
```
