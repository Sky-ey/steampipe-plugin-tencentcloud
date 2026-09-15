---
title: "Steampipe Table: tencentcloud_cam_policy - Query Tencent Cloud CAM policies using SQL"
description: "Allows users to query Tencent Cloud CAM policies, both preset (QCS) and custom (Local)."
folder: "CAM"
---

# Table: tencentcloud_cam_policy - Query Tencent Cloud CAM policies using SQL

Tencent Cloud CAM (Cloud Access Management) policies define the permissions granted to principals. Policies are either preset (QCS, managed by Tencent Cloud) or custom (Local, created by the account), and each carries a policy document, an associated service type, and optionally tags. The `tags` column is a JSON array of `{Key, Value}` objects (not a map), so tag filtering requires expanding the array. CAM is an account-scoped global service: there is no region dimension, so the table does not expose a `region` column.

## Table Usage Guide

The `tencentcloud_cam_policy` table provides you with detailed information about the policies of your Tencent Cloud account. Use this table as a security or DevOps engineer to inventory policies, distinguish preset from custom policies, inspect policy documents, count attachments, and review tag-based classifications. The list call accepts a `scope` qualifier (`All`, `Local`, `QCS`) that is pushed down to the API.

## Examples

### Basic info

List every policy and its core attributes.

```sql+postgres
select
  policy_id,
  policy_name,
  type,
  scope,
  create_mode,
  service_type,
  attachments,
  add_time
from
  tencentcloud_cam_policy;
```

```sql+sqlite
select
  policy_id,
  policy_name,
  type,
  scope,
  create_mode,
  service_type,
  attachments,
  add_time
from
  tencentcloud_cam_policy;
```

### Get a specific policy by ID

Fetch a single policy by its ID, including the policy document.

```sql+postgres
select
  policy_id,
  policy_name,
  policy_document
from
  tencentcloud_cam_policy
where
  policy_id = 123456;
```

```sql+sqlite
select
  policy_id,
  policy_name,
  policy_document
from
  tencentcloud_cam_policy
where
  policy_id = 123456;
```

### List custom (Local) policies only

Filter to policies created by the account.

```sql+postgres
select
  policy_id,
  policy_name,
  create_mode,
  description
from
  tencentcloud_cam_policy
where
  scope = 'Local';
```

```sql+sqlite
select
  policy_id,
  policy_name,
  create_mode,
  description
from
  tencentcloud_cam_policy
where
  scope = 'Local';
```

### Count policies by scope

Summarize how many preset vs custom policies exist.

```sql+postgres
select
  scope,
  count(*) as policy_count
from
  tencentcloud_cam_policy
group by
  scope;
```

```sql+sqlite
select
  scope,
  count(*) as policy_count
from
  tencentcloud_cam_policy
group by
  scope;
```

### List policies by tag

Filter policies by a specific tag key-value pair. The `tags` column is a JSON array of `{Key, Value}` objects, so expand it before matching.

```sql+postgres
select
  policy_id,
  policy_name,
  tags
from
  tencentcloud_cam_policy,
  jsonb_array_elements(tags) as tag
where
  tag ->> 'Key' = 'Environment'
  and tag ->> 'Value' = 'Production';
```

```sql+sqlite
select
  policy_id,
  policy_name,
  tags
from
  tencentcloud_cam_policy,
  json_each(tags) as tag
where
  json_extract(tag.value, '$.Key') = 'Environment'
  and json_extract(tag.value, '$.Value') = 'Production';
```
