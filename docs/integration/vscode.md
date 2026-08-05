# VS Code Integration

## Via .vscode/mcp.json

1. Install the MCP extension for VS Code (e.g., MCP Client by org.mcp or similar).
2. Create `.vscode/mcp.json` in your project root:

```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "${workspaceFolder}"]
    }
  }
}
```

> // "${{workspaceFolder}}" will be passed as the workspace path

3. Reload the VS Code window (Ctrl+Shift+P → "Reload Window").
4. Use the AI assistant — it will now know about your APIs.

## Alternative: Via VS Code Settings

You can also configure in `.vscode/settings.json`:

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

After setup, VS Code AI Assistant can work with your APIs through swag2mcp.

## Others

Don't see your client? All MCP integrations follow the same pattern:
- Set the command to `swag2mcp` with argument `mcp`
- Optionally add a workspace path: `mcp /path/to/workspace`
- Check your client's documentation for the exact config file location and format

Most MCP clients support stdio transport, and some support HTTP (SSE / Streamable HTTP).
