# LLM 会话示例

此页面展示了用户与由 swag2mcp 驱动的 LLM 智能体之间的真实对话。每个示例包括用户的请求、MCP 工具调用和响应。

## 设置

工作区配置了两个 API：

- **icanhazdadjoke** — 公共冷笑话 API
- **Open-Meteo** — 天气预报 API

---

## 发现可用的 API

**用户：** 你可以访问哪些 API？

**LLM → MCP：** `spec_list()`

**MCP → LLM：**
```
Specifications:
  dadjoke (icanhazdadjoke)
  meteo (Open-Meteo Weather API)
```

**LLM：** 我可以访问两个 API：icanhazdadjoke 用于冷笑话，Open-Meteo 用于天气数据。

---

## 获取 spec 详情

**用户：** 告诉我更多关于 dadjoke API 的信息。

**LLM → MCP：** `spec_by_id(id: "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6")`

**MCP → LLM：**
```
Spec: dadjoke (icanhazdadjoke)
Collections:
  - Jokes (3 endpoints, 1 tag)
```

**LLM：** dadjoke API 有一个名为"Jokes"的 collection，包含 3 个端点，属于一个标签。

---

## 搜索端点

**用户：** 帮我找一个获取随机冷笑话的端点。

**LLM → MCP：** `search(query: "random joke", limit: 5)`

**MCP → LLM：**
```
Endpoints:
  GET /  — Return a random dad joke
  GET /j/{id} — Return a specific dad joke by ID
  GET /search — Search for dad jokes
```

**LLM：** 我找到了三个端点。主要的是 `GET /`，它返回一个随机冷笑话。

---

## 检查端点

**用户：** 随机笑话端点需要什么参数？

**LLM → MCP：** `inspect(endpointId: "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6")`

**MCP → LLM：**
```
GET /
Base URL: https://icanhazdadjoke.com
Full URL: https://icanhazdadjoke.com/

Parameters: none

Responses:
  200: A random dad joke
    Headers: Accept: text/plain
```

**LLM：** 该端点不需要任何参数。直接调用它，你就会得到一个随机笑话。

---

## 调用 API

**用户：** 给我讲个冷笑话！

**LLM → MCP：** `invoke(endpointId: "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6")`

**MCP → LLM：**
```json
{
  "statusCode": 200,
  "body": "Why don't scientists trust atoms? Because they make up everything!"
}
```

**LLM：** 这是你的冷笑话："Why don't scientists trust atoms? Because they make up everything!"

---

## 处理大响应

**用户：** 帮我获取柏林的天气预报。

**LLM → MCP：** `search(query: "weather forecast", limit: 5)`

**MCP → LLM：**
```
Endpoints:
  GET /v1/forecast — Weather forecast
```

**LLM → MCP：** `inspect(endpointId: "e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7")`

**MCP → LLM：**
```
GET /v1/forecast
Parameters:
  latitude (query, number, required)
  longitude (query, number, required)
  hourly (query, string, optional)

Responses:
  200: Forecast data
```

**LLM → MCP：** `invoke(endpointId: "e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7", parameters: {"latitude": 52.52, "longitude": 13.41, "hourly": "temperature_2m"})`

**MCP → LLM：**
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

**LLM：** 响应很大（1.5 MB）。让我探索一下它的结构。

**LLM → MCP：** `response_outline(path: "/home/user/.swag2mcp/responses/response_x7k9m2.json")`

**MCP → LLM：**
```
Type: object
Size: 1.5 MB
Keys: latitude, longitude, hourly, daily
Arrays:
  hourly.time (168 items)
  hourly.temperature_2m (168 items)
```

**LLM → MCP：** `response_compress(path: "/home/user/.swag2mcp/responses/response_x7k9m2.json", mode: "sample_array", jsonPath: "hourly", arrayHead: 3, arrayTail: 2)`

**MCP → LLM：**
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

**LLM：** 这是柏林的天气预报。温度范围从夜间 16.9°C 到晚上约 20°C。

---

## 认证令牌检索

**用户：** 显示我的 API 的认证令牌，这样我可以用 curl 测试它。

**LLM → MCP：** `auth(specId: "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6")`

**MCP → LLM：**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**LLM：** 这是令牌。你可以在 curl 中使用它：

```bash
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." https://api.example.com/v1/users
```

---

## 运行时信息

**用户：** 当前系统状态如何？

**LLM → MCP：** `info()`

**MCP → LLM：**
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

**LLM：** 系统已运行 2 小时 15 分钟。两个 spec 处于活动状态，共有 20 个端点。HTTP 客户端有 30 秒超时和 2 KB 响应大小限制。

---

## 完整对话：计划一次旅行

**用户：** 我计划下周去东京旅行。你能帮我准备吗？

**LLM → MCP：** `search(query: "weather forecast", limit: 5)`

**MCP → LLM：**
```
Endpoints:
  GET /v1/forecast — Weather forecast
```

**LLM → MCP：** `inspect(endpointId: "e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7")`

**MCP → LLM：**
```
GET /v1/forecast
Parameters:
  latitude (query, number, required)
  longitude (query, number, required)
  daily (query, string, optional)
```

**LLM → MCP：** `invoke(endpointId: "e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7", parameters: {"latitude": 35.68, "longitude": 139.69, "daily": "temperature_2m_max,temperature_2m_min,precipitation_sum"})`

**MCP → LLM：**
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

**LLM：** 这是下周东京的天气预报：

| 日期 | 最高 | 最低 | 降雨 |
|------|------|------|------|
| 7月28日 | 32°C | 25°C | 0 mm |
| 7月29日 | 33°C | 26°C | 0.5 mm |
| 7月30日 | 31°C | 24°C | 2.1 mm |
| 7月31日 | 30°C | 23°C | 0 mm |
| 8月1日 | 29°C | 22°C | 0 mm |

看起来天气炎热，大部分时间干燥。带轻便衣服，7月30日以防万一带把伞！

---

## 完整对话：找个笑话

**用户：** 给我讲个有趣的事！

**LLM → MCP：** `spec_list()`

**MCP → LLM：**
```
Specifications:
  dadjoke (icanhazdadjoke)
  meteo (Open-Meteo Weather API)
```

**LLM → MCP：** `search(query: "random joke", limit: 5)`

**MCP → LLM：**
```
Endpoints:
  GET /  — Return a random dad joke
```

**LLM → MCP：** `invoke(endpointId: "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6")`

**MCP → LLM：**
```json
{
  "statusCode": 200,
  "body": "I told my wife she was drawing her eyebrows too high. She looked surprised."
}
```

**LLM：** 给你："I told my wife she was drawing her eyebrows too high. She looked surprised."
