# Spec Settings

Spec settings override global settings for a specific API.

## Spec Section

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    llm_instruction: "Use this API for weather forecasts and climate data"
    base_url: https://api.open-meteo.com
    disable: false
    tags: ["weather", "climate"]
    http_client:
      timeout: 10s
      max_response_size: 1024
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
    auth:
      type: bearer
      config:
        token: "my-token"
```

## Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `domain` | string | — | Unique API identifier (a-z, 0-9, _, -, max 60 chars) |
| `llm_title` | string | — | Human-readable title (5-120 chars) |
| `llm_instruction` | string | `""` | Instructions for the LLM (max 500 chars) |
| `base_url` | string | — | Base URL for API requests |
| `disable` | bool | `false` | Disable this spec |
| `tags` | array | `[]` | Tags for filtering |
| `http_client` | object | — | HTTP client override (all settings: timeout, proxy, headers, cookies, redirects, user-agent, response size, randomizer) |
| `collections` | array | — | Collection list (1-30 items) |
| `auth` | object | — | Auth settings |

## Disabling a Spec

```yaml
specs:
  - domain: old-api
    llm_title: Old API
    base_url: https://old-api.example.com
    disable: true
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

Disabled specs are not loaded or indexed.

## HTTP Client Override

```yaml
specs:
  - domain: slow-api
    llm_title: Slow API
    base_url: https://slow-api.example.com
    http_client:
      timeout: 120s
      max_response_size: 8388608
      headers:
        "X-DC": "us-east-1"
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

All `http_client` settings from the global level can be overridden at the spec level: `timeout`, `proxy`, `user_agent`, `follow_redirects`, `max_redirects`, `max_response_size`, `random`, `headers`, and `cookies`.

## Proxy Override

```yaml
specs:
  - domain: proxied-api
    llm_title: Proxied API
    base_url: https://api.example.com
    http_client:
      proxy:
        url: http://proxy.company.com:8080
        username: $(PROXY_USER)
        password: $(PROXY_PASS)
        bypass:
          - "*.local"
          - "10.0.0.0/8"
    collections:
      - llm_title: Main
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo.json
```
