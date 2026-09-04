---
title: "Steampipe Table: tencentcloud_ccn_instance - Query Tencent Cloud CCN associated instances using SQL"
description: "Allows users to query the network instances (VPCs, Direct Connect gateways, BM VPCs) attached to a Tencent Cloud CCN."
folder: "CCN"
---

# Table: tencentcloud_ccn_instance - Query Tencent Cloud CCN associated instances using SQL

A Tencent Cloud CCN (Cloud Connect Network) associates multiple network instances — VPCs, Direct Connect gateways, and BM VPCs — so that traffic can route between them across regions and accounts. Each association records the instance type, region, owner UIN, CIDR blocks, and the route table that the instance is bound to.

## Table Usage Guide

The `tencentcloud_ccn_instance` table provides you with detailed information about the network instances attached to a CCN. This is a sub-resource table that fans out from a parent CCN: pass `ccn_id` to list associations for a specific CCN, otherwise the plugin enumerates every CCN in the account. Use this table as a network administrator to inventory associations, audit cross-account attachments (via `instance_uin` vs `ccn_uin`), verify state, and inspect CIDR blocks advertised into the CCN.

## Examples

### Basic info

List every CCN-associated network instance and its core attributes.

```sql+postgres
select
  ccn_id,
  instance_id,
  instance_type,
  instance_name,
  instance_region,
  state
from
  tencentcloud_ccn_instance;
```

```sql+sqlite
select
  ccn_id,
  instance_id,
  instance_type,
  instance_name,
  instance_region,
  state
from
  tencentcloud_ccn_instance;
```

### List associated instances for a specific CCN

Filter associations by the parent CCN ID. This is the recommended fan-out pattern for sub-resource tables.

```sql+postgres
select
  instance_id,
  instance_type,
  instance_name,
  instance_region,
  state,
  cidr_block
from
  tencentcloud_ccn_instance
where
  ccn_id = 'ccn-f49l6u0z';
```

```sql+sqlite
select
  instance_id,
  instance_type,
  instance_name,
  instance_region,
  state,
  cidr_block
from
  tencentcloud_ccn_instance
where
  ccn_id = 'ccn-f49l6u0z';
```

### Filter associated instances by type

List only VPC instances attached to CCNs.

```sql+postgres
select
  ccn_id,
  instance_id,
  instance_name,
  instance_region,
  state
from
  tencentcloud_ccn_instance
where
  instance_type = 'VPC';
```

```sql+sqlite
select
  ccn_id,
  instance_id,
  instance_name,
  instance_region,
  state
from
  tencentcloud_ccn_instance
where
  instance_type = 'VPC';
```

### Count associated instances by type per CCN

Summarize how many instances of each type are attached to each CCN.

```sql+postgres
select
  ccn_id,
  instance_type,
  count(*) as count
from
  tencentcloud_ccn_instance
group by
  ccn_id,
  instance_type
order by
  ccn_id,
  count desc;
```

```sql+sqlite
select
  ccn_id,
  instance_type,
  count(*) as count
from
  tencentcloud_ccn_instance
group by
  ccn_id,
  instance_type
order by
  ccn_id,
  count desc;
```

### Inspect cross-account CCN associations

Find associations where the attached instance belongs to a different root account than the CCN owner.

```sql+postgres
select
  ccn_id,
  instance_id,
  instance_type,
  instance_uin,
  ccn_uin,
  instance_region,
  state
from
  tencentcloud_ccn_instance
where
  instance_uin <> ccn_uin;
```

```sql+sqlite
select
  ccn_id,
  instance_id,
  instance_type,
  instance_uin,
  ccn_uin,
  instance_region,
  state
from
  tencentcloud_ccn_instance
where
  instance_uin <> ccn_uin;
```
