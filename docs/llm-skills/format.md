# swag2mcp-format

The **swag2mcp-format** skill teaches your LLM how to display swag2mcp MCP tool responses in a clean, compact, human-readable markdown format. Without this skill, the LLM decides on its own how to format responses — often resulting in verbose or inconsistent output.

## What it covers

All swag2mcp MCP tools:

- `spec_list`, `spec_by_id` — specs overview and details
- `collection_by_spec`, `collection_by_id` — collections with tags
- `tag_by_spec`, `tag_by_collection`, `tag_by_id` — tag listings
- `endpoint_by_spec`, `endpoint_by_collection`, `endpoint_by_tag`, `endpoint_by_id` — endpoint lists
- `search` — search results
- `inspect` — full operation details with compact schemas
- `invoke` — API call results
- `auth` — authentication info
- `info` — runtime information

## Direct link

<https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md>

## Install via LLM agent

Copy this request to your AI-powered IDE:

```
Create the directory .agents/skills/swag2mcp-format/ and add the skill from
https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md
```

## Manual install

```bash
mkdir -p .agents/skills/swag2mcp-format
curl -o .agents/skills/swag2mcp-format/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md
```

## Restart required

After adding the skill, restart your LLM client or IDE (see [Overview](overview.md#restart-required)).
