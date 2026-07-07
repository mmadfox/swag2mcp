# API Key in Query Parameter

Sends an API key as a URL query parameter on every request. The parameter
name and value are configurable.

## What it demonstrates

- `auth.type: api-key` with `in: query`
- `key` (query param name) and `value` (query param value) fields
- API key is appended to the URL as a query parameter

## Expected behavior

- Every `invoke` call includes `?api_key=my-api-key` in the URL
- The `auth` tool returns the query parameter
