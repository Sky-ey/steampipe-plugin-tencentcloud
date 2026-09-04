---
title: "Steampipe Table: tencentcloud_vpc_security_group_policy - Query Tencent Cloud VPC security group rules using SQL"
description: "Allows users to query inbound and outbound rules of every VPC security group."
folder: "VPC"
---

# Table: tencentcloud_vpc_security_group_policy - Query Tencent Cloud VPC security group rules using SQL

A security group policy (rule) defines a single match-and-action pair — protocol, port, source/destination CIDR — applied to inbound (`INGRESS`) or outbound (`EGRESS`) traffic. The Tencent Cloud `DescribeSecurityGroupPolicies` API accepts a single security group ID per call and returns both rule sets at once, so this table fans out across every security group when no `security_group_id` qual is supplied and expands each rule into one row.

## Table Usage Guide

The `tencentcloud_vpc_security_group_policy` table provides you with detailed information about individual security group rules. Use this table to audit permissive rules (e.g. `0.0.0.0/0` on SSH/RDP), find rules with a given action (`ACCEPT`/`DROP`), inventory rules per security group, or filter by direction (`INBOUND`/`OUTBOUND`).

## Examples

### Basic info

List every security group rule with its key attributes.

```sql+postgres
select
  security_group_id,
  direction,
  protocol,
  port,
  cidr_block,
  action
from
  tencentcloud_vpc_security_group_policy;
```

```sql+sqlite
select
  security_group_id,
  direction,
  protocol,
  port,
  cidr_block,
  action
from
  tencentcloud_vpc_security_group_policy;
```

### Find permissive inbound rules

Identify inbound rules that allow traffic from any IP address.

```sql+postgres
select
  security_group_id,
  protocol,
  port,
  cidr_block,
  action
from
  tencentcloud_vpc_security_group_policy
where
  direction = 'INBOUND'
  and cidr_block = '0.0.0.0/0';
```

```sql+sqlite
select
  security_group_id,
  protocol,
  port,
  cidr_block,
  action
from
  tencentcloud_vpc_security_group_policy
where
  direction = 'INBOUND'
  and cidr_block = '0.0.0.0/0';
```

### List rules of a specific security group

Filter rules by their parent security group ID.

```sql+postgres
select
  direction,
  protocol,
  port,
  cidr_block,
  ipv6_cidr_block,
  action,
  policy_description
from
  tencentcloud_vpc_security_group_policy
where
  security_group_id = 'sg-xxxxxxxx';
```

```sql+sqlite
select
  direction,
  protocol,
  port,
  cidr_block,
  ipv6_cidr_block,
  action,
  policy_description
from
  tencentcloud_vpc_security_group_policy
where
  security_group_id = 'sg-xxxxxxxx';
```

### Count rules by action

Summarize how many rules accept vs. drop traffic.

```sql+postgres
select
  action,
  count(*) as count
from
  tencentcloud_vpc_security_group_policy
group by
  action
order by
  count desc;
```

```sql+sqlite
select
  action,
  count(*) as count
from
  tencentcloud_vpc_security_group_policy
group by
  action
order by
  count desc;
```

### Find SSH and RDP rules exposed to the internet

Identify rules that expose SSH (port 22) or RDP (port 3389) to the public internet.

```sql+postgres
select
  security_group_id,
  direction,
  protocol,
  port,
  cidr_block,
  action
from
  tencentcloud_vpc_security_group_policy
where
  direction = 'INBOUND'
  and cidr_block = '0.0.0.0/0'
  and port in ('22', '3389');
```

```sql+sqlite
select
  security_group_id,
  direction,
  protocol,
  port,
  cidr_block,
  action
from
  tencentcloud_vpc_security_group_policy
where
  direction = 'INBOUND'
  and cidr_block = '0.0.0.0/0'
  and port in ('22', '3389');
```
