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
swag2mcp init
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
mkdir ./swag2mcp; swag2mcp init ./swag2mcp
```

:::

### From ZIP

If you have a ready-made workspace (e.g., from a colleague):

```bash
swag2mcp import --from-zip workspace.zip
```

## 2. Add an API

### Option A — via CLI with YAML

Add a spec with all fields (domain, title, instruction, base URL, collections):

::: code-group

```bash [Home directory]
swag2mcp add spec --yaml - <<EOF
domain: meteo
llm_title: Open-Meteo Weather APIs
llm_instruction: "Use this API for weather forecasts and climate data"
base_url: https://api.open-meteo.com
collections:
  - llm_title: Forecast
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
  - llm_title: Air Quality
    base_url: https://air-quality-api.open-meteo.com
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
EOF
```

```bash [Project directory]
swag2mcp add spec ./swag2mcp --yaml - <<EOF
domain: meteo
llm_title: Open-Meteo Weather APIs
llm_instruction: "Use this API for weather forecasts and climate data"
base_url: https://api.open-meteo.com
collections:
  - llm_title: Forecast
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
  - llm_title: Air Quality
    base_url: https://air-quality-api.open-meteo.com
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
EOF
```

:::

### Option B — edit YAML file directly

Open the config file and add your spec:

::: code-group

```text [Home directory]
~/.swag2mcp/swag2mcp.yaml
```

```text [Project directory]
./swag2mcp/swag2mcp.yaml
```

:::

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    llm_instruction: "Use this API for weather forecasts and climate data"
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

### Add more collections to an existing spec

::: code-group

```bash [Home directory]
swag2mcp add collection --yaml - <<EOF
spec_domain: meteo
llm_title: Air Quality
base_url: https://air-quality-api.open-meteo.com
location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
EOF
```

```bash [Project directory]
swag2mcp add collection ./swag2mcp --yaml - <<EOF
spec_domain: meteo
llm_title: Air Quality
base_url: https://air-quality-api.open-meteo.com
location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
EOF
```

:::

## 3. Start MCP Server

```bash
swag2mcp mcp
```

Output:

```
MCP server listening on http://127.0.0.1:8080/mcp
```

## 4. Test It

In another terminal:

```bash
curl -X POST http://127.0.0.1:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"spec_list","arguments":{}}}'
```

## 5. LLM Client Configuration

::: code-group

```json [OpenCode]
{
  "mcp": {
    "swag2mcp": {
      "type": "local",
      "command": ["swag2mcp", "mcp"]
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

```json [Cursor]
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp"]
    }
  }
}
```

:::

## What's Next?

- [Concepts](../concepts/overview.md) — understand the architecture
- [Configuration](../configuration/config-file.md) — customize settings
- [CLI Commands](../cli/overview.md) — full command reference
