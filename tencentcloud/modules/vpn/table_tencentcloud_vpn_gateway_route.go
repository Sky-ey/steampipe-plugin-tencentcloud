package vpn

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"time"

	vpc "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/vpc/v20170312"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudVpnGatewayRoute() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_vpn_gateway_route",
		Description:       "Tencent Cloud VPN gateway destination routes that route IDC IP ranges to VPN tunnels or CCN instances.",
		GetMatrixItemFunc: buildVpnRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "vpn_gateway_id", Require: plugin.Optional},
				{Name: "destination_cidr_block", Require: plugin.Optional},
			},
			Hydrate: listVpnGatewayRoutes,
			Tags:    map[string]string{"service": "vpn", "action": "DescribeVpnGatewayRoutes"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "vpn_gateway_id", Type: proto.ColumnType_STRING, Description: "The parent VPN gateway instance ID.", Transform: transform.FromField("VpnGatewayId")},
			{Name: "route_id", Type: proto.ColumnType_STRING, Description: "The route ID.", Transform: transform.FromField("RouteId")},
			{Name: "destination_cidr_block", Type: proto.ColumnType_STRING, Description: "The destination IDC IP range.", Transform: transform.FromField("DestinationCidrBlock")},
			{Name: "instance_type", Type: proto.ColumnType_STRING, Description: "The next hop type (VPNCONN, CCN).", Transform: transform.FromField("InstanceType")},
			{Name: "instance_id", Type: proto.ColumnType_STRING, Description: "The next hop instance ID.", Transform: transform.FromField("InstanceId")},
			{Name: "priority", Type: proto.ColumnType_INT, Description: "The route priority (0 or 100).", Transform: transform.FromField("Priority")},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The route status (ENABLE, DISABLE).", Transform: transform.FromField("Status")},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The route type (VPC, CCN, Static, BGP).", Transform: transform.FromField("Type")},
			{Name: "create_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the route.", Transform: transform.FromField("CreateTime").NullIfZero()},
			{Name: "update_time", Type: proto.ColumnType_TIMESTAMP, Description: "The last update time of the route.", Transform: transform.FromField("UpdateTime").NullIfZero()},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the route resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listVpnGatewayRoutes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listVpnGatewayRoutes", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listVpnGatewayRoutes", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	gwIds, err := resolveVpnGatewayIDs(ctx, d, client)
	if err != nil {
		plugin.Logger(ctx).Error("listVpnGatewayRoutes", "resolve_gateway_ids_error", err)
		return nil, err
	}

	filters := buildVpnGatewayRouteFilters(d.EqualsQuals)

	for _, gwId := range gwIds {
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}

		pageSize := utils.PageSizeFromLimit(d, 100, 1)
		var offset int64

		for {
			d.WaitForListRateLimit(ctx)

			req := vpc.NewDescribeVpnGatewayRoutesRequest()
			req.VpnGatewayId = &gwId
			req.Offset = &offset
			req.Limit = &pageSize
			if len(filters) > 0 {
				req.Filters = filters
			}

			utils.LogRequest(ctx, "listVpnGatewayRoutes", req)
			resp, err := client.DescribeVpnGatewayRoutesWithContext(ctx, req)
			if err != nil {
				plugin.Logger(ctx).Error("listVpnGatewayRoutes", "request_error", err)
				return nil, err
			}
			if resp == nil || resp.Response == nil {
				break
			}

			routes := resp.Response.Routes
			total := utils.PtrUint64(resp.Response.TotalCount)

			for _, r := range routes {
				if r == nil {
					continue
				}
				row := toVpnGatewayRouteRow(r, gwId)
				row.Region = region
				d.StreamListItem(ctx, row)
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}

			if utils.PageDone(offset, int64(len(routes)), total) {
				break
			}
			offset += int64(len(routes))
		}
	}

	return nil, nil
}

func buildVpnGatewayRouteFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*vpc.Filter {
	return utils.BuildFilters(equalsQuals, newVpnFilter,
		utils.FilterMapping{QualName: "destination_cidr_block", FilterName: "DestinationCidr"},
	)
}

// Row Type

type vpnGatewayRouteRow struct {
	VpnGatewayId         string
	RouteId              string
	DestinationCidrBlock string
	InstanceType         string
	InstanceId           string
	Priority             int64
	Status               string
	Type                 string
	CreateTime           *time.Time
	UpdateTime           *time.Time
	Region               string
	Title                string
	Akas                 []string
}

func toVpnGatewayRouteRow(r *vpc.VpnGatewayRoute, gwId string) vpnGatewayRouteRow {
	row := vpnGatewayRouteRow{
		VpnGatewayId:         gwId,
		RouteId:              utils.PtrString(r.RouteId),
		DestinationCidrBlock: utils.PtrString(r.DestinationCidrBlock),
		InstanceType:         utils.PtrString(r.InstanceType),
		InstanceId:           utils.PtrString(r.InstanceId),
		Priority:             utils.PtrInt64(r.Priority),
		Status:               utils.PtrString(r.Status),
		Type:                 utils.PtrString(r.Type),
		CreateTime:           utils.ParseTimestamp(r.CreateTime),
		UpdateTime:           utils.ParseTimestamp(r.UpdateTime),
	}

	row.Title = row.RouteId

	if row.VpnGatewayId != "" && row.RouteId != "" {
		row.Akas = []string{"tencentcloud:vpn:gateway-route:" + row.VpnGatewayId + "/" + row.RouteId}
	}
	return row
}
