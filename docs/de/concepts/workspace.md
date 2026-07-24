# Arbeitsbereich

Der Arbeitsbereich ist das Verzeichnis, in dem swag2mcp alle seine Daten speichert — Konfiguration, zwischengespeicherte Spezifikationen, lokale Spezifikationsdateien, gespeicherte Antworten und Auth-Skripte.

## Struktur

```
~/.swag2mcp/                          # Arbeitsbereichs-Stammverzeichnis (Standard)
├── swag2mcp.yaml                     # Konfigurationsdatei
├── cache/                            # Zwischengespeicherte entfernte Spezifikationsdateien
│   ├── a1b2c3d4e5f6...spec          # Zwischengespeicherter Spezifikationsinhalt
│   └── a1b2c3d4e5f6...meta          # Cache-Metadaten (JSON)
├── specs/                            # Lokale Spezifikationsdateien
│   └── my-api.yaml
├── responses/                        # Gespeicherte API-Antworten (große Antworten)
│   ├── meteo-get-forecast-abc123.json
│   └── response-fragment-def456.json
└── auth_scripts/                     # Authentifizierungsskripte
    ├── meteo.sh                      # Unix-Shell-Skript
    └── meteo.bat                     # Windows-Batch-Skript
```

## Standardpfad

- **Linux/macOS**: `~/.swag2mcp/`
- **Windows**: `%USERPROFILE%\.swag2mcp\`

## Benutzerdefinierter Pfad

```bash
swag2mcp mcp /pfad/zu/arbeitsbereich
swag2mcp mcp ./my-workspace
```

## Verzeichnisse

### cache/

Speichert heruntergeladene entfernte Spezifikationsdateien. Jede Datei wird mit einem SHA-256-Hash ihrer URL als Dateiname zwischengespeichert:

- `{hash}.spec` — der zwischengespeicherte Spezifikationsdateiinhalt
- `{hash}.meta` — JSON-Metadaten (Quell-URL, Cache-Zeit, TTL)

Jede zwischengespeicherte Datei hat eine zufällige TTL zwischen 1 Stunde und 48 Stunden. Der Cache wird bei jedem Start automatisch überprüft — wenn ein gültiger (nicht abgelaufener) Eintrag existiert, wird er ohne Download wiederverwendet.

**Befehle:**
- `swag2mcp update` — leert den Cache und lädt alle Spezifikationen neu herunter
- `swag2mcp clean` — leert Cache und Antworten

### specs/

Speichert lokale Spezifikationsdateien, auf die Collections über `location: specs/{name}` verweisen. Dateien hier werden direkt ohne Zwischenspeicherung verwendet.

Dieses Verzeichnis wird befüllt durch:
- `swag2mcp import <source> <name>` — lädt eine entfernte Spezifikation herunter und speichert sie hier
- `swag2mcp export` — kopiert Spezifikationen von hier in das Export-ZIP
- Manuelle Platzierung — Sie können Spezifikationsdateien selbst hierher kopieren

### responses/

Speichert API-Antworten, die das Limit `max_response_size` überschreiten (Standard 1 MB). Wenn der LLM einen Endpunkt aufruft und die Antwort zu groß ist, speichert swag2mcp sie hier und gibt stattdessen einen Dateiverweis zurück.

Namenskonvention: `{domain}-{method}-{path_with_underscores}-{6char_hex}.json`

Alte Antworten werden automatisch nach 48 Stunden beim MCP-Server-Start bereinigt.

### auth_scripts/

Speichert Authentifizierungsskripte für den `script`-Auth-Typ. Jedes Skript ist nach der Domain der Spec benannt.

#### Namenskonvention

| Plattform | Dateiname | Beispiel |
|-----------|-----------|----------|
| Unix (Linux, macOS) | `{domain}.sh` | `meteo.sh` |
| Windows | `{domain}.bat` | `meteo.bat` |

Die Domain darf keine `/`- oder `\`-Zeichen enthalten.

#### Wie Skripte funktionieren

1. swag2mcp führt das Skript mit einem 30-Sekunden-Timeout aus
2. Das Skript muss gültiges JSON an die Standardausgabe ausgeben
3. swag2mcp parst das JSON und verwendet das Token für API-Anfragen

#### Erwartetes Ausgabeformat

```json
{
  "token": "ihr-token-hier",
  "expires_in": 3600
}
```

| Feld | Typ | Erforderlich | Beschreibung |
|------|-----|-------------|--------------|
| `token` | string | ✅ | Das Authentifizierungstoken |
| `access_token` | string | ❌ | Alternative zu `token` (zuerst geprüft) |
| `token_type` | string | ❌ | Token-Typ (z. B. "Bearer") |
| `expires_in` | number | ❌ | Token-Lebensdauer in Sekunden (Standard: 3600) |

#### Ausführung

| Plattform | Befehl |
|-----------|--------|
| Unix | `sh {domain}.sh` |
| Windows | `cmd /c {domain}.bat` |

#### Token-Zwischenspeicherung

Das Token wird im Speicher zwischengespeichert, bis es abläuft. Bei jedem API-Aufruf überprüft swag2mcp zuerst den Cache — das Skript wird nur ausgeführt, wenn das zwischengespeicherte Token abgelaufen ist.

#### Stub-Erstellung

Wenn Sie `auth: { type: script, config: { domain: "myapi" } }` konfigurieren, erstellt swag2mcp automatisch ein Stub-Skript:

**Unix (`auth_scripts/myapi.sh`):**
```bash
#!/bin/sh
echo '{"token": "ihr-token-hier", "expires_in": 3600}'
```

**Windows (`auth_scripts/myapi.bat`):**
```bat
@echo off
echo {"token": "ihr-token-hier", "expires_in": 3600}
```

Ersetzen Sie das Platzhalter-Token durch Ihre tatsächliche Authentifizierungslogik.

#### Bereinigung verwaister Einträge

Wenn Sie eine Spec löschen, wird ihr Auth-Skript verwaist. swag2mcp entfernt verwaiste Skripte automatisch bei:
- `swag2mcp update`
- `swag2mcp clean`

## Befehle

### update

```bash
swag2mcp update [path]
```

Validiert die Konfiguration, leert Cache und Antworten und lädt dann alle Spezifikationsdateien neu herunter. Stellt auch sicher, dass Auth-Skripte existieren, und entfernt verwaiste Skripte.

Verwenden Sie diesen Befehl nach:
- Hinzufügen oder Entfernen von Collections
- Ändern von Collection-Speicherorten
- Bearbeiten von Spezifikationsdateien, die neu zwischengespeichert werden müssen

### clean

```bash
swag2mcp clean [path]
```

Entfernt alle Inhalte von `cache/` und `responses/`, plus verwaiste Auth-Skripte. Spezifikationen werden NICHT neu zwischengespeichert — verwenden Sie dafür `update`.

### validate

```bash
swag2mcp validate [path]
```

Validiert die Konfiguration einschließlich aller Collection-Speicherorte. Siehe [CLI: validate](../cli/validate.md).

## Export und Import

```bash
# Arbeitsbereich in ZIP exportieren (Standardname: swag2mcp-backup-{date}.zip)
swag2mcp export

# In einen bestimmten Pfad exportieren
swag2mcp export /pfad/zu/arbeitsbereich /pfad/zu/sicherung.zip

# Nur bestimmte Specs exportieren
swag2mcp export --spec meteo

# Aus Backup wiederherstellen
swag2mcp import --from-zip /pfad/zu/sicherung.zip
swag2mcp import /pfad/zu/arbeitsbereich /pfad/zu/sicherung.zip
```

Export enthält: `swag2mcp.yaml`, `specs/`, `auth_scripts/`. Cache und Antworten sind ausgeschlossen (es sind lokale Daten).

## .gitignore

Wenn sich Ihr Arbeitsbereich in einem Git-Repository befindet, fügen Sie diese Einträge zu `.gitignore` hinzu:

```gitignore
# swag2mcp — nur lokale Daten
.swag2mcp/cache/
.swag2mcp/responses/
```

Die Verzeichnisse `cache/` und `responses/` enthalten lokale, maschinenspezifische Daten, die nicht committet werden sollten. Alles andere (`swag2mcp.yaml`, `specs/`, `auth_scripts/`) sollte im Repository sein, damit die Konfiguration im Team geteilt wird.
