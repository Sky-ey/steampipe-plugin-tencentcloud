---
title: "Steampipe Table: tencentcloud_tke_cluster - Query Tencent Kubernetes Engine clusters using SQL"
description: "Allows users to query TKE managed and independent Kubernetes clusters."
folder: "TKE"
---

# Table: tencentcloud_tke_cluster - Query Tencent Kubernetes Engine clusters using SQL

Tencent Kubernetes Engine (TKE) is Tencent Cloud's managed Kubernetes service. A cluster is the top-level resource that holds nodes, node pools, networking settings, and container runtime configuration. TKE clusters come in two types: `MANAGED_CLUSTER` (control plane managed by Tencent) and `INDEPENDENT_CLUSTER` (self-managed control plane).

## Table Usage Guide

The `tencentcloud_tke_cluster` table provides detailed information about each TKE cluster in a region, including its version, OS, container runtime, network settings (VPC, CIDR, service CIDR), node counts, tags, and lifecycle status. Use this table as a Kubernetes platform operator to inventory clusters, verify versions and runtime configurations, audit VPC and network settings, and find clusters by tag or status. Other TKE tables (`tencentcloud_tke_cluster_instance`) fan out from cluster IDs — join on `cluster_id` to drill into cluster instances (nodes).

## Examples

### Basic info

List every TKE cluster and its core attributes in the current region.

```sql+postgres
select
  cluster_id,
  cluster_name,
  cluster_type,
  cluster_status,
  cluster_version,
  cluster_node_num,
  region
from
  tencentcloud_tke_cluster;
```

```sql+sqlite
select
  cluster_id,
  cluster_name,
  cluster_type,
  cluster_status,
  cluster_version,
  cluster_node_num,
  region
from
  tencentcloud_tke_cluster;
```

### Filter clusters by status

Find all running clusters.

```sql+postgres
select
  cluster_id,
  cluster_name,
  cluster_version,
  vpc_id,
  cluster_cidr
from
  tencentcloud_tke_cluster
where
  cluster_status = 'Running';
```

```sql+sqlite
select
  cluster_id,
  cluster_name,
  cluster_version,
  vpc_id,
  cluster_cidr
from
  tencentcloud_tke_cluster
where
  cluster_status = 'Running';
```

### Count clusters by type and status

```sql+postgres
select
  cluster_type,
  cluster_status,
  count(*) as total
from
  tencentcloud_tke_cluster
group by
  cluster_type,
  cluster_status
order by
  total desc;
```

```sql+sqlite
select
  cluster_type,
  cluster_status,
  count(*) as total
from
  tencentcloud_tke_cluster
group by
  cluster_type,
  cluster_status
order by
  total desc;
```

### Find clusters without deletion protection

```sql+postgres
select
  cluster_id,
  cluster_name,
  cluster_status,
  deletion_protection
from
  tencentcloud_tke_cluster
where
  not deletion_protection;
```

```sql+sqlite
select
  cluster_id,
  cluster_name,
  cluster_status,
  deletion_protection
from
  tencentcloud_tke_cluster
where
  deletion_protection = 0;
```

### Lookup a specific cluster by ID

```sql+postgres
select
  cluster_id,
  cluster_name,
  cluster_description,
  cluster_version,
  container_runtime,
  runtime_version,
  network_settings,
  tags
from
  tencentcloud_tke_cluster
where
  cluster_id = 'cls-xxxxxxxx';
```

```sql+sqlite
select
  cluster_id,
  cluster_name,
  cluster_description,
  cluster_version,
  container_runtime,
  runtime_version,
  network_settings,
  tags
from
  tencentcloud_tke_cluster
where
  cluster_id = 'cls-xxxxxxxx';
```
