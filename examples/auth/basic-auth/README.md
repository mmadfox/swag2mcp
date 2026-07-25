# HTTP Basic Authentication

Uses HTTP Basic Authentication with a username and password. The credentials
are sent as an `Authorization: Basic <base64>` header on every request.

## What it demonstrates

- `auth.type: basic` configuration
- `username` and `password` fields
- Environment variable resolution with `$(VAR_NAME)`
- Credentials are sent on every API call

## Expected behavior

- Every `invoke` call includes `Authorization: Basic <base64>`
- The `auth` tool returns the `Authorization` header
- Environment variables are resolved at startup
