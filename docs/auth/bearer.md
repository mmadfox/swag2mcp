# Bearer Auth

## Purpose

Bearer Token authentication — the most common method for modern REST APIs. The token is sent in the `Authorization: Bearer <token>` header.

## When to use

- Modern REST APIs
- JWT (JSON Web Tokens)
- OAuth2 access tokens (when the token is already obtained)
- Any API that accepts a Bearer Token

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
      type: bearer
      config:
        token: "eyJhbGciOiJIUzI1NiIs..."
```

## Parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| `token` | Yes | Bearer token (JWT, OAuth2 token, etc.) |

## Notes

- The token is static — if it expires, you need to update it in the config manually
- For automatic token refresh, use `oauth2-cc` or `oauth2-pwd`
- Store the token in an environment variable: `token: "$(API_TOKEN)"`
