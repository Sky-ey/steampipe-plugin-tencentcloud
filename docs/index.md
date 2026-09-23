---
organization: TencentCloud
category: ["public cloud"]
engines: ["steampipe", "sqlite", "postgres", "export"]
icon_url: "/images/plugins/tencentcloud/tencentcloud.svg"
brand_color: "#0052D9"
display_name: "Tencent Cloud"
short_name: "tencentcloud"
description: "Steampipe plugin for querying Tencent Cloud resources."
og_description: "Query Tencent Cloud with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/tencentcloud/tencentcloud-social-graphic.png"
---

# Tencent Cloud + Steampipe

[Steampipe](https://steampipe.io) is an open-source zero-ETL engine to instantly query cloud APIs using SQL.

[Tencent Cloud](https://www.tencentcloud.com/) provides on-demand cloud computing platforms and APIs to authenticated customers on a metered pay-as-you-go basis.

For example:

```sql
select
  instance_id,
  instance_name,
  instance_state,
  instance_type,
  region
from
  tencentcloud_cvm_instance
```

```text
+----------------+---------------+----------------+---------------+--------------+
| instance_id    | instance_name | instance_state | instance_type | region       |
+----------------+---------------+----------------+---------------+--------------+
| ins-2x7hq8k1   | web-server-1  | RUNNING        | S5.MEDIUM4    | ap-guangzhou |
| ins-9mtdp3wz   | web-server-2  | RUNNING        | S5.LARGE8     | ap-guangzhou |
| ins-4kpqvn82   | batch-worker  | STOPPED        | SA2.MEDIUM4   | ap-shanghai  |
| ins-7rjcx5yd   | api-gateway   | RUNNING        | S5.MEDIUM8    | ap-singapore |
+----------------+---------------+----------------+---------------+--------------+
```

## Documentation

- [**Table definitions & examples →**](https://hub.steampipe.io/plugins/tencentcloud/tencentcloud/tables)

## Get started

### Install

Download and install the latest Tencent Cloud plugin:

```bash
steampipe plugin install tencentcloud/tencentcloud
```

### Credentials

The plugin authenticates with a Tencent Cloud [access key](https://www.tencentcloud.com/document/product/598/34227) (a `SecretId` / `SecretKey` pair). Create or manage your keys in the console under **Cloud Access Management → Access Keys → [API Keys](https://console.tencentcloud.com/cam/capi)**.

Credentials are resolved in this order:

1. Static credentials in config file (`secret_id` / `secret_key`, plus optional `token` for STS temporary credentials)
2. The `TENCENTCLOUD_SECRET_ID`, `TENCENTCLOUD_SECRET_KEY` and `TENCENTCLOUD_TOKEN` (or `TENCENTCLOUD_SECURITY_TOKEN`) environment variables
3. The TCCLI credential file `~/.tccli/default.credential`
4. The credential profile `~/.tencentcloud/credentials`
5. CVM instance role when running on a Tencent Cloud CVM with an attached CAM role

#### Required permissions

All tables are read-only: each service only requires the corresponding `Describe*` / `List*` API permissions. The simplest setup is to attach the per-service read-only preset policies (e.g. `QcloudCVMReadOnlyAccess`, `QcloudCBSReadOnlyAccess`, `QcloudVPCReadOnlyAccess`) to the CAM user or role used by the plugin. See the [CAM documentation](https://www.tencentcloud.com/document/product/598) for policy management.

#### STS temporary credentials

Besides permanent access keys, Tencent Cloud also supports [STS temporary credentials](https://www.tencentcloud.com/document/product/1150) (a `SecretId` / `SecretKey` / `Token` triple). Temporary credentials expire automatically: even if leaked, they stop working on their own once the lifetime ends, which reduces at the source the security risks of long-term credential exposure.

#### Role Assumption

To access resources in another account, set `role_arn` (or `TENCENTCLOUD_ASSUME_ROLE_ARN`) to have the plugin call [AssumeRole](https://www.tencentcloud.com/document/product/1150/47148) with the resolved source credentials. The Tencent Cloud SDK manages the returned STS temporary credentials and refreshes them as needed (refresh only applies to credentials obtained via `role_arn`, not to externally supplied ones.)

### Configuration

Installing the latest Tencent Cloud plugin will create a config file (`~/.steampipe/config/tencentcloud.spc`) with a single connection named `tencentcloud`:

```hcl
connection "tencentcloud" {
  plugin = "tencentcloud/tencentcloud"

  # ---- Credentials ----
  # The plugin resolves credentials in this order:
  #   1. Static credentials in this file (`secret_id`, `secret_key`, plus optional `token`)
  #   2. The `TENCENTCLOUD_SECRET_ID`, `TENCENTCLOUD_SECRET_KEY`, `TENCENTCLOUD_TOKEN` (or `TENCENTCLOUD_SECURITY_TOKEN`) env vars
  #   3. The TCCLI credential file `~/.tccli/default.credential`
  #   4. The credential profile `~/.tencentcloud/credentials`
  #   5. CVM instance role
  # secret_id  = "AKIDxxxxxxxxxxxxxxxxxxxxxxxx"
  # secret_key = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

  # STS token. Only needed when using STS temporary credentials
  # Also set with the `TENCENTCLOUD_TOKEN` (or `TENCENTCLOUD_SECURITY_TOKEN`) env var.
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
  # This will not applied for STS credential request or refresh
  #allow_insecure_endpoint = false

  # Set to `true` to skip TLS certificate verification. Also required (together
  # with `allow_insecure_endpoint = true`) when using an `http://` endpoint.
  # This will not applied for STS credential request or refresh
  #insecure_skip_verify = false

  # Custom DNS resolution for API hosts. Keys are hostnames
  # "module.example.com" for an exact match, "*.example.com" for suffix matching
  # Values must be valid IP addresses.
  # This will not applied for STS credential request or refresh
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
```

By default, all options are commented out in the default connection, so Steampipe will resolve your credentials from environment variables, the TCCLI credential file, the SDK credentials file, or an attached CVM instance role. This provides a quick way to get started, but you will probably want to customize your experience using the configuration options for [querying multiple regions](#multi-region-connections) and [querying multiple accounts](#multi-account-connections).

Start querying:

```bash
steampipe query
```

## Multi-Region Connections

Most tables are regional. A Tencent Cloud region determines which regional control plane the API call is routed to (it is not a filter applied to a global result set). Resource metadata is isolated per region, so a query sent to the wrong region returns an empty result instead of an error.

### The `regions` argument

`regions` lists the regions each query should target. API calls are made to all matching regions in parallel and the rows are merged into a single result set.

```hcl
connection "tencentcloud" {
  plugin = "tencentcloud/tencentcloud"

  secret_id  = "your-secret-id"
  secret_key = "your-secret-key"

  # regions = ["*"]                           # All available regions
  # regions = ["ap-*"]                        # All Asia Pacific regions
  # regions = ["ap-guangzhou", "ap-shanghai"] # Specific regions
  regions = ["ap-*", "eu-frankfurt"]          # Wildcards and exact values may be mixed
}
```

Supported wildcards:

| Pattern | Meaning | Example |
|---------|---------|---------|
| `*` | Any number of characters | `ap-*` matches all Asia Pacific regions |
| `?` | Exactly one character | `ap-?ongkong` matches `ap-hongkong` |
| `[gs]` | Any character in the set | `ap-[gs]*` matches `ap-guangzhou`, `ap-shanghai`, `ap-singapore` |
| `[^g]` | Any character not in the set | `ap-[^g]*` excludes `ap-guangzhou` |

Brace expansion (`{a,b}`) is **not** supported. Values are case-insensitive.

Patterns are matched against the regions returned by the region API for that specific product, so `regions` can only narrow the set — it cannot add regions that do not exist or that the product does not support. If `regions` matches no available region, the query fails with an explicit error rather than returning zero rows.

### When `regions` is omitted

The plugin queries only the default region `ap-singapore` and logs a warning:

```
You are currently using default region ap-singapore, instances in other regions will not be displayed here!
```

Resources in every other region are absent from the results. Set `regions` explicitly to query more than one region.

### Narrowing a single query

Every regional table exposes a `region` column. Adding a `region` filter targets that one region and skips region discovery entirely, regardless of the `regions` configuration:

```sql
select
  instance_id,
  instance_name,
  region
from
  tencentcloud_cvm_instance
where
  region = 'ap-shanghai';
```

The catalog table `tencentcloud_cvm_region` is global and does not iterate over regions.

## Multi-Account Connections

Every table exposes an `owner_uin` column holding the Tencent Cloud account UIN that owns the resource, plus the companion `caller_uin` (the calling identity) and `owner_app_id` columns.

To query multiple accounts in one place, define one connection per account and an `aggregator` connection that combines them:

```hcl
# Development account
connection "tencentcloud_dev" {
  plugin     = "tencentcloud/tencentcloud"
  secret_id  = "AKIDxxxxxxxx"
  secret_key = "xxxxxxxxxxxx"
  regions    = ["ap-guangzhou", "ap-shanghai"]
}

# Production account
connection "tencentcloud_prod" {
  plugin     = "tencentcloud/tencentcloud"
  secret_id  = "AKIDyyyyyyyy"
  secret_key = "yyyyyyyyyyyy"
  regions    = ["*"]
}

# Aggregator: merges rows from every matching child connection
connection "tencentcloud_all" {
  plugin      = "tencentcloud/tencentcloud"
  type        = "aggregator"
  connections = ["tencentcloud_*"]
}
```

Query the aggregator to see rows from every account at once:

```sql
-- Count CVM instances per account, per region
select
  owner_uin,
  region,
  count(*) as instance_count
from
  tencentcloud_all.tencentcloud_cvm_instance
group by
  owner_uin, region;

-- Filter to a single account — Steampipe skips the other child connections,
-- avoiding unnecessary API calls against the wrong account
select instance_id, instance_name, owner_uin
from tencentcloud_all.tencentcloud_cvm_instance
where owner_uin = '100012345678';
```

| Column | Source | Description |
|--------|--------|-------------|
| `owner_uin` | `OwnerUin` | The master account UIN that owns the resource. |
| `caller_uin` | `Uin` | The calling identity UIN (the master account when using root credentials), or the sub-account/role UIN when called via Role Arn. |
| `owner_app_id` | `AppId` | The numeric application ID. Matches the `<appid>` suffix of COS bucket names (e.g. `my-bucket-1250000000`). |

A sub-account connection and its master-account connection share the same `owner_uin`. Use `caller_uin` to distinguish the calling identity when needed.
