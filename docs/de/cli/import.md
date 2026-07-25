# import

## Zweck

Importiert Spezifikationsdateien in den Arbeitsbereich oder stellt einen vollständigen Arbeitsbereich aus einem ZIP-Backup wieder her. Drei Modi decken verschiedene Szenarien ab: Hinzufügen einer einzelnen Spec, Massenimport aus bestehender Konfiguration oder Wiederherstellen eines vollständigen Arbeitsbereichs.

## Wann verwenden

- Sie haben eine Spec-URL oder -Datei und möchten sie zum Arbeitsbereich hinzufügen
- Sie möchten alle in der Konfiguration referenzierten Spezifikationsdateien herunterladen
- Sie müssen einen Arbeitsbereich aus einem von `export` erstellten ZIP-Backup wiederherstellen
- Sie migrieren swag2mcp auf einen anderen Rechner

## Syntax

```bash
swag2mcp import [path] [source] [name] [flags]
```

## Argumente

| Argument | Position | Erforderlich | Beschreibung |
|----------|----------|-------------|--------------|
| `path` | 1 | Nein | Arbeitsbereichsverzeichnis. Wenn nicht angegeben, wird über die Pfadauflösungsregeln ermittelt. |
| `source` | 2 | Variiert | URL oder lokaler Pfad zu einer Spezifikationsdatei oder Pfad zu einem ZIP-Archiv |
| `name` | 3 | Variiert | Domain-Name für die neue Spec |

## Flags

| Flag | Kurzform | Typ | Standard | Beschreibung |
|------|----------|-----|----------|--------------|
| `--spec` | `-s` | `stringSlice` | `nil` | Collections aus angegebenen Specs importieren (kommagetrennt) |
| `--from-zip` | | `string` | `""` | Arbeitsbereich aus einem swag2mcp-Backup-ZIP wiederherstellen |

## Wie es funktioniert

### Modus 1 — Einzelimport von URL oder Datei

Laden Sie eine Spezifikationsdatei herunter und fügen Sie sie mit einem Domain-Namen zum Arbeitsbereich hinzu:

```bash
swag2mcp import https://example.com/spec.yaml myspec
swag2mcp import /pfad/zu/arbeitsbereich https://example.com/spec.yaml myspec
swag2mcp import ./local-spec.yaml myspec
```

Die Spezifikationsdatei wird in `specs/` gespeichert und die Konfiguration wird mit dem neuen Spec-Eintrag aktualisiert.

### Modus 2 — Massenimport aus bestehender Konfiguration

Laden Sie alle Collections für die angegebenen Domains von ihren konfigurierten URLs herunter:

```bash
swag2mcp import --spec meteo
swag2mcp import /pfad/zu/arbeitsbereich --spec meteo,store
```

Die Spezifikationsdatei jeder Collection wird heruntergeladen und in `specs/` gespeichert. Die Konfiguration wird aktualisiert, um auf die lokalen Kopien zu verweisen.

### Modus 3 — Aus ZIP-Backup wiederherstellen

Stellen Sie einen vollständigen Arbeitsbereich aus einem von `swag2mcp export` erstellten ZIP-Archiv wieder her:

```bash
swag2mcp import --from-zip /pfad/zu/sicherung.zip
swag2mcp import /pfad/zu/arbeitsbereich /pfad/zu/sicherung.zip
```

> **Das ZIP muss von `swag2mcp export` erstellt worden sein.** Beliebige ZIP-Dateien funktionieren nicht — das Archiv hat eine spezifische interne Struktur (`swag2mcp.yaml`, `specs/`, `auth_scripts/`).

## Überprüfung nach dem Befehl

```bash
# Einzel- oder Massenimport
swag2mcp ls [path]
# Die neue Spec sollte in der Liste erscheinen

# ZIP-Wiederherstellung
swag2mcp ls [path]
# Alle Specs aus dem Backup sollten erscheinen
```

## Nuancen

- **Massenmodus erfordert Konfiguration:** Bei Verwendung von `--spec` muss die Konfigurationsdatei existieren. Führen Sie bei Bedarf zuerst `init` aus.
- **Einzelimport erstellt Arbeitsbereich:** Wenn der Arbeitsbereich nicht existiert, wird er automatisch erstellt.
- **ZIP-Erkennung:** Ein Positionsargument, das auf `.zip` endet, wird als ZIP-Quelle behandelt. Das Flag `--from-zip` hat Vorrang vor der Positionserkennung.
- **`--force`:** Für die ZIP-Wiederherstellung verfügbar, um einen bestehenden Arbeitsbereich zu überschreiben.
- **HTTP-Client:** Die globalen HTTP-Client-Einstellungen aus der Konfiguration werden während des Imports angewendet (Timeout, Proxy, Header usw.).
