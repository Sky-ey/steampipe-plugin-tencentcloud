package cvm

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"

	cvm "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/cvm/v20170312"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudCvmZone() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_cvm_zone",
		Description:       "Tencent Cloud availability zones available to CVM, including the zone ID, description, and availability state.",
		GetMatrixItemFunc: buildCvmRegionList,
		List: &plugin.ListConfig{
			Hydrate: listCvmZones,
			Tags:    map[string]string{"service": "cvm", "action": "DescribeZones"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("zone"),
			Hydrate:    getCvmZone,
			Tags:       map[string]string{"service": "cvm", "action": "DescribeZones"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "zone", Type: proto.ColumnType_STRING, Description: "The ID of the availability zone (e.g., ap-guangzhou-3).", Transform: transform.FromField("Zone")},
			{Name: "zone_name", Type: proto.ColumnType_STRING, Description: "The description of the availability zone (e.g., 广州三区).", Transform: transform.FromField("ZoneName")},
			{Name: "zone_id", Type: proto.ColumnType_STRING, Description: "The numeric ID of the availability zone (e.g., 100003).", Transform: transform.FromField("ZoneId")},
			{Name: "zone_state", Type: proto.ColumnType_STRING, Description: "The availability state of the zone (e.g., AVAILABLE).", Transform: transform.FromField("ZoneState")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region the availability zone belongs to (e.g., ap-guangzhou).", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Zone")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listCvmZones(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listCvmZones", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &cvm.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listCvmZones", "client_init_error", err)
		return nil, err
	}

	d.WaitForListRateLimit(ctx)

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	req := cvm.NewDescribeZonesRequest()

	utils.LogRequest(ctx, "listCvmZones", req)
	resp, err := client.DescribeZonesWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("listCvmZones", "request_error", err)
		return nil, err
	}

	if resp == nil || resp.Response == nil {
		return nil, nil
	}
	zones := resp.Response.ZoneSet
	for _, zone := range zones {
		if zone == nil {
			continue
		}
		row := toCvmZoneRow(zone)
		row.Region = region
		d.StreamListItem(ctx, row)

		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

// Get Function

func getCvmZone(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	zone := d.EqualsQuals["zone"].GetStringValue()
	plugin.Logger(ctx).Debug("getCvmZone", "zone", zone, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &cvm.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getCvmZone", "client_init_error", err)
		return nil, err
	}

	if zone == "" {
		return nil, nil
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	req := cvm.NewDescribeZonesRequest()

	utils.LogRequest(ctx, "getCvmZone", req)
	resp, err := client.DescribeZonesWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getCvmZone", "request_error", err)
		return nil, err
	}

	if resp == nil || resp.Response == nil {
		return nil, nil
	}
	for _, info := range resp.Response.ZoneSet {
		if info == nil {
			continue
		}
		if utils.PtrString(info.Zone) == zone {
			row := toCvmZoneRow(info)
			row.Region = region
			return row, nil
		}
	}

	return nil, nil
}

// Row Type

type cvmZoneRow struct {
	Zone      string
	ZoneName  string
	ZoneId    string
	ZoneState string
	Region    string
	Akas      []string
}

func toCvmZoneRow(info *cvm.ZoneInfo) cvmZoneRow {
	row := cvmZoneRow{
		Zone:      utils.PtrString(info.Zone),
		ZoneName:  utils.PtrString(info.ZoneName),
		ZoneId:    utils.PtrString(info.ZoneId),
		ZoneState: utils.PtrString(info.ZoneState),
	}

	if row.Zone != "" {
		row.Akas = []string{"tencentcloud:cvm:zone:" + row.Zone}
	}
	return row
}
