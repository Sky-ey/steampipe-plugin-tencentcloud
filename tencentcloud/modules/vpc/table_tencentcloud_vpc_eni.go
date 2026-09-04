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

func TableTencentcloudVpcEni() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_vpc_eni",
		Description:       "Tencent Cloud VPC elastic network interfaces (ENIs) that can be attached to CVM instances.",
		GetMatrixItemFunc: buildVpcRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "network_interface_id", Require: plugin.Optional},
				{Name: "vpc_id", Require: plugin.Optional},
				{Name: "subnet_id", Require: plugin.Optional},
				{Name: "network_interface_name", Require: plugin.Optional},
				{Name: "instance_id", Require: plugin.Optional},
				{Name: "security_group_id", Require: plugin.Optional},
			},
			Hydrate: listVpcEnis,
			Tags:    map[string]string{"service": "vpc", "action": "DescribeNetworkInterfaces"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("network_interface_id"),
			Hydrate:    getVpcEni,
			Tags:       map[string]string{"service": "vpc", "action": "DescribeNetworkInterfaces"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "network_interface_id", Type: proto.ColumnType_STRING, Description: "The ENI instance ID, such as `eni-f1xjkw1b`.", Transform: transform.FromField("NetworkInterfaceId")},
			{Name: "network_interface_name", Type: proto.ColumnType_STRING, Description: "The ENI name.", Transform: transform.FromField("NetworkInterfaceName")},
			{Name: "network_interface_description", Type: proto.ColumnType_STRING, Description: "The ENI description.", Transform: transform.FromField("NetworkInterfaceDescription")},
			{Name: "vpc_id", Type: proto.ColumnType_STRING, Description: "The VPC instance ID.", Transform: transform.FromField("VpcId")},
			{Name: "subnet_id", Type: proto.ColumnType_STRING, Description: "The subnet instance ID.", Transform: transform.FromField("SubnetId")},
			{Name: "zone", Type: proto.ColumnType_STRING, Description: "The availability zone of the ENI.", Transform: transform.FromField("Zone")},
			{Name: "instance_id", Type: proto.ColumnType_STRING, Description: "The CVM instance ID the ENI is attached to, if any.", Transform: transform.FromField("InstanceId")},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "The ENI state (PENDING, AVAILABLE, ATTACHING, DETACHING, DELETING).", Transform: transform.FromField("State")},
			{Name: "primary", Type: proto.ColumnType_BOOL, Description: "Whether it is the primary ENI.", Transform: transform.FromField("Primary")},
			{Name: "eni_type", Type: proto.ColumnType_INT, Description: "The ENI type. 0: ENI; 1: EVM ENI.", Transform: transform.FromField("EniType")},
			{Name: "attach_type", Type: proto.ColumnType_INT, Description: "ENI attach type. 0 (standard); 1 (extension).", Transform: transform.FromField("AttachType")},
			{Name: "business", Type: proto.ColumnType_STRING, Description: "Type of the resource bound with the ENI (cvm, eks).", Transform: transform.FromField("Business")},
			{Name: "resource_id", Type: proto.ColumnType_STRING, Description: "The ID of the resource that retains the ENI primary IP.", Transform: transform.FromField("ResourceId")},
			{Name: "qos_level", Type: proto.ColumnType_STRING, Description: "The service level (DEFAULT, PT, AU, AG).", Transform: transform.FromField("QosLevel")},
			{Name: "mac_address", Type: proto.ColumnType_STRING, Description: "The MAC address of the ENI.", Transform: transform.FromField("MacAddress")},
			{Name: "cdc_id", Type: proto.ColumnType_STRING, Description: "The CDC instance ID associated with the ENI.", Transform: transform.FromField("CdcId")},
			{Name: "security_group_id", Type: proto.ColumnType_STRING, Description: "The first bound security group ID. See group_set for the full list.", Transform: transform.FromField("SecurityGroupId")},
			{Name: "group_set", Type: proto.ColumnType_JSON, Description: "The list of bound security group IDs.", Transform: transform.FromField("GroupSet")},
			{Name: "private_ip_address_set", Type: proto.ColumnType_JSON, Description: "The private IP information of the ENI.", Transform: transform.FromField("PrivateIpAddressSet")},
			{Name: "ipv6_address_set", Type: proto.ColumnType_JSON, Description: "The IPv6 address list of the ENI.", Transform: transform.FromField("Ipv6AddressSet")},
			{Name: "attachment", Type: proto.ColumnType_JSON, Description: "The bound CVM object of the ENI.", Transform: transform.FromField("Attachment")},
			{Name: "created_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the ENI.", Transform: transform.FromField("CreatedTime").NullIfZero()},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The tags attached to the ENI as a map of key -> value.", Transform: transform.FromField("Tags")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the ENI resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("NetworkInterfaceName")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listVpcEnis(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listVpcEnis", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listVpcEnis", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	pageSize := uint64(utils.PageSizeFromLimit(d, 100, 1))
	var offset uint64

	eniIds := utils.SingleIdFromQual(d.EqualsQuals, "network_interface_id")
	filters := buildVpcEniFilters(d.EqualsQuals)

	for {
		d.WaitForListRateLimit(ctx)

		req := vpc.NewDescribeNetworkInterfacesRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		if len(eniIds) > 0 {
			req.NetworkInterfaceIds = eniIds
		} else if len(filters) > 0 {
			req.Filters = filters
		}

		utils.LogRequest(ctx, "listVpcEnis", req)
		resp, err := client.DescribeNetworkInterfacesWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listVpcEnis", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		enis := resp.Response.NetworkInterfaceSet
		total := utils.PtrUint64(resp.Response.TotalCount)

		for _, item := range enis {
			if item == nil {
				continue
			}
			row := toVpcEniRow(item)
			row.Region = region
			d.StreamListItem(ctx, row)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if utils.PageDone(int64(offset), int64(len(enis)), total) {
			break
		}
		offset += uint64(len(enis))
	}

	return nil, nil
}

func buildVpcEniFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*vpc.Filter {
	return utils.BuildFilters(equalsQuals, newVpcFilter,
		utils.FilterMapping{QualName: "vpc_id", FilterName: "vpc-id"},
		utils.FilterMapping{QualName: "subnet_id", FilterName: "subnet-id"},
		utils.FilterMapping{QualName: "network_interface_name", FilterName: "network-interface-name"},
		utils.FilterMapping{QualName: "instance_id", FilterName: "attachment.instance-id"},
		utils.FilterMapping{QualName: "security_group_id", FilterName: "groups.security-group-id"},
	)
}

// Get Function

func getVpcEni(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	eniId := d.EqualsQuals["network_interface_id"].GetStringValue()
	plugin.Logger(ctx).Info("getVpcEni", "network_interface_id", eniId, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getVpcEni", "client_init_error", err)
		return nil, err
	}

	if eniId == "" {
		return nil, nil
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	req := vpc.NewDescribeNetworkInterfacesRequest()
	req.NetworkInterfaceIds = common.StringPtrs([]string{eniId})

	utils.LogRequest(ctx, "getVpcEni", req)
	resp, err := client.DescribeNetworkInterfacesWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getVpcEni", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.NetworkInterfaceSet) == 0 {
		return nil, nil
	}

	row := toVpcEniRow(resp.Response.NetworkInterfaceSet[0])
	row.Region = region
	return row, nil
}

// Row Type

type vpcEniRow struct {
	NetworkInterfaceId          string
	NetworkInterfaceName        string
	NetworkInterfaceDescription string
	VpcId                       string
	SubnetId                    string
	Zone                        string
	InstanceId                  string
	State                       string
	Primary                     bool
	EniType                     uint64
	AttachType                  uint64
	Business                    string
	ResourceId                  string
	QosLevel                    string
	MacAddress                  string
	CdcId                       string
	SecurityGroupId             string
	GroupSet                    []string
	PrivateIpAddressSet         []*vpc.PrivateIpAddressSpecification
	Ipv6AddressSet              []*vpc.Ipv6Address
	Attachment                  *vpc.NetworkInterfaceAttachment
	CreatedTime                 *time.Time
	Tags                        map[string]string
	Region                      string
	Akas                        []string
}

func toVpcEniRow(n *vpc.NetworkInterface) vpcEniRow {
	row := vpcEniRow{
		NetworkInterfaceId:          utils.PtrString(n.NetworkInterfaceId),
		NetworkInterfaceName:        utils.PtrString(n.NetworkInterfaceName),
		NetworkInterfaceDescription: utils.PtrString(n.NetworkInterfaceDescription),
		VpcId:                       utils.PtrString(n.VpcId),
		SubnetId:                    utils.PtrString(n.SubnetId),
		Zone:                        utils.PtrString(n.Zone),
		State:                       utils.PtrString(n.State),
		Primary:                     utils.PtrBool(n.Primary),
		EniType:                     utils.PtrUint64(n.EniType),
		AttachType:                  utils.PtrUint64(n.AttachType),
		Business:                    utils.PtrString(n.Business),
		ResourceId:                  utils.PtrString(n.ResourceId),
		QosLevel:                    utils.PtrString(n.QosLevel),
		MacAddress:                  utils.PtrString(n.MacAddress),
		CdcId:                       utils.PtrString(n.CdcId),
		GroupSet:                    utils.PtrStringSlice(n.GroupSet),
		PrivateIpAddressSet:         n.PrivateIpAddressSet,
		Ipv6AddressSet:              n.Ipv6AddressSet,
		Attachment:                  n.Attachment,
		CreatedTime:                 utils.ParseTimestamp(n.CreatedTime),
		Tags:                        TagsToMap(n.TagSet),
	}

	if n.Attachment != nil {
		row.InstanceId = utils.PtrString(n.Attachment.InstanceId)
	}
	if len(row.GroupSet) > 0 {
		row.SecurityGroupId = row.GroupSet[0]
	}

	if row.NetworkInterfaceId != "" {
		row.Akas = []string{"tencentcloud:vpc:eni:" + row.NetworkInterfaceId}
	}
	return row
}
