# Servidor MCP

El servidor MCP es el punto principal de interacción para los agentes LLM. Expone todas las APIs configuradas como herramientas MCP que el LLM puede llamar.

## Configuración

```yaml
mcp:
  transport: stdio
  addr: ":8080"
  path: "/mcp"
  auth:
    token: ""
```

## Transportes

Hay tres tipos de transporte disponibles:

| Transporte | Descripción | Cuándo usarlo |
|------------|-------------|---------------|
| `stdio` | Entrada/salida estándar | Clientes LLM locales (VS Code, Cursor, Claude Desktop) |
| `sse` | Eventos Enviados por el Servidor | Clientes remotos, comunicación basada en HTTP |
| `streamable-http` | HTTP con streaming | Clientes web, clientes MCP modernos |

### stdio (predeterminado)

El cliente LLM ejecuta swag2mcp como un proceso hijo. La comunicación ocurre a través de la entrada y salida estándar. No se necesita puerto de red.

```yaml
mcp:
  transport: stdio
```

```bash
swag2mcp mcp
```

### SSE

Transporte de Eventos Enviados por el Servidor para comunicación basada en HTTP. El servidor MCP escucha en un puerto HTTP y el cliente LLM se conecta de forma remota.

```yaml
mcp:
  transport: sse
  addr: "127.0.0.1:8080"
  path: "/mcp"
```

```bash
swag2mcp mcp --transport sse --http-addr 127.0.0.1:8080
```

### HTTP Streamable

Transporte HTTP moderno que admite respuestas en streaming. Similar a SSE pero usa un protocolo diferente.

```yaml
mcp:
  transport: streamable-http
  addr: "127.0.0.1:8080"
  path: "/mcp"
```

```bash
swag2mcp mcp --transport streamable-http --http-addr 0.0.0.0:8080
```

## Parámetros

### transport

- **Tipo:** `string`
- **Valor predeterminado:** `"stdio"`
- **Opciones:** `stdio`, `sse`, `streamable-http`
- **Efecto:** Determina cómo se comunica el servidor MCP con el cliente LLM.

### addr

- **Tipo:** `string`
- **Valor predeterminado:** `":8080"`
- **Descripción:** Dirección de escucha para los transportes SSE y HTTP Streamable. Formato: `host:port`.
- **Ejemplos:** `":8080"`, `"127.0.0.1:8080"`, `"0.0.0.0:9000"`

### path

- **Tipo:** `string`
- **Valor predeterminado:** `"/mcp"`
- **Descripción:** Ruta URL para el endpoint MCP. El cliente LLM envía solicitudes a `http://<addr><path>`.
- **Ejemplos:** `"/mcp"`, `"/api/mcp"`, `"/v1/mcp"`

### auth.token

- **Tipo:** `string`
- **Valor predeterminado:** `""` (sin autenticación)
- **Descripción:** Token Bearer para autenticación de transporte HTTP. Cuando se establece, el cliente LLM debe incluir `Authorization: Bearer <token>` en cada solicitud.
- **Nota:** Admite resolución de `$(ENV_VAR)`.

## Autenticación HTTP

Proteger el endpoint HTTP del MCP con un token bearer:

```yaml
mcp:
  auth:
    token: "my-secret-token"
```

O mediante bandera CLI:

```bash
swag2mcp mcp --auth-token "my-secret-token"
```

## Verificación de Salud

El servidor MCP proporciona un endpoint de verificación de salud que funciona sin inicialización MCP:

```bash
curl http://127.0.0.1:8080/health
# {"status":"ok","version":"v1.2.0"}
```

## Banderas de Inicio

Las banderas CLI anulan la configuración YAML. Si una bandera no se establece, el valor de la sección `mcp` en YAML se usa como respaldo.

| Bandera | Tipo | Valor predeterminado | Descripción |
|---------|------|---------------------|-------------|
| `--transport` | string | `"stdio"` | Tipo de transporte: `stdio`, `sse`, `streamable-http` |
| `--http-addr` | string | `":8080"` | Dirección del servidor HTTP (para SSE y HTTP Streamable) |
| `--http-path` | string | `"/mcp"` | Ruta URL para el controlador MCP |
| `--auth-token` | string | `""` | Token Bearer para autenticación de transporte HTTP |
| `--logfile` | string | `""` | Ruta del archivo de registro (registra en stderr si no se establece) |
| `--disable-llm-auth` | bool | `true` | Eliminar la herramienta `auth` de la lista de herramientas MCP |
| `--dump-dir` | string | `""` | Directorio para volcar solicitudes HTTP para depuración |
| `--tags` | string | `""` | Filtrar especificaciones por etiquetas (separadas por comas) |
