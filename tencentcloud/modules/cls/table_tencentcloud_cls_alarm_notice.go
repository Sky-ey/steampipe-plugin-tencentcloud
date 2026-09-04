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

func TableTencentcloudClsAlarmNotice() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_cls_alarm_notice",
		Description:       "Tencent Cloud CLS alarm notification channel groups define how and where alarm notifications are delivered (recipients, webhooks, and shipping) within a region.",
		GetMatrixItemFunc: buildClsRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "alarm_notice_id", Require: plugin.Optional},
				{Name: "name", Require: plugin.Optional},
			},
			Hydrate: listClsAlarmNotices,
			Tags:    map[string]string{"service": "cls", "action": "DescribeAlarmNotices"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("alarm_notice_id"),
			Hydrate:    getClsAlarmNotice,
			Tags:       map[string]string{"service": "cls", "action": "DescribeAlarmNotices"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "alarm_notice_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the alarm notification channel group.", Transform: transform.FromField("AlarmNoticeId")},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the alarm notification channel group.", Transform: transform.FromField("Name")},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The alarm template type (Trigger - alarm trigger, Recovery - alarm recovery, All - alarm trigger and recovery).", Transform: transform.FromField("Type")},
			{Name: "alarm_shield_status", Type: proto.ColumnType_INT, Description: "The login-free operation alarm switch (1: off, 2: on).", Transform: transform.FromField("AlarmShieldStatus")},
			{Name: "deliver_status", Type: proto.ColumnType_INT, Description: "The delivery log switch (1: disabled, 2: enabled).", Transform: transform.FromField("DeliverStatus")},
			{Name: "deliver_flag", Type: proto.ColumnType_INT, Description: "The delivery log flag (1: disabled, 2: enabled, 3: delivery exception).", Transform: transform.FromField("DeliverFlag")},
			{Name: "jump_domain", Type: proto.ColumnType_STRING, Description: "The call link domain name. Must start with http:// or https:// and must not end with /.", Transform: transform.FromField("JumpDomain")},
			{Name: "callback_prioritize", Type: proto.ColumnType_BOOL, Description: "Whether to use custom callback parameters in the notification content template to override the request header and body (true: use template params, false: use alarm policy params).", Transform: transform.FromField("CallbackPrioritize")},
			{Name: "notice_receivers", Type: proto.ColumnType_JSON, Description: "The alarm notification template recipient information.", Transform: transform.FromField("NoticeReceivers")},
			{Name: "web_callbacks", Type: proto.ColumnType_JSON, Description: "The callback information of the alarm notification template.", Transform: transform.FromField("WebCallbacks")},
			{Name: "notice_rules", Type: proto.ColumnType_JSON, Description: "The notification rules.", Transform: transform.FromField("NoticeRules")},
			{Name: "alarm_notice_deliver_config", Type: proto.ColumnType_JSON, Description: "The shipping-related information.", Transform: transform.FromField("AlarmNoticeDeliverConfig")},
			{Name: "alarm_shield_count", Type: proto.ColumnType_JSON, Description: "The alarm silence status quantity information configured for the notification channel group.", Transform: transform.FromField("AlarmShieldCount")},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The tags attached to the alarm notification channel group as a map of key -> value.", Transform: transform.FromField("Tags")},
			{Name: "create_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the alarm notification channel group.", Transform: transform.FromField("CreateTime").NullIfZero()},
			{Name: "update_time", Type: proto.ColumnType_TIMESTAMP, Description: "The latest update time of the alarm notification channel group.", Transform: transform.FromField("UpdateTime").NullIfZero()},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the alarm notification channel group resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Name")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listClsAlarmNotices(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listClsAlarmNotices", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &cls.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listClsAlarmNotices", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	filters := buildClsAlarmNoticeFilters(d.EqualsQuals)

	pageSize := utils.PageSizeFromLimit(d, 100, 1)
	var offset int64

	for {
		d.WaitForListRateLimit(ctx)

		req := cls.NewDescribeAlarmNoticesRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		if len(filters) > 0 {
			req.Filters = filters
		}

		utils.LogRequest(ctx, "listClsAlarmNotices", req)
		resp, err := client.DescribeAlarmNoticesWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listClsAlarmNotices", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		items := resp.Response.AlarmNotices
		for _, item := range items {
			if item == nil {
				continue
			}
			row := toClsAlarmNoticeRow(item)
			row.Region = region
			d.StreamListItem(ctx, row)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		total := utils.PtrInt64(resp.Response.TotalCount)
		if utils.PageDone(offset, int64(len(items)), total) {
			break
		}
		offset += int64(len(items))
	}

	return nil, nil
}

func buildClsAlarmNoticeFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*cls.Filter {
	return utils.BuildFilters(equalsQuals, newClsFilter,
		utils.FilterMapping{QualName: "alarm_notice_id", FilterName: "alarmNoticeId"},
		utils.FilterMapping{QualName: "name", FilterName: "name"},
	)
}

// Get Function

func getClsAlarmNotice(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	alarmNoticeId := d.EqualsQuals["alarm_notice_id"].GetStringValue()
	plugin.Logger(ctx).Info("getClsAlarmNotice", "alarm_notice_id", alarmNoticeId, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &cls.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getClsAlarmNotice", "client_init_error", err)
		return nil, err
	}

	if alarmNoticeId == "" {
		return nil, nil
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	req := cls.NewDescribeAlarmNoticesRequest()
	req.Filters = []*cls.Filter{newClsFilter("alarmNoticeId", alarmNoticeId)}

	utils.LogRequest(ctx, "getClsAlarmNotice", req)
	resp, err := client.DescribeAlarmNoticesWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getClsAlarmNotice", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.AlarmNotices) == 0 {
		return nil, nil
	}

	row := toClsAlarmNoticeRow(resp.Response.AlarmNotices[0])
	row.Region = region
	return row, nil
}

// Row Type

type clsAlarmNoticeRow struct {
	AlarmNoticeId            string
	Name                     string
	Type                     string
	AlarmShieldStatus        uint64
	DeliverStatus            uint64
	DeliverFlag              uint64
	JumpDomain               string
	CallbackPrioritize       bool
	NoticeReceivers          []*cls.NoticeReceiver
	WebCallbacks             []*cls.WebCallback
	NoticeRules              []*cls.NoticeRule
	AlarmNoticeDeliverConfig *cls.AlarmNoticeDeliverConfig
	AlarmShieldCount         *cls.AlarmShieldCount
	Tags                     map[string]string
	CreateTime               *time.Time
	UpdateTime               *time.Time
	Region                   string
	Akas                     []string
}

func toClsAlarmNoticeRow(n *cls.AlarmNotice) clsAlarmNoticeRow {
	row := clsAlarmNoticeRow{
		AlarmNoticeId:            utils.PtrString(n.AlarmNoticeId),
		Name:                     utils.PtrString(n.Name),
		Type:                     utils.PtrString(n.Type),
		AlarmShieldStatus:        utils.PtrUint64(n.AlarmShieldStatus),
		DeliverStatus:            utils.PtrUint64(n.DeliverStatus),
		DeliverFlag:              utils.PtrUint64(n.DeliverFlag),
		JumpDomain:               utils.PtrString(n.JumpDomain),
		CallbackPrioritize:       utils.PtrBool(n.CallbackPrioritize),
		NoticeReceivers:          n.NoticeReceivers,
		WebCallbacks:             n.WebCallbacks,
		NoticeRules:              n.NoticeRules,
		AlarmNoticeDeliverConfig: n.AlarmNoticeDeliverConfig,
		AlarmShieldCount:         n.AlarmShieldCount,
		Tags:                     TagsToMap(n.Tags),
		CreateTime:               utils.ParseTimestamp(n.CreateTime),
		UpdateTime:               utils.ParseTimestamp(n.UpdateTime),
	}

	if row.AlarmNoticeId != "" {
		row.Akas = []string{"tencentcloud:cls:alarm-notice:" + row.AlarmNoticeId}
	}
	return row
}
