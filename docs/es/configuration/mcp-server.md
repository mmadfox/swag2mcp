# Servidor MCP

El servidor MCP es el punto principal de interacción para los agentes LLM. Expone todas las APIs configuradas como herramientas MCP que el LLM puede llamar.

## Configuración

```yaml
mcp:
  transport: stdio
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
- **Descripción:** Ruta URL para el endpoint MCP. El cliente LLM envía solicitudes a `http://&lt;addr&gt;&lt;path&gt;`.
- **Ejemplos:** `"/mcp"`, `"/api/mcp"`, `"/v1/mcp"`

### auth.token

- **Tipo:** `string`
- **Valor predeterminado:** `""` (sin autenticación)
- **Descripción:** Token Bearer para autenticación de transporte HTTP. Cuando se establece, el cliente LLM debe incluir `Authorization: Bearer &lt;token&gt;` en cada solicitud.
- **Nota:** Admite resolución de `$(ENV_VAR)`.

### auth.type

- **Type:** `string`
- **Default:** `""` (no JWT auth)
- **Options:** `jwks`, `oidc`, `introspection`
- **Description:** JWT authentication type for HTTP transport. When set, enables dynamic token verification using JWKS, OIDC Discovery, or token introspection.

### auth.jwks_url

- **Type:** `string`
- **Default:** `""`
- **Description:** URL of the JWKS (JSON Web Key Set) endpoint. Required when `auth.type` is `jwks` or resolved via OIDC discovery.

### auth.issuer

- **Type:** `string`
- **Default:** `""`
- **Description:** Expected JWT issuer (`iss` claim). If set, tokens with a different issuer are rejected.

### auth.audience

- **Type:** `string`
- **Default:** `""`
- **Description:** Expected JWT audience (`aud` claim). If set, tokens without this audience are rejected.

### auth.introspection_url

- **Type:** `string`
- **Default:** `""`
- **Description:** Token introspection endpoint URL. Required when `auth.type` is `introspection`.

### auth.client_id

- **Type:** `string`
- **Default:** `""`
- **Description:** Client ID for introspection auth. Required when `auth.type` is `introspection`.

### auth.client_secret

- **Type:** `string`
- **Default:** `""`
- **Description:** Client secret for introspection auth. Supports `$(ENV_VAR)` resolution.

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

### With JWT authentication (JWKS)

Protect the MCP HTTP endpoint with JWT verification via a JWKS endpoint:

```yaml
mcp:
  auth:
    type: jwks
    jwks_url: "https://auth.example.com/.well-known/jwks.json"
    issuer: "https://auth.example.com/"
    audience: "swag2mcp"
```

```bash
swag2mcp mcp --transport sse --http-addr 0.0.0.0:8080 \
  --auth-type jwks \
  --auth-jwks-url "https://auth.example.com/.well-known/jwks.json" \
  --auth-issuer "https://auth.example.com/" \
  --auth-audience "swag2mcp"
```

### With JWT authentication (OIDC Discovery)

```yaml
mcp:
  auth:
    type: oidc
    issuer: "https://auth.example.com/"
    audience: "swag2mcp"
```

### With JWT authentication (Token Introspection)

```yaml
mcp:
  auth:
    type: introspection
    introspection_url: "https://auth.example.com/introspect"
    client_id: "my-client"
    client_secret: "$(MCP_CLIENT_SECRET)"
```

```bash
swag2mcp mcp --transport sse --http-addr 0.0.0.0:8080 \
  --auth-type introspection \
  --auth-introspection-url "https://auth.example.com/introspect" \
  --auth-client-id "my-client" \
  --auth-client-secret "$(MCP_CLIENT_SECRET)"
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
| `--auth-type` | string | `""` | JWT auth type: `jwks`, `oidc`, `introspection` |
| `--auth-jwks-url` | string | `""` | JWKS URL for JWT auth |
| `--auth-issuer` | string | `""` | JWT issuer for token validation |
| `--auth-audience` | string | `""` | JWT audience for token validation |
| `--auth-introspection-url` | string | `""` | Token introspection URL |
| `--auth-client-id` | string | `""` | Client ID for introspection auth |
| `--auth-client-secret` | string | `""` | Client secret for introspection auth |
