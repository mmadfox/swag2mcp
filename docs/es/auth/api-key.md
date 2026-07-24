# Clave de API

## Propósito

Autenticación mediante una clave de API. La clave puede enviarse como un encabezado HTTP o como un parámetro de consulta URL.

## Cuándo usarlo

- Servicios que utilizan claves de API
- Servicios meteorológicos, geodatos, APIs de traducción
- Cuando la API espera una clave en un encabezado (`X-API-Key`) o parámetro de consulta (`?api_key=...`)

## Configuración

### Clave en encabezado

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: api-key
      config:
        key: "X-API-Key"
        in: header
        value: "$(API_KEY)"
```

### Clave en parámetro de consulta

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: api-key
      config:
        key: "api_key"
        in: query
        value: "$(API_KEY)"
```

## Parámetros

| Parámetro | Requerido | Descripción |
|-----------|-----------|-------------|
| `key` | Sí | Nombre del encabezado o parámetro de consulta |
| `in` | Sí | Dónde colocar la clave: `header` o `query` |
| `value` | Sí | El valor de la clave |

## Notas

- En modo `header`, la clave se agrega como un encabezado HTTP
- En modo `query`, la clave se agrega como un parámetro URL
- Almacene el valor en una variable de entorno: `value: "$(MY_API_KEY)"`
