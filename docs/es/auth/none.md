# Ninguna

## Propósito

No se requiere autenticación. La API es accesible sin tokens ni claves.

## Cuándo usarlo

- APIs públicas (Open-Meteo, icanhazdadjoke, PokéAPI)
- Entornos de prueba y demostración
- Cuando la API no requiere autorización

## Configuración

Establezca `type: none` o simplemente omita la sección `auth`:

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: none
```

## Parámetros

Ninguno.

## Notas

- Si la sección `auth` está completamente ausente de la configuración, es equivalente a `type: none`
- No se agregan encabezados de autorización a las solicitudes
