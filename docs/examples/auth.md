# Auth Examples

## Basic Auth

```yaml
specs:
  - domain: "api.example.com"
    auth:
      type: basic
      basic:
        username: "admin"
        password: "{{PASSWORD}}"
```

## Bearer Token

```yaml
specs:
  - domain: "api.example.com"
    auth:
      type: bearer
      bearer:
        token: "{{TOKEN}}"
```

## API Key (Header)

```yaml
specs:
  - domain: "api.example.com"
    auth:
      type: api-key
      api_key:
        name: "X-API-Key"
        in: header
        value: "{{API_KEY}}"
```

## OAuth2 Client Credentials

```yaml
specs:
  - domain: "api.example.com"
    auth:
      type: oauth2-cc
      oauth2_cc:
        client_id: "{{CLIENT_ID}}"
        client_secret: "{{CLIENT_SECRET}}"
        token_url: "https://auth.example.com/oauth/token"
        scopes: ["read", "write"]
```

## HMAC (Binance-style)

```yaml
specs:
  - domain: "api.binance.com"
    auth:
      type: hmac
      hmac:
        api_key: "{{API_KEY}}"
        secret_key: "{{SECRET_KEY}}"
```

Full examples in `examples/auth/`.
