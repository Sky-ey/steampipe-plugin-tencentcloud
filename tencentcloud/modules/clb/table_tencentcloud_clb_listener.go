package clb

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"time"

	clb "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/clb/v20180317"
	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudClbListener() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_clb_listener",
		Description:       "Tencent Cloud Load Balancer (CLB) listeners.",
		GetMatrixItemFunc: buildClbRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "load_balancer_id", Require: plugin.Optional},
				{Name: "listener_id", Require: plugin.Optional},
				{Name: "protocol", Require: plugin.Optional},
				{Name: "port", Require: plugin.Optional},
			},
			Hydrate: listClbListeners,
			Tags:    map[string]string{"service": "clb", "action": "DescribeListeners"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"load_balancer_id", "listener_id"}),
			Hydrate:    getClbListener,
			Tags:       map[string]string{"service": "clb", "action": "DescribeListeners"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "load_balancer_id", Type: proto.ColumnType_STRING, Description: "The ID of the load balancer to which the listener belongs.", Transform: transform.FromField("LoadBalancerId")},
			{Name: "listener_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the listener.", Transform: transform.FromField("ListenerId")},
			{Name: "listener_name", Type: proto.ColumnType_STRING, Description: "The name of the listener."},
			{Name: "protocol", Type: proto.ColumnType_STRING, Description: "The listener protocol: TCP, UDP, HTTP, HTTPS, TCP_SSL or QUIC."},
			{Name: "port", Type: proto.ColumnType_INT, Description: "The listening port."},
			{Name: "end_port", Type: proto.ColumnType_INT, Description: "The end port when the listener uses a port range.", Transform: transform.FromField("EndPort")},
			{Name: "scheduler", Type: proto.ColumnType_STRING, Description: "The scheduling algorithm: WRR, LEAST_CONN or IP_HASH."},
			{Name: "session_expire_time", Type: proto.ColumnType_INT, Description: "The session persistence timeout in seconds."},
			{Name: "session_type", Type: proto.ColumnType_STRING, Description: "The session persistence type: NORMAL or QUIC_CID."},
			{Name: "sni_switch", Type: proto.ColumnType_INT, Description: "Whether SNI is enabled: 1 (enabled), 0 (disabled).", Transform: transform.FromField("SniSwitch")},
			{Name: "target_type", Type: proto.ColumnType_STRING, Description: "The backend target type: NODE, TARGETGROUP or SCF."},
			{Name: "keepalive_enable", Type: proto.ColumnType_INT, Description: "Whether HTTP keepalive is enabled: 1 (enabled), 0 (disabled)."},
			{Name: "toa", Type: proto.ColumnType_BOOL, Description: "Whether TOA (TCP Option Address) is enabled."},
			{Name: "deregister_target_rst", Type: proto.ColumnType_BOOL, Description: "Whether sending RST to targets on deregistration is enabled."},
			{Name: "max_conn", Type: proto.ColumnType_INT, Description: "The maximum number of connections of the listener."},
			{Name: "max_cps", Type: proto.ColumnType_INT, Description: "The maximum number of new connections per second of the listener."},
			{Name: "idle_connect_timeout", Type: proto.ColumnType_INT, Description: "The idle connection timeout in seconds."},
			{Name: "reschedule_interval", Type: proto.ColumnType_INT, Description: "The rescheduling interval in seconds."},
			{Name: "reschedule_start_time", Type: proto.ColumnType_INT, Description: "The rescheduling start time."},
			{Name: "data_compress_mode", Type: proto.ColumnType_STRING, Description: "The data compression mode of the listener."},
			{Name: "create_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the listener was created.", Transform: transform.FromField("CreateTime").NullIfZero()},
			{Name: "certificate", Type: proto.ColumnType_JSON, Description: "The certificate information of the listener."},
			{Name: "health_check", Type: proto.ColumnType_JSON, Description: "The health check configuration of the listener."},
			{Name: "rules", Type: proto.ColumnType_JSON, Description: "The forwarding rules of the HTTP/HTTPS listener."},
			{Name: "target_group", Type: proto.ColumnType_JSON, Description: "The target group bound to the listener."},
			{Name: "target_group_list", Type: proto.ColumnType_JSON, Description: "The target groups bound to the listener with weights."},
			{Name: "attr_flags", Type: proto.ColumnType_JSON, Description: "The attribute flags of the listener."},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the listener resides."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource."},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource."},
		}),
	}
}

// List Function

func listClbListeners(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listClbListeners", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &clb.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listClbListeners", "client_init_error", err)
		return nil, err
	}
	region := d.EqualsQualString(utils.MatrixKeyRegion)

	loadBalancerIDs, err := resolveClbLoadBalancerIDs(ctx, d, client)
	if err != nil {
		plugin.Logger(ctx).Error("listClbListeners", "request_error", err)
		return nil, err
	}

	for _, loadBalancerID := range loadBalancerIDs {
		d.WaitForListRateLimit(ctx)

		req := clb.NewDescribeListenersRequest()
		req.LoadBalancerId = &loadBalancerID
		if id := d.EqualsQualString("listener_id"); id != "" {
			req.ListenerIds = common.StringPtrs([]string{id})
		}
		if protocol := d.EqualsQualString("protocol"); protocol != "" {
			req.Protocol = &protocol
		}
		if qual := d.EqualsQuals["port"]; qual != nil {
			port := qual.GetInt64Value()
			req.Port = &port
		}

		utils.LogRequest(ctx, "listClbListeners", req)
		resp, err := client.DescribeListenersWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listClbListeners", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			continue
		}

		for _, listener := range resp.Response.Listeners {
			if listener == nil {
				continue
			}
			row := toClbListenerRow(listener, loadBalancerID, region)
			d.StreamListItem(ctx, row)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, nil
}

// Get Function

func getClbListener(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	loadBalancerID := d.EqualsQualString("load_balancer_id")
	listenerID := d.EqualsQualString("listener_id")
	plugin.Logger(ctx).Info("getClbListener", "load_balancer_id", loadBalancerID, "listener_id", listenerID, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &clb.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getClbListener", "client_init_error", err)
		return nil, err
	}

	if loadBalancerID == "" || listenerID == "" {
		return nil, nil
	}
	req := clb.NewDescribeListenersRequest()
	req.LoadBalancerId = &loadBalancerID
	req.ListenerIds = common.StringPtrs([]string{listenerID})

	utils.LogRequest(ctx, "getClbListener", req)
	resp, err := client.DescribeListenersWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getClbListener", "request_error", err)
		return nil, err
	}
	if len(resp.Response.Listeners) == 0 || resp.Response.Listeners[0] == nil {
		return nil, nil
	}
	return toClbListenerRow(resp.Response.Listeners[0], loadBalancerID, d.EqualsQualString(utils.MatrixKeyRegion)), nil
}

// Row Type

type clbListenerRow struct {
	LoadBalancerId      string
	ListenerId          string
	ListenerName        string
	Protocol            string
	Port                int64
	EndPort             int64
	Scheduler           string
	SessionExpireTime   int64
	SessionType         string
	SniSwitch           int64
	TargetType          string
	KeepaliveEnable     int64
	Toa                 bool
	DeregisterTargetRst bool
	MaxConn             int64
	MaxCps              int64
	IdleConnectTimeout  int64
	RescheduleInterval  uint64
	RescheduleStartTime int64
	DataCompressMode    string
	CreateTime          *time.Time
	Certificate         *clb.CertificateOutput
	HealthCheck         *clb.HealthCheck
	Rules               []*clb.RuleOutput
	TargetGroup         *clb.BasicTargetGroupInfo
	TargetGroupList     []*clb.BasicTargetGroupInfo
	AttrFlags           []string
	Region              string
	Title               string
	Akas                []string
}

func toClbListenerRow(listener *clb.Listener, loadBalancerID, region string) clbListenerRow {
	row := clbListenerRow{
		LoadBalancerId:      loadBalancerID,
		ListenerId:          utils.PtrString(listener.ListenerId),
		ListenerName:        utils.PtrString(listener.ListenerName),
		Protocol:            utils.PtrString(listener.Protocol),
		Port:                utils.PtrInt64(listener.Port),
		EndPort:             utils.PtrInt64(listener.EndPort),
		Scheduler:           utils.PtrString(listener.Scheduler),
		SessionExpireTime:   utils.PtrInt64(listener.SessionExpireTime),
		SessionType:         utils.PtrString(listener.SessionType),
		SniSwitch:           utils.PtrInt64(listener.SniSwitch),
		TargetType:          utils.PtrString(listener.TargetType),
		KeepaliveEnable:     utils.PtrInt64(listener.KeepaliveEnable),
		Toa:                 utils.PtrBool(listener.Toa),
		DeregisterTargetRst: utils.PtrBool(listener.DeregisterTargetRst),
		MaxConn:             utils.PtrInt64(listener.MaxConn),
		MaxCps:              utils.PtrInt64(listener.MaxCps),
		IdleConnectTimeout:  utils.PtrInt64(listener.IdleConnectTimeout),
		RescheduleInterval:  utils.PtrUint64(listener.RescheduleInterval),
		RescheduleStartTime: utils.PtrInt64(listener.RescheduleStartTime),
		DataCompressMode:    utils.PtrString(listener.DataCompressMode),
		CreateTime:          utils.ParseTimestamp(listener.CreateTime),
		Certificate:         listener.Certificate,
		HealthCheck:         listener.HealthCheck,
		Rules:               listener.Rules,
		TargetGroup:         listener.TargetGroup,
		TargetGroupList:     listener.TargetGroupList,
		AttrFlags:           utils.PtrStringSlice(listener.AttrFlags),
		Region:              region,
	}
	row.Title = row.ListenerName
	if row.Title == "" {
		row.Title = row.ListenerId
	}
	if row.ListenerId != "" {
		row.Akas = []string{"tencentcloud:clb:listener:" + row.LoadBalancerId + "/" + row.ListenerId}
	}
	return row
}
