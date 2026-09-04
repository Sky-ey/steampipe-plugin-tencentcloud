# 腾讯云 Steampipe 插件

> 使用 SQL 查询腾讯云资源的 Steampipe 插件。

[English](README.md) | 简体中文

## 概述

使用 SQL 查询你腾讯云账号下的 CVM 实例、CBS 云盘、VPC 网络、子网等资源。

## 快速开始

### 安装

从源码构建并安装：

```bash
git clone https://github.com/TencentCloud/steampipe-plugin-tencentcloud.git
cd steampipe-plugin-tencentcloud
make install
```

### 配置

创建或编辑 `~/.steampipe/config/tencentcloud.spc`：

```hcl
connection "tencentcloud" {
  plugin = "tencentcloud/tencentcloud"

  secret_id  = "your-secret-id"
  secret_key = "your-secret-key"

  # 未指定 region 时要查询的地域列表。支持通配符（`*`、`?`、`[gs]`）。
  # 若省略该项，SQL 中未指定 region 时仅查询默认地域 `ap-singapore`。
  regions = ["ap-*"]
}
```

也可以使用环境变量：

```bash
export TENCENTCLOUD_SECRET_ID=your-secret-id
export TENCENTCLOUD_SECRET_KEY=your-secret-key
```

区域表会查询 `regions` 匹配到的所有地域。可以通过 `where region = 'ap-guangzhou'` 这类 SQL 条件直接定位单个地域，跳过地域发现流程。完整匹配规则见 [docs/index.md](docs/index.md#multi-region-queries)。

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

## 测试

运行快速 Go 单元测试：

```bash
make test
```

针对本地 Mock 腾讯云 API 运行 SQL 查询契约测试套件。该套件覆盖每张表的 List/Get 查询和一个跨服务 join，并校验 API action、路径、请求体、地域以及请求次数：

```bash
make test-query
```

契约测试套件使用隔离的临时 Steampipe 安装目录，不会写入 `~/.steampipe`。

## 数据表

插件目前覆盖 13 个腾讯云服务，共 **94 张表**：

| 表名 | 描述 |
|------|------|
| `tencentcloud_cam_access_key` | 子账号的 CAM 访问密钥。 |
| `tencentcloud_cam_group` | CAM 用户组。 |
| `tencentcloud_cam_policy` | CAM 策略（预设策略与自定义策略）。 |
| `tencentcloud_cam_policy_attachment` | CAM 策略绑定关系（策略与主体的映射）。 |
| `tencentcloud_cam_role` | CAM 角色。 |
| `tencentcloud_cam_saml_provider` | CAM SAML 身份提供商（用于企业 SSO）。 |
| `tencentcloud_cam_user` | CAM 子账号（IAM 用户）。 |
| `tencentcloud_ccn` | CCN（云联网）实例。 |
| `tencentcloud_ccn_instance` | 关联到云联网的实例（VPC、专线网关等）。 |
| `tencentcloud_ccn_route` | 云联网路由条目。 |
| `tencentcloud_cbs_auto_snapshot_policy` | CBS 自动快照策略。 |
| `tencentcloud_cbs_disk` | CBS 云硬盘。 |
| `tencentcloud_cbs_disk_backup` | CBS 云硬盘备份点。 |
| `tencentcloud_cbs_disk_storage_pool` | CBS 专属集群的云硬盘存储池。 |
| `tencentcloud_cbs_region` | CBS 可用地域。 |
| `tencentcloud_cbs_snapshot` | CBS 快照。 |
| `tencentcloud_cbs_snapshot_group` | CBS 快照组。 |
| `tencentcloud_clb_block_ip` | CLB 封禁 IP（黑名单）条目。 |
| `tencentcloud_clb_customized_config` | CLB 自定义配置。 |
| `tencentcloud_clb_cross_target` | CLB 跨域后端目标（CCN 跨域 2.0）。 |
| `tencentcloud_clb_instance` | CLB 负载均衡实例。 |
| `tencentcloud_clb_listener` | CLB 监听器。 |
| `tencentcloud_clb_region` | CLB 可用地域。 |
| `tencentcloud_clb_rewrite` | CLB URL 重写重定向规则。 |
| `tencentcloud_clb_target` | 监听器绑定的 CLB 后端目标。 |
| `tencentcloud_clb_target_group` | CLB 目标组。 |
| `tencentcloud_clb_target_group_instance` | 目标组中注册的 CLB 后端服务器。 |
| `tencentcloud_cls_alarm` | CLS 告警策略。 |
| `tencentcloud_cls_alarm_notice` | CLS 告警通知渠道。 |
| `tencentcloud_cls_config` | CLS 采集配置。 |
| `tencentcloud_cls_logset` | CLS 日志集。 |
| `tencentcloud_cls_machine_group` | CLS 机器组。 |
| `tencentcloud_cls_region` | CLS 可用地域。 |
| `tencentcloud_cls_topic` | CLS 日志主题。 |
| `tencentcloud_cos_bucket` | 当前账号下的 COS 存储桶。 |
| `tencentcloud_cos_object` | COS 存储桶中的对象。 |
| `tencentcloud_cos_region` | COS 可用地域。 |
| `tencentcloud_cvm_auto_scaling_group` | 腾讯云弹性伸缩（AS）组。 |
| `tencentcloud_cvm_disaster_recover_group` | CVM 分散置放群组。 |
| `tencentcloud_cvm_host` | 腾讯云专用宿主机（CDH）实例。 |
| `tencentcloud_cvm_image` | CVM 镜像（公共镜像、自定义镜像与共享镜像）。 |
| `tencentcloud_cvm_instance` | 腾讯云云服务器（CVM）实例。 |
| `tencentcloud_cvm_metric_cpu_daily` | CVM 实例近 30 天的每日 CPU 利用率指标。 |
| `tencentcloud_cvm_metric_cpu_hourly` | CVM 实例近 24 小时的每小时 CPU 利用率指标。 |
| `tencentcloud_cvm_metric_disk_daily` | CVM 实例近 30 天的每日磁盘利用率指标。 |
| `tencentcloud_cvm_metric_disk_hourly` | CVM 实例近 24 小时的每小时磁盘利用率指标。 |
| `tencentcloud_cvm_metric_lan_in_daily` | CVM 实例近 30 天的每日内网入带宽指标。 |
| `tencentcloud_cvm_metric_lan_in_hourly` | CVM 实例近 24 小时的每小时内网入带宽指标。 |
| `tencentcloud_cvm_metric_lan_out_daily` | CVM 实例近 30 天的每日内网出带宽指标。 |
| `tencentcloud_cvm_metric_lan_out_hourly` | CVM 实例近 24 小时的每小时内网出带宽指标。 |
| `tencentcloud_cvm_metric_mem_daily` | CVM 实例近 30 天的每日内存利用率指标。 |
| `tencentcloud_cvm_metric_mem_hourly` | CVM 实例近 24 小时的每小时内存利用率指标。 |
| `tencentcloud_cvm_metric_wan_in_daily` | CVM 实例近 30 天的每日公网入带宽指标。 |
| `tencentcloud_cvm_metric_wan_in_hourly` | CVM 实例近 24 小时的每小时公网入带宽指标。 |
| `tencentcloud_cvm_metric_wan_out_daily` | CVM 实例近 30 天的每日公网出带宽指标。 |
| `tencentcloud_cvm_metric_wan_out_hourly` | CVM 实例近 24 小时的每小时公网出带宽指标。 |
| `tencentcloud_cvm_key_pair` | 可绑定到实例的 SSH 密钥对。 |
| `tencentcloud_cvm_launch_template` | 实例启动模板。 |
| `tencentcloud_cvm_region` | CVM 可用地域。 |
| `tencentcloud_cvm_zone` | CVM 可用区。 |
| `tencentcloud_dc_direct_connect` | DC 物理专线。 |
| `tencentcloud_dc_direct_connect_tunnel` | DC 专用通道。 |
| `tencentcloud_dc_gateway` | DC 专线网关。 |
| `tencentcloud_dc_gateway_ccn_route` | 专线网关的云联网路由。 |
| `tencentcloud_dc_internet_address` | DC 分配的互联网公网地址。 |
| `tencentcloud_dc_region` | DC 可用地域。 |
| `tencentcloud_ssl_certificate` | SSL 证书。 |
| `tencentcloud_tke_cluster` | TKE 托管与独立 Kubernetes 集群。 |
| `tencentcloud_tke_cluster_instance` | TKE 集群管理的 CVM 实例。 |
| `tencentcloud_tke_region` | TKE 可用地域。 |
| `tencentcloud_vpc` | VPC 实例。 |
| `tencentcloud_vpc_bandwidth_package` | VPC 带宽包。 |
| `tencentcloud_vpc_eip` | VPC 弹性公网 IP（EIP）。 |
| `tencentcloud_vpc_eni` | VPC 弹性网卡（ENI）。 |
| `tencentcloud_vpc_nat_gateway` | VPC NAT 网关。 |
| `tencentcloud_vpc_region` | VPC 可用地域。 |
| `tencentcloud_vpc_route` | 从路由表展开的 VPC 路由条目。 |
| `tencentcloud_vpc_route_table` | VPC 路由表。 |
| `tencentcloud_vpc_security_group` | VPC 安全组。 |
| `tencentcloud_vpc_security_group_policy` | VPC 安全组规则（入站与出站）。 |
| `tencentcloud_vpc_subnet` | VPC 子网。 |
| `tencentcloud_vpc_traffic_package` | VPC 流量包。 |
| `tencentcloud_vpn_gateway` | VPN 网关。 |
| `tencentcloud_vpn_customer_gateway` | VPN 对端网关。 |
| `tencentcloud_vpn_connection` | VPN 通道（IPsec 隧道）。 |
| `tencentcloud_vpn_gateway_route` | VPN 网关的路由条目。 |
| `tencentcloud_vpn_gateway_ccn_route` | VPN 网关的云联网路由。 |
| `tencentcloud_vpn_region` | VPN 可用地域。 |
| `tencentcloud_waf_clb_host` | WAF CLB 型防护域名。 |
| `tencentcloud_waf_custom_rule` | WAF 自定义访问控制规则。 |
| `tencentcloud_waf_domain` | WAF SaaS 型防护域名。 |
| `tencentcloud_waf_instance` | WAF 实例。 |
| `tencentcloud_waf_ip_access_control` | WAF IP 白名单/黑名单条目。 |
| `tencentcloud_waf_object` | WAF 防护对象。 |
