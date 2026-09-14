package vpn

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

func TableTencentcloudVpnGateway() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_vpn_gateway",
		Description:       "Tencent Cloud VPN gateways that establish IPsec/SSL tunnels between a VPC and on-premises networks, or CCN-attached gateways.",
		GetMatrixItemFunc: buildVpnRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "vpn_gateway_id", Require: plugin.Optional},
				{Name: "vpc_id", Require: plugin.Optional},
				{Name: "vpn_gateway_name", Require: plugin.Optional},
				{Name: "type", Require: plugin.Optional},
				{Name: "zone", Require: plugin.Optional},
				{Name: "public_ip_address", Require: plugin.Optional},
				{Name: "renew_flag", Require: plugin.Optional},
			},
			Hydrate: listVpnGateways,
			Tags:    map[string]string{"service": "vpn", "action": "DescribeVpnGateways"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("vpn_gateway_id"),
			Hydrate:    getVpnGateway,
			Tags:       map[string]string{"service": "vpn", "action": "DescribeVpnGateways"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "vpn_gateway_id", Type: proto.ColumnType_STRING, Description: "The VPN gateway instance ID, such as `vpngw-f49l6u0z`.", Transform: transform.FromField("VpnGatewayId")},
			{Name: "vpn_gateway_name", Type: proto.ColumnType_STRING, Description: "The VPN gateway name.", Transform: transform.FromField("VpnGatewayName")},
			{Name: "vpc_id", Type: proto.ColumnType_STRING, Description: "The VPC instance ID to which the VPN gateway belongs.", Transform: transform.FromField("VpcId")},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The gateway type (IPSEC, SSL, CCN).", Transform: transform.FromField("Type")},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "The gateway lifecycle state (PENDING, DELETING, AUTHENTICATING, DELETED, AVAILABLE).", Transform: transform.FromField("State")},
			{Name: "restrict_state", Type: proto.ColumnType_STRING, Description: "The gateway billing state (PROTECTIVELY_ISOLATED, NORMAL).", Transform: transform.FromField("RestrictState")},
			{Name: "public_ip_address", Type: proto.ColumnType_STRING, Description: "The public IP address of the gateway.", Transform: transform.FromField("PublicIpAddress")},
			{Name: "internet_max_bandwidth_out", Type: proto.ColumnType_INT, Description: "The maximum outbound bandwidth in Mbps.", Transform: transform.FromField("InternetMaxBandwidthOut")},
			{Name: "zone", Type: proto.ColumnType_STRING, Description: "The availability zone, such as `ap-guangzhou-2`.", Transform: transform.FromField("Zone")},
			{Name: "version", Type: proto.ColumnType_STRING, Description: "The gateway instance version.", Transform: transform.FromField("Version")},
			{Name: "max_connection", Type: proto.ColumnType_INT, Description: "The maximum number of SSL VPN clients allowed to connect.", Transform: transform.FromField("MaxConnection")},
			{Name: "network_instance_id", Type: proto.ColumnType_STRING, Description: "The CCN instance ID when Type is CCN.", Transform: transform.FromField("NetworkInstanceId")},
			{Name: "cdc_id", Type: proto.ColumnType_STRING, Description: "The CDC instance ID if the gateway is associated with a CDC.", Transform: transform.FromField("CdcId")},
			{Name: "instance_charge_type", Type: proto.ColumnType_STRING, Description: "The billing type (POSTPAID_BY_HOUR, PREPAID).", Transform: transform.FromField("InstanceChargeType")},
			{Name: "renew_flag", Type: proto.ColumnType_STRING, Description: "The renewal type (NOTIFY_AND_MANUAL_RENEW, NOTIFY_AND_AUTO_RENEW, NOT_NOTIFY_AND_NOT_RENEW).", Transform: transform.FromField("RenewFlag")},
			{Name: "new_purchase_plan", Type: proto.ColumnType_STRING, Description: "The billing plan change (PREPAID_TO_POSTPAID).", Transform: transform.FromField("NewPurchasePlan")},
			{Name: "is_address_blocked", Type: proto.ColumnType_BOOL, Description: "Whether the public IP address is blocked.", Transform: transform.FromField("IsAddressBlocked")},
			{Name: "vpn_gateway_quota_set", Type: proto.ColumnType_JSON, Description: "The bandwidth quota information of the gateway.", Transform: transform.FromField("VpnGatewayQuotaSet")},
			{Name: "created_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the gateway.", Transform: transform.FromField("CreatedTime").NullIfZero()},
			{Name: "expired_time", Type: proto.ColumnType_TIMESTAMP, Description: "The expiration time of the prepaid gateway.", Transform: transform.FromField("ExpiredTime").NullIfZero()},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the gateway resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listVpnGateways(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listVpnGateways", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listVpnGateways", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	pageSize := uint64(utils.PageSizeFromLimit(d, 100, 1))
	var offset uint64

	ids := utils.SingleIdFromQual(d.EqualsQuals, "vpn_gateway_id")
	filters := buildVpnGatewayFilters(d.EqualsQuals)

	for {
		d.WaitForListRateLimit(ctx)

		req := vpc.NewDescribeVpnGatewaysRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		if len(ids) > 0 {
			req.VpnGatewayIds = ids
		} else if len(filters) > 0 {
			req.Filters = filters
		}

		utils.LogRequest(ctx, "listVpnGateways", req)
		resp, err := client.DescribeVpnGatewaysWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listVpnGateways", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		gws := resp.Response.VpnGatewaySet
		total := utils.PtrUint64(resp.Response.TotalCount)

		for _, gw := range gws {
			if gw == nil {
				continue
			}
			row := toVpnGatewayRow(gw)
			row.Region = region
			d.StreamListItem(ctx, row)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if utils.PageDone(int64(offset), int64(len(gws)), total) {
			break
		}
		offset += uint64(len(gws))
	}

	return nil, nil
}

func buildVpnGatewayFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*vpc.FilterObject {
	return utils.BuildFilters(equalsQuals, newVpnFilterObject,
		utils.FilterMapping{QualName: "vpc_id", FilterName: "vpc-id"},
		utils.FilterMapping{QualName: "vpn_gateway_name", FilterName: "vpn-gateway-name"},
		utils.FilterMapping{QualName: "type", FilterName: "type"},
		utils.FilterMapping{QualName: "zone", FilterName: "zone"},
		utils.FilterMapping{QualName: "public_ip_address", FilterName: "public-ip-address"},
		utils.FilterMapping{QualName: "renew_flag", FilterName: "renew-flag"},
	)
}

// Get Function

func getVpnGateway(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	vpnGatewayId := d.EqualsQuals["vpn_gateway_id"].GetStringValue()
	plugin.Logger(ctx).Debug("getVpnGateway", "vpn_gateway_id", vpnGatewayId, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getVpnGateway", "client_init_error", err)
		return nil, err
	}

	if vpnGatewayId == "" {
		return nil, nil
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	req := vpc.NewDescribeVpnGatewaysRequest()
	req.VpnGatewayIds = common.StringPtrs([]string{vpnGatewayId})

	utils.LogRequest(ctx, "getVpnGateway", req)
	resp, err := client.DescribeVpnGatewaysWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getVpnGateway", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.VpnGatewaySet) == 0 {
		return nil, nil
	}

	row := toVpnGatewayRow(resp.Response.VpnGatewaySet[0])
	row.Region = region
	return row, nil
}

// Row Type

type vpnGatewayRow struct {
	VpnGatewayId            string
	VpnGatewayName          string
	VpcId                   string
	Type                    string
	State                   string
	RestrictState           string
	PublicIpAddress         string
	InternetMaxBandwidthOut uint64
	Zone                    string
	Version                 string
	MaxConnection           uint64
	NetworkInstanceId       string
	CdcId                   string
	InstanceChargeType      string
	RenewFlag               string
	NewPurchasePlan         string
	IsAddressBlocked        bool
	VpnGatewayQuotaSet      []*vpc.VpnGatewayQuota
	CreatedTime             *time.Time
	ExpiredTime             *time.Time
	Region                  string
	Title                   string
	Akas                    []string
}

func toVpnGatewayRow(g *vpc.VpnGateway) vpnGatewayRow {
	row := vpnGatewayRow{
		VpnGatewayId:            utils.PtrString(g.VpnGatewayId),
		VpnGatewayName:          utils.PtrString(g.VpnGatewayName),
		VpcId:                   utils.PtrString(g.VpcId),
		Type:                    utils.PtrString(g.Type),
		State:                   utils.PtrString(g.State),
		RestrictState:           utils.PtrString(g.RestrictState),
		PublicIpAddress:         utils.PtrString(g.PublicIpAddress),
		InternetMaxBandwidthOut: utils.PtrUint64(g.InternetMaxBandwidthOut),
		Zone:                    utils.PtrString(g.Zone),
		Version:                 utils.PtrString(g.Version),
		MaxConnection:           utils.PtrUint64(g.MaxConnection),
		NetworkInstanceId:       utils.PtrString(g.NetworkInstanceId),
		CdcId:                   utils.PtrString(g.CdcId),
		InstanceChargeType:      utils.PtrString(g.InstanceChargeType),
		RenewFlag:               utils.PtrString(g.RenewFlag),
		NewPurchasePlan:         utils.PtrString(g.NewPurchasePlan),
		IsAddressBlocked:        utils.PtrBool(g.IsAddressBlocked),
		VpnGatewayQuotaSet:      g.VpnGatewayQuotaSet,
		CreatedTime:             utils.ParseTimestamp(g.CreatedTime),
		ExpiredTime:             utils.ParseTimestamp(g.ExpiredTime),
	}

	if row.VpnGatewayName != "" {
		row.Title = row.VpnGatewayName
	} else {
		row.Title = row.VpnGatewayId
	}

	if row.VpnGatewayId != "" {
		row.Akas = []string{"tencentcloud:vpn:gateway:" + row.VpnGatewayId}
	}
	return row
}
