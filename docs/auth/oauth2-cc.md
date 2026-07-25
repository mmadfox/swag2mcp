# OAuth2 Client Credentials

OAuth2 authentication via Client Credentials Grant.

## Configuration

```yaml
specs:
  - domain: "api.example.com"
    location: "https://api.example.com/openapi.json"
    auth:
      type: oauth2-cc
      oauth2_cc:
        client_id: "{{CLIENT_ID}}"
        client_secret: "{{CLIENT_SECRET}}"
        token_url: "https://auth.example.com/oauth/token"
        scopes: ["read", "write"]
```

## How It Works

1. swag2mcp requests a token from `token_url` with client_id and client_secret
2. The Bearer token is used for all requests
3. Token is automatically refreshed on expiry

## Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `client_id` | string | Client ID |
| `client_secret` | string | Client secret |
| `token_url` | string | Token endpoint URL |
| `scopes` | array | Scope list |
