# Настройки спецификации

Настройки спецификации определяют API-сервис и переопределяют глобальные настройки для этого конкретного API. Каждая спецификация представляет один логический API (например, "Open-Meteo Weather APIs") и может содержать несколько коллекций (файлов спецификаций).

## Раздел спецификации

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    llm_instruction: "Используй этот API для прогнозов погоды и климатических данных"
    base_url: https://api.open-meteo.com
    disable: false
    tags: ["weather", "climate"]
    http_client:
      timeout: 10s
      max_response_size: 1024
    auth:
      type: bearer
      config:
        token: "$(TOKEN)"
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

## Параметры

### domain

- **Тип:** `string`
- **Обязательно:** Да
- **Описание:** Уникальный идентификатор этой спецификации API. Используется внутри системы для ссылки на спецификацию.
- **Правила:** 1-60 символов. Только строчные буквы (`a-z`), цифры (`0-9`), дефисы (`-`) и подчёркивания (`_`).
- **Пример:** `meteo`, `binance`, `my-api`

### llm_title

- **Тип:** `string`
- **Обязательно:** Да
- **Описание:** Человекочитаемое имя, которое LLM использует для обращения к этому API. Отображается в ответах MCP-инструментов.
- **Правила:** 5-120 символов. Только буквы, цифры, пробелы и базовая пунктуация.
- **Пример:** `Open-Meteo Weather APIs`, `Binance Market Data`

### llm_instruction

- **Тип:** `string`
- **По умолчанию:** `""`
- **Описание:** Инструкции для LLM о том, как использовать этот API. Описывает, что делает API и когда его использовать.
- **Правила:** Макс. 500 символов. Только буквы, цифры, пробелы и базовая пунктуация.
- **Пример:** `"Используй этот API для прогнозов погоды, текущих условий и климатических данных."`

### base_url

- **Тип:** `string`
- **Обязательно:** Да
- **Описание:** Базовый URL для всех запросов к API в этой спецификации. Пути эндпоинтов из OpenAPI-спецификации добавляются к этому URL.
- **Пример:** `https://api.open-meteo.com`, `https://api.binance.com`
- **Примечание:** Может быть переопределён на уровне коллекции, если разные коллекции используют разные базовые URL.

### disable

- **Тип:** `bool`
- **По умолчанию:** `false`
- **Описание:** Если `true`, эта спецификация исключается из MCP-инструментов. Она не загружается, не индексируется и недоступна LLM.
- **Когда использовать:** Временно отключить API без удаления из конфига. Полезно для API, которые недоступны, устарели или находятся на обслуживании.

### tags

- **Тип:** `[]string` (массив строк)
- **По умолчанию:** `[]`
- **Описание:** Теги для фильтрации спецификаций. Используются с флагом `--tags` в командах CLI (`ls`, `validate`, `mcp`, `update`).
- **Пример:** `["public", "weather"]`, `["internal", "production"]`
- **Эффект:** Когда вы запускаете `swag2mcp mcp --tags=public`, загружаются только спецификации с тегом `public`.

### http_client

- **Тип:** `object`
- **По умолчанию:** наследуется от глобального
- **Описание:** Переопределить глобальные настройки HTTP-клиента для этой спецификации. Все настройки из глобального `http_client` могут быть переопределены: `timeout`, `max_response_size`, `user_agent`, `follow_redirects`, `max_redirects`, `random`, `proxy`, `headers`, `cookies`.
- **Пример:**
  ```yaml
  http_client:
    timeout: 60s
    max_response_size: 4194304
    headers:
      "X-DC": "us-east-1"
  ```

### auth

- **Тип:** `object`
- **По умолчанию:** `none` (без аутентификации)
- **Описание:** Конфигурация аутентификации для этой спецификации. Подробнее о всех 9 методах и их параметрах: [Аутентификация](/auth/overview).
- **Пример:**
  ```yaml
  auth:
    type: bearer
    config:
      token: "$(API_TOKEN)"
  ```

### collections

- **Тип:** `[]object` (массив коллекций)
- **Обязательно:** Да (как минимум 1)
- **Описание:** Список файлов спецификаций OpenAPI/Swagger/Postman, принадлежащих этой спецификации. Каждая коллекция — это один файл спецификации.
- **Правила:** 1-30 коллекций на спецификацию.
- **См.:** [Настройки коллекции](./collection-settings) — все параметры коллекций.

## Отключение спецификации

Отключённые спецификации не загружаются и не индексируются. LLM не может их видеть или использовать.

```yaml
specs:
  - domain: old-api
    llm_title: Old API
    base_url: https://old-api.example.com
    disable: true
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Переопределение HTTP-клиента

Все настройки `http_client` с глобального уровня могут быть переопределены на уровне спецификации. Значения спецификации имеют приоритет над глобальными значениями только для этой спецификации.

```yaml
specs:
  - domain: slow-api
    llm_title: Slow API
    base_url: https://slow-api.example.com
    http_client:
      timeout: 120s
      max_response_size: 8388608
      headers:
        "X-DC": "us-east-1"
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Переопределение прокси

Если этой спецификации требуется другой прокси, чем глобальный, настройте его на уровне спецификации:

```yaml
specs:
  - domain: proxied-api
    llm_title: Proxied API
    base_url: https://api.example.com
    http_client:
      proxy:
        url: http://proxy.company.com:8080
        username: $(PROXY_USER)
        password: $(PROXY_PASS)
        bypass:
          - "*.local"
          - "10.0.0.0/8"
    collections:
      - llm_title: Main
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo.json
```
