# Переменные окружения

## Обзор

swag2mcp поддерживает подстановку переменных окружения в конфигурационном файле с использованием синтаксиса `$(VAR_NAME)`. Это позволяет хранить чувствительные данные (токены, пароли, ключи) вне YAML-файла.

## Как это работает

При запуске swag2mcp сканирует конфигурацию на предмет шаблонов `$(VAR_NAME)` и заменяет их значением соответствующей переменной окружения.

```yaml
specs:
  - domain: meteo
    auth:
      type: bearer
      config:
        token: "$(API_TOKEN)"
```

Если переменная окружения `API_TOKEN` установлена, она будет подставлена. Если не установлена, значение становится пустым.

## Где `$(VAR)` разрешается

| Поле | Пример |
|------|--------|
| Auth `token` (bearer) | `token: "$(API_TOKEN)"` |
| Auth `username` / `password` (basic, digest) | `password: "$(API_PASSWORD)"` |
| Auth `client_id` / `client_secret` (oauth2-cc, oauth2-pwd) | `client_secret: "$(OAUTH_SECRET)"` |
| Auth `api_key` / `secret_key` (hmac) | `api_key: "$(BINANCE_API_KEY)"` |
| Auth `domain` (script) | `domain: "$(AUTH_DOMAIN)"` |
| MCP server token | `token: "$(MCP_TOKEN)"` |
| HTTP client headers | `"X-API-Key": "$(API_KEY)"` |
| HTTP client cookie values | `value: "$(SESSION_TOKEN)"` |

## Где `$(VAR)` НЕ разрешается

- Базовые URL (`base_url`)
- Расположения коллекций (`location`)
- Доменные имена спецификаций (`domain`)

## Пример

```bash
export API_TOKEN="eyJhbGciOiJIUzI1NiIs..."
export MCP_TOKEN="my-secret-token"

swag2mcp mcp
```

## Лучшие практики безопасности

- **Никогда** не храните секреты напрямую в YAML-файле
- Используйте переменные окружения или внешний менеджер секретов
- Добавьте YAML-файл в `.gitignore`, если он содержит любые жёстко закодированные секреты
- Устанавливайте переменные окружения в профиле оболочки, конфигурации IDE или пайплайне развёртывания

## Детали синтаксиса

- `$(VAR_NAME)` — стандартный синтаксис
- `$( VAR_NAME )` — пробелы внутри скобок разрешены и обрезаются
- `$()` — пустое имя переменной возвращает исходную строку без изменений
- Вложенные шаблоны `$(...)` не разрешаются
