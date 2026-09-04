---
title: "Steampipe Table: tencentcloud_vpc_traffic_package - Query Tencent Cloud VPC traffic packages using SQL"
description: "Allows users to query VPC traffic packages and their usage."
folder: "VPC"
---

# Table: tencentcloud_vpc_traffic_package - Query Tencent Cloud VPC traffic packages using SQL

A traffic package is a prepaid data transfer allowance (in GB) that is deducted when your account consumes public network traffic. Each package carries a total amount, the used and remaining balance, a status (`AVAILABLE`, `EXPIRED`, `EXHAUSTED`, `REFUNDED`, `DELETED`), a creation time and an expiration deadline, and can be tagged.

## Table Usage Guide

The `tencentcloud_vpc_traffic_package` table provides detailed information about traffic packages within your Tencent Cloud account. Use this table to monitor remaining balance, identify packages that are about to expire or are already exhausted, find packages whose traffic is being deducted, and filter packages by status.

## Examples

### Basic info

List every traffic package with its key attributes.

```sql+postgres
select
  traffic_package_id,
  traffic_package_name,
  status,
  total_amount,
  used_amount,
  remaining_amount,
  deadline,
  region
from
  tencentcloud_vpc_traffic_package;
```

```sql+sqlite
select
  traffic_package_id,
  traffic_package_name,
  status,
  total_amount,
  used_amount,
  remaining_amount,
  deadline,
  region
from
  tencentcloud_vpc_traffic_package;
```

### List available packages

Find traffic packages that can still be used to offset public network traffic.

```sql+postgres
select
  traffic_package_id,
  traffic_package_name,
  total_amount,
  remaining_amount,
  deadline
from
  tencentcloud_vpc_traffic_package
where
  status = 'AVAILABLE';
```

```sql+sqlite
select
  traffic_package_id,
  traffic_package_name,
  total_amount,
  remaining_amount,
  deadline
from
  tencentcloud_vpc_traffic_package
where
  status = 'AVAILABLE';
```

### List exhausted packages

Identify traffic packages whose balance has been fully consumed.

```sql+postgres
select
  traffic_package_id,
  traffic_package_name,
  total_amount,
  used_amount,
  remaining_amount
from
  tencentcloud_vpc_traffic_package
where
  status = 'EXHAUSTED';
```

```sql+sqlite
select
  traffic_package_id,
  traffic_package_name,
  total_amount,
  used_amount,
  remaining_amount
from
  tencentcloud_vpc_traffic_package
where
  status = 'EXHAUSTED';
```

### Find packages expiring soon

List packages whose deadline falls within the next 30 days.

```sql+postgres
select
  traffic_package_id,
  traffic_package_name,
  status,
  remaining_amount,
  deadline
from
  tencentcloud_vpc_traffic_package
where
  status = 'AVAILABLE'
  and deadline < now() + interval '30 days';
```

```sql+sqlite
select
  traffic_package_id,
  traffic_package_name,
  status,
  remaining_amount,
  deadline
from
  tencentcloud_vpc_traffic_package
where
  status = 'AVAILABLE'
  and deadline < datetime('now', '+30 days');
```

### Count packages by status

Summarize how many traffic packages exist in each status.

```sql+postgres
select
  status,
  count(*) as count
from
  tencentcloud_vpc_traffic_package
group by
  status
order by
  count desc;
```

```sql+sqlite
select
  status,
  count(*) as count
from
  tencentcloud_vpc_traffic_package
group by
  status
order by
  count desc;
```

### List packages with a specific tag

Filter traffic packages by a tag key.

```sql+postgres
select
  traffic_package_id,
  traffic_package_name,
  status,
  remaining_amount
from
  tencentcloud_vpc_traffic_package
where
  tags ? 'environment';
```

```sql+sqlite
select
  traffic_package_id,
  traffic_package_name,
  status,
  remaining_amount
from
  tencentcloud_vpc_traffic_package
where
  tags ->> 'environment' is not null;
```
