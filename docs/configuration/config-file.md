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
http_client:
  timeout: 30s
  max_response_size: 1048576
  headers:
    "User-Agent": "swag2mcp/1.0"

mcp:
  transport: sse
  addr: "127.0.0.1:8080"
  path: "/mcp"
  auth:
    token: "my-secret-token"

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
    auth:
      type: bearer
      config:
        token: "$(TOKEN)"

  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Environment Variables

Use `$(VAR_NAME)` syntax in auth config fields and MCP server token:

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

`$(VAR)` is **not** resolved in headers, cookies, proxy settings, base URLs, or collection locations.

## Validation

```bash
# Validate default workspace (~/.swag2mcp)
swag2mcp validate

# Validate a custom project workspace
swag2mcp validate ./my-project
```

If the workspace is not in the home directory (e.g., inside a project repository), always specify the path when running `validate`, `update`, `mcp`, or any other command. Otherwise swag2mcp will use the default `~/.swag2mcp` workspace.
