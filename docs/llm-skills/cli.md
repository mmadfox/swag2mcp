# swag2mcp-cli

The **swag2mcp-cli** skill gives your LLM the complete swag2mcp CLI reference — every command, flag, argument, and configuration option. With this skill, the LLM can answer "how do I..." questions accurately instead of guessing.

## What it covers

All 13 CLI commands:

| Command | Purpose |
|---------|---------|
| `init` | Initialize a workspace and config |
| `add` | Add a spec or collection |
| `delete` | Remove a spec or collection |
| `ls` | List configured specs |
| `run` | Start the API explorer TUI |
| `validate` | Validate the configuration file |
| `clean` | Clear cached data |
| `update` | Update cache from configuration |
| `mcp` | Start the MCP server |
| `version` | Show version information |
| `info` | Show runtime information |
| `import` | Import spec files or restore workspace from backup |
| `export` | Export the workspace to a ZIP file |

Plus all flags, config file structure, auth methods, and advanced options.

## Direct link

<https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md>

## Install via LLM agent

Copy this request to your AI-powered IDE:

```
Create the directory .agents/skills/swag2mcp-cli/ and add the skill from
https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md
```

## Manual install

```bash
mkdir -p .agents/skills/swag2mcp-cli
curl -o .agents/skills/swag2mcp-cli/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md
```

## Restart required

After adding the skill, restart your LLM client or IDE (see [Overview](overview.md#restart-required)).
