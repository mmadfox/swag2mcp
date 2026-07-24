# export

## Propósito

Crear una copia de seguridad ZIP portátil del espacio de trabajo. El archivo contiene el archivo de configuración, todos los archivos de especificación y los scripts de autenticación — todo lo necesario para restaurar el espacio de trabajo en otra máquina.

## Cuándo usarlo

- Desea hacer una copia de seguridad de su espacio de trabajo antes de hacer cambios
- Está migrando swag2mcp a otra máquina
- Desea compartir su configuración de API con un colega
- Está preparando un entorno reproducible

## Sintaxis

```bash
swag2mcp export [path] [output] [flags]
```

## Argumentos

| Argumento | Posición | Requerido | Descripción |
|-----------|----------|-----------|-------------|
| `path` | 1 | No | Directorio del espacio de trabajo. Si se omite, se resuelve mediante reglas de resolución de ruta. |
| `output` | 2 | No | Ruta completa para el archivo ZIP de salida. Si se omite, el valor predeterminado es `./swag2mcp-backup-<timestamp>.zip`. |

## Banderas

| Bandera | Abreviatura | Tipo | Valor predeterminado | Descripción |
|---------|-------------|------|---------------------|-------------|
| `--spec` | `-s` | `stringSlice` | `nil` | Exportar solo las especificaciones indicadas (separadas por comas) |

## Cómo funciona

### Exportación predeterminada

Crea un ZIP en el directorio actual con un nombre con marca de tiempo:

```bash
swag2mcp export
# Crea ./swag2mcp-backup-2026-07-22-143022.zip
```

### Ruta de salida personalizada

```bash
swag2mcp export /path/to/workspace /path/to/backup.zip
```

### Exportar especificaciones específicas

```bash
swag2mcp export --spec meteo
swag2mcp export --spec meteo,store
```

## Qué hay en el ZIP

| Entrada | Descripción |
|---------|-------------|
| `swag2mcp.meta` | Metadatos sobre la exportación |
| `swag2mcp.yaml` | Archivo de configuración |
| `specs/` | Todos los archivos de especificación (OpenAPI/Swagger/Postman) |
| `auth_scripts/` | Scripts de autenticación |
| `cache/` | Vacío (la caché no se exporta) |
| `responses/` | Vacío (las respuestas no se exportan) |

## Restaurar

Use `import` para restaurar desde una copia de seguridad:

```bash
swag2mcp import --from-zip /path/to/backup.zip
```

## Verificación posterior al comando

Siempre verifique que el archivo ZIP se haya creado:

```bash
ls -la swag2mcp-backup-*.zip
# o para una ruta de salida personalizada:
ls -la /path/to/backup.zip
```

## Matices

- **La salida debe ser una ruta de archivo:** El argumento `[output]` debe ser una ruta de archivo completa que termine en `.zip`. **No** pase un directorio — el comando no creará un ZIP si se le da una ruta de directorio.
- **Nombre de archivo predeterminado:** `swag2mcp-backup-<YYYY-MM-DD-HHMMSS>.zip` usando marca de tiempo UTC.
- **Filtro `--spec`:** Cuando se establece, solo se incluyen las especificaciones indicadas. Otras especificaciones se excluyen del archivo.
- **No requiere configuración:** `export` funciona incluso sin un archivo de configuración válido. Exporta lo que exista en el espacio de trabajo.
- **La caché y las respuestas están excluidas:** Estos son datos transitorios que estarían obsoletos al restaurar. Solo se conservan la configuración, las especificaciones y los scripts de autenticación.
