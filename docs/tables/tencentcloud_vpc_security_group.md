---
title: "Steampipe Table: tencentcloud_vpc_security_group - Query Tencent Cloud VPC security groups using SQL"
description: "Allows users to query VPC security groups to inventory inbound and outbound rule containers."
folder: "VPC"
---

# Table: tencentcloud_vpc_security_group - Query Tencent Cloud VPC security groups using SQL

A security group is a stateful virtual firewall that controls inbound and outbound traffic for associated resources (CVM, ENI, CLB, etc.). Each security group contains a list of ingress and egress rules; the rules themselves are exposed via the `tencentcloud_vpc_security_group_policy` table.

## Table Usage Guide

The `tencentcloud_vpc_security_group` table provides you with detailed information about security groups within your Tencent Cloud account. Use this table to inventory security groups per project, identify default security groups, audit last update times, and review tag-based classifications.

## Examples

### Basic info

List every security group with its key attributes.

```sql+postgres
select
  security_group_id,
  security_group_name,
  security_group_desc,
  project_id,
  is_default,
  region
from
  tencentcloud_vpc_security_group;
```

```sql+sqlite
select
  security_group_id,
  security_group_name,
  security_group_desc,
  project_id,
  is_default,
  region
from
  tencentcloud_vpc_security_group;
```

### Find default security groups

Identify the default security group in each region.

```sql+postgres
select
  security_group_id,
  security_group_name,
  project_id,
  region
from
  tencentcloud_vpc_security_group
where
  is_default;
```

```sql+sqlite
select
  security_group_id,
  security_group_name,
  project_id,
  region
from
  tencentcloud_vpc_security_group
where
  is_default = 1;
```

### List security groups of a specific project

Filter security groups by their project ID.

```sql+postgres
select
  security_group_id,
  security_group_name,
  created_time,
  update_time
from
  tencentcloud_vpc_security_group
where
  project_id = '0';
```

```sql+sqlite
select
  security_group_id,
  security_group_name,
  created_time,
  update_time
from
  tencentcloud_vpc_security_group
where
  project_id = '0';
```

### Count security groups by region

Summarize how many security groups exist in each region.

```sql+postgres
select
  region,
  count(*) as count
from
  tencentcloud_vpc_security_group
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
  tencentcloud_vpc_security_group
group by
  region
order by
  count desc;
```

### List security groups with a specific tag

Filter security groups by a specific tag key-value pair.

```sql+postgres
select
  security_group_id,
  security_group_name,
  tags
from
  tencentcloud_vpc_security_group
where
  tags ->> 'Environment' = 'Production';
```

```sql+sqlite
select
  security_group_id,
  security_group_name,
  tags
from
  tencentcloud_vpc_security_group
where
  json_extract(tags, '$.Environment') = 'Production';
```
