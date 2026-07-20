# Collection Settings

Collection settings override spec settings for a specific endpoint group.

## collection Section

```yaml
specs:
  - domain: "api.example.com"
    location: "https://api.example.com/openapi.json"
    collections:
      - name: "users"
        tags: ["users", "auth"]
        headers:
          "X-Role": "admin"
        http_client:
          timeout: 5s
```

## Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `name` | string | — | Collection name |
| `tags` | array | `["*"]` | Tags to include |
| `headers` | map | `{}` | Headers for this collection |
| `http_client` | object | — | HTTP client override |
| `filter` | object | — | Endpoint filter |

## Tag Filtering

```yaml
collections:
  - name: "public"
    tags: ["public", "health"]
  - name: "admin"
    tags: ["admin", "management"]
```

## Endpoint Filter

Include/exclude by method and path:

```yaml
collections:
  - name: "read-only"
    tags: ["*"]
    filter:
      include:
        - method: GET
      exclude:
        - method: DELETE
        - method: POST
        - path: "/internal/*"
```

## Example: Splitting an API

```yaml
specs:
  - domain: "api.example.com"
    location: "https://api.example.com/openapi.json"
    collections:
      - name: "users"
        tags: ["users"]
        headers:
          "X-Scope": "users"
      - name: "orders"
        tags: ["orders"]
        headers:
          "X-Scope": "orders"
      - name: "admin"
        tags: ["admin"]
        headers:
          "X-API-Key": "{{ADMIN_KEY}}"
```
