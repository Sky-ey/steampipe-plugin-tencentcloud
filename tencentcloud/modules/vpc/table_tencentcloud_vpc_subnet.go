package vpc

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"strconv"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/vpc/v20170312"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudVpcSubnet() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_vpc_subnet",
		Description:       "Tencent Cloud VPC subnets that segment a VPC into smaller IP ranges per availability zone.",
		GetMatrixItemFunc: buildVpcRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "subnet_id", Require: plugin.Optional},
				{Name: "vpc_id", Require: plugin.Optional},
				{Name: "subnet_name", Require: plugin.Optional},
				{Name: "zone", Require: plugin.Optional},
				{Name: "is_default", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "cidr_block", Require: plugin.Optional},
			},
			Hydrate: listVpcSubnets,
			Tags:    map[string]string{"service": "vpc", "action": "DescribeSubnets"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("subnet_id"),
			Hydrate:    getVpcSubnet,
			Tags:       map[string]string{"service": "vpc", "action": "DescribeSubnets"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "subnet_id", Type: proto.ColumnType_STRING, Description: "The subnet instance ID, such as `subnet-bthucmmy`.", Transform: transform.FromField("SubnetId")},
			{Name: "subnet_name", Type: proto.ColumnType_STRING, Description: "The subnet name.", Transform: transform.FromField("SubnetName")},
			{Name: "vpc_id", Type: proto.ColumnType_STRING, Description: "The ID of the VPC to which the subnet belongs.", Transform: transform.FromField("VpcId")},
			{Name: "cidr_block", Type: proto.ColumnType_STRING, Description: "The IPv4 CIDR block of the subnet, such as `192.168.1.0/24`.", Transform: transform.FromField("CidrBlock")},
			{Name: "ipv6_cidr_block", Type: proto.ColumnType_STRING, Description: "The IPv6 CIDR block of the subnet.", Transform: transform.FromField("Ipv6CidrBlock")},
			{Name: "zone", Type: proto.ColumnType_STRING, Description: "The availability zone of the subnet.", Transform: transform.FromField("Zone")},
			{Name: "route_table_id", Type: proto.ColumnType_STRING, Description: "The route table ID associated with the subnet.", Transform: transform.FromField("RouteTableId")},
			{Name: "network_acl_id", Type: proto.ColumnType_STRING, Description: "The associated network ACL ID.", Transform: transform.FromField("NetworkAclId")},
			{Name: "cdc_id", Type: proto.ColumnType_STRING, Description: "The CDC instance ID if the subnet is associated with a CDC.", Transform: transform.FromField("CdcId")},
			{Name: "is_default", Type: proto.ColumnType_BOOL, Description: "Whether it is the default subnet.", Transform: transform.FromField("IsDefault")},
			{Name: "is_remote_vpc_snat", Type: proto.ColumnType_BOOL, Description: "Whether it is a VPC SNAT address pool subnet.", Transform: transform.FromField("IsRemoteVpcSnat")},
			{Name: "is_cdc_subnet", Type: proto.ColumnType_INT, Description: "Whether the subnet is associated with a CDC. 0 (no), 1 (yes).", Transform: transform.FromField("IsCdcSubnet")},
			{Name: "enable_broadcast", Type: proto.ColumnType_BOOL, Description: "Whether broadcast is enabled.", Transform: transform.FromField("EnableBroadcast")},
			{Name: "available_ip_address_count", Type: proto.ColumnType_INT, Description: "The number of available IPv4 addresses in the subnet.", Transform: transform.FromField("AvailableIpAddressCount")},
			{Name: "total_ip_address_count", Type: proto.ColumnType_INT, Description: "The total number of IPv4 addresses in the subnet.", Transform: transform.FromField("TotalIpAddressCount")},
			{Name: "created_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the subnet.", Transform: transform.FromField("CreatedTime").NullIfZero()},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The tags attached to the subnet as a map of key -> value.", Transform: transform.FromField("Tags")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the subnet resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("SubnetName")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listVpcSubnets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listVpcSubnets", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listVpcSubnets", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	pageSize := strconv.FormatInt(utils.PageSizeFromLimit(d, 100, 1), 10)
	offset := "0"

	subnetIds := utils.SingleIdFromQual(d.EqualsQuals, "subnet_id")
	filters := buildVpcSubnetFilters(d.EqualsQuals)

	for {
		d.WaitForListRateLimit(ctx)

		req := vpc.NewDescribeSubnetsRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		if len(subnetIds) > 0 {
			req.SubnetIds = subnetIds
		} else if len(filters) > 0 {
			req.Filters = filters
		}

		utils.LogRequest(ctx, "listVpcSubnets", req)
		resp, err := client.DescribeSubnetsWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listVpcSubnets", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		subnets := resp.Response.SubnetSet
		total := utils.PtrUint64(resp.Response.TotalCount)
		currentOffset, _ := strconv.ParseUint(offset, 10, 64)

		for _, item := range subnets {
			if item == nil {
				continue
			}
			row := toVpcSubnetRow(item)
			row.Region = region
			d.StreamListItem(ctx, row)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if utils.PageDone(int64(currentOffset), int64(len(subnets)), total) {
			break
		}
		offset = strconv.FormatUint(currentOffset+uint64(len(subnets)), 10)
	}

	return nil, nil
}

func buildVpcSubnetFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*vpc.Filter {
	filters := utils.BuildFilters(equalsQuals, newVpcFilter,
		utils.FilterMapping{QualName: "vpc_id", FilterName: "vpc-id"},
		utils.FilterMapping{QualName: "subnet_name", FilterName: "subnet-name"},
		utils.FilterMapping{QualName: "zone", FilterName: "zone"},
		utils.FilterMapping{QualName: "cidr_block", FilterName: "cidr-block"},
	)

	if equalsQuals["is_default"] != nil {
		value := equalsQuals["is_default"].GetBoolValue()
		filters = append(filters, newVpcFilter("is-default", strconv.FormatBool(value)))
	}

	return filters
}

// Get Function

func getVpcSubnet(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	subnetId := d.EqualsQuals["subnet_id"].GetStringValue()
	plugin.Logger(ctx).Debug("getVpcSubnet", "subnet_id", subnetId, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getVpcSubnet", "client_init_error", err)
		return nil, err
	}

	if subnetId == "" {
		return nil, nil
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	req := vpc.NewDescribeSubnetsRequest()
	req.SubnetIds = common.StringPtrs([]string{subnetId})

	utils.LogRequest(ctx, "getVpcSubnet", req)
	resp, err := client.DescribeSubnetsWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getVpcSubnet", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.SubnetSet) == 0 {
		return nil, nil
	}

	row := toVpcSubnetRow(resp.Response.SubnetSet[0])
	row.Region = region
	return row, nil
}

// Row Type

type vpcSubnetRow struct {
	SubnetId                string
	SubnetName              string
	VpcId                   string
	CidrBlock               string
	Ipv6CidrBlock           string
	Zone                    string
	RouteTableId            string
	NetworkAclId            string
	CdcId                   string
	IsDefault               bool
	IsRemoteVpcSnat         bool
	IsCdcSubnet             int64
	EnableBroadcast         bool
	AvailableIpAddressCount uint64
	TotalIpAddressCount     uint64
	CreatedTime             *time.Time
	Tags                    map[string]string
	Region                  string
	Akas                    []string
}

func toVpcSubnetRow(s *vpc.Subnet) vpcSubnetRow {
	row := vpcSubnetRow{
		SubnetId:                utils.PtrString(s.SubnetId),
		SubnetName:              utils.PtrString(s.SubnetName),
		VpcId:                   utils.PtrString(s.VpcId),
		CidrBlock:               utils.PtrString(s.CidrBlock),
		Ipv6CidrBlock:           utils.PtrString(s.Ipv6CidrBlock),
		Zone:                    utils.PtrString(s.Zone),
		RouteTableId:            utils.PtrString(s.RouteTableId),
		NetworkAclId:            utils.PtrString(s.NetworkAclId),
		CdcId:                   utils.PtrString(s.CdcId),
		IsDefault:               utils.PtrBool(s.IsDefault),
		IsRemoteVpcSnat:         utils.PtrBool(s.IsRemoteVpcSnat),
		IsCdcSubnet:             utils.PtrInt64(s.IsCdcSubnet),
		EnableBroadcast:         utils.PtrBool(s.EnableBroadcast),
		AvailableIpAddressCount: utils.PtrUint64(s.AvailableIpAddressCount),
		TotalIpAddressCount:     utils.PtrUint64(s.TotalIpAddressCount),
		CreatedTime:             utils.ParseTimestamp(s.CreatedTime),
		Tags:                    TagsToMap(s.TagSet),
	}

	if row.SubnetId != "" {
		row.Akas = []string{"tencentcloud:vpc:subnet:" + row.SubnetId}
	}
	return row
}
