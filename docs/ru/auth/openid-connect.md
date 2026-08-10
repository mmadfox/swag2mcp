# OpenID Connect (OIDC Discovery)

## Назначение

Аутентификация на основе OpenID Connect Discovery. swag2mcp загружает документ OIDC Discovery с сервера issuer, находит эндпоинт токена и получает Bearer-токен через грант **client credentials**. Токен кэшируется до истечения срока действия.

## Когда использовать

- API, которые предоставляют документ OIDC Discovery (`.well-known/openid-configuration`)
- Интеграции сервер-сервер, где у вас есть `client_id` + `client_secret`
- Когда API использует OIDC и вы хотите автоматическое получение токена

## Конфигурация

```yaml
specs:
  - domain: my-api
    llm_title: My API
    base_url: https://api.example.com
    collections:
      - llm_title: Main
        location: https://example.com/spec.yaml
    auth:
      type: openid-connect
      config:
        issuer: "https://auth.example.com"
        client_id: "$(CLIENT_ID)"
        client_secret: "$(CLIENT_SECRET)"
        scopes:
          - openid
          - profile
```

## Параметры

| Параметр | Обязательно | Описание |
|-----------|----------|-------------|
| `issuer` | Да | URL OIDC issuer. swag2mcp добавляет `/.well-known/openid-configuration` для поиска эндпоинта токена |
| `client_id` | Да | Идентификатор клиента |
| `client_secret` | Да | Секрет клиента |
| `scopes` | Нет | Список разрешений (опционально) |

## Примечания

- swag2mcp автоматически запрашивает новый токен, когда текущий истекает
- Токен кэшируется до истечения срока действия (`expires_in`)
- Если сервер не предоставляет `expires_in`, токен считается действительным в течение 1 часа
- Поля `issuer`, `client_id` и `client_secret` поддерживают `$(VAR)` для чтения из переменных окружения
- Токен получается через грант **client credentials** (сервер-сервер, без участия пользователя)
- Документ Discovery должен содержать поле `token_endpoint`
