# init

## Propósito

El comando `init` crea un **espacio de trabajo** — un directorio con un archivo de configuración `swag2mcp.yaml` y subdirectorios para caché, especificaciones, respuestas y scripts de autenticación. Este es el primer comando a ejecutar al configurar swag2mcp.

## Cuándo usarlo

- Está configurando swag2mcp por primera vez
- Desea crear un nuevo espacio de trabajo en un directorio específico
- Necesita reinicializar un espacio de trabajo corrupto o faltante

## Sintaxis

```bash
swag2mcp init [path] [flags]
```

## Argumentos

| Argumento | Posición | Requerido | Descripción |
|-----------|----------|-----------|-------------|
| `path` | 1 | No | Directorio del espacio de trabajo. Si se omite, el valor predeterminado es `~/.swag2mcp`. |

## Banderas

| Bandera | Abreviatura | Tipo | Valor predeterminado | Descripción |
|---------|-------------|------|---------------------|-------------|
| `--interactive` | `-i` | `bool` | `false` | Ejecutar el asistente TUI interactivo |
| `--force` | `-f` | `bool` | `false` | Sobrescribir la configuración existente en un directorio no vacío |

## Cómo funciona

### Modo no interactivo (predeterminado)

Crea un `swag2mcp.yaml` mínimo sin especificaciones. Usted edita el archivo manualmente después.

```bash
swag2mcp init
# Crea ~/.swag2mcp/swag2mcp.yaml

swag2mcp init ./my-project
# Crea ./my-project/swag2mcp.yaml

swag2mcp init /absolute/path
# Crea /absolute/path/swag2mcp.yaml
```

### Modo interactivo (`-i`)

Inicia un asistente TUI de 18 pasos que le guía a través de:

1. Elegir el directorio del espacio de trabajo
2. Agregar especificaciones con dominio, título, URL base
3. Configurar colecciones con URL de ubicación
4. Configurar autenticación (los 9 métodos)
5. Configurar ajustes del cliente HTTP (tiempo de espera, proxy, encabezados, etc.)

```bash
swag2mcp init -i
```

### Modo forzado (`--force`)

Por defecto, `init` se niega a ejecutarse en un directorio no vacío. Use `--force` para sobrescribir:

```bash
swag2mcp init -f
swag2mcp init ./existing-dir -f
```

## Qué se crea

```
~/.swag2mcp/
├── swag2mcp.yaml       # Archivo de configuración
├── cache/               # Archivos de especificación remotos descargados
├── specs/               # Archivos de especificación locales
├── responses/           # Respuestas de invocación de API guardadas
└── auth_scripts/        # Scripts de autenticación (para tipo ScriptAuth)
```

## Verificación posterior al comando

```bash
ls ~/.swag2mcp/swag2mcp.yaml
# Si el archivo existe, init tuvo éxito
```

## Matices

- **Resolución de ruta:** `[path]` es un **directorio de espacio de trabajo**, no una ruta de archivo. La CLI agrega `swag2mcp.yaml` automáticamente. Orden de resolución: `[path]` explícito → directorio actual (`./`) → `~/.swag2mcp/`.
- **Verificación de directorio no vacío:** Sin `--force`, `init` devuelve un error si el directorio de destino existe y no está vacío. Esto evita sobrescrituras accidentales.
- **Plantillas de scripts de autenticación:** Si alguna especificación usa `ScriptAuth`, `init` crea archivos de script de plantilla (`.sh` en Unix, `.bat` en Windows) en `auth_scripts/`.
- **Salida:** En caso de éxito, imprime la ruta de configuración y una sugerencia: `"Next step: edit swag2mcp.yaml or run 'swag2mcp ls' to list configured specs"`.
