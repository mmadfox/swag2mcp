# swag2mcp

<p>
    <a href="https://github.com/mmadfox/swag2mcp/releases"><img src="https://img.shields.io/github/release/mmadfox/swag2mcp.svg" alt="Latest Release"></a>
    <a href="https://pkg.go.dev/github.com/mmadfox/swag2mcp?tab=doc"><img src="https://godoc.org/github.com/mmadfox/swag2mcp?status.svg" alt="GoDoc"></a>
    <a href="https://github.com/mmadfox/swag2mcp/actions"><img src="https://github.com/mmadfox/swag2mcp/actions/workflows/test.yml/badge.svg?branch=main" alt="Build Status"></a>
    <a href="https://coveralls.io/github/mmadfox/swag2mcp?branch=main"><img src="https://coveralls.io/repos/github/mmadfox/swag2mcp/badge.svg?branch=main&v=3" alt="Coverage Status"></a>
</p>

> ⚠️ **Work in progress** — API may change, contributions welcome.

<a href="https://www.youtube.com/watch?v=9CcvwmfTkds" target="_blank">
  <img src="assets/cover.png" width="600" alt="Preview">
</a>

**swag2mcp** bridges OpenAPI/Swagger/Postman API specifications with LLM agents via the Model Context Protocol (MCP).

- **16 MCP tools** for discovering, inspecting, and invoking APIs
- **Interactive TUI explorer** with full-text search
- **Zero integration code** — just point to your specs and go

---

## Table of Contents

- [Installation](#installation)
- [Example LLM Queries](#example-llm-queries)
- [Documentation](#documentation)
- [Integration](#integration)
- [CLI Commands](#cli-commands)
- [Mock Server](#mock-server)
- [Configuration Examples](#configuration-examples)
- [License](#license)

---

## Installation ([For AI Agents](#for-ai-agents), [For Humans](#for-humans))

### For AI Agents

This path is for users who want to use swag2mcp through an AI agent (Opencode, Cursor, Claude, Copilot, Crush, etc.).

**Step 1 — Install the `swag2mcp-cli` skill**

The skill teaches your agent all commands, flags, config format, and real-world examples. Without it, the agent won't know how to use swag2mcp.

You can ask your agent to install the skill by providing this link:

**< PROMPTS >**
```bash
1. "Add the swag2mcp-cli skill from https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-cli/SKILL.md"
2. "Add the swag2mcp-format skill from https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md"
```

Or configure it manually — refer to your IDE's documentation on how to add custom skills.

> **Note:** Some IDEs require a restart after adding a new skill for it to take effect.

**Step 2 — Choose a setup method**

**Option A — Fully automated (agent does everything)**

**Step 1** — Tell your agent:

**< PROMPTS >**
```bash
1. "Load the swag2mcp-cli skill"
2. "Load the swag2mcp-format skill"
3. "Set up swag2mcp"
```

The agent will download and install `swag2mcp`, then create a workspace with ready-to-use specs.

> **Note:** If the `swag2mcp` MCP server is already running, restart it so the agent can pick up the new workspace. Otherwise, proceed to the next step.

**Step 2** — Configure your IDE (see [Integration](#integration) below) — the IDE will start `swag2mcp` automatically.

**Step 3** — Use MCP tools to discover and invoke APIs.

**Option B — You install, agent connects**

```bash
# Download the latest release from:
https://github.com/mmadfox/swag2mcp/releases/latest
```

Then configure your IDE (see [Integration](#integration) below).

### For Humans

Prefer the command line? Here's the manual setup.

**Option 1 — Download from GitHub Releases (recommended)**

1. Open [github.com/mmadfox/swag2mcp/releases/latest](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Find the archive for your system:

   | OS | Architecture | Archive |
   |----|-------------|---------|
   | Linux | x86_64 | `swag2mcp_<version>_linux_amd64.tar.gz` |
   | Linux | ARM64 | `swag2mcp_<version>_linux_arm64.tar.gz` |
   | macOS | Intel | `swag2mcp_<version>_darwin_amd64.tar.gz` |
   | macOS | Apple Silicon | `swag2mcp_<version>_darwin_arm64.tar.gz` |
   | Windows | x86_64 | `swag2mcp_<version>_windows_amd64.zip` |

3. Download and install:

   **Linux / macOS:**
   ```bash
   tar -xzf swag2mcp_<version>_<os>_<arch>.tar.gz
   sudo mv swag2mcp /usr/local/bin/
   swag2mcp --version
   ```

   **Windows (PowerShell):**
   ```powershell
   Expand-Archive swag2mcp_<version>_windows_amd64.zip -DestinationPath .
   move swag2mcp.exe C:\Windows\System32\
   swag2mcp --version
   ```

4. (Optional) Repeat for mock server — download `swag2mcp-mock_<version>_<os>_<arch>.tar.gz`

**Option 2 — Install with Go**

If you have Go installed:

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

**After installation:**

```bash
swag2mcp init

Via CLI:
swag2mcp add spec --yaml 'domain: binance
llm_title: Binance Market Data API
base_url: https://api.binance.com
collections:
  - llm_title: Market Data
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/binance.yaml
```

Or edit `swag2mcp.yaml` manually:

```yaml
specs:
  - domain: binance
    llm_title: Binance Market Data API
    base_url: https://api.binance.com
    collections:
      - llm_title: Market Data
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/binance.yaml
```

```bash
# 4. Validate and update
swag2mcp validate
swag2mcp update

# 5. Add to your IDE

Configure your IDE to start swag2mcp automatically (see [Integration](#integration) below). The IDE will run `swag2mcp mcp` when needed.

# 6. Explore interactively (optional)
swag2mcp run
```

---

## Example LLM Queries

Before asking your agent, make sure swag2mcp MCP is connected in your IDE (see [Integration](#integration) above). The agent needs access to the MCP tools to discover and invoke APIs.

It is recommended to always start with this question to verify swag2mcp is connected:

**< PROMPT >**
```bash
"What MCP tools do you support?"
```

If the agent lists swag2mcp tools (spec_list, search, invoke, etc.) — everything is working. If not, ask the agent to connect swag2mcp:

**< PROMPT >**
```bash
"Load swag2mcp MCP"
```

After setup, try asking your agent:

| < PROMPT > | What happens |
|-------|-------------|
| "List my configured APIs" | `spec_list` — lists available specs |
| "Show me all available APIs" | `spec_list` — lists petstore, binance, dadjoke, pokeapi |
| "What endpoints does Binance have?" | `endpoint_by_spec` — shows 4 market data endpoints |
| "Find endpoints related to pets" | `search("pet")` — finds petstore endpoints |
| "What tags are in the Petstore API?" | `tag_by_spec` — shows "pets" tag |
| "Show me the GET /pets endpoint details" | `inspect` — shows parameters and response schema |
| "Get the current BTC price from Binance" | `invoke` — real API call to Binance |
| "Get a random dad joke" | `invoke` — calls icanhazdadjoke API |

---

## Documentation

| Language | File |
|----------|------|
| English | [docs/guide.md](docs/guide.md) |
| Русский | [docs/guide.ru.md](docs/guide.ru.md) |
| Deutsch | [docs/guide.de.md](docs/guide.de.md) |
| Français | [docs/guide.fr.md](docs/guide.fr.md) |
| Español | [docs/guide.es.md](docs/guide.es.md) |
| 中文 | [docs/guide.zh.md](docs/guide.zh.md) |
| 日本語 | [docs/guide.ja.md](docs/guide.ja.md) |

### Agent Skills

| Skill | Description |
|-------|-------------|
| [swag2mcp-cli](.agents/skills/swag2mcp-cli/SKILL.md) | Complete CLI reference — all commands, flags, config format, real-world examples |
| [swag2mcp-format](https://github.com/mmadfox/skills#swag2mcp-format) | Formats MCP tool responses as human-readable markdown |

---

## Integration

swag2mcp speaks the Model Context Protocol (MCP) and works with any MCP-compatible client. All settings (tags, transport, auth) are configured in `swag2mcp.yaml` — see [examples](examples/).

> **Note:** If you initialized the workspace at a custom path (e.g. `swag2mcp init ./my-project`), you must specify the path when starting the MCP server: `swag2mcp mcp ./my-project`. The IDE configuration must also use the full path to the config file.

### Local (stdio) — agent on the same machine

| Client | Config File | Content |
|--------|-------------|---------|
| **OpenCode** | `opencode.json` | `{"mcp":{"swag2mcp":{"type":"local","command":["swag2mcp","mcp"]}}}` |
| **Cursor** | `.cursor/mcp.json` | `{"mcpServers":{"swag2mcp":{"command":"swag2mcp","args":["mcp"]}}}` |
| **Claude Desktop** | `claude_desktop_config.json` | `{"mcpServers":{"swag2mcp":{"command":"swag2mcp","args":["mcp"]}}}` |
| **VS Code** | `.vscode/mcp.json` | `{"servers":{"swag2mcp":{"type":"stdio","command":"swag2mcp","args":["mcp"]}}}` |
| **Crush** | `crush.json` | `{"mcp":{"swag2mcp":{"type":"stdio","command":"swag2mcp","args":["mcp"]}}}` |

### Remote (HTTP) — agent in cloud / different machine

Start the server with HTTP transport:

```bash
swag2mcp mcp --transport streamable-http --http-addr :8080 --auth-token my-secret
```

Or configure in `swag2mcp.yaml`:

```yaml
mcp:
  transport: streamable-http
  addr: ":8080"
  path: "/mcp"
  auth_token: $(MCP_AUTH_TOKEN)
```

| Client | Config File | Content |
|--------|-------------|---------|
| **OpenCode** | `opencode.json` | `{"mcp":{"swag2mcp":{"type":"remote","url":"http://localhost:8080/mcp","headers":{"Authorization":"Bearer ${MCP_AUTH_TOKEN}"}}}}` |
| **VS Code** | `.vscode/mcp.json` | `{"servers":{"swag2mcp":{"type":"http","url":"http://localhost:8080/mcp"}}}` |

> **Health check** (works without MCP handshake):
> ```bash
> curl http://localhost:8080/health
> # → {"status":"ok","version":"v1.1.3"}
> ```

---

## CLI Commands

All commands that accept `[path]` use the same path resolution:

```
swag2mcp <command>          → ~/.swag2mcp/
swag2mcp <command> ./       → {cwd}/.swag2mcp/
swag2mcp <command> path/to  → {cwd}/path/to/.swag2mcp/
```

| Command | Description |
|---------|-------------|
| `init` | Initialize workspace and configuration |
| `add spec` / `add collection` | Add a specification or collection |
| `delete spec` / `delete collection` | Delete a specification or collection |
| `ls` | List specifications and collections |
| `run` | Interactive TUI API explorer |
| `validate` | Validate configuration |
| `clean` | Remove cache and responses |
| `update` | Validate, clear cache, re-cache all specs |
| `mcp` | Start MCP server |
| `version` / `--version` | Print version |
| `info` | Show runtime information |
| `import` | Import spec files |
| `export` | Export workspace as ZIP |

See [docs/guide.md](docs/guide.md) for full reference.

---

## Mock Server

**swag2mcp-mock** is a built-in mock server that generates random API responses
based on your OpenAPI/Swagger schemas — no real backend required.

Use it for development, testing, or when the real API is not available.

```bash
# Install
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest

# Start mock server (reads the same config as swag2mcp)
swag2mcp-mock

# Add swag2mcp to your IDE (see Integration above) — invoke will use mock URLs
```

To enable, add to your config:

```yaml
mock_enabled: true
specs:
  - domain: petstore
    collections:
      - location: specs/petstore.json
        base_mock_url: localhost:8080
```

---

## Configuration Examples

Browse ready-to-use configuration examples:

| Category | Example | Description |
|----------|---------|-------------|
| **Basics** | [minimal-config](examples/minimal-config) | Minimal configuration — one spec, one collection, no auth |
| | [full-config](examples/full-config) | Complete configuration with all features |
| **Auth** | [no-auth](examples/auth/no-auth) | No authentication |
| | [basic-auth](examples/auth/basic-auth) | HTTP Basic Authentication |
| | [bearer-auth](examples/auth/bearer-auth) | Bearer Token Authentication |
| | [digest-auth](examples/auth/digest-auth) | HTTP Digest Authentication |
| | [oauth2-client-credentials](examples/auth/oauth2-client-credentials) | OAuth2 Client Credentials Grant |
| | [oauth2-password](examples/auth/oauth2-password) | OAuth2 Password Grant |
| | [api-key-header](examples/auth/api-key-header) | API Key in HTTP Header |
| | [api-key-query](examples/auth/api-key-query) | API Key in Query Parameter |
| | [script-auth](examples/auth/script-auth) | Script-Based Authentication |
| | [hmac-auth](examples/auth/hmac-auth) | HMAC-SHA256 Authentication |
| **Spec Features** | [llm-metadata](examples/spec-features/llm-metadata) | LLM titles and instructions |
| | [disable-spec](examples/spec-features/disable-spec) | Disabling specs and collections |
| | [tags-filtering](examples/spec-features/tags-filtering) | Tag-based filtering with `--tags` |
| | [custom-headers](examples/spec-features/custom-headers) | Custom HTTP headers |
| | [multiple-collections](examples/spec-features/multiple-collections) | Multiple collections per spec |
| | [collection-override](examples/spec-features/collection-override) | Collection-level overrides |
| | [http-client-config](examples/spec-features/http-client-config) | HTTP client configuration (headers, cookies, timeout, redirects) |
| | [proxy-config](examples/spec-features/proxy-config) | Proxy configuration (SOCKS5, HTTP, HTTPS, bypass) |
| | [random-client](examples/spec-features/random-client) | Random browser-like headers |
| **MCP Transport** | [stdio](examples/mcp-transport/stdio) | Default stdio transport |
| | [sse](examples/mcp-transport/sse) | SSE transport with HTTP and bearer token auth |
| | [streamable-http](examples/mcp-transport/streamable-http) | Streamable HTTP transport with HTTP and bearer token auth |
| **Mock Server** | [mock-server](examples/mock-server) | Mock server with random data generation and auth mock |

---

## License

MIT
