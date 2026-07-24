# Autenticación Bearer

## Propósito

Autenticación mediante Token Bearer — el método más común para APIs REST modernas. El token se envía en el encabezado `Authorization: Bearer <token>`.

## Cuándo usarlo

- APIs REST modernas
- JWT (JSON Web Tokens)
- Tokens de acceso OAuth2 (cuando el token ya se ha obtenido)
- Cualquier API que acepte un Token Bearer

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
      type: bearer
      config:
        token: "eyJhbGciOiJIUzI1NiIs..."
```

## Parámetros

| Parámetro | Requerido | Descripción |
|-----------|-----------|-------------|
| `token` | Sí | Token Bearer (JWT, token OAuth2, etc.) |

## Notas

- El token es estático — si expira, debe actualizarlo en la configuración manualmente
- Para renovación automática de tokens, use `oauth2-cc` o `oauth2-pwd`
- Almacene el token en una variable de entorno: `token: "$(API_TOKEN)"`
