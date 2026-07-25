# HMAC Auth

## Назначение

Подпись запросов HMAC-SHA256 — метод аутентификации, используемый криптовалютными биржами (Binance, Bybit и другие). Каждый запрос подписывается секретным ключом.

## Когда использовать

- API Binance и совместимые с Binance биржи
- Криптовалютные торговые платформы
- API, требующие подписи запросов

## Конфигурация

```yaml
specs:
  - domain: binance
    llm_title: Binance Market Data
    base_url: https://api.binance.com
    collections:
      - llm_title: Market Data
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/binance.yaml
    auth:
      type: hmac
      config:
        api_key: "$(BINANCE_API_KEY)"
        secret_key: "$(BINANCE_SECRET_KEY)"
```

## Параметры

| Параметр | Обязательно | Описание |
|-----------|----------|-------------|
| `api_key` | Да | Публичный API-ключ |
| `secret_key` | Да | Секретный ключ для подписи |

## Примечания

- swag2mcp автоматически добавляет временную метку (Unix в миллисекундах) к каждому запросу
- Подпись вычисляется из всех параметров запроса
- Храните ключи в переменных окружения: `api_key: "$(BINANCE_API_KEY)"`
- Этот метод совместим с API Binance и аналогичными биржами
