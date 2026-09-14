package utils

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/cvm/v20170312"
	"github.com/turbot/steampipe-plugin-sdk/v5/connection"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// Keep a compile-time check against a real generated service client. All other
// Tencent Cloud service clients use the same embedded common.Client pattern.
var _ Client = (*cvm.Client)(nil)

type recordingSDKClient struct {
	commonClient common.Client
	region       string
}

func (c *recordingSDKClient) Init(region string) *common.Client {
	c.region = region
	return c.commonClient.Init(region)
}

func TestInitClient(t *testing.T) {
	t.Setenv("TENCENTCLOUD_SECRET_ID", "")
	t.Setenv("TENCENTCLOUD_SECRET_KEY", "")

	secretID := "configured-secret-id"
	secretKey := "configured-secret-key"
	region := "ap-shanghai"
	endpoint := "cvm.internal.tencentcloudapi.com"
	d := queryDataWithConfig(TencentcloudConfig{
		SecretId:  &secretID,
		SecretKey: &secretKey,
		Endpoint:  &endpoint,
	})
	d.EqualsQuals = plugin.KeyColumnEqualsQualMap{
		MatrixKeyRegion: proto.NewQualValue(region),
	}
	client := &recordingSDKClient{}

	if err := InitClient(context.Background(), d, client); err != nil {
		t.Fatalf("initClient returned an error: %v", err)
	}
	if client.region != region {
		t.Fatalf("expected region %q, got %q", region, client.region)
	}
}

func TestInitClientUsesEnvironmentCredentials(t *testing.T) {
	t.Setenv("TENCENTCLOUD_SECRET_ID", "environment-secret-id")
	t.Setenv("TENCENTCLOUD_SECRET_KEY", "environment-secret-key")

	client := &recordingSDKClient{}
	if err := InitClient(context.Background(), queryDataWithConfig(TencentcloudConfig{}), client); err != nil {
		t.Fatalf("initClient returned an error: %v", err)
	}
	if client.region != "" {
		t.Fatalf("expected empty region for nearest-access routing, got %q", client.region)
	}
}

func TestSharedHTTPTransportReusesTransportPerConnection(t *testing.T) {
	d := queryDataWithConnectionCache(t, "shared-transport", TencentcloudConfig{})
	overridesA := map[string]string{
		"cvm.tencentcloudapi.com": "127.0.0.1",
		"*.tencentcloudapi.com":   "127.0.0.2",
	}
	overridesB := map[string]string{
		"*.tencentcloudapi.com":   "127.0.0.2",
		"cvm.tencentcloudapi.com": "127.0.0.1",
	}

	first := SharedHTTPTransport(t.Context(), d, overridesA, false)
	second := SharedHTTPTransport(t.Context(), d, overridesB, false)
	if first == nil || second == nil {
		t.Fatal("expected shared HTTP transports")
	}
	if first != second {
		t.Fatal("expected identical transport settings to reuse the connection transport")
	}
}

func TestSharedHTTPTransportIsolatesConnectionsAndSettings(t *testing.T) {
	firstConnection := queryDataWithConnectionCache(t, "first-connection", TencentcloudConfig{})
	secondConnection := queryDataWithConnectionCache(t, "second-connection", TencentcloudConfig{})

	first := SharedHTTPTransport(t.Context(), firstConnection, nil, false)
	second := SharedHTTPTransport(t.Context(), secondConnection, nil, false)
	insecure := SharedHTTPTransport(t.Context(), firstConnection, nil, true)

	if first == nil || second == nil || insecure == nil {
		t.Fatal("expected shared HTTP transports")
	}
	if first == second {
		t.Fatal("expected different connections to use isolated transports")
	}
	if first == insecure {
		t.Fatal("expected different TLS settings to use isolated transports")
	}
	if insecure.TLSClientConfig == nil || !insecure.TLSClientConfig.InsecureSkipVerify {
		t.Fatal("expected insecure transport to disable TLS certificate verification")
	}
	if got := first.IdleConnTimeout; got != HTTPTransportIdleTimeout {
		t.Fatalf("expected idle connection timeout %s, got %s", HTTPTransportIdleTimeout, got)
	}
}

func TestSharedHTTPTransportConcurrentInitialization(t *testing.T) {
	d := queryDataWithConnectionCache(t, "concurrent-transport", TencentcloudConfig{})
	const callers = 64

	start := make(chan struct{})
	transports := make(chan *http.Transport, callers)
	var wg sync.WaitGroup
	for range callers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			transports <- SharedHTTPTransport(t.Context(), d, nil, false)
		}()
	}
	close(start)
	wg.Wait()
	close(transports)

	var first *http.Transport
	for transport := range transports {
		if transport == nil {
			t.Fatal("expected shared HTTP transport")
		}
		if first == nil {
			first = transport
			continue
		}
		if first != transport {
			t.Fatal("expected concurrent initialization to return one shared transport")
		}
	}
}

func TestApplyEndpoint(t *testing.T) {
	tests := []struct {
		name           string
		endpoint       string
		expectedScheme string
		expectedHost   string
		wantError      bool
		allowInsecure  bool
	}{
		{
			name:           "hostname keeps default scheme",
			endpoint:       "cvm.internal.tencentcloudapi.com",
			expectedScheme: "HTTPS",
			expectedHost:   "cvm.internal.tencentcloudapi.com",
		},
		{
			name:           "https url accepted by default",
			endpoint:       "https://api.example.com/tencentcloud",
			expectedScheme: "HTTPS",
			expectedHost:   "api.example.com/tencentcloud",
		},
		{
			name:      "http url rejected by default",
			endpoint:  "http://127.0.0.1:8080/tencentcloud",
			wantError: true,
		},
		{
			name:           "http url allowed with allow_insecure_endpoint",
			endpoint:       "http://127.0.0.1:8080/tencentcloud",
			expectedScheme: "HTTP",
			expectedHost:   "127.0.0.1:8080/tencentcloud",
			allowInsecure:  true,
		},
		{
			name:      "reject unsupported scheme",
			endpoint:  "ftp://127.0.0.1:8080",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientProfile := profile.NewClientProfile()
			cfg := TencentcloudConfig{}
			if tt.allowInsecure {
				t := true
				cfg.AllowInsecureEndpoint = &t
			}
			err := applyEndpoint(testContext(), clientProfile, tt.endpoint, cfg)
			if tt.wantError {
				if err == nil {
					t.Fatal("expected endpoint validation error")
				}
				return
			}
			if err != nil {
				t.Fatalf("applyEndpoint returned an error: %v", err)
			}
			if clientProfile.HttpProfile.Scheme != tt.expectedScheme {
				t.Fatalf("expected scheme %q, got %q", tt.expectedScheme, clientProfile.HttpProfile.Scheme)
			}
			if clientProfile.HttpProfile.Endpoint != tt.expectedHost {
				t.Fatalf("expected endpoint %q, got %q", tt.expectedHost, clientProfile.HttpProfile.Endpoint)
			}
		})
	}
}

func TestApplyBaseUrl(t *testing.T) {
	tests := []struct {
		name             string
		endpoint         string
		baseUrl          string
		expectedRoot     string
		expectedEndpoint string
	}{
		{
			name:         "base_url sets root domain when endpoint empty",
			endpoint:     "",
			baseUrl:      "internal.example.com",
			expectedRoot: "internal.example.com",
		},
		{
			name:         "base_url with port keeps service prefix",
			endpoint:     "",
			baseUrl:      "internal.example.com:8080",
			expectedRoot: "internal.example.com:8080",
		},
		{
			name:             "endpoint takes precedence over base_url",
			endpoint:         "cvm.internal.tencentcloudapi.com",
			baseUrl:          "internal.example.com",
			expectedRoot:     "",
			expectedEndpoint: "cvm.internal.tencentcloudapi.com",
		},
		{
			name:         "empty base_url leaves root domain unset",
			endpoint:     "",
			baseUrl:      "",
			expectedRoot: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientProfile := profile.NewClientProfile()
			if err := applyEndpoint(testContext(), clientProfile, tt.endpoint, TencentcloudConfig{}); err != nil {
				t.Fatalf("applyEndpoint returned an error: %v", err)
			}
			if tt.endpoint == "" && tt.baseUrl != "" {
				clientProfile.HttpProfile.RootDomain = tt.baseUrl
			}
			if clientProfile.HttpProfile.RootDomain != tt.expectedRoot {
				t.Fatalf("expected root domain %q, got %q", tt.expectedRoot, clientProfile.HttpProfile.RootDomain)
			}
			if tt.expectedEndpoint != "" && clientProfile.HttpProfile.Endpoint != tt.expectedEndpoint {
				t.Fatalf("expected endpoint %q, got %q", tt.expectedEndpoint, clientProfile.HttpProfile.Endpoint)
			}
		})
	}
}

func TestGetConfigValidatesBaseUrl(t *testing.T) {
	t.Run("rejects base_url with scheme", func(t *testing.T) {
		bad := "http://example.com"
		conn := &plugin.Connection{Config: TencentcloudConfig{BaseUrl: &bad}}
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic for base_url with scheme")
			}
		}()
		GetConfig(conn)
	})

	t.Run("accepts bare domain base_url", func(t *testing.T) {
		good := "  internal.example.com  "
		conn := &plugin.Connection{Config: TencentcloudConfig{BaseUrl: &good}}
		cfg := GetConfig(conn)
		if cfg.BaseUrl == nil || *cfg.BaseUrl != "internal.example.com" {
			t.Fatalf("expected trimmed base_url, got %v", cfg.BaseUrl)
		}
	})

	t.Run("rejects empty base_url after trim", func(t *testing.T) {
		empty := "   "
		conn := &plugin.Connection{Config: TencentcloudConfig{BaseUrl: &empty}}
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic for empty base_url")
			}
		}()
		GetConfig(conn)
	})
}

func TestInitClientValidatesCredentials(t *testing.T) {
	t.Setenv("TENCENTCLOUD_SECRET_ID", "")
	t.Setenv("TENCENTCLOUD_SECRET_KEY", "")
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	t.Setenv("USERPROFILE", homeDir)

	tests := []struct {
		name   string
		config TencentcloudConfig
	}{
		{
			name:   "no credentials falls back to default chain",
			config: TencentcloudConfig{},
		},
		{
			name:   "partial credentials fall back to default chain",
			config: TencentcloudConfig{SecretId: ptr("configured-secret-id")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := InitClient(context.Background(), queryDataWithConfig(tt.config), &recordingSDKClient{})
			if err == nil {
				t.Fatalf("expected a missing-credentials error, got nil")
			}
			if !strings.Contains(err.Error(), "missing credentials") {
				t.Fatalf("expected error mentioning missing credentials, got %v", err)
			}
		})
	}
}

func TestResolveRoleArnFromEnvironment(t *testing.T) {
	t.Setenv("TENCENTCLOUD_ASSUME_ROLE_ARN", "  qcs::cam::uin/100000000001:roleName/environment-role  ")

	got := ResolveRoleArn(TencentcloudConfig{})
	if got != "qcs::cam::uin/100000000001:roleName/environment-role" {
		t.Fatalf("expected role ARN from environment, got %q", got)
	}
}

func TestResolveRoleArnConfigOverridesEnvironment(t *testing.T) {
	t.Setenv("TENCENTCLOUD_ASSUME_ROLE_ARN", "qcs::cam::uin/100000000001:roleName/environment-role")
	configured := "qcs::cam::uin/100000000002:roleName/configured-role"

	got := ResolveRoleArn(TencentcloudConfig{RoleArn: &configured})
	if got != configured {
		t.Fatalf("expected configured role ARN %q, got %q", configured, got)
	}
}

// 静态凭证 + STS token 应生成携带 token 的凭证(跨账号场景)
func TestResolveCredentialStaticToken(t *testing.T) {
	t.Setenv("TENCENTCLOUD_SECRET_ID", "")
	t.Setenv("TENCENTCLOUD_SECRET_KEY", "")
	t.Setenv("TENCENTCLOUD_SECURITY_TOKEN", "")

	sid := "static-id"
	skey := "static-key"
	tok := "static-token"
	cred, err := CreateCredential(TencentcloudConfig{
		SecretId:  &sid,
		SecretKey: &skey,
		Token:     &tok,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	gotID, gotKey, gotToken := cred.GetCredential()
	if gotID != sid || gotKey != skey || gotToken != tok {
		t.Fatalf("expected (%s,%s,%s), got (%s,%s,%s)", sid, skey, tok, gotID, gotKey, gotToken)
	}
}

// 环境变量 token 同样生效
func TestResolveCredentialEnvToken(t *testing.T) {
	t.Setenv("TENCENTCLOUD_SECRET_ID", "env-id")
	t.Setenv("TENCENTCLOUD_SECRET_KEY", "env-key")
	t.Setenv("TENCENTCLOUD_SECURITY_TOKEN", "env-token")

	cred, err := CreateCredential(TencentcloudConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	gotID, gotKey, gotToken := cred.GetCredential()
	if gotID != "env-id" || gotKey != "env-key" || gotToken != "env-token" {
		t.Fatalf("got (%s,%s,%s)", gotID, gotKey, gotToken)
	}
}

func TestResolveCredentialEnvTokenFallback(t *testing.T) {
	t.Setenv("TENCENTCLOUD_SECRET_ID", "env-id")
	t.Setenv("TENCENTCLOUD_SECRET_KEY", "env-key")
	t.Setenv("TENCENTCLOUD_SECURITY_TOKEN", "")
	t.Setenv("TENCENTCLOUD_TOKEN", "alt-env-token")

	cred, err := CreateCredential(TencentcloudConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	gotID, gotKey, gotToken := cred.GetCredential()
	if gotID != "env-id" || gotKey != "env-key" || gotToken != "alt-env-token" {
		t.Fatalf("got (%s,%s,%s)", gotID, gotKey, gotToken)
	}
}

func TestResolveCredentialEnvTokenPreferredOverFallback(t *testing.T) {
	t.Setenv("TENCENTCLOUD_SECRET_ID", "env-id")
	t.Setenv("TENCENTCLOUD_SECRET_KEY", "env-key")
	t.Setenv("TENCENTCLOUD_SECURITY_TOKEN", "env-token")
	t.Setenv("TENCENTCLOUD_TOKEN", "alt-env-token")

	cred, err := CreateCredential(TencentcloudConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	gotID, gotKey, gotToken := cred.GetCredential()
	if gotID != "env-id" || gotKey != "env-key" || gotToken != "env-token" {
		t.Fatalf("got (%s,%s,%s)", gotID, gotKey, gotToken)
	}
}

func TestResolveCredentialTCCLIFile(t *testing.T) {
	t.Setenv("TENCENTCLOUD_SECRET_ID", "")
	t.Setenv("TENCENTCLOUD_SECRET_KEY", "")
	t.Setenv("TENCENTCLOUD_SECURITY_TOKEN", "")
	writeTCCLICredentialFile(t, `{
		"secretId": "  tccli-id  ",
		"secretKey": "  tccli-key  ",
		"token": "  tccli-token  "
	}`)

	cred, err := CreateCredential(TencentcloudConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	gotID, gotKey, gotToken := cred.GetCredential()
	if gotID != "tccli-id" || gotKey != "tccli-key" || gotToken != "tccli-token" {
		t.Fatalf("expected TC CLI credentials, got (%s,%s,%s)", gotID, gotKey, gotToken)
	}
}

func TestResolveCredentialTCCLIFileErrors(t *testing.T) {
	t.Setenv("TENCENTCLOUD_SECRET_ID", "")
	t.Setenv("TENCENTCLOUD_SECRET_KEY", "")
	t.Setenv("TENCENTCLOUD_SECURITY_TOKEN", "")

	tests := []struct {
		name      string
		contents  string
		wantError string
	}{
		{
			name:      "invalid JSON",
			contents:  `{`,
			wantError: "parse TC CLI credential file",
		},
		{
			name:      "missing secret key",
			contents:  `{"secretId":"tccli-id"}`,
			wantError: `must contain non-empty "secretId" and "secretKey" fields`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeTCCLICredentialFile(t, tt.contents)
			_, err := CreateCredential(TencentcloudConfig{})
			if err == nil || !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("expected error containing %q, got %v", tt.wantError, err)
			}
		})
	}
}

// .spc 凭证优先于环境变量
func TestResolveCredentialSpcOverridesEnv(t *testing.T) {
	t.Setenv("TENCENTCLOUD_SECRET_ID", "env-id")
	t.Setenv("TENCENTCLOUD_SECRET_KEY", "env-key")

	sid := "spc-id"
	skey := "spc-key"
	cred, err := CreateCredential(TencentcloudConfig{SecretId: &sid, SecretKey: &skey})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	gotID, gotKey, _ := cred.GetCredential()
	if gotID != sid || gotKey != skey {
		t.Fatalf("expected spc credentials to win, got (%s,%s)", gotID, gotKey)
	}
}

func writeTCCLICredentialFile(t *testing.T, contents string) {
	t.Helper()
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	t.Setenv("USERPROFILE", homeDir)
	credentialDir := filepath.Join(homeDir, ".tccli")
	if err := os.MkdirAll(credentialDir, 0o700); err != nil {
		t.Fatalf("failed to create TC CLI credential directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(credentialDir, "default.credential"), []byte(contents), 0o600); err != nil {
		t.Fatalf("failed to write TC CLI credential file: %v", err)
	}
}

func queryDataWithConfig(config TencentcloudConfig) *plugin.QueryData {
	return &plugin.QueryData{
		Connection: &plugin.Connection{Config: config},
	}
}

func queryDataWithConnectionCache(t *testing.T, connectionName string, config TencentcloudConfig) *plugin.QueryData {
	t.Helper()
	cache, err := connection.NewConnectionCache(connectionName, 1<<20)
	if err != nil {
		t.Fatalf("failed to create connection cache: %v", err)
	}
	return &plugin.QueryData{
		Connection:      &plugin.Connection{Name: connectionName, Config: config},
		ConnectionCache: cache,
	}
}
