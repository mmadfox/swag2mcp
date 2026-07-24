# OAuth2 Client Credentials

## Назначение

OAuth2 Client Credentials Grant — аутентификация для взаимодействия сервер-сервер. Приложение получает токен, используя свой client_id и client_secret, без участия пользователя.

## Когда использовать

- Микросервисы и интеграции сервер-сервер
- Взаимодействие машина-машина
- Когда API использует OAuth2 и у вас есть client_id + client_secret

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

| Параметр | Обязательно | Описание |
|-----------|----------|-------------|
| `client_id` | Да | Идентификатор клиента |
| `client_secret` | Да | Секрет клиента |
| `token_url` | Да | URL эндпоинта токена |
| `scopes` | Нет | Список разрешений (опционально) |

## Примечания

- swag2mcp автоматически запрашивает новый токен, когда текущий истекает
- Токен кэшируется до истечения срока действия (`expires_in`)
- Если сервер не предоставляет `expires_in`, токен считается действительным в течение 1 часа
- Все параметры можно хранить в переменных окружения
