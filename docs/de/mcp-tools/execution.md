# Ausführungstools

Ausführungstools sind der Kern von swag2mcp: **search** findet Endpunkte, wenn Sie keine ID haben, **inspect** zeigt den vollständigen OpenAPI-Vertrag, und **invoke** führt den tatsächlichen API-Aufruf aus. Verwenden Sie sie immer in dieser Reihenfolge: search → inspect → invoke.

---

## search

### Zweck

Das einzige Tool zum Finden von Endpunkten, wenn Sie keine Endpunkt-ID haben. Führt eine Volltextsuche über alle Endpunkte aller Specs mit der bluge-Suchmaschine durch.

### Wann verwenden

- Wenn Sie die Endpunkt-ID nicht kennen
- Wenn Sie Endpunkte nach Schlüsselwörtern, Methode, Tag oder Pfad finden möchten
- Wenn Sie entdecken müssen, welche Endpunkte für eine bestimmte Funktion existieren

### Wie es funktioniert

Durchsucht den Volltextindex über alle Specs hinweg. Unterstützt strukturierte Abfragen mit Feldfiltern, booleschen Operatoren, unscharfer Suche, Platzhaltern und mehr.

### Parameter

| Parameter | Typ | Erforderlich | Beschreibung |
|-----------|-----|-------------|--------------|
| `query` | string | Ja | Suchabfrage (unterstützt strukturierte Syntax) |
| `limit` | int | Ja | Maximale Anzahl zurückzugebender Ergebnisse (1-50) |

### Abfragesyntax

| Beispiel | Beschreibung |
|----------|--------------|
| `pet` | Einfache Textsuche über alle Felder |
| `method:GET` | Nach HTTP-Methode filtern |
| `tag:pet` | Nach Tag-Name filtern |
| `path:"/api/v1/users"` | Exakte Pfadsuche |
| `+method:POST +tag:pet` | Beide Bedingungen müssen zutreffen |
| `-method:DELETE` | DELETE-Methoden ausschließen |
| `create~` | Unscharfe Suche (toleriert Tippfehler) |
| `path:/api/v1/*` | Platzhalter-Pfadsuche |
| `/pattern/` | Regex-Suche |
| `term^3` | Relevanz eines Begriffs erhöhen |

**Durchsuchbare Felder:** `method` (Schlüsselwort), `tag` (Schlüsselwort), `path` (Text), `summary` (Text), `_all` (Standard-Textfeld).

**Nicht unterstützt:** Klammern für Gruppierung, explizite `AND`/`OR`-Operatoren, Feldgruppierung.

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

Jedes Ergebnis enthält die vollständige Abstammung (Spec → Collection → Tag), damit der LLM zu verwandten Endpunkten navigieren kann.

### Nuancen

- `limit` muss zwischen 1 und 50 liegen (gibt sonst `validation_failed` zurück)
- `query` ist erforderlich (gibt `validation_failed` zurück, wenn leer)
- Ergebnisse werden in der Reihenfolge der Relevanz zurückgegeben (beste Übereinstimmung zuerst)
- Verwenden Sie Feldfilter (`method:GET`, `tag:pet`), um Ergebnisse einzugrenzen
- Für exakte Pfadübereinstimmung verwenden Sie Anführungszeichen: `path:"/v1/forecast"`

---

## inspect

### Zweck

Ruft das vollständige OpenAPI-Operationsobjekt für einen Endpunkt ab: alle Parameter, Anforderungstext-Schema, Antwortschemata, Basis-URL und vollständige URL. Dies ist das Tool, das **vor** `invoke` aufgerufen werden sollte, um den Vertrag des Endpunkts zu verstehen.

### Wann verwenden

- Immer vor `invoke` — Sie benötigen den vollständigen Vertrag, um einen korrekten Aufruf zu tätigen
- Wenn Sie dem Benutzer die technischen Details einer API erklären müssen
- Wenn Sie erforderliche Parameter, die Anforderungstext-Struktur oder das Antwortformat kennen müssen

### Wie es funktioniert

Sucht den Endpunkt im Index und gibt das vollständige OpenAPI-Operationsobjekt mit allen aufgelösten Schemata zurück.

### Parameter

| Parameter | Typ | Erforderlich | Beschreibung |
|-----------|-----|-------------|--------------|
| `endpointId` | string | Ja | 32-stelliger MD5-Hash des Endpunkts |

### Antwort

```json
{
  "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
  "tagId": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
  "collectionId": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
  "specId": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
  "specDomain": "meteo",
  "method": "POST",
  "path": "/pet",
  "baseUrl": "https://meteo.swagger.io/v2",
  "fullUrl": "https://meteo.swagger.io/v2/pet",
  "operation": {
    "id": "addPet",
    "tags": ["pet"],
    "summary": "Ein neues Haustier hinzufügen",
    "description": "Ein neues Haustier zum Store hinzufügen",
    "deprecated": false,
    "parameters": [
      {
        "name": "petId",
        "in": "path",
        "description": "ID des Haustiers",
        "required": true,
        "schema": {
          "type": "integer",
          "format": "int64"
        }
      }
    ],
    "requestBody": {
      "description": "Hinzuzufügendes Haustier-Objekt",
      "required": true,
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "name": { "type": "string" },
              "status": { "type": "string", "enum": ["available", "pending", "sold"] }
            },
            "required": ["name"]
          }
        }
      }
    },
    "responses": {
      "200": {
        "description": "Erfolgreiche Operation",
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/Pet"
            }
          }
        }
      },
      "405": {
        "description": "Ungültige Eingabe"
      }
    }
  }
}
```

| Feld | Typ | Beschreibung |
|------|-----|--------------|
| `baseUrl` | string | Basis-URL der API (aus der Konfiguration) |
| `fullUrl` | string | Vollständige URL des Endpunkts (Basis + Pfad) |
| `operation.parameters[]` | array | Parameter mit Name, Ort (path/query/header/cookie), Beschreibung, erforderlich-Flag und Schema |
| `operation.requestBody` | object | Anforderungstext mit Inhaltstyp und Schema |
| `operation.responses` | map | Antwortcodes mit Beschreibungen und Schemata |
| `operation.deprecated` | bool | Ob der Endpunkt veraltet ist |

### Nuancen

- Gibt `not_found` zurück, wenn der Endpunkt nicht existiert
- Dies ist das **einzige** Tool, das das vollständige OpenAPI-Operationsobjekt zurückgibt — `endpoint_by_id` gibt nur eine Zusammenfassung zurück
- Rufen Sie immer `inspect` vor `invoke` auf, um erforderliche Parameter und die Textstruktur zu verstehen
- Das `operation`-Objekt enthält `$ref`-Referenzen, die zu ihren vollständigen Schema-Definitionen aufgelöst werden

---

## invoke

### Zweck

Führt einen echten API-Aufruf an einen Endpunkt aus. Dies ist das einzige Tool, das tatsächliche HTTP-Anfragen stellt. Auth wird automatisch angewendet — Sie müssen nicht zuerst `auth` aufrufen.

### Wann verwenden

- Nur nach dem Aufruf von `inspect`, um den Vertrag des Endpunkts zu verstehen
- Nur mit expliziter Benutzerbestätigung für destruktive Operationen (POST, PUT, PATCH, DELETE)
- Wenn der Benutzer den Aufruf einer API anfordert und Sie alle erforderlichen Parameter haben

### Wie es funktioniert

1. Sucht den Endpunkt im Index
2. Setzt Pfadparameter in die URL ein
3. Fügt Abfrageparameter hinzu
4. Fügt Header und Cookies hinzu
5. Serialisiert den Anforderungstext als JSON
6. Holt automatisch Auth (Token, Header, Abfrageparameter) und wendet es an
7. Führt die HTTP-Anfrage aus
8. Gibt die Antwort zurück oder speichert sie in einer Datei, wenn sie zu groß ist

### Parameter

| Parameter | Typ | Erforderlich | Beschreibung |
|-----------|-----|-------------|--------------|
| `endpointId` | string | Ja | 32-stelliger MD5-Hash des Endpunkts |
| `parameters` | object | Nein | Pfad-, Abfrage- und Header-Parameter als Schlüssel-Wert-Paare |
| `requestBody` | object | Nein | Anforderungstext für POST/PUT/PATCH-Anfragen |
| `headers` | object | Nein | Zusätzliche HTTP-Header zum Senden |
| `cookies` | object | Nein | Zusätzliche HTTP-Cookies zum Senden |

### Antwort (inline)

```json
{
  "statusCode": 200,
  "headers": {
    "content-type": "application/json"
  },
  "body": {
    "id": 1,
    "name": "Rex",
    "status": "available"
  }
}
```

### Antwort (Dateiverweis — wenn Text die Größenbeschränkung überschreitet)

```json
{
  "statusCode": 200,
  "headers": {
    "content-type": "application/json"
  },
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

| Feld | Typ | Beschreibung |
|------|-----|--------------|
| `statusCode` | int | HTTP-Antwort-Statuscode |
| `headers` | object | HTTP-Antwort-Header |
| `body` | any | Antworttext (vorhanden, wenn innerhalb der Größenbeschränkung) |
| `fileRef` | object | Dateiverweis (vorhanden, wenn Text die Größenbeschränkung überschreitet) |

### Mit großen Antworten arbeiten

Wenn `invoke` einen `fileRef` zurückgibt, verwenden Sie die Antwort-Tools, um die Daten zu erkunden:

1. **`response_outline(path)`** — die strukturelle Zusammenfassung abrufen (Schlüssel, Typen, Array-Längen)
2. **`response_compress(path, mode)`** — die Daten komprimieren, um sie inline einzufügen
3. **`response_slice(path, jsonPath)`** — ein bestimmtes Fragment extrahieren

### Nuancen

- **Auth ist automatisch:** Das `invoke`-Tool holt automatisch die Authentifizierung aus der Auth-Konfiguration der Spec und wendet sie an. Sie müssen **nicht** zuerst `auth` aufrufen.
- **Ratenbegrenzung:** Jeder Endpunkt hat eine 10-Sekunden-Abklingzeit. Ein zweiter Aufruf desselben Endpunkts innerhalb von 10 Sekunden wird stillschweigend blockiert (gibt `rate_limit`-Fehler zurück).
- **Antwortgrößenlimit:** Standard ist 2 KB (konfigurierbar über `max_response_size`). Wenn die Antwort dieses Limit überschreitet, wird sie in `{workspace}/responses/` gespeichert und ein `FileReference` wird anstelle des inline `body` zurückgegeben.
- **Parameterbehandlung:** Pfadparameter werden in die URL eingesetzt. Abfrageparameter werden angehängt. Parameter aus der Anfrage überschreiben die Standardwerte der Operationsspezifikation.
- **Anforderungstext:** Für POST/PUT/PATCH wird der Text als JSON serialisiert. `Content-Type` wird automatisch auf `application/json` gesetzt.
- **Fehlerbehandlung:** HTTP-Fehler (nicht-2xx) werden als `invoke_error` mit dem Statuscode und dem Antworttext im Hinweis zurückgegeben.
- **Destruktive Operationen:** Rufen Sie POST/PUT/PATCH/DELETE niemals ohne explizite Benutzerbestätigung auf.
