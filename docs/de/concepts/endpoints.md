# Endpunkte

Ein Endpunkt ist eine bestimmte HTTP-Methode + Pfad, der aufgerufen werden kann (z. B. `GET /api/users/{id}`). Endpunkte sind die eigentlichen API-Operationen, die der LLM entdeckt, inspiziert und aufruft.

## Struktur

Jeder Endpunkt enthält:

- **HTTP-Methode**: GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS
- **Pfad**: `/api/v1/users/{id}`
- **Zusammenfassung**: eine kurze Beschreibung dessen, was der Endpunkt tut — sehr nützlich für den LLM, um seinen Zweck auf einen Blick zu verstehen
- **Beschreibung**: eine detaillierte Erklärung des Verhaltens, der Parameter und der Anwendungsfälle des Endpunkts
- **Parameter**: Pfad, Abfrage, Header, Cookie
- **Anforderungstext**: für POST/PUT/PATCH
- **Antworten**: Statuscodes und Antwortschemata

Die Felder `summary` und `description` stammen aus der OpenAPI/Swagger/Postman-Datei. Sie sind der primäre Weg, wie der LLM versteht, was ein Endpunkt tut. Gut geschriebene Zusammenfassungen machen die Endpunkterkennung viel effektiver.

## MCP-Tools für Endpunkte

| Tool | Beschreibung |
|------|--------------|
| `endpoint_by_spec` | Alle Endpunkte in einer Spec |
| `endpoint_by_collection` | Endpunkte in einer Collection |
| `endpoint_by_tag` | Endpunkte in einem Tag |
| `endpoint_by_id` | Kurze Endpunktzusammenfassung |
| `inspect` | Vollständige Endpunktdetails (Schemata, Parameter) |
| `invoke` | Endpunkt aufrufen |
| `search` | Endpunkte nach Text durchsuchen |

## Veraltete Endpunkte

Endpunkte, die in der Spezifikation als `deprecated` markiert sind, werden bei der Inspektion mit einem Hinweis angezeigt.

## Konfiguration

Endpunkte sind aus der Perspektive von swag2mcp **schreibgeschützt**. Es gibt keine YAML-Konfigurationseinstellungen für Endpunkte — Sie können sie in `swag2mcp.yaml` nicht hinzufügen, entfernen, umbenennen oder ändern.

Um Endpunkte zu ändern (neue hinzufügen, Zusammenfassungen aktualisieren, Parameter ändern, als veraltet markieren), bearbeiten Sie die ursprüngliche OpenAPI/Swagger/Postman-Datei und führen Sie `swag2mcp update` aus, um neu zu parsen und neu zu indizieren.

## Beispiel

```
Abfrage: "Zeige Details für GET /pet/{petId}"
→ inspect(endpointId: "abc123...")
→ Ergebnis:
  GET /pet/{petId}
  Zusammenfassung: Haustier nach ID finden
  Beschreibung: Gibt ein einzelnes Haustier anhand seiner ID zurück
  Parameter:
    - petId (path, integer, erforderlich)
  Antworten:
    - 200: Pet-Objekt
    - 400: Fehler
    - 404: Nicht gefunden
```
