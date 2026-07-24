# Basic Auth

## 用途

HTTP Basic 认证 — 使用用户名和密码进行认证的最简单方式。

## 何时使用

- 仅支持 Basic Auth 的旧版 API
- 无需复杂令牌的简单认证
- 内部服务

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
      type: basic
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

- 密码以 Base64 编码通过 `Authorization: Basic ...` 头发送 — 这**不是加密**。始终使用 HTTPS。
- 将密码存储在环境变量中：`password: "$(MY_PASSWORD)"`
