package cbs

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"strconv"
	"strings"
	"time"

	cbs "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/cbs/v20170312"
	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudCbsSnapshot() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_cbs_snapshot",
		Description:       "Tencent Cloud Block Storage snapshots.",
		GetMatrixItemFunc: buildCbsRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "snapshot_id", Require: plugin.Optional},
				{Name: "snapshot_name", Require: plugin.Optional},
				{Name: "snapshot_state", Require: plugin.Optional},
				{Name: "disk_usage", Require: plugin.Optional},
				{Name: "project_id", Require: plugin.Optional},
				{Name: "disk_id", Require: plugin.Optional},
				{Name: "encrypt", Require: plugin.Optional},
				{Name: "snapshot_type", Require: plugin.Optional},
			},
			Hydrate: listCbsSnapshots,
			Tags:    map[string]string{"service": "cbs", "action": "DescribeSnapshots"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("snapshot_id"),
			Hydrate:    getCbsSnapshot,
			Tags:       map[string]string{"service": "cbs", "action": "DescribeSnapshots"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "snapshot_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the snapshot.", Transform: transform.FromField("SnapshotId")},
			{Name: "snapshot_name", Type: proto.ColumnType_STRING, Description: "The snapshot name."},
			{Name: "snapshot_state", Type: proto.ColumnType_STRING, Description: "The current state of the snapshot (NORMAL, CREATING, ROLLBACKING, COPYING_FROM_REMOTE, CHECKING_COPIED, TORECYCLE)."},
			{Name: "snapshot_type", Type: proto.ColumnType_STRING, Description: "The snapshot type (PRIVATE_SNAPSHOT, SHARED_SNAPSHOT)."},
			{Name: "disk_id", Type: proto.ColumnType_STRING, Description: "The ID of the source disk.", Transform: transform.FromField("DiskId")},
			{Name: "disk_size", Type: proto.ColumnType_INT, Description: "The source disk size in GiB."},
			{Name: "disk_usage", Type: proto.ColumnType_STRING, Description: "The source disk usage (SYSTEM_DISK, DATA_DISK)."},
			{Name: "percent", Type: proto.ColumnType_INT, Description: "The snapshot creation progress percentage."},
			{Name: "create_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the snapshot was created.", Transform: transform.FromField("CreateTime").NullIfZero()},
			{Name: "deadline_time", Type: proto.ColumnType_TIMESTAMP, Description: "The snapshot expiration time.", Transform: transform.FromField("DeadlineTime").NullIfZero()},
			{Name: "time_start_share", Type: proto.ColumnType_TIMESTAMP, Description: "The time when snapshot sharing started.", Transform: transform.FromField("TimeStartShare").NullIfZero()},
			{Name: "is_permanent", Type: proto.ColumnType_BOOL, Description: "Whether the snapshot is retained permanently."},
			{Name: "encrypt", Type: proto.ColumnType_BOOL, Description: "Whether the snapshot was created from an encrypted disk."},
			{Name: "copy_from_remote", Type: proto.ColumnType_BOOL, Description: "Whether the snapshot was copied from another region."},
			{Name: "copying_to_regions", Type: proto.ColumnType_JSON, Description: "The destination regions to which the snapshot is being copied."},
			{Name: "images", Type: proto.ColumnType_JSON, Description: "The images associated with the snapshot."},
			{Name: "image_count", Type: proto.ColumnType_INT, Description: "The number of images associated with the snapshot."},
			{Name: "share_reference", Type: proto.ColumnType_INT, Description: "The current number of snapshot shares."},
			{Name: "zone", Type: proto.ColumnType_STRING, Description: "The availability zone of the snapshot."},
			{Name: "project_id", Type: proto.ColumnType_INT, Description: "The project ID of the snapshot.", Transform: transform.FromField("ProjectId")},
			{Name: "project_name", Type: proto.ColumnType_STRING, Description: "The project name of the snapshot."},
			{Name: "placement", Type: proto.ColumnType_JSON, Description: "The complete placement information of the snapshot."},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The tags attached to the snapshot as a map of key to value."},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the snapshot resides."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("SnapshotName")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource."},
		}),
	}
}

// List Function

func listCbsSnapshots(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listCbsSnapshots", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &cbs.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listCbsSnapshots", "client_init_error", err)
		return nil, err
	}
	filters := buildCbsSnapshotFilters(d.EqualsQuals)
	region := d.EqualsQualString(utils.MatrixKeyRegion)
	pageSize := uint64(utils.PageSizeFromLimit(d, 100, 1))
	var offset uint64
	for {
		d.WaitForListRateLimit(ctx)
		req := cbs.NewDescribeSnapshotsRequest()
		req.Offset, req.Limit, req.Filters = &offset, &pageSize, filters
		utils.LogRequest(ctx, "listCbsSnapshots", req)
		resp, err := client.DescribeSnapshotsWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listCbsSnapshots", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}
		items := resp.Response.SnapshotSet
		for _, item := range items {
			if item == nil {
				continue
			}
			d.StreamListItem(ctx, toCbsSnapshotRow(item, region))
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

// Get Function

func getCbsSnapshot(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	id := d.EqualsQualString("snapshot_id")
	plugin.Logger(ctx).Info("getCbsSnapshot", "snapshot_id", id, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &cbs.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getCbsSnapshot", "client_init_error", err)
		return nil, err
	}

	if id == "" {
		return nil, nil
	}
	req := cbs.NewDescribeSnapshotsRequest()
	req.SnapshotIds = common.StringPtrs([]string{id})
	utils.LogRequest(ctx, "getCbsSnapshot", req)
	resp, err := client.DescribeSnapshotsWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getCbsSnapshot", "request_error", err)
		return nil, err
	}
	if len(resp.Response.SnapshotSet) == 0 || resp.Response.SnapshotSet[0] == nil {
		return nil, nil
	}
	return toCbsSnapshotRow(resp.Response.SnapshotSet[0], d.EqualsQualString(utils.MatrixKeyRegion)), nil
}

func buildCbsSnapshotFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*cbs.Filter {
	filters := utils.BuildFilters(equalsQuals, newCbsFilter,
		utils.FilterMapping{QualName: "snapshot_id", FilterName: "snapshot-id"},
		utils.FilterMapping{QualName: "snapshot_name", FilterName: "snapshot-name"},
		utils.FilterMapping{QualName: "snapshot_state", FilterName: "snapshot-state"},
		utils.FilterMapping{QualName: "disk_usage", FilterName: "disk-usage"},
		utils.FilterMapping{QualName: "disk_id", FilterName: "disk-id"},
		utils.FilterMapping{QualName: "snapshot_type", FilterName: "snapshot-type"},
	)
	if qual := equalsQuals["project_id"]; qual != nil {
		filters = append(filters, newCbsFilter("project-id", strconv.FormatInt(qual.GetInt64Value(), 10)))
	}
	if qual := equalsQuals["encrypt"]; qual != nil {
		filters = append(filters, newCbsFilter("encrypt", strings.ToUpper(strconv.FormatBool(qual.GetBoolValue()))))
	}
	return filters
}

// Row Type

type cbsSnapshotRow struct {
	SnapshotId       string
	SnapshotName     string
	SnapshotState    string
	SnapshotType     string
	DiskId           string
	DiskSize         uint64
	DiskUsage        string
	Percent          uint64
	CreateTime       *time.Time
	DeadlineTime     *time.Time
	TimeStartShare   *time.Time
	IsPermanent      bool
	Encrypt          bool
	CopyFromRemote   bool
	CopyingToRegions []string
	Images           []*cbs.Image
	ImageCount       uint64
	ShareReference   uint64
	Zone             string
	ProjectId        uint64
	ProjectName      string
	Placement        *cbs.Placement
	Tags             map[string]string
	Region           string
	Akas             []string
}

func toCbsSnapshotRow(item *cbs.Snapshot, region string) cbsSnapshotRow {
	row := cbsSnapshotRow{
		SnapshotId:       utils.PtrString(item.SnapshotId),
		SnapshotName:     utils.PtrString(item.SnapshotName),
		SnapshotState:    utils.PtrString(item.SnapshotState),
		SnapshotType:     utils.PtrString(item.SnapshotType),
		DiskId:           utils.PtrString(item.DiskId),
		DiskSize:         utils.PtrUint64(item.DiskSize),
		DiskUsage:        utils.PtrString(item.DiskUsage),
		Percent:          utils.PtrUint64(item.Percent),
		CreateTime:       utils.ParseTimestamp(item.CreateTime),
		DeadlineTime:     utils.ParseTimestamp(item.DeadlineTime),
		TimeStartShare:   utils.ParseTimestamp(item.TimeStartShare),
		IsPermanent:      utils.PtrBool(item.IsPermanent),
		Encrypt:          utils.PtrBool(item.Encrypt),
		CopyFromRemote:   utils.PtrBool(item.CopyFromRemote),
		CopyingToRegions: utils.PtrStringSlice(item.CopyingToRegions),
		Images:           item.Images,
		ImageCount:       utils.PtrUint64(item.ImageCount),
		ShareReference:   utils.PtrUint64(item.ShareReference),
		Placement:        item.Placement,
		Tags:             TagsToMap(item.Tags),
		Region:           region,
	}
	if item.Placement != nil {
		row.Zone = utils.PtrString(item.Placement.Zone)
		row.ProjectId = utils.PtrUint64(item.Placement.ProjectId)
		row.ProjectName = utils.PtrString(item.Placement.ProjectName)
	}
	if row.SnapshotId != "" {
		row.Akas = []string{"tencentcloud:cbs:snapshot:" + row.SnapshotId}
	}
	return row
}
