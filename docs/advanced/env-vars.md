# Environment Variables

swag2mcp supports environment variables in configuration using `$(VAR_NAME)` syntax.

## Syntax

```yaml
specs:
  - domain: meteo
    auth:
      type: bearer
      config:
        token: "$(API_TOKEN)"
```

## Where Used

`$(VAR)` is resolved in these fields:

| Field | Example |
|-------|---------|
| Auth `token` (bearer) | `token: "$(API_TOKEN)"` |
| Auth `username` / `password` (basic, digest) | `password: "$(API_PASSWORD)"` |
| Auth `client_id` / `client_secret` (oauth2-cc, oauth2-pwd) | `client_secret: "$(OAUTH_SECRET)"` |
| Auth `api_key` / `secret_key` (hmac) | `api_key: "$(BINANCE_API_KEY)"` |
| Auth `domain` (script) | `domain: "$(AUTH_DOMAIN)"` |
| MCP server token | `token: "$(MCP_TOKEN)"` |

`$(VAR)` is **not** resolved in headers, cookies, proxy settings, base URLs, or collection locations.

## Example

```bash
export API_TOKEN="eyJhbGciOiJIUzI1NiIs..."
export MCP_TOKEN="my-secret-token"

swag2mcp mcp
```

## Security

Do not store secrets in the YAML file. Use environment variables or external secret managers.
