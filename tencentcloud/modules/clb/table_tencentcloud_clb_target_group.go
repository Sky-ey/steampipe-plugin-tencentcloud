package clb

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"time"

	clb "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/clb/v20180317"
	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudClbTargetGroup() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_clb_target_group",
		Description:       "Tencent Cloud Load Balancer (CLB) target groups.",
		GetMatrixItemFunc: buildClbRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "target_group_id", Require: plugin.Optional},
				{Name: "target_group_name", Require: plugin.Optional},
				{Name: "vpc_id", Require: plugin.Optional},
			},
			Hydrate: listClbTargetGroups,
			Tags:    map[string]string{"service": "clb", "action": "DescribeTargetGroups"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("target_group_id"),
			Hydrate:    getClbTargetGroup,
			Tags:       map[string]string{"service": "clb", "action": "DescribeTargetGroups"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "target_group_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the target group.", Transform: transform.FromField("TargetGroupId")},
			{Name: "target_group_name", Type: proto.ColumnType_STRING, Description: "The name of the target group."},
			{Name: "target_group_type", Type: proto.ColumnType_STRING, Description: "The type of the target group: v1 (legacy) or v2 (current)."},
			{Name: "protocol", Type: proto.ColumnType_STRING, Description: "The backend forwarding protocol of the target group."},
			{Name: "port", Type: proto.ColumnType_INT, Description: "The default backend port of the target group."},
			{Name: "vpc_id", Type: proto.ColumnType_STRING, Description: "The ID of the VPC to which the target group belongs.", Transform: transform.FromField("VpcId")},
			{Name: "weight", Type: proto.ColumnType_INT, Description: "The default forwarding weight of the target group."},
			{Name: "schedule_algorithm", Type: proto.ColumnType_STRING, Description: "The scheduling algorithm of the target group (WRR, LEAST_CONN, IP_HASH), valid for target groups with HTTP, HTTPS or GRPC backend protocol."},
			{Name: "full_listen_switch", Type: proto.ColumnType_BOOL, Description: "Whether the target group is bound to all listeners of the load balancer."},
			{Name: "keepalive_enable", Type: proto.ColumnType_BOOL, Description: "Whether HTTP keepalive is enabled."},
			{Name: "session_expire_time", Type: proto.ColumnType_INT, Description: "The session persistence timeout in seconds."},
			{Name: "ip_version", Type: proto.ColumnType_STRING, Description: "The IP version of the target group.", Transform: transform.FromField("IpVersion")},
			{Name: "snat_enable", Type: proto.ColumnType_BOOL, Description: "Whether SNAT is enabled for the target group.", Transform: transform.FromField("SnatEnable")},
			{Name: "associated_rule_count", Type: proto.ColumnType_INT, Description: "The number of rules associated with the target group."},
			{Name: "registered_instances_count", Type: proto.ColumnType_INT, Description: "The number of instances registered in the target group."},
			{Name: "created_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the target group was created.", Transform: transform.FromField("CreatedTime").NullIfZero()},
			{Name: "updated_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the target group was last updated.", Transform: transform.FromField("UpdatedTime").NullIfZero()},
			{Name: "health_check", Type: proto.ColumnType_JSON, Description: "The health check configuration of the target group."},
			{Name: "associated_rule", Type: proto.ColumnType_JSON, Description: "The rules associated with the target group.", Transform: transform.FromField("AssociatedRule")},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The tags attached to the target group as a map of key to value."},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the target group resides."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource."},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource."},
		}),
	}
}

// List Function

func listClbTargetGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listClbTargetGroups", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &clb.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listClbTargetGroups", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)
	filters := utils.BuildFilters(d.EqualsQuals, newClbFilter,
		utils.FilterMapping{QualName: "vpc_id", FilterName: "TargetGroupVpcId"},
		utils.FilterMapping{QualName: "target_group_name", FilterName: "TargetGroupName"},
	)
	pageSize := uint64(utils.PageSizeFromLimit(d, 100, 1))
	var offset uint64

	for {
		d.WaitForListRateLimit(ctx)
		req := clb.NewDescribeTargetGroupListRequest()
		req.Offset = &offset
		req.Limit = &pageSize

		if id := d.EqualsQualString("target_group_id"); id != "" {
			// TargetGroupIds 与 Filters 互斥,ID 优先
			req.TargetGroupIds = common.StringPtrs([]string{id})
		} else if len(filters) > 0 {
			req.Filters = filters
		}

		utils.LogRequest(ctx, "listClbTargetGroups", req)
		resp, err := client.DescribeTargetGroupListWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listClbTargetGroups", "request_error", err)
			return nil, err
		}

		if resp == nil || resp.Response == nil {
			break
		}
		targetGroups := resp.Response.TargetGroupSet
		for _, targetGroup := range targetGroups {
			if targetGroup == nil {
				continue
			}
			row := toClbTargetGroupRow(targetGroup, region)
			d.StreamListItem(ctx, row)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if utils.PageDone(int64(offset), int64(len(targetGroups)), utils.PtrUint64(resp.Response.TotalCount)) {
			break
		}
		offset += uint64(len(targetGroups))
	}
	return nil, nil
}

// Get Function

func getClbTargetGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	targetGroupID := d.EqualsQualString("target_group_id")
	plugin.Logger(ctx).Info("getClbTargetGroup", "target_group_id", targetGroupID, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &clb.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getClbTargetGroup", "client_init_error", err)
		return nil, err
	}

	if targetGroupID == "" {
		return nil, nil
	}
	req := clb.NewDescribeTargetGroupListRequest()
	req.TargetGroupIds = common.StringPtrs([]string{targetGroupID})

	utils.LogRequest(ctx, "getClbTargetGroup", req)
	resp, err := client.DescribeTargetGroupListWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getClbTargetGroup", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.TargetGroupSet) == 0 || resp.Response.TargetGroupSet[0] == nil {
		return nil, nil
	}
	return toClbTargetGroupRow(resp.Response.TargetGroupSet[0], d.EqualsQualString(utils.MatrixKeyRegion)), nil
}

// Row Type

type clbTargetGroupRow struct {
	TargetGroupId            string
	TargetGroupName          string
	TargetGroupType          string
	Protocol                 string
	Port                     uint64
	VpcId                    string
	Weight                   uint64
	ScheduleAlgorithm        string
	FullListenSwitch         bool
	KeepaliveEnable          bool
	SessionExpireTime        int64
	IpVersion                string
	SnatEnable               bool
	AssociatedRuleCount      int64
	RegisteredInstancesCount int64
	CreatedTime              *time.Time
	UpdatedTime              *time.Time
	HealthCheck              *clb.TargetGroupHealthCheck
	AssociatedRule           []*clb.AssociationItem
	Tags                     map[string]string
	Region                   string
	Title                    string
	Akas                     []string
}

func toClbTargetGroupRow(targetGroup *clb.TargetGroupInfo, region string) clbTargetGroupRow {
	row := clbTargetGroupRow{
		TargetGroupId:            utils.PtrString(targetGroup.TargetGroupId),
		TargetGroupName:          utils.PtrString(targetGroup.TargetGroupName),
		TargetGroupType:          utils.PtrString(targetGroup.TargetGroupType),
		Protocol:                 utils.PtrString(targetGroup.Protocol),
		Port:                     utils.PtrUint64(targetGroup.Port),
		VpcId:                    utils.PtrString(targetGroup.VpcId),
		Weight:                   utils.PtrUint64(targetGroup.Weight),
		ScheduleAlgorithm:        utils.PtrString(targetGroup.ScheduleAlgorithm),
		FullListenSwitch:         utils.PtrBool(targetGroup.FullListenSwitch),
		KeepaliveEnable:          utils.PtrBool(targetGroup.KeepaliveEnable),
		SessionExpireTime:        utils.PtrInt64(targetGroup.SessionExpireTime),
		IpVersion:                utils.PtrString(targetGroup.IpVersion),
		SnatEnable:               utils.PtrBool(targetGroup.SnatEnable),
		AssociatedRuleCount:      utils.PtrInt64(targetGroup.AssociatedRuleCount),
		RegisteredInstancesCount: utils.PtrInt64(targetGroup.RegisteredInstancesCount),
		CreatedTime:              utils.ParseTimestamp(targetGroup.CreatedTime),
		UpdatedTime:              utils.ParseTimestamp(targetGroup.UpdatedTime),
		HealthCheck:              targetGroup.HealthCheck,
		AssociatedRule:           targetGroup.AssociatedRule,
		Tags:                     TagsToMap(targetGroup.Tag),
		Region:                   region,
	}
	row.Title = row.TargetGroupName
	if row.Title == "" {
		row.Title = row.TargetGroupId
	}
	if row.TargetGroupId != "" {
		row.Akas = []string{"tencentcloud:clb:target-group:" + row.TargetGroupId}
	}
	return row
}
