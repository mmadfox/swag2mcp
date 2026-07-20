# Authentication

swag2mcp supports 9 authentication methods for calling protected APIs.

## Methods

| Method | Description |
|--------|-------------|
| `none` | No authentication |
| `basic` | HTTP Basic Auth |
| `bearer` | Bearer Token (JWT) |
| `api-key` | API Key (header or query) |
| `digest` | HTTP Digest Auth |
| `hmac` | HMAC-SHA256 (Binance-style) |
| `oauth2-cc` | OAuth2 Client Credentials |
| `oauth2-pwd` | OAuth2 Password Grant |
| `script` | External script |

## Configuration

Auth is configured at the spec level:

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: bearer
      config:
        token: "my-token"
```

## MCP auth Tool

LLM agents can get tokens/headers via the `auth` tool:

```
→ auth(specId: "abc123")
← Authorization: Bearer eyJhbGci...
```

!!! note
    The `auth` tool is disabled with `--disable-llm-auth` for security.
