# Configuración Global

La configuración global son los bloques de configuración de nivel superior en `swag2mcp.yaml`. Se aplican a todas las especificaciones a menos que se anulen a nivel de especificación o colección.

## Estructura

```yaml
http_client:
  # Configuración del cliente HTTP para todas las llamadas a la API

mcp:
  # Configuración del servidor MCP

mock_enabled: false
mock_auth:
  oauth2_port: 9090
  digest_port: 9091
  hmac_port: 9092

disable_ratelimiter: false
rate_limit_interval: 10s
```

## Cliente HTTP

Controla cómo swag2mcp realiza solicitudes HTTP a las APIs: tiempo de espera, límite de tamaño de respuesta, proxy, encabezados, cookies, redirecciones y agente de usuario. Estas configuraciones se transmiten en cascada a las especificaciones y colecciones.

Consulte [Cliente HTTP](./http-client) para todos los parámetros y ejemplos.

## Servidor MCP

Controla cómo el servidor MCP se comunica con los agentes LLM: tipo de transporte (stdio, SSE, HTTP Streamable), dirección, ruta y autenticación opcional con token bearer.

Consulte [Servidor MCP](./mcp-server) para todos los parámetros, transportes y banderas de inicio.

## Servidor Simulado

El servidor simulado genera respuestas de API falsas basadas en esquemas OpenAPI. Útil para pruebas sin golpear APIs reales.

```yaml
mock_enabled: true
mock_auth:
  oauth2_port: 9090
  digest_port: 9091
  hmac_port: 9092
```

### mock_enabled

- **Tipo:** `bool`
- **Valor predeterminado:** `false`
- **Efecto:** Cuando es `true`, swag2mcp inicia servidores simulados para todas las especificaciones que tienen `base_mock_url` configurado. Cada colección debe tener `base_mock_url` establecido.
- **Cuándo habilitar:** Desea probar su integración de API sin realizar llamadas HTTP reales. Los servidores simulados devuelven datos falsos basados en el esquema OpenAPI.

### mock_auth

Configuración de puertos para servidores de autenticación simulados. Se usan al probar métodos de autenticación (OAuth2, Digest, HMAC) con el servidor simulado.

| Campo | Tipo | Valor predeterminado | Descripción |
|-------|------|---------------------|-------------|
| `oauth2_port` | int | `9090` | Puerto para el servidor de token OAuth2 simulado (1024-65535) |
| `digest_port` | int | `9091` | Puerto para el servidor de autenticación Digest simulado (1024-65535) |
| `hmac_port` | int | `9092` | Puerto para el servidor de autenticación HMAC simulado (1024-65535) |

## Limitador de Velocidad

El limitador de velocidad evita que el LLM llame al mismo endpoint de API con demasiada frecuencia. Por defecto, cada endpoint puede llamarse una vez cada 10 segundos.

```yaml
disable_ratelimiter: false
rate_limit_interval: 10s
```

### disable_ratelimiter

- **Tipo:** `bool`
- **Valor predeterminado:** `false`
- **Efecto:** Cuando es `true`, el limitador de velocidad por endpoint se deshabilita por completo. El LLM puede llamar al mismo endpoint repetidamente sin esperar.
- **Cuándo habilitar:** Pruebas, depuración, o cuando necesita llamar al mismo endpoint múltiples veces en rápida sucesión.
- **Cuándo mantenerlo deshabilitado (recomendado):** Producción. El limitador de velocidad evita el abuso accidental y respeta los límites de velocidad de la API.

### rate_limit_interval

- **Tipo:** duración (formato Go: `10s`, `30s`, `1m`)
- **Valor predeterminado:** `10s`
- **Efecto:** Establece cuánto tiempo debe esperar el LLM entre llamadas al mismo endpoint.
- **Cuándo cambiar:** Aumentar para APIs con límites de velocidad estrictos. Disminuir para APIs internas donde usted controla la carga.
- **Rango:** Cualquier duración válida (por ejemplo, `5s`, `30s`, `1m`, `2m`).

## Cascada

La configuración global puede anularse a nivel de especificación y colección. Todas las configuraciones de `http_client` (tiempo de espera, proxy, agente de usuario, redirecciones, tamaño de respuesta, aleatorizador, encabezados, cookies) pueden anularse tanto a nivel de especificación como de colección.

```
Global (http_client, mock_enabled, disable_ratelimiter, rate_limit_interval)
    ↓ anula (solo http_client)
Especificación (specs[].http_client)
    ↓ anula (solo http_client)
Colección (specs[].collections[].http_client)
```

Consulte [Cascada de Configuración](./cascade) para más detalles.
