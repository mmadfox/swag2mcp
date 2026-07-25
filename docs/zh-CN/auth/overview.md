# 认证

## 概述

swag2mcp 支持 **9 种认证方法**，用于处理需要授权的 API。你在配置文件中配置一次 — 之后，通过 `invoke` 的每个 API 调用都会自动包含正确的令牌和头。

### 在哪里配置

认证在 `swag2mcp.yaml` 的 **spec** 级别设置：

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
        token: "my-token"
```

### 工作原理

- 你在配置中指定认证类型和参数
- 当你调用 `invoke` 时，swag2mcp 自动将其应用于每个请求
- 你**不需要**在调用 API 之前请求令牌 — 它会自动发生
- 如果令牌过期（OAuth2、Script），swag2mcp 会自动刷新

### 环境变量

敏感数据（令牌、密码、密钥）可以使用 `$(VAR_NAME)` 语法存储在环境变量中：

```yaml
auth:
  type: bearer
  config:
    token: "$(MY_API_TOKEN)"
```

swag2mcp 在启动时替换 `MY_API_TOKEN` 的值。

### MCP auth 工具

LLM 智能体可以通过 `auth` MCP 工具检索令牌或头 — 例如，用于构建 curl 命令或向用户显示。

在**生产环境**中，应使用 `--disable-llm-auth`（默认启用）禁用此工具，以便 LLM 永远无法访问令牌。

### 方法

| 方法 | 描述 | 最适合 |
|------|------|--------|
| [`none`](/auth/none) | 无需认证 | 公共 API |
| [`basic`](/auth/basic) | HTTP Basic（用户名 + 密码） | 旧版 API、简单认证 |
| [`bearer`](/auth/bearer) | Bearer Token（JWT、令牌） | 现代 REST API |
| [`api-key`](/auth/api-key) | 头或查询参数中的 API 密钥 | 使用 API 密钥的服务 |
| [`digest`](/auth/digest) | HTTP Digest（用户名 + 密码） | 旧版 API，比 Basic 更安全 |
| [`hmac`](/auth/hmac) | HMAC-SHA256 签名（Binance 风格） | 加密货币交易所 |
| [`oauth2-cc`](/auth/oauth2-cc) | OAuth2 客户端凭证 | 服务器到服务器、微服务 |
| [`oauth2-pwd`](/auth/oauth2-pwd) | OAuth2 密码授权 | 带用户登录的应用程序 |
| [`script`](/auth/script) | 用于获取令牌的外部脚本 | 任何自定义认证方案 |
