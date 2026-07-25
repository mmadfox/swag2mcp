# info

## 用途

以 **JSON** 格式显示 swag2mcp 运行时的全面摘要。包括版本、工作区路径、spec 摘要、HTTP 客户端设置、MCP 传输配置、认证方法和模拟模式状态。

## 何时使用

- 你想要工作区的机器可读概览
- 你需要检查运行时配置以进行调试
- 你想查看有多少 spec 和端点处于活动状态
- 你需要验证 HTTP 客户端或 MCP 传输设置

## 语法

```bash
swag2mcp info [path]
```

## 参数

| 参数 | 位置 | 必需 | 描述 |
|------|------|------|------|
| `path` | 1 | 否 | 工作区目录。如果省略，通过路径解析规则解析。 |

## 标志

无。

## 工作原理

```bash
swag2mcp info
swag2mcp info ./my-workspace
```

## 输出

输出是具有以下结构的 JSON 对象：

```json
{
  "version": "v1.2.0",
  "workspace": "~/.swag2mcp",
  "uptime": "2h 15m",
  "specs": {
    "total": 4,
    "active": 3,
    "disabled": 1,
    "collections": 6,
    "endpoints": 42
  },
  "http_client": {
    "timeout": "30s",
    "max_response_size": "2 KB",
    "proxy": "none",
    "follow_redirects": true,
    "max_redirects": 10,
    "randomize": false
  },
  "mcp": {
    "transport": "stdio",
    "addr": ":8080",
    "path": "/mcp"
  },
  "auth_methods": ["bearer", "api-key"],
  "mock_enabled": false
}
```

## 命令后验证

在启动 MCP 服务器之前，使用 `info` 确认工作区加载正确且所有 spec 处于活动状态。

## 细节

- **自动初始化：** 如果不存在配置文件，`info` 会自动先运行初始化向导。
- **仅 JSON：** 输出始终是 JSON。对于人类可读的输出，使用 `ls`。
- **`max_response_size`：** 以人类可读格式显示（例如 `"1 KB"`、`"2 MB"`）。
- **无全文索引：** `info` 禁用全文索引，因为它只需要配置和 spec 元数据。
