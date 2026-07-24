# Tests

## Befehle

```bash
# Unit-Tests
go test ./...

# Bestimmtes Paket
go test ./internal/service/...

# Integrationstests
make integration-tests

# Abdeckung
make cover

# Alle Tests
make testall
```

## Teststruktur

```
tests/
├── main_test.go              # Einstiegspunkt
├── suite_test.go             # Suite-Einrichtung
├── suite_auth_test.go        # Auth-Tests
├── suite_config_test.go      # Konfigurationstests
├── suite_mcp_tools_test.go   # MCP-Tool-Tests
├── suite_search_test.go      # Suchtests
├── suite_ratelimit_test.go   # Ratenbegrenzertests
├── suite_response_test.go    # Antworttests
├── suite_export_test.go      # Exporttests
├── suite_import_test.go      # Importtests
├── suite_parsing_test.go     # Parsing-Tests
├── suite_transport_test.go   # Transporttests
├── suite_mock_test.go        # Mock-Server-Tests
├── suite_workspace_test.go   # Arbeitsbereichstests
├── suite_errors_test.go      # Fehlertests
└── suite_version_test.go     # Versionstests
```

## Abdeckung

Ziel: 80%+ für Kernpakete:

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

Verwendet `go.uber.org/mock` für MCP-Server-Tests:

```bash
go generate ./...
```

Generiert `internal/server/mcp/mock_svc_test.go` aus `handler.go`.

## Tabellengesteuerte Tests

```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"gültige Eingabe", "hallo", "HALLO", false},
        {"leere Eingabe", "", "", true},
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
