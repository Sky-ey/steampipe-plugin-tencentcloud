package clb

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"

	clb "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/clb/v20180317"
	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudClbRewrite() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_clb_rewrite",
		Description:       "Tencent Cloud Load Balancer (CLB) URL rewrite redirection rules.",
		GetMatrixItemFunc: buildClbRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "load_balancer_id", Require: plugin.Optional},
				{Name: "listener_id", Require: plugin.Optional},
				{Name: "location_id", Require: plugin.Optional},
			},
			Hydrate: listClbRewrites,
			Tags:    map[string]string{"service": "clb", "action": "DescribeRewrite"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "load_balancer_id", Type: proto.ColumnType_STRING, Description: "The ID of the load balancer to which the rewrite rule belongs.", Transform: transform.FromField("LoadBalancerId")},
			{Name: "listener_id", Type: proto.ColumnType_STRING, Description: "The ID of the source listener of the rewrite rule.", Transform: transform.FromField("ListenerId")},
			{Name: "location_id", Type: proto.ColumnType_STRING, Description: "The ID of the source forwarding rule of the rewrite rule.", Transform: transform.FromField("LocationId")},
			{Name: "domain", Type: proto.ColumnType_STRING, Description: "The domain name of the source forwarding rule."},
			{Name: "url", Type: proto.ColumnType_STRING, Description: "The URL of the source forwarding rule.", Transform: transform.FromField("Url")},
			{Name: "target_listener_id", Type: proto.ColumnType_STRING, Description: "The listener ID of the redirection target.", Transform: transform.FromField("TargetListenerId")},
			{Name: "target_location_id", Type: proto.ColumnType_STRING, Description: "The forwarding rule ID of the redirection target.", Transform: transform.FromField("TargetLocationId")},
			{Name: "rewrite_code", Type: proto.ColumnType_INT, Description: "The redirection status code of the rewrite rule.", Transform: transform.FromField("RewriteCode")},
			{Name: "take_url", Type: proto.ColumnType_BOOL, Description: "Whether the matched URL is carried in the redirection.", Transform: transform.FromField("TakeUrl")},
			{Name: "rewrite_type", Type: proto.ColumnType_STRING, Description: "The redirection type (Manual, Auto).", Transform: transform.FromField("RewriteType")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the load balancer resides."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource."},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource."},
		}),
	}
}

// List Function

func listClbRewrites(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listClbRewrites", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &clb.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listClbRewrites", "client_init_error", err)
		return nil, err
	}
	region := d.EqualsQualString(utils.MatrixKeyRegion)

	loadBalancerIDs, err := resolveClbLoadBalancerIDs(ctx, d, client)
	if err != nil {
		plugin.Logger(ctx).Error("listClbRewrites", "request_error", err)
		return nil, err
	}

	for _, loadBalancerID := range loadBalancerIDs {
		d.WaitForListRateLimit(ctx)

		req := clb.NewDescribeRewriteRequest()
		req.LoadBalancerId = &loadBalancerID
		if id := d.EqualsQualString("listener_id"); id != "" {
			req.SourceListenerIds = common.StringPtrs([]string{id})
		}
		if id := d.EqualsQualString("location_id"); id != "" {
			req.SourceLocationIds = common.StringPtrs([]string{id})
		}

		utils.LogRequest(ctx, "listClbRewrites", req)
		resp, err := client.DescribeRewriteWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listClbRewrites", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			continue
		}

		for _, rule := range resp.Response.RewriteSet {
			if rule == nil {
				continue
			}
			row := toClbRewriteRow(rule, loadBalancerID, region)
			d.StreamListItem(ctx, row)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, nil
}

// Row Type

type clbRewriteRow struct {
	LoadBalancerId   string
	ListenerId       string
	LocationId       string
	Domain           string
	Url              string
	TargetListenerId string
	TargetLocationId string
	RewriteCode      int64
	TakeUrl          bool
	RewriteType      string
	Region           string
	Title            string
	Akas             []string
}

func toClbRewriteRow(rule *clb.RuleOutput, loadBalancerID, region string) clbRewriteRow {
	row := clbRewriteRow{
		LoadBalancerId: loadBalancerID,
		ListenerId:     utils.PtrString(rule.ListenerId),
		LocationId:     utils.PtrString(rule.LocationId),
		Domain:         utils.PtrString(rule.Domain),
		Url:            utils.PtrString(rule.Url),
		Region:         region,
	}
	if rule.RewriteTarget != nil {
		row.TargetListenerId = utils.PtrString(rule.RewriteTarget.TargetListenerId)
		row.TargetLocationId = utils.PtrString(rule.RewriteTarget.TargetLocationId)
		row.RewriteCode = utils.PtrInt64(rule.RewriteTarget.RewriteCode)
		row.TakeUrl = utils.PtrBool(rule.RewriteTarget.TakeUrl)
		row.RewriteType = utils.PtrString(rule.RewriteTarget.RewriteType)
	}
	row.Title = row.Domain
	if row.Title != "" {
		row.Title += row.Url
	}
	if row.Title == "" {
		row.Title = row.LocationId
	}
	if row.LocationId != "" {
		row.Akas = []string{"tencentcloud:clb:rewrite:" + row.LoadBalancerId + "/" + row.ListenerId + "/" + row.LocationId}
	}
	return row
}
