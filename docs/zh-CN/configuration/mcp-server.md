# MCP 服务器

MCP 服务器是 LLM 智能体的主要交互点。它将所有配置的 API 作为 MCP 工具暴露给 LLM 调用。

## 配置

```yaml
mcp:
  transport: stdio
  addr: ":8080"
  path: "/mcp"
  auth:
    token: ""
```

## 传输方式

三种传输方式可用：

| 传输方式 | 描述 | 何时使用 |
|----------|------|----------|
| `stdio` | 标准输入/输出 | 本地 LLM 客户端（VS Code、Cursor、Claude Desktop） |
| `sse` | 服务器发送事件 | 远程客户端、基于 HTTP 的通信 |
| `streamable-http` | 带流式传输的 HTTP | Web 客户端、现代 MCP 客户端 |

### stdio（默认）

LLM 客户端将 swag2mcp 作为子进程运行。通信通过标准输入和输出进行。不需要网络端口。

```yaml
mcp:
  transport: stdio
```

```bash
swag2mcp mcp
```

### SSE

用于基于 HTTP 通信的服务器发送事件传输。MCP 服务器监听 HTTP 端口，LLM 客户端远程连接。

```yaml
mcp:
  transport: sse
  addr: "127.0.0.1:8080"
  path: "/mcp"
```

```bash
swag2mcp mcp --transport sse --http-addr 127.0.0.1:8080
```

### Streamable HTTP

支持流式响应的现代 HTTP 传输。类似于 SSE，但使用不同的协议。

```yaml
mcp:
  transport: streamable-http
  addr: "127.0.0.1:8080"
  path: "/mcp"
```

```bash
swag2mcp mcp --transport streamable-http --http-addr 0.0.0.0:8080
```

## 参数

### transport

- **类型：** `string`
- **默认值：** `"stdio"`
- **选项：** `stdio`、`sse`、`streamable-http`
- **效果：** 确定 MCP 服务器如何与 LLM 客户端通信。

### addr

- **类型：** `string`
- **默认值：** `":8080"`
- **描述：** SSE 和 Streamable HTTP 传输的监听地址。格式：`host:port`。
- **示例：** `":8080"`、`"127.0.0.1:8080"`、`"0.0.0.0:9000"`

### path

- **类型：** `string`
- **默认值：** `"/mcp"`
- **描述：** MCP 端点的 URL 路径。LLM 客户端发送请求到 `http://&lt;addr&gt;&lt;path&gt;`。
- **示例：** `"/mcp"`、`"/api/mcp"`、`"/v1/mcp"`

### auth.token

- **类型：** `string`
- **默认值：** `""`（无认证）
- **描述：** HTTP 传输认证的 Bearer 令牌。设置后，LLM 客户端必须在每个请求中包含 `Authorization: Bearer &lt;token&gt;`。
- **注意：** 支持 `$(ENV_VAR)` 解析。

## HTTP 认证

使用 bearer 令牌保护 MCP HTTP 端点：

```yaml
mcp:
  auth:
    token: "my-secret-token"
```

或通过 CLI 标志：

```bash
swag2mcp mcp --auth-token "my-secret-token"
```

## 健康检查

MCP 服务器提供无需 MCP 初始化即可工作的健康检查端点：

```bash
curl http://127.0.0.1:8080/health
# {"status":"ok","version":"v1.2.0"}
```

## 启动标志

CLI 标志覆盖 YAML 配置。如果未设置标志，则使用 YAML 中 `mcp` 部分的值作为回退。

| 标志 | 类型 | 默认值 | 描述 |
|------|------|--------|------|
| `--transport` | string | `"stdio"` | 传输类型：`stdio`、`sse`、`streamable-http` |
| `--http-addr` | string | `":8080"` | HTTP 服务器地址（用于 SSE 和 Streamable HTTP） |
| `--http-path` | string | `"/mcp"` | MCP 处理程序的 URL 路径 |
| `--auth-token` | string | `""` | HTTP 传输认证的 Bearer 令牌 |
| `--logfile` | string | `""` | 日志文件路径（未设置时输出到 stderr） |
| `--disable-llm-auth` | bool | `true` | 从 MCP 工具列表中移除 `auth` 工具 |
| `--dump-dir` | string | `""` | 用于调试的 HTTP 请求转储目录 |
| `--tags` | string | `""` | 按标签过滤 spec（逗号分隔） |
