# OpenCode 集成

## stdio

```json
{
  "mcp": {
    "swag2mcp": {
      "type": "local",
      "command": ["swag2mcp", "mcp"],
      "enabled": true
    }
  }
}
```

## HTTP

```json
{
  "mcp": {
    "swag2mcp": {
      "type": "local",
      "command": ["swag2mcp", "mcp", "--transport", "sse", "--http-addr", "127.0.0.1:8080"],
      "enabled": true
    }
  }
}
```

## 示例查询

连接后，你可以询问：

- "你有什么 API？"
- "显示 petstore 中的所有端点"
- "查找用于创建用户的 API"
- "调用 GET /pet/1 并显示结果"

## 其他

没有看到你的客户端？所有 MCP 集成遵循相同的模式：
- 将命令设置为 `swag2mcp`，参数为 `mcp`
- 可选地添加工作区路径：`mcp /path/to/workspace`
- 查看客户端的文档以了解确切的配置文件位置和格式

大多数 MCP 客户端支持 stdio 传输，部分支持 HTTP（SSE / Streamable HTTP）。
