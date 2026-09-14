package dc

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

// Note: DescribeDirectConnectGateways lives in the vpc SDK package, not the dc
// package. The table is still registered under the dc module because it is a
// Direct Connect resource; service=dc in the Tags map reflects that.
func TableTencentcloudDcGateway() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_dc_gateway",
		Description:       "Tencent Cloud Direct Connect gateways that bridge Direct Connect tunnels with VPC or CCN networks.",
		GetMatrixItemFunc: buildDcRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "direct_connect_gateway_id", Require: plugin.Optional},
				{Name: "direct_connect_gateway_name", Require: plugin.Optional},
				{Name: "gateway_type", Require: plugin.Optional},
				{Name: "network_type", Require: plugin.Optional},
				{Name: "ccn_id", Require: plugin.Optional},
				{Name: "vpc_id", Require: plugin.Optional},
			},
			Hydrate: listDcGateways,
			Tags:    map[string]string{"service": "dc", "action": "DescribeDirectConnectGateways"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("direct_connect_gateway_id"),
			Hydrate:    getDcGateway,
			Tags:       map[string]string{"service": "dc", "action": "DescribeDirectConnectGateways"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "direct_connect_gateway_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the Direct Connect gateway, such as `dcg-9o233uri`.", Transform: transform.FromField("DirectConnectGatewayId")},
			{Name: "direct_connect_gateway_name", Type: proto.ColumnType_STRING, Description: "The name of the Direct Connect gateway.", Transform: transform.FromField("DirectConnectGatewayName")},
			{Name: "gateway_type", Type: proto.ColumnType_STRING, Description: "The gateway type (NORMAL, NAT).", Transform: transform.FromField("GatewayType")},
			{Name: "network_type", Type: proto.ColumnType_STRING, Description: "The network type (VPC, CCN).", Transform: transform.FromField("NetworkType")},
			{Name: "network_instance_id", Type: proto.ColumnType_STRING, Description: "The ID of the associated network instance (VPC or CCN).", Transform: transform.FromField("NetworkInstanceId")},
			{Name: "vpc_id", Type: proto.ColumnType_STRING, Description: "The ID of the VPC instance associated with the gateway.", Transform: transform.FromField("VpcId")},
			{Name: "ccn_id", Type: proto.ColumnType_STRING, Description: "The ID of the CCN instance associated with the gateway.", Transform: transform.FromField("CcnId")},
			{Name: "ccn_route_type", Type: proto.ColumnType_STRING, Description: "The CCN route-learning type (BGP, STATIC).", Transform: transform.FromField("CcnRouteType")},
			{Name: "direct_connect_gateway_ip", Type: proto.ColumnType_STRING, Description: "The Direct Connect gateway IP.", Transform: transform.FromField("DirectConnectGatewayIp")},
			{Name: "enable_bgp", Type: proto.ColumnType_BOOL, Description: "Whether BGP is enabled.", Transform: transform.FromField("EnableBGP")},
			{Name: "enable_bgp_community", Type: proto.ColumnType_BOOL, Description: "Whether the BGP community attribute is enabled.", Transform: transform.FromField("EnableBGPCommunity")},
			{Name: "nat_gateway_id", Type: proto.ColumnType_STRING, Description: "The ID of the NAT gateway bound to the gateway.", Transform: transform.FromField("NatGatewayId")},
			{Name: "mode_type", Type: proto.ColumnType_STRING, Description: "The CCN route publishing mode (standard, exquisite).", Transform: transform.FromField("ModeType")},
			{Name: "access_network_type", Type: proto.ColumnType_STRING, Description: "The access network type (VXLAN, MPLS, Hybrid).", Transform: transform.FromField("AccessNetworkType")},
			{Name: "zone", Type: proto.ColumnType_STRING, Description: "The availability zone where the gateway resides.", Transform: transform.FromField("Zone")},
			{Name: "ha_zone_list", Type: proto.ColumnType_JSON, Description: "The AZ list of the gateway with cross-AZ placement groups.", Transform: transform.FromField("HaZoneList")},
			{Name: "vxlan_support", Type: proto.ColumnType_JSON, Description: "Whether the gateway supports the VXLAN architecture.", Transform: transform.FromField("VXLANSupport")},
			{Name: "local_zone", Type: proto.ColumnType_BOOL, Description: "Whether the gateway is for an edge zone.", Transform: transform.FromField("LocalZone")},
			{Name: "new_afc", Type: proto.ColumnType_INT, Description: "Whether gateway traffic monitoring is supported (0: no, 1: yes).", Transform: transform.FromField("NewAfc")},
			{Name: "enable_flow_details", Type: proto.ColumnType_INT, Description: "The status of gateway traffic monitoring (0: disabled, 1: enabled).", Transform: transform.FromField("EnableFlowDetails")},
			{Name: "flow_details_update_time", Type: proto.ColumnType_TIMESTAMP, Description: "The last time gateway traffic monitoring was enabled or disabled.", Transform: transform.FromField("FlowDetailsUpdateTime").NullIfZero()},
			{Name: "create_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the gateway.", Transform: transform.FromField("CreateTime").NullIfZero()},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the gateway resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("DirectConnectGatewayName")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listDcGateways(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listDcGateways", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listDcGateways", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	pageSize := uint64(utils.PageSizeFromLimit(d, 100, 1))
	var offset uint64 = 0

	gatewayIds := utils.SingleIdFromQual(d.EqualsQuals, "direct_connect_gateway_id")
	filters := buildDcGatewayFilters(d.EqualsQuals)

	for {
		d.WaitForListRateLimit(ctx)

		req := vpc.NewDescribeDirectConnectGatewaysRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		if len(gatewayIds) > 0 {
			req.DirectConnectGatewayIds = gatewayIds
		} else if len(filters) > 0 {
			req.Filters = filters
		}

		utils.LogRequest(ctx, "listDcGateways", req)
		resp, err := client.DescribeDirectConnectGatewaysWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listDcGateways", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		items := resp.Response.DirectConnectGatewaySet
		total := utils.PtrUint64(resp.Response.TotalCount)

		for _, item := range items {
			if item == nil {
				continue
			}
			row := toDcGatewayRow(item)
			row.Region = region
			d.StreamListItem(ctx, row)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if utils.PageDone(int64(offset), int64(len(items)), total) {
			break
		}
		offset += uint64(len(items))
	}

	return nil, nil
}

func buildDcGatewayFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*vpc.Filter {
	filters := make([]*vpc.Filter, 0, 6)
	add := func(qualName, filterName string) {
		qual := equalsQuals[qualName]
		if qual == nil {
			return
		}
		value := qual.GetStringValue()
		if value == "" {
			return
		}
		filters = append(filters, &vpc.Filter{
			Name:   common.StringPtr(filterName),
			Values: common.StringPtrs([]string{value}),
		})
	}
	add("direct_connect_gateway_name", "direct-connect-gateway-name")
	add("gateway_type", "gateway-type")
	add("network_type", "network-type")
	add("ccn_id", "ccn-id")
	add("vpc_id", "vpc-id")
	return filters
}

// Get Function

func getDcGateway(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	gatewayId := d.EqualsQuals["direct_connect_gateway_id"].GetStringValue()
	plugin.Logger(ctx).Debug("getDcGateway", "direct_connect_gateway_id", gatewayId, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getDcGateway", "client_init_error", err)
		return nil, err
	}

	if gatewayId == "" {
		return nil, nil
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	req := vpc.NewDescribeDirectConnectGatewaysRequest()
	req.DirectConnectGatewayIds = common.StringPtrs([]string{gatewayId})

	utils.LogRequest(ctx, "getDcGateway", req)
	resp, err := client.DescribeDirectConnectGatewaysWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getDcGateway", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.DirectConnectGatewaySet) == 0 {
		return nil, nil
	}

	row := toDcGatewayRow(resp.Response.DirectConnectGatewaySet[0])
	row.Region = region
	return row, nil
}

// Row Type

type dcGatewayRow struct {
	DirectConnectGatewayId   string
	DirectConnectGatewayName string
	GatewayType              string
	NetworkType              string
	NetworkInstanceId        string
	VpcId                    string
	CcnId                    string
	CcnRouteType             string
	DirectConnectGatewayIp   string
	EnableBGP                bool
	EnableBGPCommunity       bool
	NatGatewayId             string
	ModeType                 string
	AccessNetworkType        string
	Zone                     string
	HaZoneList               []*string
	VXLANSupport             []*bool
	LocalZone                bool
	NewAfc                   uint64
	EnableFlowDetails        uint64
	FlowDetailsUpdateTime    *time.Time
	CreateTime               *time.Time
	Region                   string
	Akas                     []string
}

func toDcGatewayRow(s *vpc.DirectConnectGateway) dcGatewayRow {
	row := dcGatewayRow{
		DirectConnectGatewayId:   utils.PtrString(s.DirectConnectGatewayId),
		DirectConnectGatewayName: utils.PtrString(s.DirectConnectGatewayName),
		GatewayType:              utils.PtrString(s.GatewayType),
		NetworkType:              utils.PtrString(s.NetworkType),
		NetworkInstanceId:        utils.PtrString(s.NetworkInstanceId),
		VpcId:                    utils.PtrString(s.VpcId),
		CcnId:                    utils.PtrString(s.CcnId),
		CcnRouteType:             utils.PtrString(s.CcnRouteType),
		DirectConnectGatewayIp:   utils.PtrString(s.DirectConnectGatewayIp),
		EnableBGP:                utils.PtrBool(s.EnableBGP),
		EnableBGPCommunity:       utils.PtrBool(s.EnableBGPCommunity),
		NatGatewayId:             utils.PtrString(s.NatGatewayId),
		ModeType:                 utils.PtrString(s.ModeType),
		AccessNetworkType:        utils.PtrString(s.AccessNetworkType),
		Zone:                     utils.PtrString(s.Zone),
		HaZoneList:               s.HaZoneList,
		VXLANSupport:             s.VXLANSupport,
		LocalZone:                utils.PtrBool(s.LocalZone),
		NewAfc:                   utils.PtrUint64(s.NewAfc),
		EnableFlowDetails:        utils.PtrUint64(s.EnableFlowDetails),
		FlowDetailsUpdateTime:    utils.ParseTimestamp(s.FlowDetailsUpdateTime),
		CreateTime:               utils.ParseTimestamp(s.CreateTime),
	}

	if row.DirectConnectGatewayId != "" {
		row.Akas = []string{"tencentcloud:dc:gateway:" + row.DirectConnectGatewayId}
	}
	return row
}
