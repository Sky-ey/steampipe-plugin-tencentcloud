package cam

import (
	"context"
	"strconv"
	"time"

	cam "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/cam/v20190116"
	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
)

// Table Definition

func TableTencentcloudCamPolicy() *plugin.Table {
	return &plugin.Table{
		Name:        "tencentcloud_cam_policy",
		Description: "Tencent Cloud CAM policies, both preset (QCS) and custom (Local).",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "scope", Require: plugin.Optional},
			},
			Hydrate: listCamPolicyRows,
			Tags:    map[string]string{"service": "cam", "action": "ListPolicies"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("policy_id"),
			Hydrate:    getCamPolicy,
			Tags:       map[string]string{"service": "cam", "action": "GetPolicy"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getCamPolicyDetail,
				Tags: map[string]string{"service": "cam", "action": "GetPolicy"},
			},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "policy_id", Type: proto.ColumnType_INT, Description: "The unique ID of the policy.", Transform: transform.FromField("PolicyId")},
			{Name: "policy_name", Type: proto.ColumnType_STRING, Description: "The name of the policy.", Transform: transform.FromField("PolicyName")},
			{Name: "type", Type: proto.ColumnType_INT, Description: "The type of the policy (1 = custom policy, 2 = preset policy).", Transform: transform.FromField("Type")},
			{Name: "scope", Type: proto.ColumnType_STRING, Description: "The effective scope of the policy (QCS = preset, Local = custom).", Transform: transform.FromField("Scope")},
			{Name: "create_mode", Type: proto.ColumnType_INT, Description: "How the policy was created (1 = via console, 2 = via syntax).", Transform: transform.FromField("CreateMode")},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the policy.", Transform: transform.FromField("Description")},
			{Name: "service_type", Type: proto.ColumnType_STRING, Description: "The product associated with the policy.", Transform: transform.FromField("ServiceType")},
			{Name: "attachments", Type: proto.ColumnType_INT, Description: "The number of principals attached to the policy.", Transform: transform.FromField("Attachments")},
			{Name: "attach_entity_count", Type: proto.ColumnType_INT, Description: "The number of entities attached to the policy.", Transform: transform.FromField("AttachEntityCount")},
			{Name: "attach_entity_boundary_count", Type: proto.ColumnType_INT, Description: "The number of entities using this policy as a permissions boundary.", Transform: transform.FromField("AttachEntityBoundaryCount")},
			{Name: "is_attached", Type: proto.ColumnType_INT, Description: "Whether the caller is attached to the policy (0 = no, 1 = yes).", Transform: transform.FromField("IsAttached")},
			{Name: "deactived", Type: proto.ColumnType_INT, Description: "Whether the policy has been deactivated (0 = no, 1 = yes).", Transform: transform.FromField("Deactived")},
			{Name: "deactived_detail", Type: proto.ColumnType_JSON, Description: "The list of deactivated products.", Transform: transform.FromField("DeactivedDetail")},
			{Name: "is_service_linked_policy", Type: proto.ColumnType_INT, Description: "Whether this is a service-linked policy (0 = no, 1 = yes).", Transform: transform.FromField("IsServiceLinkedPolicy")},
			{Name: "add_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time the policy was created.", Transform: transform.FromField("AddTime").NullIfZero()},
			{Name: "update_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time the policy was last updated.", Transform: transform.FromField("UpdateTime").NullIfZero()},
			{Name: "policy_document", Type: proto.ColumnType_JSON, Description: "The policy document as JSON, only available from the GetPolicy detail call.", Hydrate: getCamPolicyDetail, Transform: transform.FromField("PolicyDocument").Transform(transform.UnmarshalYAML)},
			{Name: "preset_alias", Type: proto.ColumnType_STRING, Description: "The alias of a preset policy, only available from the GetPolicy detail call.", Hydrate: getCamPolicyDetail, Transform: transform.FromField("PresetAlias")},
			{Name: "is_service_linked_role_policy", Type: proto.ColumnType_INT, Description: "Whether this is a service-linked role policy, only available from the GetPolicy detail call.", Hydrate: getCamPolicyDetail, Transform: transform.FromField("IsServiceLinkedRolePolicy")},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The tags attached to the policy as a map of key -> value, only available from the GetPolicy detail call.", Hydrate: getCamPolicyDetail, Transform: transform.FromField("Tags")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listCamPolicyRows(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listCamPolicyRows", "quals", d.EqualsQuals)

	client := &cam.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listCamPolicyRows", "client_init_error", err)
		return nil, err
	}

	scope := d.EqualsQualString("scope")
	if scope == "" {
		scope = "All"
	}

	pageSize := uint64(utils.PageSizeFromLimit(d, int64(camDefaultPageSize), 1))
	var page uint64 = 1

	for {
		d.WaitForListRateLimit(ctx)
		req := cam.NewListPoliciesRequest()
		req.Rp = common.Uint64Ptr(pageSize)
		req.Page = common.Uint64Ptr(page)
		req.Scope = common.StringPtr(scope)

		utils.LogRequest(ctx, "listCamPolicyRows", req)
		resp, err := client.ListPoliciesWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listCamPolicyRows", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			return nil, nil
		}

		items := resp.Response.List
		for _, p := range items {
			if p == nil {
				continue
			}
			d.StreamListItem(ctx, toCamPolicyRow(p))
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
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
	return nil, nil
}

// Get Function

func getCamPolicy(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	policyID := d.EqualsQuals["policy_id"].GetInt64Value()
	plugin.Logger(ctx).Debug("getCamPolicy", "policy_id", policyID)

	client := &cam.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getCamPolicy", "client_init_error", err)
		return nil, err
	}

	if policyID == 0 {
		return nil, nil
	}

	req := cam.NewGetPolicyRequest()
	req.PolicyId = common.Uint64Ptr(uint64(policyID))

	utils.LogRequest(ctx, "getCamPolicy", req)
	resp, err := client.GetPolicyWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getCamPolicy", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil {
		return nil, nil
	}

	info := &cam.StrategyInfo{
		PolicyId:    common.Uint64Ptr(uint64(policyID)),
		PolicyName:  resp.Response.PolicyName,
		Type:        resp.Response.Type,
		AddTime:     resp.Response.AddTime,
		UpdateTime:  resp.Response.UpdateTime,
		Description: resp.Response.Description,
	}
	return toCamPolicyRow(info), nil
}

func getCamPolicyDetail(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (any, error) {
	var policyID uint64
	switch p := h.Item.(type) {
	case camPolicyRow:
		policyID = p.PolicyId
	case *cam.StrategyInfo:
		if p != nil && p.PolicyId != nil {
			policyID = *p.PolicyId
		}
	}
	if policyID == 0 {
		return nil, nil
	}

	client := &cam.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getCamPolicyDetail", "client_init_error", err)
		return nil, err
	}

	req := cam.NewGetPolicyRequest()
	req.PolicyId = common.Uint64Ptr(policyID)

	utils.LogRequest(ctx, "getCamPolicyDetail", req)
	resp, err := client.GetPolicyWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getCamPolicyDetail", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil {
		return nil, nil
	}
	return resp.Response, nil
}

// Row Type

type camPolicyRow struct {
	PolicyId                  uint64
	PolicyName                string
	Type                      uint64
	Scope                     string
	CreateMode                uint64
	Description               string
	ServiceType               string
	Attachments               uint64
	AttachEntityCount         int64
	AttachEntityBoundaryCount int64
	IsAttached                uint64
	Deactived                 uint64
	DeactivedDetail           []string
	IsServiceLinkedPolicy     uint64
	AddTime                   *time.Time
	UpdateTime                *time.Time
	Title                     string
	Akas                      []string
}

func toCamPolicyRow(p *cam.StrategyInfo) camPolicyRow {
	row := camPolicyRow{
		PolicyId:                  utils.PtrUint64(p.PolicyId),
		PolicyName:                utils.PtrString(p.PolicyName),
		Type:                      utils.PtrUint64(p.Type),
		CreateMode:                utils.PtrUint64(p.CreateMode),
		Description:               utils.PtrString(p.Description),
		ServiceType:               utils.PtrString(p.ServiceType),
		Attachments:               utils.PtrUint64(p.Attachments),
		AttachEntityCount:         utils.PtrInt64(p.AttachEntityCount),
		AttachEntityBoundaryCount: utils.PtrInt64(p.AttachEntityBoundaryCount),
		IsAttached:                utils.PtrUint64(p.IsAttached),
		Deactived:                 utils.PtrUint64(p.Deactived),
		DeactivedDetail:           utils.PtrStringSlice(p.DeactivedDetail),
		IsServiceLinkedPolicy:     utils.PtrUint64(p.IsServiceLinkedPolicy),
		AddTime:                   utils.ParseTimestamp(p.AddTime),
		UpdateTime:                utils.ParseTimestamp(p.UpdateTime),
	}

	if row.Type == 2 {
		row.Scope = "QCS"
	} else {
		row.Scope = "Local"
	}

	if row.PolicyName != "" {
		row.Title = row.PolicyName
	} else {
		row.Title = strconv.FormatUint(row.PolicyId, 10)
	}

	if row.PolicyId != 0 {
		row.Akas = []string{"tencentcloud:cam:policy:" + strconv.FormatUint(row.PolicyId, 10)}
	}
	return row
}
