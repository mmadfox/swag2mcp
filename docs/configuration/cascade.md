# Configuration Cascade

swag2mcp uses a three-level configuration cascade. Each level overrides the previous.

## Levels

```
Global (global)
    ↓ overrides
Spec (specs[].*)
    ↓ overrides
Collection (specs[].collections[].*)
```

## What Overrides What

| Parameter | Global | Spec | Collection |
|-----------|--------|------|------------|
| `http_client.timeout` | ✅ | ✅ | ✅ |
| `http_client.max_response_size` | ✅ | ✅ | ✅ |
| `http_client.proxy` | ✅ | ✅ | ❌ |
| `http_client.headers` | ✅ | ✅ | ✅ |
| `http_client.cookies` | ✅ | ✅ | ❌ |
| `headers` | ❌ | ✅ | ✅ |
| `cookies` | ❌ | ✅ | ❌ |
| `auth` | ❌ | ✅ | ❌ |
| `disabled` | ❌ | ✅ | ❌ |

## Cascade Example

```yaml
global:
  http_client:
    timeout: 30s
    max_response_size: 2048
    headers:
      "User-Agent": "swag2mcp/1.0"

specs:
  - domain: "api.example.com"
    location: "https://api.example.com/openapi.json"
    http_client:
      timeout: 10s  # overrides 30s
    headers:
      "X-API-Key": "key123"  # added to User-Agent
    collections:
      - name: "admin"
        tags: ["admin"]
        http_client:
          timeout: 5s  # overrides 10s
        headers:
          "X-Role": "admin"  # added to X-API-Key and User-Agent
```

## Effective Settings for "admin" Collection

```
timeout: 5s
max_response_size: 2048 (from global)
headers:
  - User-Agent: swag2mcp/1.0 (from global)
  - X-API-Key: key123 (from spec)
  - X-Role: admin (from collection)
```
