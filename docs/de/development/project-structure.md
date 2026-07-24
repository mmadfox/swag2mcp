# Projektstruktur

```
swag2mcp/
├── cmd/
│   ├── swag2mcp/          # Hauptbinärdatei
│   │   └── main.go
│   └── swag2mcp-mock/     # Mock-Server
│       └── main.go
├── internal/
│   ├── auth/              # 9 Auth-Methoden
│   ├── cache/             # Spec-Zwischenspeicherung
│   ├── commands/          # 13 CLI-Befehle (cobra)
│   ├── config/            # YAML-Konfiguration
│   ├── env/               # Umgebungsvariablen
│   ├── httpclient/        # HTTP-Client
│   ├── id/                # MD5-ID-Generierung
│   ├── index/             # Volltextsuche (bluge)
│   ├── model/             # Datenmodelle
│   ├── reader/            # Lesen großer Antworten
│   ├── server/
│   │   ├── mcp/           # MCP-Server (19 Tools)
│   │   └── mockserver/    # Mock-Server
│   ├── service/           # Geschäftslogik
│   ├── spec/              # Spec-Parser
│   ├── tui/               # TUI-Oberfläche
│   └── workspace/         # Arbeitsbereichsverwaltung
├── specs/                 # Beispiel-Spezifikationen
├── tests/                 # Integrationstests
├── docs/                  # Dokumentation
├── examples/              # Konfigurationsbeispiele
└── playground/            # Entwicklungssandbox
```

## Wichtige Pakete

| Paket | Beschreibung |
|-------|--------------|
| `auth` | 9 Authentifizierungsmethoden |
| `cache` | Festplattenbasierte Zwischenspeicherung mit TTL |
| `commands` | Cobra-CLI-Befehle |
| `config` | YAML-Konfiguration mit Kaskade |
| `httpclient` | Konfigurierbarer HTTP-Client |
| `index` | Volltextsuche (bluge) |
| `server/mcp` | MCP-Server (3 Transports) |
| `service` | Geschäftslogik (Kern) |
| `spec` | OpenAPI/Swagger/Postman-Parser |
| `tui` | Bubbletea-TUI |
| `workspace` | Dateiverwaltung |
