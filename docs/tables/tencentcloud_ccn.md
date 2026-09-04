---
title: "Steampipe Table: tencentcloud_ccn - Query Tencent Cloud CCN instances using SQL"
description: "Allows users to query Cloud Connect Network (CCN) instances in Tencent Cloud to inspect cross-region and cross-account VPC interconnectivity."
folder: "CCN"
---

# Table: tencentcloud_ccn - Query Tencent Cloud CCN instances using SQL

Tencent Cloud Cloud Connect Network (CCN) provides cross-region, cross-account VPC interconnectivity. A CCN instance acts as a hub that interconnects VPCs, Direct Connect gateways, and BM VPCs, with controls for QoS, bandwidth limiting, and route priority. CCN is the core resource for building enterprise-grade private networks on Tencent Cloud.

## Table Usage Guide

The `tencentcloud_ccn` table provides you with detailed information about CCN instances within your Tencent Cloud account. Use this table as a network administrator or DevOps engineer to inventory CCNs, audit QoS levels, inspect bandwidth limit types, verify route priority and route broadcast policy flags, and review tag-based classifications. CCN is an account-scoped global service, so a single query returns every CCN in the account regardless of region.

## Examples

### Basic info

List every CCN instance and its core attributes.

```sql+postgres
select
  ccn_id,
  ccn_name,
  state,
  qos_level,
  instance_count,
  bandwidth_limit_type
from
  tencentcloud_ccn;
```

```sql+sqlite
select
  ccn_id,
  ccn_name,
  state,
  qos_level,
  instance_count,
  bandwidth_limit_type
from
  tencentcloud_ccn;
```

### Filter CCNs by state

Identify CCN instances that are currently running.

```sql+postgres
select
  ccn_id,
  ccn_name,
  state,
  qos_level
from
  tencentcloud_ccn
where
  state = 'AVAILABLE';
```

```sql+sqlite
select
  ccn_id,
  ccn_name,
  state,
  qos_level
from
  tencentcloud_ccn
where
  state = 'AVAILABLE';
```

### Count CCN instances by state

Summarize how many CCN instances exist in each state.

```sql+postgres
select
  state,
  count(*) as count
from
  tencentcloud_ccn
group by
  state
order by
  count desc;
```

```sql+sqlite
select
  state,
  count(*) as count
from
  tencentcloud_ccn
group by
  state
order by
  count desc;
```

### Inspect CCNs with route priority enabled

Find CCN instances that support the route priority feature.

```sql+postgres
select
  ccn_id,
  ccn_name,
  route_priority_flag,
  route_broadcast_policy_flag,
  route_table_flag
from
  tencentcloud_ccn
where
  route_priority_flag;
```

```sql+sqlite
select
  ccn_id,
  ccn_name,
  route_priority_flag,
  route_broadcast_policy_flag,
  route_table_flag
from
  tencentcloud_ccn
where
  route_priority_flag = 1;
```

### List CCN instances by tag

Filter CCN instances by a specific tag key-value pair.

```sql+postgres
select
  ccn_id,
  ccn_name,
  tags
from
  tencentcloud_ccn
where
  tags ->> 'Environment' = 'Production';
```

```sql+sqlite
select
  ccn_id,
  ccn_name,
  tags
from
  tencentcloud_ccn
where
  json_extract(tags, '$.Environment') = 'Production';
```
