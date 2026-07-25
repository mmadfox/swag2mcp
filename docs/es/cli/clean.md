# clean

## Propósito

Eliminar las especificaciones remotas en caché y las respuestas de invocación de API guardadas. Esto libera espacio en disco y fuerza una descarga nueva de los archivos de especificación en el próximo inicio de `update` o `mcp`.

## Cuándo usarlo

- Los archivos de especificación han cambiado en el servidor remoto y desea forzar una actualización
- Desea liberar espacio en disco
- Está solucionando problemas de caché obsoleta
- Antes de ejecutar `update` para asegurar una recarga limpia

## Sintaxis

```bash
swag2mcp clean [path]
```

## Argumentos

| Argumento | Posición | Requerido | Descripción |
|-----------|----------|-----------|-------------|
| `path` | 1 | No | Directorio del espacio de trabajo. Si se omite, se resuelve mediante reglas de resolución de ruta. |

## Banderas

Ninguna.

## Cómo funciona

```bash
swag2mcp clean
swag2mcp clean ./my-workspace
```

## Qué se limpia

| Directorio | Contenido | Por qué |
|------------|-----------|---------|
| `cache/` | Archivos de especificación remotos descargados | Fuerza la redescarga en el próximo acceso |
| `responses/` | Respuestas de invocación de API guardadas | Libera espacio en disco |

## Qué se conserva

| Directorio | Contenido | Por qué |
|------------|-----------|---------|
| `specs/` | Archivos de especificación locales | Estos son sus archivos fuente, no caché |
| `auth_scripts/` | Scripts de autenticación | Estos son creados por el usuario, no caché |

## Limpieza de scripts huérfanos

Después de limpiar, `clean` también elimina los scripts de autenticación para especificaciones que ya no existen en la configuración. Esto evita que se acumulen scripts obsoletos.

## Limpieza automática

Cuando se inicia el servidor MCP (`swag2mcp mcp`), las respuestas con más de 48 horas se eliminan automáticamente. Normalmente no necesita ejecutar `clean` manualmente para el mantenimiento rutinario.

## Verificación posterior al comando

```bash
ls ~/.swag2mcp/cache
# Debería estar vacío (el directorio existe pero no tiene archivos)
```

## Matices

- **No requiere configuración:** `clean` funciona incluso sin un archivo de configuración válido. Simplemente elimina los directorios de caché y respuestas.
- **La limpieza de huérfanos es de mejor esfuerzo:** Si el archivo de configuración está corrupto o es ilegible, la limpieza de scripts huérfanos se omite (no es fatal).
- **Los directorios se conservan:** Los directorios `cache/` y `responses/` se mantienen — solo se elimina su contenido.
