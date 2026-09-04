package waf

import (
	"context"
	"strconv"
	"time"

	waf "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/waf/v20180125"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
)

// Table Definition

func TableTencentcloudWafCustomRule() *plugin.Table {
	return &plugin.Table{
		Name:        "tencentcloud_waf_custom_rule",
		Description: "Tencent Cloud WAF custom access control rules per protected domain.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "domain", Require: plugin.Optional},
				{Name: "rule_id", Require: plugin.Optional},
				{Name: "rule_name", Require: plugin.Optional},
			},
			Hydrate: listWafCustomRules,
			Tags:    map[string]string{"service": "waf", "action": "DescribeCustomRuleList"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "rule_id", Type: proto.ColumnType_STRING, Description: "The custom rule ID.", Transform: transform.FromField("RuleId")},
			{Name: "rule_name", Type: proto.ColumnType_STRING, Description: "The custom rule name.", Transform: transform.FromField("RuleName")},
			{Name: "domain", Type: proto.ColumnType_STRING, Description: "The domain name the rule belongs to. Required by the API; populated by fan-out when not specified in the query.", Transform: transform.FromField("Domain")},
			{Name: "action_type", Type: proto.ColumnType_STRING, Description: "The action type taken when the rule matches.", Transform: transform.FromField("ActionType")},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The rule status.", Transform: transform.FromField("Status")},
			{Name: "valid_status", Type: proto.ColumnType_INT, Description: "The effective status of the rule.", Transform: transform.FromField("ValidStatus")},
			{Name: "sort_id", Type: proto.ColumnType_STRING, Description: "The priority (sort ID) of the rule.", Transform: transform.FromField("SortId")},
			{Name: "bypass", Type: proto.ColumnType_STRING, Description: "The skipped policy (e.g. cc, owasp, ai, antileakage).", Transform: transform.FromField("Bypass")},
			{Name: "redirect", Type: proto.ColumnType_STRING, Description: "The redirect URL when the rule redirects.", Transform: transform.FromField("Redirect")},
			{Name: "page_id", Type: proto.ColumnType_STRING, Description: "The ID of the blocked page used by the rule.", Transform: transform.FromField("PageId")},
			{Name: "event_id", Type: proto.ColumnType_STRING, Description: "The event ID associated with the rule.", Transform: transform.FromField("EventId")},
			{Name: "source", Type: proto.ColumnType_STRING, Description: "The rule source (e.g. custom).", Transform: transform.FromField("Source")},
			{Name: "label", Type: proto.ColumnType_STRING, Description: "The custom tag indicating whether the rule is built-in or user-defined.", Transform: transform.FromField("Label")},
			{Name: "strategies", Type: proto.ColumnType_JSON, Description: "The match strategy list of the rule (field, compare function, content, etc.).", Transform: transform.FromField("Strategies")},
			{Name: "job_type", Type: proto.ColumnType_STRING, Description: "The scheduled task type of the rule.", Transform: transform.FromField("JobType")},
			{Name: "cron_type", Type: proto.ColumnType_STRING, Description: "The periodic task granularity (e.g. week, month).", Transform: transform.FromField("CronType")},
			{Name: "job_date_time", Type: proto.ColumnType_JSON, Description: "The scheduled task configuration (timed and cron entries) of the rule.", Transform: transform.FromField("JobDateTime")},
			{Name: "create_time", Type: proto.ColumnType_TIMESTAMP, Description: "The rule creation time.", Transform: transform.FromField("CreateTime").NullIfZero()},
			{Name: "modify_time", Type: proto.ColumnType_TIMESTAMP, Description: "The rule last modification time.", Transform: transform.FromField("ModifyTime").NullIfZero()},
			{Name: "expire_time", Type: proto.ColumnType_TIMESTAMP, Description: "The rule expiration time.", Transform: transform.FromField("ExpireTime").NullIfZero()},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listWafCustomRules(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listWafCustomRules", "quals", d.EqualsQuals)

	client := &waf.Client{}
	if err := utils.InitClientInRegion(ctx, d, client, wafRequestRegion(d)); err != nil {
		plugin.Logger(ctx).Error("listWafCustomRules", "client_init_error", err)
		return nil, err
	}

	domains, err := resolveWafDomainIDs(ctx, d, client)
	if err != nil {
		plugin.Logger(ctx).Error("listWafCustomRules", "resolve_domain_ids_error", err)
		return nil, err
	}

	filters := buildWafCustomRuleFilters(d.EqualsQuals)

	pageSize := uint64(utils.PageSizeFromLimit(d, 100, 1))

	for _, domain := range domains {
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
		if err := streamWafCustomRulePage(ctx, d, client, domain, filters, pageSize); err != nil {
			plugin.Logger(ctx).Error("listWafCustomRules", "request_error", err)
			return nil, err
		}
	}
	return nil, nil
}

func buildWafCustomRuleFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*waf.FiltersItemNew {
	exact := utils.BuildFilters(equalsQuals, newWafExactFilter,
		utils.FilterMapping{QualName: "rule_id", FilterName: "RuleID"},
	)
	fuzzy := utils.BuildFilters(equalsQuals, newWafFuzzyFilter,
		utils.FilterMapping{QualName: "rule_name", FilterName: "RuleName"},
	)
	return append(exact, fuzzy...)
}

func streamWafCustomRulePage(ctx context.Context, d *plugin.QueryData, client *waf.Client, domain string, filters []*waf.FiltersItemNew, pageSize uint64) error {
	var offset uint64
	for {
		d.WaitForListRateLimit(ctx)

		req := waf.NewDescribeCustomRuleListRequest()
		req.Domain = &domain
		req.Offset = &offset
		req.Limit = &pageSize
		if len(filters) > 0 {
			req.Filters = filters
		}

		utils.LogRequest(ctx, "listWafCustomRules", req)
		resp, err := client.DescribeCustomRuleListWithContext(ctx, req)
		if err != nil {
			return err
		}
		if resp == nil || resp.Response == nil {
			return nil
		}

		items := resp.Response.RuleList
		for _, item := range items {
			if item == nil {
				continue
			}
			d.StreamListItem(ctx, toWafCustomRuleRow(item, domain))
			if d.RowsRemaining(ctx) == 0 {
				return nil
			}
		}

		totalStr := utils.PtrString(resp.Response.TotalCount)
		total, parseErr := strconv.ParseUint(totalStr, 10, 64)
		if parseErr != nil {
			return nil
		}
		if utils.PageDone(int64(offset), int64(len(items)), total) {
			return nil
		}
		offset += uint64(len(items))
	}
}

// Row Type

type wafCustomRuleRow struct {
	RuleId      string
	RuleName    string
	Domain      string
	ActionType  string
	Status      string
	ValidStatus int64
	SortId      string
	Bypass      string
	Redirect    string
	PageId      string
	EventId     string
	Source      string
	Label       string
	Strategies  []*waf.Strategy
	JobType     string
	CronType    string
	JobDateTime *waf.JobDateTime
	CreateTime  *time.Time
	ModifyTime  *time.Time
	ExpireTime  *time.Time
	Title       string
	Akas        []string
}

func toWafCustomRuleRow(r *waf.DescribeCustomRulesRspRuleListItem, domain string) wafCustomRuleRow {
	row := wafCustomRuleRow{
		RuleId:      utils.PtrString(r.RuleId),
		RuleName:    utils.PtrString(r.Name),
		Domain:      domain,
		ActionType:  utils.PtrString(r.ActionType),
		Status:      utils.PtrString(r.Status),
		ValidStatus: utils.PtrInt64(r.ValidStatus),
		SortId:      utils.PtrString(r.SortId),
		Bypass:      utils.PtrString(r.Bypass),
		Redirect:    utils.PtrString(r.Redirect),
		PageId:      utils.PtrString(r.PageId),
		EventId:     utils.PtrString(r.EventId),
		Source:      utils.PtrString(r.Source),
		Label:       utils.PtrString(r.Label),
		Strategies:  r.Strategies,
		JobType:     utils.PtrString(r.JobType),
		CronType:    utils.PtrString(r.CronType),
		JobDateTime: r.JobDateTime,
		CreateTime:  utils.ParseTimestamp(r.CreateTime),
		ModifyTime:  utils.ParseTimestamp(r.ModifyTime),
		ExpireTime:  utils.ParseTimestamp(r.ExpireTime),
	}

	row.Title = row.RuleName
	if row.Title == "" {
		row.Title = row.RuleId
	}
	if row.RuleId != "" {
		row.Akas = []string{"tencentcloud:waf:custom-rule:" + domain + ":" + row.RuleId}
	}
	return row
}
