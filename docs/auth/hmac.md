# HMAC Auth

## Для чего

HMAC-SHA256 подпись запроса — метод аутентификации, используемый криптовалютными биржами (Binance, Bybit и другие). Каждый запрос подписывается секретным ключом.

## Когда использовать

- Binance API и Binance-совместимые биржи
- Криптовалютные торговые платформы
- API, требующие подписи каждого запроса

## Как настроить

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

| Параметр | Обязательный | Описание |
|-----------|-------------|----------|
| `api_key` | Да | Публичный API-ключ |
| `secret_key` | Да | Секретный ключ для подписи |

## Важные моменты

- swag2mcp автоматически добавляет timestamp (Unix в миллисекундах) к каждому запросу
- Подпись вычисляется от всех параметров запроса
- Ключи можно хранить в переменных окружения: `api_key: "$(BINANCE_API_KEY)"`
- Этот метод совместим с Binance API и аналогичными биржами
