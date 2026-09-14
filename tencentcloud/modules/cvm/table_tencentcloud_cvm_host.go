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

func TableTencentcloudCvmHost() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_cvm_host",
		Description:       "Tencent Cloud Cloud Dedicated Host (CDH) instances that provide dedicated, isolated physical servers for running CVM instances.",
		GetMatrixItemFunc: buildCvmRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "host_id", Require: plugin.Optional},
				{Name: "host_name", Require: plugin.Optional},
				{Name: "host_state", Require: plugin.Optional},
				{Name: "zone", Require: plugin.Optional},
				{Name: "project_id", Require: plugin.Optional},
			},
			Hydrate: listCvmHosts,
			Tags:    map[string]string{"service": "cvm", "action": "DescribeHosts"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("host_id"),
			Hydrate:    getCvmHost,
			Tags:       map[string]string{"service": "cvm", "action": "DescribeHosts"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "host_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the CDH instance (e.g., host-xxxxxxxx).", Transform: transform.FromField("HostId")},
			{Name: "host_name", Type: proto.ColumnType_STRING, Description: "The name of the CDH instance.", Transform: transform.FromField("HostName")},
			{Name: "host_state", Type: proto.ColumnType_STRING, Description: "The state of the CDH instance (PENDING, LAUNCH_FAILURE, RUNNING, EXPIRED).", Transform: transform.FromField("HostState")},
			{Name: "host_type", Type: proto.ColumnType_STRING, Description: "The type of the CDH instance (e.g., C3.LARGE).", Transform: transform.FromField("HostType")},
			{Name: "host_charge_type", Type: proto.ColumnType_STRING, Description: "The billing mode of the CDH instance (PREPAID, POSTPAID_BY_HOUR).", Transform: transform.FromField("HostChargeType")},
			{Name: "renew_flag", Type: proto.ColumnType_STRING, Description: "The auto-renewal flag of the CDH instance (NOTIFY_AND_MANUAL_RENEW, NOTIFY_AND_AUTO_RENEW, DISABLE_NOTIFY_AND_MANUAL_RENEW).", Transform: transform.FromField("RenewFlag")},
			{Name: "created_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the CDH instance in ISO 8601 format (UTC).", Transform: transform.FromField("CreatedTime").NullIfZero()},
			{Name: "expired_time", Type: proto.ColumnType_TIMESTAMP, Description: "The expiry time of the CDH instance in ISO 8601 format (UTC).", Transform: transform.FromField("ExpiredTime").NullIfZero()},
			{Name: "host_ip", Type: proto.ColumnType_IPADDR, Description: "The IP address of the CDH instance.", Transform: transform.FromField("HostIp")},
			{Name: "instance_ids", Type: proto.ColumnType_JSON, Description: "The IDs of the CVM instances running on the CDH instance.", Transform: transform.FromField("InstanceIds")},
			{Name: "cage_id", Type: proto.ColumnType_STRING, Description: "The cage ID of the CDH instance. Only valid for CDH instances in the cages of finance availability zones.", Transform: transform.FromField("CageId")},
			{Name: "zone", Type: proto.ColumnType_STRING, Description: "The availability zone the CDH instance is deployed in.", Transform: transform.FromField("Zone")},
			{Name: "project_id", Type: proto.ColumnType_INT, Description: "The project ID the CDH instance belongs to. 0 means the default project.", Transform: transform.FromField("ProjectId")},
			{Name: "cpu_total", Type: proto.ColumnType_INT, Description: "The total number of CPU cores of the CDH instance.", Transform: transform.FromField("CpuTotal")},
			{Name: "cpu_available", Type: proto.ColumnType_INT, Description: "The number of available CPU cores of the CDH instance.", Transform: transform.FromField("CpuAvailable")},
			{Name: "mem_total", Type: proto.ColumnType_DOUBLE, Description: "The total memory size of the CDH instance in GiB.", Transform: transform.FromField("MemTotal")},
			{Name: "mem_available", Type: proto.ColumnType_DOUBLE, Description: "The available memory size of the CDH instance in GiB.", Transform: transform.FromField("MemAvailable")},
			{Name: "disk_total", Type: proto.ColumnType_INT, Description: "The total disk size of the CDH instance in GiB.", Transform: transform.FromField("DiskTotal")},
			{Name: "disk_available", Type: proto.ColumnType_INT, Description: "The available disk size of the CDH instance in GiB.", Transform: transform.FromField("DiskAvailable")},
			{Name: "disk_type", Type: proto.ColumnType_STRING, Description: "The disk type of the CDH instance.", Transform: transform.FromField("DiskType")},
			{Name: "gpu_total", Type: proto.ColumnType_INT, Description: "The total number of GPU cards of the CDH instance.", Transform: transform.FromField("GpuTotal")},
			{Name: "gpu_available", Type: proto.ColumnType_INT, Description: "The number of available GPU cards of the CDH instance.", Transform: transform.FromField("GpuAvailable")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the CDH instance resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("HostName")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listCvmHosts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listCvmHosts", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &cvm.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listCvmHosts", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	pageSize := uint64(utils.PageSizeFromLimit(d, 100, 1))
	var offset uint64 = 0

	filters := buildCvmHostFilters(d.EqualsQuals)

	for {
		d.WaitForListRateLimit(ctx)

		req := cvm.NewDescribeHostsRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		req.Filters = filters

		utils.LogRequest(ctx, "listCvmHosts", req)
		resp, err := client.DescribeHostsWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listCvmHosts", "request_error", err)
			return nil, err
		}

		if resp == nil || resp.Response == nil {
			break
		}
		hosts := resp.Response.HostSet
		for _, host := range hosts {
			if host == nil {
				continue
			}
			row := toCvmHostRow(host)
			row.Region = region
			d.StreamListItem(ctx, row)

			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if resp.Response == nil {
			break
		}
		total := utils.PtrUint64(resp.Response.TotalCount)
		if utils.PageDone(int64(offset), int64(len(hosts)), total) {
			break
		}
		offset += uint64(len(hosts))
	}

	return nil, nil
}

func buildCvmHostFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*cvm.Filter {
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

	addFilter("host_id", "host-id")
	addFilter("host_name", "host-name")
	addFilter("host_state", "host-state")
	addFilter("zone", "zone")

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

func getCvmHost(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	hostId := d.EqualsQuals["host_id"].GetStringValue()
	plugin.Logger(ctx).Info("getCvmHost", "host_id", hostId, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &cvm.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getCvmHost", "client_init_error", err)
		return nil, err
	}

	if hostId == "" {
		return nil, nil
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	req := cvm.NewDescribeHostsRequest()
	req.Filters = []*cvm.Filter{
		{
			Name:   common.StringPtr("host-id"),
			Values: common.StringPtrs([]string{hostId}),
		},
	}

	utils.LogRequest(ctx, "getCvmHost", req)
	resp, err := client.DescribeHostsWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getCvmHost", "request_error", err)
		return nil, err
	}

	if resp == nil || resp.Response == nil || len(resp.Response.HostSet) == 0 {
		return nil, nil
	}
	row := toCvmHostRow(resp.Response.HostSet[0])
	row.Region = region
	return row, nil
}

// Row Type

type cvmHostRow struct {
	HostId         string
	HostName       string
	HostState      string
	HostType       string
	HostChargeType string
	RenewFlag      string
	CreatedTime    *time.Time
	ExpiredTime    *time.Time
	HostIp         string
	InstanceIds    []string
	CageId         string
	Zone           string
	ProjectId      int64
	CpuTotal       uint64
	CpuAvailable   uint64
	MemTotal       float64
	MemAvailable   float64
	DiskTotal      uint64
	DiskAvailable  uint64
	DiskType       string
	GpuTotal       uint64
	GpuAvailable   uint64
	Region         string
	Akas           []string
}

func toCvmHostRow(host *cvm.HostItem) cvmHostRow {
	row := cvmHostRow{
		HostId:         utils.PtrString(host.HostId),
		HostName:       utils.PtrString(host.HostName),
		HostState:      utils.PtrString(host.HostState),
		HostType:       utils.PtrString(host.HostType),
		HostChargeType: utils.PtrString(host.HostChargeType),
		RenewFlag:      utils.PtrString(host.RenewFlag),
		CreatedTime:    utils.ParseTimestamp(host.CreatedTime),
		ExpiredTime:    utils.ParseTimestamp(host.ExpiredTime),
		HostIp:         utils.PtrString(host.HostIp),
		InstanceIds:    utils.PtrStringSlice(host.InstanceIds),
		CageId:         utils.PtrString(host.CageId),
	}

	if host.Placement != nil {
		row.Zone = utils.PtrString(host.Placement.Zone)
		row.ProjectId = utils.PtrInt64(host.Placement.ProjectId)
	}
	if host.HostResource != nil {
		row.CpuTotal = utils.PtrUint64(host.HostResource.CpuTotal)
		row.CpuAvailable = utils.PtrUint64(host.HostResource.CpuAvailable)
		row.MemTotal = utils.PtrFloat64(host.HostResource.MemTotal)
		row.MemAvailable = utils.PtrFloat64(host.HostResource.MemAvailable)
		row.DiskTotal = utils.PtrUint64(host.HostResource.DiskTotal)
		row.DiskAvailable = utils.PtrUint64(host.HostResource.DiskAvailable)
		row.DiskType = utils.PtrString(host.HostResource.DiskType)
		row.GpuTotal = utils.PtrUint64(host.HostResource.GpuTotal)
		row.GpuAvailable = utils.PtrUint64(host.HostResource.GpuAvailable)
	}

	if row.HostId != "" {
		row.Akas = []string{"tencentcloud:cvm:host:" + row.HostId}
	}
	return row
}
