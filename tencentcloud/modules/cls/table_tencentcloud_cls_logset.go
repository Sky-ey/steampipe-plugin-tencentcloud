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

func TableTencentcloudClsLogset() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_cls_logset",
		Description:       "Tencent Cloud CLS logsets are containers that group related log topics and act as the top-level organizational unit for Cloud Log Service (CLS) in each region.",
		GetMatrixItemFunc: buildClsRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "logset_id", Require: plugin.Optional},
				{Name: "logset_name", Require: plugin.Optional},
			},
			Hydrate: listClsLogsets,
			Tags:    map[string]string{"service": "cls", "action": "DescribeLogsets"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("logset_id"),
			Hydrate:    getClsLogset,
			Tags:       map[string]string{"service": "cls", "action": "DescribeLogsets"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "logset_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the logset.", Transform: transform.FromField("LogsetId")},
			{Name: "logset_name", Type: proto.ColumnType_STRING, Description: "The name of the logset.", Transform: transform.FromField("LogsetName")},
			{Name: "topic_count", Type: proto.ColumnType_INT, Description: "The number of log topics in the logset.", Transform: transform.FromField("TopicCount")},
			{Name: "metric_topic_count", Type: proto.ColumnType_INT, Description: "The number of metric topics under the logset.", Transform: transform.FromField("MetricTopicCount")},
			{Name: "assumer_uin", Type: proto.ColumnType_INT, Description: "The Uin of the service provider that created the logset, if any.", Transform: transform.FromField("AssumerUin")},
			{Name: "assumer_name", Type: proto.ColumnType_STRING, Description: "The cloud product identifier that created the logset (e.g., CDN, TKE), if created by another cloud product.", Transform: transform.FromField("AssumerName")},
			{Name: "role_name", Type: proto.ColumnType_STRING, Description: "The role name of the service provider that created the logset, if any.", Transform: transform.FromField("RoleName")},
			{Name: "create_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the logset.", Transform: transform.FromField("CreateTime").NullIfZero()},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The tags attached to the logset as a map of key -> value.", Transform: transform.FromField("Tags")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the logset resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("LogsetName")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listClsLogsets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listClsLogsets", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &cls.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listClsLogsets", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	filters := buildClsLogsetFilters(d.EqualsQuals)

	pageSize := utils.PageSizeFromLimit(d, 100, 1)
	var offset int64

	for {
		d.WaitForListRateLimit(ctx)

		req := cls.NewDescribeLogsetsRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		if len(filters) > 0 {
			req.Filters = filters
		}

		utils.LogRequest(ctx, "listClsLogsets", req)
		resp, err := client.DescribeLogsetsWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listClsLogsets", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		items := resp.Response.Logsets
		for _, item := range items {
			if item == nil {
				continue
			}
			row := toClsLogsetRow(item)
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

func buildClsLogsetFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*cls.Filter {
	return utils.BuildFilters(equalsQuals, newClsFilter,
		utils.FilterMapping{QualName: "logset_id", FilterName: "logsetId"},
		utils.FilterMapping{QualName: "logset_name", FilterName: "logsetName"},
	)
}

// Get Function

func getClsLogset(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	logsetId := d.EqualsQuals["logset_id"].GetStringValue()
	plugin.Logger(ctx).Info("getClsLogset", "logset_id", logsetId, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &cls.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getClsLogset", "client_init_error", err)
		return nil, err
	}

	if logsetId == "" {
		return nil, nil
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	req := cls.NewDescribeLogsetsRequest()
	req.Filters = []*cls.Filter{newClsFilter("logsetId", logsetId)}

	utils.LogRequest(ctx, "getClsLogset", req)
	resp, err := client.DescribeLogsetsWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getClsLogset", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.Logsets) == 0 {
		return nil, nil
	}

	row := toClsLogsetRow(resp.Response.Logsets[0])
	row.Region = region
	return row, nil
}

// Row Type

type clsLogsetRow struct {
	LogsetId         string
	LogsetName       string
	TopicCount       int64
	MetricTopicCount int64
	AssumerUin       uint64
	AssumerName      string
	RoleName         string
	CreateTime       *time.Time
	Tags             map[string]string
	Region           string
	Akas             []string
}

func toClsLogsetRow(l *cls.LogsetInfo) clsLogsetRow {
	row := clsLogsetRow{
		LogsetId:         utils.PtrString(l.LogsetId),
		LogsetName:       utils.PtrString(l.LogsetName),
		TopicCount:       utils.PtrInt64(l.TopicCount),
		MetricTopicCount: utils.PtrInt64(l.MetricTopicCount),
		AssumerUin:       utils.PtrUint64(l.AssumerUin),
		AssumerName:      utils.PtrString(l.AssumerName),
		RoleName:         utils.PtrString(l.RoleName),
		CreateTime:       utils.ParseTimestamp(l.CreateTime),
		Tags:             TagsToMap(l.Tags),
	}

	if row.LogsetId != "" {
		row.Akas = []string{"tencentcloud:cls:logset:" + row.LogsetId}
	}
	return row
}
