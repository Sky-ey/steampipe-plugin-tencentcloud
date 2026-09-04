---
title: "Steampipe Table: tencentcloud_clb_rewrite - Query Tencent Cloud CLB Rewrite Rules using SQL"
description: "Query Tencent Cloud Load Balancer URL rewrite redirection rules, including source and target rules."
folder: "CLB"
---

# Table: tencentcloud_clb_rewrite - Query Tencent Cloud CLB Rewrite Rules using SQL

Tencent Cloud Load Balancer (CLB) rewrite rules redirect requests from a source forwarding rule to a target rule, enabling URL-based redirection. The `tencentcloud_clb_rewrite` table uses the `DescribeRewrite` API and returns the redirection rules of each load balancer in every region, including the source rule, the redirection target, and the rewrite code and type.

## Table Usage Guide

The `tencentcloud_clb_rewrite` table in Steampipe provides information about the URL rewrite redirection rules of CLB load balancers within your Tencent Cloud account. This table allows you, as a DevOps engineer or cloud administrator, to query redirection relationships, including the source listener and rule, the target listener and rule, the redirection status code, and whether the matched URL is carried. You can use this table to audit redirection chains, detect loops or misconfigured rewrites, and review manual versus automatic redirection.

## Examples

### Basic rewrite rule inventory

Get an overview of all rewrite rules across load balancers.

```sql+postgres
select
  load_balancer_id,
  listener_id,
  location_id,
  domain,
  url,
  target_listener_id,
  target_location_id,
  region
from
  tencentcloud_clb_rewrite;
```

```sql+sqlite
select
  load_balancer_id,
  listener_id,
  location_id,
  domain,
  url,
  target_listener_id,
  target_location_id,
  region
from
  tencentcloud_clb_rewrite;
```

### List rewrite rules of a specific load balancer

The `load_balancer_id` qual is pushed down to the API.

```sql+postgres
select
  listener_id,
  location_id,
  domain,
  url,
  target_listener_id,
  target_location_id,
  rewrite_code,
  rewrite_type,
  take_url
from
  tencentcloud_clb_rewrite
where
  load_balancer_id = 'lb-12345678'
  and region = 'ap-guangzhou';
```

```sql+sqlite
select
  listener_id,
  location_id,
  domain,
  url,
  target_listener_id,
  target_location_id,
  rewrite_code,
  rewrite_type,
  take_url
from
  tencentcloud_clb_rewrite
where
  load_balancer_id = 'lb-12345678'
  and region = 'ap-guangzhou';
```

### Count rewrite rules per load balancer

```sql+postgres
select
  load_balancer_id,
  count(*) as rewrite_count
from
  tencentcloud_clb_rewrite
group by
  load_balancer_id
order by
  rewrite_count desc;
```

```sql+sqlite
select
  load_balancer_id,
  count(*) as rewrite_count
from
  tencentcloud_clb_rewrite
group by
  load_balancer_id
order by
  rewrite_count desc;
```

### Find manual rewrites

List manually configured redirection rules (as opposed to automatic ones).

```sql+postgres
select
  load_balancer_id,
  listener_id,
  domain,
  url,
  target_listener_id,
  target_location_id,
  rewrite_code,
  region
from
  tencentcloud_clb_rewrite
where
  rewrite_type = 'Manual';
```

```sql+sqlite
select
  load_balancer_id,
  listener_id,
  domain,
  url,
  target_listener_id,
  target_location_id,
  rewrite_code,
  region
from
  tencentcloud_clb_rewrite
where
  rewrite_type = 'Manual';
```

### Find redirects that drop the URL

Identify rewrite rules that do not carry the matched URL to the target.

```sql+postgres
select
  load_balancer_id,
  listener_id,
  domain,
  url,
  target_listener_id,
  target_location_id,
  take_url,
  region
from
  tencentcloud_clb_rewrite
where
  take_url = false;
```

```sql+sqlite
select
  load_balancer_id,
  listener_id,
  domain,
  url,
  target_listener_id,
  target_location_id,
  take_url,
  region
from
  tencentcloud_clb_rewrite
where
  take_url = false;
```
