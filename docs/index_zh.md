---
organization: TencentCloud
category: ["public cloud"]
icon_url: "/images/plugins/tencentcloud/tencentcloud.svg"
brand_color: "#0052D9"
display_name: "Tencent Cloud"
short_name: "tencentcloud"
description: "用于查询腾讯云资源的 Steampipe 插件。"
og_description: "使用 SQL 查询腾讯云！开源 CLI，无需数据库。"
og_image: "/images/plugins/tencentcloud/tencentcloud-social-graphic.png"
---

# Tencent Cloud 插件

> 使用 SQL 查询腾讯云资源的 Steampipe 插件。

[English](index.md) | 简体中文

## 概述

使用 SQL 查询你腾讯云账号下的资源——单条查询即可覆盖你配置的所有地域。

## 快速开始

### 安装

从源码构建并安装：

```bash
git clone https://github.com/TencentCloud/steampipe-plugin-tencentcloud.git
cd steampipe-plugin-tencentcloud
make install
```

## 凭证

插件使用腾讯云[访问密钥](https://www.tencentcloud.com/document/product/598)进行身份认证——即 `SecretId` / `SecretKey` 密钥对。可在控制台 **访问管理 → API 密钥管理**（[CAM 控制台](https://console.tencentcloud.com/cam/capi)）中创建或管理密钥。

凭证按以下顺序解析：

1. 连接配置中的静态凭证（`secret_id` / `secret_key`，以及用于 STS 临时凭证的可选 `token`）
2. 环境变量 `TENCENTCLOUD_SECRET_ID`、`TENCENTCLOUD_SECRET_KEY` 与 `TENCENTCLOUD_TOKEN`
3. 在绑定了 CAM 角色的腾讯云 CVM 实例上运行时，使用 CVM 实例角色

### 所需权限

所有表均为只读：每个服务只需要对应的 `Describe*` / `List*` API 权限。最简单的配置方式是给插件使用的 CAM 用户或角色绑定各服务的只读预设策略（如 `QcloudCVMReadOnlyAccess`、`QcloudCBSReadOnlyAccess`、`QcloudVPCReadOnlyAccess`）。策略管理详见 [CAM 文档](https://www.tencentcloud.com/document/product/598)。

### STS 临时凭证

跨账号访问时，可通过 [STS](https://www.tencentcloud.com/document/product/598) 生成临时凭证，并设置 `secret_id`、`secret_key` 与 `token`（或 `TENCENTCLOUD_TOKEN`）。临时凭证到期后不会自动刷新——一旦过期，插件的每次 API 调用都会返回 `AuthFailure` 错误。此时需获取新的临时凭证并更新配置。

### 配置

创建或编辑 `~/.steampipe/config/tencentcloud.spc`：

```hcl
connection "tencentcloud" {
  plugin = "tencentcloud/tencentcloud"

  secret_id  = "your-secret-id"
  secret_key = "your-secret-key"

  regions = ["ap-*"]
}
```

`regions` 参数控制每次查询的目标地域。匹配规则及省略时的行为见[多地域查询](#多地域查询)。

也可以使用环境变量：

```bash
export TENCENTCLOUD_SECRET_ID=your-secret-id
export TENCENTCLOUD_SECRET_KEY=your-secret-key
```

### 开始查询

```bash
steampipe query
```

```sql
-- 列出所有 CVM 实例
select
  instance_id,
  instance_name,
  instance_state,
  instance_type,
  region
from
  tencentcloud_cvm_instance;

-- 查找指定地域下运行中的实例
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

## 多地域查询

大多数表是区域性的。腾讯云地域决定了 API 调用被路由到哪个区域的控制面——它不是作用于全局结果集的过滤器。资源元数据按地域隔离，因此发送到错误地域的查询会返回空结果而非报错。

### `regions` 参数

`regions` 列出每次查询要目标的地域。API 调用会并行发往所有匹配地域，行数据合并为单一结果集。

```hcl
connection "tencentcloud" {
  plugin = "tencentcloud/tencentcloud"

  secret_id  = "your-secret-id"
  secret_key = "your-secret-key"

  # regions = ["*"]                           # 所有可用地域
  # regions = ["ap-*"]                        # 所有亚太地域
  # regions = ["ap-guangzhou", "ap-shanghai"] # 指定地域
  regions = ["ap-*", "eu-frankfurt"]          # 通配符与精确值可以混用
}
```

支持的通配符：

| 模式 | 含义 | 示例 |
|------|------|------|
| `*` | 任意数量字符 | `ap-*` 匹配所有亚太地域 |
| `?` | 恰好一个字符 | `ap-?ongkong` 匹配 `ap-hongkong` |
| `[gs]` | 集合中的任一字符 | `ap-[gs]*` 匹配 `ap-guangzhou`、`ap-shanghai`、`ap-singapore` |
| `[^g]` | 不在集合中的任一字符 | `ap-[^g]*` 排除 `ap-guangzhou` |

**不支持**花括号展开（`{a,b}`）。取值不区分大小写。

模式匹配的对象是该产品通过地域 API 返回的地域列表，因此 `regions` 只能缩小范围——不能添加不存在的或产品不支持的地域。如果 `regions` 未匹配到任何可用地域，查询会直接报错而不是返回零行。

### 省略 `regions` 时

插件只查询默认地域 `ap-singapore` 并输出一条警告：

```
You are currently using default region ap-singapore, instances in other regions will not be displayed here!
```

其他地域的资源不会出现在结果中。要查询多个地域，请显式设置 `regions`。

### 单条查询收窄范围

每张区域表都提供 `region` 列。添加 `region` 过滤条件可以只查询该地域并完全跳过地域发现流程，不受 `regions` 配置影响：

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

目录表 `tencentcloud_cvm_region` 是全局表，不会按地域遍历。

## 多账号查询

每张表都提供 `owner_uin` 列，保存拥有该资源的腾讯云主账号 UIN（`OwnerUin`），并附带 `caller_uin`（调用方身份）与 `owner_app_id` 列。插件通过调用 CAM `GetUserAppId` API 每个连接解析一次——无论查询多少张表或多少个地域，开销都只有每连接一次 API 调用。

`owner_uin` 被注册为插件的 **connection key column**，当聚合查询按账号过滤时，Steampipe 可以据此裁剪子连接。要在一个地方查询多个账号，可为每个账号定义一个连接，并定义一个组合它们的 `aggregator` 连接：

```hcl
# 开发账号
connection "tencentcloud_dev" {
  plugin     = "tencentcloud/tencentcloud"
  secret_id  = "AKIDxxxxxxxx"
  secret_key = "xxxxxxxxxxxx"
  regions    = ["ap-guangzhou", "ap-shanghai"]
}

# 生产账号
connection "tencentcloud_prod" {
  plugin     = "tencentcloud/tencentcloud"
  secret_id  = "AKIDyyyyyyyy"
  secret_key = "yyyyyyyyyyyy"
  regions    = ["*"]
}

# 聚合连接：合并所有匹配子连接的行
connection "tencentcloud_all" {
  plugin      = "tencentcloud/tencentcloud"
  type        = "aggregator"
  connections = ["tencentcloud_*"]
}
```

查询聚合连接即可同时看到所有账号的行：

```sql
-- 按账号、按地域统计 CVM 实例数
select
  owner_uin,
  region,
  count(*) as instance_count
from
  tencentcloud_all.tencentcloud_cvm_instance
group by
  owner_uin, region;

-- 过滤到单个账号——Steampipe 会跳过其他子连接，
-- 避免向错误账号发起不必要的 API 调用
select instance_id, instance_name, owner_uin
from tencentcloud_all.tencentcloud_cvm_instance
where owner_uin = '100012345678';
```

| 列 | 来源字段 | 描述 |
|----|----------|------|
| `owner_uin` | `OwnerUin` | 拥有该资源的主账号 UIN。作为 connection key column 用于聚合连接裁剪。 |
| `caller_uin` | `Uin` | 调用方身份 UIN——使用主账号凭证时即主账号，通过 CAM 调用时为子账号/角色 UIN。 |
| `owner_app_id` | `AppId` | 数字形式的 APP ID。与 COS 存储桶名称的 `<appid>` 后缀一致（如 `my-bucket-1250000000`）。 |

子账号连接与其主账号连接的 `owner_uin` 相同（都解析到同一个 `OwnerUin`）。需要区分调用方身份时使用 `caller_uin`。

## 数据表

表文档（英文）见各表链接：

| 表名 | 描述 |
|------|------|
| [tencentcloud_cam_access_key](tables/tencentcloud_cam_access_key.md) | 子账号的 CAM 访问密钥。 |
| [tencentcloud_cam_group](tables/tencentcloud_cam_group.md) | CAM 用户组。 |
| [tencentcloud_cam_policy](tables/tencentcloud_cam_policy.md) | CAM 策略（预设策略与自定义策略）。 |
| [tencentcloud_cam_policy_attachment](tables/tencentcloud_cam_policy_attachment.md) | CAM 策略绑定关系（策略与主体的映射）。 |
| [tencentcloud_cam_role](tables/tencentcloud_cam_role.md) | CAM 角色。 |
| [tencentcloud_cam_saml_provider](tables/tencentcloud_cam_saml_provider.md) | CAM SAML 身份提供商（用于企业 SSO）。 |
| [tencentcloud_cam_user](tables/tencentcloud_cam_user.md) | CAM 子账号（IAM 用户）。 |
| [tencentcloud_ccn](tables/tencentcloud_ccn.md) | CCN（云联网）实例。 |
| [tencentcloud_ccn_instance](tables/tencentcloud_ccn_instance.md) | 关联到云联网的实例（VPC、专线网关等）。 |
| [tencentcloud_ccn_route](tables/tencentcloud_ccn_route.md) | 云联网路由条目。 |
| [tencentcloud_cbs_auto_snapshot_policy](tables/tencentcloud_cbs_auto_snapshot_policy.md) | CBS 自动快照策略。 |
| [tencentcloud_cbs_disk](tables/tencentcloud_cbs_disk.md) | CBS 云硬盘。 |
| [tencentcloud_cbs_disk_backup](tables/tencentcloud_cbs_disk_backup.md) | CBS 云硬盘备份点。 |
| [tencentcloud_cbs_disk_storage_pool](tables/tencentcloud_cbs_disk_storage_pool.md) | CBS 专属集群的云硬盘存储池。 |
| [tencentcloud_cbs_region](tables/tencentcloud_cbs_region.md) | CBS 可用地域。 |
| [tencentcloud_cbs_snapshot](tables/tencentcloud_cbs_snapshot.md) | CBS 快照。 |
| [tencentcloud_cbs_snapshot_group](tables/tencentcloud_cbs_snapshot_group.md) | CBS 快照组。 |
| [tencentcloud_clb_block_ip](tables/tencentcloud_clb_block_ip.md) | CLB 封禁 IP（黑名单）条目。 |
| [tencentcloud_clb_customized_config](tables/tencentcloud_clb_customized_config.md) | CLB 自定义配置。 |
| [tencentcloud_clb_cross_target](tables/tencentcloud_clb_cross_target.md) | CLB 跨域后端目标（CCN 跨域 2.0）。 |
| [tencentcloud_clb_instance](tables/tencentcloud_clb_instance.md) | CLB 负载均衡实例。 |
| [tencentcloud_clb_listener](tables/tencentcloud_clb_listener.md) | CLB 监听器。 |
| [tencentcloud_clb_region](tables/tencentcloud_clb_region.md) | CLB 可用地域。 |
| [tencentcloud_clb_rewrite](tables/tencentcloud_clb_rewrite.md) | CLB URL 重写重定向规则。 |
| [tencentcloud_clb_target](tables/tencentcloud_clb_target.md) | 监听器绑定的 CLB 后端目标。 |
| [tencentcloud_clb_target_group](tables/tencentcloud_clb_target_group.md) | CLB 目标组。 |
| [tencentcloud_clb_target_group_instance](tables/tencentcloud_clb_target_group_instance.md) | 目标组中注册的 CLB 后端服务器。 |
| [tencentcloud_cls_alarm](tables/tencentcloud_cls_alarm.md) | CLS 告警策略。 |
| [tencentcloud_cls_alarm_notice](tables/tencentcloud_cls_alarm_notice.md) | CLS 告警通知渠道。 |
| [tencentcloud_cls_config](tables/tencentcloud_cls_config.md) | CLS 采集配置。 |
| [tencentcloud_cls_logset](tables/tencentcloud_cls_logset.md) | CLS 日志集。 |
| [tencentcloud_cls_machine_group](tables/tencentcloud_cls_machine_group.md) | CLS 机器组。 |
| [tencentcloud_cls_region](tables/tencentcloud_cls_region.md) | CLS 可用地域。 |
| [tencentcloud_cls_topic](tables/tencentcloud_cls_topic.md) | CLS 日志主题。 |
| [tencentcloud_cos_bucket](tables/tencentcloud_cos_bucket.md) | 当前账号下的 COS 存储桶。 |
| [tencentcloud_cos_object](tables/tencentcloud_cos_object.md) | COS 存储桶中的对象。 |
| [tencentcloud_cos_region](tables/tencentcloud_cos_region.md) | COS 可用地域。 |
| [tencentcloud_cvm_auto_scaling_group](tables/tencentcloud_cvm_auto_scaling_group.md) | 腾讯云弹性伸缩（AS）组。 |
| [tencentcloud_cvm_disaster_recover_group](tables/tencentcloud_cvm_disaster_recover_group.md) | CVM 分散置放群组。 |
| [tencentcloud_cvm_host](tables/tencentcloud_cvm_host.md) | 腾讯云专用宿主机（CDH）实例。 |
| [tencentcloud_cvm_image](tables/tencentcloud_cvm_image.md) | CVM 镜像（公共镜像、自定义镜像与共享镜像）。 |
| [tencentcloud_cvm_instance](tables/tencentcloud_cvm_instance.md) | 腾讯云云服务器（CVM）实例。 |
| [tencentcloud_cvm_metric_cpu_daily](tables/tencentcloud_cvm_metric_cpu_daily.md) | CVM 实例近 30 天的每日 CPU 利用率指标。 |
| [tencentcloud_cvm_metric_cpu_hourly](tables/tencentcloud_cvm_metric_cpu_hourly.md) | CVM 实例近 24 小时的每小时 CPU 利用率指标。 |
| [tencentcloud_cvm_metric_disk_daily](tables/tencentcloud_cvm_metric_disk_daily.md) | CVM 实例近 30 天的每日磁盘利用率指标。 |
| [tencentcloud_cvm_metric_disk_hourly](tables/tencentcloud_cvm_metric_disk_hourly.md) | CVM 实例近 24 小时的每小时磁盘利用率指标。 |
| [tencentcloud_cvm_metric_lan_in_daily](tables/tencentcloud_cvm_metric_lan_in_daily.md) | CVM 实例近 30 天的每日内网入带宽指标。 |
| [tencentcloud_cvm_metric_lan_in_hourly](tables/tencentcloud_cvm_metric_lan_in_hourly.md) | CVM 实例近 24 小时的每小时内网入带宽指标。 |
| [tencentcloud_cvm_metric_lan_out_daily](tables/tencentcloud_cvm_metric_lan_out_daily.md) | CVM 实例近 30 天的每日内网出带宽指标。 |
| [tencentcloud_cvm_metric_lan_out_hourly](tables/tencentcloud_cvm_metric_lan_out_hourly.md) | CVM 实例近 24 小时的每小时内网出带宽指标。 |
| [tencentcloud_cvm_metric_mem_daily](tables/tencentcloud_cvm_metric_mem_daily.md) | CVM 实例近 30 天的每日内存利用率指标。 |
| [tencentcloud_cvm_metric_mem_hourly](tables/tencentcloud_cvm_metric_mem_hourly.md) | CVM 实例近 24 小时的每小时内存利用率指标。 |
| [tencentcloud_cvm_metric_wan_in_daily](tables/tencentcloud_cvm_metric_wan_in_daily.md) | CVM 实例近 30 天的每日公网入带宽指标。 |
| [tencentcloud_cvm_metric_wan_in_hourly](tables/tencentcloud_cvm_metric_wan_in_hourly.md) | CVM 实例近 24 小时的每小时公网入带宽指标。 |
| [tencentcloud_cvm_metric_wan_out_daily](tables/tencentcloud_cvm_metric_wan_out_daily.md) | CVM 实例近 30 天的每日公网出带宽指标。 |
| [tencentcloud_cvm_metric_wan_out_hourly](tables/tencentcloud_cvm_metric_wan_out_hourly.md) | CVM 实例近 24 小时的每小时公网出带宽指标。 |
| [tencentcloud_cvm_key_pair](tables/tencentcloud_cvm_key_pair.md) | 可绑定到实例的 SSH 密钥对。 |
| [tencentcloud_cvm_launch_template](tables/tencentcloud_cvm_launch_template.md) | 实例启动模板。 |
| [tencentcloud_cvm_region](tables/tencentcloud_cvm_region.md) | CVM 可用地域。 |
| [tencentcloud_cvm_zone](tables/tencentcloud_cvm_zone.md) | CVM 可用区。 |
| [tencentcloud_dc_direct_connect](tables/tencentcloud_dc_direct_connect.md) | DC 物理专线。 |
| [tencentcloud_dc_direct_connect_tunnel](tables/tencentcloud_dc_direct_connect_tunnel.md) | DC 专用通道。 |
| [tencentcloud_dc_gateway](tables/tencentcloud_dc_gateway.md) | DC 专线网关。 |
| [tencentcloud_dc_gateway_ccn_route](tables/tencentcloud_dc_gateway_ccn_route.md) | 专线网关的云联网路由。 |
| [tencentcloud_dc_internet_address](tables/tencentcloud_dc_internet_address.md) | DC 分配的互联网公网地址。 |
| [tencentcloud_dc_region](tables/tencentcloud_dc_region.md) | DC 可用地域。 |
| [tencentcloud_ssl_certificate](tables/tencentcloud_ssl_certificate.md) | SSL 证书。 |
| [tencentcloud_tke_cluster](tables/tencentcloud_tke_cluster.md) | TKE 托管与独立 Kubernetes 集群。 |
| [tencentcloud_tke_cluster_instance](tables/tencentcloud_tke_cluster_instance.md) | TKE 集群管理的 CVM 实例。 |
| [tencentcloud_tke_region](tables/tencentcloud_tke_region.md) | TKE 可用地域。 |
| [tencentcloud_vpc](tables/tencentcloud_vpc.md) | VPC 实例。 |
| [tencentcloud_vpc_bandwidth_package](tables/tencentcloud_vpc_bandwidth_package.md) | VPC 带宽包。 |
| [tencentcloud_vpc_eip](tables/tencentcloud_vpc_eip.md) | VPC 弹性公网 IP（EIP）。 |
| [tencentcloud_vpc_eni](tables/tencentcloud_vpc_eni.md) | VPC 弹性网卡（ENI）。 |
| [tencentcloud_vpc_nat_gateway](tables/tencentcloud_vpc_nat_gateway.md) | VPC NAT 网关。 |
| [tencentcloud_vpc_region](tables/tencentcloud_vpc_region.md) | VPC 可用地域。 |
| [tencentcloud_vpc_route](tables/tencentcloud_vpc_route.md) | 从路由表展开的 VPC 路由条目。 |
| [tencentcloud_vpc_route_table](tables/tencentcloud_vpc_route_table.md) | VPC 路由表。 |
| [tencentcloud_vpc_security_group](tables/tencentcloud_vpc_security_group.md) | VPC 安全组。 |
| [tencentcloud_vpc_security_group_policy](tables/tencentcloud_vpc_security_group_policy.md) | VPC 安全组规则（入站与出站）。 |
| [tencentcloud_vpc_subnet](tables/tencentcloud_vpc_subnet.md) | VPC 子网。 |
| [tencentcloud_vpc_traffic_package](tables/tencentcloud_vpc_traffic_package.md) | VPC 流量包。 |
| [tencentcloud_vpn_gateway](tables/tencentcloud_vpn_gateway.md) | VPN 网关。 |
| [tencentcloud_vpn_customer_gateway](tables/tencentcloud_vpn_customer_gateway.md) | VPN 对端网关。 |
| [tencentcloud_vpn_connection](tables/tencentcloud_vpn_connection.md) | VPN 通道（IPsec 隧道）。 |
| [tencentcloud_vpn_gateway_route](tables/tencentcloud_vpn_gateway_route.md) | VPN 网关的路由条目。 |
| [tencentcloud_vpn_gateway_ccn_route](tables/tencentcloud_vpn_gateway_ccn_route.md) | VPN 网关的云联网路由。 |
| [tencentcloud_vpn_region](tables/tencentcloud_vpn_region.md) | VPN 可用地域。 |
| [tencentcloud_waf_clb_host](tables/tencentcloud_waf_clb_host.md) | WAF CLB 型防护域名。 |
| [tencentcloud_waf_custom_rule](tables/tencentcloud_waf_custom_rule.md) | WAF 自定义访问控制规则。 |
| [tencentcloud_waf_domain](tables/tencentcloud_waf_domain.md) | WAF SaaS 型防护域名。 |
| [tencentcloud_waf_instance](tables/tencentcloud_waf_instance.md) | WAF 实例。 |
| [tencentcloud_waf_ip_access_control](tables/tencentcloud_waf_ip_access_control.md) | WAF IP 白名单/黑名单条目。 |
| [tencentcloud_waf_object](tables/tencentcloud_waf_object.md) | WAF 防护对象。 |
