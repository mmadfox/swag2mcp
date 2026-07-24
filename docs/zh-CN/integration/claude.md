# Claude Desktop 集成

## stdio

在 `claude_desktop_config.json` 中：

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

## 自定义工作区

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

## 使用

重启 Claude Desktop 后，你可以：

- "显示所有 API 的列表"
- "查找创建订单的端点"
- "调用莫斯科的天气 API"

## 其他

没有看到你的客户端？所有 MCP 集成遵循相同的模式：
- 将命令设置为 `swag2mcp`，参数为 `mcp`
- 可选地添加工作区路径：`mcp /path/to/workspace`
- 查看客户端的文档以了解确切的配置文件位置和格式

大多数 MCP 客户端支持 stdio 传输，部分支持 HTTP（SSE / Streamable HTTP）。
