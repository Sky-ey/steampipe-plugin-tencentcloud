---
title: "Steampipe Table: tencentcloud_waf_object - Query Tencent Cloud WAF protection objects using SQL"
description: "Allows users to query WAF protection objects (CLB listeners or cloud-native gateways bound to WAF instances)."
folder: "WAF"
---

# Table: tencentcloud_waf_object - Query Tencent Cloud WAF protection objects using SQL

Tencent Cloud WAF protection objects represent the CLB listeners or TSE cloud-native gateways that have WAF protection enabled. Each object records its bound WAF instance, protection/log switches (WAF, CLS, Bot, API), access mode, proxy configuration, VPC, and the region the underlying resource lives in. WAF is an account-scoped global service, so this table does not expose a `region` matrix column; `region` is returned as a regular business column.

## Table Usage Guide

The `tencentcloud_waf_object` table provides you with detailed information about WAF protection objects in your Tencent Cloud account. Use this table as a security or DevOps engineer to inventory protected listeners, audit protection switch states, identify objects bound to a specific WAF instance, inspect access modes, and review VPC placement. The `object_id` qual supports direct single-resource lookup via `Get`.

## Examples

### Basic info

List every WAF protection object and its core attributes.

```sql+postgres
select
  object_id,
  object_name,
  instance_id,
  instance_name,
  type,
  status,
  cls_status,
  region
from
  tencentcloud_waf_object;
```

```sql+sqlite
select
  object_id,
  object_name,
  instance_id,
  instance_name,
  type,
  status,
  cls_status,
  region
from
  tencentcloud_waf_object;
```

### Filter objects bound to a specific WAF instance

List protection objects attached to a given WAF instance ID.

```sql+postgres
select
  object_id,
  object_name,
  status,
  cls_status,
  region
from
  tencentcloud_waf_object
where
  instance_id = 'waf-3p4x5yq7';
```

```sql+sqlite
select
  object_id,
  object_name,
  status,
  cls_status,
  region
from
  tencentcloud_waf_object
where
  instance_id = 'waf-3p4x5yq7';
```

### Find objects with WAF protection disabled

Identify protection objects whose WAF switch is turned off.

```sql+postgres
select
  object_id,
  object_name,
  instance_id,
  region
from
  tencentcloud_waf_object
where
  status = 0;
```

```sql+sqlite
select
  object_id,
  object_name,
  instance_id,
  region
from
  tencentcloud_waf_object
where
  status = 0;
```

### Count objects by region

Summarize how many protection objects exist in each region.

```sql+postgres
select
  region,
  count(*) as count
from
  tencentcloud_waf_object
group by
  region
order by
  count desc;
```

```sql+sqlite
select
  region,
  count(*) as count
from
  tencentcloud_waf_object
group by
  region
order by
  count desc;
```

### Look up a single object by ID

Direct retrieval of one protection object via its `object_id`.

```sql+postgres
select
  object_id,
  object_name,
  instance_id,
  type,
  precise_domains,
  public_ip,
  private_ip,
  vpc_id,
  vpc_name
from
  tencentcloud_waf_object
where
  object_id = 'lb-xxxxxxxx';
```

```sql+sqlite
select
  object_id,
  object_name,
  instance_id,
  type,
  precise_domains,
  public_ip,
  private_ip,
  vpc_id,
  vpc_name
from
  tencentcloud_waf_object
where
  object_id = 'lb-xxxxxxxx';
```
