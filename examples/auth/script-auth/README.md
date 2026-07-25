# Script-Based Authentication

Executes an external shell script to obtain a bearer token. The script must
output a JSON object with a `token` field and optionally an `expires_in` field.
This is useful for custom or legacy auth flows.

## What it demonstrates

- `auth.type: script` configuration
- `domain` field (maps to script filename)
- Script location: `{workspaceDir}/auth_scripts/{domain}.sh` (Unix) or `.bat` (Windows)
- Script output format: `{"token": "...", "expires_in": N}`
- Token caching with configurable expiration

## Expected behavior

- First `invoke` call executes the script and caches the token
- Subsequent calls reuse the cached token
- The `auth` tool returns the bearer token
- After `expires_in` seconds, the script is re-executed
