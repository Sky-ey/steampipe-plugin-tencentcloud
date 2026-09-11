# 腾讯云 Steampipe 插件

> 使用 SQL 查询腾讯云资源的 Steampipe 插件。

[English](README.md) | 简体中文

## 概述

使用 SQL 查询你腾讯云账号下的 CVM 实例、CBS 云盘、VPC 网络、子网等资源。

## 快速开始

安装插件：

```bash
steampipe plugin install tencentcloud/tencentcloud
```

在 `~/.steampipe/config/tencentcloud.spc` 中配置凭证：

```hcl
connection "tencentcloud" {
  plugin = "tencentcloud/tencentcloud"

  secret_id  = "your-secret-id"
  secret_key = "your-secret-key"

  # 未指定 region 时要查询的地域列表。支持通配符（`*`、`?`、`[gs]`）。
  # 若省略该项，SQL 中未指定 region 时仅查询默认地域 `ap-singapore`。
  regions = ["ap-*"]
}
```

也可以使用环境变量：

```bash
export TENCENTCLOUD_SECRET_ID=your-secret-id
export TENCENTCLOUD_SECRET_KEY=your-secret-key
```

运行查询：

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

区域表会查询 `regions` 匹配到的所有地域。可以通过 `where region = 'ap-guangzhou'` 这类 SQL 条件直接定位单个地域，跳过地域发现流程。完整的匹配规则、凭证解析与多账号（聚合连接）支持见 [docs/index.md](docs/index_zh.md#多地域查询)。

## 开发

插件使用 Go 和 [Steampipe plugin SDK](https://github.com/turbot/steampipe-plugin-sdk) 构建。

克隆仓库并通过 Makefile 构建安装插件：

```bash
git clone https://github.com/TencentCloud/steampipe-plugin-tencentcloud.git
cd steampipe-plugin-tencentcloud
make install
```

运行快速 Go 单元测试：

```bash
make test
```

针对本地 Mock 腾讯云 API 运行 SQL 查询契约测试套件。该套件覆盖每张表的 List/Get 查询和一个跨服务 join，并校验 API action、路径、请求体、地域以及请求次数：

```bash
make test-query
```

契约测试套件使用隔离的临时 Steampipe 安装目录，不会写入 `~/.steampipe`。

## 开源与贡献

本仓库基于 [Apache 2.0 许可证](LICENSE)发布。

欢迎贡献——issue 反馈、文档改进与代码提交都非常感谢。入门请参阅[贡献指南](https://steampipe.io/community/contribute)，参与前请阅读[行为准则](https://steampipe.io/community/code-of-conduct)。
