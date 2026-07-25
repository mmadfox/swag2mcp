# Konfigurationsdatei

swag2mcp verwendet eine YAML-Konfigurationsdatei. Erstellt von `swag2mcp init`.

## Speicherort

- **Linux/macOS**: `~/.swag2mcp/swag2mcp.yaml`
- **Windows**: `%USERPROFILE%\.swag2mcp\swag2mcp.yaml`

## Grundstruktur

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

## Vollständiges Beispiel

```yaml
# ── Globaler HTTP-Client ──────────────────────────────────
http_client:
  timeout: 30s
  max_response_size: 1048576
  user_agent: "swag2mcp-global/1.0"
  follow_redirects: true
  max_redirects: 10
  random: false
  proxy:
    url: ""
    username: ""
    password: ""
    bypass: []
  headers:
    "Accept": "application/json"
  cookies:
    - name: "session"
      value: "abc123"

# ── MCP-Server ──────────────────────────────────────────
mcp:
  transport: stdio
  addr: ":8080"
  path: "/mcp"
  auth:
    token: ""

# ── Mock-Server ─────────────────────────────────────────
mock_enabled: false
mock_auth:
  oauth2_port: 9090
  digest_port: 9091
  hmac_port: 9092

# ── Ratenbegrenzer ────────────────────────────────────────
disable_ratelimiter: false
rate_limit_interval: 10s

# ── Specs ───────────────────────────────────────────────
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    llm_instruction: "Verwenden Sie diese API für Wettervorhersagen und Klimadaten"
    base_url: https://api.open-meteo.com
    disable: false
    tags: ["weather", "climate"]
    http_client:
      timeout: 10s
    auth:
      type: bearer
      config:
        token: "$(TOKEN)"
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
        disable: false
        http_client:
          timeout: 5s

  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Umgebungsvariablen

Verwenden Sie die Syntax `$(VAR_NAME)`, um auf Umgebungsvariablen zu verweisen. swag2mcp löst sie beim Start auf.

```yaml
specs:
  - domain: meteo
    auth:
      type: bearer
      config:
        token: "$(API_TOKEN)"

mcp:
  auth:
    token: "$(MCP_TOKEN)"
```

`$(VAR)` wird aufgelöst in:
- Auth-Konfigurationsfeldern: `token`, `username`, `password`, `client_id`, `client_secret`, `api_key`, `secret_key`, `domain`
- MCP-Server-Auth-Token: `mcp.auth.token`
- HTTP-Client-Headern und Cookie-Werten

`$(VAR)` wird **nicht** in Basis-URLs oder Collection-Speicherorten aufgelöst.

## Validierung

```bash
# Standard-Arbeitsbereich validieren (~/.swag2mcp)
swag2mcp validate

# Benutzerdefinierten Projektarbeitsbereich validieren
swag2mcp validate ./my-project
```

Wenn sich der Arbeitsbereich nicht im Home-Verzeichnis befindet (z. B. innerhalb eines Projekt-Repositorys), geben Sie den Pfad immer an, wenn Sie `validate`, `update`, `mcp` oder einen anderen Befehl ausführen. Andernfalls verwendet swag2mcp den Standard-Arbeitsbereich `~/.swag2mcp`.
