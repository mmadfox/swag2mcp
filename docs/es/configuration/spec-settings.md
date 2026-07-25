# Configuración de Especificaciones

La configuración de especificaciones define un servicio de API y anula la configuración global para esa API específica. Cada especificación representa una API lógica (por ejemplo, "Open-Meteo Weather APIs") y puede contener múltiples colecciones (archivos de especificación).

## Sección de Especificación

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    llm_instruction: "Use esta API para pronósticos meteorológicos y datos climáticos"
    base_url: https://api.open-meteo.com
    disable: false
    tags: ["weather", "climate"]
    http_client:
      timeout: 10s
      max_response_size: 1024
    auth:
      type: bearer
      config:
        token: "$(TOKEN)"
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

## Parámetros

### domain

- **Tipo:** `string`
- **Requerido:** Sí
- **Descripción:** Identificador único para esta especificación de API. Se usa internamente para referenciar la especificación.
- **Reglas:** 1-60 caracteres. Solo letras minúsculas (`a-z`), dígitos (`0-9`), guiones (`-`) y guiones bajos (`_`).
- **Ejemplo:** `meteo`, `binance`, `my-api`

### llm_title

- **Tipo:** `string`
- **Requerido:** Sí
- **Descripción:** Nombre legible por humanos que el LLM usa para referenciar esta API. Se muestra en las respuestas de las herramientas MCP.
- **Reglas:** 5-120 caracteres. Solo letras, dígitos, espacios y puntuación básica.
- **Ejemplo:** `Open-Meteo Weather APIs`, `Binance Market Data`

### llm_instruction

- **Tipo:** `string`
- **Valor predeterminado:** `""`
- **Descripción:** Instrucciones para el LLM sobre cómo usar esta API. Describe qué hace la API y cuándo usarla.
- **Reglas:** Máx. 500 caracteres. Solo letras, dígitos, espacios y puntuación básica.
- **Ejemplo:** `"Use esta API para pronósticos meteorológicos, condiciones actuales y datos climáticos."`

### base_url

- **Tipo:** `string`
- **Requerido:** Sí
- **Descripción:** URL base para todas las solicitudes de API en esta especificación. Las rutas de endpoint de la especificación OpenAPI se agregan a esta URL.
- **Ejemplo:** `https://api.open-meteo.com`, `https://api.binance.com`
- **Nota:** Puede anularse a nivel de colección si diferentes colecciones usan diferentes URL base.

### disable

- **Tipo:** `bool`
- **Valor predeterminado:** `false`
- **Descripción:** Cuando es `true`, esta especificación se excluye de las herramientas MCP. No se carga, indexa ni está disponible para el LLM.
- **Cuándo usarlo:** Deshabilitar temporalmente una API sin eliminarla de la configuración. Útil para APIs que están caídas, obsoletas o en mantenimiento.

### tags

- **Tipo:** `[]string` (arreglo de cadenas)
- **Valor predeterminado:** `[]`
- **Descripción:** Etiquetas para filtrar especificaciones. Se usan con la bandera `--tags` en comandos CLI (`ls`, `validate`, `mcp`, `update`).
- **Ejemplo:** `["public", "weather"]`, `["internal", "production"]`
- **Efecto:** Cuando ejecuta `swag2mcp mcp --tags=public`, solo se cargan las especificaciones con la etiqueta `public`.

### http_client

- **Tipo:** `object`
- **Valor predeterminado:** hereda de global
- **Descripción:** Anular la configuración global del cliente HTTP para esta especificación. Todas las configuraciones del `http_client` global pueden anularse: `timeout`, `max_response_size`, `user_agent`, `follow_redirects`, `max_redirects`, `random`, `proxy`, `headers`, `cookies`.
- **Ejemplo:**
  ```yaml
  http_client:
    timeout: 60s
    max_response_size: 4194304
    headers:
      "X-DC": "us-east-1"
  ```

### auth

- **Tipo:** `object`
- **Valor predeterminado:** `none` (sin autenticación)
- **Descripción:** Configuración de autenticación para esta especificación. Consulte la sección [Autenticación](/auth/overview) para los 9 métodos y sus parámetros.
- **Ejemplo:**
  ```yaml
  auth:
    type: bearer
    config:
      token: "$(API_TOKEN)"
  ```

### collections

- **Tipo:** `[]object` (arreglo de colecciones)
- **Requerido:** Sí (al menos 1)
- **Descripción:** Lista de archivos de especificación OpenAPI/Swagger/Postman que pertenecen a esta especificación. Cada colección es un archivo de especificación.
- **Reglas:** 1-30 colecciones por especificación.
- **Ver:** [Configuración de Colecciones](./collection-settings) para todos los parámetros de colección.

## Deshabilitar una Especificación

Las especificaciones deshabilitadas no se cargan ni indexan. El LLM no puede verlas ni usarlas.

```yaml
specs:
  - domain: old-api
    llm_title: Old API
    base_url: https://old-api.example.com
    disable: true
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Anulación del Cliente HTTP

Todas las configuraciones de `http_client` del nivel global pueden anularse a nivel de especificación. Los valores de la especificación tienen prioridad sobre los valores globales solo para esta especificación.

```yaml
specs:
  - domain: slow-api
    llm_title: Slow API
    base_url: https://slow-api.example.com
    http_client:
      timeout: 120s
      max_response_size: 8388608
      headers:
        "X-DC": "us-east-1"
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Anulación del Proxy

Si esta especificación requiere un proxy diferente al global, configúrelo a nivel de especificación:

```yaml
specs:
  - domain: proxied-api
    llm_title: Proxied API
    base_url: https://api.example.com
    http_client:
      proxy:
        url: http://proxy.company.com:8080
        username: $(PROXY_USER)
        password: $(PROXY_PASS)
        bypass:
          - "*.local"
          - "10.0.0.0/8"
    collections:
      - llm_title: Main
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo.json
```
