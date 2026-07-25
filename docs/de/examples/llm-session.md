# LLM-Sitzungsbeispiele

Diese Seite zeigt echte Dialoge zwischen einem Benutzer und einem LLM-Agenten, der von swag2mcp betrieben wird. Jedes Beispiel enthält die Anfrage des Benutzers, den MCP-Tool-Aufruf und die Antwort.

## Einrichtung

Der Arbeitsbereich hat zwei APIs konfiguriert:

- **icanhazdadjoke** — öffentliche Dad-Jokes-API
- **Open-Meteo** — Wettervorhersage-API

---

## Verfügbare APIs entdecken

**Benutzer:** Auf welche APIs hast du Zugriff?

**LLM → MCP:** `spec_list()`

**MCP → LLM:**
```
Specifications:
  dadjoke (icanhazdadjoke)
  meteo (Open-Meteo Weather API)
```

**LLM:** Ich habe Zugriff auf zwei APIs: icanhazdadjoke für Dad Jokes und Open-Meteo für Wetterdaten.

---

## Spec-Details abrufen

**Benutzer:** Erzähl mir mehr über die dadjoke-API.

**LLM → MCP:** `spec_by_id(id: "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6")`

**MCP → LLM:**
```
Spec: dadjoke (icanhazdadjoke)
Collections:
  - Jokes (3 Endpunkte, 1 Tag)
```

**LLM:** Die dadjoke-API hat eine Collection namens "Jokes" mit 3 Endpunkten unter einem einzigen Tag.

---

## Nach Endpunkten suchen

**Benutzer:** Finde einen Endpunkt, um einen zufälligen Dad Joke zu bekommen.

**LLM → MCP:** `search(query: "random joke", limit: 5)`

**MCP → LLM:**
```
Endpoints:
  GET /  — Einen zufälligen Dad Joke zurückgeben
  GET /j/{id} — Einen bestimmten Dad Joke nach ID zurückgeben
  GET /search — Nach Dad Jokes suchen
```

**LLM:** Ich habe drei Endpunkte gefunden. Der wichtigste ist `GET /`, der einen zufälligen Dad Joke zurückgibt.

---

## Einen Endpunkt inspizieren

**Benutzer:** Welche Parameter braucht der Zufallsjoke-Endpunkt?

**LLM → MCP:** `inspect(endpointId: "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6")`

**MCP → LLM:**
```
GET /
Basis-URL: https://icanhazdadjoke.com
Vollständige URL: https://icanhazdadjoke.com/

Parameter: keine

Antworten:
  200: Ein zufälliger Dad Joke
    Header: Accept: text/plain
```

**LLM:** Der Endpunkt benötigt keine Parameter. Rufen Sie ihn einfach auf, und Sie erhalten einen zufälligen Witz.

---

## Eine API aufrufen

**Benutzer:** Gib mir einen zufälligen Dad Joke!

**LLM → MCP:** `invoke(endpointId: "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6")`

**MCP → LLM:**
```json
{
  "statusCode": 200,
  "body": "Why don't scientists trust atoms? Because they make up everything!"
}
```

**LLM:** Hier ist Ihr Dad Joke: "Why don't scientists trust atoms? Because they make up everything!"

---

## Mit großen Antworten arbeiten

**Benutzer:** Hol mir die Wettervorhersage für Berlin.

**LLM → MCP:** `search(query: "weather forecast", limit: 5)`

**MCP → LLM:**
```
Endpoints:
  GET /v1/forecast — Wettervorhersage
```

**LLM → MCP:** `inspect(endpointId: "e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7")`

**MCP → LLM:**
```
GET /v1/forecast
Parameter:
  latitude (query, number, erforderlich)
  longitude (query, number, erforderlich)
  hourly (query, string, optional)

Antworten:
  200: Vorhersagedaten
```

**LLM → MCP:** `invoke(endpointId: "e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7", parameters: {"latitude": 52.52, "longitude": 13.41, "hourly": "temperature_2m"})`

**MCP → LLM:**
```json
{
  "statusCode": 200,
  "fileRef": {
    "path": "/home/user/.swag2mcp/responses/response_x7k9m2.json",
    "size": 1572864,
    "sizeHint": "1.5 MB",
    "maxSizeHint": "2 KB",
    "message": "Die Antwort überschreitet das Limit von 2 KB und wurde auf der Festplatte gespeichert."
  }
}
```

**LLM:** Die Antwort ist groß (1,5 MB). Lassen Sie mich ihre Struktur erkunden.

**LLM → MCP:** `response_outline(path: "/home/user/.swag2mcp/responses/response_x7k9m2.json")`

**MCP → LLM:**
```
Typ: object
Größe: 1.5 MB
Schlüssel: latitude, longitude, hourly, daily
Arrays:
  hourly.time (168 Elemente)
  hourly.temperature_2m (168 Elemente)
```

**LLM → MCP:** `response_compress(path: "/home/user/.swag2mcp/responses/response_x7k9m2.json", mode: "sample_array", jsonPath: "hourly", arrayHead: 3, arrayTail: 2)`

**MCP → LLM:**
```json
{
  "compressed": {
    "type": "array",
    "length": 5,
    "original": 168,
    "skipped": 163,
    "sample": [
      {"time": "2026-07-24T00:00", "temperature_2m": 18.5},
      {"time": "2026-07-24T01:00", "temperature_2m": 17.8},
      {"time": "2026-07-24T02:00", "temperature_2m": 16.9},
      {"time": "2026-07-24T22:00", "temperature_2m": 20.1},
      {"time": "2026-07-24T23:00", "temperature_2m": 19.4}
    ]
  }
}
```

**LLM:** Hier ist die Wettervorhersage für Berlin. Die Temperaturen reichen von 16,9°C in der Nacht bis etwa 20°C am Abend.

---

## Auth-Token abrufen

**Benutzer:** Zeig mir das Auth-Token für meine API, damit ich es mit curl testen kann.

**LLM → MCP:** `auth(specId: "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6")`

**MCP → LLM:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**LLM:** Hier ist das Token. Sie können es in curl verwenden:

```bash
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." https://api.example.com/v1/users
```

---

## Laufzeitinfo

**Benutzer:** Wie ist der aktuelle Systemstatus?

**LLM → MCP:** `info()`

**MCP → LLM:**
```json
{
  "version": "v1.2.0",
  "workspace": "/home/user/.swag2mcp",
  "uptime": "2h 15m",
  "specs": {
    "total": 2,
    "active": 2,
    "endpoints": 20
  },
  "http_client": {
    "timeout": "30s",
    "max_response_size": "2 KB"
  },
  "mcp": {
    "transport": "stdio"
  }
}
```

**LLM:** Das System läuft seit 2 Stunden und 15 Minuten. Zwei Specs sind aktiv mit insgesamt 20 Endpunkten. Der HTTP-Client hat ein 30-Sekunden-Timeout und ein Antwortgrößenlimit von 2 KB.

---

## Vollständiges Gespräch: Reiseplanung

**Benutzer:** Ich plane nächste Woche eine Reise nach Tokio. Kannst du mir bei der Vorbereitung helfen?

**LLM → MCP:** `search(query: "weather forecast", limit: 5)`

**MCP → LLM:**
```
Endpoints:
  GET /v1/forecast — Wettervorhersage
```

**LLM → MCP:** `inspect(endpointId: "e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7")`

**MCP → LLM:**
```
GET /v1/forecast
Parameter:
  latitude (query, number, erforderlich)
  longitude (query, number, erforderlich)
  daily (query, string, optional)
```

**LLM → MCP:** `invoke(endpointId: "e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7", parameters: {"latitude": 35.68, "longitude": 139.69, "daily": "temperature_2m_max,temperature_2m_min,precipitation_sum"})`

**MCP → LLM:**
```json
{
  "statusCode": 200,
  "body": {
    "latitude": 35.68,
    "longitude": 139.69,
    "daily": {
      "time": ["2026-07-28", "2026-07-29", "2026-07-30", "2026-07-31", "2026-08-01"],
      "temperature_2m_max": [32, 33, 31, 30, 29],
      "temperature_2m_min": [25, 26, 24, 23, 22],
      "precipitation_sum": [0.0, 0.5, 2.1, 0.0, 0.0]
    }
  }
}
```

**LLM:** Hier ist die Wettervorhersage für Tokio für nächste Woche:

| Datum | Max | Min | Regen |
|------|-----|-----|-------|
| 28. Jul | 32°C | 25°C | 0 mm |
| 29. Jul | 33°C | 26°C | 0,5 mm |
| 30. Jul | 31°C | 24°C | 2,1 mm |
| 31. Jul | 30°C | 23°C | 0 mm |
| 1. Aug | 29°C | 22°C | 0 mm |

Sieht nach heißem und meist trockenem Wetter aus. Packen Sie leichte Kleidung und für den 30. Juli einen Regenschirm ein!

---

## Vollständiges Gespräch: Einen Witz finden

**Benutzer:** Erzähl mir etwas Lustiges!

**LLM → MCP:** `spec_list()`

**MCP → LLM:**
```
Specifications:
  dadjoke (icanhazdadjoke)
  meteo (Open-Meteo Weather API)
```

**LLM → MCP:** `search(query: "random joke", limit: 5)`

**MCP → LLM:**
```
Endpoints:
  GET /  — Einen zufälligen Dad Joke zurückgeben
```

**LLM → MCP:** `invoke(endpointId: "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6")`

**MCP → LLM:**
```json
{
  "statusCode": 200,
  "body": "I told my wife she was drawing her eyebrows too high. She looked surprised."
}
```

**LLM:** Hier: "I told my wife she was drawing her eyebrows too high. She looked surprised."
