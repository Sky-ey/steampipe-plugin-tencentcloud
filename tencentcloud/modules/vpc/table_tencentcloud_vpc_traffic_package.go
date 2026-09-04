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

func TableTencentcloudVpcTrafficPackage() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_vpc_traffic_package",
		Description:       "Tencent Cloud VPC traffic packages that provide prepaid data transfer allowances for public network traffic.",
		GetMatrixItemFunc: buildVpcRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "traffic_package_id", Require: plugin.Optional},
				{Name: "traffic_package_name", Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
			},
			Hydrate: listVpcTrafficPackages,
			Tags:    map[string]string{"service": "vpc", "action": "DescribeTrafficPackages"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("traffic_package_id"),
			Hydrate:    getVpcTrafficPackage,
			Tags:       map[string]string{"service": "vpc", "action": "DescribeTrafficPackages"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "traffic_package_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the traffic package.", Transform: transform.FromField("TrafficPackageId")},
			{Name: "traffic_package_name", Type: proto.ColumnType_STRING, Description: "The name of the traffic package.", Transform: transform.FromField("TrafficPackageName")},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of the traffic package (AVAILABLE, EXPIRED, EXHAUSTED, REFUNDED, DELETED).", Transform: transform.FromField("Status")},
			{Name: "deduct_type", Type: proto.ColumnType_STRING, Description: "The deduction type of the traffic package (idle-time, full-time).", Transform: transform.FromField("DeductType")},
			{Name: "total_amount", Type: proto.ColumnType_DOUBLE, Description: "The total amount of the traffic package in GB.", Transform: transform.FromField("TotalAmount")},
			{Name: "used_amount", Type: proto.ColumnType_DOUBLE, Description: "The used amount of the traffic package in GB.", Transform: transform.FromField("UsedAmount")},
			{Name: "remaining_amount", Type: proto.ColumnType_DOUBLE, Description: "The remaining amount of the traffic package in GB.", Transform: transform.FromField("RemainingAmount")},
			{Name: "created_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the traffic package (ISO 8601, UTC).", Transform: transform.FromField("CreatedTime").NullIfZero()},
			{Name: "deadline", Type: proto.ColumnType_TIMESTAMP, Description: "The expiration time of the traffic package (ISO 8601, UTC).", Transform: transform.FromField("Deadline").NullIfZero()},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The tags attached to the traffic package as a map of key -> value.", Transform: transform.FromField("Tags")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the traffic package resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("TrafficPackageName")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listVpcTrafficPackages(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listVpcTrafficPackages", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listVpcTrafficPackages", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	pageSize := uint64(utils.PageSizeFromLimit(d, 100, 1))
	var offset uint64

	trafficPackageIds := utils.SingleIdFromQual(d.EqualsQuals, "traffic_package_id")
	filters := buildVpcTrafficPackageFilters(d.EqualsQuals)

	for {
		d.WaitForListRateLimit(ctx)

		req := vpc.NewDescribeTrafficPackagesRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		if len(trafficPackageIds) > 0 {
			req.TrafficPackageIds = trafficPackageIds
		} else if len(filters) > 0 {
			req.Filters = filters
		}

		utils.LogRequest(ctx, "listVpcTrafficPackages", req)
		resp, err := client.DescribeTrafficPackagesWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listVpcTrafficPackages", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		packages := resp.Response.TrafficPackageSet
		total := utils.PtrInt64(resp.Response.TotalCount)

		for _, item := range packages {
			if item == nil {
				continue
			}
			row := toVpcTrafficPackageRow(item)
			row.Region = region
			d.StreamListItem(ctx, row)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if utils.PageDone(int64(offset), int64(len(packages)), total) {
			break
		}
		offset += uint64(len(packages))
	}

	return nil, nil
}

func buildVpcTrafficPackageFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*vpc.Filter {
	return utils.BuildFilters(equalsQuals, newVpcFilter,
		utils.FilterMapping{QualName: "traffic_package_name", FilterName: "traffic-package-name"},
		utils.FilterMapping{QualName: "status", FilterName: "status"},
	)
}

// Get Function

func getVpcTrafficPackage(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	trafficPackageId := d.EqualsQuals["traffic_package_id"].GetStringValue()
	plugin.Logger(ctx).Info("getVpcTrafficPackage", "traffic_package_id", trafficPackageId, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getVpcTrafficPackage", "client_init_error", err)
		return nil, err
	}

	if trafficPackageId == "" {
		return nil, nil
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	req := vpc.NewDescribeTrafficPackagesRequest()
	req.TrafficPackageIds = common.StringPtrs([]string{trafficPackageId})

	utils.LogRequest(ctx, "getVpcTrafficPackage", req)
	resp, err := client.DescribeTrafficPackagesWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getVpcTrafficPackage", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.TrafficPackageSet) == 0 {
		return nil, nil
	}

	row := toVpcTrafficPackageRow(resp.Response.TrafficPackageSet[0])
	row.Region = region
	return row, nil
}

// Row Type

type vpcTrafficPackageRow struct {
	TrafficPackageId   string
	TrafficPackageName string
	Status             string
	DeductType         string
	TotalAmount        float64
	UsedAmount         float64
	RemainingAmount    float64
	CreatedTime        *time.Time
	Deadline           *time.Time
	Tags               map[string]string
	Region             string
	Akas               []string
}

func toVpcTrafficPackageRow(t *vpc.TrafficPackage) vpcTrafficPackageRow {
	row := vpcTrafficPackageRow{
		TrafficPackageId:   utils.PtrString(t.TrafficPackageId),
		TrafficPackageName: utils.PtrString(t.TrafficPackageName),
		Status:             utils.PtrString(t.Status),
		DeductType:         utils.PtrString(t.DeductType),
		TotalAmount:        utils.PtrFloat64(t.TotalAmount),
		UsedAmount:         utils.PtrFloat64(t.UsedAmount),
		RemainingAmount:    utils.PtrFloat64(t.RemainingAmount),
		CreatedTime:        utils.ParseTimestamp(t.CreatedTime),
		Deadline:           utils.ParseTimestamp(t.Deadline),
		Tags:               TagsToMap(t.TagSet),
	}

	if row.TrafficPackageId != "" {
		row.Akas = []string{"tencentcloud:vpc:traffic-package:" + row.TrafficPackageId}
	}
	return row
}
