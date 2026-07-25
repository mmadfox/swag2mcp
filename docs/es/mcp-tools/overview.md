# Herramientas MCP

## Descripción General

swag2mcp proporciona **19 herramientas MCP** que dan a un agente LLM acceso completo a sus APIs a través del Protocolo de Contexto de Modelo. Estas herramientas cubren el flujo de trabajo completo: descubrir qué APIs están disponibles, navegar por la jerarquía de especificaciones, buscar e inspeccionar endpoints, ejecutar llamadas a la API y trabajar con respuestas grandes.

### Qué resuelven las herramientas

- **Descubrimiento** — el LLM puede encontrar especificaciones, colecciones y etiquetas sin conocer los IDs de antemano
- **Navegación** — profundizar desde especificación → colección → etiqueta → endpoint en una jerarquía estructurada
- **Búsqueda** — búsqueda de texto completo en todos los endpoints cuando no tiene un ID
- **Inspección** — obtener el objeto de operación OpenAPI completo antes de hacer una llamada
- **Ejecución** — invocar llamadas reales a la API con autenticación automática
- **Manejo de respuestas grandes** — esquematizar, comprimir y segmentar respuestas demasiado grandes que no caben en línea

### Solo lectura vs Mutables

| Tipo | Cantidad | Herramientas |
|------|----------|--------------|
| **Solo lectura** | 17 | Todas las herramientas de descubrimiento, endpoints, búsqueda, inspección, información y respuesta |
| **Mutables** | 2 | `invoke` (realiza llamadas HTTP reales), `auth` (recupera tokens) |

Las herramientas de solo lectura están marcadas con `ReadOnlyHint=true` y `IdempotentHint=true` en el protocolo MCP, indicando al LLM que son seguras de llamar sin efectos secundarios.

### Manejo de errores

Todas las herramientas devuelven errores como objetos `LLMError` estructurados con un código legible por máquina y un mensaje legible por humanos que explica qué salió mal y qué hacer a continuación:

| Código de error | Significado |
|-----------------|-------------|
| `validation_failed` | Entrada inválida (formato de ID incorrecto, campos requeridos faltantes) |
| `not_found` | Entidad no encontrada en el índice o espacio de trabajo |
| `rate_limit` | Segunda llamada `invoke` dentro de 10 segundos en el mismo endpoint |
| `invoke_error` | Fallo de llamada HTTP, fallo de descarga |
| `auth_error` | Fallo de recuperación de token de autenticación |
| `config_error` | Fallo de carga o guardado del archivo de configuración |
| `parse_error` | Fallo de análisis del archivo de especificación |

## Categorías

| Categoría | Herramientas | Descripción |
|-----------|--------------|-------------|
| **Descubrimiento** | `spec_list`, `spec_by_id`, `collection_by_spec`, `collection_by_id`, `tag_by_spec`, `tag_by_collection`, `tag_by_id` | Navegar por la jerarquía de especificaciones: encontrar especificaciones, colecciones y etiquetas |
| **Endpoints** | `endpoint_by_spec`, `endpoint_by_collection`, `endpoint_by_tag`, `endpoint_by_id` | Ver endpoints en diferentes niveles de la jerarquía |
| **Ejecución** | `search`, `inspect`, `invoke` | Buscar, inspeccionar el contrato completo y llamar APIs |
| **Utilidades** | `auth`, `info`, `response_outline`, `response_compress`, `response_slice` | Tokens de autenticación, información de ejecución y manejo de respuestas grandes |
| **Habilidades** | [Guía de formato](/mcp-tools/skills) | Personalizar cómo se muestran las respuestas de las herramientas |

## Lista Completa

| Herramienta | Descripción |
|-------------|-------------|
| `spec_list` | Listar todas las especificaciones de API en el espacio de trabajo |
| `spec_by_id` | Obtener información detallada de la especificación con colecciones |
| `collection_by_spec` | Listar colecciones dentro de una especificación |
| `collection_by_id` | Obtener detalles de la colección con etiquetas |
| `tag_by_spec` | Listar todas las etiquetas en una especificación |
| `tag_by_collection` | Listar etiquetas dentro de una colección |
| `tag_by_id` | Obtener detalles de la etiqueta (ID, título, recuento de métodos) |
| `endpoint_by_spec` | Listar todos los endpoints en una especificación |
| `endpoint_by_collection` | Listar endpoints en una colección |
| `endpoint_by_tag` | Listar endpoints en una etiqueta |
| `endpoint_by_id` | Resumen rápido del endpoint (método, ruta, resumen) |
| `search` | Búsqueda de texto completo en todos los endpoints |
| `inspect` | Detalles completos de la operación OpenAPI (parámetros, esquemas) |
| `invoke` | Ejecutar una llamada real a la API |
| `auth` | Obtener token o encabezados de autenticación para una especificación |
| `info` | Información de ejecución (versión, especificaciones, configuración) |
| `response_outline` | Resumen estructural de un archivo de respuesta grande |
| `response_compress` | Comprimir una respuesta grande para que quepa en línea |
| `response_slice` | Extraer un fragmento de una respuesta grande |

## Jerarquía de Navegación

```
spec_list
  └── spec_by_id(id)
        └── collection_by_spec(specId)
              └── collection_by_id(id)
                    └── tag_by_collection(collectionId)
                          └── tag_by_id(id)
                                └── endpoint_by_tag(tagId)
                                      └── endpoint_by_id(id)
                                            └── inspect(endpointId)
                                                  └── invoke(endpointId)
```

Cuando no tiene un ID, use `search` para encontrar endpoints por consulta. Cuando `invoke` devuelve un `fileRef` (respuesta demasiado grande), use `response_outline` → `response_compress` o `response_slice` para explorar los datos.
