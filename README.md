# swag2mcp

> ⚠️ **Work in progress** — API may change, contributions welcome.

**swag2mcp** bridges OpenAPI/Swagger/Postman API specifications with LLM agents via the Model Context Protocol (MCP).

- **For LLM agents**: 14 MCP tools for discovering, inspecting, and invoking APIs
- **For humans**: Interactive TUI explorer with full-text search
- **Zero integration code**: Just point to your specs and go

## Documentation

| Language | File |
|----------|------|
| English | [docs/README.md](docs/README.md) |
| Русский | [docs/README.ru.md](docs/README.ru.md) |
| Deutsch | [docs/README.de.md](docs/README.de.md) |
| Français | [docs/README.fr.md](docs/README.fr.md) |
| Español | [docs/README.es.md](docs/README.es.md) |
| 中文 | [docs/README.zh.md](docs/README.zh.md) |
| 日本語 | [docs/README.ja.md](docs/README.ja.md) |

## Quick Start

```bash
$ go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
$ swag2mcp init
$ swag2mcp mcp
OR
$ swag2mcp mcp --tags=project-1,work,petstore
```

## License

MIT
