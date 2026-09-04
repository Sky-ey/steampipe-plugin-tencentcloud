package utils

import (
	"context"

	cam "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/cam/v20190116"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type CommonColumnData struct {
	OwnerUin uint64
	Uin      uint64
	AppId    uint64
}

var GetCallerIdentityMemoized = plugin.HydrateFunc(getCallerIdentityUncached).Memoize()

func getCallerIdentityUncached(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	client := &cam.Client{}
	if err := InitClientInRegion(ctx, d, client, DefaultRegion); err != nil {
		plugin.Logger(ctx).Error("getCallerIdentityUncached", "client_init_error", err)
		return nil, err
	}

	req := cam.NewGetUserAppIdRequest()
	LogRequest(ctx, "getCallerIdentityUncached", req)

	resp, err := client.GetUserAppIdWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getCallerIdentityUncached", "api_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil {
		return nil, nil
	}
	return resp.Response, nil
}

func GetCommonColumns(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (any, error) {
	data, err := GetCallerIdentityMemoized(ctx, d, h)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return &CommonColumnData{}, nil
	}
	resp := data.(*cam.GetUserAppIdResponseParams)
	return &CommonColumnData{
		OwnerUin: ToUint64(resp.OwnerUin),
		Uin:      ToUint64(resp.Uin),
		AppId:    PtrUint64(resp.AppId),
	}, nil
}

func OwnerUinForConnection(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (any, error) {
	data, err := GetCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	return data.(*CommonColumnData).OwnerUin, nil
}

func CommonColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "owner_uin",
			Type:        proto.ColumnType_INT,
			Description: "The Tencent Cloud owner account UIN (OwnerUin) that owns the resource. Used as the connection key column for multi-account aggregator queries.",
			Hydrate:     GetCommonColumns,
			Transform:   transform.FromField("OwnerUin"),
		},
		{
			Name:        "caller_uin",
			Type:        proto.ColumnType_INT,
			Description: "The UIN of the calling identity — the master account itself when called with root credentials, or the sub-account/role UIN when called via CAM.",
			Hydrate:     GetCommonColumns,
			Transform:   transform.FromField("Uin"),
		},
		{
			Name:        "owner_app_id",
			Type:        proto.ColumnType_INT,
			Description: "The numeric AppId of the account — the <appid> suffix of COS bucket names and the resource owner in CAM.",
			Hydrate:     GetCommonColumns,
			Transform:   transform.FromField("AppId"),
		},
	}
}

func WithCommonColumns(columns []*plugin.Column) []*plugin.Column {
	return append(columns, CommonColumns()...)
}
