# Spec Settings

Spec settings override global settings for a specific API.

## spec Section

```yaml
specs:
  - domain: "api.example.com"
    location: "https://api.example.com/openapi.json"
    disabled: false
    headers:
      "X-API-Key": "my-key"
    cookies:
      - name: "session"
        value: "abc123"
    http_client:
      timeout: 10s
      max_response_size: 1024
    collections:
      - name: "default"
        tags: ["*"]
    auth:
      type: bearer
      bearer:
        token: "my-token"
```

## Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `domain` | string | — | Unique API identifier |
| `location` | string | — | URL or path to spec |
| `disabled` | bool | `false` | Disable this spec |
| `headers` | map | `{}` | Headers for this API |
| `cookies` | array | `[]` | Cookies for this API |
| `http_client` | object | — | HTTP client override |
| `collections` | array | — | Collection list |
| `auth` | object | — | Auth settings |

## Disabling a Spec

```yaml
specs:
  - domain: "old-api.example.com"
    location: "https://old-api.example.com/swagger.json"
    disabled: true
```

Disabled specs are not loaded or indexed.

## HTTP Client Override

```yaml
specs:
  - domain: "slow-api.example.com"
    location: "https://slow-api.example.com/openapi.json"
    http_client:
      timeout: 120s
      max_response_size: 8192
```
