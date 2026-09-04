package tke

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	tke "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/tke/v20180525"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudTkeCluster() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_tke_cluster",
		Description:       "Tencent Kubernetes Engine (TKE) clusters are the top-level managed Kubernetes resource in Tencent Cloud.",
		GetMatrixItemFunc: buildTkeRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "cluster_id", Require: plugin.Optional},
				{Name: "cluster_name", Require: plugin.Optional},
				{Name: "cluster_type", Require: plugin.Optional},
				{Name: "cluster_status", Require: plugin.Optional},
				{Name: "vpc_id", Require: plugin.Optional},
			},
			Hydrate: listTkeClusters,
			Tags:    map[string]string{"service": "tke", "action": "DescribeClusters"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("cluster_id"),
			Hydrate:    getTkeCluster,
			Tags:       map[string]string{"service": "tke", "action": "DescribeClusters"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "cluster_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the TKE cluster.", Transform: transform.FromField("ClusterId")},
			{Name: "cluster_name", Type: proto.ColumnType_STRING, Description: "The name of the cluster.", Transform: transform.FromField("ClusterName")},
			{Name: "cluster_type", Type: proto.ColumnType_STRING, Description: "The cluster type (MANAGED_CLUSTER, INDEPENDENT_CLUSTER).", Transform: transform.FromField("ClusterType")},
			{Name: "cluster_status", Type: proto.ColumnType_STRING, Description: "The cluster status (Trading, Creating, Running, Deleting, Idling, Recovering, Scaling, Upgrading, WaittingForConnect, Pause, NodeUpgrading, RuntimeUpgrading, MasterScaling, ClusterLevelUpgrading, ResourceIsolate, ResourceIsolated, ResourceReverse, Abnormal).", Transform: transform.FromField("ClusterStatus")},
			{Name: "cluster_description", Type: proto.ColumnType_STRING, Description: "The description of the cluster.", Transform: transform.FromField("ClusterDescription")},
			{Name: "cluster_version", Type: proto.ColumnType_STRING, Description: "The Kubernetes version of the cluster (e.g., 1.10.5).", Transform: transform.FromField("ClusterVersion")},
			{Name: "cluster_os", Type: proto.ColumnType_STRING, Description: "The operating system of the cluster nodes.", Transform: transform.FromField("ClusterOs")},
			{Name: "cluster_node_num", Type: proto.ColumnType_INT, Description: "The current total number of nodes in the cluster.", Transform: transform.FromField("ClusterNodeNum")},
			{Name: "cluster_master_node_num", Type: proto.ColumnType_INT, Description: "The current number of master nodes in the cluster.", Transform: transform.FromField("ClusterMaterNodeNum")},
			{Name: "cluster_etcd_node_num", Type: proto.ColumnType_INT, Description: "The current number of etcd nodes in the cluster.", Transform: transform.FromField("ClusterEtcdNodeNum")},
			{Name: "project_id", Type: proto.ColumnType_INT, Description: "The ID of the project the cluster belongs to.", Transform: transform.FromField("ProjectId")},
			{Name: "image_id", Type: proto.ColumnType_STRING, Description: "The ID of the image used by the cluster.", Transform: transform.FromField("ImageId")},
			{Name: "os_customize_type", Type: proto.ColumnType_STRING, Description: "The container image tag (DOCKER_CUSTOMIZE or GENERAL).", Transform: transform.FromField("OsCustomizeType")},
			{Name: "container_runtime", Type: proto.ColumnType_STRING, Description: "The container runtime of the cluster (docker or containerd).", Transform: transform.FromField("ContainerRuntime")},
			{Name: "runtime_version", Type: proto.ColumnType_STRING, Description: "The runtime version of the cluster.", Transform: transform.FromField("RuntimeVersion")},
			{Name: "cluster_level", Type: proto.ColumnType_STRING, Description: "The cluster specification level, valid for managed clusters.", Transform: transform.FromField("ClusterLevel")},
			{Name: "auto_upgrade_cluster_level", Type: proto.ColumnType_BOOL, Description: "Whether to auto-upgrade the cluster level.", Transform: transform.FromField("AutoUpgradeClusterLevel")},
			{Name: "qgpu_share_enable", Type: proto.ColumnType_BOOL, Description: "Whether qGPU sharing is enabled.", Transform: transform.FromField("QGPUShareEnable")},
			{Name: "deletion_protection", Type: proto.ColumnType_BOOL, Description: "Whether deletion protection is enabled for the cluster.", Transform: transform.FromField("DeletionProtection")},
			{Name: "enable_external_node", Type: proto.ColumnType_BOOL, Description: "Whether the cluster supports external (registered) nodes.", Transform: transform.FromField("EnableExternalNode")},
			{Name: "cdc_id", Type: proto.ColumnType_STRING, Description: "The CDC ID the cluster belongs to, if any.", Transform: transform.FromField("CdcId")},
			{Name: "vpc_id", Type: proto.ColumnType_STRING, Description: "The VPC ID the cluster resides in.", Transform: transform.FromField("VpcId")},
			{Name: "cluster_cidr", Type: proto.ColumnType_STRING, Description: "The CIDR used to assign container and service IPs.", Transform: transform.FromField("ClusterCIDR")},
			{Name: "service_cidr", Type: proto.ColumnType_STRING, Description: "The IP range for service assignment.", Transform: transform.FromField("ServiceCIDR")},
			{Name: "max_node_pod_num", Type: proto.ColumnType_INT, Description: "The maximum number of pods per node.", Transform: transform.FromField("MaxNodePodNum")},
			{Name: "max_cluster_service_num", Type: proto.ColumnType_INT, Description: "The maximum number of services in the cluster.", Transform: transform.FromField("MaxClusterServiceNum")},
			{Name: "cilium_mode", Type: proto.ColumnType_STRING, Description: "The Cilium mode configuration (e.g., clusterIP).", Transform: transform.FromField("CiliumMode")},
			{Name: "kube_proxy_mode", Type: proto.ColumnType_STRING, Description: "The kube-proxy network mode, applicable to ipvs+bpf mode.", Transform: transform.FromField("KubeProxyMode")},
			{Name: "is_dual_stack", Type: proto.ColumnType_BOOL, Description: "Whether the cluster is a dual-stack cluster in VPC-CNI mode.", Transform: transform.FromField("IsDualStack")},
			{Name: "network_settings", Type: proto.ColumnType_JSON, Description: "The full cluster network settings.", Transform: transform.FromField("ClusterNetworkSettings")},
			{Name: "property", Type: proto.ColumnType_JSON, Description: "Cluster attributes map.", Transform: transform.FromField("Property")},
			{Name: "created_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the cluster.", Transform: transform.FromField("CreatedTime").NullIfZero()},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The tags attached to the cluster as a map of key -> value.", Transform: transform.FromField("Tags")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the cluster resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("ClusterName")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listTkeClusters(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listTkeClusters", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &tke.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listTkeClusters", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	clusterIds := utils.SingleIdFromQual(d.EqualsQuals, "cluster_id")
	filters := buildTkeClusterFilters(d.EqualsQuals)
	clusterType := d.EqualsQualString("cluster_type")

	pageSize := utils.PageSizeFromLimit(d, 100, 1)
	var offset int64

	for {
		d.WaitForListRateLimit(ctx)

		req := tke.NewDescribeClustersRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		if len(clusterIds) > 0 {
			req.ClusterIds = clusterIds
		} else if len(filters) > 0 {
			req.Filters = filters
		}
		if clusterType != "" {
			req.ClusterType = &clusterType
		}

		utils.LogRequest(ctx, "listTkeClusters", req)
		resp, err := client.DescribeClustersWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listTkeClusters", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		clusters := resp.Response.Clusters
		for _, c := range clusters {
			if c == nil {
				continue
			}
			row := toTkeClusterRow(c)
			row.Region = region
			d.StreamListItem(ctx, row)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		total := utils.PtrInt64(resp.Response.TotalCount)
		if utils.PageDone(offset, int64(len(clusters)), total) {
			break
		}
		offset += int64(len(clusters))
	}

	return nil, nil
}

func buildTkeClusterFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*tke.Filter {
	return utils.BuildFilters(equalsQuals, newTkeFilter,
		utils.FilterMapping{QualName: "cluster_name", FilterName: "ClusterName"},
		utils.FilterMapping{QualName: "cluster_status", FilterName: "ClusterStatus"},
		utils.FilterMapping{QualName: "vpc_id", FilterName: "vpc-id"},
	)
}

// Get Function

func getTkeCluster(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	clusterId := d.EqualsQuals["cluster_id"].GetStringValue()
	plugin.Logger(ctx).Info("getTkeCluster", "cluster_id", clusterId, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &tke.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getTkeCluster", "client_init_error", err)
		return nil, err
	}

	if clusterId == "" {
		return nil, nil
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	req := tke.NewDescribeClustersRequest()
	req.ClusterIds = common.StringPtrs([]string{clusterId})

	utils.LogRequest(ctx, "getTkeCluster", req)
	resp, err := client.DescribeClustersWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getTkeCluster", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.Clusters) == 0 {
		return nil, nil
	}

	row := toTkeClusterRow(resp.Response.Clusters[0])
	row.Region = region
	return row, nil
}

// Row Type

type tkeClusterRow struct {
	ClusterId               string
	ClusterName             string
	ClusterType             string
	ClusterStatus           string
	ClusterDescription      string
	ClusterVersion          string
	ClusterOs               string
	ClusterNodeNum          uint64
	ClusterMaterNodeNum     uint64
	ClusterEtcdNodeNum      uint64
	ProjectId               uint64
	ImageId                 string
	OsCustomizeType         string
	ContainerRuntime        string
	RuntimeVersion          string
	ClusterLevel            string
	AutoUpgradeClusterLevel bool
	QGPUShareEnable         bool
	DeletionProtection      bool
	EnableExternalNode      bool
	CdcId                   string
	VpcId                   string
	ClusterCIDR             string
	ServiceCIDR             string
	MaxNodePodNum           uint64
	MaxClusterServiceNum    uint64
	CiliumMode              string
	KubeProxyMode           string
	IsDualStack             bool
	ClusterNetworkSettings  *tke.ClusterNetworkSettings
	Property                *string
	CreatedTime             *time.Time
	Tags                    map[string]string
	Region                  string
	Akas                    []string
}

func toTkeClusterRow(c *tke.Cluster) tkeClusterRow {
	row := tkeClusterRow{
		ClusterId:               utils.PtrString(c.ClusterId),
		ClusterName:             utils.PtrString(c.ClusterName),
		ClusterType:             utils.PtrString(c.ClusterType),
		ClusterStatus:           utils.PtrString(c.ClusterStatus),
		ClusterDescription:      utils.PtrString(c.ClusterDescription),
		ClusterVersion:          utils.PtrString(c.ClusterVersion),
		ClusterOs:               utils.PtrString(c.ClusterOs),
		ClusterNodeNum:          utils.PtrUint64(c.ClusterNodeNum),
		ClusterMaterNodeNum:     utils.PtrUint64(c.ClusterMaterNodeNum),
		ClusterEtcdNodeNum:      utils.PtrUint64(c.ClusterEtcdNodeNum),
		ProjectId:               utils.PtrUint64(c.ProjectId),
		ImageId:                 utils.PtrString(c.ImageId),
		OsCustomizeType:         utils.PtrString(c.OsCustomizeType),
		ContainerRuntime:        utils.PtrString(c.ContainerRuntime),
		RuntimeVersion:          utils.PtrString(c.RuntimeVersion),
		ClusterLevel:            utils.PtrString(c.ClusterLevel),
		AutoUpgradeClusterLevel: utils.PtrBool(c.AutoUpgradeClusterLevel),
		QGPUShareEnable:         utils.PtrBool(c.QGPUShareEnable),
		DeletionProtection:      utils.PtrBool(c.DeletionProtection),
		EnableExternalNode:      utils.PtrBool(c.EnableExternalNode),
		CdcId:                   utils.PtrString(c.CdcId),
		ClusterNetworkSettings:  c.ClusterNetworkSettings,
		Property:                c.Property,
		CreatedTime:             utils.ParseTimestamp(c.CreatedTime),
		Tags:                    TagSpecificationToMap(c.TagSpecification),
	}

	if c.ClusterNetworkSettings != nil {
		row.VpcId = utils.PtrString(c.ClusterNetworkSettings.VpcId)
		row.ClusterCIDR = utils.PtrString(c.ClusterNetworkSettings.ClusterCIDR)
		row.ServiceCIDR = utils.PtrString(c.ClusterNetworkSettings.ServiceCIDR)
		row.MaxNodePodNum = utils.PtrUint64(c.ClusterNetworkSettings.MaxNodePodNum)
		row.MaxClusterServiceNum = utils.PtrUint64(c.ClusterNetworkSettings.MaxClusterServiceNum)
		row.CiliumMode = utils.PtrString(c.ClusterNetworkSettings.CiliumMode)
		row.KubeProxyMode = utils.PtrString(c.ClusterNetworkSettings.KubeProxyMode)
		row.IsDualStack = utils.PtrBool(c.ClusterNetworkSettings.IsDualStack)
	}

	if row.ClusterId != "" {
		row.Akas = []string{"tencentcloud:tke:cluster:" + row.ClusterId}
	}
	return row
}
