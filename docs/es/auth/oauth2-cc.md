# OAuth2 Credenciales de Cliente

## Propósito

Concesión de Credenciales de Cliente OAuth2 — autenticación para comunicación servidor a servidor. La aplicación obtiene un token usando su client_id y client_secret, sin intervención del usuario.

## Cuándo usarlo

- Microservicios e integraciones servidor a servidor
- Comunicación máquina a máquina
- Cuando la API usa OAuth2 y tiene un client_id + client_secret

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
      type: oauth2-cc
      config:
        client_id: "$(CLIENT_ID)"
        client_secret: "$(CLIENT_SECRET)"
        token_url: "https://auth.example.com/oauth/token"
        scopes:
          - read
          - write
```

## Parámetros

| Parámetro | Requerido | Descripción |
|-----------|-----------|-------------|
| `client_id` | Sí | Identificador del cliente |
| `client_secret` | Sí | Secreto del cliente |
| `token_url` | Sí | URL del endpoint de token |
| `scopes` | No | Lista de permisos (opcional) |

## Notas

- swag2mcp solicita automáticamente un nuevo token cuando el actual expira
- El token se almacena en caché hasta su tiempo de expiración (`expires_in`)
- Si el servidor no proporciona `expires_in`, el token se considera válido por 1 hora
- Todos los parámetros pueden almacenarse en variables de entorno
