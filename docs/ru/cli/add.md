# add

## Назначение

Добавить новую **спецификацию** (API-сервис) или **коллекцию** (файл OpenAPI/Swagger/Postman) в существующую конфигурацию. Это основной способ расширения рабочей области новыми API.

## Когда использовать

- У вас есть новый API для подключения к вашему LLM-агенту
- Вы нашли URL OpenAPI-спецификации и хотите её добавить
- Вы хотите добавить дополнительный файл спецификации (коллекцию) к существующей спецификации
- Вы предпочитаете писать YAML напрямую вместо использования интерактивного мастера

## Синтаксис

```bash
swag2mcp add spec [path] [flags]
swag2mcp add collection [path] [flags]
```

## Аргументы

| Аргумент | Позиция | Обязательно | Описание |
|----------|----------|-------------|----------|
| `path` | 1 | Нет | Директория рабочей области. Если не указан, разрешается по правилам разрешения пути. |

## Флаги

### `add spec`

| Флаг | Сокращение | Тип | По умолчанию | Описание |
|------|-----------|------|---------|-------------|
| `--yaml` | `-y` | `string` | `""` | YAML-ввод (строка или `-` для stdin) |
| `--example` | `-e` | `bool` | `false` | Вывести YAML-шаблон и выйти |

### `add collection`

| Флаг | Сокращение | Тип | По умолчанию | Описание |
|------|-----------|------|---------|-------------|
| `--yaml` | `-y` | `string` | `""` | YAML-ввод (строка или `-` для stdin) |
| `--example` | `-e` | `bool` | `false` | Вывести YAML-шаблон и выйти |

## Как это работает

### Интерактивный режим (по умолчанию)

Запускает TUI-мастер, который позволяет заполнить поля спецификации или коллекции шаг за шагом.

```bash
swag2mcp add spec
swag2mcp add collection
```

### YAML inline-режим

Передайте YAML напрямую строкой. **Будьте осторожны с экранированием в оболочке** — специальные символы вроде `:`, `#`, `&`, `{` могут сломать команду.

```bash
swag2mcp add spec --yaml 'domain: meteo
llm_title: Open-Meteo API
base_url: https://meteo.swagger.io/v2
collections:
  - llm_title: Main
    location: https://example.com/spec.json'
```

### YAML из stdin (рекомендуется для сложного YAML)

Передайте через файл или используйте heredoc, чтобы полностью избежать проблем с экранированием:

```bash
# Через файл
cat spec.yaml | swag2mcp add spec --yaml -

# Heredoc
swag2mcp add spec --yaml - <<EOF
domain: my-api
llm_title: My API
llm_instruction: "Use this API for X & Y # important"
base_url: https://api.example.com/v1
collections:
  - llm_title: Main
    location: https://raw.githubusercontent.com/org/repo/main/spec.yaml
EOF
```

### YAML-шаблон

Вывести ожидаемую YAML-структуру и выйти:

```bash
swag2mcp add spec --example
swag2mcp add collection --example
```

## Формат YAML

### Спецификация

```yaml
domain: meteo
llm_title: Open-Meteo API
llm_instruction: Use this API to manage pets.
base_url: https://meteo.swagger.io/v2
tags: [public, demo]
auth:
  type: bearer
  config:
    token: $(TOKEN)
collections:
  - llm_title: Open-Meteo Swagger
    location: https://example.com/spec.json
```

### Коллекция

```yaml
spec_domain: meteo
llm_title: Orders Collection
location: https://example.com/orders.json
```

## Проверка после команды

```bash
swag2mcp ls [path]
# Новая спецификация или коллекция должна появиться в списке
```

## Нюансы

- **Автоинициализация:** Если файл конфигурации не существует, `add` автоматически запускает мастер инициализации. Вам не нужно запускать `init` отдельно.
- **Экранирование в оболочке:** Inline YAML (`--yaml '...'`) ненадёжен со специальными символами. Предпочитайте `--yaml -` с heredoc или pipe для всего, кроме простых значений.
- **`--example` завершается сразу** без проверки существующего конфига или внесения изменений.
- **`add spec` vs `add collection`:** Используйте `add spec` для нового API-сервиса (новый domain). Используйте `add collection` для добавления ещё одного файла спецификации к существующей спецификации.
