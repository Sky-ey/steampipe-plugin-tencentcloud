---
title: "Steampipe Table: tencentcloud_waf_domain - Query Tencent Cloud WAF protected domains using SQL"
description: "Allows users to query SaaS WAF protected domains managed under the Tencent Cloud Web Application Firewall service."
folder: "WAF"
---

# Table: tencentcloud_waf_domain - Query Tencent Cloud WAF protected domains using SQL

Tencent Cloud Web Application Firewall (WAF) protects web domains from common exploits such as SQL injection, XSS, and CC attacks. A protected domain binds to a WAF instance and records its CNAME, edition (sparta-waf for SaaS, clb-waf for CLB), WAF switch, rule engine and AI engine protection modes, access log (CLS) shipping status, Bot and API security switches, origin server IP/domain lists, service ports, security group state, and group tags (labels). WAF is an account-scoped global service: there is no region matrix, so the table does not expose a steampipe `region` column (the API-reported `region` of each domain is exposed as a regular business column).

## Table Usage Guide

The `tencentcloud_waf_domain` table provides you with detailed information about SaaS WAF protected domains in your Tencent Cloud account. Use this table as a security or DevOps engineer to inventory protected domains, audit WAF and CLS log switches, review protection modes (observation vs. interception), inspect origin server configurations, and filter by domain, domain ID, or instance ID. The table supports optional `domain`, `domain_id`, and `instance_id` quals that are pushed down to the DescribeDomains API as exact filters.

## Examples

### Basic info

List every WAF protected domain and its core attributes.

```sql+postgres
select
  domain,
  domain_id,
  instance_id,
  edition,
  status,
  mode,
  engine,
  cls_status
from
  tencentcloud_waf_domain;
```

```sql+sqlite
select
  domain,
  domain_id,
  instance_id,
  edition,
  status,
  mode,
  engine,
  cls_status
from
  tencentcloud_waf_domain;
```

### Query a specific domain by ID

Retrieve a single protected domain by its domain ID (uses the Get hydrate).

```sql+postgres
select
  domain,
  domain_id,
  instance_id,
  cname,
  edition,
  status,
  create_time
from
  tencentcloud_waf_domain
where
  domain_id = 'xxxxxx';
```

```sql+sqlite
select
  domain,
  domain_id,
  instance_id,
  cname,
  edition,
  status,
  create_time
from
  tencentcloud_waf_domain
where
  domain_id = 'xxxxxx';
```

### List domains for a specific WAF instance

Filter protected domains by the bound WAF instance ID.

```sql+postgres
select
  domain,
  domain_id,
  edition,
  status,
  mode,
  engine
from
  tencentcloud_waf_domain
where
  instance_id = 'waf-xxxxxxxx';
```

```sql+sqlite
select
  domain,
  domain_id,
  edition,
  status,
  mode,
  engine
from
  tencentcloud_waf_domain
where
  instance_id = 'waf-xxxxxxxx';
```

### Find domains with WAF protection disabled

Identify protected domains where the WAF switch is off.

```sql+postgres
select
  domain,
  domain_id,
  instance_id,
  status,
  cls_status
from
  tencentcloud_waf_domain
where
  status = 0;
```

```sql+sqlite
select
  domain,
  domain_id,
  instance_id,
  status,
  cls_status
from
  tencentcloud_waf_domain
where
  status = 0;
```

### Count domains by instance and protection mode

Summarize how many protected domains exist per WAF instance and rule engine mode.

```sql+postgres
select
  instance_id,
  mode,
  count(*) as count
from
  tencentcloud_waf_domain
group by
  instance_id,
  mode
order by
  instance_id,
  mode;
```

```sql+sqlite
select
  instance_id,
  mode,
  count(*) as count
from
  tencentcloud_waf_domain
group by
  instance_id,
  mode
order by
  instance_id,
  mode;
```
