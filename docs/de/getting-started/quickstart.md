# Schnellstart

Starten Sie swag2mcp in 2 Minuten.

## 1. Initialisieren

### Home-Verzeichnis (empfohlen)

Einmalige Einrichtung für Ihr gesamtes System. Die Konfiguration wird in Ihrem Home-Ordner gespeichert.

::: code-group

```bash [macOS / Linux]
swag2mcp init
# Erstellt ~/.swag2mcp/swag2mcp.yaml
```

```powershell [Windows]
swag2mcp.exe init
# Erstellt %USERPROFILE%\.swag2mcp\swag2mcp.yaml
```

:::

### Projektverzeichnis

Für einen isolierten Arbeitsbereich innerhalb Ihres Projekts.

::: code-group

```bash [macOS / Linux]
mkdir -p ./swag2mcp && swag2mcp init ./swag2mcp
```

```powershell [Windows]
mkdir ./swag2mcp; swag2mcp.exe init ./swag2mcp
```

:::

### Aus ZIP

Wenn Sie einen fertigen Arbeitsbereich haben (z. B. von einem Kollegen):

```bash
swag2mcp import --from-zip workspace.zip
```

## 2. Agenten-Skills installieren (empfohlen)

Installieren Sie die swag2mcp-Skills, um Ihrem KI-Agenten alle Befehle, Flags, das Konfigurationsformat und reale Beispiele beizubringen.

Bitten Sie Ihren Agenten:

```bash
"Fügen Sie den swag2mcp-cli-Skill von https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-cli/SKILL.md hinzu"
"Fügen Sie den swag2mcp-format-Skill von https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md hinzu"
```

> Einige IDEs erfordern einen Neustart nach dem Hinzufügen von Skills.

## 3. LLM-Client / IDE-Konfiguration

Konfigurieren Sie Ihre IDE, um eine Verbindung zu swag2mcp herzustellen. Die IDE startet den MCP-Server automatisch, wenn er benötigt wird.

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

Für andere IDEs (Cursor, VS Code, JetBrains) siehe den [Integrationsleitfaden](../integration/opencode.md).

> Wenn Sie den Arbeitsbereich unter einem benutzerdefinierten Pfad initialisiert haben (z. B. `./swag2mcp`), verwenden Sie den vollständigen Pfad im Befehl:
> `"command": ["swag2mcp", "mcp", "/absoluter/pfad/zu/swag2mcp"]`

> **Starten Sie den MCP-Server nach jeder Konfigurationsänderung neu**, damit die Änderungen wirksam werden.

## 4. MCP-Server starten

### stdio (Standard) — für lokale IDE

Nichts zu konfigurieren. Ihre IDE startet swag2mcp automatisch über die obige Konfiguration.

```bash
swag2mcp mcp
```

### SSE / Streamable HTTP — für Remote-Zugriff

```bash
swag2mcp mcp --transport sse --http-addr :8080
```

Oder in `swag2mcp.yaml` konfigurieren:

```yaml
mcp:
  transport: sse
  addr: ":8080"
  path: "/mcp"
```

Siehe [MCP-Server-Referenz](../configuration/mcp-server.md) für alle Flags.

### Specs nach Tags filtern

```bash
swag2mcp mcp --tags weather,public
```

Nur Specs mit passenden Tags werden dem LLM zur Verfügung gestellt.

### Überprüfen, ob es funktioniert

Fragen Sie nach dem Verbinden Ihren LLM-Agenten:

```bash
"Welche MCP-Tools unterstützt du?"
```

Wenn der Agent swag2mcp-Tools auflistet (`spec_list`, `search`, `invoke` usw.) — funktioniert alles.

### Beispiel-Abfragen zum Ausprobieren

| Fragen Sie Ihren Agenten | Was passiert |
|-------------------------|--------------|
| "Wie ist das Wetter in New York?" | `invoke` — ruft Open-Meteo-Vorhersage-API auf |
| "Wie ist der aktuelle BTC-Kurs?" | `invoke` — ruft Binance-Ticker-API auf |
| "Erzähl mir einen Dad Joke" | `invoke` — ruft icanhazdadjoke-API auf |
| "Zeig mir Pikachu" | `invoke` — ruft PokéAPI nach Namen auf |
| "Wer ist Rick Sanchez?" | `invoke` — ruft Rick and Morty-Charakter-API auf |
| "Wie ist die Luftqualität in Peking?" | `invoke` — ruft Open-Meteo-Luftqualitäts-API auf |
| "Wie hoch sind die Wellen vor Portugal?" | `invoke` — ruft Open-Meteo-Meer-API auf |
| "Suche nach Witzen über Hunde" | `invoke` — ruft dadjoke-Such-Endpunkt auf |
| "Liste alle Pokémon auf" | `invoke` — ruft PokéAPI-Listen-Endpunkt auf |
| "Wie hoch ist der Mount Everest?" | `invoke` — ruft Open-Meteo-Höhen-API auf |

## 5. Wie geht es weiter?

- [Konzepte](../concepts/overview.md) — die Architektur verstehen
- [Konfiguration](../configuration/config-file.md) — Einstellungen anpassen
- [CLI-Befehle](../cli/overview.md) — vollständige Befehlsreferenz
