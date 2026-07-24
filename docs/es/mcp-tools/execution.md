# Herramientas de Ejecución

Las herramientas de ejecución son el núcleo de swag2mcp: **search** encuentra endpoints cuando no tiene un ID, **inspect** revela el contrato OpenAPI completo, e **invoke** realiza la llamada real a la API. Úselas siempre en este orden: search → inspect → invoke.

---

## search

### Propósito

La única herramienta para encontrar endpoints cuando no tiene un ID de endpoint. Realiza búsqueda de texto completo en todos los endpoints de todas las especificaciones usando el motor de búsqueda bluge.

### Cuándo usarlo

- Cuando no conoce el ID del endpoint
- Cuando desea encontrar endpoints por palabras clave, método, etiqueta o ruta
- Cuando necesita descubrir qué endpoints existen para una característica específica

### Cómo funciona

Busca en el índice de texto completo en todas las especificaciones. Admite consultas estructuradas con filtros de campo, operadores booleanos, coincidencia difusa, comodines y más.

### Parámetros

| Parámetro | Tipo | Requerido | Descripción |
|-----------|------|-----------|-------------|
| `query` | string | Sí | Consulta de búsqueda (admite sintaxis estructurada) |
| `limit` | int | Sí | Resultados máximos a devolver (1-50) |

### Sintaxis de consulta

| Ejemplo | Descripción |
|---------|-------------|
| `pet` | Búsqueda de texto simple en todos los campos |
| `method:GET` | Filtrar por método HTTP |
| `tag:pet` | Filtrar por nombre de etiqueta |
| `path:"/api/v1/users"` | Búsqueda de ruta exacta |
| `+method:POST +tag:pet` | Debe coincidir con ambas condiciones |
| `-method:DELETE` | Excluir métodos DELETE |
| `create~` | Búsqueda difusa (tolerante a errores tipográficos) |
| `path:/api/v1/*` | Búsqueda de ruta con comodín |
| `/pattern/` | Búsqueda regex |
| `term^3` | Aumentar la relevancia de un término |

**Campos buscables:** `method` (palabra clave), `tag` (palabra clave), `path` (texto), `summary` (texto), `_all` (campo de texto predeterminado).

**No soportado:** paréntesis para agrupación, operadores `AND`/`OR` explícitos, agrupación de campos.

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

Cada resultado incluye la ascendencia completa (especificación → colección → etiqueta) para que el LLM pueda navegar a endpoints relacionados.

### Matices

- `limit` debe estar entre 1 y 50 (devuelve `validation_failed` de lo contrario)
- `query` es requerido (devuelve `validation_failed` si está vacío)
- Los resultados se devuelven en orden de relevancia (mejor coincidencia primero)
- Use filtros de campo (`method:GET`, `tag:pet`) para acotar los resultados
- Para coincidencia exacta de ruta, use comillas: `path:"/v1/forecast"`

---

## inspect

### Propósito

Recuperar el objeto de operación OpenAPI completo para un endpoint: todos los parámetros, esquema del cuerpo de solicitud, esquemas de respuesta, URL base y URL completa. Esta es la herramienta a llamar **antes** de `invoke` para entender el contrato del endpoint.

### Cuándo usarlo

- Siempre antes de `invoke` — necesita el contrato completo para hacer una llamada correcta
- Cuando necesita explicar los detalles técnicos de una API al usuario
- Cuando necesita conocer los parámetros requeridos, la estructura del cuerpo de solicitud o el formato de respuesta

### Cómo funciona

Busca el endpoint en el índice y devuelve el objeto de operación OpenAPI completo con todos los esquemas resueltos.

### Parámetros

| Parámetro | Tipo | Requerido | Descripción |
|-----------|------|-----------|-------------|
| `endpointId` | string | Sí | Hash MD5 de 32 caracteres del endpoint |

### Respuesta

```json
{
  "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
  "tagId": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
  "collectionId": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
  "specId": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
  "specDomain": "meteo",
  "method": "POST",
  "path": "/pet",
  "baseUrl": "https://meteo.swagger.io/v2",
  "fullUrl": "https://meteo.swagger.io/v2/pet",
  "operation": {
    "id": "addPet",
    "tags": ["pet"],
    "summary": "Add a new pet",
    "description": "Add a new pet to the store",
    "deprecated": false,
    "parameters": [
      {
        "name": "petId",
        "in": "path",
        "description": "ID of the pet",
        "required": true,
        "schema": {
          "type": "integer",
          "format": "int64"
        }
      }
    ],
    "requestBody": {
      "description": "Pet object to add",
      "required": true,
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "name": { "type": "string" },
              "status": { "type": "string", "enum": ["available", "pending", "sold"] }
            },
            "required": ["name"]
          }
        }
      }
    },
    "responses": {
      "200": {
        "description": "Successful operation",
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/Pet"
            }
          }
        }
      },
      "405": {
        "description": "Invalid input"
      }
    }
  }
}
```

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `baseUrl` | string | URL base de la API (de la configuración) |
| `fullUrl` | string | URL completa del endpoint (base + path) |
| `operation.parameters[]` | array | Parámetros con nombre, ubicación (path/query/header/cookie), descripción, indicador required y esquema |
| `operation.requestBody` | object | Cuerpo de solicitud con tipo de contenido y esquema |
| `operation.responses` | map | Códigos de respuesta con descripciones y esquemas |
| `operation.deprecated` | bool | Indica si el endpoint está obsoleto |

### Matices

- Devuelve `not_found` si el endpoint no existe
- Esta es la **única** herramienta que devuelve la operación OpenAPI completa — `endpoint_by_id` devuelve solo un resumen
- Siempre llame a `inspect` antes de `invoke` para entender los parámetros requeridos y la estructura del cuerpo
- El objeto `operation` incluye referencias `$ref` que se resuelven a sus definiciones de esquema completas

---

## invoke

### Propósito

Ejecutar una llamada real a la API en un endpoint. Esta es la única herramienta que realiza solicitudes HTTP reales. La autenticación se aplica automáticamente — no necesita llamar a `auth` primero.

### Cuándo usarlo

- Solo después de llamar a `inspect` para entender el contrato del endpoint
- Solo con confirmación explícita del usuario para operaciones destructivas (POST, PUT, PATCH, DELETE)
- Cuando el usuario pide llamar a una API y usted tiene todos los parámetros requeridos

### Cómo funciona

1. Busca el endpoint en el índice
2. Sustituye los parámetros de ruta en la URL
3. Agrega los parámetros de consulta
4. Agrega encabezados y cookies
5. Serializa el cuerpo de la solicitud como JSON
6. Obtiene y aplica automáticamente la autenticación (token, encabezados, parámetros de consulta)
7. Realiza la solicitud HTTP
8. Devuelve la respuesta o la guarda en un archivo si es demasiado grande

### Parámetros

| Parámetro | Tipo | Requerido | Descripción |
|-----------|------|-----------|-------------|
| `endpointId` | string | Sí | Hash MD5 de 32 caracteres del endpoint |
| `parameters` | object | No | Parámetros de ruta, consulta y encabezado como pares clave-valor |
| `requestBody` | object | No | Cuerpo de solicitud para solicitudes POST/PUT/PATCH |
| `headers` | object | No | Encabezados HTTP adicionales a enviar |
| `cookies` | object | No | Cookies HTTP adicionales a enviar |

### Respuesta (en línea)

```json
{
  "statusCode": 200,
  "headers": {
    "content-type": "application/json"
  },
  "body": {
    "id": 1,
    "name": "Rex",
    "status": "available"
  }
}
```

### Respuesta (referencia de archivo — cuando el cuerpo excede el límite de tamaño)

```json
{
  "statusCode": 200,
  "headers": {
    "content-type": "application/json"
  },
  "fileRef": {
    "path": "/Users/user/.swag2mcp/responses/response_a1b2c3d4.json",
    "size": 1572864,
    "sizeHint": "1.5 MB",
    "maxSizeHint": "2 KB",
    "message": "La respuesta excede el límite de 2 KB y se ha guardado en disco.",
    "openCmd": "open /Users/user/.swag2mcp/responses/response_a1b2c3d4.json"
  }
}
```

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `statusCode` | int | Código de estado de la respuesta HTTP |
| `headers` | object | Encabezados de la respuesta HTTP |
| `body` | any | Cuerpo de la respuesta (presente cuando está dentro del límite de tamaño) |
| `fileRef` | object | Referencia de archivo (presente cuando el cuerpo excede el límite de tamaño) |

### Trabajar con respuestas grandes

Cuando `invoke` devuelve un `fileRef`, use las herramientas de respuesta para explorar los datos:

1. **`response_outline(path)`** — obtener el resumen estructural (claves, tipos, longitudes de arreglos)
2. **`response_compress(path, mode)`** — comprimir los datos para que quepan en línea
3. **`response_slice(path, jsonPath)`** — extraer un fragmento específico

### Matices

- **La autenticación es automática:** La herramienta `invoke` obtiene y aplica automáticamente la autenticación de la configuración de la especificación. **No** necesita llamar a `auth` primero.
- **Limitación de velocidad:** Cada endpoint tiene un enfriamiento de 10 segundos. Una segunda llamada al mismo endpoint dentro de 10 segundos se bloquea silenciosamente (devuelve error `rate_limit`).
- **Límite de tamaño de respuesta:** El valor predeterminado es 2 KB (configurable mediante `max_response_size`). Si la respuesta excede este límite, se guarda en `{workspace}/responses/` y se devuelve un `FileReference` en lugar del `body` en línea.
- **Manejo de parámetros:** Los parámetros de ruta se sustituyen en la URL. Los parámetros de consulta se agregan. Los parámetros de la solicitud anulan los valores predeterminados de la operación.
- **Cuerpo de solicitud:** Para POST/PUT/PATCH, el cuerpo se serializa como JSON. `Content-Type` se establece en `application/json` automáticamente.
- **Manejo de errores:** Los errores HTTP (no 2xx) se devuelven como `invoke_error` con el código de estado y el cuerpo de la respuesta en la sugerencia.
- **Operaciones destructivas:** Nunca invoque POST/PUT/PATCH/DELETE sin confirmación explícita del usuario.
