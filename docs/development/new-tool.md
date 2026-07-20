# Adding a New MCP Tool

## Steps

1. Add method to `Svc` interface in `internal/server/mcp/handler.go`
2. Implement in `internal/service/`
3. Add handler in `handler.go`
4. Register tool in `mcp.go`
5. Create markdown description in `internal/service/definitions/`
6. Generate mocks: `go generate ./...`
7. Write tests

## Svc Interface

```go
type Svc interface {
    // ... existing methods
    MyNewTool(ctx context.Context, params map[string]any) (*mcp.CallToolResult, error)
}
```

## Handler

```go
func (h *handler) handleMyNewTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    params := req.Arguments
    result, err := h.svc.MyNewTool(ctx, params)
    if err != nil {
        return nil, err
    }
    return mcp.NewToolResultText(result), nil
}
```

## Registration

```go
func (s *Server) registerTools() {
    s.addTool("my_new_tool", "Description of my new tool", inputSchema)
}
```

## Markdown Description

```markdown
# my_new_tool

Description of the tool.

## Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `param1` | string | Description |
```
