# swag2mcp-cli

Der **swag2mcp-cli** Skill gibt Ihrem LLM die vollständige swag2mcp CLI-Referenz — jeden Befehl, jedes Flag, jedes Argument und jede Konfigurationsoption. Mit diesem Skill kann der LLM "Wie mache ich..."-Fragen genau beantworten, ohne zu raten.

## Was wird abgedeckt

Alle 13 CLI-Befehle:

| Befehl | Zweck |
|--------|-------|
| `init` | Arbeitsbereich und Konfiguration initialisieren |
| `add` | Spezifikation oder Sammlung hinzufügen |
| `delete` | Spezifikation oder Sammlung entfernen |
| `ls` | Konfigurierte Spezifikationen auflisten |
| `run` | API Explorer TUI starten |
| `validate` | Konfigurationsdatei validieren |
| `clean` | Zwischengespeicherte Daten löschen |
| `update` | Cache aus Konfiguration aktualisieren |
| `mcp` | MCP-Server starten |
| `version` | Versionsinformation anzeigen |
| `info` | Laufzeitinformationen anzeigen |
| `import` | Arbeitsbereich aus ZIP-Datei importieren |
| `export` | Arbeitsbereich in ZIP-Datei exportieren |

Plus alle Flags, Konfigurationsdateistruktur, Authentifizierungsmethoden und erweiterte Optionen.

## Direkter Link

<https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md>

## Installation via LLM-Agent

Kopieren Sie diese Anfrage in Ihre KI-gestützte IDE:

```
Erstelle das Verzeichnis .agents/skills/swag2mcp-cli/ und füge den Skill von
https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md hinzu
```

## Manuelle Installation

```bash
mkdir -p .agents/skills/swag2mcp-cli
curl -o .agents/skills/swag2mcp-cli/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md
```

## Neustart erforderlich

Starten Sie nach dem Hinzufügen des Skills Ihren LLM-Client oder Ihre IDE neu (siehe [Übersicht](overview.md#neustart-erforderlich)).
