package waf

import (
	"strings"

	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	waf "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/waf/v20180125"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"

	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
)

func wafRequestRegion(d *plugin.QueryData) string {
	for _, r := range utils.GetConfig(d.Connection).Regions {
		if !strings.ContainsAny(r, "*?[") {
			return r
		}
	}
	return utils.DefaultRegion
}

func newWafFilter(name, value string, exact bool) *waf.FiltersItemNew {
	return &waf.FiltersItemNew{
		Name:       common.StringPtr(name),
		Values:     common.StringPtrs([]string{value}),
		ExactMatch: common.BoolPtr(exact),
	}
}

func newWafExactFilter(name, value string) *waf.FiltersItemNew {
	return newWafFilter(name, value, true)
}

func newWafFuzzyFilter(name, value string) *waf.FiltersItemNew {
	return newWafFilter(name, value, false)
}
