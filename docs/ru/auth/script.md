# Script Auth

## Назначение

Аутентификация через внешний скрипт — самый гибкий метод. Вы можете написать скрипт на любом языке (bash, Python и т.д.), который получает токен любым способом и возвращает его swag2mcp.

## Когда использовать

- Пользовательские или нестандартные схемы аутентификации
- Сложная логика получения токена (многошаговая, с дополнительными проверками)
- Когда ни один из стандартных методов не подходит

## Конфигурация

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

| Параметр | Обязательно | Описание |
|-----------|----------|-------------|
| `domain` | Да | Имя файла скрипта (без расширения) |

## Расположение скрипта

Скрипт должен быть помещён в директорию `auth_scripts` вашей рабочей области:

- **Linux / macOS:** `{workspace}/auth_scripts/{domain}.sh`
- **Windows:** `{workspace}/auth_scripts/{domain}.bat`

## Формат вывода скрипта

Скрипт должен вывести JSON в stdout с токеном и временем его действия:

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

| Поле | Обязательно | Описание |
|-------|----------|-------------|
| `token` | Да | Токен аутентификации |
| `expires_in` | Нет | Время жизни токена в секундах (по умолчанию: 3600) |

## Примечания

- swag2mcp запускает скрипт при каждом запросе, если кэшированный токен истёк
- Скрипт должен завершиться в течение 30 секунд
- Токен кэшируется до истечения срока действия
- Имя файла скрипта = `{domain}.sh` (Unix) или `{domain}.bat` (Windows)
- `domain` не должен содержать `/` или `\`
