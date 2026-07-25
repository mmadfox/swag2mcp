# OAuth2 Client Credentials

## Purpose

OAuth2 Client Credentials Grant — authentication for server-to-server communication. The application obtains a token using its client_id and client_secret, without user involvement.

## When to use

- Microservices and server-to-server integrations
- Machine-to-machine communication
- When the API uses OAuth2 and you have a client_id + client_secret

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
      type: oauth2-cc
      config:
        client_id: "$(CLIENT_ID)"
        client_secret: "$(CLIENT_SECRET)"
        token_url: "https://auth.example.com/oauth/token"
        scopes:
          - read
          - write
```

## Parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| `client_id` | Yes | Client identifier |
| `client_secret` | Yes | Client secret |
| `token_url` | Yes | Token endpoint URL |
| `scopes` | No | List of permissions (optional) |

## Notes

- swag2mcp automatically requests a new token when the current one expires
- The token is cached until its expiry time (`expires_in`)
- If the server doesn't provide `expires_in`, the token is considered valid for 1 hour
- All parameters can be stored in environment variables
