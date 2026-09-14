package cvm

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"time"

	as "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/as/v20180419"
	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudCvmAutoScalingGroup() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_cvm_auto_scaling_group",
		Description:       "Tencent Cloud Auto Scaling (AS) groups automatically adjust the number of CVM instances according to configured scaling policies.",
		GetMatrixItemFunc: buildCvmRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "auto_scaling_group_id", Require: plugin.Optional},
				{Name: "auto_scaling_group_name", Require: plugin.Optional},
				{Name: "launch_configuration_id", Require: plugin.Optional},
			},
			Hydrate: listCvmAutoScalingGroups,
			Tags:    map[string]string{"service": "as", "action": "DescribeAutoScalingGroups"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("auto_scaling_group_id"),
			Hydrate:    getCvmAutoScalingGroup,
			Tags:       map[string]string{"service": "as", "action": "DescribeAutoScalingGroups"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "auto_scaling_group_id", Type: proto.ColumnType_STRING, Description: "The ID of the auto scaling group.", Transform: transform.FromField("AutoScalingGroupId")},
			{Name: "auto_scaling_group_name", Type: proto.ColumnType_STRING, Description: "The name of the auto scaling group.", Transform: transform.FromField("AutoScalingGroupName")},
			{Name: "auto_scaling_group_status", Type: proto.ColumnType_STRING, Description: "The current status of the scaling group (NORMAL, CVM_ABNORMAL, LB_ABNORMAL, LB_LISTENER_ABNORMAL, LB_LOCATION_ABNORMAL, VPC_ABNORMAL, SUBNET_ABNORMAL, INSUFFICIENT_BALANCE, LB_BACKEND_REGION_NOT_MATCH, LB_BACKEND_VPC_NOT_MATCH).", Transform: transform.FromField("AutoScalingGroupStatus")},
			{Name: "enabled_status", Type: proto.ColumnType_STRING, Description: "The enabled status of the scaling group (ENABLED, DISABLED).", Transform: transform.FromField("EnabledStatus")},
			{Name: "launch_configuration_id", Type: proto.ColumnType_STRING, Description: "The ID of the launch configuration.", Transform: transform.FromField("LaunchConfigurationId")},
			{Name: "launch_configuration_name", Type: proto.ColumnType_STRING, Description: "The name of the launch configuration.", Transform: transform.FromField("LaunchConfigurationName")},
			{Name: "vpc_id", Type: proto.ColumnType_STRING, Description: "The VPC ID of the scaling group.", Transform: transform.FromField("VpcId")},
			{Name: "project_id", Type: proto.ColumnType_INT, Description: "The project ID the scaling group belongs to.", Transform: transform.FromField("ProjectId")},
			{Name: "desired_capacity", Type: proto.ColumnType_INT, Description: "The desired number of instances.", Transform: transform.FromField("DesiredCapacity")},
			{Name: "min_size", Type: proto.ColumnType_INT, Description: "The minimum number of instances.", Transform: transform.FromField("MinSize")},
			{Name: "max_size", Type: proto.ColumnType_INT, Description: "The maximum number of instances.", Transform: transform.FromField("MaxSize")},
			{Name: "instance_count", Type: proto.ColumnType_INT, Description: "The total number of instances in the scaling group.", Transform: transform.FromField("InstanceCount")},
			{Name: "in_service_instance_count", Type: proto.ColumnType_INT, Description: "The number of instances in IN_SERVICE status.", Transform: transform.FromField("InServiceInstanceCount")},
			{Name: "default_cooldown", Type: proto.ColumnType_INT, Description: "The default cooldown period in seconds.", Transform: transform.FromField("DefaultCooldown")},
			{Name: "in_activity_status", Type: proto.ColumnType_STRING, Description: "Whether the group is performing a scaling activity (IN_ACTIVITY, NOT_IN_ACTIVITY).", Transform: transform.FromField("InActivityStatus")},
			{Name: "retry_policy", Type: proto.ColumnType_STRING, Description: "The retry policy (IMMEDIATE_RETRY, INCREMENTAL_INTERVALS, NO_RETRY).", Transform: transform.FromField("RetryPolicy")},
			{Name: "multi_zone_subnet_policy", Type: proto.ColumnType_STRING, Description: "The multi-AZ/subnet policy (PRIORITY, EQUALITY).", Transform: transform.FromField("MultiZoneSubnetPolicy")},
			{Name: "health_check_type", Type: proto.ColumnType_STRING, Description: "The instance health check type (CVM, CLB).", Transform: transform.FromField("HealthCheckType")},
			{Name: "instance_allocation_policy", Type: proto.ColumnType_STRING, Description: "The instance allocation policy (LAUNCH_CONFIGURATION, SPOT_MIXED).", Transform: transform.FromField("InstanceAllocationPolicy")},
			{Name: "ipv6_address_count", Type: proto.ColumnType_INT, Description: "The number of IPv6 addresses an instance has.", Transform: transform.FromField("Ipv6AddressCount")},
			{Name: "load_balancer_health_check_grace_period", Type: proto.ColumnType_INT, Description: "The grace period in seconds of the CLB health check.", Transform: transform.FromField("LoadBalancerHealthCheckGracePeriod")},
			{Name: "capacity_rebalance", Type: proto.ColumnType_BOOL, Description: "Whether capacity rebalancing is enabled for spot instances.", Transform: transform.FromField("CapacityRebalance")},
			{Name: "concurrent_scale_out_for_desired_capacity", Type: proto.ColumnType_BOOL, Description: "Whether concurrent scale-out to reach desired capacity is enabled.", Transform: transform.FromField("ConcurrentScaleOutForDesiredCapacity")},
			{Name: "created_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the scaling group in UTC format.", Transform: transform.FromField("CreatedTime").NullIfZero()},
			{Name: "zone_set", Type: proto.ColumnType_JSON, Description: "The list of availability zones of the scaling group.", Transform: transform.FromField("ZoneSet")},
			{Name: "subnet_id_set", Type: proto.ColumnType_JSON, Description: "The list of subnet IDs of the scaling group.", Transform: transform.FromField("SubnetIdSet")},
			{Name: "load_balancer_id_set", Type: proto.ColumnType_JSON, Description: "The list of classic load balancer IDs.", Transform: transform.FromField("LoadBalancerIdSet")},
			{Name: "termination_policy_set", Type: proto.ColumnType_JSON, Description: "The termination policies (OLDEST_INSTANCE, NEWEST_INSTANCE).", Transform: transform.FromField("TerminationPolicySet")},
			{Name: "forward_load_balancer_set", Type: proto.ColumnType_JSON, Description: "The list of application load balancers.", Transform: transform.FromField("ForwardLoadBalancerSet")},
			{Name: "service_settings", Type: proto.ColumnType_JSON, Description: "The service settings of the scaling group.", Transform: transform.FromField("ServiceSettings")},
			{Name: "spot_mixed_allocation_policy", Type: proto.ColumnType_JSON, Description: "The spot mixed allocation policy, returned only when InstanceAllocationPolicy is SPOT_MIXED.", Transform: transform.FromField("SpotMixedAllocationPolicy")},
			{Name: "instance_name_index_settings", Type: proto.ColumnType_JSON, Description: "The instance name index settings.", Transform: transform.FromField("InstanceNameIndexSettings")},
			{Name: "host_name_index_settings", Type: proto.ColumnType_JSON, Description: "The instance host name index settings.", Transform: transform.FromField("HostNameIndexSettings")},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The tags attached to the scaling group as a map of key -> value.", Transform: transform.FromField("Tags")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the scaling group resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("AutoScalingGroupName")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listCvmAutoScalingGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listCvmAutoScalingGroups", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &as.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listCvmAutoScalingGroups", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	// Mode B: ID-first. AutoScalingGroupIds and Filters are mutually exclusive.
	groupIds := autoScalingGroupIdsFromQuals(d.EqualsQuals)
	filters := buildAutoScalingGroupFilters(d.EqualsQuals)

	pageSize := utils.PageSizeFromLimit(d, 100, 1)
	limit := uint64(pageSize)
	var offset uint64 = 0

	for {
		d.WaitForListRateLimit(ctx)

		req := as.NewDescribeAutoScalingGroupsRequest()
		req.Limit = &limit
		req.Offset = &offset

		if len(groupIds) > 0 {
			req.AutoScalingGroupIds = common.StringPtrs(groupIds)
		} else if len(filters) > 0 {
			req.Filters = filters
		}

		utils.LogRequest(ctx, "listCvmAutoScalingGroups", req)
		resp, err := client.DescribeAutoScalingGroupsWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listCvmAutoScalingGroups", "request_error", err)
			return nil, err
		}

		if resp == nil || resp.Response == nil {
			break
		}

		groups := resp.Response.AutoScalingGroupSet
		for _, g := range groups {
			if g == nil {
				continue
			}
			row := toCvmAutoScalingGroupRow(g)
			row.Region = region
			d.StreamListItem(ctx, row)

			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		var total int64
		if resp.Response.TotalCount != nil {
			total = int64(*resp.Response.TotalCount)
		}
		if utils.PageDone(int64(offset), int64(len(groups)), total) {
			break
		}
		offset += uint64(len(groups))
	}

	return nil, nil
}

func autoScalingGroupIdsFromQuals(equalsQuals plugin.KeyColumnEqualsQualMap) []string {
	if equalsQuals["auto_scaling_group_id"] == nil {
		return nil
	}
	id := equalsQuals["auto_scaling_group_id"].GetStringValue()
	if id == "" {
		return nil
	}
	return []string{id}
}

func buildAutoScalingGroupFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*as.Filter {
	var filters []*as.Filter

	addFilter := func(qualName, filterName string) {
		if equalsQuals[qualName] == nil {
			return
		}
		value := equalsQuals[qualName].GetStringValue()
		if value == "" {
			return
		}
		filters = append(filters, &as.Filter{
			Name:   common.StringPtr(filterName),
			Values: common.StringPtrs([]string{value}),
		})
	}

	addFilter("auto_scaling_group_name", "auto-scaling-group-name")
	addFilter("launch_configuration_id", "launch-configuration-id")

	return filters
}

// Get Function

func getCvmAutoScalingGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	groupId := d.EqualsQuals["auto_scaling_group_id"].GetStringValue()
	plugin.Logger(ctx).Debug("getCvmAutoScalingGroup", "auto_scaling_group_id", groupId, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &as.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getCvmAutoScalingGroup", "client_init_error", err)
		return nil, err
	}

	if groupId == "" {
		return nil, nil
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	req := as.NewDescribeAutoScalingGroupsRequest()
	req.AutoScalingGroupIds = common.StringPtrs([]string{groupId})

	utils.LogRequest(ctx, "getCvmAutoScalingGroup", req)
	resp, err := client.DescribeAutoScalingGroupsWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getCvmAutoScalingGroup", "request_error", err)
		return nil, err
	}

	if resp == nil || resp.Response == nil || len(resp.Response.AutoScalingGroupSet) == 0 {
		return nil, nil
	}

	row := toCvmAutoScalingGroupRow(resp.Response.AutoScalingGroupSet[0])
	row.Region = region
	return row, nil
}

// Row Type

type cvmAutoScalingGroupRow struct {
	AutoScalingGroupId                   string
	AutoScalingGroupName                 string
	AutoScalingGroupStatus               string
	EnabledStatus                        string
	LaunchConfigurationId                string
	LaunchConfigurationName              string
	VpcId                                string
	ProjectId                            int64
	DesiredCapacity                      int64
	MinSize                              int64
	MaxSize                              int64
	InstanceCount                        int64
	InServiceInstanceCount               int64
	DefaultCooldown                      int64
	InActivityStatus                     string
	RetryPolicy                          string
	MultiZoneSubnetPolicy                string
	HealthCheckType                      string
	InstanceAllocationPolicy             string
	Ipv6AddressCount                     int64
	LoadBalancerHealthCheckGracePeriod   int64
	CapacityRebalance                    bool
	ConcurrentScaleOutForDesiredCapacity bool
	CreatedTime                          *time.Time
	ZoneSet                              []string
	SubnetIdSet                          []string
	LoadBalancerIdSet                    []string
	TerminationPolicySet                 []string
	ForwardLoadBalancerSet               []*as.ForwardLoadBalancer
	ServiceSettings                      *as.ServiceSettings
	SpotMixedAllocationPolicy            *as.SpotMixedAllocationPolicy
	InstanceNameIndexSettings            *as.InstanceNameIndexSettings
	HostNameIndexSettings                *as.HostNameIndexSettings
	Tags                                 map[string]string
	Region                               string
	Akas                                 []string
}

func toCvmAutoScalingGroupRow(g *as.AutoScalingGroup) cvmAutoScalingGroupRow {
	row := cvmAutoScalingGroupRow{
		AutoScalingGroupId:                   utils.PtrString(g.AutoScalingGroupId),
		AutoScalingGroupName:                 utils.PtrString(g.AutoScalingGroupName),
		AutoScalingGroupStatus:               utils.PtrString(g.AutoScalingGroupStatus),
		EnabledStatus:                        utils.PtrString(g.EnabledStatus),
		LaunchConfigurationId:                utils.PtrString(g.LaunchConfigurationId),
		LaunchConfigurationName:              utils.PtrString(g.LaunchConfigurationName),
		VpcId:                                utils.PtrString(g.VpcId),
		ProjectId:                            utils.PtrInt64(g.ProjectId),
		DesiredCapacity:                      utils.PtrInt64(g.DesiredCapacity),
		MinSize:                              utils.PtrInt64(g.MinSize),
		MaxSize:                              utils.PtrInt64(g.MaxSize),
		InstanceCount:                        utils.PtrInt64(g.InstanceCount),
		InServiceInstanceCount:               utils.PtrInt64(g.InServiceInstanceCount),
		DefaultCooldown:                      utils.PtrInt64(g.DefaultCooldown),
		InActivityStatus:                     utils.PtrString(g.InActivityStatus),
		RetryPolicy:                          utils.PtrString(g.RetryPolicy),
		MultiZoneSubnetPolicy:                utils.PtrString(g.MultiZoneSubnetPolicy),
		HealthCheckType:                      utils.PtrString(g.HealthCheckType),
		InstanceAllocationPolicy:             utils.PtrString(g.InstanceAllocationPolicy),
		Ipv6AddressCount:                     utils.PtrInt64(g.Ipv6AddressCount),
		LoadBalancerHealthCheckGracePeriod:   int64(utils.PtrUint64(g.LoadBalancerHealthCheckGracePeriod)),
		CapacityRebalance:                    utils.PtrBool(g.CapacityRebalance),
		ConcurrentScaleOutForDesiredCapacity: utils.PtrBool(g.ConcurrentScaleOutForDesiredCapacity),
		CreatedTime:                          utils.ParseTimestamp(g.CreatedTime),
		ZoneSet:                              utils.PtrStringSlice(g.ZoneSet),
		SubnetIdSet:                          utils.PtrStringSlice(g.SubnetIdSet),
		LoadBalancerIdSet:                    utils.PtrStringSlice(g.LoadBalancerIdSet),
		TerminationPolicySet:                 utils.PtrStringSlice(g.TerminationPolicySet),
		ForwardLoadBalancerSet:               g.ForwardLoadBalancerSet,
		ServiceSettings:                      g.ServiceSettings,
		SpotMixedAllocationPolicy:            g.SpotMixedAllocationPolicy,
		InstanceNameIndexSettings:            g.InstanceNameIndexSettings,
		HostNameIndexSettings:                g.HostNameIndexSettings,
		Tags:                                 TagsAsToMap(g.Tags),
	}

	if row.AutoScalingGroupId != "" {
		row.Akas = []string{"tencentcloud:cvm:auto-scaling-group:" + row.AutoScalingGroupId}
	}
	return row
}
