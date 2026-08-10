# OpenID Connect (OIDC Discovery)

## Propósito

Autenticación basada en OpenID Connect Discovery. swag2mcp obtiene el documento de descubrimiento OIDC del emisor, descubre el punto final del token y obtiene un token Bearer mediante el flujo **client credentials**. El token se almacena en caché hasta su expiración.

## Cuándo usarlo

- APIs que exponen un documento de descubrimiento OIDC (`.well-known/openid-configuration`)
- Integraciones servidor a servidor con un `client_id` + `client_secret`
- Cuando la API usa OIDC y desea adquisición automática del token

## Configuración

```yaml
specs:
  - domain: my-api
    llm_title: My API
    base_url: https://api.example.com
    collections:
      - llm_title: Main
        location: https://example.com/spec.yaml
    auth:
      type: openid-connect
      config:
        issuer: "https://auth.example.com"
        client_id: "$(CLIENT_ID)"
        client_secret: "$(CLIENT_SECRET)"
        scopes:
          - openid
          - profile
```

## Parámetros

| Parámetro | Requerido | Descripción |
|-----------|----------|-------------|
| `issuer` | Sí | URL del emisor OIDC. swag2mcp agrega `/.well-known/openid-configuration` para descubrir el punto final del token |
| `client_id` | Sí | Identificador del cliente |
| `client_secret` | Sí | Secreto del cliente |
| `scopes` | No | Lista de permisos (opcional) |

## Notas

- swag2mcp solicita automáticamente un nuevo token cuando el actual expira
- El token se almacena en caché hasta su expiración (`expires_in`)
- Si el servidor no proporciona `expires_in`, el token es válido durante 1 hora
- `issuer`, `client_id` y `client_secret` admiten `$(VAR)` para variables de entorno
- El token se obtiene mediante el flujo **client credentials** (servidor a servidor, sin usuario)
- El documento de descubrimiento debe contener un campo `token_endpoint`
