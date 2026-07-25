# Примеры LLM-сессий

На этой странице показаны реальные диалоги между пользователем и LLM-агентом на базе swag2mcp. Каждый пример включает запрос пользователя, вызов MCP-инструмента и ответ.

## Настройка

В рабочей области настроены два API:

- **icanhazdadjoke** — публичное API шуток про пап
- **Open-Meteo** — API прогноза погоды

---

## Обнаружение доступных API

**Пользователь:** К каким API у тебя есть доступ?

**LLM → MCP:** `spec_list()`

**MCP → LLM:**
```
Specifications:
  dadjoke (icanhazdadjoke)
  meteo (Open-Meteo Weather API)
```

**LLM:** У меня есть доступ к двум API: icanhazdadjoke для шуток про пап и Open-Meteo для данных о погоде.

---

## Получение деталей спецификации

**Пользователь:** Расскажи подробнее об API dadjoke.

**LLM → MCP:** `spec_by_id(id: "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6")`

**MCP → LLM:**
```
Spec: dadjoke (icanhazdadjoke)
Collections:
  - Jokes (3 endpoints, 1 tag)
```

**LLM:** API dadjoke имеет одну коллекцию "Jokes" с 3 эндпоинтами в одном теге.

---

## Поиск эндпоинтов

**Пользователь:** Найди эндпоинт для получения случайной шутки про пап.

**LLM → MCP:** `search(query: "random joke", limit: 5)`

**MCP → LLM:**
```
Endpoints:
  GET /  — Return a random dad joke
  GET /j/{id} — Return a specific dad joke by ID
  GET /search — Search for dad jokes
```

**LLM:** Я нашёл три эндпоинта. Основной — `GET /`, который возвращает случайную шутку про пап.

---

## Проверка эндпоинта

**Пользователь:** Какие параметры нужны эндпоинту случайной шутки?

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

**LLM:** Эндпоинт не требует параметров. Просто вызовите его, и вы получите случайную шутку.

---

## Вызов API

**Пользователь:** Расскажи случайную шутку про пап!

**LLM → MCP:** `invoke(endpointId: "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6")`

**MCP → LLM:**
```json
{
  "statusCode": 200,
  "body": "Why don't scientists trust atoms? Because they make up everything!"
}
```

**LLM:** Вот ваша шутка: "Why don't scientists trust atoms? Because they make up everything!"

---

## Работа с большими ответами

**Пользователь:** Получи прогноз погоды для Берлина.

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

**LLM:** Ответ большой (1.5 МБ). Давайте исследуем его структуру.

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

**LLM:** Вот прогноз погоды для Берлина. Температура колеблется от 16.9°C ночью до около 20°C вечером.

---

## Получение токена аутентификации

**Пользователь:** Покажи токен аутентификации для моего API, чтобы я мог протестировать его через curl.

**LLM → MCP:** `auth(specId: "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6")`

**MCP → LLM:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**LLM:** Вот токен. Вы можете использовать его в curl:

```bash
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." https://api.example.com/v1/users
```

---

## Информация о рантайме

**Пользователь:** Каков текущий статус системы?

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

**LLM:** Система работает 2 часа 15 минут. Активны две спецификации с 20 эндпоинтами. HTTP-клиент имеет таймаут 30 секунд и лимит размера ответа 2 КБ.

---

## Полный разговор: планирование поездки

**Пользователь:** Я планирую поездку в Токио на следующей неделе. Можешь помочь подготовиться?

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

**LLM:** Вот прогноз погоды для Токио на следующую неделю:

| Дата | Макс | Мин | Дождь |
|------|------|-----|-------|
| 28 июл | 32°C | 25°C | 0 мм |
| 29 июл | 33°C | 26°C | 0.5 мм |
| 30 июл | 31°C | 24°C | 2.1 мм |
| 31 июл | 30°C | 23°C | 0 мм |
| 1 авг | 29°C | 22°C | 0 мм |

Похоже на жаркую и в основном сухую погоду. Берите лёгкую одежду и зонтик на всякий случай на 30 июля!

---

## Полный разговор: поиск шутки

**Пользователь:** Расскажи что-нибудь смешное!

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

**LLM:** Держите: "I told my wife she was drawing her eyebrows too high. She looked surprised."
