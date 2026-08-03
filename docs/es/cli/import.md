# import

## Propósito

Importar archivos de especificación al directorio `specs/` del espacio de trabajo para uso local, o restaurar un espacio de trabajo completo desde una copia de seguridad ZIP. Tres modos cubren diferentes escenarios: agregar una sola especificación, importación masiva desde configuración existente o restaurar un espacio de trabajo completo.

## Cuándo usarlo

- Tiene una URL o archivo de especificación y desea guardarlo localmente en el espacio de trabajo
- Desea descargar todos los archivos de especificación de colecciones desde la configuración y hacer que el espacio de trabajo sea autónomo
- Necesita restaurar un espacio de trabajo desde una copia de seguridad ZIP creada por `export`
- Está migrando swag2mcp a otra máquina

## Sintaxis

```bash
swag2mcp import [path] [source] [name] [flags]
```

## Argumentos

| Argumento | Posición | Requerido | Descripción |
|-----------|----------|-----------|-------------|
| `path` | 1 | No | Directorio del espacio de trabajo. Si se omite, se resuelve mediante reglas de resolución de ruta. |
| `source` | 2 | Varía | URL o ruta local a un archivo de especificación, o ruta a un archivo ZIP |
| `name` | 3 | Varía | Nombre de archivo para guardar (ej. `example-api.yaml`). Se deriva de la URL si se omite. |

## Banderas

| Bandera | Abreviatura | Tipo | Valor predeterminado | Descripción |
|---------|-------------|------|---------------------|-------------|
| `--spec` | `-s` | `string` | `""` | Descargar archivos de especificación de colecciones desde la configuración. Sin valor para todas las specs, o especificar dominios como `--spec meteo,github` |
| `--force` | `-f` | `bool` | `false` | Sobrescribir archivos de especificación existentes sin error |
| `--from-zip` | | `string` | `""` | Restaurar espacio de trabajo desde un ZIP de copia de seguridad de swag2mcp |

## Cómo funciona

### Modo 1 — Importación única desde URL o archivo

Descargue un archivo de especificación y guárdelo en `specs/`:

```bash
swag2mcp import https://example.com/spec.yaml example-api.yaml
swag2mcp import /path/to/workspace https://example.com/spec.yaml example-api.yaml
swag2mcp import ./local-spec.yaml example-api.yaml
```

Si se omite `name`, se deriva del nombre del archivo en la URL:
```bash
swag2mcp import https://example.com/specs/petstore.yaml
# → guardado como petstore.yaml
```

Sobrescribir un archivo existente con `--force`:
```bash
swag2mcp import --force https://example.com/spec.yaml example-api.yaml
```

Después de la importación, la salida muestra la ruta del espacio de trabajo, el archivo guardado y una plantilla YAML para agregar a `swag2mcp.yaml`:

```
✅ Imported to /path/to/workspace
   specs/example-api.yaml

   Add to swag2mcp.yaml:
     specs:
       - domain: <your-domain>
         collections:
           - location: specs/example-api.yaml
```

### Modo 2 — Importación masiva desde configuración existente (`--spec`)

Descargue todos los archivos de especificación de colecciones para los dominios especificados desde sus URL `location` configuradas, guárdelos en `specs/` y actualice la configuración para que apunte a las copias locales:

```bash
swag2mcp import --spec                # todas las specs
swag2mcp import --spec meteo           # spec específica
swag2mcp import --spec meteo,github    # múltiples specs
swag2mcp import /path/to/workspace --spec meteo
```

Si un dominio especificado no existe en la configuración, el comando devuelve un error:
```
Error: import_no_match
  Spec "nonexistent" not found in config.
```

Esto hace que el espacio de trabajo sea autónomo — no se necesitan URL de especificación remotas después de la importación.

### Modo 3 — Restaurar desde copia de seguridad ZIP

Restaurar un espacio de trabajo completo desde un archivo ZIP creado por `swag2mcp export`:

```bash
swag2mcp import --from-zip /path/to/backup.zip
swag2mcp import /path/to/workspace /path/to/backup.zip
```

> **El ZIP debe ser creado por `swag2mcp export`.** Los archivos ZIP arbitrarios no funcionarán — el archivo tiene una estructura interna específica (`swag2mcp.yaml`, `specs/`, `auth_scripts/`).

## Verificación posterior al comando

```bash
# Importación única o masiva
swag2mcp ls [path]
# La nueva especificación debería aparecer en la lista

# Restauración ZIP
swag2mcp ls [path]
# Todas las especificaciones de la copia de seguridad deberían aparecer
```

## Matices

- **El modo masivo requiere configuración:** Al usar `--spec`, el archivo de configuración debe existir. Ejecute `init` primero si es necesario.
- **La importación única crea el espacio de trabajo:** Si el espacio de trabajo no existe, se crea automáticamente.
- **Detección de ZIP:** Un argumento posicional que termina en `.zip` se trata como un origen ZIP. La bandera `--from-zip` tiene prioridad sobre la detección posicional.
- **Cliente HTTP:** La configuración global del cliente HTTP se aplica durante la importación (tiempo de espera, proxy, encabezados, etc.).
