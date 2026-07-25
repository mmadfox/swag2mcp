# Cursor-Integration

## stdio

Fügen Sie in den Cursor-Einstellungen den MCP-Server hinzu:

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

Nach dem Verbinden kann der Cursor-KI-Agent:

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
