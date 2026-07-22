# OAuth2 Password Grant

## Для чего

OAuth2 Resource Owner Password Grant — аутентификация по логину и паролю пользователя. Подходит для first-party приложений, где пользователь доверяет свои учётные данные приложению.

## Когда использовать

- First-party приложения (мобильные, веб)
- Интеграция с Keycloak и аналогичными Identity Provider
- Когда API поддерживает OAuth2 Password Grant

## Как настроить

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: oauth2-pwd
      config:
        client_id: "$(CLIENT_ID)"
        client_secret: "$(CLIENT_SECRET)"
        username: "$(USERNAME)"
        password: "$(PASSWORD)"
        token_url: "https://auth.example.com/oauth/token"
        scopes:
          - openid
          - profile
```

## Параметры

| Параметр | Обязательный | Описание |
|-----------|-------------|----------|
| `client_id` | Да | Идентификатор клиента |
| `username` | Да | Имя пользователя |
| `password` | Да | Пароль |
| `token_url` | Да | URL токен-эндпоинта |
| `client_secret` | Нет | Секрет клиента (опционально, для public client) |
| `scopes` | Нет | Список разрешений (опционально) |

## Важные моменты

- `client_secret` опционален — поддерживаются **public клиенты** (например, Keycloak)
- swag2mcp автоматически обновляет токен по истечении срока действия
- Токен кэшируется до expiry
- Все параметры можно хранить в переменных окружения
