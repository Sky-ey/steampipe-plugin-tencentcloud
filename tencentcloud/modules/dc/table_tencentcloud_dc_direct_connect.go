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

func TableTencentcloudDcDirectConnect() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_dc_direct_connect",
		Description:       "Tencent Cloud Direct Connect (DC) physical connections between an on-premises data center and Tencent Cloud.",
		GetMatrixItemFunc: buildDcRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "direct_connect_id", Require: plugin.Optional},
			},
			Hydrate: listDcDirectConnects,
			Tags:    map[string]string{"service": "dc", "action": "DescribeDirectConnects"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("direct_connect_id"),
			Hydrate:    getDcDirectConnect,
			Tags:       map[string]string{"service": "dc", "action": "DescribeDirectConnects"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "direct_connect_id", Type: proto.ColumnType_STRING, Description: "The Direct Connect instance ID, such as `dc-abcdefgh`.", Transform: transform.FromField("DirectConnectId")},
			{Name: "direct_connect_name", Type: proto.ColumnType_STRING, Description: "The connection name.", Transform: transform.FromField("DirectConnectName")},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "The connection status (PENDING, REJECTED, PAYPROCESSING, PAID, ALLOCATED, ALLOCATION, JAIL, AVAILABLE, ADMINSTOP, DELETED).", Transform: transform.FromField("State")},
			{Name: "access_point_id", Type: proto.ColumnType_STRING, Description: "The access point ID of the connection.", Transform: transform.FromField("AccessPointId")},
			{Name: "access_point_type", Type: proto.ColumnType_STRING, Description: "The access point type of the connection.", Transform: transform.FromField("AccessPointType")},
			{Name: "line_operator", Type: proto.ColumnType_STRING, Description: "The ISP that provides the connection (ChinaTelecom, ChinaMobile, ChinaUnicom, In-houseWiring, ChinaOther, InternationalOperator).", Transform: transform.FromField("LineOperator")},
			{Name: "port_type", Type: proto.ColumnType_STRING, Description: "The user-side port type of the connection (100Base-T, 1000Base-T, 1000Base-LX, 10GBase-T, 10GBase-LR).", Transform: transform.FromField("PortType")},
			{Name: "circuit_code", Type: proto.ColumnType_STRING, Description: "The circuit code provided by the ISP.", Transform: transform.FromField("CircuitCode")},
			{Name: "redundant_direct_connect_id", Type: proto.ColumnType_STRING, Description: "The ID of the redundant connection.", Transform: transform.FromField("RedundantDirectConnectId")},
			{Name: "location", Type: proto.ColumnType_STRING, Description: "The location of the local IDC.", Transform: transform.FromField("Location")},
			{Name: "idc_city", Type: proto.ColumnType_STRING, Description: "The IDC city.", Transform: transform.FromField("IdcCity")},
			{Name: "bandwidth", Type: proto.ColumnType_INT, Description: "The connection port bandwidth in Mbps.", Transform: transform.FromField("Bandwidth")},
			{Name: "min_bandwidth", Type: proto.ColumnType_INT, Description: "The minimum bandwidth of the connection.", Transform: transform.FromField("MinBandwidth")},
			{Name: "vlan", Type: proto.ColumnType_INT, Description: "The VLAN for connection debugging.", Transform: transform.FromField("Vlan")},
			{Name: "vlan_zero_direct_connect_tunnel_count", Type: proto.ColumnType_INT, Description: "The number of dedicated tunnels with disabled VLAN in the connection.", Transform: transform.FromField("VlanZeroDirectConnectTunnelCount")},
			{Name: "other_vlan_direct_connect_tunnel_count", Type: proto.ColumnType_INT, Description: "The number of dedicated tunnels with enabled VLAN in the connection.", Transform: transform.FromField("OtherVlanDirectConnectTunnelCount")},
			{Name: "tencent_address", Type: proto.ColumnType_STRING, Description: "The Tencent-side IP address for connection debugging.", Transform: transform.FromField("TencentAddress")},
			{Name: "customer_address", Type: proto.ColumnType_STRING, Description: "The user-side IP address for connection debugging.", Transform: transform.FromField("CustomerAddress")},
			{Name: "charge_type", Type: proto.ColumnType_STRING, Description: "The connection billing mode (NON_RECURRING_CHARGE for one-time access charge).", Transform: transform.FromField("ChargeType")},
			{Name: "charge_state", Type: proto.ColumnType_STRING, Description: "The billing status of the connection.", Transform: transform.FromField("ChargeState")},
			{Name: "sign_law", Type: proto.ColumnType_BOOL, Description: "Whether the service agreement has been signed.", Transform: transform.FromField("SignLaw")},
			{Name: "local_zone", Type: proto.ColumnType_BOOL, Description: "Whether the connection is an edge zone connection.", Transform: transform.FromField("LocalZone")},
			{Name: "created_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the connection.", Transform: transform.FromField("CreatedTime").NullIfZero()},
			{Name: "enabled_time", Type: proto.ColumnType_TIMESTAMP, Description: "The activation time of the connection.", Transform: transform.FromField("EnabledTime").NullIfZero()},
			{Name: "expired_time", Type: proto.ColumnType_TIMESTAMP, Description: "The expiration time of the connection.", Transform: transform.FromField("ExpiredTime").NullIfZero()},
			{Name: "start_time", Type: proto.ColumnType_TIMESTAMP, Description: "The activation (start) time of the connection.", Transform: transform.FromField("StartTime").NullIfZero()},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The tags attached to the connection as a map of key -> value.", Transform: transform.FromField("Tags")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the connection resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("DirectConnectName")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listDcDirectConnects(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listDcDirectConnects", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &dc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listDcDirectConnects", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	pageSize := utils.PageSizeFromLimit(d, 100, 1)
	var offset int64 = 0

	directConnectIds := utils.SingleIdFromQual(d.EqualsQuals, "direct_connect_id")

	for {
		d.WaitForListRateLimit(ctx)

		req := dc.NewDescribeDirectConnectsRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		// Mode B: DirectConnectIds takes priority. The Filter names for this
		// API are not publicly documented, so we only expose DirectConnectIds.
		if len(directConnectIds) > 0 {
			req.DirectConnectIds = directConnectIds
		}

		utils.LogRequest(ctx, "listDcDirectConnects", req)
		resp, err := client.DescribeDirectConnectsWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listDcDirectConnects", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		items := resp.Response.DirectConnectSet
		total := utils.PtrInt64(resp.Response.TotalCount)

		for _, item := range items {
			if item == nil {
				continue
			}
			row := toDcDirectConnectRow(item)
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

func getDcDirectConnect(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	directConnectId := d.EqualsQuals["direct_connect_id"].GetStringValue()
	plugin.Logger(ctx).Info("getDcDirectConnect", "direct_connect_id", directConnectId, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &dc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getDcDirectConnect", "client_init_error", err)
		return nil, err
	}

	if directConnectId == "" {
		return nil, nil
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	req := dc.NewDescribeDirectConnectsRequest()
	req.DirectConnectIds = common.StringPtrs([]string{directConnectId})

	utils.LogRequest(ctx, "getDcDirectConnect", req)
	resp, err := client.DescribeDirectConnectsWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getDcDirectConnect", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.DirectConnectSet) == 0 {
		return nil, nil
	}

	row := toDcDirectConnectRow(resp.Response.DirectConnectSet[0])
	row.Region = region
	return row, nil
}

// Row Type

type dcDirectConnectRow struct {
	DirectConnectId                   string
	DirectConnectName                 string
	State                             string
	AccessPointId                     string
	AccessPointType                   string
	LineOperator                      string
	PortType                          string
	CircuitCode                       string
	RedundantDirectConnectId          string
	Location                          string
	IdcCity                           string
	Bandwidth                         int64
	MinBandwidth                      uint64
	Vlan                              int64
	VlanZeroDirectConnectTunnelCount  uint64
	OtherVlanDirectConnectTunnelCount uint64
	TencentAddress                    string
	CustomerAddress                   string
	ChargeType                        string
	ChargeState                       string
	SignLaw                           bool
	LocalZone                         bool
	CreatedTime                       *time.Time
	EnabledTime                       *time.Time
	ExpiredTime                       *time.Time
	StartTime                         *time.Time
	Tags                              map[string]string
	Region                            string
	Akas                              []string
}

func toDcDirectConnectRow(s *dc.DirectConnect) dcDirectConnectRow {
	row := dcDirectConnectRow{
		DirectConnectId:                   utils.PtrString(s.DirectConnectId),
		DirectConnectName:                 utils.PtrString(s.DirectConnectName),
		State:                             utils.PtrString(s.State),
		AccessPointId:                     utils.PtrString(s.AccessPointId),
		AccessPointType:                   utils.PtrString(s.AccessPointType),
		LineOperator:                      utils.PtrString(s.LineOperator),
		PortType:                          utils.PtrString(s.PortType),
		CircuitCode:                       utils.PtrString(s.CircuitCode),
		RedundantDirectConnectId:          utils.PtrString(s.RedundantDirectConnectId),
		Location:                          utils.PtrString(s.Location),
		IdcCity:                           utils.PtrString(s.IdcCity),
		Bandwidth:                         utils.PtrInt64(s.Bandwidth),
		MinBandwidth:                      utils.PtrUint64(s.MinBandwidth),
		Vlan:                              utils.PtrInt64(s.Vlan),
		VlanZeroDirectConnectTunnelCount:  utils.PtrUint64(s.VlanZeroDirectConnectTunnelCount),
		OtherVlanDirectConnectTunnelCount: utils.PtrUint64(s.OtherVlanDirectConnectTunnelCount),
		TencentAddress:                    utils.PtrString(s.TencentAddress),
		CustomerAddress:                   utils.PtrString(s.CustomerAddress),
		ChargeType:                        utils.PtrString(s.ChargeType),
		ChargeState:                       utils.PtrString(s.ChargeState),
		SignLaw:                           utils.PtrBool(s.SignLaw),
		LocalZone:                         utils.PtrBool(s.LocalZone),
		CreatedTime:                       utils.ParseTimestamp(s.CreatedTime),
		EnabledTime:                       utils.ParseTimestamp(s.EnabledTime),
		ExpiredTime:                       utils.ParseTimestamp(s.ExpiredTime),
		StartTime:                         utils.ParseTimestamp(s.StartTime),
		Tags:                              TagsToMap(s.TagSet),
	}

	if row.DirectConnectId != "" {
		row.Akas = []string{"tencentcloud:dc:direct-connect:" + row.DirectConnectId}
	}
	return row
}
