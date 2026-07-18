# HMAC-SHA256 Authentication

This example demonstrates HMAC-SHA256 authentication (Binance-style).

## Configuration

```yaml
auth:
  type: hmac
  config:
    api_key: $(API_KEY)
    secret_key: $(SECRET_KEY)
```

## How it works

1. A `timestamp` query parameter is added with the current Unix millisecond
2. A `signature` query parameter is computed as HMAC-SHA256 of the query string using `secret_key`
3. An `X-MBX-APIKEY` header is sent with the `api_key` value

## Usage

```bash
export API_KEY=my-api-key
export SECRET_KEY=my-secret-key
swag2mcp mcp
```

## When to use

Use HMAC auth for APIs that require request signing, such as:
- Binance API
- Other cryptocurrency exchanges
- APIs with HMAC-based request authentication
