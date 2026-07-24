# Cursor Integration

## stdio

In Cursor settings, add the MCP server:

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

## Usage

After connecting, Cursor AI Agent can:

- Explore your APIs
- Find relevant endpoints
- Call APIs and show results
- Help debug requests

## Others

Don't see your client? All MCP integrations follow the same pattern:
- Set the command to `swag2mcp` with argument `mcp`
- Optionally add a workspace path: `mcp /path/to/workspace`
- Check your client's documentation for the exact config file location and format

Most MCP clients support stdio transport, and some support HTTP (SSE / Streamable HTTP).
