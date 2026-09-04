package vpc

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/vpc/v20170312"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudVpcNatGateway() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_vpc_nat_gateway",
		Description:       "Tencent Cloud VPC NAT gateways that provide SNAT and DNAT translation for private subnets.",
		GetMatrixItemFunc: buildVpcRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "nat_gateway_id", Require: plugin.Optional},
				{Name: "vpc_id", Require: plugin.Optional},
				{Name: "nat_gateway_name", Require: plugin.Optional},
			},
			Hydrate: listVpcNatGateways,
			Tags:    map[string]string{"service": "vpc", "action": "DescribeNatGateways"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("nat_gateway_id"),
			Hydrate:    getVpcNatGateway,
			Tags:       map[string]string{"service": "vpc", "action": "DescribeNatGateways"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "nat_gateway_id", Type: proto.ColumnType_STRING, Description: "The NAT gateway ID, such as `nat-123xx454`.", Transform: transform.FromField("NatGatewayId")},
			{Name: "nat_gateway_name", Type: proto.ColumnType_STRING, Description: "The NAT gateway name.", Transform: transform.FromField("NatGatewayName")},
			{Name: "vpc_id", Type: proto.ColumnType_STRING, Description: "The VPC instance ID to which the NAT gateway belongs.", Transform: transform.FromField("VpcId")},
			{Name: "subnet_id", Type: proto.ColumnType_STRING, Description: "The subnet ID to which the NAT gateway belongs.", Transform: transform.FromField("SubnetId")},
			{Name: "zone", Type: proto.ColumnType_STRING, Description: "The availability zone of the NAT gateway.", Transform: transform.FromField("Zone")},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "The lifecycle state (PENDING, DELETING, AVAILABLE, UPDATING, FAILED).", Transform: transform.FromField("State")},
			{Name: "network_state", Type: proto.ColumnType_STRING, Description: "The runtime state (AVAILABLE, UNAVAILABLE, INSUFFICIENT).", Transform: transform.FromField("NetworkState")},
			{Name: "restrict_state", Type: proto.ColumnType_STRING, Description: "Whether the NAT gateway is blocked (NORMAL, RESTRICTED).", Transform: transform.FromField("RestrictState")},
			{Name: "nat_product_version", Type: proto.ColumnType_INT, Description: "NAT gateway major version. 1 (Classic), 2 (Standard).", Transform: transform.FromField("NatProductVersion")},
			{Name: "internet_max_bandwidth_out", Type: proto.ColumnType_INT, Description: "The maximum outbound bandwidth of the gateway (Mbps).", Transform: transform.FromField("InternetMaxBandwidthOut")},
			{Name: "max_concurrent_connection", Type: proto.ColumnType_INT, Description: "The concurrent connections cap of the gateway.", Transform: transform.FromField("MaxConcurrentConnection")},
			{Name: "exclusive_gateway_bandwidth", Type: proto.ColumnType_INT, Description: "Bandwidth of the dedicated gateway cluster (Mbps).", Transform: transform.FromField("ExclusiveGatewayBandwidth")},
			{Name: "is_exclusive", Type: proto.ColumnType_BOOL, Description: "Whether the NAT gateway is dedicated.", Transform: transform.FromField("IsExclusive")},
			{Name: "smart_schedule_mode", Type: proto.ColumnType_BOOL, Description: "Whether EIP selection for SNAT based on destination network segment is enabled.", Transform: transform.FromField("SmartScheduleMode")},
			{Name: "public_ip_address_set", Type: proto.ColumnType_JSON, Description: "The public IP addresses bound to the NAT gateway.", Transform: transform.FromField("PublicIpAddressSet")},
			{Name: "destination_ip_port_translation_nat_rule_set", Type: proto.ColumnType_JSON, Description: "The DNAT (port forwarding) rules of the NAT gateway.", Transform: transform.FromField("DestinationIpPortTranslationNatRuleSet")},
			{Name: "source_ip_translation_nat_rule_set", Type: proto.ColumnType_JSON, Description: "The SNAT rules of the NAT gateway.", Transform: transform.FromField("SourceIpTranslationNatRuleSet")},
			{Name: "direct_connect_gateway_ids", Type: proto.ColumnType_JSON, Description: "IDs of the direct connect gateways bound to the NAT gateway.", Transform: transform.FromField("DirectConnectGatewayIds")},
			{Name: "security_group_set", Type: proto.ColumnType_JSON, Description: "The security groups bound to the NAT gateway.", Transform: transform.FromField("SecurityGroupSet")},
			{Name: "created_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the NAT gateway.", Transform: transform.FromField("CreatedTime").NullIfZero()},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The tags attached to the NAT gateway as a map of key -> value.", Transform: transform.FromField("Tags")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the NAT gateway resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("NatGatewayName")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listVpcNatGateways(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listVpcNatGateways", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listVpcNatGateways", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	pageSize := uint64(utils.PageSizeFromLimit(d, 100, 1))
	var offset uint64

	natIds := utils.SingleIdFromQual(d.EqualsQuals, "nat_gateway_id")
	filters := utils.BuildFilters(d.EqualsQuals, newVpcFilter,
		utils.FilterMapping{QualName: "vpc_id", FilterName: "vpc-id"},
		utils.FilterMapping{QualName: "nat_gateway_name", FilterName: "nat-gateway-name"},
	)

	for {
		d.WaitForListRateLimit(ctx)

		req := vpc.NewDescribeNatGatewaysRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		if len(natIds) > 0 {
			req.NatGatewayIds = natIds
		} else if len(filters) > 0 {
			req.Filters = filters
		}

		utils.LogRequest(ctx, "listVpcNatGateways", req)
		resp, err := client.DescribeNatGatewaysWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listVpcNatGateways", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		nats := resp.Response.NatGatewaySet
		total := utils.PtrUint64(resp.Response.TotalCount)

		for _, item := range nats {
			if item == nil {
				continue
			}
			row := toVpcNatGatewayRow(item)
			row.Region = region
			d.StreamListItem(ctx, row)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if utils.PageDone(int64(offset), int64(len(nats)), total) {
			break
		}
		offset += uint64(len(nats))
	}

	return nil, nil
}

// Get Function

func getVpcNatGateway(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	natId := d.EqualsQuals["nat_gateway_id"].GetStringValue()
	plugin.Logger(ctx).Info("getVpcNatGateway", "nat_gateway_id", natId, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getVpcNatGateway", "client_init_error", err)
		return nil, err
	}

	if natId == "" {
		return nil, nil
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	req := vpc.NewDescribeNatGatewaysRequest()
	req.NatGatewayIds = common.StringPtrs([]string{natId})

	utils.LogRequest(ctx, "getVpcNatGateway", req)
	resp, err := client.DescribeNatGatewaysWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getVpcNatGateway", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.NatGatewaySet) == 0 {
		return nil, nil
	}

	row := toVpcNatGatewayRow(resp.Response.NatGatewaySet[0])
	row.Region = region
	return row, nil
}

// Row Type

type vpcNatGatewayRow struct {
	NatGatewayId                           string
	NatGatewayName                         string
	VpcId                                  string
	SubnetId                               string
	Zone                                   string
	State                                  string
	NetworkState                           string
	RestrictState                          string
	NatProductVersion                      uint64
	InternetMaxBandwidthOut                uint64
	MaxConcurrentConnection                uint64
	ExclusiveGatewayBandwidth              uint64
	IsExclusive                            bool
	SmartScheduleMode                      bool
	PublicIpAddressSet                     []*vpc.NatGatewayAddress
	DestinationIpPortTranslationNatRuleSet []*vpc.DestinationIpPortTranslationNatRule
	SourceIpTranslationNatRuleSet          []*vpc.SourceIpTranslationNatRule
	DirectConnectGatewayIds                []string
	SecurityGroupSet                       []string
	CreatedTime                            *time.Time
	Tags                                   map[string]string
	Region                                 string
	Akas                                   []string
}

func toVpcNatGatewayRow(n *vpc.NatGateway) vpcNatGatewayRow {
	row := vpcNatGatewayRow{
		NatGatewayId:                           utils.PtrString(n.NatGatewayId),
		NatGatewayName:                         utils.PtrString(n.NatGatewayName),
		VpcId:                                  utils.PtrString(n.VpcId),
		SubnetId:                               utils.PtrString(n.SubnetId),
		Zone:                                   utils.PtrString(n.Zone),
		State:                                  utils.PtrString(n.State),
		NetworkState:                           utils.PtrString(n.NetworkState),
		RestrictState:                          utils.PtrString(n.RestrictState),
		NatProductVersion:                      utils.PtrUint64(n.NatProductVersion),
		InternetMaxBandwidthOut:                utils.PtrUint64(n.InternetMaxBandwidthOut),
		MaxConcurrentConnection:                utils.PtrUint64(n.MaxConcurrentConnection),
		ExclusiveGatewayBandwidth:              utils.PtrUint64(n.ExclusiveGatewayBandwidth),
		IsExclusive:                            utils.PtrBool(n.IsExclusive),
		SmartScheduleMode:                      utils.PtrBool(n.SmartScheduleMode),
		PublicIpAddressSet:                     n.PublicIpAddressSet,
		DestinationIpPortTranslationNatRuleSet: n.DestinationIpPortTranslationNatRuleSet,
		SourceIpTranslationNatRuleSet:          n.SourceIpTranslationNatRuleSet,
		DirectConnectGatewayIds:                utils.PtrStringSlice(n.DirectConnectGatewayIds),
		SecurityGroupSet:                       utils.PtrStringSlice(n.SecurityGroupSet),
		CreatedTime:                            utils.ParseTimestamp(n.CreatedTime),
		Tags:                                   TagsToMap(n.TagSet),
	}

	if row.NatGatewayId != "" {
		row.Akas = []string{"tencentcloud:vpc:nat-gateway:" + row.NatGatewayId}
	}
	return row
}
