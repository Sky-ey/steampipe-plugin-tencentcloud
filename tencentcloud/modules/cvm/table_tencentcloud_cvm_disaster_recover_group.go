package cvm

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/cvm/v20170312"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudCvmDisasterRecoverGroup() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_cvm_disaster_recover_group",
		Description:       "Tencent Cloud CVM spread placement groups distribute instances across different physical hardware to reduce the impact of hardware failures.",
		GetMatrixItemFunc: buildCvmRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "disaster_recover_group_id", Require: plugin.Optional},
				{Name: "name", Require: plugin.Optional},
			},
			Hydrate: listCvmDisasterRecoverGroups,
			Tags:    map[string]string{"service": "cvm", "action": "DescribeDisasterRecoverGroups"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("disaster_recover_group_id"),
			Hydrate:    getCvmDisasterRecoverGroup,
			Tags:       map[string]string{"service": "cvm", "action": "DescribeDisasterRecoverGroups"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "disaster_recover_group_id", Type: proto.ColumnType_STRING, Description: "The ID of the spread placement group.", Transform: transform.FromField("DisasterRecoverGroupId")},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the spread placement group (1-60 characters).", Transform: transform.FromField("Name")},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The type of the spread placement group (HOST, SW, RACK).", Transform: transform.FromField("Type")},
			{Name: "cvm_quota_total", Type: proto.ColumnType_INT, Description: "The maximum number of CVMs that can be hosted in the spread placement group.", Transform: transform.FromField("CvmQuotaTotal")},
			{Name: "current_num", Type: proto.ColumnType_INT, Description: "The current number of CVMs in the spread placement group.", Transform: transform.FromField("CurrentNum")},
			{Name: "instance_ids", Type: proto.ColumnType_JSON, Description: "The list of CVM IDs in the spread placement group.", Transform: transform.FromField("InstanceIds")},
			{Name: "create_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the spread placement group.", Transform: transform.FromField("CreateTime").NullIfZero()},
			{Name: "affinity", Type: proto.ColumnType_INT, Description: "The placement group affinity.", Transform: transform.FromField("Affinity")},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The tags attached to the spread placement group as a map of key -> value.", Transform: transform.FromField("Tags")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the spread placement group resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Name")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listCvmDisasterRecoverGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Info("listCvmDisasterRecoverGroups", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &cvm.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listCvmDisasterRecoverGroups", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	// Mode B: ID-first. DisasterRecoverGroupIds and Filters are mutually exclusive.
	groupIds := disasterRecoverGroupIdsFromQuals(d.EqualsQuals)

	pageSize := utils.PageSizeFromLimit(d, 100, 1)
	var offset int64 = 0

	for {
		d.WaitForListRateLimit(ctx)

		req := cvm.NewDescribeDisasterRecoverGroupsRequest()
		req.Offset = &offset
		req.Limit = &pageSize

		if len(groupIds) > 0 {
			req.DisasterRecoverGroupIds = common.StringPtrs(groupIds)
		} else {
			if name := d.EqualsQualString("name"); name != "" {
				req.Name = common.StringPtr(name)
			}
		}

		utils.LogRequest(ctx, "listCvmDisasterRecoverGroups", req)
		resp, err := client.DescribeDisasterRecoverGroupsWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listCvmDisasterRecoverGroups", "request_error", err)
			return nil, err
		}

		if resp == nil || resp.Response == nil {
			break
		}

		groups := resp.Response.DisasterRecoverGroupSet
		for _, g := range groups {
			if g == nil {
				continue
			}
			row := toCvmDisasterRecoverGroupRow(g)
			row.Region = region
			d.StreamListItem(ctx, row)

			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		total := utils.PtrInt64(resp.Response.TotalCount)
		if utils.PageDone(offset, int64(len(groups)), total) {
			break
		}
		offset += int64(len(groups))
	}

	return nil, nil
}

func disasterRecoverGroupIdsFromQuals(equalsQuals plugin.KeyColumnEqualsQualMap) []string {
	if equalsQuals["disaster_recover_group_id"] == nil {
		return nil
	}
	id := equalsQuals["disaster_recover_group_id"].GetStringValue()
	if id == "" {
		return nil
	}
	return []string{id}
}

// Get Function

func getCvmDisasterRecoverGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	groupId := d.EqualsQuals["disaster_recover_group_id"].GetStringValue()
	plugin.Logger(ctx).Info("getCvmDisasterRecoverGroup", "disaster_recover_group_id", groupId, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &cvm.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getCvmDisasterRecoverGroup", "client_init_error", err)
		return nil, err
	}

	if groupId == "" {
		return nil, nil
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	req := cvm.NewDescribeDisasterRecoverGroupsRequest()
	req.DisasterRecoverGroupIds = common.StringPtrs([]string{groupId})

	utils.LogRequest(ctx, "getCvmDisasterRecoverGroup", req)
	resp, err := client.DescribeDisasterRecoverGroupsWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getCvmDisasterRecoverGroup", "request_error", err)
		return nil, err
	}

	if resp == nil || resp.Response == nil || len(resp.Response.DisasterRecoverGroupSet) == 0 {
		return nil, nil
	}

	row := toCvmDisasterRecoverGroupRow(resp.Response.DisasterRecoverGroupSet[0])
	row.Region = region
	return row, nil
}

// Row Type

type cvmDisasterRecoverGroupRow struct {
	DisasterRecoverGroupId string
	Name                   string
	Type                   string
	CvmQuotaTotal          int64
	CurrentNum             int64
	InstanceIds            []string
	CreateTime             *time.Time
	Affinity               int64
	Tags                   map[string]string
	Region                 string
	Akas                   []string
}

func toCvmDisasterRecoverGroupRow(g *cvm.DisasterRecoverGroup) cvmDisasterRecoverGroupRow {
	row := cvmDisasterRecoverGroupRow{
		DisasterRecoverGroupId: utils.PtrString(g.DisasterRecoverGroupId),
		Name:                   utils.PtrString(g.Name),
		Type:                   utils.PtrString(g.Type),
		CvmQuotaTotal:          utils.PtrInt64(g.CvmQuotaTotal),
		CurrentNum:             utils.PtrInt64(g.CurrentNum),
		InstanceIds:            utils.PtrStringSlice(g.InstanceIds),
		CreateTime:             utils.ParseTimestamp(g.CreateTime),
		Affinity:               utils.PtrInt64(g.Affinity),
		Tags:                   TagsToMap(g.Tags),
	}

	if row.DisasterRecoverGroupId != "" {
		row.Akas = []string{"tencentcloud:cvm:disaster-recover-group:" + row.DisasterRecoverGroupId}
	}
	return row
}
