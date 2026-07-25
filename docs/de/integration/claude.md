# Claude Desktop-Integration

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

## Benutzerdefinierter Arbeitsbereich

```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "/pfad/zu/arbeitsbereich"]
    }
  }
}
```

## Verwendung

Nach dem Neustart von Claude Desktop können Sie:

- "Zeig mir die Liste aller APIs"
- "Finde den Endpunkt zum Erstellen einer Bestellung"
- "Rufe die Wetter-API für Moskau auf"

## Andere

Ihr Client ist nicht dabei? Alle MCP-Integrationen folgen dem gleichen Muster:
- Setzen Sie den Befehl auf `swag2mcp` mit dem Argument `mcp`
- Optional einen Arbeitsbereichspfad hinzufügen: `mcp /pfad/zu/arbeitsbereich`
- Überprüfen Sie die Dokumentation Ihres Clients für den genauen Konfigurationsdatei-Speicherort und das Format

Die meisten MCP-Clients unterstützen den stdio-Transport, und einige unterstützen HTTP (SSE / Streamable HTTP).
