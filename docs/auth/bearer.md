# Bearer Auth

Bearer Token authentication (JWT, OAuth2 tokens).

## Configuration

```yaml
specs:
  - domain: "api.example.com"
    location: "https://api.example.com/openapi.json"
    auth:
      type: bearer
      bearer:
        token: "eyJhbGciOiJIUzI1NiIs..."
```

## How It Works

Header `Authorization: Bearer <token>` is added to every request.

## Environment Variables

```yaml
bearer:
  token: "$(API_TOKEN)"
```
