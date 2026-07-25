# Skills

## Ausgabeformat anpassen

Jedes swag2mcp-MCP-Tool gibt strukturierte JSON-Daten zurück. Wie diese Daten dem **Benutzer präsentiert** werden, hängt vom Formatierungs-Skill des LLM ab — und Sie können es vollständig steuern.

### Der Standard-Formatierungs-Skill

swag2mcp enthält einen integrierten Formatierungs-Skill, der kompaktes, menschenlesbares Markdown für jede Tool-Antwort definiert:

[swag2mcp-format SKILL.md](https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md)

Dieser Skill deckt alle 19 MCP-Tools ab mit:
- Kompakten Tabellen für Listen (Specs, Collections, Tags, Endpunkte)
- Inline-Headern für Detailansichten
- Kompakter Schema-Darstellung für `inspect`
- Konsistentem Stil über alle Antworten hinweg

### Warum Skills wichtig sind

Dieselben Daten können je nach Skill radikal unterschiedlich dargestellt werden:

| Stil | Beispielausgabe |
|------|----------------|
| **Kompakte Tabellen** (Standard) | `GET /pet/{petId}` — Haustier nach ID finden |
| **Ausführlich** | `Methode: GET, Pfad: /pet/{petId}, Zusammenfassung: Haustier nach ID finden, Veraltet: false` |
| **Minimal** | `GET /pet/{petId}` |
| **Technisch** | `GET /pet/{petId} → 200: Pet-Objekt, 404: Nicht gefunden` |
| **Benutzerdefiniert** | Jedes Format, das Sie beschreiben können |

### Ihren eigenen Skill erstellen

Sie können Ihren eigenen Formatierungs-Skill schreiben, indem Sie das genaue Ausgabeformat beschreiben, das Sie wünschen. Der Skill ist eine Markdown-Datei mit Formatierungsregeln für jedes Tool. Hier sind einige Ideen:

- **JSON-Ausgabe** — rohes JSON für maschinelle Verarbeitung zurückgeben
- **CSV-ähnlich** — tabellarische Daten für den Tabellenkalkulationsimport
- **Diagrammfreundlich** — Mermaid- oder ASCII-Diagramme der API-Struktur
- **Minimal** — nur Methode und Pfad, nichts weiter
- **Dokumentationsstil** — vollständige Beschreibungen, Beispiele und Hinweise

### Die einzige Grenze ist das Modell

Die Qualität der formatierten Ausgabe hängt vollständig von der Fähigkeit des LLM ab, Ihren Formatierungsregeln zu folgen. Ein gut geschriebener Skill mit klaren Beispielen erzeugt konsistente, zuverlässige Ausgaben. Ein vager Skill erzeugt inkonsistente Ergebnisse.

Sie können:
- Den Standard-Skill unverändert verwenden
- Ihn forken und die Formatierung nach Ihrem Geschmack anpassen
- Ihren eigenen von Grund auf neu schreiben
- Je nach Aufgabe zwischen Skills wechseln

### Wie man einen Skill verwendet

Skills werden vom LLM-Client (OpenCode, Cursor, Claude Desktop usw.) als Teil seines System-Prompts oder seiner Agentenkonfiguration geladen. Lesen Sie in der Dokumentation Ihres Clients nach, wie eine Skill-Datei angehängt wird.

Für OpenCode werden Skills in `opencode.json` konfiguriert:

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
