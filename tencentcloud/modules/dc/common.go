package dc

import (
	"context"

	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"

	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	dc "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/dc/v20180410"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func buildDcRegionList(ctx context.Context, d *plugin.QueryData) []map[string]any {
	return utils.BuildRegionList(ctx, d, "dc")
}

func newDcFilter(name, value string) *dc.Filter {
	return &dc.Filter{
		Name:   common.StringPtr(name),
		Values: common.StringPtrs([]string{value}),
	}
}

// TagsToMap converts a slice of dc.Tag to a map[string]string.
func TagsToMap(tags []*dc.Tag) map[string]string {
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
