package utils

import (
	"context"
	"fmt"
	"strings"

	region "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/region/v20220627"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type regionRow struct {
	Region      string
	RegionName  string
	RegionState string
	Akas        []string
}

// TableRegion builds a region metadata table for the given service.
func TableRegion(service, product string) *plugin.Table {
	return &plugin.Table{
		Name:        fmt.Sprintf("tencentcloud_%s_region", service),
		Description: fmt.Sprintf("Tencent Cloud regions available to %s, including the region ID, description, and availability state.", strings.ToUpper(service)),
		List: &plugin.ListConfig{
			Hydrate: listRegions(product, service),
			Tags:    map[string]string{"service": service, "action": "DescribeRegions"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("region"),
			Hydrate:    getRegion(product, service),
			Tags:       map[string]string{"service": service, "action": "DescribeRegions"},
		},
		Columns: WithCommonColumns([]*plugin.Column{
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The ID of the region (e.g., ap-guangzhou).", Transform: transform.FromField("Region")},
			{Name: "region_name", Type: proto.ColumnType_STRING, Description: "The description of the region.", Transform: transform.FromField("RegionName")},
			{Name: "region_state", Type: proto.ColumnType_STRING, Description: "The availability state of the region (e.g., AVAILABLE).", Transform: transform.FromField("RegionState")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Region")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

func listRegions(product, service string) plugin.HydrateFunc {
	return func(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
		plugin.Logger(ctx).Debug("listRegions", "service", service, "product", product, "region", d.EqualsQualString(MatrixKeyRegion), "quals", d.EqualsQuals)

		client := &region.Client{}
		if err := InitClient(ctx, d, client); err != nil {
			plugin.Logger(ctx).Error("listRegions", "service", service, "client_init_error", err)
			return nil, err
		}

		d.WaitForListRateLimit(ctx)

		req := region.NewDescribeRegionsRequest()
		req.Product = &product
		LogRequest(ctx, "listRegions", req)
		resp, err := client.DescribeRegionsWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listRegions", "service", service, "request_error", err)
			return nil, err
		}

		if resp == nil || resp.Response == nil {
			return nil, nil
		}
		for _, info := range resp.Response.RegionSet {
			if info == nil {
				continue
			}
			d.StreamListItem(ctx, toRegionRow(info, service))
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		return nil, nil
	}
}

func getRegion(product, service string) plugin.HydrateFunc {
	return func(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
		regionQual := d.EqualsQualString("region")
		plugin.Logger(ctx).Debug("getRegion", "service", service, "region", regionQual)

		if regionQual == "" {
			return nil, nil
		}

		client := &region.Client{}
		if err := InitClient(ctx, d, client); err != nil {
			plugin.Logger(ctx).Error("getRegion", "service", service, "client_init_error", err)
			return nil, err
		}

		req := region.NewDescribeRegionsRequest()
		req.Product = &product
		LogRequest(ctx, "getRegion", req)
		resp, err := client.DescribeRegionsWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("getRegion", "service", service, "request_error", err)
			return nil, err
		}

		if resp == nil || resp.Response == nil {
			return nil, nil
		}
		for _, info := range resp.Response.RegionSet {
			if info == nil {
				continue
			}
			if PtrString(info.Region) == regionQual {
				return toRegionRow(info, service), nil
			}
		}
		return nil, nil
	}
}

func toRegionRow(info *region.RegionInfo, service string) regionRow {
	row := regionRow{
		Region:      PtrString(info.Region),
		RegionName:  PtrString(info.RegionName),
		RegionState: PtrString(info.RegionState),
	}
	if row.Region != "" {
		row.Akas = []string{fmt.Sprintf("tencentcloud:%s:region:%s", service, row.Region)}
	}
	return row
}
