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
