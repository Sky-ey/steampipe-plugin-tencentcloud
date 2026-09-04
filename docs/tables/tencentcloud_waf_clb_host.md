---
title: "Steampipe Table: tencentcloud_waf_clb_host - Query Tencent Cloud WAF CLB-protected domains using SQL"
description: "Allows users to query CLB WAF protected domains (hosts) managed under the Tencent Cloud Web Application Firewall service."
folder: "WAF"
---

# Table: tencentcloud_waf_clb_host - Query Tencent Cloud WAF CLB-protected domains using SQL

Tencent Cloud Web Application Firewall (WAF) protects domains whose traffic flows through a Cloud Load Balancer (CLB) — the CLB edition. A CLB-protected host records its domain, domain ID, binding status to the CLB, listener state, rule engine and AI engine protection modes, traffic mode (mirror vs. cleaning), access log (CLS) shipping switch, protection level, bound CLB instance list, and the region(s) of the bound CLB. The DescribeHosts API has no Offset/Limit pagination: a single request returns the full matching host list. WAF is an account-scoped global service: there is no region matrix, so the table does not expose a steampipe `region` column (the API-reported `region` of the bound CLB is exposed as a regular business column).

## Table Usage Guide

The `tencentcloud_waf_clb_host` table provides you with detailed information about CLB WAF protected domains in your Tencent Cloud account. Use this table as a security or DevOps engineer to inventory CLB-protected hosts, audit WAF binding and listener states, review protection modes and traffic modes, inspect CLS log switches, and filter by domain, domain ID, or WAF instance ID. The table supports optional `domain`, `domain_id`, and `instance_id` quals that are pushed down directly to the DescribeHosts request fields (not Filters). Note that `instance_id` is a filter-only column: the API uses it to filter results but does not echo it back on each host, so it is populated from the query qual.

## Examples

### Basic info

List every CLB WAF protected domain and its core attributes.

```sql+postgres
select
  domain,
  domain_id,
  main_domain,
  edition,
  status,
  state,
  mode,
  engine,
  cls_status
from
  tencentcloud_waf_clb_host;
```

```sql+sqlite
select
  domain,
  domain_id,
  main_domain,
  edition,
  status,
  state,
  mode,
  engine,
  cls_status
from
  tencentcloud_waf_clb_host;
```

### Query a specific CLB-protected domain

Retrieve protected hosts for a specific domain name.

```sql+postgres
select
  domain,
  domain_id,
  edition,
  status,
  state,
  flow_mode,
  level,
  region
from
  tencentcloud_waf_clb_host
where
  domain = 'www.example.com';
```

```sql+sqlite
select
  domain,
  domain_id,
  edition,
  status,
  state,
  flow_mode,
  level,
  region
from
  tencentcloud_waf_clb_host
where
  domain = 'www.example.com';
```

### List hosts bound to a specific WAF instance

Filter CLB-protected hosts by the WAF instance ID.

```sql+postgres
select
  domain,
  domain_id,
  status,
  state,
  mode,
  engine
from
  tencentcloud_waf_clb_host
where
  instance_id = 'waf-xxxxxxxx';
```

```sql+sqlite
select
  domain,
  domain_id,
  status,
  state,
  mode,
  engine
from
  tencentcloud_waf_clb_host
where
  instance_id = 'waf-xxxxxxxx';
```

### Find unbound hosts

Identify CLB-protected domains where WAF is not yet bound to the CLB.

```sql+postgres
select
  domain,
  domain_id,
  status,
  state,
  alb_type
from
  tencentcloud_waf_clb_host
where
  status = 0;
```

```sql+sqlite
select
  domain,
  domain_id,
  status,
  state,
  alb_type
from
  tencentcloud_waf_clb_host
where
  status = 0;
```

### Count hosts by protection mode and traffic mode

Summarize how many CLB-protected hosts exist per rule engine mode and traffic mode.

```sql+postgres
select
  mode,
  flow_mode,
  count(*) as count
from
  tencentcloud_waf_clb_host
group by
  mode,
  flow_mode
order by
  mode,
  flow_mode;
```

```sql+sqlite
select
  mode,
  flow_mode,
  count(*) as count
from
  tencentcloud_waf_clb_host
group by
  mode,
  flow_mode
order by
  mode,
  flow_mode;
```
