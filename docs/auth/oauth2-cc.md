# OAuth2 Client Credentials

## Для чего

OAuth2 Client Credentials Grant — аутентификация для server-to-server коммуникации. Приложение получает токен, используя свой client_id и client_secret, без участия пользователя.

## Когда использовать

- Микросервисы и server-to-server интеграции
- Machine-to-machine коммуникация
- Когда API использует OAuth2 и у вас есть client_id + client_secret

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
      type: oauth2-cc
      config:
        client_id: "$(CLIENT_ID)"
        client_secret: "$(CLIENT_SECRET)"
        token_url: "https://auth.example.com/oauth/token"
        scopes:
          - read
          - write
```

## Параметры

| Параметр | Обязательный | Описание |
|-----------|-------------|----------|
| `client_id` | Да | Идентификатор клиента |
| `client_secret` | Да | Секрет клиента |
| `token_url` | Да | URL токен-эндпоинта |
| `scopes` | Нет | Список разрешений (опционально) |

## Важные моменты

- swag2mcp автоматически запрашивает новый токен, когда текущий истекает
- Токен кэшируется до окончания срока действия (`expires_in`)
- Если сервер не указал `expires_in`, токен считается действительным 1 час
- Все параметры можно хранить в переменных окружения
