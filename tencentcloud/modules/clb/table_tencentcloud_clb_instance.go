package clb

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"time"

	clb "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/clb/v20180317"
	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudClbLoadBalancer() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_clb_instance",
		Description:       "Tencent Cloud Load Balancer (CLB) instances.",
		GetMatrixItemFunc: buildClbRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "load_balancer_id", Require: plugin.Optional},
				{Name: "load_balancer_name", Require: plugin.Optional},
				{Name: "load_balancer_type", Require: plugin.Optional},
				{Name: "vpc_id", Require: plugin.Optional},
				{Name: "security_group", Require: plugin.Optional},
				{Name: "project_id", Require: plugin.Optional},
			},
			Hydrate: listClbLoadBalancers,
			Tags:    map[string]string{"service": "clb", "action": "DescribeLoadBalancers"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("load_balancer_id"),
			Hydrate:    getClbLoadBalancer,
			Tags:       map[string]string{"service": "clb", "action": "DescribeLoadBalancers"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "load_balancer_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the load balancer instance.", Transform: transform.FromField("LoadBalancerId")},
			{Name: "load_balancer_name", Type: proto.ColumnType_STRING, Description: "The name of the load balancer instance."},
			{Name: "load_balancer_type", Type: proto.ColumnType_STRING, Description: "The network type of the load balancer: OPEN (public) or INTERNAL (private)."},
			{Name: "load_balancer_vips", Type: proto.ColumnType_JSON, Description: "The VIP addresses of the load balancer.", Transform: transform.FromField("LoadBalancerVips")},
			{Name: "domain", Type: proto.ColumnType_STRING, Description: "The domain name assigned to the load balancer."},
			{Name: "load_balancer_domain", Type: proto.ColumnType_STRING, Description: "The custom domain name bound to the load balancer."},
			{Name: "status", Type: proto.ColumnType_INT, Description: "The status of the load balancer: 0 (creating), 1 (running)."},
			{Name: "forward", Type: proto.ColumnType_INT, Description: "The load balancer instance type: 1 (general), 0 (classic)."},
			{Name: "create_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the load balancer was created.", Transform: transform.FromField("CreateTime").NullIfZero()},
			{Name: "status_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the load balancer status last changed.", Transform: transform.FromField("StatusTime").NullIfZero()},
			{Name: "project_id", Type: proto.ColumnType_INT, Description: "The project ID of the load balancer.", Transform: transform.FromField("ProjectId")},
			{Name: "vpc_id", Type: proto.ColumnType_STRING, Description: "The ID of the VPC to which the load balancer belongs.", Transform: transform.FromField("VpcId")},
			{Name: "numerical_vpc_id", Type: proto.ColumnType_INT, Description: "The numeric ID of the VPC.", Transform: transform.FromField("NumericalVpcId")},
			{Name: "subnet_id", Type: proto.ColumnType_STRING, Description: "The ID of the subnet to which the load balancer belongs.", Transform: transform.FromField("SubnetId")},
			{Name: "address_ip_version", Type: proto.ColumnType_STRING, Description: "The IP version: ipv4, ipv6, or IPv6FullChain.", Transform: transform.FromField("AddressIPVersion")},
			{Name: "address_ipv6", Type: proto.ColumnType_STRING, Description: "The IPv6 address of the load balancer.", Transform: transform.FromField("AddressIPv6")},
			{Name: "ip_v6_mode", Type: proto.ColumnType_STRING, Description: "The IPv6 mode of the load balancer.", Transform: transform.FromField("IPv6Mode")},
			{Name: "vip_isp", Type: proto.ColumnType_STRING, Description: "The ISP type of the VIP, such as BGP, CMCC, CTCC or CUCC.", Transform: transform.FromField("VipIsp")},
			{Name: "open_bgp", Type: proto.ColumnType_INT, Description: "Whether BGP is enabled: 1 (enabled), 0 (disabled).", Transform: transform.FromField("OpenBgp")},
			{Name: "local_bgp", Type: proto.ColumnType_BOOL, Description: "Whether the load balancer uses local BGP."},
			{Name: "anycast_zone", Type: proto.ColumnType_STRING, Description: "The anycast zone of the load balancer."},
			{Name: "egress", Type: proto.ColumnType_STRING, Description: "The network egress of the load balancer."},
			{Name: "charge_type", Type: proto.ColumnType_STRING, Description: "The billing mode of the load balancer: PREPAID or POSTPAID_BY_HOUR."},
			{Name: "sla_type", Type: proto.ColumnType_STRING, Description: "The performance capacity specification, such as clb.c2.medium.", Transform: transform.FromField("SlaType")},
			{Name: "isolation", Type: proto.ColumnType_INT, Description: "Whether the load balancer is isolated: 0 (no), 1 (yes)."},
			{Name: "isolated_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the load balancer was isolated.", Transform: transform.FromField("IsolatedTime").NullIfZero()},
			{Name: "expire_time", Type: proto.ColumnType_TIMESTAMP, Description: "The expiration time of the prepaid load balancer.", Transform: transform.FromField("ExpireTime").NullIfZero()},
			{Name: "is_block", Type: proto.ColumnType_BOOL, Description: "Whether the load balancer is blocked."},
			{Name: "is_block_time", Type: proto.ColumnType_STRING, Description: "The time when the load balancer was blocked."},
			{Name: "is_ddos", Type: proto.ColumnType_BOOL, Description: "Whether the load balancer is protected by Anti-DDoS Pro.", Transform: transform.FromField("IsDDos")},
			{Name: "snat", Type: proto.ColumnType_BOOL, Description: "Whether SNAT is enabled for the classic load balancer."},
			{Name: "snat_pro", Type: proto.ColumnType_BOOL, Description: "Whether SNAT Pro is enabled."},
			{Name: "snat_ips", Type: proto.ColumnType_JSON, Description: "The SNAT IPs of the load balancer.", Transform: transform.FromField("SnatIps")},
			{Name: "mix_ip_target", Type: proto.ColumnType_BOOL, Description: "Whether hybrid IP target binding is enabled."},
			{Name: "exclusive", Type: proto.ColumnType_INT, Description: "Whether the load balancer is exclusive: 1 (exclusive), 0 (shared)."},
			{Name: "log", Type: proto.ColumnType_STRING, Description: "The log information of the load balancer."},
			{Name: "log_set_id", Type: proto.ColumnType_STRING, Description: "The CLS log set ID for access logs.", Transform: transform.FromField("LogSetId")},
			{Name: "log_topic_id", Type: proto.ColumnType_STRING, Description: "The CLS log topic ID for access logs.", Transform: transform.FromField("LogTopicId")},
			{Name: "health_log_set_id", Type: proto.ColumnType_STRING, Description: "The CLS log set ID for health check logs.", Transform: transform.FromField("HealthLogSetId")},
			{Name: "health_log_topic_id", Type: proto.ColumnType_STRING, Description: "The CLS log topic ID for health check logs.", Transform: transform.FromField("HealthLogTopicId")},
			{Name: "config_id", Type: proto.ColumnType_STRING, Description: "The ID of the personalized configuration bound to the load balancer.", Transform: transform.FromField("ConfigId")},
			{Name: "load_balancer_pass_to_target", Type: proto.ColumnType_BOOL, Description: "Whether the load balancer passes traffic to targets without SNAT."},
			{Name: "target_count", Type: proto.ColumnType_INT, Description: "The number of backend services bound to the load balancer."},
			{Name: "associate_endpoint", Type: proto.ColumnType_STRING, Description: "The endpoint associated with the load balancer."},
			{Name: "cluster_tag", Type: proto.ColumnType_STRING, Description: "The tag of the dedicated cluster to which the load balancer belongs."},
			{Name: "cluster_ids", Type: proto.ColumnType_JSON, Description: "The IDs of the dedicated clusters to which the load balancer belongs.", Transform: transform.FromField("ClusterIds")},
			{Name: "nfv_info", Type: proto.ColumnType_STRING, Description: "The NFV information of the load balancer.", Transform: transform.FromField("NfvInfo")},
			{Name: "zones", Type: proto.ColumnType_JSON, Description: "The availability zones of the load balancer."},
			{Name: "master_zone", Type: proto.ColumnType_JSON, Description: "The primary availability zone information."},
			{Name: "backup_zone_set", Type: proto.ColumnType_JSON, Description: "The backup availability zone information."},
			{Name: "available_zone_affinity_info", Type: proto.ColumnType_JSON, Description: "The availability zone affinity information."},
			{Name: "target_region_info", Type: proto.ColumnType_JSON, Description: "The region information of the backend services."},
			{Name: "network_attributes", Type: proto.ColumnType_JSON, Description: "The network billing attributes of the load balancer."},
			{Name: "prepaid_attributes", Type: proto.ColumnType_JSON, Description: "The prepaid billing attributes of the load balancer."},
			{Name: "extra_info", Type: proto.ColumnType_JSON, Description: "The extra information of the load balancer."},
			{Name: "exclusive_cluster", Type: proto.ColumnType_JSON, Description: "The exclusive cluster information of the load balancer."},
			{Name: "secure_groups", Type: proto.ColumnType_JSON, Description: "The security group IDs bound to the load balancer."},
			{Name: "attribute_flags", Type: proto.ColumnType_JSON, Description: "The attribute flags of the load balancer."},
			{Name: "security_group", Type: proto.ColumnType_STRING, Description: "The first security group bound to the load balancer, used for filtering."},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The tags attached to the load balancer as a map of key to value."},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the load balancer resides."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource."},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource."},
		}),
	}
}

// List Function

func listClbLoadBalancers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listClbLoadBalancers", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &clb.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listClbLoadBalancers", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)
	pageSize := utils.PageSizeFromLimit(d, 100, 1)
	var offset int64
	additionalFields := []string{"TargetCount"}

	for {
		d.WaitForListRateLimit(ctx)
		req := clb.NewDescribeLoadBalancersRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		req.AdditionalFields = common.StringPtrs(additionalFields)

		if id := d.EqualsQualString("load_balancer_id"); id != "" {
			// LoadBalancerIds 与其他过滤参数互斥,ID 优先
			req.LoadBalancerIds = common.StringPtrs([]string{id})
		} else {
			if name := d.EqualsQualString("load_balancer_name"); name != "" {
				req.LoadBalancerName = &name
			}
			if lbType := d.EqualsQualString("load_balancer_type"); lbType != "" {
				req.LoadBalancerType = &lbType
			}
			if vpcID := d.EqualsQualString("vpc_id"); vpcID != "" {
				req.VpcId = &vpcID
			}
			if sg := d.EqualsQualString("security_group"); sg != "" {
				req.SecurityGroup = &sg
			}
			if qual := d.EqualsQuals["project_id"]; qual != nil {
				projectID := qual.GetInt64Value()
				req.ProjectId = &projectID
			}
		}

		utils.LogRequest(ctx, "listClbLoadBalancers", req)
		resp, err := client.DescribeLoadBalancersWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listClbLoadBalancers", "request_error", err)
			return nil, err
		}

		if resp == nil || resp.Response == nil {
			break
		}
		loadBalancers := resp.Response.LoadBalancerSet
		for _, lb := range loadBalancers {
			if lb == nil {
				continue
			}
			row := toClbLoadBalancerRow(lb, region)
			d.StreamListItem(ctx, row)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if utils.PageDone(offset, int64(len(loadBalancers)), utils.PtrUint64(resp.Response.TotalCount)) {
			break
		}
		offset += int64(len(loadBalancers))
	}
	return nil, nil
}

// Get Function

func getClbLoadBalancer(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	loadBalancerID := d.EqualsQualString("load_balancer_id")
	plugin.Logger(ctx).Debug("getClbLoadBalancer", "load_balancer_id", loadBalancerID, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &clb.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getClbLoadBalancer", "client_init_error", err)
		return nil, err
	}

	if loadBalancerID == "" {
		return nil, nil
	}
	req := clb.NewDescribeLoadBalancersRequest()
	req.LoadBalancerIds = common.StringPtrs([]string{loadBalancerID})
	req.AdditionalFields = common.StringPtrs([]string{"TargetCount"})

	utils.LogRequest(ctx, "getClbLoadBalancer", req)
	resp, err := client.DescribeLoadBalancersWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getClbLoadBalancer", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.LoadBalancerSet) == 0 || resp.Response.LoadBalancerSet[0] == nil {
		return nil, nil
	}
	return toClbLoadBalancerRow(resp.Response.LoadBalancerSet[0], d.EqualsQualString(utils.MatrixKeyRegion)), nil
}

// Row Type

type clbLoadBalancerRow struct {
	LoadBalancerId            string
	LoadBalancerName          string
	LoadBalancerType          string
	LoadBalancerVips          []string
	Domain                    string
	LoadBalancerDomain        string
	Status                    uint64
	Forward                   uint64
	CreateTime                *time.Time
	StatusTime                *time.Time
	ProjectId                 uint64
	VpcId                     string
	NumericalVpcId            uint64
	SubnetId                  string
	AddressIPVersion          string
	AddressIPv6               string
	IPv6Mode                  string
	VipIsp                    string
	OpenBgp                   uint64
	LocalBgp                  bool
	AnycastZone               string
	Egress                    string
	ChargeType                string
	SlaType                   string
	Isolation                 uint64
	IsolatedTime              *time.Time
	ExpireTime                *time.Time
	IsBlock                   bool
	IsBlockTime               *time.Time
	IsDDos                    bool
	Snat                      bool
	SnatPro                   bool
	SnatIps                   []*clb.SnatIp
	MixIpTarget               bool
	Exclusive                 uint64
	Log                       string
	LogSetId                  string
	LogTopicId                string
	HealthLogSetId            string
	HealthLogTopicId          string
	ConfigId                  string
	LoadBalancerPassToTarget  bool
	TargetCount               uint64
	AssociateEndpoint         string
	ClusterTag                string
	ClusterIds                []string
	NfvInfo                   string
	Zones                     []string
	MasterZone                *clb.ZoneInfo
	BackupZoneSet             []*clb.ZoneInfo
	AvailableZoneAffinityInfo *clb.AvailableZoneAffinityInfo
	TargetRegionInfo          *clb.TargetRegionInfo
	NetworkAttributes         *clb.InternetAccessible
	PrepaidAttributes         *clb.LBChargePrepaid
	ExtraInfo                 *clb.ExtraInfo
	ExclusiveCluster          *clb.ExclusiveCluster
	SecureGroups              []string
	AttributeFlags            []string
	SecurityGroup             string
	Tags                      map[string]string
	Region                    string
	Title                     string
	Akas                      []string
}

func toClbLoadBalancerRow(lb *clb.LoadBalancer, region string) clbLoadBalancerRow {
	row := clbLoadBalancerRow{
		LoadBalancerId:            utils.PtrString(lb.LoadBalancerId),
		LoadBalancerName:          utils.PtrString(lb.LoadBalancerName),
		LoadBalancerType:          utils.PtrString(lb.LoadBalancerType),
		LoadBalancerVips:          utils.PtrStringSlice(lb.LoadBalancerVips),
		Domain:                    utils.PtrString(lb.Domain),
		LoadBalancerDomain:        utils.PtrString(lb.LoadBalancerDomain),
		Status:                    utils.PtrUint64(lb.Status),
		Forward:                   utils.PtrUint64(lb.Forward),
		CreateTime:                utils.ParseTimestamp(lb.CreateTime),
		StatusTime:                utils.ParseTimestamp(lb.StatusTime),
		ProjectId:                 utils.PtrUint64(lb.ProjectId),
		VpcId:                     utils.PtrString(lb.VpcId),
		NumericalVpcId:            utils.PtrUint64(lb.NumericalVpcId),
		SubnetId:                  utils.PtrString(lb.SubnetId),
		AddressIPVersion:          utils.PtrString(lb.AddressIPVersion),
		AddressIPv6:               utils.PtrString(lb.AddressIPv6),
		IPv6Mode:                  utils.PtrString(lb.IPv6Mode),
		VipIsp:                    utils.PtrString(lb.VipIsp),
		OpenBgp:                   utils.PtrUint64(lb.OpenBgp),
		LocalBgp:                  utils.PtrBool(lb.LocalBgp),
		AnycastZone:               utils.PtrString(lb.AnycastZone),
		Egress:                    utils.PtrString(lb.Egress),
		ChargeType:                utils.PtrString(lb.ChargeType),
		SlaType:                   utils.PtrString(lb.SlaType),
		Isolation:                 utils.PtrUint64(lb.Isolation),
		IsolatedTime:              utils.ParseTimestamp(lb.IsolatedTime),
		ExpireTime:                utils.ParseTimestamp(lb.ExpireTime),
		IsBlock:                   utils.PtrBool(lb.IsBlock),
		IsBlockTime:               utils.ParseTimestamp(lb.IsBlockTime),
		IsDDos:                    utils.PtrBool(lb.IsDDos),
		Snat:                      utils.PtrBool(lb.Snat),
		SnatPro:                   utils.PtrBool(lb.SnatPro),
		SnatIps:                   lb.SnatIps,
		MixIpTarget:               utils.PtrBool(lb.MixIpTarget),
		Exclusive:                 utils.PtrUint64(lb.Exclusive),
		Log:                       utils.PtrString(lb.Log),
		LogSetId:                  utils.PtrString(lb.LogSetId),
		LogTopicId:                utils.PtrString(lb.LogTopicId),
		HealthLogSetId:            utils.PtrString(lb.HealthLogSetId),
		HealthLogTopicId:          utils.PtrString(lb.HealthLogTopicId),
		ConfigId:                  utils.PtrString(lb.ConfigId),
		LoadBalancerPassToTarget:  utils.PtrBool(lb.LoadBalancerPassToTarget),
		TargetCount:               utils.PtrUint64(lb.TargetCount),
		AssociateEndpoint:         utils.PtrString(lb.AssociateEndpoint),
		ClusterTag:                utils.PtrString(lb.ClusterTag),
		ClusterIds:                utils.PtrStringSlice(lb.ClusterIds),
		NfvInfo:                   utils.PtrString(lb.NfvInfo),
		Zones:                     utils.PtrStringSlice(lb.Zones),
		MasterZone:                lb.MasterZone,
		BackupZoneSet:             lb.BackupZoneSet,
		AvailableZoneAffinityInfo: lb.AvailableZoneAffinityInfo,
		TargetRegionInfo:          lb.TargetRegionInfo,
		NetworkAttributes:         lb.NetworkAttributes,
		PrepaidAttributes:         lb.PrepaidAttributes,
		ExtraInfo:                 lb.ExtraInfo,
		ExclusiveCluster:          lb.ExclusiveCluster,
		SecureGroups:              utils.PtrStringSlice(lb.SecureGroups),
		AttributeFlags:            utils.PtrStringSlice(lb.AttributeFlags),
		Tags:                      TagsToMap(lb.Tags),
		Region:                    region,
	}
	if len(row.SecureGroups) > 0 {
		row.SecurityGroup = row.SecureGroups[0]
	}
	row.Title = row.LoadBalancerName
	if row.Title == "" {
		row.Title = row.LoadBalancerId
	}
	if row.LoadBalancerId != "" {
		row.Akas = []string{"tencentcloud:clb:load-balancer:" + row.LoadBalancerId}
	}
	return row
}
