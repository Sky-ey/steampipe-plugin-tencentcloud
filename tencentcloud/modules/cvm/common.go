package cvm

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"

	as "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/as/v20180419"
	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/cvm/v20170312"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func TagsToMap(tags []*cvm.Tag) map[string]string {
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

func TagsAsToMap(tags []*as.Tag) map[string]string {
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

func buildCvmRegionList(ctx context.Context, d *plugin.QueryData) []map[string]any {
	return utils.BuildRegionList(ctx, d, "cvm")
}

func listCvmInstanceIds(ctx context.Context, d *plugin.QueryData, client *cvm.Client, requirePublicIP bool) ([]string, error) {
	pageSize := int64(100)
	var offset int64 = 0
	var ids []string
	instanceID := d.EqualsQualString("instance_id")

	for {
		d.WaitForListRateLimit(ctx)
		req := cvm.NewDescribeInstancesRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		if instanceID != "" {
			req.InstanceIds = common.StringPtrs([]string{instanceID})
		}

		utils.LogRequest(ctx, "listCvmInstanceIds", req)
		resp, err := client.DescribeInstancesWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listCvmInstanceIds", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		instances := resp.Response.InstanceSet
		for _, inst := range instances {
			if inst == nil || inst.InstanceId == nil {
				continue
			}
			if requirePublicIP && !cvmInstanceHasPublicIPAddress(inst) {
				plugin.Logger(ctx).Debug("listCvmInstanceIds", "skip_instance_without_public_ip", utils.PtrString(inst.InstanceId))
				continue
			}
			ids = append(ids, *inst.InstanceId)
		}

		total := utils.PtrInt64(resp.Response.TotalCount)
		if utils.PageDone(offset, int64(len(instances)), total) {
			break
		}
		offset += int64(len(instances))
	}

	return ids, nil
}

func cvmInstanceHasPublicIPAddress(inst *cvm.Instance) bool {
	if inst == nil {
		return false
	}
	for _, address := range inst.PublicIpAddresses {
		if utils.PtrString(address) != "" {
			return true
		}
	}
	return false
}
