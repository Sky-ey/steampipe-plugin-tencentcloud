package clb

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"strconv"
	"time"

	clb "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/clb/v20180317"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudClbTargetGroupInstance() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_clb_target_group_instance",
		Description:       "Tencent Cloud Load Balancer (CLB) backend servers registered in target groups.",
		GetMatrixItemFunc: buildClbRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "target_group_id", Require: plugin.Optional},
				{Name: "instance_id", Require: plugin.Optional},
			},
			Hydrate: listClbTargetGroupInstances,
			Tags:    map[string]string{"service": "clb", "action": "DescribeTargetGroupInstances"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "target_group_id", Type: proto.ColumnType_STRING, Description: "The ID of the target group to which the backend server is registered.", Transform: transform.FromField("TargetGroupId")},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The type of the real server (CVM, ENI).", Transform: transform.FromField("Type")},
			{Name: "instance_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the real server.", Transform: transform.FromField("InstanceId")},
			{Name: "instance_name", Type: proto.ColumnType_STRING, Description: "The name of the real server instance."},
			{Name: "port", Type: proto.ColumnType_INT, Description: "The listening port of the backend service; returns 0 for port range listeners."},
			{Name: "weight", Type: proto.ColumnType_INT, Description: "The forwarding weight of the real server, range [0, 100]."},
			{Name: "public_ip_addresses", Type: proto.ColumnType_JSON, Description: "The public network IPs of the real server."},
			{Name: "private_ip_addresses", Type: proto.ColumnType_JSON, Description: "The private network IPs of the real server."},
			{Name: "eni_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the elastic network interface.", Transform: transform.FromField("EniId")},
			{Name: "zone_id", Type: proto.ColumnType_INT, Description: "The availability zone ID of the backend service.", Transform: transform.FromField("ZoneId")},
			{Name: "registered_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the real server was bound to the target group.", Transform: transform.FromField("RegisteredTime").NullIfZero()},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the target group instance resides."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource."},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource."},
		}),
	}
}

// List Function

func listClbTargetGroupInstances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listClbTargetGroupInstances", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &clb.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listClbTargetGroupInstances", "client_init_error", err)
		return nil, err
	}
	region := d.EqualsQualString(utils.MatrixKeyRegion)

	targetGroupIDs, err := resolveClbTargetGroupIDs(ctx, d, client)
	if err != nil {
		plugin.Logger(ctx).Error("listClbTargetGroupInstances", "resolve_target_group_ids_error", err)
		return nil, err
	}

	instanceID := d.EqualsQualString("instance_id")
	pageSize := uint64(utils.PageSizeFromLimit(d, 100, 1))

	for _, targetGroupID := range targetGroupIDs {
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}

		filters := []*clb.Filter{newClbFilter("TargetGroupId", targetGroupID)}
		if instanceID != "" {
			filters = append(filters, newClbFilter("InstanceId", instanceID))
		}

		var offset uint64
		for {
			d.WaitForListRateLimit(ctx)
			req := clb.NewDescribeTargetGroupInstancesRequest()
			req.Offset = &offset
			req.Limit = &pageSize
			req.Filters = filters

			utils.LogRequest(ctx, "listClbTargetGroupInstances", req)
			resp, err := client.DescribeTargetGroupInstancesWithContext(ctx, req)
			if err != nil {
				plugin.Logger(ctx).Error("listClbTargetGroupInstances", "request_error", err)
				return nil, err
			}
			if resp == nil || resp.Response == nil {
				break
			}

			instances := resp.Response.TargetGroupInstanceSet
			for _, instance := range instances {
				if instance == nil {
					continue
				}
				row := toClbTargetGroupInstanceRow(instance, region)
				d.StreamListItem(ctx, row)
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}

			if utils.PageDone(int64(offset), int64(len(instances)), utils.PtrUint64(resp.Response.TotalCount)) {
				break
			}
			offset += uint64(len(instances))
		}
	}
	return nil, nil
}

// Row Type

type clbTargetGroupInstanceRow struct {
	TargetGroupId      string
	Type               string
	InstanceId         string
	InstanceName       string
	Port               uint64
	Weight             uint64
	PublicIpAddresses  []string
	PrivateIpAddresses []string
	EniId              string
	ZoneId             uint64
	RegisteredTime     *time.Time
	Region             string
	Title              string
	Akas               []string
}

func toClbTargetGroupInstanceRow(instance *clb.TargetGroupBackend, region string) clbTargetGroupInstanceRow {
	row := clbTargetGroupInstanceRow{
		TargetGroupId:      utils.PtrString(instance.TargetGroupId),
		Type:               utils.PtrString(instance.Type),
		InstanceId:         utils.PtrString(instance.InstanceId),
		InstanceName:       utils.PtrString(instance.InstanceName),
		Port:               utils.PtrUint64(instance.Port),
		Weight:             utils.PtrUint64(instance.Weight),
		PublicIpAddresses:  utils.PtrStringSlice(instance.PublicIpAddresses),
		PrivateIpAddresses: utils.PtrStringSlice(instance.PrivateIpAddresses),
		EniId:              utils.PtrString(instance.EniId),
		ZoneId:             utils.PtrUint64(instance.ZoneId),
		RegisteredTime:     utils.ParseTimestamp(instance.RegisteredTime),
		Region:             region,
	}
	row.Title = row.InstanceName
	if row.Title == "" {
		row.Title = row.InstanceId
	}
	if row.Title == "" && len(row.PrivateIpAddresses) > 0 {
		row.Title = row.PrivateIpAddresses[0]
	}
	if row.Title == "" && len(row.PublicIpAddresses) > 0 {
		row.Title = row.PublicIpAddresses[0]
	}

	identity := row.InstanceId
	if identity == "" {
		identity = row.EniId
	}
	if identity == "" && len(row.PrivateIpAddresses) > 0 {
		identity = row.PrivateIpAddresses[0]
	}
	if identity != "" {
		row.Akas = []string{"tencentcloud:clb:target-group-instance:" + row.TargetGroupId + "/" + identity + ":" + strconv.FormatUint(row.Port, 10)}
	}
	return row
}
