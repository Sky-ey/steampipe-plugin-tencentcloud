---
title: "Steampipe Table: tencentcloud_cam_policy_attachment - Query Tencent Cloud CAM policy attachments using SQL"
description: "Allows users to query Tencent Cloud CAM policy attachments, the mapping between policies and their attached principals."
folder: "CAM"
---

# Table: tencentcloud_cam_policy_attachment - Query Tencent Cloud CAM policy attachments using SQL

Tencent Cloud CAM (Cloud Access Management) policy attachments represent the mapping between a policy and the principals (sub-accounts or user groups) it is attached to. This is the core question of IAM audits: who has been granted what permissions. CAM is an account-scoped global service: there is no region dimension, so the table does not expose a `region` column.

## Table Usage Guide

The `tencentcloud_cam_policy_attachment` table provides you with detailed information about which principals are attached to which policies. Use this table as a security auditor to answer "who has been granted what", list all entities attached to a specific policy, count attachments per policy, and review attachment times.

## Examples

### Basic info

List every policy attachment with its entity.

```sql+postgres
select
  policy_id,
  policy_name,
  entity_id,
  entity_name,
  entity_uin,
  related_type,
  attachment_time
from
  tencentcloud_cam_policy_attachment;
```

```sql+sqlite
select
  policy_id,
  policy_name,
  entity_id,
  entity_name,
  entity_uin,
  related_type,
  attachment_time
from
  tencentcloud_cam_policy_attachment;
```

### List all entities attached to a specific policy

Answer "who has been granted this policy".

```sql+postgres
select
  policy_id,
  policy_name,
  entity_name,
  entity_uin,
  related_type
from
  tencentcloud_cam_policy_attachment
where
  policy_id = 123456;
```

```sql+sqlite
select
  policy_id,
  policy_name,
  entity_name,
  entity_uin,
  related_type
from
  tencentcloud_cam_policy_attachment
where
  policy_id = 123456;
```

### Count attachments per policy

Find the policies with the most attached principals.

```sql+postgres
select
  policy_name,
  count(*) as attachment_count
from
  tencentcloud_cam_policy_attachment
group by
  policy_name
order by
  attachment_count desc;
```

```sql+sqlite
select
  policy_name,
  count(*) as attachment_count
from
  tencentcloud_cam_policy_attachment
group by
  policy_name
order by
  attachment_count desc;
```

### List attachments for user groups only

Filter to entities of type user group (related_type = 2).

```sql+postgres
select
  policy_name,
  entity_name,
  attachment_time
from
  tencentcloud_cam_policy_attachment
where
  related_type = 2;
```

```sql+sqlite
select
  policy_name,
  entity_name,
  attachment_time
from
  tencentcloud_cam_policy_attachment
where
  related_type = 2;
```

### Join with users to find a user's policies

Join with `tencentcloud_cam_user` to resolve the username of each attached entity.

```sql+postgres
select
  a.policy_name,
  u.name as user_name,
  a.attachment_time
from
  tencentcloud_cam_policy_attachment as a
  join tencentcloud_cam_user as u on u.uin = a.entity_uin
where
  a.related_type = 1;
```

```sql+sqlite
select
  a.policy_name,
  u.name as user_name,
  a.attachment_time
from
  tencentcloud_cam_policy_attachment as a
  join tencentcloud_cam_user as u on u.uin = a.entity_uin
where
  a.related_type = 1;
```
