package vpn

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"

	vpc "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/vpc/v20170312"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudVpnGatewayCcnRoute() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_vpn_gateway_ccn_route",
		Description:       "Tencent Cloud CCN routes propagated by a CCN-type VPN gateway (IDC IP ranges advertised to the CCN).",
		GetMatrixItemFunc: buildVpnRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "vpn_gateway_id", Require: plugin.Optional},
			},
			Hydrate: listVpnGatewayCcnRoutes,
			Tags:    map[string]string{"service": "vpn", "action": "DescribeVpnGatewayCcnRoutes"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "vpn_gateway_id", Type: proto.ColumnType_STRING, Description: "The parent VPN gateway instance ID.", Transform: transform.FromField("VpnGatewayId")},
			{Name: "route_id", Type: proto.ColumnType_STRING, Description: "The route ID.", Transform: transform.FromField("RouteId")},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Whether the route is enabled (ENABLE, DISABLE).", Transform: transform.FromField("Status")},
			{Name: "destination_cidr_block", Type: proto.ColumnType_STRING, Description: "The route CIDR block.", Transform: transform.FromField("DestinationCidrBlock")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the route resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listVpnGatewayCcnRoutes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listVpnGatewayCcnRoutes", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listVpnGatewayCcnRoutes", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	gwIds, err := resolveVpnGatewayIDs(ctx, d, client)
	if err != nil {
		plugin.Logger(ctx).Error("listVpnGatewayCcnRoutes", "resolve_gateway_ids_error", err)
		return nil, err
	}

	for _, gwId := range gwIds {
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}

		pageSize := uint64(utils.PageSizeFromLimit(d, 100, 1))
		var offset uint64

		for {
			d.WaitForListRateLimit(ctx)

			req := vpc.NewDescribeVpnGatewayCcnRoutesRequest()
			req.VpnGatewayId = &gwId
			req.Offset = &offset
			req.Limit = &pageSize

			utils.LogRequest(ctx, "listVpnGatewayCcnRoutes", req)
			resp, err := client.DescribeVpnGatewayCcnRoutesWithContext(ctx, req)
			if err != nil {
				plugin.Logger(ctx).Error("listVpnGatewayCcnRoutes", "request_error", err)
				return nil, err
			}
			if resp == nil || resp.Response == nil {
				break
			}

			routes := resp.Response.RouteSet
			total := utils.PtrUint64(resp.Response.TotalCount)

			for _, r := range routes {
				if r == nil {
					continue
				}
				row := toVpnGatewayCcnRouteRow(r, gwId)
				row.Region = region
				d.StreamListItem(ctx, row)
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}

			if utils.PageDone(int64(offset), int64(len(routes)), total) {
				break
			}
			offset += uint64(len(routes))
		}
	}

	return nil, nil
}

// Row Type

type vpnGatewayCcnRouteRow struct {
	VpnGatewayId         string
	RouteId              string
	Status               string
	DestinationCidrBlock string
	Region               string
	Title                string
	Akas                 []string
}

func toVpnGatewayCcnRouteRow(r *vpc.VpngwCcnRoutes, gwId string) vpnGatewayCcnRouteRow {
	row := vpnGatewayCcnRouteRow{
		VpnGatewayId:         gwId,
		RouteId:              utils.PtrString(r.RouteId),
		Status:               utils.PtrString(r.Status),
		DestinationCidrBlock: utils.PtrString(r.DestinationCidrBlock),
	}

	row.Title = row.RouteId

	if row.VpnGatewayId != "" && row.RouteId != "" {
		row.Akas = []string{"tencentcloud:vpn:gateway-ccn-route:" + row.VpnGatewayId + "/" + row.RouteId}
	}
	return row
}
