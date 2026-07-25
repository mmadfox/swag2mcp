# Настройки коллекции

Настройки коллекции определяют один файл спецификации OpenAPI/Swagger/Postman и переопределяют настройки спецификации для этого конкретного файла. Каждая коллекция принадлежит спецификации и представляет один документ спецификации API.

## Раздел коллекции

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        llm_instruction: "Используй для текущих данных и прогнозов погоды"
        disable: false
        base_url: https://forecast-api.open-meteo.com
        base_mock_url: localhost:8081
        http_client:
          timeout: 5s
```

## Параметры

### llm_title

- **Тип:** `string`
- **Обязательно:** Нет
- **Описание:** Человекочитаемое имя для этой коллекции. Отображается в ответах MCP-инструментов.
- **Правила:** Макс. 120 символов. Только буквы, цифры, пробелы и базовая пунктуация.
- **Пример:** `Forecast`, `Air Quality`, `Market Data`

### llm_instruction

- **Тип:** `string`
- **По умолчанию:** `""`
- **Описание:** Инструкции для LLM об этой конкретной коллекции. Описывает, какие эндпоинты предоставляет эта коллекция.
- **Правила:** Макс. 360 символов. Только буквы, цифры, пробелы и базовая пунктуация.
- **Пример:** `"Используй для текущих данных и прогнозов погоды."`

### title

- **Тип:** `string`
- **По умолчанию:** `""`
- **Описание:** Исходный заголовок из файла спецификации. Заполняется автоматически во время выполнения. Обычно вам не нужно устанавливать это в YAML.

### location

- **Тип:** `string`
- **Обязательно:** Да
- **Описание:** URL или локальный путь к файлу спецификации OpenAPI 3.x, Swagger 2.0 или Postman collection.
- **Правила:** 5-250 символов.
- **Примеры:**
  - URL: `https://raw.githubusercontent.com/org/repo/main/spec.yaml`
  - Локальный: `./specs/my-api.json`
  - Локальный (абсолютный): `/home/user/.swag2mcp/specs/my-api.yaml`

### disable

- **Тип:** `bool`
- **По умолчанию:** `false`
- **Описание:** Если `true`, эта коллекция исключается из MCP-инструментов. Она не загружается и не индексируется.
- **Когда использовать:** Временно отключить коллекцию без удаления из конфига. Полезно, когда файл спецификации обновляется или версия API устарела.

### http_client

- **Тип:** `object`
- **По умолчанию:** наследуется от спецификации (или глобального)
- **Описание:** Переопределить настройки HTTP-клиента для этой коллекции. Все настройки из глобального `http_client` могут быть переопределены: `timeout`, `max_response_size`, `user_agent`, `follow_redirects`, `max_redirects`, `random`, `proxy`, `headers`, `cookies`.
- **Пример:**
  ```yaml
  http_client:
    timeout: 120s
    headers:
      "X-Custom": "value"
    cookies:
      - name: "session"
        value: "abc123"
  ```

### base_url

- **Тип:** `string`
- **По умолчанию:** `""` (наследуется от спецификации)
- **Описание:** Переопределить `base_url` уровня спецификации для этой коллекции. Используйте, когда разные коллекции в рамках одной спецификации используют разные базовые URL.
- **Пример:** Если спецификация имеет `base_url: https://api.open-meteo.com`, но одна коллекция использует `https://air-quality-api.open-meteo.com`, установите `base_url` на уровне коллекции.

### base_mock_url

- **Тип:** `string`
- **По умолчанию:** `""`
- **Описание:** Адрес mock-сервера в формате `host:port`. Обязательно, когда `mock_enabled: true` в глобальном конфиге.
- **Правила:** Host должен быть `localhost`, `127.0.0.1` или `0.0.0.0`. Port должен быть валидным номером порта.
- **Пример:** `localhost:8081`, `127.0.0.1:9000`
- **Когда использовать:** У вас `mock_enabled: true`, и вы хотите тестировать эту коллекцию с фиктивными ответами.

## Несколько коллекций из одной спецификации

Спецификация может иметь несколько коллекций — например, когда у API есть отдельные файлы спецификаций для разных сервисов:

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

## Отключение коллекции

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
        disable: true
```

## Переопределение HTTP-клиента

Все настройки `http_client` могут быть переопределены на уровне коллекции. Значения коллекции имеют приоритет над значениями спецификации и глобальными значениями только для этой коллекции.

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        http_client:
          timeout: 120s
          headers:
            "X-Custom": "value"
          cookies:
            - name: "session"
              value: "abc123"
```
