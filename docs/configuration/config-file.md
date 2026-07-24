# Configuration File

swag2mcp uses a YAML configuration file. Created by `swag2mcp init`.

## Location

- **Linux/macOS**: `~/.swag2mcp/swag2mcp.yaml`
- **Windows**: `%USERPROFILE%\.swag2mcp\swag2mcp.yaml`

## Basic Structure

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

## Full Example

```yaml
# ── Global HTTP client ──────────────────────────────────
http_client:
  timeout: 30s
  max_response_size: 1048576
  user_agent: "swag2mcp-global/1.0"
  follow_redirects: true
  max_redirects: 10
  random: false
  proxy:
    url: ""
    username: ""
    password: ""
    bypass: []
  headers:
    "Accept": "application/json"
  cookies:
    - name: "session"
      value: "abc123"

# ── MCP server ──────────────────────────────────────────
mcp:
  transport: stdio
  addr: ":8080"
  path: "/mcp"
  auth:
    token: ""

# ── Mock server ─────────────────────────────────────────
mock_enabled: false
mock_auth:
  oauth2_port: 9090
  digest_port: 9091
  hmac_port: 9092

# ── Rate limiter ────────────────────────────────────────
disable_ratelimiter: false
rate_limit_interval: 10s

# ── Specs ───────────────────────────────────────────────
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    llm_instruction: "Use this API for weather forecasts and climate data"
    base_url: https://api.open-meteo.com
    disable: false
    tags: ["weather", "climate"]
    http_client:
      timeout: 10s
    auth:
      type: bearer
      config:
        token: "$(TOKEN)"
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
        disable: false
        http_client:
          timeout: 5s

  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Environment Variables

Use `$(VAR_NAME)` syntax to reference environment variables. swag2mcp resolves them at startup.

```yaml
specs:
  - domain: meteo
    auth:
      type: bearer
      config:
        token: "$(API_TOKEN)"

mcp:
  auth:
    token: "$(MCP_TOKEN)"
```

`$(VAR)` is resolved in:
- Auth config fields: `token`, `username`, `password`, `client_id`, `client_secret`, `api_key`, `secret_key`, `domain`
- MCP server auth token: `mcp.auth.token`
- HTTP client headers and cookie values

`$(VAR)` is **not** resolved in proxy settings, base URLs, or collection locations.

## Validation

```bash
# Validate default workspace (~/.swag2mcp)
swag2mcp validate

# Validate a custom project workspace
swag2mcp validate ./my-project
```

If the workspace is not in the home directory (e.g., inside a project repository), always specify the path when running `validate`, `update`, `mcp`, or any other command. Otherwise swag2mcp will use the default `~/.swag2mcp` workspace.
