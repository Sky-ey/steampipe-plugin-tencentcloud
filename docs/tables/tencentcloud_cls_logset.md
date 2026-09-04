---
title: "Steampipe Table: tencentcloud_cls_logset - Query Tencent Cloud CLS Logsets using SQL"
description: "The tencentcloud_cls_logset table queries Cloud Log Service (CLS) logsets, the top-level container that groups related log topics within a region."
folder: "CLS"
---

# Table: tencentcloud_cls_logset - Query Tencent Cloud CLS Logsets using SQL

A CLS logset is the top-level organizational unit in Tencent Cloud Log Service (CLS). Each logset lives in a specific region and groups one or more log topics (the actual ingestion and storage units). Logsets carry tags, a topic count, and an optional cloud-product attribution (`assumer_name`) when they are created by other services such as CDN or TKE.

## Table Usage Guide

Use this table to inventory CLS logsets across regions, inspect their tag compliance, and locate logsets created by other cloud products. Common operators include cloud administrators auditing log retention scope and tag governance. The `logset_id` and `logset_name` quals push down to the `DescribeLogsets` Filters for efficient single-resource or filtered listing.

## Examples

### List all logsets in the default region

```sql+postgres
select
  logset_id,
  logset_name,
  topic_count,
  create_time,
  region
from tencentcloud_cls_logset
order by create_time desc;
```

```sql+sqlite
select
  logset_id,
  logset_name,
  topic_count,
  create_time,
  region
from tencentcloud_cls_logset
order by create_time desc;
```

### Query a specific logset by ID

```sql+postgres
select
  logset_id,
  logset_name,
  topic_count,
  metric_topic_count,
  tags
from tencentcloud_cls_logset
where logset_id = 'c5b7d2a0-xxxx-xxxx-xxxx-xxxxxxxxxxxx';
```

```sql+sqlite
select
  logset_id,
  logset_name,
  topic_count,
  metric_topic_count,
  tags
from tencentcloud_cls_logset
where logset_id = 'c5b7d2a0-xxxx-xxxx-xxxx-xxxxxxxxxxxx';
```

### Count logsets by region

```sql+postgres
select
  region,
  count(*) as logset_count,
  sum(topic_count) as total_topics
from tencentcloud_cls_logset
group by region
order by logset_count desc;
```

```sql+sqlite
select
  region,
  count(*) as logset_count,
  sum(topic_count) as total_topics
from tencentcloud_cls_logset
group by region
order by logset_count desc;
```

### Find logsets created by other cloud products

```sql+postgres
select
  logset_id,
  logset_name,
  assumer_name,
  role_name,
  assumer_uin,
  region
from tencentcloud_cls_logset
where assumer_name is not null and assumer_name <> '';
```

```sql+sqlite
select
  logset_id,
  logset_name,
  assumer_name,
  role_name,
  assumer_uin,
  region
from tencentcloud_cls_logset
where assumer_name is not null and assumer_name <> '';
```

### Filter logsets by tag key-value pair

```sql+postgres
select
  logset_id,
  logset_name,
  tags ->> 'env' as env,
  topic_count,
  region
from tencentcloud_cls_logset
where tags ->> 'env' = 'production';
```

```sql+sqlite
select
  logset_id,
  logset_name,
  json_extract(tags, '$.env') as env,
  topic_count,
  region
from tencentcloud_cls_logset
where json_extract(tags, '$.env') = 'production';
```
