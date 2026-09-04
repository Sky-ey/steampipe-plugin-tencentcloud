package ccn

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"time"

	vpc "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/vpc/v20170312"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudCcnInstance() *plugin.Table {
	return &plugin.Table{
		Name:        "tencentcloud_ccn_instance",
		Description: "Tencent Cloud CCN associated instances (VPCs, Direct Connect gateways and BM VPCs) attached to a Cloud Connect Network.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "ccn_id", Require: plugin.Optional},
				{Name: "instance_type", Require: plugin.Optional},
				{Name: "instance_id", Require: plugin.Optional},
				{Name: "instance_region", Require: plugin.Optional},
			},
			Hydrate: listCcnInstance,
			Tags:    map[string]string{"service": "ccn", "action": "DescribeCcnAttachedInstances"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "ccn_id", Type: proto.ColumnType_STRING, Description: "The CCN instance ID that the associated instance belongs to.", Transform: transform.FromField("CcnId")},
			{Name: "instance_id", Type: proto.ColumnType_STRING, Description: "The ID of the associated instance.", Transform: transform.FromField("InstanceId")},
			{Name: "instance_type", Type: proto.ColumnType_STRING, Description: "The type of the associated instance. `VPC`, `DIRECTCONNECT`, or `BMVPC`.", Transform: transform.FromField("InstanceType")},
			{Name: "instance_name", Type: proto.ColumnType_STRING, Description: "The name of the associated instance.", Transform: transform.FromField("InstanceName")},
			{Name: "instance_region", Type: proto.ColumnType_STRING, Description: "The region to which the associated instance belongs, such as `ap-guangzhou`.", Transform: transform.FromField("InstanceRegion")},
			{Name: "instance_uin", Type: proto.ColumnType_STRING, Description: "The UIN (root account) to which the associated instance belongs.", Transform: transform.FromField("InstanceUin")},
			{Name: "ccn_uin", Type: proto.ColumnType_STRING, Description: "The UIN (root account) to which the CCN belongs.", Transform: transform.FromField("CcnUin")},
			{Name: "cidr_block", Type: proto.ColumnType_JSON, Description: "The CIDR blocks of the associated instance.", Transform: transform.FromField("CidrBlock")},
			{Name: "instance_area", Type: proto.ColumnType_STRING, Description: "The general location of the associated instance, such as `CHINA_MAINLAND`.", Transform: transform.FromField("InstanceArea")},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "The status of the associated instance. `PENDING`, `ACTIVE`, `EXPIRED`, `REJECTED`, `DELETED`, `FAILED`, `ATTACHING`, `DETACHING`, `DETACHFAILED`.", Transform: transform.FromField("State")},
			{Name: "attached_time", Type: proto.ColumnType_TIMESTAMP, Description: "The association time of the associated instance.", Transform: transform.FromField("AttachedTime").NullIfZero()},
			{Name: "route_table_id", Type: proto.ColumnType_STRING, Description: "The route table ID associated with the instance, if any.", Transform: transform.FromField("RouteTableId")},
			{Name: "route_table_name", Type: proto.ColumnType_STRING, Description: "The route table name associated with the instance, if any.", Transform: transform.FromField("RouteTableName")},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the associated instance.", Transform: transform.FromField("Description")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listCcnInstance(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listCcnInstance", "quals", d.EqualsQuals)

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listCcnInstance", "client_init_error", err)
		return nil, err
	}

	ccnIds, err := resolveCcnIDs(ctx, d, client)
	if err != nil {
		plugin.Logger(ctx).Error("listCcnInstance", "resolve_ccn_ids_error", err)
		return nil, err
	}

	filters := buildCcnInstanceFilters(d.EqualsQuals)

	for _, ccnId := range ccnIds {
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}

		pageSize := uint64(utils.PageSizeFromLimit(d, 100, 1))
		var offset uint64

		for {
			d.WaitForListRateLimit(ctx)

			req := vpc.NewDescribeCcnAttachedInstancesRequest()
			req.CcnId = &ccnId
			req.Offset = &offset
			req.Limit = &pageSize
			if len(filters) > 0 {
				req.Filters = filters
			}

			utils.LogRequest(ctx, "listCcnInstance", req)
			resp, err := client.DescribeCcnAttachedInstancesWithContext(ctx, req)
			if err != nil {
				plugin.Logger(ctx).Error("listCcnInstance", "request_error", err)
				return nil, err
			}
			if resp == nil || resp.Response == nil {
				break
			}

			instances := resp.Response.InstanceSet
			if resp.Response == nil {
				break
			}
			total := utils.PtrUint64(resp.Response.TotalCount)

			for _, item := range instances {
				if item == nil {
					continue
				}
				// CcnId is passed through from the parent query parameter; the
				// response field is not reliably populated.
				row := toCcnInstanceRow(item, ccnId)
				d.StreamListItem(ctx, row)
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}

			if utils.PageDone(int64(offset), int64(len(instances)), total) {
				break
			}
			offset += uint64(len(instances))
		}
	}

	return nil, nil
}

func buildCcnInstanceFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*vpc.Filter {
	return utils.BuildFilters(equalsQuals, newCcnFilter,
		utils.FilterMapping{QualName: "instance_type", FilterName: "instance-type"},
		utils.FilterMapping{QualName: "instance_id", FilterName: "instance-id"},
		utils.FilterMapping{QualName: "instance_region", FilterName: "instance-region"},
	)
}

// Row Type

type ccnInstanceRow struct {
	CcnId          string
	InstanceId     string
	InstanceType   string
	InstanceName   string
	InstanceRegion string
	InstanceUin    string
	CcnUin         string
	CidrBlock      []string
	InstanceArea   string
	State          string
	AttachedTime   *time.Time
	RouteTableId   string
	RouteTableName string
	Description    string
	Title          string
	Akas           []string
}

func toCcnInstanceRow(i *vpc.CcnAttachedInstance, ccnId string) ccnInstanceRow {
	row := ccnInstanceRow{
		CcnId:          ccnId,
		InstanceId:     utils.PtrString(i.InstanceId),
		InstanceType:   utils.PtrString(i.InstanceType),
		InstanceName:   utils.PtrString(i.InstanceName),
		InstanceRegion: utils.PtrString(i.InstanceRegion),
		InstanceUin:    utils.PtrString(i.InstanceUin),
		CcnUin:         utils.PtrString(i.CcnUin),
		CidrBlock:      utils.PtrStringSlice(i.CidrBlock),
		InstanceArea:   utils.PtrString(i.InstanceArea),
		State:          utils.PtrString(i.State),
		AttachedTime:   utils.ParseTimestamp(i.AttachedTime),
		RouteTableId:   utils.PtrString(i.RouteTableId),
		RouteTableName: utils.PtrString(i.RouteTableName),
		Description:    utils.PtrString(i.Description),
	}

	row.Title = row.InstanceName
	if row.Title == "" {
		row.Title = row.InstanceId
	}

	if row.CcnId != "" && row.InstanceType != "" && row.InstanceId != "" {
		row.Akas = []string{"tencentcloud:ccn:instance:" + row.CcnId + "/" + row.InstanceType + "/" + row.InstanceId}
	}
	return row
}
