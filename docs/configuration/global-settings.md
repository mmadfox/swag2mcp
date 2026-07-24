# Global Settings

Global settings are the top-level configuration blocks in `swag2mcp.yaml`. They apply to all specs unless overridden at the spec or collection level.

## Structure

```yaml
http_client:
  # HTTP client settings for all API calls

mcp:
  # MCP server settings

mock_enabled: false
mock_auth:
  oauth2_port: 9090
  digest_port: 9091
  hmac_port: 9092

disable_ratelimiter: false
rate_limit_interval: 10s
```

## HTTP Client

Controls how swag2mcp makes HTTP requests to APIs: timeout, response size limit, proxy, headers, cookies, redirects, and user-agent. These settings cascade down to specs and collections.

See [HTTP Client](./http-client) for all parameters and examples.

## MCP Server

Controls how the MCP server communicates with LLM agents: transport type (stdio, SSE, Streamable HTTP), address, path, and optional bearer token auth.

See [MCP Server](./mcp-server) for all parameters, transports, and startup flags.

## Mock Server

The mock server generates fake API responses based on OpenAPI schemas. Useful for testing without hitting real APIs.

```yaml
mock_enabled: true
mock_auth:
  oauth2_port: 9090
  digest_port: 9091
  hmac_port: 9092
```

### mock_enabled

- **Type:** `bool`
- **Default:** `false`
- **Effect:** When `true`, swag2mcp starts mock servers for all specs that have `base_mock_url` configured. Each collection must have `base_mock_url` set.
- **When to enable:** You want to test your API integration without making real HTTP calls. Mock servers return fake data based on the OpenAPI schema.

### mock_auth

Port configuration for mock authentication servers. These are used when testing auth methods (OAuth2, Digest, HMAC) with the mock server.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `oauth2_port` | int | `9090` | Port for the mock OAuth2 token server (1024-65535) |
| `digest_port` | int | `9091` | Port for the mock Digest auth server (1024-65535) |
| `hmac_port` | int | `9092` | Port for the mock HMAC auth server (1024-65535) |

## Rate Limiter

The rate limiter prevents the LLM from calling the same API endpoint too frequently. By default, each endpoint can be called once every 10 seconds.

```yaml
disable_ratelimiter: false
rate_limit_interval: 10s
```

### disable_ratelimiter

- **Type:** `bool`
- **Default:** `false`
- **Effect:** When `true`, the per-endpoint rate limiter is disabled entirely. The LLM can call the same endpoint repeatedly without waiting.
- **When to enable:** Testing, debugging, or when you need to call the same endpoint multiple times in quick succession.
- **When to keep disabled (recommended):** Production. The rate limiter prevents accidental abuse and respects API rate limits.

### rate_limit_interval

- **Type:** duration (Go format: `10s`, `30s`, `1m`)
- **Default:** `10s`
- **Effect:** Sets how long the LLM must wait between calls to the same endpoint.
- **When to change:** Increase for APIs with strict rate limits. Decrease for internal APIs where you control the load.
- **Range:** Any valid duration (e.g., `5s`, `30s`, `1m`, `2m`).

## Cascade

Global settings can be overridden at the spec and collection levels. All `http_client` settings (timeout, proxy, user-agent, redirects, response size, randomizer, headers, cookies) can be overridden at both spec and collection levels.

```
Global (http_client, mock_enabled, disable_ratelimiter, rate_limit_interval)
    ↓ overrides (http_client only)
Spec (specs[].http_client)
    ↓ overrides (http_client only)
Collection (specs[].collections[].http_client)
```

See [Configuration Cascade](./cascade) for details.
