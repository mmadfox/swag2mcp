# info

## Propósito

Mostrar un resumen completo del entorno de ejecución de swag2mcp como **JSON**. Esto incluye la versión, la ruta del espacio de trabajo, el resumen de especificaciones, la configuración del cliente HTTP, la configuración del transporte MCP, los métodos de autenticación y el estado del modo simulado.

## Cuándo usarlo

- Desea una visión general legible por máquina del espacio de trabajo
- Necesita verificar la configuración de ejecución para depuración
- Desea ver cuántas especificaciones y endpoints están activos
- Necesita verificar la configuración del cliente HTTP o del transporte MCP

## Sintaxis

```bash
swag2mcp info [path]
```

## Argumentos

| Argumento | Posición | Requerido | Descripción |
|-----------|----------|-----------|-------------|
| `path` | 1 | No | Directorio del espacio de trabajo. Si se omite, se resuelve mediante reglas de resolución de ruta. |

## Banderas

Ninguna.

## Cómo funciona

```bash
swag2mcp info
swag2mcp info ./my-workspace
```

## Salida

La salida es un objeto JSON con la siguiente estructura:

```json
{
  "version": "v1.2.0",
  "workspace": "~/.swag2mcp",
  "uptime": "2h 15m",
  "specs": {
    "total": 4,
    "active": 3,
    "disabled": 1,
    "collections": 6,
    "endpoints": 42
  },
  "http_client": {
    "timeout": "30s",
    "max_response_size": "2 KB",
    "proxy": "none",
    "follow_redirects": true,
    "max_redirects": 10,
    "randomize": false
  },
  "mcp": {
    "transport": "stdio",
    "addr": ":8080",
    "path": "/mcp"
  },
  "auth_methods": ["bearer", "api-key"],
  "mock_enabled": false
}
```

## Verificación posterior al comando

Use `info` para confirmar que el espacio de trabajo se cargó correctamente y que todas las especificaciones están activas antes de iniciar el servidor MCP.

## Matices

- **Auto-inicio:** Si no existe un archivo de configuración, `info` ejecuta automáticamente el asistente de inicio primero.
- **Solo JSON:** La salida siempre es JSON. Para una salida legible por humanos, use `ls`.
- **`max_response_size`:** Se muestra en formato legible por humanos (por ejemplo, `"1 KB"`, `"2 MB"`).
- **Sin índice de texto completo:** `info` deshabilita la indexación de texto completo ya que solo necesita metadatos de configuración y especificaciones.
