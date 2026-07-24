# Pruebas

## Comandos

```bash
# Pruebas unitarias
go test ./...

# Paquete específico
go test ./internal/service/...

# Pruebas de integración
make integration-tests

# Cobertura
make cover

# Todas las pruebas
make testall
```

## Estructura de Pruebas

```
tests/
├── main_test.go              # Punto de entrada
├── suite_test.go             # Configuración del conjunto
├── suite_auth_test.go        # Pruebas de autenticación
├── suite_config_test.go      # Pruebas de configuración
├── suite_mcp_tools_test.go   # Pruebas de herramientas MCP
├── suite_search_test.go      # Pruebas de búsqueda
├── suite_ratelimit_test.go   # Pruebas de límite de velocidad
├── suite_response_test.go    # Pruebas de respuesta
├── suite_export_test.go      # Pruebas de exportación
├── suite_import_test.go      # Pruebas de importación
├── suite_parsing_test.go     # Pruebas de análisis
├── suite_transport_test.go   # Pruebas de transporte
├── suite_mock_test.go        # Pruebas del servidor simulado
├── suite_workspace_test.go   # Pruebas del espacio de trabajo
├── suite_errors_test.go      # Pruebas de errores
└── suite_version_test.go     # Pruebas de versión
```

## Cobertura

Objetivo: 80%+ para paquetes principales:

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

Usa `go.uber.org/mock` para pruebas del servidor MCP:

```bash
go generate ./...
```

Genera `internal/server/mcp/mock_svc_test.go` a partir de `handler.go`.

## Pruebas Basadas en Tablas

```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"entrada válida", "hello", "HELLO", false},
        {"entrada vacía", "", "", true},
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
