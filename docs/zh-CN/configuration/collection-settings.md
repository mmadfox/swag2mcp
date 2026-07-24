# Collection 设置

Collection 设置定义单个 OpenAPI/Swagger/Postman 规范文件，并为该特定文件覆盖 spec 设置。每个 collection 属于一个 spec，代表一个 API 规范文档。

## Collection 部分

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        llm_instruction: "Use for current and forecast weather data"
        disable: false
        base_url: https://forecast-api.open-meteo.com
        base_mock_url: localhost:8081
        http_client:
          timeout: 5s
```

## 参数

### llm_title

- **类型：** `string`
- **必需：** 否
- **描述：** 此 collection 的人类可读名称。显示在 MCP 工具响应中。
- **规则：** 最多 120 字符。仅限字母、数字、空格和基本标点。
- **示例：** `Forecast`、`Air Quality`、`Market Data`

### llm_instruction

- **类型：** `string`
- **默认值：** `""`
- **描述：** 关于此特定 collection 的 LLM 指令。描述此 collection 提供哪些端点。
- **规则：** 最多 360 字符。仅限字母、数字、空格和基本标点。
- **示例：** `"Use for current and forecast weather data."`

### title

- **类型：** `string`
- **默认值：** `""`
- **描述：** 来自规范文件的原始标题。在运行时自动填充。你通常不需要在 YAML 中设置此字段。

### location

- **类型：** `string`
- **必需：** 是
- **描述：** OpenAPI 3.x、Swagger 2.0 或 Postman collection 规范文件的 URL 或本地文件路径。
- **规则：** 5-250 字符。
- **示例：**
  - URL：`https://raw.githubusercontent.com/org/repo/main/spec.yaml`
  - 本地：`./specs/my-api.json`
  - 本地（绝对路径）：`/home/user/.swag2mcp/specs/my-api.yaml`

### disable

- **类型：** `bool`
- **默认值：** `false`
- **描述：** 当为 `true` 时，此 collection 被排除在 MCP 工具之外。它不会被加载或索引。
- **何时使用：** 临时禁用一个 collection 而不从配置中删除它。当规范文件正在更新或 API 版本已弃用时很有用。

### http_client

- **类型：** `object`
- **默认值：** 继承自 spec（或全局）
- **描述：** 为此 collection 覆盖 HTTP 客户端设置。全局 `http_client` 中的所有设置都可以被覆盖：`timeout`、`max_response_size`、`user_agent`、`follow_redirects`、`max_redirects`、`random`、`proxy`、`headers`、`cookies`。
- **示例：**
  ```yaml
  http_client:
    timeout: 120s
    headers:
      "X-Custom": "value"
    cookies:
      - name: "session"
        value: "abc123"
  ```

### base_url

- **类型：** `string`
- **默认值：** `""`（继承自 spec）
- **描述：** 为此 collection 覆盖 spec 级别的 `base_url`。当同一 spec 中的不同 collection 使用不同的基础 URL 时使用。
- **示例：** 如果 spec 有 `base_url: https://api.open-meteo.com` 但一个 collection 使用 `https://air-quality-api.open-meteo.com`，在 collection 级别设置 `base_url`。

### base_mock_url

- **类型：** `string`
- **默认值：** `""`
- **描述：** 模拟服务器地址，格式为 `host:port`。当全局配置中 `mock_enabled: true` 时必需。
- **规则：** Host 必须是 `localhost`、`127.0.0.1` 或 `0.0.0.0`。端口必须是有效的端口号。
- **示例：** `localhost:8081`、`127.0.0.1:9000`
- **何时使用：** 你设置了 `mock_enabled: true` 并想用模拟响应测试此 collection。

## 一个 Spec 中的多个 Collection

一个 spec 可以有多个 collection — 例如，当一个 API 的不同服务有单独的规范文件时：

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
      - llm_title: Marine
        base_url: https://marine-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
```

## 禁用一个 Collection

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
        disable: true
```

## HTTP 客户端覆盖

所有 `http_client` 设置可以在 collection 级别覆盖。Collection 值仅对此 collection 优先于 spec 和全局值。

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        http_client:
          timeout: 120s
          headers:
            "X-Custom": "value"
          cookies:
            - name: "session"
              value: "abc123"
```
