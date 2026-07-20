# Environment Variables

swag2mcp supports environment variables in configuration.

## Syntax

```yaml
http_client:
  headers:
    "Authorization": "Bearer $(MY_TOKEN)"
    "X-API-Key": "$(API_KEY)"
```

Variables use `$(VAR_NAME)` format.

## Example

```bash
export MY_TOKEN="eyJhbGciOiJIUzI1NiIs..."
export API_KEY="abc123def456"

swag2mcp mcp
```

## Where Used

- HTTP headers
- Cookies
- Auth parameters
- Proxy URL
- File paths

## Priority

Environment variables take precedence over values in the YAML file.

## Security

!!! warning
    Do not store secrets in the YAML file. Use environment variables or external secret managers.
