# CLI-рабочий процесс

На этой странице показаны реальные примеры использования swag2mcp из терминала — от инициализации до повседневных операций.

## Быстрый старт

```bash
# 1. Инициализация рабочей области
mkdir -p .swag2mcp && swag2mcp init ./.swag2mcp

# 2. Список спецификаций
swag2mcp ls
```

## Добавление спецификации через YAML

### Простая спецификация (публичное API)

```bash
swag2mcp add spec --yaml - <<EOF
domain: meteo
llm_title: Open-Meteo Weather API
base_url: https://api.open-meteo.com
collections:
  - llm_title: Weather Forecast
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
EOF
```

### Спецификация с аутентификацией (bearer-токен из env)

```bash
swag2mcp add spec --yaml - <<EOF
domain: my-api
llm_title: My Protected API
base_url: https://api.example.com/v1
auth:
  type: bearer
  config:
    token: \$(MY_TOKEN)
collections:
  - llm_title: Users
    location: https://raw.githubusercontent.com/my-org/my-api/main/users.yaml
EOF
```

### Спецификация с несколькими коллекциями

```bash
swag2mcp add spec --yaml - <<EOF
domain: meteo
llm_title: Open-Meteo APIs
base_url: https://api.open-meteo.com
collections:
  - llm_title: Forecast
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
  - llm_title: Air Quality
    base_url: https://air-quality-api.open-meteo.com
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
EOF
```

## Добавление коллекции к существующей спецификации

```bash
swag2mcp add collection --yaml - <<EOF
spec_domain: meteo
llm_title: Marine Weather
location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
EOF
```

## Список спецификаций

```bash
$ swag2mcp ls
Specifications:
  dadjoke (https://icanhazdadjoke.com)
    jokes (3 endpoints)
  meteo (https://api.open-meteo.com)
    forecast (5 endpoints)
    air-quality (8 endpoints)
    marine (4 endpoints)
```

### Фильтр по тегам

```bash
swag2mcp ls --tags=public
```

## Просмотр информации о рантайме

```bash
$ swag2mcp info
{
  "version": "v1.2.0",
  "workspace": "/home/user/.swag2mcp",
  "specs": {
    "total": 2,
    "active": 2,
    "disabled": 0,
    "collections": 4,
    "endpoints": 20
  },
  "http_client": {
    "timeout": "30s",
    "max_response_size": "2 KB",
    "follow_redirects": true
  },
  "mcp": {
    "transport": "stdio"
  },
  "auth": {
    "methods": ["bearer"]
  }
}
```

## Валидация конфигурации

```bash
$ swag2mcp validate
✅ Configuration is valid.
✓ Spec dadjoke: OK
✓ Spec meteo: OK
```

## Запуск MCP-сервера

### stdio (для интеграции с IDE)

```bash
swag2mcp mcp
```

### HTTP (для удалённого доступа)

```bash
swag2mcp mcp --transport sse --http-addr :8080
```

### С фильтром по тегам

```bash
swag2mcp mcp --tags=public
```

## Обновление спецификаций

Обновление всех кэшированных файлов спецификаций:

```bash
swag2mcp update
```

## Очистка кэша

```bash
swag2mcp clean
```

## Экспорт и импорт

### Резервное копирование рабочей области

```bash
swag2mcp export --output ~/backups/swag2mcp-2026-07-24.zip
```

### Восстановление на другой машине

```bash
# На новой машине
swag2mcp import --from-zip swag2mcp-2026-07-24.zip
```

## Интерактивный TUI-обозреватель

```bash
swag2mcp run
```

Открывает полноэкранный терминальный интерфейс для поиска, просмотра и вызова API.

## Мок-сервер

```bash
# Установка мок-бинарника
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest

# Запуск мок-серверов
swag2mcp-mock mockserver
```
