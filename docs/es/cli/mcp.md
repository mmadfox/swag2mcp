# mcp

## Propósito

Iniciar el **servidor MCP (Protocolo de Contexto de Modelo)** — el modo principal para la integración con LLM. Esto es lo que ejecuta para dar a un agente de IA (Claude, Cursor, OpenCode, etc.) acceso a sus APIs a través de 16 herramientas MCP.

## Cuándo usarlo

- Desea conectar un agente LLM a sus APIs
- Está configurando un IDE (VS Code, Cursor, JetBrains) o aplicación de escritorio (Claude Desktop)
- Necesita exponer sus APIs a través del protocolo MCP
- Está probando el servidor MCP antes de la integración

## Sintaxis

```bash
swag2mcp mcp [path] [flags]
```

## Argumentos

| Argumento | Posición | Requerido | Descripción |
|-----------|----------|-----------|-------------|
| `path` | 1 | No | Directorio del espacio de trabajo. Si se omite, se resuelve mediante reglas de resolución de ruta. |

## Banderas

| Bandera | Abreviatura | Tipo | Valor predeterminado | Descripción |
|---------|-------------|------|---------------------|-------------|
| `--transport` | | `string` | `"stdio"` | Transporte MCP: `stdio`, `sse`, `streamable-http` |
| `--http-addr` | | `string` | `":8080"` | Dirección del servidor HTTP (para `sse` y `streamable-http`) |
| `--http-path` | | `string` | `"/mcp"` | Ruta HTTP para el controlador MCP |
| `--auth-token` | | `string` | `""` | Token Bearer para autenticación de transporte HTTP |
| `--logfile` | `-f` | `string` | `""` | Ruta del archivo de registro. Si no se establece, registra en stderr. |
| `--disable-llm-auth` | | `bool` | `true` | Eliminar la herramienta `auth` de la lista de herramientas MCP |
| `--dump-dir` | | `string` | `""` | Directorio para volcar solicitudes HTTP para depuración |
| `--tags` | `-t` | `string` | `""` | Filtrar especificaciones por etiquetas (separadas por comas) |

## Cómo funciona

### Transporte stdio (predeterminado)

Se usa cuando el servidor MCP se inicia como un subproceso por el cliente LLM (IDE, Claude Desktop, etc.). El servidor se comunica a través de la entrada/salida estándar.

```bash
swag2mcp mcp
```

### Transporte SSE

Transporte de Eventos Enviados por el Servidor para comunicación basada en HTTP. Requiere la secuencia de protocolo de enlace MCP.

```bash
swag2mcp mcp --transport sse --http-addr :8080
```

### Transporte HTTP Streamable

Transporte HTTP moderno que admite respuestas en streaming.

```bash
swag2mcp mcp --transport streamable-http --http-addr 0.0.0.0:8080
```

### Con autenticación

Proteger el endpoint HTTP con un token bearer:

```bash
swag2mcp mcp --transport sse --http-addr :8080 --auth-token "my-secret"
```

### Con filtrado por etiquetas

Cargar solo especificaciones con etiquetas específicas:

```bash
swag2mcp mcp --tags=public
```

### Con herramienta auth habilitada (modo depuración)

Permitir que el LLM solicite tokens frescos a través de la herramienta `auth`:

```bash
swag2mcp mcp --disable-llm-auth=false
```

### Con directorio de volcado de solicitudes

Guardar todas las solicitudes HTTP para depuración:

```bash
swag2mcp mcp --dump-dir ./dumps
```

## Transporte HTTP MCP — Protocolo de Enlace

Al usar `sse` o `streamable-http`, el protocolo MCP requiere un protocolo de enlace específico. Las llamadas a herramientas fallarán antes de la inicialización:

```
Paso 1: POST /mcp → {"method":"initialize", ...}
Paso 2: POST /mcp → {"method":"notifications/initialized"}
Paso 3: POST /mcp → {"method":"tools/list", ...}   ← ahora funciona
```

### Verificación de salud

Funciona sin inicialización:

```bash
curl http://localhost:8080/health
# → {"status":"ok","version":"v1.2.0"}
```

## Ejemplos de Configuración del IDE

### VS Code (`.vscode/settings.json` o configuración global)

```json
{
  "mcp": {
    "servers": {
      "swag2mcp": {
        "command": "swag2mcp",
        "args": ["mcp", "/absolute/path/to/.swag2mcp"]
      }
    }
  }
}
```

### Cursor / Windsurf (`~/.cursor/mcp.json` o proyecto `.cursor/mcp.json`)

```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "/absolute/path/to/.swag2mcp"]
    }
  }
}
```

### Claude Desktop (`~/Library/Application Support/Claude/claude_desktop_config.json` en macOS)

```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "/absolute/path/to/.swag2mcp"]
    }
  }
}
```

### IDEs JetBrains (Configuración → Herramientas → MCP)

- Nombre: `swag2mcp`
- Comando: `swag2mcp`
- Argumentos: `mcp /absolute/path/to/.swag2mcp`

> **Use siempre una ruta absoluta** al directorio del espacio de trabajo en la configuración del IDE. Las rutas relativas pueden fallar dependiendo del directorio de trabajo del IDE.

## Salida

En caso de éxito, el servidor imprime:

```
MCP server listening on http://127.0.0.1:8080/mcp
```

## Matices

- **Sin auto-inicio:** Si el archivo de configuración no existe, `mcp` devuelve un error: `"configuration not found at &lt;path&gt;"`. Ejecute `init` primero.
- **`--disable-llm-auth` (predeterminado: `true`):** Cuando está habilitado, la herramienta `auth` se elimina por completo de la lista de herramientas MCP. El LLM no puede ver ni solicitar tokens. La autenticación sigue funcionando — los tokens se obtienen a través del mecanismo de configuración estándar, no a través del LLM. Este modo se recomienda para **producción**. Para **depuración** o cuando se usan tokens de corta duración, establezca `--disable-llm-auth=false` para permitir que el LLM solicite tokens frescos a través de la herramienta `auth`.
- **Respaldo de configuración YAML:** Si una bandera CLI no se establece explícitamente, el valor se toma de la sección `mcp` en `swag2mcp.yaml` (si está presente). Esto le permite configurar el servidor en el archivo de configuración en lugar de pasar banderas cada vez.
- **Limpieza de respuestas:** Al iniciar, las respuestas con más de 48 horas se eliminan automáticamente del directorio `responses/`.
- **Advertencia de resolución de ruta:** Cuando se omite `[path]`, `mcp` busca `swag2mcp.yaml` en el directorio actual primero, luego recurre a `~/.swag2mcp/`. Si ejecuta el comando desde el directorio incorrecto, puede cargar un espacio de trabajo diferente al previsto. **Siempre especifique `[path]` explícitamente cuando ejecute como servicio o en la configuración del IDE.**
