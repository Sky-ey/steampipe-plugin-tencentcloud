---
title: "Steampipe Table: tencentcloud_cvm_instance - Query Tencent Cloud CVM Instances using SQL"
description: "Allows users to query CVM instances in Tencent Cloud to retrieve information about virtual machine resources."
folder: "CVM"
---

# Table: tencentcloud_cvm_instance - Query Tencent Cloud CVM Instances using SQL

Tencent Cloud Virtual Machine (CVM) provides scalable, on-demand cloud computing capacity. CVM instances allow you to deploy applications, run services, and manage workloads in the cloud without investing in physical hardware. Instances come in a variety of types optimized for different use cases — from general-purpose computing to memory-intensive and GPU-accelerated workloads — and support both prepaid (subscription) and postpaid (pay-as-you-go) billing models.

## Table Usage Guide

The `tencentcloud_cvm_instance` table in Steampipe provides you with detailed information about CVM instances within your Tencent Cloud account. This table allows you, as a DevOps engineer or cloud administrator, to query instance-specific details, including the instance ID, state, type, CPU and memory configuration, network settings, and associated metadata. You can utilize this table to gather insights on your instances, such as identifying running vs. stopped instances, checking billing methods and expiry times, auditing security group associations, and reviewing tag-based classifications. The schema outlines the various attributes of the CVM instance, including the instance ID, creation time, operating system, VPC and subnet information, IP addresses, and associated tags.

## Examples

### Basic info

Get an overview of all CVM instances in your account, including their names, states, instance types, and regions. This is useful for inventory management and understanding the current state of your cloud infrastructure.

```sql+postgres
select
  instance_id,
  instance_name,
  instance_state,
  instance_type,
  zone
from
  tencentcloud_cvm_instance;
```

```sql+sqlite
select
  instance_id,
  instance_name,
  instance_state,
  instance_type,
  zone
from
  tencentcloud_cvm_instance;
```

### List running instances

Identify all running CVM instances to understand which workloads are currently active. This helps in monitoring operational status and ensuring that expected services are up and running.

```sql+postgres
select
  instance_id,
  instance_name,
  instance_type,
  private_ip_address,
  public_ip_address
from
  tencentcloud_cvm_instance
where
  instance_state = 'RUNNING';
```

```sql+sqlite
select
  instance_id,
  instance_name,
  instance_type,
  private_ip_address,
  public_ip_address
from
  tencentcloud_cvm_instance
where
  instance_state = 'RUNNING';
```

### Count instances by state

Summarize the number of instances in each state to quickly assess the overall health and distribution of your compute resources. This aids in identifying stopped instances that may need attention or unnecessary resources that can be decommissioned.

```sql+postgres
select
  instance_state,
  count(*) as count
from
  tencentcloud_cvm_instance
group by
  instance_state
order by
  count desc;
```

```sql+sqlite
select
  instance_state,
  count(*) as count
from
  tencentcloud_cvm_instance
group by
  instance_state
order by
  count desc;
```

### Find instances with public IPs

Discover instances that have public IP addresses assigned, which are accessible from the internet. This is important for security auditing to ensure that only intended instances are exposed publicly and that appropriate security groups are in place.

```sql+postgres
select
  instance_id,
  instance_name,
  public_ip_address,
  instance_type,
  security_group_ids
from
  tencentcloud_cvm_instance
where
  public_ip_address is not null;
```

```sql+sqlite
select
  instance_id,
  instance_name,
  public_ip_address,
  instance_type,
  security_group_ids
from
  tencentcloud_cvm_instance
where
  public_ip_address is not null;
```

### List prepaid instances approaching expiry

Identify prepaid (subscription-based) instances that are nearing their expiry time so you can take action on renewal or migration. This helps avoid unexpected service interruptions due to expired subscriptions.

```sql+postgres
select
  instance_id,
  instance_name,
  instance_charge_type,
  expired_time,
  renew_flag
from
  tencentcloud_cvm_instance
where
  instance_charge_type = 'PREPAID'
  and expired_time is not null
order by
  expired_time asc;
```

```sql+sqlite
select
  instance_id,
  instance_name,
  instance_charge_type,
  expired_time,
  renew_flag
from
  tencentcloud_cvm_instance
where
  instance_charge_type = 'PREPAID'
  and expired_time is not null
order by
  expired_time asc;
```

### Find instances with specific tags

Query instances that have been tagged with a specific key-value pair, such as an environment classification. Tag-based filtering is useful for organizing resources and applying policies based on metadata.

```sql+postgres
select
  instance_id,
  instance_name,
  tags
from
  tencentcloud_cvm_instance
where
  tags ->> 'Environment' = 'Production';
```

```sql+sqlite
select
  instance_id,
  instance_name,
  tags
from
  tencentcloud_cvm_instance
where
  json_extract(tags, '$.Environment') = 'Production';
```

### List instances in a specific VPC and subnet

Filter instances by their VPC and subnet to understand the network topology and verify that instances are deployed in the correct network segments. This is essential for network segmentation audits and troubleshooting connectivity issues.

```sql+postgres
select
  instance_id,
  instance_name,
  vpc_id,
  subnet_id,
  private_ip_address,
  instance_state
from
  tencentcloud_cvm_instance
where
  vpc_id = 'vpc-xxxxxxxx'
  and subnet_id = 'subnet-xxxxxxxx';
```

```sql+sqlite
select
  instance_id,
  instance_name,
  vpc_id,
  subnet_id,
  private_ip_address,
  instance_state
from
  tencentcloud_cvm_instance
where
  vpc_id = 'vpc-xxxxxxxx'
  and subnet_id = 'subnet-xxxxxxxx';
```

### Check instances with failed latest operations

Identify instances where the most recent operation (such as a stop, start, or reset) has failed. This is critical for operational monitoring and quick incident response, allowing you to investigate and remediate issues promptly.

```sql+postgres
select
  instance_id,
  instance_name,
  latest_operation,
  latest_operation_state
from
  tencentcloud_cvm_instance
where
  latest_operation_state = 'FAILED';
```

```sql+sqlite
select
  instance_id,
  instance_name,
  latest_operation,
  latest_operation_state
from
  tencentcloud_cvm_instance
where
  latest_operation_state = 'FAILED';
```

### Review network billing and bandwidth configuration

Inspect the network billing method and outbound bandwidth cap of instances, useful for cost analysis and identifying instances billed by traffic with unexpectedly high bandwidth limits.

```sql+postgres
select
  instance_id,
  instance_name,
  internet_charge_type,
  internet_max_bandwidth_out,
  public_ip_assigned,
  bandwidth_package_id
from
  tencentcloud_cvm_instance
where
  public_ip_assigned;
```

```sql+sqlite
select
  instance_id,
  instance_name,
  internet_charge_type,
  internet_max_bandwidth_out,
  public_ip_assigned,
  bandwidth_package_id
from
  tencentcloud_cvm_instance
where
  public_ip_assigned;
```

### List instances by system disk type

Identify instances running on specific system disk types (for example, premium or SSD cloud disks) to plan storage performance upgrades or cost optimization.

```sql+postgres
select
  instance_id,
  instance_name,
  system_disk_id,
  system_disk_type,
  system_disk_size
from
  tencentcloud_cvm_instance
where
  system_disk_type = 'CLOUD_PREMIUM';
```

```sql+sqlite
select
  instance_id,
  instance_name,
  system_disk_id,
  system_disk_type,
  system_disk_size
from
  tencentcloud_cvm_instance
where
  system_disk_type = 'CLOUD_PREMIUM';
```

### List GPU instances and their GPU models

Find GPU-type instances along with the GPU model and count, useful for tracking accelerated computing resources.

```sql+postgres
select
  instance_id,
  instance_name,
  instance_type,
  gpu_type,
  gpu_count
from
  tencentcloud_cvm_instance
where
  gpu_count > 0;
```

```sql+sqlite
select
  instance_id,
  instance_name,
  instance_type,
  gpu_type,
  gpu_count
from
  tencentcloud_cvm_instance
where
  gpu_count > 0;
```
