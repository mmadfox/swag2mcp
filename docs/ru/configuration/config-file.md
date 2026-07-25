# Файл конфигурации

swag2mcp использует YAML-файл конфигурации. Создаётся командой `swag2mcp init`.

## Расположение

- **Linux/macOS**: `~/.swag2mcp/swag2mcp.yaml`
- **Windows**: `%USERPROFILE%\.swag2mcp\swag2mcp.yaml`

## Базовая структура

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

## Полный пример

```yaml
# ── Глобальный HTTP-клиент ──────────────────────────────────
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

# ── MCP-сервер ──────────────────────────────────────────
mcp:
  transport: stdio
  addr: ":8080"
  path: "/mcp"
  auth:
    token: ""

# ── Mock-сервер ─────────────────────────────────────────
mock_enabled: false
mock_auth:
  oauth2_port: 9090
  digest_port: 9091
  hmac_port: 9092

# ── Ограничитель запросов ────────────────────────────────────────
disable_ratelimiter: false
rate_limit_interval: 10s

# ── Спецификации ───────────────────────────────────────────────
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    llm_instruction: "Используй этот API для прогнозов погоды и климатических данных"
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

## Переменные окружения

Используйте синтаксис `$(VAR_NAME)` для ссылки на переменные окружения. swag2mcp разрешает их при запуске.

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

`$(VAR)` разрешается в:
- Полях конфига auth: `token`, `username`, `password`, `client_id`, `client_secret`, `api_key`, `secret_key`, `domain`
- Токене аутентификации MCP-сервера: `mcp.auth.token`
- Заголовках и значениях кук HTTP-клиента

`$(VAR)` **не** разрешается в базовых URL и location коллекций.

## Валидация

```bash
# Проверка рабочей области по умолчанию (~/.swag2mcp)
swag2mcp validate

# Проверка пользовательской рабочей области проекта
swag2mcp validate ./my-project
```

Если рабочая область находится не в домашней директории (например, внутри репозитория проекта), всегда указывайте путь при запуске `validate`, `update`, `mcp` или любой другой команды. В противном случае swag2mcp будет использовать рабочую область по умолчанию `~/.swag2mcp`.
