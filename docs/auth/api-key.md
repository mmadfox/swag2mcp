# API Key

## Для чего

Аутентификация по API-ключу. Ключ может передаваться в заголовке запроса или как параметр URL.

## Когда использовать

- Сервисы с API-ключами
- Погодные сервисы, геоданные, переводчики
- Когда API требует ключ в заголовке (`X-API-Key`) или в параметре запроса (`?api_key=...`)

## Как настроить

### Ключ в заголовке

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: api-key
      config:
        key: "X-API-Key"
        in: header
        value: "$(API_KEY)"
```

### Ключ в параметре запроса

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: api-key
      config:
        key: "api_key"
        in: query
        value: "$(API_KEY)"
```

## Параметры

| Параметр | Обязательный | Описание |
|-----------|-------------|----------|
| `key` | Да | Имя заголовка или параметра запроса |
| `in` | Да | Куда поместить ключ: `header` или `query` |
| `value` | Да | Значение ключа |

## Важные моменты

- В режиме `header` ключ добавляется как HTTP-заголовок
- В режиме `query` ключ добавляется как параметр URL
- Значение можно хранить в переменной окружения: `value: "$(MY_API_KEY)"`
