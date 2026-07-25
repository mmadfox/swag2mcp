# FAQ

## General

### What is swag2mcp and what problem does it solve?

swag2mcp bridges OpenAPI/Swagger/Postman API specifications with LLM agents via the Model Context Protocol (MCP). Instead of writing custom code to connect each API to an AI agent, you configure it once in a YAML file and the LLM gets 19 tools to discover, inspect, and call your APIs.

### How is it different from other API-to-LLM tools?

- **No coding required** — configure APIs in YAML, no integration code needed
- **19 MCP tools** — complete toolkit from discovery to invocation to large response handling
- **9 auth methods** — works with any API authentication scheme
- **Full-text search** — bluge-powered search across all endpoints
- **TUI explorer** — interactive terminal interface for browsing and testing
- **Mock server** — test without real API calls

### What API specification formats are supported?

OpenAPI 3.x, Swagger 2.0, and Postman Collections v2.1.

### What is the difference between a spec and a collection?

A **spec** represents a logical API service (e.g., "Open-Meteo Weather APIs"). A **collection** is one OpenAPI/Swagger/Postman file. A spec can have multiple collections — for example, when an API has separate spec files for different services (forecast, air quality, marine).

### What MCP transports are supported?

Three transports: `stdio` (default, for local LLM clients), `sse` (Server-Sent Events for remote clients), and `streamable-http` (modern HTTP streaming).

### Can I use swag2mcp with any LLM?

Yes, any LLM client that supports the MCP protocol: Claude Desktop, VS Code, Cursor, Windsurf, JetBrains IDEs, OpenCode, and others.

## Installation

### How do I install swag2mcp?

```bash
# Option 1: Download from GitHub Releases
# Go to https://github.com/mmadfox/swag2mcp/releases/latest
# Download the archive for your OS and architecture

# Option 2: Install with Go
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

### Do I need Go installed?

No. Pre-built binaries are available for Linux (amd64, arm64), macOS (amd64, arm64), and Windows (amd64) on the [GitHub Releases page](https://github.com/mmadfox/swag2mcp/releases).

### How do I install the mock server?

The mock server is a separate binary:

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

Or download `swag2mcp-mock_&lt;version&gt;_&lt;os&gt;_&lt;arch&gt;.tar.gz` from GitHub Releases.

## Getting Started

### How do I quickly get started?

```bash
# 1. Initialize a workspace
mkdir -p .swag2mcp && swag2mcp init ./.swag2mcp

# 2. Start the MCP server (public example specs are included after init)
swag2mcp mcp
```

After `init`, the workspace already includes several public example specs (icanhazdadjoke, Open-Meteo, Binance, PokéAPI). You can start the MCP server immediately — no need to add specs manually.

If you want to add your own API instead:

```bash
swag2mcp add spec --yaml - <<EOF
domain: dadjoke
llm_title: icanhazdadjoke API
base_url: https://icanhazdadjoke.com
collections:
  - llm_title: Jokes
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
EOF
```

### How do I connect swag2mcp to my IDE?

**VS Code** (`.vscode/settings.json`):
```json
{
  "mcp": {
    "servers": {
      "swag2mcp": {
        "command": "swag2mcp",
        "args": ["mcp", "/absolute/path/to/.swag2mcp"]
      }
    }
  }
}
```

**Cursor** (`~/.cursor/mcp.json`):
```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "/absolute/path/to/.swag2mcp"]
    }
  }
}
```

**Claude Desktop** (`claude_desktop_config.json`):
```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "/absolute/path/to/.swag2mcp"]
    }
  }
}
```

Always use an absolute path to the workspace directory.

## Configuration

### Where is the config file located?

Default: `~/.swag2mcp/swag2mcp.yaml`. You can also create it in any directory and pass the path to commands.

### How do I add an API?

```bash
# Interactive mode
swag2mcp add spec

# With YAML (recommended for scripting)
swag2mcp add spec --yaml - <<EOF
domain: my-api
llm_title: My API
base_url: https://api.example.com/v1
collections:
  - llm_title: Main
    location: https://example.com/spec.yaml
EOF
```

### How do I add a collection to an existing spec?

```bash
swag2mcp add collection --yaml - <<EOF
spec_domain: meteo
llm_title: Air Quality
location: https://example.com/air-quality.yaml
EOF
```

### How do I disable a spec temporarily?

Set `disable: true` in the spec config. The spec will not be loaded or indexed.

### Can I filter which specs are loaded?

Yes, use the `--tags` flag: `swag2mcp mcp --tags=public`. Only specs with matching tags will be loaded.

### How do I use environment variables for secrets?

Use `$(VAR_NAME)` syntax in auth fields:

```yaml
auth:
  type: bearer
  config:
    token: "$(MY_API_TOKEN)"
```

Set the variable before starting: `export MY_API_TOKEN="eyJhbGci..."`

## Authentication

### What auth methods are supported?

Nine methods: `none`, `basic`, `bearer`, `digest`, `hmac`, `oauth2-cc` (client credentials), `oauth2-pwd` (password grant), `api-key`, and `script`.

### How do I pass a token?

Via the config file or environment variables:

```yaml
auth:
  type: bearer
  config:
    token: "$(MY_TOKEN)"
```

### Do I need to call auth before invoke?

No. The `invoke` tool automatically applies authentication from the spec's config. You only need the `auth` MCP tool if you want to show the token to the user (e.g., for a curl command).

### Why is the auth tool not showing up?

The `auth` tool is disabled by default (`--disable-llm-auth=true`). This is a security measure for production. To enable it: `swag2mcp mcp --disable-llm-auth=false`.

### How do OAuth2 tokens refresh?

OAuth2 Client Credentials and Password Grant tokens are automatically refreshed when they expire. Bearer tokens are static and must be updated manually.

## MCP Server

### How do I start the MCP server?

```bash
# Default (stdio transport)
swag2mcp mcp

# With HTTP transport
swag2mcp mcp --transport sse --http-addr :8080
```

### How do I change the port?

```bash
swag2mcp mcp --transport sse --http-addr 0.0.0.0:9090
```

### How do I secure the MCP HTTP endpoint?

Set a bearer token:

```bash
swag2mcp mcp --transport sse --http-addr :8080 --auth-token "my-secret"
```

The LLM client must include `Authorization: Bearer my-secret` in every request.

### What is the MCP handshake for HTTP transport?

For SSE and Streamable HTTP transports, the MCP protocol requires a three-step handshake:

```
Step 1: POST /mcp → {"method":"initialize", ...}
Step 2: POST /mcp → {"method":"notifications/initialized"}
Step 3: POST /mcp → {"method":"tools/list", ...}  ← now works
```

Tool calls will fail before initialization.

## Usage

### How do I search for endpoints?

Use the `search` MCP tool or the TUI (`swag2mcp run`). The search supports field filters (`method:GET`, `tag:pets`), fuzzy search, wildcards, and boolean operators.

### How do I call an API?

The LLM uses the `invoke` MCP tool. Always inspect the endpoint first to understand required parameters:

```
inspect(endpointId: "...")  → understand the contract
invoke(endpointId: "...", parameters: {...})  → make the call
```

### What happens if a response is too large?

Responses exceeding `max_response_size` (default 1 MB) are saved to disk. The LLM receives a file reference and can explore it with `response_outline`, `response_compress`, and `response_slice` tools.

### How does the rate limiter work?

Each endpoint has a 10-second cooldown. If the LLM calls the same endpoint twice within 10 seconds, the second call is silently blocked. You can disable or adjust this in the config.

### Can I test without making real API calls?

Yes, use the mock server:

```bash
swag2mcp-mock mockserver
```

It generates fake responses based on OpenAPI schemas.

## Workspace Management

### How do I backup my configuration?

```bash
swag2mcp export --output ~/backups/swag2mcp-2026-07-24.zip
```

### How do I transfer to another machine?

```bash
# On the old machine
swag2mcp export --output swag2mcp.zip

# Copy the ZIP, then on the new machine
swag2mcp import --from-zip swag2mcp.zip
```

### How do I update spec files?

```bash
swag2mcp update
```

This re-validates the config, clears the cache, and re-downloads all spec files.

### How do I clean up disk space?

```bash
swag2mcp clean
```

Removes cached spec files and saved API responses. Old responses (>48h) are also cleaned automatically on MCP server start.

## TUI

### What is the TUI and how do I use it?

The TUI (Terminal User Interface) is an interactive API explorer. Launch it with `swag2mcp run`. It has three modes: Search (full-text search), Browse (tree navigation: Spec → Collection → Tag → Endpoint), and Auth (view tokens).

### What are the keyboard shortcuts?

| Key | Action |
|-----|--------|
| `↑/↓` | Navigate |
| `Enter` | Select |
| `Esc` | Back |
| `Tab` | Switch modes |
| `/` | Search |
| `N/P` | Next/previous page |
| `q` | Quit |

## Advanced

### Can I use a proxy?

Yes, configure it in `http_client.proxy`:

```yaml
http_client:
  proxy:
    url: "http://proxy.company.com:8080"
    username: "$(PROXY_USER)"
    password: "$(PROXY_PASS)"
    bypass:
      - "localhost"
      - "*.internal.com"
```

### Can I add a custom auth method?

Yes, implement the `Authenticator` interface in `internal/auth/` and register it in the config parser. See the Development section for details.

### Can I add a custom MCP tool?

Yes, add a method to the `Svc` interface, implement it in the service layer, add a handler, and register it. See the Development section for details.

### What is the difference between `swag2mcp` and `swag2mcp-mock`?

`swag2mcp` is the main binary with CLI commands and the MCP server. `swag2mcp-mock` is a separate binary that starts mock servers for testing without real API calls.
