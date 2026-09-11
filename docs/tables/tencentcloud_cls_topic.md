---
title: "Steampipe Table: tencentcloud_cls_topic - Query Tencent Cloud CLS Topics using SQL"
description: "The tencentcloud_cls_topic table queries Cloud Log Service (CLS) log topics, the ingestion and storage units that hold log data, indexing, lifecycle and storage configuration within a logset."
folder: "CLS"
---

# Table: tencentcloud_cls_topic - Query Tencent Cloud CLS Topics using SQL

A CLS topic is the core ingestion and storage unit in Tencent Cloud Log Service. Each topic belongs to a logset and carries its own partition count, storage type (standard/infrequent), lifecycle, indexing flag, auto-split settings, and optional KMS encryption. A topic may be a log topic (`biz_type=0`) or a metric topic (`biz_type=1`).

## Table Usage Guide

Use this table to inventory CLS topics across regions, audit lifecycle and storage configuration, verify indexing status, and trace topics back to their parent logset. The `topic_id`, `topic_name`, and `logset_id` quals push down to the `DescribeTopics` Filters for efficient filtered listing and single-topic lookups.

## Examples

### Basic info

```sql+postgres
select
  topic_id,
  topic_name,
  logset_id,
  storage_type,
  period,
  index_enabled,
  create_time,
  region
from tencentcloud_cls_topic
order by create_time desc;
```

```sql+sqlite
select
  topic_id,
  topic_name,
  logset_id,
  storage_type,
  period,
  index_enabled,
  create_time,
  region
from tencentcloud_cls_topic
order by create_time desc;
```

### Query a specific topic by ID

```sql+postgres
select
  topic_id,
  topic_name,
  logset_id,
  partition_count,
  auto_split,
  max_split_partitions,
  billing_mode,
  tags
from tencentcloud_cls_topic
where topic_id = 'd3c8e1f2-xxxx-xxxx-xxxx-xxxxxxxxxxxx';
```

```sql+sqlite
select
  topic_id,
  topic_name,
  logset_id,
  partition_count,
  auto_split,
  max_split_partitions,
  billing_mode,
  tags
from tencentcloud_cls_topic
where topic_id = 'd3c8e1f2-xxxx-xxxx-xxxx-xxxxxxxxxxxx';
```

### List all topics under a given logset

```sql+postgres
select
  topic_id,
  topic_name,
  storage_type,
  period,
  index_enabled
from tencentcloud_cls_topic
where logset_id = 'c5b7d2a0-xxxx-xxxx-xxxx-xxxxxxxxxxxx';
```

```sql+sqlite
select
  topic_id,
  topic_name,
  storage_type,
  period,
  index_enabled
from tencentcloud_cls_topic
where logset_id = 'c5b7d2a0-xxxx-xxxx-xxxx-xxxxxxxxxxxx';
```

### Find topics with indexing disabled

```sql+postgres
select
  topic_id,
  topic_name,
  logset_id,
  region
from tencentcloud_cls_topic
where index_enabled = false
order by region;
```

```sql+sqlite
select
  topic_id,
  topic_name,
  logset_id,
  region
from tencentcloud_cls_topic
where index_enabled = false
order by region;
```

### Count topics by storage type and region

```sql+postgres
select
  region,
  storage_type,
  count(*) as topic_count,
  sum(partition_count) as total_partitions
from tencentcloud_cls_topic
where biz_type = 0
group by region, storage_type
order by topic_count desc;
```

```sql+sqlite
select
  region,
  storage_type,
  count(*) as topic_count,
  sum(partition_count) as total_partitions
from tencentcloud_cls_topic
where biz_type = 0
group by region, storage_type
order by topic_count desc;
```
