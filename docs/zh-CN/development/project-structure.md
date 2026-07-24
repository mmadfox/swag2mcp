# 项目结构

```
swag2mcp/
├── cmd/
│   ├── swag2mcp/          # 主二进制文件
│   │   └── main.go
│   └── swag2mcp-mock/     # 模拟服务器
│       └── main.go
├── internal/
│   ├── auth/              # 9 种认证方法
│   ├── cache/             # 规范缓存
│   ├── commands/          # 13 个 CLI 命令（cobra）
│   ├── config/            # YAML 配置
│   ├── env/               # 环境变量
│   ├── httpclient/        # HTTP 客户端
│   ├── id/                # MD5 ID 生成
│   ├── index/             # 全文搜索（bluge）
│   ├── model/             # 数据模型
│   ├── reader/            # 大响应读取
│   ├── server/
│   │   ├── mcp/           # MCP 服务器（19 个工具）
│   │   └── mockserver/    # 模拟服务器
│   ├── service/           # 业务逻辑
│   ├── spec/              # 规范解析器
│   ├── tui/               # TUI 界面
│   └── workspace/         # 工作区管理
├── specs/                 # 示例规范
├── tests/                 # 集成测试
├── docs/                  # 文档
├── examples/              # 配置示例
└── playground/            # 开发沙箱
```

## 关键包

| 包 | 描述 |
|------|------|
| `auth` | 9 种认证方法 |
| `cache` | 基于磁盘的缓存，带 TTL |
| `commands` | Cobra CLI 命令 |
| `config` | 带级联的 YAML 配置 |
| `httpclient` | 可配置的 HTTP 客户端 |
| `index` | 全文搜索（bluge） |
| `server/mcp` | MCP 服务器（3 种传输方式） |
| `service` | 业务逻辑（核心） |
| `spec` | OpenAPI/Swagger/Postman 解析器 |
| `tui` | Bubbletea TUI |
| `workspace` | 文件管理 |
