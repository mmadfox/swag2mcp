# init

## Zweck

Der Befehl `init` erstellt einen **Arbeitsbereich** — ein Verzeichnis mit einer `swag2mcp.yaml`-Konfigurationsdatei und Unterverzeichnissen für Cache, Spezifikationen, Antworten und Auth-Skripte. Dies ist der erste Befehl, der bei der Einrichtung von swag2mcp ausgeführt wird.

## Wann verwenden

- Sie richten swag2mcp zum ersten Mal ein
- Sie möchten einen neuen Arbeitsbereich in einem bestimmten Verzeichnis erstellen
- Sie müssen einen beschädigten oder fehlenden Arbeitsbereich neu initialisieren

## Syntax

```bash
swag2mcp init [path] [flags]
```

## Argumente

| Argument | Position | Erforderlich | Beschreibung |
|----------|----------|-------------|--------------|
| `path` | 1 | Nein | Arbeitsbereichsverzeichnis. Wenn nicht angegeben, Standard: `~/.swag2mcp`. |

## Flags

| Flag | Kurzform | Typ | Standard | Beschreibung |
|------|----------|-----|----------|--------------|
| `--interactive` | `-i` | `bool` | `false` | Interaktiven TUI-Assistenten ausführen |
| `--force` | `-f` | `bool` | `false` | Bestehende Konfiguration in einem nicht-leeren Verzeichnis überschreiben |

## Wie es funktioniert

### Nicht-interaktiver Modus (Standard)

Erstellt eine minimale `swag2mcp.yaml` ohne Specs. Sie bearbeiten die Datei anschließend manuell.

```bash
swag2mcp init
# Erstellt ~/.swag2mcp/swag2mcp.yaml

swag2mcp init ./my-project
# Erstellt ./my-project/swag2mcp.yaml

swag2mcp init /absoluter/pfad
# Erstellt /absoluter/pfad/swag2mcp.yaml
```

### Interaktiver Modus (`-i`)

Startet einen 18-schrittigen TUI-Assistenten, der Sie durch Folgendes führt:

1. Auswahl des Arbeitsbereichsverzeichnisses
2. Hinzufügen von Specs mit Domain, Titel, Basis-URL
3. Konfigurieren von Collections mit Speicherort-URLs
4. Einrichten der Authentifizierung (alle 9 Methoden)
5. Konfigurieren von HTTP-Client-Einstellungen (Timeout, Proxy, Header usw.)

```bash
swag2mcp init -i
```

### Force-Modus (`--force`)

Standardmäßig weigert sich `init`, in einem nicht-leeren Verzeichnis ausgeführt zu werden. Verwenden Sie `--force` zum Überschreiben:

```bash
swag2mcp init -f
swag2mcp init ./existing-dir -f
```

## Was erstellt wird

```
~/.swag2mcp/
├── swag2mcp.yaml       # Konfigurationsdatei
├── cache/               # Heruntergeladene entfernte Spezifikationsdateien
├── specs/               # Lokale Spezifikationsdateien
├── responses/           # Gespeicherte API-Aufrufantworten
└── auth_scripts/        # Authentifizierungsskripte (für ScriptAuth-Typ)
```

## Überprüfung nach dem Befehl

```bash
ls ~/.swag2mcp/swag2mcp.yaml
# Wenn die Datei existiert, war init erfolgreich
```

## Nuancen

- **Pfadauflösung:** `[path]` ist ein **Arbeitsbereichsverzeichnis**, kein Dateipfad. Die CLI hängt `swag2mcp.yaml` automatisch an. Auflösungsreihenfolge: expliziter `[path]` → aktuelles Verzeichnis (`./`) → `~/.swag2mcp/`.
- **Prüfung auf nicht-leeres Verzeichnis:** Ohne `--force` gibt `init` einen Fehler zurück, wenn das Zielverzeichnis existiert und nicht leer ist. Dies verhindert versehentliches Überschreiben.
- **Auth-Skript-Stubs:** Wenn eine Spec `ScriptAuth` verwendet, erstellt `init` Stub-Skriptdateien (`.sh` unter Unix, `.bat` unter Windows) in `auth_scripts/`.
- **Ausgabe:** Bei Erfolg wird der Konfigurationspfad und ein Hinweis ausgegeben: `"Nächster Schritt: swag2mcp.yaml bearbeiten oder 'swag2mcp ls' ausführen, um konfigurierte Specs aufzulisten"`.
