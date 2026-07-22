# Digest Auth

## Для чего

HTTP Digest Access Authentication — более безопасная альтернатива Basic Auth. Пароль не передаётся в открытом виде, вместо этого используется MD5-хеш.

## Когда использовать

- Legacy API, которые поддерживают только Digest
- Когда нужна аутентификация без передачи пароля в открытом виде
- Внутренние корпоративные системы

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
      type: digest
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

- swag2mcp сначала отправляет запрос без аутентификации, получает от сервера challenge (HTTP 401), вычисляет ответ и повторяет запрос с заголовком `Authorization: Digest ...`
- Challenge кэшируется на 5 минут — повторные запросы не требуют дополнительного round-trip
- Пароль можно хранить в переменной окружения: `password: "$(API_PASSWORD)"`
