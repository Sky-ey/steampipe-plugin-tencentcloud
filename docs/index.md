---
organization: TencentCloud
category: ["public cloud"]
icon_url: "/images/plugins/tencentcloud/tencentcloud.svg"
brand_color: "#0052D9"
display_name: "Tencent Cloud"
short_name: "tencentcloud"
description: "Steampipe plugin for querying Tencent Cloud resources."
og_description: "Query Tencent Cloud with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/tencentcloud/tencentcloud-social-graphic.png"
---

# Tencent Cloud Plugin

> Steampipe plugin to query Tencent Cloud resources using SQL.

English | [简体中文](index_zh.md)

## Overview

Use SQL to query from your Tencent Cloud account — across every region you configure in a single query.

## Getting Started

### Installation

Build and install from source:

```bash
git clone https://github.com/TencentCloud/steampipe-plugin-tencentcloud.git
cd steampipe-plugin-tencentcloud
make install
```

## Credentials

The plugin authenticates with a Tencent Cloud [access key](https://www.tencentcloud.com/document/product/598) — a `SecretId` / `SecretKey` pair. Create or manage your keys in the console under **Access Management → API Key Management** ([CAM console](https://console.tencentcloud.com/cam/capi)).

Credentials are resolved in this order:

1. Static credentials in the connection config (`secret_id` / `secret_key`, plus optional `token` for STS temporary credentials)
2. The `TENCENTCLOUD_SECRET_ID`, `TENCENTCLOUD_SECRET_KEY` and `TENCENTCLOUD_TOKEN` environment variables
3. CVM instance role when running on a Tencent Cloud CVM with an attached CAM role

### Required permissions

All tables are read-only: each service only requires the corresponding `Describe*` / `List*` API permissions. The simplest setup is to attach the per-service read-only preset policies (e.g. `QcloudCVMReadOnlyAccess`, `QcloudCBSReadOnlyAccess`, `QcloudVPCReadOnlyAccess`) to the CAM user or role used by the plugin. See the [CAM documentation](https://www.tencentcloud.com/document/product/598) for policy management.

### STS temporary credentials

For cross-account access, generate temporary credentials via [STS](https://www.tencentcloud.com/document/product/598) and set `secret_id`, `secret_key` and `token` (or `TENCENTCLOUD_TOKEN`). Temporary credentials expire on their own schedule and are not refreshed automatically — once expired, the plugin returns an `AuthFailure` error for every API call. Obtain a fresh token and update the configuration.

### Configuration

Create or edit `~/.steampipe/config/tencentcloud.spc`:

```hcl
connection "tencentcloud" {
  plugin = "tencentcloud/tencentcloud"

  secret_id  = "your-secret-id"
  secret_key = "your-secret-key"

  regions = ["ap-*"]
}
```

The `regions` argument controls which regions each query targets. See [Multi-region queries](#multi-region-queries) for the matching rules and the behaviour when it is omitted.

Alternatively, use environment variables:

```bash
export TENCENTCLOUD_SECRET_ID=your-secret-id
export TENCENTCLOUD_SECRET_KEY=your-secret-key
```

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

## Multi-region queries

Most tables are regional. A Tencent Cloud region determines which regional control plane the API call is routed to — it is not a filter applied to a global result set. Resource metadata is isolated per region, so a query sent to the wrong region returns an empty result instead of an error.

### The `regions` argument

`regions` lists the regions each query should target. API calls are made to all matching regions in parallel and the rows are merged into a single result set.

```hcl
connection "tencentcloud" {
  plugin = "tencentcloud/tencentcloud"

  secret_id  = "your-secret-id"
  secret_key = "your-secret-key"

  # regions = ["*"]                           # All available regions
  # regions = ["ap-*"]                        # All Asia Pacific regions
  # regions = ["ap-guangzhou", "ap-shanghai"] # Specific regions
  regions = ["ap-*", "eu-frankfurt"]          # Wildcards and exact values may be mixed
}
```

Supported wildcards:

| Pattern | Meaning | Example |
|---------|---------|---------|
| `*` | Any number of characters | `ap-*` matches all Asia Pacific regions |
| `?` | Exactly one character | `ap-?ongkong` matches `ap-hongkong` |
| `[gs]` | Any character in the set | `ap-[gs]*` matches `ap-guangzhou`, `ap-shanghai`, `ap-singapore` |
| `[^g]` | Any character not in the set | `ap-[^g]*` excludes `ap-guangzhou` |

Brace expansion (`{a,b}`) is **not** supported. Values are case-insensitive.

Patterns are matched against the regions returned by the region API for that specific product, so `regions` can only narrow the set — it cannot add regions that do not exist or that the product does not support. If `regions` matches no available region, the query fails with an explicit error rather than returning zero rows.

### When `regions` is omitted

The plugin queries only the default region `ap-singapore` and logs a warning:

```
You are currently using default region ap-singapore, instances in other regions will not be displayed here!
```

Resources in every other region are absent from the results. Set `regions` explicitly to query more than one region.

### Narrowing a single query

Every regional table exposes a `region` column. Adding a `region` filter targets that one region and skips region discovery entirely, regardless of the `regions` configuration:

```sql
select
  instance_id,
  instance_name,
  region
from
  tencentcloud_cvm_instance
where
  region = 'ap-shanghai';
```

The catalog table `tencentcloud_cvm_region` is global and does not iterate over regions.

## Multi-account queries

Every table exposes an `owner_uin` column holding the Tencent Cloud owner account UIN (`OwnerUin`) that owns the resource, plus the companion `caller_uin` (the calling identity) and `owner_app_id` columns. The plugin resolves these once per connection by calling the CAM `GetUserAppId` API — the cost is a single API call per connection regardless of how many tables or regions are queried.

`owner_uin` is registered as the plugin's **connection key column**, which lets Steampipe prune child connections when an aggregator query filters by account. To query multiple accounts in one place, define one connection per account and an `aggregator` connection that combines them:

```hcl
# Development account
connection "tencentcloud_dev" {
  plugin     = "tencentcloud/tencentcloud"
  secret_id  = "AKIDxxxxxxxx"
  secret_key = "xxxxxxxxxxxx"
  regions    = ["ap-guangzhou", "ap-shanghai"]
}

# Production account
connection "tencentcloud_prod" {
  plugin     = "tencentcloud/tencentcloud"
  secret_id  = "AKIDyyyyyyyy"
  secret_key = "yyyyyyyyyyyy"
  regions    = ["*"]
}

# Aggregator: merges rows from every matching child connection
connection "tencentcloud_all" {
  plugin      = "tencentcloud/tencentcloud"
  type        = "aggregator"
  connections = ["tencentcloud_*"]
}
```

Query the aggregator to see rows from every account at once:

```sql
-- Count CVM instances per account, per region
select
  owner_uin,
  region,
  count(*) as instance_count
from
  tencentcloud_all.tencentcloud_cvm_instance
group by
  owner_uin, region;

-- Filter to a single account — Steampipe skips the other child connections,
-- avoiding unnecessary API calls against the wrong account
select instance_id, instance_name, owner_uin
from tencentcloud_all.tencentcloud_cvm_instance
where owner_uin = '100012345678';
```

| Column | Source | Description |
|--------|--------|-------------|
| `owner_uin` | `OwnerUin` | The master account UIN that owns the resource. Connection key column used for aggregator pruning. |
| `caller_uin` | `Uin` | The calling identity UIN — the master account when using root credentials, or the sub-account/role UIN when called via CAM. |
| `owner_app_id` | `AppId` | The numeric application ID. Matches the `<appid>` suffix of COS bucket names (e.g. `my-bucket-1250000000`). |

A sub-account connection and its master-account connection share the same `owner_uin` (both resolve to the same `OwnerUin`). Use `caller_uin` to distinguish the calling identity when needed.

## Tables

| Table | Description |
|-------|-------------|
| [tencentcloud_cam_access_key](tables/tencentcloud_cam_access_key.md) | CAM access keys of sub-accounts. |
| [tencentcloud_cam_group](tables/tencentcloud_cam_group.md) | CAM user groups. |
| [tencentcloud_cam_policy](tables/tencentcloud_cam_policy.md) | CAM policies (preset and custom). |
| [tencentcloud_cam_policy_attachment](tables/tencentcloud_cam_policy_attachment.md) | CAM policy attachments (policy-to-principal mapping). |
| [tencentcloud_cam_role](tables/tencentcloud_cam_role.md) | CAM roles. |
| [tencentcloud_cam_saml_provider](tables/tencentcloud_cam_saml_provider.md) | CAM SAML identity providers for enterprise SSO. |
| [tencentcloud_cam_user](tables/tencentcloud_cam_user.md) | CAM sub-accounts (IAM users). |
| [tencentcloud_ccn](tables/tencentcloud_ccn.md) | CCN (Cloud Connect Network) instances. |
| [tencentcloud_ccn_instance](tables/tencentcloud_ccn_instance.md) | Instances (VPCs, DC gateways) attached to a CCN. |
| [tencentcloud_ccn_route](tables/tencentcloud_ccn_route.md) | Route entries of a CCN. |
| [tencentcloud_cbs_auto_snapshot_policy](tables/tencentcloud_cbs_auto_snapshot_policy.md) | CBS auto snapshot policies. |
| [tencentcloud_cbs_disk](tables/tencentcloud_cbs_disk.md) | CBS cloud disks. |
| [tencentcloud_cbs_disk_backup](tables/tencentcloud_cbs_disk_backup.md) | CBS disk backup points. |
| [tencentcloud_cbs_disk_storage_pool](tables/tencentcloud_cbs_disk_storage_pool.md) | CBS dedicated cluster disk storage pools. |
| [tencentcloud_cbs_region](tables/tencentcloud_cbs_region.md) | Regions available to CBS. |
| [tencentcloud_cbs_snapshot](tables/tencentcloud_cbs_snapshot.md) | CBS snapshots. |
| [tencentcloud_cbs_snapshot_group](tables/tencentcloud_cbs_snapshot_group.md) | CBS snapshot groups. |
| [tencentcloud_clb_block_ip](tables/tencentcloud_clb_block_ip.md) | CLB blocked IP (blacklist) entries. |
| [tencentcloud_clb_customized_config](tables/tencentcloud_clb_customized_config.md) | CLB customized configurations. |
| [tencentcloud_clb_cross_target](tables/tencentcloud_clb_cross_target.md) | CLB cross-VPC backend targets (CCN cross-domain 2.0). |
| [tencentcloud_clb_instance](tables/tencentcloud_clb_instance.md) | CLB load balancer instances. |
| [tencentcloud_clb_listener](tables/tencentcloud_clb_listener.md) | CLB listeners. |
| [tencentcloud_clb_region](tables/tencentcloud_clb_region.md) | Regions available to CLB. |
| [tencentcloud_clb_rewrite](tables/tencentcloud_clb_rewrite.md) | CLB URL rewrite redirection rules. |
| [tencentcloud_clb_target](tables/tencentcloud_clb_target.md) | CLB backend targets bound to listeners. |
| [tencentcloud_clb_target_group](tables/tencentcloud_clb_target_group.md) | CLB target groups. |
| [tencentcloud_clb_target_group_instance](tables/tencentcloud_clb_target_group_instance.md) | CLB backend servers registered in target groups. |
| [tencentcloud_cls_alarm](tables/tencentcloud_cls_alarm.md) | CLS alarm policies. |
| [tencentcloud_cls_alarm_notice](tables/tencentcloud_cls_alarm_notice.md) | CLS alarm notification channels. |
| [tencentcloud_cls_config](tables/tencentcloud_cls_config.md) | CLS collection configurations. |
| [tencentcloud_cls_logset](tables/tencentcloud_cls_logset.md) | CLS logsets. |
| [tencentcloud_cls_machine_group](tables/tencentcloud_cls_machine_group.md) | CLS machine groups. |
| [tencentcloud_cls_region](tables/tencentcloud_cls_region.md) | Regions available to CLS. |
| [tencentcloud_cls_topic](tables/tencentcloud_cls_topic.md) | CLS log topics. |
| [tencentcloud_cos_bucket](tables/tencentcloud_cos_bucket.md) | COS buckets under the current account. |
| [tencentcloud_cos_object](tables/tencentcloud_cos_object.md) | Objects stored in COS buckets. |
| [tencentcloud_cos_region](tables/tencentcloud_cos_region.md) | Regions available to COS. |
| [tencentcloud_cvm_auto_scaling_group](tables/tencentcloud_cvm_auto_scaling_group.md) | Tencent Cloud Auto Scaling (AS) groups. |
| [tencentcloud_cvm_disaster_recover_group](tables/tencentcloud_cvm_disaster_recover_group.md) | CVM spread placement groups. |
| [tencentcloud_cvm_host](tables/tencentcloud_cvm_host.md) | Tencent Cloud Cloud Dedicated Host (CDH) instances. |
| [tencentcloud_cvm_image](tables/tencentcloud_cvm_image.md) | CVM images (public, custom and shared). |
| [tencentcloud_cvm_instance](tables/tencentcloud_cvm_instance.md) | Tencent Cloud Virtual Machine (CVM) instances. |
| [tencentcloud_cvm_metric_cpu_daily](tables/tencentcloud_cvm_metric_cpu_daily.md) | Daily CPU utilization metric of CVM instances over the last 30 days. |
| [tencentcloud_cvm_metric_cpu_hourly](tables/tencentcloud_cvm_metric_cpu_hourly.md) | Hourly CPU utilization metric of CVM instances over the last 24 hours. |
| [tencentcloud_cvm_metric_disk_daily](tables/tencentcloud_cvm_metric_disk_daily.md) | Daily disk utilization metric of CVM instances over the last 30 days. |
| [tencentcloud_cvm_metric_disk_hourly](tables/tencentcloud_cvm_metric_disk_hourly.md) | Hourly disk utilization metric of CVM instances over the last 24 hours. |
| [tencentcloud_cvm_metric_lan_in_daily](tables/tencentcloud_cvm_metric_lan_in_daily.md) | Daily private network inbound bandwidth metric of CVM instances over the last 30 days. |
| [tencentcloud_cvm_metric_lan_in_hourly](tables/tencentcloud_cvm_metric_lan_in_hourly.md) | Hourly private network inbound bandwidth metric of CVM instances over the last 24 hours. |
| [tencentcloud_cvm_metric_lan_out_daily](tables/tencentcloud_cvm_metric_lan_out_daily.md) | Daily private network outbound bandwidth metric of CVM instances over the last 30 days. |
| [tencentcloud_cvm_metric_lan_out_hourly](tables/tencentcloud_cvm_metric_lan_out_hourly.md) | Hourly private network outbound bandwidth metric of CVM instances over the last 24 hours. |
| [tencentcloud_cvm_metric_mem_daily](tables/tencentcloud_cvm_metric_mem_daily.md) | Daily memory utilization metric of CVM instances over the last 30 days. |
| [tencentcloud_cvm_metric_mem_hourly](tables/tencentcloud_cvm_metric_mem_hourly.md) | Hourly memory utilization metric of CVM instances over the last 24 hours. |
| [tencentcloud_cvm_metric_wan_in_daily](tables/tencentcloud_cvm_metric_wan_in_daily.md) | Daily public network inbound bandwidth metric of CVM instances over the last 30 days. |
| [tencentcloud_cvm_metric_wan_in_hourly](tables/tencentcloud_cvm_metric_wan_in_hourly.md) | Hourly public network inbound bandwidth metric of CVM instances over the last 24 hours. |
| [tencentcloud_cvm_metric_wan_out_daily](tables/tencentcloud_cvm_metric_wan_out_daily.md) | Daily public network outbound bandwidth metric of CVM instances over the last 30 days. |
| [tencentcloud_cvm_metric_wan_out_hourly](tables/tencentcloud_cvm_metric_wan_out_hourly.md) | Hourly public network outbound bandwidth metric of CVM instances over the last 24 hours. |
| [tencentcloud_cvm_key_pair](tables/tencentcloud_cvm_key_pair.md) | SSH key pairs that can be bound to instances. |
| [tencentcloud_cvm_launch_template](tables/tencentcloud_cvm_launch_template.md) | Instance launch templates. |
| [tencentcloud_cvm_region](tables/tencentcloud_cvm_region.md) | Regions available to CVM. |
| [tencentcloud_cvm_zone](tables/tencentcloud_cvm_zone.md) | Availability zones available to CVM. |
| [tencentcloud_dc_direct_connect](tables/tencentcloud_dc_direct_connect.md) | DC physical direct connect lines. |
| [tencentcloud_dc_direct_connect_tunnel](tables/tencentcloud_dc_direct_connect_tunnel.md) | DC dedicated tunnels. |
| [tencentcloud_dc_gateway](tables/tencentcloud_dc_gateway.md) | DC gateways. |
| [tencentcloud_dc_gateway_ccn_route](tables/tencentcloud_dc_gateway_ccn_route.md) | CCN routes of a DC gateway. |
| [tencentcloud_dc_internet_address](tables/tencentcloud_dc_internet_address.md) | Internet public addresses allocated by DC. |
| [tencentcloud_dc_region](tables/tencentcloud_dc_region.md) | Regions available to DC. |
| [tencentcloud_ssl_certificate](tables/tencentcloud_ssl_certificate.md) | SSL certificates. |
| [tencentcloud_tke_cluster](tables/tencentcloud_tke_cluster.md) | TKE managed and independent Kubernetes clusters. |
| [tencentcloud_tke_cluster_instance](tables/tencentcloud_tke_cluster_instance.md) | CVM instances managed by TKE clusters. |
| [tencentcloud_tke_region](tables/tencentcloud_tke_region.md) | Regions available to TKE. |
| [tencentcloud_vpc](tables/tencentcloud_vpc.md) | VPC instances. |
| [tencentcloud_vpc_bandwidth_package](tables/tencentcloud_vpc_bandwidth_package.md) | VPC bandwidth packages. |
| [tencentcloud_vpc_eip](tables/tencentcloud_vpc_eip.md) | VPC elastic public IP (EIP) addresses. |
| [tencentcloud_vpc_eni](tables/tencentcloud_vpc_eni.md) | VPC elastic network interfaces (ENI). |
| [tencentcloud_vpc_nat_gateway](tables/tencentcloud_vpc_nat_gateway.md) | VPC NAT gateways. |
| [tencentcloud_vpc_region](tables/tencentcloud_vpc_region.md) | Regions available to VPC. |
| [tencentcloud_vpc_route](tables/tencentcloud_vpc_route.md) | VPC routing entries expanded from route tables. |
| [tencentcloud_vpc_route_table](tables/tencentcloud_vpc_route_table.md) | VPC route tables. |
| [tencentcloud_vpc_security_group](tables/tencentcloud_vpc_security_group.md) | VPC security groups. |
| [tencentcloud_vpc_security_group_policy](tables/tencentcloud_vpc_security_group_policy.md) | VPC security group rules (ingress and egress). |
| [tencentcloud_vpc_subnet](tables/tencentcloud_vpc_subnet.md) | VPC subnets. |
| [tencentcloud_vpc_traffic_package](tables/tencentcloud_vpc_traffic_package.md) | VPC traffic packages. |
| [tencentcloud_vpn_gateway](tables/tencentcloud_vpn_gateway.md) | VPN gateways. |
| [tencentcloud_vpn_customer_gateway](tables/tencentcloud_vpn_customer_gateway.md) | VPN customer gateways. |
| [tencentcloud_vpn_connection](tables/tencentcloud_vpn_connection.md) | VPN connections (IPsec tunnels). |
| [tencentcloud_vpn_gateway_route](tables/tencentcloud_vpn_gateway_route.md) | Routes of a VPN gateway. |
| [tencentcloud_vpn_gateway_ccn_route](tables/tencentcloud_vpn_gateway_ccn_route.md) | CCN routes of a VPN gateway. |
| [tencentcloud_vpn_region](tables/tencentcloud_vpn_region.md) | Regions available to VPN. |
| [tencentcloud_waf_clb_host](tables/tencentcloud_waf_clb_host.md) | WAF CLB-type protected domains. |
| [tencentcloud_waf_custom_rule](tables/tencentcloud_waf_custom_rule.md) | WAF custom access control rules. |
| [tencentcloud_waf_domain](tables/tencentcloud_waf_domain.md) | WAF SaaS-type protected domains. |
| [tencentcloud_waf_instance](tables/tencentcloud_waf_instance.md) | WAF instances. |
| [tencentcloud_waf_ip_access_control](tables/tencentcloud_waf_ip_access_control.md) | WAF IP allowlist/blocklist entries. |
| [tencentcloud_waf_object](tables/tencentcloud_waf_object.md) | WAF protection objects. |
