---
title: "Steampipe Table: tencentcloud_tke_cluster_instance - Query TKE cluster instances using SQL"
description: "Allows users to query the CVM instances managed by Tencent Kubernetes Engine clusters."
folder: "TKE"
---

# Table: tencentcloud_tke_cluster_instance - Query TKE cluster instances using SQL

A TKE cluster instance is a CVM instance that has been registered to a TKE cluster and runs the Kubernetes kubelet and container runtime. Each instance has a role (`MASTER`, `WORKER`, `ETCD`, `MASTER_ETCD`), a lifecycle state, a private IP, and may belong to a node pool and an auto-scaling group.

## Table Usage Guide

The `tencentcloud_tke_cluster_instance` table provides information about every cluster instance across all clusters in a region. This is a sub-resource table that fans out from a parent cluster: pass `cluster_id` to list instances for a specific cluster, otherwise the plugin enumerates every cluster and lists instances for each. Use this table as a Kubernetes operator to inventory cluster nodes, check node health (`instance_state`, `failed_reason`), verify node roles, find nodes by IP, and identify which node pool or ASG a node belongs to. Join with `tencentcloud_tke_cluster` on `cluster_id` to correlate node health with cluster-level configuration.

## Examples

### Basic info

List every cluster instance across all clusters in the current region.

```sql+postgres
select
  instance_id,
  cluster_id,
  instance_role,
  instance_state,
  lan_ip,
  node_pool_id,
  region
from
  tencentcloud_tke_cluster_instance;
```

```sql+sqlite
select
  instance_id,
  cluster_id,
  instance_role,
  instance_state,
  lan_ip,
  node_pool_id,
  region
from
  tencentcloud_tke_cluster_instance;
```

### List instances for a specific cluster

Filter by the parent cluster ID. This is the recommended fan-out pattern.

```sql+postgres
select
  instance_id,
  instance_role,
  instance_state,
  lan_ip,
  failed_reason
from
  tencentcloud_tke_cluster_instance
where
  cluster_id = 'cls-xxxxxxxx';
```

```sql+sqlite
select
  instance_id,
  instance_role,
  instance_state,
  lan_ip,
  failed_reason
from
  tencentcloud_tke_cluster_instance
where
  cluster_id = 'cls-xxxxxxxx';
```

### Find failed or unhealthy nodes

```sql+postgres
select
  instance_id,
  cluster_id,
  instance_state,
  failed_reason,
  lan_ip
from
  tencentcloud_tke_cluster_instance
where
  instance_state != 'running';
```

```sql+sqlite
select
  instance_id,
  cluster_id,
  instance_state,
  failed_reason,
  lan_ip
from
  tencentcloud_tke_cluster_instance
where
  instance_state != 'running';
```

### Count instances by role in each cluster

```sql+postgres
select
  cluster_id,
  instance_role,
  count(*) as total
from
  tencentcloud_tke_cluster_instance
group by
  cluster_id,
  instance_role
order by
  cluster_id,
  instance_role;
```

```sql+sqlite
select
  cluster_id,
  instance_role,
  count(*) as total
from
  tencentcloud_tke_cluster_instance
group by
  cluster_id,
  instance_role
order by
  cluster_id,
  instance_role;
```

### Find instances by node pool

```sql+postgres
select
  instance_id,
  cluster_id,
  instance_role,
  lan_ip,
  autoscaling_group_id
from
  tencentcloud_tke_cluster_instance
where
  node_pool_id = 'np-xxxxxxxx';
```

```sql+sqlite
select
  instance_id,
  cluster_id,
  instance_role,
  lan_ip,
  autoscaling_group_id
from
  tencentcloud_tke_cluster_instance
where
  node_pool_id = 'np-xxxxxxxx';
```
