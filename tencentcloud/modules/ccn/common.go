package ccn

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"strings"

	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/vpc/v20170312"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// CCN is an account-scoped global service.

func ccnRequestRegion(d *plugin.QueryData) string {
	for _, r := range utils.GetConfig(d.Connection).Regions {
		if !strings.ContainsAny(r, "*?[") {
			return r
		}
	}
	return utils.DefaultRegion
}

func newCcnFilter(name, value string) *vpc.Filter {
	return &vpc.Filter{
		Name:   common.StringPtr(name),
		Values: common.StringPtrs([]string{value}),
	}
}

func listCcnIDs(ctx context.Context, d *plugin.QueryData, client *vpc.Client) ([]string, error) {
	pageSize := uint64(100)
	var offset uint64
	var ids []string

	for {
		d.WaitForListRateLimit(ctx)
		req := vpc.NewDescribeCcnsRequest()
		req.Offset = &offset
		req.Limit = &pageSize

		utils.LogRequest(ctx, "listCcnIDs", req)
		resp, err := client.DescribeCcnsWithContext(ctx, req)
		if err != nil {
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		ccns := resp.Response.CcnSet
		for _, ccn := range ccns {
			if ccn == nil {
				continue
			}
			if id := utils.PtrString(ccn.CcnId); id != "" {
				ids = append(ids, id)
			}
		}

		if resp.Response == nil {
			break
		}
		total := utils.PtrUint64(resp.Response.TotalCount)
		if utils.PageDone(int64(offset), int64(len(ccns)), total) {
			break
		}
		offset += uint64(len(ccns))
	}
	return ids, nil
}

func resolveCcnIDs(ctx context.Context, d *plugin.QueryData, client *vpc.Client) ([]string, error) {
	if id := d.EqualsQualString("ccn_id"); id != "" {
		return []string{id}, nil
	}
	return listCcnIDs(ctx, d, client)
}
