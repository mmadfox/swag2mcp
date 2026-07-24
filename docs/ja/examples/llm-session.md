# LLM セッション例

このページでは、swag2mcp を利用した LLM エージェントとユーザーの実際の対話を示します。各例には、ユーザーのリクエスト、MCP ツール呼び出し、およびレスポンスが含まれています。

## セットアップ

ワークスペースには 2 つの API が設定されています：

- **icanhazdadjoke** — 公開ダジャレ API
- **Open-Meteo** — 天気予報 API

---

## 利用可能な API の発見

**ユーザー：** どの API にアクセスできますか？

**LLM → MCP：** `spec_list()`

**MCP → LLM：**
```
Specifications:
  dadjoke (icanhazdadjoke)
  meteo (Open-Meteo Weather API)
```

**LLM：** 2 つの API にアクセスできます：ダジャレ用の icanhazdadjoke と天気データ用の Open-Meteo です。

---

## スペック詳細の取得

**ユーザー：** dadjoke API について詳しく教えてください。

**LLM → MCP：** `spec_by_id(id: "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6")`

**MCP → LLM：**
```
Spec: dadjoke (icanhazdadjoke)
Collections:
  - Jokes (3 endpoints, 1 tag)
```

**LLM：** dadjoke API には「Jokes」という 1 つのコレクションがあり、1 つのタグの下に 3 つのエンドポイントがあります。

---

## エンドポイントの検索

**ユーザー：** ランダムなダジャレを取得するエンドポイントを探してください。

**LLM → MCP：** `search(query: "random joke", limit: 5)`

**MCP → LLM：**
```
Endpoints:
  GET /  — Return a random dad joke
  GET /j/{id} — Return a specific dad joke by ID
  GET /search — Search for dad jokes
```

**LLM：** 3 つのエンドポイントが見つかりました。メインは `GET /` で、ランダムなダジャレを返します。

---

## エンドポイントの検査

**ユーザー：** ランダムジョークエンドポイントにはどのようなパラメータが必要ですか？

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

**LLM：** このエンドポイントはパラメータを必要としません。呼び出すだけでランダムなジョークが得られます。

---

## API の呼び出し

**ユーザー：** ランダムなダジャレを教えて！

**LLM → MCP：** `invoke(endpointId: "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6")`

**MCP → LLM：**
```json
{
  "statusCode": 200,
  "body": "Why don't scientists trust atoms? Because they make up everything!"
}
```

**LLM：** こちらがダジャレです：「Why don't scientists trust atoms? Because they make up everything!」

---

## 大規模レスポンスの操作

**ユーザー：** ベルリンの天気予報を取得してください。

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

**LLM：** レスポンスが大きいです（1.5 MB）。構造を調べてみましょう。

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

**LLM：** こちらがベルリンの天気予報です。気温は夜の 16.9°C から夕方の約 20°C の範囲です。

---

## 認証トークンの取得

**ユーザー：** curl でテストできるように、API の認証トークンを表示してください。

**LLM → MCP：** `auth(specId: "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6")`

**MCP → LLM：**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**LLM：** こちらがトークンです。curl で使用できます：

```bash
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." https://api.example.com/v1/users
```

---

## ランタイム情報

**ユーザー：** 現在のシステムステータスは？

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

**LLM：** システムは 2 時間 15 分間稼働しています。2 つのスペックがアクティブで、合計 20 のエンドポイントがあります。HTTP クライアントのタイムアウトは 30 秒、レスポンスサイズ制限は 2 KB です。

---

## 完全な会話：旅行の計画

**ユーザー：** 来週東京に旅行に行く予定です。準備を手伝ってくれますか？

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

**LLM：** こちらが来週の東京の天気予報です：

| 日付 | 最高 | 最低 | 雨 |
|------|-----|-----|------|
| 7/28 | 32°C | 25°C | 0 mm |
| 7/29 | 33°C | 26°C | 0.5 mm |
| 7/30 | 31°C | 24°C | 2.1 mm |
| 7/31 | 30°C | 23°C | 0 mm |
| 8/1 | 29°C | 22°C | 0 mm |

暑くてほぼ乾燥した天気のようです。軽装で、7月30日に備えて傘も持っていきましょう！

---

## 完全な会話：ジョークを探す

**ユーザー：** 何か面白いことを教えて！

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

**LLM：** どうぞ：「I told my wife she was drawing her eyebrows too high. She looked surprised.」
