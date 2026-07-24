# ls

## Propósito

Listar todas las **especificaciones** configuradas y sus **colecciones** en un formato legible por humanos. Esta es la forma principal de inspeccionar qué APIs están disponibles en su espacio de trabajo.

## Cuándo usarlo

- Desea ver qué APIs están configuradas
- Necesita encontrar un ID de especificación o colección
- Desea verificar cuántos endpoints tiene cada colección
- Desea filtrar especificaciones por etiquetas

## Sintaxis

```bash
swag2mcp ls [path] [flags]
```

## Argumentos

| Argumento | Posición | Requerido | Descripción |
|-----------|----------|-----------|-------------|
| `path` | 1 | No | Directorio del espacio de trabajo. Si se omite, se resuelve mediante reglas de resolución de ruta. |

## Banderas

| Bandera | Abreviatura | Tipo | Valor predeterminado | Descripción |
|---------|-------------|------|---------------------|-------------|
| `--tags` | `-t` | `string` | `""` | Filtrar especificaciones por etiquetas (separadas por comas) |

## Cómo funciona

### Listar todas las especificaciones

Muestra cada especificación con su dominio, colecciones y recuentos de endpoints:

```bash
swag2mcp ls
```

Ejemplo de salida:

```
Specifications:
  dadjoke (https://icanhazdadjoke.com)
    jokes (3 endpoints)
  meteo (https://meteo.swagger.io/v2)
    forecast (5 endpoints)
    current (8 endpoints)
  binance (https://api.binance.com)
    market-data (12 endpoints)
```

### Filtrar por etiquetas

Mostrar solo las especificaciones que tienen las etiquetas especificadas:

```bash
swag2mcp ls --tags=public
swag2mcp ls --tags=public,internal
```

## Verificación posterior al comando

Use `ls` después de `add`, `delete`, `update` o `import` para confirmar que el estado del espacio de trabajo coincida con sus expectativas.

## Matices

- **Auto-inicio:** Si no existe un archivo de configuración, `ls` ejecuta automáticamente el asistente de inicio primero.
- **Filtrado por etiquetas:** Las etiquetas están separadas por comas. Se muestran las especificaciones que coinciden con **cualquiera** de las etiquetas especificadas (lógica OR).
- **Formato de salida:** La salida es texto plano, no JSON. Para una salida legible por máquina, use `info`.
