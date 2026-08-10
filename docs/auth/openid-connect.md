# OpenID Connect (OIDC Discovery)

## Purpose

OpenID Connect Discovery-based authentication. swag2mcp fetches the OIDC discovery document from the issuer, discovers the token endpoint, and obtains a Bearer token via the **client credentials** grant. The token is cached until expiry.

## When to use

- APIs that expose an OIDC discovery document (`.well-known/openid-configuration`)
- Server-to-server integrations where you have a `client_id` + `client_secret`
- When the API uses OIDC and you want automatic token acquisition

## Configuration

```yaml
specs:
  - domain: my-api
    llm_title: My API
    base_url: https://api.example.com
    collections:
      - llm_title: Main
        location: https://example.com/spec.yaml
    auth:
      type: openid-connect
      config:
        issuer: "https://auth.example.com"
        client_id: "$(CLIENT_ID)"
        client_secret: "$(CLIENT_SECRET)"
        scopes:
          - openid
          - profile
```

## Parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| `issuer` | Yes | OIDC issuer URL. swag2mcp appends `/.well-known/openid-configuration` to discover the token endpoint |
| `client_id` | Yes | Client identifier |
| `client_secret` | Yes | Client secret |
| `scopes` | No | List of permissions (optional) |

## Notes

- swag2mcp automatically requests a new token when the current one expires
- The token is cached until its expiry time (`expires_in`)
- If the server doesn't provide `expires_in`, the token is considered valid for 1 hour
- `issuer`, `client_id`, and `client_secret` support `$(VAR)` syntax for environment variables
- The token is obtained via the **client credentials** grant (server-to-server, no user involvement)
- The discovery document must include a `token_endpoint` field
