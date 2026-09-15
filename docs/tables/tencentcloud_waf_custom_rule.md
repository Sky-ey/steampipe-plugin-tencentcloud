---
title: "Steampipe Table: tencentcloud_waf_custom_rule - Query Tencent Cloud WAF custom access control rules using SQL"
description: "Allows users to query WAF custom access control rules per protected domain."
folder: "WAF"
---

# Table: tencentcloud_waf_custom_rule - Query Tencent Cloud WAF custom access control rules using SQL

Tencent Cloud WAF custom access control rules let you define fine-grained match conditions on HTTP request fields (IP, URL, User-Agent, Referer, GET/POST parameters, headers, etc.) and apply a configurable action when the conditions are met. Each rule carries a priority, action type, optional scheduled-task configuration, expiration, and a list of match strategies. WAF is an account-scoped global service, so this table does not expose a `region` matrix column.

## Table Usage Guide

The `tencentcloud_waf_custom_rule` table provides you with detailed information about WAF custom rules. Use this table as a security or DevOps engineer to inventory custom rules per domain, audit rule status, identify expired or soon-to-expire rules, inspect match strategies, and review scheduled-task rules. Because the underlying `DescribeCustomRuleList` API requires a `Domain`, this table fans out over every protected domain (plus the `global` namespace) when no `domain` qual is supplied. The `rule_id` and `rule_name` quals are pushed down as `RuleID` and `RuleName` filters.

## Examples

### Basic info

List custom rules across all domains with their core attributes.

```sql+postgres
select
  domain,
  rule_id,
  rule_name,
  action_type,
  status,
  valid_status,
  sort_id
from
  tencentcloud_waf_custom_rule;
```

```sql+sqlite
select
  domain,
  rule_id,
  rule_name,
  action_type,
  status,
  valid_status,
  sort_id
from
  tencentcloud_waf_custom_rule;
```

### Find rules for a specific domain

List all custom rules attached to a particular domain.

```sql+postgres
select
  rule_id,
  rule_name,
  action_type,
  status,
  create_time,
  modify_time
from
  tencentcloud_waf_custom_rule
where
  domain = 'www.example.com';
```

```sql+sqlite
select
  rule_id,
  rule_name,
  action_type,
  status,
  create_time,
  modify_time
from
  tencentcloud_waf_custom_rule
where
  domain = 'www.example.com';
```

### Look up a single rule by ID

Direct retrieval of one custom rule via its `rule_id`.

```sql+postgres
select
  rule_id,
  rule_name,
  domain,
  action_type,
  status,
  strategies,
  bypass,
  redirect
from
  tencentcloud_waf_custom_rule
where
  rule_id = '1101205922';
```

```sql+sqlite
select
  rule_id,
  rule_name,
  domain,
  action_type,
  status,
  strategies,
  bypass,
  redirect
from
  tencentcloud_waf_custom_rule
where
  rule_id = '1101205922';
```

### Find expired or soon-to-expire rules

Identify rules whose expiration time is in the past, or within the next 7 days.

```sql+postgres
select
  domain,
  rule_id,
  rule_name,
  expire_time,
  status
from
  tencentcloud_waf_custom_rule
where
  expire_time is not null
  and expire_time < now() + interval '7 days';
```

```sql+sqlite
select
  domain,
  rule_id,
  rule_name,
  expire_time,
  status
from
  tencentcloud_waf_custom_rule
where
  expire_time is not null
  and expire_time < datetime('now', '+7 days');
```

### Count rules by domain

Summarize how many custom rules each domain has.

```sql+postgres
select
  domain,
  count(*) as count
from
  tencentcloud_waf_custom_rule
group by
  domain
order by
  count desc;
```

```sql+sqlite
select
  domain,
  count(*) as count
from
  tencentcloud_waf_custom_rule
group by
  domain
order by
  count desc;
```

### Review scheduled-task rules

Inspect rules that use a scheduled task (timed or periodic) rather than permanent enforcement.

```sql+postgres
select
  domain,
  rule_id,
  rule_name,
  job_type,
  cron_type,
  job_date_time
from
  tencentcloud_waf_custom_rule
where
  job_type <> '';
```

```sql+sqlite
select
  domain,
  rule_id,
  rule_name,
  job_type,
  cron_type,
  job_date_time
from
  tencentcloud_waf_custom_rule
where
  job_type <> '';
```
