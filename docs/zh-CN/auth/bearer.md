# Bearer Auth

## 用途

Bearer 令牌认证 — 现代 REST API 最常用的方法。令牌通过 `Authorization: Bearer &lt;token&gt;` 头发送。

## 何时使用

- 现代 REST API
- JWT（JSON Web Tokens）
- OAuth2 访问令牌（当令牌已获取时）
- 任何接受 Bearer Token 的 API

## 配置

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: bearer
      config:
        token: "eyJhbGciOiJIUzI1NiIs..."
```

## 参数

| 参数 | 必需 | 描述 |
|------|------|------|
| `token` | 是 | Bearer 令牌（JWT、OAuth2 令牌等） |

## 说明

- 令牌是静态的 — 如果过期，你需要在配置中手动更新
- 对于自动令牌刷新，使用 `oauth2-cc` 或 `oauth2-pwd`
- 将令牌存储在环境变量中：`token: "$(API_TOKEN)"`
