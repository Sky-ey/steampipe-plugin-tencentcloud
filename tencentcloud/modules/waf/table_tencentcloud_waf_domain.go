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

func TableTencentcloudWafDomain() *plugin.Table {
	return &plugin.Table{
		Name:        "tencentcloud_waf_domain",
		Description: "Tencent Cloud WAF protected domains (SaaS edition) managed under the Web Application Firewall service (document/api/627).",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "domain", Require: plugin.Optional},
				{Name: "domain_id", Require: plugin.Optional},
				{Name: "instance_id", Require: plugin.Optional},
			},
			Hydrate: listWafDomains,
			Tags:    map[string]string{"service": "waf", "action": "DescribeDomains"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("domain_id"),
			Hydrate:    getWafDomain,
			Tags:       map[string]string{"service": "waf", "action": "DescribeDomains"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "domain", Type: proto.ColumnType_STRING, Description: "The domain name.", Transform: transform.FromField("Domain")},
			{Name: "domain_id", Type: proto.ColumnType_STRING, Description: "The unique domain ID.", Transform: transform.FromField("DomainId")},
			{Name: "instance_id", Type: proto.ColumnType_STRING, Description: "The instance ID bound to the domain.", Transform: transform.FromField("InstanceId")},
			{Name: "instance_name", Type: proto.ColumnType_STRING, Description: "The instance name bound to the domain.", Transform: transform.FromField("InstanceName")},
			{Name: "cname", Type: proto.ColumnType_STRING, Description: "The CNAME address of the domain.", Transform: transform.FromField("Cname")},
			{Name: "edition", Type: proto.ColumnType_STRING, Description: "The instance type of the domain (sparta-waf = SaaS WAF, clb-waf = CLB WAF, cdc-clb-waf = CLB WAF in CDC environment).", Transform: transform.FromField("Edition")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The domain region reported by the API (business column, not a steampipe matrix region).", Transform: transform.FromField("Region")},
			{Name: "status", Type: proto.ColumnType_INT, Description: "The WAF switch status (0 = disabled, 1 = enabled).", Transform: transform.FromField("Status")},
			{Name: "cls_status", Type: proto.ColumnType_INT, Description: "The access log switch status (0 = disabled, 1 = enabled).", Transform: transform.FromField("ClsStatus")},
			{Name: "flow_mode", Type: proto.ColumnType_INT, Description: "The CLB WAF usage mode (0 = image mode, 1 = cleaning mode).", Transform: transform.FromField("FlowMode")},
			{Name: "mode", Type: proto.ColumnType_INT, Description: "The rule engine protection mode (0 = observation, 1 = interception).", Transform: transform.FromField("Mode")},
			{Name: "engine", Type: proto.ColumnType_INT, Description: "The joint status of rule engine and AI engine protection modes (1, 10, 11, 12, 20, 21, 22).", Transform: transform.FromField("Engine")},
			{Name: "state", Type: proto.ColumnType_INT, Description: "The LB listener status of the CLB WAF protected domain (0, 4, 6, 7, 8, 10).", Transform: transform.FromField("State")},
			{Name: "level", Type: proto.ColumnType_INT, Description: "The instance version (101 = Small and Micro Agile, 102 = Ultra-light, 2 = Advanced, 3 = Enterprise, 4 = Ultimate, 6 = Exclusive).", Transform: transform.FromField("Level")},
			{Name: "ipv6_status", Type: proto.ColumnType_INT, Description: "The IPv6 switch status (0 = disabled, 1 = enabled).", Transform: transform.FromField("Ipv6Status")},
			{Name: "bot_status", Type: proto.ColumnType_INT, Description: "The Bot switch status (0, 1 = disabled, 2, 3 = enabled).", Transform: transform.FromField("BotStatus")},
			{Name: "api_status", Type: proto.ColumnType_INT, Description: "The API security switch status (0 = disabled, 1 = enabled).", Transform: transform.FromField("ApiStatus")},
			{Name: "post_cls_status", Type: proto.ColumnType_INT, Description: "The CLS shipping status (0 = disabled, 1 = enabled).", Transform: transform.FromField("PostCLSStatus")},
			{Name: "post_ckafka_status", Type: proto.ColumnType_INT, Description: "The CKafka shipping status (0 = disabled, 1 = enabled).", Transform: transform.FromField("PostCKafkaStatus")},
			{Name: "sg_state", Type: proto.ColumnType_INT, Description: "The security group status (0 = not display, 1 = non-Tencent Cloud origin, 2 = binding failed, 3 = changed).", Transform: transform.FromField("SgState")},
			{Name: "sg_id", Type: proto.ColumnType_STRING, Description: "The security group ID.", Transform: transform.FromField("SgID")},
			{Name: "sg_detail", Type: proto.ColumnType_STRING, Description: "The detailed explanation of the security group status.", Transform: transform.FromField("SgDetail")},
			{Name: "access_status", Type: proto.ColumnType_INT, Description: "The CLB WAF access status.", Transform: transform.FromField("AccessStatus")},
			{Name: "alb_type", Type: proto.ColumnType_STRING, Description: "The application-based CLB type (clb = layer-7 CLB, apisix = APISIX gateway).", Transform: transform.FromField("AlbType")},
			{Name: "cloud_type", Type: proto.ColumnType_STRING, Description: "The domain cloud environment (hybrid = hybrid cloud, public = public cloud).", Transform: transform.FromField("CloudType")},
			{Name: "cdc_clusters", Type: proto.ColumnType_STRING, Description: "The CDC cluster information accessed by the CDC instance domain (ignored for non-CDC instances).", Transform: transform.FromField("CdcClusters")},
			{Name: "note", Type: proto.ColumnType_STRING, Description: "The domain remarks.", Transform: transform.FromField("Note")},
			{Name: "create_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time.", Transform: transform.FromField("CreateTime").NullIfZero()},
			{Name: "app_id", Type: proto.ColumnType_INT, Description: "The user AppId.", Transform: transform.FromField("AppId")},
			{Name: "cc_list", Type: proto.ColumnType_JSON, Description: "The sandbox cluster origin-pull outbound IP list.", Transform: transform.FromField("CCList")},
			{Name: "rs_list", Type: proto.ColumnType_JSON, Description: "The production cluster origin-pull outbound IP list.", Transform: transform.FromField("RsList")},
			{Name: "src_list", Type: proto.ColumnType_JSON, Description: "The SaaS WAF origin server IP list.", Transform: transform.FromField("SrcList")},
			{Name: "upstream_domain_list", Type: proto.ColumnType_JSON, Description: "The SaaS WAF origin server domain name list.", Transform: transform.FromField("UpstreamDomainList")},
			{Name: "ports", Type: proto.ColumnType_JSON, Description: "The service port configuration (PortInfo list).", Transform: transform.FromField("Ports")},
			{Name: "load_balancer_set", Type: proto.ColumnType_JSON, Description: "The CLB-related configuration (LoadBalancerPackageNew list).", Transform: transform.FromField("LoadBalancerSet")},
			{Name: "labels", Type: proto.ColumnType_JSON, Description: "The domain group tags (labels) attached to the domain as a list of strings.", Transform: transform.FromField("Labels")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listWafDomains(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listWafDomains", "quals", d.EqualsQuals)

	client := &waf.Client{}
	if err := utils.InitClientInRegion(ctx, d, client, wafRequestRegion(d)); err != nil {
		plugin.Logger(ctx).Error("listWafDomains", "client_init_error", err)
		return nil, err
	}

	filters := utils.BuildFilters(d.EqualsQuals, newWafExactFilter,
		utils.FilterMapping{QualName: "domain", FilterName: "Domain"},
		utils.FilterMapping{QualName: "domain_id", FilterName: "DomainId"},
		utils.FilterMapping{QualName: "instance_id", FilterName: "InstanceId"},
	)

	pageSize := uint64(utils.PageSizeFromLimit(d, 100, 1))
	offset := uint64(0)

	for {
		d.WaitForListRateLimit(ctx)

		req := waf.NewDescribeDomainsRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		if len(filters) > 0 {
			req.Filters = filters
		}

		utils.LogRequest(ctx, "listWafDomains", req)
		resp, err := client.DescribeDomainsWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listWafDomains", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		items := resp.Response.Domains
		total := utils.PtrUint64(resp.Response.Total)

		for _, item := range items {
			if item == nil {
				continue
			}
			d.StreamListItem(ctx, toWafDomainRow(item))
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

func getWafDomain(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	domainId := d.EqualsQuals["domain_id"].GetStringValue()
	plugin.Logger(ctx).Info("getWafDomain", "domain_id", domainId)

	client := &waf.Client{}
	if err := utils.InitClientInRegion(ctx, d, client, wafRequestRegion(d)); err != nil {
		plugin.Logger(ctx).Error("getWafDomain", "client_init_error", err)
		return nil, err
	}

	if domainId == "" {
		return nil, nil
	}

	req := waf.NewDescribeDomainsRequest()
	req.Filters = []*waf.FiltersItemNew{
		newWafFilter("DomainId", domainId, true),
	}

	utils.LogRequest(ctx, "getWafDomain", req)
	resp, err := client.DescribeDomainsWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getWafDomain", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.Domains) == 0 {
		return nil, nil
	}

	return toWafDomainRow(resp.Response.Domains[0]), nil
}

// Row Type

type wafDomainRow struct {
	Domain             string
	DomainId           string
	InstanceId         string
	InstanceName       string
	Cname              string
	Edition            string
	Region             string
	Status             uint64
	ClsStatus          uint64
	FlowMode           uint64
	Mode               uint64
	Engine             uint64
	State              int64
	Level              int64
	Ipv6Status         int64
	BotStatus          int64
	ApiStatus          int64
	PostCLSStatus      int64
	PostCKafkaStatus   int64
	SgState            int64
	SgID               string
	SgDetail           string
	AccessStatus       int64
	AlbType            string
	CloudType          string
	CdcClusters        string
	Note               string
	CreateTime         *time.Time
	AppId              uint64
	CCList             []string
	RsList             []string
	SrcList            []string
	UpstreamDomainList []string
	Ports              []*waf.PortInfo
	LoadBalancerSet    []*waf.LoadBalancerPackageNew
	Labels             []string
	Title              string
	Akas               []string
}

func toWafDomainRow(d *waf.DomainInfo) wafDomainRow {
	row := wafDomainRow{
		Domain:             utils.PtrString(d.Domain),
		DomainId:           utils.PtrString(d.DomainId),
		InstanceId:         utils.PtrString(d.InstanceId),
		InstanceName:       utils.PtrString(d.InstanceName),
		Cname:              utils.PtrString(d.Cname),
		Edition:            utils.PtrString(d.Edition),
		Region:             utils.PtrString(d.Region),
		Status:             utils.PtrUint64(d.Status),
		ClsStatus:          utils.PtrUint64(d.ClsStatus),
		FlowMode:           utils.PtrUint64(d.FlowMode),
		Mode:               utils.PtrUint64(d.Mode),
		Engine:             utils.PtrUint64(d.Engine),
		State:              utils.PtrInt64(d.State),
		Level:              utils.PtrInt64(d.Level),
		Ipv6Status:         utils.PtrInt64(d.Ipv6Status),
		BotStatus:          utils.PtrInt64(d.BotStatus),
		ApiStatus:          utils.PtrInt64(d.ApiStatus),
		PostCLSStatus:      utils.PtrInt64(d.PostCLSStatus),
		PostCKafkaStatus:   utils.PtrInt64(d.PostCKafkaStatus),
		SgState:            utils.PtrInt64(d.SgState),
		SgID:               utils.PtrString(d.SgID),
		SgDetail:           utils.PtrString(d.SgDetail),
		AccessStatus:       utils.PtrInt64(d.AccessStatus),
		AlbType:            utils.PtrString(d.AlbType),
		CloudType:          utils.PtrString(d.CloudType),
		CdcClusters:        utils.PtrString(d.CdcClusters),
		Note:               utils.PtrString(d.Note),
		CreateTime:         utils.ParseTimestamp(d.CreateTime),
		AppId:              utils.PtrUint64(d.AppId),
		CCList:             utils.PtrStringSlice(d.CCList),
		RsList:             utils.PtrStringSlice(d.RsList),
		SrcList:            utils.PtrStringSlice(d.SrcList),
		UpstreamDomainList: utils.PtrStringSlice(d.UpstreamDomainList),
		Ports:              d.Ports,
		LoadBalancerSet:    d.LoadBalancerSet,
		Labels:             utils.PtrStringSlice(d.Labels),
	}

	if row.Domain != "" {
		row.Title = row.Domain
	} else {
		row.Title = row.DomainId
	}

	if row.DomainId != "" {
		row.Akas = []string{"tencentcloud:waf:domain:" + row.DomainId}
	}
	return row
}
