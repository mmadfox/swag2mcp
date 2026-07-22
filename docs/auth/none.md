# None

## Для чего

Аутентификация не требуется. API доступен без токенов и ключей.

## Когда использовать

- Публичные API (Open-Meteo, icanhazdadjoke, PokéAPI)
- Тестовые и демо-окружения
- Когда API не требует авторизации

## Как настроить

Укажите `type: none` или просто опустите секцию `auth`:

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

## Важные моменты

- Если секция `auth` полностью отсутствует в конфиге, это равносильно `type: none`
- Никакие заголовки авторизации не добавляются к запросам
