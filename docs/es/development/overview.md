# Resumen de Desarrollo

## Acerca de este proyecto

swag2mcp es un proyecto Go que conecta especificaciones OpenAPI/Swagger/Postman con agentes LLM mediante el Protocolo de Contexto de Modelo (MCP). Está construido con Go 1.23+ y sigue estrictas convenciones de codificación aplicadas por más de 80 linters.

Esta sección está escrita para **ingenieros** que quieran entender el código base, contribuir o ampliar swag2mcp con nuevos métodos de autenticación, herramientas MCP o integraciones.

## Habilidades de desarrollo

El proyecto incluye dos habilidades de desarrollo que codifican las convenciones y patrones del proyecto. Puede usarlas o ignorarlas — son herramientas, no reglas.

### godeveloper

La [habilidad godeveloper](https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/godeveloper/SKILL.md) define cada convención de código en el proyecto:

- **Nomenclatura** — paquetes, archivos, tipos, interfaces, receptores, constantes
- **Formato** — gofmt/gofumpt/goimports/gci, límite de 120 líneas, orden de importaciones
- **Manejo de errores** — `LLMError` con 8 códigos de error, errores centinela, envoltura de errores
- **Interfaces** — interfaces pequeñas, composición, definiciones del lado del consumidor
- **Concurrencia** — granularidad de mutex, ciclos de vida de gorutinas, paso de contexto
- **Pruebas** — pruebas basadas en tablas, ayudantes `newTestService()`/`seedTestData()`, generación de mocks
- **Patrones del proyecto** — capa de servicio, estructuras de solicitud/respuesta, opciones funcionales, patrón de controlador MCP

### swag2mcp-cli

La [habilidad swag2mcp-cli](https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-cli/SKILL.md) documenta cada comando CLI con sintaxis, banderas, argumentos y ejemplos. Útil cuando se trabaja en comandos CLI o se escribe documentación.

## Decisiones arquitectónicas clave

### Patrón de capa de servicio

Cada característica sigue el mismo patrón de tres pasos:

1. **Validar** la solicitud con `s.validateRequest(req)` (usa `go-playground/validator`)
2. **Buscar** entidades en el índice en memoria (devuelve `LLMError` con código `not_found`)
3. **Ejecutar** la lógica de negocio y devolver una respuesta tipada o `LLMError`

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

### Estructuras de solicitud/respuesta

Cada método tiene una estructura `{Method}Request` y `{Method}Response` dedicada. Las estructuras de solicitud usan etiquetas `validate` para validación y etiquetas `jsonschema` para documentación:

```go
type SearchRequest struct {
    Query string `json:"query" validate:"required,min=1" jsonschema:"description=Search query supporting field filters"`
    Limit int    `json:"limit" validate:"required,min=1,max=50" jsonschema:"description=Maximum results"`
}

type SearchResponse struct {
    Results []EndpointSearchItem `json:"results"`
}
```

### Opciones funcionales

La configuración usa el patrón de opciones funcionales:

```go
type Option func(*Service)

func New(opts ...Option) (*Service, error)

func WithDisableLLMAuth(disable bool) Option {
    return func(s *Service) {
        s.disableLLMAuth.Store(disable)
    }
}
```

### Patrón de controlador MCP

El servidor MCP usa un patrón de interfaz compuesta. La interfaz `Svc` en `internal/server/mcp/handler.go` se compone de interfaces más pequeñas (`CatalogReader`, `EndpointExplorer`, `EndpointExecutor`, `SystemInfo`, `ResponseManager`). Cada método de controlador delega en la capa de servicio:

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

Todos los errores devueltos al LLM usan el tipo `LLMError` con uno de 8 códigos:

| Código | Cuándo |
|-------|--------|
| `validation_failed` | Entrada inválida (formato de ID incorrecto, campos requeridos faltantes) |
| `not_found` | Entidad no encontrada en el índice |
| `rate_limit` | Enfriamiento de 10s por endpoint excedido |
| `invoke_error` | Fallos de solicitud/respuesta HTTP |
| `config_error` | Fallo de carga o validación de configuración |
| `workspace_error` | Fallo de operación de directorio o archivo del espacio de trabajo |
| `parse_error` | Fallo de análisis de archivo de especificación |
| `auth_error` | Fallo de recuperación de token de autenticación |

Los mensajes deben explicar qué salió mal Y qué hacer a continuación, en lenguaje sencillo adecuado para un consumidor LLM.

### Generación de IDs

Todos los IDs son hashes MD5 deterministas:

```go
id.Domain("meteo")                          // 32 caracteres hexadecimales
id.Collection("meteo", "Forecast")          // 32 caracteres hexadecimales
id.Tag("meteo", "Forecast", "pets")         // 32 caracteres hexadecimales
id.Method("meteo", "Forecast", "pets", "GET", "/v2/pet/{petId}")
```

### Cascada de configuración

La configuración se transmite en cascada a través de tres niveles: **global → especificación → colección**. Cada nivel anula el anterior. Todas las configuraciones de `http_client` pueden anularse en cada nivel. Los encabezados y cookies se fusionan; los valores simples se reemplazan.

## Referencia rápida

| Área | Convención |
|------|------------|
| **Versión de Go** | 1.23+ |
| **Formateadores** | gofmt, gofumpt, goimports, gci |
| **Longitud de línea** | 120 caracteres |
| **Linters** | 80+ en `.golangci.yml` |
| **Tipo de error** | `LLMError` con 8 códigos |
| **Framework de mocks** | `go.uber.org/mock` |
| **Ayudantes de prueba** | `newTestService()`, `seedTestData()` |
| **Formato de configuración** | YAML con cascada |
| **Despacho de autenticación** | `UnmarshalYAML` lee el campo `type` |
| **Generación de IDs** | Basado en MD5 (`id.Domain()`, `id.Collection()`, etc.) |
| **Límite de velocidad** | 10s por endpoint para `invoke` |
| **Tamaño de respuesta** | 1 MB predeterminado, guardado en archivo cuando se excede |
| **Objetivo de cobertura** | 80%+ para paquetes principales |
| **Compilación** | `make build` |
| **Lint** | `make lint` |
| **Prueba** | `go test ./...` |
| **Generar** | `go generate ./...` |
