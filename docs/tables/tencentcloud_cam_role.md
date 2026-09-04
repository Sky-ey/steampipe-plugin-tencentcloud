---
title: "Steampipe Table: tencentcloud_cam_role - Query Tencent Cloud CAM roles using SQL"
description: "Allows users to query Tencent Cloud CAM roles, including cross-account and service-linked roles."
folder: "CAM"
---

# Table: tencentcloud_cam_role - Query Tencent Cloud CAM roles using SQL

Tencent Cloud CAM (Cloud Access Management) roles are permission carriers that can be assumed by users, services, or cross-account principals. Each role has a role ID and name, a trust policy document that declares who can assume it, an optional session duration, and tags. CAM is an account-scoped global service: there is no region dimension, so the table does not expose a `region` column.

## Table Usage Guide

The `tencentcloud_cam_role` table provides you with detailed information about the roles of your Tencent Cloud account. Use this table as a security or DevOps engineer to inventory roles, audit which roles can log in to the console, inspect trust policy documents, identify service-linked roles, and review tag-based classifications.

## Examples

### Basic info

List every role and its core attributes.

```sql+postgres
select
  role_id,
  role_name,
  role_type,
  role_arn,
  description,
  add_time
from
  tencentcloud_cam_role;
```

```sql+sqlite
select
  role_id,
  role_name,
  role_type,
  role_arn,
  description,
  add_time
from
  tencentcloud_cam_role;
```

### Get a specific role by name

Fetch a single role by its name.

```sql+postgres
select
  *
from
  tencentcloud_cam_role
where
  role_name = 'QcloudCVMReadOnlyAccess';
```

```sql+sqlite
select
  *
from
  tencentcloud_cam_role
where
  role_name = 'QcloudCVMReadOnlyAccess';
```

### List service-linked roles

Identify roles created by Tencent Cloud services.

```sql+postgres
select
  role_id,
  role_name,
  description,
  deletion_task_id
from
  tencentcloud_cam_role
where
  role_type = 'service_linked';
```

```sql+sqlite
select
  role_id,
  role_name,
  description,
  deletion_task_id
from
  tencentcloud_cam_role
where
  role_type = 'service_linked';
```

### List roles that can log in to the console

Find roles with console login enabled.

```sql+postgres
select
  role_name,
  role_id,
  session_duration
from
  tencentcloud_cam_role
where
  console_login = 1;
```

```sql+sqlite
select
  role_name,
  role_id,
  session_duration
from
  tencentcloud_cam_role
where
  console_login = 1;
```

### List roles by tag

Filter roles by a specific tag key-value pair.

```sql+postgres
select
  role_name,
  role_id,
  tags
from
  tencentcloud_cam_role
where
  tags ->> 'Environment' = 'Production';
```

```sql+sqlite
select
  role_name,
  role_id,
  tags
from
  tencentcloud_cam_role
where
  json_extract(tags, '$.Environment') = 'Production';
```
