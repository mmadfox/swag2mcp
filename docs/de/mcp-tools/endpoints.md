# Endpunkt-Tools

Endpunkt-Tools ermöglichen es dem LLM, API-Endpunkte auf verschiedenen Ebenen der Hierarchie anzuzeigen: alle Endpunkte in einer Spec, in einer Collection, in einem Tag oder eine einzelne Endpunkt-Zusammenfassung. Verwenden Sie diese, um verfügbare Operationen zu entdecken, bevor Sie inspizieren oder aufrufen.

---

## endpoint_by_spec

### Zweck

Listet alle Endpunkte über eine gesamte Spec hinweg auf, über alle Collections und Tags hinweg. Gibt die umfassendste Ansicht zurück — jeden Endpunkt in der Spec mit seinem vollständigen Kontext (Tag, Collection, Spec).

### Wann verwenden

- Wenn Sie jeden in einer Spec verfügbaren Endpunkt sehen möchten
- Wenn Sie nicht wissen, welche Collection oder welcher Tag den benötigten Endpunkt enthält
- Nach `spec_by_id`, um die vollständige Endpunkt-Liste zu erhalten

### Parameter

| Parameter | Typ | Erforderlich | Beschreibung |
|-----------|-----|-------------|--------------|
| `specId` | string | Ja | 32-stelliger MD5-Hash der Spec |

### Antwort

```json
{
  "endpoints": [
    {
      "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
      "tagId": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "tagName": "forecast",
      "collectionId": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
      "collectionTitle": "Weather Forecast",
      "specId": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
      "specDomain": "meteo",
      "method": "GET",
      "path": "/v1/forecast",
      "summary": "Wettervorhersage für einen Ort abrufen"
    }
  ]
}
```

| Feld | Typ | Beschreibung |
|------|-----|--------------|
| `id` | string | Endpunkt-Identifikator |
| `tagId` | string | Identifikator des übergeordneten Tags |
| `tagName` | string | Menschenlesbarer Tag-Name |
| `collectionId` | string | Identifikator der übergeordneten Collection |
| `collectionTitle` | string | Menschenlesbarer Collection-Titel |
| `specId` | string | Identifikator der übergeordneten Spec |
| `specDomain` | string | Spec-Domain-Name |
| `method` | string | HTTP-Methode (GET, POST, PUT, DELETE usw.) |
| `path` | string | API-Pfad (z. B. /v1/forecast) |
| `summary` | string | Menschenlesbare Zusammenfassung der Funktion des Endpunkts |

### Nuancen

- Gibt `not_found` zurück, wenn die Spec nicht existiert
- Jeder Endpunkt enthält seine vollständige Abstammung (Spec → Collection → Tag) für den Kontext
- Für eine kurze Zusammenfassung eines einzelnen Endpunkts verwenden Sie `endpoint_by_id`

---

## endpoint_by_collection

### Zweck

Listet alle Endpunkte innerhalb einer bestimmten Collection auf, unabhängig von ihrem Tag. Gibt Endpunkte gruppiert nach Collection mit Spec- und Collection-Metadaten zurück.

### Wann verwenden

- Nach `collection_by_id`, um alle Endpunkte in einer Collection zu sehen
- Wenn Sie die vollständige API-Oberfläche einer Collection erkunden möchten

### Parameter

| Parameter | Typ | Erforderlich | Beschreibung |
|-----------|-----|-------------|--------------|
| `collectionId` | string | Ja | 32-stelliger MD5-Hash der Collection |

### Antwort

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collection": {
    "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
    "title": "Weather Forecast",
    "countMethods": 12
  },
  "endpoints": [
    {
      "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
      "tagId": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "tagName": "forecast",
      "method": "GET",
      "path": "/v1/forecast",
      "summary": "Wettervorhersage für einen Ort abrufen"
    }
  ]
}
```

### Nuancen

- Gibt `not_found` zurück, wenn die Collection nicht existiert
- Enthält Spec- und Collection-Metadaten für den Kontext
- Endpunkte aus allen Tags innerhalb der Collection werden zusammen zurückgegeben

---

## endpoint_by_tag

### Zweck

Listet alle Endpunkte auf, die unter einem bestimmten Tag gruppiert sind. Dies ist die fokussierteste Ansicht — Endpunkte in einem Tag innerhalb einer Collection.

### Wann verwenden

- Nach `tag_by_id`, um die tatsächlichen Endpunkte in einem Tag zu sehen
- Wenn Sie den Tag kennen und seine Operationen sehen möchten

### Parameter

| Parameter | Typ | Erforderlich | Beschreibung |
|-----------|-----|-------------|--------------|
| `tagId` | string | Ja | 32-stelliger MD5-Hash des Tags |

### Antwort

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collection": {
    "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
    "title": "Weather Forecast",
    "countMethods": 12
  },
  "tag": {
    "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
    "title": "forecast",
    "countMethods": 5
  },
  "endpoints": [
    {
      "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
      "method": "GET",
      "path": "/v1/forecast",
      "summary": "Wettervorhersage für einen Ort abrufen"
    }
  ]
}
```

### Nuancen

- Gibt `not_found` zurück, wenn der Tag nicht existiert
- Enthält vollständigen Kontext: Spec-, Collection- und Tag-Metadaten
- Endpunkte sind auf einen einzelnen Tag innerhalb einer einzelnen Collection beschränkt

---

## endpoint_by_id

### Zweck

Ruft eine kurze Zusammenfassung eines einzelnen Endpunkts ab: Methode, Pfad, Zusammenfassung und Veraltungsstatus. Dies ist ein leichtgewichtiges Tool — für das vollständige OpenAPI-Operationsobjekt (Parameter, Anforderungstext, Antwortschemata) verwenden Sie `inspect`.

### Wann verwenden

- Wenn Sie eine Endpunkt-ID haben und eine schnelle Erinnerung an seine Funktion wünschen
- Bevor Sie entscheiden, ob Sie `inspect` für vollständige Details aufrufen
- Wenn Sie die Methode und den Pfad vor dem Aufruf bestätigen müssen

### Parameter

| Parameter | Typ | Erforderlich | Beschreibung |
|-----------|-----|-------------|--------------|
| `id` | string | Ja | 32-stelliger MD5-Hash des Endpunkts |

### Antwort

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collection": {
    "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
    "title": "Weather Forecast",
    "countMethods": 12
  },
  "tag": {
    "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
    "title": "forecast",
    "countMethods": 5
  },
  "endpoint": {
    "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
    "method": "GET",
    "path": "/v1/forecast",
    "summary": "Wettervorhersage für einen Ort abrufen"
  }
}
```

| Feld | Typ | Beschreibung |
|------|-----|--------------|
| `endpoint.id` | string | Endpunkt-Identifikator |
| `endpoint.method` | string | HTTP-Methode |
| `endpoint.path` | string | API-Pfad |
| `endpoint.summary` | string | Menschenlesbare Zusammenfassung |

### Nuancen

- Gibt `not_found` zurück, wenn der Endpunkt nicht existiert
- Dies ist eine **kurze Zusammenfassung** — sie gibt keine Parameter, Anforderungstext oder Antwortschemata zurück
- Für vollständige technische Details (vor `invoke` erforderlich) verwenden Sie `inspect`
