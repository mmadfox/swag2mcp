# Mock Server

## Overview

The mock server generates fake API responses based on your OpenAPI schemas. It lets you test your API integration without making real HTTP calls. This is useful for development, testing LLM agents, and demonstrations.

The mock server is a **separate binary** — `swag2mcp-mock`. It is not included in the main `swag2mcp` binary and must be installed separately.

## Installation

```bash
# Option 1: Download from GitHub Releases
# Look for swag2mcp-mock_<version>_<os>_<arch>.tar.gz

# Option 2: Install with Go
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

## Configuration

Enable the mock server in your config:

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

### mock_enabled

- **Type:** `bool`
- **Default:** `false`
- **Effect:** When `true`, every active collection must have `base_mock_url` set. The mock server starts HTTP servers for each collection.

### mock_auth

Ports for mock authentication servers. These simulate OAuth2, Digest, and HMAC auth endpoints so you can test authenticated APIs without real credentials.

| Field | Default | Description |
|-------|---------|-------------|
| `oauth2_port` | `9090` | Port for the mock OAuth2 token server |
| `digest_port` | `9091` | Port for the mock Digest auth server |
| `hmac_port` | `9092` | Port for the mock HMAC auth server |

### base_mock_url (per collection)

- **Type:** `string`
- **Required:** Yes (when `mock_enabled: true`)
- **Format:** `host:port` (e.g., `localhost:8080`, `127.0.0.1:9000`)
- **Effect:** Each collection gets its own HTTP server on this address. The server responds to all endpoints defined in the spec with randomly generated data.

## Starting the mock server

```bash
# Start with default config
swag2mcp-mock mockserver

# Start with TLS
swag2mcp-mock mockserver --tls

# Start with custom TLS certificate
swag2mcp-mock mockserver --tls --tls-cert cert.pem --tls-key key.pem
```

### TLS flags

| Flag | Description |
|------|-------------|
| `--tls` | Enable TLS with a self-signed certificate |
| `--tls-cert` | Path to TLS certificate file |
| `--tls-key` | Path to TLS key file |

If `--tls` is set without `--tls-cert` and `--tls-key`, a self-signed certificate is generated automatically for `localhost`.

## What the mock server does

When you start the mock server, it:

1. **Parses all spec files** — reads each collection's OpenAPI/Swagger spec
2. **Registers handlers** — creates an HTTP handler for every path and method defined in the spec
3. **Generates fake data** — responds with randomly generated data that matches the response schema (correct types, formats, and structure)
4. **Starts auth servers** — simulates OAuth2, Digest, and HMAC auth endpoints for testing

### Testing the mock

```bash
# In one terminal:
swag2mcp-mock mockserver

# In another terminal:
curl http://localhost:8080/pets
# → [{"id":1,"name":"Pet_name","status":"available"}]
```

## How fake data is generated

The mock server generates realistic fake data based on the OpenAPI schema:

- **Strings** — random words, sentences, or format-specific values (email, URL, UUID, date, phone, etc.)
- **Numbers** — random integers and floats within the specified range
- **Booleans** — random true/false
- **Arrays** — 1 to 3 random items
- **Objects** — all properties filled with random values
- **Enums** — random value from the enum list
- **Nullable fields** — sometimes returns `null` (~10% chance)

## Use cases

- **Development** — test your integration without real API access
- **Testing LLM agents** — verify the LLM can discover, inspect, and invoke endpoints
- **Demonstrations** — show swag2mcp working without configuring real APIs
- **Load testing** — test the MCP server under load without hitting real APIs

## Important notes

- **Separate binary** — `swag2mcp-mock` is not included in the main `swag2mcp` binary. Install it separately.
- **Each collection gets its own port** — configure `base_mock_url` per collection
- **Auth mock servers are global** — OAuth2, Digest, and HMAC servers run on the configured ports regardless of how many collections you have
- **Spec parse failures are non-fatal** — if a collection's spec cannot be parsed, it is skipped with a warning
- **Self-signed TLS** — when using `--tls` without certificates, a self-signed certificate is generated for localhost only
