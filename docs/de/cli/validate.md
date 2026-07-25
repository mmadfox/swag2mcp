# validate

## Zweck

Überprüft die Konfigurationsdatei und alle referenzierten Spezifikationsdateien auf Fehler. Dies ist ein **schreibgeschützter** Diagnosebefehl — er ändert niemals etwas.

## Wann verwenden

- Nach manueller Bearbeitung von `swag2mcp.yaml`
- Vor dem Ausführen von `mcp` oder `update`, um Probleme frühzeitig zu erkennen
- Bei der Fehlerbehebung, warum eine Spec nicht geladen wird
- In CI/CD-Pipelines zur Validierung von Konfigurationsänderungen

## Syntax

```bash
swag2mcp validate [path] [flags]
```

## Argumente

| Argument | Position | Erforderlich | Beschreibung |
|----------|----------|-------------|--------------|
| `path` | 1 | Nein | Arbeitsbereichsverzeichnis. Wenn nicht angegeben, wird über die Pfadauflösungsregeln ermittelt. |

## Flags

| Flag | Kurzform | Typ | Standard | Beschreibung |
|------|----------|-----|----------|--------------|
| `--tags` | `-t` | `string` | `""` | Nur Specs mit passenden Tags validieren (kommagetrennt) |

## Wie es funktioniert

```bash
swag2mcp validate
swag2mcp validate ./my-workspace
swag2mcp validate --tags=public
```

## Was geprüft wird

| Prüfung | Beschreibung |
|---------|--------------|
| YAML-Syntax | Die Konfigurationsdatei muss gültiges YAML sein |
| Konfigurationsstruktur | Alle erforderlichen Felder vorhanden, Typen korrekt |
| Domain-Eindeutigkeit | Keine doppelten Domains |
| Domain-Format | Nur Kleinbuchstaben, Ziffern, Bindestriche |
| Spezifikationsdatei-Existenz | Die `location`-Datei oder URL muss erreichbar sein |
| Spezifikationsformat | Die Datei muss gültiges OpenAPI 3.x, Swagger 2.0 oder Postman-Collection sein |
| Auth-Einstellungen | Auth-Typ und Konfiguration sind für die ausgewählte Methode gültig |
| HTTP-Client | HTTP-Client-Einstellungen sind gültig |

## Was NICHT geprüft wird

| Nicht geprüft | Grund |
|---------------|-------|
| Authentifizierungsendpunkte | `validate` prüft die Auth-Konfigurationssyntax, testet aber keine Anmeldung/Token-Austausch |
| API-Endpunktverfügbarkeit | Nur die Spezifikationsdatei-URL wird geprüft, nicht die `base_url` |
| `base_url`-Korrektheit | Das Format wird validiert, aber es wird keine Testanfrage gestellt |
| Mock-Server-Konfiguration | `base_mock_url` wird nicht auf Konnektivität überprüft |

## Beispielausgabe

```
✅ Konfiguration ist gültig.
✓ Spec petstore: OK
✓ Spec meteo: OK
✗ Spec old-api: Datei nicht gefunden
```

## Überprüfung nach dem Befehl

Wenn die Validierung bestanden wird, ist die Konfiguration bereit für `mcp`, `update` oder `run`.

## Nuancen

- **Kein Auto-Init:** Anders als `add`, `ls` oder `run` führt `validate` **keine** automatische Initialisierung durch, wenn die Konfiguration fehlt. Es gibt einen Fehler zurück: `"Konfiguration nicht gefunden unter &lt;path&gt;"`.
- **Netzwerkzugriff:** Entfernte Spec-URLs werden während der Validierung abgerufen. Der Befehl kann länger dauern, wenn Spezifikationen auf langsamen Servern gehostet werden.
- **Tag-Filterung:** Wenn `--tags` gesetzt ist, werden nur Specs mit den angegebenen Tags validiert. Andere Specs werden übersprungen.
