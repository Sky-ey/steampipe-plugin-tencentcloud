package dc

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	dc "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/dc/v20180410"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudDcDirectConnectTunnel() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_dc_direct_connect_tunnel",
		Description:       "Tencent Cloud Direct Connect dedicated tunnels that connect a Direct Connect to a VPC, BMVPC, or CCN.",
		GetMatrixItemFunc: buildDcRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "direct_connect_tunnel_id", Require: plugin.Optional},
				{Name: "direct_connect_id", Require: plugin.Optional},
				{Name: "direct_connect_tunnel_name", Require: plugin.Optional},
			},
			Hydrate: listDcDirectConnectTunnels,
			Tags:    map[string]string{"service": "dc", "action": "DescribeDirectConnectTunnels"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("direct_connect_tunnel_id"),
			Hydrate:    getDcDirectConnectTunnel,
			Tags:       map[string]string{"service": "dc", "action": "DescribeDirectConnectTunnels"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "direct_connect_tunnel_id", Type: proto.ColumnType_STRING, Description: "The dedicated tunnel ID, such as `dcx-abcdefgh`.", Transform: transform.FromField("DirectConnectTunnelId")},
			{Name: "direct_connect_tunnel_name", Type: proto.ColumnType_STRING, Description: "The dedicated tunnel name.", Transform: transform.FromField("DirectConnectTunnelName")},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "The dedicated tunnel status (AVAILABLE, PENDING, ALLOCATING, ALLOCATED, ALTERING, DELETING, DELETED, COMFIRMING, REJECTED).", Transform: transform.FromField("State")},
			{Name: "direct_connect_id", Type: proto.ColumnType_STRING, Description: "The parent Direct Connect ID.", Transform: transform.FromField("DirectConnectId")},
			{Name: "direct_connect_owner_account", Type: proto.ColumnType_STRING, Description: "The Direct Connect owner account ID.", Transform: transform.FromField("DirectConnectOwnerAccount")},
			{Name: "owner_account", Type: proto.ColumnType_STRING, Description: "The dedicated tunnel owner account ID.", Transform: transform.FromField("OwnerAccount")},
			{Name: "network_type", Type: proto.ColumnType_STRING, Description: "The network type (VPC, BMVPC, CCN).", Transform: transform.FromField("NetworkType")},
			{Name: "network_region", Type: proto.ColumnType_STRING, Description: "The network region of the VPC, such as `ap-guangzhou`.", Transform: transform.FromField("NetworkRegion")},
			{Name: "vpc_id", Type: proto.ColumnType_STRING, Description: "The unified VPC or BMVPC ID.", Transform: transform.FromField("VpcId")},
			{Name: "vpc_name", Type: proto.ColumnType_STRING, Description: "The associated VPC name.", Transform: transform.FromField("VpcName")},
			{Name: "vpc_region", Type: proto.ColumnType_STRING, Description: "The VPC region abbreviation, such as `gz`.", Transform: transform.FromField("VpcRegion")},
			{Name: "direct_connect_gateway_id", Type: proto.ColumnType_STRING, Description: "The associated Direct Connect gateway ID.", Transform: transform.FromField("DirectConnectGatewayId")},
			{Name: "direct_connect_gateway_name", Type: proto.ColumnType_STRING, Description: "The associated Direct Connect gateway name.", Transform: transform.FromField("DirectConnectGatewayName")},
			{Name: "route_type", Type: proto.ColumnType_STRING, Description: "The routing type (BGP, STATIC).", Transform: transform.FromField("RouteType")},
			{Name: "vlan", Type: proto.ColumnType_INT, Description: "The VLAN of the dedicated tunnel.", Transform: transform.FromField("Vlan")},
			{Name: "bandwidth", Type: proto.ColumnType_INT, Description: "The bandwidth value of the dedicated tunnel.", Transform: transform.FromField("Bandwidth")},
			{Name: "tencent_address", Type: proto.ColumnType_STRING, Description: "The Tencent-side IP address.", Transform: transform.FromField("TencentAddress")},
			{Name: "tencent_backup_address", Type: proto.ColumnType_STRING, Description: "The backup IP address on the Tencent side.", Transform: transform.FromField("TencentBackupAddress")},
			{Name: "customer_address", Type: proto.ColumnType_STRING, Description: "The user-side IP address.", Transform: transform.FromField("CustomerAddress")},
			{Name: "nat_type", Type: proto.ColumnType_INT, Description: "Whether it is a NAT tunnel.", Transform: transform.FromField("NatType")},
			{Name: "bfd_enable", Type: proto.ColumnType_INT, Description: "Whether BFD is enabled.", Transform: transform.FromField("BfdEnable")},
			{Name: "enable_bgp_community", Type: proto.ColumnType_BOOL, Description: "Whether the BGP community attribute is enabled.", Transform: transform.FromField("EnableBGPCommunity")},
			{Name: "access_point_type", Type: proto.ColumnType_STRING, Description: "The access point type of the dedicated tunnel.", Transform: transform.FromField("AccessPointType")},
			{Name: "cloud_attach_id", Type: proto.ColumnType_STRING, Description: "The Cloud Attached Connection Service ID.", Transform: transform.FromField("CloudAttachId")},
			{Name: "sign_law", Type: proto.ColumnType_BOOL, Description: "Whether the connection associated with the tunnel has the service agreement signed.", Transform: transform.FromField("SignLaw")},
			{Name: "net_detect_id", Type: proto.ColumnType_STRING, Description: "The associated custom network probe ID.", Transform: transform.FromField("NetDetectId")},
			{Name: "bgp_peer", Type: proto.ColumnType_JSON, Description: "The user-side BGP configuration, including Asn and AuthKey.", Transform: transform.FromField("BgpPeer")},
			{Name: "route_filter_prefixes", Type: proto.ColumnType_JSON, Description: "The user-side IP ranges (route filter prefixes).", Transform: transform.FromField("RouteFilterPrefixes")},
			{Name: "created_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the dedicated tunnel.", Transform: transform.FromField("CreatedTime").NullIfZero()},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The tags attached to the dedicated tunnel as a map of key -> value.", Transform: transform.FromField("Tags")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the dedicated tunnel resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("DirectConnectTunnelName")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listDcDirectConnectTunnels(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listDcDirectConnectTunnels", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &dc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listDcDirectConnectTunnels", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	pageSize := utils.PageSizeFromLimit(d, 100, 1)
	var offset int64 = 0

	tunnelIds := utils.SingleIdFromQual(d.EqualsQuals, "direct_connect_tunnel_id")
	filters := utils.BuildFilters(d.EqualsQuals, newDcFilter,
		utils.FilterMapping{QualName: "direct_connect_tunnel_name", FilterName: "direct-connect-tunnel-name"},
		utils.FilterMapping{QualName: "direct_connect_tunnel_id", FilterName: "direct-connect-tunnel-id"},
		utils.FilterMapping{QualName: "direct_connect_id", FilterName: "direct-connect-id"},
	)

	for {
		d.WaitForListRateLimit(ctx)

		req := dc.NewDescribeDirectConnectTunnelsRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		if len(tunnelIds) > 0 {
			req.DirectConnectTunnelIds = tunnelIds
		} else if len(filters) > 0 {
			req.Filters = filters
		}

		utils.LogRequest(ctx, "listDcDirectConnectTunnels", req)
		resp, err := client.DescribeDirectConnectTunnelsWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listDcDirectConnectTunnels", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		items := resp.Response.DirectConnectTunnelSet
		total := utils.PtrInt64(resp.Response.TotalCount)

		for _, item := range items {
			if item == nil {
				continue
			}
			row := toDcDirectConnectTunnelRow(item)
			row.Region = region
			d.StreamListItem(ctx, row)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if utils.PageDone(offset, int64(len(items)), total) {
			break
		}
		offset += int64(len(items))
	}

	return nil, nil
}

// Get Function

func getDcDirectConnectTunnel(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	tunnelId := d.EqualsQuals["direct_connect_tunnel_id"].GetStringValue()
	plugin.Logger(ctx).Info("getDcDirectConnectTunnel", "direct_connect_tunnel_id", tunnelId, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &dc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getDcDirectConnectTunnel", "client_init_error", err)
		return nil, err
	}

	if tunnelId == "" {
		return nil, nil
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	req := dc.NewDescribeDirectConnectTunnelsRequest()
	req.DirectConnectTunnelIds = common.StringPtrs([]string{tunnelId})

	utils.LogRequest(ctx, "getDcDirectConnectTunnel", req)
	resp, err := client.DescribeDirectConnectTunnelsWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getDcDirectConnectTunnel", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.DirectConnectTunnelSet) == 0 {
		return nil, nil
	}

	row := toDcDirectConnectTunnelRow(resp.Response.DirectConnectTunnelSet[0])
	row.Region = region
	return row, nil
}

// Row Type

type dcDirectConnectTunnelRow struct {
	DirectConnectTunnelId     string
	DirectConnectTunnelName   string
	State                     string
	DirectConnectId           string
	DirectConnectOwnerAccount string
	OwnerAccount              string
	NetworkType               string
	NetworkRegion             string
	VpcId                     string
	VpcName                   string
	VpcRegion                 string
	DirectConnectGatewayId    string
	DirectConnectGatewayName  string
	RouteType                 string
	Vlan                      int64
	Bandwidth                 int64
	TencentAddress            string
	TencentBackupAddress      string
	CustomerAddress           string
	NatType                   int64
	BfdEnable                 int64
	EnableBGPCommunity        bool
	AccessPointType           string
	CloudAttachId             string
	SignLaw                   bool
	NetDetectId               string
	BgpPeer                   *dc.BgpPeer
	RouteFilterPrefixes       []*dc.RouteFilterPrefix
	CreatedTime               *time.Time
	Tags                      map[string]string
	Region                    string
	Akas                      []string
}

func toDcDirectConnectTunnelRow(s *dc.DirectConnectTunnel) dcDirectConnectTunnelRow {
	row := dcDirectConnectTunnelRow{
		DirectConnectTunnelId:     utils.PtrString(s.DirectConnectTunnelId),
		DirectConnectTunnelName:   utils.PtrString(s.DirectConnectTunnelName),
		State:                     utils.PtrString(s.State),
		DirectConnectId:           utils.PtrString(s.DirectConnectId),
		DirectConnectOwnerAccount: utils.PtrString(s.DirectConnectOwnerAccount),
		OwnerAccount:              utils.PtrString(s.OwnerAccount),
		NetworkType:               utils.PtrString(s.NetworkType),
		NetworkRegion:             utils.PtrString(s.NetworkRegion),
		VpcId:                     utils.PtrString(s.VpcId),
		VpcName:                   utils.PtrString(s.VpcName),
		VpcRegion:                 utils.PtrString(s.VpcRegion),
		DirectConnectGatewayId:    utils.PtrString(s.DirectConnectGatewayId),
		DirectConnectGatewayName:  utils.PtrString(s.DirectConnectGatewayName),
		RouteType:                 utils.PtrString(s.RouteType),
		Vlan:                      utils.PtrInt64(s.Vlan),
		Bandwidth:                 utils.PtrInt64(s.Bandwidth),
		TencentAddress:            utils.PtrString(s.TencentAddress),
		TencentBackupAddress:      utils.PtrString(s.TencentBackupAddress),
		CustomerAddress:           utils.PtrString(s.CustomerAddress),
		NatType:                   utils.PtrInt64(s.NatType),
		BfdEnable:                 utils.PtrInt64(s.BfdEnable),
		EnableBGPCommunity:        utils.PtrBool(s.EnableBGPCommunity),
		AccessPointType:           utils.PtrString(s.AccessPointType),
		CloudAttachId:             utils.PtrString(s.CloudAttachId),
		SignLaw:                   utils.PtrBool(s.SignLaw),
		NetDetectId:               utils.PtrString(s.NetDetectId),
		BgpPeer:                   s.BgpPeer,
		RouteFilterPrefixes:       s.RouteFilterPrefixes,
		CreatedTime:               utils.ParseTimestamp(s.CreatedTime),
		Tags:                      TagsToMap(s.TagSet),
	}

	if row.DirectConnectTunnelId != "" {
		row.Akas = []string{"tencentcloud:dc:direct-connect-tunnel:" + row.DirectConnectTunnelId}
	}
	return row
}
