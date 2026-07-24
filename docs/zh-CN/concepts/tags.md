# 标签

标签是 collection 内对相关端点进行分组的类别。标签可能存在也可能不存在 — 并非所有 collection 都有标签，一个 collection 可以有任意数量的标签。

标签来自 OpenAPI/Swagger/Postman 文件本身。标签**没有 YAML 配置设置** — 你不能在 `swag2mcp.yaml` 中创建、重命名或删除标签。更改标签的唯一方法是编辑原始规范文件。

## 层次结构

```
Spec (domain, e.g. "meteo")
  └── Collection (spec file, e.g. forecast.yml)
        └── Tag "weather"
              └── GET /forecast
              └── GET /forecast/hourly
        └── Tag "alerts"
              └── GET /alerts
```

## 标签的创建方式

标签在解析过程中从规范文档中提取：

**OpenAPI 3.x / Swagger 2.0** — 每个操作的 `tags` 列表成为标签：

```yaml
paths:
  /pet:
    get:
      tags: ["pets"]
      summary: "Find pet by ID"
    post:
      tags: ["pets"]
      summary: "Add a new pet"
  /pet/{petId}/uploadImage:
    post:
      tags: ["pet_images"]
      summary: "Uploads an image"
```

**Postman** — 每个顶级文件夹成为一个标签。嵌套文件夹使用最后一个文件夹名称。

如果端点没有标签，它会被放在 `"default"` 标签下。

## 用途

标签帮助 LLM 找到相关端点组。LLM 无需搜索 collection 中的每个端点，而是可以先找到正确的标签，然后只列出其中的端点。

## 标签的 MCP 工具

| 工具 | 描述 |
|------|------|
| `tag_by_spec` | 整个 spec 中的所有标签 |
| `tag_by_collection` | 特定 collection 中的标签 |
| `tag_by_id` | 标签详情（标题、方法计数） |
| `endpoint_by_tag` | 按标签分组的端点 |

## 示例

```
查询："Show all tags in the pet collection"
→ tag_by_collection(collectionId: "...")
→ 结果：pets (5 methods), pet_images (1 method)
```

## 限制

- 从配置角度来看，标签是只读的。要添加、重命名或删除标签，请编辑原始的 OpenAPI/Swagger/Postman 文件并运行 `swag2mcp update`。
- 不能在 YAML 配置中按 collection 过滤或禁用标签。
