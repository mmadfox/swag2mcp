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
