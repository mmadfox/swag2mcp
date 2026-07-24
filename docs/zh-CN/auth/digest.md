# Digest Auth

## 用途

HTTP Digest 访问认证 — 比 Basic Auth 更安全的替代方案。密码不以明文发送，而是使用 MD5 哈希。

## 何时使用

- 仅支持 Digest 的旧版 API
- 需要在不以明文发送密码的情况下进行认证时
- 内部企业系统

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
      type: digest
      config:
        username: "admin"
        password: "$(PASSWORD)"
```

## 参数

| 参数 | 必需 | 描述 |
|------|------|------|
| `username` | 是 | 用户名 |
| `password` | 是 | 密码 |

## 说明

- swag2mcp 首先发送不带认证的请求，从服务器接收挑战（HTTP 401），计算响应，然后使用 `Authorization: Digest ...` 头重试
- 挑战被缓存 5 分钟 — 后续请求不需要额外的往返
- 将密码存储在环境变量中：`password: "$(API_PASSWORD)"`
