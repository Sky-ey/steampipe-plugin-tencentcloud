---
title: "Steampipe Table: tencentcloud_cos_bucket - Query Tencent Cloud Object Storage buckets using SQL"
description: "Allows users to query all COS buckets under the current account."
folder: "COS"
---

# Table: tencentcloud_cos_bucket - Query Tencent Cloud Object Storage buckets using SQL

Tencent Cloud Object Storage (COS) is a distributed storage service that stores files (objects) in buckets. Each bucket has a globally unique name (in the format `<name>-<appid>`), a home region, a creation date, and a type. Bucket listing is a global Service API — a single call returns every bucket in the account across all regions.

## Table Usage Guide

The `tencentcloud_cos_bucket` table provides information about every COS bucket under the current account. This is a global table — it does not fan out by region. Use this table as a storage administrator to inventory buckets, find buckets by region or creation date, audit for stale or old buckets, and use `tag_key`/`tag_value` to filter by tag. To list the objects inside a bucket, query `tencentcloud_cos_object` and join on `bucket_name` and `region`.

## Examples

### Basic info

List every COS bucket under the current account.

```sql+postgres
select
  name,
  region,
  creation_date,
  bucket_type
from
  tencentcloud_cos_bucket;
```

```sql+sqlite
select
  name,
  region,
  creation_date,
  bucket_type
from
  tencentcloud_cos_bucket;
```

### Find buckets in a specific region

```sql+postgres
select
  name,
  creation_date,
  bucket_type
from
  tencentcloud_cos_bucket
where
  region = 'ap-guangzhou';
```

```sql+sqlite
select
  name,
  creation_date,
  bucket_type
from
  tencentcloud_cos_bucket
where
  region = 'ap-guangzhou';
```

### Count buckets by region

```sql+postgres
select
  region,
  count(*) as total
from
  tencentcloud_cos_bucket
group by
  region
order by
  total desc;
```

```sql+sqlite
select
  region,
  count(*) as total
from
  tencentcloud_cos_bucket
group by
  region
order by
  total desc;
```

### Find old buckets created over a year ago

```sql+postgres
select
  name,
  region,
  creation_date
from
  tencentcloud_cos_bucket
where
  creation_date < now() - interval '1 year'
order by
  creation_date;
```

```sql+sqlite
select
  name,
  region,
  creation_date
from
  tencentcloud_cos_bucket
where
  creation_date < datetime('now', '-1 year')
order by
  creation_date;
```

### Filter buckets by tag

```sql+postgres
select
  name,
  region,
  creation_date
from
  tencentcloud_cos_bucket
where
  tag_key = 'env'
  and tag_value = 'prod';
```

```sql+sqlite
select
  name,
  region,
  creation_date
from
  tencentcloud_cos_bucket
where
  tag_key = 'env'
  and tag_value = 'prod';
```
