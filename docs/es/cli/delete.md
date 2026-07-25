# delete

## Propósito

Eliminar una **especificación** (servicio de API) o **colección** (archivo de especificación) de la configuración. Es la operación inversa de `add`.

## Cuándo usarlo

- Una API ya no es necesaria
- Desea eliminar un archivo de especificación específico de una especificación
- Está limpiando su espacio de trabajo

## Sintaxis

```bash
swag2mcp delete spec [path]
swag2mcp delete collection [path]
```

## Argumentos

| Argumento | Posición | Requerido | Descripción |
|-----------|----------|-----------|-------------|
| `path` | 1 | No | Directorio del espacio de trabajo. Si se omite, se resuelve mediante reglas de resolución de ruta. |

## Banderas

Ninguna. Ambos subcomandos son puramente interactivos.

## Cómo funciona

### Eliminar una especificación

Le solicita que seleccione una especificación de una lista, luego pide confirmación antes de eliminar.

```bash
swag2mcp delete spec
```

### Eliminar una colección

Le solicita que seleccione una especificación, luego una colección dentro de esa especificación, luego pide confirmación.

```bash
swag2mcp delete collection
```

## Encontrar IDs

Los mensajes interactivos muestran nombres legibles por humanos, no IDs. Si necesita IDs como referencia:

```bash
# Listar todas las especificaciones con sus IDs
swag2mcp ls

# Listar colecciones para una especificación específica
swag2mcp ls --tags
```

## Verificación posterior al comando

```bash
swag2mcp ls [path]
# La especificación o colección eliminada ya no debería aparecer
```

## Matices

- **Se requiere TTY:** Ambos comandos requieren una terminal interactiva. **No** funcionarán en pipelines CI/CD, tareas cron o scripts no interactivos.
- **No hay `--force` o `--yes`:** No hay forma de omitir el mensaje de confirmación. Esto es intencional para evitar eliminaciones accidentales.
- **Auto-inicio:** Si no existe un archivo de configuración, `delete` ejecuta automáticamente el asistente de inicio primero.
- **Sin modo YAML:** A diferencia de `add`, no hay una bandera `--yaml`. La eliminación siempre es interactiva.
