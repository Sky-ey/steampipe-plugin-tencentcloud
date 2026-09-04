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

func TableTencentcloudWafIpAccessControl() *plugin.Table {
	return &plugin.Table{
		Name:        "tencentcloud_waf_ip_access_control",
		Description: "Tencent Cloud WAF IP allowlist/blocklist entries grouped by domain and action type.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "domain", Require: plugin.Optional},
				{Name: "action_type", Require: plugin.Optional},
			},
			Hydrate: listWafIpAccessControls,
			Tags:    map[string]string{"service": "waf", "action": "DescribeIpAccessControl"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "rule_id", Type: proto.ColumnType_INT, Description: "The IP allowlist/blocklist rule ID.", Transform: transform.FromField("RuleId")},
			{Name: "domain", Type: proto.ColumnType_STRING, Description: "The domain name the rule belongs to; 'global' for global rules. Required by the API; populated by fan-out when not specified in the query.", Transform: transform.FromField("Domain")},
			{Name: "action_type", Type: proto.ColumnType_INT, Description: "The action type (40: allowlist, 42: blocklist).", Transform: transform.FromField("ActionType")},
			{Name: "ip", Type: proto.ColumnType_STRING, Description: "The IP or IP range the rule applies to.", Transform: transform.FromField("Ip")},
			{Name: "ip_list", Type: proto.ColumnType_JSON, Description: "The IP list associated with the rule.", Transform: transform.FromField("IpList")},
			{Name: "note", Type: proto.ColumnType_STRING, Description: "The remark of the rule.", Transform: transform.FromField("Note")},
			{Name: "source", Type: proto.ColumnType_STRING, Description: "The source of the rule (e.g. custom).", Transform: transform.FromField("Source")},
			{Name: "valid_status", Type: proto.ColumnType_INT, Description: "The effective status of the rule.", Transform: transform.FromField("ValidStatus")},
			{Name: "valid_ts", Type: proto.ColumnType_TIMESTAMP, Description: "The expiration timestamp of the rule (seconds since epoch).", Transform: transform.FromField("ValidTs").NullIfZero()},
			{Name: "ts_version", Type: proto.ColumnType_TIMESTAMP, Description: "The update timestamp of the rule (seconds since epoch).", Transform: transform.FromField("TsVersion").NullIfZero()},
			{Name: "create_time", Type: proto.ColumnType_TIMESTAMP, Description: "The rule creation timestamp (seconds since epoch).", Transform: transform.FromField("CreateTime").NullIfZero()},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The internal MongoDB auto-increment ID of the rule.", Transform: transform.FromField("Id")},
			{Name: "job_type", Type: proto.ColumnType_STRING, Description: "The scheduled task type of the rule.", Transform: transform.FromField("JobType")},
			{Name: "cron_type", Type: proto.ColumnType_STRING, Description: "The periodic task granularity (e.g. week, month).", Transform: transform.FromField("CronType")},
			{Name: "job_date_time", Type: proto.ColumnType_JSON, Description: "The scheduled task configuration (timed and cron entries) of the rule.", Transform: transform.FromField("JobDateTime")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listWafIpAccessControls(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listWafIpAccessControls", "quals", d.EqualsQuals)

	client := &waf.Client{}
	if err := utils.InitClientInRegion(ctx, d, client, wafRequestRegion(d)); err != nil {
		plugin.Logger(ctx).Error("listWafIpAccessControls", "client_init_error", err)
		return nil, err
	}

	domains, err := resolveWafDomainIDs(ctx, d, client)
	if err != nil {
		plugin.Logger(ctx).Error("listWafIpAccessControls", "resolve_domain_ids_error", err)
		return nil, err
	}

	actionTypes := []uint64{40, 42}
	if v := d.EqualsQuals["action_type"]; v != nil {
		if a := uint64(v.GetInt64Value()); a == 40 || a == 42 {
			actionTypes = []uint64{a}
		}
	}

	count := uint64(1)
	pageSize := uint64(utils.PageSizeFromLimit(d, 100, 1))

	for _, domain := range domains {
		for _, actionType := range actionTypes {
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
			if err := streamWafIpAccessControlPage(ctx, d, client, domain, actionType, count, pageSize); err != nil {
				plugin.Logger(ctx).Error("listWafIpAccessControls", "request_error", err)
				return nil, err
			}
		}
	}
	return nil, nil
}

func streamWafIpAccessControlPage(ctx context.Context, d *plugin.QueryData, client *waf.Client, domain string, actionType, count, pageSize uint64) error {
	var offset uint64
	for {
		d.WaitForListRateLimit(ctx)

		req := waf.NewDescribeIpAccessControlRequest()
		req.Domain = &domain
		req.Count = &count
		req.ActionType = &actionType
		req.OffSet = &offset
		req.Limit = &pageSize

		utils.LogRequest(ctx, "listWafIpAccessControls", req)
		resp, err := client.DescribeIpAccessControlWithContext(ctx, req)
		if err != nil {
			return err
		}
		if resp == nil || resp.Response == nil || resp.Response.Data == nil {
			return nil
		}

		items := resp.Response.Data.Res
		for _, item := range items {
			if item == nil {
				continue
			}
			d.StreamListItem(ctx, toWafIpAccessControlRow(item, domain))
			if d.RowsRemaining(ctx) == 0 {
				return nil
			}
		}

		total := utils.PtrUint64(resp.Response.Data.TotalCount)
		if utils.PageDone(int64(offset), int64(len(items)), total) {
			return nil
		}
		offset += uint64(len(items))
	}
}

// Row Type

type wafIpAccessControlRow struct {
	RuleId      uint64
	Domain      string
	ActionType  uint64
	Ip          string
	IpList      []*string
	Note        string
	Source      string
	ValidStatus int64
	ValidTs     uint64
	TsVersion   uint64
	CreateTime  uint64
	Id          string
	JobType     string
	CronType    string
	JobDateTime *waf.JobDateTime
	Title       string
	Akas        []string
}

func toWafIpAccessControlRow(i *waf.IpAccessControlItem, domain string) wafIpAccessControlRow {
	row := wafIpAccessControlRow{
		RuleId:      utils.PtrUint64(i.RuleId),
		Domain:      domain,
		ActionType:  utils.PtrUint64(i.ActionType),
		Ip:          utils.PtrString(i.Ip),
		IpList:      i.IpList,
		Note:        utils.PtrString(i.Note),
		Source:      utils.PtrString(i.Source),
		ValidStatus: utils.PtrInt64(i.ValidStatus),
		ValidTs:     utils.PtrUint64(i.ValidTs),
		TsVersion:   utils.PtrUint64(i.TsVersion),
		CreateTime:  utils.PtrUint64(i.CreateTime),
		Id:          utils.PtrString(i.Id),
		JobType:     utils.PtrString(i.JobType),
		CronType:    utils.PtrString(i.CronType),
		JobDateTime: i.JobDateTime,
	}

	row.Title = row.Ip
	if row.Title == "" {
		row.Title = row.Id
	}
	if row.Id != "" {
		row.Akas = []string{"tencentcloud:waf:ip-access-control:" + domain + ":" + row.Id}
	}
	return row
}

func listWafDomainIDs(ctx context.Context, d *plugin.QueryData, client *waf.Client) ([]string, error) {
	pageSize := uint64(100)
	var offset uint64
	var ids []string

	for {
		d.WaitForListRateLimit(ctx)
		req := waf.NewDescribeDomainsRequest()
		req.Offset = &offset
		req.Limit = &pageSize

		utils.LogRequest(ctx, "listWafDomainIDs", req)
		resp, err := client.DescribeDomainsWithContext(ctx, req)
		if err != nil {
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		domains := resp.Response.Domains
		for _, dm := range domains {
			if dm == nil {
				continue
			}
			if id := utils.PtrString(dm.Domain); id != "" {
				ids = append(ids, id)
			}
		}

		total := utils.PtrUint64(resp.Response.Total)
		if utils.PageDone(int64(offset), int64(len(domains)), total) {
			break
		}
		offset += uint64(len(domains))
	}

	ids = append(ids, "global")
	return ids, nil
}

func resolveWafDomainIDs(ctx context.Context, d *plugin.QueryData, client *waf.Client) ([]string, error) {
	if id := d.EqualsQualString("domain"); id != "" {
		return []string{id}, nil
	}
	return listWafDomainIDs(ctx, d, client)
}
