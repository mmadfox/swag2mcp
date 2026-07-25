# Archivo de Configuración

swag2mcp utiliza un archivo de configuración YAML. Creado por `swag2mcp init`.

## Ubicación

- **Linux/macOS**: `~/.swag2mcp/swag2mcp.yaml`
- **Windows**: `%USERPROFILE%\.swag2mcp\swag2mcp.yaml`

## Estructura Básica

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

## Ejemplo Completo

```yaml
# ── Cliente HTTP global ──────────────────────────────────
http_client:
  timeout: 30s
  max_response_size: 1048576
  user_agent: "swag2mcp-global/1.0"
  follow_redirects: true
  max_redirects: 10
  random: false
  proxy:
    url: ""
    username: ""
    password: ""
    bypass: []
  headers:
    "Accept": "application/json"
  cookies:
    - name: "session"
      value: "abc123"

# ── Servidor MCP ──────────────────────────────────────────
mcp:
  transport: stdio
  addr: ":8080"
  path: "/mcp"
  auth:
    token: ""

# ── Servidor simulado ─────────────────────────────────────
mock_enabled: false
mock_auth:
  oauth2_port: 9090
  digest_port: 9091
  hmac_port: 9092

# ── Limitador de velocidad ────────────────────────────────
disable_ratelimiter: false
rate_limit_interval: 10s

# ── Especificaciones ─────────────────────────────────────
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    llm_instruction: "Use esta API para pronósticos meteorológicos y datos climáticos"
    base_url: https://api.open-meteo.com
    disable: false
    tags: ["weather", "climate"]
    http_client:
      timeout: 10s
    auth:
      type: bearer
      config:
        token: "$(TOKEN)"
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
        disable: false
        http_client:
          timeout: 5s

  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Variables de Entorno

Use la sintaxis `$(VAR_NAME)` para referenciar variables de entorno. swag2mcp las resuelve al iniciar.

```yaml
specs:
  - domain: meteo
    auth:
      type: bearer
      config:
        token: "$(API_TOKEN)"

mcp:
  auth:
    token: "$(MCP_TOKEN)"
```

`$(VAR)` se resuelve en:
- Campos de configuración de autenticación: `token`, `username`, `password`, `client_id`, `client_secret`, `api_key`, `secret_key`, `domain`
- Token de autenticación del servidor MCP: `mcp.auth.token`
- Encabezados del cliente HTTP y valores de cookies

`$(VAR)` **no** se resuelve en URL base ni ubicaciones de colecciones.

## Validación

```bash
# Validar espacio de trabajo predeterminado (~/.swag2mcp)
swag2mcp validate

# Validar un espacio de trabajo de proyecto personalizado
swag2mcp validate ./my-project
```

Si el espacio de trabajo no está en el directorio personal (por ejemplo, dentro de un repositorio de proyecto), siempre especifique la ruta al ejecutar `validate`, `update`, `mcp` o cualquier otro comando. De lo contrario, swag2mcp usará el espacio de trabajo predeterminado `~/.swag2mcp`.
