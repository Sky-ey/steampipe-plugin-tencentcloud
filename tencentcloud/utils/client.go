package utils

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"maps"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/profile"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// Client Common client interface
type Client interface {
	Init(region string) *common.Client
}

const (
	HTTPTransportCacheKeyPrefix = "tencentcloud-http-transport-v1-"
	HTTPTransportIdleTimeout    = 30 * time.Second
)

var sharedHTTPTransportMu sync.Mutex

// InitClient Initialize client for module
func InitClient(ctx context.Context, d *plugin.QueryData, client Client) error {
	// With current request region
	return InitClientInRegion(ctx, d, client, d.EqualsQualString(MatrixKeyRegion))
}

// InitClientInRegion Initialize client in region
func InitClientInRegion(ctx context.Context, d *plugin.QueryData, client Client, region string) error {
	cfg := GetConfig(d.Connection)

	// Endpoint Request URL for all module.
	endpoint := strings.TrimSpace(os.Getenv("TENCENTCLOUD_ENDPOINT"))
	// BaseURL Root domain for URL concatenation.
	baseUrl := strings.TrimSpace(os.Getenv("TENCENTCLOUD_BASE_URL"))
	// InsecureSkipVerify Skip TLS certificate verification.
	insecure := strings.TrimSpace(os.Getenv("TENCENTCLOUD_INSECURE_SKIP_VERIFY"))
	// DNS Override, Format: host=ip[,host=ip] Host support `*.` format.
	// Example: TENCENTCLOUD_DNS_OVERRIDE="cvm.tencentcloudapi.com=127.0.0.1,*.tencentcloudapi.com=10.0.0.5"
	overrides := ParseDNSOverrides(os.Getenv("TENCENTCLOUD_DNS_OVERRIDE"))

	maxRetries := 3
	clientProfile := profile.NewClientProfile()

	if cfg.Endpoint != nil && *cfg.Endpoint != "" {
		endpoint = *cfg.Endpoint
	}
	if cfg.BaseUrl != nil && *cfg.BaseUrl != "" {
		baseUrl = *cfg.BaseUrl
	}
	if cfg.DNSOverride != nil {
		maps.Copy(overrides, cfg.DNSOverride)
	}
	if cfg.Timeout != nil && *cfg.Timeout > 0 {
		clientProfile.HttpProfile.ReqTimeout = *cfg.Timeout
	}
	if cfg.MaxRetry != nil && *cfg.MaxRetry > 0 {
		maxRetries = *cfg.MaxRetry
	}
	if endpoint == "" && baseUrl != "" {
		clientProfile.HttpProfile.RootDomain = baseUrl
	}

	if err := applyEndpoint(ctx, clientProfile, endpoint, cfg); err != nil {
		return err
	}
	clientProfile.NetworkFailureMaxRetries = maxRetries
	clientProfile.RateLimitExceededMaxRetries = maxRetries
	clientProfile.NetworkFailureRetryDuration = profile.ExponentialBackoff
	clientProfile.RateLimitExceededRetryDuration = profile.ExponentialBackoff

	// Parse Credential
	credential, err := CreateCredential(cfg)
	if err != nil {
		return err
	}
	c := client.Init(NormalizeRegion(region)).WithCredential(credential).WithProfile(clientProfile)

	skipVerify := PtrBool(cfg.InsecureSkipVerify) || insecure == "1" || strings.EqualFold(insecure, "true")

	if skipVerify {
		plugin.Logger(ctx).Warn("tencentcloud: TLS certificate verification is disabled (insecure_skip_verify=true); requests are vulnerable to man-in-the-middle attacks. Use only for local mock/testing environments.")
	}

	if transport := SharedHTTPTransport(ctx, d, overrides, skipVerify); transport != nil {
		c.WithHttpTransport(transport)
	}
	return nil
}

// ResolveCredentialInfo Resolve credential info from config.
// ConfigFile > Env > TCCLI credential > Credential Profile > CvmRole
func ResolveCredentialInfo(cfg TencentcloudConfig) (secretId, secretKey, token string, err error) {
	secretId = strings.TrimSpace(os.Getenv("TENCENTCLOUD_SECRET_ID"))
	secretKey = strings.TrimSpace(os.Getenv("TENCENTCLOUD_SECRET_KEY"))
	token = strings.TrimSpace(os.Getenv("TENCENTCLOUD_TOKEN"))
	if token == "" {
		token = strings.TrimSpace(os.Getenv("TENCENTCLOUD_SECURITY_TOKEN"))
	}

	if cfg.SecretId != nil && strings.TrimSpace(*cfg.SecretId) != "" {
		secretId = strings.TrimSpace(*cfg.SecretId)
	}
	if cfg.SecretKey != nil && strings.TrimSpace(*cfg.SecretKey) != "" {
		secretKey = strings.TrimSpace(*cfg.SecretKey)
	}
	if cfg.Token != nil && strings.TrimSpace(*cfg.Token) != "" {
		token = strings.TrimSpace(*cfg.Token)
	}

	if secretId != "" && secretKey != "" {
		return secretId, secretKey, token, nil
	}

	secretId, secretKey, token, found, err := loadTCCLICredential()
	if err != nil {
		return "", "", "", err
	}
	if found {
		return secretId, secretKey, token, nil
	}

	cred, err := common.DefaultProviderChain().GetCredential()
	if err != nil {
		return "", "", "", stderrors.New("tencentcloud: missing credentials, set secret_id/secret_key[/token] in the connection config")
	}
	return cred.GetSecretId(), cred.GetSecretKey(), cred.GetToken(), nil
}

type tccliCredential struct {
	SecretID  string `json:"secretId"`
	SecretKey string `json:"secretKey"`
	Token     string `json:"token"`
}

func loadTCCLICredential() (secretId, secretKey, token string, found bool, err error) {
	homeDir, err := os.UserHomeDir()
	if err != nil || strings.TrimSpace(homeDir) == "" {
		return "", "", "", false, nil
	}

	path := filepath.Join(homeDir, ".tccli", "default.credential")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return "", "", "", false, nil
	}
	if err != nil {
		return "", "", "", false, fmt.Errorf("tencentcloud: failed to read TCCLI credential file %s: %w", path, err)
	}

	var credential tccliCredential
	if err = json.Unmarshal(data, &credential); err != nil {
		return "", "", "", false, fmt.Errorf("tencentcloud: failed to parse TCCLI credential file %s: %w", path, err)
	}

	secretId = strings.TrimSpace(credential.SecretID)
	secretKey = strings.TrimSpace(credential.SecretKey)
	token = strings.TrimSpace(credential.Token)
	if secretId == "" || secretKey == "" {
		return "", "", "", false, nil
	}
	return secretId, secretKey, token, true, nil
}

// CreateCredential Create credential with credential info.
// ConfigFile > Env > TCCLI credential > Credential Profile > CvmRole
func CreateCredential(cfg TencentcloudConfig) (common.CredentialIface, error) {
	secretId, secretKey, token, err := ResolveCredentialInfo(cfg)
	if err != nil {
		return nil, err
	}

	if roleArn := ResolveRoleArn(cfg); roleArn != "" {
		credential := common.NewCredential(secretId, secretKey)
		if token != "" {
			credential = common.NewTokenCredential(secretId, secretKey, token)
		}
		return common.DefaultRoleArnProviderWithCredential(credential, roleArn).GetCredential()
	}
	if token != "" {
		return common.NewTokenCredential(secretId, secretKey, token), nil
	}
	return common.NewCredential(secretId, secretKey), nil
}

func ResolveRoleArn(cfg TencentcloudConfig) string {
	roleArn := strings.TrimSpace(os.Getenv("TENCENTCLOUD_ASSUME_ROLE_ARN"))
	if cfg.RoleArn != nil && strings.TrimSpace(*cfg.RoleArn) != "" {
		roleArn = strings.TrimSpace(*cfg.RoleArn)
	}
	return roleArn
}

func ParseDNSOverrides(raw string) map[string]string {
	out := map[string]string{}
	for pair := range strings.SplitSeq(raw, ",") {
		kv := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(kv) != 2 {
			continue
		}
		host := NormalizeRegion(kv[0])
		ip := strings.TrimSpace(kv[1])
		if host == "" || ip == "" {
			continue
		}
		if net.ParseIP(ip) == nil {
			continue
		}
		out[host] = ip
	}
	return out
}

func matchDNSOverride(host string, overrides map[string]string) (string, bool) {
	if ip, ok := overrides[host]; ok {
		return ip, true
	}
	bestIP, bestLen := "", -1
	for pattern, ip := range overrides {
		if strings.HasPrefix(pattern, "*.") {
			suffix := pattern[1:] // Match ".tencentcloudapi.com"
			if strings.HasSuffix(host, suffix) && len(suffix) > bestLen {
				bestIP, bestLen = ip, len(suffix)
			}
		}
	}
	if bestLen < 0 {
		return "", false
	}
	return bestIP, true
}

// BuildCustomTransport Build customize HTTP transport
func BuildCustomTransport(overrides map[string]string, insecureSkipVerify bool) http.RoundTripper {
	if len(overrides) == 0 && !insecureSkipVerify {
		return nil
	}
	return buildHTTPTransport(overrides, insecureSkipVerify)
}

// SharedHTTPTransport returns an HTTP transport shared by all clients for the same Steampipe connection.
func SharedHTTPTransport(ctx context.Context, d *plugin.QueryData, overrides map[string]string, insecureSkipVerify bool) *http.Transport {
	if d == nil || d.ConnectionCache == nil {
		return buildHTTPTransport(overrides, insecureSkipVerify)
	}

	cacheKey := getTransportCacheKey(overrides, insecureSkipVerify)
	if transport := cachedHTTPTransport(ctx, d, cacheKey); transport != nil {
		return transport
	}

	sharedHTTPTransportMu.Lock()
	defer sharedHTTPTransportMu.Unlock()

	// Another query may have populated the connection cache while this query was waiting for the initialization lock.
	if transport := cachedHTTPTransport(ctx, d, cacheKey); transport != nil {
		return transport
	}

	transport := buildHTTPTransport(overrides, insecureSkipVerify)
	if transport == nil {
		return nil
	}
	if err := d.ConnectionCache.SetWithTTL(ctx, cacheKey, transport, 0); err != nil {
		plugin.Logger(ctx).Warn("tencentcloud: failed to cache shared HTTP transport; using an uncached transport", "error", err)
	}
	return transport
}

func cachedHTTPTransport(ctx context.Context, d *plugin.QueryData, cacheKey string) *http.Transport {
	cached, ok := d.ConnectionCache.Get(ctx, cacheKey)
	if !ok {
		return nil
	}
	transport, ok := cached.(*http.Transport)
	if ok {
		return transport
	}

	// Delete the invalid cache entry
	d.ConnectionCache.Delete(ctx, cacheKey)
	return nil
}

func getTransportCacheKey(overrides map[string]string, insecureSkipVerify bool) string {
	hosts := make([]string, 0, len(overrides))
	for host := range overrides {
		hosts = append(hosts, host)
	}
	sort.Strings(hosts)

	hasher := sha256.New()
	if insecureSkipVerify {
		hasher.Write([]byte{1})
	} else {
		hasher.Write([]byte{0})
	}
	for _, host := range hosts {
		hasher.Write([]byte(host))
		hasher.Write([]byte{0})
		hasher.Write([]byte(overrides[host]))
		hasher.Write([]byte{0})
	}
	return HTTPTransportCacheKeyPrefix + hex.EncodeToString(hasher.Sum(nil))
}

func buildHTTPTransport(overrides map[string]string, insecureSkipVerify bool) *http.Transport {
	base, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return nil
	}
	tr := base.Clone()
	tr.IdleConnTimeout = HTTPTransportIdleTimeout

	if len(overrides) > 0 {
		dialer := &net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}
		tr.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}
			if ip, ok := matchDNSOverride(host, overrides); ok {
				addr = net.JoinHostPort(ip, port)
			}
			return dialer.DialContext(ctx, network, addr)
		}
	}

	if insecureSkipVerify {
		tlsCfg := tr.TLSClientConfig
		if tlsCfg == nil {
			tlsCfg = &tls.Config{}
		}
		tlsCfg = tlsCfg.Clone()
		tlsCfg.InsecureSkipVerify = true
		tr.TLSClientConfig = tlsCfg
	}

	return tr
}

func applyEndpoint(ctx context.Context, clientProfile *profile.ClientProfile, endpoint string, cfg TencentcloudConfig) error {
	if !strings.Contains(endpoint, "://") {
		clientProfile.HttpProfile.Endpoint = endpoint
		return nil
	}

	parsed, err := url.Parse(endpoint)
	if err != nil {
		return stderrors.New("invalid Tencent Cloud endpoint: " + err.Error())
	}
	if parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return stderrors.New("invalid Tencent Cloud endpoint: expected an http or https URL")
	}

	if parsed.Scheme == "http" {
		if !PtrBool(cfg.AllowInsecureEndpoint) && !PtrBool(cfg.InsecureSkipVerify) {
			return stderrors.New("insecure Tencent Cloud endpoint: http:// scheme is not allowed by default; set allow_insecure_endpoint = true (and insecure_skip_verify = true for TLS bypass) in the connection config to override, or use https://")
		}
		plugin.Logger(ctx).Warn("tencentcloud: using insecure http:// endpoint — API requests will be sent in cleartext, including the Authorization header with the SecretId. Use https:// for production.")
	}

	clientProfile.HttpProfile.Scheme = strings.ToUpper(parsed.Scheme)
	clientProfile.HttpProfile.Endpoint = parsed.Host + parsed.EscapedPath()
	return nil
}
