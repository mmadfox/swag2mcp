# MCP Server

The MCP server is the main interaction point for LLM agents. It exposes all configured APIs as MCP tools that the LLM can call.

## Configuration

```yaml
mcp:
  transport: stdio
```

## Transports

Three transport types are available:

| Transport | Description | When to Use |
|-----------|-------------|-------------|
| `stdio` | Standard input/output | Local LLM clients (VS Code, Cursor, Claude Desktop) |
| `sse` | Server-Sent Events | Remote clients, HTTP-based communication |
| `streamable-http` | HTTP with streaming | Web clients, modern MCP clients |

### stdio (default)

The LLM client runs swag2mcp as a child process. Communication happens over standard input and output. No network port is needed.

```yaml
mcp:
  transport: stdio
```

```bash
swag2mcp mcp
```

### SSE

Server-Sent Events transport for HTTP-based communication. The MCP server listens on an HTTP port and the LLM client connects remotely.

```yaml
mcp:
  transport: sse
  addr: "127.0.0.1:8080"
  path: "/mcp"
```

```bash
swag2mcp mcp --transport sse --http-addr 127.0.0.1:8080
```

### Streamable HTTP

Modern HTTP transport that supports streaming responses. Similar to SSE but uses a different protocol.

```yaml
mcp:
  transport: streamable-http
  addr: "127.0.0.1:8080"
  path: "/mcp"
```

```bash
swag2mcp mcp --transport streamable-http --http-addr 0.0.0.0:8080
```

## Parameters

### transport

- **Type:** `string`
- **Default:** `"stdio"`
- **Options:** `stdio`, `sse`, `streamable-http`
- **Effect:** Determines how the MCP server communicates with the LLM client.

### addr

- **Type:** `string`
- **Default:** `":8080"`
- **Description:** Listen address for SSE and Streamable HTTP transports. Format: `host:port`.
- **Examples:** `":8080"`, `"127.0.0.1:8080"`, `"0.0.0.0:9000"`

### path

- **Type:** `string`
- **Default:** `"/mcp"`
- **Description:** URL path for the MCP endpoint. The LLM client sends requests to `http://&lt;addr&gt;&lt;path&gt;`.
- **Examples:** `"/mcp"`, `"/api/mcp"`, `"/v1/mcp"`

### auth.token

- **Type:** `string`
- **Default:** `""` (no auth)
- **Description:** Bearer token for HTTP transport authentication. When set, the LLM client must include `Authorization: Bearer &lt;token&gt;` in every request.
- **Note:** Supports `$(ENV_VAR)` resolution.


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

## HTTP Authentication

Protect the MCP HTTP endpoint with a bearer token:

```yaml
mcp:
  auth:
    token: "my-secret-token"
```

Or via CLI flag:

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

## Health Check

The MCP server provides a health check endpoint that works without MCP initialization:

```bash
curl http://127.0.0.1:8080/health
# {"status":"ok","version":"v1.2.0"}
```

## Startup Flags

CLI flags override the YAML configuration. If a flag is not set, the value from `mcp` section in YAML is used as fallback.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--transport` | string | `"stdio"` | Transport type: `stdio`, `sse`, `streamable-http` |
| `--http-addr` | string | `":8080"` | HTTP server address (for SSE and Streamable HTTP) |
| `--http-path` | string | `"/mcp"` | URL path for the MCP handler |
| `--auth-token` | string | `""` | Bearer token for HTTP transport authentication |
| `--logfile` | string | `""` | Log file path (logs to stderr if unset) |
| `--disable-llm-auth` | bool | `true` | Remove the `auth` tool from the MCP tool list |
| `--dump-dir` | string | `""` | Directory to dump HTTP requests for debugging |
| `--tags` | string | `""` | Filter specs by tags (comma-separated) |
| `--auth-type` | string | `""` | JWT auth type: `jwks`, `oidc`, `introspection` |
| `--auth-jwks-url` | string | `""` | JWKS URL for JWT auth |
| `--auth-issuer` | string | `""` | JWT issuer for token validation |
| `--auth-audience` | string | `""` | JWT audience for token validation |
| `--auth-introspection-url` | string | `""` | Token introspection URL |
| `--auth-client-id` | string | `""` | Client ID for introspection auth |
| `--auth-client-secret` | string | `""` | Client secret for introspection auth |
