# Configuration Cascade

swag2mcp uses a three-level configuration cascade. Each level overrides the previous.

## Levels

```
Global (http_client, mcp)
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
| `http_client.headers` | ✅ | ✅ | ✅ |
| `http_client.cookies` | ✅ | ✅ | ❌ |
| `http_client.proxy` | ✅ | ✅ | ❌ |
| `http_client.user_agent` | ✅ | ✅ | ❌ |
| `http_client.follow_redirects` | ✅ | ✅ | ❌ |
| `http_client.max_redirects` | ✅ | ✅ | ❌ |
| `http_client.random` | ✅ | ✅ | ❌ |
| `base_url` | ❌ | ✅ | ✅ |
| `auth` | ❌ | ✅ | ❌ |
| `disable` | ❌ | ✅ | ✅ |

## Cascade Example

```yaml
http_client:
  timeout: 30s
  max_response_size: 1048576
  headers:
    "User-Agent": "swag2mcp/1.0"

specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    http_client:
      timeout: 10s  # overrides 30s
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        http_client:
          timeout: 5s  # overrides 10s
          headers:
            "X-Custom": "value"  # added to User-Agent
```

## Effective Settings for "Forecast" Collection

```
timeout: 5s
max_response_size: 1048576 (from global)
headers:
  - User-Agent: swag2mcp/1.0 (from global)
  - X-Custom: value (from collection)
```
