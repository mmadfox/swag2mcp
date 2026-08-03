# OAuth2 Concesión de Contraseña

## Propósito

Concesión de Contraseña del Propietario del Recurso OAuth2 — autenticación usando el nombre de usuario y contraseña de un usuario. Adecuado para aplicaciones de primera parte donde el usuario confía en la aplicación con sus credenciales.

## Cuándo usarlo

- Aplicaciones de primera parte (móvil, web)
- Integración con Keycloak y Proveedores de Identidad similares
- Cuando la API admite la Concesión de Contraseña OAuth2

## Configuración

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: oauth2-pwd
      config:
        client_id: "$(CLIENT_ID)"
        client_secret: "$(CLIENT_SECRET)"
        username: "$(USERNAME)"
        password: "$(PASSWORD)"
        token_url: "https://auth.example.com/oauth/token"
        scopes:
          - openid
          - profile
        request_format: form
```

## Parámetros

| Parámetro | Requerido | Descripción |
|-----------|-----------|-------------|
| `client_id` | Sí | Identificador del cliente |
| `username` | Sí | Nombre de usuario |
| `password` | Sí | Contraseña |
| `token_url` | Sí | URL del endpoint de token |
| `client_secret` | No | Secreto del cliente (opcional, para clientes públicos) |
| `scopes` | No | Lista de permisos (opcional) |
| `request_format` | No | Formato del cuerpo: `form` (predeterminado, `application/x-www-form-urlencoded`) o `json` (`application/json`) |

## Notas

- `client_secret` es opcional — se admiten **clientes públicos** (por ejemplo, Keycloak)
- swag2mcp renueva automáticamente el token cuando expira
- El token se almacena en caché hasta su expiración
- Los campos `client_id`, `client_secret`, `username`, `password` y `token_url` admiten variables de entorno mediante `$(VAR)`
- `request_format: json` envía la solicitud de token como cuerpo JSON — úselo cuando el endpoint de token requiera `Content-Type: application/json`
