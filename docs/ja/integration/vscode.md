# VS Code Integration

## Via VS Code Settings

In `.vscode/settings.json`:

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

## Via Extension

Install the MCP extension for VS Code and add:

```json
{
  "mcp.servers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp"]
    }
  }
}
```

## Usage

After setup, VS Code AI Assistant can work with your APIs through swag2mcp.

## Others

Don't see your client? All MCP integrations follow the same pattern:
- Set the command to `swag2mcp` with argument `mcp`
- Optionally add a workspace path: `mcp /path/to/workspace`
- Check your client's documentation for the exact config file location and format

Most MCP clients support stdio transport, and some support HTTP (SSE / Streamable HTTP).
