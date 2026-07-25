# HTTP Client Configuration

Demonstrates how to configure the HTTP client at the global, spec, and
collection levels. Settings cascade: collection overrides spec, spec overrides
global.

## What it demonstrates

- `http_client` at global level — applies to all specs and collections
- `http_client` at spec level — overrides global settings for that spec
- `http_client` at collection level — overrides spec and global for that collection
- `headers` — custom HTTP headers
- `cookies` — HTTP cookies with Name, Value, Domain, Path, Secure, HttpOnly
- `timeout` — request timeout (Go duration format: 30s, 5m, etc.)
- `follow_redirects` — enable/disable following redirects
- `max_redirects` — maximum number of redirects to follow

## Expected behavior

- All specs inherit `Accept: application/json` and `timeout: 30s` from global
- "api-a" adds custom headers
- "api-b" has cookies configured
- "api-b" collection "Billing" has its own headers
