package cam

import (
	"context"

	cam "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/cam/v20190116"
	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"

	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
)

// CAM is an account-scoped global service.
const camDefaultPageSize uint64 = 200

func RoleTagsToMap(tags []*cam.RoleTags) map[string]string {
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

func listCamUsers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listCamUsers", "quals", d.EqualsQuals)

	client := &cam.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listCamUsers", "client_init_error", err)
		return nil, err
	}

	d.WaitForListRateLimit(ctx)
	req := cam.NewListUsersRequest()
	utils.LogRequest(ctx, "listCamUsers", req)
	resp, err := client.ListUsersWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("listCamUsers", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil {
		return nil, nil
	}

	for _, u := range resp.Response.Data {
		if u == nil {
			continue
		}
		d.StreamListItem(ctx, u)
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}

func listAllCamPolicies(ctx context.Context, d *plugin.QueryData, client *cam.Client) ([]*cam.StrategyInfo, error) {
	scope := d.EqualsQualString("scope")
	if scope == "" {
		scope = "All"
	}

	pageSize := camDefaultPageSize
	var page uint64 = 1
	var out []*cam.StrategyInfo

	for {
		d.WaitForListRateLimit(ctx)
		req := cam.NewListPoliciesRequest()
		req.Rp = common.Uint64Ptr(pageSize)
		req.Page = common.Uint64Ptr(page)
		req.Scope = common.StringPtr(scope)

		utils.LogRequest(ctx, "listAllCamPolicies", req)
		resp, err := client.ListPoliciesWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listAllCamPolicies", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			return out, nil
		}

		items := resp.Response.List
		for _, p := range items {
			if p == nil {
				continue
			}
			out = append(out, p)
		}

		offset := int64((page - 1) * pageSize)
		fetched := int64(len(items))
		total := int64(0)
		if resp.Response.TotalNum != nil {
			total = int64(*resp.Response.TotalNum)
		}
		if utils.PageDone(offset, fetched, total) {
			break
		}
		page++
	}
	return out, nil
}
