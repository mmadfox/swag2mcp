# Exportación e Importación

## Descripción General

swag2mcp admite el ciclo completo del espacio de trabajo mediante archivos ZIP. Puede exportar todo su espacio de trabajo (configuración, archivos de especificación, scripts de autenticación) a un archivo ZIP y restaurarlo en otra máquina.

## Exportación

Crea una copia de seguridad ZIP portátil de su espacio de trabajo.

```bash
# Exportar al archivo predeterminado (swag2mcp-backup-&lt;timestamp&gt;.zip)
swag2mcp export

# Exportar con ruta personalizada
swag2mcp export ~/backups/swag2mcp-backup.zip

# Exportar solo especificaciones específicas
swag2mcp export --spec meteo
swag2mcp export --spec meteo,store
```

### Qué se incluye en la exportación

| Elemento | Descripción |
|----------|-------------|
| `swag2mcp.yaml` | Archivo de configuración |
| `specs/` | Todos los archivos de especificación (OpenAPI/Swagger/Postman) |
| `auth_scripts/` | Scripts de autenticación |
| `swag2mcp.meta` | Metadatos (información de versión para compatibilidad) |

La caché y las respuestas **no** se exportan — son transitorias y estarían obsoletas al restaurar.

### Nombre de archivo predeterminado

Si no especifica una ruta de salida, el archivo se guarda como `swag2mcp-backup-<YYYY-MM-DD-HHMMSS>.zip` en el directorio actual (marca de tiempo UTC).

## Importación

Restaura un espacio de trabajo desde una copia de seguridad ZIP o importa archivos de especificación.

### Restaurar desde ZIP

```bash
# Restaurar espacio de trabajo completo
swag2mcp import --from-zip /path/to/backup.zip

# Restaurar con sobrescritura
swag2mcp import --from-zip /path/to/backup.zip --force
```

El ZIP debe ser creado por `swag2mcp export` — los archivos ZIP arbitrarios no funcionarán.

### Importar un solo archivo de especificación

Descargue un archivo de especificación y guárdelo en `specs/`:

```bash
swag2mcp import https://example.com/spec.yaml example-api.yaml
swag2mcp import /path/to/workspace https://example.com/spec.yaml example-api.yaml
```

### Importación masiva desde configuración existente (`--spec`)

Descargue todos los archivos de especificación de colecciones para los dominios especificados desde sus URL `location` configuradas, guárdelos en `specs/` y actualice la configuración para que apunte a las copias locales:

```bash
swag2mcp import --spec meteo
swag2mcp import /path/to/workspace --spec meteo,github
```

Esto hace que el espacio de trabajo sea autónomo — no se necesitan URL de especificación remotas después de la importación.

## Casos de uso

### Copia de seguridad

```bash
swag2mcp export swag2mcp-$(date +%Y-%m-%d).zip
```

### Transferencia a otra máquina

```bash
# En la máquina antigua
swag2mcp export swag2mcp.zip

# Copie el ZIP a la máquina nueva, luego:
swag2mcp import --from-zip swag2mcp.zip
```

### Compartir configuración

```bash
swag2mcp init
swag2mcp export template.zip
# Comparta template.zip con un colega
```

## Verificación posterior a la exportación

Siempre verifique que el archivo ZIP se haya creado:

```bash
ls -la swag2mcp-backup-*.zip
```

## Notas importantes

- **La salida debe ser una ruta de archivo que termine en `.zip`** — no pase un directorio
- **La caché y las respuestas están excluidas** — solo se conservan la configuración, las especificaciones y los scripts de autenticación
- **El ZIP es autónomo** — puede restaurarse en cualquier máquina con swag2mcp instalado
- **Filtro de especificaciones** — use `--spec` para exportar o importar solo especificaciones específicas
