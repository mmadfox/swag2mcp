# Cursor Integration

## stdio

### Via Settings UI

1. Open Cursor Settings (Cmd+, / Ctrl+,)
2. Go to **MCP Servers**
3. Click **Add new server**
4. Fill in:
   - **Name:** `swag2mcp`
   - **Type:** `command`
   - **Command:** `swag2mcp mcp`
5. Click **Save**

### Via config file

In `~/.cursor/mcp.json`:

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
