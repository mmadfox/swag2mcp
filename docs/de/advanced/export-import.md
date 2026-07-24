# Export und Import

## Übersicht

swag2mcp unterstützt die vollständige Arbeitsbereichs-Roundtrip-Funktionalität über ZIP-Archive. Sie können Ihren gesamten Arbeitsbereich (Konfiguration, Spezifikationsdateien, Authentifizierungsskripte) in eine ZIP-Datei exportieren und auf einem anderen Rechner wiederherstellen.

## Export

Erstellt ein portables ZIP-Backup Ihres Arbeitsbereichs.

```bash
# In Standarddatei exportieren (swag2mcp-backup-&lt;timestamp&gt;.zip)
swag2mcp export

# Mit benutzerdefiniertem Pfad exportieren
swag2mcp export --output ~/backups/swag2mcp-backup.zip

# Nur bestimmte Specs exportieren
swag2mcp export --spec meteo
swag2mcp export --spec meteo,store
```

### Was im Export enthalten ist

| Element | Beschreibung |
|---------|--------------|
| `swag2mcp.yaml` | Konfigurationsdatei |
| `specs/` | Alle Spezifikationsdateien (OpenAPI/Swagger/Postman) |
| `auth_scripts/` | Authentifizierungsskripte |
| `swag2mcp.meta` | Metadaten (Versionsinfo für Kompatibilität) |

Cache und Antworten werden **nicht** exportiert — sie sind flüchtig und wären bei der Wiederherstellung veraltet.

### Standard-Dateiname

Wenn Sie keinen Ausgabepfad angeben, wird die Datei als `swag2mcp-backup-<YYYY-MM-DD-HHMMSS>.zip` im aktuellen Verzeichnis gespeichert (UTC-Zeitstempel).

## Import

Stellen Sie einen Arbeitsbereich aus einem ZIP-Backup wieder her oder importieren Sie Spezifikationsdateien.

### Aus ZIP wiederherstellen

```bash
# Vollständigen Arbeitsbereich wiederherstellen
swag2mcp import --from-zip /pfad/zu/sicherung.zip

# Mit Überschreiben wiederherstellen
swag2mcp import --from-zip /pfad/zu/sicherung.zip -f
```

Das ZIP muss von `swag2mcp export` erstellt worden sein — beliebige ZIP-Dateien funktionieren nicht.

### Einzelne Spezifikationsdatei importieren

Laden Sie eine Spezifikationsdatei herunter und fügen Sie sie zum Arbeitsbereich hinzu:

```bash
swag2mcp import https://example.com/spec.yaml myspec
swag2mcp import /pfad/zu/arbeitsbereich https://example.com/spec.yaml myspec
```

### Massenimport aus bestehender Konfiguration

Laden Sie alle Collection-Spezifikationsdateien für die angegebenen Specs (Domains) herunter:

```bash
swag2mcp import --spec meteo
swag2mcp import /pfad/zu/arbeitsbereich --spec meteo,store
```

Dies lädt die Spezifikationsdatei jeder Collection herunter, speichert sie in `specs/` und aktualisiert die Konfiguration, um auf die lokale Kopie zu verweisen.

## Anwendungsfälle

### Backup

```bash
swag2mcp export --output swag2mcp-$(date +%Y-%m-%d).zip
```

### Auf einen anderen Rechner übertragen

```bash
# Auf altem Rechner
swag2mcp export --output swag2mcp.zip

# ZIP auf den neuen Rechner kopieren, dann:
swag2mcp import --from-zip swag2mcp.zip
```

### Konfiguration teilen

```bash
swag2mcp init
swag2mcp export --output template.zip
# template.zip mit einem Kollegen teilen
```

## Überprüfung nach dem Export

Überprüfen Sie immer, ob die ZIP-Datei erstellt wurde:

```bash
ls -la swag2mcp-backup-*.zip
```

## Wichtige Hinweise

- **Die Ausgabe muss ein Dateipfad sein, der auf `.zip` endet** — übergeben Sie kein Verzeichnis
- **Cache und Antworten sind ausgeschlossen** — nur die Konfiguration, Spezifikationen und Authentifizierungsskripte werden gespeichert
- **Das ZIP ist in sich geschlossen** — es kann auf jedem Rechner mit installiertem swag2mcp wiederhergestellt werden
- **Spec-Filter** — verwenden Sie `--spec`, um nur bestimmte Specs zu exportieren oder zu importieren
