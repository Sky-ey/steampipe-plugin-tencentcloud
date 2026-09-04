package vpn

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"

	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/vpc/v20170312"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func buildVpnRegionList(ctx context.Context, d *plugin.QueryData) []map[string]any {
	return utils.BuildRegionList(ctx, d, "vpc")
}

func newVpnFilter(name, value string) *vpc.Filter {
	return &vpc.Filter{
		Name:   common.StringPtr(name),
		Values: common.StringPtrs([]string{value}),
	}
}

func newVpnFilterObject(name, value string) *vpc.FilterObject {
	return &vpc.FilterObject{
		Name:   common.StringPtr(name),
		Values: common.StringPtrs([]string{value}),
	}
}

func listVpnGatewayIDs(ctx context.Context, d *plugin.QueryData, client *vpc.Client) ([]string, error) {
	pageSize := uint64(100)
	var offset uint64
	var ids []string

	for {
		d.WaitForListRateLimit(ctx)
		req := vpc.NewDescribeVpnGatewaysRequest()
		req.Offset = &offset
		req.Limit = &pageSize

		utils.LogRequest(ctx, "listVpnGatewayIDs", req)
		resp, err := client.DescribeVpnGatewaysWithContext(ctx, req)
		if err != nil {
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		gws := resp.Response.VpnGatewaySet
		for _, gw := range gws {
			if gw == nil {
				continue
			}
			if id := utils.PtrString(gw.VpnGatewayId); id != "" {
				ids = append(ids, id)
			}
		}

		if resp.Response == nil {
			break
		}
		total := utils.PtrUint64(resp.Response.TotalCount)
		if utils.PageDone(int64(offset), int64(len(gws)), total) {
			break
		}
		offset += uint64(len(gws))
	}
	return ids, nil
}

func resolveVpnGatewayIDs(ctx context.Context, d *plugin.QueryData, client *vpc.Client) ([]string, error) {
	if id := d.EqualsQualString("vpn_gateway_id"); id != "" {
		return []string{id}, nil
	}
	return listVpnGatewayIDs(ctx, d, client)
}
