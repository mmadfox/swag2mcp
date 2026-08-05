# Cursor 集成

## stdio

### Via Settings UI

1. 打开 Cursor 设置 (Cmd+, / Ctrl+,)
2. Go to **MCP 服务器**
3. Click **添加新服务器**
4. Fill in:
   - **名称:** `swag2mcp`
   - **类型:** `command`
   - **命令:** `swag2mcp mcp`
5. Click **保存**

### 通过配置文件

In `~/.cursor/mcp.json`:

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

After connecting, Cursor AI Agent can:

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
