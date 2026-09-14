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

func TableTencentcloudCamGroup() *plugin.Table {
	return &plugin.Table{
		Name:        "tencentcloud_cam_group",
		Description: "Tencent Cloud CAM user groups.",
		List: &plugin.ListConfig{
			Hydrate: listCamGroups,
			Tags:    map[string]string{"service": "cam", "action": "ListGroups"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("group_id"),
			Hydrate:    getCamGroup,
			Tags:       map[string]string{"service": "cam", "action": "GetGroup"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getCamGroup,
				Tags: map[string]string{"service": "cam", "action": "GetGroup"},
			},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "group_id", Type: proto.ColumnType_INT, Description: "The unique ID of the user group.", Transform: transform.FromField("GroupId")},
			{Name: "group_name", Type: proto.ColumnType_STRING, Description: "The name of the user group.", Transform: transform.FromField("GroupName")},
			{Name: "remark", Type: proto.ColumnType_STRING, Description: "The description of the user group.", Transform: transform.FromField("Remark")},
			{Name: "create_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time the user group was created.", Transform: transform.FromField("CreateTime").NullIfZero()},
			{Name: "group_num", Type: proto.ColumnType_INT, Description: "The number of members in the user group, only available from the GetGroup detail call.", Hydrate: getCamGroup, Transform: transform.FromField("GroupNum")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listCamGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listCamGroups", "quals", d.EqualsQuals)

	client := &cam.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listCamGroups", "client_init_error", err)
		return nil, err
	}

	pageSize := uint64(utils.PageSizeFromLimit(d, int64(camDefaultPageSize), 1))
	var page uint64 = 1

	for {
		d.WaitForListRateLimit(ctx)
		req := cam.NewListGroupsRequest()
		req.Rp = common.Uint64Ptr(pageSize)
		req.Page = common.Uint64Ptr(page)

		utils.LogRequest(ctx, "listCamGroups", req)
		resp, err := client.ListGroupsWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listCamGroups", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			return nil, nil
		}

		items := resp.Response.GroupInfo
		for _, g := range items {
			if g == nil {
				continue
			}
			d.StreamListItem(ctx, toCamGroupRow(g))
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

func getCamGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (any, error) {
	isSupplement := h.Item != nil
	var groupID uint64
	if isSupplement {
		if row, ok := h.Item.(camGroupRow); ok {
			groupID = row.GroupId
		}
	} else {
		groupID = uint64(d.EqualsQuals["group_id"].GetInt64Value())
	}
	plugin.Logger(ctx).Debug("getCamGroup", "group_id", groupID, "supplement", isSupplement)

	if groupID == 0 {
		return nil, nil
	}

	client := &cam.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getCamGroup", "client_init_error", err)
		return nil, err
	}

	req := cam.NewGetGroupRequest()
	req.GroupId = common.Uint64Ptr(groupID)

	utils.LogRequest(ctx, "getCamGroup", req)
	resp, err := client.GetGroupWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getCamGroup", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil {
		return nil, nil
	}

	if isSupplement {
		return resp.Response, nil
	}

	row := camGroupRow{
		GroupId:    groupID,
		GroupName:  utils.PtrString(resp.Response.GroupName),
		Remark:     utils.PtrString(resp.Response.Remark),
		CreateTime: utils.ParseTimestamp(resp.Response.CreateTime),
		GroupNum:   utils.PtrUint64(resp.Response.GroupNum),
	}
	if row.GroupName != "" {
		row.Title = row.GroupName
	} else {
		row.Title = strconv.FormatUint(row.GroupId, 10)
	}
	if row.GroupId != 0 {
		row.Akas = []string{"tencentcloud:cam:group:" + strconv.FormatUint(row.GroupId, 10)}
	}
	return row, nil
}

// Row Type

type camGroupRow struct {
	GroupId    uint64
	GroupName  string
	Remark     string
	CreateTime *time.Time
	GroupNum   uint64
	Title      string
	Akas       []string
}

func toCamGroupRow(g *cam.GroupInfo) camGroupRow {
	row := camGroupRow{
		GroupId:    utils.PtrUint64(g.GroupId),
		GroupName:  utils.PtrString(g.GroupName),
		Remark:     utils.PtrString(g.Remark),
		CreateTime: utils.ParseTimestamp(g.CreateTime),
	}

	if row.GroupName != "" {
		row.Title = row.GroupName
	} else {
		row.Title = strconv.FormatUint(row.GroupId, 10)
	}

	if row.GroupId != 0 {
		row.Akas = []string{"tencentcloud:cam:group:" + strconv.FormatUint(row.GroupId, 10)}
	}
	return row
}
