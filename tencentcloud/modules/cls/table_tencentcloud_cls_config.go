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

func TableTencentcloudClsConfig() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_cls_config",
		Description:       "Tencent Cloud CLS collection configs define how logs are collected from paths, parsed, and shipped to a log topic within a region.",
		GetMatrixItemFunc: buildClsRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "config_id", Require: plugin.Optional},
				{Name: "config_name", Require: plugin.Optional},
				{Name: "topic_id", Require: plugin.Optional},
			},
			Hydrate: listClsConfigs,
			Tags:    map[string]string{"service": "cls", "action": "DescribeConfigs"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("config_id"),
			Hydrate:    getClsConfig,
			Tags:       map[string]string{"service": "cls", "action": "DescribeConfigs"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "config_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the collection configuration.", Transform: transform.FromField("ConfigId")},
			{Name: "config_name", Type: proto.ColumnType_STRING, Description: "The name of the collection configuration.", Transform: transform.FromField("Name")},
			{Name: "topic_id", Type: proto.ColumnType_STRING, Description: "The log topic ID (SDK field Output) that this collection configuration ships to.", Transform: transform.FromField("Output")},
			{Name: "log_type", Type: proto.ColumnType_STRING, Description: "The type of log collected (json_log, delimiter_log, minimalist_log, fullregex_log, multiline_log, multiline_fullregex_log, user_define_log, service_syslog, windows_event_log).", Transform: transform.FromField("LogType")},
			{Name: "log_format", Type: proto.ColumnType_STRING, Description: "The log formatting method.", Transform: transform.FromField("LogFormat")},
			{Name: "path", Type: proto.ColumnType_STRING, Description: "The log collection path.", Transform: transform.FromField("Path")},
			{Name: "input_type", Type: proto.ColumnType_STRING, Description: "The log input type (file, windows_event, syslog).", Transform: transform.FromField("InputType")},
			{Name: "extract_rule", Type: proto.ColumnType_JSON, Description: "The extraction rule for the collected log.", Transform: transform.FromField("ExtractRule")},
			{Name: "exclude_paths", Type: proto.ColumnType_JSON, Description: "The collection path blocklist.", Transform: transform.FromField("ExcludePaths")},
			{Name: "user_define_rule", Type: proto.ColumnType_STRING, Description: "The custom parsing string used in combined parsing extraction mode.", Transform: transform.FromField("UserDefineRule")},
			{Name: "advanced_config", Type: proto.ColumnType_JSON, Description: "The advanced collection configuration as a JSON string (ClsAgentFileTimeout, ClsAgentMaxDepth, ClsAgentParseFailMerge).", Transform: transform.FromField("AdvancedConfig")},
			{Name: "create_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the collection configuration.", Transform: transform.FromField("CreateTime").NullIfZero()},
			{Name: "update_time", Type: proto.ColumnType_TIMESTAMP, Description: "The latest update time of the collection configuration.", Transform: transform.FromField("UpdateTime").NullIfZero()},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the collection configuration resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Name")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listClsConfigs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listClsConfigs", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &cls.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listClsConfigs", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	filters := buildClsConfigFilters(d.EqualsQuals)

	pageSize := utils.PageSizeFromLimit(d, 100, 1)
	var offset int64

	for {
		d.WaitForListRateLimit(ctx)

		req := cls.NewDescribeConfigsRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		if len(filters) > 0 {
			req.Filters = filters
		}

		utils.LogRequest(ctx, "listClsConfigs", req)
		resp, err := client.DescribeConfigsWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listClsConfigs", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		items := resp.Response.Configs
		for _, item := range items {
			if item == nil {
				continue
			}
			row := toClsConfigRow(item)
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

func buildClsConfigFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*cls.Filter {
	return utils.BuildFilters(equalsQuals, newClsFilter,
		utils.FilterMapping{QualName: "config_id", FilterName: "configId"},
		utils.FilterMapping{QualName: "config_name", FilterName: "configName"},
		utils.FilterMapping{QualName: "topic_id", FilterName: "topicId"},
	)
}

// Get Function

func getClsConfig(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	configId := d.EqualsQuals["config_id"].GetStringValue()
	plugin.Logger(ctx).Debug("getClsConfig", "config_id", configId, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &cls.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getClsConfig", "client_init_error", err)
		return nil, err
	}

	if configId == "" {
		return nil, nil
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	req := cls.NewDescribeConfigsRequest()
	req.Filters = []*cls.Filter{newClsFilter("configId", configId)}

	utils.LogRequest(ctx, "getClsConfig", req)
	resp, err := client.DescribeConfigsWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getClsConfig", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.Configs) == 0 {
		return nil, nil
	}

	row := toClsConfigRow(resp.Response.Configs[0])
	row.Region = region
	return row, nil
}

// Row Type

type clsConfigRow struct {
	ConfigId       string
	Name           string
	Output         string
	LogType        string
	LogFormat      string
	Path           string
	InputType      string
	ExtractRule    *cls.ExtractRuleInfo
	ExcludePaths   []*cls.ExcludePathInfo
	UserDefineRule string
	AdvancedConfig *string
	CreateTime     *time.Time
	UpdateTime     *time.Time
	Region         string
	Akas           []string
}

func toClsConfigRow(c *cls.ConfigInfo) clsConfigRow {
	row := clsConfigRow{
		ConfigId:       utils.PtrString(c.ConfigId),
		Name:           utils.PtrString(c.Name),
		Output:         utils.PtrString(c.Output),
		LogType:        utils.PtrString(c.LogType),
		LogFormat:      utils.PtrString(c.LogFormat),
		Path:           utils.PtrString(c.Path),
		InputType:      utils.PtrString(c.InputType),
		ExtractRule:    c.ExtractRule,
		ExcludePaths:   c.ExcludePaths,
		UserDefineRule: utils.PtrString(c.UserDefineRule),
		AdvancedConfig: c.AdvancedConfig,
		CreateTime:     utils.ParseTimestamp(c.CreateTime),
		UpdateTime:     utils.ParseTimestamp(c.UpdateTime),
	}

	if row.ConfigId != "" {
		row.Akas = []string{"tencentcloud:cls:config:" + row.ConfigId}
	}
	return row
}
