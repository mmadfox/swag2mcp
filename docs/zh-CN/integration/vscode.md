# VS Code 集成

## Via .vscode/mcp.json

1. 安装 VS Code 的 MCP 扩展（例如 org.mcp 的 MCP Client 或类似扩展）。
2. 在项目根目录创建 `.vscode/mcp.json`：

```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "${workspaceFolder}"]
    }
  }
}
```

> // "${{workspaceFolder}}" 将作为工作区路径传递

3. 重新加载 VS Code 窗口（Ctrl+Shift+P → "Reload Window"）。
4. 使用 AI 助手 — 现在它将了解你的 API。

## 替代方案：通过 VS Code 设置

也可以在 `.vscode/settings.json` 中配置：

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

设置后，VS Code AI 助手可以通过 swag2mcp 与你的 API 一起工作。

## 其他

没有看到你的客户端？所有 MCP 集成遵循相同的模式：
- 将命令设置为 `swag2mcp`，参数为 `mcp`
- 可选地添加工作区路径：`mcp /path/to/workspace`
- 查看客户端的文档以了解确切的配置文件位置和格式

大多数 MCP 客户端支持 stdio 传输，部分支持 HTTP（SSE / Streamable HTTP）。
