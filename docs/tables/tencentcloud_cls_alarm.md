---
title: "Steampipe Table: tencentcloud_cls_alarm - Query Tencent Cloud CLS Alarm Policies using SQL"
description: "The tencentcloud_cls_alarm table queries Cloud Log Service (CLS) alarm policies, which define trigger conditions, monitoring schedules, and notification channels for log-based alerts."
folder: "CLS"
---

# Table: tencentcloud_cls_alarm - Query Tencent Cloud CLS Alarm Policies using SQL

A CLS alarm policy evaluates log queries on a schedule and fires alerts when trigger conditions are met. Each alarm policy attaches to one or more monitored log topics, references a notification channel group, and carries an alarm level (Warn, Information, Critical). Alarm policies are region-scoped and identified by a unique `alarm_id`. The `alarm_id`, `alarm_name`, and `topic_id` quals push down to the `DescribeAlarms` Filters for efficient single-resource or filtered listing.

## Table Usage Guide

Use this table to inventory CLS alarm policies across regions, inspect their trigger conditions and notification channels, and locate alarms that monitor a specific log topic. Common operators include cloud administrators auditing alerting coverage and troubleshooting alarm configuration.

## Examples

### List all alarm policies in the default region

```sql+postgres
select
  alarm_id,
  alarm_name,
  status,
  alarm_level,
  trigger_count,
  alarm_period,
  create_time,
  region
from tencentcloud_cls_alarm
order by create_time desc;
```

```sql+sqlite
select
  alarm_id,
  alarm_name,
  status,
  alarm_level,
  trigger_count,
  alarm_period,
  create_time,
  region
from tencentcloud_cls_alarm
order by create_time desc;
```

### Query a specific alarm policy by ID

```sql+postgres
select
  alarm_id,
  alarm_name,
  condition,
  alarm_targets,
  alarm_notice_ids,
  tags,
  region
from tencentcloud_cls_alarm
where alarm_id = 'alarm-b60cf034-c3d6-4b01-xxxx-4e877ebb4751';
```

```sql+sqlite
select
  alarm_id,
  alarm_name,
  condition,
  alarm_targets,
  alarm_notice_ids,
  tags,
  region
from tencentcloud_cls_alarm
where alarm_id = 'alarm-b60cf034-c3d6-4b01-xxxx-4e877ebb4751';
```

### Find alarms monitoring a specific log topic

```sql+postgres
select
  alarm_id,
  alarm_name,
  status,
  alarm_level,
  region
from tencentcloud_cls_alarm
where topic_id = '6766f83d-659e-xxxx-a8f7-9104a1012743';
```

```sql+sqlite
select
  alarm_id,
  alarm_name,
  status,
  alarm_level,
  region
from tencentcloud_cls_alarm
where topic_id = '6766f83d-659e-xxxx-a8f7-9104a1012743';
```

### Count alarms by level and status

```sql+postgres
select
  alarm_level,
  status,
  count(*) as alarm_count
from tencentcloud_cls_alarm
group by alarm_level, status
order by alarm_level, status;
```

```sql+sqlite
select
  alarm_level,
  status,
  count(*) as alarm_count
from tencentcloud_cls_alarm
group by alarm_level, status
order by alarm_level, status;
```

### Find disabled alarm policies

```sql+postgres
select
  alarm_id,
  alarm_name,
  alarm_level,
  create_time,
  region
from tencentcloud_cls_alarm
where status = false;
```

```sql+sqlite
select
  alarm_id,
  alarm_name,
  alarm_level,
  create_time,
  region
from tencentcloud_cls_alarm
where status = false;
```
