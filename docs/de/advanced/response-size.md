# Antwortgrößenverwaltung

## Übersicht

API-Antworten können sehr groß sein — manchmal zu groß, um in das Kontextfenster des LLM zu passen. swag2mcp verwaltet Antwortgrößen automatisch, indem übergroße Antworten auf der Festplatte gespeichert und Tools zu deren Erkundung bereitgestellt werden.

## Wie es funktioniert

1. **Sie rufen `invoke` auf** — swag2mcp führt die API-Anfrage aus
2. **Wenn die Antwort klein ist** (innerhalb des Limits) — wird sie inline an den LLM zurückgegeben
3. **Wenn die Antwort zu groß ist** (überschreitet das Limit) — wird sie als JSON-Datei in `{workspace}/responses/` gespeichert. Der LLM erhält einen Dateiverweis anstelle der vollständigen Antwort

### Beispiel: kleine Antwort (inline)

```json
{
  "statusCode": 200,
  "body": {
    "id": 1,
    "name": "Rex",
    "status": "available"
  }
}
```

### Beispiel: große Antwort (Dateiverweis)

```json
{
  "statusCode": 200,
  "fileRef": {
    "path": "/Users/user/.swag2mcp/responses/response_a1b2c3d4.json",
    "size": 1572864,
    "sizeHint": "1.5 MB",
    "maxSizeHint": "2 KB",
    "message": "Die Antwort überschreitet das Limit von 2 KB und wurde auf der Festplatte gespeichert.",
    "openCmd": "open /Users/user/.swag2mcp/responses/response_a1b2c3d4.json"
  }
}
```

## Konfiguration

```yaml
http_client:
  max_response_size: 1048576  # 1 MB in Bytes
```

### max_response_size

- **Typ:** `int` (Bytes)
- **Standard:** `1048576` (1 MB)
- **Bereich:** 256 bis 10.485.760 Bytes (10 MB)
- **Wirkung:** Antworten, die größer als dieser Wert sind, werden auf der Festplatte gespeichert, anstatt inline zurückgegeben zu werden
- **Wann erhöhen:** APIs, die große Datensätze zurückgeben (Berichte, Protokolle, Analysen)
- **Wann verringern:** Begrenztes LLM-Kontextfenster oder wenn Sie dateibasierten Zugriff bevorzugen

## Mit großen Antworten arbeiten

Wenn `invoke` einen `fileRef` zurückgibt, verwenden Sie diese drei Tools, um die Daten zu erkunden:

### 1. response_outline — die Struktur verstehen

Ruft eine strukturelle Zusammenfassung der Antwort ab: Schlüssel, Typen, Array-Längen und Navigationshinweise.

```json
→ response_outline(path: "/pfad/zu/datei.json")
← {
    "type": "object",
    "size": 1572864,
    "keys": ["data", "meta"],
    "itemCount": 500,
    "compressionHints": [
      "response_compress(path, 'first_of_array', 'data')",
      "response_compress(path, 'sample_array', 'data', arrayHead=3, arrayTail=2)"
    ]
  }
```

### 2. response_compress — eine kleinere Version erhalten

Komprimiert die Daten, um sie inline einzufügen. Mehrere Komprimierungsmodi ermöglichen Ihnen den richtigen Kompromiss.

| Modus | Beschreibung | Am besten für |
|-------|-------------|---------------|
| `first_of_array` | Nur das erste Element eines Arrays behalten | Wenn alle Elemente dieselbe Struktur haben |
| `sample_array` | Kopf (3) und Ende (2) eines Arrays behalten | Wenn Sie die Wertspanne sehen müssen |
| `truncate_strings` | Jede Zeichenfolge auf N Zeichen kürzen | Wenn Zeichenfolgen sehr lang sind |
| `keys_only` | Werte durch ihre Typnamen ersetzen | Wenn Sie nur die Struktur benötigen |
| `select_keys` | Nur bestimmte Schlüssel behalten | Wenn Sie bestimmte Felder benötigen |

```json
→ response_compress(path: "/pfad/zu/datei.json", mode: "first_of_array", jsonPath: "data")
← {
    "body": [{ "id": 1, "name": "Rex" }],
    "hint": "Array von 500 auf 1 Element mit first_of_array-Modus komprimiert"
  }
```

### 3. response_slice — ein bestimmtes Fragment extrahieren

Ruft ein bestimmtes Element oder einen Wert per JSON-Pfad oder Zeilenbereich ab.

```json
→ response_slice(path: "/pfad/zu/datei.json", jsonPath: "data.0")
← {
    "slice": {
      "value": { "id": 1, "name": "Rex" },
      "jsonPath": "data.0",
      "nextPath": "data.1",
      "prevPath": null
    }
  }
```

## Vollständiger Arbeitsablauf

```
1. invoke(endpoint) → fileRef (Antwort ist 1.5 MB)
2. response_outline(path) → Struktur: { data: Array(500) }
3. response_compress(path, mode: "first_of_array", jsonPath: "data") → erstes Element
4. response_slice(path, jsonPath: "data.0") → vollständige Details des ersten Elements
5. response_slice(path, jsonPath: "data.1") → zweites Element
```

## Automatische Bereinigung

Beim Start des MCP-Servers (`swag2mcp mcp`) werden Antwortdateien, die älter als 48 Stunden sind, automatisch entfernt. Sie können sie auch manuell bereinigen:

```bash
swag2mcp clean
```

## Wichtige Hinweise

- **Das Limit ist in Bytes** — `1048576` = 1 MB, `2097152` = 2 MB usw.
- **Dateiverweise enthalten einen Öffnungsbefehl** — unter macOS ist es `open`, unter Linux `xdg-open`
- **Antwortdateien werden mit zufälligen Suffixen benannt** — keine Konflikte zwischen gleichzeitigen Aufrufen
- **Das Antwortverzeichnis wird automatisch erstellt** — keine manuelle Einrichtung erforderlich
