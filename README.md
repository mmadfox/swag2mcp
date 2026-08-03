![swag2mcp](docs/logo.svg)

<p>
    <a href="https://github.com/mmadfox/swag2mcp/releases"><img src="https://img.shields.io/github/release/mmadfox/swag2mcp.svg" alt="Latest Release"></a>
    <a href="https://coveralls.io/github/mmadfox/swag2mcp?branch=main"><img src="https://coveralls.io/repos/github/mmadfox/swag2mcp/badge.svg?branch=main&v=3" alt="Coverage Status"></a>
</p>

**swag2mcp** is a local-first bridge between OpenAPI/Swagger/Postman API specifications and LLM agents via the Model Context Protocol (MCP).

Not every API speaks MCP — private endpoints, internal services, legacy systems, and third-party APIs rarely do. swag2mcp wraps any REST API in an MCP interface, giving LLM agents instant access to your entire API surface without modifying a single line of server code. Through live API calls, the LLM gains real-world knowledge to make informed decisions, automate workflows, and act on your data — not just guess.

- **16 MCP tools** for discovering, inspecting, and invoking APIs
- **Interactive TUI explorer** with full-text search
- **Zero integration code** — just point to your specs and go

---

- <a href="https://mmadfox.github.io/swag2mcp/getting-started/installation" target="_blank" rel="noopener noreferrer">Installation</a>
- <a href="https://mmadfox.github.io/swag2mcp/getting-started/quickstart" target="_blank" rel="noopener noreferrer">Quickstart</a>
- <a href="https://mmadfox.github.io/swag2mcp" target="_blank" rel="noopener noreferrer">Documentation</a>

---

<p style="font-size:1.2em;font-weight:bold">🎬 <a href="https://swag2mcp.io/concepts/example" target="_blank" rel="noopener noreferrer">Live session</a></p>

---

## Quick Start

### Install

**macOS (Homebrew):**
```bash
brew install mmadfox/tap/swag2mcp
```

**macOS / Linux (one-liner):**
```bash
curl -fsSL https://swag2mcp.io/install.sh | bash
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

