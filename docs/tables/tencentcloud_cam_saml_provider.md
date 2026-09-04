---
title: "Steampipe Table: tencentcloud_cam_saml_provider - Query Tencent Cloud CAM SAML identity providers using SQL"
description: "Allows users to query Tencent Cloud CAM SAML identity providers used for enterprise SSO."
folder: "CAM"
---

# Table: tencentcloud_cam_saml_provider - Query Tencent Cloud CAM SAML identity providers using SQL

Tencent Cloud CAM (Cloud Access Management) SAML identity providers enable enterprise SSO: an external IdP can be registered so that users of the enterprise can log in to Tencent Cloud with their corporate credentials. Each provider has a name, a description, creation and modification times, and a SAML metadata document. CAM is an account-scoped global service: there is no region dimension, so the table does not expose a `region` column.

## Table Usage Guide

The `tencentcloud_cam_saml_provider` table provides you with detailed information about the SAML identity providers of your Tencent Cloud account. Use this table as a security or DevOps engineer to inventory SSO providers, audit provider metadata, and review when each provider was last modified.

## Examples

### Basic info

List every SAML identity provider and its core attributes.

```sql+postgres
select
  name,
  description,
  create_time,
  modify_time
from
  tencentcloud_cam_saml_provider;
```

```sql+sqlite
select
  name,
  description,
  create_time,
  modify_time
from
  tencentcloud_cam_saml_provider;
```

### Get a specific provider by name

Fetch a single provider, including its SAML metadata document.

```sql+postgres
select
  name,
  description,
  saml_metadata
from
  tencentcloud_cam_saml_provider
where
  name = 'corp-idp';
```

```sql+sqlite
select
  name,
  description,
  saml_metadata
from
  tencentcloud_cam_saml_provider
where
  name = 'corp-idp';
```

### Count providers

Summarize the total number of SAML identity providers.

```sql+postgres
select
  count(*) as provider_count
from
  tencentcloud_cam_saml_provider;
```

```sql+sqlite
select
  count(*) as provider_count
from
  tencentcloud_cam_saml_provider;
```

### List recently modified providers

Find SSO providers that were modified in the last 30 days.

```sql+postgres
select
  name,
  modify_time,
  description
from
  tencentcloud_cam_saml_provider
where
  modify_time >= now() - interval '30 days';
```

```sql+sqlite
select
  name,
  modify_time,
  description
from
  tencentcloud_cam_saml_provider
where
  modify_time >= datetime('now', '-30 days');
```

### List providers without a description

Find SSO providers that have no description.

```sql+postgres
select
  name,
  create_time
from
  tencentcloud_cam_saml_provider
where
  description is null
  or description = '';
```

```sql+sqlite
select
  name,
  create_time
from
  tencentcloud_cam_saml_provider
where
  description is null
  or description = '';
```
