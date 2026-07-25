# Examples

## Minimal Configuration

```yaml
specs:
  - domain: "api.example.com"
    location: "https://api.example.com/openapi.json"
    collections:
      - name: "default"
        tags: ["*"]
```

Full example in `examples/minimal-config/`.

## Full Configuration

```yaml
global:
  http_client:
    timeout: 30s
    max_response_size: 2048
    proxy: ""
    headers:
      "User-Agent": "swag2mcp/1.0"
  mcp:
    transport: stdio
    http_addr: "127.0.0.1:8080"
    http_path: "/mcp"

specs:
  - domain: "petstore.swagger.io"
    location: "https://petstore.swagger.io/v2/swagger.json"
    disabled: false
    headers:
      "X-API-Key": "{{API_KEY}}"
    collections:
      - name: "pets"
        tags: ["pet"]
      - name: "store"
        tags: ["store"]
      - name: "users"
        tags: ["user"]
    auth:
      type: bearer
      bearer:
        token: "{{TOKEN}}"
```

Full example in `examples/full-config/`.

## Authentication

All 9 auth methods with examples in `examples/auth/`:

- `no-auth/` — no auth
- `basic-auth/` — HTTP Basic
- `bearer-auth/` — Bearer Token
- `api-key-header/` — API Key in header
- `api-key-query/` — API Key in query
- `digest-auth/` — HTTP Digest
- `hmac-auth/` — HMAC-SHA256
- `oauth2-client-credentials/` — OAuth2 CC
- `oauth2-password/` — OAuth2 Password
- `script-auth/` — external script

## Transport

Examples in `examples/mcp-transport/`:

- `stdio/` — standard input/output
- `sse/` — Server-Sent Events
- `streamable-http/` — Streamable HTTP

## Mock Server

Example in `examples/mock-server/`:

```yaml
specs:
  - domain: "api.example.com"
    location: "https://api.example.com/openapi.json"
    mock:
      enabled: true
      delay: 100ms
      error_rate: 0.05
```
