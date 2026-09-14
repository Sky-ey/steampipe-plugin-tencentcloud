package cbs

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"time"

	cbs "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/cbs/v20170312"
	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudCbsDiskStoragePool() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_cbs_disk_storage_pool",
		Description:       "Tencent Cloud Block Storage dedicated cluster disk storage pools.",
		GetMatrixItemFunc: buildCbsRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "cdc_id", Require: plugin.Optional},
				{Name: "zone", Require: plugin.Optional},
				{Name: "cage_id", Require: plugin.Optional},
				{Name: "disk_type", Require: plugin.Optional},
			},
			Hydrate: listCbsDiskStoragePools,
			Tags:    map[string]string{"service": "cbs", "action": "DescribeDiskStoragePool"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("cdc_id"),
			Hydrate:    getCbsDiskStoragePool,
			Tags:       map[string]string{"service": "cbs", "action": "DescribeDiskStoragePool"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "cdc_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the dedicated cluster (e.g., cdc-xxx).", Transform: transform.FromField("CdcId")},
			{Name: "cdc_name", Type: proto.ColumnType_STRING, Description: "The name of the dedicated cluster."},
			{Name: "cdc_state", Type: proto.ColumnType_STRING, Description: "The state of the dedicated cluster (NORMAL, CLOSED, FAULT, ISOLATED)."},
			{Name: "zone", Type: proto.ColumnType_STRING, Description: "The availability zone where the dedicated cluster resides."},
			{Name: "cage_id", Type: proto.ColumnType_STRING, Description: "The ID of the cage to which the dedicated cluster belongs."},
			{Name: "disk_type", Type: proto.ColumnType_STRING, Description: "The media type of the cloud disks in the cluster (CLOUD_BASIC, CLOUD_PREMIUM, CLOUD_SSD)."},
			{Name: "disk_total", Type: proto.ColumnType_INT, Description: "The total capacity of the dedicated cluster in GiB.", Transform: transform.FromField("DiskTotal")},
			{Name: "disk_available", Type: proto.ColumnType_INT, Description: "The available capacity of the dedicated cluster in GiB.", Transform: transform.FromField("DiskAvailable")},
			{Name: "disk_number", Type: proto.ColumnType_INT, Description: "The number of cloud disks created in the cluster."},
			{Name: "expired_time", Type: proto.ColumnType_TIMESTAMP, Description: "The expiry time of the dedicated cluster.", Transform: transform.FromField("ExpiredTime").NullIfZero()},
			{Name: "created_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the storage pool.", Transform: transform.FromField("CreatedTime").NullIfZero()},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the dedicated cluster resides."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("CdcName")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource."},
		}),
	}
}

// List Function

func listCbsDiskStoragePools(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listCbsDiskStoragePools", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &cbs.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listCbsDiskStoragePools", "client_init_error", err)
		return nil, err
	}

	cdcIds := utils.SingleIdFromQual(d.EqualsQuals, "cdc_id")
	filters := buildDiskStoragePoolFilters(d.EqualsQuals)
	region := d.EqualsQualString(utils.MatrixKeyRegion)
	pageSize := uint64(utils.PageSizeFromLimit(d, 100, 1))
	var offset uint64
	for {
		d.WaitForListRateLimit(ctx)
		req := cbs.NewDescribeDiskStoragePoolRequest()
		req.Offset, req.Limit = &offset, &pageSize
		if len(cdcIds) > 0 {
			req.CdcIds = cdcIds
		} else {
			req.Filters = filters
		}
		utils.LogRequest(ctx, "listCbsDiskStoragePools", req)
		resp, err := client.DescribeDiskStoragePoolWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listCbsDiskStoragePools", "request_error", err)
			return nil, err
		}
		items := diskStoragePoolSet(resp.Response)
		for _, item := range items {
			if item == nil {
				continue
			}
			d.StreamListItem(ctx, toCbsDiskStoragePoolRow(item, region))
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if utils.PageDone(int64(offset), int64(len(items)), utils.PtrUint64(resp.Response.TotalCount)) {
			break
		}
		offset += uint64(len(items))
	}
	return nil, nil
}

func buildDiskStoragePoolFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*cbs.Filter {
	return utils.BuildFilters(equalsQuals, newCbsFilter,
		utils.FilterMapping{QualName: "zone", FilterName: "zone"},
		utils.FilterMapping{QualName: "cage_id", FilterName: "cage-id"},
		utils.FilterMapping{QualName: "disk_type", FilterName: "disk-type"},
	)
}

// Get Function

func getCbsDiskStoragePool(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	id := d.EqualsQualString("cdc_id")
	plugin.Logger(ctx).Debug("getCbsDiskStoragePool", "cdc_id", id, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &cbs.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getCbsDiskStoragePool", "client_init_error", err)
		return nil, err
	}

	if id == "" {
		return nil, nil
	}
	req := cbs.NewDescribeDiskStoragePoolRequest()
	req.CdcIds = common.StringPtrs([]string{id})
	utils.LogRequest(ctx, "getCbsDiskStoragePool", req)
	resp, err := client.DescribeDiskStoragePoolWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getCbsDiskStoragePool", "request_error", err)
		return nil, err
	}
	items := diskStoragePoolSet(resp.Response)
	if len(items) == 0 || items[0] == nil {
		return nil, nil
	}
	return toCbsDiskStoragePoolRow(items[0], d.EqualsQualString(utils.MatrixKeyRegion)), nil
}

// Note: the response of DescribeDiskStoragePool carries the dedicated
// cluster list in both CdcSet and DiskStoragePoolSet; official examples
// populate CdcSet, so fall back to it when DiskStoragePoolSet is empty.

func diskStoragePoolSet(params *cbs.DescribeDiskStoragePoolResponseParams) []*cbs.Cdc {
	if len(params.DiskStoragePoolSet) > 0 {
		return params.DiskStoragePoolSet
	}
	return params.CdcSet
}

// Row Type

type cbsDiskStoragePoolRow struct {
	CdcId         string
	CdcName       string
	CdcState      string
	Zone          string
	CageId        string
	DiskType      string
	DiskTotal     uint64
	DiskAvailable uint64
	DiskNumber    uint64
	ExpiredTime   *time.Time
	CreatedTime   *time.Time
	Region        string
	Akas          []string
}

func toCbsDiskStoragePoolRow(item *cbs.Cdc, region string) cbsDiskStoragePoolRow {
	row := cbsDiskStoragePoolRow{
		CdcId:       utils.PtrString(item.CdcId),
		CdcName:     utils.PtrString(item.CdcName),
		CdcState:    utils.PtrString(item.CdcState),
		Zone:        utils.PtrString(item.Zone),
		CageId:      utils.PtrString(item.CageId),
		DiskType:    utils.PtrString(item.DiskType),
		DiskNumber:  utils.PtrUint64(item.DiskNumber),
		ExpiredTime: utils.ParseTimestamp(item.ExpiredTime),
		CreatedTime: utils.ParseTimestamp(item.CreatedTime),
		Region:      region,
	}
	if item.CdcResource != nil {
		row.DiskTotal = utils.PtrUint64(item.CdcResource.DiskTotal)
		row.DiskAvailable = utils.PtrUint64(item.CdcResource.DiskAvailable)
	}
	if row.CdcId != "" {
		row.Akas = []string{"tencentcloud:cbs:disk-storage-pool:" + row.CdcId}
	}
	return row
}
