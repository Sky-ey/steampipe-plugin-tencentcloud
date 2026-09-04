package tencentcloud

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/modules/cam"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/modules/cbs"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/modules/ccn"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/modules/clb"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/modules/cls"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/modules/cos"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/modules/cvm"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/modules/dc"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/modules/ssl"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/modules/tke"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/modules/vpc"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/modules/vpn"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/modules/waf"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func Plugin(_ context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-tencentcloud",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: utils.ConfigInstance,
		},
		ConnectionKeyColumns: []plugin.ConnectionKeyColumn{
			{
				Name:    "owner_uin",
				Hydrate: utils.OwnerUinForConnection,
			},
		},
		DefaultTransform: transform.FromGo().NullIfZero(),
		DefaultIgnoreConfig: &plugin.IgnoreConfig{
			ShouldIgnoreErrorFunc: isNotFoundError,
		},
		DefaultRetryConfig: &plugin.RetryConfig{
			ShouldRetryErrorFunc: shouldRetryError,
		},
		TableMap: map[string]*plugin.Table{
			// CAM
			"tencentcloud_cam_access_key":        cam.TableTencentcloudCamAccessKey(),
			"tencentcloud_cam_group":             cam.TableTencentcloudCamGroup(),
			"tencentcloud_cam_policy":            cam.TableTencentcloudCamPolicy(),
			"tencentcloud_cam_policy_attachment": cam.TableTencentcloudCamPolicyAttachment(),
			"tencentcloud_cam_role":              cam.TableTencentcloudCamRole(),
			"tencentcloud_cam_saml_provider":     cam.TableTencentcloudCamSamlProvider(),
			"tencentcloud_cam_user":              cam.TableTencentcloudCamUser(),
			// CVM
			"tencentcloud_cvm_instance":               cvm.TableTencentcloudCvmInstance(),
			"tencentcloud_cvm_host":                   cvm.TableTencentcloudCvmHost(),
			"tencentcloud_cvm_auto_scaling_group":     cvm.TableTencentcloudCvmAutoScalingGroup(),
			"tencentcloud_cvm_disaster_recover_group": cvm.TableTencentcloudCvmDisasterRecoverGroup(),
			"tencentcloud_cvm_image":                  cvm.TableTencentcloudCvmImage(),
			"tencentcloud_cvm_metric_cpu_daily":       cvm.TableTencentcloudCvmMetricCpuDaily(),
			"tencentcloud_cvm_metric_cpu_hourly":      cvm.TableTencentcloudCvmMetricCpuHourly(),
			"tencentcloud_cvm_metric_disk_daily":      cvm.TableTencentcloudCvmMetricDiskDaily(),
			"tencentcloud_cvm_metric_disk_hourly":     cvm.TableTencentcloudCvmMetricDiskHourly(),
			"tencentcloud_cvm_metric_lan_in_daily":    cvm.TableTencentcloudCvmMetricLanInDaily(),
			"tencentcloud_cvm_metric_lan_in_hourly":   cvm.TableTencentcloudCvmMetricLanInHourly(),
			"tencentcloud_cvm_metric_lan_out_daily":   cvm.TableTencentcloudCvmMetricLanOutDaily(),
			"tencentcloud_cvm_metric_lan_out_hourly":  cvm.TableTencentcloudCvmMetricLanOutHourly(),
			"tencentcloud_cvm_metric_mem_daily":       cvm.TableTencentcloudCvmMetricMemDaily(),
			"tencentcloud_cvm_metric_mem_hourly":      cvm.TableTencentcloudCvmMetricMemHourly(),
			"tencentcloud_cvm_metric_wan_in_daily":    cvm.TableTencentcloudCvmMetricWanInDaily(),
			"tencentcloud_cvm_metric_wan_in_hourly":   cvm.TableTencentcloudCvmMetricWanInHourly(),
			"tencentcloud_cvm_metric_wan_out_daily":   cvm.TableTencentcloudCvmMetricWanOutDaily(),
			"tencentcloud_cvm_metric_wan_out_hourly":  cvm.TableTencentcloudCvmMetricWanOutHourly(),
			"tencentcloud_cvm_key_pair":               cvm.TableTencentcloudCvmKeyPair(),
			"tencentcloud_cvm_launch_template":        cvm.TableTencentcloudCvmLaunchTemplate(),
			"tencentcloud_cvm_region":                 cvm.TableTencentcloudCvmRegion(),
			"tencentcloud_cvm_zone":                   cvm.TableTencentcloudCvmZone(),
			// CBS
			"tencentcloud_cbs_auto_snapshot_policy": cbs.TableTencentcloudCbsAutoSnapshotPolicy(),
			"tencentcloud_cbs_disk":                 cbs.TableTencentcloudCbsDisk(),
			"tencentcloud_cbs_disk_backup":          cbs.TableTencentcloudCbsDiskBackup(),
			"tencentcloud_cbs_disk_storage_pool":    cbs.TableTencentcloudCbsDiskStoragePool(),
			"tencentcloud_cbs_snapshot":             cbs.TableTencentcloudCbsSnapshot(),
			"tencentcloud_cbs_snapshot_group":       cbs.TableTencentcloudCbsSnapshotGroup(),
			"tencentcloud_cbs_region":               cbs.TableTencentcloudCbsRegion(),
			// CLB
			"tencentcloud_clb_block_ip":              clb.TableTencentcloudClbBlockIP(),
			"tencentcloud_clb_customized_config":     clb.TableTencentcloudClbCustomizedConfig(),
			"tencentcloud_clb_cross_target":          clb.TableTencentcloudClbCrossTarget(),
			"tencentcloud_clb_instance":              clb.TableTencentcloudClbLoadBalancer(),
			"tencentcloud_clb_listener":              clb.TableTencentcloudClbListener(),
			"tencentcloud_clb_rewrite":               clb.TableTencentcloudClbRewrite(),
			"tencentcloud_clb_target":                clb.TableTencentcloudClbTarget(),
			"tencentcloud_clb_target_group":          clb.TableTencentcloudClbTargetGroup(),
			"tencentcloud_clb_target_group_instance": clb.TableTencentcloudClbTargetGroupInstance(),
			"tencentcloud_clb_region":                clb.TableTencentcloudClbRegion(),
			// VPC
			"tencentcloud_vpc":                       vpc.TableTencentcloudVpc(),
			"tencentcloud_vpc_bandwidth_package":     vpc.TableTencentcloudVpcBandwidthPackage(),
			"tencentcloud_vpc_traffic_package":       vpc.TableTencentcloudVpcTrafficPackage(),
			"tencentcloud_vpc_subnet":                vpc.TableTencentcloudVpcSubnet(),
			"tencentcloud_vpc_route_table":           vpc.TableTencentcloudVpcRouteTable(),
			"tencentcloud_vpc_route":                 vpc.TableTencentcloudVpcRoute(),
			"tencentcloud_vpc_security_group":        vpc.TableTencentcloudVpcSecurityGroup(),
			"tencentcloud_vpc_security_group_policy": vpc.TableTencentcloudVpcSecurityGroupPolicy(),
			"tencentcloud_vpc_eni":                   vpc.TableTencentcloudVpcEni(),
			"tencentcloud_vpc_eip":                   vpc.TableTencentcloudVpcEip(),
			"tencentcloud_vpc_nat_gateway":           vpc.TableTencentcloudVpcNatGateway(),
			"tencentcloud_vpc_region":                vpc.TableTencentcloudVpcRegion(),
			// DC
			"tencentcloud_dc_direct_connect":        dc.TableTencentcloudDcDirectConnect(),
			"tencentcloud_dc_direct_connect_tunnel": dc.TableTencentcloudDcDirectConnectTunnel(),
			"tencentcloud_dc_gateway":               dc.TableTencentcloudDcGateway(),
			"tencentcloud_dc_gateway_ccn_route":     dc.TableTencentcloudDcGatewayCcnRoute(),
			"tencentcloud_dc_internet_address":      dc.TableTencentcloudDcInternetAddress(),
			"tencentcloud_dc_region":                dc.TableTencentcloudDcRegion(),
			// VPN
			"tencentcloud_vpn_gateway":           vpn.TableTencentcloudVpnGateway(),
			"tencentcloud_vpn_customer_gateway":  vpn.TableTencentcloudVpnCustomerGateway(),
			"tencentcloud_vpn_connection":        vpn.TableTencentcloudVpnConnection(),
			"tencentcloud_vpn_gateway_route":     vpn.TableTencentcloudVpnGatewayRoute(),
			"tencentcloud_vpn_gateway_ccn_route": vpn.TableTencentcloudVpnGatewayCcnRoute(),
			"tencentcloud_vpn_region":            vpn.TableTencentcloudVpnRegion(),
			// CCN
			"tencentcloud_ccn":          ccn.TableTencentcloudCcn(),
			"tencentcloud_ccn_instance": ccn.TableTencentcloudCcnInstance(),
			"tencentcloud_ccn_route":    ccn.TableTencentcloudCcnRoute(),
			// SSL
			"tencentcloud_ssl_certificate": ssl.TableTencentcloudSslCertificate(),
			// TKE
			"tencentcloud_tke_cluster":          tke.TableTencentcloudTkeCluster(),
			"tencentcloud_tke_cluster_instance": tke.TableTencentcloudTkeClusterInstance(),
			"tencentcloud_tke_region":           tke.TableTencentcloudTkeRegion(),
			// COS
			"tencentcloud_cos_bucket": cos.TableTencentcloudCosBucket(),
			"tencentcloud_cos_object": cos.TableTencentcloudCosObject(),
			"tencentcloud_cos_region": cos.TableTencentcloudCosRegion(),
			// WAF
			"tencentcloud_waf_instance":          waf.TableTencentcloudWafInstance(),
			"tencentcloud_waf_domain":            waf.TableTencentcloudWafDomain(),
			"tencentcloud_waf_clb_host":          waf.TableTencentcloudWafClbHost(),
			"tencentcloud_waf_object":            waf.TableTencentcloudWafObject(),
			"tencentcloud_waf_ip_access_control": waf.TableTencentcloudWafIpAccessControl(),
			"tencentcloud_waf_custom_rule":       waf.TableTencentcloudWafCustomRule(),
			// CLS
			"tencentcloud_cls_logset":        cls.TableTencentcloudClsLogset(),
			"tencentcloud_cls_topic":         cls.TableTencentcloudClsTopic(),
			"tencentcloud_cls_machine_group": cls.TableTencentcloudClsMachineGroup(),
			"tencentcloud_cls_config":        cls.TableTencentcloudClsConfig(),
			"tencentcloud_cls_alarm":         cls.TableTencentcloudClsAlarm(),
			"tencentcloud_cls_alarm_notice":  cls.TableTencentcloudClsAlarmNotice(),
			"tencentcloud_cls_region":        cls.TableTencentcloudClsRegion(),
		},
	}
	return p
}
