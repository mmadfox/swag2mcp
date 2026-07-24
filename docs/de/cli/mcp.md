# mcp

## Zweck

Startet den **MCP-Server (Model Context Protocol)** — den primären Modus für die LLM-Integration. Dies führen Sie aus, um einem KI-Agenten (Claude, Cursor, OpenCode usw.) über 16 MCP-Tools Zugriff auf Ihre APIs zu geben.

## Wann verwenden

- Sie möchten einen LLM-Agenten mit Ihren APIs verbinden
- Sie konfigurieren eine IDE (VS Code, Cursor, JetBrains) oder eine Desktop-App (Claude Desktop)
- Sie müssen Ihre APIs über das MCP-Protokoll bereitstellen
- Sie testen den MCP-Server vor der Integration

## Syntax

```bash
swag2mcp mcp [path] [flags]
```

## Argumente

| Argument | Position | Erforderlich | Beschreibung |
|----------|----------|-------------|--------------|
| `path` | 1 | Nein | Arbeitsbereichsverzeichnis. Wenn nicht angegeben, wird über die Pfadauflösungsregeln ermittelt. |

## Flags

| Flag | Kurzform | Typ | Standard | Beschreibung |
|------|----------|-----|----------|--------------|
| `--transport` | | `string` | `"stdio"` | MCP-Transport: `stdio`, `sse`, `streamable-http` |
| `--http-addr` | | `string` | `":8080"` | HTTP-Serveradresse (für `sse` und `streamable-http`) |
| `--http-path` | | `string` | `"/mcp"` | HTTP-Pfad für den MCP-Handler |
| `--auth-token` | | `string` | `""` | Bearer-Token für HTTP-Transport-Authentifizierung |
| `--logfile` | `-f` | `string` | `""` | Log-Dateipfad. Wenn nicht gesetzt, wird nach stderr geloggt. |
| `--disable-llm-auth` | | `bool` | `true` | Entfernt das `auth`-Tool aus der MCP-Tool-Liste |
| `--dump-dir` | | `string` | `""` | Verzeichnis zum Speichern von HTTP-Anfragen zum Debuggen |
| `--tags` | `-t` | `string` | `""` | Specs nach Tags filtern (kommagetrennt) |

## Wie es funktioniert

### stdio-Transport (Standard)

Wird verwendet, wenn der MCP-Server als Unterprozess vom LLM-Client (IDE, Claude Desktop usw.) gestartet wird. Der Server kommuniziert über Standardein-/ausgabe.

```bash
swag2mcp mcp
```

### SSE-Transport

Server-Sent-Events-Transport für HTTP-basierte Kommunikation. Erfordert die MCP-Handshake-Sequenz.

```bash
swag2mcp mcp --transport sse --http-addr :8080
```

### Streamable-HTTP-Transport

Moderner HTTP-Transport, der Streaming-Antworten unterstützt.

```bash
swag2mcp mcp --transport streamable-http --http-addr 0.0.0.0:8080
```

### Mit Authentifizierung

Schützen Sie den HTTP-Endpunkt mit einem Bearer-Token:

```bash
swag2mcp mcp --transport sse --http-addr :8080 --auth-token "my-secret"
```

### Mit Tag-Filterung

Nur Specs mit bestimmten Tags laden:

```bash
swag2mcp mcp --tags=public
```

### Mit aktiviertem Auth-Tool (Debug-Modus)

Dem LLM erlauben, frische Tokens über das `auth`-Tool anzufordern:

```bash
swag2mcp mcp --disable-llm-auth=false
```

### Mit Anfrage-Dump-Verzeichnis

Alle HTTP-Anfragen zum Debuggen speichern:

```bash
swag2mcp mcp --dump-dir ./dumps
```

## MCP-HTTP-Transport — Handshake-Protokoll

Bei Verwendung von `sse` oder `streamable-http` erfordert das MCP-Protokoll einen bestimmten Handshake. Tool-Aufrufe schlagen vor der Initialisierung fehl:

```
Schritt 1: POST /mcp → {"method":"initialize", ...}
Schritt 2: POST /mcp → {"method":"notifications/initialized"}
Schritt 3: POST /mcp → {"method":"tools/list", ...}   ← jetzt funktioniert es
```

### Health-Check

Funktioniert ohne Initialisierung:

```bash
curl http://localhost:8080/health
# → {"status":"ok","version":"v1.2.0"}
```

## IDE-Konfigurationsbeispiele

### VS Code (`.vscode/settings.json` oder globale Einstellungen)

```json
{
  "mcp": {
    "servers": {
      "swag2mcp": {
        "command": "swag2mcp",
        "args": ["mcp", "/absoluter/pfad/zu/.swag2mcp"]
      }
    }
  }
}
```

### Cursor / Windsurf (`~/.cursor/mcp.json` oder Projekt `.cursor/mcp.json`)

```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "/absoluter/pfad/zu/.swag2mcp"]
    }
  }
}
```

### Claude Desktop (`~/Library/Application Support/Claude/claude_desktop_config.json` unter macOS)

```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "/absoluter/pfad/zu/.swag2mcp"]
    }
  }
}
```

### JetBrains-IDE (Einstellungen → Tools → MCP)

- Name: `swag2mcp`
- Befehl: `swag2mcp`
- Argumente: `mcp /absoluter/pfad/zu/.swag2mcp`

> **Verwenden Sie immer einen absoluten Pfad** zum Arbeitsbereichsverzeichnis in der IDE-Konfiguration. Relative Pfade können je nach Arbeitsverzeichnis der IDE fehlschlagen.

## Ausgabe

Bei Erfolg gibt der Server aus:

```
MCP server listening on http://127.0.0.1:8080/mcp
```

## Nuancen

- **Kein Auto-Init:** Wenn die Konfigurationsdatei nicht existiert, gibt `mcp` einen Fehler zurück: `"Konfiguration nicht gefunden unter &lt;path&gt;"`. Führen Sie zuerst `init` aus.
- **`--disable-llm-auth` (Standard: `true`):** Wenn aktiviert, wird das `auth`-Tool vollständig aus der MCP-Tool-Liste entfernt. Der LLM kann Tokens weder sehen noch anfordern. Auth funktioniert weiterhin — Tokens werden über den Standard-Konfigurationsmechanismus bezogen, nicht über den LLM. Dieser Modus wird für die **Produktion** empfohlen. Für das **Debuggen** oder bei kurzlebigen Tokens setzen Sie `--disable-llm-auth=false`, damit der LLM frische Tokens über das `auth`-Tool anfordern kann.
- **YAML-Konfigurations-Fallback:** Wenn ein CLI-Flag nicht explizit gesetzt ist, wird der Wert aus dem `mcp`-Abschnitt in `swag2mcp.yaml` übernommen (falls vorhanden). Dies ermöglicht es Ihnen, den Server in der Konfigurationsdatei zu konfigurieren, anstatt jedes Mal Flags zu übergeben.
- **Antwortbereinigung:** Beim Start werden Antworten, die älter als 48 Stunden sind, automatisch aus dem Verzeichnis `responses/` entfernt.
- **Warnung zur Pfadauflösung:** Wenn `[path]` weggelassen wird, sucht `mcp` zuerst im aktuellen Verzeichnis nach `swag2mcp.yaml` und fällt dann auf `~/.swag2mcp/` zurück. Wenn Sie den Befehl aus dem falschen Verzeichnis ausführen, wird möglicherweise ein anderer Arbeitsbereich als beabsichtigt geladen. **Geben Sie `[path]` immer explizit an, wenn Sie es als Dienst oder in der IDE-Konfiguration ausführen.**
