# FAQ

## Allgemein

### Was ist swag2mcp und welches Problem löst es?

swag2mcp verbindet OpenAPI/Swagger/Postman-API-Spezifikationen mit LLM-Agenten über das Model Context Protocol (MCP). Anstatt benutzerdefinierten Code zu schreiben, um jede API mit einem KI-Agenten zu verbinden, konfigurieren Sie es einmal in einer YAML-Datei, und der LLM erhält 19 Tools, um Ihre APIs zu entdecken, zu inspizieren und aufzurufen.

### Wie unterscheidet es sich von anderen API-zu-LLM-Tools?

- **Keine Programmierung erforderlich** — APIs in YAML konfigurieren, kein Integrationscode nötig
- **19 MCP-Tools** — komplettes Toolkit von der Erkennung bis zum Aufruf und zur Verarbeitung großer Antworten
- **9 Authentifizierungsmethoden** — funktioniert mit jedem API-Authentifizierungsschema
- **Volltextsuche** — bluge-gestützte Suche über alle Endpunkte
- **TUI-Explorer** — interaktive Terminal-Oberfläche zum Durchsuchen und Testen
- **Mock-Server** — Testen ohne echte API-Aufrufe

### Welche API-Spezifikationsformate werden unterstützt?

OpenAPI 3.x, Swagger 2.0 und Postman Collections v2.1.

### Was ist der Unterschied zwischen einer Spec und einer Collection?

Eine **Spec** repräsentiert einen logischen API-Dienst (z. B. "Open-Meteo Weather APIs"). Eine **Collection** ist eine einzelne OpenAPI/Swagger/Postman-Datei. Eine Spec kann mehrere Collections haben — zum Beispiel, wenn eine API separate Spezifikationsdateien für verschiedene Dienste (Vorhersage, Luftqualität, Meer) hat.

### Welche MCP-Transports werden unterstützt?

Drei Transports: `stdio` (Standard, für lokale LLM-Clients), `sse` (Server-Sent Events für entfernte Clients) und `streamable-http` (modernes HTTP-Streaming).

### Kann ich swag2mcp mit jedem LLM verwenden?

Ja, mit jedem LLM-Client, der das MCP-Protokoll unterstützt: Claude Desktop, VS Code, Cursor, Windsurf, JetBrains-IDE, OpenCode und andere.

## Installation

### Wie installiere ich swag2mcp?

```bash
# Option 1: Von GitHub Releases herunterladen
# Gehen Sie zu https://github.com/mmadfox/swag2mcp/releases/latest
# Laden Sie das Archiv für Ihr Betriebssystem und Ihre Architektur herunter

# Option 2: Mit Go installieren
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

### Muss Go installiert sein?

Nein. Vorgefertigte Binärdateien sind für Linux (amd64, arm64), macOS (amd64, arm64) und Windows (amd64) auf der [GitHub Releases-Seite](https://github.com/mmadfox/swag2mcp/releases) verfügbar.

### Wie installiere ich den Mock-Server?

Der Mock-Server ist eine separate Binärdatei:

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

Oder laden Sie `swag2mcp-mock_<version>_<os>_<arch>.tar.gz` von GitHub Releases herunter.

## Erste Schritte

### Wie komme ich schnell los?

```bash
# 1. Arbeitsbereich initialisieren
mkdir -p .swag2mcp && swag2mcp init ./.swag2mcp

# 2. MCP-Server starten (öffentliche Beispiel-Spezifikationen sind nach init enthalten)
swag2mcp mcp
```

Nach `init` enthält der Arbeitsbereich bereits mehrere öffentliche Beispiel-Spezifikationen (icanhazdadjoke, Open-Meteo, Binance, PokéAPI). Sie können den MCP-Server sofort starten — es ist nicht nötig, manuell Spezifikationen hinzuzufügen.

Wenn Sie stattdessen Ihre eigene API hinzufügen möchten:

```bash
swag2mcp add spec --yaml - <<EOF
domain: dadjoke
llm_title: icanhazdadjoke API
base_url: https://icanhazdadjoke.com
collections:
  - llm_title: Jokes
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
EOF
```

### Wie verbinde ich swag2mcp mit meiner IDE?

**VS Code** (`.vscode/settings.json`):
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

**Cursor** (`~/.cursor/mcp.json`):
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

**Claude Desktop** (`claude_desktop_config.json`):
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

Verwenden Sie immer einen absoluten Pfad zum Arbeitsbereichsverzeichnis.

## Konfiguration

### Wo befindet sich die Konfigurationsdatei?

Standard: `~/.swag2mcp/swag2mcp.yaml`. Sie können sie auch in einem beliebigen Verzeichnis erstellen und den Pfad an Befehle übergeben.

### Wie füge ich eine API hinzu?

```bash
# Interaktiver Modus
swag2mcp add spec

# Mit YAML (für Skripte empfohlen)
swag2mcp add spec --yaml - <<EOF
domain: my-api
llm_title: My API
base_url: https://api.example.com/v1
collections:
  - llm_title: Main
    location: https://example.com/spec.yaml
EOF
```

### Wie füge ich eine Collection zu einer bestehenden Spec hinzu?

```bash
swag2mcp add collection --yaml - <<EOF
spec_domain: meteo
llm_title: Air Quality
location: https://example.com/air-quality.yaml
EOF
```

### Wie deaktiviere ich eine Spec vorübergehend?

Setzen Sie `disable: true` in der Spec-Konfiguration. Die Spec wird nicht geladen oder indiziert.

### Kann ich filtern, welche Specs geladen werden?

Ja, verwenden Sie das Flag `--tags`: `swag2mcp mcp --tags=public`. Es werden nur Specs mit passenden Tags geladen.

### Wie verwende ich Umgebungsvariablen für Geheimnisse?

Verwenden Sie die Syntax `$(VAR_NAME)` in Authentifizierungsfeldern:

```yaml
auth:
  type: bearer
  config:
    token: "$(MY_API_TOKEN)"
```

Setzen Sie die Variable vor dem Start: `export MY_API_TOKEN="eyJhbGci..."`

## Authentifizierung

### Welche Authentifizierungsmethoden werden unterstützt?

Neun Methoden: `none`, `basic`, `bearer`, `digest`, `hmac`, `oauth2-cc` (Client Credentials), `oauth2-pwd` (Password Grant), `api-key` und `script`.

### Wie übergebe ich ein Token?

Über die Konfigurationsdatei oder Umgebungsvariablen:

```yaml
auth:
  type: bearer
  config:
    token: "$(MY_TOKEN)"
```

### Muss ich vor dem Aufruf die Authentifizierung durchführen?

Nein. Das `invoke`-Tool wendet die Authentifizierung automatisch aus der Konfiguration der Spec an. Sie benötigen das `auth`-MCP-Tool nur, wenn Sie das Token dem Benutzer anzeigen möchten (z. B. für einen curl-Befehl).

### Warum wird das Auth-Tool nicht angezeigt?

Das `auth`-Tool ist standardmäßig deaktiviert (`--disable-llm-auth=true`). Dies ist eine Sicherheitsmaßnahme für die Produktion. Um es zu aktivieren: `swag2mcp mcp --disable-llm-auth=false`.

### Wie werden OAuth2-Tokens erneuert?

OAuth2-Client-Credentials- und Password-Grant-Tokens werden automatisch erneuert, wenn sie ablaufen. Bearer-Tokens sind statisch und müssen manuell aktualisiert werden.

## MCP-Server

### Wie starte ich den MCP-Server?

```bash
# Standard (stdio-Transport)
swag2mcp mcp

# Mit HTTP-Transport
swag2mcp mcp --transport sse --http-addr :8080
```

### Wie ändere ich den Port?

```bash
swag2mcp mcp --transport sse --http-addr 0.0.0.0:9090
```

### Wie schütze ich den MCP-HTTP-Endpunkt?

Setzen Sie ein Bearer-Token:

```bash
swag2mcp mcp --transport sse --http-addr :8080 --auth-token "my-secret"
```

Der LLM-Client muss bei jeder Anfrage `Authorization: Bearer my-secret` mitsenden.

### Was ist der MCP-Handshake für den HTTP-Transport?

Für SSE- und Streamable-HTTP-Transports erfordert das MCP-Protokoll einen dreistufigen Handshake:

```
Schritt 1: POST /mcp → {"method":"initialize", ...}
Schritt 2: POST /mcp → {"method":"notifications/initialized"}
Schritt 3: POST /mcp → {"method":"tools/list", ...}  ← jetzt funktioniert es
```

Tool-Aufrufe schlagen vor der Initialisierung fehl.

## Verwendung

### Wie suche ich nach Endpunkten?

Verwenden Sie das `search`-MCP-Tool oder die TUI (`swag2mcp run`). Die Suche unterstützt Feldfilter (`method:GET`, `tag:pets`), unscharfe Suche, Platzhalter und boolesche Operatoren.

### Wie rufe ich eine API auf?

Der LLM verwendet das `invoke`-MCP-Tool. Inspizieren Sie immer zuerst den Endpunkt, um die erforderlichen Parameter zu verstehen:

```
inspect(endpointId: "...")  → Vertrag verstehen
invoke(endpointId: "...", parameters: {...})  → Aufruf durchführen
```

### Was passiert, wenn eine Antwort zu groß ist?

Antworten, die `max_response_size` überschreiten (Standard 1 MB), werden auf die Festplatte gespeichert. Der LLM erhält einen Dateiverweis und kann ihn mit den Tools `response_outline`, `response_compress` und `response_slice` erkunden.

### Wie funktioniert der Ratenbegrenzer?

Jeder Endpunkt hat eine 10-Sekunden-Abklingzeit. Wenn der LLM denselben Endpunkt zweimal innerhalb von 10 Sekunden aufruft, wird der zweite Aufruf stillschweigend blockiert. Sie können dies in der Konfiguration deaktivieren oder anpassen.

### Kann ich testen, ohne echte API-Aufrufe zu tätigen?

Ja, verwenden Sie den Mock-Server:

```bash
swag2mcp-mock mockserver
```

Er generiert gefälschte Antworten basierend auf OpenAPI-Schemata.

## Arbeitsbereichsverwaltung

### Wie sichere ich meine Konfiguration?

```bash
swag2mcp export --output ~/backups/swag2mcp-2026-07-24.zip
```

### Wie übertrage ich auf einen anderen Rechner?

```bash
# Auf dem alten Rechner
swag2mcp export --output swag2mcp.zip

# ZIP kopieren, dann auf dem neuen Rechner
swag2mcp import --from-zip swag2mcp.zip
```

### Wie aktualisiere ich Spezifikationsdateien?

```bash
swag2mcp update
```

Dies validiert die Konfiguration erneut, leert den Cache und lädt alle Spezifikationsdateien neu herunter.

### Wie räume ich Speicherplatz auf?

```bash
swag2mcp clean
```

Entfernt zwischengespeicherte Spezifikationsdateien und gespeicherte API-Antworten. Alte Antworten (>48h) werden auch automatisch beim Start des MCP-Servers bereinigt.

## TUI

### Was ist die TUI und wie verwende ich sie?

Die TUI (Terminal User Interface) ist ein interaktiver API-Explorer. Starten Sie sie mit `swag2mcp run`. Sie hat drei Modi: Suche (Volltextsuche), Durchsuchen (Baumnavigation: Spec → Collection → Tag → Endpunkt) und Auth (Tokens anzeigen).

### Was sind die Tastenkürzel?

| Taste | Aktion |
|-------|--------|
| `↑/↓` | Navigieren |
| `Enter` | Auswählen |
| `Esc` | Zurück |
| `Tab` | Modi wechseln |
| `/` | Suchen |
| `N/P` | Nächste/vorherige Seite |
| `q` | Beenden |

## Erweitert

### Kann ich einen Proxy verwenden?

Ja, konfigurieren Sie ihn in `http_client.proxy`:

```yaml
http_client:
  proxy:
    url: "http://proxy.company.com:8080"
    username: "$(PROXY_USER)"
    password: "$(PROXY_PASS)"
    bypass:
      - "localhost"
      - "*.internal.com"
```

### Kann ich eine benutzerdefinierte Authentifizierungsmethode hinzufügen?

Ja, implementieren Sie das `Authenticator`-Interface in `internal/auth/` und registrieren Sie es im Konfigurationsparser. Siehe den Entwicklungsabschnitt für Details.

### Kann ich ein benutzerdefiniertes MCP-Tool hinzufügen?

Ja, fügen Sie eine Methode zum `Svc`-Interface hinzu, implementieren Sie sie in der Service-Schicht, fügen Sie einen Handler hinzu und registrieren Sie ihn. Siehe den Entwicklungsabschnitt für Details.

### Was ist der Unterschied zwischen `swag2mcp` und `swag2mcp-mock`?

`swag2mcp` ist die Hauptbinärdatei mit CLI-Befehlen und dem MCP-Server. `swag2mcp-mock` ist eine separate Binärdatei, die Mock-Server zum Testen ohne echte API-Aufrufe startet.
