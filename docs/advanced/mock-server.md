# Mock Server

swag2mcp-mock is a separate binary for running a mock server.

## Installation

The mock server is included with the main binary:

```bash
swag2mcp-mock mockserver
```

## Configuration

```yaml
mock_enabled: true
mock_auth:
  oauth2_port: 9090
  digest_port: 9091
  hmac_port: 9092

specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
        base_mock_url: "127.0.0.1:9090"
```

## Parameters

| Parameter | Description |
|-----------|-------------|
| `mock_enabled` | Enable mock mode globally |
| `mock_auth.oauth2_port` | OAuth2 mock server port (default: 9090) |
| `mock_auth.digest_port` | Digest mock server port (default: 9091) |
| `mock_auth.hmac_port` | HMAC mock server port (default: 9092) |
| `base_mock_url` | Mock server address per collection (host:port) |

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
