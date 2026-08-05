# VS Code-Integration

## Via .vscode/mcp.json

1. Installieren Sie die MCP-Erweiterung für VS Code (z.B. MCP Client von org.mcp oder ähnlich).
2. Erstellen Sie `.vscode/mcp.json` im Projektstamm:

```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "${workspaceFolder}"]
    }
  }
}
```

> // "${{workspaceFolder}}" wird als Arbeitsbereichspfad übergeben

3. Laden Sie das VS Code-Fenster neu (Strg+Umschalt+P → "Reload Window").
4. Nutzen Sie den KI-Assistenten — er kennt jetzt Ihre APIs.

## Alternative: Über VS Code-Einstellungen

Sie können auch in `.vscode/settings.json` konfigurieren:

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

Nach der Einrichtung kann der VS Code-KI-Assistent über swag2mcp mit Ihren APIs arbeiten.

## Andere

Ihr Client ist nicht dabei? Alle MCP-Integrationen folgen dem gleichen Muster:
- Setzen Sie den Befehl auf `swag2mcp` mit dem Argument `mcp`
- Optional einen Arbeitsbereichspfad hinzufügen: `mcp /pfad/zu/arbeitsbereich`
- Überprüfen Sie die Dokumentation Ihres Clients für den genauen Konfigurationsdatei-Speicherort und das Format

Die meisten MCP-Clients unterstützen den stdio-Transport, und einige unterstützen HTTP (SSE / Streamable HTTP).
