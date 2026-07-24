# mcp

## Purpose

Start the **MCP (Model Context Protocol) server** — the primary mode for LLM integration. This is what you run to give an AI agent (Claude, Cursor, OpenCode, etc.) access to your APIs through 16 MCP tools.

## When to use

- You want to connect an LLM agent to your APIs
- You are configuring an IDE (VS Code, Cursor, JetBrains) or desktop app (Claude Desktop)
- You need to expose your APIs via the MCP protocol
- You are testing the MCP server before integration

## Syntax

```bash
swag2mcp mcp [path] [flags]
```

## Arguments

| Argument | Position | Required | Description |
|----------|----------|----------|-------------|
| `path` | 1 | No | Workspace directory. If omitted, resolves via path resolution rules. |

## Flags

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--transport` | | `string` | `"stdio"` | MCP transport: `stdio`, `sse`, `streamable-http` |
| `--http-addr` | | `string` | `":8080"` | HTTP server address (for `sse` and `streamable-http`) |
| `--http-path` | | `string` | `"/mcp"` | HTTP path for the MCP handler |
| `--auth-token` | | `string` | `""` | Bearer token for HTTP transport authentication |
| `--logfile` | `-f` | `string` | `""` | Log file path. If unset, logs to stderr. |
| `--disable-llm-auth` | | `bool` | `true` | Remove the `auth` tool from the MCP tool list |
| `--dump-dir` | | `string` | `""` | Directory to dump HTTP requests for debugging |
| `--tags` | `-t` | `string` | `""` | Filter specs by tags (comma-separated) |

## How it works

### stdio transport (default)

Used when the MCP server is launched as a subprocess by the LLM client (IDE, Claude Desktop, etc.). The server communicates over standard input/output.

```bash
swag2mcp mcp
```

### SSE transport

Server-Sent Events transport for HTTP-based communication. Requires the MCP handshake sequence.

```bash
swag2mcp mcp --transport sse --http-addr :8080
```

### Streamable HTTP transport

Modern HTTP transport that supports streaming responses.

```bash
swag2mcp mcp --transport streamable-http --http-addr 0.0.0.0:8080
```

### With authentication

Protect the HTTP endpoint with a bearer token:

```bash
swag2mcp mcp --transport sse --http-addr :8080 --auth-token "my-secret"
```

### With tag filtering

Only load specs with specific tags:

```bash
swag2mcp mcp --tags=public
```

### With auth tool enabled (debug mode)

Allow the LLM to request fresh tokens via the `auth` tool:

```bash
swag2mcp mcp --disable-llm-auth=false
```

### With request dump directory

Save all HTTP requests for debugging:

```bash
swag2mcp mcp --dump-dir ./dumps
```

## MCP HTTP Transport — Handshake Protocol

When using `sse` or `streamable-http`, the MCP protocol requires a specific handshake. Tool calls will fail before initialization:

```
Step 1: POST /mcp → {"method":"initialize", ...}
Step 2: POST /mcp → {"method":"notifications/initialized"}
Step 3: POST /mcp → {"method":"tools/list", ...}   ← now works
```

### Health check

Works without initialization:

```bash
curl http://localhost:8080/health
# → {"status":"ok","version":"v1.2.0"}
```

## IDE Configuration Examples

### VS Code (`.vscode/settings.json` or global settings)

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

### Cursor / Windsurf (`~/.cursor/mcp.json` or project `.cursor/mcp.json`)

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

### Claude Desktop (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS)

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

### JetBrains IDEs (Settings → Tools → MCP)

- Name: `swag2mcp`
- Command: `swag2mcp`
- Arguments: `mcp /absolute/path/to/.swag2mcp`

> **Always use an absolute path** to the workspace directory in IDE config. Relative paths may fail depending on the IDE's working directory.

## Output

On success, the server prints:

```
MCP server listening on http://127.0.0.1:8080/mcp
```

## Nuances

- **No auto-init:** If the config file does not exist, `mcp` returns an error: `"configuration not found at <path>"`. Run `init` first.
- **`--disable-llm-auth` (default: `true`):** When enabled, the `auth` tool is removed from the MCP tool list entirely. The LLM cannot see or request tokens. Auth still works — tokens are obtained through the standard config mechanism, not via the LLM. This mode is recommended for **production**. For **debugging** or when using short-lived tokens, set `--disable-llm-auth=false` to let the LLM request fresh tokens via the `auth` tool.
- **YAML config fallback:** If a CLI flag is not explicitly set, the value is taken from the `mcp` section in `swag2mcp.yaml` (if present). This allows you to configure the server in the config file instead of passing flags every time.
- **Response cleanup:** On startup, responses older than 48 hours are automatically removed from the `responses/` directory.
- **Path resolution warning:** When `[path]` is omitted, `mcp` searches for `swag2mcp.yaml` in the current directory first, then falls back to `~/.swag2mcp/`. If you run the command from the wrong directory, it may load a different workspace than intended. **Always specify `[path]` explicitly when running as a service or in IDE config.**
