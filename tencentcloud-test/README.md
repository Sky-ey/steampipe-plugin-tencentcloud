# 查询契约测试

该目录采用“SQL + expected JSON + requests JSON”结构，通过本地 `httptest.Server` 模拟 Tencent Cloud HTTP API，并验证查询结果及请求的 method、path、action、region、body 子集和精确次数。

## 覆盖范围

- 覆盖 70 张表（contract tests 目录按表分目录）。
- 除 `tencentcloud_cos_object` 外全部有契约：COS 的 Object 请求通过 `cos.NewBucketURL` 固定发往 `<bucket>.cos.<region>.myqcloud.com` 真实域名（`endpoint` 仅覆盖 `ServiceURL`），无法被本地 mock 服务器捕获，故该表不写契约；`tencentcloud_cos_bucket` 的 `GetService` 请求走 `ServiceURL`，可 mock（XML 响应经 `cosFixtures` 按 `method+path` 路由，见 query_test.go）。
- 跨服务同名 action（如 CVM 与 WAF 的 `DescribeInstances`/`DescribeHosts`）在 `responseFixtures` 中合并返回两套字段，互不影响。

测试使用临时 Steampipe install-dir，不会安装插件到 `~/.steampipe`。运行：

```sh
go test -tags=querytest ./tencentcloud-test -v
```

未安装 `steampipe` 时测试会自动跳过；默认 `go test ./...` 不包含该测试。macOS 下若 steampipe 报 `failed to get pids`，是运行环境禁止 `sysctl kern.proc` 进程枚举所致，需在不受限的终端中运行。
