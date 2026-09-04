package vpc

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"strconv"
	"time"

	vpc "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/vpc/v20170312"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudVpcRoute() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_vpc_route",
		Description:       "Tencent Cloud VPC routing entries, expanded one row per route policy from each route table.",
		GetMatrixItemFunc: buildVpcRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "route_table_id", Require: plugin.Optional},
				{Name: "vpc_id", Require: plugin.Optional},
			},
			Hydrate: listVpcRoutes,
			Tags:    map[string]string{"service": "vpc", "action": "DescribeRouteTables"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "route_table_id", Type: proto.ColumnType_STRING, Description: "The route table instance ID to which the route belongs.", Transform: transform.FromField("RouteTableId")},
			{Name: "vpc_id", Type: proto.ColumnType_STRING, Description: "The VPC instance ID to which the route belongs.", Transform: transform.FromField("VpcId")},
			{Name: "route_id", Type: proto.ColumnType_INT, Description: "The routing policy ID. IPv6 routes always return 0; prefer route_item_id as the unique key.", Transform: transform.FromField("RouteId")},
			{Name: "route_item_id", Type: proto.ColumnType_STRING, Description: "The unique routing policy ID.", Transform: transform.FromField("RouteItemId")},
			{Name: "destination_cidr_block", Type: proto.ColumnType_STRING, Description: "The destination IPv4 CIDR block, such as `112.20.51.0/24`.", Transform: transform.FromField("DestinationCidrBlock")},
			{Name: "destination_ipv6_cidr_block", Type: proto.ColumnType_STRING, Description: "The destination IPv6 CIDR block.", Transform: transform.FromField("DestinationIpv6CidrBlock")},
			{Name: "gateway_type", Type: proto.ColumnType_STRING, Description: "Type of the next hop (CVM, VPN, DIRECTCONNECT, PEERCONNECTION, HAVIP, NAT, NORMAL_CVM, EIP, LOCAL_GATEWAY).", Transform: transform.FromField("GatewayType")},
			{Name: "gateway_id", Type: proto.ColumnType_STRING, Description: "The next hop address (gateway ID or, for NORMAL_CVM, the private IP of the instance).", Transform: transform.FromField("GatewayId")},
			{Name: "route_type", Type: proto.ColumnType_STRING, Description: "The route type (USER, NETD, CCN).", Transform: transform.FromField("RouteType")},
			{Name: "route_description", Type: proto.ColumnType_STRING, Description: "The description of the routing policy.", Transform: transform.FromField("RouteDescription")},
			{Name: "enabled", Type: proto.ColumnType_BOOL, Description: "Whether the routing policy is enabled.", Transform: transform.FromField("Enabled")},
			{Name: "published_to_vbc", Type: proto.ColumnType_BOOL, Description: "Whether the routing policy is published to CCN.", Transform: transform.FromField("PublishedToVbc")},
			{Name: "created_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the routing policy.", Transform: transform.FromField("CreatedTime").NullIfZero()},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the route resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listVpcRoutes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listVpcRoutes", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listVpcRoutes", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	pageSize := strconv.FormatInt(utils.PageSizeFromLimit(d, 100, 1), 10)
	offset := "0"

	routeTableIds := utils.SingleIdFromQual(d.EqualsQuals, "route_table_id")
	filters := utils.BuildFilters(d.EqualsQuals, newVpcFilter,
		utils.FilterMapping{QualName: "vpc_id", FilterName: "vpc-id"},
	)

	for {
		d.WaitForListRateLimit(ctx)

		req := vpc.NewDescribeRouteTablesRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		if len(routeTableIds) > 0 {
			req.RouteTableIds = routeTableIds
		} else if len(filters) > 0 {
			req.Filters = filters
		}

		utils.LogRequest(ctx, "listVpcRoutes", req)
		resp, err := client.DescribeRouteTablesWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listVpcRoutes", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		rts := resp.Response.RouteTableSet
		total := utils.PtrUint64(resp.Response.TotalCount)
		currentOffset, _ := strconv.ParseUint(offset, 10, 64)

		for _, rt := range rts {
			if rt == nil {
				continue
			}
			vpcId := utils.PtrString(rt.VpcId)
			routeTableId := utils.PtrString(rt.RouteTableId)
			for _, route := range rt.RouteSet {
				if route == nil {
					continue
				}
				row := toVpcRouteRow(route, routeTableId, vpcId, region)
				d.StreamListItem(ctx, row)
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}

		if utils.PageDone(int64(currentOffset), int64(len(rts)), total) {
			break
		}
		offset = strconv.FormatUint(currentOffset+uint64(len(rts)), 10)
	}

	return nil, nil
}

// Row Type

type vpcRouteRow struct {
	RouteTableId             string
	VpcId                    string
	RouteId                  uint64
	RouteItemId              string
	DestinationCidrBlock     string
	DestinationIpv6CidrBlock string
	GatewayType              string
	GatewayId                string
	RouteType                string
	RouteDescription         string
	Enabled                  bool
	PublishedToVbc           bool
	CreatedTime              *time.Time
	Region                   string
	Title                    string
	Akas                     []string
}

func toVpcRouteRow(r *vpc.Route, routeTableId, vpcId, region string) vpcRouteRow {
	row := vpcRouteRow{
		RouteTableId:             routeTableId,
		VpcId:                    vpcId,
		RouteId:                  utils.PtrUint64(r.RouteId),
		RouteItemId:              utils.PtrString(r.RouteItemId),
		DestinationCidrBlock:     utils.PtrString(r.DestinationCidrBlock),
		DestinationIpv6CidrBlock: utils.PtrString(r.DestinationIpv6CidrBlock),
		GatewayType:              utils.PtrString(r.GatewayType),
		GatewayId:                utils.PtrString(r.GatewayId),
		RouteType:                utils.PtrString(r.RouteType),
		RouteDescription:         utils.PtrString(r.RouteDescription),
		Enabled:                  utils.PtrBool(r.Enabled),
		PublishedToVbc:           utils.PtrBool(r.PublishedToVbc),
		CreatedTime:              utils.ParseTimestamp(r.CreatedTime),
		Region:                   region,
	}

	row.Title = row.DestinationCidrBlock
	if row.Title == "" {
		row.Title = row.DestinationIpv6CidrBlock
	}
	if row.Title == "" {
		row.Title = row.RouteItemId
	}

	identity := row.RouteItemId
	if identity == "" {
		identity = strconv.FormatUint(row.RouteId, 10)
	}
	if identity != "" {
		row.Akas = []string{"tencentcloud:vpc:route:" + routeTableId + "/" + identity}
	}
	return row
}
