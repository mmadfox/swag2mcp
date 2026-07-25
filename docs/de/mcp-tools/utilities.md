# Hilfsprogramm-Tools

Hilfsprogramm-Tools bieten unterstützende Funktionen: Abrufen von Auth-Tokens, Abrufen von Laufzeitinformationen und Arbeiten mit großen API-Antworten, die nicht inline passen.

---

## auth

### Zweck

Ruft ein Authentifizierungstoken, Header oder Abfrageparameter für eine bestimmte Spec ab. Dies gibt dem LLM Zugriff auf Anmeldeinformationen, die außerhalb von swag2mcp verwendet werden können (z. B. zum Generieren eines curl-Befehls).

### Wann verwenden

- Nur wenn der Benutzer explizit nach dem rohen Token oder den Anmeldeinformationen fragt
- Beim Generieren eines curl-Befehls oder Code-Snippets, das Auth benötigt
- Wenn der Benutzer sehen möchte, welche Auth-Methode konfiguriert ist

### Wann NICHT verwenden

- **Rufen Sie `auth` nicht** vor `inspect` oder `invoke` auf — `invoke` holt und wendet die Authentifizierung automatisch an
- **Rufen Sie `auth` nicht** nur auf, um zu prüfen, ob Auth konfiguriert ist — verwenden Sie stattdessen `info`

### Wie es funktioniert

Sucht die Auth-Konfiguration der Spec und führt den Auth-Ablauf aus (Token-Austausch, Skriptausführung usw.), um die aktuellen Anmeldeinformationen zu erhalten.

### Parameter

| Parameter | Typ | Erforderlich | Beschreibung |
|-----------|-----|-------------|--------------|
| `specId` | string | Ja | 32-stelliger MD5-Hash der Spec |

### Antwort

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "headers": {
    "Authorization": "Bearer eyJhbGciOiJIUzI1NiIs...",
    "X-API-Key": "my-api-key"
  },
  "queryParams": {
    "api_key": "my-api-key"
  }
}
```

| Feld | Typ | Beschreibung |
|------|-----|--------------|
| `token` | string | Roher Token-Wert (Bearer-Token, API-Schlüssel usw.) |
| `headers` | object | HTTP-Header, die in Anfragen eingefügt werden sollen |
| `queryParams` | object | Abfrageparameter, die in Anfragen eingefügt werden sollen |

### Nuancen

- **Standardmäßig in der Produktion deaktiviert:** Das Flag `--disable-llm-auth` (Standard: `true`) entfernt das `auth`-Tool vollständig aus der MCP-Tool-Liste. Der LLM kann Tokens weder sehen noch anfordern. Setzen Sie `--disable-llm-auth=false`, um es für Debugging oder kurzlebige Tokens zu aktivieren.
- **`invoke` behandelt Auth automatisch:** Sie müssen `auth` nicht vor `invoke` aufrufen. Der Invoke-Dienst holt und wendet automatisch die korrekte Authentifizierung an.
- **Unterstützt 9 Auth-Methoden:** `none`, `basic`, `bearer`, `digest`, `hmac`, `oauth2-cc` (Client Credentials), `oauth2-pwd` (Passwort), `api-key`, `script`.
- Gibt `auth_error` zurück, wenn die Auth-Methode fehlschlägt (z. B. OAuth2-Token-Endpunkt nicht erreichbar, Skriptausführung fehlgeschlagen).

---

## info

### Zweck

Gibt eine umfassende Zusammenfassung der swag2mcp-Laufzeit zurück: Version, Arbeitsbereichspfad, aktive Specs, HTTP-Client-Einstellungen, MCP-Transportkonfiguration, Auth-Methoden und Mock-Modus-Status.

### Wann verwenden

- Wenn der Benutzer nach der Systemkonfiguration fragt
- Wenn Sie Laufzeiteinstellungen überprüfen müssen (Timeout, Antwortgrößenlimit, Transport)
- Wenn Sie wissen müssen, welche Auth-Methoden verfügbar sind
- Bei der Fehlerbehebung von Konfigurationsproblemen

### Wie es funktioniert

Gibt eine vorberechnete Momentaufnahme des Laufzeitzustands zurück. Keine Parameter erforderlich.

### Parameter

Keine.

### Antwort

```json
{
  "version": "v1.2.0",
  "workspace": "~/.swag2mcp",
  "uptime": "2h 15m",
  "specs": {
    "total": 4,
    "active": 3,
    "disabled": 1,
    "collections": 6,
    "endpoints": 42
  },
  "http_client": {
    "timeout": "30s",
    "max_response_size": "2 KB",
    "follow_redirects": true,
    "max_redirects": 10,
    "randomize": false,
    "proxy": null,
    "headers": {},
    "cookies": []
  },
  "mcp": {
    "transport": "stdio",
    "addr": ":8080",
    "path": "/mcp",
    "auth_enabled": false
  },
  "auth": {
    "methods": ["bearer", "api-key"]
  },
  "mock": {
    "enabled": false
  }
}
```

| Feld | Typ | Beschreibung |
|------|-----|--------------|
| `version` | string | swag2mcp-Version |
| `workspace` | string | Pfad zum Arbeitsbereichsverzeichnis |
| `uptime` | string | Server-Betriebszeit (menschenlesbar) |
| `specs` | object | Spec-Zusammenfassung: gesamt, aktiv, deaktiviert, Collections, Endpunkte |
| `http_client` | object | HTTP-Client-Konfiguration |
| `http_client.max_response_size` | string | Maximale Antwortgröße in menschenlesbarem Format (z. B. "2 KB") |
| `mcp` | object | MCP-Server-Konfiguration |
| `auth` | object | Verfügbare Auth-Methoden |
| `mock` | object | Mock-Server-Status |

### Nuancen

- `max_response_size` wird in menschenlesbarem Format angezeigt (z. B. `"1 KB"`, `"2 MB"`)
- `uptime` wird aus der Server-Startzeit berechnet
- Die Daten sind eine Momentaufnahme zum Zeitpunkt des Bootstraps — sie spiegeln den Zustand wider, als der MCP-Server gestartet wurde

---

## response_outline

### Zweck

Ruft eine strukturelle Zusammenfassung auf hoher Ebene einer großen JSON-Antwortdatei ab, die von `invoke` auf der Festplatte gespeichert wurde. Sie gibt die Form der Daten zurück — Schlüssel, Typen, Array-Längen und Navigationshinweise — ohne die tatsächlichen Werte zurückzugeben.

### Wann verwenden

- Unmittelbar nachdem `invoke` einen `fileRef` zurückgegeben hat (Antwort zu groß für inline)
- Dies ist der **obligatorische erste Schritt** im Arbeitsablauf für große Antworten

### Wie es funktioniert

Liest die gespeicherte Antwortdatei und analysiert ihre Struktur: Typ der obersten Ebene, Schlüssel, Array-Längen, Verschachtelungstiefe und Komprimierungshinweise.

### Parameter

| Parameter | Typ | Erforderlich | Beschreibung |
|-----------|-----|-------------|--------------|
| `path` | string | Ja | Absoluter Pfad von `fileRef.path` |
| `maxDepth` | int | Nein | Maximale Rekursionstiefe (Standard: 3) |
| `maxArrayItems` | int | Nein | Wie viele Array-Elemente inspiziert werden sollen (Standard: 5) |

### Antwort

```json
{
  "outline": {
    "type": "object",
    "size": 1572864,
    "lineCount": 12500,
    "depth": 3,
    "structure": {
      "type": "object",
      "keys": ["data", "meta", "error"],
      "data": {
        "type": "array",
        "length": 500,
        "items": {
          "type": "object",
          "keys": ["id", "name", "status", "createdAt"]
        }
      }
    },
    "schemaHint": "Objekt mit 3 Schlüsseln: data (array[500]), meta (object), error (null)",
    "keys": ["data", "meta", "error"],
    "itemCount": 500,
    "itemType": "object",
    "compressionHints": [
      "response_compress(path, 'first_of_array', 'data')",
      "response_compress(path, 'sample_array', 'data', arrayHead=3, arrayTail=2)",
      "response_compress(path, 'keys_only', 'data')",
      "response_compress(path, 'select_keys', 'data', selectKeys=[id, name])"
    ],
    "navigationHints": {
      "paths": ["data", "meta", "error"],
      "arrays": [
        {"path": "data", "length": 500}
      ]
    }
  }
}
```

| Feld | Typ | Beschreibung |
|------|-----|--------------|
| `type` | string | Typ der obersten Ebene: "object" oder "array" |
| `size` | int | Dateigröße in Bytes |
| `lineCount` | int | Anzahl der Zeilen in der Datei |
| `depth` | int | Maximale inspizierte Verschachtelungstiefe |
| `structure` | object | Rekursive Struktur mit Schlüsseln, Typen, Array-Längen |
| `schemaHint` | string | Einzeilige Zusammenfassung der Form der obersten Ebene |
| `keys` | array | Schlüssel der obersten Ebene (für Objekte) |
| `itemCount` | int | Array-Länge (für Arrays) |
| `compressionHints` | array | Vorgeschlagene `response_compress`-Aufrufe mit Parametern |
| `navigationHints` | object | Pfade und Arrays der obersten Ebene mit Längen |

### Nuancen

- Gibt `validation_failed` zurück, wenn der Pfad ungültig ist oder nicht im Antwortverzeichnis liegt
- Gibt `not_found` zurück, wenn die Datei nicht existiert
- Gibt `validation_failed` zurück, wenn die Datei kein gültiges JSON ist
- Das Feld `compressionHints` bietet gebrauchsfertige Vorschläge für `response_compress`-Aufrufe

---

## response_compress

### Zweck

Reduziert einen JSON-Wert innerhalb einer gespeicherten Antwortdatei, sodass er innerhalb des Antwortgrößenlimits liegt und inline an den LLM zurückgegeben werden kann. Mehrere Komprimierungsmodi ermöglichen Ihnen den richtigen Kompromiss zwischen Größe und Informationsgehalt.

### Wann verwenden

- Nach `response_outline`, um die Struktur zu verstehen
- Wenn Sie Daten aus einer großen Antwort inline benötigen
- Wenn `response_slice` zu eng ist und Sie eine breitere Ansicht benötigen

### Wie es funktioniert

Liest die gespeicherte Antwortdatei, navigiert zum angegebenen JSON-Pfad, wendet den Komprimierungsmodus an und gibt das komprimierte Ergebnis zurück. Wenn das Ergebnis immer noch die Größenbeschränkung überschreitet, wird es in einer neuen Datei gespeichert.

### Parameter

| Parameter | Typ | Erforderlich | Beschreibung |
|-----------|-----|-------------|--------------|
| `path` | string | Ja | Absoluter Pfad von `fileRef.path` |
| `jsonPath` | string | Nein | Pfad zum zu komprimierenden Wert (z. B. `data` oder `data.0`) |
| `mode` | string | Ja | Komprimierungsmodus (siehe Tabelle unten) |
| `arrayHead` | int | Nein | Führende Elemente, die im `sample_array`-Modus behalten werden sollen (Standard: 3) |
| `arrayTail` | int | Nein | Nachfolgende Elemente, die im `sample_array`-Modus behalten werden sollen (Standard: 2) |
| `stringLen` | int | Nein | Maximale Zeichenfolgenlänge im `truncate_strings`-Modus (Standard: 80) |
| `selectKeys` | array | Nein | Schlüssel, die im `select_keys`-Modus behalten werden sollen |

### Komprimierungsmodi

| Modus | Beschreibung | Am besten für |
|-------|-------------|---------------|
| `first_of_array` | Nur das erste Element eines Arrays behalten | Wenn alle Elemente dieselbe Struktur haben |
| `sample_array` | Kopf und Ende eines Arrays behalten | Wenn Sie die Wertspanne sehen müssen |
| `truncate_strings` | Jede Zeichenfolge auf `stringLen` Zeichen kürzen | Wenn Zeichenfolgen sehr lang sind, aber die Struktur wichtig ist |
| `keys_only` | Objektwerte durch ihre Typnamen ersetzen | Wenn Sie nur die Struktur benötigen |
| `select_keys` | Nur bestimmte Schlüssel in jedem Objekt behalten | Wenn Sie bestimmte Felder aus vielen Objekten benötigen |

### Antwort

```json
{
  "body": [
    { "id": 1, "name": "Rex", "status": "available" },
    { "id": 2, "name": "Max", "status": "pending" }
  ],
  "hint": "Array von 500 auf 2 Elemente mit first_of_array-Modus komprimiert"
}
```

| Feld | Typ | Beschreibung |
|------|-----|--------------|
| `body` | any | Komprimierter JSON-Wert (vorhanden, wenn innerhalb der Größenbeschränkung) |
| `fileRef` | object | Dateiverweis (vorhanden, wenn immer noch zu groß) |
| `hint` | string | Erklärung, was komprimiert wurde |

### Nuancen

- Wenn das komprimierte Ergebnis immer noch `max_response_size` überschreitet, wird es in einer neuen Datei gespeichert und ein `FileReference` zurückgegeben
- Standardwerte: `arrayHead=3`, `arrayTail=2`, `stringLen=80`
- Gibt `validation_failed` für ungültigen Pfad, ungültigen JSONPath oder Nicht-JSON-Datei zurück
- Gibt `not_found` zurück, wenn die Datei nicht existiert oder JSONPath nicht übereinstimmt

---

## response_slice

### Zweck

Extrahiert ein bestimmtes Fragment einer gespeicherten JSON-Antwortdatei nach logischem JSON-Pfad oder nach Zeilenbereich. Im Gegensatz zu `response_compress` erhalten Sie hier die rohen, unveränderten Daten.

### Wann verwenden

- Wenn Sie ein bestimmtes Element oder einen bestimmten Wert aus einer großen Antwort benötigen
- Wenn `response_compress` nicht genügend Details liefert
- Wenn Sie Schritt für Schritt durch eine Antwort navigieren möchten

### Wie es funktioniert

Liest die gespeicherte Antwortdatei und extrahiert ein Fragment nach JSON-Pfad (z. B. `data.3.name`) oder nach Zeilenbereich (z. B. `120-240`). Gibt Navigationshinweise zum Durchlaufen von Arrays und Objekten zurück.

### Parameter

| Parameter | Typ | Erforderlich | Beschreibung |
|-----------|-----|-------------|--------------|
| `path` | string | Ja | Absoluter Pfad von `fileRef.path` |
| `jsonPath` | string | Nein | Logischer Pfad zum Wert (z. B. `data.3.name`) |
| `line` | int | Nein | 1-basierte Zeilennummer, um die das Fragment zentriert werden soll |
| `range` | string | Nein | Zeilenbereich als `start-end` (z. B. `120-240`) |
| `around` | int | Nein | Zeilen, die um `line` herum eingeschlossen werden sollen (Standard: 20) |

### Antwort

```json
{
  "slice": {
    "lines": [120, 130],
    "fragment": "{\n  \"id\": 1,\n  \"name\": \"Rex\"\n}",
    "value": {
      "id": 1,
      "name": "Rex"
    },
    "jsonPath": "data.0",
    "context": "object",
    "isComplete": true,
    "nextLine": 131,
    "prevLine": 119,
    "nextPath": "data.1",
    "prevPath": null
  }
}
```

| Feld | Typ | Beschreibung |
|------|-----|--------------|
| `lines` | array | 1-basierter Zeilenbereich [start, end] |
| `fragment` | string | Roher JSON-Text (wenn klein genug) |
| `value` | any | Extrahierter JSON-Wert |
| `jsonPath` | string | Der verwendete JSON-Pfad |
| `context` | string | "object", "array" oder "value" |
| `isComplete` | bool | Wahr, wenn der Wert ein gültiges JSON-Fragment ist |
| `nextLine` | int | Vorgeschlagene nächste Zeile für zeilenbasierte Navigation |
| `prevLine` | int | Vorgeschlagene vorherige Zeile |
| `nextPath` | string | Vorgeschlagener nächster JSON-Pfad für Array-Navigation |
| `prevPath` | string | Vorgeschlagener vorheriger JSON-Pfad |

### Nuancen

- **Bevorzugen Sie `jsonPath` gegenüber Zeilennummern** — JSON-Pfade sind stabil und beschreibend, Zeilennummern ändern sich, wenn die Datei neu generiert wird
- Wenn das extrahierte Fragment `max_response_size` überschreitet, wird es in einer neuen Datei gespeichert und ein `FileReference` zurückgegeben
- Standard `around` ist 20 Zeilen
- Die Antwort enthält `nextPath`/`prevPath` zum Durchlaufen von Arrays und `nextLine`/`prevLine` für zeilenbasierte Navigation
- Gibt `validation_failed` für ungültigen Pfad, ungültigen JSONPath, ungültige Zeile/ Bereich oder Nicht-JSON-Datei zurück
- Gibt `not_found` zurück, wenn die Datei nicht existiert oder JSONPath nicht übereinstimmt
