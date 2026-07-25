# Herramientas de Utilidad

Las herramientas de utilidad proporcionan funcionalidad de apoyo: recuperar tokens de autenticación, obtener información de ejecución y trabajar con respuestas grandes de API que no caben en línea.

---

## auth

### Propósito

Recuperar un token de autenticación, encabezados o parámetros de consulta para una especificación específica. Esto da al LLM acceso a credenciales que pueden usarse fuera de swag2mcp (por ejemplo, generar un comando curl).

### Cuándo usarlo

- Solo cuando el usuario solicita explícitamente el token o las credenciales sin procesar
- Al generar un comando curl o fragmento de código que necesita autenticación
- Cuando el usuario quiere ver qué método de autenticación está configurado

### Cuándo NO usarlo

- **No** llame a `auth` antes de `inspect` o `invoke` — `invoke` obtiene y aplica la autenticación automáticamente
- **No** llame a `auth` solo para verificar si la autenticación está configurada — use `info` en su lugar

### Cómo funciona

Busca la configuración de autenticación de la especificación y ejecuta el flujo de autenticación (intercambio de tokens, ejecución de script, etc.) para obtener las credenciales actuales.

### Parámetros

| Parámetro | Tipo | Requerido | Descripción |
|-----------|------|-----------|-------------|
| `specId` | string | Sí | Hash MD5 de 32 caracteres de la especificación |

### Respuesta

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "headers": {
    "Authorization": "Bearer eyJhbGciOiJIUzI1NiIs...",
    "X-API-Key": "my-api-key"
  },
  "queryParams": {
    "api_key": "my-api-key"
  }
}
```

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `token` | string | Valor del token sin procesar (token bearer, clave de API, etc.) |
| `headers` | object | Encabezados HTTP a incluir en las solicitudes |
| `queryParams` | object | Parámetros de consulta a incluir en las solicitudes |

### Matices

- **Deshabilitado por defecto en producción:** La bandera `--disable-llm-auth` (predeterminado: `true`) elimina la herramienta `auth` de la lista de herramientas MCP por completo. El LLM no puede ver ni solicitar tokens. Establezca `--disable-llm-auth=false` para habilitarla para depuración o tokens de corta duración.
- **`invoke` maneja la autenticación automáticamente:** No necesita llamar a `auth` antes de `invoke`. El servicio invoke obtiene y aplica automáticamente la autenticación correcta.
- **Admite 9 métodos de autenticación:** `none`, `basic`, `bearer`, `digest`, `hmac`, `oauth2-cc` (credenciales de cliente), `oauth2-pwd` (contraseña), `api-key`, `script`.
- Devuelve `auth_error` si el método de autenticación falla (por ejemplo, endpoint de token OAuth2 inalcanzable, fallo de ejecución de script).

---

## info

### Propósito

Devolver un resumen completo del entorno de ejecución de swag2mcp: versión, ruta del espacio de trabajo, especificaciones activas, configuración del cliente HTTP, configuración del transporte MCP, métodos de autenticación y estado del modo simulado.

### Cuándo usarlo

- Cuando el usuario pregunta sobre la configuración del sistema
- Cuando necesita verificar la configuración de ejecución (tiempo de espera, límite de tamaño de respuesta, transporte)
- Cuando necesita saber qué métodos de autenticación están disponibles
- Al solucionar problemas de configuración

### Cómo funciona

Devuelve una instantánea precalculada del estado de ejecución. No se necesitan parámetros.

### Parámetros

Ninguno.

### Respuesta

```json
{
  "version": "v1.2.0",
  "workspace": "~/.swag2mcp",
  "uptime": "2h 15m",
  "specs": {
    "total": 4,
    "active": 3,
    "disabled": 1,
    "collections": 6,
    "endpoints": 42
  },
  "http_client": {
    "timeout": "30s",
    "max_response_size": "2 KB",
    "follow_redirects": true,
    "max_redirects": 10,
    "randomize": false,
    "proxy": null,
    "headers": {},
    "cookies": []
  },
  "mcp": {
    "transport": "stdio",
    "addr": ":8080",
    "path": "/mcp",
    "auth_enabled": false
  },
  "auth": {
    "methods": ["bearer", "api-key"]
  },
  "mock": {
    "enabled": false
  }
}
```

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `version` | string | Versión de swag2mcp |
| `workspace` | string | Ruta del directorio del espacio de trabajo |
| `uptime` | string | Tiempo de actividad del servidor (legible por humanos) |
| `specs` | object | Resumen de especificaciones: total, activas, deshabilitadas, colecciones, endpoints |
| `http_client` | object | Configuración del cliente HTTP |
| `http_client.max_response_size` | string | Tamaño máximo de respuesta en formato legible por humanos (por ejemplo, "2 KB") |
| `mcp` | object | Configuración del servidor MCP |
| `auth` | object | Métodos de autenticación disponibles |
| `mock` | object | Estado del servidor simulado |

### Matices

- `max_response_size` se muestra en formato legible por humanos (por ejemplo, `"1 KB"`, `"2 MB"`)
- `uptime` se calcula desde la hora de inicio del servidor
- Los datos son una instantánea tomada en el momento del arranque — reflejan el estado cuando se inició el servidor MCP

---

## response_outline

### Propósito

Obtener un resumen estructural de alto nivel de un archivo de respuesta JSON grande que fue guardado en disco por `invoke`. Devuelve la forma de los datos — claves, tipos, longitudes de arreglos y sugerencias de navegación — sin devolver los valores reales.

### Cuándo usarlo

- Inmediatamente después de que `invoke` devuelva un `fileRef` (respuesta demasiado grande para en línea)
- Este es el **primer paso obligatorio** en el flujo de trabajo de respuestas grandes

### Cómo funciona

Lee el archivo de respuesta guardado y analiza su estructura: tipo de nivel superior, claves, longitudes de arreglos, profundidad de anidamiento y sugerencias de compresión.

### Parámetros

| Parámetro | Tipo | Requerido | Descripción |
|-----------|------|-----------|-------------|
| `path` | string | Sí | Ruta absoluta de `fileRef.path` |
| `maxDepth` | int | No | Profundidad máxima de recursión (predeterminado: 3) |
| `maxArrayItems` | int | No | Cuántos elementos del arreglo inspeccionar (predeterminado: 5) |

### Respuesta

```json
{
  "outline": {
    "type": "object",
    "size": 1572864,
    "lineCount": 12500,
    "depth": 3,
    "structure": {
      "type": "object",
      "keys": ["data", "meta", "error"],
      "data": {
        "type": "array",
        "length": 500,
        "items": {
          "type": "object",
          "keys": ["id", "name", "status", "createdAt"]
        }
      }
    },
    "schemaHint": "object with 3 keys: data (array[500]), meta (object), error (null)",
    "keys": ["data", "meta", "error"],
    "itemCount": 500,
    "itemType": "object",
    "compressionHints": [
      "response_compress(path, 'first_of_array', 'data')",
      "response_compress(path, 'sample_array', 'data', arrayHead=3, arrayTail=2)",
      "response_compress(path, 'keys_only', 'data')",
      "response_compress(path, 'select_keys', 'data', selectKeys=[id, name])"
    ],
    "navigationHints": {
      "paths": ["data", "meta", "error"],
      "arrays": [
        {"path": "data", "length": 500}
      ]
    }
  }
}
```

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `type` | string | Tipo de nivel superior: "object" o "array" |
| `size` | int | Tamaño del archivo en bytes |
| `lineCount` | int | Número de líneas en el archivo |
| `depth` | int | Profundidad máxima de anidamiento inspeccionada |
| `structure` | object | Estructura recursiva con claves, tipos, longitudes de arreglos |
| `schemaHint` | string | Resumen de una línea de la forma de nivel superior |
| `keys` | array | Claves de nivel superior (para objetos) |
| `itemCount` | int | Longitud del arreglo (para arreglos) |
| `compressionHints` | array | Llamadas sugeridas a `response_compress` con parámetros |
| `navigationHints` | object | Rutas de nivel superior y arreglos con longitudes |

### Matices

- Devuelve `validation_failed` si la ruta es inválida o no está dentro del directorio de respuestas
- Devuelve `not_found` si el archivo no existe
- Devuelve `validation_failed` si el archivo no es JSON válido
- El campo `compressionHints` proporciona sugerencias listas para usar para llamadas a `response_compress`

---

## response_compress

### Propósito

Reducir un valor JSON dentro de un archivo de respuesta guardado para que quepa dentro del límite de tamaño de respuesta y pueda devolverse al LLM en línea. Múltiples modos de compresión le permiten elegir el equilibrio adecuado entre tamaño e información.

### Cuándo usarlo

- Después de `response_outline` para entender la estructura
- Cuando necesita obtener datos de una respuesta grande en línea
- Cuando `response_slice` es demasiado limitado y necesita una vista más amplia

### Cómo funciona

Lee el archivo de respuesta guardado, navega a la ruta JSON especificada, aplica el modo de compresión y devuelve el resultado comprimido. Si el resultado aún excede el límite de tamaño, se guarda en un nuevo archivo.

### Parámetros

| Parámetro | Tipo | Requerido | Descripción |
|-----------|------|-----------|-------------|
| `path` | string | Sí | Ruta absoluta de `fileRef.path` |
| `jsonPath` | string | No | Ruta al valor a comprimir (por ejemplo, `data` o `data.0`) |
| `mode` | string | Sí | Modo de compresión (ver tabla abajo) |
| `arrayHead` | int | No | Elementos iniciales a conservar en modo `sample_array` (predeterminado: 3) |
| `arrayTail` | int | No | Elementos finales a conservar en modo `sample_array` (predeterminado: 2) |
| `stringLen` | int | No | Longitud máxima de cadena en modo `truncate_strings` (predeterminado: 80) |
| `selectKeys` | array | No | Claves a conservar en modo `select_keys` |

### Modos de compresión

| Modo | Descripción | Mejor para |
|------|-------------|------------|
| `first_of_array` | Conservar solo el primer elemento de un arreglo | Cuando todos los elementos tienen la misma estructura |
| `sample_array` | Conservar cabeza y cola de un arreglo | Cuando necesita ver el rango de valores |
| `truncate_strings` | Acortar cada cadena a `stringLen` caracteres | Cuando las cadenas son muy largas pero la estructura importa |
| `keys_only` | Reemplazar valores de objetos con sus nombres de tipo | Cuando solo necesita la estructura |
| `select_keys` | Conservar solo las claves especificadas en cada objeto | Cuando necesita campos específicos de muchos objetos |

### Respuesta

```json
{
  "body": [
    { "id": 1, "name": "Rex", "status": "available" },
    { "id": 2, "name": "Max", "status": "pending" }
  ],
  "hint": "Arreglo comprimido de 500 a 2 elementos usando el modo first_of_array"
}
```

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `body` | any | Valor JSON comprimido (presente cuando está dentro del límite de tamaño) |
| `fileRef` | object | Referencia de archivo (presente cuando aún es demasiado grande) |
| `hint` | string | Explicación de lo que se comprimió |

### Matices

- Si el resultado comprimido aún excede `max_response_size`, se guarda en un nuevo archivo y se devuelve un `FileReference`
- Valores predeterminados: `arrayHead=3`, `arrayTail=2`, `stringLen=80`
- Devuelve `validation_failed` para ruta inválida, JSONPath inválido o archivo no JSON
- Devuelve `not_found` si el archivo no existe o JSONPath no coincide

---

## response_slice

### Propósito

Extraer un fragmento específico de un archivo de respuesta JSON guardado por ruta JSON lógica o por rango de líneas. A diferencia de `response_compress`, esto le da los datos sin procesar y sin modificar.

### Cuándo usarlo

- Cuando necesita un elemento o valor específico de una respuesta grande
- Cuando `response_compress` no le da suficiente detalle
- Cuando desea navegar a través de una respuesta paso a paso

### Cómo funciona

Lee el archivo de respuesta guardado y extrae un fragmento por ruta JSON (por ejemplo, `data.3.name`) o por rango de líneas (por ejemplo, `120-240`). Devuelve sugerencias de navegación para recorrer arreglos y objetos.

### Parámetros

| Parámetro | Tipo | Requerido | Descripción |
|-----------|------|-----------|-------------|
| `path` | string | Sí | Ruta absoluta de `fileRef.path` |
| `jsonPath` | string | No | Ruta lógica al valor (por ejemplo, `data.3.name`) |
| `line` | int | No | Número de línea (basado en 1) para centrar el fragmento |
| `range` | string | No | Rango de líneas como `inicio-fin` (por ejemplo, `120-240`) |
| `around` | int | No | Líneas a incluir alrededor de `line` (predeterminado: 20) |

### Respuesta

```json
{
  "slice": {
    "lines": [120, 130],
    "fragment": "{\n  \"id\": 1,\n  \"name\": \"Rex\"\n}",
    "value": {
      "id": 1,
      "name": "Rex"
    },
    "jsonPath": "data.0",
    "context": "object",
    "isComplete": true,
    "nextLine": 131,
    "prevLine": 119,
    "nextPath": "data.1",
    "prevPath": null
  }
}
```

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `lines` | array | Rango de líneas basado en 1 [inicio, fin] |
| `fragment` | string | Texto JSON sin procesar (cuando es lo suficientemente pequeño) |
| `value` | any | Valor JSON extraído |
| `jsonPath` | string | La ruta JSON utilizada |
| `context` | string | "object", "array" o "value" |
| `isComplete` | bool | Verdadero cuando el valor es un fragmento JSON válido |
| `nextLine` | int | Siguiente línea sugerida para navegación basada en líneas |
| `prevLine` | int | Línea anterior sugerida |
| `nextPath` | string | Siguiente ruta JSON sugerida para navegación de arreglos |
| `prevPath` | string | Ruta JSON anterior sugerida |

### Matices

- **Prefiera `jsonPath` sobre números de línea** — las rutas JSON son estables y descriptivas, los números de línea cambian si el archivo se regenera
- Si el fragmento extraído excede `max_response_size`, se guarda en un nuevo archivo y se devuelve un `FileReference`
- El valor predeterminado de `around` es 20 líneas
- La respuesta incluye `nextPath`/`prevPath` para recorrer arreglos y `nextLine`/`prevLine` para navegación basada en líneas
- Devuelve `validation_failed` para ruta inválida, JSONPath inválido, línea/rango inválido o archivo no JSON
- Devuelve `not_found` si el archivo no existe o JSONPath no coincide
