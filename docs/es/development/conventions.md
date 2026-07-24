# Convenciones de Código

## Go

- **Go 1.26+**
- **gofmt** / **gofumpt** / **goimports** / **gci**
- **120 caracteres** por línea
- **Cláusulas de guarda** en lugar de ifs anidados
- **Nomenclatura**: `camelCase` para privado, `PascalCase` para exportado

## Errores

Use `LLMError` para errores visibles por el LLM:

```go
type LLMError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}
```

Códigos de error:
- `validation_failed` — parámetros inválidos
- `not_found` — recurso no encontrado
- `rate_limit` — límite de velocidad excedido
- `invoke_error` — error de llamada a la API

## Interfaces

- Interfaces pequeñas (1-3 métodos)
- Composición de interfaces
- Opciones funcionales para configuración

## Pruebas

- Pruebas basadas en tablas
- Ayudantes de prueba (`newTestService()`, `seedTestData()`)
- Mocks mediante `go.uber.org/mock`
- 80%+ de cobertura para paquetes principales

## Configuración

- Formato YAML
- Cascada: global → especificación → colección
- Validación mediante `go-playground/validator`
- Variables de entorno mediante `$(VAR)`
