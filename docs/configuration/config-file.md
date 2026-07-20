# Configuration File

swag2mcp uses a YAML configuration file. Created by `swag2mcp init`.

## Location

Default: `~/.swag2mcp/swag2mcp.yaml`

## Basic Structure

```yaml
specs:
  - domain: "api.example.com"
    location: "https://api.example.com/openapi.json"
    collections:
      - name: "default"
        tags: ["*"]
```

## Full Example

```yaml
global:
  http_client:
    timeout: 30s
    max_response_size: 2048
    proxy: "http://proxy:8080"
    headers:
      "User-Agent": "swag2mcp/1.0"
  mcp:
    transport: "stdio"
    http_addr: "127.0.0.1:8080"
    http_path: "/mcp"
    auth_token: "my-secret-token"

specs:
  - domain: "petstore.swagger.io"
    location: "https://petstore.swagger.io/v2/swagger.json"
    disabled: false
    headers:
      "X-API-Key": "{{API_KEY}}"
    collections:
      - name: "pets"
        tags: ["pet"]
      - name: "store"
        tags: ["store"]
    auth:
      type: api-key
      api_key:
        name: "X-API-Key"
        in: header
        value: "{{API_KEY}}"
```

## Environment Variables

Use `$(VAR_NAME)` syntax:

```yaml
headers:
  "Authorization": "Bearer $(MY_TOKEN)"
```

## Validation

```bash
swag2mcp validate
```
