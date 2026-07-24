# Environment Variables

## Overview

swag2mcp supports environment variable substitution in the configuration file using `$(VAR_NAME)` syntax. This lets you keep sensitive data (tokens, passwords, keys) out of the YAML file.

## How it works

When swag2mcp starts, it scans the configuration for `$(VAR_NAME)` patterns and replaces them with the value of the corresponding environment variable.

```yaml
specs:
  - domain: meteo
    auth:
      type: bearer
      config:
        token: "$(API_TOKEN)"
```

If the environment variable `API_TOKEN` is set, it will be substituted. If it is not set, the value becomes empty.

## Where `$(VAR)` is resolved

| Field | Example |
|-------|---------|
| Auth `token` (bearer) | `token: "$(API_TOKEN)"` |
| Auth `username` / `password` (basic, digest) | `password: "$(API_PASSWORD)"` |
| Auth `client_id` / `client_secret` (oauth2-cc, oauth2-pwd) | `client_secret: "$(OAUTH_SECRET)"` |
| Auth `api_key` / `secret_key` (hmac) | `api_key: "$(BINANCE_API_KEY)"` |
| Auth `domain` (script) | `domain: "$(AUTH_DOMAIN)"` |
| MCP server token | `token: "$(MCP_TOKEN)"` |
| HTTP client headers | `"X-API-Key": "$(API_KEY)"` |
| HTTP client cookie values | `value: "$(SESSION_TOKEN)"` |

## Where `$(VAR)` is NOT resolved

- Base URLs (`base_url`)
- Collection locations (`location`)
- Spec domain names (`domain`)

## Example

```bash
export API_TOKEN="eyJhbGciOiJIUzI1NiIs..."
export MCP_TOKEN="my-secret-token"

swag2mcp mcp
```

## Security best practices

- **Never** store secrets directly in the YAML file
- Use environment variables or an external secret manager
- Add the YAML file to `.gitignore` if it contains any hardcoded secrets
- Set environment variables in your shell profile, IDE configuration, or deployment pipeline

## Syntax details

- `$(VAR_NAME)` — standard syntax
- `$( VAR_NAME )` — whitespace inside parentheses is allowed and trimmed
- `$()` — empty variable name returns the original string unchanged
- Nested `$(...)` patterns are not resolved
