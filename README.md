![Tencent Cloud + Steampipe](https://steampipe.io/images/plugins/tencentcloud/tencentcloud-social-graphic.png)

# Tencent Cloud + Steampipe

[Steampipe](https://steampipe.io) is an open-source zero-ETL engine to instantly query cloud APIs using SQL.

Use SQL to query CVM instances, CBS disks, VPC networks, subnets, and more from your Tencent Cloud account.

## Quick start

Install the plugin:

```bash
steampipe plugin install tencentcloud/tencentcloud
```

Configure your credentials in `~/.steampipe/config/tencentcloud.spc`:

```hcl
connection "tencentcloud" {
  plugin = "tencentcloud/tencentcloud"

  secret_id  = "your-secret-id"
  secret_key = "your-secret-key"

  # Regions to query when no region is specified. Supports wildcards (`*`, `?`, `[gs]`).
  regions = ["ap-*"]
}
```

Alternatively, use environment variables:

```bash
export TENCENTCLOUD_SECRET_ID=your-secret-id
export TENCENTCLOUD_SECRET_KEY=your-secret-key
export TENCENTCLOUD_REGION=ap-guangzhou
```

Run a query:

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

## Developing

The plugin is built with Go and the [Steampipe plugin SDK](https://github.com/turbot/steampipe-plugin-sdk).

Clone the repository and build & install the plugin via the Makefile:

```bash
git clone https://github.com/TencentCloud/steampipe-plugin-tencentcloud.git
cd steampipe-plugin-tencentcloud
make install
```

Run the fast Go unit tests:

```bash
make test
```

Run the SQL query contract suite against the local mock Tencent Cloud API. This covers List/Get queries for every table, a cross-service join, and validates API actions, paths, request bodies, regions, and request counts:

```bash
make test-query
```

The contract suite uses an isolated temporary Steampipe install directory and never writes to `~/.steampipe`.

## Open Source & Contributing

This repository is published under the [Apache 2.0 license](LICENSE).

We welcome all kinds of contributions (bug reports, documentation improvements, and code).

To get started, see [contributing guidelines](https://steampipe.io/community/contribute). Please read the [code of conduct](https://steampipe.io/community/code-of-conduct) before participating.
