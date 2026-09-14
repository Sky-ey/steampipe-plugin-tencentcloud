package cbs

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"time"

	cbs "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/cbs/v20170312"
	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudCbsDiskBackup() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_cbs_disk_backup",
		Description:       "Tencent Cloud Block Storage disk backup points.",
		GetMatrixItemFunc: buildCbsRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "disk_backup_id", Require: plugin.Optional},
				{Name: "disk_id", Require: plugin.Optional},
				{Name: "disk_usage", Require: plugin.Optional},
			},
			Hydrate: listCbsDiskBackups,
			Tags:    map[string]string{"service": "cbs", "action": "DescribeDiskBackups"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("disk_backup_id"),
			Hydrate:    getCbsDiskBackup,
			Tags:       map[string]string{"service": "cbs", "action": "DescribeDiskBackups"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "disk_backup_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the disk backup point.", Transform: transform.FromField("DiskBackupId")},
			{Name: "disk_backup_name", Type: proto.ColumnType_STRING, Description: "The name of the disk backup point."},
			{Name: "disk_backup_state", Type: proto.ColumnType_STRING, Description: "The state of the disk backup point (NORMAL, CREATING, ROLLBACKING)."},
			{Name: "disk_id", Type: proto.ColumnType_STRING, Description: "The ID of the source disk.", Transform: transform.FromField("DiskId")},
			{Name: "disk_size", Type: proto.ColumnType_INT, Description: "The source disk size in GiB."},
			{Name: "disk_usage", Type: proto.ColumnType_STRING, Description: "The source disk usage (SYSTEM_DISK, DATA_DISK)."},
			{Name: "percent", Type: proto.ColumnType_INT, Description: "The backup creation progress percentage."},
			{Name: "create_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the backup point was created.", Transform: transform.FromField("CreateTime").NullIfZero()},
			{Name: "encrypt", Type: proto.ColumnType_BOOL, Description: "Whether the source disk is encrypted."},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the backup point resides."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("DiskBackupName")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource."},
		}),
	}
}

// List Function

func listCbsDiskBackups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listCbsDiskBackups", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &cbs.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listCbsDiskBackups", "client_init_error", err)
		return nil, err
	}
	filters := utils.BuildFilters(d.EqualsQuals, newCbsFilter,
		utils.FilterMapping{QualName: "disk_backup_id", FilterName: "disk-backup-id"},
		utils.FilterMapping{QualName: "disk_id", FilterName: "disk-id"},
		utils.FilterMapping{QualName: "disk_usage", FilterName: "disk-usage"},
	)
	region := d.EqualsQualString(utils.MatrixKeyRegion)
	// DescribeDiskBackups requires Limit in [10, 100]; minSize must match the API lower bound.
	pageSize := uint64(utils.PageSizeFromLimit(d, 100, 10))
	var offset uint64
	for {
		d.WaitForListRateLimit(ctx)
		req := cbs.NewDescribeDiskBackupsRequest()
		req.Offset, req.Limit, req.Filters = &offset, &pageSize, filters
		utils.LogRequest(ctx, "listCbsDiskBackups", req)
		resp, err := client.DescribeDiskBackupsWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listCbsDiskBackups", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}
		items := resp.Response.DiskBackupSet
		for _, item := range items {
			if item == nil {
				continue
			}
			d.StreamListItem(ctx, toCbsDiskBackupRow(item, region))
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

func getCbsDiskBackup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	id := d.EqualsQualString("disk_backup_id")
	plugin.Logger(ctx).Debug("getCbsDiskBackup", "disk_backup_id", id, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &cbs.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getCbsDiskBackup", "client_init_error", err)
		return nil, err
	}

	if id == "" {
		return nil, nil
	}
	req := cbs.NewDescribeDiskBackupsRequest()
	req.DiskBackupIds = common.StringPtrs([]string{id})
	utils.LogRequest(ctx, "getCbsDiskBackup", req)
	resp, err := client.DescribeDiskBackupsWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getCbsDiskBackup", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.DiskBackupSet) == 0 || resp.Response.DiskBackupSet[0] == nil {
		return nil, nil
	}
	return toCbsDiskBackupRow(resp.Response.DiskBackupSet[0], d.EqualsQualString(utils.MatrixKeyRegion)), nil
}

// Row Type

type cbsDiskBackupRow struct {
	DiskBackupId    string
	DiskId          string
	DiskSize        uint64
	DiskUsage       string
	DiskBackupName  string
	DiskBackupState string
	Percent         uint64
	CreateTime      *time.Time
	Encrypt         bool
	Region          string
	Akas            []string
}

func toCbsDiskBackupRow(item *cbs.DiskBackup, region string) cbsDiskBackupRow {
	row := cbsDiskBackupRow{
		DiskBackupId:    utils.PtrString(item.DiskBackupId),
		DiskId:          utils.PtrString(item.DiskId),
		DiskSize:        utils.PtrUint64(item.DiskSize),
		DiskUsage:       utils.PtrString(item.DiskUsage),
		DiskBackupName:  utils.PtrString(item.DiskBackupName),
		DiskBackupState: utils.PtrString(item.DiskBackupState),
		Percent:         utils.PtrUint64(item.Percent),
		CreateTime:      utils.ParseTimestamp(item.CreateTime),
		Encrypt:         utils.PtrBool(item.Encrypt),
		Region:          region,
	}
	if row.DiskBackupId != "" {
		row.Akas = []string{"tencentcloud:cbs:disk-backup:" + row.DiskBackupId}
	}
	return row
}
