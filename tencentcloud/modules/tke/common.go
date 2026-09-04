package tke

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"

	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	tke "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/tke/v20180525"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func buildTkeRegionList(ctx context.Context, d *plugin.QueryData) []map[string]any {
	return utils.BuildRegionList(ctx, d, "tke")
}

func newTkeFilter(name, value string) *tke.Filter {
	return &tke.Filter{
		Name:   common.StringPtr(name),
		Values: common.StringPtrs([]string{value}),
	}
}

func TagsToMap(tags []*tke.Tag) map[string]string {
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

func TagSpecificationToMap(specs []*tke.TagSpecification) map[string]string {
	if len(specs) == 0 {
		return nil
	}
	out := make(map[string]string)
	for _, spec := range specs {
		if spec == nil {
			continue
		}
		for _, t := range spec.Tags {
			if t == nil {
				continue
			}
			k := utils.PtrString(t.Key)
			if k == "" {
				continue
			}
			out[k] = utils.PtrString(t.Value)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func listTkeClusterIDs(ctx context.Context, d *plugin.QueryData, client *tke.Client) ([]string, error) {
	pageSize := int64(100)
	var offset int64
	var ids []string

	for {
		d.WaitForListRateLimit(ctx)
		req := tke.NewDescribeClustersRequest()
		req.Offset = &offset
		req.Limit = &pageSize

		utils.LogRequest(ctx, "listTkeClusterIDs", req)
		resp, err := client.DescribeClustersWithContext(ctx, req)
		if err != nil {
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		clusters := resp.Response.Clusters
		for _, c := range clusters {
			if c == nil {
				continue
			}
			if id := utils.PtrString(c.ClusterId); id != "" {
				ids = append(ids, id)
			}
		}

		if resp.Response == nil {
			break
		}
		total := utils.PtrInt64(resp.Response.TotalCount)
		if utils.PageDone(offset, int64(len(clusters)), total) {
			break
		}
		offset += int64(len(clusters))
	}
	return ids, nil
}

func resolveTkeClusterIDs(ctx context.Context, d *plugin.QueryData, client *tke.Client) ([]string, error) {
	if id := d.EqualsQualString("cluster_id"); id != "" {
		return []string{id}, nil
	}
	return listTkeClusterIDs(ctx, d, client)
}
