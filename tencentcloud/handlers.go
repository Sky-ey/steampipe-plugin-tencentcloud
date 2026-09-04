package tencentcloud

import (
	"context"
	stderrors "errors"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"strings"

	tcerrors "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/errors"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// Error Handling

// GetDefaultErrorIgnoreCodes 获取默认忽略的错误码
func GetDefaultErrorIgnoreCodes() []string {
	return []string{
		"InvalidInstanceId.NotFound",
		"InvalidInstanceId.Malformed",
		"InvalidZone.NotFound",
		"InvalidParameter.UserNotExist",
		"InvalidParameter.GroupNotExist",
		"InvalidParameter.PolicyIdNotExist",
		"InvalidParameter.RoleNotExist",
		"ResourceNotFound",
		"NotFound",
	}
}

func isNotFoundError(_ context.Context, d *plugin.QueryData, _ *plugin.HydrateData, err error) bool {
	if err == nil {
		return false
	}

	// 默认忽略码 + 用户在 .spc 中追加的 ignore_error_codes
	ignoreCodes := append(append([]string{}, GetDefaultErrorIgnoreCodes()...), utils.GetConfig(d.Connection).IgnoreErrorCodes...)

	var sdkErr *tcerrors.TencentCloudSDKError
	if stderrors.As(err, &sdkErr) {
		code := sdkErr.Code
		for _, c := range ignoreCodes {
			if c == "" {
				continue
			}
			if code == c || strings.HasPrefix(code, c) {
				return true
			}
		}
		return false
	}

	// Fallback
	msg := err.Error()
	for _, s := range ignoreCodes {
		if s == "" {
			continue
		}
		if strings.Contains(msg, s) {
			return true
		}
	}
	return false
}

// GetDefaultErrorRetryCodes 获取默认重试的错误码
func GetDefaultErrorRetryCodes() []string {
	return []string{
		"RequestLimitExceeded",
		"ServiceUnavailable",
		"TooManyRequests",
		"InternalError",
		"Throttling",
		"429",
		"500",
		"503",
		"504",
	}
}

// shouldRetryError returns true if the error should trigger a Steampipe-level retry
func shouldRetryError(_ context.Context, d *plugin.QueryData, _ *plugin.HydrateData, err error) bool {
	if err == nil {
		return false
	}

	// 默认重试码 + 用户在 .spc 中追加的 retry_error_codes
	retryCodes := append(append([]string{}, GetDefaultErrorRetryCodes()...), utils.GetConfig(d.Connection).RetryErrorCodes...)

	var sdkErr *tcerrors.TencentCloudSDKError
	if stderrors.As(err, &sdkErr) {
		code := sdkErr.Code
		for _, c := range retryCodes {
			if c == "" {
				continue
			}
			if code == c || strings.HasPrefix(code, c) {
				return true
			}
		}
		return false
	}

	// Fallback
	msg := err.Error()
	for _, s := range retryCodes {
		if s == "" {
			continue
		}
		if strings.Contains(msg, s) {
			return true
		}
	}
	return false
}
