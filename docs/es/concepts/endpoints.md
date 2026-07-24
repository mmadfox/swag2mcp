# Endpoints

Un endpoint es un método HTTP + ruta específico que puede ser invocado (por ejemplo, `GET /api/users/{id}`). Los endpoints son las operaciones de API reales que el LLM descubre, inspecciona y llama.

## Estructura

Cada endpoint contiene:

- **Método HTTP**: GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS
- **Ruta**: `/api/v1/users/{id}`
- **Resumen**: una descripción corta de lo que hace el endpoint — muy útil para que el LLM entienda su propósito de un vistazo
- **Descripción**: una explicación detallada del comportamiento, parámetros y casos de uso del endpoint
- **Parámetros**: ruta, consulta, encabezado, cookie
- **Cuerpo de solicitud**: para POST/PUT/PATCH
- **Respuestas**: códigos de estado y esquemas de respuesta

Los campos `summary` y `description` provienen del archivo OpenAPI/Swagger/Postman. Son la forma principal en que el LLM entiende lo que hace un endpoint. Los resúmenes bien escritos hacen que el descubrimiento de endpoints sea mucho más efectivo.

## Herramientas MCP para Endpoints

| Herramienta | Descripción |
|-------------|-------------|
| `endpoint_by_spec` | Todos los endpoints en una especificación |
| `endpoint_by_collection` | Endpoints en una colección |
| `endpoint_by_tag` | Endpoints en una etiqueta |
| `endpoint_by_id` | Resumen rápido del endpoint |
| `inspect` | Detalles completos del endpoint (esquemas, parámetros) |
| `invoke` | Llamar al endpoint |
| `search` | Buscar endpoints por texto |

## Endpoints Obsoletos

Los endpoints marcados como `deprecated` en la especificación se muestran con un aviso al ser inspeccionados.

## Configuración

Los endpoints son **de solo lectura** desde la perspectiva de swag2mcp. No hay configuraciones YAML para endpoints — no puede agregar, eliminar, renombrar ni modificarlos en `swag2mcp.yaml`.

Para cambiar endpoints (agregar nuevos, actualizar resúmenes, modificar parámetros, marcar como obsoletos), edite el archivo OpenAPI/Swagger/Postman original y ejecute `swag2mcp update` para volver a analizar y reindexar.

## Ejemplo

```
Consulta: "Muestra detalles para GET /pet/{petId}"
→ inspect(endpointId: "abc123...")
→ Resultado:
  GET /pet/{petId}
  Resumen: Find pet by ID
  Descripción: Returns a single pet by its ID
  Parámetros:
    - petId (path, integer, required)
  Respuestas:
    - 200: Pet object
    - 400: Error
    - 404: Not found
```
