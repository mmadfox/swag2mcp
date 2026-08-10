# OpenID Connect (OIDC Discovery)

## 用途

基于 OpenID Connect Discovery 的认证。swag2mcp 从 issuer 获取 OIDC 发现文档，发现令牌端点，并通过 **client credentials** 授权获取 Bearer 令牌。令牌会缓存到过期为止。

## 何时使用

- 提供 OIDC 发现文档（`.well-known/openid-configuration`）的 API
- 拥有 `client_id` + `client_secret` 的服务器到服务器集成
- 当 API 使用 OIDC 且需要自动获取令牌时

## 配置

```yaml
specs:
  - domain: my-api
    llm_title: My API
    base_url: https://api.example.com
    collections:
      - llm_title: Main
        location: https://example.com/spec.yaml
    auth:
      type: openid-connect
      config:
        issuer: "https://auth.example.com"
        client_id: "$(CLIENT_ID)"
        client_secret: "$(CLIENT_SECRET)"
        scopes:
          - openid
          - profile
```

## 参数

| 参数 | 必填 | 说明 |
|-----------|----------|-------------|
| `issuer` | 是 | OIDC issuer URL。swag2mcp 会附加 `/.well-known/openid-configuration` 来发现令牌端点 |
| `client_id` | 是 | 客户端标识符 |
| `client_secret` | 是 | 客户端密钥 |
| `scopes` | 否 | 权限列表（可选） |

## 说明

- 当前令牌过期时，swag2mcp 会自动请求新令牌
- 令牌会缓存到过期（`expires_in`）为止
- 如果服务器未提供 `expires_in`，令牌有效期为 1 小时
- `issuer`、`client_id` 和 `client_secret` 支持 `$(VAR)` 环境变量语法
- 令牌通过 **client credentials** 授权获取（服务器到服务器，无需用户参与）
- 发现文档必须包含 `token_endpoint` 字段
