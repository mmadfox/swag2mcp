# clean

## Zweck

Entfernt zwischengespeicherte entfernte Spezifikationen und gespeicherte API-Aufrufantworten. Dies gibt Speicherplatz frei und erzwingt einen frischen Download der Spezifikationsdateien beim nächsten `update` oder `mcp`-Start.

## Wann verwenden

- Spezifikationsdateien haben sich auf dem entfernten Server geändert und Sie möchten eine Aktualisierung erzwingen
- Sie möchten Speicherplatz freigeben
- Sie beheben Probleme mit veraltetem Cache
- Vor dem Ausführen von `update`, um eine saubere Neucachesicherzustellen

## Syntax

```bash
swag2mcp clean [path]
```

## Argumente

| Argument | Position | Erforderlich | Beschreibung |
|----------|----------|-------------|--------------|
| `path` | 1 | Nein | Arbeitsbereichsverzeichnis. Wenn nicht angegeben, wird über die Pfadauflösungsregeln ermittelt. |

## Flags

Keine.

## Wie es funktioniert

```bash
swag2mcp clean
swag2mcp clean ./my-workspace
```

## Was bereinigt wird

| Verzeichnis | Inhalt | Warum |
|-------------|--------|-------|
| `cache/` | Heruntergeladene entfernte Spezifikationsdateien | Erzwingt erneuten Download beim nächsten Zugriff |
| `responses/` | Gespeicherte API-Aufrufantworten | Gibt Speicherplatz frei |

## Was erhalten bleibt

| Verzeichnis | Inhalt | Warum |
|-------------|--------|-------|
| `specs/` | Lokale Spezifikationsdateien | Dies sind Ihre Quelldateien, kein Cache |
| `auth_scripts/` | Authentifizierungsskripte | Diese wurden vom Benutzer erstellt, kein Cache |

## Bereinigung verwaister Auth-Skripte

Nach der Bereinigung entfernt `clean` auch Auth-Skripte für Specs, die nicht mehr in der Konfiguration existieren. Dies verhindert die Ansammlung veralteter Skripte.

## Automatische Bereinigung

Beim Start des MCP-Servers (`swag2mcp mcp`) werden Antworten, die älter als 48 Stunden sind, automatisch entfernt. Normalerweise müssen Sie `clean` für die routinemäßige Wartung nicht manuell ausführen.

## Überprüfung nach dem Befehl

```bash
ls ~/.swag2mcp/cache
# Sollte leer sein (Verzeichnis existiert, hat aber keine Dateien)
```

## Nuancen

- **Keine Konfiguration erforderlich:** `clean` funktioniert auch ohne gültige Konfigurationsdatei. Es entfernt einfach die Cache- und Antwortverzeichnisse.
- **Bereinigung verwaister Einträge nach bestem Bemühen:** Wenn die Konfigurationsdatei beschädigt oder nicht lesbar ist, wird die Bereinigung verwaister Auth-Skripte übersprungen (nicht fatal).
- **Verzeichnisse bleiben erhalten:** Die Verzeichnisse `cache/` und `responses/` selbst bleiben erhalten — nur ihr Inhalt wird entfernt.
