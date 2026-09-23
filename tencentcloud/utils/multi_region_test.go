package utils

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/turbot/steampipe-plugin-sdk/v5/connection"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/context_key"
)

func testContext() context.Context {
	logger := hclog.New(&hclog.LoggerOptions{Name: "test", Output: nil})
	return context.WithValue(context.Background(), context_key.Logger, logger)
}

func queryDataWithRegionQual(region string) *plugin.QueryData {
	d := queryDataWithConfig(TencentcloudConfig{})
	d.QueryContext = &plugin.QueryContext{
		UnsafeQuals: map[string]*proto.Quals{
			MatrixKeyRegion: {
				Quals: []*proto.Qual{
					{
						FieldName: MatrixKeyRegion,
						Operator:  &proto.Qual_StringValue{StringValue: "="},
						Value:     proto.NewQualValue(region),
					},
				},
			},
		},
	}
	return d
}

func newTestConnectionManager(t *testing.T) *connection.Manager {
	t.Helper()
	cache, err := connection.NewConnectionCache("test-connection", 1<<20)
	if err != nil {
		t.Fatalf("failed to create connection cache: %v", err)
	}
	return connection.NewManager(cache)
}

// queryDataWithCachedRegions 构造一个已缓存地域发现结果的 QueryData,
// 使 BuildRegionList 无需真实调用 DescribeRegions 即可走完过滤逻辑。
func queryDataWithCachedRegions(t *testing.T, config TencentcloudConfig, product string, regions []string) *plugin.QueryData {
	t.Helper()
	d := queryDataWithConfig(config)
	d.QueryContext = &plugin.QueryContext{}
	d.ConnectionManager = newTestConnectionManager(t)
	d.ConnectionManager.Cache.Set(fmt.Sprintf("tencentcloud-regions-%s", product), regions)
	return d
}

func matrixRegions(matrix []map[string]any) []string {
	out := make([]string, 0, len(matrix))
	for _, item := range matrix {
		out = append(out, item[MatrixKeyRegion].(string))
	}
	return out
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestGetQueryRegionReturnsProvidedRegion(t *testing.T) {
	d := queryDataWithRegionQual("ap-shanghai")

	if got := getQueryRegion(d); got != "ap-shanghai" {
		t.Fatalf("expected ap-shanghai, got %q", got)
	}
}

func TestGetQueryRegionEmptyWhenNoQual(t *testing.T) {
	d := queryDataWithConfig(TencentcloudConfig{})
	if got := getQueryRegion(d); got != "" {
		t.Fatalf("expected empty region, got %q", got)
	}
}

func TestGetQueryRegionEmptyWhenNoQueryContext(t *testing.T) {
	d := queryDataWithConfig(TencentcloudConfig{})
	d.QueryContext = nil
	if got := getQueryRegion(d); got != "" {
		t.Fatalf("expected empty region for nil QueryContext, got %q", got)
	}
}

// 查询里显式给出 region 时应短路,不做地域发现也不受 regions 配置影响
func TestBuildRegionListUsesProvidedRegion(t *testing.T) {
	d := queryDataWithRegionQual("ap-shanghai")

	matrix := BuildRegionList(t.Context(), d, "cvm")

	if got := matrixRegions(matrix); !equalStringSlices(got, []string{"ap-shanghai"}) {
		t.Fatalf("expected only ap-shanghai, got %v", got)
	}
}

func TestBuildRegionListNormalizesProvidedRegion(t *testing.T) {
	d := queryDataWithRegionQual("AP-Shanghai")

	matrix := BuildRegionList(t.Context(), d, "cvm")

	if got := matrixRegions(matrix); !equalStringSlices(got, []string{"ap-shanghai"}) {
		t.Fatalf("expected normalized ap-shanghai, got %v", got)
	}
}

func TestBuildRegionListUsesRegionEnvironmentVariable(t *testing.T) {
	t.Setenv("TENCENTCLOUD_REGION", "AP-Shanghai")
	d := queryDataWithCachedRegions(t,
		TencentcloudConfig{},
		"cvm",
		[]string{"ap-guangzhou", "ap-shanghai"},
	)

	matrix := BuildRegionList(testContext(), d, "cvm")

	if got := matrixRegions(matrix); !equalStringSlices(got, []string{"ap-shanghai"}) {
		t.Fatalf("expected only ap-shanghai, got %v", got)
	}
}

// Environment-only connections have nil Config in the SDK, not an empty config struct.
func TestBuildRegionListWithNilConfig(t *testing.T) {
	for _, tt := range []struct {
		name      string
		region    string
		want      string
		wantPanic string
	}{
		{name: "valid environment region", region: " AP-Shanghai ", want: "ap-shanghai"},
		{name: "nonexistent environment region", region: "ap-nonexistent", wantPanic: "does not match any region available to the cvm product"},
		{name: "malformed environment pattern", region: "ap-[", wantPanic: "invalid pattern"},
		{name: "empty environment defaults", want: DefaultRegion},
		{name: "blank environment defaults", region: "   ", want: DefaultRegion},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("TENCENTCLOUD_REGION", tt.region)
			d := queryDataWithCachedRegions(t, TencentcloudConfig{}, "cvm", []string{"ap-guangzhou", "ap-shanghai"})
			d.Connection.Config = nil

			if tt.wantPanic != "" {
				defer func() {
					r := recover()
					if r == nil || !strings.Contains(fmt.Sprint(r), tt.wantPanic) {
						t.Fatalf("expected panic containing %q, got %v", tt.wantPanic, r)
					}
				}()
			}

			matrix := BuildRegionList(testContext(), d, "cvm")
			if tt.wantPanic == "" {
				if got := matrixRegions(matrix); !equalStringSlices(got, []string{tt.want}) {
					t.Fatalf("expected [%s], got %v", tt.want, got)
				}
			}
		})
	}
}

// 未配置 regions 和 TENCENTCLOUD_REGION 时退化为兜底地域,且不应触发地域发现
func TestBuildRegionListDefaultsToSingapore(t *testing.T) {
	t.Setenv("TENCENTCLOUD_SECRET_ID", "")
	t.Setenv("TENCENTCLOUD_SECRET_KEY", "")
	t.Setenv("TENCENTCLOUD_REGION", "")

	d := queryDataWithConfig(TencentcloudConfig{})
	d.QueryContext = &plugin.QueryContext{}

	matrix := BuildRegionList(testContext(), d, "cvm")

	if got := matrixRegions(matrix); !equalStringSlices(got, []string{DefaultRegion}) {
		t.Fatalf("expected only %s, got %v", DefaultRegion, got)
	}
	if DefaultRegion != "ap-singapore" {
		t.Fatalf("expected the default region to be ap-singapore, got %q", DefaultRegion)
	}
}

func TestBuildRegionListFiltersDiscoveredRegionsByGlob(t *testing.T) {
	discovered := []string{
		"ap-guangzhou", "ap-shanghai", "ap-singapore", "ap-hongkong",
		"na-ashburn", "na-siliconvalley", "eu-frankfurt",
	}

	tests := []struct {
		name     string
		patterns []string
		expected []string
	}{
		{
			name:     "match all",
			patterns: []string{"*"},
			expected: discovered,
		},
		{
			name:     "prefix wildcard",
			patterns: []string{"ap-*"},
			expected: []string{"ap-guangzhou", "ap-shanghai", "ap-singapore", "ap-hongkong"},
		},
		{
			name:     "exact regions",
			patterns: []string{"ap-shanghai", "eu-frankfurt"},
			expected: []string{"ap-shanghai", "eu-frankfurt"},
		},
		{
			name:     "wildcard and exact mixed and deduplicated",
			patterns: []string{"na-*", "na-ashburn"},
			expected: []string{"na-ashburn", "na-siliconvalley"},
		},
		{
			name:     "single character wildcard",
			patterns: []string{"ap-?ongkong"},
			expected: []string{"ap-hongkong"},
		},
		{
			name:     "character class",
			patterns: []string{"ap-[gs]*"},
			expected: []string{"ap-guangzhou", "ap-shanghai", "ap-singapore"},
		},
		{
			name:     "config values are case insensitive",
			patterns: []string{"AP-Shanghai"},
			expected: []string{"ap-shanghai"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := queryDataWithCachedRegions(t, TencentcloudConfig{Regions: tt.patterns}, "cvm", discovered)
			// 走 GetConfig 的归一化路径,保证配置值大小写不敏感
			d.Connection.Config = TencentcloudConfig{Regions: tt.patterns}

			matrix := BuildRegionList(testContext(), d, "cvm")

			if got := matrixRegions(matrix); !equalStringSlices(got, tt.expected) {
				t.Fatalf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

// regions 只能收窄发现结果,不能扩展出云上不存在的地域
func TestBuildRegionListIgnoresRegionsNotAvailableToProduct(t *testing.T) {
	d := queryDataWithCachedRegions(t,
		TencentcloudConfig{Regions: []string{"ap-shanghai", "ap-nonexistent"}},
		"cvm",
		[]string{"ap-guangzhou", "ap-shanghai"},
	)

	matrix := BuildRegionList(testContext(), d, "cvm")

	if got := matrixRegions(matrix); !equalStringSlices(got, []string{"ap-shanghai"}) {
		t.Fatalf("expected only ap-shanghai, got %v", got)
	}
}

// 交集为空必须显式失败,否则用户只会看到「查询成功,0 行」
func TestBuildRegionListPanicsWhenNoRegionMatches(t *testing.T) {
	d := queryDataWithCachedRegions(t,
		TencentcloudConfig{Regions: []string{"us-*"}},
		"cvm",
		[]string{"ap-guangzhou", "ap-shanghai"},
	)

	defer func() {
		if recover() == nil {
			t.Fatal("expected a panic when the regions config matches no available region")
		}
	}()

	BuildRegionList(testContext(), d, "cvm")
}

// 非法通配符应在读取配置阶段就报错,不必等到地域发现成功
func TestBuildRegionListPanicsOnInvalidPattern(t *testing.T) {
	d := queryDataWithCachedRegions(t,
		TencentcloudConfig{Regions: []string{"ap-["}},
		"cvm",
		[]string{"ap-guangzhou"},
	)

	defer func() {
		if recover() == nil {
			t.Fatal("expected a panic for a malformed glob pattern")
		}
	}()

	BuildRegionList(testContext(), d, "cvm")
}

func TestBuildRegionListUsesCachedRegions(t *testing.T) {
	product := "cvm"
	cached := []string{"ap-guangzhou", "ap-shanghai", "ap-guangzhou", ""}
	d := queryDataWithCachedRegions(t, TencentcloudConfig{Regions: []string{"*"}}, product, cached)

	matrix := BuildRegionList(testContext(), d, product)

	if got := matrixRegions(matrix); !equalStringSlices(got, []string{"ap-guangzhou", "ap-shanghai"}) {
		t.Fatalf("expected deduplicated regions, got %v", got)
	}
}

// 配置了 regions 但发现失败时必须 panic,不能静默返回空矩阵
func TestBuildRegionListPanicsWhenDiscoveryFails(t *testing.T) {
	cleanCredentialEnvironment(t)

	d := queryDataWithConfig(TencentcloudConfig{Regions: []string{"*"}})
	d.QueryContext = &plugin.QueryContext{}
	d.ConnectionManager = newTestConnectionManager(t)

	defer func() {
		r := recover()
		if r == nil || !strings.Contains(fmt.Sprint(r), "missing credentials") {
			t.Fatalf("expected a missing-credentials panic when region discovery fails, got %v", r)
		}
	}()

	BuildRegionList(testContext(), d, "cvm")
}

func TestToRegionMatrix(t *testing.T) {
	matrix := toRegionMatrix([]string{"ap-shanghai", "", "ap-shanghai", "ap-guangzhou"})

	if got := matrixRegions(matrix); !equalStringSlices(got, []string{"ap-shanghai", "ap-guangzhou"}) {
		t.Fatalf("expected deduplicated non-empty regions, got %v", got)
	}
}
