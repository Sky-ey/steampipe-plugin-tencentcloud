package ccn

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

func TableTencentcloudCcn() *plugin.Table {
	return &plugin.Table{
		Name:        "tencentcloud_ccn",
		Description: "Tencent Cloud CCN (Cloud Connect Network) instances that interconnect VPCs, Direct Connect gateways and BM VPCs across regions.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "ccn_id", Require: plugin.Optional},
				{Name: "ccn_name", Require: plugin.Optional},
				{Name: "state", Require: plugin.Optional},
			},
			Hydrate: listCcn,
			Tags:    map[string]string{"service": "ccn", "action": "DescribeCcns"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("ccn_id"),
			Hydrate:    getCcn,
			Tags:       map[string]string{"service": "ccn", "action": "DescribeCcns"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "ccn_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the CCN instance, such as `ccn-f49l6u0z`.", Transform: transform.FromField("CcnId")},
			{Name: "ccn_name", Type: proto.ColumnType_STRING, Description: "The name of the CCN instance.", Transform: transform.FromField("CcnName")},
			{Name: "ccn_description", Type: proto.ColumnType_STRING, Description: "The description of the CCN instance.", Transform: transform.FromField("CcnDescription")},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "The instance status. `ISOLATED`: isolated (arrears, service suspended); `AVAILABLE`: running.", Transform: transform.FromField("State")},
			{Name: "qos_level", Type: proto.ColumnType_STRING, Description: "The service quality. `PT`: Platinum, `AU`: Gold, `AG`: Silver.", Transform: transform.FromField("QosLevel")},
			{Name: "instance_charge_type", Type: proto.ColumnType_STRING, Description: "The billing method. `POSTPAID` indicates postpaid.", Transform: transform.FromField("InstanceChargeType")},
			{Name: "bandwidth_limit_type", Type: proto.ColumnType_STRING, Description: "The bandwidth limit type. `INTER_REGION_LIMIT`: between regions; `OUTER_REGION_LIMIT`: region egress limit.", Transform: transform.FromField("BandwidthLimitType")},
			{Name: "instance_count", Type: proto.ColumnType_INT, Description: "The number of associated instances.", Transform: transform.FromField("InstanceCount")},
			{Name: "create_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the CCN instance.", Transform: transform.FromField("CreateTime").NullIfZero()},
			{Name: "route_table_count", Type: proto.ColumnType_INT, Description: "The number of route tables associated with the CCN instance.", Transform: transform.FromField("RouteTableCount")},
			{Name: "route_table_flag", Type: proto.ColumnType_BOOL, Description: "Whether the multiple route tables feature is enabled for the CCN instance.", Transform: transform.FromField("RouteTableFlag")},
			{Name: "route_priority_flag", Type: proto.ColumnType_BOOL, Description: "Whether the CCN route priority feature is supported.", Transform: transform.FromField("RoutePriorityFlag")},
			{Name: "route_broadcast_policy_flag", Type: proto.ColumnType_BOOL, Description: "Whether the CCN route broadcasting policy is enabled.", Transform: transform.FromField("RouteBroadcastPolicyFlag")},
			{Name: "is_security_lock", Type: proto.ColumnType_BOOL, Description: "Whether the instance is blocked (traffic blocked).", Transform: transform.FromField("IsSecurityLock")},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The tags attached to the CCN instance as a map of key -> value.", Transform: transform.FromField("Tags")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listCcn(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listCcn", "quals", d.EqualsQuals)

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listCcn", "client_init_error", err)
		return nil, err
	}

	pageSize := uint64(utils.PageSizeFromLimit(d, 100, 1))
	var offset uint64

	ccnIds := utils.SingleIdFromQual(d.EqualsQuals, "ccn_id")
	filters := buildCcnFilters(d.EqualsQuals)

	for {
		d.WaitForListRateLimit(ctx)

		req := vpc.NewDescribeCcnsRequest()
		req.Offset = &offset
		req.Limit = &pageSize

		if len(ccnIds) > 0 {
			req.CcnIds = ccnIds
		} else if len(filters) > 0 {
			req.Filters = filters
		}

		utils.LogRequest(ctx, "listCcn", req)
		resp, err := client.DescribeCcnsWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listCcn", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		ccns := resp.Response.CcnSet
		total := utils.PtrUint64(resp.Response.TotalCount)

		for _, item := range ccns {
			if item == nil {
				continue
			}
			row := toCcnRow(item)
			d.StreamListItem(ctx, row)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if utils.PageDone(int64(offset), int64(len(ccns)), total) {
			break
		}
		offset += uint64(len(ccns))
	}

	return nil, nil
}

func buildCcnFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*vpc.Filter {
	return utils.BuildFilters(equalsQuals, newCcnFilter,
		utils.FilterMapping{QualName: "ccn_name", FilterName: "ccn-name"},
		utils.FilterMapping{QualName: "state", FilterName: "state"},
	)
}

// Get Function

func getCcn(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	ccnId := d.EqualsQuals["ccn_id"].GetStringValue()
	plugin.Logger(ctx).Info("getCcn", "ccn_id", ccnId)

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getCcn", "client_init_error", err)
		return nil, err
	}

	if ccnId == "" {
		return nil, nil
	}

	req := vpc.NewDescribeCcnsRequest()
	req.CcnIds = common.StringPtrs([]string{ccnId})

	utils.LogRequest(ctx, "getCcn", req)
	resp, err := client.DescribeCcnsWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getCcn", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.CcnSet) == 0 {
		return nil, nil
	}

	row := toCcnRow(resp.Response.CcnSet[0])
	return row, nil
}

// Row Type

type ccnRow struct {
	CcnId                    string
	CcnName                  string
	CcnDescription           string
	State                    string
	QosLevel                 string
	InstanceChargeType       string
	BandwidthLimitType       string
	InstanceCount            uint64
	CreateTime               *time.Time
	RouteTableCount          uint64
	RouteTableFlag           bool
	RoutePriorityFlag        bool
	RouteBroadcastPolicyFlag bool
	IsSecurityLock           bool
	Tags                     map[string]string
	Title                    string
	Akas                     []string
}

func toCcnRow(c *vpc.CCN) ccnRow {
	row := ccnRow{
		CcnId:                    utils.PtrString(c.CcnId),
		CcnName:                  utils.PtrString(c.CcnName),
		CcnDescription:           utils.PtrString(c.CcnDescription),
		State:                    utils.PtrString(c.State),
		QosLevel:                 utils.PtrString(c.QosLevel),
		InstanceChargeType:       utils.PtrString(c.InstanceChargeType),
		BandwidthLimitType:       utils.PtrString(c.BandwidthLimitType),
		InstanceCount:            utils.PtrUint64(c.InstanceCount),
		CreateTime:               utils.ParseTimestamp(c.CreateTime),
		RouteTableCount:          utils.PtrUint64(c.RouteTableCount),
		RouteTableFlag:           utils.PtrBool(c.RouteTableFlag),
		RoutePriorityFlag:        utils.PtrBool(c.RoutePriorityFlag),
		RouteBroadcastPolicyFlag: utils.PtrBool(c.RouteBroadcastPolicyFlag),
		IsSecurityLock:           utils.PtrBool(c.IsSecurityLock),
		Tags:                     vpcmod.TagsToMap(c.TagSet),
	}

	row.Title = row.CcnName
	if row.Title == "" {
		row.Title = row.CcnId
	}

	if row.CcnId != "" {
		row.Akas = []string{"tencentcloud:ccn:ccn:" + row.CcnId}
	}
	return row
}
