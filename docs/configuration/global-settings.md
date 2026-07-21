# Global Settings

Global settings are the top-level configuration blocks in `swag2mcp.yaml`. They apply to all specs unless overridden at the spec or collection level.

## Structure

There are two global blocks:

```yaml
http_client:
  # HTTP client settings for all API calls

mcp:
  # MCP server settings
```

### HTTP Client

Controls how swag2mcp makes HTTP requests to APIs: timeout, response size limit, proxy, headers, cookies, redirects, and user-agent. These settings cascade down to specs and collections.

See [HTTP Client](./http-client) for all parameters and examples.

### MCP Server

Controls how the MCP server communicates with LLM agents: transport type (stdio, SSE, Streamable HTTP), address, path, and optional bearer token auth.

See [MCP Server](./mcp-server) for all parameters, transports, and startup flags.

## Cascade

Global settings can be overridden at the spec and collection levels, but **only `headers` and `cookies`** can be overridden. All other HTTP settings (timeout, proxy, user-agent, redirects, response size, randomizer) are global only.

```
Global (http_client)
    ↓ overrides (headers, cookies only)
Spec (specs[].http_client)
    ↓ overrides (headers only)
Collection (specs[].collections[].http_client)
```

See [Configuration Cascade](./cascade) for details.
