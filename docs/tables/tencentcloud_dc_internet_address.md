---
title: "Steampipe Table: tencentcloud_dc_internet_address - Query Tencent Cloud Direct Connect internet addresses using SQL"
description: "Allows users to query Direct Connect internet address (public IP) resources allocated to internet tunnels."
folder: "DC"
---

# Table: tencentcloud_dc_internet_address - Query Tencent Cloud Direct Connect internet addresses using SQL

Tencent Cloud Direct Connect internet addresses are public IP address ranges allocated to internet tunnels. Each address records the subnet, mask length, address type (BGP, China Telecom, China Mobile, or China Unicom), address protocol (IPv4 or IPv6), status, and the application, stop, and release timestamps. These addresses are the public-facing endpoints used by internet-facing Direct Connect tunnels.

## Table Usage Guide

The `tencentcloud_dc_internet_address` table provides you with detailed information about public IP address ranges allocated to Direct Connect internet tunnels. Use this table as a network administrator to inventory internet addresses across regions, audit address type (ISP) distribution, verify IPv4 vs IPv6 allocations, inspect status (in use, disabled, returned), and review retention periods for released addresses.

## Examples

### Basic info

List every Direct Connect internet address and its core attributes.

```sql+postgres
select
  instance_id,
  subnet,
  mask_len,
  addr_type,
  addr_proto,
  status,
  resource_region,
  region
from
  tencentcloud_dc_internet_address;
```

```sql+sqlite
select
  instance_id,
  subnet,
  mask_len,
  addr_type,
  addr_proto,
  status,
  resource_region,
  region
from
  tencentcloud_dc_internet_address;
```

### Filter addresses by status

List internet addresses that are currently in use.

```sql+postgres
select
  instance_id,
  subnet,
  mask_len,
  addr_type,
  apply_time,
  region
from
  tencentcloud_dc_internet_address
where
  status = 0;
```

```sql+sqlite
select
  instance_id,
  subnet,
  mask_len,
  addr_type,
  apply_time,
  region
from
  tencentcloud_dc_internet_address
where
  status = 0;
```

### Count addresses by ISP type

Summarize how many internet addresses are allocated to each ISP type.

```sql+postgres
select
  addr_type,
  count(*) as count
from
  tencentcloud_dc_internet_address
group by
  addr_type
order by
  count desc;
```

```sql+sqlite
select
  addr_type,
  count(*) as count
from
  tencentcloud_dc_internet_address
group by
  addr_type
order by
  count desc;
```

### Inspect IPv6 internet addresses

Find internet addresses that are allocated as IPv6.

```sql+postgres
select
  instance_id,
  subnet,
  mask_len,
  addr_type,
  status,
  region
from
  tencentcloud_dc_internet_address
where
  addr_proto = 1;
```

```sql+sqlite
select
  instance_id,
  subnet,
  mask_len,
  addr_type,
  status,
  region
from
  tencentcloud_dc_internet_address
where
  addr_proto = 1;
```

### List disabled or returned addresses

Find internet addresses that are disabled or have been returned.

```sql+postgres
select
  instance_id,
  subnet,
  addr_type,
  status,
  stop_time,
  release_time,
  reserve_time,
  region
from
  tencentcloud_dc_internet_address
where
  status in (1, 2);
```

```sql+sqlite
select
  instance_id,
  subnet,
  addr_type,
  status,
  stop_time,
  release_time,
  reserve_time,
  region
from
  tencentcloud_dc_internet_address
where
  status in (1, 2);
```
