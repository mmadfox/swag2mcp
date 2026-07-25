# MCP Server

The MCP server is the main interaction point for LLM agents.

## Transports

Three transport types:

| Transport | Description | When to Use |
|-----------|-------------|-------------|
| **stdio** | Standard input/output | Local LLM clients |
| **SSE** | Server-Sent Events | Remote clients |
| **Streamable HTTP** | HTTP with streaming | Web clients |

## stdio

```yaml
global:
  mcp:
    transport: stdio
```

Default. LLM client runs swag2mcp as a child process.

## SSE

```yaml
global:
  mcp:
    transport: sse
    http_addr: "127.0.0.1:8080"
    http_path: "/mcp"
```

```bash
swag2mcp mcp --transport sse --http-addr 127.0.0.1:8080
```

## Streamable HTTP

```yaml
global:
  mcp:
    transport: streamable-http
    http_addr: "127.0.0.1:8080"
    http_path: "/mcp"
```

## HTTP Authentication

```yaml
global:
  mcp:
    auth_token: "my-secret-token"
```

Or via flag:

```bash
swag2mcp mcp --auth-token "my-secret-token"
```

## Health Check

```bash
curl http://127.0.0.1:8080/health
# {"status":"ok"}
```

## Startup Flags

| Flag | Description |
|------|-------------|
| `--transport` | Transport type (stdio, sse, streamable-http) |
| `--http-addr` | HTTP server address |
| `--http-path` | MCP endpoint path |
| `--auth-token` | Bearer token for HTTP |
| `--logfile` | Log file path |
| `--disable-llm-auth` | Disable auth tool |
| `--dump-dir` | Dump directory |
| `--tags` | Tag filter |
