---
title: "Steampipe Table: tencentcloud_cam_group - Query Tencent Cloud CAM user groups using SQL"
description: "Allows users to query Tencent Cloud CAM user groups."
folder: "CAM"
---

# Table: tencentcloud_cam_group - Query Tencent Cloud CAM user groups using SQL

Tencent Cloud CAM (Cloud Access Management) user groups bundle sub-accounts so that policies can be attached to many users at once. Each group has an ID, a name, a remark, a creation time, and a member count. CAM is an account-scoped global service: there is no region dimension, so the table does not expose a `region` column.

## Table Usage Guide

The `tencentcloud_cam_group` table provides you with detailed information about the user groups of your Tencent Cloud account. Use this table as a security or DevOps engineer to inventory user groups, audit group member counts, and review group descriptions.

## Examples

### Basic info

List every user group and its core attributes.

```sql+postgres
select
  group_id,
  group_name,
  remark,
  create_time,
  group_num
from
  tencentcloud_cam_group;
```

```sql+sqlite
select
  group_id,
  group_name,
  remark,
  create_time,
  group_num
from
  tencentcloud_cam_group;
```

### Get a specific group by ID

Fetch a single user group by its ID, including the member count.

```sql+postgres
select
  *
from
  tencentcloud_cam_group
where
  group_id = 12345;
```

```sql+sqlite
select
  *
from
  tencentcloud_cam_group
where
  group_id = 12345;
```

### List the largest groups

Find user groups with the most members.

```sql+postgres
select
  group_name,
  group_num
from
  tencentcloud_cam_group
order by
  group_num desc;
```

```sql+sqlite
select
  group_name,
  group_num
from
  tencentcloud_cam_group
order by
  group_num desc;
```

### Count groups

Summarize the total number of user groups.

```sql+postgres
select
  count(*) as group_count
from
  tencentcloud_cam_group;
```

```sql+sqlite
select
  count(*) as group_count
from
  tencentcloud_cam_group;
```

### List groups without a remark

Find user groups that have no description.

```sql+postgres
select
  group_id,
  group_name,
  create_time
from
  tencentcloud_cam_group
where
  remark is null
  or remark = '';
```

```sql+sqlite
select
  group_id,
  group_name,
  create_time
from
  tencentcloud_cam_group
where
  remark is null
  or remark = '';
```
