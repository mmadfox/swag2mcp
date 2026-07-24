# 发现工具

发现工具让 LLM 导航 spec 层次结构：查找所有 spec、深入 spec 查看其 collection，以及探索 collection 中的标签。从 `spec_list` 开始查看可用的 API，然后使用 ID 深入探索。

---

## spec_list

### 用途

列出工作区中注册的所有 API 规范。这是任何会话的起点 — LLM 首先调用它以发现可用的 API。

### 何时使用

- 在会话开始时查看配置了哪些 API
- 添加或删除 spec 后刷新列表
- 当你需要其他工具的 spec ID 时

### 工作原理

返回所有 spec 的列表，包含其唯一 ID 和域名。无需参数。

### 参数

无。

### 响应

```json
{
  "specs": [
    {
      "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
      "domain": "meteo"
    },
    {
      "id": "b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7",
      "domain": "dadjoke"
    }
  ]
}
```

| 字段 | 类型 | 描述 |
|------|------|------|
| `id` | string | 32 字符 MD5 哈希，spec 的唯一标识符 |
| `domain` | string | spec 的域名（例如"meteo"、"dadjoke"） |

### 细节

- 只返回 `id` 和 `domain` — 要获取完整详情（collection、标签），使用 `spec_by_id`
- 所有 ID 是 32 字符的 MD5 十六进制字符串（`^[0-9a-f]{32}$`）
- 如果没有配置 spec，返回空数组

---

## spec_by_id

### 用途

获取特定 spec 的详细信息：其域、所有 collection 及其统计信息（标签计数、方法计数）。

### 何时使用

- 在 `spec_list` 之后查看 spec 内的 collection
- 当你需要 collection ID 以进一步导航时

### 工作原理

接受 spec ID 并返回 spec 元数据及其所有 collection 的计数。

### 参数

| 参数 | 类型 | 必需 | 描述 |
|------|------|------|------|
| `id` | string | 是 | spec 的 32 字符 MD5 哈希 |

### 响应

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collections": [
    {
      "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
      "title": "Weather Forecast",
      "llmTitle": "Forecast API",
      "countTags": 3,
      "countMethods": 12
    }
  ]
}
```

| 字段 | 类型 | 描述 |
|------|------|------|
| `spec.id` | string | Spec 标识符 |
| `spec.domain` | string | Spec 域名 |
| `collections[].id` | string | Collection 标识符 |
| `collections[].title` | string | 人类可读的标题 |
| `collections[].llmTitle` | string | 对 LLM 友好的标题（可选） |
| `collections[].countTags` | int | collection 中的标签数量 |
| `collections[].countMethods` | int | collection 中的 HTTP 方法数量 |

### 细节

- 如果 spec ID 不存在，返回 `not_found` 错误
- `id` 必须是有效的 32 字符 MD5 十六进制字符串

---

## collection_by_spec

### 用途

列出特定 spec 中的所有 collection。类似于 `spec_by_id`，但只返回 collection 列表，不包含额外的 spec 元数据。

### 何时使用

- 当你已有 spec ID，只需要 collection 列表时
- 作为 `spec_by_id` 的轻量级替代

### 参数

| 参数 | 类型 | 必需 | 描述 |
|------|------|------|------|
| `specId` | string | 是 | spec 的 32 字符 MD5 哈希 |

### 响应

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collections": [
    {
      "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
      "title": "Weather Forecast",
      "llmTitle": "Forecast API",
      "countTags": 3,
      "countMethods": 12
    }
  ]
}
```

### 细节

- 如果 spec 不存在，返回 `not_found`
- 与 `spec_by_id` 数据相同，但没有额外的 spec 包装

---

## collection_by_id

### 用途

获取特定 collection 的详细信息：其元数据、父 spec 以及 collection 中的所有标签。

### 何时使用

- 在 `collection_by_spec` 之后查看 collection 内的标签
- 当你需要 `tag_by_id` 或 `endpoint_by_tag` 的标签 ID 时

### 参数

| 参数 | 类型 | 必需 | 描述 |
|------|------|------|------|
| `id` | string | 是 | collection 的 32 字符 MD5 哈希 |

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
  "tags": [
    {
      "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "title": "forecast",
      "countMethods": 5
    },
    {
      "id": "e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7",
      "title": "current",
      "countMethods": 7
    }
  ]
}
```

| 字段 | 类型 | 描述 |
|------|------|------|
| `spec` | object | 父 spec（id、domain） |
| `collection` | object | Collection 元数据（id、title、countMethods） |
| `tags[]` | array | 标签列表，包含 id、title、countMethods |

### 细节

- 如果 collection ID 不存在，返回 `not_found`
- 标签随其 ID 一起返回 — 使用 `endpoint_by_tag(tagId)` 查看实际端点

---

## tag_by_spec

### 用途

列出整个 spec 中的所有标签，跨越所有 collection。用于鸟瞰所有可用标签。

### 何时使用

- 当你想查看 spec 中的所有标签，而无需深入每个 collection 时
- 当你不知道哪个 collection 包含你需要的标签时

### 参数

| 参数 | 类型 | 必需 | 描述 |
|------|------|------|------|
| `specId` | string | 是 | spec 的 32 字符 MD5 哈希 |

### 响应

```json
{
  "tags": [
    {
      "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "title": "forecast",
      "countMethods": 5
    },
    {
      "id": "e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7",
      "title": "current",
      "countMethods": 7
    }
  ]
}
```

### 细节

- 如果 spec 不存在，返回 `not_found`
- 标签从 spec 中的所有 collection 聚合

---

## tag_by_collection

### 用途

列出特定 collection 中的所有标签。与 `tag_by_spec` 不同，此工具还返回父 spec 和 collection 元数据。

### 何时使用

- 在 `collection_by_id` 之后确认标签列表
- 当你需要完整上下文（spec + collection + tags）时

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
  "tags": [
    {
      "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "title": "forecast",
      "countMethods": 5
    }
  ]
}
```

### 细节

- 如果 collection 不存在，返回 `not_found`
- 与 `tag_by_spec` 的标签数据相同，但限定在一个 collection 内

---

## tag_by_id

### 用途

获取单个标签的信息：其 ID、标题以及包含的方法数量。这告诉你关于标签本身的信息 — 要查看实际端点，使用 `endpoint_by_tag`。

### 何时使用

- 当你有一个标签 ID 并想确认其名称和大小
- 在调用 `endpoint_by_tag` 之前了解预期有多少端点

### 参数

| 参数 | 类型 | 必需 | 描述 |
|------|------|------|------|
| `id` | string | 是 | 标签的 32 字符 MD5 哈希 |

### 响应

```json
{
  "tag": {
    "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
    "title": "forecast",
    "countMethods": 5
  }
}
```

| 字段 | 类型 | 描述 |
|------|------|------|
| `tag.id` | string | 标签标识符 |
| `tag.title` | string | 人类可读的标签名称 |
| `tag.countMethods` | int | 此标签中的 HTTP 方法数量 |

### 细节

- 如果标签不存在，返回 `not_found`
- 此工具只返回标签元数据 — 使用 `endpoint_by_tag` 获取实际端点列表
