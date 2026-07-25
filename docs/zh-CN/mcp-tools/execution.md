# 执行工具

执行工具是 swag2mcp 的核心：**search** 在你没有 ID 时查找端点，**inspect** 揭示完整的 OpenAPI 契约，**invoke** 执行实际的 API 调用。始终按此顺序使用它们：search → inspect → invoke。

---

## search

### 用途

唯一一个在你没有端点 ID 时查找端点的工具。使用 bluge 搜索引擎在所有 spec 的所有端点中执行全文搜索。

### 何时使用

- 当你不知道端点 ID 时
- 当你想按关键字、方法、标签或路径查找端点时
- 当你需要发现特定功能存在哪些端点时

### 工作原理

搜索所有 spec 的全文索引。支持带有字段过滤器、布尔运算符、模糊匹配、通配符等的结构化查询。

### 参数

| 参数 | 类型 | 必需 | 描述 |
|------|------|------|------|
| `query` | string | 是 | 搜索查询（支持结构化语法） |
| `limit` | int | 是 | 最大返回结果数（1-50） |

### 查询语法

| 示例 | 描述 |
|------|------|
| `pet` | 跨所有字段的简单文本搜索 |
| `method:GET` | 按 HTTP 方法过滤 |
| `tag:pet` | 按标签名称过滤 |
| `path:"/api/v1/users"` | 精确路径搜索 |
| `+method:POST +tag:pet` | 必须匹配两个条件 |
| `-method:DELETE` | 排除 DELETE 方法 |
| `create~` | 模糊搜索（容错） |
| `path:/api/v1/*` | 通配符路径搜索 |
| `/pattern/` | 正则表达式搜索 |
| `term^3` | 提升术语的相关性 |

**可搜索字段：** `method`（关键字）、`tag`（关键字）、`path`（文本）、`summary`（文本）、`_all`（默认文本字段）。

**不支持：** 括号分组、显式 `AND`/`OR` 运算符、字段分组。

### 响应

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
      "summary": "Get weather forecast for a location"
    }
  ]
}
```

每个结果包含完整谱系（spec → collection → tag），以便 LLM 可以导航到相关端点。

### 细节

- `limit` 必须在 1 到 50 之间（否则返回 `validation_failed`）
- `query` 是必需的（如果为空，返回 `validation_failed`）
- 结果按相关性顺序返回（最佳匹配优先）
- 使用字段过滤器（`method:GET`、`tag:pet`）缩小结果范围
- 对于精确路径匹配，使用引号：`path:"/v1/forecast"`

---

## inspect

### 用途

检索端点的完整 OpenAPI 操作对象：所有参数、请求体模式、响应模式、基础 URL 和完整 URL。这是在 `invoke` **之前**调用的工具，用于了解端点的契约。

### 何时使用

- 始终在 `invoke` 之前 — 你需要完整的契约才能正确调用
- 当你需要向用户解释 API 的技术细节时
- 当你需要了解必需参数、请求体结构或响应格式时

### 工作原理

在索引中查找端点并返回完整的 OpenAPI 操作对象，所有模式已解析。

### 参数

| 参数 | 类型 | 必需 | 描述 |
|------|------|------|------|
| `endpointId` | string | 是 | 端点的 32 字符 MD5 哈希 |

### 响应

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
    "summary": "Add a new pet",
    "description": "Add a new pet to the store",
    "deprecated": false,
    "parameters": [
      {
        "name": "petId",
        "in": "path",
        "description": "ID of the pet",
        "required": true,
        "schema": {
          "type": "integer",
          "format": "int64"
        }
      }
    ],
    "requestBody": {
      "description": "Pet object to add",
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
        "description": "Successful operation",
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/Pet"
            }
          }
        }
      },
      "405": {
        "description": "Invalid input"
      }
    }
  }
}
```

| 字段 | 类型 | 描述 |
|------|------|------|
| `baseUrl` | string | API 的基础 URL（来自配置） |
| `fullUrl` | string | 端点的完整 URL（base + path） |
| `operation.parameters[]` | array | 参数，包含名称、位置（path/query/header/cookie）、描述、必需标志和模式 |
| `operation.requestBody` | object | 请求体，包含内容类型和模式 |
| `operation.responses` | map | 响应码，包含描述和模式 |
| `operation.deprecated` | bool | 端点是否已弃用 |

### 细节

- 如果端点不存在，返回 `not_found`
- 这是**唯一**返回完整 OpenAPI 操作的工具 — `endpoint_by_id` 只返回摘要
- 始终在 `invoke` 之前调用 `inspect` 以了解必需参数和主体结构
- `operation` 对象包含已解析为其完整模式定义的 `$ref` 引用

---

## invoke

### 用途

对端点执行真实的 API 调用。这是唯一发出实际 HTTP 请求的工具。认证自动应用 — 你不需要先调用 `auth`。

### 何时使用

- 仅在调用 `inspect` 了解端点的契约之后
- 仅在对破坏性操作（POST、PUT、PATCH、DELETE）有明确的用户确认时
- 当用户要求调用 API 并且你有所有必需参数时

### 工作原理

1. 在索引中查找端点
2. 将路径参数替换到 URL 中
3. 追加查询参数
4. 添加头和 cookie
5. 将请求体序列化为 JSON
6. 自动获取并应用认证（令牌、头、查询参数）
7. 发起 HTTP 请求
8. 返回响应，如果太大则保存到文件

### 参数

| 参数 | 类型 | 必需 | 描述 |
|------|------|------|------|
| `endpointId` | string | 是 | 端点的 32 字符 MD5 哈希 |
| `parameters` | object | 否 | 路径、查询和头参数，键值对 |
| `requestBody` | object | 否 | POST/PUT/PATCH 请求的请求体 |
| `headers` | object | 否 | 要发送的额外 HTTP 头 |
| `cookies` | object | 否 | 要发送的额外 HTTP cookie |

### 响应（内联）

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

### 响应（文件引用 — 当主体超过大小限制时）

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
    "message": "Response exceeds the 2 KB limit and has been saved to disk.",
    "openCmd": "open /Users/user/.swag2mcp/responses/response_a1b2c3d4.json"
  }
}
```

| 字段 | 类型 | 描述 |
|------|------|------|
| `statusCode` | int | HTTP 响应状态码 |
| `headers` | object | HTTP 响应头 |
| `body` | any | 响应体（在大小限制内时存在） |
| `fileRef` | object | 文件引用（当主体超过大小限制时存在） |

### 处理大响应

当 `invoke` 返回 `fileRef` 时，使用响应工具探索数据：

1. **`response_outline(path)`** — 获取结构摘要（键、类型、数组长度）
2. **`response_compress(path, mode)`** — 压缩数据以适应内联
3. **`response_slice(path, jsonPath)`** — 提取特定片段

### 细节

- **认证是自动的：** `invoke` 工具自动从 spec 的认证配置中获取并应用认证。你**不需要**先调用 `auth`。
- **速率限制：** 每个端点有 10 秒冷却时间。10 秒内对同一端点的第二次调用被静默阻止（返回 `rate_limit` 错误）。
- **响应大小限制：** 默认 2 KB（可通过 `max_response_size` 配置）。如果响应超过此限制，保存到 `{workspace}/responses/` 并返回 `FileReference` 而不是内联 `body`。
- **参数处理：** 路径参数替换到 URL 中。查询参数被追加。请求中的参数覆盖操作规范的默认值。
- **请求体：** 对于 POST/PUT/PATCH，主体序列化为 JSON。`Content-Type` 自动设置为 `application/json`。
- **错误处理：** HTTP 错误（非 2xx）作为 `invoke_error` 返回，状态码和响应体在提示中。
- **破坏性操作：** 在没有明确用户确认的情况下，永远不要调用 POST/PUT/PATCH/DELETE。
