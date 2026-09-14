package waf

import (
	"context"

	waf "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/waf/v20180125"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
)

// Table Definition

func TableTencentcloudWafObject() *plugin.Table {
	return &plugin.Table{
		Name:        "tencentcloud_waf_object",
		Description: "Tencent Cloud WAF protection objects (CLB listeners or cloud-native gateways bound to WAF).",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "object_id", Require: plugin.Optional},
				{Name: "instance_id", Require: plugin.Optional},
				{Name: "domain", Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
				{Name: "cls_status", Require: plugin.Optional},
			},
			Hydrate: listWafObjects,
			Tags:    map[string]string{"service": "waf", "action": "DescribeObjects"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("object_id"),
			Hydrate:    getWafObject,
			Tags:       map[string]string{"service": "waf", "action": "DescribeObjects"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "object_id", Type: proto.ColumnType_STRING, Description: "The protection object ID (CLB instance ID).", Transform: transform.FromField("ObjectId")},
			{Name: "object_name", Type: proto.ColumnType_STRING, Description: "The protection object name.", Transform: transform.FromField("ObjectName")},
			{Name: "instance_id", Type: proto.ColumnType_STRING, Description: "The WAF instance ID bound to the object.", Transform: transform.FromField("InstanceId")},
			{Name: "instance_name", Type: proto.ColumnType_STRING, Description: "The WAF instance name bound to the object.", Transform: transform.FromField("InstanceName")},
			{Name: "instance_level", Type: proto.ColumnType_INT, Description: "The WAF instance level; 0 if no instance is bound.", Transform: transform.FromField("InstanceLevel")},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The object type (CLB: Load Balancer, TSE: Cloud-native Gateway).", Transform: transform.FromField("Type")},
			{Name: "status", Type: proto.ColumnType_INT, Description: "The WAF protection switch status (0: disabled, 1: enabled).", Transform: transform.FromField("Status")},
			{Name: "cls_status", Type: proto.ColumnType_INT, Description: "The WAF log switch status (0: disabled, 1: enabled).", Transform: transform.FromField("ClsStatus")},
			{Name: "post_cls_status", Type: proto.ColumnType_INT, Description: "The CLB log shipping switch status (0: disabled, 1: enabled).", Transform: transform.FromField("PostCLSStatus")},
			{Name: "post_c_kafka_status", Type: proto.ColumnType_INT, Description: "The Kafka log shipping switch status (0: disabled, 1: enabled).", Transform: transform.FromField("PostCKafkaStatus")},
			{Name: "bot_status", Type: proto.ColumnType_INT, Description: "The Bot protection switch status (0: disabled, 1: enabled).", Transform: transform.FromField("BotStatus")},
			{Name: "api_status", Type: proto.ColumnType_INT, Description: "The API protection switch status (0: disabled, 1: enabled).", Transform: transform.FromField("ApiStatus")},
			{Name: "object_flow_mode", Type: proto.ColumnType_INT, Description: "The object access mode (0: image mode, 1: cleaning mode, 2: examination mode).", Transform: transform.FromField("ObjectFlowMode")},
			{Name: "proxy", Type: proto.ColumnType_INT, Description: "The proxy status (0: disabled, 1: use first XFF IP, 2: use remote_addr, 3: use header field from ip_headers).", Transform: transform.FromField("Proxy")},
			{Name: "ip_headers", Type: proto.ColumnType_JSON, Description: "The header fields for obtaining the client IP; effective when Proxy is 3.", Transform: transform.FromField("IpHeaders")},
			{Name: "virtual_domain", Type: proto.ColumnType_STRING, Description: "The virtual domain name corresponding to the CLB object.", Transform: transform.FromField("VirtualDomain")},
			{Name: "domain", Type: proto.ColumnType_STRING, Description: "The precise domain name filter used in the query; populated when specified.", Transform: transform.FromQual("domain")},
			{Name: "precise_domains", Type: proto.ColumnType_JSON, Description: "The precise domain list protected under the object.", Transform: transform.FromField("PreciseDomains")},
			{Name: "public_ip", Type: proto.ColumnType_JSON, Description: "The public network addresses of the object.", Transform: transform.FromField("PublicIp")},
			{Name: "private_ip", Type: proto.ColumnType_JSON, Description: "The private network addresses of the object.", Transform: transform.FromField("PrivateIp")},
			{Name: "vpc_id", Type: proto.ColumnType_STRING, Description: "The VPC ID of the object.", Transform: transform.FromField("Vpc")},
			{Name: "vpc_name", Type: proto.ColumnType_STRING, Description: "The VPC name of the object.", Transform: transform.FromField("VpcName")},
			{Name: "numerical_vpc_id", Type: proto.ColumnType_INT, Description: "The VPC ID in numerical format.", Transform: transform.FromField("NumericalVpcId")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the object resides (returned by the API as a business field, not a matrix column).", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("ObjectId")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listWafObjects(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listWafObjects", "quals", d.EqualsQuals)

	client := &waf.Client{}
	if err := utils.InitClientInRegion(ctx, d, client, wafRequestRegion(d)); err != nil {
		plugin.Logger(ctx).Error("listWafObjects", "client_init_error", err)
		return nil, err
	}

	filters := buildWafObjectFilters(d.EqualsQuals)

	d.WaitForListRateLimit(ctx)
	req := waf.NewDescribeObjectsRequest()
	if len(filters) > 0 {
		req.Filters = filters
	}

	utils.LogRequest(ctx, "listWafObjects", req)
	resp, err := client.DescribeObjectsWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("listWafObjects", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil {
		return nil, nil
	}

	for _, item := range resp.Response.ClbObjects {
		if item == nil {
			continue
		}
		d.StreamListItem(ctx, toWafObjectRow(item))
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}

func buildWafObjectFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*waf.FiltersItemNew {
	return utils.BuildFilters(equalsQuals, newWafExactFilter,
		utils.FilterMapping{QualName: "object_id", FilterName: "ObjectId"},
		utils.FilterMapping{QualName: "instance_id", FilterName: "InstanceId"},
		utils.FilterMapping{QualName: "domain", FilterName: "Domain"},
		utils.FilterMapping{QualName: "status", FilterName: "Status"},
		utils.FilterMapping{QualName: "cls_status", FilterName: "ClsStatus"},
	)
}

// Get Function

func getWafObject(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	objectId := d.EqualsQuals["object_id"].GetStringValue()
	plugin.Logger(ctx).Debug("getWafObject", "object_id", objectId)

	client := &waf.Client{}
	if err := utils.InitClientInRegion(ctx, d, client, wafRequestRegion(d)); err != nil {
		plugin.Logger(ctx).Error("getWafObject", "client_init_error", err)
		return nil, err
	}

	if objectId == "" {
		return nil, nil
	}

	req := waf.NewDescribeObjectsRequest()
	req.Filters = []*waf.FiltersItemNew{
		newWafFilter("ObjectId", objectId, true),
	}

	utils.LogRequest(ctx, "getWafObject", req)
	resp, err := client.DescribeObjectsWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getWafObject", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.ClbObjects) == 0 {
		return nil, nil
	}

	return toWafObjectRow(resp.Response.ClbObjects[0]), nil
}

// Row Type

type wafObjectRow struct {
	ObjectId         string
	ObjectName       string
	InstanceId       string
	InstanceName     string
	InstanceLevel    int64
	Type             string
	Status           int64
	ClsStatus        int64
	PostCLSStatus    int64
	PostCKafkaStatus int64
	BotStatus        int64
	ApiStatus        int64
	ObjectFlowMode   int64
	Proxy            uint64
	IpHeaders        []*string
	VirtualDomain    string
	Domain           string
	PreciseDomains   []*string
	PublicIp         []*string
	PrivateIp        []*string
	Vpc              string
	VpcName          string
	NumericalVpcId   int64
	Region           string
	Akas             []string
}

func toWafObjectRow(o *waf.ClbObject) wafObjectRow {
	row := wafObjectRow{
		ObjectId:         utils.PtrString(o.ObjectId),
		ObjectName:       utils.PtrString(o.ObjectName),
		InstanceId:       utils.PtrString(o.InstanceId),
		InstanceName:     utils.PtrString(o.InstanceName),
		InstanceLevel:    utils.PtrInt64(o.InstanceLevel),
		Type:             utils.PtrString(o.Type),
		Status:           utils.PtrInt64(o.Status),
		ClsStatus:        utils.PtrInt64(o.ClsStatus),
		PostCLSStatus:    utils.PtrInt64(o.PostCLSStatus),
		PostCKafkaStatus: utils.PtrInt64(o.PostCKafkaStatus),
		BotStatus:        utils.PtrInt64(o.BotStatus),
		ApiStatus:        utils.PtrInt64(o.ApiStatus),
		ObjectFlowMode:   utils.PtrInt64(o.ObjectFlowMode),
		Proxy:            utils.PtrUint64(o.Proxy),
		IpHeaders:        o.IpHeaders,
		VirtualDomain:    utils.PtrString(o.VirtualDomain),
		PreciseDomains:   o.PreciseDomains,
		PublicIp:         o.PublicIp,
		PrivateIp:        o.PrivateIp,
		Vpc:              utils.PtrString(o.Vpc),
		VpcName:          utils.PtrString(o.VpcName),
		NumericalVpcId:   utils.PtrInt64(o.NumericalVpcId),
		Region:           utils.PtrString(o.Region),
	}

	if row.ObjectId != "" {
		row.Akas = []string{"tencentcloud:waf:object:" + row.ObjectId}
	}
	return row
}
