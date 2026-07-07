# API Key in Header

Sends an API key as an HTTP header on every request. The key name and value
are configurable.

## What it demonstrates

- `auth.type: api-key` with `in: header`
- `key` (header name) and `value` (header value) fields
- API key is sent as a custom header

## Expected behavior

- Every `invoke` call includes `X-API-Key: sk-abc123`
- The `auth` tool returns the `X-API-Key` header
