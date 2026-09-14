package cam

import (
	"context"
	"time"

	cam "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/cam/v20190116"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
)

// Table Definition

func TableTencentcloudCamAccessKey() *plugin.Table {
	return &plugin.Table{
		Name:        "tencentcloud_cam_access_key",
		Description: "Tencent Cloud CAM access keys of sub-accounts.",
		List: &plugin.ListConfig{
			ParentHydrate: listCamUsers,
			Hydrate:       listCamAccessKeys,
			Tags:          map[string]string{"service": "cam", "action": "ListAccessKeys"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "access_key_id", Type: proto.ColumnType_STRING, Description: "The access key ID (SecretId).", Transform: transform.FromField("AccessKeyId")},
			{Name: "user_name", Type: proto.ColumnType_STRING, Description: "The username of the sub-account that owns the key.", Transform: transform.FromField("UserName")},
			{Name: "user_uin", Type: proto.ColumnType_INT, Description: "The UIN of the sub-account that owns the key.", Transform: transform.FromField("UserUin")},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of the key (Active, Inactive).", Transform: transform.FromField("Status")},
			{Name: "create_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time the key was created.", Transform: transform.FromField("CreateTime").NullIfZero()},
			{Name: "last_used_date", Type: proto.ColumnType_STRING, Description: "The date the key was last used (YYYY-MM-DD, reported with a one-day delay).", Transform: transform.FromField("LastUsedDate")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listCamAccessKeys(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (any, error) {
	user, ok := h.Item.(*cam.SubAccountInfo)
	if !ok || user == nil || user.Uin == nil {
		return nil, nil
	}
	userUin := *user.Uin

	plugin.Logger(ctx).Debug("listCamAccessKeys", "user_uin", userUin)

	client := &cam.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listCamAccessKeys", "client_init_error", err)
		return nil, err
	}

	d.WaitForListRateLimit(ctx)
	req := cam.NewListAccessKeysRequest()
	req.TargetUin = user.Uin

	utils.LogRequest(ctx, "listCamAccessKeys", req)
	resp, err := client.ListAccessKeysWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("listCamAccessKeys", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil {
		return nil, nil
	}

	keys := resp.Response.AccessKeys
	if len(keys) == 0 {
		return nil, nil
	}

	// Batch-fetch last-used dates (the API accepts up to 10 SecretIds at once).
	secretIds := make([]*string, 0, len(keys))
	for _, k := range keys {
		if k == nil || k.AccessKeyId == nil {
			continue
		}
		secretIds = append(secretIds, k.AccessKeyId)
	}
	lastUsedMap := map[string]string{}
	if len(secretIds) > 0 {
		d.WaitForListRateLimit(ctx)
		luReq := cam.NewGetSecurityLastUsedRequest()
		luReq.SecretIdList = secretIds

		utils.LogRequest(ctx, "listCamAccessKeys", luReq)
		luResp, luErr := client.GetSecurityLastUsedWithContext(ctx, luReq)
		if luErr != nil {
			plugin.Logger(ctx).Error("listCamAccessKeys", "last_used_request_error", luErr)
		} else if luResp != nil && luResp.Response != nil {
			for _, row := range luResp.Response.SecretIdLastUsedRows {
				if row == nil || row.SecretId == nil || row.LastUsedDate == nil {
					continue
				}
				lastUsedMap[*row.SecretId] = *row.LastUsedDate
			}
		}
	}

	userName := utils.PtrString(user.Name)
	for _, k := range keys {
		if k == nil {
			continue
		}
		row := toCamAccessKeyRow(userUin, userName, k)
		if k.AccessKeyId != nil {
			row.LastUsedDate = lastUsedMap[*k.AccessKeyId]
		}
		d.StreamListItem(ctx, row)
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}

// Row Type

type camAccessKeyRow struct {
	AccessKeyId  string
	UserName     string
	UserUin      int64
	Status       string
	CreateTime   *time.Time
	LastUsedDate string
	Title        string
	Akas         []string
}

func toCamAccessKeyRow(userUin uint64, userName string, k *cam.AccessKey) camAccessKeyRow {
	row := camAccessKeyRow{
		AccessKeyId: utils.PtrString(k.AccessKeyId),
		UserName:    userName,
		UserUin:     int64(userUin),
		Status:      utils.PtrString(k.Status),
		CreateTime:  utils.ParseTimestamp(k.CreateTime),
	}

	if row.AccessKeyId != "" {
		row.Title = row.AccessKeyId
		row.Akas = []string{"tencentcloud:cam:access-key:" + row.AccessKeyId}
	}
	return row
}
