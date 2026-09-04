package tke

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"time"

	tke "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/tke/v20180525"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudTkeClusterInstance() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_tke_cluster_instance",
		Description:       "TKE cluster instances (CVM instances managed by a Tencent Kubernetes Engine cluster).",
		GetMatrixItemFunc: buildTkeRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "cluster_id", Require: plugin.Optional},
				{Name: "instance_id", Require: plugin.Optional},
				{Name: "instance_role", Require: plugin.Optional},
				{Name: "node_pool_id", Require: plugin.Optional},
			},
			Hydrate: listTkeClusterInstances,
			Tags:    map[string]string{"service": "tke", "action": "DescribeClusterInstances"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "instance_id", Type: proto.ColumnType_STRING, Description: "The ID of the cluster instance (CVM instance).", Transform: transform.FromField("InstanceId")},
			{Name: "cluster_id", Type: proto.ColumnType_STRING, Description: "The TKE cluster ID the instance belongs to.", Transform: transform.FromField("ClusterId")},
			{Name: "instance_role", Type: proto.ColumnType_STRING, Description: "The instance role (MASTER, WORKER, ETCD, MASTER_ETCD, ALL).", Transform: transform.FromField("InstanceRole")},
			{Name: "instance_state", Type: proto.ColumnType_STRING, Description: "The instance status (running, initializing, or failed).", Transform: transform.FromField("InstanceState")},
			{Name: "failed_reason", Type: proto.ColumnType_STRING, Description: "The reason for instance exception or initialization failure.", Transform: transform.FromField("FailedReason")},
			{Name: "drain_status", Type: proto.ColumnType_STRING, Description: "Whether the instance is drained.", Transform: transform.FromField("DrainStatus")},
			{Name: "lan_ip", Type: proto.ColumnType_STRING, Description: "The private IP of the instance.", Transform: transform.FromField("LanIP")},
			{Name: "node_pool_id", Type: proto.ColumnType_STRING, Description: "The node pool ID the instance belongs to.", Transform: transform.FromField("NodePoolId")},
			{Name: "autoscaling_group_id", Type: proto.ColumnType_STRING, Description: "The auto-scaling group ID of the instance.", Transform: transform.FromField("AutoscalingGroupId")},
			{Name: "instance_advanced_settings", Type: proto.ColumnType_JSON, Description: "Advanced instance settings (desired pods, GPU args, taints, user scripts, etc.).", Transform: transform.FromField("InstanceAdvancedSettings")},
			{Name: "created_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the instance.", Transform: transform.FromField("CreatedTime").NullIfZero()},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the instance resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("InstanceId")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listTkeClusterInstances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listTkeClusterInstances", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &tke.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listTkeClusterInstances", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	clusterIds, err := resolveTkeClusterIDs(ctx, d, client)
	if err != nil {
		plugin.Logger(ctx).Error("listTkeClusterInstances", "resolve_cluster_ids_error", err)
		return nil, err
	}

	instanceIds := utils.SingleIdFromQual(d.EqualsQuals, "instance_id")
	instanceRole := d.EqualsQualString("instance_role")
	filters := buildTkeClusterInstanceFilters(d.EqualsQuals)

	pageSize := utils.PageSizeFromLimit(d, 100, 1)

	for _, clusterId := range clusterIds {
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}

		var offset int64
		for {
			d.WaitForListRateLimit(ctx)

			req := tke.NewDescribeClusterInstancesRequest()
			req.ClusterId = &clusterId
			req.Offset = &offset
			req.Limit = &pageSize
			if len(instanceIds) > 0 {
				req.InstanceIds = instanceIds
			}
			if instanceRole != "" {
				req.InstanceRole = &instanceRole
			}
			if len(filters) > 0 {
				req.Filters = filters
			}

			utils.LogRequest(ctx, "listTkeClusterInstances", req)
			resp, err := client.DescribeClusterInstancesWithContext(ctx, req)
			if err != nil {
				plugin.Logger(ctx).Error("listTkeClusterInstances", "request_error", err)
				return nil, err
			}
			if resp == nil || resp.Response == nil {
				break
			}

			instances := resp.Response.InstanceSet
			for _, inst := range instances {
				if inst == nil {
					continue
				}
				row := toTkeClusterInstanceRow(inst, clusterId)
				row.Region = region
				d.StreamListItem(ctx, row)
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}

			total := utils.PtrUint64(resp.Response.TotalCount)
			if utils.PageDone(offset, int64(len(instances)), total) {
				break
			}
			offset += int64(len(instances))
		}
	}

	return nil, nil
}

func buildTkeClusterInstanceFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*tke.Filter {
	return utils.BuildFilters(equalsQuals, newTkeFilter,
		utils.FilterMapping{QualName: "node_pool_id", FilterName: "nodepool-id"},
	)
}

// Row Type

type tkeClusterInstanceRow struct {
	InstanceId               string
	ClusterId                string
	InstanceRole             string
	InstanceState            string
	FailedReason             string
	DrainStatus              string
	LanIP                    string
	NodePoolId               string
	AutoscalingGroupId       string
	InstanceAdvancedSettings *tke.InstanceAdvancedSettings
	CreatedTime              *time.Time
	Region                   string
	Akas                     []string
}

func toTkeClusterInstanceRow(i *tke.Instance, clusterId string) tkeClusterInstanceRow {
	row := tkeClusterInstanceRow{
		InstanceId:               utils.PtrString(i.InstanceId),
		ClusterId:                clusterId,
		InstanceRole:             utils.PtrString(i.InstanceRole),
		InstanceState:            utils.PtrString(i.InstanceState),
		FailedReason:             utils.PtrString(i.FailedReason),
		DrainStatus:              utils.PtrString(i.DrainStatus),
		LanIP:                    utils.PtrString(i.LanIP),
		NodePoolId:               utils.PtrString(i.NodePoolId),
		AutoscalingGroupId:       utils.PtrString(i.AutoscalingGroupId),
		InstanceAdvancedSettings: i.InstanceAdvancedSettings,
		CreatedTime:              utils.ParseTimestamp(i.CreatedTime),
	}

	if row.InstanceId != "" && row.ClusterId != "" {
		row.Akas = []string{"tencentcloud:tke:cluster-instance:" + row.ClusterId + "/" + row.InstanceId}
	}
	return row
}
