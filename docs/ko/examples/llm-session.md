# LLM 세션 예시

이 페이지는 swag2mcp로 구동되는 LLM 에이전트와 사용자 간의 실제 대화를 보여줍니다. 각 예시에는 사용자 요청, MCP 도구 호출, 응답이 포함됩니다.

## 설정

워크스페이스에는 두 개의 API가 설정되어 있습니다:

- **icanhazdadjoke** — 공개 아재개그 API
- **Open-Meteo** — 날씨 예보 API

---

## 사용 가능한 API 발견

**사용자:** 어떤 API에 접근할 수 있나요?

**LLM → MCP:** `spec_list()`

**MCP → LLM:**
```
Specifications:
  dadjoke (icanhazdadjoke)
  meteo (Open-Meteo Weather API)
```

**LLM:** 두 개의 API에 접근할 수 있습니다: 아재개그용 icanhazdadjoke와 날씨 데이터용 Open-Meteo입니다.

---

## Spec 세부 정보 보기

**사용자:** dadjoke API에 대해 더 알려주세요.

**LLM → MCP:** `spec_by_id(id: "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6")`

**MCP → LLM:**
```
Spec: dadjoke (icanhazdadjoke)
Collections:
  - Jokes (3 endpoints, 1 tag)
```

**LLM:** dadjoke API에는 "Jokes"라는 하나의 collection이 있으며, 단일 태그 아래에 3개의 엔드포인트가 있습니다.

---

## 엔드포인트 검색

**사용자:** 무작위 아재개그를 가져오는 엔드포인트를 찾아주세요.

**LLM → MCP:** `search(query: "random joke", limit: 5)`

**MCP → LLM:**
```
Endpoints:
  GET /  — 무작위 아재개그 반환
  GET /j/{id} — ID로 특정 아재개그 반환
  GET /search — 아재개그 검색
```

**LLM:** 세 개의 엔드포인트를 찾았습니다. 주요 엔드포인트는 `GET /`이며 무작위 아재개그를 반환합니다.

---

## 엔드포인트 검사

**사용자:** 무작위 농담 엔드포인트에 어떤 매개변수가 필요하나요?

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

**LLM:** 엔드포인트는 매개변수가 필요하지 않습니다. 그냥 호출하면 무작위 농담을 받을 수 있습니다.

---

## API 호출

**사용자:** 아재개그를 들려주세요!

**LLM → MCP:** `invoke(endpointId: "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6")`

**MCP → LLM:**
```json
{
  "statusCode": 200,
  "body": "과학자들이 원자를 믿지 않는 이유는 무엇일까요? 그들이 모든 것을 만들어내기 때문입니다!"
}
```

**LLM:** 아재개그입니다: "과학자들이 원자를 믿지 않는 이유는 무엇일까요? 그들이 모든 것을 만들어내기 때문입니다!"

---

## 큰 응답 작업

**사용자:** 베를린의 날씨 예보를 가져와주세요.

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
    "message": "응답이 2 KB 제한을 초과하여 디스크에 저장되었습니다."
  }
}
```

**LLM:** 응답이 큽니다(1.5 MB). 구조를 살펴보겠습니다.

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

**LLM:** 베를린 날씨 예보입니다. 기온은 야간 16.9°C에서 저녁 약 20°C까지입니다.

---

## 인증 토큰 검색

**사용자:** curl로 테스트할 수 있도록 내 API의 인증 토큰을 보여주세요.

**LLM → MCP:** `auth(specId: "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6")`

**MCP → LLM:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**LLM:** 토큰입니다. curl에서 사용할 수 있습니다:

```bash
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." https://api.example.com/v1/users
```

---

## 런타임 정보

**사용자:** 현재 시스템 상태는 무엇인가요?

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

**LLM:** 시스템이 2시간 15분 동안 실행되었습니다. 두 개의 spec이 활성화되어 있으며 총 20개의 엔드포인트가 있습니다. HTTP 클라이언트는 30초 타임아웃과 2 KB 응답 크기 제한이 있습니다.

---

## 전체 대화: 여행 계획

**사용자:** 다음 주에 도쿄 여행을 계획 중입니다. 준비를 도와주실 수 있나요?

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

**LLM:** 다음 주 도쿄 날씨 예보입니다:

| 날짜 | 최고 | 최저 | 강수량 |
|------|------|------|--------|
| 7/28 | 32°C | 25°C | 0 mm |
| 7/29 | 33°C | 26°C | 0.5 mm |
| 7/30 | 31°C | 24°C | 2.1 mm |
| 7/31 | 30°C | 23°C | 0 mm |
| 8/1 | 29°C | 22°C | 0 mm |

덥고 대체로 건조한 날씨입니다. 가벼운 옷을 챙기고 7월 30일을 대비해 우산도 준비하세요!

---

## 전체 대화: 농담 찾기

**사용자:** 재미있는 이야기 해주세요!

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
  "body": "아내한테 눈썹을 너무 높게 그리고 있다고 말했어요. 그녀가 놀란 표정을 지었습니다."
}
```

**LLM:** 여기 있습니다: "아내한테 눈썹을 너무 높게 그리고 있다고 말했어요. 그녀가 놀란 표정을 지었습니다."
