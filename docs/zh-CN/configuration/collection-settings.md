# Collection Settings

Collection settings define a single OpenAPI/Swagger/Postman spec file and override spec settings for that specific file. Each collection belongs to a spec and represents one API specification document.

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
        disable: false
        base_url: https://forecast-api.open-meteo.com
        base_mock_url: localhost:8081
        http_client:
          timeout: 5s
```

## Parameters

### llm_title

- **Type:** `string`
- **Required:** No
- **Description:** Human-readable name for this collection. Shown in MCP tool responses.
- **Rules:** Max 120 characters. Letters, digits, spaces, and basic punctuation only.
- **Example:** `Forecast`, `Air Quality`, `Market Data`

### llm_instruction

- **Type:** `string`
- **Default:** `""`
- **Description:** Instructions for the LLM about this specific collection. Describes what endpoints this collection provides.
- **Rules:** Max 360 characters. Letters, digits, spaces, and basic punctuation only.
- **Example:** `"Use for current and forecast weather data."`

### title

- **Type:** `string`
- **Default:** `""`
- **Description:** Raw title from the spec file. Populated automatically at runtime. You typically don't need to set this in YAML.

### location

- **Type:** `string`
- **Required:** Yes
- **Description:** URL or local file path to the OpenAPI 3.x, Swagger 2.0, or Postman collection spec file.
- **Rules:** 5-250 characters.
- **Examples:**
  - URL: `https://raw.githubusercontent.com/org/repo/main/spec.yaml`
  - Local: `./specs/my-api.json`
  - Local (absolute): `/home/user/.swag2mcp/specs/my-api.yaml`

### disable

- **Type:** `bool`
- **Default:** `false`
- **Description:** When `true`, this collection is excluded from MCP tools. It is not loaded or indexed.
- **When to use:** Temporarily disable a collection without removing it from the config. Useful when a spec file is being updated or an API version is deprecated.

### http_client

- **Type:** `object`
- **Default:** inherits from spec (or global)
- **Description:** Override HTTP client settings for this collection. All settings from the global `http_client` can be overridden: `timeout`, `max_response_size`, `user_agent`, `follow_redirects`, `max_redirects`, `random`, `proxy`, `headers`, `cookies`.
- **Example:**
  ```yaml
  http_client:
    timeout: 120s
    headers:
      "X-Custom": "value"
    cookies:
      - name: "session"
        value: "abc123"
  ```

### base_url

- **Type:** `string`
- **Default:** `""` (inherits from spec)
- **Description:** Override the spec-level `base_url` for this collection. Use when different collections within the same spec use different base URLs.
- **Example:** If the spec has `base_url: https://api.open-meteo.com` but one collection uses `https://air-quality-api.open-meteo.com`, set `base_url` at the collection level.

### base_mock_url

- **Type:** `string`
- **Default:** `""`
- **Description:** Mock server address in `host:port` format. Required when `mock_enabled: true` in the global config.
- **Rules:** Host must be `localhost`, `127.0.0.1`, or `0.0.0.0`. Port must be a valid port number.
- **Example:** `localhost:8081`, `127.0.0.1:9000`
- **When to use:** You have `mock_enabled: true` and want to test this collection with fake responses.

## Multiple Collections from One Spec

A spec can have multiple collections — for example, when an API has separate spec files for different services:

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

All `http_client` settings can be overridden at the collection level. Collection values take precedence over spec and global values for this collection only.

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        http_client:
          timeout: 120s
          headers:
            "X-Custom": "value"
          cookies:
            - name: "session"
              value: "abc123"
```
