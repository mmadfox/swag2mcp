# swag2mcp-format

Der **swag2mcp-format** Skill bringt Ihrem LLM bei, swag2mcp MCP-Tool-Antworten in einem sauberen, kompakten, lesbaren Markdown-Format anzuzeigen. Ohne diesen Skill entscheidet der LLM selbst, wie Antworten formatiert werden — oft ausführlich und inkonsistent.

## Was wird abgedeckt

Alle swag2mcp MCP-Tools:

- `spec_list`, `spec_by_id` — Spezifikationsübersicht und -details
- `collection_by_spec`, `collection_by_id` — Sammlungen mit Tags
- `tag_by_spec`, `tag_by_collection`, `tag_by_id` — Tag-Listen
- `endpoint_by_spec`, `endpoint_by_collection`, `endpoint_by_tag`, `endpoint_by_id` — Endpunkt-Listen
- `search` — Suchergebnisse
- `inspect` — vollständige Operationsdetails mit kompakten Schemata
- `invoke` — API-Aufrufergebnisse
- `auth` — Authentifizierungsinformationen
- `info` — Laufzeitinformationen

## Direkter Link

<https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md>

## Installation via LLM-Agent

Kopieren Sie diese Anfrage in Ihre KI-gestützte IDE:

```
Erstelle das Verzeichnis .agents/skills/swag2mcp-format/ und füge den Skill von
https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md hinzu
```

## Manuelle Installation

```bash
mkdir -p .agents/skills/swag2mcp-format
curl -o .agents/skills/swag2mcp-format/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md
```

## Neustart erforderlich

Starten Sie nach dem Hinzufügen des Skills Ihren LLM-Client oder Ihre IDE neu (siehe [Übersicht](overview.md#neustart-erforderlich)).
