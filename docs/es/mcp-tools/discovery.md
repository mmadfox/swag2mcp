# Herramientas de Descubrimiento

Las herramientas de descubrimiento permiten al LLM navegar por la jerarquía de especificaciones: encontrar todas las especificaciones, profundizar en una especificación para ver sus colecciones y explorar etiquetas dentro de una colección. Comience con `spec_list` para ver qué APIs están disponibles, luego use IDs para profundizar más.

---

## spec_list

### Propósito

Listar todas las especificaciones de API registradas en el espacio de trabajo. Este es el punto de partida para cualquier sesión — el LLM lo llama primero para descubrir qué APIs están disponibles.

### Cuándo usarlo

- Al inicio de una sesión para ver qué APIs están configuradas
- Después de agregar o eliminar especificaciones para actualizar la lista
- Cuando necesita un ID de especificación para otras herramientas

### Cómo funciona

Devuelve una lista de todas las especificaciones con su ID único y nombre de dominio. No se necesitan parámetros.

### Parámetros

Ninguno.

### Respuesta

```json
{
  "specs": [
    {
      "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
      "domain": "meteo"
    },
    {
      "id": "b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7",
      "domain": "dadjoke"
    }
  ]
}
```

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `id` | string | Hash MD5 de 32 caracteres, identificador único de la especificación |
| `domain` | string | Nombre de dominio de la especificación (por ejemplo, "meteo", "dadjoke") |

### Matices

- Devuelve solo `id` y `domain` — para detalles completos (colecciones, etiquetas), use `spec_by_id`
- Todos los IDs son cadenas hexadecimales MD5 de 32 caracteres (`^[0-9a-f]{32}$`)
- Si no hay especificaciones configuradas, devuelve un arreglo vacío

---

## spec_by_id

### Propósito

Obtener información detallada sobre una especificación específica: su dominio, todas las colecciones y sus estadísticas (recuento de etiquetas, recuento de métodos).

### Cuándo usarlo

- Después de `spec_list` para ver las colecciones dentro de una especificación
- Cuando necesita IDs de colección para navegación adicional

### Cómo funciona

Toma un ID de especificación y devuelve los metadatos de la especificación más todas sus colecciones con recuentos.

### Parámetros

| Parámetro | Tipo | Requerido | Descripción |
|-----------|------|-----------|-------------|
| `id` | string | Sí | Hash MD5 de 32 caracteres de la especificación |

### Respuesta

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collections": [
    {
      "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
      "title": "Weather Forecast",
      "llmTitle": "Forecast API",
      "countTags": 3,
      "countMethods": 12
    }
  ]
}
```

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `spec.id` | string | Identificador de la especificación |
| `spec.domain` | string | Nombre de dominio de la especificación |
| `collections[].id` | string | Identificador de la colección |
| `collections[].title` | string | Título legible por humanos |
| `collections[].llmTitle` | string | Título amigable para LLM (opcional) |
| `collections[].countTags` | int | Número de etiquetas en la colección |
| `collections[].countMethods` | int | Número de métodos HTTP en la colección |

### Matices

- Devuelve error `not_found` si el ID de la especificación no existe
- El `id` debe ser una cadena hexadecimal MD5 válida de 32 caracteres

---

## collection_by_spec

### Propósito

Listar todas las colecciones dentro de una especificación específica. Similar a `spec_by_id` pero devuelve solo la lista de colecciones sin metadatos adicionales de la especificación.

### Cuándo usarlo

- Cuando ya tiene el ID de la especificación y solo necesita la lista de colecciones
- Como una alternativa más ligera a `spec_by_id`

### Parámetros

| Parámetro | Tipo | Requerido | Descripción |
|-----------|------|-----------|-------------|
| `specId` | string | Sí | Hash MD5 de 32 caracteres de la especificación |

### Respuesta

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collections": [
    {
      "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
      "title": "Weather Forecast",
      "llmTitle": "Forecast API",
      "countTags": 3,
      "countMethods": 12
    }
  ]
}
```

### Matices

- Devuelve `not_found` si la especificación no existe
- Mismos datos que `spec_by_id` pero sin el envoltorio adicional de la especificación

---

## collection_by_id

### Propósito

Obtener información detallada sobre una colección específica: sus metadatos, la especificación padre y todas las etiquetas dentro de la colección.

### Cuándo usarlo

- Después de `collection_by_spec` para ver las etiquetas dentro de una colección
- Cuando necesita IDs de etiqueta para `tag_by_id` o `endpoint_by_tag`

### Parámetros

| Parámetro | Tipo | Requerido | Descripción |
|-----------|------|-----------|-------------|
| `id` | string | Sí | Hash MD5 de 32 caracteres de la colección |

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
  "tags": [
    {
      "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "title": "forecast",
      "countMethods": 5
    },
    {
      "id": "e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7",
      "title": "current",
      "countMethods": 7
    }
  ]
}
```

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `spec` | object | Especificación padre (id, domain) |
| `collection` | object | Metadatos de la colección (id, title, countMethods) |
| `tags[]` | array | Lista de etiquetas con id, title, countMethods |

### Matices

- Devuelve `not_found` si el ID de la colección no existe
- Las etiquetas se devuelven con sus IDs — use `endpoint_by_tag(tagId)` para ver los endpoints reales

---

## tag_by_spec

### Propósito

Listar todas las etiquetas en una especificación completa, abarcando todas las colecciones. Útil para obtener una vista general de todas las etiquetas disponibles.

### Cuándo usarlo

- Cuando desea ver todas las etiquetas en una especificación sin profundizar en cada colección
- Cuando no sabe qué colección contiene la etiqueta que necesita

### Parámetros

| Parámetro | Tipo | Requerido | Descripción |
|-----------|------|-----------|-------------|
| `specId` | string | Sí | Hash MD5 de 32 caracteres de la especificación |

### Respuesta

```json
{
  "tags": [
    {
      "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "title": "forecast",
      "countMethods": 5
    },
    {
      "id": "e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7",
      "title": "current",
      "countMethods": 7
    }
  ]
}
```

### Matices

- Devuelve `not_found` si la especificación no existe
- Las etiquetas se agregan de todas las colecciones en la especificación

---

## tag_by_collection

### Propósito

Listar todas las etiquetas dentro de una colección específica. A diferencia de `tag_by_spec`, también devuelve la especificación padre y los metadatos de la colección.

### Cuándo usarlo

- Después de `collection_by_id` para confirmar la lista de etiquetas
- Cuando necesita el contexto completo (especificación + colección + etiquetas)

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
  "tags": [
    {
      "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "title": "forecast",
      "countMethods": 5
    }
  ]
}
```

### Matices

- Devuelve `not_found` si la colección no existe
- Mismos datos de etiqueta que `tag_by_spec` pero limitados a una colección

---

## tag_by_id

### Propósito

Obtener información sobre una sola etiqueta: su ID, título y cuántos métodos contiene. Esto le informa sobre la etiqueta en sí — para ver los endpoints reales, use `endpoint_by_tag`.

### Cuándo usarlo

- Cuando tiene un ID de etiqueta y desea confirmar su nombre y tamaño
- Antes de llamar a `endpoint_by_tag` para entender cuántos endpoints esperar

### Parámetros

| Parámetro | Tipo | Requerido | Descripción |
|-----------|------|-----------|-------------|
| `id` | string | Sí | Hash MD5 de 32 caracteres de la etiqueta |

### Respuesta

```json
{
  "tag": {
    "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
    "title": "forecast",
    "countMethods": 5
  }
}
```

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `tag.id` | string | Identificador de la etiqueta |
| `tag.title` | string | Nombre de etiqueta legible por humanos |
| `tag.countMethods` | int | Número de métodos HTTP en esta etiqueta |

### Matices

- Devuelve `not_found` si la etiqueta no existe
- Esta herramienta devuelve solo metadatos de la etiqueta — use `endpoint_by_tag` para obtener la lista real de endpoints
