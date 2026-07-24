# Claude Desktop Integration

## stdio

In `claude_desktop_config.json`:

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

## Custom Workspace

```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "/path/to/workspace"]
    }
  }
}
```

## Usage

After restarting Claude Desktop, you can:

- "Show me the list of all APIs"
- "Find the endpoint for creating an order"
- "Call the weather API for Moscow"

## Others

Don't see your client? All MCP integrations follow the same pattern:
- Set the command to `swag2mcp` with argument `mcp`
- Optionally add a workspace path: `mcp /path/to/workspace`
- Check your client's documentation for the exact config file location and format

Most MCP clients support stdio transport, and some support HTTP (SSE / Streamable HTTP).
