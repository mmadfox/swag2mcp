# Any MCP Client

swag2mcp is an **MCP server** (Model Context Protocol). This means it works with **any MCP client** — not just the ones listed in this section. If your editor, IDE, or agent supports the MCP protocol, you can connect swag2mcp to it.

## Universal pattern

Every MCP client follows the same basic setup. Add swag2mcp as an MCP server with:

- **Command:** `swag2mcp`
- **Arguments:** `mcp` (plus an optional workspace path: `mcp /path/to/workspace`)

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

The exact location and format of the config file (JSON, TOML, GUI settings) vary from client to client — **check your client's MCP documentation** for where to put this.

## Transports

- **stdio** — works everywhere; most MCP clients support it
- **HTTP (SSE / Streamable HTTP)** — supported by clients that offer an HTTP transport option

See the [`mcp` command](/cli/mcp) reference for transport flags.

## Tested integrations

| Client | Guide |
|--------|-------|
| OpenCode | [OpenCode](/integration/opencode) |
| Cursor | [Cursor](/integration/cursor) |
| Claude Desktop | [Claude Desktop](/integration/claude) |
| VS Code | [VS Code](/integration/vscode) |
| Crush | [Crush](/integration/crush) |

> If your client is not in the list, that does **not** mean it is unsupported. As long as it speaks the MCP protocol, use the universal pattern above and follow your client's manual for the config location.
