# Examples

## Minimal Configuration

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Full Configuration

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

specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    disable: false
    http_client:
      headers:
        "X-Custom": "value"
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

  - domain: pokemon
    llm_title: PokeAPI
    base_url: https://pokeapi.co
    collections:
      - llm_title: Pokemon
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/pokeapi.yaml
```

## Authentication

All 9 auth methods with examples in `examples/auth/`:

- `no-auth/` — no auth
- `basic-auth/` — HTTP Basic
- `bearer-auth/` — Bearer Token
- `api-key-header/` — API Key in header
- `api-key-query/` — API Key in query
- `digest-auth/` — HTTP Digest
- `hmac-auth/` — HMAC-SHA256
- `oauth2-client-credentials/` — OAuth2 CC
- `oauth2-password/` — OAuth2 Password
- `script-auth/` — external script

## Transport

Examples in `examples/mcp-transport/`:

- `stdio/` — standard input/output
- `sse/` — Server-Sent Events
- `streamable-http/` — Streamable HTTP

## Mock Server

Example in `examples/mock-server/`:

```yaml
mock_enabled: true

specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
        base_mock_url: "127.0.0.1:9090"
```
