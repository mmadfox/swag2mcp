# OpenCode-Integration

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

## Beispiel-Abfragen

Sobald verbunden, können Sie fragen:

- "Welche APIs hast du?"
- "Zeige alle Endpunkte in petstore"
- "Finde eine API zum Erstellen eines Benutzers"
- "Rufe GET /pet/1 auf und zeige das Ergebnis"

## Andere

Ihr Client ist nicht dabei? Alle MCP-Integrationen folgen dem gleichen Muster:
- Setzen Sie den Befehl auf `swag2mcp` mit dem Argument `mcp`
- Optional einen Arbeitsbereichspfad hinzufügen: `mcp /pfad/zu/arbeitsbereich`
- Überprüfen Sie die Dokumentation Ihres Clients für den genauen Konfigurationsdatei-Speicherort und das Format

Die meisten MCP-Clients unterstützen den stdio-Transport, und einige unterstützen HTTP (SSE / Streamable HTTP).
