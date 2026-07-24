# Preguntas Frecuentes

## Generales

### ¿Qué es swag2mcp y qué problema resuelve?

swag2mcp conecta especificaciones de API OpenAPI/Swagger/Postman con agentes LLM mediante el Protocolo de Contexto de Modelo (MCP). En lugar de escribir código personalizado para conectar cada API a un agente de IA, lo configura una vez en un archivo YAML y el LLM obtiene 19 herramientas para descubrir, inspeccionar y llamar a sus APIs.

### ¿En qué se diferencia de otras herramientas API-a-LLM?

- **Sin necesidad de programación** — configure APIs en YAML, no necesita código de integración
- **19 herramientas MCP** — conjunto completo desde descubrimiento hasta invocación y manejo de respuestas grandes
- **9 métodos de autenticación** — funciona con cualquier esquema de autenticación de API
- **Búsqueda de texto completo** — búsqueda impulsada por bluge en todos los endpoints
- **Explorador TUI** — interfaz de terminal interactiva para navegar y probar
- **Servidor simulado** — pruebe sin llamadas reales a la API

### ¿Qué formatos de especificación de API son compatibles?

OpenAPI 3.x, Swagger 2.0 y Postman Collections v2.1.

### ¿Cuál es la diferencia entre una especificación y una colección?

Una **especificación** representa un servicio de API lógico (por ejemplo, "Open-Meteo Weather APIs"). Una **colección** es un archivo OpenAPI/Swagger/Postman. Una especificación puede tener múltiples colecciones — por ejemplo, cuando una API tiene archivos de especificación separados para diferentes servicios (pronóstico, calidad del aire, marino).

### ¿Qué transportes MCP son compatibles?

Tres transportes: `stdio` (predeterminado, para clientes LLM locales), `sse` (Eventos Enviados por el Servidor para clientes remotos) y `streamable-http` (transmisión HTTP moderna).

### ¿Puedo usar swag2mcp con cualquier LLM?

Sí, cualquier cliente LLM que admita el protocolo MCP: Claude Desktop, VS Code, Cursor, Windsurf, IDEs JetBrains, OpenCode y otros.

## Instalación

### ¿Cómo instalo swag2mcp?

```bash
# Opción 1: Descargar desde GitHub Releases
# Vaya a https://github.com/mmadfox/swag2mcp/releases/latest
# Descargue el archivo para su sistema operativo y arquitectura

# Opción 2: Instalar con Go
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

### ¿Necesito tener Go instalado?

No. Los binarios precompilados están disponibles para Linux (amd64, arm64), macOS (amd64, arm64) y Windows (amd64) en la [página de GitHub Releases](https://github.com/mmadfox/swag2mcp/releases).

### ¿Cómo instalo el servidor simulado?

El servidor simulado es un binario separado:

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

O descargue `swag2mcp-mock_<version>_<os>_<arch>.tar.gz` desde GitHub Releases.

## Primeros Pasos

### ¿Cómo empiezo rápidamente?

```bash
# 1. Inicialice un espacio de trabajo
mkdir -p .swag2mcp && swag2mcp init ./.swag2mcp

# 2. Inicie el servidor MCP (las especificaciones de ejemplo público están incluidas después de init)
swag2mcp mcp
```

Después de `init`, el espacio de trabajo ya incluye varias especificaciones de ejemplo público (icanhazdadjoke, Open-Meteo, Binance, PokéAPI). Puede iniciar el servidor MCP inmediatamente — no es necesario agregar especificaciones manualmente.

Si desea agregar su propia API:

```bash
swag2mcp add spec --yaml - <<EOF
domain: dadjoke
llm_title: icanhazdadjoke API
base_url: https://icanhazdadjoke.com
collections:
  - llm_title: Jokes
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
EOF
```

### ¿Cómo conecto swag2mcp a mi IDE?

**VS Code** (`.vscode/settings.json`):
```json
{
  "mcp": {
    "servers": {
      "swag2mcp": {
        "command": "swag2mcp",
        "args": ["mcp", "/absolute/path/to/.swag2mcp"]
      }
    }
  }
}
```

**Cursor** (`~/.cursor/mcp.json`):
```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "/absolute/path/to/.swag2mcp"]
    }
  }
}
```

**Claude Desktop** (`claude_desktop_config.json`):
```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "/absolute/path/to/.swag2mcp"]
    }
  }
}
```

Use siempre una ruta absoluta al directorio del espacio de trabajo.

## Configuración

### ¿Dónde se encuentra el archivo de configuración?

Predeterminado: `~/.swag2mcp/swag2mcp.yaml`. También puede crearlo en cualquier directorio y pasar la ruta a los comandos.

### ¿Cómo agrego una API?

```bash
# Modo interactivo
swag2mcp add spec

# Con YAML (recomendado para scripting)
swag2mcp add spec --yaml - <<EOF
domain: my-api
llm_title: My API
base_url: https://api.example.com/v1
collections:
  - llm_title: Main
    location: https://example.com/spec.yaml
EOF
```

### ¿Cómo agrego una colección a una especificación existente?

```bash
swag2mcp add collection --yaml - <<EOF
spec_domain: meteo
llm_title: Air Quality
location: https://example.com/air-quality.yaml
EOF
```

### ¿Cómo deshabilito una especificación temporalmente?

Establezca `disable: true` en la configuración de la especificación. La especificación no se cargará ni indexará.

### ¿Puedo filtrar qué especificaciones se cargan?

Sí, use la bandera `--tags`: `swag2mcp mcp --tags=public`. Solo se cargarán las especificaciones con etiquetas coincidentes.

### ¿Cómo uso variables de entorno para secretos?

Use la sintaxis `$(VAR_NAME)` en los campos de autenticación:

```yaml
auth:
  type: bearer
  config:
    token: "$(MY_API_TOKEN)"
```

Establezca la variable antes de iniciar: `export MY_API_TOKEN="eyJhbGci..."`

## Autenticación

### ¿Qué métodos de autenticación son compatibles?

Nueve métodos: `none`, `basic`, `bearer`, `digest`, `hmac`, `oauth2-cc` (credenciales de cliente), `oauth2-pwd` (concesión de contraseña), `api-key` y `script`.

### ¿Cómo paso un token?

A través del archivo de configuración o variables de entorno:

```yaml
auth:
  type: bearer
  config:
    token: "$(MY_TOKEN)"
```

### ¿Necesito llamar a auth antes de invoke?

No. La herramienta `invoke` aplica automáticamente la autenticación desde la configuración de la especificación. Solo necesita la herramienta MCP `auth` si desea mostrar el token al usuario (por ejemplo, para un comando curl).

### ¿Por qué no aparece la herramienta auth?

La herramienta `auth` está deshabilitada por defecto (`--disable-llm-auth=true`). Esto es una medida de seguridad para producción. Para habilitarla: `swag2mcp mcp --disable-llm-auth=false`.

### ¿Cómo se renuevan los tokens OAuth2?

Los tokens de OAuth2 de Credenciales de Cliente y Concesión de Contraseña se renuevan automáticamente cuando expiran. Los tokens Bearer son estáticos y deben actualizarse manualmente.

## Servidor MCP

### ¿Cómo inicio el servidor MCP?

```bash
# Predeterminado (transporte stdio)
swag2mcp mcp

# Con transporte HTTP
swag2mcp mcp --transport sse --http-addr :8080
```

### ¿Cómo cambio el puerto?

```bash
swag2mcp mcp --transport sse --http-addr 0.0.0.0:9090
```

### ¿Cómo aseguro el endpoint HTTP del MCP?

Establezca un token bearer:

```bash
swag2mcp mcp --transport sse --http-addr :8080 --auth-token "my-secret"
```

El cliente LLM debe incluir `Authorization: Bearer my-secret` en cada solicitud.

### ¿Cuál es el protocolo de enlace MCP para transporte HTTP?

Para los transportes SSE y Streamable HTTP, el protocolo MCP requiere un protocolo de enlace de tres pasos:

```
Paso 1: POST /mcp → {"method":"initialize", ...}
Paso 2: POST /mcp → {"method":"notifications/initialized"}
Paso 3: POST /mcp → {"method":"tools/list", ...}  ← ahora funciona
```

Las llamadas a herramientas fallarán antes de la inicialización.

## Uso

### ¿Cómo busco endpoints?

Use la herramienta MCP `search` o la TUI (`swag2mcp run`). La búsqueda admite filtros de campo (`method:GET`, `tag:pets`), búsqueda difusa, comodines y operadores booleanos.

### ¿Cómo llamo a una API?

El LLM usa la herramienta MCP `invoke`. Siempre inspeccione el endpoint primero para entender los parámetros requeridos:

```
inspect(endpointId: "...")  → entender el contrato
invoke(endpointId: "...", parameters: {...})  → hacer la llamada
```

### ¿Qué sucede si una respuesta es demasiado grande?

Las respuestas que exceden `max_response_size` (predeterminado 1 MB) se guardan en disco. El LLM recibe una referencia de archivo y puede explorarla con las herramientas `response_outline`, `response_compress` y `response_slice`.

### ¿Cómo funciona el limitador de velocidad?

Cada endpoint tiene un período de enfriamiento de 10 segundos. Si el LLM llama al mismo endpoint dos veces en 10 segundos, la segunda llamada se bloquea silenciosamente. Puede deshabilitarlo o ajustarlo en la configuración.

### ¿Puedo probar sin hacer llamadas reales a la API?

Sí, use el servidor simulado:

```bash
swag2mcp-mock mockserver
```

Genera respuestas falsas basadas en esquemas OpenAPI.

## Gestión del Espacio de Trabajo

### ¿Cómo hago una copia de seguridad de mi configuración?

```bash
swag2mcp export --output ~/backups/swag2mcp-2026-07-24.zip
```

### ¿Cómo transfiero a otra máquina?

```bash
# En la máquina antigua
swag2mcp export --output swag2mcp.zip

# Copie el ZIP, luego en la máquina nueva
swag2mcp import --from-zip swag2mcp.zip
```

### ¿Cómo actualizo los archivos de especificación?

```bash
swag2mcp update
```

Esto revalida la configuración, limpia la caché y vuelve a descargar todos los archivos de especificación.

### ¿Cómo libero espacio en disco?

```bash
swag2mcp clean
```

Elimina los archivos de especificación en caché y las respuestas de API guardadas. Las respuestas antiguas (>48h) también se limpian automáticamente al iniciar el servidor MCP.

## TUI

### ¿Qué es la TUI y cómo la uso?

La TUI (Interfaz de Usuario de Terminal) es un explorador de API interactivo. Inícielo con `swag2mcp run`. Tiene tres modos: Búsqueda (búsqueda de texto completo), Navegación (navegación en árbol: Especificación → Colección → Etiqueta → Endpoint) y Autenticación (ver tokens).

### ¿Cuáles son los atajos de teclado?

| Tecla | Acción |
|-------|--------|
| `↑/↓` | Navegar |
| `Enter` | Seleccionar |
| `Esc` | Retroceder |
| `Tab` | Cambiar modos |
| `/` | Buscar |
| `N/P` | Siguiente/página anterior |
| `q` | Salir |

## Avanzado

### ¿Puedo usar un proxy?

Sí, configúrelo en `http_client.proxy`:

```yaml
http_client:
  proxy:
    url: "http://proxy.company.com:8080"
    username: "$(PROXY_USER)"
    password: "$(PROXY_PASS)"
    bypass:
      - "localhost"
      - "*.internal.com"
```

### ¿Puedo agregar un método de autenticación personalizado?

Sí, implemente la interfaz `Authenticator` en `internal/auth/` y regístrela en el analizador de configuración. Consulte la sección de Desarrollo para más detalles.

### ¿Puedo agregar una herramienta MCP personalizada?

Sí, agregue un método a la interfaz `Svc`, impleméntelo en la capa de servicio, agregue un controlador y regístrelo. Consulte la sección de Desarrollo para más detalles.

### ¿Cuál es la diferencia entre `swag2mcp` y `swag2mcp-mock`?

`swag2mcp` es el binario principal con comandos CLI y el servidor MCP. `swag2mcp-mock` es un binario separado que inicia servidores simulados para pruebas sin llamadas reales a la API.
