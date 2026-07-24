# Каскад конфигурации

swag2mcp использует трёхуровневый каскад конфигурации. Каждый уровень переопределяет предыдущий. Это позволяет задавать разумные значения по умолчанию глобально и точно настраивать параметры для конкретных спецификаций или коллекций.

## Уровни

```
Global (http_client, mcp, mock_enabled, disable_ratelimiter, rate_limit_interval)
    ↓ переопределяет
Spec (specs[].http_client, specs[].auth, specs[].base_url, specs[].disable, specs[].tags)
    ↓ переопределяет
Collection (specs[].collections[].http_client, specs[].collections[].base_url, specs[].collections[].disable)
```

## Что переопределяет что

| Параметр | Global | Spec | Collection |
|-----------|--------|------|------------|
| `http_client.timeout` | ✅ | ✅ | ✅ |
| `http_client.max_response_size` | ✅ | ✅ | ✅ |
| `http_client.user_agent` | ✅ | ✅ | ✅ |
| `http_client.follow_redirects` | ✅ | ✅ | ✅ |
| `http_client.max_redirects` | ✅ | ✅ | ✅ |
| `http_client.proxy` | ✅ | ✅ | ✅ |
| `http_client.random` | ✅ | ✅ | ✅ |
| `http_client.headers` | ✅ | ✅ | ✅ |
| `http_client.cookies` | ✅ | ✅ | ✅ |
| `base_url` | ❌ | ✅ | ✅ |
| `auth` | ❌ | ✅ | ❌ |
| `disable` | ❌ | ✅ | ✅ |
| `tags` | ❌ | ✅ | ❌ |
| `mock_enabled` | ✅ | ❌ | ❌ |
| `disable_ratelimiter` | ✅ | ❌ | ❌ |
| `rate_limit_interval` | ✅ | ❌ | ❌ |

Все настройки `http_client` могут быть переопределены на каждом уровне. Настройки уровня коллекции имеют полный приоритет над спецификацией и глобальным уровнем.

## Пример каскада

```yaml
http_client:
  timeout: 30s
  max_response_size: 1048576
  headers:
    "User-Agent": "swag2mcp/1.0"

specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    http_client:
      timeout: 60s  # переопределяет глобальный timeout
      headers:
        "X-API-Version": "2"  # добавляется к глобальным заголовкам
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        http_client:
          timeout: 120s  # переопределяет timeout спецификации
          headers:
            "X-Custom": "value"  # добавляется к заголовкам спецификации + глобальным
```

## Результирующие настройки для коллекции "Forecast"

```
timeout: 120s (из коллекции, переопределяет 60s спецификации и 30s глобальные)
max_response_size: 1048576 (из глобальных)
headers:
  - User-Agent: swag2mcp/1.0 (из глобальных)
  - X-API-Version: 2 (из спецификации)
  - X-Custom: value (из коллекции)
```

## Как работает объединение

### Настройки HTTP-клиента

Простые значения (`timeout`, `max_response_size`, `user_agent`, `follow_redirects`, `max_redirects`, `random`) **заменяются** на каждом уровне. Если спецификация устанавливает `timeout: 60s`, это полностью заменяет глобальное `30s`.

### Заголовки

Заголовки **объединяются** по уровням. Заголовки всех трёх уровней комбинируются. Если один и тот же ключ заголовка появляется на нескольких уровнях, побеждает самый нижний уровень.

### Куки

Куки **объединяются** по уровням. Если одно и то же имя куки появляется на нескольких уровнях, побеждает самый нижний уровень.

### Прокси

Прокси **заменяется** на каждом уровне. Если спецификация устанавливает прокси, он полностью заменяет глобальный прокси для этой спецификации.
