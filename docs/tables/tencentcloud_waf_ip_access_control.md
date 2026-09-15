---
title: "Steampipe Table: tencentcloud_waf_ip_access_control - Query Tencent Cloud WAF IP allowlist/blocklist using SQL"
description: "Allows users to query WAF IP allowlist and blocklist entries across all protected domains."
folder: "WAF"
---

# Table: tencentcloud_waf_ip_access_control - Query Tencent Cloud WAF IP allowlist/blocklist using SQL

Tencent Cloud WAF IP allowlist/blocklist entries control which source IPs are always allowed (40) or always blocked (42) for a specific protected domain. Each entry carries the rule ID, the IP or IP list, remark, source (e.g. `custom`), effective status, optional expiration timestamp, and optional scheduled task configuration. Global rules are stored under the virtual domain name `global`. WAF is an account-scoped global service, so this table does not expose a `region` matrix column.

## Table Usage Guide

The `tencentcloud_waf_ip_access_control` table provides you with detailed information about WAF IP allowlist/blocklist rules. Use this table as a security or DevOps engineer to audit blocklist entries, identify soon-to-expire allowlist rules, review scheduled task configurations, and verify that the right IPs are allowlisted per domain. Because the underlying `DescribeIpAccessControl` API requires a `Domain`, this table fans out over every protected domain (plus the `global` namespace) when no `domain` qual is supplied, and queries both action types (`40` allowlist and `42` blocklist) when `action_type` is not specified.

## Examples

### Basic info

List IP allowlist/blocklist entries across all domains.

```sql+postgres
select
  domain,
  action_type,
  ip,
  note,
  source,
  valid_status,
  create_time
from
  tencentcloud_waf_ip_access_control;
```

```sql+sqlite
select
  domain,
  action_type,
  ip,
  note,
  source,
  valid_status,
  create_time
from
  tencentcloud_waf_ip_access_control;
```

### List blocklist entries for a specific domain

Filter to blocklist (action_type = 42) rules on a particular domain.

```sql+postgres
select
  rule_id,
  ip,
  note,
  valid_status,
  create_time
from
  tencentcloud_waf_ip_access_control
where
  domain = 'www.example.com'
  and action_type = 42;
```

```sql+sqlite
select
  rule_id,
  ip,
  note,
  valid_status,
  create_time
from
  tencentcloud_waf_ip_access_control
where
  domain = 'www.example.com'
  and action_type = 42;
```

### Find soon-to-expire allowlist rules

Identify allowlist rules (action_type = 40) that will expire within the next 7 days.

```sql+postgres
select
  domain,
  rule_id,
  ip,
  note,
  valid_ts
from
  tencentcloud_waf_ip_access_control
where
  action_type = 40
  and valid_ts is not null
  and valid_ts < now() + interval '7 days';
```

```sql+sqlite
select
  domain,
  rule_id,
  ip,
  note,
  valid_ts
from
  tencentcloud_waf_ip_access_control
where
  action_type = 40
  and valid_ts is not null
  and valid_ts < datetime('now', '+7 days');
```

### Count rules by domain and action type

Summarize how many allowlist and blocklist rules each domain has.

```sql+postgres
select
  domain,
  case action_type when 40 then 'allowlist' when 42 then 'blocklist' end as action,
  count(*) as count
from
  tencentcloud_waf_ip_access_control
group by
  domain,
  action_type
order by
  domain,
  action_type;
```

```sql+sqlite
select
  domain,
  case action_type when 40 then 'allowlist' when 42 then 'blocklist' end as action,
  count(*) as count
from
  tencentcloud_waf_ip_access_control
group by
  domain,
  action_type
order by
  domain,
  action_type;
```

### Review scheduled-task rules

Inspect rules that use a scheduled task (timed or periodic) rather than permanent enforcement.

```sql+postgres
select
  domain,
  rule_id,
  ip,
  job_type,
  cron_type,
  job_date_time
from
  tencentcloud_waf_ip_access_control
where
  job_type <> '';
```

```sql+sqlite
select
  domain,
  rule_id,
  ip,
  job_type,
  cron_type,
  job_date_time
from
  tencentcloud_waf_ip_access_control
where
  job_type <> '';
```
