---
title: "Steampipe Table: tencentcloud_ssl_certificate - Query Tencent Cloud SSL certificates using SQL"
description: "Allows users to query SSL certificates managed under the Tencent Cloud SSL Certificates Service (document/api/400)."
folder: "SSL"
---

# Table: tencentcloud_ssl_certificate - Query Tencent Cloud SSL certificates using SQL

Tencent Cloud SSL Certificates Service issues and manages SSL/TLS certificates for cloud resources. Each certificate records the primary domain, SANs, certificate type (CA or SVR), validity window, status, bound cloud resources, project, and tags. The SSL service is account-scoped and global: there is no region dimension, so the table does not expose a `region` column.

## Table Usage Guide

The `tencentcloud_ssl_certificate` table provides you with detailed information about SSL certificates in your Tencent Cloud account. Use this table as a security or DevOps engineer to inventory certificates, audit expiration dates, identify certificates that are about to expire (within 30 days), verify deployability, inspect bound cloud resources, and review tag-based classifications.

## Examples

### Basic info

List every SSL certificate and its core attributes.

```sql+postgres
select
  certificate_id,
  alias,
  domain,
  certificate_type,
  status_name,
  cert_end_time,
  project_id
from
  tencentcloud_ssl_certificate;
```

```sql+sqlite
select
  certificate_id,
  alias,
  domain,
  certificate_type,
  status_name,
  cert_end_time,
  project_id
from
  tencentcloud_ssl_certificate;
```

### Find certificates that are about to expire

Identify certificates that will expire within 30 days.

```sql+postgres
select
  certificate_id,
  alias,
  domain,
  cert_end_time,
  is_expiring
from
  tencentcloud_ssl_certificate
where
  is_expiring;
```

```sql+sqlite
select
  certificate_id,
  alias,
  domain,
  cert_end_time,
  is_expiring
from
  tencentcloud_ssl_certificate
where
  is_expiring = 1;
```

### List expired certificates

Find certificates whose validity end time is in the past.

```sql+postgres
select
  certificate_id,
  alias,
  domain,
  cert_end_time,
  status_name
from
  tencentcloud_ssl_certificate
where
  cert_end_time < now();
```

```sql+sqlite
select
  certificate_id,
  alias,
  domain,
  cert_end_time,
  status_name
from
  tencentcloud_ssl_certificate
where
  cert_end_time < datetime('now');
```

### Count certificates by type

Summarize how many certificates exist of each type.

```sql+postgres
select
  certificate_type,
  count(*) as count
from
  tencentcloud_ssl_certificate
group by
  certificate_type
order by
  count desc;
```

```sql+sqlite
select
  certificate_type,
  count(*) as count
from
  tencentcloud_ssl_certificate
group by
  certificate_type
order by
  count desc;
```

### List certificates by tag

Filter SSL certificates by a specific tag key-value pair.

```sql+postgres
select
  certificate_id,
  alias,
  domain,
  tags
from
  tencentcloud_ssl_certificate
where
  tags ->> 'Environment' = 'Production';
```

```sql+sqlite
select
  certificate_id,
  alias,
  domain,
  tags
from
  tencentcloud_ssl_certificate
where
  json_extract(tags, '$.Environment') = 'Production';
```
