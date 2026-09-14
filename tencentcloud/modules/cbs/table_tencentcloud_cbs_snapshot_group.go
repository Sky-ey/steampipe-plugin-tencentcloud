package cbs

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"time"

	cbs "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/cbs/v20170312"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudCbsSnapshotGroup() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_cbs_snapshot_group",
		Description:       "Tencent Cloud Block Storage snapshot groups.",
		GetMatrixItemFunc: buildCbsRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "snapshot_group_id", Require: plugin.Optional},
				{Name: "snapshot_group_state", Require: plugin.Optional},
				{Name: "snapshot_group_name", Require: plugin.Optional},
				{Name: "snapshot_id", Require: plugin.Optional},
			},
			Hydrate: listCbsSnapshotGroups,
			Tags:    map[string]string{"service": "cbs", "action": "DescribeSnapshotGroups"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("snapshot_group_id"),
			Hydrate:    getCbsSnapshotGroup,
			Tags:       map[string]string{"service": "cbs", "action": "DescribeSnapshotGroups"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "snapshot_group_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the snapshot group.", Transform: transform.FromField("SnapshotGroupId")},
			{Name: "snapshot_group_name", Type: proto.ColumnType_STRING, Description: "The snapshot group name."},
			{Name: "snapshot_group_type", Type: proto.ColumnType_STRING, Description: "The snapshot group type (NORMAL: normal snapshot group, non-consistency snapshot)."},
			{Name: "snapshot_group_state", Type: proto.ColumnType_STRING, Description: "The current state of the snapshot group (NORMAL, CREATING, ROLLBACKING)."},
			{Name: "contain_root_snapshot", Type: proto.ColumnType_BOOL, Description: "Whether the group contains a system disk snapshot."},
			{Name: "snapshot_id", Type: proto.ColumnType_STRING, Description: "The snapshot ID used to filter snapshot groups; populated when specified in the query.", Transform: transform.FromQual("snapshot_id")},
			{Name: "snapshot_id_set", Type: proto.ColumnType_JSON, Description: "The snapshot IDs contained in the group.", Transform: transform.FromField("SnapshotIdSet")},
			{Name: "percent", Type: proto.ColumnType_INT, Description: "The snapshot group creation progress percentage."},
			{Name: "create_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the snapshot group was created.", Transform: transform.FromField("CreateTime").NullIfZero()},
			{Name: "modify_time", Type: proto.ColumnType_TIMESTAMP, Description: "The latest snapshot group modification time.", Transform: transform.FromField("ModifyTime").NullIfZero()},
			{Name: "images", Type: proto.ColumnType_JSON, Description: "The images associated with the snapshot group."},
			{Name: "image_count", Type: proto.ColumnType_INT, Description: "The number of images associated with the group."},
			{Name: "is_permanent", Type: proto.ColumnType_BOOL, Description: "Whether the snapshot group is retained permanently."},
			{Name: "deadline_time", Type: proto.ColumnType_TIMESTAMP, Description: "The snapshot group expiration time.", Transform: transform.FromField("DeadlineTime").NullIfZero()},
			{Name: "auto_snapshot_policy_id", Type: proto.ColumnType_STRING, Description: "The auto snapshot policy that created the group.", Transform: transform.FromField("AutoSnapshotPolicyId")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the snapshot group resides."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("SnapshotGroupName")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource."},
		}),
	}
}

// List Function

func listCbsSnapshotGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listCbsSnapshotGroups", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &cbs.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listCbsSnapshotGroups", "client_init_error", err)
		return nil, err
	}
	filters := buildCbsSnapshotGroupFilters(d.EqualsQuals)
	region := d.EqualsQualString(utils.MatrixKeyRegion)
	pageSize := uint64(utils.PageSizeFromLimit(d, 100, 1))
	var offset uint64
	for {
		d.WaitForListRateLimit(ctx)
		req := cbs.NewDescribeSnapshotGroupsRequest()
		req.Offset, req.Limit, req.Filters = &offset, &pageSize, filters
		utils.LogRequest(ctx, "listCbsSnapshotGroups", req)
		resp, err := client.DescribeSnapshotGroupsWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listCbsSnapshotGroups", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}
		items := resp.Response.SnapshotGroupSet
		for _, item := range items {
			if item == nil {
				continue
			}
			d.StreamListItem(ctx, toCbsSnapshotGroupRow(item, region))
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if utils.PageDone(int64(offset), int64(len(items)), utils.PtrUint64(resp.Response.TotalCount)) {
			break
		}
		offset += uint64(len(items))
	}
	return nil, nil
}

func buildCbsSnapshotGroupFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*cbs.Filter {
	return utils.BuildFilters(equalsQuals, newCbsFilter,
		utils.FilterMapping{QualName: "snapshot_group_id", FilterName: "snapshot-group-id"},
		utils.FilterMapping{QualName: "snapshot_group_state", FilterName: "snapshot-group-state"},
		utils.FilterMapping{QualName: "snapshot_group_name", FilterName: "snapshot-group-name"},
		utils.FilterMapping{QualName: "snapshot_id", FilterName: "snapshot-id"},
	)
}

// Get Function

func getCbsSnapshotGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	id := d.EqualsQualString("snapshot_group_id")
	plugin.Logger(ctx).Info("getCbsSnapshotGroup", "snapshot_group_id", id, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &cbs.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getCbsSnapshotGroup", "client_init_error", err)
		return nil, err
	}

	if id == "" {
		return nil, nil
	}
	req := cbs.NewDescribeSnapshotGroupsRequest()
	req.Filters = []*cbs.Filter{newCbsFilter("snapshot-group-id", id)}
	utils.LogRequest(ctx, "getCbsSnapshotGroup", req)
	resp, err := client.DescribeSnapshotGroupsWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getCbsSnapshotGroup", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.SnapshotGroupSet) == 0 || resp.Response.SnapshotGroupSet[0] == nil {
		return nil, nil
	}
	return toCbsSnapshotGroupRow(resp.Response.SnapshotGroupSet[0], d.EqualsQualString(utils.MatrixKeyRegion)), nil
}

// Row Type

type cbsSnapshotGroupRow struct {
	SnapshotGroupId      string
	SnapshotGroupName    string
	SnapshotGroupType    string
	SnapshotGroupState   string
	ContainRootSnapshot  bool
	SnapshotIdSet        []string
	Percent              uint64
	CreateTime           *time.Time
	ModifyTime           *time.Time
	Images               []*cbs.Image
	ImageCount           uint64
	IsPermanent          bool
	DeadlineTime         *time.Time
	AutoSnapshotPolicyId string
	Region               string
	Akas                 []string
}

func toCbsSnapshotGroupRow(item *cbs.SnapshotGroup, region string) cbsSnapshotGroupRow {
	row := cbsSnapshotGroupRow{
		SnapshotGroupId:      utils.PtrString(item.SnapshotGroupId),
		SnapshotGroupName:    utils.PtrString(item.SnapshotGroupName),
		SnapshotGroupType:    utils.PtrString(item.SnapshotGroupType),
		SnapshotGroupState:   utils.PtrString(item.SnapshotGroupState),
		ContainRootSnapshot:  utils.PtrBool(item.ContainRootSnapshot),
		SnapshotIdSet:        utils.PtrStringSlice(item.SnapshotIdSet),
		Percent:              utils.PtrUint64(item.Percent),
		CreateTime:           utils.ParseTimestamp(item.CreateTime),
		ModifyTime:           utils.ParseTimestamp(item.ModifyTime),
		Images:               item.Images,
		ImageCount:           utils.PtrUint64(item.ImageCount),
		IsPermanent:          utils.PtrBool(item.IsPermanent),
		DeadlineTime:         utils.ParseTimestamp(item.DeadlineTime),
		AutoSnapshotPolicyId: utils.PtrString(item.AutoSnapshotPolicyId),
		Region:               region,
	}
	if row.SnapshotGroupId != "" {
		row.Akas = []string{"tencentcloud:cbs:snapshot-group:" + row.SnapshotGroupId}
	}
	return row
}
