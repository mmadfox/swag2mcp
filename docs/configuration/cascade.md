# Configuration Cascade

swag2mcp uses a three-level configuration cascade. Each level overrides the previous.

## Levels

```
Global (http_client, mcp)
    ↓ overrides
Spec (specs[].http_client)
    ↓ overrides
Collection (specs[].collections[].http_client)
```

## What Overrides What

| Parameter | Global | Spec | Collection |
|-----------|--------|------|------------|
| `http_client.timeout` | ✅ | ✅ | ❌ |
| `http_client.max_response_size` | ✅ | ✅ | ❌ |
| `http_client.user_agent` | ✅ | ✅ | ❌ |
| `http_client.follow_redirects` | ✅ | ✅ | ❌ |
| `http_client.max_redirects` | ✅ | ✅ | ❌ |
| `http_client.proxy` | ✅ | ✅ | ❌ |
| `http_client.random` | ✅ | ✅ | ❌ |
| `http_client.headers` | ✅ | ✅ | ✅ |
| `http_client.cookies` | ✅ | ✅ | ❌ |
| `base_url` | ❌ | ✅ | ✅ |
| `auth` | ❌ | ✅ | ❌ |
| `disable` | ❌ | ✅ | ✅ |

Transport settings (timeout, proxy, user-agent, redirects, response size, randomizer) can be overridden at the spec level. Collection level only supports `headers` and `cookies`.

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
      timeout: 60s  # overrides global timeout
      headers:
        "X-API-Version": "2"  # added to global headers
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        http_client:
          headers:
            "X-Custom": "value"  # added to spec + global headers
```

## Effective Settings for "Forecast" Collection

```
timeout: 60s (from spec, overrides global 30s)
max_response_size: 1048576 (from global)
headers:
  - User-Agent: swag2mcp/1.0 (from global)
  - X-API-Version: 2 (from spec)
  - X-Custom: value (from collection)
```
