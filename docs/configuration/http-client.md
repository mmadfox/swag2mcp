# HTTP Client

swag2mcp uses a configurable HTTP client for API calls.

## Configuration

```yaml
global:
  http_client:
    timeout: 30s
    max_response_size: 2048
    proxy: ""
    insecure_skip_verify: false
    headers:
      "User-Agent": "swag2mcp/1.0"
    cookies:
      - name: "session"
        value: "abc123"
```

## Parameters

| Parameter | Description |
|-----------|-------------|
| `timeout` | HTTP request timeout |
| `max_response_size` | Max response size in bytes |
| `proxy` | HTTP proxy URL |
| `insecure_skip_verify` | Disable TLS certificate verification |
| `headers` | Headers added to every request |
| `cookies` | Cookies added to every request |

## Randomizer

Add random browser-like headers:

```yaml
global:
  http_client:
    randomize_headers: true
```

Adds random `User-Agent`, `Accept`, `Accept-Language` headers.

## Proxy

```yaml
global:
  http_client:
    proxy: "http://user:pass@proxy.example.com:8080"
```

## Custom Headers

```yaml
specs:
  - domain: "api.example.com"
    headers:
      "X-API-Key": "{{API_KEY}}"
      "X-Correlation-ID": "swag2mcp"
```

## Cookies

```yaml
specs:
  - domain: "api.example.com"
    cookies:
      - name: "session"
        value: "abc123"
      - name: "csrf"
        value: "{{CSRF_TOKEN}}"
```

## Cascade

HTTP client settings cascade:

1. Global (`global.http_client`)
2. Spec (`specs[].http_client`)
3. Collection (`specs[].collections[].http_client`)

Each level overrides the previous.
