# Global Settings

Global settings apply to all specs unless overridden at spec or collection level.

## HTTP Client

```yaml
http_client:
  timeout: 30s
  max_response_size: 1048576
  proxy:
    url: ""
  headers: {}
  cookies: []
  user_agent: "swag2mcp/1.0"
  follow_redirects: true
  max_redirects: 10
```

## MCP Server

```yaml
mcp:
  transport: stdio
  addr: "127.0.0.1:8080"
  path: "/mcp"
  auth:
    token: ""
```

## Parameters

### HTTP Client

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `timeout` | duration | `30s` | Request timeout |
| `max_response_size` | int | `1048576` | Max response size in bytes (1 MB) |
| `user_agent` | string | `swag2mcp-global/1.0` | User-Agent header |
| `follow_redirects` | bool | `true` | Follow HTTP redirects |
| `max_redirects` | int | `10` | Max redirects to follow |
| `proxy.url` | string | `""` | HTTP proxy URL |
| `proxy.username` | string | `""` | Proxy username |
| `proxy.password` | string | `""` | Proxy password |
| `proxy.bypass` | array | `[]` | Domains to bypass proxy |
| `headers` | map | `{}` | Headers for all requests |
| `cookies` | array | `[]` | Cookies for all requests |
| `random` | bool | `false` | Randomize browser-like headers |

### MCP

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `transport` | string | `stdio` | Transport type (stdio, sse, streamable-http) |
| `addr` | string | `127.0.0.1:8080` | HTTP server address |
| `path` | string | `/mcp` | MCP endpoint path |
| `auth.token` | string | `""` | Bearer token for HTTP auth |

## Example

```yaml
http_client:
  timeout: 60s
  max_response_size: 4194304
  proxy:
    url: "http://corporate-proxy:8080"
  headers:
    "User-Agent": "MyApp/1.0"

mcp:
  transport: sse
  addr: "0.0.0.0:8080"
  path: "/api/mcp"
  auth:
    token: "my-secret-token"
```
