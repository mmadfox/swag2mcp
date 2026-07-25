# HMAC Auth

## Purpose

HMAC-SHA256 request signing — the authentication method used by cryptocurrency exchanges (Binance, Bybit, and others). Each request is signed with a secret key.

## When to use

- Binance API and Binance-compatible exchanges
- Cryptocurrency trading platforms
- APIs that require request signing

## Configuration

```yaml
specs:
  - domain: binance
    llm_title: Binance Market Data
    base_url: https://api.binance.com
    collections:
      - llm_title: Market Data
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/binance.yaml
    auth:
      type: hmac
      config:
        api_key: "$(BINANCE_API_KEY)"
        secret_key: "$(BINANCE_SECRET_KEY)"
```

## Parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| `api_key` | Yes | Public API key |
| `secret_key` | Yes | Secret key for signing |

## Notes

- swag2mcp automatically adds a timestamp (Unix in milliseconds) to every request
- The signature is computed from all request parameters
- Store keys in environment variables: `api_key: "$(BINANCE_API_KEY)"`
- This method is compatible with Binance API and similar exchanges
