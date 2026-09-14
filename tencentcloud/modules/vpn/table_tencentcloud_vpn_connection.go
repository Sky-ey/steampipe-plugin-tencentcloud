package vpn

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/vpc/v20170312"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	vpcmod "github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/modules/vpc"
)

// Table Definition

func TableTencentcloudVpnConnection() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_vpn_connection",
		Description:       "Tencent Cloud VPN connections (IPsec tunnels) that bind a VPN gateway to a customer gateway.",
		GetMatrixItemFunc: buildVpnRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "vpn_connection_id", Require: plugin.Optional},
				{Name: "vpc_id", Require: plugin.Optional},
				{Name: "vpn_gateway_id", Require: plugin.Optional},
				{Name: "customer_gateway_id", Require: plugin.Optional},
				{Name: "vpn_connection_name", Require: plugin.Optional},
			},
			Hydrate: listVpnConnections,
			Tags:    map[string]string{"service": "vpn", "action": "DescribeVpnConnections"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("vpn_connection_id"),
			Hydrate:    getVpnConnection,
			Tags:       map[string]string{"service": "vpn", "action": "DescribeVpnConnections"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "vpn_connection_id", Type: proto.ColumnType_STRING, Description: "The VPN tunnel instance ID, such as `vpnx-f49l6u0z`.", Transform: transform.FromField("VpnConnectionId")},
			{Name: "vpn_connection_name", Type: proto.ColumnType_STRING, Description: "The VPN tunnel name.", Transform: transform.FromField("VpnConnectionName")},
			{Name: "vpc_id", Type: proto.ColumnType_STRING, Description: "The VPC instance ID associated with the tunnel.", Transform: transform.FromField("VpcId")},
			{Name: "vpn_gateway_id", Type: proto.ColumnType_STRING, Description: "The VPN gateway instance ID.", Transform: transform.FromField("VpnGatewayId")},
			{Name: "customer_gateway_id", Type: proto.ColumnType_STRING, Description: "The customer gateway instance ID.", Transform: transform.FromField("CustomerGatewayId")},
			{Name: "vpn_proto", Type: proto.ColumnType_STRING, Description: "The tunnel transmission protocol.", Transform: transform.FromField("VpnProto")},
			{Name: "encrypt_proto", Type: proto.ColumnType_STRING, Description: "The tunnel encryption protocol.", Transform: transform.FromField("EncryptProto")},
			{Name: "route_type", Type: proto.ColumnType_STRING, Description: "The route type of the tunnel.", Transform: transform.FromField("RouteType")},
			{Name: "negotiation_type", Type: proto.ColumnType_STRING, Description: "The negotiation type (main/aggresive).", Transform: transform.FromField("NegotiationType")},
			{Name: "pre_share_key", Type: proto.ColumnType_STRING, Description: "The pre-shared key (sensitive, may be hidden by IAM).", Transform: transform.FromField("PreShareKey")},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "The lifecycle state of the tunnel (PENDING, AVAILABLE, DELETING).", Transform: transform.FromField("State")},
			{Name: "net_status", Type: proto.ColumnType_STRING, Description: "The connectivity status of the tunnel (AVAILABLE means connected).", Transform: transform.FromField("NetStatus")},
			{Name: "enable_health_check", Type: proto.ColumnType_BOOL, Description: "Whether the health check is enabled.", Transform: transform.FromField("EnableHealthCheck")},
			{Name: "health_check_local_ip", Type: proto.ColumnType_STRING, Description: "The local IP address used for the health check.", Transform: transform.FromField("HealthCheckLocalIp")},
			{Name: "health_check_remote_ip", Type: proto.ColumnType_STRING, Description: "The peer IP address used for the health check.", Transform: transform.FromField("HealthCheckRemoteIp")},
			{Name: "health_check_status", Type: proto.ColumnType_STRING, Description: "The tunnel health check status (AVAILABLE, UNAVAILABLE).", Transform: transform.FromField("HealthCheckStatus")},
			{Name: "dpd_enable", Type: proto.ColumnType_INT, Description: "Whether DPD is enabled (0 disabled, 1 enabled).", Transform: transform.FromField("DpdEnable")},
			{Name: "dpd_timeout", Type: proto.ColumnType_STRING, Description: "The DPD timeout period.", Transform: transform.FromField("DpdTimeout")},
			{Name: "dpd_action", Type: proto.ColumnType_STRING, Description: "The action on DPD timeout (clear, restart).", Transform: transform.FromField("DpdAction")},
			{Name: "security_policy_database_set", Type: proto.ColumnType_JSON, Description: "The SPD (Security Policy Database) set of the tunnel.", Transform: transform.FromField("SecurityPolicyDatabaseSet")},
			{Name: "ike_options", Type: proto.ColumnType_JSON, Description: "The IKE options specification of the tunnel.", Transform: transform.FromField("IKEOptionsSpecification")},
			{Name: "ipsec_options", Type: proto.ColumnType_JSON, Description: "The IPSEC options specification of the tunnel.", Transform: transform.FromField("IPSECOptionsSpecification")},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The tags attached to the tunnel as a map of key -> value.", Transform: transform.FromField("Tags")},
			{Name: "created_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the tunnel.", Transform: transform.FromField("CreatedTime").NullIfZero()},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the tunnel resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listVpnConnections(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listVpnConnections", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listVpnConnections", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	pageSize := uint64(utils.PageSizeFromLimit(d, 100, 1))
	var offset uint64

	ids := utils.SingleIdFromQual(d.EqualsQuals, "vpn_connection_id")
	filters := buildVpnConnectionFilters(d.EqualsQuals)

	for {
		d.WaitForListRateLimit(ctx)

		req := vpc.NewDescribeVpnConnectionsRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		if len(ids) > 0 {
			req.VpnConnectionIds = ids
		} else if len(filters) > 0 {
			req.Filters = filters
		}

		utils.LogRequest(ctx, "listVpnConnections", req)
		resp, err := client.DescribeVpnConnectionsWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listVpnConnections", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		conns := resp.Response.VpnConnectionSet
		total := utils.PtrUint64(resp.Response.TotalCount)

		for _, c := range conns {
			if c == nil {
				continue
			}
			row := toVpnConnectionRow(c)
			row.Region = region
			d.StreamListItem(ctx, row)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if utils.PageDone(int64(offset), int64(len(conns)), total) {
			break
		}
		offset += uint64(len(conns))
	}

	return nil, nil
}

func buildVpnConnectionFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*vpc.Filter {
	return utils.BuildFilters(equalsQuals, newVpnFilter,
		utils.FilterMapping{QualName: "vpc_id", FilterName: "vpc-id"},
		utils.FilterMapping{QualName: "vpn_gateway_id", FilterName: "vpn-gateway-id"},
		utils.FilterMapping{QualName: "customer_gateway_id", FilterName: "customer-gateway-id"},
		utils.FilterMapping{QualName: "vpn_connection_name", FilterName: "vpn-connection-name"},
	)
}

// Get Function

func getVpnConnection(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	vpnConnectionId := d.EqualsQuals["vpn_connection_id"].GetStringValue()
	plugin.Logger(ctx).Debug("getVpnConnection", "vpn_connection_id", vpnConnectionId, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getVpnConnection", "client_init_error", err)
		return nil, err
	}

	if vpnConnectionId == "" {
		return nil, nil
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	req := vpc.NewDescribeVpnConnectionsRequest()
	req.VpnConnectionIds = common.StringPtrs([]string{vpnConnectionId})

	utils.LogRequest(ctx, "getVpnConnection", req)
	resp, err := client.DescribeVpnConnectionsWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getVpnConnection", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.VpnConnectionSet) == 0 {
		return nil, nil
	}

	row := toVpnConnectionRow(resp.Response.VpnConnectionSet[0])
	row.Region = region
	return row, nil
}

// Row Type

type vpnConnectionRow struct {
	VpnConnectionId           string
	VpnConnectionName         string
	VpcId                     string
	VpnGatewayId              string
	CustomerGatewayId         string
	VpnProto                  string
	EncryptProto              string
	RouteType                 string
	NegotiationType           string
	PreShareKey               string
	State                     string
	NetStatus                 string
	EnableHealthCheck         bool
	HealthCheckLocalIp        string
	HealthCheckRemoteIp       string
	HealthCheckStatus         string
	DpdEnable                 int64
	DpdTimeout                string
	DpdAction                 string
	SecurityPolicyDatabaseSet []*vpc.SecurityPolicyDatabase
	IKEOptionsSpecification   *vpc.IKEOptionsSpecification
	IPSECOptionsSpecification *vpc.IPSECOptionsSpecification
	Tags                      map[string]string
	CreatedTime               *time.Time
	Region                    string
	Title                     string
	Akas                      []string
}

func toVpnConnectionRow(c *vpc.VpnConnection) vpnConnectionRow {
	row := vpnConnectionRow{
		VpnConnectionId:           utils.PtrString(c.VpnConnectionId),
		VpnConnectionName:         utils.PtrString(c.VpnConnectionName),
		VpcId:                     utils.PtrString(c.VpcId),
		VpnGatewayId:              utils.PtrString(c.VpnGatewayId),
		CustomerGatewayId:         utils.PtrString(c.CustomerGatewayId),
		VpnProto:                  utils.PtrString(c.VpnProto),
		EncryptProto:              utils.PtrString(c.EncryptProto),
		RouteType:                 utils.PtrString(c.RouteType),
		NegotiationType:           utils.PtrString(c.NegotiationType),
		PreShareKey:               utils.PtrString(c.PreShareKey),
		State:                     utils.PtrString(c.State),
		NetStatus:                 utils.PtrString(c.NetStatus),
		EnableHealthCheck:         utils.PtrBool(c.EnableHealthCheck),
		HealthCheckLocalIp:        utils.PtrString(c.HealthCheckLocalIp),
		HealthCheckRemoteIp:       utils.PtrString(c.HealthCheckRemoteIp),
		HealthCheckStatus:         utils.PtrString(c.HealthCheckStatus),
		DpdEnable:                 utils.PtrInt64(c.DpdEnable),
		DpdTimeout:                utils.PtrString(c.DpdTimeout),
		DpdAction:                 utils.PtrString(c.DpdAction),
		SecurityPolicyDatabaseSet: c.SecurityPolicyDatabaseSet,
		IKEOptionsSpecification:   c.IKEOptionsSpecification,
		IPSECOptionsSpecification: c.IPSECOptionsSpecification,
		Tags:                      vpcmod.TagsToMap(c.TagSet),
		CreatedTime:               utils.ParseTimestamp(c.CreatedTime),
	}

	if row.VpnConnectionName != "" {
		row.Title = row.VpnConnectionName
	} else {
		row.Title = row.VpnConnectionId
	}

	if row.VpnConnectionId != "" {
		row.Akas = []string{"tencentcloud:vpn:connection:" + row.VpnConnectionId}
	}
	return row
}
