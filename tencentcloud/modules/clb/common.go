package clb

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"

	clb "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/clb/v20180317"
	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func buildClbRegionList(ctx context.Context, d *plugin.QueryData) []map[string]any {
	return utils.BuildRegionList(ctx, d, "clb")
}

// listClbLoadBalancerIDs 拉取当前 region 下所有 CLB 实例 ID。
func listClbLoadBalancerIDs(ctx context.Context, d *plugin.QueryData, client *clb.Client) ([]string, error) {
	pageSize := int64(100)
	var offset int64
	var ids []string

	for {
		d.WaitForListRateLimit(ctx)
		req := clb.NewDescribeLoadBalancersRequest()
		req.Offset = &offset
		req.Limit = &pageSize

		utils.LogRequest(ctx, "listClbLoadBalancerIDs", req)
		resp, err := client.DescribeLoadBalancersWithContext(ctx, req)
		if err != nil {
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		loadBalancers := resp.Response.LoadBalancerSet
		for _, lb := range loadBalancers {
			if lb == nil {
				continue
			}
			if id := utils.PtrString(lb.LoadBalancerId); id != "" {
				ids = append(ids, id)
			}
		}

		if utils.PageDone(offset, int64(len(loadBalancers)), utils.PtrUint64(resp.Response.TotalCount)) {
			break
		}
		offset += int64(len(loadBalancers))
	}
	return ids, nil
}

// resolveClbLoadBalancerIDs 解析待查询的 CLB 实例 ID 列表
func resolveClbLoadBalancerIDs(ctx context.Context, d *plugin.QueryData, client *clb.Client) ([]string, error) {
	if id := d.EqualsQualString("load_balancer_id"); id != "" {
		return []string{id}, nil
	}
	return listClbLoadBalancerIDs(ctx, d, client)
}

// listClbTargetGroupIDs 拉取当前 region 下所有 CLB 目标组 ID。
func listClbTargetGroupIDs(ctx context.Context, d *plugin.QueryData, client *clb.Client) ([]string, error) {
	pageSize := uint64(100)
	var offset uint64
	var ids []string

	for {
		d.WaitForListRateLimit(ctx)
		req := clb.NewDescribeTargetGroupListRequest()
		req.Offset = &offset
		req.Limit = &pageSize

		utils.LogRequest(ctx, "listClbTargetGroupIDs", req)
		resp, err := client.DescribeTargetGroupListWithContext(ctx, req)
		if err != nil {
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		targetGroups := resp.Response.TargetGroupSet
		for _, tg := range targetGroups {
			if tg == nil {
				continue
			}
			if id := utils.PtrString(tg.TargetGroupId); id != "" {
				ids = append(ids, id)
			}
		}

		if utils.PageDone(int64(offset), int64(len(targetGroups)), utils.PtrUint64(resp.Response.TotalCount)) {
			break
		}
		offset += uint64(len(targetGroups))
	}
	return ids, nil
}

// resolveClbTargetGroupIDs 解析待查询的 CLB 目标组 ID 列表。
// 优先使用 target_group_id qual;未指定时枚举当前 region 下所有目标组。
func resolveClbTargetGroupIDs(ctx context.Context, d *plugin.QueryData, client *clb.Client) ([]string, error) {
	if id := d.EqualsQualString("target_group_id"); id != "" {
		return []string{id}, nil
	}
	return listClbTargetGroupIDs(ctx, d, client)
}

func newClbFilter(name, value string) *clb.Filter {
	return &clb.Filter{
		Name:   common.StringPtr(name),
		Values: common.StringPtrs([]string{value}),
	}
}

func TagsToMap(tags []*clb.TagInfo) map[string]string {
	if len(tags) == 0 {
		return nil
	}
	result := make(map[string]string, len(tags))
	for _, tag := range tags {
		if tag == nil {
			continue
		}
		key := utils.PtrString(tag.TagKey)
		if key != "" {
			result[key] = utils.PtrString(tag.TagValue)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}
