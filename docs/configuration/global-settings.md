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

### HTTP Client

Controls how swag2mcp makes HTTP requests to APIs: timeout, response size limit, proxy, headers, cookies, redirects, and user-agent. These settings cascade down to specs and collections.

See [HTTP Client](./http-client) for all parameters and examples.

### MCP Server

Controls how the MCP server communicates with LLM agents: transport type (stdio, SSE, Streamable HTTP), address, path, and optional bearer token auth.

See [MCP Server](./mcp-server) for all parameters, transports, and startup flags.

### Mock Server

```yaml
mock_enabled: true
mock_auth:
  oauth2_port: 9090
  digest_port: 9091
  hmac_port: 9092
```

- `mock_enabled` — enables mock server mode. When `true`, each collection must have `base_mock_url` set.
- `mock_auth` — port configuration for mock auth servers (OAuth2, Digest, HMAC).

### Rate Limiter

```yaml
disable_ratelimiter: false
rate_limit_interval: 10s
```

- `disable_ratelimiter` — disables the per-endpoint 10-second rate limiter for the `invoke` tool. Set to `true` when testing or when you need to call the same endpoint repeatedly.
- `rate_limit_interval` — custom rate limit interval (Go duration format: `10s`, `30s`, `1m`). Default: `10s`.

## Cascade

Global settings can be overridden at the spec and collection levels. All `http_client` settings (timeout, proxy, user-agent, redirects, response size, randomizer, headers, cookies) can be overridden at both spec and collection levels.

```
Global (http_client)
    ↓ overrides (all settings)
Spec (specs[].http_client)
    ↓ overrides (all settings)
Collection (specs[].collections[].http_client)
```

See [Configuration Cascade](./cascade) for details.
