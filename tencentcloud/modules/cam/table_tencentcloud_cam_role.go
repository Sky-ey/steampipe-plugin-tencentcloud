package cam

import (
	"context"
	"time"

	cam "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/cam/v20190116"
	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
)

// Table Definition

func TableTencentcloudCamRole() *plugin.Table {
	return &plugin.Table{
		Name:        "tencentcloud_cam_role",
		Description: "Tencent Cloud CAM roles, including cross-account and service-linked roles.",
		List: &plugin.ListConfig{
			Hydrate: listCamRoles,
			Tags:    map[string]string{"service": "cam", "action": "DescribeRoleList"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("role_name"),
			Hydrate:    getCamRole,
			Tags:       map[string]string{"service": "cam", "action": "GetRole"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "role_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the role.", Transform: transform.FromField("RoleId")},
			{Name: "role_name", Type: proto.ColumnType_STRING, Description: "The name of the role.", Transform: transform.FromField("RoleName")},
			{Name: "role_type", Type: proto.ColumnType_STRING, Description: "The type of the role (user, system, service_linked).", Transform: transform.FromField("RoleType")},
			{Name: "role_arn", Type: proto.ColumnType_STRING, Description: "The ARN of the role.", Transform: transform.FromField("RoleArn")},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the role.", Transform: transform.FromField("Description")},
			{Name: "console_login", Type: proto.ColumnType_INT, Description: "Whether the role is allowed to log in to the console (0 = no, 1 = yes).", Transform: transform.FromField("ConsoleLogin")},
			{Name: "session_duration", Type: proto.ColumnType_INT, Description: "The maximum session duration in seconds for the role.", Transform: transform.FromField("SessionDuration")},
			{Name: "policy_document", Type: proto.ColumnType_JSON, Description: "The trust policy document of the role, i.e. who can assume this role.", Transform: transform.FromField("PolicyDocument").Transform(transform.UnmarshalYAML)},
			{Name: "deletion_task_id", Type: proto.ColumnType_STRING, Description: "The task identifier used to track deletion of a service-linked role.", Transform: transform.FromField("DeletionTaskId")},
			{Name: "add_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time the role was created.", Transform: transform.FromField("AddTime").NullIfZero()},
			{Name: "update_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time the role was last updated.", Transform: transform.FromField("UpdateTime").NullIfZero()},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The tags attached to the role as a map of key -> value.", Transform: transform.FromField("Tags")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listCamRoles(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listCamRoles", "quals", d.EqualsQuals)

	client := &cam.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listCamRoles", "client_init_error", err)
		return nil, err
	}

	pageSize := uint64(utils.PageSizeFromLimit(d, int64(camDefaultPageSize), 1))
	var page uint64 = 1

	for {
		d.WaitForListRateLimit(ctx)
		req := cam.NewDescribeRoleListRequest()
		req.Rp = common.Uint64Ptr(pageSize)
		req.Page = common.Uint64Ptr(page)

		utils.LogRequest(ctx, "listCamRoles", req)
		resp, err := client.DescribeRoleListWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listCamRoles", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			return nil, nil
		}

		items := resp.Response.List
		for _, r := range items {
			if r == nil {
				continue
			}
			d.StreamListItem(ctx, toCamRoleRow(r))
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

func getCamRole(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	roleName := d.EqualsQuals["role_name"].GetStringValue()
	plugin.Logger(ctx).Info("getCamRole", "role_name", roleName)

	client := &cam.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getCamRole", "client_init_error", err)
		return nil, err
	}

	if roleName == "" {
		return nil, nil
	}

	req := cam.NewGetRoleRequest()
	req.RoleName = common.StringPtr(roleName)

	utils.LogRequest(ctx, "getCamRole", req)
	resp, err := client.GetRoleWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getCamRole", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || resp.Response.RoleInfo == nil {
		return nil, nil
	}

	return toCamRoleRow(resp.Response.RoleInfo), nil
}

// Row Type

type camRoleRow struct {
	RoleId          string
	RoleName        string
	RoleType        string
	RoleArn         string
	Description     string
	ConsoleLogin    uint64
	SessionDuration uint64
	PolicyDocument  string
	DeletionTaskId  string
	AddTime         *time.Time
	UpdateTime      *time.Time
	Tags            map[string]string
	Title           string
	Akas            []string
}

func toCamRoleRow(r *cam.RoleInfo) camRoleRow {
	row := camRoleRow{
		RoleId:          utils.PtrString(r.RoleId),
		RoleName:        utils.PtrString(r.RoleName),
		RoleType:        utils.PtrString(r.RoleType),
		RoleArn:         utils.PtrString(r.RoleArn),
		Description:     utils.PtrString(r.Description),
		ConsoleLogin:    utils.PtrUint64(r.ConsoleLogin),
		SessionDuration: utils.PtrUint64(r.SessionDuration),
		PolicyDocument:  utils.PtrString(r.PolicyDocument),
		DeletionTaskId:  utils.PtrString(r.DeletionTaskId),
		AddTime:         utils.ParseTimestamp(r.AddTime),
		UpdateTime:      utils.ParseTimestamp(r.UpdateTime),
		Tags:            RoleTagsToMap(r.Tags),
	}

	if row.RoleName != "" {
		row.Title = row.RoleName
	} else {
		row.Title = row.RoleId
	}

	if row.RoleId != "" {
		row.Akas = []string{"tencentcloud:cam:role:" + row.RoleId}
	}
	return row
}
