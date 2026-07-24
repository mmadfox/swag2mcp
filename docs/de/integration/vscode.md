# VS Code-Integration

## Über VS Code-Einstellungen

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

## Über Erweiterung

Installieren Sie die MCP-Erweiterung für VS Code und fügen Sie hinzu:

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

## Verwendung

Nach der Einrichtung kann der VS Code-KI-Assistent über swag2mcp mit Ihren APIs arbeiten.

## Andere

Ihr Client ist nicht dabei? Alle MCP-Integrationen folgen dem gleichen Muster:
- Setzen Sie den Befehl auf `swag2mcp` mit dem Argument `mcp`
- Optional einen Arbeitsbereichspfad hinzufügen: `mcp /pfad/zu/arbeitsbereich`
- Überprüfen Sie die Dokumentation Ihres Clients für den genauen Konfigurationsdatei-Speicherort und das Format

Die meisten MCP-Clients unterstützen den stdio-Transport, und einige unterstützen HTTP (SSE / Streamable HTTP).
