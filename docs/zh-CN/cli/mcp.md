# mcp

## 用途

启动 **MCP（模型上下文协议）服务器** — LLM 集成的主要模式。这是你运行以让 AI 智能体（Claude、Cursor、OpenCode 等）通过 16 个 MCP 工具访问你的 API 的方式。

## 何时使用

- 你想将 LLM 智能体连接到你的 API
- 你正在配置 IDE（VS Code、Cursor、JetBrains）或桌面应用（Claude Desktop）
- 你需要通过 MCP 协议暴露你的 API
- 你正在集成前测试 MCP 服务器

## 语法

```bash
swag2mcp mcp [path] [flags]
```

## 参数

| 参数 | 位置 | 必需 | 描述 |
|------|------|------|------|
| `path` | 1 | 否 | 工作区目录。如果省略，通过路径解析规则解析。 |

## 标志

| 标志 | 简写 | 类型 | 默认值 | 描述 |
|------|------|------|--------|------|
| `--transport` | | `string` | `"stdio"` | MCP 传输：`stdio`、`sse`、`streamable-http` |
| `--http-addr` | | `string` | `":8080"` | HTTP 服务器地址（用于 `sse` 和 `streamable-http`） |
| `--http-path` | | `string` | `"/mcp"` | MCP 处理程序的 HTTP 路径 |
| `--auth-token` | | `string` | `""` | HTTP 传输认证的 Bearer 令牌 |
| `--logfile` | `-f` | `string` | `""` | 日志文件路径。如果未设置，日志输出到 stderr。 |
| `--disable-llm-auth` | | `bool` | `true` | 从 MCP 工具列表中移除 `auth` 工具 |
| `--dump-dir` | | `string` | `""` | 用于调试的 HTTP 请求转储目录 |
| `--tags` | `-t` | `string` | `""` | 按标签过滤 spec（逗号分隔） |

## 工作原理

### stdio 传输（默认）

当 MCP 服务器作为子进程由 LLM 客户端（IDE、Claude Desktop 等）启动时使用。服务器通过标准输入/输出进行通信。

```bash
swag2mcp mcp
```

### SSE 传输

用于基于 HTTP 通信的服务器发送事件传输。需要 MCP 握手序列。

```bash
swag2mcp mcp --transport sse --http-addr :8080
```

### Streamable HTTP 传输

支持流式响应的现代 HTTP 传输。

```bash
swag2mcp mcp --transport streamable-http --http-addr 0.0.0.0:8080
```

### 带认证

使用 bearer 令牌保护 HTTP 端点：

```bash
swag2mcp mcp --transport sse --http-addr :8080 --auth-token "my-secret"
```

### 带标签过滤

仅加载具有特定标签的 spec：

```bash
swag2mcp mcp --tags=public
```

### 启用 auth 工具（调试模式）

允许 LLM 通过 `auth` 工具请求新令牌：

```bash
swag2mcp mcp --disable-llm-auth=false
```

### 带请求转储目录

保存所有 HTTP 请求以进行调试：

```bash
swag2mcp mcp --dump-dir ./dumps
```

## MCP HTTP 传输 — 握手协议

使用 `sse` 或 `streamable-http` 时，MCP 协议需要特定的握手。在初始化之前，工具调用将失败：

```
步骤 1：POST /mcp → {"method":"initialize", ...}
步骤 2：POST /mcp → {"method":"notifications/initialized"}
步骤 3：POST /mcp → {"method":"tools/list", ...}   ← 现在可以工作
```

### 健康检查

无需初始化即可工作：

```bash
curl http://localhost:8080/health
# → {"status":"ok","version":"v1.2.0"}
```

## IDE 配置示例

### VS Code（`.vscode/settings.json` 或全局设置）

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

### Cursor / Windsurf（`~/.cursor/mcp.json` 或项目 `.cursor/mcp.json`）

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

### Claude Desktop（macOS 上为 `~/Library/Application Support/Claude/claude_desktop_config.json`）

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

### JetBrains IDE（设置 → 工具 → MCP）

- 名称：`swag2mcp`
- 命令：`swag2mcp`
- 参数：`mcp /absolute/path/to/.swag2mcp`

> **始终使用工作区目录的绝对路径**。相对路径可能因 IDE 的工作目录而失败。

## 输出

成功时，服务器打印：

```
MCP server listening on http://127.0.0.1:8080/mcp
```

## 细节

- **无自动初始化：** 如果配置文件不存在，`mcp` 返回错误：`"configuration not found at <path>"`。先运行 `init`。
- **`--disable-llm-auth`（默认：`true`）：** 启用时，`auth` 工具从 MCP 工具列表中完全移除。LLM 无法看到或请求令牌。认证仍然有效 — 令牌通过标准配置机制获取，而不是通过 LLM。此模式推荐用于**生产环境**。对于**调试**或使用短期令牌时，设置 `--disable-llm-auth=false` 让 LLM 通过 `auth` 工具请求新令牌。
- **YAML 配置回退：** 如果未显式设置 CLI 标志，则从 `swag2mcp.yaml` 中的 `mcp` 部分取值（如果存在）。这允许你在配置文件中配置服务器，而不是每次都传递标志。
- **响应清理：** 启动时，超过 48 小时的响应会自动从 `responses/` 目录中删除。
- **路径解析警告：** 当省略 `[path]` 时，`mcp` 首先在当前目录中搜索 `swag2mcp.yaml`，然后回退到 `~/.swag2mcp/`。如果你从错误的目录运行命令，可能会加载与预期不同的工作区。**作为服务运行或在 IDE 配置中时，始终显式指定 `[path]`。**
