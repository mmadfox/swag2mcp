# Authentication

## Overview

swag2mcp поддерживает **9 методов аутентификации** для работы с API, которые требуют авторизации. Настройка выполняется один раз в конфигурационном файле — после этого все вызовы API через `invoke` будут автоматически содержать нужные токены и заголовки.

### Где настраивается

Аутентификация указывается на уровне **spec** в `swag2mcp.yaml`:

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
- Вам **не нужно** дополнительно запрашивать токен перед вызовом API — всё происходит автоматически
- Если токен истекает (OAuth2, Script), swag2mcp обновляет его самостоятельно

### Переменные окружения

Чувствительные данные (токены, пароли, ключи) можно хранить в переменных окружения, используя синтаксис `$(VAR_NAME)`:

```yaml
auth:
  type: bearer
  config:
    token: "$(MY_API_TOKEN)"
```

swag2mcp подставит значение переменной `MY_API_TOKEN` при запуске.

### MCP инструмент auth

LLM-агент может получить токен или заголовки через MCP инструмент `auth` — например, чтобы сформировать curl-запрос или показать пользователю.

В **production** этот инструмент рекомендуется отключать флагом `--disable-llm-auth` (включён по умолчанию), чтобы LLM не имела доступа к токенам.

### Методы

| Метод | Описание | Для каких API |
|-------|----------|---------------|
| `none` | Без аутентификации | Публичные API |
| `basic` | HTTP Basic (логин + пароль) | Legacy API, простая аутентификация |
| `bearer` | Bearer Token (JWT, токен) | Современные REST API |
| `api-key` | API-ключ в заголовке или параметре запроса | Сервисы с API-ключами |
| `digest` | HTTP Digest (логин + пароль) | Legacy API, безопаснее Basic |
| `hmac` | HMAC-SHA256 подпись (Binance-style) | Криптовалютные биржи |
| `oauth2-cc` | OAuth2 Client Credentials | Server-to-server, микросервисы |
| `oauth2-pwd` | OAuth2 Password Grant | Приложения с логином пользователя |
| `script` | Внешний скрипт для получения токена | Любая кастомная схема |
