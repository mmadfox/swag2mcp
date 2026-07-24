# Erkennungstools

Erkennungstools ermöglichen es dem LLM, die Spec-Hierarchie zu navigieren: alle Specs finden, in eine Spec hineinzoomen, um ihre Collections zu sehen, und Tags innerhalb einer Collection erkunden. Beginnen Sie mit `spec_list`, um zu sehen, welche APIs verfügbar sind, und verwenden Sie dann IDs, um tiefer zu gehen.

---

## spec_list

### Zweck

Listet alle im Arbeitsbereich registrierten API-Spezifikationen auf. Dies ist der Ausgangspunkt für jede Sitzung — der LLM ruft es zuerst auf, um zu entdecken, welche APIs verfügbar sind.

### Wann verwenden

- Zu Beginn einer Sitzung, um zu sehen, welche APIs konfiguriert sind
- Nach dem Hinzufügen oder Entfernen von Specs, um die Liste zu aktualisieren
- Wenn Sie eine Spec-ID für andere Tools benötigen

### Wie es funktioniert

Gibt eine Liste aller Specs mit ihrer eindeutigen ID und ihrem Domain-Namen zurück. Keine Parameter erforderlich.

### Parameter

Keine.

### Antwort

```json
{
  "specs": [
    {
      "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
      "domain": "meteo"
    },
    {
      "id": "b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7",
      "domain": "dadjoke"
    }
  ]
}
```

| Feld | Typ | Beschreibung |
|------|-----|--------------|
| `id` | string | 32-stelliger MD5-Hash, eindeutiger Identifikator für die Spec |
| `domain` | string | Domain-Name der Spec (z. B. "meteo", "dadjoke") |

### Nuancen

- Gibt nur `id` und `domain` zurück — für vollständige Details (Collections, Tags) verwenden Sie `spec_by_id`
- Alle IDs sind 32-stellige MD5-Hex-Zeichenfolgen (`^[0-9a-f]{32}$`)
- Wenn keine Specs konfiguriert sind, wird ein leeres Array zurückgegeben

---

## spec_by_id

### Zweck

Ruft detaillierte Informationen über eine bestimmte Spec ab: ihre Domain, alle Collections und deren Statistiken (Tag-Anzahl, Methoden-Anzahl).

### Wann verwenden

- Nach `spec_list`, um die Collections innerhalb einer Spec zu sehen
- Wenn Sie Collection-IDs für die weitere Navigation benötigen

### Wie es funktioniert

Nimmt eine Spec-ID entgegen und gibt die Spec-Metadaten plus alle ihre Collections mit Zählwerten zurück.

### Parameter

| Parameter | Typ | Erforderlich | Beschreibung |
|-----------|-----|-------------|--------------|
| `id` | string | Ja | 32-stelliger MD5-Hash der Spec |

### Antwort

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collections": [
    {
      "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
      "title": "Weather Forecast",
      "llmTitle": "Forecast API",
      "countTags": 3,
      "countMethods": 12
    }
  ]
}
```

| Feld | Typ | Beschreibung |
|------|-----|--------------|
| `spec.id` | string | Spec-Identifikator |
| `spec.domain` | string | Spec-Domain-Name |
| `collections[].id` | string | Collection-Identifikator |
| `collections[].title` | string | Menschenlesbarer Titel |
| `collections[].llmTitle` | string | LLM-freundlicher Titel (optional) |
| `collections[].countTags` | int | Anzahl der Tags in der Collection |
| `collections[].countMethods` | int | Anzahl der HTTP-Methoden in der Collection |

### Nuancen

- Gibt einen `not_found`-Fehler zurück, wenn die Spec-ID nicht existiert
- Die `id` muss eine gültige 32-stellige MD5-Hex-Zeichenfolge sein

---

## collection_by_spec

### Zweck

Listet alle Collections innerhalb einer bestimmten Spec auf. Ähnlich wie `spec_by_id`, gibt aber nur die Collection-Liste ohne zusätzliche Spec-Metadaten zurück.

### Wann verwenden

- Wenn Sie bereits die Spec-ID haben und nur die Collection-Liste benötigen
- Als leichtere Alternative zu `spec_by_id`

### Parameter

| Parameter | Typ | Erforderlich | Beschreibung |
|-----------|-----|-------------|--------------|
| `specId` | string | Ja | 32-stelliger MD5-Hash der Spec |

### Antwort

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collections": [
    {
      "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
      "title": "Weather Forecast",
      "llmTitle": "Forecast API",
      "countTags": 3,
      "countMethods": 12
    }
  ]
}
```

### Nuancen

- Gibt `not_found` zurück, wenn die Spec nicht existiert
- Gleiche Daten wie `spec_by_id`, aber ohne den zusätzlichen Spec-Wrapper

---

## collection_by_id

### Zweck

Ruft detaillierte Informationen über eine bestimmte Collection ab: ihre Metadaten, die übergeordnete Spec und alle Tags innerhalb der Collection.

### Wann verwenden

- Nach `collection_by_spec`, um die Tags innerhalb einer Collection zu sehen
- Wenn Sie Tag-IDs für `tag_by_id` oder `endpoint_by_tag` benötigen

### Parameter

| Parameter | Typ | Erforderlich | Beschreibung |
|-----------|-----|-------------|--------------|
| `id` | string | Ja | 32-stelliger MD5-Hash der Collection |

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
  "tags": [
    {
      "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "title": "forecast",
      "countMethods": 5
    },
    {
      "id": "e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7",
      "title": "current",
      "countMethods": 7
    }
  ]
}
```

| Feld | Typ | Beschreibung |
|------|-----|--------------|
| `spec` | object | Übergeordnete Spec (id, domain) |
| `collection` | object | Collection-Metadaten (id, title, countMethods) |
| `tags[]` | array | Liste der Tags mit id, title, countMethods |

### Nuancen

- Gibt `not_found` zurück, wenn die Collection-ID nicht existiert
- Tags werden mit ihren IDs zurückgegeben — verwenden Sie `endpoint_by_tag(tagId)`, um die tatsächlichen Endpunkte zu sehen

---

## tag_by_spec

### Zweck

Listet alle Tags über eine gesamte Spec hinweg auf, über alle Collections hinweg. Nützlich für eine Vogelperspektive aller verfügbaren Tags.

### Wann verwenden

- Wenn Sie alle Tags in einer Spec sehen möchten, ohne in jede Collection hineinzuzoomen
- Wenn Sie nicht wissen, welche Collection den benötigten Tag enthält

### Parameter

| Parameter | Typ | Erforderlich | Beschreibung |
|-----------|-----|-------------|--------------|
| `specId` | string | Ja | 32-stelliger MD5-Hash der Spec |

### Antwort

```json
{
  "tags": [
    {
      "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "title": "forecast",
      "countMethods": 5
    },
    {
      "id": "e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7",
      "title": "current",
      "countMethods": 7
    }
  ]
}
```

### Nuancen

- Gibt `not_found` zurück, wenn die Spec nicht existiert
- Tags werden aus allen Collections in der Spec aggregiert

---

## tag_by_collection

### Zweck

Listet alle Tags innerhalb einer bestimmten Collection auf. Im Gegensatz zu `tag_by_spec` gibt dies auch die übergeordnete Spec und Collection-Metadaten zurück.

### Wann verwenden

- Nach `collection_by_id`, um die Tag-Liste zu bestätigen
- Wenn Sie den vollständigen Kontext benötigen (Spec + Collection + Tags)

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
  "tags": [
    {
      "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "title": "forecast",
      "countMethods": 5
    }
  ]
}
```

### Nuancen

- Gibt `not_found` zurück, wenn die Collection nicht existiert
- Gleiche Tag-Daten wie `tag_by_spec`, aber auf eine Collection beschränkt

---

## tag_by_id

### Zweck

Ruft Informationen über einen einzelnen Tag ab: seine ID, seinen Titel und wie viele Methoden er enthält. Dies gibt Auskunft über den Tag selbst — um die tatsächlichen Endpunkte zu sehen, verwenden Sie `endpoint_by_tag`.

### Wann verwenden

- Wenn Sie eine Tag-ID haben und deren Namen und Größe bestätigen möchten
- Vor dem Aufruf von `endpoint_by_tag`, um zu verstehen, wie viele Endpunkte zu erwarten sind

### Parameter

| Parameter | Typ | Erforderlich | Beschreibung |
|-----------|-----|-------------|--------------|
| `id` | string | Ja | 32-stelliger MD5-Hash des Tags |

### Antwort

```json
{
  "tag": {
    "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
    "title": "forecast",
    "countMethods": 5
  }
}
```

| Feld | Typ | Beschreibung |
|------|-----|--------------|
| `tag.id` | string | Tag-Identifikator |
| `tag.title` | string | Menschenlesbarer Tag-Name |
| `tag.countMethods` | int | Anzahl der HTTP-Methoden in diesem Tag |

### Nuancen

- Gibt `not_found` zurück, wenn der Tag nicht existiert
- Dieses Tool gibt nur Tag-Metadaten zurück — verwenden Sie `endpoint_by_tag`, um die tatsächliche Liste der Endpunkte zu erhalten
