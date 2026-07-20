# Collection Settings

Collection settings override spec settings for a specific endpoint group.

## Collection Section

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        llm_instruction: "Use for current and forecast weather data"
        http_client:
          timeout: 5s
```

## Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `llm_title` | string | — | Collection display name (max 120 chars) |
| `llm_instruction` | string | `""` | Instructions for the LLM (max 360 chars) |
| `title` | string | `""` | Original spec title override |
| `location` | string | — | URL or path to spec file (5-250 chars) |
| `disable` | bool | `false` | Disable this collection |
| `http_client` | object | — | HTTP client override |
| `base_url` | string | `""` | Override base URL for this collection |
| `base_mock_url` | string | `""` | Mock server address (host:port) |

## Multiple Collections from One Spec

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
      - llm_title: Marine
        base_url: https://marine-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
```

## Disabling a Collection

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
        disable: true
```

## HTTP Client Override

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        http_client:
          timeout: 5s
          headers:
            "X-Custom": "value"
```
