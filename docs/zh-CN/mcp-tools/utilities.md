# 实用工具

实用工具提供辅助功能：检索认证令牌、获取运行时信息以及处理不适合内联的大 API 响应。

---

## auth

### 用途

检索特定 spec 的认证令牌、头或查询参数。这使 LLM 可以访问可在 swag2mcp 之外使用的凭据（例如，生成 curl 命令）。

### 何时使用

- 仅当用户明确要求原始令牌或凭据时
- 当生成需要认证的 curl 命令或代码片段时
- 当用户想查看配置了哪种认证方法时

### 何时不使用

- **不要**在 `inspect` 或 `invoke` 之前调用 `auth` — `invoke` 会自动获取并应用认证
- **不要**仅仅为了检查是否配置了认证而调用 `auth` — 使用 `info` 代替

### 工作原理

查找 spec 的认证配置并执行认证流程（令牌交换、脚本执行等）以获取当前凭据。

### 参数

| 参数 | 类型 | 必需 | 描述 |
|------|------|------|------|
| `specId` | string | 是 | spec 的 32 字符 MD5 哈希 |

### 响应

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "headers": {
    "Authorization": "Bearer eyJhbGciOiJIUzI1NiIs...",
    "X-API-Key": "my-api-key"
  },
  "queryParams": {
    "api_key": "my-api-key"
  }
}
```

| 字段 | 类型 | 描述 |
|------|------|------|
| `token` | string | 原始令牌值（bearer 令牌、API 密钥等） |
| `headers` | object | 请求中包含的 HTTP 头 |
| `queryParams` | object | 请求中包含的查询参数 |

### 细节

- **生产环境默认禁用：** `--disable-llm-auth` 标志（默认：`true`）从 MCP 工具列表中完全移除 `auth` 工具。LLM 无法看到或请求令牌。设置 `--disable-llm-auth=false` 以在调试或使用短期令牌时启用。
- **`invoke` 自动处理认证：** 你不需要在 `invoke` 之前调用 `auth`。invoke 服务自动获取并应用正确的认证。
- **支持 9 种认证方法：** `none`、`basic`、`bearer`、`digest`、`hmac`、`oauth2-cc`（客户端凭证）、`oauth2-pwd`（密码）、`api-key`、`script`。
- 如果认证方法失败（例如 OAuth2 令牌端点不可达、脚本执行失败），返回 `auth_error`。

---

## info

### 用途

返回 swag2mcp 运行时的全面摘要：版本、工作区路径、活动 spec、HTTP 客户端设置、MCP 传输配置、认证方法和模拟模式状态。

### 何时使用

- 当用户询问系统配置时
- 当你需要检查运行时设置（超时、响应大小限制、传输方式）时
- 当你想知道哪些认证方法可用时
- 排查配置问题时

### 工作原理

返回运行时状态的预计算快照。无需参数。

### 参数

无。

### 响应

```json
{
  "version": "v1.2.0",
  "workspace": "~/.swag2mcp",
  "uptime": "2h 15m",
  "specs": {
    "total": 4,
    "active": 3,
    "disabled": 1,
    "collections": 6,
    "endpoints": 42
  },
  "http_client": {
    "timeout": "30s",
    "max_response_size": "2 KB",
    "follow_redirects": true,
    "max_redirects": 10,
    "randomize": false,
    "proxy": null,
    "headers": {},
    "cookies": []
  },
  "mcp": {
    "transport": "stdio",
    "addr": ":8080",
    "path": "/mcp",
    "auth_enabled": false
  },
  "auth": {
    "methods": ["bearer", "api-key"]
  },
  "mock": {
    "enabled": false
  }
}
```

| 字段 | 类型 | 描述 |
|------|------|------|
| `version` | string | swag2mcp 版本 |
| `workspace` | string | 工作区目录路径 |
| `uptime` | string | 服务器运行时间（人类可读） |
| `specs` | object | Spec 摘要：总数、活动、禁用、collection、端点 |
| `http_client` | object | HTTP 客户端配置 |
| `http_client.max_response_size` | string | 最大响应大小，人类可读格式（例如"2 KB"） |
| `mcp` | object | MCP 服务器配置 |
| `auth` | object | 可用的认证方法 |
| `mock` | object | 模拟服务器状态 |

### 细节

- `max_response_size` 以人类可读格式显示（例如 `"1 KB"`、`"2 MB"`）
- `uptime` 从服务器启动时间计算
- 数据是在引导时拍摄的快照 — 它反映了 MCP 服务器启动时的状态

---

## response_outline

### 用途

获取由 `invoke` 保存到磁盘的大 JSON 响应文件的高级结构摘要。它返回数据的形状 — 键、类型、数组长度和导航提示 — 而不返回实际值。

### 何时使用

- 在 `invoke` 返回 `fileRef`（响应太大无法内联）后立即使用
- 这是大响应工作流程中的**强制第一步**

### 工作原理

读取保存的响应文件并分析其结构：顶级类型、键、数组长度、嵌套深度和压缩提示。

### 参数

| 参数 | 类型 | 必需 | 描述 |
|------|------|------|------|
| `path` | string | 是 | 来自 `fileRef.path` 的绝对路径 |
| `maxDepth` | int | 否 | 最大递归深度（默认：3） |
| `maxArrayItems` | int | 否 | 要检查的数组项数（默认：5） |

### 响应

```json
{
  "outline": {
    "type": "object",
    "size": 1572864,
    "lineCount": 12500,
    "depth": 3,
    "structure": {
      "type": "object",
      "keys": ["data", "meta", "error"],
      "data": {
        "type": "array",
        "length": 500,
        "items": {
          "type": "object",
          "keys": ["id", "name", "status", "createdAt"]
        }
      }
    },
    "schemaHint": "object with 3 keys: data (array[500]), meta (object), error (null)",
    "keys": ["data", "meta", "error"],
    "itemCount": 500,
    "itemType": "object",
    "compressionHints": [
      "response_compress(path, 'first_of_array', 'data')",
      "response_compress(path, 'sample_array', 'data', arrayHead=3, arrayTail=2)",
      "response_compress(path, 'keys_only', 'data')",
      "response_compress(path, 'select_keys', 'data', selectKeys=[id, name])"
    ],
    "navigationHints": {
      "paths": ["data", "meta", "error"],
      "arrays": [
        {"path": "data", "length": 500}
      ]
    }
  }
}
```

| 字段 | 类型 | 描述 |
|------|------|------|
| `type` | string | 顶级类型："object"或"array" |
| `size` | int | 文件大小（字节） |
| `lineCount` | int | 文件行数 |
| `depth` | int | 检查的最大嵌套深度 |
| `structure` | object | 递归结构，包含键、类型、数组长度 |
| `schemaHint` | string | 顶级形状的一行摘要 |
| `keys` | array | 顶级键（对于对象） |
| `itemCount` | int | 数组长度（对于数组） |
| `compressionHints` | array | 建议的 `response_compress` 调用及参数 |
| `navigationHints` | object | 顶级路径和数组及其长度 |

### 细节

- 如果路径无效或不在响应目录内，返回 `validation_failed`
- 如果文件不存在，返回 `not_found`
- 如果文件不是有效的 JSON，返回 `validation_failed`
- `compressionHints` 字段提供即用型的 `response_compress` 调用建议

---

## response_compress

### 用途

缩小保存的响应文件中的 JSON 值，使其适合响应大小限制并可以内联返回给 LLM。多种压缩模式让你在大小和信息之间选择合适的权衡。

### 何时使用

- 在 `response_outline` 之后了解结构
- 当你需要从大响应中内联获取数据时
- 当 `response_slice` 太窄，你需要更广泛的视图时

### 工作原理

读取保存的响应文件，导航到指定的 JSON 路径，应用压缩模式，并返回压缩结果。如果结果仍然超过大小限制，则保存到新文件。

### 参数

| 参数 | 类型 | 必需 | 描述 |
|------|------|------|------|
| `path` | string | 是 | 来自 `fileRef.path` 的绝对路径 |
| `jsonPath` | string | 否 | 要压缩的值的路径（例如 `data` 或 `data.0`） |
| `mode` | string | 是 | 压缩模式（见下表） |
| `arrayHead` | int | 否 | `sample_array` 模式下保留的前导项数（默认：3） |
| `arrayTail` | int | 否 | `sample_array` 模式下保留的尾部项数（默认：2） |
| `stringLen` | int | 否 | `truncate_strings` 模式下的最大字符串长度（默认：80） |
| `selectKeys` | array | 否 | `select_keys` 模式下要保留的键 |

### 压缩模式

| 模式 | 描述 | 最适合 |
|------|------|--------|
| `first_of_array` | 只保留数组的第一个元素 | 所有元素结构相同时 |
| `sample_array` | 保留数组的头部和尾部 | 需要查看值的范围时 |
| `truncate_strings` | 将每个字符串缩短到 `stringLen` 字符 | 字符串非常长但结构重要时 |
| `keys_only` | 将对象值替换为类型名称 | 只需要结构时 |
| `select_keys` | 在每个对象中只保留指定的键 | 需要来自多个对象的特定字段时 |

### 响应

```json
{
  "body": [
    { "id": 1, "name": "Rex", "status": "available" },
    { "id": 2, "name": "Max", "status": "pending" }
  ],
  "hint": "Compressed array from 500 to 2 items using first_of_array mode"
}
```

| 字段 | 类型 | 描述 |
|------|------|------|
| `body` | any | 压缩的 JSON 值（在大小限制内时存在） |
| `fileRef` | object | 文件引用（仍然太大时存在） |
| `hint` | string | 压缩内容的说明 |

### 细节

- 如果压缩结果仍然超过 `max_response_size`，则保存到新文件并返回 `FileReference`
- 默认值：`arrayHead=3`、`arrayTail=2`、`stringLen=80`
- 对于无效路径、无效 JSONPath 或非 JSON 文件，返回 `validation_failed`
- 如果文件不存在或 JSONPath 不匹配，返回 `not_found`

---

## response_slice

### 用途

通过逻辑 JSON 路径或行范围提取保存的 JSON 响应文件的特定片段。与 `response_compress` 不同，这给你原始、未修改的数据。

### 何时使用

- 当你需要大响应中的特定元素或值时
- 当 `response_compress` 没有提供足够的细节时
- 当你想逐步浏览响应时

### 工作原理

读取保存的响应文件并通过 JSON 路径（例如 `data.3.name`）或行范围（例如 `120-240`）提取片段。返回逐步浏览数组和对象的导航提示。

### 参数

| 参数 | 类型 | 必需 | 描述 |
|------|------|------|------|
| `path` | string | 是 | 来自 `fileRef.path` 的绝对路径 |
| `jsonPath` | string | 否 | 值的逻辑路径（例如 `data.3.name`） |
| `line` | int | 否 | 以片段为中心的基于 1 的行号 |
| `range` | string | 否 | 行范围，格式为 `start-end`（例如 `120-240`） |
| `around` | int | 否 | `line` 周围包含的行数（默认：20） |

### 响应

```json
{
  "slice": {
    "lines": [120, 130],
    "fragment": "{\n  \"id\": 1,\n  \"name\": \"Rex\"\n}",
    "value": {
      "id": 1,
      "name": "Rex"
    },
    "jsonPath": "data.0",
    "context": "object",
    "isComplete": true,
    "nextLine": 131,
    "prevLine": 119,
    "nextPath": "data.1",
    "prevPath": null
  }
}
```

| 字段 | 类型 | 描述 |
|------|------|------|
| `lines` | array | 基于 1 的行范围 [start, end] |
| `fragment` | string | 原始 JSON 文本（足够小时） |
| `value` | any | 提取的 JSON 值 |
| `jsonPath` | string | 使用的 JSON 路径 |
| `context` | string | "object"、"array"或"value" |
| `isComplete` | bool | 值为有效的 JSON 片段时为 true |
| `nextLine` | int | 基于行的导航的建议下一行 |
| `prevLine` | int | 建议的上一行 |
| `nextPath` | string | 数组导航的建议下一个 JSON 路径 |
| `prevPath` | string | 建议的上一个 JSON 路径 |

### 细节

- **优先使用 `jsonPath` 而不是行号** — JSON 路径稳定且具有描述性，行号在文件重新生成时会变化
- 如果提取的片段超过 `max_response_size`，则保存到新文件并返回 `FileReference`
- 默认 `around` 为 20 行
- 响应包含用于逐步浏览数组的 `nextPath`/`prevPath` 和用于基于行的导航的 `nextLine`/`prevLine`
- 对于无效路径、无效 JSONPath、无效行/范围或非 JSON 文件，返回 `validation_failed`
- 如果文件不存在或 JSONPath 不匹配，返回 `not_found`
