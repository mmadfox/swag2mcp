# swag2mcp Examples

This directory contains ready-to-use configuration examples for swag2mcp.
Each example demonstrates a specific feature or use case.

## How to use

1. Copy the `config.yaml` from any example into your `swag2mcp.yaml`
2. Adjust values (URLs, tokens, paths) to match your environment
3. Run `swag2mcp validate` to check the config
4. Run `swag2mcp update` to cache remote specs
5. Start the MCP server: `swag2mcp mcp`

---

## Basics

| Example | Description |
|---------|-------------|
| [minimal-config](minimal-config) | One spec, one collection, no auth — the absolute minimum |
| [full-config](full-config) | Every feature in a single file — reference config |

## Auth

| Example | Description |
|---------|-------------|
| [no-auth](auth/no-auth) | No authentication |
| [basic-auth](auth/basic-auth) | HTTP Basic Authentication |
| [bearer-auth](auth/bearer-auth) | Bearer Token Authentication |
| [digest-auth](auth/digest-auth) | HTTP Digest Authentication |
| [oauth2-client-credentials](auth/oauth2-client-credentials) | OAuth2 Client Credentials Grant |
| [oauth2-password](auth/oauth2-password) | OAuth2 Password Grant |
| [api-key-header](auth/api-key-header) | API Key in HTTP Header |
| [api-key-query](auth/api-key-query) | API Key in Query Parameter |
| [script-auth](auth/script-auth) | Script-Based Authentication |
| [hmac-auth](auth/hmac-auth) | HMAC-SHA256 Authentication (Binance-style) |

## Spec Features

| Example | Description |
|---------|-------------|
| [llm-metadata](spec-features/llm-metadata) | LLM titles and instructions |
| [disable-spec](spec-features/disable-spec) | Disabling specs and collections |
| [tags-filtering](spec-features/tags-filtering) | Tag-based filtering with `--tags` |
| [custom-headers](spec-features/custom-headers) | Custom HTTP headers |
| [multiple-collections](spec-features/multiple-collections) | Multiple collections per spec |
| [collection-override](spec-features/collection-override) | Collection-level overrides |
| [http-client-config](spec-features/http-client-config) | HTTP client configuration (headers, cookies, timeout, redirects) |
| [proxy-config](spec-features/proxy-config) | Proxy configuration (SOCKS5, HTTP, HTTPS, bypass) |
| [random-client](spec-features/random-client) | Random browser-like headers |

## MCP Transport

| Example | Description |
|---------|-------------|
| [stdio](mcp-transport/stdio) | Default stdio transport (local agent) |
| [sse](mcp-transport/sse) | SSE transport with HTTP and bearer token auth |
| [streamable-http](mcp-transport/streamable-http) | Streamable HTTP transport with HTTP and bearer token auth |

## Mock Server

| Example | Description |
|---------|-------------|
| [mock-server](mock-server) | Mock server with random data generation and auth mock |
