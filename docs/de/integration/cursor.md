# Cursor-Integration

## stdio

### Via Settings UI

1. Cursor-Einstellungen öffnen (Cmd+, / Ctrl+,)
2. Go to **MCP-Server**
3. Click **Neuen Server hinzufügen**
4. Fill in:
   - **Name:** `swag2mcp`
   - **Typ:** `command`
   - **Befehl:** `swag2mcp mcp`
5. Click **Speichern**

### Über Konfigurationsdatei

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

## Verwendung

After connecting, Cursor AI Agent can:

- Ihre APIs erkunden
- Relevante Endpunkte finden
- APIs aufrufen und Ergebnisse anzeigen
- Beim Debuggen von Anfragen helfen

## Andere

Ihr Client ist nicht dabei? Alle MCP-Integrationen folgen dem gleichen Muster:
- Setzen Sie den Befehl auf `swag2mcp` mit dem Argument `mcp`
- Optional einen Arbeitsbereichspfad hinzufügen: `mcp /pfad/zu/arbeitsbereich`
- Überprüfen Sie die Dokumentation Ihres Clients für den genauen Konfigurationsdatei-Speicherort und das Format

Die meisten MCP-Clients unterstützen den stdio-Transport, und einige unterstützen HTTP (SSE / Streamable HTTP).
