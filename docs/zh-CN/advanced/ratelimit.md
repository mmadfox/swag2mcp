# 速率限制

## 概述

swag2mcp 内置了速率限制器，防止 LLM 过于频繁地调用同一 API 端点。这可以防止意外的重复调用，并尊重 API 的速率限制。

## 工作原理

每个端点都有一个冷却期。如果 LLM 在冷却期内尝试再次调用同一端点，该调用将被拒绝并返回结构化错误。

```
t=0s  → invoke(endpoint) → 执行
t=2s  → invoke(endpoint) → 被拒绝，返回 rate_limit 错误
t=12s → invoke(endpoint) → 执行（冷却期已过）
```

### 默认行为

- **冷却时间：** 每个端点 10 秒
- **范围：** 按端点 — 调用端点 A 不影响端点 B
- **错误响应：** LLM 收到代码为 `rate_limit` 的 `LLMError`，并附带指示等待时间的消息
- **重置：** 该端点 10 秒无活动后

### 错误格式

当触发速率限制时，LLM 会收到：

```json
{
  "code": "rate_limit",
  "message": "rate limit exceeded for endpoint \"abc123\": try again in 8 seconds",
  "hint": "Wait for the cooldown period to expire, then try invoking the endpoint again. Use the search tool to find other endpoints you can call in the meantime."
}
```

LLM 可以使用此信息等待并重试，或切换到不同的端点。

### 为什么存在

- **防止意外重复调用** — LLM 可能在短时间内多次调用同一端点
- **防止 API 速率限制** — 许多 API 有自己的速率限制，触发它们会导致错误
- **节省资源** — 减少不必要的网络流量

## 配置

你可以禁用速率限制器或更改冷却间隔：

```yaml
# 完全禁用速率限制器
disable_ratelimiter: true

# 自定义冷却间隔
rate_limit_interval: 30s
```

### disable_ratelimiter

- **类型：** `bool`
- **默认值：** `false`
- **效果：** 当为 `true` 时，按端点的速率限制器被禁用。LLM 可以重复调用同一端点而无需等待。
- **何时启用：** 测试、调试，或需要快速连续多次调用同一端点时。
- **何时保持禁用（推荐）：** 生产环境。速率限制器防止意外滥用。

### rate_limit_interval

- **类型：** 持续时间（Go 格式：`10s`、`30s`、`1m`）
- **默认值：** `10s`
- **效果：** 设置对同一端点的两次调用之间的冷却期。
- **何时增加：** 具有严格速率限制的 API（例如每分钟 10 次请求）。
- **何时减少：** 你可以控制负载的内部 API。
- **示例：** `5s`、`30s`、`1m`、`2m`

## 重要说明

- **按端点跟踪** — 每个端点独立跟踪。调用一个端点不影响其他端点。
- **错误返回给 LLM** — 冷却期内的第二次调用被拒绝，返回 `rate_limit` 错误。LLM 收到冷却持续时间，可以在等待后重试。
- **无需清理** — 速率限制器自动跟踪端点，无需维护。
