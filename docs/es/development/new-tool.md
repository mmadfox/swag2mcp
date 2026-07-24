# Agregar una Nueva Herramienta MCP

## Pasos

1. **Agregar una constante de nombre de herramienta** en `internal/service/service.go`
2. **Crear tipos de solicitud/respuesta** en `internal/service/types.go`
3. **Implementar el servicio** en `internal/service/` (nuevo archivo o agregar a uno existente)
4. **Crear una definición markdown** en `internal/service/definitions/` — esto es lo que lee `MakeToolDefinitions`
5. **Agregar método a la interfaz `Svc`** en `internal/server/mcp/handler.go`
6. **Agregar controlador** en `handler.go`
7. **Registrar herramienta** en `registerTools` en `mcp.go`
8. **Generar mocks**: `go generate ./...`
9. **Escribir pruebas**

## 1. Constante de nombre de herramienta

Agregue una constante en `internal/service/service.go`:

```go
const MyNewTool = "my_new_tool"
```

## 2. Tipos de solicitud/respuesta

Defina en `internal/service/types.go`:

```go
type MyNewToolRequest struct {
    Param1 string `json:"param1" validate:"required" jsonschema:"required,Description of param1"`
}

type MyNewToolResponse struct {
    Result string `json:"result"`
}
```

## 3. Implementación del servicio

Cree `internal/service/my_new_tool.go` o agregue a un archivo de servicio existente. Siga el patrón de servicio estándar: validar → buscar → ejecutar → devolver:

```go
func (s *Service) MyNewTool(ctx context.Context, req MyNewToolRequest) (MyNewToolResponse, error) {
    if err := s.validateRequest(req); err != nil {
        return MyNewToolResponse{}, NewLLMError(validationFailedErrCode, err.Error())
    }
    // lógica de negocio
    return MyNewToolResponse{Result: "ok"}, nil
}
```

## 4. Definición markdown

Cree `internal/service/definitions/my_new_tool.md`. Este archivo es leído por `MakeToolDefinitions()` y se incrusta en el binario. El campo `name:` del frontmatter debe coincidir con la constante:

```markdown
---
name: my_new_tool
---

# my_new_tool

Descripción de la herramienta.

## Parámetros

| Parámetro | Tipo | Descripción |
|-----------|------|-------------|
| `param1` | string | Descripción |
```

La función `MakeToolDefinitions()` en `tools.go` lee todos los archivos `.md` del directorio `definitions/` incrustado, analiza el frontmatter YAML para el campo `name` y usa el cuerpo como la descripción de la herramienta. El archivo `instruction.md` se trata de forma especial — se convierte en la instrucción del sistema para el LLM.

## 5. Interfaz Svc

Agregue un método a la interfaz compuesta `Svc` en `handler.go`:

```go
type Svc interface {
    // ... métodos existentes
    MyNewTool(ctx context.Context, req service.MyNewToolRequest) (service.MyNewToolResponse, error)
}
```

## 6. Controlador

Agregue un método de controlador en `handler.go`. El controlador delega en el servicio y envuelve el resultado en `StructuredContent`:

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

## 7. Registro

Registre la herramienta en la función `registerTools` en `mcp.go`. Agregue una entrada al mapa `toolRegistrations`:

```go
service.MyNewTool: {
    addTool[service.MyNewToolRequest](mcpServer, h.handleMyNewTool),
    true, // false si la herramienta es mutable (como invoke o auth)
},
```

La firma de la función `registerTools` es:

```go
func registerTools(mcpServer *sdkmcp.Server, tools []service.Tool, h handler) {
```

Itera sobre las definiciones de herramientas devueltas por `MakeToolDefinitions()` y registra cada una con su controlador tipado. El mapa `toolRegistrations` conecta las constantes de nombre de herramienta con sus controladores.
