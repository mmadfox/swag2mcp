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
        request_format: form
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
| `request_format` | 否 | 请求体格式：`form`（默认，`application/x-www-form-urlencoded`）或 `json`（`application/json`） |

## 说明

- `client_secret` 是可选的 — 支持**公共客户端**（例如 Keycloak）
- swag2mcp 在令牌过期时自动刷新
- 令牌被缓存直到过期
- `client_id`、`client_secret`、`username` 和 `password` 支持 `$(VAR)` 语法用于环境变量
- `token_url` 和 `scopes` 按原样使用（不解析环境变量）
- `request_format: json` 以 JSON 正文发送令牌请求 — 当令牌端点需要 `Content-Type: application/json` 时使用
