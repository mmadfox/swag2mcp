# Variables de Entorno

## Descripción General

swag2mcp admite la sustitución de variables de entorno en el archivo de configuración usando la sintaxis `$(VAR_NAME)`. Esto le permite mantener datos sensibles (tokens, contraseñas, claves) fuera del archivo YAML.

## Cómo funciona

Cuando swag2mcp se inicia, escanea la configuración en busca de patrones `$(VAR_NAME)` y los reemplaza con el valor de la variable de entorno correspondiente.

```yaml
specs:
  - domain: meteo
    auth:
      type: bearer
      config:
        token: "$(API_TOKEN)"
```

Si la variable de entorno `API_TOKEN` está establecida, se sustituirá. Si no está establecida, el valor se vuelve vacío.

## Dónde se resuelve `$(VAR)`

| Campo | Ejemplo |
|-------|---------|
| `token` de autenticación (bearer) | `token: "$(API_TOKEN)"` |
| `username` / `password` de autenticación (basic, digest) | `password: "$(API_PASSWORD)"` |
| `client_id` / `client_secret` de autenticación (oauth2-cc, oauth2-pwd) | `client_secret: "$(OAUTH_SECRET)"` |
| `api_key` / `secret_key` de autenticación (hmac) | `api_key: "$(BINANCE_API_KEY)"` |
| `domain` de autenticación (script) | `domain: "$(AUTH_DOMAIN)"` |
| Token del servidor MCP | `token: "$(MCP_TOKEN)"` |
| Encabezados del cliente HTTP | `"X-API-Key": "$(API_KEY)"` |
| Valores de cookies del cliente HTTP | `value: "$(SESSION_TOKEN)"` |

## Dónde NO se resuelve `$(VAR)`

- URLs base (`base_url`)
- Ubicaciones de colecciones (`location`)
- Nombres de dominio de especificaciones (`domain`)

## Ejemplo

```bash
export API_TOKEN="eyJhbGciOiJIUzI1NiIs..."
export MCP_TOKEN="my-secret-token"

swag2mcp mcp
```

## Mejores prácticas de seguridad

- **Nunca** almacene secretos directamente en el archivo YAML
- Use variables de entorno o un gestor de secretos externo
- Agregue el archivo YAML a `.gitignore` si contiene algún secreto codificado
- Establezca las variables de entorno en su perfil de shell, configuración del IDE o pipeline de despliegue

## Detalles de sintaxis

- `$(VAR_NAME)` — sintaxis estándar
- `$( VAR_NAME )` — se permiten y recortan espacios en blanco dentro de los paréntesis
- `$()` — un nombre de variable vacío devuelve la cadena original sin cambios
- Los patrones `$(...)` anidados no se resuelven
