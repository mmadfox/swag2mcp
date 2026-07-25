# 端点

端点是特定的 HTTP 方法 + 路径，可以被调用（例如 `GET /api/users/{id}`）。端点是 LLM 发现、检查和调用的实际 API 操作。

## 结构

每个端点包含：

- **HTTP 方法**：GET、POST、PUT、PATCH、DELETE、HEAD、OPTIONS
- **路径**：`/api/v1/users/{id}`
- **摘要**：端点功能的简短描述 — 对 LLM 快速理解其用途非常有用
- **描述**：端点行为、参数和用例的详细说明
- **参数**：路径、查询、头、cookie
- **请求体**：用于 POST/PUT/PATCH
- **响应**：状态码和响应模式

`summary` 和 `description` 字段来自 OpenAPI/Swagger/Postman 文件。它们是 LLM 理解端点功能的主要方式。编写良好的摘要使端点发现更加高效。

## 端点的 MCP 工具

| 工具 | 描述 |
|------|------|
| `endpoint_by_spec` | spec 中的所有端点 |
| `endpoint_by_collection` | collection 中的端点 |
| `endpoint_by_tag` | 标签中的端点 |
| `endpoint_by_id` | 快速端点摘要 |
| `inspect` | 完整端点详情（模式、参数） |
| `invoke` | 调用端点 |
| `search` | 按文本搜索端点 |

## 弃用的端点

在规范中标记为 `deprecated` 的端点在检查时会显示通知。

## 配置

从 swag2mcp 的角度来看，端点是**只读的**。端点没有 YAML 配置设置 — 你不能在 `swag2mcp.yaml` 中添加、删除、重命名或修改它们。

要更改端点（添加新端点、更新摘要、修改参数、标记为已弃用），请编辑原始的 OpenAPI/Swagger/Postman 文件并运行 `swag2mcp update` 以重新解析和重新索引。

## 示例

```
查询："Show details for GET /pet/{petId}"
→ inspect(endpointId: "abc123...")
→ 结果：
  GET /pet/{petId}
  摘要：Find pet by ID
  描述：Returns a single pet by its ID
  参数：
    - petId (path, integer, required)
  响应：
    - 200: Pet object
    - 400: Error
    - 404: Not found
```
