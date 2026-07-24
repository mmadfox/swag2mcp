# OpenCode Integration

## stdio

```json
{
  "mcp": {
    "swag2mcp": {
      "type": "local",
      "command": ["swag2mcp", "mcp"],
      "enabled": true
    }
  }
}
```

## HTTP

```json
{
  "mcp": {
    "swag2mcp": {
      "type": "local",
      "command": ["swag2mcp", "mcp", "--transport", "sse", "--http-addr", "127.0.0.1:8080"],
      "enabled": true
    }
  }
}
```

## Example Queries

Once connected, you can ask:

- "What APIs do you have?"
- "Show all endpoints in petstore"
- "Find an API for creating a user"
- "Call GET /pet/1 and show the result"

## Others

Don't see your client? All MCP integrations follow the same pattern:
- Set the command to `swag2mcp` with argument `mcp`
- Optionally add a workspace path: `mcp /path/to/workspace`
- Check your client's documentation for the exact config file location and format

Most MCP clients support stdio transport, and some support HTTP (SSE / Streamable HTTP).
