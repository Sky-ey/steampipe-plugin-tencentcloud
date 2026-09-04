package cls

import (
	"context"

	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"

	cls "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/cls/v20201016"
	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func buildClsRegionList(ctx context.Context, d *plugin.QueryData) []map[string]any {
	return utils.BuildRegionList(ctx, d, "cls")
}

func newClsFilter(name, value string) *cls.Filter {
	return &cls.Filter{
		Key:    common.StringPtr(name),
		Values: common.StringPtrs([]string{value}),
	}
}

func TagsToMap(tags []*cls.Tag) map[string]string {
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
