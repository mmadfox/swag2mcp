# Spec Settings

Spec settings define an API service and override global settings for that specific API. Each spec represents one logical API (e.g., "Open-Meteo Weather APIs") and can contain multiple collections (spec files).

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
    auth:
      type: bearer
      config:
        token: "$(TOKEN)"
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

## Parameters

### domain

- **Type:** `string`
- **Required:** Yes
- **Description:** Unique identifier for this API spec. Used internally to reference the spec.
- **Rules:** 1-60 characters. Only lowercase letters (`a-z`), digits (`0-9`), hyphens (`-`), and underscores (`_`).
- **Example:** `meteo`, `binance`, `my-api`

### llm_title

- **Type:** `string`
- **Required:** Yes
- **Description:** Human-readable name that the LLM uses to reference this API. Shown in MCP tool responses.
- **Rules:** 5-120 characters. Letters, digits, spaces, and basic punctuation only.
- **Example:** `Open-Meteo Weather APIs`, `Binance Market Data`

### llm_instruction

- **Type:** `string`
- **Default:** `""`
- **Description:** Instructions for the LLM on how to use this API. Describes what the API does and when to use it.
- **Rules:** Max 500 characters. Letters, digits, spaces, and basic punctuation only.
- **Example:** `"Use this API for weather forecasts, current conditions, and climate data."`

### base_url

- **Type:** `string`
- **Required:** Yes
- **Description:** Base URL for all API requests in this spec. The endpoint paths from the OpenAPI spec are appended to this URL.
- **Example:** `https://api.open-meteo.com`, `https://api.binance.com`
- **Note:** Can be overridden at the collection level if different collections use different base URLs.

### disable

- **Type:** `bool`
- **Default:** `false`
- **Description:** When `true`, this spec is excluded from MCP tools. It is not loaded, indexed, or available to the LLM.
- **When to use:** Temporarily disable an API without removing it from the config. Useful for APIs that are down, deprecated, or under maintenance.

### tags

- **Type:** `[]string` (array of strings)
- **Default:** `[]`
- **Description:** Tags for filtering specs. Used with the `--tags` flag in CLI commands (`ls`, `validate`, `mcp`, `update`).
- **Example:** `["public", "weather"]`, `["internal", "production"]`
- **Effect:** When you run `swag2mcp mcp --tags=public`, only specs with the `public` tag are loaded.

### http_client

- **Type:** `object`
- **Default:** inherits from global
- **Description:** Override global HTTP client settings for this spec. All settings from the global `http_client` can be overridden: `timeout`, `max_response_size`, `user_agent`, `follow_redirects`, `max_redirects`, `random`, `proxy`, `headers`, `cookies`.
- **Example:**
  ```yaml
  http_client:
    timeout: 60s
    max_response_size: 4194304
    headers:
      "X-DC": "us-east-1"
  ```

### auth

- **Type:** `object`
- **Default:** `none` (no authentication)
- **Description:** Authentication configuration for this spec. See the [Authentication](/auth/overview) section for all 9 methods and their parameters.
- **Example:**
  ```yaml
  auth:
    type: bearer
    config:
      token: "$(API_TOKEN)"
  ```

### collections

- **Type:** `[]object` (array of collections)
- **Required:** Yes (at least 1)
- **Description:** List of OpenAPI/Swagger/Postman spec files that belong to this spec. Each collection is one spec file.
- **Rules:** 1-30 collections per spec.
- **See:** [Collection Settings](./collection-settings) for all collection parameters.

## Disabling a Spec

Disabled specs are not loaded or indexed. The LLM cannot see or use them.

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

## HTTP Client Override

All `http_client` settings from the global level can be overridden at the spec level. The spec values take precedence over global values for this spec only.

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

## Proxy Override

If this spec requires a different proxy than the global one, configure it at the spec level:

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
