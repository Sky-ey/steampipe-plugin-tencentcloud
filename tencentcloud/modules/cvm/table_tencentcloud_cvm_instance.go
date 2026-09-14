package cvm

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"strconv"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/cvm/v20170312"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudCvmInstance() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_cvm_instance",
		Description:       "Tencent Cloud Virtual Machine (CVM) instances provide scalable computing capacity in the cloud.",
		GetMatrixItemFunc: buildCvmRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "zone", Require: plugin.Optional},
				{Name: "project_id", Require: plugin.Optional},
				{Name: "vpc_id", Require: plugin.Optional},
				{Name: "subnet_id", Require: plugin.Optional},
				{Name: "uuid", Require: plugin.Optional},
				{Name: "instance_name", Require: plugin.Optional},
				{Name: "instance_charge_type", Require: plugin.Optional},
				{Name: "instance_state", Require: plugin.Optional},
				{Name: "private_ip_address", Require: plugin.Optional},
				{Name: "public_ip_address", Require: plugin.Optional},
			},
			Hydrate: listCvmInstances,
			Tags:    map[string]string{"service": "cvm", "action": "DescribeInstances"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("instance_id"),
			Hydrate:    getCvmInstance,
			Tags:       map[string]string{"service": "cvm", "action": "DescribeInstances"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "instance_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the CVM instance.", Transform: transform.FromField("InstanceId")},
			{Name: "cam_role_name", Type: proto.ColumnType_STRING, Description: "The CAM role name attached to the instance, if any.", Transform: transform.FromField("CamRoleName")},
			{Name: "cpu", Type: proto.ColumnType_INT, Description: "The number of CPU cores allocated to the instance.", Transform: transform.FromField("CPU")},
			{Name: "created_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the instance in ISO 8601 format (UTC).", Transform: transform.FromField("CreatedTime").NullIfZero()},
			{Name: "expired_time", Type: proto.ColumnType_TIMESTAMP, Description: "The expiry time of the instance in ISO 8601 format (UTC). Null for postpaid instances.", Transform: transform.FromField("ExpiredTime").NullIfZero()},
			{Name: "image_id", Type: proto.ColumnType_STRING, Description: "The ID of the image used to launch the instance.", Transform: transform.FromField("ImageId")},
			{Name: "instance_charge_type", Type: proto.ColumnType_STRING, Description: "The billing method of the instance (PREPAID, POSTPAID_BY_HOUR, CDHPAID, SPOTPAID, CDCPAID).", Transform: transform.FromField("InstanceChargeType")},
			{Name: "instance_name", Type: proto.ColumnType_STRING, Description: "The name of the CVM instance.", Transform: transform.FromField("InstanceName")},
			{Name: "instance_state", Type: proto.ColumnType_STRING, Description: "The current state of the instance (PENDING, RUNNING, STARTING, STOPPING, STOPPED, SHUTDOWN, TERMINATING).", Transform: transform.FromField("InstanceState")},
			{Name: "instance_type", Type: proto.ColumnType_STRING, Description: "The instance type (e.g., S5.MEDIUM4, S5.LARGE8).", Transform: transform.FromField("InstanceType")},
			{Name: "ipv6_addresses", Type: proto.ColumnType_JSON, Description: "The IPv6 addresses assigned to the instance.", Transform: transform.FromField("IPv6Addresses")},
			{Name: "latest_operation", Type: proto.ColumnType_STRING, Description: "The latest operation performed on the instance (e.g., StopInstances, ResetInstance).", Transform: transform.FromField("LatestOperation")},
			{Name: "latest_operation_state", Type: proto.ColumnType_STRING, Description: "The state of the latest operation (SUCCESS, OPERATING, FAILED).", Transform: transform.FromField("LatestOperationState")},
			{Name: "memory", Type: proto.ColumnType_INT, Description: "The memory size in GiB allocated to the instance.", Transform: transform.FromField("Memory")},
			{Name: "os_name", Type: proto.ColumnType_STRING, Description: "The operating system name of the instance.", Transform: transform.FromField("OsName")},
			{Name: "private_ip_address", Type: proto.ColumnType_IPADDR, Description: "The primary private IP address of the instance.", Transform: transform.FromField("PrivateIPAddress")},
			{Name: "project_id", Type: proto.ColumnType_INT, Description: "The project ID the instance belongs to. 0 means the default project.", Transform: transform.FromField("ProjectId")},
			{Name: "public_ip_address", Type: proto.ColumnType_IPADDR, Description: "The primary public IP address of the instance, if assigned.", Transform: transform.FromField("PublicIPAddress")},
			{Name: "renew_flag", Type: proto.ColumnType_STRING, Description: "Auto-renewal flag for prepaid instances (NOTIFY_AND_MANUAL_RENEW, NOTIFY_AND_AUTO_RENEW, DISABLE_NOTIFY_AND_MANUAL_RENEW).", Transform: transform.FromField("RenewFlag")},
			{Name: "restrict_state", Type: proto.ColumnType_STRING, Description: "The business state of the instance (NORMAL, EXPIRED, PROTECTIVELY_ISOLATED).", Transform: transform.FromField("RestrictState")},
			{Name: "security_group_ids", Type: proto.ColumnType_JSON, Description: "The IDs of the security groups associated with the instance.", Transform: transform.FromField("SecurityGroupIds")},
			{Name: "stop_charging_mode", Type: proto.ColumnType_STRING, Description: "The charging mode when the instance is stopped (KEEP_CHARGING, STOP_CHARGING, NOT_APPLICABLE).", Transform: transform.FromField("StopChargingMode")},
			{Name: "subnet_id", Type: proto.ColumnType_STRING, Description: "The ID of the subnet the instance belongs to.", Transform: transform.FromField("SubnetId")},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The tags attached to the instance as a map of key -> value.", Transform: transform.FromField("Tags")},
			{Name: "uuid", Type: proto.ColumnType_STRING, Description: "The globally unique ID of the instance.", Transform: transform.FromField("Uuid")},
			{Name: "vpc_id", Type: proto.ColumnType_STRING, Description: "The ID of the VPC the instance belongs to.", Transform: transform.FromField("VpcId")},
			{Name: "zone", Type: proto.ColumnType_STRING, Description: "The availability zone the instance is deployed in.", Transform: transform.FromField("Zone")},
			{Name: "system_disk", Type: proto.ColumnType_JSON, Description: "Information on the system disk of the instance.", Transform: transform.FromField("SystemDisk")},
			{Name: "system_disk_id", Type: proto.ColumnType_STRING, Description: "The ID of the system disk.", Transform: transform.FromField("SystemDiskId")},
			{Name: "system_disk_type", Type: proto.ColumnType_STRING, Description: "The type of the system disk (CLOUD_BASIC, CLOUD_PREMIUM, CLOUD_SSD, CLOUD_HSSD, CLOUD_BSSD, CLOUD_TSSD).", Transform: transform.FromField("SystemDiskType")},
			{Name: "system_disk_size", Type: proto.ColumnType_INT, Description: "The size of the system disk in GiB.", Transform: transform.FromField("SystemDiskSize")},
			{Name: "data_disks", Type: proto.ColumnType_JSON, Description: "Information of the instance data disks.", Transform: transform.FromField("DataDisks")},
			{Name: "internet_accessible", Type: proto.ColumnType_JSON, Description: "Information on instance bandwidth, including the network billing method and bandwidth cap.", Transform: transform.FromField("InternetAccessible")},
			{Name: "internet_charge_type", Type: proto.ColumnType_STRING, Description: "The network billing method of the instance (BANDWIDTH_PREPAID, TRAFFIC_POSTPAID_BY_HOUR, BANDWIDTH_POSTPAID_BY_HOUR, BANDWIDTH_PACKAGE).", Transform: transform.FromField("InternetChargeType")},
			{Name: "internet_max_bandwidth_out", Type: proto.ColumnType_INT, Description: "The maximum outbound bandwidth of the instance in Mbps.", Transform: transform.FromField("InternetMaxBandwidthOut")},
			{Name: "public_ip_assigned", Type: proto.ColumnType_BOOL, Description: "Whether a public IP address is assigned to the instance.", Transform: transform.FromField("PublicIpAssigned")},
			{Name: "bandwidth_package_id", Type: proto.ColumnType_STRING, Description: "The ID of the bandwidth package associated with the instance's public network.", Transform: transform.FromField("BandwidthPackageId")},
			{Name: "login_settings", Type: proto.ColumnType_JSON, Description: "Login settings of the instance, including the associated key pair IDs.", Transform: transform.FromField("LoginSettings")},
			{Name: "key_ids", Type: proto.ColumnType_JSON, Description: "The IDs of the key pairs associated with the instance.", Transform: transform.FromField("KeyIds")},
			{Name: "placement", Type: proto.ColumnType_JSON, Description: "Location of the instance, including zone, project ID and dedicated host information.", Transform: transform.FromField("Placement")},
			{Name: "host_id", Type: proto.ColumnType_STRING, Description: "The ID of the dedicated host (CDH) where the instance resides, if any.", Transform: transform.FromField("HostId")},
			{Name: "virtual_private_cloud", Type: proto.ColumnType_JSON, Description: "Information on the VPC where the instance resides, including VPC ID, subnet ID, public gateway flag and private IP addresses.", Transform: transform.FromField("VirtualPrivateCloud")},
			{Name: "private_ip_addresses", Type: proto.ColumnType_JSON, Description: "List of private IPs of the instance's primary ENI.", Transform: transform.FromField("PrivateIpAddresses")},
			{Name: "public_ip_addresses", Type: proto.ColumnType_JSON, Description: "List of public IPs of the instance's primary ENI.", Transform: transform.FromField("PublicIpAddresses")},
			{Name: "rdma_ip_addresses", Type: proto.ColumnType_JSON, Description: "IP list of the HPC cluster (RDMA).", Transform: transform.FromField("RdmaIpAddresses")},
			{Name: "public_ipv6_addresses", Type: proto.ColumnType_JSON, Description: "List of public IPv6 addresses bound to the instance.", Transform: transform.FromField("PublicIPv6Addresses")},
			{Name: "gpu_info", Type: proto.ColumnType_JSON, Description: "GPU information, returned only for GPU type instances.", Transform: transform.FromField("GPUInfo")},
			{Name: "gpu_count", Type: proto.ColumnType_DOUBLE, Description: "The number of GPUs of the instance. Only returned for GPU type instances.", Transform: transform.FromField("GPUCount")},
			{Name: "gpu_type", Type: proto.ColumnType_STRING, Description: "The GPU model of the instance (e.g., V100, T4). Only returned for GPU type instances.", Transform: transform.FromField("GPUType")},
			{Name: "gpu_ids", Type: proto.ColumnType_JSON, Description: "The IDs of the GPUs attached to the instance.", Transform: transform.FromField("GPUIds")},
			{Name: "metadata", Type: proto.ColumnType_JSON, Description: "Custom metadata specified when creating the CVM.", Transform: transform.FromField("Metadata")},
			{Name: "disaster_recover_group_id", Type: proto.ColumnType_STRING, Description: "Spread placement group ID.", Transform: transform.FromField("DisasterRecoverGroupId")},
			{Name: "hpc_cluster_id", Type: proto.ColumnType_STRING, Description: "High-performance computing cluster ID.", Transform: transform.FromField("HpcClusterId")},
			{Name: "dedicated_cluster_id", Type: proto.ColumnType_STRING, Description: "Dedicated cluster ID where the instance is located.", Transform: transform.FromField("DedicatedClusterId")},
			{Name: "isolated_source", Type: proto.ColumnType_STRING, Description: "Instance isolation type (ARREAR, EXPIRE, MANMADE, NOTISOLATED).", Transform: transform.FromField("IsolatedSource")},
			{Name: "license_type", Type: proto.ColumnType_STRING, Description: "Instance OS license type. Default value: TencentCloud.", Transform: transform.FromField("LicenseType")},
			{Name: "default_login_user", Type: proto.ColumnType_STRING, Description: "Default login user of the instance.", Transform: transform.FromField("DefaultLoginUser")},
			{Name: "default_login_port", Type: proto.ColumnType_INT, Description: "Default login port of the instance.", Transform: transform.FromField("DefaultLoginPort")},
			{Name: "disable_api_termination", Type: proto.ColumnType_BOOL, Description: "Instance destruction protection flag. true indicates deletion through APIs is not allowed.", Transform: transform.FromField("DisableApiTermination")},
			{Name: "latest_operation_request_id", Type: proto.ColumnType_STRING, Description: "Unique request ID for the last operation of the instance.", Transform: transform.FromField("LatestOperationRequestId")},
			{Name: "latest_operation_error_msg", Type: proto.ColumnType_STRING, Description: "Latest operation error message of the instance.", Transform: transform.FromField("LatestOperationErrorMsg")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the instance resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("InstanceName")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listCvmInstances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listCvmInstances", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &cvm.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listCvmInstances", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	pageSize := utils.PageSizeFromLimit(d, 100, 1)
	var offset int64 = 0

	filters := buildInstanceFilters(d.EqualsQuals)

	for {
		d.WaitForListRateLimit(ctx)

		req := cvm.NewDescribeInstancesRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		req.Filters = filters

		utils.LogRequest(ctx, "listCvmInstances", req)
		resp, err := client.DescribeInstancesWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listCvmInstances", "request_error", err)
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
			row := toCvmInstanceRow(inst)
			row.Region = region
			d.StreamListItem(ctx, row)

			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		total := utils.PtrInt64(resp.Response.TotalCount)
		if utils.PageDone(offset, int64(len(instances)), total) {
			break
		}
		offset += int64(len(instances))
	}

	return nil, nil
}

func buildInstanceFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*cvm.Filter {
	var filters []*cvm.Filter

	addFilter := func(qualName, filterName string) {
		if equalsQuals[qualName] == nil {
			return
		}
		value := equalsQuals[qualName].GetStringValue()
		if value == "" {
			return
		}
		filters = append(filters, &cvm.Filter{
			Name:   common.StringPtr(filterName),
			Values: common.StringPtrs([]string{value}),
		})
	}

	addFilter("zone", "zone")
	addFilter("vpc_id", "vpc-id")
	addFilter("subnet_id", "subnet-id")
	addFilter("uuid", "uuid")
	addFilter("instance_name", "instance-name")
	addFilter("instance_charge_type", "instance-charge-type")
	addFilter("instance_state", "instance-state")

	addInetFilter := func(qualName, filterName string) {
		if equalsQuals[qualName] == nil {
			return
		}
		inet := equalsQuals[qualName].GetInetValue()
		if inet == nil || inet.GetAddr() == "" {
			return
		}
		filters = append(filters, &cvm.Filter{
			Name:   common.StringPtr(filterName),
			Values: common.StringPtrs([]string{inet.GetAddr()}),
		})
	}
	addInetFilter("private_ip_address", "private-ip-address")
	addInetFilter("public_ip_address", "public-ip-address")

	if equalsQuals["project_id"] != nil {
		projectId := equalsQuals["project_id"].GetInt64Value()
		filters = append(filters, &cvm.Filter{
			Name:   common.StringPtr("project-id"),
			Values: common.StringPtrs([]string{strconv.FormatInt(projectId, 10)}),
		})
	}

	return filters
}

// Get Function

func getCvmInstance(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	instanceId := d.EqualsQuals["instance_id"].GetStringValue()
	plugin.Logger(ctx).Debug("getCvmInstance", "instance_id", instanceId, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &cvm.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getCvmInstance", "client_init_error", err)
		return nil, err
	}

	if instanceId == "" {
		return nil, nil
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	req := cvm.NewDescribeInstancesRequest()
	req.InstanceIds = common.StringPtrs([]string{instanceId})

	utils.LogRequest(ctx, "getCvmInstance", req)
	resp, err := client.DescribeInstancesWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getCvmInstance", "request_error", err)
		return nil, err
	}

	if resp == nil || resp.Response == nil || len(resp.Response.InstanceSet) == 0 {
		return nil, nil
	}
	row := toCvmInstanceRow(resp.Response.InstanceSet[0])
	row.Region = region
	return row, nil
}

// Row Type

type cvmInstanceRow struct {
	InstanceId               string
	InstanceName             string
	InstanceState            string
	InstanceType             string
	CPU                      int64
	Memory                   int64
	InstanceChargeType       string
	RestrictState            string
	Zone                     string
	ProjectId                int64
	ImageId                  string
	OsName                   string
	VpcId                    string
	SubnetId                 string
	PrivateIPAddress         string
	PublicIPAddress          string
	IPv6Addresses            []string
	SecurityGroupIds         []string
	RenewFlag                string
	CreatedTime              *time.Time
	ExpiredTime              *time.Time
	StopChargingMode         string
	Uuid                     string
	CamRoleName              string
	LatestOperation          string
	LatestOperationState     string
	LatestOperationRequestId string
	DisasterRecoverGroupId   string
	HpcClusterId             string
	DedicatedClusterId       string
	IsolatedSource           string
	LicenseType              string
	DefaultLoginUser         string
	DefaultLoginPort         int64
	DisableApiTermination    bool
	LatestOperationErrorMsg  string
	InternetChargeType       string
	InternetMaxBandwidthOut  int64
	PublicIpAssigned         bool
	BandwidthPackageId       string
	SystemDiskId             string
	SystemDiskType           string
	SystemDiskSize           int64
	KeyIds                   []string
	GPUCount                 float64
	GPUType                  string
	GPUIds                   []string
	HostId                   string
	SystemDisk               *cvm.SystemDisk
	DataDisks                []*cvm.DataDisk
	InternetAccessible       *cvm.InternetAccessible
	LoginSettings            *cvm.LoginSettings
	Placement                *cvm.Placement
	VirtualPrivateCloud      *cvm.VirtualPrivateCloud
	PrivateIpAddresses       []string
	PublicIpAddresses        []string
	RdmaIpAddresses          []string
	PublicIPv6Addresses      []string
	GPUInfo                  *cvm.GPUInfo
	Metadata                 *cvm.Metadata
	Tags                     map[string]string
	Region                   string
	Akas                     []string
}

func toCvmInstanceRow(inst *cvm.Instance) cvmInstanceRow {
	row := cvmInstanceRow{
		InstanceId:               utils.PtrString(inst.InstanceId),
		InstanceName:             utils.PtrString(inst.InstanceName),
		InstanceState:            utils.PtrString(inst.InstanceState),
		InstanceType:             utils.PtrString(inst.InstanceType),
		CPU:                      utils.PtrInt64(inst.CPU),
		Memory:                   utils.PtrInt64(inst.Memory),
		InstanceChargeType:       utils.PtrString(inst.InstanceChargeType),
		RestrictState:            utils.PtrString(inst.RestrictState),
		ImageId:                  utils.PtrString(inst.ImageId),
		OsName:                   utils.PtrString(inst.OsName),
		RenewFlag:                utils.PtrString(inst.RenewFlag),
		CreatedTime:              utils.ParseTimestamp(inst.CreatedTime),
		ExpiredTime:              utils.ParseTimestamp(inst.ExpiredTime),
		StopChargingMode:         utils.PtrString(inst.StopChargingMode),
		Uuid:                     utils.PtrString(inst.Uuid),
		CamRoleName:              utils.PtrString(inst.CamRoleName),
		LatestOperation:          utils.PtrString(inst.LatestOperation),
		LatestOperationState:     utils.PtrString(inst.LatestOperationState),
		LatestOperationRequestId: utils.PtrString(inst.LatestOperationRequestId),
		DisasterRecoverGroupId:   utils.PtrString(inst.DisasterRecoverGroupId),
		HpcClusterId:             utils.PtrString(inst.HpcClusterId),
		DedicatedClusterId:       utils.PtrString(inst.DedicatedClusterId),
		IsolatedSource:           utils.PtrString(inst.IsolatedSource),
		LicenseType:              utils.PtrString(inst.LicenseType),
		DefaultLoginUser:         utils.PtrString(inst.DefaultLoginUser),
		DefaultLoginPort:         utils.PtrInt64(inst.DefaultLoginPort),
		DisableApiTermination:    utils.PtrBool(inst.DisableApiTermination),
		LatestOperationErrorMsg:  utils.PtrString(inst.LatestOperationErrorMsg),
		IPv6Addresses:            utils.PtrStringSlice(inst.IPv6Addresses),
		SecurityGroupIds:         utils.PtrStringSlice(inst.SecurityGroupIds),
		PrivateIpAddresses:       utils.PtrStringSlice(inst.PrivateIpAddresses),
		PublicIpAddresses:        utils.PtrStringSlice(inst.PublicIpAddresses),
		RdmaIpAddresses:          utils.PtrStringSlice(inst.RdmaIpAddresses),
		PublicIPv6Addresses:      utils.PtrStringSlice(inst.PublicIPv6Addresses),
		SystemDisk:               inst.SystemDisk,
		DataDisks:                inst.DataDisks,
		InternetAccessible:       inst.InternetAccessible,
		LoginSettings:            inst.LoginSettings,
		Placement:                inst.Placement,
		VirtualPrivateCloud:      inst.VirtualPrivateCloud,
		GPUInfo:                  inst.GPUInfo,
		Metadata:                 inst.Metadata,
		Tags:                     TagsToMap(inst.Tags),
	}

	if inst.Placement != nil {
		row.Zone = utils.PtrString(inst.Placement.Zone)
		row.ProjectId = utils.PtrInt64(inst.Placement.ProjectId)
		row.HostId = utils.PtrString(inst.Placement.HostId)
	}
	if inst.InternetAccessible != nil {
		row.InternetChargeType = utils.PtrString(inst.InternetAccessible.InternetChargeType)
		row.InternetMaxBandwidthOut = utils.PtrInt64(inst.InternetAccessible.InternetMaxBandwidthOut)
		row.PublicIpAssigned = utils.PtrBool(inst.InternetAccessible.PublicIpAssigned)
		row.BandwidthPackageId = utils.PtrString(inst.InternetAccessible.BandwidthPackageId)
	}
	if inst.SystemDisk != nil {
		row.SystemDiskId = utils.PtrString(inst.SystemDisk.DiskId)
		row.SystemDiskType = utils.PtrString(inst.SystemDisk.DiskType)
		row.SystemDiskSize = utils.PtrInt64(inst.SystemDisk.DiskSize)
	}
	if inst.LoginSettings != nil {
		row.KeyIds = utils.PtrStringSlice(inst.LoginSettings.KeyIds)
	}
	if inst.GPUInfo != nil {
		row.GPUCount = utils.PtrFloat64(inst.GPUInfo.GPUCount)
		row.GPUType = utils.PtrString(inst.GPUInfo.GPUType)
		row.GPUIds = utils.PtrStringSlice(inst.GPUInfo.GPUId)
	}
	if inst.VirtualPrivateCloud != nil {
		row.VpcId = utils.PtrString(inst.VirtualPrivateCloud.VpcId)
		row.SubnetId = utils.PtrString(inst.VirtualPrivateCloud.SubnetId)
	}
	if len(inst.PrivateIpAddresses) > 0 {
		row.PrivateIPAddress = utils.PtrString(inst.PrivateIpAddresses[0])
	}
	if len(inst.PublicIpAddresses) > 0 {
		row.PublicIPAddress = utils.PtrString(inst.PublicIpAddresses[0])
	}

	if row.InstanceId != "" {
		row.Akas = []string{"tencentcloud:cvm:instance:" + row.InstanceId}
	}
	return row
}
