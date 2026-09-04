package cos

import (
	"context"
	"fmt"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"net/http"
	"net/url"
	"os"
	"strings"

	cos "github.com/tencentyun/cos-go-sdk-v5"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func buildCosClient(ctx context.Context, d *plugin.QueryData, bucketName, region string) (*cos.Client, error) {
	cfg := utils.GetConfig(d.Connection)

	secretId, secretKey, token, err := utils.ResolveCredentialInfo(cfg)
	if err != nil {
		return nil, err
	}

	overrides := utils.ParseDNSOverrides(os.Getenv("TENCENTCLOUD_DNS_OVERRIDE"))
	if cfg.DNSOverride != nil {
		for k, v := range cfg.DNSOverride {
			overrides[k] = v
		}
	}
	skipVerify := utils.PtrBool(cfg.InsecureSkipVerify) ||
		strings.EqualFold(strings.TrimSpace(os.Getenv("TENCENTCLOUD_INSECURE_SKIP_VERIFY")), "true") ||
		strings.TrimSpace(os.Getenv("TENCENTCLOUD_INSECURE_SKIP_VERIFY")) == "1"

	if skipVerify {
		plugin.Logger(ctx).Warn("tencentcloud: TLS certificate verification is disabled (insecure_skip_verify=true) for COS client; requests are vulnerable to man-in-the-middle attacks.")
	}

	inner := utils.SharedHTTPTransport(ctx, d, overrides, skipVerify)

	authTransport := &cos.AuthorizationTransport{
		SecretID:     secretId,
		SecretKey:    secretKey,
		SessionToken: token,
		Transport:    inner,
	}
	httpClient := &http.Client{Transport: authTransport}

	baseURL := &cos.BaseURL{}
	if bucketName != "" && region != "" {
		bucketURL, err := cos.NewBucketURL(bucketName, region, true)
		if err != nil {
			return nil, err
		}
		baseURL.BucketURL = bucketURL
	}

	if cfg.Endpoint != nil && strings.TrimSpace(*cfg.Endpoint) != "" {
		endpoint := strings.TrimSpace(*cfg.Endpoint)
		if !strings.Contains(endpoint, "://") {
			endpoint = "https://" + endpoint
		}
		if u, err := url.Parse(endpoint); err == nil && u.Host != "" {
			baseURL.ServiceURL = u
		}
	}

	return cos.NewClient(baseURL, httpClient), nil
}

// maxKeysFromLimit returns the COS list page size
func maxKeysFromLimit(d *plugin.QueryData) int64 {
	maxKeys := int64(1000)
	if d.QueryContext.Limit != nil && *d.QueryContext.Limit < maxKeys {
		maxKeys = *d.QueryContext.Limit
	}
	if maxKeys < 1 {
		maxKeys = 1
	}
	return maxKeys
}

func listCosBucketNames(ctx context.Context, d *plugin.QueryData) ([]string, []string, error) {
	client, err := buildCosClient(ctx, d, "", "")
	if err != nil {
		return nil, nil, err
	}

	var names, regions []string
	var marker string
	for {
		d.WaitForListRateLimit(ctx)
		opt := &cos.ServiceGetOptions{Marker: marker, MaxKeys: 1000}
		utils.LogRequest(ctx, "listCosBucketNames", opt)
		res, _, err := client.Service.Get(ctx, opt)
		if err != nil {
			return nil, nil, err
		}
		if res == nil {
			break
		}
		for _, b := range res.Buckets {
			names = append(names, b.Name)
			regions = append(regions, b.Region)
		}
		if !res.IsTruncated || res.NextMarker == "" {
			break
		}
		marker = res.NextMarker
	}
	return names, regions, nil
}

func resolveCosBucketNames(ctx context.Context, d *plugin.QueryData) ([]string, []string, error) {
	if name := d.EqualsQualString("bucket_name"); name != "" {
		region := d.EqualsQualString("region")
		if region == "" {
			region = d.EqualsQualString(utils.MatrixKeyRegion)
		}

		if region == "" {
			names, regions, err := listCosBucketNames(ctx, d)
			if err != nil {
				return nil, nil, err
			}
			for i, n := range names {
				if n == name {
					return []string{name}, []string{regions[i]}, nil
				}
			}
			return nil, nil, fmt.Errorf("querying tencentcloud_cos_object by bucket_name requires also specifying region, and the bucket %q was not found in any region", name)
		}
		return []string{name}, []string{region}, nil
	}
	return listCosBucketNames(ctx, d)
}
