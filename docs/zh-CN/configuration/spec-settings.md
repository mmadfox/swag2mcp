# Spec 设置

Spec 设置定义 API 服务，并为该特定 API 覆盖全局设置。每个 spec 代表一个逻辑 API（例如"Open-Meteo Weather APIs"），可以包含多个 collection（规范文件）。

## Spec 部分

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    llm_instruction: "Use this API for weather forecasts and climate data"
    base_url: https://api.open-meteo.com
    disable: false
    tags: ["weather", "climate"]
    http_client:
      timeout: 10s
      max_response_size: 1024
    auth:
      type: bearer
      config:
        token: "$(TOKEN)"
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

## 参数

### domain

- **类型：** `string`
- **必需：** 是
- **描述：** 此 API spec 的唯一标识符。内部用于引用 spec。
- **规则：** 1-60 字符。仅限小写字母（`a-z`）、数字（`0-9`）、连字符（`-`）和下划线（`_`）。
- **示例：** `meteo`、`binance`、`my-api`

### llm_title

- **类型：** `string`
- **必需：** 是
- **描述：** LLM 用于引用此 API 的人类可读名称。显示在 MCP 工具响应中。
- **规则：** 5-120 字符。仅限字母、数字、空格和基本标点。
- **示例：** `Open-Meteo Weather APIs`、`Binance Market Data`

### llm_instruction

- **类型：** `string`
- **默认值：** `""`
- **描述：** 关于如何使用此 API 的 LLM 指令。描述 API 的功能和何时使用。
- **规则：** 最多 500 字符。仅限字母、数字、空格和基本标点。
- **示例：** `"Use this API for weather forecasts, current conditions, and climate data."`

### base_url

- **类型：** `string`
- **必需：** 是
- **描述：** 此 spec 中所有 API 请求的基础 URL。OpenAPI 规范中的端点路径会追加到此 URL。
- **示例：** `https://api.open-meteo.com`、`https://api.binance.com`
- **注意：** 如果不同的 collection 使用不同的基础 URL，可以在 collection 级别覆盖。

### disable

- **类型：** `bool`
- **默认值：** `false`
- **描述：** 当为 `true` 时，此 spec 被排除在 MCP 工具之外。它不会被加载、索引或提供给 LLM。
- **何时使用：** 临时禁用一个 API 而不从配置中删除它。适用于已关闭、已弃用或正在维护的 API。

### tags

- **类型：** `[]string`（字符串数组）
- **默认值：** `[]`
- **描述：** 用于过滤 spec 的标签。与 CLI 命令中的 `--tags` 标志一起使用（`ls`、`validate`、`mcp`、`update`）。
- **示例：** `["public", "weather"]`、`["internal", "production"]`
- **效果：** 当你运行 `swag2mcp mcp --tags=public` 时，只加载具有 `public` 标签的 spec。

### http_client

- **类型：** `object`
- **默认值：** 继承自全局
- **描述：** 为此 spec 覆盖全局 HTTP 客户端设置。全局 `http_client` 中的所有设置都可以被覆盖：`timeout`、`max_response_size`、`user_agent`、`follow_redirects`、`max_redirects`、`random`、`proxy`、`headers`、`cookies`。
- **示例：**
  ```yaml
  http_client:
    timeout: 60s
    max_response_size: 4194304
    headers:
      "X-DC": "us-east-1"
  ```

### auth

- **类型：** `object`
- **默认值：** `none`（无认证）
- **描述：** 此 spec 的认证配置。请参见[认证](/auth/overview)部分了解所有 9 种方法及其参数。
- **示例：**
  ```yaml
  auth:
    type: bearer
    config:
      token: "$(API_TOKEN)"
  ```

### collections

- **类型：** `[]object`（collection 数组）
- **必需：** 是（至少 1 个）
- **描述：** 属于此 spec 的 OpenAPI/Swagger/Postman 规范文件列表。每个 collection 是一个规范文件。
- **规则：** 每个 spec 1-30 个 collection。
- **参见：** [Collection 设置](./collection-settings)了解所有 collection 参数。

## 禁用一个 Spec

禁用的 spec 不会被加载或索引。LLM 无法看到或使用它们。

```yaml
specs:
  - domain: old-api
    llm_title: Old API
    base_url: https://old-api.example.com
    disable: true
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## HTTP 客户端覆盖

全局级别的所有 `http_client` 设置可以在 spec 级别被覆盖。Spec 值仅对此 spec 优先于全局值。

```yaml
specs:
  - domain: slow-api
    llm_title: Slow API
    base_url: https://slow-api.example.com
    http_client:
      timeout: 120s
      max_response_size: 8388608
      headers:
        "X-DC": "us-east-1"
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## 代理覆盖

如果此 spec 需要与全局不同的代理，在 spec 级别配置：

```yaml
specs:
  - domain: proxied-api
    llm_title: Proxied API
    base_url: https://api.example.com
    http_client:
      proxy:
        url: http://proxy.company.com:8080
        username: $(PROXY_USER)
        password: $(PROXY_PASS)
        bypass:
          - "*.local"
          - "10.0.0.0/8"
    collections:
      - llm_title: Main
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo.json
```
