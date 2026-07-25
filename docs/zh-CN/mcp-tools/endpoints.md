# 端点工具

端点工具让 LLM 在层次结构的不同级别查看 API 端点：spec 中的所有端点、collection 中的端点、标签中的端点或单个端点摘要。在检查或调用之前，使用这些工具发现可用的操作。

---

## endpoint_by_spec

### 用途

列出整个 spec 中的所有端点，跨越所有 collection 和标签。返回最全面的视图 — spec 中的每个端点及其完整上下文（标签、collection、spec）。

### 何时使用

- 当你想查看 spec 中可用的每个端点时
- 当你不知道哪个 collection 或标签包含你需要的端点时
- 在 `spec_by_id` 之后获取完整的端点列表

### 参数

| 参数 | 类型 | 必需 | 描述 |
|------|------|------|------|
| `specId` | string | 是 | spec 的 32 字符 MD5 哈希 |

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

| 字段 | 类型 | 描述 |
|------|------|------|
| `id` | string | 端点标识符 |
| `tagId` | string | 父标签标识符 |
| `tagName` | string | 人类可读的标签名称 |
| `collectionId` | string | 父 collection 标识符 |
| `collectionTitle` | string | 人类可读的 collection 标题 |
| `specId` | string | 父 spec 标识符 |
| `specDomain` | string | Spec 域名 |
| `method` | string | HTTP 方法（GET、POST、PUT、DELETE 等） |
| `path` | string | API 路径（例如 /v1/forecast） |
| `summary` | string | 端点功能的人类可读摘要 |

### 细节

- 如果 spec 不存在，返回 `not_found`
- 每个端点包含其完整谱系（spec → collection → tag）以提供上下文
- 要快速查看单个端点的摘要，使用 `endpoint_by_id`

---

## endpoint_by_collection

### 用途

列出特定 collection 中的所有端点，无论其标签如何。返回按 collection 分组的端点，包含 spec 和 collection 元数据。

### 何时使用

- 在 `collection_by_id` 之后查看 collection 中的所有端点
- 当你想探索 collection 的完整 API 表面时

### 参数

| 参数 | 类型 | 必需 | 描述 |
|------|------|------|------|
| `collectionId` | string | 是 | collection 的 32 字符 MD5 哈希 |

### 响应

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collection": {
    "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
    "title": "Weather Forecast",
    "countMethods": 12
  },
  "endpoints": [
    {
      "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
      "tagId": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "tagName": "forecast",
      "method": "GET",
      "path": "/v1/forecast",
      "summary": "Get weather forecast for a location"
    }
  ]
}
```

### 细节

- 如果 collection 不存在，返回 `not_found`
- 包含 spec 和 collection 元数据以提供上下文
- collection 中所有标签的端点一起返回

---

## endpoint_by_tag

### 用途

列出特定标签下分组的所有端点。这是最集中的视图 — 一个 collection 中一个标签内的端点。

### 何时使用

- 在 `tag_by_id` 之后查看标签中的实际端点
- 当你知道标签并想查看其操作时

### 参数

| 参数 | 类型 | 必需 | 描述 |
|------|------|------|------|
| `tagId` | string | 是 | 标签的 32 字符 MD5 哈希 |

### 响应

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collection": {
    "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
    "title": "Weather Forecast",
    "countMethods": 12
  },
  "tag": {
    "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
    "title": "forecast",
    "countMethods": 5
  },
  "endpoints": [
    {
      "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
      "method": "GET",
      "path": "/v1/forecast",
      "summary": "Get weather forecast for a location"
    }
  ]
}
```

### 细节

- 如果标签不存在，返回 `not_found`
- 包含完整上下文：spec、collection 和标签元数据
- 端点限定在单个 collection 中的单个标签内

---

## endpoint_by_id

### 用途

获取单个端点的快速摘要：方法、路径、摘要和弃用状态。这是一个轻量级工具 — 要获取完整的 OpenAPI 操作对象（参数、请求体、响应模式），使用 `inspect`。

### 何时使用

- 当你有一个端点 ID 并想快速提醒其功能时
- 在决定是否调用 `inspect` 获取完整详情之前
- 当你需要在调用之前确认方法和路径时

### 参数

| 参数 | 类型 | 必需 | 描述 |
|------|------|------|------|
| `id` | string | 是 | 端点的 32 字符 MD5 哈希 |

### 响应

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collection": {
    "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
    "title": "Weather Forecast",
    "countMethods": 12
  },
  "tag": {
    "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
    "title": "forecast",
    "countMethods": 5
  },
  "endpoint": {
    "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
    "method": "GET",
    "path": "/v1/forecast",
    "summary": "Get weather forecast for a location"
  }
}
```

| 字段 | 类型 | 描述 |
|------|------|------|
| `endpoint.id` | string | 端点标识符 |
| `endpoint.method` | string | HTTP 方法 |
| `endpoint.path` | string | API 路径 |
| `endpoint.summary` | string | 人类可读的摘要 |

### 细节

- 如果端点不存在，返回 `not_found`
- 这是**快速摘要** — 它不返回参数、请求体或响应模式
- 要获取完整的技术细节（在 `invoke` 之前必需），使用 `inspect`
