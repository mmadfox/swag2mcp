# Autenticación Digest

## Propósito

Autenticación de Acceso Digest HTTP — una alternativa más segura a la Autenticación Básica. La contraseña no se envía en texto plano; en su lugar, se utilizan hashes MD5.

## Cuándo usarlo

- APIs heredadas que solo admiten Digest
- Cuando necesita autenticación sin enviar la contraseña en texto plano
- Sistemas empresariales internos

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
      type: digest
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

- swag2mcp primero envía una solicitud sin autenticación, recibe un desafío del servidor (HTTP 401), calcula la respuesta y reintenta con el encabezado `Authorization: Digest ...`
- El desafío se almacena en caché durante 5 minutos — las solicitudes subsiguientes no necesitan un viaje de ida y vuelta adicional
- Almacene la contraseña en una variable de entorno: `password: "$(API_PASSWORD)"`
