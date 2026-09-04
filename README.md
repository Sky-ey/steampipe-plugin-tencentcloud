# Tencent Cloud Steampipe Plugin

> Steampipe plugin to query Tencent Cloud resources using SQL.

English | [简体中文](README_zh.md)

## Overview

Use SQL to query CVM instances, CBS disks, VPC networks, subnets, and more from your Tencent Cloud account.

## Getting Started

### Installation

Build and install from source:

```bash
git clone https://github.com/TencentCloud/steampipe-plugin-tencentcloud.git
cd steampipe-plugin-tencentcloud
make install
```

### Configuration

Create or edit `~/.steampipe/config/tencentcloud.spc`:

```hcl
connection "tencentcloud" {
  plugin = "tencentcloud/tencentcloud"

  secret_id  = "your-secret-id"
  secret_key = "your-secret-key"

  # Regions to query when no region is specified. Supports wildcards (`*`, `?`, `[gs]`).
  # If omitted, only the default region `ap-singapore` is queried when no region specified in SQL.
  regions = ["ap-*"]
}
```

Alternatively, use environment variables:

```bash
export TENCENTCLOUD_SECRET_ID=your-secret-id
export TENCENTCLOUD_SECRET_KEY=your-secret-key
```

Regional tables query every region matched by `regions`. Use a SQL condition such as `where region = 'ap-guangzhou'` to target one region directly and skip region discovery. See [docs/index.md](docs/index.md#multi-region-queries) for the full matching rules.

### Start Querying

```bash
steampipe query
```

```sql
-- List all CVM instances
select
  instance_id,
  instance_name,
  instance_state,
  instance_type,
  region
from
  tencentcloud_cvm_instance;

-- Find running instances in a specific region
select
  instance_id,
  instance_name,
  instance_type
from
  tencentcloud_cvm_instance
where
  region = 'ap-guangzhou'
  and instance_state = 'RUNNING';
```

## Testing

Run the fast Go unit tests:

```bash
make test
```

Run the SQL query contract suite against the local mock Tencent Cloud API. This covers List/Get queries for every table, a cross-service join, and validates API actions, paths, request bodies, regions, and request counts:

```bash
make test-query
```

The contract suite uses an isolated temporary Steampipe install directory and never writes to `~/.steampipe`.

## Tables

The plugin currently exposes **94 tables** across 13 Tencent Cloud services:

| Table | Description |
|-------|-------------|
| `tencentcloud_cam_access_key` | CAM access keys of sub-accounts. |
| `tencentcloud_cam_group` | CAM user groups. |
| `tencentcloud_cam_policy` | CAM policies (preset and custom). |
| `tencentcloud_cam_policy_attachment` | CAM policy attachments (policy-to-principal mapping). |
| `tencentcloud_cam_role` | CAM roles. |
| `tencentcloud_cam_saml_provider` | CAM SAML identity providers for enterprise SSO. |
| `tencentcloud_cam_user` | CAM sub-accounts (IAM users). |
| `tencentcloud_ccn` | CCN (Cloud Connect Network) instances. |
| `tencentcloud_ccn_instance` | Instances (VPCs, DC gateways) attached to a CCN. |
| `tencentcloud_ccn_route` | Route entries of a CCN. |
| `tencentcloud_cbs_auto_snapshot_policy` | CBS auto snapshot policies. |
| `tencentcloud_cbs_disk` | CBS cloud disks. |
| `tencentcloud_cbs_disk_backup` | CBS disk backup points. |
| `tencentcloud_cbs_disk_storage_pool` | CBS dedicated cluster disk storage pools. |
| `tencentcloud_cbs_region` | Regions available to CBS. |
| `tencentcloud_cbs_snapshot` | CBS snapshots. |
| `tencentcloud_cbs_snapshot_group` | CBS snapshot groups. |
| `tencentcloud_clb_block_ip` | CLB blocked IP (blacklist) entries. |
| `tencentcloud_clb_customized_config` | CLB customized configurations. |
| `tencentcloud_clb_cross_target` | CLB cross-VPC backend targets (CCN cross-domain 2.0). |
| `tencentcloud_clb_instance` | CLB load balancer instances. |
| `tencentcloud_clb_listener` | CLB listeners. |
| `tencentcloud_clb_region` | Regions available to CLB. |
| `tencentcloud_clb_rewrite` | CLB URL rewrite redirection rules. |
| `tencentcloud_clb_target` | CLB backend targets bound to listeners. |
| `tencentcloud_clb_target_group` | CLB target groups. |
| `tencentcloud_clb_target_group_instance` | CLB backend servers registered in target groups. |
| `tencentcloud_cls_alarm` | CLS alarm policies. |
| `tencentcloud_cls_alarm_notice` | CLS alarm notification channels. |
| `tencentcloud_cls_config` | CLS collection configurations. |
| `tencentcloud_cls_logset` | CLS logsets. |
| `tencentcloud_cls_machine_group` | CLS machine groups. |
| `tencentcloud_cls_region` | Regions available to CLS. |
| `tencentcloud_cls_topic` | CLS log topics. |
| `tencentcloud_cos_bucket` | COS buckets under the current account. |
| `tencentcloud_cos_object` | Objects stored in COS buckets. |
| `tencentcloud_cos_region` | Regions available to COS. |
| `tencentcloud_cvm_auto_scaling_group` | Tencent Cloud Auto Scaling (AS) groups. |
| `tencentcloud_cvm_disaster_recover_group` | CVM spread placement groups. |
| `tencentcloud_cvm_host` | Tencent Cloud Cloud Dedicated Host (CDH) instances. |
| `tencentcloud_cvm_image` | CVM images (public, custom and shared). |
| `tencentcloud_cvm_instance` | Tencent Cloud Virtual Machine (CVM) instances. |
| `tencentcloud_cvm_metric_cpu_daily` | Daily CPU utilization metric of CVM instances over the last 30 days. |
| `tencentcloud_cvm_metric_cpu_hourly` | Hourly CPU utilization metric of CVM instances over the last 24 hours. |
| `tencentcloud_cvm_metric_disk_daily` | Daily disk utilization metric of CVM instances over the last 30 days. |
| `tencentcloud_cvm_metric_disk_hourly` | Hourly disk utilization metric of CVM instances over the last 24 hours. |
| `tencentcloud_cvm_metric_lan_in_daily` | Daily private network inbound bandwidth metric of CVM instances over the last 30 days. |
| `tencentcloud_cvm_metric_lan_in_hourly` | Hourly private network inbound bandwidth metric of CVM instances over the last 24 hours. |
| `tencentcloud_cvm_metric_lan_out_daily` | Daily private network outbound bandwidth metric of CVM instances over the last 30 days. |
| `tencentcloud_cvm_metric_lan_out_hourly` | Hourly private network outbound bandwidth metric of CVM instances over the last 24 hours. |
| `tencentcloud_cvm_metric_mem_daily` | Daily memory utilization metric of CVM instances over the last 30 days. |
| `tencentcloud_cvm_metric_mem_hourly` | Hourly memory utilization metric of CVM instances over the last 24 hours. |
| `tencentcloud_cvm_metric_wan_in_daily` | Daily public network inbound bandwidth metric of CVM instances over the last 30 days. |
| `tencentcloud_cvm_metric_wan_in_hourly` | Hourly public network inbound bandwidth metric of CVM instances over the last 24 hours. |
| `tencentcloud_cvm_metric_wan_out_daily` | Daily public network outbound bandwidth metric of CVM instances over the last 30 days. |
| `tencentcloud_cvm_metric_wan_out_hourly` | Hourly public network outbound bandwidth metric of CVM instances over the last 24 hours. |
| `tencentcloud_cvm_key_pair` | SSH key pairs that can be bound to instances. |
| `tencentcloud_cvm_launch_template` | Instance launch templates. |
| `tencentcloud_cvm_region` | Regions available to CVM. |
| `tencentcloud_cvm_zone` | Availability zones available to CVM. |
| `tencentcloud_dc_direct_connect` | DC physical direct connect lines. |
| `tencentcloud_dc_direct_connect_tunnel` | DC dedicated tunnels. |
| `tencentcloud_dc_gateway` | DC gateways. |
| `tencentcloud_dc_gateway_ccn_route` | CCN routes of a DC gateway. |
| `tencentcloud_dc_internet_address` | Internet public addresses allocated by DC. |
| `tencentcloud_dc_region` | Regions available to DC. |
| `tencentcloud_ssl_certificate` | SSL certificates. |
| `tencentcloud_tke_cluster` | TKE managed and independent Kubernetes clusters. |
| `tencentcloud_tke_cluster_instance` | CVM instances managed by TKE clusters. |
| `tencentcloud_tke_region` | Regions available to TKE. |
| `tencentcloud_vpc` | VPC instances. |
| `tencentcloud_vpc_bandwidth_package` | VPC bandwidth packages. |
| `tencentcloud_vpc_eip` | VPC elastic public IP (EIP) addresses. |
| `tencentcloud_vpc_eni` | VPC elastic network interfaces (ENI). |
| `tencentcloud_vpc_nat_gateway` | VPC NAT gateways. |
| `tencentcloud_vpc_region` | Regions available to VPC. |
| `tencentcloud_vpc_route` | VPC routing entries expanded from route tables. |
| `tencentcloud_vpc_route_table` | VPC route tables. |
| `tencentcloud_vpc_security_group` | VPC security groups. |
| `tencentcloud_vpc_security_group_policy` | VPC security group rules (ingress and egress). |
| `tencentcloud_vpc_subnet` | VPC subnets. |
| `tencentcloud_vpc_traffic_package` | VPC traffic packages. |
| `tencentcloud_vpn_gateway` | VPN gateways. |
| `tencentcloud_vpn_customer_gateway` | VPN customer gateways. |
| `tencentcloud_vpn_connection` | VPN connections (IPsec tunnels). |
| `tencentcloud_vpn_gateway_route` | Routes of a VPN gateway. |
| `tencentcloud_vpn_gateway_ccn_route` | CCN routes of a VPN gateway. |
| `tencentcloud_vpn_region` | Regions available to VPN. |
| `tencentcloud_waf_clb_host` | WAF CLB-type protected domains. |
| `tencentcloud_waf_custom_rule` | WAF custom access control rules. |
| `tencentcloud_waf_domain` | WAF SaaS-type protected domains. |
| `tencentcloud_waf_instance` | WAF instances. |
| `tencentcloud_waf_ip_access_control` | WAF IP allowlist/blocklist entries. |
| `tencentcloud_waf_object` | WAF protection objects. |
