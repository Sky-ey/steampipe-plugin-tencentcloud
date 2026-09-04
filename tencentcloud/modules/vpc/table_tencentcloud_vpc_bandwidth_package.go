package vpc

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/vpc/v20170312"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudVpcBandwidthPackage() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_vpc_bandwidth_package",
		Description:       "Tencent Cloud VPC bandwidth packages that aggregate and share bandwidth across multiple public IP resources.",
		GetMatrixItemFunc: buildVpcRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "bandwidth_package_id", Require: plugin.Optional},
				{Name: "bandwidth_package_name", Require: plugin.Optional},
				{Name: "network_type", Require: plugin.Optional},
				{Name: "charge_type", Require: plugin.Optional},
				{Name: "resource_type", Require: plugin.Optional},
				{Name: "resource_id", Require: plugin.Optional},
				{Name: "address_ip", Require: plugin.Optional},
			},
			Hydrate: listVpcBandwidthPackages,
			Tags:    map[string]string{"service": "vpc", "action": "DescribeBandwidthPackages"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("bandwidth_package_id"),
			Hydrate:    getVpcBandwidthPackage,
			Tags:       map[string]string{"service": "vpc", "action": "DescribeBandwidthPackages"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "bandwidth_package_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the bandwidth package.", Transform: transform.FromField("BandwidthPackageId")},
			{Name: "bandwidth_package_name", Type: proto.ColumnType_STRING, Description: "The name of the bandwidth package.", Transform: transform.FromField("BandwidthPackageName")},
			{Name: "network_type", Type: proto.ColumnType_STRING, Description: "The bandwidth package type (BGP, HIGH_QUALITY_BGP, ANYCAST, SINGLEISP_CMCC, SINGLEISP_CTCC, SINGLEISP_CUCC).", Transform: transform.FromField("NetworkType")},
			{Name: "charge_type", Type: proto.ColumnType_STRING, Description: "The bandwidth package billing type (ENHANCED95_POSTPAID_BY_MONTH, PRIMARY_TRAFFIC_POSTPAID_BY_HOUR, BANDWIDTH_POSTPAID_BY_DAY, PEAK_BANDWIDTH_POSTPAID_BY_DAY, TOP5_POSTPAID_BY_MONTH).", Transform: transform.FromField("ChargeType")},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of the bandwidth package (CREATING, CREATED, DELETING, DELETED).", Transform: transform.FromField("Status")},
			{Name: "bandwidth", Type: proto.ColumnType_INT, Description: "The limit of the bandwidth package in Mbps; -1 indicates no limit.", Transform: transform.FromField("Bandwidth")},
			{Name: "egress", Type: proto.ColumnType_STRING, Description: "The network egress of the bandwidth package (center_egress1, center_egress2, center_egress3).", Transform: transform.FromField("Egress")},
			{Name: "resource_set", Type: proto.ColumnType_JSON, Description: "The resources (EIPs or load balancers) attached to the bandwidth package, each with resource type, resource ID and IP.", Transform: transform.FromField("ResourceSet")},
			{Name: "resource_type", Type: proto.ColumnType_STRING, Description: "The type of the first resource attached to the package (Address, LoadBalance); mirrors the resource.resource-type filter.", Transform: transform.FromField("ResourceType")},
			{Name: "resource_id", Type: proto.ColumnType_STRING, Description: "The ID of the first resource attached to the package (e.g. `eip-xxxx`, `lb-xxxx`); mirrors the resource.resource-id filter.", Transform: transform.FromField("ResourceId")},
			{Name: "address_ip", Type: proto.ColumnType_IPADDR, Description: "The IP of the first resource attached to the package; mirrors the resource.address-ip filter.", Transform: transform.FromField("AddressIp")},
			{Name: "created_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the bandwidth package (ISO 8601, UTC).", Transform: transform.FromField("CreatedTime").NullIfZero()},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the bandwidth package resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("BandwidthPackageName")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listVpcBandwidthPackages(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listVpcBandwidthPackages", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listVpcBandwidthPackages", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	pageSize := uint64(utils.PageSizeFromLimit(d, 100, 1))
	var offset uint64

	bandwidthPackageIds := utils.SingleIdFromQual(d.EqualsQuals, "bandwidth_package_id")
	filters := buildVpcBandwidthPackageFilters(d.EqualsQuals)

	for {
		d.WaitForListRateLimit(ctx)

		req := vpc.NewDescribeBandwidthPackagesRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		if len(bandwidthPackageIds) > 0 {
			req.BandwidthPackageIds = bandwidthPackageIds
		} else if len(filters) > 0 {
			req.Filters = filters
		}

		utils.LogRequest(ctx, "listVpcBandwidthPackages", req)
		resp, err := client.DescribeBandwidthPackagesWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listVpcBandwidthPackages", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		packages := resp.Response.BandwidthPackageSet
		total := utils.PtrUint64(resp.Response.TotalCount)

		for _, item := range packages {
			if item == nil {
				continue
			}
			row := toVpcBandwidthPackageRow(item)
			row.Region = region
			d.StreamListItem(ctx, row)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if utils.PageDone(int64(offset), int64(len(packages)), total) {
			break
		}
		offset += uint64(len(packages))
	}

	return nil, nil
}

func buildVpcBandwidthPackageFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*vpc.Filter {
	return utils.BuildFilters(equalsQuals, newVpcFilter,
		utils.FilterMapping{QualName: "bandwidth_package_name", FilterName: "bandwidth-package-name"},
		utils.FilterMapping{QualName: "network_type", FilterName: "network-type"},
		utils.FilterMapping{QualName: "charge_type", FilterName: "charge-type"},
		utils.FilterMapping{QualName: "resource_type", FilterName: "resource.resource-type"},
		utils.FilterMapping{QualName: "resource_id", FilterName: "resource.resource-id"},
		utils.FilterMapping{QualName: "address_ip", FilterName: "resource.address-ip"},
	)
}

// Get Function

func getVpcBandwidthPackage(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	bandwidthPackageId := d.EqualsQuals["bandwidth_package_id"].GetStringValue()
	plugin.Logger(ctx).Info("getVpcBandwidthPackage", "bandwidth_package_id", bandwidthPackageId, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getVpcBandwidthPackage", "client_init_error", err)
		return nil, err
	}

	if bandwidthPackageId == "" {
		return nil, nil
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	req := vpc.NewDescribeBandwidthPackagesRequest()
	req.BandwidthPackageIds = common.StringPtrs([]string{bandwidthPackageId})

	utils.LogRequest(ctx, "getVpcBandwidthPackage", req)
	resp, err := client.DescribeBandwidthPackagesWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getVpcBandwidthPackage", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.BandwidthPackageSet) == 0 {
		return nil, nil
	}

	row := toVpcBandwidthPackageRow(resp.Response.BandwidthPackageSet[0])
	row.Region = region
	return row, nil
}

// Row Type

type vpcBandwidthPackageRow struct {
	BandwidthPackageId   string
	BandwidthPackageName string
	NetworkType          string
	ChargeType           string
	Status               string
	Bandwidth            int64
	Egress               string
	ResourceSet          []*vpc.Resource
	ResourceType         string
	ResourceId           string
	AddressIp            string
	CreatedTime          *time.Time
	Region               string
	Akas                 []string
}

func toVpcBandwidthPackageRow(b *vpc.BandwidthPackage) vpcBandwidthPackageRow {
	row := vpcBandwidthPackageRow{
		BandwidthPackageId:   utils.PtrString(b.BandwidthPackageId),
		BandwidthPackageName: utils.PtrString(b.BandwidthPackageName),
		NetworkType:          utils.PtrString(b.NetworkType),
		ChargeType:           utils.PtrString(b.ChargeType),
		Status:               utils.PtrString(b.Status),
		Bandwidth:            utils.PtrInt64(b.Bandwidth),
		Egress:               utils.PtrString(b.Egress),
		ResourceSet:          b.ResourceSet,
		CreatedTime:          utils.ParseTimestamp(b.CreatedTime),
	}

	// The filter-only quals (resource_type, resource_id, address_ip) are
	// echoed from the first attached resource so that the key columns
	// reference existing columns in the table.
	if len(b.ResourceSet) > 0 {
		row.ResourceType = utils.PtrString(b.ResourceSet[0].ResourceType)
		row.ResourceId = utils.PtrString(b.ResourceSet[0].ResourceId)
		row.AddressIp = utils.PtrString(b.ResourceSet[0].AddressIp)
	}

	if row.BandwidthPackageId != "" {
		row.Akas = []string{"tencentcloud:vpc:bandwidth-package:" + row.BandwidthPackageId}
	}
	return row
}
