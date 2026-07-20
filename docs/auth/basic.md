# Basic Auth

HTTP Basic Authentication.

## Configuration

```yaml
specs:
  - domain: "api.example.com"
    location: "https://api.example.com/openapi.json"
    auth:
      type: basic
      basic:
        username: "admin"
        password: "{{PASSWORD}}"
```

## How It Works

Header `Authorization: Basic base64(username:password)` is added to every request.

## Environment Variables

```yaml
basic:
  username: "admin"
  password: "$(API_PASSWORD)"
```
