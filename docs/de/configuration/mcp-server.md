# MCP Server

The MCP server is the main interaction point for LLM agents. It exposes all configured APIs as MCP tools that the LLM can call.

## Configuration

```yaml
mcp:
  transport: stdio
  addr: ":8080"
  path: "/mcp"
  auth:
    token: ""
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
- **Description:** URL path for the MCP endpoint. The LLM client sends requests to `http://<addr><path>`.
- **Examples:** `"/mcp"`, `"/api/mcp"`, `"/v1/mcp"`

### auth.token

- **Type:** `string`
- **Default:** `""` (no auth)
- **Description:** Bearer token for HTTP transport authentication. When set, the LLM client must include `Authorization: Bearer <token>` in every request.
- **Note:** Supports `$(ENV_VAR)` resolution.

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
