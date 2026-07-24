# 全局设置

全局设置是 `swag2mcp.yaml` 中的顶级配置块。它们适用于所有 spec，除非在 spec 或 collection 级别被覆盖。

## 结构

```yaml
http_client:
  # 所有 API 调用的 HTTP 客户端设置

mcp:
  # MCP 服务器设置

mock_enabled: false
mock_auth:
  oauth2_port: 9090
  digest_port: 9091
  hmac_port: 9092

disable_ratelimiter: false
rate_limit_interval: 10s
```

## HTTP 客户端

控制 swag2mcp 如何向 API 发出 HTTP 请求：超时、响应大小限制、代理、头、cookie、重定向和用户代理。这些设置级联到 spec 和 collection。

所有参数和示例请参见[HTTP 客户端](./http-client)。

## MCP 服务器

控制 MCP 服务器如何与 LLM 智能体通信：传输类型（stdio、SSE、Streamable HTTP）、地址、路径和可选的 bearer 令牌认证。

所有参数、传输方式和启动标志请参见[MCP 服务器](./mcp-server)。

## 模拟服务器

模拟服务器基于 OpenAPI 模式生成模拟 API 响应。用于无需访问真实 API 即可进行测试。

```yaml
mock_enabled: true
mock_auth:
  oauth2_port: 9090
  digest_port: 9091
  hmac_port: 9092
```

### mock_enabled

- **类型：** `bool`
- **默认值：** `false`
- **效果：** 当为 `true` 时，swag2mcp 为所有配置了 `base_mock_url` 的 spec 启动模拟服务器。每个 collection 必须设置 `base_mock_url`。
- **何时启用：** 你想在不进行真实 HTTP 调用的情况下测试 API 集成。模拟服务器基于 OpenAPI 模式返回模拟数据。

### mock_auth

模拟认证服务器的端口配置。这些用于在模拟服务器上测试认证方法（OAuth2、Digest、HMAC）。

| 字段 | 类型 | 默认值 | 描述 |
|------|------|--------|------|
| `oauth2_port` | int | `9090` | 模拟 OAuth2 令牌服务器的端口（1024-65535） |
| `digest_port` | int | `9091` | 模拟 Digest 认证服务器的端口（1024-65535） |
| `hmac_port` | int | `9092` | 模拟 HMAC 认证服务器的端口（1024-65535） |

## 速率限制器

速率限制器防止 LLM 过于频繁地调用同一 API 端点。默认情况下，每个端点每 10 秒可被调用一次。

```yaml
disable_ratelimiter: false
rate_limit_interval: 10s
```

### disable_ratelimiter

- **类型：** `bool`
- **默认值：** `false`
- **效果：** 当为 `true` 时，按端点的速率限制器完全禁用。LLM 可以重复调用同一端点而无需等待。
- **何时启用：** 测试、调试，或需要快速连续多次调用同一端点时。
- **何时保持禁用（推荐）：** 生产环境。速率限制器防止意外滥用并尊重 API 速率限制。

### rate_limit_interval

- **类型：** 持续时间（Go 格式：`10s`、`30s`、`1m`）
- **默认值：** `10s`
- **效果：** 设置 LLM 在两次调用同一端点之间必须等待的时间。
- **何时更改：** 对于具有严格速率限制的 API 增加。对于你可以控制负载的内部 API 减少。
- **范围：** 任何有效的持续时间（例如 `5s`、`30s`、`1m`、`2m`）。

## 级联

全局设置可以在 spec 和 collection 级别被覆盖。所有 `http_client` 设置（超时、代理、用户代理、重定向、响应大小、随机化器、头、cookie）可以在 spec 和 collection 级别被覆盖。

```
全局 (http_client, mock_enabled, disable_ratelimiter, rate_limit_interval)
    ↓ 覆盖（仅 http_client）
Spec (specs[].http_client)
    ↓ 覆盖（仅 http_client）
Collection (specs[].collections[].http_client)
```

详情请参见[配置级联](./cascade)。
