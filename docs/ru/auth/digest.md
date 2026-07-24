# Digest Auth

## Назначение

HTTP Digest Access Authentication — более безопасная альтернатива Basic Auth. Пароль не передаётся в открытом виде; вместо этого используются MD5-хеши.

## Когда использовать

- Устаревшие API, поддерживающие только Digest
- Когда нужна аутентификация без отправки пароля в открытом виде
- Внутренние корпоративные системы

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
      type: digest
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

- swag2mcp сначала отправляет запрос без аутентификации, получает challenge от сервера (HTTP 401), вычисляет ответ и повторяет запрос с заголовком `Authorization: Digest ...`
- Challenge кэшируется на 5 минут — последующие запросы не требуют дополнительного round-trip
- Храните пароль в переменной окружения: `password: "$(API_PASSWORD)"`
