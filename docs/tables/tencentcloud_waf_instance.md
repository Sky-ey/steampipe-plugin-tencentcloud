---
title: "Steampipe Table: tencentcloud_waf_instance - Query Tencent Cloud WAF instances using SQL"
description: "Allows users to query WAF instances (SaaS and CLB editions) managed under the Tencent Cloud Web Application Firewall service."
folder: "WAF"
---

# Table: tencentcloud_waf_instance - Query Tencent Cloud WAF instances using SQL

Tencent Cloud Web Application Firewall (WAF) keeps web services safe from common web exploits and attacks. A WAF instance is the billable resource that carries the protection capability and comes in two editions: SaaS WAF (domain routed to WAF via CNAME) and CLB WAF (traffic flowing through a load balancer). Each instance records its edition, package level, QPS/bandwidth entitlements, domain quotas, elastic billing switch, expiration window, and status. WAF is an account-scoped global service: there is no region matrix, so the table does not expose a steampipe `region` column (the API-reported `region` of each instance is exposed as a regular business column).

## Table Usage Guide

The `tencentcloud_waf_instance` table provides you with detailed information about WAF instances in your Tencent Cloud account. Use this table as a security or DevOps engineer to inventory WAF instances, audit package levels and expiration dates, review QPS and bandwidth entitlements, inspect elastic billing and auto-renewal flags, and filter by instance ID or name. The table supports optional `instance_id` and `instance_name` quals that are pushed down to the DescribeInstances API as exact filters.

## Examples

### Basic info

List every WAF instance and its core attributes.

```sql+postgres
select
  instance_id,
  instance_name,
  edition,
  level,
  status,
  region,
  valid_time
from
  tencentcloud_waf_instance;
```

```sql+sqlite
select
  instance_id,
  instance_name,
  edition,
  level,
  status,
  region,
  valid_time
from
  tencentcloud_waf_instance;
```

### Query a specific instance by ID

Retrieve a single WAF instance by its instance ID (uses the Get hydrate).

```sql+postgres
select
  instance_id,
  instance_name,
  edition,
  level,
  max_qps,
  domain_count,
  sub_domain_limit
from
  tencentcloud_waf_instance
where
  instance_id = 'waf-xxxxxxxx';
```

```sql+sqlite
select
  instance_id,
  instance_name,
  edition,
  level,
  max_qps,
  domain_count,
  sub_domain_limit
from
  tencentcloud_waf_instance
where
  instance_id = 'waf-xxxxxxxx';
```

### Find instances about to expire

Identify WAF instances whose validity window ends within the next 30 days.

```sql+postgres
select
  instance_id,
  instance_name,
  edition,
  valid_time,
  renew_flag
from
  tencentcloud_waf_instance
where
  valid_time <= now() + interval '30 days'
order by
  valid_time;
```

```sql+sqlite
select
  instance_id,
  instance_name,
  edition,
  valid_time,
  renew_flag
from
  tencentcloud_waf_instance
where
  valid_time <= datetime('now', '+30 days')
order by
  valid_time;
```

### Count instances by edition and level

Summarize how many WAF instances exist per edition and package level.

```sql+postgres
select
  edition,
  level,
  count(*) as count
from
  tencentcloud_waf_instance
group by
  edition,
  level
order by
  edition,
  level;
```

```sql+sqlite
select
  edition,
  level,
  count(*) as count
from
  tencentcloud_waf_instance
group by
  edition,
  level
order by
  edition,
  level;
```

### List instances with elastic billing enabled

Find WAF instances that have elastic QPS billing turned on.

```sql+postgres
select
  instance_id,
  instance_name,
  edition,
  mode,
  elastic_billing,
  qps_standard,
  bandwidth_standard
from
  tencentcloud_waf_instance
where
  mode = 1;
```

```sql+sqlite
select
  instance_id,
  instance_name,
  edition,
  mode,
  elastic_billing,
  qps_standard,
  bandwidth_standard
from
  tencentcloud_waf_instance
where
  mode = 1;
```
