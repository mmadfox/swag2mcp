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
        request_format: form
```

## 参数

| 参数 | 必需 | 描述 |
|------|------|------|
| `client_id` | 是 | 客户端标识符 |
| `client_secret` | 是 | 客户端密钥 |
| `token_url` | 是 | 令牌端点 URL |
| `scopes` | 否 | 权限列表（可选） |
| `request_format` | 否 | 请求体格式：`form`（默认，`application/x-www-form-urlencoded`）或 `json`（`application/json`） |

## 说明

- swag2mcp 在当前令牌过期时自动请求新令牌
- 令牌被缓存直到其过期时间（`expires_in`）
- 如果服务器未提供 `expires_in`，令牌被视为有效 1 小时
- `client_id` 和 `client_secret` 支持 `$(VAR)` 语法用于环境变量
- `token_url` 和 `scopes` 按原样使用（不解析环境变量）
- `request_format: json` 以 JSON 正文发送令牌请求 — 当令牌端点需要 `Content-Type: application/json` 时使用
