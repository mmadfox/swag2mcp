# swag2mcp

<p>
    <a href="https://github.com/mmadfox/swag2mcp/releases"><img src="https://img.shields.io/github/release/mmadfox/swag2mcp.svg" alt="Latest Release"></a>
    <a href="https://pkg.go.dev/github.com/mmadfox/swag2mcp?tab=doc"><img src="https://godoc.org/github.com/mmadfox/swag2mcp?status.svg" alt="GoDoc"></a>
    <a href="https://github.com/mmadfox/swag2mcp/actions"><img src="https://github.com/mmadfox/swag2mcp/actions/workflows/test.yml/badge.svg?branch=main" alt="Build Status"></a>
    <a href="https://coveralls.io/github/mmadfox/swag2mcp?branch=main"><img src="https://coveralls.io/repos/github/mmadfox/swag2mcp/badge.svg?branch=main" alt="Coverage Status"></a>
</p>

> ⚠️ **Work in progress** — API may change, contributions welcome.

**swag2mcp** bridges OpenAPI/Swagger/Postman API specifications with LLM agents via the Model Context Protocol (MCP).

- **For LLM agents**: 14 MCP tools for discovering, inspecting, and invoking APIs
- **For humans**: Interactive TUI explorer with full-text search
- **Zero integration code**: Just point to your specs and go

## Documentation

| Language | File |
|----------|------|
| English | [docs/README.md](docs/README.md) |
| Русский | [docs/README.ru.md](docs/README.ru.md) |
| Deutsch | [docs/README.de.md](docs/README.de.md) |
| Français | [docs/README.fr.md](docs/README.fr.md) |
| Español | [docs/README.es.md](docs/README.es.md) |
| 中文 | [docs/README.zh.md](docs/README.zh.md) |
| 日本語 | [docs/README.ja.md](docs/README.ja.md) |

## Quick Start

```bash
$ go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
$ swag2mcp init
$ swag2mcp mcp
OR
$ swag2mcp mcp --tags=project-1,work,petstore
```

## Examples

This directory contains example configurations for swag2mcp. Each example
demonstrates a specific feature or use case.

| Category | Example | Description |
|----------|---------|-------------|
| **Basics** | [minimal-config](examples/minimal-config) | Minimal configuration — one spec, one collection, no auth |
| | [full-config](examples/full-config) | Complete configuration with all features |
| **Auth** | [no-auth](examples/auth/no-auth) | No authentication |
| | [basic-auth](examples/auth/basic-auth) | HTTP Basic Authentication |
| | [bearer-auth](examples/auth/bearer-auth) | Bearer Token Authentication |
| | [digest-auth](examples/auth/digest-auth) | HTTP Digest Authentication |
| | [oauth2-client-credentials](examples/auth/oauth2-client-credentials) | OAuth2 Client Credentials Grant |
| | [oauth2-password](examples/auth/oauth2-password) | OAuth2 Password Grant |
| | [api-key-header](examples/auth/api-key-header) | API Key in HTTP Header |
| | [api-key-query](examples/auth/api-key-query) | API Key in Query Parameter |
| | [script-auth](examples/auth/script-auth) | Script-Based Authentication |
| **Spec Features** | [llm-metadata](examples/spec-features/llm-metadata) | LLM titles and instructions |
| | [disable-spec](examples/spec-features/disable-spec) | Disabling specs and collections |
| | [tags-filtering](examples/spec-features/tags-filtering) | Tag-based filtering with `--tags` |
| | [custom-headers](examples/spec-features/custom-headers) | Custom HTTP headers |
| | [multiple-collections](examples/spec-features/multiple-collections) | Multiple collections per spec |
| | [collection-override](examples/spec-features/collection-override) | Collection-level overrides |
| | [http-client-config](examples/spec-features/http-client-config) | HTTP client configuration (headers, cookies, timeout, redirects) |
| | [proxy-config](examples/spec-features/proxy-config) | Proxy configuration (SOCKS5, HTTP, HTTPS, bypass) |
| **MCP Transport** | [stdio](examples/mcp-transport/stdio) | Default stdio transport |
| | [sse](examples/mcp-transport/sse) | SSE transport with HTTP and bearer token auth |
| | [streamable-http](examples/mcp-transport/streamable-http) | Streamable HTTP transport with HTTP and bearer token auth |

## License

MIT
