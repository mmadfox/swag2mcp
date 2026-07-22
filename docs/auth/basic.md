# Basic Auth

## Для чего

HTTP Basic Authentication — самый простой способ аутентификации по логину и паролю.

## Когда использовать

- Legacy API, которые поддерживают только Basic Auth
- Простая аутентификация без сложных токенов
- Внутренние сервисы

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
      type: basic
      config:
        username: "admin"
        password: "$(PASSWORD)"
```

## Параметры

| Параметр | Обязательный | Описание |
|-----------|-------------|----------|
| `username` | Да | Имя пользователя |
| `password` | Да | Пароль |

## Важные моменты

- Пароль передаётся в заголовке `Authorization: Basic ...` в кодировке Base64 — это **не шифрование**. Всегда используйте HTTPS.
- Пароль можно хранить в переменной окружения: `password: "$(MY_PASSWORD)"`
