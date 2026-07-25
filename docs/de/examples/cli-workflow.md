# CLI-Workflow

Diese Seite zeigt echte Beispiele für die Verwendung von swag2mcp vom Terminal — von der Initialisierung bis zu täglichen Operationen.

## Schnellstart

```bash
# 1. Arbeitsbereich initialisieren
mkdir -p .swag2mcp && swag2mcp init ./.swag2mcp

# 2. Specs auflisten
swag2mcp ls
```

## Hinzufügen einer Spec mit YAML

### Einfache Spec (öffentliche API)

```bash
swag2mcp add spec --yaml - <<EOF
domain: meteo
llm_title: Open-Meteo Weather API
base_url: https://api.open-meteo.com
collections:
  - llm_title: Weather Forecast
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
EOF
```

### Spec mit Auth (Bearer-Token aus Umgebungsvariable)

```bash
swag2mcp add spec --yaml - <<EOF
domain: my-api
llm_title: My Protected API
base_url: https://api.example.com/v1
auth:
  type: bearer
  config:
    token: \$(MY_TOKEN)
collections:
  - llm_title: Users
    location: https://raw.githubusercontent.com/my-org/my-api/main/users.yaml
EOF
```

### Spec mit mehreren Collections

```bash
swag2mcp add spec --yaml - <<EOF
domain: meteo
llm_title: Open-Meteo APIs
base_url: https://api.open-meteo.com
collections:
  - llm_title: Forecast
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
  - llm_title: Air Quality
    base_url: https://air-quality-api.open-meteo.com
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
EOF
```

## Hinzufügen einer Collection zu einer bestehenden Spec

```bash
swag2mcp add collection --yaml - <<EOF
spec_domain: meteo
llm_title: Marine Weather
location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
EOF
```

## Specs auflisten

```bash
$ swag2mcp ls
Specifications:
  dadjoke (https://icanhazdadjoke.com)
    jokes (3 Endpunkte)
  meteo (https://api.open-meteo.com)
    forecast (5 Endpunkte)
    air-quality (8 Endpunkte)
    marine (4 Endpunkte)
```

### Nach Tags filtern

```bash
swag2mcp ls --tags=public
```

## Laufzeitinfo anzeigen

```bash
$ swag2mcp info
{
  "version": "v1.2.0",
  "workspace": "/home/user/.swag2mcp",
  "specs": {
    "total": 2,
    "active": 2,
    "disabled": 0,
    "collections": 4,
    "endpoints": 20
  },
  "http_client": {
    "timeout": "30s",
    "max_response_size": "2 KB",
    "follow_redirects": true
  },
  "mcp": {
    "transport": "stdio"
  },
  "auth": {
    "methods": ["bearer"]
  }
}
```

## Konfiguration validieren

```bash
$ swag2mcp validate
✅ Konfiguration ist gültig.
✓ Spec dadjoke: OK
✓ Spec meteo: OK
```

## MCP-Server starten

### stdio (für IDE-Integration)

```bash
swag2mcp mcp
```

### HTTP (für Remote-Zugriff)

```bash
swag2mcp mcp --transport sse --http-addr :8080
```

### Mit Tag-Filter

```bash
swag2mcp mcp --tags=public
```

## Specs aktualisieren

Alle zwischengespeicherten Spezifikationsdateien aktualisieren:

```bash
swag2mcp update
```

## Cache bereinigen

```bash
swag2mcp clean
```

## Export und Import

### Arbeitsbereich sichern

```bash
swag2mcp export --output ~/backups/swag2mcp-2026-07-24.zip
```

### Auf einem anderen Rechner wiederherstellen

```bash
# Auf dem neuen Rechner
swag2mcp import --from-zip swag2mcp-2026-07-24.zip
```

## Interaktiver TUI-Explorer

```bash
swag2mcp run
```

Öffnet eine Vollbild-Terminal-Oberfläche zum Suchen, Durchsuchen und Aufrufen von APIs.

## Mock-Server

```bash
# Mock-Binärdatei installieren
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest

# Mock-Server starten
swag2mcp-mock mockserver
```
