# 任意 MCP 客户端

swag2mcp 是一个 **MCP 服务器**（Model Context Protocol）。这意味着它适用于**任何 MCP 客户端** — 不仅是本节列出的那些。只要您的编辑器、IDE 或代理支持 MCP 协议，就可以将 swag2mcp 连接到它。

## 通用模式

每个 MCP 客户端都使用相同的基本设置。将 swag2mcp 添加为 MCP 服务器：

- **命令：** `swag2mcp`
- **参数：** `mcp`（可选加工作区路径：`mcp /path/to/workspace`）

```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "/path/to/workspace"]
    }
  }
}
```

配置文件的确切位置和格式（JSON、TOML、GUI 设置）因客户端而异 — **请查阅您客户端的 MCP 文档**，了解将其放在何处。

## 传输方式

- **stdio** — 随处可用；大多数 MCP 客户端都支持
- **HTTP（SSE / Streamable HTTP）** — 支持提供 HTTP 传输选项的客户端

有关传输标志，请参阅 [`mcp` 命令](/zh-CN/cli/mcp) 参考。

## 已测试的集成

| 客户端 | 指南 |
|--------|------|
| OpenCode | [OpenCode](/zh-CN/integration/opencode) |
| Cursor | [Cursor](/zh-CN/integration/cursor) |
| Claude Desktop | [Claude Desktop](/zh-CN/integration/claude) |
| VS Code | [VS Code](/zh-CN/integration/vscode) |
| Crush | [Crush](/zh-CN/integration/crush) |

> 如果您的客户端不在列表中，**并不**意味着它不受支持。只要它使用 MCP 协议，就使用上面的通用模式，并按照客户端的操作手册查找配置位置。
