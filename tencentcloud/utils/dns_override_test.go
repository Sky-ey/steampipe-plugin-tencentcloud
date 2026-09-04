package utils

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseDNSOverrides(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want map[string]string
	}{
		{
			name: "single entry",
			raw:  "clb.tencentcloudapi.com=127.0.0.1",
			want: map[string]string{"clb.tencentcloudapi.com": "127.0.0.1"},
		},
		{
			name: "multiple entries with wildcard",
			raw:  "cvm.tencentcloudapi.com=10.0.0.1, *.tencentcloudapi.com=10.0.0.5",
			want: map[string]string{
				"cvm.tencentcloudapi.com": "10.0.0.1",
				"*.tencentcloudapi.com":   "10.0.0.5",
			},
		},
		{
			name: "host normalized to lower case",
			raw:  "CLB.TencentCloudAPI.COM=127.0.0.1",
			want: map[string]string{"clb.tencentcloudapi.com": "127.0.0.1"},
		},
		{
			name: "empty raw string",
			raw:  "",
			want: map[string]string{},
		},
		{
			name: "malformed pairs are skipped",
			raw:  "no-equals-sign, =empty-host, host=, ok=1.2.3.4",
			want: map[string]string{"ok": "1.2.3.4"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseDNSOverrides(tt.raw)
			if len(got) != len(tt.want) {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
			for k, v := range tt.want {
				if got[k] != v {
					t.Fatalf("expected %v, got %v", tt.want, got)
				}
			}
		})
	}
}

func TestMatchDNSOverride(t *testing.T) {
	overrides := map[string]string{
		"clb.tencentcloudapi.com": "127.0.0.1",
		"*.tencentcloudapi.com":   "10.0.0.5",
		"*.com":                   "9.9.9.9",
	}

	tests := []struct {
		name    string
		host    string
		wantIP  string
		wantHit bool
	}{
		{
			name:    "exact match wins over wildcard",
			host:    "clb.tencentcloudapi.com",
			wantIP:  "127.0.0.1",
			wantHit: true,
		},
		{
			name:    "wildcard with longest suffix wins",
			host:    "cvm.tencentcloudapi.com",
			wantIP:  "10.0.0.5",
			wantHit: true,
		},
		{
			name:    "shorter wildcard matches unrelated host",
			host:    "example.com",
			wantIP:  "9.9.9.9",
			wantHit: true,
		},
		{
			name:    "no match",
			host:    "example.org",
			wantHit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip, hit := matchDNSOverride(tt.host, overrides)
			if hit != tt.wantHit {
				t.Fatalf("expected hit=%v, got %v", tt.wantHit, hit)
			}
			if hit && ip != tt.wantIP {
				t.Fatalf("expected ip %q, got %q", tt.wantIP, ip)
			}
		})
	}
}

func TestBuildCustomTransportNilWithoutOverrides(t *testing.T) {
	if tr := BuildCustomTransport(nil, false); tr != nil {
		t.Fatalf("expected nil transport when nothing to override, got %T", tr)
	}
}

// 端到端验证: 请求仍按域名发出(Host 头保持域名), 但连接实际建立到覆盖后的地址
func TestDNSOverrideTransportEndToEnd(t *testing.T) {
	// 隔离本机系统代理环境变量, 确保请求直连测试服务器
	t.Setenv("NO_PROXY", "*")
	t.Setenv("no_proxy", "*")

	var gotHost string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHost = r.Host
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	_, port, err := net.SplitHostPort(strings.TrimPrefix(srv.URL, "http://"))
	if err != nil {
		t.Fatalf("parse test server addr: %v", err)
	}

	tr, ok := BuildCustomTransport(map[string]string{"example.test": "127.0.0.1"}, false).(*http.Transport)
	if !ok {
		t.Fatalf("expected *http.Transport, got %T", tr)
	}

	req, err := http.NewRequestWithContext(context.Background(), "GET", "http://example.test:"+port+"/", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	resp, err := tr.RoundTrip(req)
	if err != nil {
		t.Fatalf("round trip failed: %v", err)
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	if gotHost != "example.test:"+port {
		t.Fatalf("expected Host header %q, got %q (domain must be preserved)", "example.test:"+port, gotHost)
	}
}

// 未命中的域名回落到系统 DNS, 不应被覆盖
func TestDNSOverrideTransportMissFallsThrough(t *testing.T) {
	var gotHost string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHost = r.Host
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	_, port, err := net.SplitHostPort(strings.TrimPrefix(srv.URL, "http://"))
	if err != nil {
		t.Fatalf("parse test server addr: %v", err)
	}

	// 只覆盖 example.test, 请求 localhost 应不受影响
	tr, ok := BuildCustomTransport(map[string]string{"example.test": "127.0.0.1"}, false).(*http.Transport)
	if !ok {
		t.Fatalf("expected *http.Transport, got %T", tr)
	}

	req, err := http.NewRequestWithContext(context.Background(), "GET", "http://127.0.0.1:"+port+"/", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	resp, err := tr.RoundTrip(req)
	if err != nil {
		t.Fatalf("round trip failed: %v", err)
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	if gotHost != "127.0.0.1:"+port {
		t.Fatalf("expected Host header %q, got %q", "127.0.0.1:"+port, gotHost)
	}
}

func TestInsecureSkipVerify(t *testing.T) {
	// 隔离本机系统代理环境变量, 确保请求直连测试服务器
	t.Setenv("NO_PROXY", "*")
	t.Setenv("no_proxy", "*")

	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	_, port, err := net.SplitHostPort(strings.TrimPrefix(srv.URL, "https://"))
	if err != nil {
		t.Fatalf("parse test server addr: %v", err)
	}

	// 默认校验: 证书不匹配请求域名, 必须失败
	strictTR, ok := BuildCustomTransport(map[string]string{"example.test": "127.0.0.1"}, false).(*http.Transport)
	if !ok {
		t.Fatalf("expected *http.Transport, got %T", strictTR)
	}
	strictReq, _ := http.NewRequestWithContext(context.Background(), "GET", "https://example.test:"+port+"/", nil)
	if _, err := strictTR.RoundTrip(strictReq); err == nil {
		t.Fatal("expected TLS verification failure without insecure_skip_verify")
	}

	// 跳过校验后成功
	skipTR, ok := BuildCustomTransport(map[string]string{"example.test": "127.0.0.1"}, true).(*http.Transport)
	if !ok {
		t.Fatalf("expected *http.Transport, got %T", skipTR)
	}
	skipReq, _ := http.NewRequestWithContext(context.Background(), "GET", "https://example.test:"+port+"/", nil)
	resp, err := skipTR.RoundTrip(skipReq)
	if err != nil {
		t.Fatalf("round trip with insecure_skip_verify failed: %v", err)
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
}
