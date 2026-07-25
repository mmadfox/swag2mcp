# Autenticación Básica

## Propósito

Autenticación HTTP Básica — la forma más simple de autenticarse con un nombre de usuario y contraseña.

## Cuándo usarlo

- APIs heredadas que solo admiten Autenticación Básica
- Autenticación simple sin tokens complejos
- Servicios internos

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
      type: basic
      config:
        username: "admin"
        password: "$(PASSWORD)"
```

## Parámetros

| Parámetro | Requerido | Descripción |
|-----------|-----------|-------------|
| `username` | Sí | Nombre de usuario |
| `password` | Sí | Contraseña |

## Notas

- La contraseña se envía en el encabezado `Authorization: Basic ...` codificada en Base64 — esto **no es cifrado**. Use siempre HTTPS.
- Almacene la contraseña en una variable de entorno: `password: "$(MY_PASSWORD)"`
