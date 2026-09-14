package cam

import (
	"context"
	stderrors "errors"
	"strconv"
	"time"

	cam "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/cam/v20190116"
	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	tcerrors "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/errors"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
)

// Table Definition

func TableTencentcloudCamPolicyAttachment() *plugin.Table {
	return &plugin.Table{
		Name:        "tencentcloud_cam_policy_attachment",
		Description: "Tencent Cloud CAM policy attachments, the mapping between policies and the principals (sub-accounts, user groups) they are attached to.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "policy_id", Require: plugin.Optional},
			},
			Hydrate: listCamPolicyAttachments,
			Tags:    map[string]string{"service": "cam", "action": "ListEntitiesForPolicy"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "policy_id", Type: proto.ColumnType_INT, Description: "The ID of the policy.", Transform: transform.FromField("PolicyId")},
			{Name: "policy_name", Type: proto.ColumnType_STRING, Description: "The name of the policy.", Transform: transform.FromField("PolicyName")},
			{Name: "entity_id", Type: proto.ColumnType_STRING, Description: "The ID of the attached entity (sub-account ID or user group ID).", Transform: transform.FromField("EntityId")},
			{Name: "entity_name", Type: proto.ColumnType_STRING, Description: "The name of the attached entity.", Transform: transform.FromField("EntityName")},
			{Name: "entity_uin", Type: proto.ColumnType_INT, Description: "The UIN of the attached entity.", Transform: transform.FromField("EntityUin")},
			{Name: "related_type", Type: proto.ColumnType_INT, Description: "The type of the attached entity (1 = sub-account, 2 = user group).", Transform: transform.FromField("RelatedType")},
			{Name: "attachment_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time the policy was attached to the entity.", Transform: transform.FromField("AttachmentTime").NullIfZero()},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listCamPolicyAttachments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listCamPolicyAttachments", "quals", d.EqualsQuals)

	client := &cam.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listCamPolicyAttachments", "client_init_error", err)
		return nil, err
	}

	policies, err := resolveCamAttachmentPolicies(ctx, d, client)
	if err != nil {
		plugin.Logger(ctx).Error("listCamPolicyAttachments", "resolve_policies_error", err)
		return nil, err
	}

	for _, policy := range policies {
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
		if policy == nil || policy.PolicyId == nil {
			continue
		}
		if camPolicyPositivelyUnattached(policy) {
			continue
		}
		if err := listCamPolicyAttachmentsForPolicy(ctx, d, client, policy); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func resolveCamAttachmentPolicies(ctx context.Context, d *plugin.QueryData, client *cam.Client) ([]*cam.StrategyInfo, error) {
	if qual := d.EqualsQuals["policy_id"]; qual != nil {
		policyID := qual.GetInt64Value()
		if policyID <= 0 {
			return nil, nil
		}
		plugin.Logger(ctx).Debug("resolveCamAttachmentPolicies", "policy_id", policyID)

		d.WaitForListRateLimit(ctx)
		req := cam.NewGetPolicyRequest()
		req.PolicyId = common.Uint64Ptr(uint64(policyID))

		utils.LogRequest(ctx, "resolveCamAttachmentPolicies", req)
		resp, err := client.GetPolicyWithContext(ctx, req)
		if err != nil {
			if isCamPolicyNotFoundError(err) {
				plugin.Logger(ctx).Warn("resolveCamAttachmentPolicies", "policy_not_found", policyID)
				return nil, nil
			}
			plugin.Logger(ctx).Error("resolveCamAttachmentPolicies", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			return nil, nil
		}
		return []*cam.StrategyInfo{{
			PolicyId:   common.Uint64Ptr(uint64(policyID)),
			PolicyName: resp.Response.PolicyName,
		}}, nil
	}

	return listAllCamPolicies(ctx, d, client)
}

func camPolicyPositivelyUnattached(p *cam.StrategyInfo) bool {
	if p.AttachEntityCount != nil {
		return *p.AttachEntityCount == 0
	}
	if p.Attachments != nil {
		return *p.Attachments == 0
	}
	return false
}

func isCamPolicyNotFoundError(err error) bool {
	var sdkErr *tcerrors.TencentCloudSDKError
	if stderrors.As(err, &sdkErr) {
		return sdkErr.Code == "ResourceNotFound.PolicyIdNotFound" || sdkErr.Code == "InvalidParameter.PolicyIdError"
	}
	return false
}

func listCamPolicyAttachmentsForPolicy(ctx context.Context, d *plugin.QueryData, client *cam.Client, policy *cam.StrategyInfo) error {
	plugin.Logger(ctx).Debug("listCamPolicyAttachmentsForPolicy", "policy_id", *policy.PolicyId)

	pageSize := uint64(utils.PageSizeFromLimit(d, int64(camDefaultPageSize), 1))
	var page uint64 = 1

	for {
		d.WaitForListRateLimit(ctx)
		req := cam.NewListEntitiesForPolicyRequest()
		req.PolicyId = policy.PolicyId
		req.Rp = common.Uint64Ptr(pageSize)
		req.Page = common.Uint64Ptr(page)

		utils.LogRequest(ctx, "listCamPolicyAttachmentsForPolicy", req)
		resp, err := client.ListEntitiesForPolicyWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listCamPolicyAttachmentsForPolicy", "request_error", err)
			return err
		}
		if resp == nil || resp.Response == nil {
			return nil
		}

		items := resp.Response.List
		for _, e := range items {
			if e == nil {
				continue
			}
			d.StreamListItem(ctx, toCamPolicyAttachmentRow(policy, e))
			if d.RowsRemaining(ctx) == 0 {
				return nil
			}
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
	return nil
}

// Row Type

type camPolicyAttachmentRow struct {
	PolicyId       uint64
	PolicyName     string
	EntityId       string
	EntityName     string
	EntityUin      uint64
	RelatedType    uint64
	AttachmentTime *time.Time
	Title          string
	Akas           []string
}

func toCamPolicyAttachmentRow(p *cam.StrategyInfo, e *cam.AttachEntityOfPolicy) camPolicyAttachmentRow {
	row := camPolicyAttachmentRow{
		PolicyId:       utils.PtrUint64(p.PolicyId),
		PolicyName:     utils.PtrString(p.PolicyName),
		EntityId:       utils.PtrString(e.Id),
		EntityName:     utils.PtrString(e.Name),
		EntityUin:      utils.PtrUint64(e.Uin),
		RelatedType:    utils.PtrUint64(e.RelatedType),
		AttachmentTime: utils.ParseTimestamp(e.AttachmentTime),
	}

	if row.EntityName != "" {
		row.Title = row.EntityName
	} else {
		row.Title = row.EntityId
	}

	if row.PolicyId != 0 && row.EntityId != "" {
		row.Akas = []string{"tencentcloud:cam:policy-attachment:" + strconv.FormatUint(row.PolicyId, 10) + "/" + row.EntityId}
	}
	return row
}
