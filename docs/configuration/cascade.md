# Configuration Cascade

swag2mcp uses a three-level configuration cascade. Each level overrides the previous. This lets you set sensible defaults globally and fine-tune settings for specific specs or collections.

## Levels

```
Global (http_client, mcp, mock_enabled, disable_ratelimiter, rate_limit_interval)
    ↓ overrides
Spec (specs[].http_client, specs[].auth, specs[].base_url, specs[].disable, specs[].tags)
    ↓ overrides
Collection (specs[].collections[].http_client, specs[].collections[].base_url, specs[].collections[].disable)
```

## What Overrides What

| Parameter | Global | Spec | Collection |
|-----------|--------|------|------------|
| `http_client.timeout` | ✅ | ✅ | ✅ |
| `http_client.max_response_size` | ✅ | ✅ | ✅ |
| `http_client.user_agent` | ✅ | ✅ | ✅ |
| `http_client.follow_redirects` | ✅ | ✅ | ✅ |
| `http_client.max_redirects` | ✅ | ✅ | ✅ |
| `http_client.proxy` | ✅ | ✅ | ✅ |
| `http_client.random` | ✅ | ✅ | ✅ |
| `http_client.headers` | ✅ | ✅ | ✅ |
| `http_client.cookies` | ✅ | ✅ | ✅ |
| `base_url` | ❌ | ✅ | ✅ |
| `auth` | ❌ | ✅ | ❌ |
| `disable` | ❌ | ✅ | ✅ |
| `tags` | ❌ | ✅ | ❌ |
| `mock_enabled` | ✅ | ❌ | ❌ |
| `disable_ratelimiter` | ✅ | ❌ | ❌ |
| `rate_limit_interval` | ✅ | ❌ | ❌ |

All `http_client` settings can be overridden at every level. Collection-level settings take full precedence over spec and global.

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
          timeout: 120s  # overrides spec timeout
          headers:
            "X-Custom": "value"  # added to spec + global headers
```

## Effective Settings for "Forecast" Collection

```
timeout: 120s (from collection, overrides spec 60s and global 30s)
max_response_size: 1048576 (from global)
headers:
  - User-Agent: swag2mcp/1.0 (from global)
  - X-API-Version: 2 (from spec)
  - X-Custom: value (from collection)
```

## How Merging Works

### HTTP Client Settings

Simple values (`timeout`, `max_response_size`, `user_agent`, `follow_redirects`, `max_redirects`, `random`) are **replaced** at each level. If a spec sets `timeout: 60s`, it completely replaces the global `30s`.

### Headers

Headers are **merged** across levels. All three levels' headers are combined. If the same header key appears at multiple levels, the lowest level wins.

### Cookies

Cookies are **merged** across levels. If the same cookie name appears at multiple levels, the lowest level wins.

### Proxy

Proxy is **replaced** at each level. If a spec sets a proxy, it completely replaces the global proxy for that spec.
