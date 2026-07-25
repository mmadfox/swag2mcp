# 概念

## 架构

swag2mcp 充当 API 规范和 LLM 智能体之间的桥梁：

<img src="/architecture.svg" width="800" alt="swag2mcp 架构">

## 核心概念

**Spec** — 代表 API 域或服务的逻辑容器（例如 YouTube、Binance、Open-Meteo）。每个 spec 有唯一的 `domain`、`base_url`、可选的 `auth`，并包含一个或多个 collection。你还可以设置 `llm_instruction` — 注入到 swag2mcp 系统提示中的简短提示，告诉 LLM 此 spec 的用途以及何时使用。了解更多：[Specs](./specs)。

**Collection** — 描述特定 API 的单个 OpenAPI/Swagger/Postman 文件。它指向一个 `location`（URL 或本地文件路径）。一个 spec 可以有多个 collection — 例如，"meteo" spec 可能有"Forecast"、"Air Quality"和"Marine" collection，每个指向不同的规范文件。了解更多：[Collections](./collections)。

**Tag** — collection 内端点的类别。帮助 LLM 更精确地找到正确的操作。了解更多：[Tags](./tags)。

**Endpoint** — 特定的 HTTP 方法 + 路径（例如 `GET /api/users`）。LLM 可以通过描述找到端点，检查其参数和模式，然后调用它。了解更多：[Endpoints](./endpoints)。

**Workspace** — swag2mcp 存储配置、规范缓存、保存的响应和认证脚本的目录。了解更多：[Workspace](./workspace)。

## 工作原理

1. **添加 spec 或 collection** — 在 YAML 配置中定义它（`~/.swag2mcp/swag2mcp.yaml`）。例如：

   ```yaml
   specs:
     - domain: jokes
       llm_title: Dad Joke API
       base_url: https://icanhazdadjoke.com
       collections:
         - llm_title: Jokes
           location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
   ```
2. **swag2mcp 解析每个 collection** — 创建标签和端点，为搜索建立索引。
3. **LLM 找到正确的端点** — 通过 MCP 工具（`search`、`endpoint_by_tag`、`inspect`），LLM 通过描述搜索匹配的端点，查看其参数和请求模式。
4. **LLM 调用端点** — 通过 MCP 工具 `invoke`，LLM 发送请求。swag2mcp 在发起调用之前，会针对端点的 OpenAPI 模式验证每个输入参数（路径参数、查询参数、头、请求体）。如果某些内容与模式不匹配，LLM 会收到清晰的错误说明问题所在。验证通过后，swag2mcp 执行真实的 HTTP 调用并返回结果。
5. **结果返回给 LLM** — API 响应被传回给智能体。大响应保存到工作区，可以使用三个专用的 MCP 工具进行探索：`response_outline`（查看结构）、`response_compress`（缩小为代表性样本）和 `response_slice`（提取特定片段）。

swag2mcp 是 LLM 和 API 世界之间的桥梁。你添加 API 规范，LLM — 通过 MCP 协议 — 找到正确的端点，检查其文档，并调用它们。你只需要添加一个 spec 并启动 MCP 服务器。

> **配置随时可编辑。** YAML 配置文件（`~/.swag2mcp/swag2mcp.yaml`）可以手动编辑 — 添加 spec、更改认证、调整设置。每次编辑后，重启 MCP 服务器（`swag2mcp mcp`）以使更改生效。

## 层次结构

```
Spec (domain, e.g. "meteo")
  └── Collection 1 (spec file, e.g. forecast.yml)
        └── Tag 1 (category)
              └── Endpoint (GET /api/forecast)
              └── Endpoint (POST /api/forecast)
        └── Tag 2
              └── Endpoint (GET /api/forecast/{id})
  └── Collection 2 (spec file, e.g. air-quality.yml)
        └── Tag 3
              └── Endpoint (GET /api/air-quality)
```
