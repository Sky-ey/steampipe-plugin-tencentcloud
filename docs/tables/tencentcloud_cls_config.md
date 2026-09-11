---
title: "Steampipe Table: tencentcloud_cls_config - Query Tencent Cloud CLS Collection Configs using SQL"
description: "The tencentcloud_cls_config table queries Cloud Log Service (CLS) collection configurations, the rules that define how logs are collected from paths, parsed, and shipped to a log topic."
folder: "CLS"
---

# Table: tencentcloud_cls_config - Query Tencent Cloud CLS Collection Configs using SQL

A CLS collection configuration (config) defines the rules for collecting logs from a specified path, parsing them according to a log type (JSON, delimiter, minimalist, full regex, multiline, syslog, Windows event), and shipping them to a log topic. Each config lives in a specific region and is identified by a unique `config_id`. The `config_id`, `config_name`, and `topic_id` quals push down to the `DescribeConfigs` Filters for efficient single-resource or filtered listing.

## Table Usage Guide

Use this table to inventory CLS collection configurations across regions, inspect their parsing rules and collection paths, and locate configs that ship to a specific log topic. Common operators include cloud administrators auditing log collection coverage and troubleshooting parsing issues.

## Examples

### Basic info

```sql+postgres
select
  config_id,
  config_name,
  log_type,
  path,
  topic_id,
  create_time,
  region
from tencentcloud_cls_config
order by create_time desc;
```

```sql+sqlite
select
  config_id,
  config_name,
  log_type,
  path,
  topic_id,
  create_time,
  region
from tencentcloud_cls_config
order by create_time desc;
```

### Query a specific config by ID

```sql+postgres
select
  config_id,
  config_name,
  log_type,
  extract_rule,
  advanced_config,
  region
from tencentcloud_cls_config
where config_id = '3581a3be-aa41-423b-995a-54ec84da6264';
```

```sql+sqlite
select
  config_id,
  config_name,
  log_type,
  extract_rule,
  advanced_config,
  region
from tencentcloud_cls_config
where config_id = '3581a3be-aa41-423b-995a-54ec84da6264';
```

### Find configs shipping to a specific log topic

```sql+postgres
select
  config_id,
  config_name,
  path,
  log_type,
  region
from tencentcloud_cls_config
where topic_id = '3b83f9d6-3a4d-47f9-9b7f-285c868b2f9a';
```

```sql+sqlite
select
  config_id,
  config_name,
  path,
  log_type,
  region
from tencentcloud_cls_config
where topic_id = '3b83f9d6-3a4d-47f9-9b7f-285c868b2f9a';
```

### Count configs by log type

```sql+postgres
select
  log_type,
  count(*) as config_count
from tencentcloud_cls_config
group by log_type
order by config_count desc;
```

```sql+sqlite
select
  log_type,
  count(*) as config_count
from tencentcloud_cls_config
group by log_type
order by config_count desc;
```

### Find configs using multiline or regex extraction

```sql+postgres
select
  config_id,
  config_name,
  log_type,
  path,
  region
from tencentcloud_cls_config
where log_type in ('multiline_log', 'multiline_fullregex_log', 'fullregex_log');
```

```sql+sqlite
select
  config_id,
  config_name,
  log_type,
  path,
  region
from tencentcloud_cls_config
where log_type in ('multiline_log', 'multiline_fullregex_log', 'fullregex_log');
```
