# Jeder MCP-Client

swag2mcp ist ein **MCP-Server** (Model Context Protocol). Das bedeutet, dass es mit **jedem MCP-Client** funktioniert — nicht nur mit den in diesem Abschnitt aufgeführten. Wenn Ihr Editor, Ihre IDE oder Ihr Agent das MCP-Protokoll unterstützt, können Sie swag2mcp daran anschließen.

## Universelles Muster

Jeder MCP-Client verwendet die gleiche grundlegende Einrichtung. Fügen Sie swag2mcp als MCP-Server hinzu mit:

- **Befehl:** `swag2mcp`
- **Argumente:** `mcp` (plus optionaler Workspace-Pfad: `mcp /path/to/workspace`)

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

Der genaue Speicherort und das Format der Konfigurationsdatei (JSON, TOML, GUI-Einstellungen) variieren je nach Client — **lesen Sie die MCP-Dokumentation Ihres Clients**, wo Sie sie ablegen.

## Transporte

- **stdio** — funktioniert überall; die meisten MCP-Clients unterstützen es
- **HTTP (SSE / Streamable HTTP)** — unterstützt von Clients mit HTTP-Transportoption

Siehe die Referenz zum [`mcp`-Befehl](/de/cli/mcp) für Transport-Flags.

## Getestete Integrationen

| Client | Anleitung |
|--------|-----------|
| OpenCode | [OpenCode](/de/integration/opencode) |
| Cursor | [Cursor](/de/integration/cursor) |
| Claude Desktop | [Claude Desktop](/de/integration/claude) |
| VS Code | [VS Code](/de/integration/vscode) |
| Crush | [Crush](/de/integration/crush) |

> Wenn Ihr Client nicht in der Liste steht, bedeutet das **nicht**, dass er nicht unterstützt wird. Solange er das MCP-Protokoll spricht, verwenden Sie das obige universelle Muster und folgen Sie dem Handbuch Ihres Clients zum Konfigspeicherort.
