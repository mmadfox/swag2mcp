# Global Settings

Global settings apply to all specs unless overridden at spec or collection level.

## global Section

```yaml
global:
  http_client:
    timeout: 30s
    max_response_size: 2048
    proxy: ""
    headers: {}
    cookies: []
    insecure_skip_verify: false
  mcp:
    transport: stdio
    http_addr: "127.0.0.1:8080"
    http_path: "/mcp"
    auth_token: ""
```

## HTTP Client Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `timeout` | duration | `30s` | Request timeout |
| `max_response_size` | int | `2048` | Max response size in bytes |
| `proxy` | string | `""` | HTTP proxy URL |
| `insecure_skip_verify` | bool | `false` | Disable TLS verification |
| `headers` | map | `{}` | Headers for all requests |
| `cookies` | array | `[]` | Cookies for all requests |

## MCP Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `transport` | string | `stdio` | Transport type |
| `http_addr` | string | `127.0.0.1:8080` | HTTP server address |
| `http_path` | string | `/mcp` | MCP endpoint path |
| `auth_token` | string | `""` | Bearer token for HTTP |

## Example

```yaml
global:
  http_client:
    timeout: 60s
    max_response_size: 4096
    proxy: "http://corporate-proxy:8080"
    headers:
      "User-Agent": "MyApp/1.0"
  mcp:
    transport: sse
    http_addr: "0.0.0.0:8080"
    http_path: "/api/mcp"
    auth_token: "my-secret-token"
```
