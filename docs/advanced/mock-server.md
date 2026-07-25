# Mock Server

swag2mcp-mock is a separate binary for running a mock server.

## Installation

The mock server is included with the main binary:

```bash
swag2mcp-mock mockserver
```

## Configuration

```yaml
# mock-config.yaml
specs:
  - domain: "api.example.com"
    location: "https://api.example.com/openapi.json"
    mock:
      enabled: true
      delay: 100ms
      error_rate: 0.05
```

## Parameters

| Parameter | Description |
|-----------|-------------|
| `enabled` | Enable mock mode |
| `delay` | Response delay |
| `error_rate` | Error probability (0.0 - 1.0) |

## Flags

| Flag | Description |
|------|-------------|
| `--tls` | Enable TLS |
| `--tls-cert` | TLS cert path |
| `--tls-key` | TLS key path |

## Usage

```bash
# Start mock server
swag2mcp-mock mockserver

# With TLS
swag2mcp-mock mockserver --tls --tls-cert cert.pem --tls-key key.pem
```

## Use Cases

- Development without real API
- Testing LLM agents
- Demonstrations
- Load testing
