package ccn

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

func TableTencentcloudCcnRoute() *plugin.Table {
	return &plugin.Table{
		Name:        "tencentcloud_ccn_route",
		Description: "Tencent Cloud CCN routing policies that route traffic between associated instances of a Cloud Connect Network.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "ccn_id", Require: plugin.Optional},
				{Name: "route_id", Require: plugin.Optional},
			},
			Hydrate: listCcnRoute,
			Tags:    map[string]string{"service": "ccn", "action": "DescribeCcnRoutes"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "ccn_id", Type: proto.ColumnType_STRING, Description: "The CCN instance ID that the route belongs to.", Transform: transform.FromField("CcnId")},
			{Name: "route_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the routing policy, such as `ccnr-f49l6u0z`.", Transform: transform.FromField("RouteId")},
			{Name: "destination_cidr_block", Type: proto.ColumnType_STRING, Description: "The destination CIDR block of the route.", Transform: transform.FromField("DestinationCidrBlock")},
			{Name: "instance_type", Type: proto.ColumnType_STRING, Description: "The next hop (associated instance) type. `VPC` or `DIRECTCONNECT`.", Transform: transform.FromField("InstanceType")},
			{Name: "instance_id", Type: proto.ColumnType_STRING, Description: "The next hop (associated instance) ID.", Transform: transform.FromField("InstanceId")},
			{Name: "instance_name", Type: proto.ColumnType_STRING, Description: "The next hop (associated instance) name.", Transform: transform.FromField("InstanceName")},
			{Name: "instance_region", Type: proto.ColumnType_STRING, Description: "The region of the next hop (associated instance).", Transform: transform.FromField("InstanceRegion")},
			{Name: "instance_uin", Type: proto.ColumnType_STRING, Description: "The UIN (root account) to which the associated instance belongs.", Transform: transform.FromField("InstanceUin")},
			{Name: "instance_extra_name", Type: proto.ColumnType_STRING, Description: "The next hop port name (associated instance's port name).", Transform: transform.FromField("InstanceExtraName")},
			{Name: "enabled", Type: proto.ColumnType_BOOL, Description: "Whether the route is enabled.", Transform: transform.FromField("Enabled")},
			{Name: "extra_state", Type: proto.ColumnType_STRING, Description: "The additional status of the route.", Transform: transform.FromField("ExtraState")},
			{Name: "is_bgp", Type: proto.ColumnType_BOOL, Description: "Whether it is a dynamic (BGP) route.", Transform: transform.FromField("IsBgp")},
			{Name: "route_priority", Type: proto.ColumnType_INT, Description: "The priority of the route.", Transform: transform.FromField("RoutePriority")},
			{Name: "update_time", Type: proto.ColumnType_TIMESTAMP, Description: "The last update time of the route.", Transform: transform.FromField("UpdateTime").NullIfZero()},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listCcnRoute(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listCcnRoute", "quals", d.EqualsQuals)

	client := &vpc.Client{}
	if err := utils.InitClientInRegion(ctx, d, client, ccnRequestRegion(d)); err != nil {
		plugin.Logger(ctx).Error("listCcnRoute", "client_init_error", err)
		return nil, err
	}

	ccnIds, err := resolveCcnIDs(ctx, d, client)
	if err != nil {
		plugin.Logger(ctx).Error("listCcnRoute", "resolve_ccn_ids_error", err)
		return nil, err
	}

	routeIds := utils.SingleIdFromQual(d.EqualsQuals, "route_id")

	for _, ccnId := range ccnIds {
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}

		pageSize := uint64(utils.PageSizeFromLimit(d, 100, 1))
		var offset uint64

		for {
			d.WaitForListRateLimit(ctx)

			req := vpc.NewDescribeCcnRoutesRequest()
			req.CcnId = &ccnId
			req.Offset = &offset
			req.Limit = &pageSize
			if len(routeIds) > 0 {
				req.RouteIds = routeIds
			}

			utils.LogRequest(ctx, "listCcnRoute", req)
			resp, err := client.DescribeCcnRoutesWithContext(ctx, req)
			if err != nil {
				plugin.Logger(ctx).Error("listCcnRoute", "request_error", err)
				return nil, err
			}
			if resp == nil || resp.Response == nil {
				break
			}

			routes := resp.Response.RouteSet
			if resp.Response == nil {
				break
			}
			total := utils.PtrUint64(resp.Response.TotalCount)

			for _, item := range routes {
				if item == nil {
					continue
				}
				row := toCcnRouteRow(item, ccnId)
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

type ccnRouteRow struct {
	CcnId                string
	RouteId              string
	DestinationCidrBlock string
	InstanceType         string
	InstanceId           string
	InstanceName         string
	InstanceRegion       string
	InstanceUin          string
	InstanceExtraName    string
	Enabled              bool
	ExtraState           string
	IsBgp                bool
	RoutePriority        uint64
	UpdateTime           *time.Time
	Title                string
	Akas                 []string
}

func toCcnRouteRow(r *vpc.CcnRoute, ccnId string) ccnRouteRow {
	row := ccnRouteRow{
		CcnId:                ccnId,
		RouteId:              utils.PtrString(r.RouteId),
		DestinationCidrBlock: utils.PtrString(r.DestinationCidrBlock),
		InstanceType:         utils.PtrString(r.InstanceType),
		InstanceId:           utils.PtrString(r.InstanceId),
		InstanceName:         utils.PtrString(r.InstanceName),
		InstanceRegion:       utils.PtrString(r.InstanceRegion),
		InstanceUin:          utils.PtrString(r.InstanceUin),
		InstanceExtraName:    utils.PtrString(r.InstanceExtraName),
		Enabled:              utils.PtrBool(r.Enabled),
		ExtraState:           utils.PtrString(r.ExtraState),
		IsBgp:                utils.PtrBool(r.IsBgp),
		RoutePriority:        utils.PtrUint64(r.RoutePriority),
		UpdateTime:           utils.ParseTimestamp(r.UpdateTime),
	}

	row.Title = row.RouteId
	if row.Title == "" {
		row.Title = row.DestinationCidrBlock
	}

	if row.CcnId != "" && row.RouteId != "" {
		row.Akas = []string{"tencentcloud:ccn:route:" + row.CcnId + "/" + row.RouteId}
	}
	return row
}
