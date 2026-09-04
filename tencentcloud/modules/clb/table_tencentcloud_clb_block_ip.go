package clb

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"time"

	clb "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/clb/v20180317"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudClbBlockIP() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_clb_block_ip",
		Description:       "Tencent Cloud Load Balancer (CLB) blocked IP (blacklist) entries.",
		GetMatrixItemFunc: buildClbRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "load_balancer_id", Require: plugin.Optional},
			},
			Hydrate: listClbBlockIPs,
			Tags:    map[string]string{"service": "clb", "action": "DescribeBlockIPList"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "load_balancer_id", Type: proto.ColumnType_STRING, Description: "The ID of the load balancer to which the blocklist belongs.", Transform: transform.FromField("LoadBalancerId")},
			{Name: "ip", Type: proto.ColumnType_IPADDR, Description: "The blocked (blacklisted) IP.", Transform: transform.FromField("IP")},
			{Name: "client_ip_field", Type: proto.ColumnType_STRING, Description: "The field used to obtain the real client IP.", Transform: transform.FromField("ClientIPField")},
			{Name: "create_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the IP was added to the blocklist.", Transform: transform.FromField("CreateTime").NullIfZero()},
			{Name: "expire_time", Type: proto.ColumnType_TIMESTAMP, Description: "The expiration time of the blocked IP.", Transform: transform.FromField("ExpireTime").NullIfZero()},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the load balancer resides."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource."},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource."},
		}),
	}
}

// List Function

func listClbBlockIPs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listClbBlockIPs", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &clb.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listClbBlockIPs", "client_init_error", err)
		return nil, err
	}
	region := d.EqualsQualString(utils.MatrixKeyRegion)

	loadBalancerIDs, err := resolveClbLoadBalancerIDs(ctx, d, client)
	if err != nil {
		plugin.Logger(ctx).Error("listClbBlockIPs", "request_error", err)
		return nil, err
	}
	pageSize := uint64(utils.PageSizeFromLimit(d, 100, 1))

	for _, loadBalancerID := range loadBalancerIDs {
		var offset uint64
		for {
			d.WaitForListRateLimit(ctx)

			req := clb.NewDescribeBlockIPListRequest()
			req.LoadBalancerId = &loadBalancerID
			req.Offset = &offset
			req.Limit = &pageSize

			utils.LogRequest(ctx, "listClbBlockIPs", req)
			resp, err := client.DescribeBlockIPListWithContext(ctx, req)
			if err != nil {
				plugin.Logger(ctx).Error("listClbBlockIPs", "request_error", err)
				return nil, err
			}
			if resp == nil || resp.Response == nil {
				break
			}

			blockedIPs := resp.Response.BlockedIPList
			if resp.Response == nil {
				break
			}
			clientIPField := utils.PtrString(resp.Response.ClientIPField)
			for _, blockedIP := range blockedIPs {
				if blockedIP == nil {
					continue
				}
				row := toClbBlockIPRow(blockedIP, loadBalancerID, clientIPField, region)
				d.StreamListItem(ctx, row)
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}

			if utils.PageDone(int64(offset), int64(len(blockedIPs)), utils.PtrUint64(resp.Response.BlockedIPCount)) {
				break
			}
			offset += uint64(len(blockedIPs))
		}
	}
	return nil, nil
}

// Row Type

type clbBlockIPRow struct {
	LoadBalancerId string
	IP             string
	ClientIPField  string
	CreateTime     *time.Time
	ExpireTime     *time.Time
	Region         string
	Title          string
	Akas           []string
}

func toClbBlockIPRow(blockedIP *clb.BlockedIP, loadBalancerID, clientIPField, region string) clbBlockIPRow {
	row := clbBlockIPRow{
		LoadBalancerId: loadBalancerID,
		IP:             utils.PtrString(blockedIP.IP),
		ClientIPField:  clientIPField,
		CreateTime:     utils.ParseTimestamp(blockedIP.CreateTime),
		ExpireTime:     utils.ParseTimestamp(blockedIP.ExpireTime),
		Region:         region,
	}
	row.Title = row.IP
	if row.Title == "" {
		row.Title = row.LoadBalancerId
	}
	if row.IP != "" {
		row.Akas = []string{"tencentcloud:clb:block-ip:" + row.LoadBalancerId + "/" + row.IP}
	}
	return row
}
