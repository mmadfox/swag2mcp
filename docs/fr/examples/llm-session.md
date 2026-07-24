# LLM Session Examples

This page shows real dialogues between a user and an LLM agent powered by swag2mcp. Each example includes the user's request, the MCP tool call, and the response.

## Setup

The workspace has two APIs configured:

- **icanhazdadjoke** — public dad jokes API
- **Open-Meteo** — weather forecast API

---

## Discovering available APIs

**User:** What APIs do you have access to?

**LLM → MCP:** `spec_list()`

**MCP → LLM:**
```
Specifications:
  dadjoke (icanhazdadjoke)
  meteo (Open-Meteo Weather API)
```

**LLM:** I have access to two APIs: icanhazdadjoke for dad jokes and Open-Meteo for weather data.

---

## Getting spec details

**User:** Tell me more about the dadjoke API.

**LLM → MCP:** `spec_by_id(id: "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6")`

**MCP → LLM:**
```
Spec: dadjoke (icanhazdadjoke)
Collections:
  - Jokes (3 endpoints, 1 tag)
```

**LLM:** The dadjoke API has one collection called "Jokes" with 3 endpoints under a single tag.

---

## Searching for endpoints

**User:** Find me an endpoint to get a random dad joke.

**LLM → MCP:** `search(query: "random joke", limit: 5)`

**MCP → LLM:**
```
Endpoints:
  GET /  — Return a random dad joke
  GET /j/{id} — Return a specific dad joke by ID
  GET /search — Search for dad jokes
```

**LLM:** I found three endpoints. The main one is `GET /` which returns a random dad joke.

---

## Inspecting an endpoint

**User:** What parameters does the random joke endpoint need?

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

**LLM:** The endpoint takes no parameters. Just call it and you'll get a random joke.

---

## Invoking an API

**User:** Give me a random dad joke!

**LLM → MCP:** `invoke(endpointId: "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6")`

**MCP → LLM:**
```json
{
  "statusCode": 200,
  "body": "Why don't scientists trust atoms? Because they make up everything!"
}
```

**LLM:** Here's your dad joke: "Why don't scientists trust atoms? Because they make up everything!"

---

## Working with large responses

**User:** Get me the weather forecast for Berlin.

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
    "message": "Response exceeds the 2 KB limit and has been saved to disk."
  }
}
```

**LLM:** The response is large (1.5 MB). Let me explore its structure.

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

**LLM:** Here's the Berlin weather forecast. Temperatures range from 16.9°C at night to around 20°C in the evening.

---

## Auth token retrieval

**User:** Show me the auth token for my API so I can test it with curl.

**LLM → MCP:** `auth(specId: "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6")`

**MCP → LLM:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**LLM:** Here's the token. You can use it in curl:

```bash
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." https://api.example.com/v1/users
```

---

## Runtime info

**User:** What's the current system status?

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

**LLM:** The system has been running for 2 hours and 15 minutes. Two specs are active with 20 endpoints total. The HTTP client has a 30-second timeout and a 2 KB response size limit.

---

## Full conversation: planning a trip

**User:** I'm planning a trip to Tokyo next week. Can you help me prepare?

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

**LLM:** Here's the Tokyo weather forecast for next week:

| Date | Max | Min | Rain |
|------|-----|-----|------|
| Jul 28 | 32°C | 25°C | 0 mm |
| Jul 29 | 33°C | 26°C | 0.5 mm |
| Jul 30 | 31°C | 24°C | 2.1 mm |
| Jul 31 | 30°C | 23°C | 0 mm |
| Aug 1 | 29°C | 22°C | 0 mm |

Looks like hot and mostly dry weather. Pack light clothes and an umbrella just in case for July 30!

---

## Full conversation: finding a joke

**User:** Tell me something funny!

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

**LLM:** Here you go: "I told my wife she was drawing her eyebrows too high. She looked surprised."
