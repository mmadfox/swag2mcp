# 常见问题

## 通用

### 什么是 swag2mcp，它解决了什么问题？

swag2mcp 通过模型上下文协议（MCP）将 OpenAPI/Swagger/Postman API 规范与 LLM 智能体连接起来。无需为每个 API 编写自定义代码来连接 AI 智能体，你只需在 YAML 文件中配置一次，LLM 就能获得 19 个工具来发现、检查和调用你的 API。

### 它与其他 API 到 LLM 的工具有什么不同？

- **无需编码** — 在 YAML 中配置 API，无需集成代码
- **19 个 MCP 工具** — 从发现到调用再到大响应处理的完整工具包
- **9 种认证方法** — 适用于任何 API 认证方案
- **全文搜索** — 基于 bluge 的跨所有端点的搜索
- **TUI 浏览器** — 用于浏览和测试的交互式终端界面
- **模拟服务器** — 无需真实 API 调用即可测试

### 支持哪些 API 规范格式？

OpenAPI 3.x、Swagger 2.0 和 Postman Collections v2.1。

### spec 和 collection 有什么区别？

**Spec** 代表一个逻辑 API 服务（例如"Open-Meteo 天气 API"）。**Collection** 是一个 OpenAPI/Swagger/Postman 文件。一个 spec 可以有多个 collection — 例如，当一个 API 的不同服务（天气预报、空气质量、海洋）有单独的规范文件时。

### 支持哪些 MCP 传输方式？

三种传输方式：`stdio`（默认，用于本地 LLM 客户端）、`sse`（服务器发送事件，用于远程客户端）和 `streamable-http`（现代 HTTP 流式传输）。

### 我可以将 swag2mcp 与任何 LLM 一起使用吗？

可以，任何支持 MCP 协议的 LLM 客户端：Claude Desktop、VS Code、Cursor、Windsurf、JetBrains IDE、OpenCode 等。

## 安装

### 如何安装 swag2mcp？

```bash
# 选项 1：从 GitHub Releases 下载
# 访问 https://github.com/mmadfox/swag2mcp/releases/latest
# 下载适用于你的操作系统和架构的压缩包

# 选项 2：使用 Go 安装
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

### 我需要安装 Go 吗？

不需要。预构建的二进制文件适用于 Linux（amd64、arm64）、macOS（amd64、arm64）和 Windows（amd64），可在 [GitHub Releases 页面](https://github.com/mmadfox/swag2mcp/releases) 获取。

### 如何安装模拟服务器？

模拟服务器是一个独立的二进制文件：

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

或者从 GitHub Releases 下载 `swag2mcp-mock_<version>_<os>_<arch>.tar.gz`。

## 快速入门

### 如何快速开始？

```bash
# 1. 初始化工作区
mkdir -p .swag2mcp && swag2mcp init ./.swag2mcp

# 2. 启动 MCP 服务器（初始化后包含公共示例规范）
swag2mcp mcp
```

执行 `init` 后，工作区已包含多个公共示例规范（icanhazdadjoke、Open-Meteo、Binance、PokéAPI）。你可以立即启动 MCP 服务器 — 无需手动添加规范。

如果你想添加自己的 API：

```bash
swag2mcp add spec --yaml - <<EOF
domain: dadjoke
llm_title: icanhazdadjoke API
base_url: https://icanhazdadjoke.com
collections:
  - llm_title: Jokes
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
EOF
```

### 如何将 swag2mcp 连接到我的 IDE？

**VS Code**（`.vscode/settings.json`）：
```json
{
  "mcp": {
    "servers": {
      "swag2mcp": {
        "command": "swag2mcp",
        "args": ["mcp", "/absolute/path/to/.swag2mcp"]
      }
    }
  }
}
```

**Cursor**（`~/.cursor/mcp.json`）：
```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "/absolute/path/to/.swag2mcp"]
    }
  }
}
```

**Claude Desktop**（`claude_desktop_config.json`）：
```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "/absolute/path/to/.swag2mcp"]
    }
  }
}
```

始终使用工作区目录的绝对路径。

## 配置

### 配置文件在哪里？

默认位置：`~/.swag2mcp/swag2mcp.yaml`。你也可以在任何目录中创建它，并将路径传递给命令。

### 如何添加 API？

```bash
# 交互模式
swag2mcp add spec

# 使用 YAML（推荐用于脚本）
swag2mcp add spec --yaml - <<EOF
domain: my-api
llm_title: My API
base_url: https://api.example.com/v1
collections:
  - llm_title: Main
    location: https://example.com/spec.yaml
EOF
```

### 如何向现有 spec 添加 collection？

```bash
swag2mcp add collection --yaml - <<EOF
spec_domain: meteo
llm_title: Air Quality
location: https://example.com/air-quality.yaml
EOF
```

### 如何临时禁用一个 spec？

在 spec 配置中设置 `disable: true`。该 spec 将不会被加载或索引。

### 我可以过滤加载哪些 spec 吗？

可以，使用 `--tags` 标志：`swag2mcp mcp --tags=public`。只有具有匹配标签的 spec 才会被加载。

### 如何使用环境变量存储密钥？

在 auth 字段中使用 `$(VAR_NAME)` 语法：

```yaml
auth:
  type: bearer
  config:
    token: "$(MY_API_TOKEN)"
```

在启动前设置环境变量：`export MY_API_TOKEN="eyJhbGci..."`

## 认证

### 支持哪些认证方法？

九种方法：`none`、`basic`、`bearer`、`digest`、`hmac`、`oauth2-cc`（客户端凭证）、`oauth2-pwd`（密码授权）、`api-key` 和 `script`。

### 如何传递令牌？

通过配置文件或环境变量：

```yaml
auth:
  type: bearer
  config:
    token: "$(MY_TOKEN)"
```

### 在调用 invoke 之前需要调用 auth 吗？

不需要。`invoke` 工具会自动应用 spec 配置中的认证。只有在你想向用户显示令牌时（例如用于 curl 命令），才需要使用 `auth` MCP 工具。

### 为什么 auth 工具没有显示？

`auth` 工具默认是禁用的（`--disable-llm-auth=true`）。这是生产环境的安全措施。要启用它：`swag2mcp mcp --disable-llm-auth=false`。

### OAuth2 令牌如何刷新？

OAuth2 客户端凭证和密码授权令牌在过期时会自动刷新。Bearer 令牌是静态的，必须手动更新。

## MCP 服务器

### 如何启动 MCP 服务器？

```bash
# 默认（stdio 传输）
swag2mcp mcp

# 使用 HTTP 传输
swag2mcp mcp --transport sse --http-addr :8080
```

### 如何更改端口？

```bash
swag2mcp mcp --transport sse --http-addr 0.0.0.0:9090
```

### 如何保护 MCP HTTP 端点？

设置 bearer 令牌：

```bash
swag2mcp mcp --transport sse --http-addr :8080 --auth-token "my-secret"
```

LLM 客户端必须在每个请求中包含 `Authorization: Bearer my-secret`。

### HTTP 传输的 MCP 握手是什么？

对于 SSE 和 Streamable HTTP 传输，MCP 协议需要三步握手：

```
步骤 1：POST /mcp → {"method":"initialize", ...}
步骤 2：POST /mcp → {"method":"notifications/initialized"}
步骤 3：POST /mcp → {"method":"tools/list", ...}  ← 现在可以工作
```

在初始化之前，工具调用将失败。

## 使用

### 如何搜索端点？

使用 `search` MCP 工具或 TUI（`swag2mcp run`）。搜索支持字段过滤器（`method:GET`、`tag:pets`）、模糊搜索、通配符和布尔运算符。

### 如何调用 API？

LLM 使用 `invoke` MCP 工具。始终先检查端点以了解所需参数：

```
inspect(endpointId: "...")  → 了解契约
invoke(endpointId: "...", parameters: {...})  → 发起调用
```

### 如果响应太大怎么办？

超过 `max_response_size`（默认 1 MB）的响应会保存到磁盘。LLM 收到文件引用，可以使用 `response_outline`、`response_compress` 和 `response_slice` 工具进行探索。

### 速率限制器如何工作？

每个端点有 10 秒的冷却时间。如果 LLM 在 10 秒内两次调用同一端点，第二次调用将被静默阻止。你可以在配置中禁用或调整此设置。

### 我可以在不进行真实 API 调用的情况下测试吗？

可以，使用模拟服务器：

```bash
swag2mcp-mock mockserver
```

它基于 OpenAPI 模式生成模拟响应。

## 工作区管理

### 如何备份我的配置？

```bash
swag2mcp export --output ~/backups/swag2mcp-2026-07-24.zip
```

### 如何迁移到另一台机器？

```bash
# 在旧机器上
swag2mcp export --output swag2mcp.zip

# 复制 ZIP，然后在新机器上
swag2mcp import --from-zip swag2mcp.zip
```

### 如何更新规范文件？

```bash
swag2mcp update
```

这会重新验证配置、清除缓存并重新下载所有规范文件。

### 如何清理磁盘空间？

```bash
swag2mcp clean
```

删除缓存的规范文件和保存的 API 响应。旧响应（超过 48 小时）也会在 MCP 服务器启动时自动清理。

## TUI

### 什么是 TUI，如何使用它？

TUI（终端用户界面）是一个交互式 API 浏览器。使用 `swag2mcp run` 启动。它有三种模式：搜索（全文搜索）、浏览（树形导航：Spec → Collection → Tag → Endpoint）和认证（查看令牌）。

### 键盘快捷键有哪些？

| 按键 | 操作 |
|------|------|
| `↑/↓` | 导航 |
| `Enter` | 选择 |
| `Esc` | 返回 |
| `Tab` | 切换模式 |
| `/` | 搜索 |
| `N/P` | 下一页/上一页 |
| `q` | 退出 |

## 高级

### 我可以使用代理吗？

可以，在 `http_client.proxy` 中配置：

```yaml
http_client:
  proxy:
    url: "http://proxy.company.com:8080"
    username: "$(PROXY_USER)"
    password: "$(PROXY_PASS)"
    bypass:
      - "localhost"
      - "*.internal.com"
```

### 我可以添加自定义认证方法吗？

可以，在 `internal/auth/` 中实现 `Authenticator` 接口，并在配置解析器中注册。详情请参阅开发部分。

### 我可以添加自定义 MCP 工具吗？

可以，向 `Svc` 接口添加方法，在服务层实现它，添加处理程序并注册。详情请参阅开发部分。

### `swag2mcp` 和 `swag2mcp-mock` 有什么区别？

`swag2mcp` 是主二进制文件，包含 CLI 命令和 MCP 服务器。`swag2mcp-mock` 是一个独立的二进制文件，用于启动模拟服务器，无需真实 API 调用即可进行测试。
