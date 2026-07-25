# Mock Server Examples

## Basic Configuration

```yaml
specs:
  - domain: "api.example.com"
    location: "https://api.example.com/openapi.json"
    mock:
      enabled: true
```

## With Delay and Errors

```yaml
specs:
  - domain: "api.example.com"
    location: "https://api.example.com/openapi.json"
    mock:
      enabled: true
      delay: 200ms
      error_rate: 0.1
```

## Launch

```bash
swag2mcp-mock mockserver
```

With TLS:

```bash
swag2mcp-mock mockserver --tls --tls-cert cert.pem --tls-key key.pem
```

Full example in `examples/mock-server/`.
