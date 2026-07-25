# Testing

## Commands

```bash
# Unit tests
go test ./...

# Specific package
go test ./internal/service/...

# Integration tests
make integration-tests

# Coverage
make cover

# All tests
make testall
```

## Test Structure

```
tests/
├── main_test.go              # Entry point
├── suite_test.go             # Suite setup
├── suite_auth_test.go        # Auth tests
├── suite_config_test.go      # Config tests
├── suite_mcp_tools_test.go   # MCP tools tests
├── suite_search_test.go      # Search tests
├── suite_ratelimit_test.go   # Rate limit tests
├── suite_response_test.go    # Response tests
├── suite_export_test.go      # Export tests
├── suite_import_test.go      # Import tests
├── suite_parsing_test.go     # Parsing tests
├── suite_transport_test.go   # Transport tests
├── suite_mock_test.go        # Mock server tests
├── suite_workspace_test.go   # Workspace tests
├── suite_errors_test.go      # Error tests
└── suite_version_test.go     # Version tests
```

## Coverage

Target: 80%+ for core packages:

- `auth`
- `cache`
- `config`
- `env`
- `httpclient`
- `id`
- `index`
- `server/mcp`
- `service`
- `spec`
- `workspace`

## Mocks

Uses `go.uber.org/mock` for MCP server tests:

```bash
go generate ./...
```

Generates `internal/server/mcp/mock_svc_test.go` from `handler.go`.

## Table-Driven Tests

```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "hello", "HELLO", false},
        {"empty input", "", "", true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := DoSomething(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            assert.Equal(t, tt.want, got)
        })
    }
}
```
