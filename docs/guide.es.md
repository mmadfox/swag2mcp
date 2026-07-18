# swag2mcp

**swag2mcp** es una herramienta CLI y servidor MCP (Model Context Protocol) que conecta especificaciones OpenAPI/Swagger/Postman con agentes LLM (Opencode, Crush, Copilot, Cursor, etc.).

Indexa sus especificaciones API en un motor de búsqueda de texto completo, las expone a través de 16 herramientas MCP y permite a los LLM descubrir, inspeccionar e invocar endpoints API reales — sin escribir una sola línea de código de integración.

---

## Tabla de contenidos

- [Inicio rápido](#inicio-rápido)
- [Configuración](#configuración)
- [Comandos CLI](#comandos-cli)
- [Servidor MCP](#servidor-mcp)
- [Búsqueda](#búsqueda)
- [Espacio de trabajo (Workspace)](#espacio-de-trabajo-workspace)
- [Caché](#caché)
- [Desarrollo](#desarrollo)

---

## Inicio rápido

### Opción 1 — Descargar desde GitHub Releases (recomendado)

1. Abra https://github.com/mmadfox/swag2mcp/releases/latest
2. Encuentre el archivo para su sistema:

   | SO | Arquitectura | Archivo |
   |----|-------------|---------|
   | Linux | x86_64 | `swag2mcp_<version>_linux_amd64.tar.gz` |
   | Linux | ARM64 | `swag2mcp_<version>_linux_arm64.tar.gz` |
   | macOS | Intel | `swag2mcp_<version>_darwin_amd64.tar.gz` |
   | macOS | Apple Silicon | `swag2mcp_<version>_darwin_arm64.tar.gz` |
   | Windows | x86_64 | `swag2mcp_<version>_windows_amd64.zip` |

3. Descargue e instale:

   **Linux / macOS:**
   ```bash
   tar -xzf swag2mcp_<version>_<os>_<arch>.tar.gz
   sudo mv swag2mcp /usr/local/bin/
   swag2mcp --version
   ```

   **Windows (PowerShell):**
   ```powershell
   Expand-Archive swag2mcp_<version>_windows_amd64.zip -DestinationPath .
   move swag2mcp.exe C:\Windows\System32\
   swag2mcp --version
   ```

4. (Opcional) Repita para el servidor mock — descargue `swag2mcp-mock_<version>_<os>_<arch>.tar.gz`

### Opción 2 — Instalar con Go

Si tiene Go instalado:

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

### Después de la instalación

```bash
# Inicializar el espacio de trabajo
swag2mcp init

# Iniciar el servidor MCP (para agentes LLM)
swag2mcp mcp

# O usar el explorador interactivo
swag2mcp run
```---

## Example LLM Queries

After setup, try asking your agent:

| Query | What happens |
|-------|-------------|
| "Show me all available APIs" | `spec_list` — lists petstore, binance, dadjoke, pokeapi |
| "What endpoints does Binance have?" | `endpoint_by_spec` — shows 4 market data endpoints |
| "Find endpoints related to pets" | `search("pet")` — finds petstore endpoints |
| "What tags are in the Petstore API?" | `tag_by_spec` — shows "pets" tag |
| "Show me the GET /pets endpoint details" | `inspect` — shows parameters and response schema |
| "Get the current BTC price from Binance" | `invoke` — real API call to Binance |
| "Get a random dad joke" | `invoke` — calls icanhazdadjoke API |

---

---

## Configuración

### Esquema YAML

```yaml
mock_enabled: true                    # opcional, activa el modo de servidor mock

http_client:                        # opcional, valores predeterminados HTTP globales
  headers:                          # opcional
    X-API-Version: "2"
  cookies: []                       # opcional
  user_agent: ""                    # opcional
  timeout: 0s                       # opcional
  follow_redirects: true            # opcional
  max_redirects: 10                 # opcional
  max_response_size: 1048           # opcional, bytes (predet. 1KB, máx 1MB)

specs:
  - domain: petstore                    # obligatorio, 1-60 car., [a-zA-Z0-9_-]
    llm_title: Petstore API             # obligatorio, 5-120 car.
    llm_instruction: |                  # opcional, máx 500 car.
      Usa esta API para gestionar mascotas, pedidos y usuarios.
    base_url: https://petstore.swagger.io/v2  # obligatorio, URL válida
    disable: false                      # opcional
    tags: [public, demo]                # opcional, para filtrar
    http_client:                        # opcional, sobrescribe el global
      headers:
        X-API-Version: "2"
    auth:                               # opcional
      type: bearer                      # ver Métodos de autenticación
      config:
        token: $(TOKEN_AUTH)
    collections:
      - llm_title: Petstore Swagger     # opcional, máx 120 car.
        llm_instruction: |             # opcional, máx 360 car.
          Endpoints principales de Petstore
        title: ""                      # opcional, auto-rellenado desde la spec
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/petstore.json  # obligatorio, 5-250 car.
        disable: false                  # opcional
        base_url: ""                    # opcional, sobrescribe base_url de la spec
        base_mock_url: localhost:8080   # opcional, formato "host:port" o "host:port/path"
        http_client: {}                 # opcional, sobrescribe la spec
```

### Tags — Filtrado de especificaciones por proyecto

Los tags permiten agrupar especificaciones por proyecto, entorno o equipo. Al iniciar el servidor MCP, use `--tags` para cargar solo las especificaciones coincidentes:

```bash
# Iniciar servidor solo con especificaciones públicas
swag2mcp mcp --tags=public

# Iniciar servidor con múltiples tags
swag2mcp mcp --tags=public,internal

# Ejecutar múltiples servidores para diferentes proyectos
swag2mcp mcp --tags=project-alpha --logfile=/tmp/swag2mcp-alpha.log
swag2mcp mcp --tags=project-beta  --logfile=/tmp/swag2mcp-beta.log
```

Esto permite ejecutar servidores MCP separados para diferentes proyectos desde un solo archivo de configuración.

### Métodos de autenticación

| Tipo | Campos | Ejemplo de configuración |
|------|--------|--------------------------|
| `none` | — | `type: none` |
| `basic` | `username`, `password` | `username: $(USER)`, `password: $(PASS)` |
| `bearer` | `token` | `token: $(TOKEN)` |
| `digest` | `username`, `password` | `username: admin`, `password: secret` |
| `hmac` | `api_key`, `secret_key` | `api_key: $(API_KEY)`, `secret_key: $(SECRET_KEY)` |
| `api-key` | `key`, `value`, `in` (header/query) | `key: X-API-Key`, `value: $(KEY)`, `in: header` |
| `oauth2-cc` | `client_id`, `client_secret`, `token_url`, `scopes` | `client_id: $(ID)`, `token_url: https://auth.example.com/token` |
| `oauth2-pwd` | `username`, `password`, `client_id`, `client_secret`, `token_url`, `scopes` | `username: $(USER)`, `token_url: https://auth.example.com/token` |
| `script` | `source` | `source: ruta/a/auth.sh` |

Todos los campos de cadena admiten la sintaxis `$(ENV_VAR)` — se resuelve en tiempo de ejecución desde variables de entorno.

---

## Comandos CLI

Todos los comandos que aceptan `[path]` usan la misma resolución de ruta:

```
swag2mcp <command>          → ~/.swag2mcp/
swag2mcp <command> ./       → {cwd}/.swag2mcp/
swag2mcp <command> path/to  → {cwd}/path/to/.swag2mcp/
```

### `init [path]`

Inicializar espacio de trabajo y configuración.

| Banderas | Corto | Predet. | Descripción |
|----------|-------|---------|-------------|
| `--interactive` | `-i` | `false` | Asistente interactivo |
| `--force` | `-f` | `false` | Sobrescribir configuración existente |

```bash
swag2mcp init              # crear ~/.swag2mcp/swag2mcp.yaml
swag2mcp init ./           # crear ./.swag2mcp/swag2mcp.yaml
swag2mcp init -i           # asistente interactivo
```

### `add spec [path]` / `add collection [path]`

Agregar una especificación o colección a la configuración.

| Banderas | Corto | Predet. | Descripción |
|----------|-------|---------|-------------|
| `--yaml` | `-y` | `""` | Entrada YAML (use `-` para stdin) |
| `--example` | `-e` | `false` | Mostrar ejemplo YAML |

```bash
swag2mcp add spec
swag2mcp add spec --yaml 'domain: petstore\nllm_title: Petstore API\nbase_url: https://...'
cat spec.yaml | swag2mcp add spec --yaml -
swag2mcp add spec --example
```

### `delete spec [path]` / `delete collection [path]`

Eliminar una especificación o colección. Selección interactiva.

```bash
swag2mcp delete spec
swag2mcp delete collection
```

### `ls [path]`

Listar especificaciones y colecciones.

| Banderas | Corto | Predet. | Descripción |
|----------|-------|---------|-------------|
| `--tags` | `-t` | `""` | Filtrar por tags (separados por comas) |

```bash
swag2mcp ls
swag2mcp ls --tags=public,internal
```

### `run [path]`

Explorador API interactivo (TUI). Buscar, navegar, inspeccionar y guardar endpoints.

```bash
swag2mcp run
```

### `validate [path]`

Validar la configuración y verificar que todas las ubicaciones de colecciones sean accesibles.

| Banderas | Corto | Predet. | Descripción |
|----------|-------|---------|-------------|
| `--tags` | `-t` | `""` | Filtrar especificaciones por tags |

```bash
swag2mcp validate
swag2mcp validate --tags=public
```

### `clean [path]`

Eliminar todo el contenido de los directorios `cache/` y `responses/`.

```bash
swag2mcp clean
```

### `update [path]`

Validar configuración, limpiar caché, volver a cachear todos los archivos de especificación.

```bash
swag2mcp update
```

### `mcp [path]`

Iniciar el servidor MCP en modo headless (transporte stdio). Este es el comando de producción principal para la integración con LLM.

| Banderas | Corto | Predet. | Descripción |
|----------|-------|---------|-------------|
| `--logfile` | `-f` | `""` | Ruta del archivo de registro |
| `--tags` | `-t` | `""` | Filtrar especificaciones por tags |
| `--disable-llm-auth` | | `true` | `true` — autenticación en segundo plano (LLM nunca ve tokens). `false` — LLM puede solicitar tokens mediante la herramienta `auth` |
| `--dump-dir` | | `""` | Directorio para volcar solicitudes HTTP (depuración) |

```bash
swag2mcp mcp
swag2mcp mcp --tags=public --logfile=/var/log/swag2mcp.log
swag2mcp mcp --disable-llm-auth=false
swag2mcp mcp --dump-dir=/tmp/dump
```

### `mockserver [path]`

Inicia servidores HTTP mock para todas las especificaciones API. Cada colección obtiene su propio
servidor HTTP que genera datos aleatorios que coinciden con los esquemas de respuesta OpenAPI.

| Banderas | Predet. | Descripción |
|----------|---------|-------------|
| `--tls` | `false` | Habilitar TLS con certificado autofirmado |
| `--tls-cert` | `""` | Ruta del archivo de certificado TLS |
| `--tls-key` | `""` | Ruta del archivo de clave TLS |

```bash
swag2mcp-mock
swag2mcp-mock --tls
```

**Flujo de trabajo:**
1. Agregue `mock_enabled: true` y `base_mock_url` a su configuración
2. Inicie el servidor mock: `swag2mcp-mock`
3. Inicie el servidor MCP: `swag2mcp mcp` — invoke usará `base_mock_url` en lugar de `base_url`
4. La autenticación se aplica automáticamente: OAuth2/Digest usan servidores mock en puertos 9090/9091; otros tipos aplican credenciales directamente

### Autenticación mock

Cuando `auth` está configurado en una especificación, el servidor MCP aplica
la autenticación automáticamente. Solo dos tipos de autenticación necesitan
un servidor mock dedicado:

| Tipo de auth | Endpoint mock | Comportamiento |
|--------------|---------------|----------------|
| `oauth2-cc` / `oauth2-pwd` | `POST /token` en puerto 9090 | Acepta cualquier `client_id`/`username`+`password`, devuelve `{"access_token":"<random>","token_type":"Bearer","expires_in":3600}` |
| `digest` | `GET /` en puerto 9091 | Envía un desafío 401 con `algorithm=MD5`, acepta cualquier respuesta Digest, devuelve `{"status":"authenticated","method":"digest"}` |

Otros tipos (`basic`, `bearer`, `api-key`, `hmac`, `script`) **no requieren** un
servidor mock — el servidor MCP aplica las credenciales configuradas a cada
solicitud automáticamente.

---

## Integración

swag2mcp habla el Protocolo de Contexto de Modelo (MCP) y funciona con cualquier cliente compatible con MCP.

### Local (stdio) — agente en la misma máquina

Iniciar el servidor:

```bash
swag2mcp mcp
```

| Cliente | Archivo de configuración | Contenido |
|---------|-------------------------|-----------|
| **OpenCode** | `opencode.json` | `{"mcp":{"swag2mcp":{"type":"local","command":["swag2mcp","mcp"]}}}` |
| **Cursor** | `.cursor/mcp.json` | `{"mcpServers":{"swag2mcp":{"command":"swag2mcp","args":["mcp"]}}}` |
| **Claude Desktop** | `claude_desktop_config.json` | `{"mcpServers":{"swag2mcp":{"command":"swag2mcp","args":["mcp"]}}}` |
| **VS Code** | `.vscode/mcp.json` | `{"servers":{"swag2mcp":{"type":"stdio","command":"swag2mcp","args":["mcp"]}}}` |
| **Crush** | `crush.json` | `{"mcp":{"swag2mcp":{"type":"stdio","command":"swag2mcp","args":["mcp"]}}}` |

### Remoto (HTTP) — agente en la nube / otra máquina

Iniciar el servidor con transporte HTTP:

```bash
swag2mcp mcp --transport streamable-http --http-addr :8080 --auth-token my-secret
```


> **Note:** If you initialized the workspace at a custom path (e.g. `swag2mcp init ./my-project`), you must specify the path when starting the MCP server: `swag2mcp mcp ./my-project`. The IDE configuration must also use the full path to the config file.

O configurar en `swag2mcp.yaml`:

```yaml
mcp:
  transport: streamable-http
  addr: ":8080"
  path: "/mcp"
  auth_token: $(MCP_AUTH_TOKEN)
```

| Cliente | Archivo de configuración | Contenido |
|---------|-------------------------|-----------|
| **OpenCode** | `opencode.json` | `{"mcp":{"swag2mcp":{"type":"remote","url":"http://localhost:8080/mcp","headers":{"Authorization":"Bearer ${MCP_AUTH_TOKEN}"}}}}` |
| **VS Code** | `.vscode/mcp.json` | `{"servers":{"swag2mcp":{"type":"http","url":"http://localhost:8080/mcp"}}}` |

> **Verificación de salud** (funciona sin handshake MCP):
> ```bash
> curl http://localhost:8080/health
> # → {"status":"ok","version":"v1.1.3"}
> ```

---

## Servidor MCP

El servidor MCP expone 16 herramientas a través del transporte stdio o HTTP. Los agentes LLM (Opencode, Cursor, Claude, Copilot, Crush, etc.) se conectan automáticamente cuando están configurados.

### Jerarquía de herramientas

```
spec_list                       — listar todas las especificaciones disponibles
  └─ spec_by_id                 — detalles de especificación por ID
       └─ collection_by_spec    — colecciones en una especificación
            └─ tag_by_collection     — tags en una colección
                 └─ endpoint_by_tag  — endpoints bajo un tag
                      └─ inspect          — operación OpenAPI completa
                           └─ invoke       — ejecutar llamada API

search                          — búsqueda de texto completo en todos los endpoints
```

### Referencia de herramientas

| Herramienta | Argumentos | Devuelve | Descripción |
|-------------|------------|----------|-------------|
| `spec_list` | — | `Spec[]` | Todas las especificaciones disponibles |
| `spec_by_id` | `id` | Spec + Collections | Detalles de especificación |
| `collection_by_spec` | `specId` | Collections | Colecciones en una especificación |
| `collection_by_id` | `id` | Collection + Tags | Detalles de colección |
| `tag_by_collection` | `collectionId` | Tags | Tags en una colección |
| `tag_by_spec` | `specId` | Tags | Todos los tags de una especificación |
| `tag_by_id` | `id` | Tag | Metadatos de un tag |
| `endpoint_by_tag` | `tagId` | Endpoints | Endpoints bajo un tag |
| `endpoint_by_collection` | `collectionId` | Endpoints | Todos los endpoints de una colección |
| `endpoint_by_spec` | `specId` | Endpoints | Todos los endpoints de una especificación |
| `endpoint_by_id` | `id` | Endpoint | Resumen rápido de endpoint |
| `search` | `query`, `limit` | Endpoints | Búsqueda de texto completo |
| `inspect` | `endpointId` | Full Operation | Objeto de operación OpenAPI completo |
| `invoke` | `endpointId`, `parameters`, `requestBody` | Response | Ejecuta llamada API real |
| `auth` | `specId` | Token | Obtener token de autenticación para una especificación |

---

## Búsqueda

### Sintaxis de consultas

| Característica | Sintaxis | Ejemplo |
|----------------|----------|---------|
| Término | `término` | `mascotas` |
| Frase | `"frase"` | `"agregar mascota"` |
| Campo: method | `method:término` | `method:post` |
| Campo: tag | `tag:término` | `tag:auth` |
| Campo: path | `path:término` | `path:/users` |
| Campo: summary | `summary:término` | `summary:login` |
| Requerido (AND) | `+término` | `+method:post +tag:user` |
| Excluido (NOT) | `-término` | `-deprecated` |
| Comodín | `*` | `path:*/v2/*` |
| Difuso | `término~` | `watex~` |
| Regex | `/patrón/` | `/user(s\|sessions)/` |
| Potenciación | `término^N` | `tag:pet^5` |
| Todo | `*` | `*` |

### Ejemplos

```
# Encontrar endpoints POST en el tag auth
+method:post +tag:auth

# Buscar endpoints relacionados con inicio de sesión
summary:"login"~

# Encontrar todas las rutas de usuarios, excluir obsoletas
path:*/users/* -deprecated

# Consulta compleja
+method:get +tag:pet summary:"find by status"
```

### Campos indexados

| Campo | Tipo | Contenido |
|-------|------|-----------|
| `method` | text | Método HTTP (minúsculas) |
| `tag` | text | Nombre del tag (minúsculas) |
| `path` | text | Ruta API (minúsculas) |
| `summary` | text (analizado) | Resumen/descripción del endpoint (minúsculas) |
| `_all` | text (analizado) | method + path + tag + summary |

---

## Espacio de trabajo (Workspace)

### Estructura de directorios

```
~/.swag2mcp/                    # o {proyecto}/.swag2mcp/
├── swag2mcp.yaml               # Archivo de configuración
├── cache/                      # Especificaciones remotas en caché
│   ├── {hash}.spec             # Contenido del archivo de especificación
│   └── {hash}.meta             # Metadatos JSON
├── specs/                      # Archivos de especificación locales (gestionados por el usuario)
├── responses/                  # Archivos de respuesta de invocaciones
└── auth_scripts/               # Scripts de autenticación
```

### Resolución de ruta

```
swag2mcp <command>          → ~/.swag2mcp/
swag2mcp <command> ./       → {cwd}/.swag2mcp/
swag2mcp <command> path/to  → {cwd}/path/to/.swag2mcp/
```

### .gitignore

Solo los datos temporales deben ignorarse:

```
.swag2mcp/cache/*
.swag2mcp/responses/*
```

La configuración `.swag2mcp/swag2mcp.yaml` y los archivos de especificación en `.swag2mcp/specs/` **deben estar en el repositorio**.

### Recomendación

Mantenga todos los archivos de especificación en `.swag2mcp/specs/` — esta es la única forma de garantizar que se usen directamente sin copiarse al caché.

---

## Caché

### Reglas

| Fuente | Comportamiento |
|--------|----------------|
| URL HTTP/HTTPS | Siempre se almacena en caché. TTL: aleatorio 1-48h. |
| Ruta local dentro de `specs/` | Se usa directamente, no se almacena en caché. |
| Ruta local fuera de `specs/` | Se copia al caché en el primer acceso. |
| URL `file://` | Se trata como ruta local. |

### Clave de caché

Hash SHA-256 de la ubicación normalizada (primeros 16 bytes = 32 caracteres hex).

### Lógica de acierto de caché

1. Leer archivo `.meta` — caducado o faltante → fallo
2. Para fuentes locales: `ModTime` cambiado → fallo
3. Archivo `.spec` faltante → fallo
4. De lo contrario → acierto

---

## Desarrollo

```bash
# Compilar
go build ./cmd/swag2mcp/

# Pruebas
go test ./...

# Linter
make lint

# Ejecutar
go run ./cmd/swag2mcp/main.go
```
