# MCP 工具

## 概述

swag2mcp 提供 **19 个 MCP 工具**，让 LLM 智能体通过模型上下文协议完全访问你的 API。这些工具涵盖完整的工作流程：发现可用的 API、导航 spec 层次结构、搜索和检查端点、执行 API 调用以及处理大响应。

### 工具解决的问题

- **发现** — LLM 无需事先知道 ID 即可找到 spec、collection 和标签
- **导航** — 在结构化层次结构中从 spec → collection → tag → endpoint 深入
- **搜索** — 当你没有 ID 时，跨所有端点进行全文搜索
- **检查** — 在调用之前获取完整的 OpenAPI 操作对象
- **执行** — 使用自动认证调用真实的 API
- **大响应处理** — 概述、压缩和切片不适合内联的过大响应

### 只读 vs 可变

| 类型 | 数量 | 工具 |
|------|------|------|
| **只读** | 17 | 所有发现、端点、搜索、检查、信息和响应工具 |
| **可变** | 2 | `invoke`（发出真实 HTTP 调用）、`auth`（检索令牌） |

只读工具在 MCP 协议中标记为 `ReadOnlyHint=true` 和 `IdempotentHint=true`，向 LLM 表明它们可以安全调用，没有副作用。

### 错误处理

所有工具以结构化的 `LLMError` 对象返回错误，包含机器可读的代码和人类可读的消息，解释出了什么问题以及下一步该怎么做：

| 错误代码 | 含义 |
|----------|------|
| `validation_failed` | 无效输入（错误的 ID 格式、缺少必需字段） |
| `not_found` | 在索引或工作区中未找到实体 |
| `rate_limit` | 10 秒内对同一端点的第二次 `invoke` 调用 |
| `invoke_error` | HTTP 调用失败、下载失败 |
| `auth_error` | 认证令牌检索失败 |
| `config_error` | 配置文件加载或保存失败 |
| `parse_error` | 规范文件解析失败 |

## 分类

| 类别 | 工具 | 描述 |
|------|------|------|
| **发现** | `spec_list`、`spec_by_id`、`collection_by_spec`、`collection_by_id`、`tag_by_spec`、`tag_by_collection`、`tag_by_id` | 导航 spec 层次结构：查找 spec、collection 和标签 |
| **端点** | `endpoint_by_spec`、`endpoint_by_collection`、`endpoint_by_tag`、`endpoint_by_id` | 在层次结构的不同级别查看端点 |
| **执行** | `search`、`inspect`、`invoke` | 搜索、检查完整契约和调用 API |
| **实用工具** | `auth`、`info`、`response_outline`、`response_compress`、`response_slice` | 认证令牌、运行时信息和大响应处理 |
| **技能** | [格式化指南](/mcp-tools/skills) | 自定义工具响应的显示方式 |

## 完整列表

| 工具 | 描述 |
|------|------|
| `spec_list` | 列出工作区中的所有 API 规范 |
| `spec_by_id` | 获取包含 collection 的详细 spec 信息 |
| `collection_by_spec` | 列出 spec 中的 collection |
| `collection_by_id` | 获取包含标签的 collection 详情 |
| `tag_by_spec` | 列出 spec 中的所有标签 |
| `tag_by_collection` | 列出 collection 中的标签 |
| `tag_by_id` | 获取标签详情（ID、标题、方法计数） |
| `endpoint_by_spec` | 列出 spec 中的所有端点 |
| `endpoint_by_collection` | 列出 collection 中的端点 |
| `endpoint_by_tag` | 列出标签中的端点 |
| `endpoint_by_id` | 快速端点摘要（方法、路径、摘要） |
| `search` | 跨所有端点的全文搜索 |
| `inspect` | 完整的 OpenAPI 操作详情（参数、模式） |
| `invoke` | 执行真实的 API 调用 |
| `auth` | 获取 spec 的认证令牌或头 |
| `info` | 运行时信息（版本、spec、配置） |
| `response_outline` | 大响应文件的结构摘要 |
| `response_compress` | 压缩大响应以适应内联 |
| `response_slice` | 提取大响应的片段 |

## 导航层次结构

```
spec_list
  └── spec_by_id(id)
        └── collection_by_spec(specId)
              └── collection_by_id(id)
                    └── tag_by_collection(collectionId)
                          └── tag_by_id(id)
                                └── endpoint_by_tag(tagId)
                                      └── endpoint_by_id(id)
                                            └── inspect(endpointId)
                                                  └── invoke(endpointId)
```

当你没有 ID 时，使用 `search` 按查询查找端点。当 `invoke` 返回 `fileRef`（响应太大）时，使用 `response_outline` → `response_compress` 或 `response_slice` 探索数据。
