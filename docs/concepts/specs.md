# Specs

A spec is a file describing an API. swag2mcp supports three formats.

## Supported Formats

| Format | Extensions | Versions |
|--------|------------|----------|
| OpenAPI 3.x | `.json`, `.yaml`, `.yml` | 3.0.0 – 3.1.1 |
| Swagger 2.0 | `.json`, `.yaml`, `.yml` | 2.0 |
| Postman Collection | `.json` | v2.1 |

## Sources

Specs can be:

- **URL**: `https://api.example.com/openapi.json`
- **Local file**: `./specs/my-api.yaml`
- **Local file outside workspace**: `/home/user/api.yaml`

## Identification

Each spec gets a unique MD5 hash based on its domain:

```go
id = md5(domain)
```

## Management

```bash
# Add a spec
swag2mcp add https://api.example.com/openapi.json

# List all specs
swag2mcp ls

# Delete a spec
swag2mcp delete <id>

# Update a spec
swag2mcp update <id>
```

## Spec Configuration

In YAML config:

```yaml
specs:
  - domain: "api.example.com"
    location: "https://api.example.com/openapi.json"
    collections:
      - name: "users"
        tags: ["users", "auth"]
    headers:
      "X-API-Key": "{{API_KEY}}"
    disabled: false
```
