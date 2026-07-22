# Script Auth

## Для чего

Аутентификация через внешний скрипт — максимально гибкий метод. Вы можете написать скрипт на любом языке (bash, Python, и т.д.), который получит токен любым способом и вернёт его swag2mcp.

## Когда использовать

- Кастомные или нестандартные схемы аутентификации
- Сложная логика получения токена (многошаговая, с дополнительными проверками)
- Когда ни один из стандартных методов не подходит

## Как настроить

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: script
      config:
        domain: "my-auth"
```

## Параметры

| Параметр | Обязательный | Описание |
|-----------|-------------|----------|
| `domain` | Да | Имя файла скрипта (без расширения) |

## Где хранить скрипт

Скрипт должен находиться в директории `auth_scripts` вашего workspace:

- **Linux / macOS:** `{workspace}/auth_scripts/{domain}.sh`
- **Windows:** `{workspace}/auth_scripts/{domain}.bat`

## Формат вывода скрипта

Скрипт должен вывести в stdout JSON с токеном и сроком действия:

```bash
#!/bin/bash
# auth_scripts/my-auth.sh

TOKEN=$(curl -s -X POST https://auth.example.com/token \
  -d "grant_type=client_credentials" \
  -d "client_id=$CLIENT_ID" \
  -d "client_secret=$CLIENT_SECRET" | jq -r '.access_token')

echo "{\"token\": \"$TOKEN\", \"expires_in\": 3600}"
```

### Поля JSON

| Поле | Обязательное | Описание |
|------|-------------|----------|
| `token` | Да | Токен для аутентификации |
| `expires_in` | Нет | Срок действия в секундах (по умолчанию 3600) |

## Важные моменты

- swag2mcp запускает скрипт при каждом запросе, если кэшированный токен истёк
- Скрипт должен завершиться за 30 секунд
- Токен кэшируется до окончания срока действия
- Имя файла скрипта = `{domain}.sh` (Unix) или `{domain}.bat` (Windows)
- `domain` не должен содержать `/` или `\`
