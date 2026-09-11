---
title: "Steampipe Table: tencentcloud_cls_alarm_notice - Query Tencent Cloud CLS Alarm Notification Channel Groups using SQL"
description: "The tencentcloud_cls_alarm_notice table queries Cloud Log Service (CLS) alarm notification channel groups, which define how and where alarm notifications are delivered."
folder: "CLS"
---

# Table: tencentcloud_cls_alarm_notice - Query Tencent Cloud CLS Alarm Notification Channel Groups using SQL

A CLS alarm notification channel group defines the delivery targets for alarm notifications: recipient users/groups, webhooks, and optional log shipping. Each channel group lives in a specific region, is identified by a unique `alarm_notice_id`, and is referenced by one or more alarm policies. The `alarm_notice_id` and `name` quals push down to the `DescribeAlarmNotices` Filters for efficient single-resource or filtered listing.

## Table Usage Guide

Use this table to inventory CLS alarm notification channel groups across regions, inspect their recipient and webhook configuration, and audit delivery status. Common operators include cloud administrators verifying alert routing and troubleshooting notification delivery failures.

## Examples

### Basic info

```sql+postgres
select
  alarm_notice_id,
  name,
  type,
  deliver_status,
  create_time,
  region
from tencentcloud_cls_alarm_notice
order by create_time desc;
```

```sql+sqlite
select
  alarm_notice_id,
  name,
  type,
  deliver_status,
  create_time,
  region
from tencentcloud_cls_alarm_notice
order by create_time desc;
```

### Query a specific notification channel group by ID

```sql+postgres
select
  alarm_notice_id,
  name,
  type,
  notice_receivers,
  web_callbacks,
  tags,
  region
from tencentcloud_cls_alarm_notice
where alarm_notice_id = 'notice-5281f1d2-6275-4e56-9ec3-a1eb19d8bc2f';
```

```sql+sqlite
select
  alarm_notice_id,
  name,
  type,
  notice_receivers,
  web_callbacks,
  tags,
  region
from tencentcloud_cls_alarm_notice
where alarm_notice_id = 'notice-5281f1d2-6275-4e56-9ec3-a1eb19d8bc2f';
```

### Find notification channel groups by name

```sql+postgres
select
  alarm_notice_id,
  name,
  type,
  deliver_status,
  region
from tencentcloud_cls_alarm_notice
where name = 'test-notice';
```

```sql+sqlite
select
  alarm_notice_id,
  name,
  type,
  deliver_status,
  region
from tencentcloud_cls_alarm_notice
where name = 'test-notice';
```

### Count notification channel groups by type and delivery status

```sql+postgres
select
  type,
  deliver_status,
  count(*) as notice_count
from tencentcloud_cls_alarm_notice
group by type, deliver_status
order by type, deliver_status;
```

```sql+sqlite
select
  type,
  deliver_status,
  count(*) as notice_count
from tencentcloud_cls_alarm_notice
group by type, deliver_status
order by type, deliver_status;
```

### Find notification channel groups with delivery exceptions

```sql+postgres
select
  alarm_notice_id,
  name,
  deliver_flag,
  alarm_notice_deliver_config,
  region
from tencentcloud_cls_alarm_notice
where deliver_flag = 3;
```

```sql+sqlite
select
  alarm_notice_id,
  name,
  deliver_flag,
  alarm_notice_deliver_config,
  region
from tencentcloud_cls_alarm_notice
where deliver_flag = 3;
```
