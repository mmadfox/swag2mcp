# MCP-Tools

## Übersicht

swag2mcp bietet **19 MCP-Tools**, die einem LLM-Agenten über das Model Context Protocol vollen Zugriff auf Ihre APIs geben. Diese Tools decken den gesamten Arbeitsablauf ab: Entdecken, welche APIs verfügbar sind, Navigieren in der Spec-Hierarchie, Suchen und Inspizieren von Endpunkten, Ausführen von API-Aufrufen und Arbeiten mit großen Antworten.

### Was die Tools lösen

- **Erkennung** — der LLM kann Specs, Collections und Tags finden, ohne IDs im Voraus zu kennen
- **Navigation** — von Spec → Collection → Tag → Endpunkt in einer strukturierten Hierarchie hineinzoomen
- **Suche** — Volltextsuche über alle Endpunkte, wenn Sie keine ID haben
- **Inspektion** — das vollständige OpenAPI-Operationsobjekt vor einem Aufruf abrufen
- **Ausführung** — echte API-Aufrufe mit automatischer Authentifizierung durchführen
- **Große Antworten verarbeiten** — übergroße Antworten, die nicht inline passen, gliedern, komprimieren und aufteilen

### Schreibgeschützt vs. Veränderlich

| Typ | Anzahl | Tools |
|-----|--------|-------|
| **Schreibgeschützt** | 17 | Alle Erkennungs-, Endpunkt-, Such-, Inspektions-, Info- und Antwort-Tools |
| **Veränderlich** | 2 | `invoke` (führt echte HTTP-Aufrufe durch), `auth` (ruft Tokens ab) |

Schreibgeschützte Tools sind mit `ReadOnlyHint=true` und `IdempotentHint=true` im MCP-Protokoll markiert, was dem LLM signalisiert, dass sie ohne Nebenwirkungen sicher aufgerufen werden können.

### Fehlerbehandlung

Alle Tools geben Fehler als strukturierte `LLMError`-Objekte mit einem maschinenlesbaren Code und einer menschenlesbaren Nachricht zurück, die erklärt, was schiefgelaufen ist und was als nächstes zu tun ist:

| Fehlercode | Bedeutung |
|------------|-----------|
| `validation_failed` | Ungültige Eingabe (falsches ID-Format, fehlende Pflichtfelder) |
| `not_found` | Entität nicht im Index oder Arbeitsbereich gefunden |
| `rate_limit` | Zweiter `invoke`-Aufruf innerhalb von 10 Sekunden auf demselben Endpunkt |
| `invoke_error` | HTTP-Aufruffehler, Download-Fehler |
| `auth_error` | Fehler beim Abrufen des Auth-Tokens |
| `config_error` | Fehler beim Laden oder Speichern der Konfigurationsdatei |
| `parse_error` | Fehler beim Parsen der Spezifikationsdatei |

## Kategorien

| Kategorie | Tools | Beschreibung |
|-----------|-------|--------------|
| **Erkennung** | `spec_list`, `spec_by_id`, `collection_by_spec`, `collection_by_id`, `tag_by_spec`, `tag_by_collection`, `tag_by_id` | In der Spec-Hierarchie navigieren: Specs, Collections und Tags finden |
| **Endpunkte** | `endpoint_by_spec`, `endpoint_by_collection`, `endpoint_by_tag`, `endpoint_by_id` | Endpunkte auf verschiedenen Ebenen der Hierarchie anzeigen |
| **Ausführung** | `search`, `inspect`, `invoke` | Suchen, den vollständigen Vertrag inspizieren und APIs aufrufen |
| **Hilfsprogramme** | `auth`, `info`, `response_outline`, `response_compress`, `response_slice` | Auth-Tokens, Laufzeitinfo und Verarbeitung großer Antworten |
| **Skills** | [Formatierungsleitfaden](/mcp-tools/skills) | Anpassen, wie Tool-Antworten angezeigt werden |

## Vollständige Liste

| Tool | Beschreibung |
|------|--------------|
| `spec_list` | Alle API-Spezifikationen im Arbeitsbereich auflisten |
| `spec_by_id` | Detaillierte Spec-Informationen mit Collections abrufen |
| `collection_by_spec` | Collections innerhalb einer Spec auflisten |
| `collection_by_id` | Collection-Details mit Tags abrufen |
| `tag_by_spec` | Alle Tags über eine Spec hinweg auflisten |
| `tag_by_collection` | Tags innerhalb einer Collection auflisten |
| `tag_by_id` | Tag-Details abrufen (ID, Titel, Methodenanzahl) |
| `endpoint_by_spec` | Alle Endpunkte in einer Spec auflisten |
| `endpoint_by_collection` | Endpunkte in einer Collection auflisten |
| `endpoint_by_tag` | Endpunkte in einem Tag auflisten |
| `endpoint_by_id` | Kurze Endpunkt-Zusammenfassung (Methode, Pfad, Zusammenfassung) |
| `search` | Volltextsuche über alle Endpunkte |
| `inspect` | Vollständige OpenAPI-Operationsdetails (Parameter, Schemata) |
| `invoke` | Echten API-Aufruf ausführen |
| `auth` | Auth-Token oder Header für eine Spec abrufen |
| `info` | Laufzeitinformationen (Version, Specs, Konfiguration) |
| `response_outline` | Strukturelle Zusammenfassung einer großen Antwortdatei |
| `response_compress` | Große Antwort komprimieren, um sie inline einzufügen |
| `response_slice` | Ein Fragment einer großen Antwort extrahieren |

## Navigationshierarchie

```
spec_list
  └── spec_by_id(id)
        └── collection_by_spec(specId)
              └── collection_by_id(id)
                    └── tag_by_collection(collectionId)
                          └── tag_by_id(id)
                                └── endpoint_by_tag(tagId)
                                      └── endpoint_by_id(id)
                                            └── inspect(endpointId)
                                                  └── invoke(endpointId)
```

Wenn Sie keine ID haben, verwenden Sie `search`, um Endpunkte per Abfrage zu finden. Wenn `invoke` einen `fileRef` zurückgibt (Antwort zu groß), verwenden Sie `response_outline` → `response_compress` oder `response_slice`, um die Daten zu erkunden.
