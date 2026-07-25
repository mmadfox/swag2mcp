# Project Structure

```
swag2mcp/
├── cmd/
│   ├── swag2mcp/          # Main binary
│   │   └── main.go
│   └── swag2mcp-mock/     # Mock server
│       └── main.go
├── internal/
│   ├── auth/              # 9 auth methods
│   ├── cache/             # Spec caching
│   ├── commands/          # 13 CLI commands (cobra)
│   ├── config/            # YAML configuration
│   ├── env/               # Environment variables
│   ├── httpclient/        # HTTP client
│   ├── id/                # MD5 ID generation
│   ├── index/             # Full-text search (bluge)
│   ├── model/             # Data models
│   ├── reader/            # Large response reading
│   ├── server/
│   │   ├── mcp/           # MCP server (19 tools)
│   │   └── mockserver/    # Mock server
│   ├── service/           # Business logic
│   ├── spec/              # Spec parsers
│   ├── tui/               # TUI interface
│   └── workspace/         # Workspace management
├── specs/                 # Sample specs
├── tests/                 # Integration tests
├── docs/                  # Documentation
├── examples/              # Config examples
└── playground/            # Development sandbox
```

## Key Packages

| Package | Description |
|---------|-------------|
| `auth` | 9 authentication methods |
| `cache` | Disk-based caching with TTL |
| `commands` | Cobra CLI commands |
| `config` | YAML config with cascade |
| `httpclient` | Configurable HTTP client |
| `index` | Full-text search (bluge) |
| `server/mcp` | MCP server (3 transports) |
| `service` | Business logic (core) |
| `spec` | OpenAPI/Swagger/Postman parsers |
| `tui` | Bubbletea TUI |
| `workspace` | File management |
