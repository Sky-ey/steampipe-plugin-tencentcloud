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

// DescribeDirectConnectGatewayCcnRoutes lives in the vpc SDK package and is a
// fan-out sub-resource: DirectConnectGatewayId is a required single-value
// parameter (not an array), with no Filters or IDs array. When the qual is not
// provided we enumerate every Direct Connect gateway in the region first and
// then query routes for each parent gateway.
func TableTencentcloudDcGatewayCcnRoute() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_dc_gateway_ccn_route",
		Description:       "Tencent Cloud Direct Connect gateway CCN routes (IDC IP ranges) learned by each Direct Connect gateway.",
		GetMatrixItemFunc: buildDcRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "direct_connect_gateway_id", Require: plugin.Optional},
				{Name: "ccn_route_type", Require: plugin.Optional},
			},
			Hydrate: listDcGatewayCcnRoutes,
			Tags:    map[string]string{"service": "dc", "action": "DescribeDirectConnectGatewayCcnRoutes"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "route_id", Type: proto.ColumnType_STRING, Description: "The route ID.", Transform: transform.FromField("RouteId")},
			{Name: "direct_connect_gateway_id", Type: proto.ColumnType_STRING, Description: "The parent Direct Connect gateway ID that owns the route.", Transform: transform.FromField("DirectConnectGatewayId")},
			{Name: "destination_cidr_block", Type: proto.ColumnType_STRING, Description: "The IDC IP range (CIDR block).", Transform: transform.FromField("DestinationCidrBlock")},
			{Name: "as_path", Type: proto.ColumnType_JSON, Description: "The AS-Path attribute of BGP.", Transform: transform.FromField("ASPath")},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The remarks of the route.", Transform: transform.FromField("Description")},
			{Name: "ccn_route_type", Type: proto.ColumnType_STRING, Description: "The route-learning type (BGP, STATIC).", Transform: transform.FromField("CcnRouteType")},
			{Name: "update_time", Type: proto.ColumnType_TIMESTAMP, Description: "The last updated time of the route.", Transform: transform.FromField("UpdateTime").NullIfZero()},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the route resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("RouteId")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listDcGatewayCcnRoutes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listDcGatewayCcnRoutes", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listDcGatewayCcnRoutes", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)
	ccnRouteType := d.EqualsQualString("ccn_route_type")

	gatewayIds, err := resolveDcGatewayIds(ctx, d, client)
	if err != nil {
		plugin.Logger(ctx).Error("listDcGatewayCcnRoutes", "resolve_gateway_ids_error", err)
		return nil, err
	}

	for _, gatewayId := range gatewayIds {
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
		if gatewayId == "" {
			continue
		}

		if err := listDcGatewayCcnRoutesForGateway(ctx, d, client, gatewayId, ccnRouteType, region); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func resolveDcGatewayIds(ctx context.Context, d *plugin.QueryData, client *vpc.Client) ([]string, error) {
	if id := d.EqualsQualString("direct_connect_gateway_id"); id != "" {
		return []string{id}, nil
	}

	pageSize := uint64(100)
	var offset uint64 = 0
	var ids []string

	for {
		d.WaitForListRateLimit(ctx)

		req := vpc.NewDescribeDirectConnectGatewaysRequest()
		req.Offset = &offset
		req.Limit = &pageSize

		utils.LogRequest(ctx, "resolveDcGatewayIds", req)
		resp, err := client.DescribeDirectConnectGatewaysWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("resolveDcGatewayIds", "request_error", err)
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
			if id := utils.PtrString(item.DirectConnectGatewayId); id != "" {
				ids = append(ids, id)
			}
		}

		if utils.PageDone(int64(offset), int64(len(items)), total) {
			break
		}
		offset += uint64(len(items))
	}
	return ids, nil
}

func listDcGatewayCcnRoutesForGateway(ctx context.Context, d *plugin.QueryData, client *vpc.Client, gatewayId, ccnRouteType, region string) error {
	pageSize := uint64(utils.PageSizeFromLimit(d, 100, 1))
	var offset uint64 = 0

	for {
		d.WaitForListRateLimit(ctx)

		req := vpc.NewDescribeDirectConnectGatewayCcnRoutesRequest()
		req.DirectConnectGatewayId = common.StringPtr(gatewayId)
		req.Offset = &offset
		req.Limit = &pageSize
		if ccnRouteType != "" {
			req.CcnRouteType = common.StringPtr(ccnRouteType)
		}

		utils.LogRequest(ctx, "listDcGatewayCcnRoutes", req)
		resp, err := client.DescribeDirectConnectGatewayCcnRoutesWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listDcGatewayCcnRoutes", "request_error", err, "gateway_id", gatewayId)
			return err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		items := resp.Response.RouteSet
		total := utils.PtrUint64(resp.Response.TotalCount)

		for _, item := range items {
			if item == nil {
				continue
			}
			row := toDcGatewayCcnRouteRow(item, gatewayId, ccnRouteType, region)
			d.StreamListItem(ctx, row)
			if d.RowsRemaining(ctx) == 0 {
				return nil
			}
		}

		if utils.PageDone(int64(offset), int64(len(items)), total) {
			break
		}
		offset += uint64(len(items))
	}
	return nil
}

// Row Type

type dcGatewayCcnRouteRow struct {
	RouteId                string
	DirectConnectGatewayId string
	DestinationCidrBlock   string
	ASPath                 []*string
	Description            string
	CcnRouteType           string
	UpdateTime             *time.Time
	Region                 string
	Akas                   []string
}

func toDcGatewayCcnRouteRow(s *vpc.DirectConnectGatewayCcnRoute, gatewayId, ccnRouteType, region string) dcGatewayCcnRouteRow {
	row := dcGatewayCcnRouteRow{
		RouteId:                utils.PtrString(s.RouteId),
		DirectConnectGatewayId: gatewayId,
		DestinationCidrBlock:   utils.PtrString(s.DestinationCidrBlock),
		ASPath:                 s.ASPath,
		Description:            utils.PtrString(s.Description),
		CcnRouteType:           ccnRouteType,
		UpdateTime:             utils.ParseTimestamp(s.UpdateTime),
		Region:                 region,
	}

	if row.RouteId != "" {
		row.Akas = []string{"tencentcloud:dc:gateway-ccn-route:" + gatewayId + "/" + row.RouteId}
	}
	return row
}
