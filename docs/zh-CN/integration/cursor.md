# Cursor 集成

## stdio

在 Cursor 设置中，添加 MCP 服务器：

```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp"]
    }
  }
}
```

## 使用

连接后，Cursor AI 智能体可以：

- 探索你的 API
- 查找相关端点
- 调用 API 并显示结果
- 帮助调试请求

## 其他

没有看到你的客户端？所有 MCP 集成遵循相同的模式：
- 将命令设置为 `swag2mcp`，参数为 `mcp`
- 可选地添加工作区路径：`mcp /path/to/workspace`
- 查看客户端的文档以了解确切的配置文件位置和格式

大多数 MCP 客户端支持 stdio 传输，部分支持 HTTP（SSE / Streamable HTTP）。
