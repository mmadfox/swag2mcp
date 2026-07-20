# FAQ

## General

### What is swag2mcp?

swag2mcp bridges OpenAPI/Swagger/Postman API specifications with LLM agents via the Model Context Protocol (MCP).

### How is it different?

- **Zero integration code** — no coding needed to connect APIs
- **19 MCP tools** — complete API interaction toolkit
- **9 auth methods** — support for any API
- **Full-text search** — fast endpoint discovery
- **TUI** — interactive interface

### What formats are supported?

OpenAPI 3.x, Swagger 2.0, Postman Collections v2.1.

## Installation

### How to install?

```bash
go install github.com/mmadfox/swag2mcp@latest
```

Or download from [GitHub Releases](https://github.com/mmadfox/swag2mcp/releases).

### Do I need Go?

No, you can download pre-built binaries from GitHub Releases.

## Usage

### How to add an API?

```bash
swag2mcp add https://api.example.com/openapi.json
```

### How to start the MCP server?

```bash
swag2mcp mcp
```

### How to change the port?

```bash
swag2mcp mcp --transport sse --http-addr 0.0.0.0:9090
```

## Authentication

### What auth methods are supported?

Basic, Bearer, Digest, HMAC, OAuth2 CC, OAuth2 Password, API Key, Script, None.

### How to pass a token?

Via config or environment variables:

```yaml
auth:
  type: bearer
  bearer:
    token: "$(MY_TOKEN)"
```

## Search

### How to search endpoints?

Via the `search` MCP tool:

```
search(query: "find pet by status", limit: 5)
```

### What search syntax is supported?

Fuzzy search, wildcards, field filters, boolean operators.

## Troubleshooting

### "command not found"

Ensure `$GOPATH/bin` is in `$PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### "config not found"

Create a config:

```bash
swag2mcp init
```

### Port already in use

Change the port:

```bash
swag2mcp mcp --http-addr 127.0.0.1:9090
```
