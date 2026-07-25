# Herramientas de Endpoints

Las herramientas de endpoints permiten al LLM ver los endpoints de API en diferentes niveles de la jerarquía: todos los endpoints en una especificación, en una colección, en una etiqueta o un resumen de un solo endpoint. Úselas para descubrir operaciones disponibles antes de inspeccionar o invocar.

---

## endpoint_by_spec

### Propósito

Listar todos los endpoints en una especificación completa, abarcando todas las colecciones y etiquetas. Devuelve la vista más completa — cada endpoint en la especificación con su contexto completo (etiqueta, colección, especificación).

### Cuándo usarlo

- Cuando desea ver cada endpoint disponible en una especificación
- Cuando no sabe qué colección o etiqueta contiene el endpoint que necesita
- Después de `spec_by_id` para obtener la lista completa de endpoints

### Parámetros

| Parámetro | Tipo | Requerido | Descripción |
|-----------|------|-----------|-------------|
| `specId` | string | Sí | Hash MD5 de 32 caracteres de la especificación |

### Respuesta

```json
{
  "endpoints": [
    {
      "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
      "tagId": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "tagName": "forecast",
      "collectionId": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
      "collectionTitle": "Weather Forecast",
      "specId": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
      "specDomain": "meteo",
      "method": "GET",
      "path": "/v1/forecast",
      "summary": "Get weather forecast for a location"
    }
  ]
}
```

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `id` | string | Identificador del endpoint |
| `tagId` | string | Identificador de la etiqueta padre |
| `tagName` | string | Nombre de etiqueta legible por humanos |
| `collectionId` | string | Identificador de la colección padre |
| `collectionTitle` | string | Título de colección legible por humanos |
| `specId` | string | Identificador de la especificación padre |
| `specDomain` | string | Nombre de dominio de la especificación |
| `method` | string | Método HTTP (GET, POST, PUT, DELETE, etc.) |
| `path` | string | Ruta de la API (por ejemplo, /v1/forecast) |
| `summary` | string | Resumen legible por humanos de lo que hace el endpoint |

### Matices

- Devuelve `not_found` si la especificación no existe
- Cada endpoint incluye su ascendencia completa (especificación → colección → etiqueta) para contexto
- Para un resumen rápido de un solo endpoint, use `endpoint_by_id`

---

## endpoint_by_collection

### Propósito

Listar todos los endpoints dentro de una colección específica, independientemente de su etiqueta. Devuelve endpoints agrupados por colección con metadatos de especificación y colección.

### Cuándo usarlo

- Después de `collection_by_id` para ver todos los endpoints en una colección
- Cuando desea explorar la superficie de API completa de una colección

### Parámetros

| Parámetro | Tipo | Requerido | Descripción |
|-----------|------|-----------|-------------|
| `collectionId` | string | Sí | Hash MD5 de 32 caracteres de la colección |

### Respuesta

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collection": {
    "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
    "title": "Weather Forecast",
    "countMethods": 12
  },
  "endpoints": [
    {
      "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
      "tagId": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "tagName": "forecast",
      "method": "GET",
      "path": "/v1/forecast",
      "summary": "Get weather forecast for a location"
    }
  ]
}
```

### Matices

- Devuelve `not_found` si la colección no existe
- Incluye metadatos de especificación y colección para contexto
- Los endpoints de todas las etiquetas dentro de la colección se devuelven juntos

---

## endpoint_by_tag

### Propósito

Listar todos los endpoints agrupados bajo una etiqueta específica. Esta es la vista más enfocada — endpoints en una etiqueta dentro de una colección.

### Cuándo usarlo

- Después de `tag_by_id` para ver los endpoints reales en una etiqueta
- Cuando conoce la etiqueta y desea ver sus operaciones

### Parámetros

| Parámetro | Tipo | Requerido | Descripción |
|-----------|------|-----------|-------------|
| `tagId` | string | Sí | Hash MD5 de 32 caracteres de la etiqueta |

### Respuesta

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collection": {
    "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
    "title": "Weather Forecast",
    "countMethods": 12
  },
  "tag": {
    "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
    "title": "forecast",
    "countMethods": 5
  },
  "endpoints": [
    {
      "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
      "method": "GET",
      "path": "/v1/forecast",
      "summary": "Get weather forecast for a location"
    }
  ]
}
```

### Matices

- Devuelve `not_found` si la etiqueta no existe
- Incluye contexto completo: metadatos de especificación, colección y etiqueta
- Los endpoints están limitados a una sola etiqueta dentro de una sola colección

---

## endpoint_by_id

### Propósito

Obtener un resumen rápido de un solo endpoint: método, ruta, resumen y estado de obsolescencia. Esta es una herramienta ligera — para el objeto de operación OpenAPI completo (parámetros, cuerpo de solicitud, esquemas de respuesta), use `inspect`.

### Cuándo usarlo

- Cuando tiene un ID de endpoint y desea un recordatorio rápido de lo que hace
- Antes de decidir si llamar a `inspect` para detalles completos
- Cuando necesita confirmar el método y la ruta antes de invocar

### Parámetros

| Parámetro | Tipo | Requerido | Descripción |
|-----------|------|-----------|-------------|
| `id` | string | Sí | Hash MD5 de 32 caracteres del endpoint |

### Respuesta

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collection": {
    "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
    "title": "Weather Forecast",
    "countMethods": 12
  },
  "tag": {
    "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
    "title": "forecast",
    "countMethods": 5
  },
  "endpoint": {
    "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
    "method": "GET",
    "path": "/v1/forecast",
    "summary": "Get weather forecast for a location"
  }
}
```

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `endpoint.id` | string | Identificador del endpoint |
| `endpoint.method` | string | Método HTTP |
| `endpoint.path` | string | Ruta de la API |
| `endpoint.summary` | string | Resumen legible por humanos |

### Matices

- Devuelve `not_found` si el endpoint no existe
- Este es un **resumen rápido** — no devuelve parámetros, cuerpo de solicitud ni esquemas de respuesta
- Para detalles técnicos completos (requeridos antes de `invoke`), use `inspect`
