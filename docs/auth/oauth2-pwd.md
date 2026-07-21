# OAuth2 Password

OAuth2 authentication via Password Grant (Resource Owner Password Credentials).

## Configuration

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: oauth2-pwd
      config:
        client_id: "$(CLIENT_ID)"
        client_secret: "$(CLIENT_SECRET)"
        username: "$(USERNAME)"
        password: "$(PASSWORD)"
        token_url: "https://auth.example.com/oauth/token"
        scopes: ["openid", "profile"]
```

## How It Works

1. swag2mcp sends username + password to `token_url`
2. The Bearer token is used for all requests
3. Token is automatically refreshed on expiry

## Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `client_id` | string | Client ID |
| `client_secret` | string | Client secret (optional, for public client) |
| `username` | string | Username |
| `password` | string | Password |
| `token_url` | string | Token endpoint URL |
| `scopes` | array | Scope list (optional) |

!!! tip "Public Client"
    `client_secret` is optional — public clients are supported (e.g., Keycloak).
