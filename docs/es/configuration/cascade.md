# Cascada de Configuración

swag2mcp utiliza una cascada de configuración de tres niveles. Cada nivel anula el anterior. Esto le permite establecer valores predeterminados sensatos a nivel global y ajustar la configuración para especificaciones o colecciones específicas.

## Niveles

```
Global (http_client, mcp, mock_enabled, disable_ratelimiter, rate_limit_interval)
    ↓ anula
Especificación (specs[].http_client, specs[].auth, specs[].base_url, specs[].disable, specs[].tags)
    ↓ anula
Colección (specs[].collections[].http_client, specs[].collections[].base_url, specs[].collections[].disable)
```

## Qué Anula Qué

| Parámetro | Global | Especificación | Colección |
|-----------|--------|----------------|-----------|
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

Todas las configuraciones de `http_client` pueden anularse en cada nivel. Las configuraciones a nivel de colección tienen prioridad total sobre las de especificación y global.

## Ejemplo de Cascada

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
      timeout: 60s  # anula el tiempo de espera global
      headers:
        "X-API-Version": "2"  # se agrega a los encabezados globales
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        http_client:
          timeout: 120s  # anula el tiempo de espera de la especificación
          headers:
            "X-Custom": "value"  # se agrega a los encabezados de especificación + globales
```

## Configuración Efectiva para la Colección "Forecast"

```
timeout: 120s (de la colección, anula especificación 60s y global 30s)
max_response_size: 1048576 (de global)
headers:
  - User-Agent: swag2mcp/1.0 (de global)
  - X-API-Version: 2 (de especificación)
  - X-Custom: value (de colección)
```

## Cómo Funciona la Fusión

### Configuraciones del Cliente HTTP

Los valores simples (`timeout`, `max_response_size`, `user_agent`, `follow_redirects`, `max_redirects`, `random`) se **reemplazan** en cada nivel. Si una especificación establece `timeout: 60s`, reemplaza completamente el global `30s`.

### Encabezados

Los encabezados se **fusionan** entre niveles. Los encabezados de los tres niveles se combinan. Si la misma clave de encabezado aparece en múltiples niveles, el nivel más bajo gana.

### Cookies

Las cookies se **fusionan** entre niveles. Si el mismo nombre de cookie aparece en múltiples niveles, el nivel más bajo gana.

### Proxy

El proxy se **reemplaza** en cada nivel. Si una especificación establece un proxy, reemplaza completamente el proxy global para esa especificación.
