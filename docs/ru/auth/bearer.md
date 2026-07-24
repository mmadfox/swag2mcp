# Bearer Auth

## Назначение

Bearer Token аутентификация — самый распространённый метод для современных REST API. Токен отправляется в заголовке `Authorization: Bearer <token>`.

## Когда использовать

- Современные REST API
- JWT (JSON Web Tokens)
- Токены доступа OAuth2 (когда токен уже получен)
- Любой API, принимающий Bearer Token

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
      type: bearer
      config:
        token: "eyJhbGciOiJIUzI1NiIs..."
```

## Параметры

| Параметр | Обязательно | Описание |
|-----------|----------|-------------|
| `token` | Да | Bearer-токен (JWT, OAuth2-токен и т.д.) |

## Примечания

- Токен статичен — если он истекает, вам нужно обновить его в конфиге вручную
- Для автоматического обновления токена используйте `oauth2-cc` или `oauth2-pwd`
- Храните токен в переменной окружения: `token: "$(API_TOKEN)"`
