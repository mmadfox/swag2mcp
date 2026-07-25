# Аутентификация

## Обзор

swag2mcp поддерживает **9 методов аутентификации** для работы с API, требующими авторизации. Вы настраиваете это один раз в файле конфигурации — после этого каждый вызов API через `invoke` автоматически включает правильные токены и заголовки.

### Где настраивать

Аутентификация задаётся на уровне **спецификации** в `swag2mcp.yaml`:

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: bearer
      config:
        token: "my-token"
```

### Как это работает

- Вы указываете тип аутентификации и параметры в конфиге
- swag2mcp автоматически применяет их к каждому запросу при вызове `invoke`
- Вам **не нужно** запрашивать токен перед вызовом API — это происходит автоматически
- Если срок действия токена истекает (OAuth2, Script), swag2mcp обновляет его самостоятельно

### Переменные окружения

Конфиденциальные данные (токены, пароли, ключи) можно хранить в переменных окружения, используя синтаксис `$(VAR_NAME)`:

```yaml
auth:
  type: bearer
  config:
    token: "$(MY_API_TOKEN)"
```

swag2mcp подставляет значение `MY_API_TOKEN` при запуске.

### MCP-инструмент auth

LLM-агент может получить токен или заголовки через MCP-инструмент `auth` — например, чтобы сформировать команду curl или показать пользователю.

В **продакшене** этот инструмент следует отключать с помощью `--disable-llm-auth` (включён по умолчанию), чтобы LLM никогда не имела доступа к токенам.

### Методы

| Метод | Описание | Для чего подходит |
|--------|-------------|----------|
| [`none`](/auth/none) | Без аутентификации | Публичные API |
| [`basic`](/auth/basic) | HTTP Basic (username + password) | Устаревшие API, простая аутентификация |
| [`bearer`](/auth/bearer) | Bearer Token (JWT, токен) | Современные REST API |
| [`api-key`](/auth/api-key) | API-ключ в заголовке или параметре запроса | Сервисы с API-ключами |
| [`digest`](/auth/digest) | HTTP Digest (username + password) | Устаревшие API, безопаснее Basic |
| [`hmac`](/auth/hmac) | Подпись HMAC-SHA256 (стиль Binance) | Криптовалютные биржи |
| [`oauth2-cc`](/auth/oauth2-cc) | OAuth2 Client Credentials | Сервер-сервер, микросервисы |
| [`oauth2-pwd`](/auth/oauth2-pwd) | OAuth2 Password Grant | Приложения с входом пользователя |
| [`script`](/auth/script) | Внешний скрипт для получения токена | Любая пользовательская схема auth |
