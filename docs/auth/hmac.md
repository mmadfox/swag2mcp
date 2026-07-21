# HMAC Auth

HMAC-SHA256 authentication (Binance-style).

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

## How It Works

1. Query string is built from parameters
2. HMAC-SHA256 signature is computed
3. `X-MBX-APIKEY` and `signature` headers are added

## Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `api_key` | string | Public API key |
| `secret_key` | string | Secret key for signing |
