---
title: "Steampipe Table: tencentcloud_cam_access_key - Query Tencent Cloud CAM access keys using SQL"
description: "Allows users to query Tencent Cloud CAM access keys of sub-accounts."
folder: "CAM"
---

# Table: tencentcloud_cam_access_key - Query Tencent Cloud CAM access keys using SQL

Tencent Cloud CAM (Cloud Access Management) access keys are the API credentials of sub-accounts. Each key has a SecretId, a status (Active or Inactive), a creation time, and a last-used date. Access keys are a common target of credential-leakage audits. CAM is an account-scoped global service: there is no region dimension, so the table does not expose a `region` column.

## Table Usage Guide

The `tencentcloud_cam_access_key` table provides you with detailed information about the access keys of all sub-accounts in your Tencent Cloud account. Use this table as a security auditor to inventory API keys, find inactive keys, identify keys that have never been used, and correlate keys with their owning users. The table fans out across all sub-accounts.

## Examples

### Basic info

List every access key with its owner and status.

```sql+postgres
select
  access_key_id,
  user_name,
  user_uin,
  status,
  create_time,
  last_used_date
from
  tencentcloud_cam_access_key;
```

```sql+sqlite
select
  access_key_id,
  user_name,
  user_uin,
  status,
  create_time,
  last_used_date
from
  tencentcloud_cam_access_key;
```

### List active keys

Find keys that are currently active.

```sql+postgres
select
  access_key_id,
  user_name,
  user_uin,
  create_time
from
  tencentcloud_cam_access_key
where
  status = 'Active';
```

```sql+sqlite
select
  access_key_id,
  user_name,
  user_uin,
  create_time
from
  tencentcloud_cam_access_key
where
  status = 'Active';
```

### List keys that have never been used

Identify potentially leaked or forgotten keys.

```sql+postgres
select
  access_key_id,
  user_name,
  user_uin,
  create_time
from
  tencentcloud_cam_access_key
where
  last_used_date is null;
```

```sql+sqlite
select
  access_key_id,
  user_name,
  user_uin,
  create_time
from
  tencentcloud_cam_access_key
where
  last_used_date is null;
```

### Count keys per user

Find users with the most API keys.

```sql+postgres
select
  user_name,
  count(*) as key_count
from
  tencentcloud_cam_access_key
group by
  user_name
order by
  key_count desc;
```

```sql+sqlite
select
  user_name,
  count(*) as key_count
from
  tencentcloud_cam_access_key
group by
  user_name
order by
  key_count desc;
```

### Join with users for full user context

Join with `tencentcloud_cam_user` to include user email and console login status.

```sql+postgres
select
  k.access_key_id,
  k.status,
  u.email,
  u.console_login
from
  tencentcloud_cam_access_key as k
  join tencentcloud_cam_user as u on u.uin = k.user_uin;
```

```sql+sqlite
select
  k.access_key_id,
  k.status,
  u.email,
  u.console_login
from
  tencentcloud_cam_access_key as k
  join tencentcloud_cam_user as u on u.uin = k.user_uin;
```
