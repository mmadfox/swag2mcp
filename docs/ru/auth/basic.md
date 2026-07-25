# Basic Auth

## Назначение

HTTP Basic Authentication — самый простой способ аутентификации с использованием имени пользователя и пароля.

## Когда использовать

- Устаревшие API, поддерживающие только Basic Auth
- Простая аутентификация без сложных токенов
- Внутренние сервисы

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
      type: basic
      config:
        username: "admin"
        password: "$(PASSWORD)"
```

## Параметры

| Параметр | Обязательно | Описание |
|-----------|----------|-------------|
| `username` | Да | Имя пользователя |
| `password` | Да | Пароль |

## Примечания

- Пароль отправляется в заголовке `Authorization: Basic ...` в кодировке Base64 — это **не шифрование**. Всегда используйте HTTPS.
- Храните пароль в переменной окружения: `password: "$(MY_PASSWORD)"`
