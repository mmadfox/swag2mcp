# info

## Zweck

Zeigt eine umfassende Zusammenfassung der swag2mcp-Laufzeit als **JSON** an. Dies umfasst Version, Arbeitsbereichspfad, Spec-Zusammenfassung, HTTP-Client-Einstellungen, MCP-Transportkonfiguration, Authentifizierungsmethoden und Mock-Modus-Status.

## Wann verwenden

- Sie möchten eine maschinenlesbare Übersicht über den Arbeitsbereich
- Sie müssen die Laufzeitkonfiguration zum Debuggen überprüfen
- Sie möchten sehen, wie viele Specs und Endpunkte aktiv sind
- Sie müssen HTTP-Client- oder MCP-Transporteinstellungen überprüfen

## Syntax

```bash
swag2mcp info [path]
```

## Argumente

| Argument | Position | Erforderlich | Beschreibung |
|----------|----------|-------------|--------------|
| `path` | 1 | Nein | Arbeitsbereichsverzeichnis. Wenn nicht angegeben, wird über die Pfadauflösungsregeln ermittelt. |

## Flags

Keine.

## Wie es funktioniert

```bash
swag2mcp info
swag2mcp info ./my-workspace
```

## Ausgabe

Die Ausgabe ist ein JSON-Objekt mit der folgenden Struktur:

```json
{
  "version": "v1.2.0",
  "workspace": "~/.swag2mcp",
  "uptime": "2h 15m",
  "specs": {
    "total": 4,
    "active": 3,
    "disabled": 1,
    "collections": 6,
    "endpoints": 42
  },
  "http_client": {
    "timeout": "30s",
    "max_response_size": "2 KB",
    "proxy": "none",
    "follow_redirects": true,
    "max_redirects": 10,
    "randomize": false
  },
  "mcp": {
    "transport": "stdio",
    "addr": ":8080",
    "path": "/mcp"
  },
  "auth_methods": ["bearer", "api-key"],
  "mock_enabled": false
}
```

## Überprüfung nach dem Befehl

Verwenden Sie `info`, um zu bestätigen, dass der Arbeitsbereich korrekt geladen wurde und alle Specs aktiv sind, bevor Sie den MCP-Server starten.

## Nuancen

- **Auto-Init:** Wenn keine Konfigurationsdatei existiert, führt `info` automatisch zuerst den Init-Assistenten aus.
- **Nur JSON:** Die Ausgabe ist immer JSON. Für menschenlesbare Ausgabe verwenden Sie `ls`.
- **`max_response_size`:** Wird in menschenlesbarem Format angezeigt (z. B. `"1 KB"`, `"2 MB"`).
- **Kein Volltextindex:** `info` deaktiviert die Volltextindizierung, da es nur Konfigurations- und Spec-Metadaten benötigt.
