# Espacio de Trabajo

El espacio de trabajo es el directorio donde swag2mcp almacena todos sus datos — configuración, especificaciones en caché, archivos de especificación locales, respuestas guardadas y scripts de autenticación.

## Estructura

```
~/.swag2mcp/                          # Raíz del espacio de trabajo (predeterminado)
├── swag2mcp.yaml                     # Archivo de configuración
├── cache/                            # Archivos de especificación remotos en caché
│   ├── a1b2c3d4e5f6...spec          # Contenido de especificación en caché
│   └── a1b2c3d4e5f6...meta          # Metadatos de la caché (JSON)
├── specs/                            # Archivos de especificación locales
│   └── my-api.yaml
├── responses/                        # Respuestas de API guardadas (respuestas grandes)
│   ├── meteo-get-forecast-abc123.json
│   └── response-fragment-def456.json
└── auth_scripts/                     # Scripts de autenticación
    ├── meteo.sh                      # Script de shell Unix
    └── meteo.bat                     # Script de batch Windows
```

## Ruta Predeterminada

- **Linux/macOS**: `~/.swag2mcp/`
- **Windows**: `%USERPROFILE%\.swag2mcp\`

## Ruta Personalizada

```bash
swag2mcp mcp /path/to/workspace
swag2mcp mcp ./my-workspace
```

## Directorios

### cache/

Almacena archivos de especificación remotos descargados. Cada archivo se almacena en caché con un hash SHA-256 de su URL como nombre de archivo:

- `{hash}.spec` — el contenido del archivo de especificación en caché
- `{hash}.meta` — metadatos JSON (URL de origen, tiempo de caché, TTL)

Cada archivo en caché tiene un TTL aleatorio entre 1 hora y 48 horas. La caché se verifica automáticamente en cada inicio — si existe una entrada válida (no expirada), se reutiliza sin descargar.

**Comandos:**
- `swag2mcp update` — limpia la caché y vuelve a descargar todas las especificaciones
- `swag2mcp clean` — limpia la caché y las respuestas

### specs/

Almacena archivos de especificación locales a los que las colecciones apuntan mediante `location: specs/{name}`. Los archivos aquí se usan directamente sin almacenamiento en caché.

Este directorio se llena mediante:
- `swag2mcp import <source> <name>` — descarga una especificación remota y la guarda aquí
- `swag2mcp export` — copia las especificaciones aquí en el ZIP de exportación
- Colocación manual — puede copiar archivos de especificación aquí usted mismo

### responses/

Almacena respuestas de API que exceden el límite de `max_response_size` (predeterminado 1 MB). Cuando el LLM invoca un endpoint y la respuesta es demasiado grande, swag2mcp la guarda aquí y devuelve una referencia de archivo en su lugar.

Convención de nomenclatura: `{domain}-{method}-{path_with_underscores}-{6char_hex}.json`

Las respuestas antiguas se limpian automáticamente después de 48 horas al iniciar el servidor MCP.

### auth_scripts/

Almacena scripts de autenticación para el tipo de autenticación `script`. Cada script se nombra según el dominio de la especificación.

#### Convención de Nomenclatura

| Plataforma | Nombre de archivo | Ejemplo |
|------------|-------------------|---------|
| Unix (Linux, macOS) | `{domain}.sh` | `meteo.sh` |
| Windows | `{domain}.bat` | `meteo.bat` |

El dominio no debe contener caracteres `/` o `\`.

#### Cómo Funcionan los Scripts

1. swag2mcp ejecuta el script con un tiempo de espera de 30 segundos
2. El script debe generar JSON válido en stdout
3. swag2mcp analiza el JSON y usa el token para las solicitudes de API

#### Formato de Salida Esperado

```json
{
  "token": "su-token-aquí",
  "expires_in": 3600
}
```

| Campo | Tipo | Requerido | Descripción |
|-------|------|-----------|-------------|
| `token` | string | ✅ | El token de autenticación |
| `access_token` | string | ❌ | Alternativa a `token` (se verifica primero) |
| `token_type` | string | ❌ | Tipo de token (por ejemplo, "Bearer") |
| `expires_in` | number | ❌ | Vida útil del token en segundos (predeterminado: 3600) |

#### Ejecución

| Plataforma | Comando |
|------------|---------|
| Unix | `sh {domain}.sh` |
| Windows | `cmd /c {domain}.bat` |

#### Almacenamiento en Caché del Token

El token se almacena en caché en memoria hasta que expira. En cada llamada a la API, swag2mcp verifica primero la caché — el script solo se ejecuta cuando el token en caché ha expirado.

#### Creación de Plantillas

Cuando configura `auth: { type: script, config: { domain: "myapi" } }`, swag2mcp crea un script de plantilla automáticamente:

**Unix (`auth_scripts/myapi.sh`):**
```bash
#!/bin/sh
echo '{"token": "your-token-here", "expires_in": 3600}'
```

**Windows (`auth_scripts/myapi.bat`):**
```bat
@echo off
echo {"token": "your-token-here", "expires_in": 3600}
```

Reemplace el token de marcador de posición con su lógica de autenticación real.

#### Limpieza de Huérfanos

Cuando elimina una especificación, su script de autenticación se vuelve huérfano. swag2mcp elimina automáticamente los scripts huérfanos en:
- `swag2mcp update`
- `swag2mcp clean`

## Comandos

### update

```bash
swag2mcp update [path]
```

Valida la configuración, limpia la caché y las respuestas, luego vuelve a descargar todos los archivos de especificación. También asegura que los scripts de autenticación existan y elimina los scripts huérfanos.

Use este comando después de:
- Agregar o eliminar colecciones
- Cambiar las ubicaciones de las colecciones
- Editar archivos de especificación que necesitan recarga en caché

### clean

```bash
swag2mcp clean [path]
```

Elimina todo el contenido de `cache/` y `responses/`, más los scripts de autenticación huérfanos. NO vuelve a almacenar en caché las especificaciones — use `update` para eso.

### validate

```bash
swag2mcp validate [path]
```

Valida la configuración incluyendo todas las ubicaciones de colecciones. Consulte [CLI: validate](../cli/validate.md).

## Exportación e Importación

```bash
# Exportar espacio de trabajo a ZIP (nombre predeterminado: swag2mcp-backup-{date}.zip)
swag2mcp export

# Exportar a una ruta específica
swag2mcp export /path/to/workspace /path/to/backup.zip

# Exportar solo especificaciones específicas
swag2mcp export --spec meteo

# Restaurar desde copia de seguridad
swag2mcp import --from-zip /path/to/backup.zip
swag2mcp import /path/to/workspace /path/to/backup.zip
```

La exportación incluye: `swag2mcp.yaml`, `specs/`, `auth_scripts/`. La caché y las respuestas están excluidas (son datos locales).

## .gitignore

Si su espacio de trabajo está dentro de un repositorio Git, agregue estas entradas a `.gitignore`:

```gitignore
# swag2mcp — solo datos locales
.swag2mcp/cache/
.swag2mcp/responses/
```

Los directorios `cache/` y `responses/` contienen datos locales específicos de la máquina que no deben ser confirmados. Todo lo demás (`swag2mcp.yaml`, `specs/`, `auth_scripts/`) debe estar en el repositorio para que la configuración se comparta en todo el equipo.
