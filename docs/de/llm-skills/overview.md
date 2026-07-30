# Skills for LLM

Skills sind Markdown-Dateien, die Ihrem LLM-Agenten beibringen, effektiver mit swag2mcp zu arbeiten. Sie werden als Teil des System-Prompts geladen und geben dem LLM präzise Anweisungen zur Formatierung von Antworten und zum Verständnis von CLI-Befehlen.

## Verfügbare Skills

| Skill | Beschreibung | Download |
|-------|-------------|----------|
| **swag2mcp-format** | Formatiert MCP-Tool-Antworten in kompakte, lesbare Markdown-Tabellen | [SKILL.md](https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md) |
| **swag2mcp-cli** | Vollständige CLI-Referenz — der LLM kennt jeden Befehl, jedes Flag und jede Konfigurationsoption | [SKILL.md](https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md) |

## Warum Skills wichtig sind

Ohne Formatierungs-Skill entscheidet der LLM selbst, wie Tool-Ergebnisse angezeigt werden — oft ausführlich und inkonsistent. Der Formatierungs-Skill sorgt für einen einheitlichen, sauberen Stil: kompakte Tabellen für Listen, Inline-Überschriften für Details und kompakte Schemata.

Der CLI-Skill ermöglicht es dem LLM, "Wie mache ich..."-Fragen zu swag2mcp-Befehlen genau zu beantworten, ohne zu raten.

## Installation via LLM-Agent

Kopieren Sie diese Anfrage in Ihre KI-gestützte IDE (OpenCode, Cursor, Claude Desktop, VS Code usw.):

```
Füge die swag2mcp-Skills zu meinem Projekt hinzu:

1. Erstelle das Verzeichnis .agents/skills/swag2mcp-format/ und füge den Skill von https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md hinzu
2. Erstelle das Verzeichnis .agents/skills/swag2mcp-cli/ und füge den Skill von https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md hinzu
```

Der Agent lädt beide Skill-Dateien herunter und platziert sie in den richtigen Verzeichnissen.

## Manuelle Installation

Falls Ihr LLM-Client keine agentenbasierte Einrichtung unterstützt, laden Sie die Dateien manuell herunter:

```bash
mkdir -p .agents/skills/swag2mcp-format
mkdir -p .agents/skills/swag2mcp-cli

curl -o .agents/skills/swag2mcp-format/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md
curl -o .agents/skills/swag2mcp-cli/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md
```

## LLM-Client konfigurieren

Jeder LLM-Client und jede IDE hat seine eigene Methode zur Installation von Skills. Das folgende Beispiel ist für **OpenCode** — lesen Sie in der Dokumentation Ihres Clients nach, wie Skills korrekt installiert werden.

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

## Neustart erforderlich

**Starten Sie nach dem Hinzufügen von Skills Ihren LLM-Client oder Ihre IDE neu.** Einige Tools laden Skills nur beim Start. Wenn die Skills nicht wirken, versuchen Sie:

- **OpenCode**: Starten Sie die Anwendung neu oder führen Sie den opencode-Befehl erneut aus
- **Cursor**: Schließen und öffnen Sie das Fenster (`Cmd+Shift+W` / `Ctrl+Shift+W`)
- **Claude Desktop**: Beenden und neu starten
- **VS Code**: Fenster neu laden (`Ctrl+Shift+P` → "Developer: Reload Window")
