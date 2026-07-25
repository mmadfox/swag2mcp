# OAuth2 Password Grant

## Purpose

OAuth2 Resource Owner Password Grant — authentication using a user's username and password. Suitable for first-party applications where the user trusts the app with their credentials.

## When to use

- First-party applications (mobile, web)
- Integration with Keycloak and similar Identity Providers
- When the API supports OAuth2 Password Grant

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
        scopes:
          - openid
          - profile
```

## Parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| `client_id` | Yes | Client identifier |
| `username` | Yes | Username |
| `password` | Yes | Password |
| `token_url` | Yes | Token endpoint URL |
| `client_secret` | No | Client secret (optional, for public clients) |
| `scopes` | No | List of permissions (optional) |

## Notes

- `client_secret` is optional — **public clients** are supported (e.g., Keycloak)
- swag2mcp automatically refreshes the token when it expires
- The token is cached until expiry
- All parameters can be stored in environment variables
