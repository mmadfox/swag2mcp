# add

## Propósito

Agregar una nueva **especificación** (servicio de API) o **colección** (archivo OpenAPI/Swagger/Postman) a una configuración existente. Esta es la forma principal de ampliar su espacio de trabajo con nuevas APIs.

## Cuándo usarlo

- Tiene una nueva API para conectar a su agente LLM
- Encontró una URL de especificación OpenAPI y desea agregarla
- Desea agregar un archivo de especificación adicional (colección) a una especificación existente
- Prefiere escribir YAML directamente en lugar de usar el asistente interactivo

## Sintaxis

```bash
swag2mcp add spec [path] [flags]
swag2mcp add collection [path] [flags]
```

## Argumentos

| Argumento | Posición | Requerido | Descripción |
|-----------|----------|-----------|-------------|
| `path` | 1 | No | Directorio del espacio de trabajo. Si se omite, se resuelve mediante reglas de resolución de ruta. |

## Banderas

### `add spec`

| Bandera | Abreviatura | Tipo | Valor predeterminado | Descripción |
|---------|-------------|------|---------------------|-------------|
| `--yaml` | `-y` | `string` | `""` | Entrada YAML en línea o `-` para stdin |
| `--example` | `-e` | `bool` | `false` | Imprimir una plantilla YAML y salir |

### `add collection`

| Bandera | Abreviatura | Tipo | Valor predeterminado | Descripción |
|---------|-------------|------|---------------------|-------------|
| `--yaml` | `-y` | `string` | `""` | Entrada YAML en línea o `-` para stdin |
| `--example` | `-e` | `bool` | `false` | Imprimir una plantilla YAML y salir |

## Cómo funciona

### Modo interactivo (predeterminado)

Inicia un asistente TUI que le permite completar los campos de la especificación o colección paso a paso.

```bash
swag2mcp add spec
swag2mcp add collection
```

### Modo YAML en línea

Pase el YAML directamente como una cadena. **Tenga cuidado con las comillas del shell** — los caracteres especiales como `:`, `#`, `&`, `{` pueden romper el comando.

```bash
swag2mcp add spec --yaml 'domain: meteo
llm_title: Open-Meteo API
base_url: https://meteo.swagger.io/v2
collections:
  - llm_title: Main
    location: https://example.com/spec.json'
```

### YAML desde stdin (recomendado para YAML complejo)

Redirija desde un archivo o use un heredoc para evitar problemas de comillas del shell por completo:

```bash
# Redirigir desde archivo
cat spec.yaml | swag2mcp add spec --yaml -

# Heredoc
swag2mcp add spec --yaml - <<EOF
domain: my-api
llm_title: My API
llm_instruction: "Use this API for X & Y # important"
base_url: https://api.example.com/v1
collections:
  - llm_title: Main
    location: https://raw.githubusercontent.com/org/repo/main/spec.yaml
EOF
```

### Plantilla YAML

Imprimir la estructura YAML esperada y salir:

```bash
swag2mcp add spec --example
swag2mcp add collection --example
```

## Formato YAML

### Especificación

```yaml
domain: meteo
llm_title: Open-Meteo API
llm_instruction: Use this API to manage pets.
base_url: https://meteo.swagger.io/v2
tags: [public, demo]
auth:
  type: bearer
  config:
    token: $(TOKEN)
collections:
  - llm_title: Open-Meteo Swagger
    location: https://example.com/spec.json
```

### Colección

```yaml
spec_domain: meteo
llm_title: Orders Collection
location: https://example.com/orders.json
```

## Verificación posterior al comando

```bash
swag2mcp ls [path]
# La nueva especificación o colección debería aparecer en la lista
```

## Matices

- **Auto-inicio:** Si no existe un archivo de configuración, `add` ejecuta automáticamente el asistente de inicio primero. No necesita ejecutar `init` por separado.
- **Comillas del shell:** El YAML en línea (`--yaml '...'`) es frágil con caracteres especiales. Prefiera `--yaml -` con un heredoc o tubería para cualquier cosa más allá de valores simples.
- **`--example` sale inmediatamente** sin verificar una configuración existente ni modificar nada.
- **`add spec` vs `add collection`:** Use `add spec` para un nuevo servicio de API (nuevo dominio). Use `add collection` para agregar otro archivo de especificación a una especificación existente.
