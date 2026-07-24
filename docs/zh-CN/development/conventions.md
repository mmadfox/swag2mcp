# Code Conventions

## Go

- **Go 1.26+**
- **gofmt** / **gofumpt** / **goimports** / **gci**
- **120 characters** per line
- **Guard clauses** instead of nested ifs
- **Naming**: `camelCase` for private, `PascalCase` for exported

## Errors

Use `LLMError` for LLM-visible errors:

```go
type LLMError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}
```

Error codes:
- `validation_failed` — invalid parameters
- `not_found` — resource not found
- `rate_limit` — rate limit exceeded
- `invoke_error` — API call error

## Interfaces

- Small interfaces (1-3 methods)
- Interface composition
- Functional options for configuration

## Testing

- Table-driven tests
- Test helpers (`newTestService()`, `seedTestData()`)
- Mocks via `go.uber.org/mock`
- 80%+ coverage for core packages

## Configuration

- YAML format
- Cascade: global → spec → collection
- Validation via `go-playground/validator`
- Environment variables via `$(VAR)`
