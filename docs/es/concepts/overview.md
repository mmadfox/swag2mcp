# Conceptos

## Arquitectura

swag2mcp actúa como un puente entre las especificaciones de API y los agentes LLM:

<img src="/architecture.svg" width="800" alt="Arquitectura de swag2mcp">

## Conceptos Clave

**Especificación** — un contenedor lógico que representa un dominio o servicio de API (por ejemplo, YouTube, Binance, Open-Meteo). Cada especificación tiene un `domain` único, una `base_url`, `auth` opcional y contiene una o más colecciones. También puede establecer `llm_instruction` — una sugerencia corta inyectada en el prompt del sistema de swag2mcp que le dice al LLM para qué sirve esta especificación y cuándo usarla. Más información: [Especificaciones](./specs).

**Colección** — un único archivo OpenAPI/Swagger/Postman que describe una API específica. Apunta a una `location` (URL o ruta de archivo local). Una especificación puede tener múltiples colecciones — por ejemplo, la especificación "meteo" podría tener colecciones "Pronóstico", "Calidad del Aire" y "Marino", cada una apuntando a un archivo de especificación diferente. Más información: [Colecciones](./collections).

**Etiqueta** — una categoría de endpoints dentro de una colección. Ayuda al LLM a encontrar las operaciones correctas con más precisión. Más información: [Etiquetas](./tags).

**Endpoint** — un método HTTP + ruta específico (por ejemplo, `GET /api/users`). El LLM puede encontrar un endpoint por descripción, inspeccionar sus parámetros y esquemas, y luego invocarlo. Más información: [Endpoints](./endpoints).

**Espacio de trabajo** — el directorio donde swag2mcp almacena la configuración, la caché de especificaciones, las respuestas guardadas y los scripts de autenticación. Más información: [Espacio de Trabajo](./workspace).

## Cómo funciona

1. **Agregue una especificación o colección** — defínala en la configuración YAML (`~/.swag2mcp/swag2mcp.yaml`). Por ejemplo:

   ```yaml
   specs:
     - domain: jokes
       llm_title: Dad Joke API
       base_url: https://icanhazdadjoke.com
       collections:
         - llm_title: Jokes
           location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
   ```
2. **swag2mcp analiza cada colección** — crea Etiquetas y Endpoints, los indexa para búsqueda.
3. **El LLM encuentra el endpoint correcto** — a través de herramientas MCP (`search`, `endpoint_by_tag`, `inspect`), el LLM busca un endpoint que coincida por descripción, revisa sus parámetros y esquema de solicitud.
4. **El LLM invoca el endpoint** — a través de la herramienta MCP `invoke`, el LLM envía la solicitud. swag2mcp valida cada parámetro de entrada contra el esquema OpenAPI del endpoint (parámetros de ruta, parámetros de consulta, encabezados, cuerpo de solicitud) antes de realizar la llamada. Si algo no coincide con el esquema, el LLM recibe un error claro explicando qué está mal. Una vez validado, swag2mcp ejecuta la llamada HTTP real y devuelve el resultado.
5. **El resultado vuelve al LLM** — la respuesta de la API se pasa de vuelta al agente. Las respuestas grandes se guardan en el espacio de trabajo y pueden explorarse con tres herramientas MCP dedicadas: `response_outline` (ver la estructura), `response_compress` (reducir a una muestra representativa) y `response_slice` (extraer fragmentos específicos).

swag2mcp es un puente entre los LLM y el mundo de las APIs. Usted agrega especificaciones de API, y el LLM — a través del protocolo MCP — encuentra los endpoints correctos, inspecciona su documentación y los llama. Todo lo que necesita hacer es agregar una especificación e iniciar el servidor MCP.

> **La configuración se puede editar en cualquier momento.** El archivo de configuración YAML (`~/.swag2mcp/swag2mcp.yaml`) se puede editar a mano — agregue especificaciones, cambie la autenticación, ajuste la configuración. Después de cada edición, reinicie el servidor MCP (`swag2mcp mcp`) para que los cambios surtan efecto.

## Jerarquía

```
Especificación (dominio, por ejemplo "meteo")
  └── Colección 1 (archivo de especificación, por ejemplo forecast.yml)
        └── Etiqueta 1 (categoría)
              └── Endpoint (GET /api/forecast)
              └── Endpoint (POST /api/forecast)
        └── Etiqueta 2
              └── Endpoint (GET /api/forecast/{id})
  └── Colección 2 (archivo de especificación, por ejemplo air-quality.yml)
        └── Etiqueta 3
              └── Endpoint (GET /api/air-quality)
```
