# Autenticación HMAC

## Propósito

Firma de solicitudes HMAC-SHA256 — el método de autenticación utilizado por exchanges de criptomonedas (Binance, Bybit y otros). Cada solicitud se firma con una clave secreta.

## Cuándo usarlo

- API de Binance y exchanges compatibles con Binance
- Plataformas de trading de criptomonedas
- APIs que requieren firma de solicitudes

## Configuración

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

## Parámetros

| Parámetro | Requerido | Descripción |
|-----------|-----------|-------------|
| `api_key` | Sí | Clave de API pública |
| `secret_key` | Sí | Clave secreta para firmar |

## Notas

- swag2mcp agrega automáticamente una marca de tiempo (Unix en milisegundos) a cada solicitud
- La firma se calcula a partir de todos los parámetros de la solicitud
- Almacene las claves en variables de entorno: `api_key: "$(BINANCE_API_KEY)"`
- Este método es compatible con la API de Binance y exchanges similares
