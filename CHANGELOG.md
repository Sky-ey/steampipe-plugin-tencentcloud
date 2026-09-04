# Changelog

## v0.0.1

Initial release — **94 tables across 13 Tencent Cloud services**.

### Tables

CAM (Access Management) 7, CCN (Cloud Connect Network) 3, CBS (Cloud Block Storage) 7, CLB (Cloud Load Balancer) 10, CLS (Cloud Log Service) 7, COS (Cloud Object Storage) 3, CVM (Cloud Virtual Machine) + Auto Scaling + Monitor 23, DC (Direct Connect) 6, SSL Certificates 1, TKE (Tencent Kubernetes Engine) 3, VPC (Virtual Private Cloud) 12, VPN Gateway 6, WAF (Web Application Firewall) 6

#### CAM (Access Management) (7 tables)

- `tencentcloud_cam_access_key` — CAM access keys of sub-accounts.
- `tencentcloud_cam_group` — CAM user groups.
- `tencentcloud_cam_policy` — CAM policies (preset and custom).
- `tencentcloud_cam_policy_attachment` — CAM policy attachments (policy-to-principal mapping).
- `tencentcloud_cam_role` — CAM roles.
- `tencentcloud_cam_saml_provider` — CAM SAML identity providers for enterprise SSO.
- `tencentcloud_cam_user` — CAM sub-accounts (IAM users).

#### CCN (Cloud Connect Network) (3 tables)

- `tencentcloud_ccn` — CCN (Cloud Connect Network) instances.
- `tencentcloud_ccn_instance` — Instances (VPCs, DC gateways) attached to a CCN.
- `tencentcloud_ccn_route` — Route entries of a CCN.

#### CBS (Cloud Block Storage) (7 tables)

- `tencentcloud_cbs_auto_snapshot_policy` — CBS auto snapshot policies.
- `tencentcloud_cbs_disk` — CBS cloud disks.
- `tencentcloud_cbs_disk_backup` — CBS disk backup points.
- `tencentcloud_cbs_disk_storage_pool` — CBS dedicated cluster disk storage pools.
- `tencentcloud_cbs_region` — Regions available to CBS.
- `tencentcloud_cbs_snapshot` — CBS snapshots.
- `tencentcloud_cbs_snapshot_group` — CBS snapshot groups.

#### CLB (Cloud Load Balancer) (10 tables)

- `tencentcloud_clb_block_ip` — CLB blocked IP (blacklist) entries.
- `tencentcloud_clb_customized_config` — CLB customized configurations.
- `tencentcloud_clb_cross_target` — CLB cross-VPC backend targets (CCN cross-domain 2.0).
- `tencentcloud_clb_instance` — CLB load balancer instances.
- `tencentcloud_clb_listener` — CLB listeners.
- `tencentcloud_clb_region` — Regions available to CLB.
- `tencentcloud_clb_rewrite` — CLB URL rewrite redirection rules.
- `tencentcloud_clb_target` — CLB backend targets bound to listeners.
- `tencentcloud_clb_target_group` — CLB target groups.
- `tencentcloud_clb_target_group_instance` — CLB backend servers registered in target groups.

#### CLS (Cloud Log Service) (7 tables)

- `tencentcloud_cls_alarm` — CLS alarm policies.
- `tencentcloud_cls_alarm_notice` — CLS alarm notification channels.
- `tencentcloud_cls_config` — CLS collection configurations.
- `tencentcloud_cls_logset` — CLS logsets.
- `tencentcloud_cls_machine_group` — CLS machine groups.
- `tencentcloud_cls_region` — Regions available to CLS.
- `tencentcloud_cls_topic` — CLS log topics.

#### COS (Cloud Object Storage) (3 tables)

- `tencentcloud_cos_bucket` — COS buckets under the current account.
- `tencentcloud_cos_object` — Objects stored in COS buckets.
- `tencentcloud_cos_region` — Regions available to COS.

#### CVM (Cloud Virtual Machine) + Auto Scaling + Monitor (23 tables)

- `tencentcloud_cvm_auto_scaling_group` — Tencent Cloud Auto Scaling (AS) groups.
- `tencentcloud_cvm_disaster_recover_group` — CVM spread placement groups.
- `tencentcloud_cvm_host` — Tencent Cloud Cloud Dedicated Host (CDH) instances.
- `tencentcloud_cvm_image` — CVM images (public, custom and shared).
- `tencentcloud_cvm_instance` — Tencent Cloud Virtual Machine (CVM) instances.
- `tencentcloud_cvm_metric_cpu_daily` — Daily CPU utilization metric of CVM instances over the last 30 days.
- `tencentcloud_cvm_metric_cpu_hourly` — Hourly CPU utilization metric of CVM instances over the last 24 hours.
- `tencentcloud_cvm_metric_disk_daily` — Daily disk utilization metric of CVM instances over the last 30 days.
- `tencentcloud_cvm_metric_disk_hourly` — Hourly disk utilization metric of CVM instances over the last 24 hours.
- `tencentcloud_cvm_metric_lan_in_daily` — Daily private network inbound bandwidth metric of CVM instances over the last 30 days.
- `tencentcloud_cvm_metric_lan_in_hourly` — Hourly private network inbound bandwidth metric of CVM instances over the last 24 hours.
- `tencentcloud_cvm_metric_lan_out_daily` — Daily private network outbound bandwidth metric of CVM instances over the last 30 days.
- `tencentcloud_cvm_metric_lan_out_hourly` — Hourly private network outbound bandwidth metric of CVM instances over the last 24 hours.
- `tencentcloud_cvm_metric_mem_daily` — Daily memory utilization metric of CVM instances over the last 30 days.
- `tencentcloud_cvm_metric_mem_hourly` — Hourly memory utilization metric of CVM instances over the last 24 hours.
- `tencentcloud_cvm_metric_wan_in_daily` — Daily public network inbound bandwidth metric of CVM instances over the last 30 days.
- `tencentcloud_cvm_metric_wan_in_hourly` — Hourly public network inbound bandwidth metric of CVM instances over the last 24 hours.
- `tencentcloud_cvm_metric_wan_out_daily` — Daily public network outbound bandwidth metric of CVM instances over the last 30 days.
- `tencentcloud_cvm_metric_wan_out_hourly` — Hourly public network outbound bandwidth metric of CVM instances over the last 24 hours.
- `tencentcloud_cvm_key_pair` — SSH key pairs that can be bound to instances.
- `tencentcloud_cvm_launch_template` — Instance launch templates.
- `tencentcloud_cvm_region` — Regions available to CVM.
- `tencentcloud_cvm_zone` — Availability zones available to CVM.

#### DC (Direct Connect) (6 tables)

- `tencentcloud_dc_direct_connect` — DC physical direct connect lines.
- `tencentcloud_dc_direct_connect_tunnel` — DC dedicated tunnels.
- `tencentcloud_dc_gateway` — DC gateways.
- `tencentcloud_dc_gateway_ccn_route` — CCN routes of a DC gateway.
- `tencentcloud_dc_internet_address` — Internet public addresses allocated by DC.
- `tencentcloud_dc_region` — Regions available to DC.

#### SSL Certificates (1 tables)

- `tencentcloud_ssl_certificate` — SSL certificates.

#### TKE (Tencent Kubernetes Engine) (3 tables)

- `tencentcloud_tke_cluster` — TKE managed and independent Kubernetes clusters.
- `tencentcloud_tke_cluster_instance` — CVM instances managed by TKE clusters.
- `tencentcloud_tke_region` — Regions available to TKE.

#### VPC (Virtual Private Cloud) (12 tables)

- `tencentcloud_vpc` — VPC instances.
- `tencentcloud_vpc_bandwidth_package` — VPC bandwidth packages.
- `tencentcloud_vpc_eip` — VPC elastic public IP (EIP) addresses.
- `tencentcloud_vpc_eni` — VPC elastic network interfaces (ENI).
- `tencentcloud_vpc_nat_gateway` — VPC NAT gateways.
- `tencentcloud_vpc_region` — Regions available to VPC.
- `tencentcloud_vpc_route` — VPC routing entries expanded from route tables.
- `tencentcloud_vpc_route_table` — VPC route tables.
- `tencentcloud_vpc_security_group` — VPC security groups.
- `tencentcloud_vpc_security_group_policy` — VPC security group rules (ingress and egress).
- `tencentcloud_vpc_subnet` — VPC subnets.
- `tencentcloud_vpc_traffic_package` — VPC traffic packages.

#### VPN Gateway (6 tables)

- `tencentcloud_vpn_gateway` — VPN gateways.
- `tencentcloud_vpn_customer_gateway` — VPN customer gateways.
- `tencentcloud_vpn_connection` — VPN connections (IPsec tunnels).
- `tencentcloud_vpn_gateway_route` — Routes of a VPN gateway.
- `tencentcloud_vpn_gateway_ccn_route` — CCN routes of a VPN gateway.
- `tencentcloud_vpn_region` — Regions available to VPN.

#### WAF (Web Application Firewall) (6 tables)

- `tencentcloud_waf_clb_host` — WAF CLB-type protected domains.
- `tencentcloud_waf_custom_rule` — WAF custom access control rules.
- `tencentcloud_waf_domain` — WAF SaaS-type protected domains.
- `tencentcloud_waf_instance` — WAF instances.
- `tencentcloud_waf_ip_access_control` — WAF IP allowlist/blocklist entries.
- `tencentcloud_waf_object` — WAF protection objects.

### Features

- **Multi-region queries** — the `regions` config argument accepts exact region IDs and glob patterns (`*`, `?`, `[gs]`); matching regions are queried concurrently and results merged. A `region` filter in SQL targets a single region directly.
- **Multi-account queries** — every table exposes `owner_uin` / `caller_uin` / `owner_app_id` account columns; `owner_uin` is registered as the connection key column so aggregator connections prune child connections on account filters.
- **Credential chain** — connection config > environment variables (`TENCENTCLOUD_SECRET_ID` / `TENCENTCLOUD_SECRET_KEY` / `TENCENTCLOUD_TOKEN`) > SDK default provider chain (shared credentials file `~/.tencentcloud/credentials`, CVM instance role). STS temporary credentials supported via `token`.
- **Network customization** — `endpoint` / `base_url` routing to custom or internal gateways, per-host `dns_override`, and `insecure_skip_verify` for test environments.
- **Reliability controls** — configurable `timeout`, automatic retry with exponential backoff (`auto_retry`, `max_retry_time`), plus `ignore_error_codes` / `retry_error_codes` overrides.
- **CVM monitoring metrics** — 7 metrics (CPU, memory, disk, LAN/WAN in/out) × 2 granularities (daily 30 days, hourly 24 hours) via the `GetMonitorData` API.
