//go:build integration

// STS 测试, 配置来源优先级:
//  1. 环境变量 TENCENTCLOUD_SECRET_ID / TENCENTCLOUD_SECRET_KEY / TENCENTCLOUD_SECURITY_TOKEN
//     TENCENTCLOUD_BASE_URL / TENCENTCLOUD_DNS_OVERRIDE / TENCENTCLOUD_INSECURE_SKIP_VERIFY
//  2. ~/.steampipe/config/tencentcloud.spc (或 TENCENTCLOUD_SPC 指定路径)

package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/cvm/v20170312"
	sts "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/sts/v20180813"
)

type stsTestConfig struct {
	SecretID           string
	SecretKey          string
	Token              string
	BaseURL            string
	DNSOverride        map[string]string
	InsecureSkipVerify bool
}

func loadSTSTestConfig(t *testing.T) stsTestConfig {
	t.Helper()
	cfg := stsTestConfig{
		SecretID:    strings.TrimSpace(os.Getenv("TENCENTCLOUD_SECRET_ID")),
		SecretKey:   strings.TrimSpace(os.Getenv("TENCENTCLOUD_SECRET_KEY")),
		Token:       strings.TrimSpace(os.Getenv("TENCENTCLOUD_SECURITY_TOKEN")),
		BaseURL:     strings.TrimSpace(os.Getenv("TENCENTCLOUD_BASE_URL")),
		DNSOverride: parseDNSOverrides(os.Getenv("TENCENTCLOUD_DNS_OVERRIDE")),
		InsecureSkipVerify: os.Getenv("TENCENTCLOUD_INSECURE_SKIP_VERIFY") == "1" ||
			strings.EqualFold(os.Getenv("TENCENTCLOUD_INSECURE_SKIP_VERIFY"), "true"),
	}
	if cfg.SecretID != "" && cfg.SecretKey != "" {
		return cfg
	}

	spcPath := strings.TrimSpace(os.Getenv("TENCENTCLOUD_SPC"))
	if spcPath == "" {
		spcPath = os.Getenv("HOME") + "/.steampipe/config/tencentcloud.spc"
	}
	data, err := os.ReadFile(spcPath)
	if err != nil {
		t.Skipf("no credentials in env and cannot read .spc %s: %v", spcPath, err)
	}
	mergeSpc(t, &cfg, string(data))
	if cfg.SecretID == "" || cfg.SecretKey == "" {
		t.Skipf("no credentials found in env or .spc %s", spcPath)
	}
	return cfg
}

// mergeSpc 用正则从 .spc 文本里提取字段。.spc 是本项目固定 HCL,字段简单。
func mergeSpc(t *testing.T, cfg *stsTestConfig, spc string) {
	t.Helper()
	if v := spcString(spc, "secret_id"); v != "" {
		cfg.SecretID = v
	}
	if v := spcString(spc, "secret_key"); v != "" {
		cfg.SecretKey = v
	}
	if v := spcString(spc, "token"); v != "" {
		cfg.Token = v
	}
	if v := spcString(spc, "base_url"); v != "" {
		cfg.BaseURL = v
	}
	if v := spcBool(spc, "insecure_skip_verify"); v {
		cfg.InsecureSkipVerify = true
	}
	if d := spcDNSOverride(spc); len(d) > 0 {
		cfg.DNSOverride = d
	}
}

func spcString(spc, key string) string {
	re := regexp.MustCompile(`(?m)^\s*` + regexp.QuoteMeta(key) + `\s*=\s*"([^"]*)"`)
	m := re.FindStringSubmatch(spc)
	if len(m) < 2 {
		return ""
	}
	return strings.TrimSpace(m[1])
}

// spcBool 匹配无引号的布尔字段(key = true/false)。
func spcBool(spc, key string) bool {
	re := regexp.MustCompile(`(?m)^\s*` + regexp.QuoteMeta(key) + `\s*=\s*(true|false)`)
	m := re.FindStringSubmatch(spc)
	return len(m) >= 2 && m[1] == "true"
}

func spcDNSOverride(spc string) map[string]string {
	out := map[string]string{}
	start := strings.Index(spc, "dns_override")
	if start < 0 {
		return out
	}
	block := spc[start:]
	open := strings.Index(block, "{")
	close := strings.Index(block, "}")
	if open < 0 || close < 0 || close < open {
		return out
	}
	body := block[open:close]
	re := regexp.MustCompile(`"([^"]+)"\s*=\s*"([^"]+)"`)
	for _, m := range re.FindAllStringSubmatch(body, -1) {
		if len(m) < 3 {
			continue
		}
		out[strings.TrimSpace(m[1])] = strings.TrimSpace(m[2])
	}
	return out
}

// buildSTSClient 构造一个配置好 profile + transport 的 STS client。
// region 透传给 SDK client;测试环境的 STS 实现要求带 Region 参数,公网全局服务可传空。
func buildSTSClient(t *testing.T, cfg stsTestConfig, cred common.CredentialIface, region string) *sts.Client {
	t.Helper()
	cpf := profile.NewClientProfile()
	if cfg.BaseURL != "" {
		cpf.HttpProfile.RootDomain = cfg.BaseURL
	}
	cpf.NetworkFailureMaxRetries = 1
	cpf.RateLimitExceededMaxRetries = 1
	cpf.NetworkFailureRetryDuration = profile.ExponentialBackoff
	cpf.RateLimitExceededRetryDuration = profile.ExponentialBackoff

	c, err := sts.NewClient(cred, region, cpf)
	if err != nil {
		t.Fatalf("sts.NewClient: %v", err)
	}
	if tr := BuildCustomTransport(cfg.DNSOverride, cfg.InsecureSkipVerify); tr != nil {
		c.WithHttpTransport(tr)
	}
	return c
}

func buildCVMClient(t *testing.T, cfg stsTestConfig, cred common.CredentialIface, region string) *cvm.Client {
	t.Helper()
	cpf := profile.NewClientProfile()
	if cfg.BaseURL != "" {
		cpf.HttpProfile.RootDomain = cfg.BaseURL
	}
	cpf.NetworkFailureMaxRetries = 1
	cpf.RateLimitExceededMaxRetries = 1
	c, err := cvm.NewClient(cred, region, cpf)
	if err != nil {
		t.Fatalf("cvm.NewClient: %v", err)
	}
	if tr := BuildCustomTransport(cfg.DNSOverride, cfg.InsecureSkipVerify); tr != nil {
		c.WithHttpTransport(tr)
	}
	return c
}

// printCallerIdentity 打印 GetCallerIdentity 结果,返回 Type 用于断言。
func printCallerIdentity(t *testing.T, label string, resp *sts.GetCallerIdentityResponse) string {
	t.Helper()
	arn, acct, uid, pid, typ := "", "", "", "", ""
	if resp.Response.Arn != nil {
		arn = *resp.Response.Arn
	}
	if resp.Response.AccountId != nil {
		acct = *resp.Response.AccountId
	}
	if resp.Response.UserId != nil {
		uid = *resp.Response.UserId
	}
	if resp.Response.PrincipalId != nil {
		pid = *resp.Response.PrincipalId
	}
	if resp.Response.Type != nil {
		typ = *resp.Response.Type
	}
	reqId := ""
	if resp.Response.RequestId != nil {
		reqId = *resp.Response.RequestId
	}
	t.Logf("[%s] Arn=%s AccountId=%s UserId=%s PrincipalId=%s Type=%s RequestId=%s",
		label, arn, acct, uid, pid, typ, reqId)
	return typ
}

// TestSTSEndToEnd 用当前密钥验证 STS 临时 token 全链路:
//  1. 永久凭证 GetCallerIdentity(对照组)
//  2. 永久凭证 GetFederationToken 取临时凭证
//  3. 临时凭证 GetCallerIdentity(证明 token 被服务端接受)
//  4. 临时凭证 CVM DescribeRegions(token 用于业务 API,失败只 warn)
func TestSTSEndToEnd(t *testing.T) {
	cfg := loadSTSTestConfig(t)
	t.Logf("config: base_url=%q dns_override=%v insecure_skip_verify=%v secret_id=%s...(%d)",
		cfg.BaseURL, cfg.DNSOverride, cfg.InsecureSkipVerify, cfg.SecretID[:8], len(cfg.SecretID))

	// 测试环境的 STS 实现要求请求带 Region 参数(公网全局服务不需要),这里用一个有效 region。
	stsRegion := "ap-guangzhou"

	// ---- Step 1: 永久凭证 GetCallerIdentity ----
	permCred := common.NewCredential(cfg.SecretID, cfg.SecretKey)
	permSTS := buildSTSClient(t, cfg, permCred, stsRegion)
	permIdent, err := permSTS.GetCallerIdentity(sts.NewGetCallerIdentityRequest())
	if err != nil {
		t.Fatalf("Step1 永久凭证 GetCallerIdentity 失败(测试环境 STS 是否可达?): %v", err)
	}
	permType := printCallerIdentity(t, "永久凭证", permIdent)

	// ---- Step 2: GetFederationToken 取临时凭证 ----
	// 最小策略:允许 sts:GetCallerIdentity(让临时凭证自证身份)+ cvm:DescribeRegions
	// (验证临时 token 可用于业务 API)。根账号签发的联合 token 受此策略约束。
	policy := `{"version":"2.0","statement":[{"effect":"allow","action":["sts:GetCallerIdentity","cvm:DescribeRegions"],"resource":["*"]}]}`
	name := "steampipetest"
	duration := uint64(900)
	req := sts.NewGetFederationTokenRequest()
	req.Name = &name
	req.Policy = &policy
	req.DurationSeconds = &duration

	resp, err := permSTS.GetFederationToken(req)
	if err != nil {
		t.Fatalf("Step2 GetFederationToken 失败(密钥是否有 STS:GetFederationToken 权限?): %v", err)
	}
	cred := resp.Response.Credentials
	if cred == nil || cred.TmpSecretId == nil || cred.TmpSecretKey == nil || cred.Token == nil {
		t.Fatalf("Step2 返回的 Credentials 为空: %+v", resp.Response)
	}
	tmpSid, tmpKey, tmpTok := *cred.TmpSecretId, *cred.TmpSecretKey, *cred.Token
	expired := uint64(0)
	if resp.Response.ExpiredTime != nil {
		expired = *resp.Response.ExpiredTime
	}
	t.Logf("Step2 取得临时凭证: TmpSecretId=%s...(%d) TmpSecretKey=***(%d) Token=***(%d) ExpiredTime(unix)=%d",
		tmpSid[:8], len(tmpSid), len(tmpKey), len(tmpTok), expired)

	// ---- Step 3: 临时凭证 GetCallerIdentity(关键断言) ----
	tmpCred := common.NewTokenCredential(tmpSid, tmpKey, tmpTok)
	tmpSTS := buildSTSClient(t, cfg, tmpCred, stsRegion)
	tmpIdent, err := tmpSTS.GetCallerIdentity(sts.NewGetCallerIdentityRequest())
	if err != nil {
		t.Fatalf("Step3 临时凭证 GetCallerIdentity 失败(STS token 未被服务端接受): %v", err)
	}
	tmpType := printCallerIdentity(t, "临时凭证", tmpIdent)

	// 联合身份的 UserId 形如 "uin:federatedUserName",Type 通常为 "federated"
	if tmpType == permType {
		t.Logf("注意:临时凭证 Type=%q 与永久凭证相同,联合身份可能未区分(测试环境行为)", tmpType)
	} else {
		t.Logf("Step3 通过:临时凭证身份 Type=%q,与永久凭证 %q 不同,STS token 端到端生效", tmpType, permType)
	}

	// ---- Step 4: 临时凭证 CVM DescribeRegions(业务 API,失败只 warn) ----
	cvmClient := buildCVMClient(t, cfg, tmpCred, "ap-guangzhou")
	regions, err := cvmClient.DescribeRegions(cvm.NewDescribeRegionsRequest())
	if err != nil {
		t.Logf("Step4 临时凭证 CVM DescribeRegions 失败(测试环境可能未部署 CVM,不影响 STS 结论): %v", err)
	} else {
		n := 0
		if regions.Response.TotalCount != nil {
			n = int(*regions.Response.TotalCount)
		}
		t.Logf("Step4 通过:临时凭证成功调用 CVM DescribeRegions,TotalCount=%d,STS token 可用于业务 API", n)
	}

	// 汇总
	t.Logf("=== STS 端到端测试结论 ===")
	t.Logf("  永久凭证身份: Type=%s", permType)
	t.Logf("  临时凭证身份: Type=%s", tmpType)
	t.Logf("  token 长度=%d, 过期 unix=%d", len(tmpTok), expired)
	b, _ := json.Marshal(struct {
		PermanentType string `json:"permanent_type"`
		FederatedType string `json:"federated_type"`
		TokenLen      int    `json:"token_len"`
		ExpiredUnix   int64  `json:"expired_unix"`
	}{permType, tmpType, len(tmpTok), int64(expired)})
	fmt.Printf("STS_E2E_RESULT_JSON=%s\n", string(b))
}
