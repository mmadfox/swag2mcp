# Bearer Token Authentication

Uses a static Bearer token sent as `Authorization: Bearer <token>` on every
request. The token can be provided directly or via an environment variable.

## What it demonstrates

- `auth.type: bearer` configuration
- `token` field with `$(ENV_VAR)` syntax
- Token is sent on every API call

## Expected behavior

- Every `invoke` call includes `Authorization: Bearer <token>`
- The `auth` tool returns the bearer token
- Environment variable is resolved at startup
