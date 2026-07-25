# Ejemplos de Sesión LLM

Esta página muestra diálogos reales entre un usuario y un agente LLM impulsado por swag2mcp. Cada ejemplo incluye la solicitud del usuario, la llamada a la herramienta MCP y la respuesta.

## Configuración

El espacio de trabajo tiene dos APIs configuradas:

- **icanhazdadjoke** — API pública de chistes de papá
- **Open-Meteo** — API de pronóstico meteorológico

---

## Descubriendo APIs disponibles

**Usuario:** ¿Qué APIs tienes disponibles?

**LLM → MCP:** `spec_list()`

**MCP → LLM:**
```
Specifications:
  dadjoke (icanhazdadjoke)
  meteo (Open-Meteo Weather API)
```

**LLM:** Tengo acceso a dos APIs: icanhazdadjoke para chistes de papá y Open-Meteo para datos meteorológicos.

---

## Obteniendo detalles de la especificación

**Usuario:** Cuéntame más sobre la API de dadjoke.

**LLM → MCP:** `spec_by_id(id: "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6")`

**MCP → LLM:**
```
Spec: dadjoke (icanhazdadjoke)
Collections:
  - Jokes (3 endpoints, 1 tag)
```

**LLM:** La API de dadjoke tiene una colección llamada "Jokes" con 3 endpoints bajo una sola etiqueta.

---

## Buscando endpoints

**Usuario:** Encuéntrame un endpoint para obtener un chiste de papá aleatorio.

**LLM → MCP:** `search(query: "random joke", limit: 5)`

**MCP → LLM:**
```
Endpoints:
  GET /  — Return a random dad joke
  GET /j/{id} — Return a specific dad joke by ID
  GET /search — Search for dad jokes
```

**LLM:** Encontré tres endpoints. El principal es `GET /` que devuelve un chiste de papá aleatorio.

---

## Inspeccionando un endpoint

**Usuario:** ¿Qué parámetros necesita el endpoint de chiste aleatorio?

**LLM → MCP:** `inspect(endpointId: "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6")`

**MCP → LLM:**
```
GET /
Base URL: https://icanhazdadjoke.com
Full URL: https://icanhazdadjoke.com/

Parameters: none

Responses:
  200: A random dad joke
    Headers: Accept: text/plain
```

**LLM:** El endpoint no requiere parámetros. Solo llámalo y obtendrás un chiste aleatorio.

---

## Invocando una API

**Usuario:** ¡Dame un chiste de papá aleatorio!

**LLM → MCP:** `invoke(endpointId: "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6")`

**MCP → LLM:**
```json
{
  "statusCode": 200,
  "body": "Why don't scientists trust atoms? Because they make up everything!"
}
```

**LLM:** Aquí tienes tu chiste de papá: "Why don't scientists trust atoms? Because they make up everything!"

---

## Trabajando con respuestas grandes

**Usuario:** Dame el pronóstico meteorológico para Berlín.

**LLM → MCP:** `search(query: "weather forecast", limit: 5)`

**MCP → LLM:**
```
Endpoints:
  GET /v1/forecast — Weather forecast
```

**LLM → MCP:** `inspect(endpointId: "e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7")`

**MCP → LLM:**
```
GET /v1/forecast
Parameters:
  latitude (query, number, required)
  longitude (query, number, required)
  hourly (query, string, optional)

Responses:
  200: Forecast data
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
    "message": "La respuesta excede el límite de 2 KB y se ha guardado en disco."
  }
}
```

**LLM:** La respuesta es grande (1.5 MB). Déjame explorar su estructura.

**LLM → MCP:** `response_outline(path: "/home/user/.swag2mcp/responses/response_x7k9m2.json")`

**MCP → LLM:**
```
Type: object
Size: 1.5 MB
Keys: latitude, longitude, hourly, daily
Arrays:
  hourly.time (168 items)
  hourly.temperature_2m (168 items)
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

**LLM:** Aquí está el pronóstico meteorológico para Berlín. Las temperaturas oscilan entre 16.9°C por la noche y alrededor de 20°C por la tarde.

---

## Recuperación de token de autenticación

**Usuario:** Muéstrame el token de autenticación para mi API para que pueda probarlo con curl.

**LLM → MCP:** `auth(specId: "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6")`

**MCP → LLM:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**LLM:** Aquí está el token. Puedes usarlo en curl:

```bash
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." https://api.example.com/v1/users
```

---

## Información de ejecución

**Usuario:** ¿Cuál es el estado actual del sistema?

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

**LLM:** El sistema ha estado funcionando durante 2 horas y 15 minutos. Dos especificaciones están activas con 20 endpoints en total. El cliente HTTP tiene un tiempo de espera de 30 segundos y un límite de tamaño de respuesta de 2 KB.

---

## Conversación completa: planificando un viaje

**Usuario:** Estoy planeando un viaje a Tokio la próxima semana. ¿Puedes ayudarme a prepararme?

**LLM → MCP:** `search(query: "weather forecast", limit: 5)`

**MCP → LLM:**
```
Endpoints:
  GET /v1/forecast — Weather forecast
```

**LLM → MCP:** `inspect(endpointId: "e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7")`

**MCP → LLM:**
```
GET /v1/forecast
Parameters:
  latitude (query, number, required)
  longitude (query, number, required)
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

**LLM:** Aquí está el pronóstico meteorológico para Tokio la próxima semana:

| Fecha | Máx | Mín | Lluvia |
|------|-----|-----|--------|
| Jul 28 | 32°C | 25°C | 0 mm |
| Jul 29 | 33°C | 26°C | 0.5 mm |
| Jul 30 | 31°C | 24°C | 2.1 mm |
| Jul 31 | 30°C | 23°C | 0 mm |
| Ago 1 | 29°C | 22°C | 0 mm |

¡Parece un clima mayormente seco y caluroso! Empaca ropa ligera y un paraguas por si acaso para el 30 de julio.

---

## Conversación completa: encontrando un chiste

**Usuario:** ¡Dime algo gracioso!

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
  GET /  — Return a random dad joke
```

**LLM → MCP:** `invoke(endpointId: "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6")`

**MCP → LLM:**
```json
{
  "statusCode": 200,
  "body": "I told my wife she was drawing her eyebrows too high. She looked surprised."
}
```

**LLM:** Aquí tienes: "I told my wife she was drawing her eyebrows too high. She looked surprised."
