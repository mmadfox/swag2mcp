# swag2mcp

**swag2mcp** ist ein CLI-Tool und MCP (Model Context Protocol) Server, der OpenAPI/Swagger/Postman API-Spezifikationen mit LLM-Agenten (Opencode, Crush, Copilot, Cursor, etc.) verbindet.

Es indexiert Ihre API-Spezifikationen in eine Volltext-Suchmaschine, stellt sie über 16 MCP-Tools bereit und ermöglicht LLMs, echte API-Endpunkte zu entdecken, zu inspizieren und aufzurufen — ohne eine einzige Zeile Integrationscode.

---

## Inhaltsverzeichnis

- [Schnellstart](#schnellstart)
- [Konfiguration](#konfiguration)
- [CLI-Befehle](#cli-befehle)
- [MCP Server](#mcp-server)
- [Suche](#suche)
- [Arbeitsverzeichnis (Workspace)](#arbeitsverzeichnis-workspace)
- [Caching](#caching)
- [Entwicklung](#entwicklung)

---

## Schnellstart

### Option 1 — Von GitHub Releases herunterladen (empfohlen)

1. Öffnen Sie https://github.com/mmadfox/swag2mcp/releases/latest
2. Finden Sie das Archiv für Ihr System:

   | OS | Architektur | Archiv |
   |----|-------------|--------|
   | Linux | x86_64 | `swag2mcp_<version>_linux_amd64.tar.gz` |
   | Linux | ARM64 | `swag2mcp_<version>_linux_arm64.tar.gz` |
   | macOS | Intel | `swag2mcp_<version>_darwin_amd64.tar.gz` |
   | macOS | Apple Silicon | `swag2mcp_<version>_darwin_arm64.tar.gz` |
   | Windows | x86_64 | `swag2mcp_<version>_windows_amd64.zip` |

3. Herunterladen und installieren:

   **Linux / macOS:**
   ```bash
   tar -xzf swag2mcp_<version>_<os>_<arch>.tar.gz
   sudo mv swag2mcp /usr/local/bin/
   swag2mcp --version
   ```

   **Windows (PowerShell):**
   ```powershell
   Expand-Archive swag2mcp_<version>_windows_amd64.zip -DestinationPath .
   move swag2mcp.exe C:\Windows\System32\
   swag2mcp --version
   ```

4. (Optional) Wiederholen für Mock-Server — `swag2mcp-mock_<version>_<os>_<arch>.tar.gz`

### Option 2 — Mit Go installieren

Wenn Go installiert ist:

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

### Nach der Installation

```bash
# Arbeitsverzeichnis initialisieren
swag2mcp init

# MCP Server starten (für LLM-Agenten)
swag2mcp mcp

# Oder interaktiven Explorer nutzen
swag2mcp run
```---

## Example LLM Queries

After setup, try asking your agent:

| Query | What happens |
|-------|-------------|
| "Show me all available APIs" | `spec_list` — lists petstore, binance, countries |
| "What endpoints does Binance have?" | `endpoint_by_spec` — shows 4 market data endpoints |
| "Find endpoints related to pets" | `search("pet")` — finds petstore endpoints |
| "What tags are in the Petstore API?" | `tag_by_spec` — shows "pets" tag |
| "Show me the GET /pets endpoint details" | `inspect` — shows parameters and response schema |
| "Get the current BTC price from Binance" | `invoke` — real API call to Binance |
| "Find countries in Europe" | `invoke` — calls REST Countries API |

---

---

## Konfiguration

### YAML Schema

```yaml
mock_enabled: true                    # optional, aktiviert den Mock-Server-Modus

http_client:                        # optional, globale HTTP-Standards
  headers:                          # optional
    X-API-Version: "2"
  cookies: []                       # optional
  user_agent: ""                    # optional
  timeout: 0s                       # optional
  follow_redirects: true            # optional
  max_redirects: 10                 # optional
  max_response_size: 1048           # optional, Bytes (Standard 1KB, max 1MB)

specs:
  - domain: petstore                    # erforderlich, 1-60 Zeichen, [a-zA-Z0-9_-]
    llm_title: Petstore API             # erforderlich, 5-120 Zeichen
    llm_instruction: |                  # optional, max 500 Zeichen
      Verwende diese API für Haustiere, Bestellungen und Benutzer.
    base_url: https://petstore.swagger.io/v2  # erforderlich, gültige URL
    disable: false                      # optional
    tags: [public, demo]                # optional, zum Filtern
    http_client:                        # optional, überschreibt global
      headers:
        X-API-Version: "2"
    auth:                               # optional
      type: bearer                      # siehe Authentifizierungsmethoden
      config:
        token: $(TOKEN_AUTH)
    collections:
      - llm_title: Petstore Swagger     # optional, max 120 Zeichen
        llm_instruction: |             # optional, max 360 Zeichen
          Hauptendpunkte von Petstore
        title: ""                      # optional, wird automatisch aus Spec befüllt
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/petstore.json  # erforderlich, 5-250 Zeichen
        disable: false                  # optional
        base_url: ""                    # optional, überschreibt base_url der Spec
        base_mock_url: localhost:8080   # optional, Format "host:port" oder "host:port/path"
        http_client: {}                 # optional, überschreibt Spec
```

### Tags — Spezifikationen nach Projekt filtern

Tags ermöglichen die Gruppierung von Spezifikationen nach Projekt, Umgebung oder Team. Beim Start des MCP Servers mit `--tags` werden nur passende Spezifikationen geladen:

```bash
# Server nur mit öffentlichen Spezifikationen starten
swag2mcp mcp --tags=public

# Server mit mehreren Tags starten
swag2mcp mcp --tags=public,internal

# Mehrere Server für verschiedene Projekte ausführen
swag2mcp mcp --tags=project-alpha --logfile=/tmp/swag2mcp-alpha.log
swag2mcp mcp --tags=project-beta  --logfile=/tmp/swag2mcp-beta.log
```

Dies ermöglicht den Betrieb separater MCP Server für verschiedene Projekte aus einer einzigen Konfigurationsdatei.

### Authentifizierungsmethoden

| Typ | Felder | Konfigurationsbeispiel |
|-----|--------|------------------------|
| `none` | — | `type: none` |
| `basic` | `username`, `password` | `username: $(USER)`, `password: $(PASS)` |
| `bearer` | `token` | `token: $(TOKEN)` |
| `digest` | `username`, `password` | `username: admin`, `password: secret` |
| `hmac` | `api_key`, `secret_key` | `api_key: $(API_KEY)`, `secret_key: $(SECRET_KEY)` |
| `api-key` | `key`, `value`, `in` (header/query) | `key: X-API-Key`, `value: $(KEY)`, `in: header` |
| `oauth2-cc` | `client_id`, `client_secret`, `token_url`, `scopes` | `client_id: $(ID)`, `token_url: https://auth.example.com/token` |
| `oauth2-pwd` | `username`, `password`, `client_id`, `client_secret`, `token_url`, `scopes` | `username: $(USER)`, `token_url: https://auth.example.com/token` |
| `script` | `source` | `source: pfad/zu/auth.sh` |

Alle Zeichenfolgenfelder unterstützen die Syntax `$(ENV_VAR)` — wird zur Laufzeit aus Umgebungsvariablen aufgelöst.

---

## CLI-Befehle

Alle Befehle, die `[path]` akzeptieren, verwenden dieselbe Pfadauflösung:

```
swag2mcp <command>          → ~/.swag2mcp/
swag2mcp <command> ./       → {cwd}/.swag2mcp/
swag2mcp <command> path/to  → {cwd}/path/to/.swag2mcp/
```

### `init [path]`

Arbeitsverzeichnis und Konfiguration initialisieren.

| Flag | Kurz | Standard | Beschreibung |
|------|------|----------|-------------|
| `--interactive` | `-i` | `false` | Interaktiven Assistenten starten |
| `--force` | `-f` | `false` | Vorhandene Konfiguration überschreiben |

```bash
swag2mcp init              # ~/.swag2mcp/swag2mcp.yaml erstellen
swag2mcp init ./           # ./.swag2mcp/swag2mcp.yaml erstellen
swag2mcp init -i           # interaktiver Assistent
```

### `add spec [path]` / `add collection [path]`

Eine Spezifikation oder Sammlung zur Konfiguration hinzufügen.

| Flag | Kurz | Standard | Beschreibung |
|------|------|----------|-------------|
| `--yaml` | `-y` | `""` | YAML-Eingabe (`-` für stdin) |
| `--example` | `-e` | `false` | YAML-Beispiel anzeigen |

```bash
swag2mcp add spec
swag2mcp add spec --yaml 'domain: petstore\nllm_title: Petstore API\nbase_url: https://...'
cat spec.yaml | swag2mcp add spec --yaml -
swag2mcp add spec --example
```

### `delete spec [path]` / `delete collection [path]`

Eine Spezifikation oder Sammlung löschen. Interaktive Auswahl.

```bash
swag2mcp delete spec
swag2mcp delete collection
```

### `ls [path]`

Spezifikationen und Sammlungen auflisten.

| Flag | Kurz | Standard | Beschreibung |
|------|------|----------|-------------|
| `--tags` | `-t` | `""` | Nach Tags filtern (kommagetrennt) |

```bash
swag2mcp ls
swag2mcp ls --tags=public,internal
```

### `run [path]`

Interaktiver API-Explorer (TUI). Endpunkte suchen, durchsuchen, inspizieren und speichern.

```bash
swag2mcp run
```

### `validate [path]`

Konfiguration validieren und Prüfung, ob alle Sammlungsorte erreichbar sind.

| Flag | Kurz | Standard | Beschreibung |
|------|------|----------|-------------|
| `--tags` | `-t` | `""` | Spezifikationen nach Tags filtern |

```bash
swag2mcp validate
swag2mcp validate --tags=public
```

### `clean [path]`

Alle Inhalte der Verzeichnisse `cache/` und `responses/` entfernen.

```bash
swag2mcp clean
```

### `update [path]`

Konfiguration validieren, Cache leeren, alle Spec-Dateien neu cachen.

```bash
swag2mcp update
```

### `mcp [path]`

MCP Server im Headless-Modus starten (stdio Transport). Dies ist der primäre Produktionsbefehl für die LLM-Integration.

| Flag | Kurz | Standard | Beschreibung |
|------|------|----------|-------------|
| `--logfile` | `-f` | `""` | Pfad zur Logdatei |
| `--tags` | `-t` | `""` | Spezifikationen nach Tags filtern |
| `--disable-llm-auth` | | `true` | `true` — Authentifizierung im Hintergrund (LLM sieht keine Tokens). `false` — LLM kann Tokens über das `auth`-Tool anfordern |
| `--dump-dir` | | `""` | Verzeichnis zum Dumpen von HTTP-Anfragen (Debugging) |

```bash
swag2mcp mcp
swag2mcp mcp --tags=public --logfile=/var/log/swag2mcp.log
swag2mcp mcp --disable-llm-auth=false
swag2mcp mcp --dump-dir=/tmp/dump
```

### `mockserver [path]`

Startet Mock-HTTP-Server für alle API-Spezifikationen. Jede Sammlung erhält einen eigenen
HTTP-Server, der Zufallsdaten generiert, die den OpenAPI-Antwortschemata entsprechen.

| Flag | Standard | Beschreibung |
|------|----------|-------------|
| `--tls` | `false` | TLS mit selbstsigniertem Zertifikat aktivieren |
| `--tls-cert` | `""` | Pfad zur TLS-Zertifikatsdatei |
| `--tls-key` | `""` | Pfad zur TLS-Schlüsseldatei |

```bash
swag2mcp-mock
swag2mcp-mock --tls
```

**Workflow:**
1. Fügen Sie `mock_enabled: true` und `base_mock_url` zu Ihrer Konfiguration hinzu
2. Starten Sie den Mock-Server: `swag2mcp-mock`
3. Starten Sie den MCP-Server: `swag2mcp mcp` — invoke verwendet `base_mock_url` anstelle von `base_url`
4. Authentifizierung erfolgt automatisch: OAuth2/Digest nutzen Mock-Server auf Ports 9090/9091; andere Typen wenden Anmeldedaten direkt an

### Mock-Authentifizierung

Wenn `auth` in einer Spezifikation konfiguriert ist, wendet der MCP-Server
die Authentifizierung automatisch an. Nur zwei Auth-Typen benötigen einen
dedizierten Mock-Server:

| Auth-Typ | Mock-Endpunkt | Verhalten |
|----------|---------------|-----------|
| `oauth2-cc` / `oauth2-pwd` | `POST /token` auf Port 9090 | Akzeptiert beliebige `client_id`/`username`+`password`, gibt `{"access_token":"<random>","token_type":"Bearer","expires_in":3600}` zurück |
| `digest` | `GET /` auf Port 9091 | Sendet eine 401-Challenge mit `algorithm=MD5`, akzeptiert jede Digest-Antwort, gibt `{"status":"authenticated","method":"digest"}` zurück |

Andere Auth-Typen (`basic`, `bearer`, `api-key`, `hmac`, `script`) benötigen **keinen**
Mock-Server — der MCP-Server wendet die konfigurierten Anmeldedaten automatisch
auf jede Anfrage an.

---

## Integration

swag2mcp spricht das Model Context Protocol (MCP) und funktioniert mit jedem MCP-kompatiblen Client.

### Lokal (stdio) — Agent auf demselben Rechner

Server starten:

```bash
swag2mcp mcp
```

| Client | Konfigurationsdatei | Inhalt |
|--------|-------------------|--------|
| **OpenCode** | `opencode.json` | `{"mcp":{"swag2mcp":{"type":"local","command":["swag2mcp","mcp"]}}}` |
| **Cursor** | `.cursor/mcp.json` | `{"mcpServers":{"swag2mcp":{"command":"swag2mcp","args":["mcp"]}}}` |
| **Claude Desktop** | `claude_desktop_config.json` | `{"mcpServers":{"swag2mcp":{"command":"swag2mcp","args":["mcp"]}}}` |
| **VS Code** | `.vscode/mcp.json` | `{"servers":{"swag2mcp":{"type":"stdio","command":"swag2mcp","args":["mcp"]}}}` |
| **Crush** | `crush.json` | `{"mcp":{"swag2mcp":{"type":"stdio","command":"swag2mcp","args":["mcp"]}}}` |

### Remote (HTTP) — Agent in der Cloud / anderem Rechner

Server mit HTTP-Transport starten:

```bash
swag2mcp mcp --transport streamable-http --http-addr :8080 --auth-token my-secret
```

Oder in `swag2mcp.yaml` konfigurieren:

```yaml
mcp:
  transport: streamable-http
  addr: ":8080"
  path: "/mcp"
  auth_token: $(MCP_AUTH_TOKEN)
```

| Client | Konfigurationsdatei | Inhalt |
|--------|-------------------|--------|
| **OpenCode** | `opencode.json` | `{"mcp":{"swag2mcp":{"type":"remote","url":"http://localhost:8080/mcp","headers":{"Authorization":"Bearer ${MCP_AUTH_TOKEN}"}}}}` |
| **VS Code** | `.vscode/mcp.json` | `{"servers":{"swag2mcp":{"type":"http","url":"http://localhost:8080/mcp"}}}` |

> **Health-Check** (funktioniert ohne MCP-Handshake):
> ```bash
> curl http://localhost:8080/health
> # → {"status":"ok","version":"v1.1.3"}
> ```

---

## MCP Server

Der MCP Server stellt 16 Tools über stdio oder HTTP Transport bereit. LLM-Agenten (Opencode, Cursor, Claude, Copilot, Crush, etc.) verbinden sich automatisch nach der Konfiguration.

### Tool-Hierarchie

```
spec_list                       — alle verfügbaren Spezifikationen auflisten
  └─ spec_by_id                 — Spezifikationsdetails per ID
       └─ collection_by_spec    — Sammlungen in einer Spezifikation
            └─ tag_by_collection     — Tags in einer Sammlung
                 └─ endpoint_by_tag  — Endpunkte unter einem Tag
                      └─ inspect          — vollständige OpenAPI-Operation
                           └─ invoke       — API-Aufruf ausführen

search                          — Volltextsuche über alle Endpunkte
```

### Tool-Referenz

| Tool | Argumente | Rückgabe | Beschreibung |
|------|-----------|----------|-------------|
| `spec_list` | — | `Spec[]` | Alle verfügbaren Spezifikationen |
| `spec_by_id` | `id` | Spec + Collections | Spezifikationsdetails |
| `collection_by_spec` | `specId` | Collections | Sammlungen in einer Spezifikation |
| `collection_by_id` | `id` | Collection + Tags | Sammlungsdetails |
| `tag_by_collection` | `collectionId` | Tags | Tags in einer Sammlung |
| `tag_by_spec` | `specId` | Tags | Alle Tags einer Spezifikation |
| `tag_by_id` | `id` | Tag | Einzelne Tag-Metadaten |
| `endpoint_by_tag` | `tagId` | Endpoints | Endpunkte unter einem Tag |
| `endpoint_by_collection` | `collectionId` | Endpoints | Alle Endpunkte einer Sammlung |
| `endpoint_by_spec` | `specId` | Endpoints | Alle Endpunkte einer Spezifikation |
| `endpoint_by_id` | `id` | Endpoint | Kurze Endpunktzusammenfassung |
| `search` | `query`, `limit` | Endpoints | Volltextsuche |
| `inspect` | `endpointId` | Full Operation | Vollständiges OpenAPI-Operationsobjekt |
| `invoke` | `endpointId`, `parameters`, `requestBody` | Response | Führt echten API-Aufruf aus |
| `auth` | `specId` | Token | Auth-Token für eine Spezifikation abrufen |

---

## Suche

### Abfragesyntax

| Funktion | Syntax | Beispiel |
|----------|--------|---------|
| Begriff | `begriff` | `haustiere` |
| Phrase | `"phrase"` | `"neues haustier hinzufügen"` |
| Feld: method | `method:begriff` | `method:post` |
| Feld: tag | `tag:begriff` | `tag:auth` |
| Feld: path | `path:begriff` | `path:/users` |
| Feld: summary | `summary:begriff` | `summary:login` |
| Erforderlich (AND) | `+begriff` | `+method:post +tag:user` |
| Ausgeschlossen (NOT) | `-begriff` | `-deprecated` |
| Wildcard | `*` | `path:*/v2/*` |
| Unscharf | `begriff~` | `watex~` |
| Regex | `/muster/` | `/user(s\|sessions)/` |
| Gewichtung | `begriff^N` | `tag:pet^5` |
| Alle treffen | `*` | `*` |

### Beispiele

```
# POST-Endpunkte im auth-Tag finden
+method:post +tag:auth

# Nach login-bezogenen Endpunkten suchen
summary:"login"~

# Alle benutzerbezogenen Pfade finden, veraltete ausschließen
path:*/users/* -deprecated

# Komplexe Abfrage
+method:get +tag:pet summary:"find by status"
```

### Indizierte Felder

| Feld | Typ | Inhalt |
|------|-----|--------|
| `method` | text | HTTP-Methode (kleingeschrieben) |
| `tag` | text | Tag-Name (kleingeschrieben) |
| `path` | text | API-Pfad (kleingeschrieben) |
| `summary` | text (analysiert) | Endpunkt-Zusammenfassung/Beschreibung (kleingeschrieben) |
| `_all` | text (analysiert) | method + path + tag + summary |

---

## Arbeitsverzeichnis (Workspace)

### Verzeichnisstruktur

```
~/.swag2mcp/                    # oder {projekt}/.swag2mcp/
├── swag2mcp.yaml               # Konfigurationsdatei
├── cache/                      # Zwischengespeicherte entfernte Spezifikationen
│   ├── {hash}.spec             # Spec-Dateiinhalt
│   └── {hash}.meta             # JSON-Metadaten
├── specs/                      # Lokale Spec-Dateien (benutzerverwaltet)
├── responses/                  # Antwortdateien von Aufrufen
└── auth_scripts/               # Authentifizierungsskripte
```

### Pfadauflösung

```
swag2mcp <command>          → ~/.swag2mcp/
swag2mcp <command> ./       → {cwd}/.swag2mcp/
swag2mcp <command> path/to  → {cwd}/path/to/.swag2mcp/
```

### .gitignore

Nur temporäre Daten sollten ignoriert werden:

```
.swag2mcp/cache/*
.swag2mcp/responses/*
```

Die Konfiguration `.swag2mcp/swag2mcp.yaml` und Spec-Dateien in `.swag2mcp/specs/` **müssen im Repository sein**.

### Empfehlung

Alle Spec-Dateien in `.swag2mcp/specs/` aufbewahren — nur so wird sichergestellt, dass sie direkt genutzt und nicht in den Cache kopiert werden.

---

## Caching

### Regeln

| Quelle | Verhalten |
|--------|-----------|
| HTTP/HTTPS URL | Wird immer gecached. TTL: zufällig 1-48h. |
| Lokaler Pfad innerhalb `specs/` | Direkt genutzt, nicht gecached. |
| Lokaler Pfad außerhalb `specs/` | Beim ersten Zugriff in Cache kopiert. |
| `file://` URL | Wird wie lokaler Pfad behandelt. |

### Cache-Schlüssel

SHA-256 Hash des normalisierten Speicherorts (erste 16 Bytes = 32 Hex-Zeichen).

### Cache-Treffer-Logik

1. `.meta`-Datei lesen — abgelaufen oder fehlt → Fehltreffer
2. Bei lokalen Quellen: `ModTime` geändert → Fehltreffer
3. `.spec`-Datei fehlt → Fehltreffer
4. Sonst → Treffer

---

## Entwicklung

```bash
# Bauen
go build ./cmd/swag2mcp/

# Tests
go test ./...

# Linter
make lint

# Ausführen
go run ./cmd/swag2mcp/main.go
```
