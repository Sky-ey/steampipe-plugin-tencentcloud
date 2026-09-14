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

func TableTencentcloudCamUser() *plugin.Table {
	return &plugin.Table{
		Name:        "tencentcloud_cam_user",
		Description: "Tencent Cloud CAM sub-accounts (IAM users).",
		List: &plugin.ListConfig{
			Hydrate: listCamUserRows,
			Tags:    map[string]string{"service": "cam", "action": "ListUsers"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getCamUser,
			Tags:       map[string]string{"service": "cam", "action": "GetUser"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "uin", Type: proto.ColumnType_STRING, Description: "The unique UIN of the sub-account.", Transform: transform.FromField("Uin")},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The username of the sub-account.", Transform: transform.FromField("Name")},
			{Name: "uid", Type: proto.ColumnType_INT, Description: "The UID of the sub-account, unique among users who are message recipients.", Transform: transform.FromField("Uid")},
			{Name: "remark", Type: proto.ColumnType_STRING, Description: "The remark of the sub-account.", Transform: transform.FromField("Remark")},
			{Name: "console_login", Type: proto.ColumnType_INT, Description: "Whether the sub-account is allowed to log in to the console (0 = no, 1 = yes).", Transform: transform.FromField("ConsoleLogin")},
			{Name: "phone_num", Type: proto.ColumnType_STRING, Description: "The mobile number bound to the sub-account.", Transform: transform.FromField("PhoneNum")},
			{Name: "country_code", Type: proto.ColumnType_STRING, Description: "The country/area code of the bound mobile number.", Transform: transform.FromField("CountryCode")},
			{Name: "email", Type: proto.ColumnType_STRING, Description: "The email address bound to the sub-account.", Transform: transform.FromField("Email")},
			{Name: "create_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time the sub-account was created.", Transform: transform.FromField("CreateTime").NullIfZero()},
			{Name: "nick_name", Type: proto.ColumnType_STRING, Description: "The nickname of the sub-account.", Transform: transform.FromField("NickName")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listCamUserRows(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listCamUserRows", "quals", d.EqualsQuals)

	client := &cam.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listCamUserRows", "client_init_error", err)
		return nil, err
	}

	d.WaitForListRateLimit(ctx)
	req := cam.NewListUsersRequest()
	utils.LogRequest(ctx, "listCamUserRows", req)
	resp, err := client.ListUsersWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("listCamUserRows", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil {
		return nil, nil
	}

	for _, u := range resp.Response.Data {
		if u == nil {
			continue
		}
		d.StreamListItem(ctx, toCamUserRow(u))
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}

// Get Function

func getCamUser(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	name := d.EqualsQuals["name"].GetStringValue()
	plugin.Logger(ctx).Debug("getCamUser", "name", name)

	client := &cam.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getCamUser", "client_init_error", err)
		return nil, err
	}

	if name == "" {
		return nil, nil
	}

	req := cam.NewGetUserRequest()
	req.Name = common.StringPtr(name)

	utils.LogRequest(ctx, "getCamUser", req)
	resp, err := client.GetUserWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getCamUser", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil {
		return nil, nil
	}

	sub := &cam.SubAccountInfo{
		Uin:          resp.Response.Uin,
		Name:         resp.Response.Name,
		Uid:          resp.Response.Uid,
		Remark:       resp.Response.Remark,
		ConsoleLogin: resp.Response.ConsoleLogin,
		PhoneNum:     resp.Response.PhoneNum,
		CountryCode:  resp.Response.CountryCode,
		Email:        resp.Response.Email,
	}
	return toCamUserRow(sub), nil
}

// Row Type

type camUserRow struct {
	Uin          uint64
	Name         string
	Uid          uint64
	Remark       string
	ConsoleLogin uint64
	PhoneNum     string
	CountryCode  string
	Email        string
	CreateTime   *time.Time
	NickName     string
	Title        string
	Akas         []string
}

func toCamUserRow(u *cam.SubAccountInfo) camUserRow {
	row := camUserRow{
		Uin:          utils.PtrUint64(u.Uin),
		Name:         utils.PtrString(u.Name),
		Uid:          utils.PtrUint64(u.Uid),
		Remark:       utils.PtrString(u.Remark),
		ConsoleLogin: utils.PtrUint64(u.ConsoleLogin),
		PhoneNum:     utils.PtrString(u.PhoneNum),
		CountryCode:  utils.PtrString(u.CountryCode),
		Email:        utils.PtrString(u.Email),
		CreateTime:   utils.ParseTimestamp(u.CreateTime),
		NickName:     utils.PtrString(u.NickName),
	}

	if row.Name != "" {
		row.Title = row.Name
	} else {
		row.Title = strconv.FormatUint(row.Uin, 10)
	}

	if row.Uin != 0 {
		row.Akas = []string{"tencentcloud:cam:user:" + strconv.FormatUint(row.Uin, 10)}
	}
	return row
}
