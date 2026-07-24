# None

## Назначение

Аутентификация не требуется. API доступен без токенов или ключей.

## Когда использовать

- Публичные API (Open-Meteo, icanhazdadjoke, PokéAPI)
- Тестовые и демонстрационные среды
- Когда API не требует авторизации

## Конфигурация

Установите `type: none` или просто опустите раздел `auth`:

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: none
```

## Параметры

Нет.

## Примечания

- Если раздел `auth` полностью отсутствует в конфиге, это эквивалентно `type: none`
- Заголовки авторизации не добавляются к запросам
