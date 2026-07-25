# Specs

Spec 是代表 API 域或服务的逻辑容器（例如 YouTube、Binance、Open-Meteo）。每个 spec 有唯一的 `domain`、`base_url`、可选的 `auth`，并包含一个或多个 collection。

[Collections](./collections) 指向 OpenAPI/Swagger/Postman 文件 — spec 本身不是文件，而是围绕它们的分组。

## 域 — 命名规则

`domain` 是 spec 的唯一标识符。它作为整个系统的主键使用。

| 规则 | 约束 |
|------|------|
| 字符 | 仅限 `a-z`、`0-9`、`_`、`-` |
| 长度 | 1–60 字符 |
| 唯一性 | **不允许重复** — 两个活动的 spec 不能共享相同的 domain |

**有效示例：** `meteo`、`binance`、`github-api`、`my_service`、`openai-v1`

**无效示例：** `Meteo`（大写）、`my api`（空格）、`my.api`（点）、`a-very-long-domain-name-that-exceeds-sixty-characters`（太长）

## Spec 字段

| 字段 | YAML 键 | 必需 | 描述 |
|------|---------|------|------|
| [域](#域--命名规则) | `domain` | ✅ | 唯一 API 标识符（1–60 字符，`a-z0-9_-`） |
| LLM 标题 | `llm_title` | ✅ | LLM 用于引用此 API 的人类可读名称（5–120 字符） |
| [LLM 指令](#llm-instruction) | `llm_instruction` | ❌ | 注入到 swag2mcp 系统提示中的简短提示（最多 500 字符） |
| 基础 URL | `base_url` | ✅ | 所有 API 请求的基础 URL（有效 URL） |
| [禁用](#disable) | `disable` | ❌ | 加载和索引时跳过此 spec |
| [标签](#tags) | `tags` | ❌ | 用于过滤的标签（例如 `["public", "demo"]`） |
| [认证](#auth) | `auth` | ❌ | 认证配置 |
| [HTTP 客户端](#http-client) | `http_client` | ❌ | 每个 spec 的 HTTP 设置（头、cookie） |
| [Collections](./collections) | `collections` | ✅ | 1–30 个 collection 的列表 |

## 验证

当 swag2mcp 验证配置时，会为每个 spec 检查以下规则：

| 检查项 | 规则 |
|--------|------|
| **重复域** | 没有两个活动的 spec 可以共享相同的 `domain` |
| **域格式** | 必须匹配 `^[a-z0-9_-]{1,60}$` |
| **LLM 标题** | 必需，5–120 字符，字母/数字/空格/基本标点 |
| **LLM 指令** | 最多 500 字符，与标题相同的字符集 |
| **基础 URL** | 必需，必须是有效的 URL |
| **Collections** | 必需，1–30 项 |
| **认证** | 按认证类型验证（例如 bearer 需要 `token`，basic 需要 `username` + `password`） |
| **位置** | 每个 collection 的 `location` 必须是有效的 URL 或文件路径（5–250 字符） |

验证在每次 `swag2mcp mcp` 启动时运行。如果失败，MCP 服务器将不会启动 — 在某些 IDE 中，这意味着服务器根本无法连接，LLM 会收到清晰的错误消息，说明需要修复什么。

要在启动服务器之前诊断问题，使用 [`validate`](../cli/validate.md) 命令：

```bash
# 验证默认工作区（~/.swag2mcp）
swag2mcp validate

# 验证自定义项目工作区
swag2mcp validate ./my-project
```

## LLM 指令

建议在每个 spec 上设置 `llm_instruction` — 一个简短的提示（最多 500 字符），告诉 LLM 此 API 的用途以及何时使用。此指令被注入到 swag2mcp 系统提示中，帮助 LLM 无需额外上下文即可理解 spec 的用途。

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    llm_instruction: "Use this API to get random dad jokes or search for specific jokes by keyword."
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

Collection 也可以有自己的 `llm_instruction`（最多 360 字符）以提供更具体的指导。

## 认证

认证在 spec 级别配置，并应用于其所有 collection。swag2mcp 支持 9 种认证方法：

| 方法 | YAML 类型 | 关键字段 |
|------|-----------|----------|
| [None](../auth/none.md) | `none` | — |
| [Basic](../auth/basic.md) | `basic` | `username`、`password` |
| [Bearer](../auth/bearer.md) | `bearer` | `token` |
| [Digest](../auth/digest.md) | `digest` | `username`、`password` |
| [OAuth2 Client Credentials](../auth/oauth2-cc.md) | `oauth2-cc` | `client_id`、`client_secret`、`token_url` |
| [OAuth2 Password](../auth/oauth2-pwd.md) | `oauth2-pwd` | `username`、`password`、`client_id`、`token_url` |
| [API Key](../auth/api-key.md) | `api-key` | `key`、`value`、`in`（`header` 或 `query`） |
| [HMAC](../auth/hmac.md) | `hmac` | `api_key`、`secret_key` |
| [Script](../auth/script.md) | `script` | `domain` |

每种方法的完整详情请参见[认证概述](../auth/overview.md)。

## HTTP 客户端

你可以在 spec 级别覆盖 HTTP 设置。这些设置适用于此 spec 的 collection 发出的所有请求。

```yaml
specs:
  - domain: slow-api
    llm_title: Slow API
    base_url: https://slow-api.example.com
    http_client:
      headers:
        X-API-Version: "2"
      cookies:
        - name: session
          value: abc123
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

设置级联：全局 → spec → collection。详情请参见[配置级联](../configuration/cascade.md)。

## 标签

标签让你按类别过滤 spec。在 `swag2mcp ls` 或引导过程中使用 `--tags` 标志。

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    tags: ["weather", "public"]
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

```bash
# 仅列出标记为 "weather" 的 spec
swag2mcp ls --tags weather
```

## 禁用

设置 `disable: true` 完全跳过 spec。它不会被加载、索引或提供给 LLM。

```yaml
specs:
  - domain: old-api
    llm_title: Old API (Deprecated)
    base_url: https://old-api.example.com
    disable: true
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## 示例

### 最小 Spec

```yaml
specs:
  - domain: dadjokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

### 带认证的 Spec

```yaml
specs:
  - domain: binance
    llm_title: Binance Market Data API
    base_url: https://api.binance.com
    auth:
      type: hmac
      config:
        api_key: $(BINANCE_API_KEY)
        secret_key: $(BINANCE_SECRET_KEY)
    collections:
      - llm_title: Market Data
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/binance.yaml
```

### 带多个 Collection 的 Spec

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
      - llm_title: Marine
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
```

### 带 LLM 指令和标签的 Spec

```yaml
specs:
  - domain: rickandmorty
    llm_title: Rick and Morty API
    llm_instruction: "Use this API to get information about characters, episodes, and locations from the Rick and Morty show."
    base_url: https://rickandmortyapi.com/api
    tags: ["entertainment", "public"]
    collections:
      - llm_title: Characters
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/rick-and-morty.json
```

## 相关

- [Spec 设置（配置）](../configuration/spec-settings.md) — 完整 YAML 参考
- [配置级联](../configuration/cascade.md) — 设置如何相互覆盖
- [认证概述](../auth/overview.md) — 全部 9 种认证方法
- [HTTP 客户端](../configuration/http-client.md) — HTTP 客户端配置
- [Collections](./collections) — spec 内的规范文件
