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

func TableTencentcloudVpnCustomerGateway() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_vpn_customer_gateway",
		Description:       "Tencent Cloud VPN customer gateways that register the public IP address of the on-premises VPN peer.",
		GetMatrixItemFunc: buildVpnRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "customer_gateway_id", Require: plugin.Optional},
				{Name: "customer_gateway_name", Require: plugin.Optional},
				{Name: "ip_address", Require: plugin.Optional},
			},
			Hydrate: listVpnCustomerGateways,
			Tags:    map[string]string{"service": "vpn", "action": "DescribeCustomerGateways"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("customer_gateway_id"),
			Hydrate:    getVpnCustomerGateway,
			Tags:       map[string]string{"service": "vpn", "action": "DescribeCustomerGateways"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "customer_gateway_id", Type: proto.ColumnType_STRING, Description: "The customer gateway instance ID, such as `cgw-2wqq41m9`.", Transform: transform.FromField("CustomerGatewayId")},
			{Name: "customer_gateway_name", Type: proto.ColumnType_STRING, Description: "The customer gateway name.", Transform: transform.FromField("CustomerGatewayName")},
			{Name: "ip_address", Type: proto.ColumnType_STRING, Description: "The public IP address of the on-premises peer.", Transform: transform.FromField("IpAddress")},
			{Name: "created_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the customer gateway.", Transform: transform.FromField("CreatedTime").NullIfZero()},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the customer gateway resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listVpnCustomerGateways(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listVpnCustomerGateways", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listVpnCustomerGateways", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	pageSize := uint64(utils.PageSizeFromLimit(d, 100, 1))
	var offset uint64

	ids := utils.SingleIdFromQual(d.EqualsQuals, "customer_gateway_id")
	filters := buildVpnCustomerGatewayFilters(d.EqualsQuals)

	for {
		d.WaitForListRateLimit(ctx)

		req := vpc.NewDescribeCustomerGatewaysRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		if len(ids) > 0 {
			req.CustomerGatewayIds = ids
		} else if len(filters) > 0 {
			req.Filters = filters
		}

		utils.LogRequest(ctx, "listVpnCustomerGateways", req)
		resp, err := client.DescribeCustomerGatewaysWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listVpnCustomerGateways", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		gws := resp.Response.CustomerGatewaySet
		total := utils.PtrUint64(resp.Response.TotalCount)

		for _, gw := range gws {
			if gw == nil {
				continue
			}
			row := toVpnCustomerGatewayRow(gw)
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

func buildVpnCustomerGatewayFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*vpc.Filter {
	return utils.BuildFilters(equalsQuals, newVpnFilter,
		utils.FilterMapping{QualName: "customer_gateway_name", FilterName: "customer-gateway-name"},
		utils.FilterMapping{QualName: "ip_address", FilterName: "ip-address"},
	)
}

// Get Function

func getVpnCustomerGateway(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	customerGatewayId := d.EqualsQuals["customer_gateway_id"].GetStringValue()
	plugin.Logger(ctx).Debug("getVpnCustomerGateway", "customer_gateway_id", customerGatewayId, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getVpnCustomerGateway", "client_init_error", err)
		return nil, err
	}

	if customerGatewayId == "" {
		return nil, nil
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	req := vpc.NewDescribeCustomerGatewaysRequest()
	req.CustomerGatewayIds = common.StringPtrs([]string{customerGatewayId})

	utils.LogRequest(ctx, "getVpnCustomerGateway", req)
	resp, err := client.DescribeCustomerGatewaysWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getVpnCustomerGateway", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.CustomerGatewaySet) == 0 {
		return nil, nil
	}

	row := toVpnCustomerGatewayRow(resp.Response.CustomerGatewaySet[0])
	row.Region = region
	return row, nil
}

// Row Type

type vpnCustomerGatewayRow struct {
	CustomerGatewayId   string
	CustomerGatewayName string
	IpAddress           string
	CreatedTime         *time.Time
	Region              string
	Title               string
	Akas                []string
}

func toVpnCustomerGatewayRow(g *vpc.CustomerGateway) vpnCustomerGatewayRow {
	row := vpnCustomerGatewayRow{
		CustomerGatewayId:   utils.PtrString(g.CustomerGatewayId),
		CustomerGatewayName: utils.PtrString(g.CustomerGatewayName),
		IpAddress:           utils.PtrString(g.IpAddress),
		CreatedTime:         utils.ParseTimestamp(g.CreatedTime),
	}

	if row.CustomerGatewayName != "" {
		row.Title = row.CustomerGatewayName
	} else {
		row.Title = row.CustomerGatewayId
	}

	if row.CustomerGatewayId != "" {
		row.Akas = []string{"tencentcloud:vpn:customer-gateway:" + row.CustomerGatewayId}
	}
	return row
}
