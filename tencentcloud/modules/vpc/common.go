package vpc

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"strconv"

	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/vpc/v20170312"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func buildVpcRegionList(ctx context.Context, d *plugin.QueryData) []map[string]any {
	return utils.BuildRegionList(ctx, d, "vpc")
}

func newVpcFilter(name, value string) *vpc.Filter {
	return &vpc.Filter{
		Name:   common.StringPtr(name),
		Values: common.StringPtrs([]string{value}),
	}
}

func listVpcSecurityGroupIDs(ctx context.Context, d *plugin.QueryData, client *vpc.Client) ([]string, error) {
	pageSize := "100"
	offset := "0"
	var ids []string

	for {
		d.WaitForListRateLimit(ctx)
		req := vpc.NewDescribeSecurityGroupsRequest()
		req.Offset = &offset
		req.Limit = &pageSize

		utils.LogRequest(ctx, "listVpcSecurityGroupIDs", req)
		resp, err := client.DescribeSecurityGroupsWithContext(ctx, req)
		if err != nil {
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		groups := resp.Response.SecurityGroupSet
		for _, sg := range groups {
			if sg == nil {
				continue
			}
			if id := utils.PtrString(sg.SecurityGroupId); id != "" {
				ids = append(ids, id)
			}
		}

		if resp.Response == nil {
			break
		}
		total := utils.PtrUint64(resp.Response.TotalCount)
		currentOffset, _ := strconv.ParseUint(offset, 10, 64)
		if utils.PageDone(int64(currentOffset), int64(len(groups)), total) {
			break
		}
		offset = strconv.FormatUint(currentOffset+uint64(len(groups)), 10)
	}
	return ids, nil
}

func resolveVpcSecurityGroupIDs(ctx context.Context, d *plugin.QueryData, client *vpc.Client) ([]string, error) {
	if id := d.EqualsQualString("security_group_id"); id != "" {
		return []string{id}, nil
	}
	return listVpcSecurityGroupIDs(ctx, d, client)
}

func TagsToMap(tags []*vpc.Tag) map[string]string {
	if len(tags) == 0 {
		return nil
	}
	out := make(map[string]string, len(tags))
	for _, t := range tags {
		if t == nil {
			continue
		}
		k := utils.PtrString(t.Key)
		if k == "" {
			continue
		}
		out[k] = utils.PtrString(t.Value)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
