# Skills for LLM

Skills are markdown files that teach your LLM agent how to work with swag2mcp more effectively. They are loaded as part of the agent's system prompt and give the LLM precise instructions for formatting responses and understanding CLI commands.

## Available skills

| Skill | Description | Download |
|-------|-------------|----------|
| **swag2mcp-format** | Formats all MCP tool responses into compact, human-readable markdown tables | [SKILL.md](https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md) |
| **swag2mcp-cli** | Full CLI reference — tells the LLM every command, flag, and config option | [SKILL.md](https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md) |

## Why skills matter

Without a formatting skill, the LLM decides on its own how to display tool results — often verbose and inconsistent. The format skill ensures every response follows the same clean pattern: tight tables for lists, inline headers for details, and compact schemas.

The CLI skill lets the LLM answer any "how do I..." question about swag2mcp commands without guessing or making things up.

## Install via LLM agent

Copy this request to your AI-powered IDE (OpenCode, Cursor, Claude Desktop, VS Code, etc.):

```
Add the swag2mcp skills to my project:

1. Create the directory .agents/skills/swag2mcp-format/ and add the skill from https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md
2. Create the directory .agents/skills/swag2mcp-cli/ and add the skill from https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md
```

The agent will download both skill files and place them in the correct locations.

## Manual install

If your LLM client doesn't support agent-based setup, download the files manually:

```bash
# Create the skills directories
mkdir -p .agents/skills/swag2mcp-format
mkdir -p .agents/skills/swag2mcp-cli

# Download the skill files
curl -o .agents/skills/swag2mcp-format/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md
curl -o .agents/skills/swag2mcp-cli/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md
```

## Configure your LLM client

Each LLM client and IDE has its own way of installing skills. The example below is for **OpenCode** — check your client's documentation for the correct method.

For OpenCode, add the skills to `opencode.json`:

```json
{
  "skills": [
    {
      "name": "swag2mcp-format",
      "sourceURL": "https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md"
    },
    {
      "name": "swag2mcp-cli",
      "sourceURL": "https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md"
    }
  ]
}
```

## Restart required

**After adding skills, restart your LLM client or IDE.** Some tools load skills only at startup. If the skills don't seem to take effect, try:

- **OpenCode**: Restart the application or run the `opencode` command again
- **Cursor**: Close and reopen the window (`Cmd+Shift+W` / `Ctrl+Shift+W`)
- **Claude Desktop**: Quit and relaunch the app
- **VS Code**: Reload the window (`Ctrl+Shift+P` → "Developer: Reload Window")
