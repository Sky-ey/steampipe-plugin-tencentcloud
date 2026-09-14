connection "tencentcloud" {
  plugin = "tencentcloud/tencentcloud"

  # ---- Credentials ----
  # The plugin resolves credentials in this order:
  #   1. Static credentials in this file (`secret_id`, `secret_key`, plus optional `token`)
  #   2. The `TENCENTCLOUD_SECRET_ID`, `TENCENTCLOUD_SECRET_KEY`, `TENCENTCLOUD_SECURITY_TOKEN` env vars
  #   3. The TCCLI credential file `~/.tccli/default.credential`
  #   4. The credential profile `~/.tencentcloud/credentials`
  #   5. CVM instance role
  # secret_id  = "AKIDxxxxxxxxxxxxxxxxxxxxxxxx"
  # secret_key = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

  # STS token. Only needed when using STS temporary credentials
  # Also set with the `TENCENTCLOUD_SECURITY_TOKEN` env var.
  # token = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

  # Role Arn. When set, the plugin calls STS AssumeRole with the source credentials resolved above
  # and uses the returned STS temporary credentials to access resources in specified role.
  # It can also be set with the `TENCENTCLOUD_ASSUME_ROLE_ARN` environment variable.
  # role_arn = "qcs::cam::uin/100000000001:roleName/steampipe-read-only"

  # ---- Regions ----
  # `regions` specifies the list of regions to query.
  # SQL conditions like `WHERE region = 'ap-guangzhou'` ignore this setting.
  # The plugin queries these regions concurrently and merges the results.
  # Wildcards are supported: `*` matches any number of characters,
  # `?` matches a single character, `[gs]` matches a set of characters,
  # and `[^g]` means exclusion (note: use `^` instead of `!`).
  # Brace expansion like `{a,b}` is not supported.
  #
  # The plugin resolves the region in this order:
  #   1. The `regions` list in the connection config
  #   2. The `TENCENTCLOUD_REGION` environment variable
  #   3. The default region `ap-singapore`, resources in other regions will not appear in query results
  #
  # regions = ["*"]                           # all available regions
  # regions = ["ap-*"]                        # all Asia Pacific regions
  # regions = ["ap-guangzhou", "ap-shanghai"] # specific regions
  # regions = ["ap-*", "eu-frankfurt"]        # mix wildcards with exact values

  # ---- Networking ----
  # Custom API endpoint (a full domain).
  # Mutually exclusive with `base_url`; when both are set, `endpoint` takes precedence.
  # http:// endpoints are rejected unless `allow_insecure_endpoint = true`
  #endpoint = "cvm.tencentcloudapi.com"

  # Custom base domain (a bare domain, e.g. "tencentcloudapi.com" — no scheme, no path).
  # Each service resolves its endpoint as `<service>.<base_url>`.
  # Mutually exclusive with `endpoint`; `endpoint` takes precedence when both are set.
  #base_url = "tencentcloudapi.com"

  # Set to `true` to allow plain `http://` endpoints in `endpoint`. Off by default
  # to prevent accidentally sending credentials over an unencrypted connection.
  #allow_insecure_endpoint = false

  # Set to `true` to skip TLS certificate verification. Also required (together
  # with `allow_insecure_endpoint = true`) when using an `http://` endpoint.
  #insecure_skip_verify = false

  # Custom DNS resolution for API hosts. Keys are hostnames
  # "module.example.com" for an exact match, "*.example.com" for suffix matching
  # Values must be valid IP addresses.
  #dns_override = {
  #  "cvm.tencentcloudapi.com" = "10.0.0.1"
  #}

  # ---- Timeout & retry ----
  # Request timeout in seconds. Defaults to the SDK built-in value.
  #timeout = 10

  # Maximum number of retry attempts for failing API calls (>= 1). Defaults to 3.
  # Applied to both network failures and rate limit errors with exponential backoff.
  #max_retry = 5

  # ---- Error handling ----
  # Additional Tencent Cloud error codes to ignore for all queries, on top of the
  # plugin's default not-found codes (e.g. `ResourceNotFound`, `.NotFound`).
  #ignore_error_codes = ["AuthFailure", "UnauthorizedOperation"]

  # Additional Tencent Cloud error codes to retry for all queries, on top of the
  # plugin's default retryable codes (e.g. `RequestLimitExceeded`, `InternalError`,
  # `ServiceUnavailable`). Useful for flaky internal services that report a
  # business error code you want to retry automatically.
  #retry_error_codes = ["LimitExceeded", "ResourceInUse"]
}

# ---- Multi-account configuration (optional) ----
# Define one connection per account, then an aggregator that merges them.
# Every table exposes `owner_uin` (OwnerUin) as a connection key column, so
# filtering an aggregator query by `where owner_uin = '...'` skips the
# non-matching child connections entirely.
#
# connection "tencentcloud_dev" {
#   plugin     = "tencentcloud/tencentcloud"
#   secret_id  = "AKIDxxxxxxxx"
#   secret_key = "xxxxxxxxxxxx"
#   regions    = ["ap-guangzhou", "ap-shanghai"]
# }
#
# connection "tencentcloud_prod" {
#   plugin     = "tencentcloud/tencentcloud"
#   secret_id  = "AKIDyyyyyyyy"
#   secret_key = "yyyyyyyyyyyy"
#   regions    = ["*"]
# }
#
# connection "tencentcloud_all" {
#   plugin      = "tencentcloud/tencentcloud"
#   type        = "aggregator"
#   connections = ["tencentcloud_*"]
# }
#
# Example query against the aggregator:
#
#   select owner_uin, region, count(*)
#   from tencentcloud_all.tencentcloud_cvm_instance
#   group by owner_uin, region;
