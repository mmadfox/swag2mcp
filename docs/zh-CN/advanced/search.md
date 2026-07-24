# 全文搜索

## 概述

swag2mcp 包含一个内置的全文搜索引擎（bluge），用于索引所有 spec 中的所有端点。LLM 可以按方法、路径、摘要或标签搜索端点 — 即使不知道端点 ID。

## 索引工作原理

当添加或更新 spec 时，每个端点都会被索引。以下字段可搜索：

| 字段 | 描述 | 示例 |
|------|------|------|
| `method` | HTTP 方法 | `GET`、`POST`、`PUT` |
| `path` | API 端点路径 | `/api/v1/users/{id}` |
| `summary` | OpenAPI 摘要 | "Find pet by ID" |
| `tag` | 端点类别 | "pets"、"users" |
| `_all` | 所有字段组合 | method + path + tag + summary |

索引在每次 MCP 服务器启动时重建。它存储在内存中以实现快速搜索。

## 查询语法

搜索支持丰富的查询语法，用于精确过滤：

| 示例 | 描述 |
|------|------|
| `pet` | 跨所有字段的简单文本搜索 |
| `method:GET` | 查找所有 GET 端点 |
| `tag:pets` | 查找 "pets" 标签中的端点 |
| `path:"/api/v1/users"` | 精确路径匹配 |
| `+method:POST +tag:pet` | 必须匹配两个条件 |
| `-method:DELETE` | 排除 DELETE 方法 |
| `create~` | 模糊搜索（容错） |
| `cr*` | 通配符搜索 |
| `"find pet"` | 短语搜索 |
| `+summary:pet -method:DELETE` | 摘要中包含 "pet"，排除 DELETE |

### 字段特定搜索

你可以使用 `field:value` 语法在特定字段内搜索：

```
method:GET
tag:pets
path:"/pet/findByStatus"
summary:"find pet by status"
```

### 布尔运算符

- `+` — 必须匹配（AND）
- `-` — 必须不匹配（NOT）
- 术语之间的空格 — OR（任何术语可以匹配）

### 模糊和通配符

- `term~` — 模糊搜索（匹配相似单词，处理拼写错误）
- `te*` — 通配符（匹配任意字符）
- `te?t` — 单字符通配符

## 示例

```
# 查找所有 GET 请求
method:GET

# 查找 pet 标签中的 POST 请求
+method:POST +tag:pet

# 按精确路径查找端点
path:"/pet/findByStatus"

# 按描述查找
"find pet by status"

# 查找除 DELETE 之外的所有内容
+summary:pet -method:DELETE

# "create" 的模糊搜索（处理拼写错误）
create~
```

## MCP 工具

`search` MCP 工具将搜索引擎暴露给 LLM：

```
→ search(query: "find pet by status", limit: 5)
← GET /pet/findByStatus — Finds Pets by status
   GET /pet/{petId} — Find pet by ID
```

### 参数

| 参数 | 必需 | 描述 |
|------|------|------|
| `query` | 是 | 搜索查询（支持结构化语法） |
| `limit` | 是 | 最大结果数（1-50） |

## 重要说明

- **索引在内存中** — 每次 MCP 服务器启动时重建。没有持久的索引文件。
- **所有字段小写** — 搜索不区分大小写
- **限制上限为 50** — 你不能请求超过 50 个结果
- **无效的查询语法** 返回带有示例的帮助性错误消息
- **`_all` 字段** 组合 method、path、tag 和 summary 用于简单文本搜索
