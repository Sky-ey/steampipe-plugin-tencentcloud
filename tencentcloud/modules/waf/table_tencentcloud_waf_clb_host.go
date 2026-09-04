package waf

import (
	"context"

	waf "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/waf/v20180125"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
)

// Table Definition

func TableTencentcloudWafClbHost() *plugin.Table {
	return &plugin.Table{
		Name:        "tencentcloud_waf_clb_host",
		Description: "Tencent Cloud WAF CLB-protected domains (CLB edition) managed under the Web Application Firewall service (document/api/627).",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "domain", Require: plugin.Optional},
				{Name: "domain_id", Require: plugin.Optional},
				{Name: "instance_id", Require: plugin.Optional},
			},
			Hydrate: listWafClbHosts,
			Tags:    map[string]string{"service": "waf", "action": "DescribeHosts"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "domain", Type: proto.ColumnType_STRING, Description: "The protection domain name.", Transform: transform.FromField("Domain")},
			{Name: "domain_id", Type: proto.ColumnType_STRING, Description: "The unique protection domain ID.", Transform: transform.FromField("DomainId")},
			{Name: "main_domain", Type: proto.ColumnType_STRING, Description: "The primary domain (empty upon input).", Transform: transform.FromField("MainDomain")},
			{Name: "instance_id", Type: proto.ColumnType_STRING, Description: "The WAF instance ID used to filter hosts; populated when specified in the query.", Transform: transform.FromQual("instance_id")},
			{Name: "edition", Type: proto.ColumnType_STRING, Description: "The instance type of the domain (clb-waf for CLB WAF).", Transform: transform.FromField("Edition")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region of the CLB instance bound to the domain (comma-separated for multiple regions; business column, not a steampipe matrix region).", Transform: transform.FromField("Region")},
			{Name: "status", Type: proto.ColumnType_INT, Description: "The binding relationship between WAF and CLB (0 = unbound, 1 = bound).", Transform: transform.FromField("Status")},
			{Name: "state", Type: proto.ColumnType_INT, Description: "The domain listener status of CLB WAF (0, 4, 6, 7, 8, 10).", Transform: transform.FromField("State")},
			{Name: "mode", Type: proto.ColumnType_INT, Description: "The rule engine protection mode (0 = observation, 1 = interception).", Transform: transform.FromField("Mode")},
			{Name: "engine", Type: proto.ColumnType_INT, Description: "The joint status of rule engine and AI engine protection modes (1, 10, 11, 12, 20, 21, 22).", Transform: transform.FromField("Engine")},
			{Name: "engine_type", Type: proto.ColumnType_INT, Description: "The rule engine type (1 = menshen, 2 = tiga).", Transform: transform.FromField("EngineType")},
			{Name: "flow_mode", Type: proto.ColumnType_INT, Description: "The traffic mode of CLB WAF protected domains (0 = mirror, 1 = cleaning).", Transform: transform.FromField("FlowMode")},
			{Name: "cls_status", Type: proto.ColumnType_INT, Description: "The access logging switch (0 = disabled, 1 = enabled).", Transform: transform.FromField("ClsStatus")},
			{Name: "is_cdn", Type: proto.ColumnType_INT, Description: "Whether a layer-7 proxy is deployed before WAF (0 = no proxy, 1 = XFF, 2 = remote_addr, 3 = custom header).", Transform: transform.FromField("IsCdn")},
			{Name: "level", Type: proto.ColumnType_INT, Description: "The protection level (100, 200, 300).", Transform: transform.FromField("Level")},
			{Name: "alb_type", Type: proto.ColumnType_STRING, Description: "The application CLB type (clb = layer-7 CLB, tsegw = API Gateway, scf = Serverless, apisix = other gateway).", Transform: transform.FromField("AlbType")},
			{Name: "cloud_type", Type: proto.ColumnType_STRING, Description: "The cloud type (public, private, hybrid).", Transform: transform.FromField("CloudType")},
			{Name: "note", Type: proto.ColumnType_STRING, Description: "The domain remarks.", Transform: transform.FromField("Note")},
			{Name: "ip_headers", Type: proto.ColumnType_JSON, Description: "The custom header(s) used to obtain the client IP when IsCdn=3.", Transform: transform.FromField("IpHeaders")},
			{Name: "cdc_clusters", Type: proto.ColumnType_JSON, Description: "The list of CDC clusters to which the domain is delivered (CDC scenes only).", Transform: transform.FromField("CdcClusters")},
			{Name: "load_balancer_set", Type: proto.ColumnType_JSON, Description: "The list of bound CLB instance information (LoadBalancer list).", Transform: transform.FromField("LoadBalancerSet")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listWafClbHosts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listWafClbHosts", "quals", d.EqualsQuals)

	client := &waf.Client{}
	if err := utils.InitClientInRegion(ctx, d, client, wafRequestRegion(d)); err != nil {
		plugin.Logger(ctx).Error("listWafClbHosts", "client_init_error", err)
		return nil, err
	}

	req := waf.NewDescribeHostsRequest()
	if v := d.EqualsQuals["domain"]; v != nil {
		if s := v.GetStringValue(); s != "" {
			req.Domain = &s
		}
	}
	if v := d.EqualsQuals["domain_id"]; v != nil {
		if s := v.GetStringValue(); s != "" {
			req.DomainId = &s
		}
	}
	if v := d.EqualsQuals["instance_id"]; v != nil {
		if s := v.GetStringValue(); s != "" {
			req.InstanceID = &s
		}
	}

	utils.LogRequest(ctx, "listWafClbHosts", req)
	resp, err := client.DescribeHostsWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("listWafClbHosts", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil {
		return nil, nil
	}

	for _, item := range resp.Response.HostList {
		if item == nil {
			continue
		}
		d.StreamListItem(ctx, toWafClbHostRow(item))
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

// Row Type

type wafClbHostRow struct {
	Domain          string
	DomainId        string
	MainDomain      string
	Edition         string
	Region          string
	Status          uint64
	State           uint64
	Mode            uint64
	Engine          uint64
	EngineType      int64
	FlowMode        uint64
	ClsStatus       uint64
	IsCdn           uint64
	Level           uint64
	AlbType         string
	CloudType       string
	Note            string
	IpHeaders       []string
	CdcClusters     []string
	LoadBalancerSet []*waf.LoadBalancer
	Title           string
	Akas            []string
}

func toWafClbHostRow(h *waf.HostRecord) wafClbHostRow {
	row := wafClbHostRow{
		Domain:          utils.PtrString(h.Domain),
		DomainId:        utils.PtrString(h.DomainId),
		MainDomain:      utils.PtrString(h.MainDomain),
		Edition:         utils.PtrString(h.Edition),
		Region:          utils.PtrString(h.Region),
		Status:          utils.PtrUint64(h.Status),
		State:           utils.PtrUint64(h.State),
		Mode:            utils.PtrUint64(h.Mode),
		Engine:          utils.PtrUint64(h.Engine),
		EngineType:      utils.PtrInt64(h.EngineType),
		FlowMode:        utils.PtrUint64(h.FlowMode),
		ClsStatus:       utils.PtrUint64(h.ClsStatus),
		IsCdn:           utils.PtrUint64(h.IsCdn),
		Level:           utils.PtrUint64(h.Level),
		AlbType:         utils.PtrString(h.AlbType),
		CloudType:       utils.PtrString(h.CloudType),
		Note:            utils.PtrString(h.Note),
		IpHeaders:       utils.PtrStringSlice(h.IpHeaders),
		CdcClusters:     utils.PtrStringSlice(h.CdcClusters),
		LoadBalancerSet: h.LoadBalancerSet,
	}

	if row.Domain != "" {
		row.Title = row.Domain
	} else {
		row.Title = row.DomainId
	}

	if row.DomainId != "" {
		row.Akas = []string{"tencentcloud:waf:clb-host:" + row.DomainId}
	}
	return row
}
