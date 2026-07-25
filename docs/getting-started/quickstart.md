# Quick Start

Get swag2mcp running in 2 minutes.

## 1. Initialize

```bash
swag2mcp init
```

Creates `~/.swag2mcp/swag2mcp.yaml`.

## 2. Add an API

Add the Open-Meteo weather API:

```bash
swag2mcp add https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

## 3. Start MCP Server

```bash
swag2mcp mcp
```

Output:

```
MCP server listening on http://127.0.0.1:8080/mcp
```

## 4. Test It

In another terminal:

```bash
curl -X POST http://127.0.0.1:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"spec_list","arguments":{}}}'
```

## 5. LLM Client Configuration

=== "OpenCode"
    ```json
    {
      "mcp": {
        "swag2mcp": {
          "type": "local",
          "command": ["swag2mcp", "mcp"]
        }
      }
    }
    ```

=== "Claude Desktop"
    ```json
    {
      "mcpServers": {
        "swag2mcp": {
          "command": "swag2mcp",
          "args": ["mcp"]
        }
      }
    }
    ```

=== "Cursor"
    ```json
    {
      "mcpServers": {
        "swag2mcp": {
          "command": "swag2mcp",
          "args": ["mcp"]
        }
      }
    }
    ```

## What's Next?

- [Concepts](../concepts/overview.md) — understand the architecture
- [Configuration](../configuration/config-file.md) — customize settings
- [CLI Commands](../cli/overview.md) — full command reference
