![swag2mcp](docs/logo.svg)

<p>
    <a href="https://github.com/mmadfox/swag2mcp/releases"><img src="https://img.shields.io/github/release/mmadfox/swag2mcp.svg" alt="Latest Release"></a>
    <a href="https://coveralls.io/github/mmadfox/swag2mcp?branch=main"><img src="https://coveralls.io/repos/github/mmadfox/swag2mcp/badge.svg?branch=main&v=3" alt="Coverage Status"></a>
</p>

**swag2mcp** is a local-first bridge between OpenAPI/Swagger/Postman API specifications and LLM agents via the Model Context Protocol (MCP).

- **16 MCP tools** for discovering, inspecting, and invoking APIs
- **Interactive TUI explorer** with full-text search
- **Zero integration code** — just point to your specs and go

---

- <a href="https://mmadfox.github.io/swag2mcp/getting-started/installation" target="_blank" rel="noopener noreferrer">Installation</a>
- <a href="https://mmadfox.github.io/swag2mcp/getting-started/quickstart" target="_blank" rel="noopener noreferrer">Quickstart</a>
- <a href="https://mmadfox.github.io/swag2mcp" target="_blank" rel="noopener noreferrer">Documentation</a>

---

## Quick Start

### Install

**macOS / Linux (one-liner):**
```bash
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install.sh | bash
```

**Windows (Scoop):**
```powershell
scoop bucket add mmadfox https://github.com/mmadfox/scoop-bucket
scoop install mmadfox/swag2mcp
```

**All platforms (go install):**
```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

**All platforms (Docker):**
```bash
docker pull ghcr.io/mmadfox/swag2mcp:latest
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp
```

> For detailed installation instructions, see [Installation](https://mmadfox.github.io/swag2mcp/getting-started/installation).

