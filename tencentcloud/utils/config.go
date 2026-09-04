package utils

import (
	"fmt"
	"net"
	"path"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// TencentcloudConfig 插件配置
type TencentcloudConfig struct {
	// Credential
	SecretId  *string `hcl:"secret_id,optional"`
	SecretKey *string `hcl:"secret_key,optional"`
	Token     *string `hcl:"token,optional"`

	// Region
	Regions []string `hcl:"regions,optional"`

	// Networks
	Endpoint              *string           `hcl:"endpoint,optional"`
	BaseUrl               *string           `hcl:"base_url,optional"`
	InsecureSkipVerify    *bool             `hcl:"insecure_skip_verify,optional"`
	AllowInsecureEndpoint *bool             `hcl:"allow_insecure_endpoint,optional"`
	DNSOverride           map[string]string `hcl:"dns_override,optional"`

	// Retry
	Timeout  *int `hcl:"timeout,optional"`
	MaxRetry *int `hcl:"max_retry,optional"`

	// Error Code
	IgnoreErrorCodes []string `hcl:"ignore_error_codes,optional"`
	RetryErrorCodes  []string `hcl:"retry_error_codes,optional"`
}

func ConfigInstance() any {
	return &TencentcloudConfig{}
}

func GetConfig(connection *plugin.Connection) TencentcloudConfig {
	if connection == nil || connection.Config == nil {
		return TencentcloudConfig{}
	}
	config, _ := connection.Config.(TencentcloudConfig)

	if config.Regions != nil {
		if len(config.Regions) == 0 {
			panic(fmt.Sprintf("Connection %s has an invalid value for 'regions', it must contain at least 1 region. Edit your connection configuration file and then restart Steampipe", connection.Name))
		}

		normalized := make([]string, 0, len(config.Regions))
		for _, region := range config.Regions {
			region = NormalizeRegion(region)
			if region == "" {
				continue
			}
			// Check pattern
			if _, err := path.Match(region, ""); err != nil {
				panic(fmt.Sprintf("Connection %s has an invalid pattern %q in 'regions': %v. Edit your connection configuration file and then restart Steampipe", connection.Name, region, err))
			}
			normalized = append(normalized, region)
		}
		if len(normalized) == 0 {
			panic(fmt.Sprintf("Connection %s has an invalid value for 'regions', it must contain at least 1 non-empty region. Edit your connection configuration file and then restart Steampipe", connection.Name))
		}
		config.Regions = normalized
	}

	if config.DNSOverride != nil {
		normalized := make(map[string]string, len(config.DNSOverride))
		for host, ip := range config.DNSOverride {
			host = NormalizeRegion(host)
			ip = strings.TrimSpace(ip)
			if host == "" {
				panic(fmt.Sprintf("Connection %s has an invalid entry in 'dns_override': the domain name must not be empty. Edit your connection configuration file and then restart Steampipe", connection.Name))
			}
			if ip == "" {
				panic(fmt.Sprintf("Connection %s has an invalid entry in 'dns_override': %q has no target address. Edit your connection configuration file and then restart Steampipe", connection.Name, host))
			}
			if net.ParseIP(ip) == nil {
				panic(fmt.Sprintf("Connection %s has an invalid entry in 'dns_override': %q has target %q which is not a valid IP address. Only IP addresses are supported as dns_override targets. Edit your connection configuration file and then restart Steampipe", connection.Name, host, ip))
			}
			normalized[host] = ip
		}
		config.DNSOverride = normalized
	}

	if config.BaseUrl != nil {
		raw := strings.TrimSpace(*config.BaseUrl)
		if raw == "" {
			panic(fmt.Sprintf("Connection %s has an invalid value for 'base_url', it must not be empty. Edit your connection configuration file and then restart Steampipe", connection.Name))
		}
		if strings.Contains(raw, "://") {
			panic(fmt.Sprintf("Connection %s has an invalid value for 'base_url' %q: it must be a bare domain (e.g. tencentcloudapi.com), not a URL with scheme. Edit your connection configuration file and then restart Steampipe", connection.Name, raw))
		}
		config.BaseUrl = &raw
	}

	return config
}

func NormalizeRegion(region string) string {
	return strings.ToLower(strings.TrimSpace(region))
}
