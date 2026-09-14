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

func TableTencentcloudClbCustomizedConfig() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_clb_customized_config",
		Description:       "Tencent Cloud Load Balancer (CLB) customized configurations.",
		GetMatrixItemFunc: buildClbRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "uconfig_id", Require: plugin.Optional},
				{Name: "config_type", Require: plugin.Optional},
				{Name: "config_name", Require: plugin.Optional},
				{Name: "load_balancer_id", Require: plugin.Optional},
				{Name: "vip", Require: plugin.Optional},
			},
			Hydrate: listClbCustomizedConfigs,
			Tags:    map[string]string{"service": "clb", "action": "DescribeCustomizedConfigList"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("uconfig_id"),
			Hydrate:    getClbCustomizedConfig,
			Tags:       map[string]string{"service": "clb", "action": "DescribeCustomizedConfigList"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "uconfig_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the customized configuration.", Transform: transform.FromField("UconfigId")},
			{Name: "config_type", Type: proto.ColumnType_STRING, Description: "The type of the configuration (CLB, SERVER, LOCATION)."},
			{Name: "config_name", Type: proto.ColumnType_STRING, Description: "The name of the configuration."},
			{Name: "config_content", Type: proto.ColumnType_STRING, Description: "The content of the configuration."},
			{Name: "create_timestamp", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the configuration, format YYYY-MM-DD HH:MM:ss.", Transform: transform.FromField("CreateTimestamp").NullIfZero()},
			{Name: "update_timestamp", Type: proto.ColumnType_TIMESTAMP, Description: "The modification time of the configuration, format YYYY-MM-DD HH:MM:ss.", Transform: transform.FromField("UpdateTimestamp").NullIfZero()},
			{Name: "load_balancer_id", Type: proto.ColumnType_STRING, Description: "Filter-only qual: the CLB instance ID used to filter customized configs (passed through to DescribeCustomizedConfigList filter loadbalancer-id). Echoes the qual value when provided; not part of the API response payload.", Transform: transform.FromField("LoadBalancerId")},
			{Name: "vip", Type: proto.ColumnType_STRING, Description: "Filter-only qual: the CLB VIP used to filter customized configs (passed through to DescribeCustomizedConfigList filter vip). Echoes the qual value when provided; not part of the API response payload.", Transform: transform.FromField("Vip")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the customized configuration resides."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource."},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource."},
		}),
	}
}

// List Function

func listClbCustomizedConfigs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listClbCustomizedConfigs", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &clb.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listClbCustomizedConfigs", "client_init_error", err)
		return nil, err
	}
	region := d.EqualsQualString(utils.MatrixKeyRegion)
	loadBalancerID := d.EqualsQualString("load_balancer_id")
	vip := d.EqualsQualString("vip")
	uconfigID := d.EqualsQualString("uconfig_id")

	if uconfigID != "" {
		return nil, streamClbCustomizedConfigs(ctx, d, client, region, loadBalancerID, vip, func(req *clb.DescribeCustomizedConfigListRequest) {
			req.UconfigIds = common.StringPtrs([]string{uconfigID})
		})
	}

	filters := utils.BuildFilters(d.EqualsQuals, newClbFilter,
		utils.FilterMapping{QualName: "load_balancer_id", FilterName: "loadbalancer-id"},
		utils.FilterMapping{QualName: "vip", FilterName: "vip"},
	)
	configName := d.EqualsQualString("config_name")
	configTypes := []string{"CLB", "SERVER", "LOCATION"}
	if ct := d.EqualsQualString("config_type"); ct != "" {
		configTypes = []string{ct}
	}

	for _, configType := range configTypes {
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
		ct := configType
		if err := streamClbCustomizedConfigs(ctx, d, client, region, loadBalancerID, vip, func(req *clb.DescribeCustomizedConfigListRequest) {
			req.ConfigType = &ct
			if len(filters) > 0 {
				req.Filters = filters
			}
			if configName != "" {
				req.ConfigName = &configName
			}
		}); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

// streamClbCustomizedConfigs 执行 DescribeCustomizedConfigList 的分页拉取并 stream 行。
// configureReq 回调负责设置每页请求的业务字段(UconfigIds 或 ConfigType+Filters)。
func streamClbCustomizedConfigs(ctx context.Context, d *plugin.QueryData, client *clb.Client,
	region, loadBalancerID, vip string, configureReq func(*clb.DescribeCustomizedConfigListRequest)) error {

	pageSize := utils.PageSizeFromLimit(d, 100, 1)
	var offset int64
	for {
		d.WaitForListRateLimit(ctx)
		req := clb.NewDescribeCustomizedConfigListRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		configureReq(req)

		utils.LogRequest(ctx, "listClbCustomizedConfigs", req)
		resp, err := client.DescribeCustomizedConfigListWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listClbCustomizedConfigs", "request_error", err)
			return err
		}
		if resp == nil || resp.Response == nil {
			return nil
		}

		configs := resp.Response.ConfigList
		for _, config := range configs {
			if config == nil {
				continue
			}
			row := toClbCustomizedConfigRow(config, region, loadBalancerID, vip)
			d.StreamListItem(ctx, row)
			if d.RowsRemaining(ctx) == 0 {
				return nil
			}
		}

		if utils.PageDone(offset, int64(len(configs)), utils.PtrInt64(resp.Response.TotalCount)) {
			return nil
		}
		offset += int64(len(configs))
	}
}

// Get Function

func getClbCustomizedConfig(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	uconfigID := d.EqualsQualString("uconfig_id")
	plugin.Logger(ctx).Debug("getClbCustomizedConfig", "uconfig_id", uconfigID, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &clb.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getClbCustomizedConfig", "client_init_error", err)
		return nil, err
	}

	if uconfigID == "" {
		return nil, nil
	}
	req := clb.NewDescribeCustomizedConfigListRequest()
	req.UconfigIds = common.StringPtrs([]string{uconfigID})

	utils.LogRequest(ctx, "getClbCustomizedConfig", req)
	resp, err := client.DescribeCustomizedConfigListWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getClbCustomizedConfig", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.ConfigList) == 0 || resp.Response.ConfigList[0] == nil {
		return nil, nil
	}
	return toClbCustomizedConfigRow(resp.Response.ConfigList[0], d.EqualsQualString(utils.MatrixKeyRegion), "", ""), nil
}

// Row Type

type clbCustomizedConfigRow struct {
	UconfigId       string
	ConfigType      string
	ConfigName      string
	ConfigContent   string
	CreateTimestamp *time.Time
	UpdateTimestamp *time.Time
	LoadBalancerId  string
	Vip             string
	Region          string
	Title           string
	Akas            []string
}

func toClbCustomizedConfigRow(config *clb.ConfigListItem, region, loadBalancerID, vip string) clbCustomizedConfigRow {
	row := clbCustomizedConfigRow{
		UconfigId:       utils.PtrString(config.UconfigId),
		ConfigType:      utils.PtrString(config.ConfigType),
		ConfigName:      utils.PtrString(config.ConfigName),
		ConfigContent:   utils.PtrString(config.ConfigContent),
		CreateTimestamp: utils.ParseTimestamp(config.CreateTimestamp),
		UpdateTimestamp: utils.ParseTimestamp(config.UpdateTimestamp),
		LoadBalancerId:  loadBalancerID,
		Vip:             vip,
		Region:          region,
	}
	row.Title = row.ConfigName
	if row.Title == "" {
		row.Title = row.UconfigId
	}
	if row.UconfigId != "" {
		row.Akas = []string{"tencentcloud:clb:customized-config:" + row.UconfigId}
	}
	return row
}
