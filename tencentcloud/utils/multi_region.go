package utils

import (
	"context"
	"fmt"
	"path"

	region "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/region/v20220627"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

const (
	MatrixKeyRegion = "region"
	DefaultRegion   = "ap-singapore"
)

// Get related region from query
func getQueryRegion(d *plugin.QueryData) string {
	if d.QueryContext == nil {
		return ""
	}

	regionQuals := plugin.NewKeyColumnQualValueMap(
		d.QueryContext.UnsafeQuals,
		plugin.KeyColumnSlice{{Name: MatrixKeyRegion, Operators: []string{"="}}},
	).ToEqualsQualValueMap()

	if reg, ok := regionQuals[MatrixKeyRegion]; ok {
		return reg.GetStringValue()
	}
	return ""
}

func BuildRegionList(ctx context.Context, d *plugin.QueryData, product string) []map[string]any {
	if reg := getQueryRegion(d); reg != "" {
		return []map[string]any{{MatrixKeyRegion: NormalizeRegion(reg)}}
	}

	patterns := GetConfig(d.Connection).Regions

	if len(patterns) == 0 {
		plugin.Logger(ctx).Warn("You are currently using the default region " + DefaultRegion + ". Instances in other regions will not be displayed.")
		return []map[string]any{{MatrixKeyRegion: DefaultRegion}}
	}

	regions := discoverRegions(ctx, d, product)

	matched := make([]string, 0, len(regions))
	for _, pattern := range patterns {
		for _, r := range regions {
			if ok, _ := path.Match(pattern, r); ok {
				matched = append(matched, r)
			}
		}
	}

	matrix := toRegionMatrix(matched)
	if len(matrix) == 0 {
		panic(fmt.Sprintf("The 'regions' config %v for connection %s does not match any region available to the %s product (available: %v). Edit your connection configuration file and then restart Steampipe", patterns, d.Connection.Name, product, regions))
	}

	plugin.Logger(ctx).Debug("BuildRegionList", "product", product, "patterns", patterns, "matrix", matrix)
	return matrix
}

// discoverRegions Describe product regions.
func discoverRegions(ctx context.Context, d *plugin.QueryData, product string) []string {
	cacheKey := fmt.Sprintf("tencentcloud-regions-%s", product)

	// Read cache
	if d.ConnectionManager != nil {
		if cached, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
			if cachedRegions, ok := cached.([]string); ok && len(cachedRegions) > 0 {
				plugin.Logger(ctx).Debug("discoverRegions: cache hit", "cache_key", cacheKey, "regions", cachedRegions)
				return cachedRegions
			}
		}
	}

	// Request for regions
	client := &region.Client{}
	if err := InitClientInRegion(ctx, d, client, DefaultRegion); err != nil {
		plugin.Logger(ctx).Error("BuildRegionList", "product", product, "stage", "init_client", "error", err)
		panic(fmt.Sprintf("Failed to discover regions for the %s product: %v", product, err))
	}

	req := region.NewDescribeRegionsRequest()
	req.Product = &product

	resp, err := client.DescribeRegionsWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("BuildRegionList", "product", product, "stage", "describe_regions", "error", err)
		panic(fmt.Sprintf("Failed to discover regions for the %s product: %v", product, err))
	}
	if resp == nil || resp.Response == nil {
		plugin.Logger(ctx).Error("BuildRegionList", "product", product, "stage", "describe_regions", "error", "empty response")
		panic(fmt.Sprintf("Failed to discover regions for the %s product: DescribeRegions returned an empty response", product))
	}

	// Parse region
	regions := make([]string, 0, len(resp.Response.RegionSet))
	for _, info := range resp.Response.RegionSet {
		if info == nil {
			continue
		}
		if r := NormalizeRegion(PtrString(info.Region)); r != "" {
			regions = append(regions, r)
		}
	}
	if len(regions) == 0 {
		plugin.Logger(ctx).Error("BuildRegionList", "product", product, "stage", "describe_regions", "error", "no valid regions")
		panic(fmt.Sprintf("Failed to discover regions for the %s product: the API returned no valid regions", product))
	}

	if d.ConnectionManager != nil {
		d.ConnectionManager.Cache.Set(cacheKey, regions)
	}

	plugin.Logger(ctx).Debug("BuildRegionList", "product", product, "discovered_regions", regions)
	return regions
}

func toRegionMatrix(regions []string) []map[string]any {
	matrix := make([]map[string]any, 0, len(regions))
	seen := make(map[string]struct{}, len(regions))
	for _, r := range regions {
		if r == "" {
			continue
		}
		if _, ok := seen[r]; ok {
			continue
		}
		seen[r] = struct{}{}
		matrix = append(matrix, map[string]any{MatrixKeyRegion: r})
	}
	return matrix
}
