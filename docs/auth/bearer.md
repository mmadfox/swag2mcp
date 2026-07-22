# Bearer Auth

## Для чего

Аутентификация по Bearer Token — самый распространённый метод для современных REST API. Токен передаётся в заголовке `Authorization: Bearer <token>`.

## Когда использовать

- Современные REST API
- JWT (JSON Web Tokens)
- OAuth2 access tokens (когда токен уже получен)
- Любой API, который принимает Bearer Token

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
      type: bearer
      config:
        token: "eyJhbGciOiJIUzI1NiIs..."
```

## Параметры

| Параметр | Обязательный | Описание |
|-----------|-------------|----------|
| `token` | Да | Bearer токен (JWT, OAuth2 token и т.д.) |

## Важные моменты

- Токен статический — если он истекает, нужно обновить его в конфиге вручную
- Для автоматического обновления токенов используйте `oauth2-cc` или `oauth2-pwd`
- Токен можно хранить в переменной окружения: `token: "$(API_TOKEN)"`
