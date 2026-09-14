package cls

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"time"

	cls "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/cls/v20201016"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudClsAlarm() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_cls_alarm",
		Description:       "Tencent Cloud CLS alarm policies define trigger conditions, monitoring schedules, and notification channels for log-based alerts within a region.",
		GetMatrixItemFunc: buildClsRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "alarm_id", Require: plugin.Optional},
				{Name: "alarm_name", Require: plugin.Optional},
				{Name: "topic_id", Require: plugin.Optional},
			},
			Hydrate: listClsAlarms,
			Tags:    map[string]string{"service": "cls", "action": "DescribeAlarms"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("alarm_id"),
			Hydrate:    getClsAlarm,
			Tags:       map[string]string{"service": "cls", "action": "DescribeAlarms"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "alarm_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the alarm policy.", Transform: transform.FromField("AlarmId")},
			{Name: "alarm_name", Type: proto.ColumnType_STRING, Description: "The name of the alarm policy.", Transform: transform.FromField("Name")},
			{Name: "status", Type: proto.ColumnType_BOOL, Description: "The enablement status of the alarm policy (true: enabled, false: disabled).", Transform: transform.FromField("Status")},
			{Name: "alarm_level", Type: proto.ColumnType_INT, Description: "The alarm level (0: Warn, 1: Information, 2: Critical).", Transform: transform.FromField("AlarmLevel")},
			{Name: "condition", Type: proto.ColumnType_STRING, Description: "The single trigger condition for whether to trigger an alarm. Mutually exclusive with MultiConditions.", Transform: transform.FromField("Condition")},
			{Name: "trigger_count", Type: proto.ColumnType_INT, Description: "The alarm persistence cycle. An alarm triggers only after the condition is met for this many times. Range: 1-10.", Transform: transform.FromField("TriggerCount")},
			{Name: "alarm_period", Type: proto.ColumnType_INT, Description: "The repeated alarm interval in minutes. Range: 0-1440.", Transform: transform.FromField("AlarmPeriod")},
			{Name: "monitor_object_type", Type: proto.ColumnType_INT, Description: "The monitored object type (0: shared monitored object for execution statements, 1: separate monitored object for each execution statement).", Transform: transform.FromField("MonitorObjectType")},
			{Name: "message_template", Type: proto.ColumnType_STRING, Description: "The custom notification template.", Transform: transform.FromField("MessageTemplate")},
			{Name: "group_trigger_status", Type: proto.ColumnType_BOOL, Description: "Whether group triggering is enabled (true: enabled, false: disabled).", Transform: transform.FromField("GroupTriggerStatus")},
			{Name: "topic_id", Type: proto.ColumnType_STRING, Description: "The log topic ID used to filter alarm policies; populated when specified in the query.", Transform: transform.FromQual("topic_id")},
			{Name: "alarm_targets", Type: proto.ColumnType_JSON, Description: "The list of monitoring objects (logset/topic/query) attached to the alarm policy.", Transform: transform.FromField("AlarmTargets")},
			{Name: "monitor_time", Type: proto.ColumnType_JSON, Description: "The monitoring task execution time point.", Transform: transform.FromField("MonitorTime")},
			{Name: "alarm_notice_ids", Type: proto.ColumnType_JSON, Description: "The list of associated alarm notification channel group IDs. Mutually exclusive with MonitorNotice.", Transform: transform.FromField("AlarmNoticeIds")},
			{Name: "call_back", Type: proto.ColumnType_JSON, Description: "The custom callback template.", Transform: transform.FromField("CallBack")},
			{Name: "analysis", Type: proto.ColumnType_JSON, Description: "The multidimensional analysis settings.", Transform: transform.FromField("Analysis")},
			{Name: "group_trigger_condition", Type: proto.ColumnType_JSON, Description: "The grouping trigger conditions.", Transform: transform.FromField("GroupTriggerCondition")},
			{Name: "classifications", Type: proto.ColumnType_JSON, Description: "Additional classification fields for the alarm.", Transform: transform.FromField("Classifications")},
			{Name: "multi_conditions", Type: proto.ColumnType_JSON, Description: "The multiple trigger conditions. Mutually exclusive with Condition.", Transform: transform.FromField("MultiConditions")},
			{Name: "monitor_notice", Type: proto.ColumnType_JSON, Description: "Tencent Cloud observability platform channel-related information. Mutually exclusive with AlarmNoticeIds.", Transform: transform.FromField("MonitorNotice")},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The tags attached to the alarm policy as a map of key -> value.", Transform: transform.FromField("Tags")},
			{Name: "create_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the alarm policy.", Transform: transform.FromField("CreateTime").NullIfZero()},
			{Name: "update_time", Type: proto.ColumnType_TIMESTAMP, Description: "The latest update time of the alarm policy.", Transform: transform.FromField("UpdateTime").NullIfZero()},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the alarm policy resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Name")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listClsAlarms(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listClsAlarms", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &cls.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listClsAlarms", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	filters := buildClsAlarmFilters(d.EqualsQuals)

	pageSize := utils.PageSizeFromLimit(d, 100, 1)
	var offset int64

	for {
		d.WaitForListRateLimit(ctx)

		req := cls.NewDescribeAlarmsRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		if len(filters) > 0 {
			req.Filters = filters
		}

		utils.LogRequest(ctx, "listClsAlarms", req)
		resp, err := client.DescribeAlarmsWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listClsAlarms", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		items := resp.Response.Alarms
		for _, item := range items {
			if item == nil {
				continue
			}
			row := toClsAlarmRow(item)
			row.Region = region
			d.StreamListItem(ctx, row)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if resp.Response == nil {
			break
		}
		total := utils.PtrInt64(resp.Response.TotalCount)
		if utils.PageDone(offset, int64(len(items)), total) {
			break
		}
		offset += int64(len(items))
	}

	return nil, nil
}

func buildClsAlarmFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*cls.Filter {
	return utils.BuildFilters(equalsQuals, newClsFilter,
		utils.FilterMapping{QualName: "alarm_id", FilterName: "alarmId"},
		utils.FilterMapping{QualName: "alarm_name", FilterName: "name"},
		utils.FilterMapping{QualName: "topic_id", FilterName: "topicId"},
	)
}

// Get Function

func getClsAlarm(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	alarmId := d.EqualsQuals["alarm_id"].GetStringValue()
	plugin.Logger(ctx).Debug("getClsAlarm", "alarm_id", alarmId, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &cls.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getClsAlarm", "client_init_error", err)
		return nil, err
	}

	if alarmId == "" {
		return nil, nil
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	req := cls.NewDescribeAlarmsRequest()
	req.Filters = []*cls.Filter{newClsFilter("alarmId", alarmId)}

	utils.LogRequest(ctx, "getClsAlarm", req)
	resp, err := client.DescribeAlarmsWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getClsAlarm", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.Alarms) == 0 {
		return nil, nil
	}

	row := toClsAlarmRow(resp.Response.Alarms[0])
	row.Region = region
	return row, nil
}

// Row Type

type clsAlarmRow struct {
	AlarmId               string
	Name                  string
	Status                bool
	AlarmLevel            uint64
	Condition             string
	TriggerCount          int64
	AlarmPeriod           int64
	MonitorObjectType     uint64
	MessageTemplate       string
	GroupTriggerStatus    bool
	AlarmTargets          []*cls.AlarmTargetInfo
	MonitorTime           *cls.MonitorTime
	AlarmNoticeIds        []string
	CallBack              *cls.CallBackInfo
	Analysis              []*cls.AnalysisDimensional
	GroupTriggerCondition []string
	Classifications       []*cls.AlarmClassification
	MultiConditions       []*cls.MultiCondition
	MonitorNotice         *cls.MonitorNotice
	Tags                  map[string]string
	CreateTime            *time.Time
	UpdateTime            *time.Time
	Region                string
	Akas                  []string
}

func toClsAlarmRow(a *cls.AlarmInfo) clsAlarmRow {
	row := clsAlarmRow{
		AlarmId:               utils.PtrString(a.AlarmId),
		Name:                  utils.PtrString(a.Name),
		Status:                utils.PtrBool(a.Status),
		AlarmLevel:            utils.PtrUint64(a.AlarmLevel),
		Condition:             utils.PtrString(a.Condition),
		TriggerCount:          utils.PtrInt64(a.TriggerCount),
		AlarmPeriod:           utils.PtrInt64(a.AlarmPeriod),
		MonitorObjectType:     utils.PtrUint64(a.MonitorObjectType),
		MessageTemplate:       utils.PtrString(a.MessageTemplate),
		GroupTriggerStatus:    utils.PtrBool(a.GroupTriggerStatus),
		AlarmTargets:          a.AlarmTargets,
		MonitorTime:           a.MonitorTime,
		AlarmNoticeIds:        utils.PtrStringSlice(a.AlarmNoticeIds),
		CallBack:              a.CallBack,
		Analysis:              a.Analysis,
		GroupTriggerCondition: utils.PtrStringSlice(a.GroupTriggerCondition),
		Classifications:       a.Classifications,
		MultiConditions:       a.MultiConditions,
		MonitorNotice:         a.MonitorNotice,
		Tags:                  TagsToMap(a.Tags),
		CreateTime:            utils.ParseTimestamp(a.CreateTime),
		UpdateTime:            utils.ParseTimestamp(a.UpdateTime),
	}

	if row.AlarmId != "" {
		row.Akas = []string{"tencentcloud:cls:alarm:" + row.AlarmId}
	}
	return row
}
