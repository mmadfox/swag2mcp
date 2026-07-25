# HMAC Auth

HMAC-SHA256 authentication (Binance-style).

## Configuration

```yaml
specs:
  - domain: "api.binance.com"
    location: "https://api.binance.com/openapi.json"
    auth:
      type: hmac
      hmac:
        api_key: "{{BINANCE_API_KEY}}"
        secret_key: "{{BINANCE_SECRET_KEY}}"
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
