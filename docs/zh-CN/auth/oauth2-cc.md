# OAuth2 Client Credentials

## 用途

OAuth2 客户端凭证授权 — 服务器到服务器通信的认证。应用程序使用其 client_id 和 client_secret 获取令牌，无需用户参与。

## 何时使用

- 微服务和服务器到服务器集成
- 机器对机器通信
- 当 API 使用 OAuth2 并且你有 client_id + client_secret 时

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
      type: oauth2-cc
      config:
        client_id: "$(CLIENT_ID)"
        client_secret: "$(CLIENT_SECRET)"
        token_url: "https://auth.example.com/oauth/token"
        scopes:
          - read
          - write
```

## 参数

| 参数 | 必需 | 描述 |
|------|------|------|
| `client_id` | 是 | 客户端标识符 |
| `client_secret` | 是 | 客户端密钥 |
| `token_url` | 是 | 令牌端点 URL |
| `scopes` | 否 | 权限列表（可选） |

## 说明

- swag2mcp 在当前令牌过期时自动请求新令牌
- 令牌被缓存直到其过期时间（`expires_in`）
- 如果服务器未提供 `expires_in`，令牌被视为有效 1 小时
- 所有参数可以存储在环境变量中
