# import

## Propósito

Importar archivos de especificación al espacio de trabajo o restaurar un espacio de trabajo completo desde una copia de seguridad ZIP. Tres modos cubren diferentes escenarios: agregar una sola especificación, importación masiva desde configuración existente o restaurar un espacio de trabajo completo.

## Cuándo usarlo

- Tiene una URL o archivo de especificación y desea agregarlo al espacio de trabajo
- Desea descargar todos los archivos de especificación referenciados en la configuración
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
| `name` | 3 | Varía | Nombre de dominio para la nueva especificación |

## Banderas

| Bandera | Abreviatura | Tipo | Valor predeterminado | Descripción |
|---------|-------------|------|---------------------|-------------|
| `--spec` | `-s` | `stringSlice` | `nil` | Importar colecciones de las especificaciones indicadas (separadas por comas) |
| `--from-zip` | | `string` | `""` | Restaurar espacio de trabajo desde un ZIP de copia de seguridad de swag2mcp |

## Cómo funciona

### Modo 1 — Importación única desde URL o archivo

Descargue un archivo de especificación y agréguelo al espacio de trabajo con un nombre de dominio:

```bash
swag2mcp import https://example.com/spec.yaml myspec
swag2mcp import /path/to/workspace https://example.com/spec.yaml myspec
swag2mcp import ./local-spec.yaml myspec
```

El archivo de especificación se guarda en `specs/` y la configuración se actualiza con la nueva entrada de especificación.

### Modo 2 — Importación masiva desde configuración existente

Descargue todas las colecciones para los dominios especificados desde sus URL configuradas:

```bash
swag2mcp import --spec meteo
swag2mcp import /path/to/workspace --spec meteo,store
```

El archivo de especificación de cada colección se descarga y se guarda en `specs/`. La configuración se actualiza para apuntar a las copias locales.

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
- **`--force`:** Disponible para la restauración ZIP para sobrescribir un espacio de trabajo existente.
- **Cliente HTTP:** La configuración global del cliente HTTP se aplica durante la importación (tiempo de espera, proxy, encabezados, etc.).
