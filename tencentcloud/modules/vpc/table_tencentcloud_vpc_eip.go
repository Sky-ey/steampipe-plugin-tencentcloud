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

func TableTencentcloudVpcEip() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_vpc_eip",
		Description:       "Tencent Cloud VPC elastic public IP (EIP) addresses that can be bound to CVM/NAT/ENI/CLB resources.",
		GetMatrixItemFunc: buildVpcRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "address_id", Require: plugin.Optional},
				{Name: "address_name", Require: plugin.Optional},
				{Name: "address_status", Require: plugin.Optional},
				{Name: "instance_id", Require: plugin.Optional},
				{Name: "instance_type", Require: plugin.Optional},
				{Name: "network_interface_id", Require: plugin.Optional},
			},
			Hydrate: listVpcEips,
			Tags:    map[string]string{"service": "vpc", "action": "DescribeAddresses"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("address_id"),
			Hydrate:    getVpcEip,
			Tags:       map[string]string{"service": "vpc", "action": "DescribeAddresses"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "address_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the EIP, such as `eip-11112222`.", Transform: transform.FromField("AddressId")},
			{Name: "address_name", Type: proto.ColumnType_STRING, Description: "The EIP name.", Transform: transform.FromField("AddressName")},
			{Name: "address_status", Type: proto.ColumnType_STRING, Description: "The EIP status (CREATING, BINDING, BIND, UNBINDING, UNBIND, OFFLINING, BIND_ENI).", Transform: transform.FromField("AddressStatus")},
			{Name: "address_ip", Type: proto.ColumnType_IPADDR, Description: "The public IP address of the EIP.", Transform: transform.FromField("AddressIp")},
			{Name: "address_type", Type: proto.ColumnType_STRING, Description: "The EIP resource type (CalcIP, WanIP, EIP, AnycastEIP, AntiDDoSEIP).", Transform: transform.FromField("AddressType")},
			{Name: "instance_id", Type: proto.ColumnType_STRING, Description: "The ID of the bound resource instance.", Transform: transform.FromField("InstanceId")},
			{Name: "instance_type", Type: proto.ColumnType_STRING, Description: "The type of bound resource (CVM, NAT, HAVIP, ENI, CLB, DHCPIP, EKS, VPCE, WAF).", Transform: transform.FromField("InstanceType")},
			{Name: "network_interface_id", Type: proto.ColumnType_STRING, Description: "The bound elastic network interface ID, if any.", Transform: transform.FromField("NetworkInterfaceId")},
			{Name: "private_address_ip", Type: proto.ColumnType_STRING, Description: "The bound resource internal IP, if any.", Transform: transform.FromField("PrivateAddressIp")},
			{Name: "bandwidth", Type: proto.ColumnType_INT, Description: "The bandwidth value of the EIP in Mbps.", Transform: transform.FromField("Bandwidth")},
			{Name: "internet_charge_type", Type: proto.ColumnType_STRING, Description: "Network billing mode (BANDWIDTH_PREPAID_BY_MONTH, TRAFFIC_POSTPAID_BY_HOUR, BANDWIDTH_POSTPAID_BY_HOUR, BANDWIDTH_PACKAGE).", Transform: transform.FromField("InternetChargeType")},
			{Name: "internet_service_provider", Type: proto.ColumnType_STRING, Description: "EIP ISP information (CMCC, CTCC, CUCC, BGP).", Transform: transform.FromField("InternetServiceProvider")},
			{Name: "egress", Type: proto.ColumnType_STRING, Description: "Static single-line IP network egress.", Transform: transform.FromField("Egress")},
			{Name: "is_arrears", Type: proto.ColumnType_BOOL, Description: "Whether the EIP is isolated (overdue).", Transform: transform.FromField("IsArrears")},
			{Name: "is_blocked", Type: proto.ColumnType_BOOL, Description: "Whether the EIP is blocked.", Transform: transform.FromField("IsBlocked")},
			{Name: "is_eip_direct_connection", Type: proto.ColumnType_BOOL, Description: "Whether the EIP supports direct connection mode.", Transform: transform.FromField("IsEipDirectConnection")},
			{Name: "cascade_release", Type: proto.ColumnType_BOOL, Description: "Whether the EIP is automatically released after being unbound.", Transform: transform.FromField("CascadeRelease")},
			{Name: "local_bg", Type: proto.ColumnType_BOOL, Description: "Whether the EIP is a local bandwidth EIP.", Transform: transform.FromField("LocalBgp")},
			{Name: "anti_ddos_package_id", Type: proto.ColumnType_STRING, Description: "The Anti-DDoS service package ID, if the EIP is an Anti-DDoS EIP.", Transform: transform.FromField("AntiDDoSPackageId")},
			{Name: "renew_flag", Type: proto.ColumnType_STRING, Description: "Auto-renewal flag for monthly prepaid bandwidth EIPs (NOTIFY_AND_MANUAL_RENEW, NOTIFY_AND_AUTO_RENEW, DISABLE_NOTIFY_AND_MANUAL_RENEW).", Transform: transform.FromField("RenewFlag")},
			{Name: "eip_alg_type", Type: proto.ColumnType_JSON, Description: "EIP ALG protocol types.", Transform: transform.FromField("EipAlgType")},
			{Name: "deadline_date", Type: proto.ColumnType_TIMESTAMP, Description: "Prepaid monthly subscription bandwidth IP expiration time.", Transform: transform.FromField("DeadlineDate").NullIfZero()},
			{Name: "created_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the EIP (ISO 8601, UTC).", Transform: transform.FromField("CreatedTime").NullIfZero()},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The tags attached to the EIP as a map of key -> value.", Transform: transform.FromField("Tags")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the EIP resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("AddressName")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listVpcEips(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listVpcEips", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listVpcEips", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	pageSize := utils.PageSizeFromLimit(d, 100, 1)
	var offset int64

	addressIds := utils.SingleIdFromQual(d.EqualsQuals, "address_id")
	filters := buildVpcEipFilters(d.EqualsQuals)

	for {
		d.WaitForListRateLimit(ctx)

		req := vpc.NewDescribeAddressesRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		if len(addressIds) > 0 {
			req.AddressIds = addressIds
		} else if len(filters) > 0 {
			req.Filters = filters
		}

		utils.LogRequest(ctx, "listVpcEips", req)
		resp, err := client.DescribeAddressesWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listVpcEips", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		addrs := resp.Response.AddressSet
		total := utils.PtrInt64(resp.Response.TotalCount)

		for _, item := range addrs {
			if item == nil {
				continue
			}
			row := toVpcEipRow(item)
			row.Region = region
			d.StreamListItem(ctx, row)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if utils.PageDone(offset, int64(len(addrs)), total) {
			break
		}
		offset += int64(len(addrs))
	}

	return nil, nil
}

func buildVpcEipFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*vpc.Filter {
	return utils.BuildFilters(equalsQuals, newVpcFilter,
		utils.FilterMapping{QualName: "address_name", FilterName: "address-name"},
		utils.FilterMapping{QualName: "address_status", FilterName: "address-status"},
		utils.FilterMapping{QualName: "instance_id", FilterName: "instance-id"},
		utils.FilterMapping{QualName: "instance_type", FilterName: "instance-type"},
		utils.FilterMapping{QualName: "network_interface_id", FilterName: "network-interface-id"},
	)
}

// Get Function

func getVpcEip(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	addressId := d.EqualsQuals["address_id"].GetStringValue()
	plugin.Logger(ctx).Debug("getVpcEip", "address_id", addressId, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getVpcEip", "client_init_error", err)
		return nil, err
	}

	if addressId == "" {
		return nil, nil
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	req := vpc.NewDescribeAddressesRequest()
	req.AddressIds = common.StringPtrs([]string{addressId})

	utils.LogRequest(ctx, "getVpcEip", req)
	resp, err := client.DescribeAddressesWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getVpcEip", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.AddressSet) == 0 {
		return nil, nil
	}

	row := toVpcEipRow(resp.Response.AddressSet[0])
	row.Region = region
	return row, nil
}

// Row Type

type vpcEipRow struct {
	AddressId               string
	AddressName             string
	AddressStatus           string
	AddressIp               string
	AddressType             string
	InstanceId              string
	InstanceType            string
	NetworkInterfaceId      string
	PrivateAddressIp        string
	Bandwidth               uint64
	InternetChargeType      string
	InternetServiceProvider string
	Egress                  string
	IsArrears               bool
	IsBlocked               bool
	IsEipDirectConnection   bool
	CascadeRelease          bool
	LocalBgp                bool
	AntiDDoSPackageId       string
	RenewFlag               string
	EipAlgType              *vpc.AlgType
	DeadlineDate            *time.Time
	CreatedTime             *time.Time
	Tags                    map[string]string
	Region                  string
	Akas                    []string
}

func toVpcEipRow(a *vpc.Address) vpcEipRow {
	row := vpcEipRow{
		AddressId:               utils.PtrString(a.AddressId),
		AddressName:             utils.PtrString(a.AddressName),
		AddressStatus:           utils.PtrString(a.AddressStatus),
		AddressIp:               utils.PtrString(a.AddressIp),
		AddressType:             utils.PtrString(a.AddressType),
		InstanceId:              utils.PtrString(a.InstanceId),
		InstanceType:            utils.PtrString(a.InstanceType),
		NetworkInterfaceId:      utils.PtrString(a.NetworkInterfaceId),
		PrivateAddressIp:        utils.PtrString(a.PrivateAddressIp),
		Bandwidth:               utils.PtrUint64(a.Bandwidth),
		InternetChargeType:      utils.PtrString(a.InternetChargeType),
		InternetServiceProvider: utils.PtrString(a.InternetServiceProvider),
		Egress:                  utils.PtrString(a.Egress),
		IsArrears:               utils.PtrBool(a.IsArrears),
		IsBlocked:               utils.PtrBool(a.IsBlocked),
		IsEipDirectConnection:   utils.PtrBool(a.IsEipDirectConnection),
		CascadeRelease:          utils.PtrBool(a.CascadeRelease),
		LocalBgp:                utils.PtrBool(a.LocalBgp),
		AntiDDoSPackageId:       utils.PtrString(a.AntiDDoSPackageId),
		RenewFlag:               utils.PtrString(a.RenewFlag),
		EipAlgType:              a.EipAlgType,
		DeadlineDate:            utils.ParseTimestamp(a.DeadlineDate),
		CreatedTime:             utils.ParseTimestamp(a.CreatedTime),
		Tags:                    TagsToMap(a.TagSet),
	}

	if row.AddressId != "" {
		row.Akas = []string{"tencentcloud:vpc:eip:" + row.AddressId}
	}
	return row
}
