# Development Overview

## About this project

swag2mcp is a Go project that bridges OpenAPI/Swagger/Postman specifications with LLM agents via the Model Context Protocol (MCP). It is built with Go 1.23+ and follows strict coding conventions enforced by 80+ linters.

This section is written for **engineers** who want to understand the codebase, contribute, or extend swag2mcp with new auth methods, MCP tools, or integrations.

## Development skills

The project ships with two development skills that encode the project's conventions and patterns. You can use them or ignore them — they are tools, not rules.

### godeveloper

The [godeveloper skill](https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/godeveloper/SKILL.md) defines every code convention in the project:

- **Naming** — packages, files, types, interfaces, receivers, constants
- **Formatting** — gofmt/gofumpt/goimports/gci, 120-line limit, import ordering
- **Error handling** — `LLMError` with 8 error codes, sentinel errors, error wrapping
- **Interfaces** — small interfaces, composition, consumer-side definitions
- **Concurrency** — mutex granularity, goroutine lifetimes, context passing
- **Testing** — table-driven tests, `newTestService()`/`seedTestData()` helpers, mock generation
- **Project patterns** — service layer, request/response structs, functional options, MCP handler pattern

### swag2mcp-cli

The [swag2mcp-cli skill](https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-cli/SKILL.md) documents every CLI command with syntax, flags, arguments, and examples. Useful when working on CLI commands or writing documentation.

## Key architectural decisions

### Service layer pattern

Every feature follows the same three-step pattern:

1. **Validate** the request with `s.validateRequest(req)` (uses `go-playground/validator`)
2. **Look up** entities from the in-memory index (returns `LLMError` with `not_found` code)
3. **Execute** business logic and return a typed response or `LLMError`

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

### Request/Response structs

Each method has a dedicated `{Method}Request` and `{Method}Response` struct. Request structs use `validate` tags for validation and `jsonschema` tags for documentation:

```go
type SearchRequest struct {
    Query string `json:"query" validate:"required,min=1" jsonschema:"description=Search query supporting field filters"`
    Limit int    `json:"limit" validate:"required,min=1,max=50" jsonschema:"description=Maximum results"`
}

type SearchResponse struct {
    Results []EndpointSearchItem `json:"results"`
}
```

### Functional options

Configuration uses the functional options pattern:

```go
type Option func(*Service)

func New(opts ...Option) (*Service, error)

func WithDisableLLMAuth(disable bool) Option {
    return func(s *Service) {
        s.disableLLMAuth.Store(disable)
    }
}
```

### MCP handler pattern

The MCP server uses a composed interface pattern. The `Svc` interface in `internal/server/mcp/handler.go` is composed from smaller interfaces (`CatalogReader`, `EndpointExplorer`, `EndpointExecutor`, `SystemInfo`, `ResponseManager`). Each handler method delegates to the service layer:

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

All errors returned to the LLM use the `LLMError` type with one of 8 codes:

| Code | When |
|------|------|
| `validation_failed` | Invalid input (wrong ID format, missing required fields) |
| `not_found` | Entity not found in index |
| `rate_limit` | Per-endpoint 10s cooldown exceeded |
| `invoke_error` | HTTP request/response failures |
| `config_error` | Configuration loading or validation failure |
| `workspace_error` | Workspace directory or file operation failure |
| `parse_error` | Spec file parsing failure |
| `auth_error` | Authentication token retrieval failure |

Messages must explain what went wrong AND what to do next, in plain language suitable for an LLM consumer.

### ID generation

All IDs are deterministic MD5 hashes:

```go
id.Domain("meteo")                          // 32-char hex
id.Collection("meteo", "Forecast")          // 32-char hex
id.Tag("meteo", "Forecast", "pets")         // 32-char hex
id.Method("meteo", "Forecast", "pets", "GET", "/v2/pet/{petId}")
```

### Config cascade

Configuration cascades through three levels: **global → spec → collection**. Each level overrides the previous. All `http_client` settings can be overridden at every level. Headers and cookies are merged; simple values are replaced.

## Quick reference

| Area | Convention |
|------|------------|
| **Go version** | 1.23+ |
| **Formatters** | gofmt, gofumpt, goimports, gci |
| **Line length** | 120 characters |
| **Linters** | 80+ in `.golangci.yml` |
| **Error type** | `LLMError` with 8 codes |
| **Mock framework** | `go.uber.org/mock` |
| **Test helpers** | `newTestService()`, `seedTestData()` |
| **Config format** | YAML with cascade |
| **Auth dispatch** | `UnmarshalYAML` reads `type` field |
| **ID generation** | MD5-based (`id.Domain()`, `id.Collection()`, etc.) |
| **Rate limit** | 10s per endpoint for `invoke` |
| **Response size** | 1 MB default, saved to file when exceeded |
| **Coverage target** | 80%+ for core packages |
| **Build** | `make build` |
| **Lint** | `make lint` |
| **Test** | `go test ./...` |
| **Generate** | `go generate ./...` |
