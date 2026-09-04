//go:build querytest

package tencentcloud_test

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

const requestID = "request-test"

const (
	buildTimeout        = 5 * time.Minute
	serviceTimeout      = 5 * time.Minute
	queryTimeout        = 2 * time.Minute
	heartbeatInterval   = 10 * time.Second
	dbStartTimeoutSecs  = 120
	maxLoggedOutputRune = 4000
)

type recordedRequest struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Action string `json:"action"`
	Region string `json:"region"`
	Body   any    `json:"body"`
}

type requestRecorder struct {
	t        *testing.T
	mu       sync.Mutex
	requests []recordedRequest
}

func TestQueryContracts(t *testing.T) {
	steampipe, err := exec.LookPath("steampipe")
	if err != nil {
		t.Skip("steampipe 不在 PATH 中，跳过查询契约测试")
	}
	t.Logf("使用 steampipe: %s", steampipe)

	testDir := testSourceDir(t)
	installDir := t.TempDir()
	t.Logf("临时 install-dir: %s", installDir)

	dbPort, err := freePort()
	if err != nil {
		t.Fatalf("分配空闲数据库端口: %v", err)
	}
	t.Logf("为临时 Steampipe 服务分配数据库端口: %d", dbPort)

	t.Cleanup(func() {
		t.Logf("停止临时 Steampipe 服务 (install-dir=%s)", installDir)
		ctx, cancel := context.WithTimeout(context.Background(), serviceTimeout)
		defer cancel()
		cmd := exec.CommandContext(ctx, steampipe, "service", "stop", "--install-dir", installDir)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Logf("停止服务失败(忽略): %v\n%s", err, out)
		}
	})

	pluginsDir := filepath.Join(installDir, "plugins")
	pluginPath := filepath.Join(pluginsDir, "local", "tencentcloud", "tencentcloud.plugin")
	if err := os.MkdirAll(filepath.Dir(pluginPath), 0o755); err != nil {
		t.Fatalf("创建插件目录: %v", err)
	}

	runCommand(t, testDir, "构建插件", buildTimeout, "go", "build", "-o", pluginPath, "..")
	versions := `{
  "plugins": {
    "local/tencentcloud": {
      "name": "local/tencentcloud",
      "version": "0.0.0",
      "struct_version": 20230502
    }
  },
  "struct_version": 20220411
}
`
	if err := os.WriteFile(filepath.Join(pluginsDir, "versions.json"), []byte(versions), 0o600); err != nil {
		t.Fatalf("写入本地插件版本注册表: %v", err)
	}

	recorder := &requestRecorder{t: t}
	server := httptest.NewServer(http.HandlerFunc(recorder.handle))
	defer server.Close()
	t.Logf("模拟 Tencent Cloud API 服务已启动: %s", server.URL)

	configDir := filepath.Join(installDir, "config")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("创建配置目录: %v", err)
	}
	// 显式声明 regions,让地域矩阵来自 mock 的 DescribeRegions 结果(ap-guangzhou);
	// 若省略 regions,插件会退化为只查兜底地域 ap-singapore,与各契约期望的 region 不符。
	// allow_insecure_endpoint 放行 mock server 的 http:// scheme(纯本地回环,无 TLS,无需 insecure_skip_verify)。
	config := fmt.Sprintf(`connection "tencentcloud" {
  plugin     = "local/tencentcloud"
  secret_id  = "test-secret-id"
  secret_key = "test-secret-key"
  endpoint    = %q
  allow_insecure_endpoint = true
  regions     = ["*"]
}
`, server.URL)
	if err := os.WriteFile(filepath.Join(configDir, "tencentcloud.spc"), []byte(config), 0o600); err != nil {
		t.Fatalf("写入临时 Steampipe 配置: %v", err)
	}

	// 独占端口 + 仅监听 localhost，避免与用户已有的 steampipe 服务(默认 9193)冲突导致长时间阻塞。
	// 同时关闭查询缓存，否则前序契约的结果会被复用，导致"请求次数"断言不可复现。
	options := fmt.Sprintf(`options "database" {
  port          = %d
  listen        = "local"
  start_timeout = %d
  cache         = false
}

options "general" {
  update_check = false
}
`, dbPort, dbStartTimeoutSecs)
	if err := os.WriteFile(filepath.Join(configDir, "default.spc"), []byte(options), 0o600); err != nil {
		t.Fatalf("写入临时 Steampipe options: %v", err)
	}

	queryFiles, err := filepath.Glob(filepath.Join(testDir, "tests", "*", "*-query.sql"))
	if err != nil {
		t.Fatalf("扫描 SQL 契约: %v", err)
	}
	sort.Strings(queryFiles)
	if len(queryFiles) == 0 {
		t.Fatal("未找到 tests/*/*-query.sql")
	}
	t.Logf("共发现 %d 个查询契约", len(queryFiles))

	// 预先拉起服务，首次启动会初始化 Postgres，耗时较长，单独打日志便于观测。
	runCommand(t, testDir, "启动 Steampipe 服务", serviceTimeout, steampipe,
		"service", "start", "--install-dir", installDir, "--database-port", strconv.Itoa(dbPort), "--database-listen", "local")

	for i, queryFile := range queryFiles {
		queryFile := queryFile
		relative, err := filepath.Rel(filepath.Join(testDir, "tests"), queryFile)
		if err != nil {
			t.Fatalf("计算测试名称: %v", err)
		}
		name := strings.TrimSuffix(filepath.ToSlash(relative), "-query.sql")
		t.Run(name, func(t *testing.T) {
			t.Logf("[%d/%d] 开始契约: %s", i+1, len(queryFiles), name)
			query := readFile(t, queryFile)
			expectedRows := readJSON(t, strings.TrimSuffix(queryFile, "-query.sql")+"-expected.json")
			expectedRequests := readRequests(t, strings.TrimSuffix(queryFile, "-query.sql")+"-requests.json")

			recorder.reset()
			stdout := runCommand(t, testDir, "执行 Steampipe 查询", queryTimeout, steampipe,
				"query", "--install-dir", installDir, "--output", "json", strings.TrimSpace(query))

			var output struct {
				Rows any `json:"rows"`
			}
			if err := json.Unmarshal(stdout, &output); err != nil {
				t.Fatalf("解析 Steampipe JSON 输出: %v\nstdout:\n%s", err, stdout)
			}
			if !reflect.DeepEqual(expectedRows, output.Rows) {
				t.Errorf("查询结果不匹配\nexpected:\n%s\nactual:\n%s", prettyJSON(expectedRows), prettyJSON(output.Rows))
			}

			actualRequests := recorder.snapshot()
			t.Logf("捕获 %d 条上游请求", len(actualRequests))
			if err := matchRequests(expectedRequests, actualRequests); err != nil {
				t.Errorf("请求契约不匹配: %v\nexpected:\n%s\nactual:\n%s", err, prettyJSON(expectedRequests), prettyJSON(actualRequests))
			}
		})
	}
}

// freePort 让内核分配一个空闲端口，避免测试与本机已运行的 steampipe 服务抢占 9193。
func freePort() (int, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		return 0, fmt.Errorf("非预期的监听地址类型 %T", listener.Addr())
	}
	return addr.Port, nil
}

func testSourceDir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("无法定位 query_test.go")
	}
	return filepath.Dir(file)
}

func runCommand(t *testing.T, dir, description string, timeout time.Duration, name string, args ...string) []byte {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	// 关闭 stdin，防止 steampipe 在缺少输入时进入交互模式而永久阻塞。
	cmd.Stdin = nil

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	// stderr 实时透传到测试日志，便于观察 steampipe 的启动/报错过程。
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		t.Fatalf("%s: 创建 stderr 管道失败: %v", description, err)
	}

	start := time.Now()
	t.Logf("%s: 执行 %s %s (超时 %s)", description, name, strings.Join(args, " "), timeout)
	if err := cmd.Start(); err != nil {
		t.Fatalf("%s启动失败: %v", description, err)
	}

	var streamWG sync.WaitGroup
	streamWG.Add(1)
	go func() {
		defer streamWG.Done()
		scanner := bufio.NewScanner(stderrPipe)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for scanner.Scan() {
			line := scanner.Text()
			stderr.WriteString(line)
			stderr.WriteString("\n")
			t.Logf("%s[stderr] %s", description, line)
		}
	}()

	done := make(chan error, 1)
	go func() {
		streamWG.Wait()
		done <- cmd.Wait()
	}()

	// 心跳日志：长时间无输出时也能看到测试仍在推进，而不是“完全卡住无 log”。
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()
	var runErr error
	for waiting := true; waiting; {
		select {
		case runErr = <-done:
			waiting = false
		case <-ticker.C:
			t.Logf("%s: 仍在运行，已耗时 %s", description, time.Since(start).Round(time.Second))
		}
	}

	elapsed := time.Since(start).Round(time.Millisecond)
	if runErr != nil {
		if ctx.Err() == context.DeadlineExceeded {
			t.Fatalf("%s超时(%s)后被终止\nstdout:\n%s\nstderr:\n%s",
				description, timeout, truncate(stdout.String()), truncate(stderr.String()))
		}
		t.Fatalf("%s失败(耗时 %s): %v\nstdout:\n%s\nstderr:\n%s",
			description, elapsed, runErr, truncate(stdout.String()), truncate(stderr.String()))
	}
	t.Logf("%s: 完成，耗时 %s，stdout %d 字节", description, elapsed, stdout.Len())
	return stdout.Bytes()
}

func truncate(s string) string {
	if len(s) <= maxLoggedOutputRune {
		return s
	}
	return s[:maxLoggedOutputRune] + fmt.Sprintf("\n...(已截断，共 %d 字节)", len(s))
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("读取 %s: %v", path, err)
	}
	return string(data)
}

func readJSON(t *testing.T, path string) any {
	t.Helper()
	data := readFile(t, path)
	var value any
	if err := json.Unmarshal([]byte(data), &value); err != nil {
		t.Fatalf("解析 %s: %v", path, err)
	}
	return value
}

func readRequests(t *testing.T, path string) []recordedRequest {
	t.Helper()
	data := readFile(t, path)
	var requests []recordedRequest
	if err := json.Unmarshal([]byte(data), &requests); err != nil {
		t.Fatalf("解析 %s: %v", path, err)
	}
	return requests
}

func (r *requestRecorder) reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.requests = nil
}

func (r *requestRecorder) snapshot() []recordedRequest {
	r.mu.Lock()
	defer r.mu.Unlock()
	filtered := make([]recordedRequest, 0, len(r.requests))
	for _, req := range r.requests {
		if req.Action == "GetUserAppId" {
			continue
		}
		filtered = append(filtered, req)
	}
	return filtered
}

func (r *requestRecorder) handle(w http.ResponseWriter, req *http.Request) {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		r.logf("模拟 API 读取请求体失败: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	var body any = map[string]any{}
	if len(bytes.TrimSpace(bodyBytes)) > 0 {
		if err := json.Unmarshal(bodyBytes, &body); err != nil {
			body = map[string]any{"invalid_json": string(bodyBytes)}
		}
	}

	record := recordedRequest{
		Method: req.Method,
		Path:   req.URL.Path,
		Action: req.Header.Get("X-TC-Action"),
		Region: req.Header.Get("X-TC-Region"),
		Body:   body,
	}
	r.mu.Lock()
	r.requests = append(r.requests, record)
	r.mu.Unlock()
	r.logf("模拟 API 收到请求: action=%s region=%s path=%s", record.Action, record.Region, record.Path)

	// COS is a RESTful (non-TC3) API: requests carry no X-TC-Action header and
	// responses are XML. Route by method+path when the action header is empty.
	fixture, ok := responseFixtures[record.Action]
	contentType := "application/json"
	if !ok {
		fixture, ok = cosFixtures[record.Method+" "+record.Path]
		if ok {
			contentType = "application/xml"
		}
	}
	if !ok {
		r.logf("模拟 API 缺少 action=%s 的响应桩", record.Action)
		http.Error(w, "unsupported action: "+record.Action, http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", contentType)
	_, _ = io.WriteString(w, fixture)
}

// logf 在测试结束后调用是不安全的，因此这里只做尽力而为的日志输出。
func (r *requestRecorder) logf(format string, args ...any) {
	if r.t == nil {
		return
	}
	defer func() { _ = recover() }()
	r.t.Logf(format, args...)
}

func matchRequests(expected, actual []recordedRequest) error {
	if len(actual) != len(expected) {
		return fmt.Errorf("请求次数应为 %d，实际为 %d", len(expected), len(actual))
	}

	used := make([]bool, len(actual))
	for _, want := range expected {
		matched := false
		for i, got := range actual {
			if used[i] {
				continue
			}
			if want.Method != got.Method || want.Path != got.Path || want.Action != got.Action || want.Region != got.Region {
				continue
			}
			if !isSubset(want.Body, got.Body) {
				continue
			}
			used[i] = true
			matched = true
			break
		}
		if !matched {
			return fmt.Errorf("没有实际请求匹配 %s", prettyJSON(want))
		}
	}
	return nil
}

func isSubset(expected, actual any) bool {
	switch want := expected.(type) {
	case map[string]any:
		got, ok := actual.(map[string]any)
		if !ok {
			return false
		}
		for key, wantValue := range want {
			gotValue, ok := got[key]
			if !ok || !isSubset(wantValue, gotValue) {
				return false
			}
		}
		return true
	case []any:
		got, ok := actual.([]any)
		if !ok || len(want) > len(got) {
			return false
		}
		used := make([]bool, len(got))
		for _, wantValue := range want {
			matched := false
			for i, gotValue := range got {
				if !used[i] && isSubset(wantValue, gotValue) {
					used[i] = true
					matched = true
					break
				}
			}
			if !matched {
				return false
			}
		}
		return true
	default:
		return reflect.DeepEqual(expected, actual)
	}
}

func prettyJSON(value any) string {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Sprintf("%#v", value)
	}
	return string(data)
}

var responseFixtures = map[string]string{
	"DescribeInstances":            `{"Response":{"TotalCount":1,"Total":1,"InstanceSet":[{"InstanceId":"ins-test","InstanceName":"test-instance","PublicIpAddresses":["203.0.113.10"],"Placement":{"Zone":"ap-guangzhou-3","ProjectId":0}}],"Instances":[{"InstanceId":"wafins-test","InstanceName":"test-waf","Region":"ap-guangzhou","Edition":"saas"}],"RequestId":"` + requestID + `"}}`,
	"DescribeZones":                `{"Response":{"ZoneSet":[{"Zone":"ap-guangzhou-3","ZoneName":"test-zone","ZoneId":"100003","ZoneState":"AVAILABLE"}],"RequestId":"` + requestID + `"}}`,
	"DescribeKeyPairs":             `{"Response":{"TotalCount":1,"KeyPairSet":[{"KeyId":"skey-test","KeyName":"test-key-pair"}],"RequestId":"` + requestID + `"}}`,
	"DescribeLaunchTemplates":      `{"Response":{"TotalCount":1,"LaunchTemplateSet":[{"LaunchTemplateId":"lt-test","LaunchTemplateName":"test-launch-template"}],"RequestId":"` + requestID + `"}}`,
	"DescribeImages":               `{"Response":{"TotalCount":1,"ImageSet":[{"ImageId":"img-test","ImageName":"test-image"}],"RequestId":"` + requestID + `"}}`,
	"DescribeRegions":              `{"Response":{"RegionSet":[{"Region":"ap-guangzhou","RegionName":"test-region","RegionState":"AVAILABLE"}],"RequestId":"` + requestID + `"}}`,
	"DescribeDisks":                `{"Response":{"TotalCount":1,"DiskSet":[{"DiskId":"disk-test","DiskName":"test-disk","InstanceId":"ins-test","Placement":{"Zone":"ap-guangzhou-3","ProjectId":0}}],"RequestId":"` + requestID + `"}}`,
	"DescribeDiskBackups":          `{"Response":{"TotalCount":1,"DiskBackupSet":[{"DiskBackupId":"dbp-test","DiskBackupName":"test-disk-backup","DiskId":"disk-test"}],"RequestId":"` + requestID + `"}}`,
	"DescribeSnapshots":            `{"Response":{"TotalCount":1,"SnapshotSet":[{"SnapshotId":"snap-test","SnapshotName":"test-snapshot","DiskId":"disk-test","Placement":{"Zone":"ap-guangzhou-3","ProjectId":0}}],"RequestId":"` + requestID + `"}}`,
	"DescribeAutoSnapshotPolicies": `{"Response":{"TotalCount":1,"AutoSnapshotPolicySet":[{"AutoSnapshotPolicyId":"asp-test","AutoSnapshotPolicyName":"test-auto-snapshot-policy"}],"RequestId":"` + requestID + `"}}`,
	"DescribeSnapshotGroups":       `{"Response":{"TotalCount":1,"SnapshotGroupSet":[{"SnapshotGroupId":"csnap-test","SnapshotGroupName":"test-snapshot-group","SnapshotIdSet":["snap-test"]}],"RequestId":"` + requestID + `"}}`,
	// Monitor
	"GetMonitorData": `{"Response":{"Period":86400,"MetricName":"CPUUsage","DataPoints":[{"Dimensions":[{"Name":"InstanceId","Value":"ins-test"}],"Timestamps":[1691193600.0],"Values":[42.5]}],"StartTime":"2023-07-06T00:00:00Z","EndTime":"2023-08-05T00:00:00Z","RequestId":"` + requestID + `"}}`,
	// CVM
	"DescribeHosts": `{"Response":{"TotalCount":1,"HostSet":[{"HostId":"host-test","HostName":"test-host","Placement":{"Zone":"ap-guangzhou-3","ProjectId":0}}],"HostList":[{"Domain":"test.example.com","DomainId":"waf-host-test","Edition":"clb-waf","Region":"ap-guangzhou","Status":1}],"RequestId":"` + requestID + `"}}`,
	// CBS
	"DescribeDiskStoragePool": `{"Response":{"TotalCount":1,"DiskStoragePoolSet":[{"CdcId":"cdc-test","CdcName":"test-cdc"}],"RequestId":"` + requestID + `"}}`,
	// CLB
	"DescribeLoadBalancers":        `{"Response":{"TotalCount":1,"LoadBalancerSet":[{"LoadBalancerId":"lb-test","LoadBalancerName":"test-lb"}],"RequestId":"` + requestID + `"}}`,
	"DescribeListeners":            `{"Response":{"Listeners":[{"ListenerId":"listener-test","ListenerName":"test-listener","LoadBalancerId":"lb-test"}],"RequestId":"` + requestID + `"}}`,
	"DescribeTargetGroupList":      `{"Response":{"TotalCount":1,"TargetGroupSet":[{"TargetGroupId":"tg-test","TargetGroupName":"test-tg"}],"RequestId":"` + requestID + `"}}`,
	"DescribeCustomizedConfigList": `{"Response":{"TotalCount":1,"ConfigList":[{"UconfigId":"uconfig-test","ConfigName":"test-config","ConfigType":"CLB"}],"RequestId":"` + requestID + `"}}`,
	"DescribeCrossTargets":         `{"Response":{"TotalCount":1,"CrossTargetSet":[{"IP":"1.2.3.4","VpcId":"vpc-test"}],"RequestId":"` + requestID + `"}}`,
	"DescribeBlockIPList":          `{"Response":{"BlockedIPCount":1,"BlockedIPList":[{"IP":"1.2.3.4"}],"ClientIPField":"X-Forwarded-For","RequestId":"` + requestID + `"}}`,
	"DescribeRewrite":              `{"Response":{"RewriteSet":[{"ListenerId":"listener-test","LocationId":"loc-test","RewriteTarget":{"TargetListenerId":"listener-test2"}}],"RequestId":"` + requestID + `"}}`,
	"DescribeTargets":              `{"Response":{"Listeners":[{"ListenerId":"listener-test","Protocol":"TCP","Port":80,"Targets":[{"InstanceId":"ins-test","Port":8080,"Weight":10}]}],"RequestId":"` + requestID + `"}}`,
	"DescribeTargetGroupInstances": `{"Response":{"TotalCount":1,"TargetGroupInstanceSet":[{"TargetGroupId":"tg-test","InstanceId":"ins-test"}],"RequestId":"` + requestID + `"}}`,
	// VPC
	"DescribeVpcs":                  `{"Response":{"TotalCount":1,"VpcSet":[{"VpcId":"vpc-test","VpcName":"test-vpc"}],"RequestId":"` + requestID + `"}}`,
	"DescribeSubnets":               `{"Response":{"TotalCount":1,"SubnetSet":[{"SubnetId":"subnet-test","SubnetName":"test-subnet"}],"RequestId":"` + requestID + `"}}`,
	"DescribeRouteTables":           `{"Response":{"TotalCount":1,"RouteTableSet":[{"RouteTableId":"rtb-test","RouteTableName":"test-rt","VpcId":"vpc-test","RouteSet":[{"RouteItemId":"route-item-test","DestinationCidrBlock":"10.0.0.0/24","GatewayType":"NAT","GatewayId":"nat-test"}]}],"RequestId":"` + requestID + `"}}`,
	"DescribeSecurityGroups":        `{"Response":{"TotalCount":1,"SecurityGroupSet":[{"SecurityGroupId":"sg-test","SecurityGroupName":"test-sg"}],"RequestId":"` + requestID + `"}}`,
	"DescribeSecurityGroupPolicies": `{"Response":{"SecurityGroupPolicySet":{"Ingress":[{"PolicyIndex":0,"Protocol":"TCP","Port":"80","CidrBlock":"10.0.0.0/8","Action":"ACCEPT"}],"Egress":[{"PolicyIndex":0,"Protocol":"ALL","Port":"all","CidrBlock":"0.0.0.0/0","Action":"ACCEPT"}]},"RequestId":"` + requestID + `"}}`,
	"DescribeNetworkInterfaces":     `{"Response":{"TotalCount":1,"NetworkInterfaceSet":[{"NetworkInterfaceId":"eni-test","NetworkInterfaceName":"test-eni"}],"RequestId":"` + requestID + `"}}`,
	"DescribeAddresses":             `{"Response":{"TotalCount":1,"AddressSet":[{"AddressId":"eip-test","AddressName":"test-eip"}],"RequestId":"` + requestID + `"}}`,
	"DescribeNatGateways":           `{"Response":{"TotalCount":1,"NatGatewaySet":[{"NatGatewayId":"nat-test","NatGatewayName":"test-nat"}],"RequestId":"` + requestID + `"}}`,
	"DescribeBandwidthPackages":     `{"Response":{"TotalCount":1,"BandwidthPackageSet":[{"BandwidthPackageId":"bp-test","BandwidthPackageName":"test-bp"}],"RequestId":"` + requestID + `"}}`,
	"DescribeTrafficPackages":       `{"Response":{"TotalCount":1,"TrafficPackageSet":[{"TrafficPackageId":"tp-test","TrafficPackageName":"test-tp"}],"RequestId":"` + requestID + `"}}`,
	// DC
	"DescribeDirectConnects":                `{"Response":{"TotalCount":1,"DirectConnectSet":[{"DirectConnectId":"dc-test","DirectConnectName":"test-dc"}],"RequestId":"` + requestID + `"}}`,
	"DescribeDirectConnectTunnels":          `{"Response":{"TotalCount":1,"DirectConnectTunnelSet":[{"DirectConnectTunnelId":"dcx-test","DirectConnectTunnelName":"test-dcx"}],"RequestId":"` + requestID + `"}}`,
	"DescribeDirectConnectGateways":         `{"Response":{"TotalCount":1,"DirectConnectGatewaySet":[{"DirectConnectGatewayId":"dcg-test","DirectConnectGatewayName":"test-dcg"}],"RequestId":"` + requestID + `"}}`,
	"DescribeDirectConnectGatewayCcnRoutes": `{"Response":{"TotalCount":1,"RouteSet":[{"RouteId":"route-test","DestinationCidrBlock":"10.0.0.0/8"}],"RequestId":"` + requestID + `"}}`,
	"DescribeInternetAddress":               `{"Response":{"TotalCount":1,"Subnets":[{"InstanceId":"dc-ip-test","Subnet":"1.2.3.0/24","AddrType":0}],"RequestId":"` + requestID + `"}}`,
	// VPN
	"DescribeVpnGateways":         `{"Response":{"TotalCount":1,"VpnGatewaySet":[{"VpnGatewayId":"vpngw-test","VpnGatewayName":"test-vpngw"}],"RequestId":"` + requestID + `"}}`,
	"DescribeCustomerGateways":    `{"Response":{"TotalCount":1,"CustomerGatewaySet":[{"CustomerGatewayId":"cgw-test","CustomerGatewayName":"test-cgw"}],"RequestId":"` + requestID + `"}}`,
	"DescribeVpnConnections":      `{"Response":{"TotalCount":1,"VpnConnectionSet":[{"VpnConnectionId":"vpnx-test","VpnConnectionName":"test-vpnx"}],"RequestId":"` + requestID + `"}}`,
	"DescribeVpnGatewayRoutes":    `{"Response":{"TotalCount":1,"Routes":[{"RouteId":"vpnr-test","DestinationCidrBlock":"10.0.0.0/8"}],"RequestId":"` + requestID + `"}}`,
	"DescribeVpnGatewayCcnRoutes": `{"Response":{"TotalCount":1,"RouteSet":[{"RouteId":"vpnr-test","DestinationCidrBlock":"10.0.0.0/8"}],"RequestId":"` + requestID + `"}}`,
	// CCN
	"DescribeCcns":                 `{"Response":{"TotalCount":1,"CcnSet":[{"CcnId":"ccn-test","CcnName":"test-ccn"}],"RequestId":"` + requestID + `"}}`,
	"DescribeCcnAttachedInstances": `{"Response":{"TotalCount":1,"InstanceSet":[{"InstanceId":"vpc-test","InstanceType":"VPC","InstanceName":"test-vpc"}],"RequestId":"` + requestID + `"}}`,
	"DescribeCcnRoutes":            `{"Response":{"TotalCount":1,"RouteSet":[{"RouteId":"ccnr-test","DestinationCidrBlock":"10.0.0.0/8"}],"RequestId":"` + requestID + `"}}`,
	// SSL
	"DescribeCertificates": `{"Response":{"TotalCount":1,"Certificates":[{"CertificateId":"cert-test","Alias":"test-cert","Domain":"example.com"}],"RequestId":"` + requestID + `"}}`,
	// CAM
	"GetUserAppId":          `{"Response":{"Uin":"100012345678","OwnerUin":"100012345678","AppId":1250000000,"RequestId":"` + requestID + `"}}`,
	"ListUsers":             `{"Response":{"Data":[{"Uin":123,"Uid":1,"Name":"test-user","Remark":"test-remark","ConsoleLogin":1}],"RequestId":"` + requestID + `"}}`,
	"GetUser":               `{"Response":{"Uin":123,"Uid":1,"Name":"test-user","Remark":"test-remark","ConsoleLogin":1},"RequestId":"` + requestID + `"}`,
	"ListAccessKeys":        `{"Response":{"AccessKeys":[{"AccessKeyId":"AKID-test","Status":"Active","CreateTime":"2024-01-01T00:00:00Z"}],"RequestId":"` + requestID + `"}}`,
	"GetSecurityLastUsed":   `{"Response":{"SecretIdLastUsedRows":[{"SecretId":"AKID-test","LastUsedDate":"2024-01-01"}],"RequestId":"` + requestID + `"}}`,
	"ListGroups":            `{"Response":{"TotalNum":1,"GroupInfo":[{"GroupId":1,"GroupName":"test-group","Remark":"test-remark"}],"RequestId":"` + requestID + `"}}`,
	"GetGroup":              `{"Response":{"GroupId":1,"GroupName":"test-group","Remark":"test-remark","GroupNum":3},"RequestId":"` + requestID + `"}`,
	"ListPolicies":          `{"Response":{"TotalNum":1,"List":[{"PolicyId":1,"PolicyName":"test-policy","Type":2,"CreateMode":1}],"RequestId":"` + requestID + `"}}`,
	"GetPolicy":             `{"Response":{"PolicyId":1,"PolicyName":"test-policy","Type":2,"Description":"test-description"},"RequestId":"` + requestID + `"}`,
	"ListEntitiesForPolicy": `{"Response":{"TotalNum":1,"List":[{"Id":"user-test","Name":"test-user","Uin":123,"RelatedType":1}],"RequestId":"` + requestID + `"}}`,
	"DescribeRoleList":      `{"Response":{"TotalNum":1,"List":[{"RoleId":"role-test","RoleName":"test-role","RoleType":"user"}],"RequestId":"` + requestID + `"}}`,
	"GetRole":               `{"Response":{"RoleInfo":{"RoleId":"role-test","RoleName":"test-role","RoleType":"user"}},"RequestId":"` + requestID + `"}`,
	"ListSAMLProviders":     `{"Response":{"SAMLProviderSet":[{"Name":"test-provider","Description":"test-description"}],"RequestId":"` + requestID + `"}}`,
	"GetSAMLProvider":       `{"Response":{"Name":"test-provider","Description":"test-description","SAMLMetadata":"<md:EntityDescriptor/>"},"RequestId":"` + requestID + `"}`,
	// CLS
	"DescribeLogsets":       `{"Response":{"TotalCount":1,"Logsets":[{"LogsetId":"logset-test","LogsetName":"test-logset","TopicCount":1}],"RequestId":"` + requestID + `"}}`,
	"DescribeTopics":        `{"Response":{"TotalCount":1,"Topics":[{"TopicId":"topic-test","TopicName":"test-topic","LogsetId":"logset-test","Status":true}],"RequestId":"` + requestID + `"}}`,
	"DescribeConfigs":       `{"Response":{"TotalCount":1,"Configs":[{"ConfigId":"config-test","Name":"test-config","Output":"topic-test","LogType":"json_log"}],"RequestId":"` + requestID + `"}}`,
	"DescribeMachineGroups": `{"Response":{"TotalCount":1,"MachineGroups":[{"GroupId":"mg-test","GroupName":"test-mg","OSType":0}],"RequestId":"` + requestID + `"}}`,
	"DescribeAlarms":        `{"Response":{"TotalCount":1,"Alarms":[{"AlarmId":"alarm-test","Name":"test-alarm","Status":true}],"RequestId":"` + requestID + `"}}`,
	"DescribeAlarmNotices":  `{"Response":{"TotalCount":1,"AlarmNotices":[{"AlarmNoticeId":"notice-test","Name":"test-notice","Type":"Trigger"}],"RequestId":"` + requestID + `"}}`,
	// TKE
	"DescribeClusters":         `{"Response":{"TotalCount":1,"Clusters":[{"ClusterId":"cls-test","ClusterName":"test-cluster","ClusterType":"MANAGED_CLUSTER","ClusterStatus":"Running"}],"RequestId":"` + requestID + `"}}`,
	"DescribeClusterInstances": `{"Response":{"TotalCount":1,"InstanceSet":[{"InstanceId":"ins-test","InstanceRole":"WORKER","InstanceState":"running","LanIP":"10.0.0.10"}],"RequestId":"` + requestID + `"}}`,
	// WAF
	"DescribeDomains":         `{"Response":{"Total":1,"Domains":[{"Domain":"test.example.com","DomainId":"waf-domain-test","InstanceId":"wafins-test","Edition":"sparta-waf","Region":"ap-guangzhou"}],"RequestId":"` + requestID + `"}}`,
	"DescribeObjects":         `{"Response":{"ClbObjects":[{"ObjectId":"clb-test","ObjectName":"test-clb","InstanceId":"wafins-test","Type":"CLB","Status":1}],"RequestId":"` + requestID + `"}}`,
	"DescribeIpAccessControl": `{"Response":{"Data":{"TotalCount":1,"Res":[{"RuleId":1,"Ip":"1.2.3.4","ActionType":40,"Note":"test-rule"}]},"RequestId":"` + requestID + `"}}`,
	"DescribeCustomRuleList":  `{"Response":{"TotalCount":"1","RuleList":[{"RuleId":"rule-test","Name":"test-rule","ActionType":"1","Status":"1","SortId":"0"}],"RequestId":"` + requestID + `"}}`,
	// CVM spread placement groups
	"DescribeDisasterRecoverGroups": `{"Response":{"TotalCount":1,"DisasterRecoverGroupSet":[{"DisasterRecoverGroupId":"drg-test","Name":"test-drg","Type":"HOST","CvmQuotaTotal":50,"CurrentNum":1,"InstanceIds":["ins-test"],"CreateTime":"2024-01-01T00:00:00Z","Affinity":1,"Tags":[]}],"RequestId":"` + requestID + `"}}`,
	// AS (Auto Scaling)
	"DescribeAutoScalingGroups": `{"Response":{"TotalCount":1,"AutoScalingGroupSet":[{"AutoScalingGroupId":"asg-test","AutoScalingGroupName":"test-asg","AutoScalingGroupStatus":"NORMAL","CreatedTime":"2024-01-01T00:00:00Z","DefaultCooldown":300,"DesiredCapacity":2,"EnabledStatus":"ENABLED","InstanceCount":2,"InServiceInstanceCount":2,"LaunchConfigurationId":"asc-test","LaunchConfigurationName":"test-asc","LoadBalancerIdSet":[],"MaxSize":5,"MinSize":1,"ProjectId":0,"SubnetIdSet":["subnet-test"],"TerminationPolicySet":["OLDEST_INSTANCE"],"VpcId":"vpc-test","ZoneSet":["ap-guangzhou-3"],"RetryPolicy":"IMMEDIATE_RETRY","InActivityStatus":"NOT_IN_ACTIVITY","Tags":[],"Ipv6AddressCount":0,"MultiZoneSubnetPolicy":"PRIORITY","HealthCheckType":"CLB","LoadBalancerHealthCheckGracePeriod":0,"InstanceAllocationPolicy":"LAUNCH_CONFIGURATION","CapacityRebalance":false,"ConcurrentScaleOutForDesiredCapacity":false}],"RequestId":"` + requestID + `"}}`,
}

// cosFixtures routes COS RESTful requests (no X-TC-Action header) by
// method+path and returns XML responses, mirroring the real COS service.
var cosFixtures = map[string]string{
	"GET /": `<?xml version="1.0" encoding="UTF-8"?><ListAllMyBucketsResult><Owner><ID>1250000000</ID><DisplayName>test</DisplayName></Owner><Buckets><Bucket><Name>bucket-test-1250000000</Name><Location>ap-guangzhou</Location><CreationDate>2024-01-01T00:00:00Z</CreationDate><Type>cos</Type><BucketType>cos</BucketType></Bucket></Buckets></ListAllMyBucketsResult>`,
}
