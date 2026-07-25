# Adding a New MCP Tool

## Steps

1. **Add a tool name constant** in `internal/service/service.go`
2. **Create request/response types** in `internal/service/types.go`
3. **Implement the service** in `internal/service/` (new file or add to existing)
4. **Create a markdown definition** in `internal/service/definitions/` — this is what `MakeToolDefinitions` reads
5. **Add method to `Svc` interface** in `internal/server/mcp/handler.go`
6. **Add handler** in `handler.go`
7. **Register tool** in `registerTools` in `mcp.go`
8. **Generate mocks**: `go generate ./...`
9. **Write tests**

## 1. Tool name constant

Add a constant in `internal/service/service.go`:

```go
const MyNewTool = "my_new_tool"
```

## 2. Request/Response types

Define in `internal/service/types.go`:

```go
type MyNewToolRequest struct {
    Param1 string `json:"param1" validate:"required" jsonschema:"required,Description of param1"`
}

type MyNewToolResponse struct {
    Result string `json:"result"`
}
```

## 3. Service implementation

Create `internal/service/my_new_tool.go` or add to an existing service file. Follow the standard service pattern: validate → lookup → execute → return:

```go
func (s *Service) MyNewTool(ctx context.Context, req MyNewToolRequest) (MyNewToolResponse, error) {
    if err := s.validateRequest(req); err != nil {
        return MyNewToolResponse{}, NewLLMError(validationFailedErrCode, err.Error())
    }
    // business logic
    return MyNewToolResponse{Result: "ok"}, nil
}
```

## 4. Markdown definition

Create `internal/service/definitions/my_new_tool.md`. This file is read by `MakeToolDefinitions()` and embedded into the binary. The frontmatter `name:` field must match the constant:

```markdown
---
name: my_new_tool
---

# my_new_tool

Description of the tool.

## Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `param1` | string | Description |
```

The `MakeToolDefinitions()` function in `tools.go` reads all `.md` files from the embedded `definitions/` directory, parses the YAML frontmatter for the `name` field, and uses the body as the tool description. The `instruction.md` file is treated specially — it becomes the system instruction for the LLM.

## 5. Svc Interface

Add a method to the composed `Svc` interface in `handler.go`:

```go
type Svc interface {
    // ... existing methods
    MyNewTool(ctx context.Context, req service.MyNewToolRequest) (service.MyNewToolResponse, error)
}
```

## 6. Handler

Add a handler method on `handler` in `handler.go`. The handler delegates to the service and wraps the result in `StructuredContent`:

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

## 7. Registration

Register the tool in the `registerTools` function in `mcp.go`. Add an entry to the `toolRegistrations` map:

```go
service.MyNewTool: {
    addTool[service.MyNewToolRequest](mcpServer, h.handleMyNewTool),
    true, // false if the tool is mutable (like invoke or auth)
},
```

The `registerTools` function signature is:

```go
func registerTools(mcpServer *sdkmcp.Server, tools []service.Tool, h handler) {
```

It iterates over the tool definitions returned by `MakeToolDefinitions()` and registers each one with its typed handler. The `toolRegistrations` map connects tool name constants to their handlers.
