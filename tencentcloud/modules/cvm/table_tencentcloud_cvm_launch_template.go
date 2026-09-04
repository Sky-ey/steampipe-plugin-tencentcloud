package cvm

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/cvm/v20170312"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudCvmLaunchTemplate() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_cvm_launch_template",
		Description:       "Tencent Cloud CVM launch templates that store instance configuration parameters (image, instance type, security groups, etc.) for quickly creating instances with a predefined setup.",
		GetMatrixItemFunc: buildCvmRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "launch_template_id", Require: plugin.Optional},
				{Name: "launch_template_name", Require: plugin.Optional},
			},
			Hydrate: listCvmLaunchTemplates,
			Tags:    map[string]string{"service": "cvm", "action": "DescribeLaunchTemplates"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("launch_template_id"),
			Hydrate:    getCvmLaunchTemplate,
			Tags:       map[string]string{"service": "cvm", "action": "DescribeLaunchTemplates"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "launch_template_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the launch template (e.g., lt-xxxxxxxx).", Transform: transform.FromField("LaunchTemplateId")},
			{Name: "launch_template_name", Type: proto.ColumnType_STRING, Description: "The name of the launch template.", Transform: transform.FromField("LaunchTemplateName")},
			{Name: "default_version_number", Type: proto.ColumnType_INT, Description: "The default version number of the launch template.", Transform: transform.FromField("DefaultVersionNumber")},
			{Name: "latest_version_number", Type: proto.ColumnType_INT, Description: "The latest version number of the launch template.", Transform: transform.FromField("LatestVersionNumber")},
			{Name: "launch_template_version_count", Type: proto.ColumnType_INT, Description: "The total number of versions the launch template contains.", Transform: transform.FromField("LaunchTemplateVersionCount")},
			{Name: "created_by", Type: proto.ColumnType_STRING, Description: "The UIN of the account that created the launch template.", Transform: transform.FromField("CreatedBy")},
			{Name: "creation_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the launch template in ISO 8601 format (UTC).", Transform: transform.FromField("CreationTime").NullIfZero()},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the launch template resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("LaunchTemplateName")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listCvmLaunchTemplates(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listCvmLaunchTemplates", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &cvm.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listCvmLaunchTemplates", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	pageSize := utils.PageSizeFromLimit(d, 100, 1)
	var offset int64 = 0

	launchTemplateIds := utils.SingleIdFromQual(d.EqualsQuals, "launch_template_id")
	filters := buildCvmLaunchTemplateFilters(d.EqualsQuals)

	for {
		d.WaitForListRateLimit(ctx)

		req := cvm.NewDescribeLaunchTemplatesRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		if len(launchTemplateIds) > 0 {
			req.LaunchTemplateIds = launchTemplateIds
		} else {
			req.Filters = filters
		}

		utils.LogRequest(ctx, "listCvmLaunchTemplates", req)
		resp, err := client.DescribeLaunchTemplatesWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listCvmLaunchTemplates", "request_error", err)
			return nil, err
		}

		if resp == nil || resp.Response == nil {
			break
		}
		launchTemplates := resp.Response.LaunchTemplateSet
		total := utils.PtrInt64(resp.Response.TotalCount)

		for _, launchTemplate := range launchTemplates {
			if launchTemplate == nil {
				continue
			}
			row := toCvmLaunchTemplateRow(launchTemplate)
			row.Region = region
			d.StreamListItem(ctx, row)

			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if utils.PageDone(offset, int64(len(launchTemplates)), total) {
			break
		}
		offset += int64(len(launchTemplates))
	}

	return nil, nil
}

func buildCvmLaunchTemplateFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*cvm.Filter {
	var filters []*cvm.Filter

	if equalsQuals["launch_template_name"] != nil {
		if value := equalsQuals["launch_template_name"].GetStringValue(); value != "" {
			filters = append(filters, &cvm.Filter{
				Name:   common.StringPtr("LaunchTemplateName"),
				Values: common.StringPtrs([]string{value}),
			})
		}
	}

	return filters
}

// Get Function

func getCvmLaunchTemplate(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	launchTemplateId := d.EqualsQuals["launch_template_id"].GetStringValue()
	plugin.Logger(ctx).Info("getCvmLaunchTemplate", "launch_template_id", launchTemplateId, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &cvm.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getCvmLaunchTemplate", "client_init_error", err)
		return nil, err
	}

	if launchTemplateId == "" {
		return nil, nil
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	req := cvm.NewDescribeLaunchTemplatesRequest()
	req.LaunchTemplateIds = common.StringPtrs([]string{launchTemplateId})

	utils.LogRequest(ctx, "getCvmLaunchTemplate", req)
	resp, err := client.DescribeLaunchTemplatesWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getCvmLaunchTemplate", "request_error", err)
		return nil, err
	}

	if len(resp.Response.LaunchTemplateSet) == 0 {
		return nil, nil
	}

	if resp.Response == nil {
		return nil, nil
	}
	row := toCvmLaunchTemplateRow(resp.Response.LaunchTemplateSet[0])
	row.Region = region
	return row, nil
}

// Row Type

type cvmLaunchTemplateRow struct {
	LaunchTemplateId           string
	LaunchTemplateName         string
	DefaultVersionNumber       uint64
	LatestVersionNumber        uint64
	LaunchTemplateVersionCount uint64
	CreatedBy                  string
	CreationTime               *time.Time
	Region                     string
	Akas                       []string
}

func toCvmLaunchTemplateRow(launchTemplate *cvm.LaunchTemplateInfo) cvmLaunchTemplateRow {
	row := cvmLaunchTemplateRow{
		LaunchTemplateId:           utils.PtrString(launchTemplate.LaunchTemplateId),
		LaunchTemplateName:         utils.PtrString(launchTemplate.LaunchTemplateName),
		DefaultVersionNumber:       utils.PtrUint64(launchTemplate.DefaultVersionNumber),
		LatestVersionNumber:        utils.PtrUint64(launchTemplate.LatestVersionNumber),
		LaunchTemplateVersionCount: utils.PtrUint64(launchTemplate.LaunchTemplateVersionCount),
		CreatedBy:                  utils.PtrString(launchTemplate.CreatedBy),
		CreationTime:               utils.ParseTimestamp(launchTemplate.CreationTime),
	}

	if row.LaunchTemplateId != "" {
		row.Akas = []string{"tencentcloud:cvm:launch-template:" + row.LaunchTemplateId}
	}
	return row
}
