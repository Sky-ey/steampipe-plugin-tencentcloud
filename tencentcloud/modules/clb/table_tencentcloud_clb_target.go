package clb

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"strconv"
	"time"

	clb "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/clb/v20180317"
	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudClbTarget() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_clb_target",
		Description:       "Tencent Cloud Load Balancer (CLB) backend targets bound to listeners.",
		GetMatrixItemFunc: buildClbRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "load_balancer_id", Require: plugin.Optional},
				{Name: "listener_id", Require: plugin.Optional},
				{Name: "listener_protocol", Require: plugin.Optional},
				{Name: "listener_port", Require: plugin.Optional},
			},
			Hydrate: listClbTargets,
			Tags:    map[string]string{"service": "clb", "action": "DescribeTargets"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "load_balancer_id", Type: proto.ColumnType_STRING, Description: "The ID of the load balancer to which the target is bound.", Transform: transform.FromField("LoadBalancerId")},
			{Name: "listener_id", Type: proto.ColumnType_STRING, Description: "The ID of the listener to which the target is bound.", Transform: transform.FromField("ListenerId")},
			{Name: "listener_protocol", Type: proto.ColumnType_STRING, Description: "The protocol of the listener."},
			{Name: "listener_port", Type: proto.ColumnType_INT, Description: "The port of the listener."},
			{Name: "location_id", Type: proto.ColumnType_STRING, Description: "The ID of the forwarding rule, only set for HTTP/HTTPS listeners.", Transform: transform.FromField("LocationId")},
			{Name: "domain", Type: proto.ColumnType_STRING, Description: "The domain of the forwarding rule, only set for HTTP/HTTPS listeners."},
			{Name: "url", Type: proto.ColumnType_STRING, Description: "The URL of the forwarding rule, only set for HTTP/HTTPS listeners.", Transform: transform.FromField("Url")},
			{Name: "target_type", Type: proto.ColumnType_STRING, Description: "The type of the backend target (CVM, ENI, CCN, EVM, GLOBALROUTE, NAT, SRV)."},
			{Name: "instance_id", Type: proto.ColumnType_STRING, Description: "The ID of the backend CVM instance.", Transform: transform.FromField("InstanceId")},
			{Name: "instance_name", Type: proto.ColumnType_STRING, Description: "The name of the backend CVM instance."},
			{Name: "port", Type: proto.ColumnType_INT, Description: "The service port of the backend target."},
			{Name: "weight", Type: proto.ColumnType_INT, Description: "The forwarding weight of the backend target."},
			{Name: "public_ip_addresses", Type: proto.ColumnType_JSON, Description: "The public IP addresses of the backend target."},
			{Name: "private_ip_addresses", Type: proto.ColumnType_JSON, Description: "The private IP addresses of the backend target."},
			{Name: "eni_id", Type: proto.ColumnType_STRING, Description: "The ID of the elastic network interface of the backend target.", Transform: transform.FromField("EniId")},
			{Name: "zone", Type: proto.ColumnType_STRING, Description: "The availability zone of the backend target."},
			{Name: "tag", Type: proto.ColumnType_STRING, Description: "The tag of the backend target."},
			{Name: "registered_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the target was bound to the listener.", Transform: transform.FromField("RegisteredTime").NullIfZero()},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the target resides."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource."},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource."},
		}),
	}
}

// List Function

func listClbTargets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listClbTargets", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &clb.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listClbTargets", "client_init_error", err)
		return nil, err
	}
	region := d.EqualsQualString(utils.MatrixKeyRegion)

	loadBalancerIDs, err := resolveClbLoadBalancerIDs(ctx, d, client)
	if err != nil {
		plugin.Logger(ctx).Error("listClbTargets", "request_error", err)
		return nil, err
	}

	for _, loadBalancerID := range loadBalancerIDs {
		d.WaitForListRateLimit(ctx)

		req := clb.NewDescribeTargetsRequest()
		req.LoadBalancerId = &loadBalancerID
		if id := d.EqualsQualString("listener_id"); id != "" {
			req.ListenerIds = common.StringPtrs([]string{id})
		}
		if protocol := d.EqualsQualString("listener_protocol"); protocol != "" {
			req.Protocol = &protocol
		}
		if qual := d.EqualsQuals["listener_port"]; qual != nil {
			port := qual.GetInt64Value()
			req.Port = &port
		}

		utils.LogRequest(ctx, "listClbTargets", req)
		resp, err := client.DescribeTargetsWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listClbTargets", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			continue
		}

		for _, backend := range resp.Response.Listeners {
			if backend == nil {
				continue
			}
			listenerID := utils.PtrString(backend.ListenerId)
			protocol := utils.PtrString(backend.Protocol)
			listenerPort := utils.PtrInt64(backend.Port)

			// 四层监听器:后端服务直接挂在监听器下
			for _, target := range backend.Targets {
				if target == nil {
					continue
				}
				row := toClbTargetRow(target, loadBalancerID, listenerID, protocol, listenerPort, nil, region)
				d.StreamListItem(ctx, row)
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}

			// 七层监听器:后端服务挂在转发规则下
			for _, rule := range backend.Rules {
				if rule == nil {
					continue
				}
				for _, target := range rule.Targets {
					if target == nil {
						continue
					}
					row := toClbTargetRow(target, loadBalancerID, listenerID, protocol, listenerPort, rule, region)
					d.StreamListItem(ctx, row)
					if d.RowsRemaining(ctx) == 0 {
						return nil, nil
					}
				}
			}
		}
	}
	return nil, nil
}

// Row Type

type clbTargetRow struct {
	LoadBalancerId     string
	ListenerId         string
	ListenerProtocol   string
	ListenerPort       int64
	LocationId         string
	Domain             string
	Url                string
	TargetType         string
	InstanceId         string
	InstanceName       string
	Port               int64
	Weight             int64
	PublicIpAddresses  []string
	PrivateIpAddresses []string
	EniId              string
	Zone               string
	Tag                string
	RegisteredTime     *time.Time
	Region             string
	Title              string
	Akas               []string
}

func toClbTargetRow(target *clb.Backend, loadBalancerID, listenerID, protocol string, listenerPort int64, rule *clb.RuleTargets, region string) clbTargetRow {
	row := clbTargetRow{
		LoadBalancerId:     loadBalancerID,
		ListenerId:         listenerID,
		ListenerProtocol:   protocol,
		ListenerPort:       listenerPort,
		TargetType:         utils.PtrString(target.Type),
		InstanceId:         utils.PtrString(target.InstanceId),
		InstanceName:       utils.PtrString(target.InstanceName),
		Port:               utils.PtrInt64(target.Port),
		Weight:             utils.PtrInt64(target.Weight),
		PublicIpAddresses:  utils.PtrStringSlice(target.PublicIpAddresses),
		PrivateIpAddresses: utils.PtrStringSlice(target.PrivateIpAddresses),
		EniId:              utils.PtrString(target.EniId),
		Zone:               utils.PtrString(target.Zone),
		Tag:                utils.PtrString(target.Tag),
		RegisteredTime:     utils.ParseTimestamp(target.RegisteredTime),
		Region:             region,
	}
	if rule != nil {
		row.LocationId = utils.PtrString(rule.LocationId)
		row.Domain = utils.PtrString(rule.Domain)
		row.Url = utils.PtrString(rule.Url)
	}

	row.Title = row.InstanceId
	if row.Title == "" {
		row.Title = row.EniId
	}
	if row.Title == "" && len(row.PrivateIpAddresses) > 0 {
		row.Title = row.PrivateIpAddresses[0]
	}
	if row.Title == "" && len(row.PublicIpAddresses) > 0 {
		row.Title = row.PublicIpAddresses[0]
	}

	identity := row.InstanceId
	if identity == "" {
		identity = row.EniId
	}
	if identity == "" && len(row.PrivateIpAddresses) > 0 {
		identity = row.PrivateIpAddresses[0]
	}
	if identity != "" {
		aka := "tencentcloud:clb:target:" + row.LoadBalancerId + "/" + row.ListenerId
		if row.LocationId != "" {
			aka += "/" + row.LocationId
		}
		aka += "/" + identity + ":" + strconv.FormatInt(row.Port, 10)
		row.Akas = []string{aka}
	}
	return row
}
