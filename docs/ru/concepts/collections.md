# Коллекции

Коллекция — это один файл OpenAPI/Swagger/Postman, описывающий конкретный API. Она указывает на `location` (URL или локальный путь к файлу) и принадлежит спецификации (domain).

Одна спецификация может иметь несколько коллекций — например, спецификация "meteo" может содержать коллекции "Прогноз", "Качество воздуха" и "Морские данные", каждая из которых указывает на свой файл спецификации.

## Поля коллекции

| Поле | YAML-ключ | Обязательно | Описание |
|-------|----------|-------------|----------|
| [LLM Title](#llm-instruction) | `llm_title` | ❌ | Отображаемое имя коллекции для LLM (макс. 120 символов). Автозаполняется из документа спецификации, если не задано |
| [LLM Instruction](#llm-instruction) | `llm_instruction` | ❌ | Короткая подсказка для LLM (макс. 360 символов). Автозаполняется из документа спецификации, если не задано |
| Title | `title` | ❌ | Переопределение оригинального заголовка спецификации (автозаполняется из разобранного документа) |
| [Location](#location--как-разрешаются-файлы-спецификаций) | `location` | ✅ | URL или путь к файлу спецификации (5–250 символов) |
| [Disable](#disable) | `disable` | ❌ | Пропустить эту коллекцию при загрузке |
| [HTTP Client](#http-клиент-переопределение) | `http_client` | ❌ | HTTP-настройки для коллекции (заголовки, куки) |
| [Base URL](#base-url-переопределение) | `base_url` | ❌ | Переопределить базовый URL спецификации для этой коллекции |
| [Mock-сервер](#mock-сервер) | `base_mock_url` | ❌ | Адрес mock-сервера в формате `host:port`. Обязательно, когда `mock_enabled: true` |

## Location — как разрешаются файлы спецификаций

Поле `location` сообщает swag2mcp, где найти файл OpenAPI/Swagger/Postman. Поддерживается несколько типов источников:

| Источник | Пример | Описание |
|--------|---------|-------------|
| **Удалённый URL** | `https://raw.githubusercontent.com/.../spec.yaml` | Скачивается и кэшируется |
| **Локальный файл (абсолютный)** | `/home/user/my-api.yaml` | Читается из файловой системы, кэшируется |
| **Локальный файл (относительный)** | `./my-api.yaml` | Преобразуется в абсолютный путь, кэшируется |
| **Локальный файл в workspace** | `specs/my-api.yaml` | Хранится в `~/.swag2mcp/specs/`, используется напрямую (не кэшируется) |
| **URI file://** | `file:///home/user/spec.yaml` | Преобразуется в локальный путь, кэшируется |

swag2mcp автоматически определяет тип источника:

- `https://` или `http://` → удалённый URL (кэшируется)
- `file://` → локальный файл (преобразуется в путь файловой системы)
- Всё остальное → локальный файл (с подстановкой `~` для домашней директории)

### Удалённые URL

При использовании удалённого URL swag2mcp скачивает файл и кэширует его локально. Кэш используется при последующих запусках, чтобы избежать повторных загрузок.

### Локальные файлы

Локальные файлы читаются напрямую из файловой системы. Если файл находится вне директории `specs/` рабочей области, он копируется в кэш для единообразия.

### Локальные файлы в workspace

Директория `specs/` внутри рабочей области (`~/.swag2mcp/specs/`) — рекомендуемое место для локальных файлов спецификаций. Файлы, хранящиеся здесь, используются напрямую без кэширования. Используйте относительный путь, начинающийся с `specs/`, для ссылки на них.

> **Примечание:** `specs/` — это просто имя директории (как `cache/` или `responses/`), а не понятие "спецификация". В ней хранятся фактические файлы OpenAPI/Swagger/Postman, на которые указывают коллекции.

```bash
# Импорт файла спецификации в рабочую область
swag2mcp import https://example.com/api.yaml myspec

# После импорта location становится:
# specs/myspec.yaml
```

## Система кэширования

swag2mcp кэширует удалённые файлы спецификаций, чтобы не загружать их при каждом запуске.

### Как это работает

1. Когда загружается коллекция с удалённым URL, swag2mcp проверяет кэш
2. Если существует валидная (непросроченная) запись в кэше, она используется напрямую
3. Если нет, файл скачивается, разбирается и сохраняется в кэш

### Структура кэша

```
~/.swag2mcp/
  cache/
    {sha256_hash}.spec    # Содержимое кэшированного файла спецификации
    {sha256_hash}.meta    # Метаданные кэша (JSON)
```

Каждый кэшированный файл имеет файл метаданных, содержащий:

```json
{
  "source": "https://example.com/api.yaml",
  "source_type": "url",
  "cached_at": "2024-01-01T00:00:00Z",
  "mod_time": "2024-01-01T00:00:00Z",
  "ttl_sec": 3600
}
```

### TTL кэша

Каждый кэшированный файл получает **случайный TTL** от 1 часа до 48 часов. Это предотвращает одновременное истечение срока действия всех кэшированных файлов (проблема "толпящегося стада").

### Ключ кэша

Ключ кэша — это SHA-256 хеш исходной строки location (первые 16 байт = 32 шестнадцатеричных символа).

### Управление кэшем

```bash
# Очистить кэш и ответы, перезагрузить все файлы спецификаций
swag2mcp update

# Очистить только кэш и ответы
swag2mcp clean
```

- `swag2mcp update` — проверяет конфиг, очищает `cache/` и `responses/`, затем перекэширует все location коллекций
- `swag2mcp clean` — удаляет всё содержимое `cache/` и `responses/`, а также осиротевшие скрипты аутентификации
- Старые ответы автоматически очищаются через 48 часов при запуске MCP-сервера

## Валидация

Каждая коллекция проверяется при загрузке конфига. Валидация запускается при каждом старте `swag2mcp mcp`. Если она не пройдена, MCP-сервер не запустится — в некоторых IDE это означает, что сервер просто не подключится, и LLM получит понятное сообщение об ошибке с объяснением, что исправить.

| Проверка | Правило |
|-------|------|
| **Location** | Обязательно, 5–250 символов |
| **Доступность location** | Должен быть доступным URL или существующим файлом |
| **Корректность location** | Должен быть валидным файлом OpenAPI 3.x, Swagger 2.0 или Postman |
| **LLM Title** | Макс. 120 символов, буквы/цифры/базовая пунктуация |
| **LLM Instruction** | Макс. 360 символов, тот же набор символов, что и title |
| **Base URL** | Должен быть валидным URL, если задан |
| **Base Mock URL** | Должен быть в формате `host:port` или `host:port/path`, где host — `localhost`, `127.0.0.1` или `0.0.0.0` |
| **Mock обязателен** | Если `mock_enabled: true`, каждая коллекция должна иметь `base_mock_url` |
| **Дублирующиеся mock-порты** | Никакие две коллекции не могут использовать один и тот же mock-порт |

Для диагностики проблем перед запуском сервера используйте команду [`validate`](../cli/validate.md):

```bash
# Проверка рабочей области по умолчанию (~/.swag2mcp)
swag2mcp validate

# Проверка пользовательской рабочей области проекта
swag2mcp validate ./my-project
```

## Добавление коллекций

### Через YAML-конфиг

Отредактируйте `~/.swag2mcp/swag2mcp.yaml` напрямую:

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

После редактирования перезапустите MCP-сервер (`swag2mcp mcp`), чтобы изменения вступили в силу.

### Через CLI

```bash
# Интерактивный режим
swag2mcp add collection

# Неинтерактивный с YAML
swag2mcp add collection --yaml 'spec_domain: meteo
llm_title: Forecast
location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml'

# Через stdin
cat collection.yaml | swag2mcp add collection --yaml -

# Показать YAML-пример
swag2mcp add collection --example
```

### Через импорт

```bash
# Импорт файла спецификации в рабочую область
swag2mcp import https://example.com/api.yaml
```

## LLM Instruction

Коллекции могут иметь собственный `llm_instruction` (до 360 символов) для более точных указаний. Он внедряется в системный промпт swag2mcp вместе с инструкцией уровня спецификации.

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        llm_instruction: "Используй эту коллекцию для текущей погоды и ежедневных прогнозов."
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        llm_instruction: "Используй эту коллекцию для индекса качества воздуха и данных о загрязнении."
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
```

Если `llm_title` не задан, он автоматически заполняется из поля `title` документа спецификации. Если `llm_instruction` не задан, он заполняется из поля `description` документа спецификации.

## Disable

Установите `disable: true`, чтобы пропустить коллекцию. Она не будет загружена, проиндексирована или доступна LLM.

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
        disable: true
```

## Base URL переопределение

Каждая коллекция может переопределить `base_url` спецификации. Это полезно, когда разные коллекции в рамках одной спецификации используют разные API-эндпоинты.

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
      - llm_title: Marine
        base_url: https://marine-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
```

## HTTP-клиент переопределение

Коллекции могут переопределять HTTP-настройки (заголовки, куки) с уровня спецификации и глобального уровня.

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        http_client:
          headers:
            X-API-Version: "2"
          cookies:
            - name: session
              value: abc123
```

Настройки каскадируются: глобальные → спецификация → коллекция. Подробнее: [Каскад конфигурации](../configuration/cascade.md).

## Mock-сервер

Когда `mock_enabled: true` установлено на уровне конфига, каждая коллекция должна иметь `base_mock_url`. Это сообщает swag2mcp, где запущен mock-сервер для этой коллекции.

```yaml
mock_enabled: true
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        base_mock_url: localhost:8080
```

Подробнее: [Mock-сервер](../advanced/mock-server.md).

## Примеры

### Минимальная коллекция

```yaml
specs:
  - domain: dadjokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

### Полная коллекция со всеми полями

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        llm_instruction: "Используй для текущей погоды и ежедневных прогнозов."
        title: "Custom Title"
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        disable: false
        base_url: https://forecast-api.open-meteo.com
        base_mock_url: localhost:8080
        http_client:
          headers:
            X-Custom: value
```

### Несколько коллекций в одной спецификации

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
      - llm_title: Marine
        base_url: https://marine-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
```

### Локальный файл в workspace (директория specs/)

```yaml
specs:
  - domain: myapi
    llm_title: My Internal API
    base_url: https://api.mycompany.com
    collections:
      - llm_title: Users
        location: specs/users.openapi.json
      - llm_title: Orders
        location: specs/orders.openapi.json
```

### Отключённая коллекция

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
        disable: true
```

## Связанные разделы

- [Настройки коллекции (конфиг)](../configuration/collection-settings.md) — полный YAML-справочник
- [Каскад конфигурации](../configuration/cascade.md) — как настройки переопределяют друг друга
- [Спецификации](./specs) — логические контейнеры для коллекций
- [HTTP-клиент](../configuration/http-client.md) — настройка HTTP-клиента
- [Mock-сервер](../advanced/mock-server.md) — настройка mock-сервера
- [CLI: validate](../cli/validate.md) — справочник команды validate
- [CLI: update](../cli/update.md) — справочник команды update
- [CLI: clean](../cli/clean.md) — справочник команды clean
