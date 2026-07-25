# Specs

Eine Spec ist ein logischer Container, der eine API-Domain oder einen API-Dienst repräsentiert (z. B. YouTube, Binance, Open-Meteo). Jede Spec hat eine eindeutige `domain`, eine `base_url`, optional `auth` und enthält eine oder mehrere Collections.

[Collections](./collections) verweisen auf OpenAPI/Swagger/Postman-Dateien — die Spec selbst ist keine Datei, sondern die Gruppierung um sie herum.

## Domain — Namensregeln

Die `domain` ist der eindeutige Identifikator einer Spec. Sie wird als Primärschlüssel im gesamten System verwendet.

| Regel | Einschränkung |
|-------|---------------|
| Zeichen | Nur `a-z`, `0-9`, `_`, `-` |
| Länge | 1–60 Zeichen |
| Eindeutigkeit | **Keine Duplikate erlaubt** — zwei aktive Specs können nicht dieselbe Domain teilen |

**Gültige Beispiele:** `meteo`, `binance`, `github-api`, `my_service`, `openai-v1`

**Ungültige Beispiele:** `Meteo` (Großbuchstaben), `my api` (Leerzeichen), `my.api` (Punkt), `a-very-long-domain-name-that-exceeds-sixty-characters` (zu lang)

## Spec-Felder

| Feld | YAML-Schlüssel | Erforderlich | Beschreibung |
|------|---------------|-------------|--------------|
| [Domain](#domain--namensregeln) | `domain` | ✅ | Eindeutiger API-Identifikator (1–60 Zeichen, `a-z0-9_-`) |
| LLM-Titel | `llm_title` | ✅ | Menschenlesbarer Name, den der LLM zur Referenzierung dieser API verwendet (5–120 Zeichen) |
| [LLM-Instruktion](#llm-instruction) | `llm_instruction` | ❌ | Kurzer Hinweis, der in den swag2mcp-System-Prompt eingefügt wird (max. 500 Zeichen) |
| Basis-URL | `base_url` | ✅ | Basis-URL für alle API-Anfragen (gültige URL) |
| [Deaktivieren](#disable) | `disable` | ❌ | Diese Spec beim Laden und Indizieren überspringen |
| [Tags](#tags) | `tags` | ❌ | Tags zum Filtern (z. B. `["public", "demo"]`) |
| [Auth](#auth) | `auth` | ❌ | Authentifizierungskonfiguration |
| [HTTP-Client](#http-client) | `http_client` | ❌ | Pro-Spec-HTTP-Einstellungen (Header, Cookies) |
| [Collections](./collections) | `collections` | ✅ | Liste von 1–30 Collections |

## Validierung

Wenn swag2mcp die Konfiguration validiert, werden diese Regeln für jede Spec geprüft:

| Prüfung | Regel |
|---------|-------|
| **Doppelte Domains** | Keine zwei aktiven Specs dürfen dieselbe `domain` teilen |
| **Domain-Format** | Muss `^[a-z0-9_-]{1,60}$` entsprechen |
| **LLM-Titel** | Erforderlich, 5–120 Zeichen, Buchstaben/Ziffern/Leerzeichen/einfache Satzzeichen |
| **LLM-Instruktion** | Max. 500 Zeichen, gleicher Zeichensatz wie Titel |
| **Basis-URL** | Erforderlich, muss eine gültige URL sein |
| **Collections** | Erforderlich, 1–30 Elemente |
| **Auth** | Pro Auth-Typ validiert (z. B. bearer erfordert `token`, basic erfordert `username` + `password`) |
| **Speicherort** | Jede Collection muss eine gültige URL oder einen gültigen Dateipfad haben (5–250 Zeichen) |

Die Validierung läuft bei jedem `swag2mcp mcp`-Start. Wenn sie fehlschlägt, startet der MCP-Server nicht — in manchen IDEs bedeutet dies, dass der Server einfach keine Verbindung herstellt, und der LLM erhält eine klare Fehlermeldung, die erklärt, was zu beheben ist.

Verwenden Sie den Befehl [`validate`](../cli/validate.md), um Probleme vor dem Start des Servers zu diagnostizieren:

```bash
# Standard-Arbeitsbereich validieren (~/.swag2mcp)
swag2mcp validate

# Benutzerdefinierten Projektarbeitsbereich validieren
swag2mcp validate ./my-project
```

## LLM-Instruktion

Es wird empfohlen, für jede Spec `llm_instruction` zu setzen — einen kurzen Hinweis (bis zu 500 Zeichen), der dem LLM sagt, wofür diese API ist und wann sie verwendet werden soll. Diese Instruktion wird in den swag2mcp-System-Prompt eingefügt und hilft dem LLM, den Zweck der Spec ohne zusätzlichen Kontext zu verstehen.

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    llm_instruction: "Verwenden Sie diese API, um zufällige Dad Jokes zu erhalten oder nach bestimmten Witzen zu suchen."
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

Collections können auch eine eigene `llm_instruction` (bis zu 360 Zeichen) für spezifischere Anleitungen haben.

## Auth

Die Authentifizierung wird auf Spec-Ebene konfiguriert und gilt für alle ihre Collections. swag2mcp unterstützt 9 Authentifizierungsmethoden:

| Methode | YAML-Typ | Schlüsselfelder |
|---------|-----------|-----------------|
| [Keine](../auth/none.md) | `none` | — |
| [Basic](../auth/basic.md) | `basic` | `username`, `password` |
| [Bearer](../auth/bearer.md) | `bearer` | `token` |
| [Digest](../auth/digest.md) | `digest` | `username`, `password` |
| [OAuth2 Client Credentials](../auth/oauth2-cc.md) | `oauth2-cc` | `client_id`, `client_secret`, `token_url` |
| [OAuth2 Password](../auth/oauth2-pwd.md) | `oauth2-pwd` | `username`, `password`, `client_id`, `token_url` |
| [API-Schlüssel](../auth/api-key.md) | `api-key` | `key`, `value`, `in` (`header` oder `query`) |
| [HMAC](../auth/hmac.md) | `hmac` | `api_key`, `secret_key` |
| [Skript](../auth/script.md) | `script` | `domain` |

Siehe [Auth-Übersicht](../auth/overview.md) für vollständige Details zu jeder Methode.

## HTTP-Client

Sie können HTTP-Einstellungen auf Spec-Ebene überschreiben. Diese gelten für alle Anfragen, die von den Collections dieser Spec gestellt werden.

```yaml
specs:
  - domain: slow-api
    llm_title: Slow API
    base_url: https://slow-api.example.com
    http_client:
      headers:
        X-API-Version: "2"
      cookies:
        - name: session
          value: abc123
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

Einstellungen kaskadieren: global → spec → collection. Siehe [Konfigurationskaskade](../configuration/cascade.md) für Details.

## Tags

Tags ermöglichen es Ihnen, Specs nach Kategorie zu filtern. Verwenden Sie sie mit dem Flag `--tags` bei `swag2mcp ls` oder während des Bootstraps.

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    tags: ["weather", "public"]
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

```bash
# Nur Specs mit Tag "weather" auflisten
swag2mcp ls --tags weather
```

## Deaktivieren

Setzen Sie `disable: true`, um eine Spec vollständig zu überspringen. Sie wird nicht geladen, indiziert oder dem LLM zur Verfügung gestellt.

```yaml
specs:
  - domain: old-api
    llm_title: Old API (Deprecated)
    base_url: https://old-api.example.com
    disable: true
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Beispiele

### Minimale Spec

```yaml
specs:
  - domain: dadjokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

### Spec mit Auth

```yaml
specs:
  - domain: binance
    llm_title: Binance Market Data API
    base_url: https://api.binance.com
    auth:
      type: hmac
      config:
        api_key: $(BINANCE_API_KEY)
        secret_key: $(BINANCE_SECRET_KEY)
    collections:
      - llm_title: Market Data
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/binance.yaml
```

### Spec mit mehreren Collections

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
      - llm_title: Marine
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
```

### Spec mit LLM-Instruktion und Tags

```yaml
specs:
  - domain: rickandmorty
    llm_title: Rick and Morty API
    llm_instruction: "Verwenden Sie diese API, um Informationen über Charaktere, Episoden und Orte aus der Rick and Morty-Show zu erhalten."
    base_url: https://rickandmortyapi.com/api
    tags: ["entertainment", "public"]
    collections:
      - llm_title: Characters
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/rick-and-morty.json
```

## Verwandte Themen

- [Spec-Einstellungen (Konfiguration)](../configuration/spec-settings.md) — vollständige YAML-Referenz
- [Konfigurationskaskade](../configuration/cascade.md) — wie Einstellungen einander überschreiben
- [Auth-Übersicht](../auth/overview.md) — alle 9 Authentifizierungsmethoden
- [HTTP-Client](../configuration/http-client.md) — HTTP-Client-Konfiguration
- [Collections](./collections) — Spezifikationsdateien innerhalb einer Spec
