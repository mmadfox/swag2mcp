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

## Notas

- `client_secret` es opcional — se admiten **clientes públicos** (por ejemplo, Keycloak)
- swag2mcp renueva automáticamente el token cuando expira
- El token se almacena en caché hasta su expiración
- Todos los parámetros pueden almacenarse en variables de entorno
