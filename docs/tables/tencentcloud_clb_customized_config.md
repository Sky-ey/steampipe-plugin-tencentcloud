---
title: "Steampipe Table: tencentcloud_clb_customized_config - Query Tencent Cloud CLB Customized Configurations using SQL"
description: "Query Tencent Cloud Load Balancer customized configurations, including type, name, and content."
folder: "CLB"
---

# Table: tencentcloud_clb_customized_config - Query Tencent Cloud CLB Customized Configurations using SQL

Tencent Cloud Load Balancer (CLB) customized configurations let you define reusable HTTP header and content customization settings at the instance (CLB), domain (SERVER), or forwarding rule (LOCATION) dimension. The `tencentcloud_clb_customized_config` table uses the `DescribeCustomizedConfigList` API and returns the configurations in each region, including their ID, type, name, and content.

## Table Usage Guide

The `tencentcloud_clb_customized_config` table in Steampipe provides information about CLB customized configurations within your Tencent Cloud account. This table allows you, as a DevOps engineer or cloud administrator, to query configuration details, including the config ID, type, name, and content, as well as creation and modification times. You can use this table to audit HTTP header customization settings, find configurations by dimension, and locate stale configs.

## Examples

### Basic configuration inventory

Get an overview of all customized configurations.

```sql+postgres
select
  uconfig_id,
  config_name,
  config_type,
  region
from
  tencentcloud_clb_customized_config;
```

```sql+sqlite
select
  uconfig_id,
  config_name,
  config_type,
  region
from
  tencentcloud_clb_customized_config;
```

### List instance-level configurations

The `config_type` qual is pushed down to the API.

```sql+postgres
select
  uconfig_id,
  config_name,
  config_content,
  region
from
  tencentcloud_clb_customized_config
where
  config_type = 'CLB';
```

```sql+sqlite
select
  uconfig_id,
  config_name,
  config_content,
  region
from
  tencentcloud_clb_customized_config
where
  config_type = 'CLB';
```

### Count configurations by type

```sql+postgres
select
  config_type,
  count(*) as config_count
from
  tencentcloud_clb_customized_config
group by
  config_type;
```

```sql+sqlite
select
  config_type,
  count(*) as config_count
from
  tencentcloud_clb_customized_config
group by
  config_type;
```

### Look up a configuration by name

The `config_name` qual supports fuzzy matching on the API side.

```sql+postgres
select
  uconfig_id,
  config_name,
  config_type,
  config_content,
  update_timestamp,
  region
from
  tencentcloud_clb_customized_config
where
  config_name = 'custom-header';
```

```sql+sqlite
select
  uconfig_id,
  config_name,
  config_type,
  config_content,
  update_timestamp,
  region
from
  tencentcloud_clb_customized_config
where
  config_name = 'custom-header';
```

### Look up a configuration by ID

The `uconfig_id` qual is pushed down to the API.

```sql+postgres
select
  uconfig_id,
  config_name,
  config_type,
  config_content,
  create_timestamp,
  update_timestamp,
  region
from
  tencentcloud_clb_customized_config
where
  uconfig_id = 'uconfig-12345678'
  and region = 'ap-guangzhou';
```

```sql+sqlite
select
  uconfig_id,
  config_name,
  config_type,
  config_content,
  create_timestamp,
  update_timestamp,
  region
from
  tencentcloud_clb_customized_config
where
  uconfig_id = 'uconfig-12345678'
  and region = 'ap-guangzhou';
```

### Find configurations by load balancer or VIP

The `load_balancer_id` and `vip` quals are pushed down to the API as filters (DescribeCustomizedConfigList `loadbalancer-id` / `vip`). These columns are filter-only — they echo the qual value you supplied and are not part of the DescribeCustomizedConfigList response payload.

```sql+postgres
select
  uconfig_id,
  config_name,
  config_type,
  load_balancer_id,
  vip,
  region
from
  tencentcloud_clb_customized_config
where
  load_balancer_id = 'lb-xxxxxxxx';
```

```sql+sqlite
select
  uconfig_id,
  config_name,
  config_type,
  load_balancer_id,
  vip,
  region
from
  tencentcloud_clb_customized_config
where
  load_balancer_id = 'lb-xxxxxxxx';
```
