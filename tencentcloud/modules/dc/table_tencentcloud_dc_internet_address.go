package dc

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"strconv"
	"time"

	dc "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/dc/v20180410"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudDcInternetAddress() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_dc_internet_address",
		Description:       "Tencent Cloud Direct Connect internet address (public IP) resources allocated to internet tunnels.",
		GetMatrixItemFunc: buildDcRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "instance_id", Require: plugin.Optional},
				{Name: "addr_type", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "addr_proto", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "status", Require: plugin.Optional, Operators: []string{"="}},
			},
			Hydrate: listDcInternetAddresses,
			Tags:    map[string]string{"service": "dc", "action": "DescribeInternetAddress"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "instance_id", Type: proto.ColumnType_STRING, Description: "The internet tunnel IP address ID.", Transform: transform.FromField("InstanceId")},
			{Name: "subnet", Type: proto.ColumnType_STRING, Description: "The internet tunnel network address.", Transform: transform.FromField("Subnet")},
			{Name: "mask_len", Type: proto.ColumnType_INT, Description: "The mask length of the network address.", Transform: transform.FromField("MaskLen")},
			{Name: "addr_type", Type: proto.ColumnType_INT, Description: "The address type (0: BGP, 1: China Telecom, 2: China Mobile, 3: China Unicom).", Transform: transform.FromField("AddrType")},
			{Name: "addr_proto", Type: proto.ColumnType_INT, Description: "The address protocol (0: IPv4, 1: IPv6).", Transform: transform.FromField("AddrProto")},
			{Name: "status", Type: proto.ColumnType_INT, Description: "The address status (0: in use, 1: disabled, 2: returned).", Transform: transform.FromField("Status")},
			{Name: "resource_region", Type: proto.ColumnType_STRING, Description: "The region of the internet address as reported by the API.", Transform: transform.FromField("ResourceRegion")},
			{Name: "app_id", Type: proto.ColumnType_STRING, Description: "The user (app) ID that owns the address.", Transform: transform.FromField("AppId")},
			{Name: "reserve_time", Type: proto.ColumnType_INT, Description: "The retention period (in days) of a released IP address.", Transform: transform.FromField("ReserveTime")},
			{Name: "apply_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time the address was applied for.", Transform: transform.FromField("ApplyTime").NullIfZero()},
			{Name: "stop_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time the address was disabled.", Transform: transform.FromField("StopTime").NullIfZero()},
			{Name: "release_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time the address was returned.", Transform: transform.FromField("ReleaseTime").NullIfZero()},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the address resides (matrix region).", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("InstanceId")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listDcInternetAddresses(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listDcInternetAddresses", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &dc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listDcInternetAddresses", "client_init_error", err)
		return nil, err
	}

	matrixRegion := d.EqualsQualString(utils.MatrixKeyRegion)

	pageSize := utils.PageSizeFromLimit(d, 100, 1)
	var offset int64 = 0

	filters := buildDcInternetAddressFilters(d.EqualsQuals)

	for {
		d.WaitForListRateLimit(ctx)

		req := dc.NewDescribeInternetAddressRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		if len(filters) > 0 {
			req.Filters = filters
		}

		utils.LogRequest(ctx, "listDcInternetAddresses", req)
		resp, err := client.DescribeInternetAddressWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listDcInternetAddresses", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		items := resp.Response.Subnets
		total := utils.PtrInt64(resp.Response.TotalCount)

		for _, item := range items {
			if item == nil {
				continue
			}
			row := toDcInternetAddressRow(item)
			row.Region = matrixRegion
			d.StreamListItem(ctx, row)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if utils.PageDone(offset, int64(len(items)), total) {
			break
		}
		offset += int64(len(items))
	}

	return nil, nil
}

func buildDcInternetAddressFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*dc.Filter {
	filters := make([]*dc.Filter, 0, 4)

	if q := equalsQuals["instance_id"]; q != nil {
		if v := q.GetStringValue(); v != "" {
			filters = append(filters, newDcFilter("InstanceIds", v))
		}
	}
	if q := equalsQuals["addr_type"]; q != nil {
		filters = append(filters, newDcFilter("AddrType", strconv.FormatInt(q.GetInt64Value(), 10)))
	}
	if q := equalsQuals["addr_proto"]; q != nil {
		filters = append(filters, newDcFilter("AddrProto", strconv.FormatInt(q.GetInt64Value(), 10)))
	}
	if q := equalsQuals["status"]; q != nil {
		filters = append(filters, newDcFilter("Status", strconv.FormatInt(q.GetInt64Value(), 10)))
	}
	return filters
}

// Row Type

type dcInternetAddressRow struct {
	InstanceId     string
	Subnet         string
	MaskLen        int64
	AddrType       int64
	AddrProto      int64
	Status         int64
	ResourceRegion string
	AppId          int64
	ReserveTime    int64
	ApplyTime      *time.Time
	StopTime       *time.Time
	ReleaseTime    *time.Time
	Region         string
	Akas           []string
}

func toDcInternetAddressRow(s *dc.InternetAddressDetail) dcInternetAddressRow {
	row := dcInternetAddressRow{
		InstanceId:     utils.PtrString(s.InstanceId),
		Subnet:         utils.PtrString(s.Subnet),
		MaskLen:        utils.PtrInt64(s.MaskLen),
		AddrType:       utils.PtrInt64(s.AddrType),
		AddrProto:      utils.PtrInt64(s.AddrProto),
		Status:         utils.PtrInt64(s.Status),
		ResourceRegion: utils.PtrString(s.Region),
		AppId:          utils.PtrInt64(s.AppId),
		ReserveTime:    utils.PtrInt64(s.ReserveTime),
		ApplyTime:      utils.ParseTimestamp(s.ApplyTime),
		StopTime:       utils.ParseTimestamp(s.StopTime),
		ReleaseTime:    utils.ParseTimestamp(s.ReleaseTime),
	}

	if row.InstanceId != "" {
		row.Akas = []string{"tencentcloud:dc:internet-address:" + row.InstanceId}
	}
	return row
}
