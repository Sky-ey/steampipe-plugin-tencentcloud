---
title: "Steampipe Table: tencentcloud_cvm_host - Query Tencent Cloud CDH Hosts using SQL"
description: "Allows users to query Cloud Dedicated Host (CDH) instances in Tencent Cloud to retrieve information about dedicated physical servers."
folder: "CVM"
---

# Table: tencentcloud_cvm_host - Query Tencent Cloud CDH Hosts using SQL

Tencent Cloud Cloud Dedicated Host (CDH) provides dedicated, isolated physical servers that are exclusively used by a single customer. CDH instances give you full control over the underlying hardware resources — CPU cores, memory, disk and GPU cards — so you can deploy CVM instances on a host you fully own, meeting compliance and performance-isolation requirements. This table exposes the host inventory, including capacity and availability of each resource, the CVM instances placed on each host, and the host billing and lifecycle state.

## Table Usage Guide

The `tencentcloud_cvm_host` table in Steampipe provides you with detailed information about CDH hosts within your Tencent Cloud account. This table allows you, as a DevOps engineer or cloud administrator, to query host-specific details, including the host ID, name, state, type, billing mode, resource capacity (CPU/memory/disk/GPU totals and available amounts), zone and project placement, and the CVM instance IDs running on each host. You can utilize this table to gather insights on your dedicated hosts, such as identifying hosts with exhausted CPU or memory capacity, checking hosts approaching expiry, reviewing hosts by availability zone or project, and auditing the instances placed on each host.

## Examples

### Basic info

Get an overview of all CDH hosts in your account, including their names, states, types, and resource capacity. This is useful for inventory management and understanding the current state of your dedicated host fleet.

```sql+postgres
select
  host_id,
  host_name,
  host_state,
  host_type,
  zone,
  region
from
  tencentcloud_cvm_host;
```

```sql+sqlite
select
  host_id,
  host_name,
  host_state,
  host_type,
  zone,
  region
from
  tencentcloud_cvm_host;
```

### List running hosts

Identify all running CDH hosts to understand which dedicated physical servers are currently active. This helps in monitoring operational status and ensuring that expected capacity is available.

```sql+postgres
select
  host_id,
  host_name,
  host_type,
  host_ip
from
  tencentcloud_cvm_host
where
  host_state = 'RUNNING';
```

```sql+sqlite
select
  host_id,
  host_name,
  host_type,
  host_ip
from
  tencentcloud_cvm_host
where
  host_state = 'RUNNING';
```

### Count hosts by state

Summarize the number of hosts in each state to quickly assess the overall health and distribution of your dedicated host resources. This aids in identifying expired or failed hosts that need attention.

```sql+postgres
select
  host_state,
  count(*) as count
from
  tencentcloud_cvm_host
group by
  host_state
order by
  count desc;
```

```sql+sqlite
select
  host_state,
  count(*) as count
from
  tencentcloud_cvm_host
group by
  host_state
order by
  count desc;
```

### Find hosts with low available CPU or memory

Discover hosts whose available CPU cores or memory are running low, which helps you plan capacity — either by migrating instances to another host or by provisioning additional dedicated hosts before resources are exhausted.

```sql+postgres
select
  host_id,
  host_name,
  cpu_total,
  cpu_available,
  mem_total,
  mem_available
from
  tencentcloud_cvm_host
where
  cpu_available <= 4
  or mem_available <= 16;
```

```sql+sqlite
select
  host_id,
  host_name,
  cpu_total,
  cpu_available,
  mem_total,
  mem_available
from
  tencentcloud_cvm_host
where
  cpu_available <= 4
  or mem_available <= 16;
```

### List hosts approaching expiry

Identify prepaid hosts that are nearing their expiry time so you can take action on renewal or migration. This helps avoid unexpected service interruptions due to expired subscriptions.

```sql+postgres
select
  host_id,
  host_name,
  host_charge_type,
  expired_time,
  renew_flag
from
  tencentcloud_cvm_host
where
  expired_time is not null
order by
  expired_time asc;
```

```sql+sqlite
select
  host_id,
  host_name,
  host_charge_type,
  expired_time,
  renew_flag
from
  tencentcloud_cvm_host
where
  expired_time is not null
order by
  expired_time asc;
```

### List hosts in a specific zone

Filter hosts by their availability zone to understand the distribution of dedicated hardware across zones and verify that workloads are deployed in the intended locations.

```sql+postgres
select
  host_id,
  host_name,
  host_state,
  cpu_available,
  mem_available
from
  tencentcloud_cvm_host
where
  zone = 'ap-guangzhou-6';
```

```sql+sqlite
select
  host_id,
  host_name,
  host_state,
  cpu_available,
  mem_available
from
  tencentcloud_cvm_host
where
  zone = 'ap-guangzhou-6';
```

### List GPU hosts

Find dedicated hosts that provide GPU cards, along with how many GPU cards are still available for placement. This is useful for GPU-heavy workloads such as training and inference.

```sql+postgres
select
  host_id,
  host_name,
  gpu_total,
  gpu_available,
  instance_ids
from
  tencentcloud_cvm_host
where
  gpu_total > 0;
```

```sql+sqlite
select
  host_id,
  host_name,
  gpu_total,
  gpu_available,
  instance_ids
from
  tencentcloud_cvm_host
where
  gpu_total > 0;
```

### Count instances placed on each host

List the number of CVM instances running on each host to understand how full the hosts are and detect potential over-placement or under-utilization.

```sql+postgres
select
  host_id,
  host_name,
  json_array_length(instance_ids) as instance_count
from
  tencentcloud_cvm_host;
```

```sql+sqlite
select
  host_id,
  host_name,
  json_array_length(instance_ids) as instance_count
from
  tencentcloud_cvm_host;
```
