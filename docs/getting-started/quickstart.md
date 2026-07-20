# Quick Start

Get swag2mcp running in 2 minutes.

## 1. Initialize

### Home directory (recommended)

One-time setup for your entire system. Config is stored in your home folder.

::: code-group

```bash [macOS / Linux]
swag2mcp init
# Creates ~/.swag2mcp/swag2mcp.yaml
```

```powershell [Windows]
swag2mcp.exe init
# Creates %USERPROFILE%\.swag2mcp\swag2mcp.yaml
```

:::

### Project directory

For an isolated workspace inside your project.

::: code-group

```bash [macOS / Linux]
mkdir -p ./swag2mcp && swag2mcp init ./swag2mcp
```

```powershell [Windows]
mkdir ./swag2mcp; swag2mcp.exe init ./swag2mcp
```

:::

### From ZIP

If you have a ready-made workspace (e.g., from a colleague):

```bash
swag2mcp import --from-zip workspace.zip
```

## 2. LLM Client / IDE Configuration

Configure your IDE to connect to swag2mcp. The IDE will start the MCP server automatically when needed.

::: code-group

```json [OpenCode]
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

```json [Claude Desktop]
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp"]
    }
  }
}
```

```json [Crush]
{
  "mcp": {
    "swag2mcp": {
      "type": "stdio",
      "command": "swag2mcp",
      "args": ["mcp"]
    }
  }
}
```

:::

For other IDEs (Cursor, VS Code, JetBrains) see the [Integration guide](../integration/opencode.md).

> If you initialized the workspace at a custom path (e.g. `./swag2mcp`), use the full path in the command:
> `"command": ["swag2mcp", "mcp", "/absolute/path/to/swag2mcp"]`

> **After any config change, restart the MCP server** for the changes to take effect.

## 3. Start MCP Server

### stdio (default) — for local IDE

Nothing to configure. Your IDE starts swag2mcp automatically via the config above.

```bash
swag2mcp mcp
```

### SSE / Streamable HTTP — for remote access

```bash
swag2mcp mcp --transport sse --http-addr :8080
```

Or configure in `swag2mcp.yaml`:

```yaml
mcp:
  transport: sse
  addr: ":8080"
  path: "/mcp"
```

See [MCP Server reference](../configuration/mcp-server.md) for all flags.

### Filter specs by tags

```bash
swag2mcp mcp --tags weather,public
```

Only specs with matching tags will be available to the LLM.

### Verify it's working

After connecting, ask your LLM agent:

```bash
"What MCP tools do you support?"
```

If the agent lists swag2mcp tools (`spec_list`, `search`, `invoke`, etc.) — everything is working.

### Example queries to try

| Ask your agent | What happens |
|-------|-------------|
| "What's the weather in New York?" | `invoke` — calls Open-Meteo forecast API |
| "What's the current BTC price?" | `invoke` — calls Binance ticker API |
| "Tell me a dad joke" | `invoke` — calls icanhazdadjoke API |
| "Show me Pikachu" | `invoke` — calls PokéAPI by name |
| "Who is Rick Sanchez?" | `invoke` — calls Rick and Morty character API |
| "What's the air quality in Beijing?" | `invoke` — calls Open-Meteo air quality API |
| "How high are the waves near Portugal?" | `invoke` — calls Open-Meteo marine API |
| "Search for jokes about dogs" | `invoke` — calls dadjoke search endpoint |
| "List all Pokémon" | `invoke` — calls PokéAPI list endpoint |
| "What's the elevation of Mount Everest?" | `invoke` — calls Open-Meteo elevation API |

## 5. Install Agent Skills (recommended)

Install the swag2mcp skills to teach your AI agent all commands, flags, config format, and real-world examples.

Ask your agent:

```bash
"Add the swag2mcp-cli skill from https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-cli/SKILL.md"
"Add the swag2mcp-format skill from https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md"
```

> Some IDEs require a restart after adding skills.

## What's Next?

- [Concepts](../concepts/overview.md) — understand the architecture
- [Configuration](../configuration/config-file.md) — customize settings
- [CLI Commands](../cli/overview.md) — full command reference
