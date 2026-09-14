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

func TableTencentcloudCbsDisk() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_cbs_disk",
		Description:       "Tencent Cloud Block Storage (CBS) disks.",
		GetMatrixItemFunc: buildCbsRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "disk_id", Require: plugin.Optional},
				{Name: "disk_usage", Require: plugin.Optional},
				{Name: "disk_charge_type", Require: plugin.Optional},
				{Name: "portable", Require: plugin.Optional},
				{Name: "project_id", Require: plugin.Optional},
				{Name: "disk_name", Require: plugin.Optional},
				{Name: "disk_type", Require: plugin.Optional},
				{Name: "disk_state", Require: plugin.Optional},
				{Name: "instance_id", Require: plugin.Optional},
				{Name: "zone", Require: plugin.Optional},
			},
			Hydrate: listCbsDisks,
			Tags:    map[string]string{"service": "cbs", "action": "DescribeDisks"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("disk_id"),
			Hydrate:    getCbsDisk,
			Tags:       map[string]string{"service": "cbs", "action": "DescribeDisks"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "disk_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the disk.", Transform: transform.FromField("DiskId")},
			{Name: "disk_name", Type: proto.ColumnType_STRING, Description: "The disk name."},
			{Name: "disk_type", Type: proto.ColumnType_STRING, Description: "The storage media type of the disk (CLOUD_BASIC, CLOUD_PREMIUM, CLOUD_BSSD, CLOUD_SSD, CLOUD_HSSD, CLOUD_TSSD)."},
			{Name: "disk_state", Type: proto.ColumnType_STRING, Description: "The current state of the disk (UNATTACHED, ATTACHING, ATTACHED, DETACHING, EXPANDING, ROLLBACKING, TORECYCLE, DUMPING)."},
			{Name: "disk_usage", Type: proto.ColumnType_STRING, Description: "The disk usage (SYSTEM_DISK, DATA_DISK)."},
			{Name: "disk_charge_type", Type: proto.ColumnType_STRING, Description: "The disk billing mode (PREPAID, POSTPAID_BY_HOUR)."},
			{Name: "disk_size", Type: proto.ColumnType_INT, Description: "The disk size in GiB."},
			{Name: "instance_id", Type: proto.ColumnType_STRING, Description: "The ID of the instance to which the disk is attached.", Transform: transform.FromField("InstanceId")},
			{Name: "instance_id_list", Type: proto.ColumnType_JSON, Description: "The IDs of instances to which a shared disk is attached.", Transform: transform.FromField("InstanceIdList")},
			{Name: "instance_type", Type: proto.ColumnType_STRING, Description: "The type of instance to which the disk is attached."},
			{Name: "last_attach_instance_id", Type: proto.ColumnType_STRING, Description: "The ID of the instance to which the disk was most recently attached.", Transform: transform.FromField("LastAttachInsId")},
			{Name: "attach_mode", Type: proto.ColumnType_STRING, Description: "The disk attachment mode (PF, VF)."},
			{Name: "attached", Type: proto.ColumnType_BOOL, Description: "Whether the disk is attached."},
			{Name: "portable", Type: proto.ColumnType_BOOL, Description: "Whether the disk is portable."},
			{Name: "shareable", Type: proto.ColumnType_BOOL, Description: "Whether the disk is shareable."},
			{Name: "encrypt", Type: proto.ColumnType_BOOL, Description: "Whether the disk is encrypted."},
			{Name: "encrypt_type", Type: proto.ColumnType_STRING, Description: "The disk encryption technology type."},
			{Name: "kms_key_id", Type: proto.ColumnType_STRING, Description: "The KMS key ID used by the encrypted disk.", Transform: transform.FromField("KmsKeyId")},
			{Name: "delete_with_instance", Type: proto.ColumnType_BOOL, Description: "Whether the disk is deleted with its attached instance."},
			{Name: "backup_disk", Type: proto.ColumnType_BOOL, Description: "Whether a backup snapshot is created when the disk is destroyed."},
			{Name: "delete_snapshot", Type: proto.ColumnType_INT, Description: "Whether non-permanent snapshots are deleted with the disk: 1 (deleted), 0 (not deleted)."},
			{Name: "renew_flag", Type: proto.ColumnType_STRING, Description: "The prepaid disk renewal setting (NOTIFY_AND_AUTO_RENEW, NOTIFY_AND_MANUAL_RENEW, DISABLE_NOTIFY_AND_MANUAL_RENEW)."},
			{Name: "create_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the disk was created.", Transform: transform.FromField("CreateTime").NullIfZero()},
			{Name: "deadline_time", Type: proto.ColumnType_TIMESTAMP, Description: "The disk expiration time.", Transform: transform.FromField("DeadlineTime").NullIfZero()},
			{Name: "differ_days_of_deadline", Type: proto.ColumnType_INT, Description: "The number of days until a prepaid disk expires."},
			{Name: "snapshot_count", Type: proto.ColumnType_INT, Description: "The number of snapshots owned by the disk."},
			{Name: "snapshot_size", Type: proto.ColumnType_INT, Description: "The total size of disk snapshots in MiB."},
			{Name: "snapshot_ability", Type: proto.ColumnType_BOOL, Description: "Whether snapshots can be created for the disk."},
			{Name: "auto_snapshot_policy_ids", Type: proto.ColumnType_JSON, Description: "The auto snapshot policy IDs bound to the disk."},
			{Name: "disk_backup_quota", Type: proto.ColumnType_INT, Description: "The maximum number of disk backup points."},
			{Name: "disk_backup_count", Type: proto.ColumnType_INT, Description: "The number of disk backup points in use."},
			{Name: "rollbacking", Type: proto.ColumnType_BOOL, Description: "Whether a snapshot rollback is in progress."},
			{Name: "rollback_percent", Type: proto.ColumnType_INT, Description: "The snapshot rollback progress percentage."},
			{Name: "migrating", Type: proto.ColumnType_BOOL, Description: "Whether a disk type migration is in progress."},
			{Name: "migrate_percent", Type: proto.ColumnType_INT, Description: "The disk type migration progress percentage."},
			{Name: "throughput_performance", Type: proto.ColumnType_INT, Description: "The additional disk throughput in MiB/s."},
			{Name: "burst_performance", Type: proto.ColumnType_BOOL, Description: "Whether performance bursting is enabled."},
			{Name: "auto_renew_flag_error", Type: proto.ColumnType_BOOL, Description: "Whether the disk renewal setting conflicts with the attached instance."},
			{Name: "deadline_error", Type: proto.ColumnType_BOOL, Description: "Whether the disk expires before its attached prepaid instance."},
			{Name: "is_returnable", Type: proto.ColumnType_BOOL, Description: "Whether the prepaid disk can be returned."},
			{Name: "return_fail_code", Type: proto.ColumnType_INT, Description: "The reason code when a prepaid disk cannot be returned."},
			{Name: "error_prompt", Type: proto.ColumnType_STRING, Description: "The latest disk operation error message."},
			{Name: "zone", Type: proto.ColumnType_STRING, Description: "The availability zone of the disk."},
			{Name: "project_id", Type: proto.ColumnType_INT, Description: "The project ID of the disk.", Transform: transform.FromField("ProjectId")},
			{Name: "project_name", Type: proto.ColumnType_STRING, Description: "The project name of the disk."},
			{Name: "cage_id", Type: proto.ColumnType_STRING, Description: "The storage cage ID of the disk.", Transform: transform.FromField("CageId")},
			{Name: "cdc_id", Type: proto.ColumnType_STRING, Description: "The dedicated cloud cluster ID of the disk.", Transform: transform.FromField("CdcId")},
			{Name: "cdc_name", Type: proto.ColumnType_STRING, Description: "The dedicated cloud cluster name of the disk."},
			{Name: "dedicated_cluster_id", Type: proto.ColumnType_STRING, Description: "The dedicated cluster ID of the disk.", Transform: transform.FromField("DedicatedClusterId")},
			{Name: "placement", Type: proto.ColumnType_JSON, Description: "The complete placement information of the disk."},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The tags attached to the disk as a map of key to value."},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the disk resides."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("DiskName")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource."},
		}),
	}
}

// List Function

func listCbsDisks(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listCbsDisks", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &cbs.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listCbsDisks", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)
	filters := buildCbsDiskFilters(d.EqualsQuals)
	pageSize := uint64(utils.PageSizeFromLimit(d, 100, 1))
	var offset uint64
	returnPolicies := true

	for {
		d.WaitForListRateLimit(ctx)
		req := cbs.NewDescribeDisksRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		req.Filters = filters
		req.ReturnBindAutoSnapshotPolicy = &returnPolicies

		utils.LogRequest(ctx, "listCbsDisks", req)
		resp, err := client.DescribeDisksWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listCbsDisks", "request_error", err)
			return nil, err
		}

		if resp == nil || resp.Response == nil {
			break
		}
		disks := resp.Response.DiskSet
		for _, disk := range disks {
			if disk == nil {
				continue
			}
			row := toCbsDiskRow(disk, region)
			d.StreamListItem(ctx, row)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if utils.PageDone(int64(offset), int64(len(disks)), utils.PtrUint64(resp.Response.TotalCount)) {
			break
		}
		offset += uint64(len(disks))
	}
	return nil, nil
}

func buildCbsDiskFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*cbs.Filter {
	filters := utils.BuildFilters(equalsQuals, newCbsFilter,
		utils.FilterMapping{QualName: "disk_id", FilterName: "disk-id"},
		utils.FilterMapping{QualName: "disk_usage", FilterName: "disk-usage"},
		utils.FilterMapping{QualName: "disk_charge_type", FilterName: "disk-charge-type"},
		utils.FilterMapping{QualName: "disk_name", FilterName: "disk-name"},
		utils.FilterMapping{QualName: "disk_type", FilterName: "disk-type"},
		utils.FilterMapping{QualName: "disk_state", FilterName: "disk-state"},
		utils.FilterMapping{QualName: "instance_id", FilterName: "instance-id"},
		utils.FilterMapping{QualName: "zone", FilterName: "zone"},
	)
	if qual := equalsQuals["portable"]; qual != nil {
		filters = append(filters, newCbsFilter("portable", strings.ToUpper(strconv.FormatBool(qual.GetBoolValue()))))
	}
	if qual := equalsQuals["project_id"]; qual != nil {
		filters = append(filters, newCbsFilter("project-id", strconv.FormatInt(qual.GetInt64Value(), 10)))
	}
	return filters
}

// Get Function

func getCbsDisk(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	diskID := d.EqualsQualString("disk_id")
	plugin.Logger(ctx).Info("getCbsDisk", "disk_id", diskID, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &cbs.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getCbsDisk", "client_init_error", err)
		return nil, err
	}

	if diskID == "" {
		return nil, nil
	}
	returnPolicies := true
	req := cbs.NewDescribeDisksRequest()
	req.DiskIds = common.StringPtrs([]string{diskID})
	req.ReturnBindAutoSnapshotPolicy = &returnPolicies
	utils.LogRequest(ctx, "getCbsDisk", req)
	resp, err := client.DescribeDisksWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getCbsDisk", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.DiskSet) == 0 || resp.Response.DiskSet[0] == nil {
		return nil, nil
	}
	return toCbsDiskRow(resp.Response.DiskSet[0], d.EqualsQualString(utils.MatrixKeyRegion)), nil
}

// Row Type

type cbsDiskRow struct {
	DiskId                string
	DiskName              string
	DiskType              string
	DiskState             string
	DiskUsage             string
	DiskChargeType        string
	DiskSize              uint64
	InstanceId            string
	InstanceIdList        []string
	InstanceType          string
	LastAttachInsId       string
	AttachMode            string
	Attached              bool
	Portable              bool
	Shareable             bool
	Encrypt               bool
	EncryptType           string
	KmsKeyId              string
	DeleteWithInstance    bool
	BackupDisk            bool
	DeleteSnapshot        int64
	RenewFlag             string
	CreateTime            *time.Time
	DeadlineTime          *time.Time
	DifferDaysOfDeadline  int64
	SnapshotCount         int64
	SnapshotSize          uint64
	SnapshotAbility       bool
	AutoSnapshotPolicyIds []string
	DiskBackupQuota       uint64
	DiskBackupCount       uint64
	Rollbacking           bool
	RollbackPercent       uint64
	Migrating             bool
	MigratePercent        uint64
	ThroughputPerformance uint64
	BurstPerformance      bool
	AutoRenewFlagError    bool
	DeadlineError         bool
	IsReturnable          bool
	ReturnFailCode        int64
	ErrorPrompt           string
	Zone                  string
	ProjectId             uint64
	ProjectName           string
	CageId                string
	CdcId                 string
	CdcName               string
	DedicatedClusterId    string
	Placement             *cbs.Placement
	Tags                  map[string]string
	Region                string
	Akas                  []string
}

func toCbsDiskRow(disk *cbs.Disk, region string) cbsDiskRow {
	row := cbsDiskRow{
		DiskId:                utils.PtrString(disk.DiskId),
		DiskName:              utils.PtrString(disk.DiskName),
		DiskType:              utils.PtrString(disk.DiskType),
		DiskState:             utils.PtrString(disk.DiskState),
		DiskUsage:             utils.PtrString(disk.DiskUsage),
		DiskChargeType:        utils.PtrString(disk.DiskChargeType),
		DiskSize:              utils.PtrUint64(disk.DiskSize),
		InstanceId:            utils.PtrString(disk.InstanceId),
		InstanceIdList:        utils.PtrStringSlice(disk.InstanceIdList),
		InstanceType:          utils.PtrString(disk.InstanceType),
		LastAttachInsId:       utils.PtrString(disk.LastAttachInsId),
		AttachMode:            utils.PtrString(disk.AttachMode),
		Attached:              utils.PtrBool(disk.Attached),
		Portable:              utils.PtrBool(disk.Portable),
		Shareable:             utils.PtrBool(disk.Shareable),
		Encrypt:               utils.PtrBool(disk.Encrypt),
		EncryptType:           utils.PtrString(disk.EncryptType),
		KmsKeyId:              utils.PtrString(disk.KmsKeyId),
		DeleteWithInstance:    utils.PtrBool(disk.DeleteWithInstance),
		BackupDisk:            utils.PtrBool(disk.BackupDisk),
		DeleteSnapshot:        utils.PtrInt64(disk.DeleteSnapshot),
		RenewFlag:             utils.PtrString(disk.RenewFlag),
		CreateTime:            utils.ParseTimestamp(disk.CreateTime),
		DeadlineTime:          utils.ParseTimestamp(disk.DeadlineTime),
		DifferDaysOfDeadline:  utils.PtrInt64(disk.DifferDaysOfDeadline),
		SnapshotCount:         utils.PtrInt64(disk.SnapshotCount),
		SnapshotSize:          utils.PtrUint64(disk.SnapshotSize),
		SnapshotAbility:       utils.PtrBool(disk.SnapshotAbility),
		AutoSnapshotPolicyIds: utils.PtrStringSlice(disk.AutoSnapshotPolicyIds),
		DiskBackupQuota:       utils.PtrUint64(disk.DiskBackupQuota),
		DiskBackupCount:       utils.PtrUint64(disk.DiskBackupCount),
		Rollbacking:           utils.PtrBool(disk.Rollbacking),
		RollbackPercent:       utils.PtrUint64(disk.RollbackPercent),
		Migrating:             utils.PtrBool(disk.Migrating),
		MigratePercent:        utils.PtrUint64(disk.MigratePercent),
		ThroughputPerformance: utils.PtrUint64(disk.ThroughputPerformance),
		BurstPerformance:      utils.PtrBool(disk.BurstPerformance),
		AutoRenewFlagError:    utils.PtrBool(disk.AutoRenewFlagError),
		DeadlineError:         utils.PtrBool(disk.DeadlineError),
		IsReturnable:          utils.PtrBool(disk.IsReturnable),
		ReturnFailCode:        utils.PtrInt64(disk.ReturnFailCode),
		ErrorPrompt:           utils.PtrString(disk.ErrorPrompt),
		Placement:             disk.Placement,
		Tags:                  TagsToMap(disk.Tags),
		Region:                region,
	}
	if disk.Placement != nil {
		row.Zone = utils.PtrString(disk.Placement.Zone)
		row.ProjectId = utils.PtrUint64(disk.Placement.ProjectId)
		row.ProjectName = utils.PtrString(disk.Placement.ProjectName)
		row.CageId = utils.PtrString(disk.Placement.CageId)
		row.CdcId = utils.PtrString(disk.Placement.CdcId)
		row.CdcName = utils.PtrString(disk.Placement.CdcName)
		row.DedicatedClusterId = utils.PtrString(disk.Placement.DedicatedClusterId)
	}
	if row.DiskId != "" {
		row.Akas = []string{"tencentcloud:cbs:disk:" + row.DiskId}
	}
	return row
}
