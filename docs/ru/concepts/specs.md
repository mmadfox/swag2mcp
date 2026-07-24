# Спецификации

Спецификация — это логический контейнер, представляющий домен или сервис API (например, YouTube, Binance, Open-Meteo). Каждая спецификация имеет уникальный `domain`, `base_url`, опциональную `auth` и содержит одну или несколько коллекций.

[Коллекции](./collections) указывают на файлы OpenAPI/Swagger/Postman — сама спецификация — это не файл, а группировка вокруг них.

## Domain — правила именования

`domain` — это уникальный идентификатор спецификации. Он используется в качестве первичного ключа во всей системе.

| Правило | Ограничение |
|------|------------|
| Символы | Только `a-z`, `0-9`, `_`, `-` |
| Длина | 1–60 символов |
| Уникальность | **Дубликаты запрещены** — две активные спецификации не могут иметь одинаковый domain |

**Примеры:** `meteo`, `binance`, `github-api`, `my_service`, `openai-v1`

**Неправильные примеры:** `Meteo` (заглавные), `my api` (пробел), `my.api` (точка), `a-very-long-domain-name-that-exceeds-sixty-characters` (слишком длинный)

## Поля спецификации

| Поле | YAML-ключ | Обязательно | Описание |
|-------|----------|-------------|----------|
| [Domain](#domain--правила-именования) | `domain` | ✅ | Уникальный идентификатор API (1–60 символов, `a-z0-9_-`) |
| LLM Title | `llm_title` | ✅ | Человекочитаемое имя, которое LLM использует для обращения к этому API (5–120 символов) |
| [LLM Instruction](#llm-instruction) | `llm_instruction` | ❌ | Короткая подсказка, внедряемая в системный промпт swag2mcp (макс. 500 символов) |
| Base URL | `base_url` | ✅ | Базовый URL для всех запросов к API (валидный URL) |
| [Disable](#disable) | `disable` | ❌ | Пропустить эту спецификацию при загрузке и индексации |
| [Tags](#теги) | `tags` | ❌ | Теги для фильтрации (например, `["public", "demo"]`) |
| [Auth](#auth) | `auth` | ❌ | Конфигурация аутентификации |
| [HTTP Client](#http-клиент) | `http_client` | ❌ | HTTP-настройки для спецификации (заголовки, куки) |
| [Коллекции](./collections) | `collections` | ✅ | Список из 1–30 коллекций |

## Валидация

При проверке конфига swag2mcp проверяет следующие правила для каждой спецификации:

| Проверка | Правило |
|-------|------|
| **Дублирующиеся домены** | Никакие две активные спецификации не могут иметь одинаковый `domain` |
| **Формат domain** | Должен соответствовать `^[a-z0-9_-]{1,60}$` |
| **LLM Title** | Обязательно, 5–120 символов, буквы/цифры/пробелы/базовая пунктуация |
| **LLM Instruction** | Макс. 500 символов, тот же набор символов, что и title |
| **Base URL** | Обязательно, должен быть валидным URL |
| **Коллекции** | Обязательно, от 1 до 30 элементов |
| **Auth** | Проверяется для каждого типа (например, bearer требует `token`, basic требует `username` + `password`) |
| **Location** | Каждая коллекция должна иметь валидный URL или путь к файлу (5–250 символов) |

Валидация запускается при каждом старте `swag2mcp mcp`. Если она не пройдена, MCP-сервер не запустится — в некоторых IDE это означает, что сервер просто не подключится, и LLM получит понятное сообщение об ошибке с объяснением, что исправить.

Для диагностики проблем перед запуском сервера используйте команду [`validate`](../cli/validate.md):

```bash
# Проверка рабочей области по умолчанию (~/.swag2mcp)
swag2mcp validate

# Проверка пользовательской рабочей области проекта
swag2mcp validate ./my-project
```

## LLM Instruction

Рекомендуется задавать `llm_instruction` для каждой спецификации — короткую подсказку (до 500 символов), которая сообщает LLM, для чего предназначен этот API и когда его использовать. Эта инструкция внедряется в системный промпт swag2mcp, помогая LLM понять назначение спецификации без дополнительного контекста.

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    llm_instruction: "Используй этот API для получения случайных шуток про пап или поиска конкретных шуток по ключевому слову."
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

Коллекции также могут иметь собственный `llm_instruction` (до 360 символов) для более точных указаний.

## Auth

Аутентификация настраивается на уровне спецификации и применяется ко всем её коллекциям. swag2mcp поддерживает 9 методов аутентификации:

| Метод | YAML-тип | Ключевые поля |
|--------|-----------|------------|
| [None](../auth/none.md) | `none` | — |
| [Basic](../auth/basic.md) | `basic` | `username`, `password` |
| [Bearer](../auth/bearer.md) | `bearer` | `token` |
| [Digest](../auth/digest.md) | `digest` | `username`, `password` |
| [OAuth2 Client Credentials](../auth/oauth2-cc.md) | `oauth2-cc` | `client_id`, `client_secret`, `token_url` |
| [OAuth2 Password](../auth/oauth2-pwd.md) | `oauth2-pwd` | `username`, `password`, `client_id`, `token_url` |
| [API Key](../auth/api-key.md) | `api-key` | `key`, `value`, `in` (`header` или `query`) |
| [HMAC](../auth/hmac.md) | `hmac` | `api_key`, `secret_key` |
| [Script](../auth/script.md) | `script` | `domain` |

Подробнее о каждом методе: [Обзор аутентификации](../auth/overview.md).

## HTTP-клиент

Вы можете переопределить HTTP-настройки на уровне спецификации. Они применяются ко всем запросам, выполняемым коллекциями этой спецификации.

```yaml
specs:
  - domain: slow-api
    llm_title: Slow API
    base_url: https://slow-api.example.com
    http_client:
      headers:
        X-API-Version: "2"
      cookies:
        - name: session
          value: abc123
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

Настройки каскадируются: глобальные → спецификация → коллекция. Подробнее: [Каскад конфигурации](../configuration/cascade.md).

## Теги

Теги позволяют фильтровать спецификации по категориям. Используйте их с флагом `--tags` в `swag2mcp ls` или при загрузке.

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    tags: ["weather", "public"]
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

```bash
# Показать только спецификации с тегом "weather"
swag2mcp ls --tags weather
```

## Disable

Установите `disable: true`, чтобы полностью пропустить спецификацию. Она не будет загружена, проиндексирована или доступна LLM.

```yaml
specs:
  - domain: old-api
    llm_title: Old API (Deprecated)
    base_url: https://old-api.example.com
    disable: true
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Примеры

### Минимальная спецификация

```yaml
specs:
  - domain: dadjokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

### Спецификация с аутентификацией

```yaml
specs:
  - domain: binance
    llm_title: Binance Market Data API
    base_url: https://api.binance.com
    auth:
      type: hmac
      config:
        api_key: $(BINANCE_API_KEY)
        secret_key: $(BINANCE_SECRET_KEY)
    collections:
      - llm_title: Market Data
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/binance.yaml
```

### Спецификация с несколькими коллекциями

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
      - llm_title: Marine
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
```

### Спецификация с LLM Instruction и тегами

```yaml
specs:
  - domain: rickandmorty
    llm_title: Rick and Morty API
    llm_instruction: "Используй этот API для получения информации о персонажах, эпизодах и локациях из мультсериала Рик и Морти."
    base_url: https://rickandmortyapi.com/api
    tags: ["entertainment", "public"]
    collections:
      - llm_title: Characters
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/rick-and-morty.json
```

## Связанные разделы

- [Настройки спецификации (конфиг)](../configuration/spec-settings.md) — полный YAML-справочник
- [Каскад конфигурации](../configuration/cascade.md) — как настройки переопределяют друг друга
- [Обзор аутентификации](../auth/overview.md) — все 9 методов
- [HTTP-клиент](../configuration/http-client.md) — настройка HTTP-клиента
- [Коллекции](./collections) — файлы спецификаций внутри спецификации
