# 构建

## 要求

- Go 1.26+
- Make

## 命令

```bash
# 构建
make build

# 带版本构建
make build VERSION=1.0.0

# 代码检查
make lint

# 测试
go test ./...

# 所有测试
make testall
```

## GoReleaser

用于发布：

```bash
goreleaser release --snapshot --clean
```

## 平台

| 平台 | 架构 |
|------|------|
| Linux | amd64, arm64 |
| macOS | amd64, arm64 |
| Windows | amd64 |

## 代码检查

```bash
make lint
```

使用 `golangci-lint`，启用 80+ 个检查器。配置在 `.golangci.yml` 中。
