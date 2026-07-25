# Búsqueda de Texto Completo

## Descripción General

swag2mcp incluye un motor de búsqueda de texto completo incorporado (bluge) que indexa todos los endpoints de todas las especificaciones. El LLM puede buscar endpoints por método, ruta, resumen o etiqueta — incluso sin conocer el ID del endpoint.

## Cómo funciona la indexación

Cuando se agrega o actualiza una especificación, cada endpoint se indexa. Los siguientes campos se pueden buscar:

| Campo | Descripción | Ejemplo |
|-------|-------------|---------|
| `method` | Método HTTP | `GET`, `POST`, `PUT` |
| `path` | Ruta del endpoint de API | `/api/v1/users/{id}` |
| `summary` | Resumen de OpenAPI | "Find pet by ID" |
| `tag` | Categoría del endpoint | "pets", "users" |
| `_all` | Todos los campos combinados | method + path + tag + summary |

El índice se reconstruye en cada inicio del servidor MCP. Se almacena en memoria para búsquedas rápidas.

## Sintaxis de consulta

La búsqueda admite una sintaxis de consulta enriquecida para un filtrado preciso:

| Ejemplo | Descripción |
|---------|-------------|
| `pet` | Búsqueda de texto simple en todos los campos |
| `method:GET` | Encontrar todos los endpoints GET |
| `tag:pets` | Encontrar endpoints en la etiqueta "pets" |
| `path:"/api/v1/users"` | Coincidencia exacta de ruta |
| `+method:POST +tag:pet` | Debe coincidir con ambas condiciones |
| `-method:DELETE` | Excluir métodos DELETE |
| `create~` | Búsqueda difusa (tolerante a errores tipográficos) |
| `cr*` | Búsqueda con comodín |
| `"find pet"` | Búsqueda de frase |
| `+summary:pet -method:DELETE` | Incluir "pet" en resumen, excluir DELETE |

### Búsqueda por campo específico

Puede buscar dentro de campos específicos usando la sintaxis `field:value`:

```
method:GET
tag:pets
path:"/pet/findByStatus"
summary:"find pet by status"
```

### Operadores booleanos

- `+` — el término debe coincidir (AND)
- `-` — el término no debe coincidir (NOT)
- Espacio entre términos — OR (cualquier término puede coincidir)

### Búsqueda difusa y comodines

- `term~` — búsqueda difusa (coincide con palabras similares, maneja errores tipográficos)
- `te*` — comodín (coincide con cualquier carácter)
- `te?t` — comodín de un solo carácter

## Ejemplos

```
# Encontrar todas las solicitudes GET
method:GET

# Encontrar solicitudes POST en la etiqueta pet
+method:POST +tag:pet

# Encontrar endpoints por ruta exacta
path:"/pet/findByStatus"

# Encontrar por descripción
"find pet by status"

# Encontrar todo excepto DELETE
+summary:pet -method:DELETE

# Búsqueda difusa para "create" (maneja errores tipográficos)
create~
```

## Herramienta MCP

La herramienta MCP `search` expone el motor de búsqueda al LLM:

```
→ search(query: "find pet by status", limit: 5)
← GET /pet/findByStatus — Finds Pets by status
   GET /pet/{petId} — Find pet by ID
```

### Parámetros

| Parámetro | Requerido | Descripción |
|-----------|-----------|-------------|
| `query` | Sí | Consulta de búsqueda (admite sintaxis estructurada) |
| `limit` | Sí | Resultados máximos (1-50) |

## Notas importantes

- **El índice está en memoria** — se reconstruye cada vez que se inicia el servidor MCP. No hay un archivo de índice persistente.
- **Todos los campos están en minúsculas** — las búsquedas no distinguen entre mayúsculas y minúsculas
- **El límite máximo es 50** — no puede solicitar más de 50 resultados
- **La sintaxis de consulta inválida** devuelve un mensaje de error útil con ejemplos
- **El campo `_all`** combina method, path, tag y summary para búsquedas de texto simples
