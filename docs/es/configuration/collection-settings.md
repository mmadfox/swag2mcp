# Configuración de Colecciones

La configuración de colecciones define un único archivo de especificación OpenAPI/Swagger/Postman y anula la configuración de la especificación para ese archivo específico. Cada colección pertenece a una especificación y representa un documento de especificación de API.

## Sección de Colección

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        llm_instruction: "Use para datos climáticos actuales y de pronóstico"
        disable: false
        base_url: https://forecast-api.open-meteo.com
        base_mock_url: localhost:8081
        http_client:
          timeout: 5s
```

## Parámetros

### llm_title

- **Tipo:** `string`
- **Requerido:** No
- **Descripción:** Nombre legible por humanos para esta colección. Se muestra en las respuestas de las herramientas MCP.
- **Reglas:** Máx. 120 caracteres. Solo letras, dígitos, espacios y puntuación básica.
- **Ejemplo:** `Forecast`, `Air Quality`, `Market Data`

### llm_instruction

- **Tipo:** `string`
- **Valor predeterminado:** `""`
- **Descripción:** Instrucciones para el LLM sobre esta colección específica. Describe qué endpoints proporciona esta colección.
- **Reglas:** Máx. 360 caracteres. Solo letras, dígitos, espacios y puntuación básica.
- **Ejemplo:** `"Use para datos climáticos actuales y de pronóstico."`

### title

- **Tipo:** `string`
- **Valor predeterminado:** `""`
- **Descripción:** Título original del archivo de especificación. Se rellena automáticamente en tiempo de ejecución. Normalmente no necesita establecer esto en YAML.

### location

- **Tipo:** `string`
- **Requerido:** Sí
- **Descripción:** URL o ruta de archivo local al archivo de especificación OpenAPI 3.x, Swagger 2.0 o Postman collection.
- **Reglas:** 5-250 caracteres.
- **Ejemplos:**
  - URL: `https://raw.githubusercontent.com/org/repo/main/spec.yaml`
  - Local: `./specs/my-api.json`
  - Local (absoluto): `/home/user/.swag2mcp/specs/my-api.yaml`

### disable

- **Tipo:** `bool`
- **Valor predeterminado:** `false`
- **Descripción:** Cuando es `true`, esta colección se excluye de las herramientas MCP. No se carga ni indexa.
- **Cuándo usarlo:** Deshabilitar temporalmente una colección sin eliminarla de la configuración. Útil cuando un archivo de especificación se está actualizando o una versión de API está obsoleta.

### http_client

- **Tipo:** `object`
- **Valor predeterminado:** hereda de la especificación (o global)
- **Descripción:** Anular la configuración del cliente HTTP para esta colección. Todas las configuraciones del `http_client` global pueden anularse: `timeout`, `max_response_size`, `user_agent`, `follow_redirects`, `max_redirects`, `random`, `proxy`, `headers`, `cookies`.
- **Ejemplo:**
  ```yaml
  http_client:
    timeout: 120s
    headers:
      "X-Custom": "value"
    cookies:
      - name: "session"
        value: "abc123"
  ```

### base_url

- **Tipo:** `string`
- **Valor predeterminado:** `""` (hereda de la especificación)
- **Descripción:** Anular la `base_url` a nivel de especificación para esta colección. Use cuando diferentes colecciones dentro de la misma especificación usen diferentes URL base.
- **Ejemplo:** Si la especificación tiene `base_url: https://api.open-meteo.com` pero una colección usa `https://air-quality-api.open-meteo.com`, establezca `base_url` a nivel de colección.

### base_mock_url

- **Tipo:** `string`
- **Valor predeterminado:** `""`
- **Descripción:** Dirección del servidor simulado en formato `host:port`. Requerido cuando `mock_enabled: true` en la configuración global.
- **Reglas:** El host debe ser `localhost`, `127.0.0.1` o `0.0.0.0`. El puerto debe ser un número de puerto válido.
- **Ejemplo:** `localhost:8081`, `127.0.0.1:9000`
- **Cuándo usarlo:** Tiene `mock_enabled: true` y desea probar esta colección con respuestas falsas.

## Múltiples Colecciones de una Especificación

Una especificación puede tener múltiples colecciones — por ejemplo, cuando una API tiene archivos de especificación separados para diferentes servicios:

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
      - llm_title: Marine
        base_url: https://marine-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
```

## Deshabilitar una Colección

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
        disable: true
```

## Anulación del Cliente HTTP

Todas las configuraciones de `http_client` pueden anularse a nivel de colección. Los valores de la colección tienen prioridad sobre los valores de especificación y global solo para esta colección.

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        http_client:
          timeout: 120s
          headers:
            "X-Custom": "value"
          cookies:
            - name: "session"
              value: "abc123"
```
