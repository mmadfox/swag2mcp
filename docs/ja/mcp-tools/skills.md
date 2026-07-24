# Skills

## Customizing Output Format

Every swag2mcp MCP tool returns structured JSON data. How this data is **presented** to the user depends on the LLM's formatting skill — and you can control it completely.

### The default format skill

swag2mcp ships with a built-in formatting skill that defines compact, human-readable markdown for every tool response:

[swag2mcp-format SKILL.md](https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md)

This skill covers all 19 MCP tools with:
- Tight tables for lists (specs, collections, tags, endpoints)
- Inline headers for detail views
- Compact schema representation for `inspect`
- Consistent styling across all responses

### Why skills matter

The same data can be presented in radically different ways depending on the skill:

| Style | Example output |
|-------|---------------|
| **Compact tables** (default) | `GET /pet/{petId}` — Find pet by ID |
| **Verbose** | `Method: GET, Path: /pet/{petId}, Summary: Find pet by ID, Deprecated: false` |
| **Minimal** | `GET /pet/{petId}` |
| **Technical** | `GET /pet/{petId} → 200: Pet object, 404: Not found` |
| **Custom** | Any format you can describe |

### Creating your own skill

You can write your own formatting skill by describing the exact output format you want. The skill is a markdown file with formatting rules for each tool. Here are some ideas:

- **JSON output** — return raw JSON for machine consumption
- **CSV-style** — tabular data for spreadsheet import
- **Diagram-friendly** — Mermaid or ASCII diagrams of API structure
- **Minimal** — just method and path, nothing else
- **Documentation-style** — full descriptions, examples, and notes

### The only limit is the model

The quality of the formatted output depends entirely on the LLM's ability to follow your formatting rules. A well-written skill with clear examples produces consistent, reliable output. A vague skill produces inconsistent results.

You can:
- Use the default skill as-is
- Fork it and tweak the formatting to your taste
- Write your own from scratch
- Switch between skills depending on the task

### How to use a skill

Skills are loaded by the LLM client (OpenCode, Cursor, Claude Desktop, etc.) as part of its system prompt or agent configuration. Refer to your client's documentation for how to attach a skill file.

For OpenCode, skills are configured in `opencode.json`:

```json
{
  "skills": [
    {
      "name": "swag2mcp-format",
      "sourceURL": "https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md"
    }
  ]
}
```
