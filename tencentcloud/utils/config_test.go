package utils

import (
	"testing"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func TestNormalizeRegion(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "ap-shanghai", expected: "ap-shanghai"},
		{input: "AP-Shanghai", expected: "ap-shanghai"},
		{input: "  ap-guangzhou  ", expected: "ap-guangzhou"},
		{input: "", expected: ""},
		{input: "   ", expected: ""},
	}

	for _, tt := range tests {
		if got := NormalizeRegion(tt.input); got != tt.expected {
			t.Fatalf("NormalizeRegion(%q) = %q, expected %q", tt.input, got, tt.expected)
		}
	}
}

func TestGetConfigNormalizesRegions(t *testing.T) {
	connection := &plugin.Connection{
		Name:   "tencentcloud",
		Config: TencentcloudConfig{Regions: []string{"AP-Shanghai", " ap-Guangzhou ", "AP-*"}},
	}

	got := GetConfig(connection).Regions
	expected := []string{"ap-shanghai", "ap-guangzhou", "ap-*"}

	if len(got) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
	for i := range expected {
		if got[i] != expected[i] {
			t.Fatalf("expected %v, got %v", expected, got)
		}
	}
}

func TestGetConfigKeepsNilRegions(t *testing.T) {
	t.Setenv("TENCENTCLOUD_REGION", "")
	connection := &plugin.Connection{
		Name:   "tencentcloud",
		Config: TencentcloudConfig{},
	}

	if got := GetConfig(connection).Regions; got != nil {
		t.Fatalf("expected nil regions when both the argument and environment variable are omitted, got %v", got)
	}
}

func TestGetConfigUsesRegionEnvironmentVariable(t *testing.T) {
	t.Setenv("TENCENTCLOUD_REGION", " AP-Shanghai ")
	connection := &plugin.Connection{
		Name:   "tencentcloud",
		Config: TencentcloudConfig{},
	}

	got := GetConfig(connection).Regions
	expected := []string{"ap-shanghai"}
	if len(got) != len(expected) || got[0] != expected[0] {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

func TestGetConfigRegionsTakePrecedenceOverEnvironment(t *testing.T) {
	t.Setenv("TENCENTCLOUD_REGION", "ap-shanghai")
	connection := &plugin.Connection{
		Name:   "tencentcloud",
		Config: TencentcloudConfig{Regions: []string{"ap-guangzhou"}},
	}

	got := GetConfig(connection).Regions
	expected := []string{"ap-guangzhou"}
	if len(got) != len(expected) || got[0] != expected[0] {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

// 显式写空数组与省略 regions 语义不同却结果雷同,属配置失误必须报错
func TestGetConfigPanicsOnEmptyRegions(t *testing.T) {
	tests := []struct {
		name    string
		regions []string
	}{
		{name: "empty slice", regions: []string{}},
		{name: "only blank values", regions: []string{"", "   "}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("expected a panic for an empty regions config")
				}
			}()

			GetConfig(&plugin.Connection{
				Name:   "tencentcloud",
				Config: TencentcloudConfig{Regions: tt.regions},
			})
		})
	}
}

// 非法通配符应在读取配置时立即报错
func TestGetConfigPanicsOnInvalidRegionPattern(t *testing.T) {
	for _, pattern := range []string{"ap-[", "[", "ap-\\"} {
		t.Run(pattern, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatalf("expected a panic for the malformed pattern %q", pattern)
				}
			}()

			GetConfig(&plugin.Connection{
				Name:   "tencentcloud",
				Config: TencentcloudConfig{Regions: []string{pattern}},
			})
		})
	}
}

func TestConfigInstance(t *testing.T) {
	instance := ConfigInstance()
	cfg, ok := instance.(*TencentcloudConfig)
	if !ok {
		t.Fatalf("expected *TencentcloudConfig, got %T", instance)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config instance")
	}
}

func TestGetConfig(t *testing.T) {
	t.Setenv("TENCENTCLOUD_REGION", "")
	secretID := "sid"
	secretKey := "skey"
	endpoint := "cvm.tencentcloudapi.com"
	want := TencentcloudConfig{
		SecretId:  &secretID,
		SecretKey: &secretKey,
		Endpoint:  &endpoint,
	}

	tests := []struct {
		name       string
		connection *plugin.Connection
		wantEmpty  bool
	}{
		{
			name:       "nil connection returns empty config",
			connection: nil,
			wantEmpty:  true,
		},
		{
			name:       "nil config returns empty config",
			connection: &plugin.Connection{Config: nil},
			wantEmpty:  true,
		},
		{
			name:       "wrong config type returns empty config",
			connection: &plugin.Connection{Config: "not-a-config"},
			wantEmpty:  true,
		},
		{
			name:       "valid config is returned",
			connection: &plugin.Connection{Config: want},
			wantEmpty:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetConfig(tt.connection)
			if tt.wantEmpty {
				if got.SecretId != nil || got.SecretKey != nil || got.Endpoint != nil {
					t.Fatalf("expected empty config, got %+v", got)
				}
				return
			}
			if PtrString(got.SecretId) != secretID ||
				PtrString(got.SecretKey) != secretKey ||
				PtrString(got.Endpoint) != endpoint {
				t.Fatalf("expected %+v, got %+v", want, got)
			}
		})
	}
}
