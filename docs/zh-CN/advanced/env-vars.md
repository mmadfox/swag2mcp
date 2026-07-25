# 环境变量

## 概述

swag2mcp 支持在配置文件中使用 `$(VAR_NAME)` 语法进行环境变量替换。这让你可以将敏感数据（令牌、密码、密钥）保留在 YAML 文件之外。

## 工作原理

当 swag2mcp 启动时，它会扫描配置中的 `$(VAR_NAME)` 模式，并将其替换为相应环境变量的值。

```yaml
specs:
  - domain: meteo
    auth:
      type: bearer
      config:
        token: "$(API_TOKEN)"
```

如果设置了环境变量 `API_TOKEN`，它将被替换。如果未设置，该值将变为空。

## `$(VAR)` 的解析位置

| 字段 | 示例 |
|------|------|
| Auth `token`（bearer） | `token: "$(API_TOKEN)"` |
| Auth `username` / `password`（basic、digest） | `password: "$(API_PASSWORD)"` |
| Auth `client_id` / `client_secret`（oauth2-cc、oauth2-pwd） | `client_secret: "$(OAUTH_SECRET)"` |
| Auth `api_key` / `secret_key`（hmac） | `api_key: "$(BINANCE_API_KEY)"` |
| Auth `domain`（script） | `domain: "$(AUTH_DOMAIN)"` |
| MCP 服务器令牌 | `token: "$(MCP_TOKEN)"` |
| HTTP 客户端头 | `"X-API-Key": "$(API_KEY)"` |
| HTTP 客户端 cookie 值 | `value: "$(SESSION_TOKEN)"` |

## `$(VAR)` 不会被解析的位置

- 基础 URL（`base_url`）
- Collection 位置（`location`）
- Spec 域名（`domain`）

## 示例

```bash
export API_TOKEN="eyJhbGciOiJIUzI1NiIs..."
export MCP_TOKEN="my-secret-token"

swag2mcp mcp
```

## 安全最佳实践

- **永远不要**将密钥直接存储在 YAML 文件中
- 使用环境变量或外部密钥管理器
- 如果 YAML 文件包含任何硬编码的密钥，将其添加到 `.gitignore`
- 在 shell 配置文件、IDE 配置或部署管道中设置环境变量

## 语法细节

- `$(VAR_NAME)` — 标准语法
- `$( VAR_NAME )` — 括号内的空格允许且会被修剪
- `$()` — 空变量名返回原始字符串不变
- 嵌套的 `$(...)` 模式不会被解析
