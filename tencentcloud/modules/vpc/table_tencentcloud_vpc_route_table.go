package vpc

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"strconv"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/vpc/v20170312"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudVpcRouteTable() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_vpc_route_table",
		Description:       "Tencent Cloud VPC route tables that control how traffic is forwarded within a VPC.",
		GetMatrixItemFunc: buildVpcRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "route_table_id", Require: plugin.Optional},
				{Name: "route_table_name", Require: plugin.Optional},
				{Name: "vpc_id", Require: plugin.Optional},
				{Name: "main", Require: plugin.Optional, Operators: []string{"="}},
			},
			Hydrate: listVpcRouteTables,
			Tags:    map[string]string{"service": "vpc", "action": "DescribeRouteTables"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("route_table_id"),
			Hydrate:    getVpcRouteTable,
			Tags:       map[string]string{"service": "vpc", "action": "DescribeRouteTables"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "route_table_id", Type: proto.ColumnType_STRING, Description: "The route table instance ID, such as `rtb-azd4dt1c`.", Transform: transform.FromField("RouteTableId")},
			{Name: "route_table_name", Type: proto.ColumnType_STRING, Description: "The route table name.", Transform: transform.FromField("RouteTableName")},
			{Name: "vpc_id", Type: proto.ColumnType_STRING, Description: "The VPC instance ID to which the route table belongs.", Transform: transform.FromField("VpcId")},
			{Name: "main", Type: proto.ColumnType_BOOL, Description: "Whether it is the default (main) route table of the VPC.", Transform: transform.FromField("Main")},
			{Name: "association_set", Type: proto.ColumnType_JSON, Description: "The associations between the route table and subnets.", Transform: transform.FromField("AssociationSet")},
			{Name: "route_set", Type: proto.ColumnType_JSON, Description: "The routing policy entries of the route table.", Transform: transform.FromField("RouteSet")},
			{Name: "local_cidr_for_ccn", Type: proto.ColumnType_JSON, Description: "Local CIDR blocks published to CCN.", Transform: transform.FromField("LocalCidrForCcn")},
			{Name: "created_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the route table.", Transform: transform.FromField("CreatedTime").NullIfZero()},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The tags attached to the route table as a map of key -> value.", Transform: transform.FromField("Tags")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the route table resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("RouteTableName")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listVpcRouteTables(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listVpcRouteTables", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listVpcRouteTables", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	pageSize := strconv.FormatInt(utils.PageSizeFromLimit(d, 100, 1), 10)
	offset := "0"

	routeTableIds := utils.SingleIdFromQual(d.EqualsQuals, "route_table_id")
	filters := buildVpcRouteTableFilters(d.EqualsQuals)

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

		utils.LogRequest(ctx, "listVpcRouteTables", req)
		resp, err := client.DescribeRouteTablesWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listVpcRouteTables", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		rts := resp.Response.RouteTableSet
		total := utils.PtrUint64(resp.Response.TotalCount)
		currentOffset, _ := strconv.ParseUint(offset, 10, 64)

		for _, item := range rts {
			if item == nil {
				continue
			}
			row := toVpcRouteTableRow(item)
			row.Region = region
			d.StreamListItem(ctx, row)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if utils.PageDone(int64(currentOffset), int64(len(rts)), total) {
			break
		}
		offset = strconv.FormatUint(currentOffset+uint64(len(rts)), 10)
	}

	return nil, nil
}

func buildVpcRouteTableFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*vpc.Filter {
	filters := utils.BuildFilters(equalsQuals, newVpcFilter,
		utils.FilterMapping{QualName: "route_table_name", FilterName: "route-table-name"},
		utils.FilterMapping{QualName: "vpc_id", FilterName: "vpc-id"},
	)

	if equalsQuals["main"] != nil {
		value := equalsQuals["main"].GetBoolValue()
		filters = append(filters, newVpcFilter("association.main", strconv.FormatBool(value)))
	}

	return filters
}

// Get Function

func getVpcRouteTable(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	routeTableId := d.EqualsQuals["route_table_id"].GetStringValue()
	plugin.Logger(ctx).Info("getVpcRouteTable", "route_table_id", routeTableId, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getVpcRouteTable", "client_init_error", err)
		return nil, err
	}

	if routeTableId == "" {
		return nil, nil
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	req := vpc.NewDescribeRouteTablesRequest()
	req.RouteTableIds = common.StringPtrs([]string{routeTableId})

	utils.LogRequest(ctx, "getVpcRouteTable", req)
	resp, err := client.DescribeRouteTablesWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getVpcRouteTable", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.RouteTableSet) == 0 {
		return nil, nil
	}

	row := toVpcRouteTableRow(resp.Response.RouteTableSet[0])
	row.Region = region
	return row, nil
}

// Row Type

type vpcRouteTableRow struct {
	RouteTableId    string
	RouteTableName  string
	VpcId           string
	Main            bool
	AssociationSet  []*vpc.RouteTableAssociation
	RouteSet        []*vpc.Route
	LocalCidrForCcn []*vpc.CidrForCcn
	CreatedTime     *time.Time
	Tags            map[string]string
	Region          string
	Akas            []string
}

func toVpcRouteTableRow(rt *vpc.RouteTable) vpcRouteTableRow {
	row := vpcRouteTableRow{
		RouteTableId:    utils.PtrString(rt.RouteTableId),
		RouteTableName:  utils.PtrString(rt.RouteTableName),
		VpcId:           utils.PtrString(rt.VpcId),
		Main:            utils.PtrBool(rt.Main),
		AssociationSet:  rt.AssociationSet,
		RouteSet:        rt.RouteSet,
		LocalCidrForCcn: rt.LocalCidrForCcn,
		CreatedTime:     utils.ParseTimestamp(rt.CreatedTime),
		Tags:            TagsToMap(rt.TagSet),
	}

	if row.RouteTableId != "" {
		row.Akas = []string{"tencentcloud:vpc:route-table:" + row.RouteTableId}
	}
	return row
}
