# API Key

## Назначение

Аутентификация через API-ключ. Ключ может быть отправлен как HTTP-заголовок или как параметр URL-запроса.

## Когда использовать

- Сервисы, использующие API-ключи
- Погодные сервисы, геоданные, API перевода
- Когда API ожидает ключ в заголовке (`X-API-Key`) или параметре запроса (`?api_key=...`)

## Конфигурация

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

| Параметр | Обязательно | Описание |
|-----------|----------|-------------|
| `key` | Да | Имя заголовка или параметра запроса |
| `in` | Да | Куда поместить ключ: `header` или `query` |
| `value` | Да | Значение ключа |

## Примечания

- В режиме `header` ключ добавляется как HTTP-заголовок
- В режиме `query` ключ добавляется как параметр URL
- Храните значение в переменной окружения: `value: "$(MY_API_KEY)"`
