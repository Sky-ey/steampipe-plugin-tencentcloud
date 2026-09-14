---
organization: TencentCloud
category: ["public cloud"]
engines: ["steampipe", "sqlite", "postgres", "export"]
icon_url: "/images/plugins/tencentcloud/tencentcloud.svg"
brand_color: "#0052D9"
display_name: "Tencent Cloud"
short_name: "tencentcloud"
description: "用于查询腾讯云资源的 Steampipe 插件。"
og_description: "使用 SQL 查询腾讯云！开源 CLI，无需数据库。"
og_image: "/images/plugins/tencentcloud/tencentcloud-social-graphic.png"
---

# Tencent Cloud 插件

[Steampipe](https://steampipe.io) 是一个开源的零 ETL 引擎，可以使用 SQL 即时查询云 API。

[腾讯云](https://www.tencentcloud.com/)面向通过身份验证的客户提供按需云计算平台和 API，按用量计费。

例如：

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

## 文档

- [**表定义与示例 →**](https://hub.steampipe.io/plugins/tencentcloud/tencentcloud/tables)

## 快速开始

### 安装

下载并安装最新的腾讯云插件：

```bash
steampipe plugin install tencentcloud/tencentcloud
```

### 凭证

插件使用腾讯云[访问密钥](https://www.tencentcloud.com/document/product/598/34227)进行身份认证——即 `SecretId` / `SecretKey` 密钥对。可在控制台 **访问管理 → API 密钥管理**（[CAM 控制台](https://console.tencentcloud.com/cam/capi)）中创建或管理密钥。

凭证按以下顺序解析：

1. 连接配置中的静态凭证（`secret_id` / `secret_key`，以及用于 STS 临时凭证的可选 `token`）
2. 环境变量 `TENCENTCLOUD_SECRET_ID`、`TENCENTCLOUD_SECRET_KEY` 与 `TENCENTCLOUD_SECURITY_TOKEN`
3. 在绑定了 CAM 角色的腾讯云 CVM 实例上运行时，使用 CVM 实例角色

#### 所需权限

所有表均为只读：每个服务只需要对应的 `Describe*` / `List*` API 权限。最简单的配置方式是给插件使用的 CAM 用户或角色绑定各服务的只读预设策略（如 `QcloudCVMReadOnlyAccess`、`QcloudCBSReadOnlyAccess`、`QcloudVPCReadOnlyAccess`）。策略管理详见 [CAM 文档](https://www.tencentcloud.com/document/product/598)。

#### STS 临时凭证

跨账号访问时，可通过 [STS](https://www.tencentcloud.com/document/product/1150) 生成临时凭证，并设置 `secret_id`、`secret_key` 与 `token`（或 `TENCENTCLOUD_SECURITY_TOKEN`）。临时凭证到期后不会自动刷新——一旦过期，插件的每次 API 调用都会返回 `AuthFailure` 错误。此时需获取新的临时凭证并更新配置。

### 配置

安装最新的腾讯云插件后，会生成一个配置文件（`~/.steampipe/config/tencentcloud.spc`），其中包含一个名为 `tencentcloud` 的连接：

```hcl
connection "tencentcloud" {
  plugin = "tencentcloud/tencentcloud"

  # ---- 凭证 ----
  # 插件按以下顺序解析凭证：
  #   1. 配置文件中的静态凭证（`secret_id`、`secret_key`，以及可选的 `token`）
  #   2. 环境变量 `TENCENTCLOUD_SECRET_ID`、`TENCENTCLOUD_SECRET_KEY`、`TENCENTCLOUD_SECURITY_TOKEN`
  #   3. SDK DefaultProviderChain：环境变量 -> 凭证文件 -> CVM 实例角色
  # 省略 `secret_id`/`secret_key` 时由默认凭证链接管——这是生产环境的推荐模式
  # （使用共享凭证文件或绑定的 CVM 实例角色）。
  # secret_id  = "AKIDxxxxxxxxxxxxxxxxxxxxxxxx"
  # secret_key = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

  # STS 临时令牌，与 `secret_id`/`secret_key` 配合用于跨账号访问（通过 STS 角色扮演）。
  # 也可通过 `TENCENTCLOUD_SECURITY_TOKEN` 环境变量设置。
  # token = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

  # ---- 地域 ----
  # `regions` 指定需要查询的地域列表。
  # SQL 条件（如 `WHERE region = 'ap-guangzhou'`）会忽略此设置。
  # 插件会并发查询这些地域并合并结果。
  # 支持通配符：`*` 匹配任意数量字符，
  # `?` 匹配单个字符，`[gs]` 匹配一组字符，
  # `[^g]` 表示排除。
  #
  # 插件按以下顺序解析凭证：
  #   1. 配置文件中的 `regions` 列表
  #   2. `TENCENTCLOUD_REGION` 环境变量
  #   3. 默认地域 `ap-singapore`, 其他地域的资源不会显示
  #
  # regions = ["*"]                           # 所有可用地域
  # regions = ["ap-*"]                        # 所有亚太地域
  # regions = ["ap-guangzhou", "ap-shanghai"] # 指定地域
  # regions = ["ap-*", "eu-frankfurt"]        # 通配符与精确值可以混用

  # ---- 网络 ----
  # 自定义 API endpoint（完整域名）。
  # 与 `base_url` 互斥；两者都设置时 `endpoint` 优先。
  # 除非设置 `allow_insecure_endpoint = true`，否则拒绝 http:// endpoint。
  #endpoint = "cvm.tencentcloudapi.com"

  # 自定义基础域名（裸域名，如 "tencentcloudapi.com"——不含协议、不含路径）。
  # 每个服务的 endpoint 解析为 `<service>.<base_url>`。
  # 与 `endpoint` 互斥；两者都设置时 `endpoint` 优先。
  #base_url = "tencentcloudapi.com"

  # 设为 `true` 允许在 `endpoint` 中使用明文 http://。默认关闭，
  # 以防止意外通过未加密连接发送凭证。
  #allow_insecure_endpoint = false

  # 设为 `true` 跳过 TLS 证书校验。使用 http:// endpoint 时也需要
  # （与 `allow_insecure_endpoint = true` 一起）设置此项。
  #insecure_skip_verify = false

  # 为 API 主机自定义 DNS 解析。键为主机名——
  # "module.example.com" 为精确匹配，"*.example.com" 为后缀匹配。
  # 值必须是合法 IP 地址。
  #dns_override = {
  #  "cvm.tencentcloudapi.com" = "10.0.0.1"
  #}

  # ---- 超时与重试 ----
  # 请求超时时间（秒）。默认为 SDK 内置值。
  #timeout = 10

  # 失败 API 调用的最大重试次数（>= 1）。默认为 3。
  # 对网络失败和限流错误均生效，采用指数退避。
  #max_retry = 5

  # ---- 错误处理 ----
  # 在插件默认忽略的 not-found 错误码（如 `ResourceNotFound`、`.NotFound`）之外，
  # 对所有查询额外忽略的腾讯云错误码。
  #ignore_error_codes = ["AuthFailure", "UnauthorizedOperation"]

  # 在插件默认可重试错误码（如 `RequestLimitExceeded`、`InternalError`、
  # `ServiceUnavailable`）之外，对所有查询额外重试的腾讯云错误码。
  # 适用于内部不稳定服务返回某个业务错误码、希望自动重试的场景。
  #retry_error_codes = ["LimitExceeded", "ResourceInUse"]
}
```

默认连接中所有配置项都被注释掉，因此 Steampipe 会使用与腾讯云 SDK 相同的机制解析凭证（环境变量、默认凭证文件或绑定的 CVM 实例角色）。这可以作为快速上手的方式，但你可能希望通过[多地域查询](#多地域查询)和[多账号查询](#多账号查询)的配置选项来自定义使用体验。

也可以使用环境变量：

```bash
export TENCENTCLOUD_SECRET_ID=your-secret-id
export TENCENTCLOUD_SECRET_KEY=your-secret-key
export TENCENTCLOUD_REGION=ap-guangzhou
```

开始查询：

```bash
steampipe query
```

## 多地域查询

大多数表是区域性的。腾讯云地域决定了 API 调用被路由到哪个区域的控制面——它不是作用于全局结果集的过滤器。资源元数据按地域隔离，因此发送到错误地域的查询会返回空结果而非报错。

### `regions` 参数

`regions` 列出每次查询要目标的地域。API 调用会并行发往所有匹配地域，行数据合并为单一结果集。

```hcl
connection "tencentcloud" {
  plugin = "tencentcloud/tencentcloud"

  secret_id  = "your-secret-id"
  secret_key = "your-secret-key"

  # regions = ["*"]                           # 所有可用地域
  # regions = ["ap-*"]                        # 所有亚太地域
  # regions = ["ap-guangzhou", "ap-shanghai"] # 指定地域
  regions = ["ap-*", "eu-frankfurt"]          # 通配符与精确值可以混用
}
```

支持的通配符：

| 模式 | 含义 | 示例 |
|------|------|------|
| `*` | 任意数量字符 | `ap-*` 匹配所有亚太地域 |
| `?` | 恰好一个字符 | `ap-?ongkong` 匹配 `ap-hongkong` |
| `[gs]` | 集合中的任一字符 | `ap-[gs]*` 匹配 `ap-guangzhou`、`ap-shanghai`、`ap-singapore` |
| `[^g]` | 不在集合中的任一字符 | `ap-[^g]*` 排除 `ap-guangzhou` |

**不支持**花括号展开（`{a,b}`）。取值不区分大小写。

模式匹配的对象是该产品通过地域 API 返回的地域列表，因此 `regions` 只能缩小范围——不能添加不存在的或产品不支持的地域。如果 `regions` 未匹配到任何可用地域，查询会直接报错而不是返回零行。

### 省略 `regions` 时

插件只查询默认地域 `ap-singapore` 并输出一条警告：

```
You are currently using default region ap-singapore, instances in other regions will not be displayed here!
```

其他地域的资源不会出现在结果中。要查询多个地域，请显式设置 `regions`。

### 单条查询收窄范围

每张区域表都提供 `region` 列。添加 `region` 过滤条件可以只查询该地域并完全跳过地域发现流程，不受 `regions` 配置影响：

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

目录表 `tencentcloud_cvm_region` 是全局表，不会按地域遍历。

## 多账号查询

每张表都提供 `owner_uin` 列，保存拥有该资源的腾讯云主账号 UIN（`OwnerUin`），并附带 `caller_uin`（调用方身份）与 `owner_app_id` 列。插件通过调用 CAM `GetUserAppId` API 每个连接解析一次——无论查询多少张表或多少个地域，开销都只有每连接一次 API 调用。

`owner_uin` 被注册为插件的 **connection key column**，当聚合查询按账号过滤时，Steampipe 可以据此裁剪子连接。要在一个地方查询多个账号，可为每个账号定义一个连接，并定义一个组合它们的 `aggregator` 连接：

```hcl
# 开发账号
connection "tencentcloud_dev" {
  plugin     = "tencentcloud/tencentcloud"
  secret_id  = "AKIDxxxxxxxx"
  secret_key = "xxxxxxxxxxxx"
  regions    = ["ap-guangzhou", "ap-shanghai"]
}

# 生产账号
connection "tencentcloud_prod" {
  plugin     = "tencentcloud/tencentcloud"
  secret_id  = "AKIDyyyyyyyy"
  secret_key = "yyyyyyyyyyyy"
  regions    = ["*"]
}

# 聚合连接：合并所有匹配子连接的行
connection "tencentcloud_all" {
  plugin      = "tencentcloud/tencentcloud"
  type        = "aggregator"
  connections = ["tencentcloud_*"]
}
```

查询聚合连接即可同时看到所有账号的行：

```sql
-- 按账号、按地域统计 CVM 实例数
select
  owner_uin,
  region,
  count(*) as instance_count
from
  tencentcloud_all.tencentcloud_cvm_instance
group by
  owner_uin, region;

-- 过滤到单个账号——Steampipe 会跳过其他子连接，
-- 避免向错误账号发起不必要的 API 调用
select instance_id, instance_name, owner_uin
from tencentcloud_all.tencentcloud_cvm_instance
where owner_uin = '100012345678';
```

| 列 | 来源字段 | 描述 |
|----|----------|------|
| `owner_uin` | `OwnerUin` | 拥有该资源的主账号 UIN。作为 connection key column 用于聚合连接裁剪。 |
| `caller_uin` | `Uin` | 调用方身份 UIN——使用主账号凭证时即主账号，通过 CAM 调用时为子账号/角色 UIN。 |
| `owner_app_id` | `AppId` | 数字形式的 APP ID。与 COS 存储桶名称的 `<appid>` 后缀一致（如 `my-bucket-1250000000`）。 |

子账号连接与其主账号连接的 `owner_uin` 相同（都解析到同一个 `OwnerUin`）。需要区分调用方身份时使用 `caller_uin`。
