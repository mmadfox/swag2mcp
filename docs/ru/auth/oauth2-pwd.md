# OAuth2 Password Grant

## Назначение

OAuth2 Resource Owner Password Grant — аутентификация с использованием имени пользователя и пароля. Подходит для собственных приложений, где пользователь доверяет приложению свои учётные данные.

## Когда использовать

- Собственные приложения (мобильные, веб)
- Интеграция с Keycloak и аналогичными Identity Provider
- Когда API поддерживает OAuth2 Password Grant

## Конфигурация

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

| Параметр | Обязательно | Описание |
|-----------|----------|-------------|
| `client_id` | Да | Идентификатор клиента |
| `username` | Да | Имя пользователя |
| `password` | Да | Пароль |
| `token_url` | Да | URL эндпоинта токена |
| `client_secret` | Нет | Секрет клиента (опционально, для публичных клиентов) |
| `scopes` | Нет | Список разрешений (опционально) |

## Примечания

- `client_secret` опционален — поддерживаются **публичные клиенты** (например, Keycloak)
- swag2mcp автоматически обновляет токен при истечении срока действия
- Токен кэшируется до истечения срока действия
- Все параметры можно хранить в переменных окружения
