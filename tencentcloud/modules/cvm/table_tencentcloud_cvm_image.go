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

func TableTencentcloudCvmImage() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_cvm_image",
		Description:       "Tencent Cloud CVM images (public, private/custom and shared images) that can be used to launch instances.",
		GetMatrixItemFunc: buildCvmRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "image_id", Require: plugin.Optional},
				{Name: "image_type", Require: plugin.Optional},
				{Name: "image_name", Require: plugin.Optional},
				{Name: "platform", Require: plugin.Optional},
			},
			Hydrate: listCvmImages,
			Tags:    map[string]string{"service": "cvm", "action": "DescribeImages"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("image_id"),
			Hydrate:    getCvmImage,
			Tags:       map[string]string{"service": "cvm", "action": "DescribeImages"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "image_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the image (e.g., img-gvbnzy6f).", Transform: transform.FromField("ImageId")},
			{Name: "architecture", Type: proto.ColumnType_STRING, Description: "The architecture of the image (x86_64, arm, i386).", Transform: transform.FromField("Architecture")},
			{Name: "created_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the image in ISO 8601 format (UTC).", Transform: transform.FromField("CreatedTime").NullIfZero()},
			{Name: "image_creator", Type: proto.ColumnType_STRING, Description: "The creator of the image.", Transform: transform.FromField("ImageCreator")},
			{Name: "image_deprecated", Type: proto.ColumnType_BOOL, Description: "Indicates whether the image is deprecated.", Transform: transform.FromField("ImageDeprecated")},
			{Name: "image_description", Type: proto.ColumnType_STRING, Description: "The description of the image.", Transform: transform.FromField("ImageDescription")},
			{Name: "image_family", Type: proto.ColumnType_STRING, Description: "The image family the image belongs to.", Transform: transform.FromField("ImageFamily")},
			{Name: "image_name", Type: proto.ColumnType_STRING, Description: "The name of the image.", Transform: transform.FromField("ImageName")},
			{Name: "image_size", Type: proto.ColumnType_INT, Description: "The size of the image in GiB.", Transform: transform.FromField("ImageSize")},
			{Name: "image_source", Type: proto.ColumnType_STRING, Description: "The source of the image (OFFICIAL, CREATE_IMAGE, EXTERNAL_IMPORT).", Transform: transform.FromField("ImageSource")},
			{Name: "image_state", Type: proto.ColumnType_STRING, Description: "The state of the image (CREATING, NORMAL, CREATEFAILED, SYNCING, IMPORTING, IMPORTFAILED).", Transform: transform.FromField("ImageState")},
			{Name: "image_type", Type: proto.ColumnType_STRING, Description: "The type of the image (PUBLIC_IMAGE, PRIVATE_IMAGE, SHARED_IMAGE).", Transform: transform.FromField("ImageType")},
			{Name: "is_support_cloudinit", Type: proto.ColumnType_BOOL, Description: "Indicates whether the image supports cloud-init.", Transform: transform.FromField("IsSupportCloudinit")},
			{Name: "license_type", Type: proto.ColumnType_STRING, Description: "The license type of the image (TencentCloud, BYOL).", Transform: transform.FromField("LicenseType")},
			{Name: "os_name", Type: proto.ColumnType_STRING, Description: "The operating system name of the image.", Transform: transform.FromField("OsName")},
			{Name: "platform", Type: proto.ColumnType_STRING, Description: "The platform of the image (e.g., TencentOS, CentOS, Windows, Ubuntu, Debian, Fedora).", Transform: transform.FromField("Platform")},
			{Name: "snapshot_set", Type: proto.ColumnType_JSON, Description: "The snapshots associated with the image.", Transform: transform.FromField("SnapshotSet")},
			{Name: "sync_percent", Type: proto.ColumnType_INT, Description: "The sync progress percentage of the image. Only meaningful while the image is being copied.", Transform: transform.FromField("SyncPercent")},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The tags attached to the image as a map of key -> value.", Transform: transform.FromField("Tags")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the image resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("ImageName")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listCvmImages(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listCvmImages", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &cvm.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listCvmImages", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	pageSize := uint64(utils.PageSizeFromLimit(d, 100, 1))
	var offset uint64 = 0

	filters := buildImageFilters(d.EqualsQuals)

	for {
		d.WaitForListRateLimit(ctx)

		req := cvm.NewDescribeImagesRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		req.Filters = filters

		utils.LogRequest(ctx, "listCvmImages", req)
		resp, err := client.DescribeImagesWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listCvmImages", "request_error", err)
			return nil, err
		}

		if resp == nil || resp.Response == nil {
			break
		}
		images := resp.Response.ImageSet
		for _, image := range images {
			if image == nil {
				continue
			}
			row := toCvmImageRow(image)
			row.Region = region
			d.StreamListItem(ctx, row)

			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		total := utils.PtrInt64(resp.Response.TotalCount)
		if utils.PageDone(int64(offset), int64(len(images)), total) {
			break
		}
		offset += uint64(len(images))
	}

	return nil, nil
}

func buildImageFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*cvm.Filter {
	var filters []*cvm.Filter

	addFilter := func(qualName, filterName string) {
		if equalsQuals[qualName] == nil {
			return
		}
		value := equalsQuals[qualName].GetStringValue()
		if value == "" {
			return
		}
		filters = append(filters, &cvm.Filter{
			Name:   common.StringPtr(filterName),
			Values: common.StringPtrs([]string{value}),
		})
	}

	addFilter("image_id", "image-id")
	addFilter("image_type", "image-type")
	addFilter("image_name", "image-name")
	addFilter("platform", "platform")

	return filters
}

// Get Function

func getCvmImage(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	imageId := d.EqualsQuals["image_id"].GetStringValue()
	plugin.Logger(ctx).Info("getCvmImage", "image_id", imageId, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &cvm.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getCvmImage", "client_init_error", err)
		return nil, err
	}

	if imageId == "" {
		return nil, nil
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	req := cvm.NewDescribeImagesRequest()
	req.ImageIds = common.StringPtrs([]string{imageId})

	utils.LogRequest(ctx, "getCvmImage", req)
	resp, err := client.DescribeImagesWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getCvmImage", "request_error", err)
		return nil, err
	}

	if resp == nil || resp.Response == nil || len(resp.Response.ImageSet) == 0 {
		return nil, nil
	}
	row := toCvmImageRow(resp.Response.ImageSet[0])
	row.Region = region
	return row, nil
}

// Row Type

type cvmImageRow struct {
	ImageId            string
	OsName             string
	ImageType          string
	CreatedTime        *time.Time
	ImageName          string
	ImageDescription   string
	ImageSize          int64
	Architecture       string
	ImageState         string
	Platform           string
	ImageCreator       string
	ImageSource        string
	SyncPercent        int64
	IsSupportCloudinit bool
	SnapshotSet        []cvmImageSnapshotRow
	Tags               map[string]string
	LicenseType        string
	ImageFamily        string
	ImageDeprecated    bool
	Region             string
	Akas               []string
}

type cvmImageSnapshotRow struct {
	SnapshotId string
	DiskUsage  string
	DiskSize   int64
}

func toCvmImageRow(image *cvm.Image) cvmImageRow {
	row := cvmImageRow{
		ImageId:            utils.PtrString(image.ImageId),
		OsName:             utils.PtrString(image.OsName),
		ImageType:          utils.PtrString(image.ImageType),
		CreatedTime:        utils.ParseTimestamp(image.CreatedTime),
		ImageName:          utils.PtrString(image.ImageName),
		ImageDescription:   utils.PtrString(image.ImageDescription),
		ImageSize:          utils.PtrInt64(image.ImageSize),
		Architecture:       utils.PtrString(image.Architecture),
		ImageState:         utils.PtrString(image.ImageState),
		Platform:           utils.PtrString(image.Platform),
		ImageCreator:       utils.PtrString(image.ImageCreator),
		ImageSource:        utils.PtrString(image.ImageSource),
		SyncPercent:        utils.PtrInt64(image.SyncPercent),
		IsSupportCloudinit: utils.PtrBool(image.IsSupportCloudinit),
		LicenseType:        utils.PtrString(image.LicenseType),
		ImageFamily:        utils.PtrString(image.ImageFamily),
		ImageDeprecated:    utils.PtrBool(image.ImageDeprecated),
		Tags:               TagsToMap(image.Tags),
	}

	for _, snapshot := range image.SnapshotSet {
		if snapshot == nil {
			continue
		}
		row.SnapshotSet = append(row.SnapshotSet, cvmImageSnapshotRow{
			SnapshotId: utils.PtrString(snapshot.SnapshotId),
			DiskUsage:  utils.PtrString(snapshot.DiskUsage),
			DiskSize:   utils.PtrInt64(snapshot.DiskSize),
		})
	}

	if row.ImageId != "" {
		row.Akas = []string{"tencentcloud:cvm:image:" + row.ImageId}
	}
	return row
}
