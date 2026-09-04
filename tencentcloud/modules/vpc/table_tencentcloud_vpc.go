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

func TableTencentcloudVpc() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_vpc",
		Description:       "Tencent Cloud Virtual Private Cloud (VPC) instances that isolate cloud resources in a user-defined network.",
		GetMatrixItemFunc: buildVpcRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "vpc_id", Require: plugin.Optional},
				{Name: "vpc_name", Require: plugin.Optional},
				{Name: "is_default", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "cidr_block", Require: plugin.Optional},
			},
			Hydrate: listVpcVpcs,
			Tags:    map[string]string{"service": "vpc", "action": "DescribeVpcs"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("vpc_id"),
			Hydrate:    getVpcVpc,
			Tags:       map[string]string{"service": "vpc", "action": "DescribeVpcs"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "vpc_id", Type: proto.ColumnType_STRING, Description: "The VPC instance ID, such as `vpc-azd4dt1c`.", Transform: transform.FromField("VpcId")},
			{Name: "vpc_name", Type: proto.ColumnType_STRING, Description: "The VPC name.", Transform: transform.FromField("VpcName")},
			{Name: "cidr_block", Type: proto.ColumnType_STRING, Description: "The IPv4 CIDR block of the VPC.", Transform: transform.FromField("CidrBlock")},
			{Name: "is_default", Type: proto.ColumnType_BOOL, Description: "Whether it is the default VPC.", Transform: transform.FromField("IsDefault")},
			{Name: "enable_multicast", Type: proto.ColumnType_BOOL, Description: "Whether multicast is enabled.", Transform: transform.FromField("EnableMulticast")},
			{Name: "enable_dhcp", Type: proto.ColumnType_BOOL, Description: "Whether DHCP is enabled.", Transform: transform.FromField("EnableDhcp")},
			{Name: "enable_route_vpc_publish", Type: proto.ColumnType_BOOL, Description: "Whether the VPC publishes CIDR routes to CCN (true) or subnet routes (false).", Transform: transform.FromField("EnableRouteVpcPublish")},
			{Name: "domain_name", Type: proto.ColumnType_STRING, Description: "The DHCP domain name option value.", Transform: transform.FromField("DomainName")},
			{Name: "dhcp_options_id", Type: proto.ColumnType_STRING, Description: "The DHCP options set ID.", Transform: transform.FromField("DhcpOptionsId")},
			{Name: "ipv6_cidr_block", Type: proto.ColumnType_STRING, Description: "The IPv6 CIDR block of the VPC.", Transform: transform.FromField("Ipv6CidrBlock")},
			{Name: "dns_servers", Type: proto.ColumnType_JSON, Description: "The DNS server list of the VPC.", Transform: transform.FromField("DnsServerSet")},
			{Name: "assistant_cidr_set", Type: proto.ColumnType_JSON, Description: "The auxiliary CIDR blocks attached to the VPC.", Transform: transform.FromField("AssistantCidrSet")},
			{Name: "ipv6_cidr_block_set", Type: proto.ColumnType_JSON, Description: "The multi-operator IPv6 CIDR blocks of the VPC.", Transform: transform.FromField("Ipv6CidrBlockSet")},
			{Name: "created_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the VPC.", Transform: transform.FromField("CreatedTime").NullIfZero()},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The tags attached to the VPC as a map of key -> value.", Transform: transform.FromField("Tags")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the VPC resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("VpcName")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listVpcVpcs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listVpcVpcs", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listVpcVpcs", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	pageSize := strconv.FormatInt(utils.PageSizeFromLimit(d, 100, 1), 10)
	offset := "0"

	vpcIds := utils.SingleIdFromQual(d.EqualsQuals, "vpc_id")
	filters := buildVpcFilters(d.EqualsQuals)

	for {
		d.WaitForListRateLimit(ctx)

		req := vpc.NewDescribeVpcsRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		if len(vpcIds) > 0 {
			req.VpcIds = vpcIds
		} else if len(filters) > 0 {
			req.Filters = filters
		}

		utils.LogRequest(ctx, "listVpcVpcs", req)
		resp, err := client.DescribeVpcsWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listVpcVpcs", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		vpcs := resp.Response.VpcSet
		total := utils.PtrUint64(resp.Response.TotalCount)
		currentOffset, _ := strconv.ParseUint(offset, 10, 64)

		for _, item := range vpcs {
			if item == nil {
				continue
			}
			row := toVpcRow(item)
			row.Region = region
			d.StreamListItem(ctx, row)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if utils.PageDone(int64(currentOffset), int64(len(vpcs)), total) {
			break
		}
		offset = strconv.FormatUint(currentOffset+uint64(len(vpcs)), 10)
	}

	return nil, nil
}

func buildVpcFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*vpc.Filter {
	filters := utils.BuildFilters(equalsQuals, newVpcFilter,
		utils.FilterMapping{QualName: "vpc_name", FilterName: "vpc-name"},
		utils.FilterMapping{QualName: "cidr_block", FilterName: "cidr-block"},
	)

	if equalsQuals["is_default"] != nil {
		value := equalsQuals["is_default"].GetBoolValue()
		filters = append(filters, newVpcFilter("is-default", strconv.FormatBool(value)))
	}

	return filters
}

// Get Function

func getVpcVpc(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	vpcId := d.EqualsQuals["vpc_id"].GetStringValue()
	plugin.Logger(ctx).Info("getVpcVpc", "vpc_id", vpcId, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getVpcVpc", "client_init_error", err)
		return nil, err
	}

	if vpcId == "" {
		return nil, nil
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	req := vpc.NewDescribeVpcsRequest()
	req.VpcIds = common.StringPtrs([]string{vpcId})

	utils.LogRequest(ctx, "getVpcVpc", req)
	resp, err := client.DescribeVpcsWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getVpcVpc", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.VpcSet) == 0 {
		return nil, nil
	}

	row := toVpcRow(resp.Response.VpcSet[0])
	row.Region = region
	return row, nil
}

// Row Type

type vpcVpcRow struct {
	VpcId                 string
	VpcName               string
	CidrBlock             string
	IsDefault             bool
	EnableMulticast       bool
	EnableDhcp            bool
	EnableRouteVpcPublish bool
	DomainName            string
	DhcpOptionsId         string
	Ipv6CidrBlock         string
	DnsServerSet          []string
	AssistantCidrSet      []*vpc.AssistantCidr
	Ipv6CidrBlockSet      []*vpc.ISPIPv6CidrBlock
	CreatedTime           *time.Time
	Tags                  map[string]string
	Region                string
	Akas                  []string
}

func toVpcRow(v *vpc.Vpc) vpcVpcRow {
	row := vpcVpcRow{
		VpcId:                 utils.PtrString(v.VpcId),
		VpcName:               utils.PtrString(v.VpcName),
		CidrBlock:             utils.PtrString(v.CidrBlock),
		IsDefault:             utils.PtrBool(v.IsDefault),
		EnableMulticast:       utils.PtrBool(v.EnableMulticast),
		EnableDhcp:            utils.PtrBool(v.EnableDhcp),
		EnableRouteVpcPublish: utils.PtrBool(v.EnableRouteVpcPublish),
		DomainName:            utils.PtrString(v.DomainName),
		DhcpOptionsId:         utils.PtrString(v.DhcpOptionsId),
		Ipv6CidrBlock:         utils.PtrString(v.Ipv6CidrBlock),
		DnsServerSet:          utils.PtrStringSlice(v.DnsServerSet),
		AssistantCidrSet:      v.AssistantCidrSet,
		Ipv6CidrBlockSet:      v.Ipv6CidrBlockSet,
		CreatedTime:           utils.ParseTimestamp(v.CreatedTime),
		Tags:                  TagsToMap(v.TagSet),
	}

	if row.VpcId != "" {
		row.Akas = []string{"tencentcloud:vpc:vpc:" + row.VpcId}
	}
	return row
}
