# API Key

## 用途

通过 API 密钥进行认证。密钥可以作为 HTTP 头或 URL 查询参数发送。

## 何时使用

- 使用 API 密钥的服务
- 天气服务、地理数据、翻译 API
- 当 API 期望在头（`X-API-Key`）或查询参数（`?api_key=...`）中提供密钥时

## 配置

### 密钥在头中

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: api-key
      config:
        key: "X-API-Key"
        in: header
        value: "$(API_KEY)"
```

### 密钥在查询参数中

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: api-key
      config:
        key: "api_key"
        in: query
        value: "$(API_KEY)"
```

## 参数

| 参数 | 必需 | 描述 |
|------|------|------|
| `key` | 是 | 头或查询参数的名称 |
| `in` | 是 | 密钥放置位置：`header` 或 `query` |
| `value` | 是 | 密钥值 |

## 说明

- 在 `header` 模式下，密钥作为 HTTP 头添加
- 在 `query` 模式下，密钥作为 URL 参数添加
- 将值存储在环境变量中：`value: "$(MY_API_KEY)"`
