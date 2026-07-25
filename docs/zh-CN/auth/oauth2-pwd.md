# OAuth2 Password Grant

## 用途

OAuth2 资源所有者密码授权 — 使用用户的用户名和密码进行认证。适用于用户信任应用程序使用其凭据的第一方应用程序。

## 何时使用

- 第一方应用程序（移动端、Web）
- 与 Keycloak 和类似身份提供商的集成
- 当 API 支持 OAuth2 Password Grant 时

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
      type: oauth2-pwd
      config:
        client_id: "$(CLIENT_ID)"
        client_secret: "$(CLIENT_SECRET)"
        username: "$(USERNAME)"
        password: "$(PASSWORD)"
        token_url: "https://auth.example.com/oauth/token"
        scopes:
          - openid
          - profile
```

## 参数

| 参数 | 必需 | 描述 |
|------|------|------|
| `client_id` | 是 | 客户端标识符 |
| `username` | 是 | 用户名 |
| `password` | 是 | 密码 |
| `token_url` | 是 | 令牌端点 URL |
| `client_secret` | 否 | 客户端密钥（可选，用于公共客户端） |
| `scopes` | 否 | 权限列表（可选） |

## 说明

- `client_secret` 是可选的 — 支持**公共客户端**（例如 Keycloak）
- swag2mcp 在令牌过期时自动刷新
- 令牌被缓存直到过期
- 所有参数可以存储在环境变量中
