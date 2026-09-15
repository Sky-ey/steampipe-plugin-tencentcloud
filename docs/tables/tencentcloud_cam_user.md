---
title: "Steampipe Table: tencentcloud_cam_user - Query Tencent Cloud CAM sub-accounts using SQL"
description: "Allows users to query Tencent Cloud CAM sub-accounts, the IAM users of an account."
folder: "CAM"
---

# Table: tencentcloud_cam_user - Query Tencent Cloud CAM sub-accounts using SQL

Tencent Cloud CAM (Cloud Access Management) sub-accounts are the IAM users of an account. Each sub-account has a UIN and a unique username, and can be granted console login and API access with policies. CAM is an account-scoped global service: there is no region dimension, so the table does not expose a `region` column.

## Table Usage Guide

The `tencentcloud_cam_user` table provides you with detailed information about the sub-accounts of your Tencent Cloud account. Use this table as a security or DevOps engineer to inventory users, audit which users can log in to the console, find users bound to specific emails or phone numbers, and review account creation dates.

## Examples

### Basic info

List every sub-account and its core attributes.

```sql+postgres
select
  uin,
  name,
  uid,
  remark,
  console_login,
  create_time
from
  tencentcloud_cam_user;
```

```sql+sqlite
select
  uin,
  name,
  uid,
  remark,
  console_login,
  create_time
from
  tencentcloud_cam_user;
```

### List users who can log in to the console

Identify sub-accounts that have console login enabled.

```sql+postgres
select
  name,
  uin,
  phone_num,
  email
from
  tencentcloud_cam_user
where
  console_login = 1;
```

```sql+sqlite
select
  name,
  uin,
  phone_num,
  email
from
  tencentcloud_cam_user
where
  console_login = 1;
```

### Get a specific user by name

Fetch a single sub-account by its username.

```sql+postgres
select
  *
from
  tencentcloud_cam_user
where
  name = 'alice';
```

```sql+sqlite
select
  *
from
  tencentcloud_cam_user
where
  name = 'alice';
```

### Count users by console login capability

Summarize how many sub-accounts can log in to the console.

```sql+postgres
select
  console_login,
  count(*) as user_count
from
  tencentcloud_cam_user
group by
  console_login;
```

```sql+sqlite
select
  console_login,
  count(*) as user_count
from
  tencentcloud_cam_user
group by
  console_login;
```

### List users without a phone number

Find sub-accounts that have no bound mobile number, useful for MFA audits.

```sql+postgres
select
  name,
  uin,
  email
from
  tencentcloud_cam_user
where
  phone_num = '';
```

```sql+sqlite
select
  name,
  uin,
  email
from
  tencentcloud_cam_user
where
  phone_num = '';
```
