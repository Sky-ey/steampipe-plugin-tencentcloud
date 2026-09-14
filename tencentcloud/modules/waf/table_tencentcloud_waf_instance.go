package waf

import (
	"context"
	"time"

	waf "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/waf/v20180125"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
)

// Table Definition

func TableTencentcloudWafInstance() *plugin.Table {
	return &plugin.Table{
		Name:        "tencentcloud_waf_instance",
		Description: "Tencent Cloud WAF instances (SaaS and CLB editions) managed under the Web Application Firewall service (document/api/627).",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "instance_id", Require: plugin.Optional},
				{Name: "instance_name", Require: plugin.Optional},
			},
			Hydrate: listWafInstances,
			Tags:    map[string]string{"service": "waf", "action": "DescribeInstances"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("instance_id"),
			Hydrate:    getWafInstance,
			Tags:       map[string]string{"service": "waf", "action": "DescribeInstances"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "instance_id", Type: proto.ColumnType_STRING, Description: "The unique instance ID.", Transform: transform.FromField("InstanceId")},
			{Name: "instance_name", Type: proto.ColumnType_STRING, Description: "The instance name.", Transform: transform.FromField("InstanceName")},
			{Name: "resource_ids", Type: proto.ColumnType_STRING, Description: "The resource ID corresponding to the instance, used for billing.", Transform: transform.FromField("ResourceIds")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The instance region reported by the API (business column, not a steampipe matrix region).", Transform: transform.FromField("Region")},
			{Name: "status", Type: proto.ColumnType_INT, Description: "The instance status.", Transform: transform.FromField("Status")},
			{Name: "edition", Type: proto.ColumnType_STRING, Description: "The instance edition (clb, saas).", Transform: transform.FromField("Edition")},
			{Name: "level", Type: proto.ColumnType_INT, Description: "The instance package version (101 = Small and Micro Edition, 102 = Ultra-light Edition, 2 = Advanced Edition, 3 = Enterprise Edition, 4 = Ultimate Edition, 6 = Exclusive Edition).", Transform: transform.FromField("Level")},
			{Name: "pay_mode", Type: proto.ColumnType_INT, Description: "The payment mode.", Transform: transform.FromField("PayMode")},
			{Name: "renew_flag", Type: proto.ColumnType_INT, Description: "The auto-renewal flag (0 = disabled, 1 = enabled).", Transform: transform.FromField("RenewFlag")},
			{Name: "mode", Type: proto.ColumnType_INT, Description: "The elastic billing switch (0 = disabled, 1 = enabled).", Transform: transform.FromField("Mode")},
			{Name: "valid_time", Type: proto.ColumnType_TIMESTAMP, Description: "The instance expiration time.", Transform: transform.FromField("ValidTime").NullIfZero()},
			{Name: "begin_time", Type: proto.ColumnType_TIMESTAMP, Description: "The instance start time.", Transform: transform.FromField("BeginTime").NullIfZero()},
			{Name: "domain_count", Type: proto.ColumnType_INT, Description: "The configured domain count.", Transform: transform.FromField("DomainCount")},
			{Name: "sub_domain_limit", Type: proto.ColumnType_INT, Description: "The maximum domain count (sub-domain limit).", Transform: transform.FromField("SubDomainLimit")},
			{Name: "main_domain_count", Type: proto.ColumnType_INT, Description: "The configured primary domain count.", Transform: transform.FromField("MainDomainCount")},
			{Name: "main_domain_limit", Type: proto.ColumnType_INT, Description: "The maximum primary domain count.", Transform: transform.FromField("MainDomainLimit")},
			{Name: "max_qps", Type: proto.ColumnType_INT, Description: "The instance QPS peak within the last 30 days.", Transform: transform.FromField("MaxQPS")},
			{Name: "qps", Type: proto.ColumnType_JSON, Description: "The QPS expansion package information (QPSPackageNew).", Transform: transform.FromField("QPS")},
			{Name: "domain_pkg", Type: proto.ColumnType_JSON, Description: "The domain extension package information (DomainPackageNew).", Transform: transform.FromField("DomainPkg")},
			{Name: "app_id", Type: proto.ColumnType_STRING, Description: "The user AppId.", Transform: transform.FromField("AppId")},
			{Name: "elastic_billing", Type: proto.ColumnType_INT, Description: "The QPS elastic billing cap.", Transform: transform.FromField("ElasticBilling")},
			{Name: "attack_log_post", Type: proto.ColumnType_INT, Description: "The attack log shipping switch.", Transform: transform.FromField("AttackLogPost")},
			{Name: "max_bandwidth", Type: proto.ColumnType_INT, Description: "The peak bandwidth in B/s (bytes per second).", Transform: transform.FromField("MaxBandwidth")},
			{Name: "api_security", Type: proto.ColumnType_INT, Description: "Whether API security is purchased.", Transform: transform.FromField("APISecurity")},
			{Name: "qps_standard", Type: proto.ColumnType_INT, Description: "The purchased QPS specification.", Transform: transform.FromField("QpsStandard")},
			{Name: "bandwidth_standard", Type: proto.ColumnType_INT, Description: "The purchased bandwidth specification.", Transform: transform.FromField("BandwidthStandard")},
			{Name: "free_delay_flag", Type: proto.ColumnType_INT, Description: "The flag for delay of instance deletion.", Transform: transform.FromField("FreeDelayFlag")},
			{Name: "sandbox_qps", Type: proto.ColumnType_INT, Description: "The instance sandbox QPS value.", Transform: transform.FromField("SandboxQps")},
			{Name: "is_api_security_trial", Type: proto.ColumnType_INT, Description: "Whether API security is on trial.", Transform: transform.FromField("IsAPISecurityTrial")},
			{Name: "mini_qps_standard", Type: proto.ColumnType_INT, Description: "The mini program QPS specification.", Transform: transform.FromField("MiniQpsStandard")},
			{Name: "mini_max_qps", Type: proto.ColumnType_INT, Description: "The mini program QPS peak.", Transform: transform.FromField("MiniMaxQPS")},
			{Name: "billing_item", Type: proto.ColumnType_STRING, Description: "The billing item.", Transform: transform.FromField("BillingItem")},
			{Name: "last_qps_exceed_time", Type: proto.ColumnType_STRING, Description: "The last QPS overage time.", Transform: transform.FromField("LastQpsExceedTime")},
			{Name: "fraud_pkg", Type: proto.ColumnType_JSON, Description: "The business security package (FraudPkg).", Transform: transform.FromField("FraudPkg")},
			{Name: "bot_pkg", Type: proto.ColumnType_JSON, Description: "The bot resource package (BotPkg).", Transform: transform.FromField("BotPkg")},
			{Name: "bot_qps", Type: proto.ColumnType_JSON, Description: "The bot QPS details (BotQPS).", Transform: transform.FromField("BotQPS")},
			{Name: "major_events_pkg", Type: proto.ColumnType_JSON, Description: "The premium package (MajorEventsPkg).", Transform: transform.FromField("MajorEventsPkg")},
			{Name: "hybrid_pkg", Type: proto.ColumnType_JSON, Description: "The hybrid cloud sub-node package (HybridPkg).", Transform: transform.FromField("HybridPkg")},
			{Name: "api_pkg", Type: proto.ColumnType_JSON, Description: "The API security resource package (ApiPkg).", Transform: transform.FromField("ApiPkg")},
			{Name: "mini_pkg", Type: proto.ColumnType_JSON, Description: "The MMPS acceleration package (MiniPkg).", Transform: transform.FromField("MiniPkg")},
			{Name: "mini_extend_pkg", Type: proto.ColumnType_JSON, Description: "The ID quantity expansion package for secure mini program access (MiniExtendPkg).", Transform: transform.FromField("MiniExtendPkg")},
			// No tags column: InstanceInfo has no Tags field (see waf/common.go note).
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listWafInstances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listWafInstances", "quals", d.EqualsQuals)

	client := &waf.Client{}
	if err := utils.InitClientInRegion(ctx, d, client, wafRequestRegion(d)); err != nil {
		plugin.Logger(ctx).Error("listWafInstances", "client_init_error", err)
		return nil, err
	}

	filters := utils.BuildFilters(d.EqualsQuals, newWafExactFilter,
		utils.FilterMapping{QualName: "instance_id", FilterName: "InstanceId"},
		utils.FilterMapping{QualName: "instance_name", FilterName: "InstanceName"},
	)

	pageSize := uint64(utils.PageSizeFromLimit(d, 100, 1))
	offset := uint64(0)

	for {
		d.WaitForListRateLimit(ctx)

		req := waf.NewDescribeInstancesRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		if len(filters) > 0 {
			req.Filters = filters
		}

		utils.LogRequest(ctx, "listWafInstances", req)
		resp, err := client.DescribeInstancesWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listWafInstances", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		items := resp.Response.Instances
		total := utils.PtrUint64(resp.Response.Total)

		for _, item := range items {
			if item == nil {
				continue
			}
			d.StreamListItem(ctx, toWafInstanceRow(item))
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if utils.PageDone(int64(offset), int64(len(items)), total) {
			break
		}
		offset += uint64(len(items))
	}

	return nil, nil
}

// Get Function

func getWafInstance(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	instanceId := d.EqualsQuals["instance_id"].GetStringValue()
	plugin.Logger(ctx).Debug("getWafInstance", "instance_id", instanceId)

	client := &waf.Client{}
	if err := utils.InitClientInRegion(ctx, d, client, wafRequestRegion(d)); err != nil {
		plugin.Logger(ctx).Error("getWafInstance", "client_init_error", err)
		return nil, err
	}

	if instanceId == "" {
		return nil, nil
	}

	req := waf.NewDescribeInstancesRequest()
	req.Filters = []*waf.FiltersItemNew{
		newWafFilter("InstanceId", instanceId, true),
	}

	utils.LogRequest(ctx, "getWafInstance", req)
	resp, err := client.DescribeInstancesWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getWafInstance", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.Instances) == 0 {
		return nil, nil
	}

	return toWafInstanceRow(resp.Response.Instances[0]), nil
}

// Row Type

type wafInstanceRow struct {
	InstanceId         string
	InstanceName       string
	ResourceIds        string
	Region             string
	Status             uint64
	Edition            string
	Level              uint64
	PayMode            uint64
	RenewFlag          uint64
	Mode               uint64
	ValidTime          *time.Time
	BeginTime          *time.Time
	DomainCount        uint64
	SubDomainLimit     uint64
	MainDomainCount    uint64
	MainDomainLimit    uint64
	MaxQPS             uint64
	QPS                *waf.QPSPackageNew
	DomainPkg          *waf.DomainPackageNew
	AppId              uint64
	ElasticBilling     uint64
	AttackLogPost      int64
	MaxBandwidth       uint64
	APISecurity        uint64
	QpsStandard        uint64
	BandwidthStandard  uint64
	FreeDelayFlag      uint64
	SandboxQps         uint64
	IsAPISecurityTrial uint64
	MiniQpsStandard    uint64
	MiniMaxQPS         uint64
	BillingItem        string
	LastQpsExceedTime  string
	FraudPkg           *waf.FraudPkg
	BotPkg             *waf.BotPkg
	BotQPS             *waf.BotQPS
	MajorEventsPkg     *waf.MajorEventsPkg
	HybridPkg          *waf.HybridPkg
	ApiPkg             *waf.ApiPkg
	MiniPkg            *waf.MiniPkg
	MiniExtendPkg      *waf.MiniExtendPkg
	Title              string
	Akas               []string
}

func toWafInstanceRow(i *waf.InstanceInfo) wafInstanceRow {
	row := wafInstanceRow{
		InstanceId:         utils.PtrString(i.InstanceId),
		InstanceName:       utils.PtrString(i.InstanceName),
		ResourceIds:        utils.PtrString(i.ResourceIds),
		Region:             utils.PtrString(i.Region),
		Status:             utils.PtrUint64(i.Status),
		Edition:            utils.PtrString(i.Edition),
		Level:              utils.PtrUint64(i.Level),
		PayMode:            utils.PtrUint64(i.PayMode),
		RenewFlag:          utils.PtrUint64(i.RenewFlag),
		Mode:               utils.PtrUint64(i.Mode),
		ValidTime:          utils.ParseTimestamp(i.ValidTime),
		BeginTime:          utils.ParseTimestamp(i.BeginTime),
		DomainCount:        utils.PtrUint64(i.DomainCount),
		SubDomainLimit:     utils.PtrUint64(i.SubDomainLimit),
		MainDomainCount:    utils.PtrUint64(i.MainDomainCount),
		MainDomainLimit:    utils.PtrUint64(i.MainDomainLimit),
		MaxQPS:             utils.PtrUint64(i.MaxQPS),
		QPS:                i.QPS,
		DomainPkg:          i.DomainPkg,
		AppId:              utils.PtrUint64(i.AppId),
		ElasticBilling:     utils.PtrUint64(i.ElasticBilling),
		AttackLogPost:      utils.PtrInt64(i.AttackLogPost),
		MaxBandwidth:       utils.PtrUint64(i.MaxBandwidth),
		APISecurity:        utils.PtrUint64(i.APISecurity),
		QpsStandard:        utils.PtrUint64(i.QpsStandard),
		BandwidthStandard:  utils.PtrUint64(i.BandwidthStandard),
		FreeDelayFlag:      utils.PtrUint64(i.FreeDelayFlag),
		SandboxQps:         utils.PtrUint64(i.SandboxQps),
		IsAPISecurityTrial: utils.PtrUint64(i.IsAPISecurityTrial),
		MiniQpsStandard:    utils.PtrUint64(i.MiniQpsStandard),
		MiniMaxQPS:         utils.PtrUint64(i.MiniMaxQPS),
		BillingItem:        utils.PtrString(i.BillingItem),
		LastQpsExceedTime:  utils.PtrString(i.LastQpsExceedTime),
		FraudPkg:           i.FraudPkg,
		BotPkg:             i.BotPkg,
		BotQPS:             i.BotQPS,
		MajorEventsPkg:     i.MajorEventsPkg,
		HybridPkg:          i.HybridPkg,
		ApiPkg:             i.ApiPkg,
		MiniPkg:            i.MiniPkg,
		MiniExtendPkg:      i.MiniExtendPkg,
	}

	if row.InstanceName != "" {
		row.Title = row.InstanceName
	} else {
		row.Title = row.InstanceId
	}

	if row.InstanceId != "" {
		row.Akas = []string{"tencentcloud:waf:instance:" + row.InstanceId}
	}
	return row
}
