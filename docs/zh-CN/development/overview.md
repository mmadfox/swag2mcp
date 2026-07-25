# 开发概述

## 关于此项目

swag2mcp 是一个 Go 项目，通过模型上下文协议（MCP）将 OpenAPI/Swagger/Postman 规范与 LLM 智能体连接起来。它使用 Go 1.23+ 构建，并遵循由 80+ 个检查器强制执行的严格编码规范。

本节是为想要理解代码库、贡献或使用新的认证方法、MCP 工具和集成来扩展 swag2mcp 的**工程师**编写的。

## 开发技能

项目附带两个开发技能，编码了项目的规范和模式。你可以使用它们或忽略它们 — 它们是工具，不是规则。

### godeveloper

[godeveloper 技能](https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/godeveloper/SKILL.md) 定义了项目中的每个代码规范：

- **命名** — 包、文件、类型、接口、接收器、常量
- **格式化** — gofmt/gofumpt/goimports/gci、120 行限制、导入排序
- **错误处理** — 带有 8 个错误代码的 `LLMError`、哨兵错误、错误包装
- **接口** — 小接口、组合、消费者端定义
- **并发** — 互斥锁粒度、goroutine 生命周期、上下文传递
- **测试** — 表格驱动测试、`newTestService()`/`seedTestData()` 辅助函数、模拟生成
- **项目模式** — 服务层、请求/响应结构体、函数选项、MCP 处理程序模式

### swag2mcp-cli

[swag2mcp-cli 技能](https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-cli/SKILL.md) 记录了每个 CLI 命令的语法、标志、参数和示例。在处理 CLI 命令或编写文档时很有用。

## 关键架构决策

### 服务层模式

每个功能遵循相同的三步模式：

1. **验证** 使用 `s.validateRequest(req)` 验证请求（使用 `go-playground/validator`）
2. **查找** 从内存索引中查找实体（返回带有 `not_found` 代码的 `LLMError`）
3. **执行** 执行业务逻辑并返回类型化响应或 `LLMError`

```go
func (s *Service) Search(ctx context.Context, req SearchRequest) (SearchResponse, error) {
    if err := s.validateRequest(req); err != nil {
        return SearchResponse{}, NewLLMError(validationFailedErrCode, err.Error())
    }
    results, err := s.index.Search(req.Query, req.Limit)
    if err != nil {
        return SearchResponse{}, NewLLMError(invokeErrorCode, err.Error())
    }
    return SearchResponse{Results: results}, nil
}
```

### 请求/响应结构体

每个方法有专用的 `{Method}Request` 和 `{Method}Response` 结构体。请求结构体使用 `validate` 标签进行验证，使用 `jsonschema` 标签进行文档：

```go
type SearchRequest struct {
    Query string `json:"query" validate:"required,min=1" jsonschema:"description=Search query supporting field filters"`
    Limit int    `json:"limit" validate:"required,min=1,max=50" jsonschema:"description=Maximum results"`
}

type SearchResponse struct {
    Results []EndpointSearchItem `json:"results"`
}
```

### 函数选项

配置使用函数选项模式：

```go
type Option func(*Service)

func New(opts ...Option) (*Service, error)

func WithDisableLLMAuth(disable bool) Option {
    return func(s *Service) {
        s.disableLLMAuth.Store(disable)
    }
}
```

### MCP 处理程序模式

MCP 服务器使用组合接口模式。`internal/server/mcp/handler.go` 中的 `Svc` 接口由较小的接口（`CatalogReader`、`EndpointExplorer`、`EndpointExecutor`、`SystemInfo`、`ResponseManager`）组合而成。每个处理程序方法委托给服务层：

```go
type handler struct {
    service Svc
}

func (h *handler) handleSearch(ctx context.Context, _ *sdkmcp.CallToolRequest, req service.SearchRequest) (*sdkmcp.CallToolResult, any, error) {
    resp, err := h.service.Search(ctx, req)
    if err != nil {
        return nil, nil, err
    }
    return &sdkmcp.CallToolResult{StructuredContent: resp}, nil, nil
}
```

### LLMError

所有返回给 LLM 的错误使用 `LLMError` 类型，带有 8 个代码之一：

| 代码 | 何时使用 |
|------|----------|
| `validation_failed` | 无效输入（错误的 ID 格式、缺少必需字段） |
| `not_found` | 索引中未找到实体 |
| `rate_limit` | 超过每个端点 10 秒冷却时间 |
| `invoke_error` | HTTP 请求/响应失败 |
| `config_error` | 配置加载或验证失败 |
| `workspace_error` | 工作区目录或文件操作失败 |
| `parse_error` | 规范文件解析失败 |
| `auth_error` | 认证令牌检索失败 |

消息必须用适合 LLM 消费者的通俗语言解释出了什么问题以及下一步该怎么做。

### ID 生成

所有 ID 是确定性的 MD5 哈希：

```go
id.Domain("meteo")                          // 32 字符十六进制
id.Collection("meteo", "Forecast")          // 32 字符十六进制
id.Tag("meteo", "Forecast", "pets")         // 32 字符十六进制
id.Method("meteo", "Forecast", "pets", "GET", "/v2/pet/{petId}")
```

### 配置级联

配置通过三个级别级联：**全局 → spec → collection**。每个级别覆盖前一个级别。所有 `http_client` 设置可以在每个级别被覆盖。头和 cookie 被合并；简单值被替换。

## 快速参考

| 领域 | 规范 |
|------|------|
| **Go 版本** | 1.23+ |
| **格式化工具** | gofmt、gofumpt、goimports、gci |
| **行长度** | 120 字符 |
| **检查器** | `.golangci.yml` 中 80+ 个 |
| **错误类型** | 带有 8 个代码的 `LLMError` |
| **模拟框架** | `go.uber.org/mock` |
| **测试辅助函数** | `newTestService()`、`seedTestData()` |
| **配置格式** | 带级联的 YAML |
| **认证分发** | `UnmarshalYAML` 读取 `type` 字段 |
| **ID 生成** | 基于 MD5（`id.Domain()`、`id.Collection()` 等） |
| **速率限制** | 每个端点 10 秒用于 `invoke` |
| **响应大小** | 默认 1 MB，超过时保存到文件 |
| **覆盖率目标** | 核心包 80%+ |
| **构建** | `make build` |
| **代码检查** | `make lint` |
| **测试** | `go test ./...` |
| **生成** | `go generate ./...` |
