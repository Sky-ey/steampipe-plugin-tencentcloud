---
title: "Steampipe Table: tencentcloud_cvm_launch_template - Query Tencent Cloud CVM Launch Templates using SQL"
description: "Allows users to query CVM launch templates in Tencent Cloud to retrieve information about reusable instance configuration templates."
folder: "CVM"
---

# Table: tencentcloud_cvm_launch_template - Query Tencent Cloud CVM Launch Templates using SQL

Tencent Cloud CVM launch templates store instance configuration parameters (such as image, instance type, security groups, and login settings) so that instances can be created quickly with a predefined setup. Each template supports multiple versions, with one version designated as the default used when no specific version is requested.

## Table Usage Guide

The `tencentcloud_cvm_launch_template` table in Steampipe provides you with detailed information about CVM launch templates within your Tencent Cloud account. This table allows you, as a DevOps engineer or cloud administrator, to query launch template details, including the template ID, name, default and latest version numbers, version count, creator, and creation time. You can utilize this table to audit template inventory, track version drift, and verify ownership of instance configurations.

Note: This table lists template metadata only. The per-version configuration parameters are available through the DescribeLaunchTemplateVersions API and are not exposed here.

## Examples

### Basic info

Get an overview of all launch templates in your account, including their names, IDs, and version information. This is useful for maintaining an inventory of reusable instance configurations.

```sql+postgres
select
  launch_template_id,
  launch_template_name,
  default_version_number,
  latest_version_number,
  launch_template_version_count,
  creation_time
from
  tencentcloud_cvm_launch_template;
```

```sql+sqlite
select
  launch_template_id,
  launch_template_name,
  default_version_number,
  latest_version_number,
  launch_template_version_count,
  creation_time
from
  tencentcloud_cvm_launch_template;
```

### Search launch templates by name

Filter launch templates by name. The filter is pushed down to the DescribeLaunchTemplates API, avoiding a full client-side scan.

```sql+postgres
select
  launch_template_id,
  launch_template_name,
  default_version_number,
  created_by
from
  tencentcloud_cvm_launch_template
where
  launch_template_name = 'my-web-template';
```

```sql+sqlite
select
  launch_template_id,
  launch_template_name,
  default_version_number,
  created_by
from
  tencentcloud_cvm_launch_template
where
  launch_template_name = 'my-web-template';
```

### Find launch templates whose default version lags the latest version

Identify templates where the default version is behind the latest version. New instances launched without an explicit version use the default, so these templates may roll out stale configurations.

```sql+postgres
select
  launch_template_id,
  launch_template_name,
  default_version_number,
  latest_version_number
from
  tencentcloud_cvm_launch_template
where
  default_version_number < latest_version_number;
```

```sql+sqlite
select
  launch_template_id,
  launch_template_name,
  default_version_number,
  latest_version_number
from
  tencentcloud_cvm_launch_template
where
  default_version_number < latest_version_number;
```

### Get a specific launch template

Retrieve details of a single launch template by its ID. The `launch_template_id` qualifier is pushed down to the DescribeLaunchTemplates API.

```sql+postgres
select
  launch_template_id,
  launch_template_name,
  launch_template_version_count,
  created_by,
  creation_time
from
  tencentcloud_cvm_launch_template
where
  launch_template_id = 'lt-xxxxxxxx';
```

```sql+sqlite
select
  launch_template_id,
  launch_template_name,
  launch_template_version_count,
  created_by,
  creation_time
from
  tencentcloud_cvm_launch_template
where
  launch_template_id = 'lt-xxxxxxxx';
```
