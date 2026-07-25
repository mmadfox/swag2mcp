# Gestión del Tamaño de Respuesta

## Descripción General

Las respuestas de API pueden ser muy grandes — a veces demasiado grandes para caber en la ventana de contexto del LLM. swag2mcp gestiona automáticamente los tamaños de respuesta guardando las respuestas demasiado grandes en disco y proporcionando herramientas para explorarlas.

## Cómo funciona

1. **Usted llama a `invoke`** — swag2mcp realiza la solicitud a la API
2. **Si la respuesta es pequeña** (dentro del límite) — se devuelve en línea al LLM
3. **Si la respuesta es demasiado grande** (excede el límite) — se guarda en `{workspace}/responses/` como un archivo JSON. El LLM recibe una referencia de archivo en lugar de la respuesta completa

### Ejemplo: respuesta pequeña (en línea)

```json
{
  "statusCode": 200,
  "body": {
    "id": 1,
    "name": "Rex",
    "status": "available"
  }
}
```

### Ejemplo: respuesta grande (referencia de archivo)

```json
{
  "statusCode": 200,
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

## Configuración

```yaml
http_client:
  max_response_size: 1048576  # 1 MB en bytes
```

### max_response_size

- **Tipo:** `int` (bytes)
- **Valor predeterminado:** `1048576` (1 MB)
- **Rango:** 256 a 10,485,760 bytes (10 MB)
- **Efecto:** Las respuestas más grandes que esto se guardan en disco en lugar de devolverse en línea
- **Cuándo aumentar:** APIs que devuelven conjuntos de datos grandes (informes, registros, análisis)
- **Cuándo disminuir:** Ventana de contexto LLM limitada, o cuando prefiere acceso basado en archivos

## Trabajar con respuestas grandes

Cuando `invoke` devuelve un `fileRef`, use estas tres herramientas para explorar los datos:

### 1. response_outline — entender la estructura

Obtenga un resumen estructural de la respuesta: claves, tipos, longitudes de arreglos y sugerencias de navegación.

```json
→ response_outline(path: "/path/to/file.json")
← {
    "type": "object",
    "size": 1572864,
    "keys": ["data", "meta"],
    "itemCount": 500,
    "compressionHints": [
      "response_compress(path, 'first_of_array', 'data')",
      "response_compress(path, 'sample_array', 'data', arrayHead=3, arrayTail=2)"
    ]
  }
```

### 2. response_compress — obtener una versión más pequeña

Comprima los datos para que quepan en línea. Múltiples modos de compresión le permiten elegir el equilibrio adecuado.

| Modo | Descripción | Mejor para |
|------|-------------|------------|
| `first_of_array` | Conservar solo el primer elemento de un arreglo | Cuando todos los elementos tienen la misma estructura |
| `sample_array` | Conservar cabeza (3) y cola (2) de un arreglo | Cuando necesita ver el rango de valores |
| `truncate_strings` | Acortar cada cadena a N caracteres | Cuando las cadenas son muy largas |
| `keys_only` | Reemplazar valores con sus nombres de tipo | Cuando solo necesita la estructura |
| `select_keys` | Conservar solo las claves especificadas | Cuando necesita campos específicos |

```json
→ response_compress(path: "/path/to/file.json", mode: "first_of_array", jsonPath: "data")
← {
    "body": [{ "id": 1, "name": "Rex" }],
    "hint": "Arreglo comprimido de 500 a 1 elemento usando el modo first_of_array"
  }
```

### 3. response_slice — extraer un fragmento específico

Obtenga un elemento o valor específico por ruta JSON o rango de líneas.

```json
→ response_slice(path: "/path/to/file.json", jsonPath: "data.0")
← {
    "slice": {
      "value": { "id": 1, "name": "Rex" },
      "jsonPath": "data.0",
      "nextPath": "data.1",
      "prevPath": null
    }
  }
```

## Flujo de trabajo completo

```
1. invoke(endpoint) → fileRef (la respuesta es de 1.5 MB)
2. response_outline(path) → estructura: { data: Array(500) }
3. response_compress(path, mode: "first_of_array", jsonPath: "data") → primer elemento
4. response_slice(path, jsonPath: "data.0") → detalles completos del primer elemento
5. response_slice(path, jsonPath: "data.1") → segundo elemento
```

## Limpieza automática

Cuando se inicia el servidor MCP (`swag2mcp mcp`), los archivos de respuesta con más de 48 horas se eliminan automáticamente. También puede limpiarlos manualmente:

```bash
swag2mcp clean
```

## Notas importantes

- **El límite está en bytes** — `1048576` = 1 MB, `2097152` = 2 MB, etc.
- **Las referencias de archivo incluyen un comando de apertura** — en macOS es `open`, en Linux es `xdg-open`
- **Los archivos de respuesta se nombran con sufijos aleatorios** — sin conflictos entre llamadas concurrentes
- **El directorio de respuestas se crea automáticamente** — no se necesita configuración manual
