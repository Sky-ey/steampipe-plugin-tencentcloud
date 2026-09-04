package clb

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"

	clb "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/clb/v20180317"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudClbCrossTarget() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_clb_cross_target",
		Description:       "Tencent Cloud Load Balancer (CLB) cross-VPC backend targets available via CCN cross-domain 2.0.",
		GetMatrixItemFunc: buildClbRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "vpc_id", Require: plugin.Optional},
				{Name: "ip", Require: plugin.Optional},
				{Name: "listener_id", Require: plugin.Optional},
				{Name: "location_id", Require: plugin.Optional},
			},
			Hydrate: listClbCrossTargets,
			Tags:    map[string]string{"service": "clb", "action": "DescribeCrossTargets"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "local_vpc_id", Type: proto.ColumnType_STRING, Description: "The VPC ID of the CLB instance.", Transform: transform.FromField("LocalVpcId")},
			{Name: "vpc_id", Type: proto.ColumnType_STRING, Description: "The VPC ID of the CVM or ENI backend instance.", Transform: transform.FromField("VpcId")},
			{Name: "vpc_name", Type: proto.ColumnType_STRING, Description: "The VPC name of the CVM or ENI backend instance."},
			{Name: "ip", Type: proto.ColumnType_IPADDR, Description: "The IP address of the CVM or ENI backend instance.", Transform: transform.FromField("IP")},
			{Name: "instance_id", Type: proto.ColumnType_STRING, Description: "The ID of the backend instance.", Transform: transform.FromField("InstanceId")},
			{Name: "instance_name", Type: proto.ColumnType_STRING, Description: "The name of the backend instance."},
			{Name: "eni_id", Type: proto.ColumnType_STRING, Description: "The ENI ID of the backend instance.", Transform: transform.FromField("EniId")},
			{Name: "target_region", Type: proto.ColumnType_STRING, Description: "The region where the backend CVM or ENI instance resides.", Transform: transform.FromField("Region")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region of the load balancer control plane used for the query.", Transform: transform.FromField("QueryRegion")},
			{Name: "listener_id", Type: proto.ColumnType_STRING, Description: "Filter-only qual: the CLB listener ID used to scope the cross-target query (passed through to DescribeCrossTargets filter listener-id). Echoes the qual value when provided; not part of the API response payload.", Transform: transform.FromField("ListenerId")},
			{Name: "location_id", Type: proto.ColumnType_STRING, Description: "Filter-only qual: the CLB 7-layer location ID used to scope the cross-target query (passed through to DescribeCrossTargets filter location-id). Echoes the qual value when provided; not part of the API response payload.", Transform: transform.FromField("LocationId")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource."},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource."},
		}),
	}
}

// List Function

func listClbCrossTargets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listClbCrossTargets", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &clb.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listClbCrossTargets", "client_init_error", err)
		return nil, err
	}
	region := d.EqualsQualString(utils.MatrixKeyRegion)
	listenerID := d.EqualsQualString("listener_id")
	locationID := d.EqualsQualString("location_id")

	// DescribeCrossTargets 的各 filter (vpc-id/ip/listener-id/location-id) 均为可选,
	// 但传入空 Filters 数组 `[]` 会被服务端拒绝 (InvalidParameterValue)。
	// 仅当至少有一个 filter 时才赋值;为空时保持 nil,SDK 的 omitnil 标签会省略该字段。
	filters := utils.BuildFilters(d.EqualsQuals, newClbFilter,
		utils.FilterMapping{QualName: "vpc_id", FilterName: "vpc-id"},
		utils.FilterMapping{QualName: "ip", FilterName: "ip"},
		utils.FilterMapping{QualName: "listener_id", FilterName: "listener-id"},
		utils.FilterMapping{QualName: "location_id", FilterName: "location-id"},
	)
	pageSize := uint64(utils.PageSizeFromLimit(d, 100, 1))
	var offset uint64

	for {
		d.WaitForListRateLimit(ctx)
		req := clb.NewDescribeCrossTargetsRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		if len(filters) > 0 {
			req.Filters = filters
		}

		utils.LogRequest(ctx, "listClbCrossTargets", req)
		resp, err := client.DescribeCrossTargetsWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listClbCrossTargets", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		crossTargets := resp.Response.CrossTargetSet
		for _, crossTarget := range crossTargets {
			if crossTarget == nil {
				continue
			}
			row := toClbCrossTargetRow(crossTarget, region, listenerID, locationID)
			d.StreamListItem(ctx, row)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if utils.PageDone(int64(offset), int64(len(crossTargets)), utils.PtrUint64(resp.Response.TotalCount)) {
			break
		}
		offset += uint64(len(crossTargets))
	}
	return nil, nil
}

// Row Type

type clbCrossTargetRow struct {
	LocalVpcId   string
	VpcId        string
	VpcName      string
	IP           string
	InstanceId   string
	InstanceName string
	EniId        string
	Region       string
	QueryRegion  string
	ListenerId   string
	LocationId   string
	Title        string
	Akas         []string
}

func toClbCrossTargetRow(crossTarget *clb.CrossTargets, queryRegion, listenerID, locationID string) clbCrossTargetRow {
	row := clbCrossTargetRow{
		LocalVpcId:   utils.PtrString(crossTarget.LocalVpcId),
		VpcId:        utils.PtrString(crossTarget.VpcId),
		VpcName:      utils.PtrString(crossTarget.VpcName),
		IP:           utils.PtrString(crossTarget.IP),
		InstanceId:   utils.PtrString(crossTarget.InstanceId),
		InstanceName: utils.PtrString(crossTarget.InstanceName),
		EniId:        utils.PtrString(crossTarget.EniId),
		Region:       utils.PtrString(crossTarget.Region),
		QueryRegion:  queryRegion,
		ListenerId:   listenerID,
		LocationId:   locationID,
	}
	row.Title = row.InstanceName
	if row.Title == "" {
		row.Title = row.InstanceId
	}
	if row.Title == "" {
		row.Title = row.IP
	}
	identity := row.InstanceId
	if identity == "" {
		identity = row.EniId
	}
	if identity == "" {
		identity = row.IP
	}
	if identity != "" {
		row.Akas = []string{"tencentcloud:clb:cross-target:" + row.LocalVpcId + "/" + row.VpcId + "/" + identity}
	}
	return row
}
