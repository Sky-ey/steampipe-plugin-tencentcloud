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

func TableTencentcloudCbsAutoSnapshotPolicy() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_cbs_auto_snapshot_policy",
		Description:       "Tencent Cloud Block Storage auto snapshot policies.",
		GetMatrixItemFunc: buildCbsRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "auto_snapshot_policy_id", Require: plugin.Optional},
				{Name: "auto_snapshot_policy_state", Require: plugin.Optional},
				{Name: "auto_snapshot_policy_name", Require: plugin.Optional},
			},
			Hydrate: listCbsAutoSnapshotPolicies,
			Tags:    map[string]string{"service": "cbs", "action": "DescribeAutoSnapshotPolicies"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("auto_snapshot_policy_id"),
			Hydrate:    getCbsAutoSnapshotPolicy,
			Tags:       map[string]string{"service": "cbs", "action": "DescribeAutoSnapshotPolicies"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "auto_snapshot_policy_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the auto snapshot policy.", Transform: transform.FromField("AutoSnapshotPolicyId")},
			{Name: "auto_snapshot_policy_name", Type: proto.ColumnType_STRING, Description: "The policy name."},
			{Name: "auto_snapshot_policy_state", Type: proto.ColumnType_STRING, Description: "The current state of the policy (NORMAL, ISOLATED)."},
			{Name: "is_activated", Type: proto.ColumnType_BOOL, Description: "Whether the policy is activated."},
			{Name: "policy", Type: proto.ColumnType_JSON, Description: "The auto snapshot execution schedule."},
			{Name: "create_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the policy was created.", Transform: transform.FromField("CreateTime").NullIfZero()},
			{Name: "next_trigger_time", Type: proto.ColumnType_TIMESTAMP, Description: "The next time the policy will trigger.", Transform: transform.FromField("NextTriggerTime").NullIfZero()},
			{Name: "is_permanent", Type: proto.ColumnType_BOOL, Description: "Whether snapshots created by the policy are retained permanently."},
			{Name: "retention_days", Type: proto.ColumnType_INT, Description: "The number of days snapshots are retained."},
			{Name: "retention_months", Type: proto.ColumnType_INT, Description: "The number of months snapshots are retained."},
			{Name: "retention_amount", Type: proto.ColumnType_INT, Description: "The maximum number of snapshots retained."},
			{Name: "advanced_retention_policy", Type: proto.ColumnType_JSON, Description: "The advanced snapshot retention policy."},
			{Name: "disk_id_set", Type: proto.ColumnType_JSON, Description: "The disk IDs bound to the policy.", Transform: transform.FromField("DiskIdSet")},
			{Name: "instance_id_set", Type: proto.ColumnType_JSON, Description: "The instance IDs bound to the policy.", Transform: transform.FromField("InstanceIdSet")},
			{Name: "is_copy_to_remote", Type: proto.ColumnType_INT, Description: "Whether this is a cross-account snapshot copy policy."},
			{Name: "copy_to_account_uin", Type: proto.ColumnType_STRING, Description: "The destination account ID for snapshot copies."},
			{Name: "copy_from_account_uin", Type: proto.ColumnType_STRING, Description: "The source account ID for a copied policy."},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The tags attached to the policy as a map of key to value."},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the policy resides."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("AutoSnapshotPolicyName")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource."},
		}),
	}
}

// List Function

func listCbsAutoSnapshotPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listCbsAutoSnapshotPolicies", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &cbs.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listCbsAutoSnapshotPolicies", "client_init_error", err)
		return nil, err
	}
	filters := utils.BuildFilters(d.EqualsQuals, newCbsFilter,
		utils.FilterMapping{QualName: "auto_snapshot_policy_id", FilterName: "auto-snapshot-policy-id"},
		utils.FilterMapping{QualName: "auto_snapshot_policy_state", FilterName: "auto-snapshot-policy-state"},
		utils.FilterMapping{QualName: "auto_snapshot_policy_name", FilterName: "auto-snapshot-policy-name"},
	)
	region := d.EqualsQualString(utils.MatrixKeyRegion)
	pageSize := uint64(utils.PageSizeFromLimit(d, 100, 1))
	var offset uint64
	for {
		d.WaitForListRateLimit(ctx)
		req := cbs.NewDescribeAutoSnapshotPoliciesRequest()
		req.Offset, req.Limit, req.Filters = &offset, &pageSize, filters
		utils.LogRequest(ctx, "listCbsAutoSnapshotPolicies", req)
		resp, err := client.DescribeAutoSnapshotPoliciesWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listCbsAutoSnapshotPolicies", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}
		items := resp.Response.AutoSnapshotPolicySet
		for _, item := range items {
			if item == nil {
				continue
			}
			d.StreamListItem(ctx, toCbsAutoSnapshotPolicyRow(item, region))
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

func getCbsAutoSnapshotPolicy(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	id := d.EqualsQualString("auto_snapshot_policy_id")
	plugin.Logger(ctx).Info("getCbsAutoSnapshotPolicy", "auto_snapshot_policy_id", id, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &cbs.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getCbsAutoSnapshotPolicy", "client_init_error", err)
		return nil, err
	}

	if id == "" {
		return nil, nil
	}
	req := cbs.NewDescribeAutoSnapshotPoliciesRequest()
	req.AutoSnapshotPolicyIds = common.StringPtrs([]string{id})
	utils.LogRequest(ctx, "getCbsAutoSnapshotPolicy", req)
	resp, err := client.DescribeAutoSnapshotPoliciesWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getCbsAutoSnapshotPolicy", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.AutoSnapshotPolicySet) == 0 || resp.Response.AutoSnapshotPolicySet[0] == nil {
		return nil, nil
	}
	return toCbsAutoSnapshotPolicyRow(resp.Response.AutoSnapshotPolicySet[0], d.EqualsQualString(utils.MatrixKeyRegion)), nil
}

// Row Type

type cbsAutoSnapshotPolicyRow struct {
	AutoSnapshotPolicyId    string
	AutoSnapshotPolicyName  string
	AutoSnapshotPolicyState string
	IsActivated             bool
	Policy                  []*cbs.Policy
	CreateTime              *time.Time
	NextTriggerTime         *time.Time
	IsPermanent             bool
	RetentionDays           uint64
	RetentionMonths         uint64
	RetentionAmount         uint64
	AdvancedRetentionPolicy *cbs.AdvancedRetentionPolicy
	DiskIdSet               []string
	InstanceIdSet           []string
	IsCopyToRemote          uint64
	CopyToAccountUin        string
	CopyFromAccountUin      string
	Tags                    map[string]string
	Region                  string
	Akas                    []string
}

func toCbsAutoSnapshotPolicyRow(item *cbs.AutoSnapshotPolicy, region string) cbsAutoSnapshotPolicyRow {
	row := cbsAutoSnapshotPolicyRow{
		AutoSnapshotPolicyId:    utils.PtrString(item.AutoSnapshotPolicyId),
		AutoSnapshotPolicyName:  utils.PtrString(item.AutoSnapshotPolicyName),
		AutoSnapshotPolicyState: utils.PtrString(item.AutoSnapshotPolicyState),
		IsActivated:             utils.PtrBool(item.IsActivated),
		Policy:                  item.Policy,
		CreateTime:              utils.ParseTimestamp(item.CreateTime),
		NextTriggerTime:         utils.ParseTimestamp(item.NextTriggerTime),
		IsPermanent:             utils.PtrBool(item.IsPermanent),
		RetentionDays:           utils.PtrUint64(item.RetentionDays),
		RetentionMonths:         utils.PtrUint64(item.RetentionMonths),
		RetentionAmount:         utils.PtrUint64(item.RetentionAmount),
		AdvancedRetentionPolicy: item.AdvancedRetentionPolicy,
		DiskIdSet:               utils.PtrStringSlice(item.DiskIdSet),
		InstanceIdSet:           utils.PtrStringSlice(item.InstanceIdSet),
		IsCopyToRemote:          utils.PtrUint64(item.IsCopyToRemote),
		CopyToAccountUin:        utils.PtrString(item.CopyToAccountUin),
		CopyFromAccountUin:      utils.PtrString(item.CopyFromAccountUin),
		Tags:                    TagsToMap(item.Tags),
		Region:                  region,
	}
	if row.AutoSnapshotPolicyId != "" {
		row.Akas = []string{"tencentcloud:cbs:auto-snapshot-policy:" + row.AutoSnapshotPolicyId}
	}
	return row
}
