# export

## Zweck

Erstellt ein portables ZIP-Backup des Arbeitsbereichs. Das Archiv enthält die Konfigurationsdatei, alle Spezifikationsdateien und Auth-Skripte — alles, was zur Wiederherstellung des Arbeitsbereichs auf einem anderen Rechner benötigt wird.

## Wann verwenden

- Sie möchten Ihren Arbeitsbereich vor Änderungen sichern
- Sie migrieren swag2mcp auf einen anderen Rechner
- Sie möchten Ihre API-Konfiguration mit einem Kollegen teilen
- Sie bereiten eine reproduzierbare Umgebung vor

## Syntax

```bash
swag2mcp export [path] [output] [flags]
```

## Argumente

| Argument | Position | Erforderlich | Beschreibung |
|----------|----------|-------------|--------------|
| `path` | 1 | Nein | Arbeitsbereichsverzeichnis. Wenn nicht angegeben, wird über die Pfadauflösungsregeln ermittelt. |
| `output` | 2 | Nein | Vollständiger Pfad für die Ausgabe-ZIP-Datei. Wenn nicht angegeben, Standard: `./swag2mcp-backup-&lt;timestamp&gt;.zip`. |

## Flags

| Flag | Kurzform | Typ | Standard | Beschreibung |
|------|----------|-----|----------|--------------|
| `--spec` | `-s` | `stringSlice` | `nil` | Nur angegebene Specs exportieren (kommagetrennt) |

## Wie es funktioniert

### Standard-Export

Erstellt eine ZIP im aktuellen Verzeichnis mit einem Zeitstempelnamen:

```bash
swag2mcp export
# Erstellt ./swag2mcp-backup-2026-07-22-143022.zip
```

### Benutzerdefinierter Ausgabepfad

```bash
swag2mcp export /pfad/zu/arbeitsbereich /pfad/zu/sicherung.zip
```

### Bestimmte Specs exportieren

```bash
swag2mcp export --spec meteo
swag2mcp export --spec meteo,store
```

## Was in der ZIP ist

| Eintrag | Beschreibung |
|---------|--------------|
| `swag2mcp.meta` | Metadaten über den Export |
| `swag2mcp.yaml` | Konfigurationsdatei |
| `specs/` | Alle Spezifikationsdateien (OpenAPI/Swagger/Postman) |
| `auth_scripts/` | Authentifizierungsskripte |
| `cache/` | Leer (Cache wird nicht exportiert) |
| `responses/` | Leer (Antworten werden nicht exportiert) |

## Wiederherstellung

Verwenden Sie `import`, um aus einem Backup wiederherzustellen:

```bash
swag2mcp import --from-zip /pfad/zu/sicherung.zip
```

## Überprüfung nach dem Befehl

Überprüfen Sie immer, ob die ZIP-Datei erstellt wurde:

```bash
ls -la swag2mcp-backup-*.zip
# oder für einen benutzerdefinierten Ausgabepfad:
ls -la /pfad/zu/sicherung.zip
```

## Nuancen

- **Ausgabe muss ein Dateipfad sein:** Das Argument `[output]` muss ein vollständiger Dateipfad sein, der auf `.zip` endet. Übergeben Sie **kein** Verzeichnis — der Befehl erstellt keine ZIP, wenn ein Verzeichnispfad angegeben wird.
- **Standard-Dateiname:** `swag2mcp-backup-<YYYY-MM-DD-HHMMSS>.zip` mit UTC-Zeitstempel.
- **`--spec`-Filter:** Wenn gesetzt, werden nur die angegebenen Specs eingeschlossen. Andere Specs werden vom Archiv ausgeschlossen.
- **Keine Konfiguration erforderlich:** `export` funktioniert auch ohne gültige Konfigurationsdatei. Es exportiert, was im Arbeitsbereich existiert.
- **Cache und Antworten sind ausgeschlossen:** Dies sind flüchtige Daten, die bei der Wiederherstellung veraltet wären. Nur die Konfiguration, Spezifikationen und Auth-Skripte werden gespeichert.
