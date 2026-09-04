package ssl

import (
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"

	ssl "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/ssl/v20191205"
)

func TagsToMap(tags []*ssl.Tags) map[string]string {
	if len(tags) == 0 {
		return nil
	}
	out := make(map[string]string, len(tags))
	for _, t := range tags {
		if t == nil {
			continue
		}
		k := utils.PtrString(t.TagKey)
		if k == "" {
			continue
		}
		out[k] = utils.PtrString(t.TagValue)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
