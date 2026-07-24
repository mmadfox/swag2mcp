# Comandos CLI

## Descripción General

La CLI de `swag2mcp` es el punto de entrada único para todas las operaciones — desde inicializar un espacio de trabajo y gestionar especificaciones de API hasta iniciar un servidor MCP para integración con LLM. Proporciona **13 comandos** que cubren el ciclo de vida completo del trabajo con especificaciones OpenAPI/Swagger/Postman.

### Qué resuelve la CLI

- **Ciclo de vida del espacio de trabajo** — crear (`init`), inspeccionar (`info`, `ls`), limpiar (`clean`), actualizar (`update`) y eliminar (`delete`) espacios de trabajo y su contenido
- **Gestión de especificaciones y colecciones** — agregar (`add`), listar (`ls`) y eliminar (`delete`) especificaciones de API y sus colecciones
- **Modos de ejecución** — iniciar el servidor MCP para acceso a herramientas LLM (`mcp`) o lanzar el explorador TUI interactivo (`run`)
- **Diagnóstico** — validar configuración (`validate`), mostrar versión (`version`), mostrar información de ejecución (`info`)
- **Copia de seguridad y restauración** — ciclo completo del espacio de trabajo mediante ZIP (`export`, `import`)

### Matices clave

- **Resolución de ruta** — los comandos que aceptan `[path]` esperan un **directorio de espacio de trabajo** (no una ruta de archivo). Orden de resolución: `[path]` explícito → directorio actual (`./`) → `~/.swag2mcp/`. La CLI agrega `swag2mcp.yaml` automáticamente. Siempre pase una ruta explícita cuando ejecute como servicio o en la configuración del IDE para evitar cargar el espacio de trabajo incorrecto.
- **Especificación vs Colección** — una **especificación** representa un servicio de API lógico (por ejemplo, "Open-Meteo API"), mientras que una **colección** es un archivo OpenAPI/Swagger/Postman. Una especificación puede tener múltiples colecciones.
- **`--version`** se admite tanto como bandera (`swag2mcp --version`) como subcomando (`swag2mcp version`).
- **`add spec` / `add collection`** aceptan entrada YAML a través de `--yaml` (cadena en línea o `-` para stdin). Redirigir desde un archivo o heredoc evita problemas de comillas del shell con caracteres especiales.
- **`delete`** requiere un TTY (terminal interactiva). No hay una bandera `--force` o `--yes` — siempre solicita selección y confirmación.
- **`mcp`** es el comando principal para la integración con LLM. Admite tres transportes: `stdio` (predeterminado), `sse` y `streamable-http`. La bandera `--disable-llm-auth` (predeterminado: `true`) elimina la herramienta `auth` de la lista de herramientas MCP, evitando que el LLM vea o solicite tokens. La autenticación sigue funcionando — los tokens se obtienen a través del mecanismo de configuración estándar, no a través del LLM. Este modo se recomienda para **producción** (el LLM nunca tiene acceso a las credenciales). Para **depuración** o cuando se usan tokens de corta duración, establezca `--disable-llm-auth=false` para permitir que el LLM solicite tokens frescos a través de la herramienta `auth`.
- **`validate`** verifica la sintaxis YAML, la estructura de configuración, la existencia de archivos de especificación, el alcance de URL, el formato de especificación (OpenAPI/Swagger/Postman), la configuración de autenticación y la corrección del cliente HTTP. **No** prueba los endpoints de autenticación ni la disponibilidad de endpoints de API.
- **`export` / `import`** proporcionan un ciclo completo del espacio de trabajo — el archivo de configuración, los archivos de especificación y los scripts de autenticación se incluyen en el archivo ZIP.
- **`clean`** elimina los directorios `cache/` y `responses/` pero conserva `specs/` y `auth_scripts/`. Las respuestas antiguas (>48h) también se limpian automáticamente al iniciar `mcp`.

## Comandos

| Comando | Descripción |
|---------|-------------|
| [`init`](/cli/init) | Inicializar un directorio de espacio de trabajo con configuración predeterminada |
| [`add`](/cli/add) | Agregar una especificación o colección a la configuración |
| [`delete`](/cli/delete) | Eliminar una especificación o colección de forma interactiva |
| [`ls`](/cli/ls) | Listar todas las especificaciones y sus colecciones |
| [`run`](/cli/run) | Iniciar el explorador de API TUI interactivo |
| [`validate`](/cli/validate) | Validar la configuración y los archivos de especificación |
| [`clean`](/cli/clean) | Limpiar especificaciones en caché y respuestas de invocación |
| [`update`](/cli/update) | Revalidar, recargar en caché y reindexar todas las especificaciones |
| [`mcp`](/cli/mcp) | Iniciar el servidor MCP para acceso a herramientas LLM |
| [`version`](/cli/version) | Imprimir la versión de swag2mcp |
| [`info`](/cli/info) | Mostrar información detallada de configuración y ejecución |
| [`import`](/cli/import) | Importar archivos de especificación o restaurar espacio de trabajo desde ZIP |
| [`export`](/cli/export) | Exportar espacio de trabajo como copia de seguridad ZIP portátil |

## Banderas Globales

| Bandera | Descripción |
|---------|-------------|
| `--version` | Mostrar versión (igual que el subcomando `version`) |
| `--help` | Mostrar ayuda para cualquier comando |
