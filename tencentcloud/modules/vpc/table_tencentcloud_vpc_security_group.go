package vpc

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"strconv"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/vpc/v20170312"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudVpcSecurityGroup() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_vpc_security_group",
		Description:       "Tencent Cloud VPC security groups that control inbound and outbound traffic for associated resources.",
		GetMatrixItemFunc: buildVpcRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "security_group_id", Require: plugin.Optional},
				{Name: "security_group_name", Require: plugin.Optional},
				{Name: "project_id", Require: plugin.Optional},
			},
			Hydrate: listVpcSecurityGroups,
			Tags:    map[string]string{"service": "vpc", "action": "DescribeSecurityGroups"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("security_group_id"),
			Hydrate:    getVpcSecurityGroup,
			Tags:       map[string]string{"service": "vpc", "action": "DescribeSecurityGroups"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "security_group_id", Type: proto.ColumnType_STRING, Description: "The security group instance ID, such as `sg-ohuuioma`.", Transform: transform.FromField("SecurityGroupId")},
			{Name: "security_group_name", Type: proto.ColumnType_STRING, Description: "The security group name (up to 60 characters).", Transform: transform.FromField("SecurityGroupName")},
			{Name: "security_group_desc", Type: proto.ColumnType_STRING, Description: "The remarks for the security group (up to 100 characters).", Transform: transform.FromField("SecurityGroupDesc")},
			{Name: "project_id", Type: proto.ColumnType_STRING, Description: "The project ID. 0 means the default project.", Transform: transform.FromField("ProjectId")},
			{Name: "is_default", Type: proto.ColumnType_BOOL, Description: "Whether it is the default security group (cannot be deleted).", Transform: transform.FromField("IsDefault")},
			{Name: "created_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the security group.", Transform: transform.FromField("CreatedTime").NullIfZero()},
			{Name: "update_time", Type: proto.ColumnType_TIMESTAMP, Description: "The last update time of the security group.", Transform: transform.FromField("UpdateTime").NullIfZero()},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The tags attached to the security group as a map of key -> value.", Transform: transform.FromField("Tags")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the security group resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("SecurityGroupName")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listVpcSecurityGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listVpcSecurityGroups", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listVpcSecurityGroups", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	pageSize := strconv.FormatInt(utils.PageSizeFromLimit(d, 100, 1), 10)
	offset := "0"

	sgIds := utils.SingleIdFromQual(d.EqualsQuals, "security_group_id")
	filters := buildVpcSecurityGroupFilters(d.EqualsQuals)

	for {
		d.WaitForListRateLimit(ctx)

		req := vpc.NewDescribeSecurityGroupsRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		if len(sgIds) > 0 {
			req.SecurityGroupIds = sgIds
		} else if len(filters) > 0 {
			req.Filters = filters
		}

		utils.LogRequest(ctx, "listVpcSecurityGroups", req)
		resp, err := client.DescribeSecurityGroupsWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listVpcSecurityGroups", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		groups := resp.Response.SecurityGroupSet
		total := utils.PtrUint64(resp.Response.TotalCount)
		currentOffset, _ := strconv.ParseUint(offset, 10, 64)

		for _, item := range groups {
			if item == nil {
				continue
			}
			row := toVpcSecurityGroupRow(item)
			row.Region = region
			d.StreamListItem(ctx, row)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if utils.PageDone(int64(currentOffset), int64(len(groups)), total) {
			break
		}
		offset = strconv.FormatUint(currentOffset+uint64(len(groups)), 10)
	}

	return nil, nil
}

func buildVpcSecurityGroupFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*vpc.Filter {
	filters := utils.BuildFilters(equalsQuals, newVpcFilter,
		utils.FilterMapping{QualName: "security_group_name", FilterName: "security-group-name"},
	)

	if equalsQuals["project_id"] != nil {
		value := equalsQuals["project_id"].GetInt64Value()
		filters = append(filters, newVpcFilter("project-id", strconv.FormatInt(value, 10)))
	}

	return filters
}

// Get Function

func getVpcSecurityGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	sgId := d.EqualsQuals["security_group_id"].GetStringValue()
	plugin.Logger(ctx).Debug("getVpcSecurityGroup", "security_group_id", sgId, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getVpcSecurityGroup", "client_init_error", err)
		return nil, err
	}

	if sgId == "" {
		return nil, nil
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	req := vpc.NewDescribeSecurityGroupsRequest()
	req.SecurityGroupIds = common.StringPtrs([]string{sgId})

	utils.LogRequest(ctx, "getVpcSecurityGroup", req)
	resp, err := client.DescribeSecurityGroupsWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getVpcSecurityGroup", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.SecurityGroupSet) == 0 {
		return nil, nil
	}

	row := toVpcSecurityGroupRow(resp.Response.SecurityGroupSet[0])
	row.Region = region
	return row, nil
}

// Row Type

type vpcSecurityGroupRow struct {
	SecurityGroupId   string
	SecurityGroupName string
	SecurityGroupDesc string
	ProjectId         string
	IsDefault         bool
	CreatedTime       *time.Time
	UpdateTime        *time.Time
	Tags              map[string]string
	Region            string
	Akas              []string
}

func toVpcSecurityGroupRow(sg *vpc.SecurityGroup) vpcSecurityGroupRow {
	row := vpcSecurityGroupRow{
		SecurityGroupId:   utils.PtrString(sg.SecurityGroupId),
		SecurityGroupName: utils.PtrString(sg.SecurityGroupName),
		SecurityGroupDesc: utils.PtrString(sg.SecurityGroupDesc),
		ProjectId:         utils.PtrString(sg.ProjectId),
		IsDefault:         utils.PtrBool(sg.IsDefault),
		CreatedTime:       utils.ParseTimestamp(sg.CreatedTime),
		UpdateTime:        utils.ParseTimestamp(sg.UpdateTime),
		Tags:              TagsToMap(sg.TagSet),
	}

	if row.SecurityGroupId != "" {
		row.Akas = []string{"tencentcloud:vpc:security-group:" + row.SecurityGroupId}
	}
	return row
}
