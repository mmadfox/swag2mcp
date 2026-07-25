# Tags

Ein Tag ist eine Kategorie, die verwandte Endpunkte innerhalb einer Collection gruppiert. Tags können vorhanden sein oder nicht — nicht alle Collections haben sie, und eine Collection kann beliebig viele Tags haben.

Tags stammen aus der OpenAPI/Swagger/Postman-Datei selbst. Es gibt **keine YAML-Konfigurationseinstellungen** für Tags — Sie können Tags in `swag2mcp.yaml` nicht erstellen, umbenennen oder löschen. Die einzige Möglichkeit, Tags zu ändern, besteht darin, die ursprüngliche Spezifikationsdatei zu bearbeiten.

## Hierarchie

```
Spec (domain, z. B. "meteo")
  └── Collection (Spezifikationsdatei, z. B. forecast.yml)
        └── Tag "weather"
              └── GET /forecast
              └── GET /forecast/hourly
        └── Tag "alerts"
              └── GET /alerts
```

## Wie Tags erstellt werden

Tags werden während des Parsens aus dem Spec-Dokument extrahiert:

**OpenAPI 3.x / Swagger 2.0** — die `tags`-Liste jeder Operation wird zu Tags:

```yaml
paths:
  /pet:
    get:
      tags: ["pets"]
      summary: "Haustier nach ID finden"
    post:
      tags: ["pets"]
      summary: "Ein neues Haustier hinzufügen"
  /pet/{petId}/uploadImage:
    post:
      tags: ["pet_images"]
      summary: "Lädt ein Bild hoch"
```

**Postman** — jeder Ordner der obersten Ebene wird zu einem Tag. Verschachtelte Ordner verwenden den Namen des letzten Ordners.

Wenn ein Endpunkt keine Tags hat, wird er unter einem `"default"`-Tag platziert.

## Zweck

Tags helfen dem LLM, Gruppen verwandter Endpunkte zu finden. Anstatt jeden Endpunkt in einer Collection zu durchsuchen, kann der LLM zuerst den richtigen Tag finden und dann nur die Endpunkte darin auflisten.

## MCP-Tools für Tags

| Tool | Beschreibung |
|------|--------------|
| `tag_by_spec` | Alle Tags über eine gesamte Spec hinweg |
| `tag_by_collection` | Tags innerhalb einer bestimmten Collection |
| `tag_by_id` | Tag-Details (Titel, Methodenanzahl) |
| `endpoint_by_tag` | Endpunkte, die unter einem Tag gruppiert sind |

## Beispiel

```
Abfrage: "Zeige alle Tags in der pet-Collection"
→ tag_by_collection(collectionId: "...")
→ Ergebnis: pets (5 Methoden), pet_images (1 Methode)
```

## Einschränkungen

- Tags sind aus Konfigurationssicht schreibgeschützt. Um Tags hinzuzufügen, umzubenennen oder zu entfernen, bearbeiten Sie die ursprüngliche OpenAPI/Swagger/Postman-Datei und führen Sie `swag2mcp update` aus.
- Tags können in der YAML-Konfiguration nicht pro Collection gefiltert oder deaktiviert werden.
