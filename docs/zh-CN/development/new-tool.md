# 添加新的 MCP 工具

## 步骤

1. **在 `internal/service/service.go` 中添加工具名称常量**
2. **在 `internal/service/types.go` 中创建请求/响应类型**
3. **在 `internal/service/` 中实现服务**（新文件或添加到现有文件）
4. **在 `internal/service/definitions/` 中创建 markdown 定义** — 这是 `MakeToolDefinitions` 读取的内容
5. **在 `internal/server/mcp/handler.go` 中向 `Svc` 接口添加方法**
6. **在 `handler.go` 中添加处理程序**
7. **在 `mcp.go` 的 `registerTools` 中注册工具**
8. **生成模拟**：`go generate ./...`
9. **编写测试**

## 1. 工具名称常量

在 `internal/service/service.go` 中添加常量：

```go
const MyNewTool = "my_new_tool"
```

## 2. 请求/响应类型

在 `internal/service/types.go` 中定义：

```go
type MyNewToolRequest struct {
    Param1 string `json:"param1" validate:"required" jsonschema:"required,Description of param1"`
}

type MyNewToolResponse struct {
    Result string `json:"result"`
}
```

## 3. 服务实现

创建 `internal/service/my_new_tool.go` 或添加到现有服务文件。遵循标准服务模式：验证 → 查找 → 执行 → 返回：

```go
func (s *Service) MyNewTool(ctx context.Context, req MyNewToolRequest) (MyNewToolResponse, error) {
    if err := s.validateRequest(req); err != nil {
        return MyNewToolResponse{}, NewLLMError(validationFailedErrCode, err.Error())
    }
    // 业务逻辑
    return MyNewToolResponse{Result: "ok"}, nil
}
```

## 4. Markdown 定义

创建 `internal/service/definitions/my_new_tool.md`。此文件由 `MakeToolDefinitions()` 读取并嵌入到二进制文件中。frontmatter 的 `name:` 字段必须与常量匹配：

```markdown
---
name: my_new_tool
---

# my_new_tool

工具的说明。

## 参数

| 参数 | 类型 | 描述 |
|------|------|------|
| `param1` | string | 描述 |
```

`tools.go` 中的 `MakeToolDefinitions()` 函数从嵌入的 `definitions/` 目录读取所有 `.md` 文件，解析 YAML frontmatter 中的 `name` 字段，并将正文用作工具描述。`instruction.md` 文件被特殊处理 — 它成为 LLM 的系统指令。

## 5. Svc 接口

在 `handler.go` 中向组合的 `Svc` 接口添加方法：

```go
type Svc interface {
    // ... 现有方法
    MyNewTool(ctx context.Context, req service.MyNewToolRequest) (service.MyNewToolResponse, error)
}
```

## 6. 处理程序

在 `handler.go` 中添加 `handler` 上的处理程序方法。处理程序委托给服务并将结果包装在 `StructuredContent` 中：

```go
func (h *handler) handleMyNewTool(
    ctx context.Context,
    _ *sdkmcp.CallToolRequest,
    req service.MyNewToolRequest,
) (*sdkmcp.CallToolResult, any, error) {
    resp, err := h.service.MyNewTool(ctx, req)
    if err != nil {
        return nil, nil, err
    }
    return &sdkmcp.CallToolResult{
        StructuredContent: resp,
    }, nil, nil
}
```

## 7. 注册

在 `mcp.go` 的 `registerTools` 函数中注册工具。向 `toolRegistrations` 映射添加条目：

```go
service.MyNewTool: {
    addTool[service.MyNewToolRequest](mcpServer, h.handleMyNewTool),
    true, // 如果工具是可变的（如 invoke 或 auth），则为 false
},
```

`registerTools` 函数签名是：

```go
func registerTools(mcpServer *sdkmcp.Server, tools []service.Tool, h handler) {
```

它遍历 `MakeToolDefinitions()` 返回的工具定义，并使用其类型化处理程序注册每个工具。`toolRegistrations` 映射将工具名称常量连接到它们的处理程序。
