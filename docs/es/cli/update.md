# update

## Propósito

Revalidar la configuración, limpiar la caché y volver a descargar todos los archivos de especificación. Esto es una **actualización completa** del espacio de trabajo — asegura que todas las especificaciones en caché estén actualizadas y el índice se reconstruya.

## Cuándo usarlo

- Los archivos de especificación remotos han cambiado y desea la versión más reciente
- Después de editar `swag2mcp.yaml` para agregar o cambiar ubicaciones de especificaciones
- Al solucionar problemas de caché obsoleta o corrupta
- Antes de ejecutar `mcp` para asegurarse de que todo esté actualizado

## Sintaxis

```bash
swag2mcp update [path]
```

## Argumentos

| Argumento | Posición | Requerido | Descripción |
|-----------|----------|-----------|-------------|
| `path` | 1 | No | Directorio del espacio de trabajo. Si se omite, se resuelve mediante reglas de resolución de ruta. |

## Banderas

Ninguna.

## Cómo funciona

El comando `update` ejecuta un pipeline de operaciones:

1. **Cargar configuración** — lee `swag2mcp.yaml` del espacio de trabajo
2. **Validar** — ejecuta las mismas comprobaciones que `validate` (sintaxis YAML, estructura, alcance del archivo de especificación, formato, autenticación, cliente HTTP)
3. **Limpiar** — elimina todo el contenido de `cache/` y `responses/`
4. **Recargar en caché** — descarga todos los archivos de especificación remotos y copia los archivos de especificación locales a la caché
5. **Reindexar** — reconstruye el índice de búsqueda de texto completo para todos los endpoints
6. **Scripts de autenticación** — crea scripts de autenticación de plantilla para especificaciones que usan `ScriptAuth`
7. **Limpieza de huérfanos** — elimina scripts de autenticación para especificaciones que ya no existen

```bash
swag2mcp update
swag2mcp update ./my-workspace
```

## Qué sucede con las colecciones deshabilitadas

Las colecciones con `disable: true` se omiten por completo — no se almacenan en caché ni se indexan.

## Verificación posterior al comando

```bash
swag2mcp ls [path]
# Todas las especificaciones deberían seguir listadas y accesibles
```

## Matices

- **Sin auto-inicio:** Si el archivo de configuración no existe, `update` devuelve un error: `"configuration not found at <path>"`. Ejecute `init` primero.
- **Dependencia de red:** Todas las URL de especificaciones remotas deben ser accesibles. Si alguna descarga falla, toda la actualización falla con un mensaje de error claro.
- **Creación de scripts de autenticación:** Si una especificación usa `ScriptAuth` y el script de plantilla no existe, `update` lo crea. Si la creación falla, la actualización falla.
- **`update` vs `clean`:** `clean` solo elimina la caché. `update` elimina la caché **y** vuelve a descargar todo. Use `clean` cuando solo quiera liberar espacio; use `update` cuando quiera actualizar.
