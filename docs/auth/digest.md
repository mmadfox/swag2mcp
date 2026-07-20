# Digest Auth

HTTP Digest Access Authentication.

## Configuration

```yaml
specs:
  - domain: "api.example.com"
    location: "https://api.example.com/openapi.json"
    auth:
      type: digest
      digest:
        username: "admin"
        password: "{{PASSWORD}}"
```

## How It Works

1. swag2mcp sends a request without auth
2. Server responds 401 with `WWW-Authenticate: Digest ...`
3. swag2mcp computes MD5 hashes and retries with `Authorization: Digest ...`

## Environment Variables

```yaml
digest:
  username: "admin"
  password: "$(API_PASSWORD)"
```
