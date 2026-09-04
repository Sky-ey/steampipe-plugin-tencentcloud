package cbs

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"

	cbs "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/cbs/v20170312"
	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func buildCbsRegionList(ctx context.Context, d *plugin.QueryData) []map[string]any {
	return utils.BuildRegionList(ctx, d, "cbs")
}

func newCbsFilter(name, value string) *cbs.Filter {
	return &cbs.Filter{
		Name:   common.StringPtr(name),
		Values: common.StringPtrs([]string{value}),
	}
}

func TagsToMap(tags []*cbs.Tag) map[string]string {
	if len(tags) == 0 {
		return nil
	}
	result := make(map[string]string, len(tags))
	for _, tag := range tags {
		if tag == nil {
			continue
		}
		key := utils.PtrString(tag.Key)
		if key != "" {
			result[key] = utils.PtrString(tag.Value)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}
