# Especificaciones

Una especificación es un contenedor lógico que representa un dominio o servicio de API (por ejemplo, YouTube, Binance, Open-Meteo). Cada especificación tiene un `domain` único, una `base_url`, `auth` opcional y contiene una o más colecciones.

Las [colecciones](./collections) apuntan a archivos OpenAPI/Swagger/Postman — la especificación en sí no es un archivo, es la agrupación alrededor de ellos.

## Dominio — Reglas de Nomenclatura

El `domain` es el identificador único de una especificación. Se utiliza como clave principal en todo el sistema.

| Regla | Restricción |
|-------|-------------|
| Caracteres | Solo `a-z`, `0-9`, `_`, `-` |
| Longitud | 1–60 caracteres |
| Unicidad | **No se permiten duplicados** — dos especificaciones activas no pueden compartir el mismo dominio |

**Ejemplos válidos:** `meteo`, `binance`, `github-api`, `my_service`, `openai-v1`

**Ejemplos inválidos:** `Meteo` (mayúsculas), `my api` (espacio), `my.api` (punto), `a-very-long-domain-name-that-exceeds-sixty-characters` (demasiado largo)

## Campos de la Especificación

| Campo | Clave YAML | Requerido | Descripción |
|-------|------------|-----------|-------------|
| [Dominio](#dominio--reglas-de-nomenclatura) | `domain` | ✅ | Identificador único de API (1–60 caracteres, `a-z0-9_-`) |
| Título LLM | `llm_title` | ✅ | Nombre legible por humanos que el LLM usa para referenciar esta API (5–120 caracteres) |
| [Instrucción LLM](#instrucción-llm) | `llm_instruction` | ❌ | Sugerencia corta inyectada en el prompt del sistema de swag2mcp (máx. 500 caracteres) |
| URL Base | `base_url` | ✅ | URL base para todas las solicitudes de API (URL válida) |
| [Deshabilitar](#deshabilitar) | `disable` | ❌ | Omitir esta especificación durante la carga e indexación |
| [Etiquetas](#etiquetas) | `tags` | ❌ | Etiquetas para filtrado (por ejemplo, `["public", "demo"]`) |
| [Autenticación](#autenticación) | `auth` | ❌ | Configuración de autenticación |
| [Cliente HTTP](#cliente-http) | `http_client` | ❌ | Configuración HTTP por especificación (encabezados, cookies) |
| [Colecciones](./collections) | `collections` | ✅ | Lista de 1–30 colecciones |

## Validación

Cuando swag2mcp valida la configuración, estas reglas se verifican para cada especificación:

| Verificación | Regla |
|--------------|-------|
| **Dominios duplicados** | No dos especificaciones activas pueden compartir el mismo `domain` |
| **Formato de dominio** | Debe coincidir con `^[a-z0-9_-]{1,60}$` |
| **Título LLM** | Requerido, 5–120 caracteres, letras/dígitos/espacios/puntuación básica |
| **Instrucción LLM** | Máx. 500 caracteres, mismo conjunto de caracteres que el título |
| **URL Base** | Requerida, debe ser una URL válida |
| **Colecciones** | Requeridas, 1–30 elementos |
| **Autenticación** | Validada por tipo de autenticación (por ejemplo, bearer requiere `token`, basic requiere `username` + `password`) |
| **Ubicación** | La `location` de cada colección debe ser una URL o ruta de archivo válida (5–250 caracteres) |

La validación se ejecuta en cada inicio de `swag2mcp mcp`. Si falla, el servidor MCP no se iniciará — en algunos IDEs esto significa que el servidor simplemente no se conectará, y el LLM recibe un mensaje de error claro explicando qué corregir.

Para diagnosticar problemas antes de iniciar el servidor, use el comando [`validate`](../cli/validate.md):

```bash
# Validar espacio de trabajo predeterminado (~/.swag2mcp)
swag2mcp validate

# Validar un espacio de trabajo de proyecto personalizado
swag2mcp validate ./my-project
```

## Instrucción LLM

Se recomienda establecer `llm_instruction` en cada especificación — una sugerencia corta (hasta 500 caracteres) que le dice al LLM para qué sirve esta API y cuándo usarla. Esta instrucción se inyecta en el prompt del sistema de swag2mcp, ayudando al LLM a entender el propósito de la especificación sin contexto adicional.

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    llm_instruction: "Use esta API para obtener chistes de papá aleatorios o buscar chistes específicos por palabra clave."
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

Las colecciones también pueden tener su propia `llm_instruction` (hasta 360 caracteres) para una guía más específica.

## Autenticación

La autenticación se configura a nivel de especificación y se aplica a todas sus colecciones. swag2mcp admite 9 métodos de autenticación:

| Método | Tipo YAML | Campos clave |
|--------|-----------|--------------|
| [Ninguna](../auth/none.md) | `none` | — |
| [Básica](../auth/basic.md) | `basic` | `username`, `password` |
| [Bearer](../auth/bearer.md) | `bearer` | `token` |
| [Digest](../auth/digest.md) | `digest` | `username`, `password` |
| [Credenciales de Cliente OAuth2](../auth/oauth2-cc.md) | `oauth2-cc` | `client_id`, `client_secret`, `token_url` |
| [Contraseña OAuth2](../auth/oauth2-pwd.md) | `oauth2-pwd` | `username`, `password`, `client_id`, `token_url` |
| [Clave de API](../auth/api-key.md) | `api-key` | `key`, `value`, `in` (`header` o `query`) |
| [HMAC](../auth/hmac.md) | `hmac` | `api_key`, `secret_key` |
| [Script](../auth/script.md) | `script` | `domain` |

Consulte [Resumen de Autenticación](../auth/overview.md) para más detalles sobre cada método.

## Cliente HTTP

Puede anular la configuración HTTP a nivel de especificación. Esto se aplica a todas las solicitudes realizadas por las colecciones de esta especificación.

```yaml
specs:
  - domain: slow-api
    llm_title: Slow API
    base_url: https://slow-api.example.com
    http_client:
      headers:
        X-API-Version: "2"
      cookies:
        - name: session
          value: abc123
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

La configuración en cascada: global → especificación → colección. Consulte [Cascada de Configuración](../configuration/cascade.md) para más detalles.

## Etiquetas

Las etiquetas le permiten filtrar especificaciones por categoría. Úselas con la bandera `--tags` en `swag2mcp ls` o durante el arranque.

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    tags: ["weather", "public"]
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

```bash
# Listar solo especificaciones etiquetadas como "weather"
swag2mcp ls --tags weather
```

## Deshabilitar

Establezca `disable: true` para omitir una especificación por completo. No se cargará, indexará ni estará disponible para el LLM.

```yaml
specs:
  - domain: old-api
    llm_title: Old API (Deprecated)
    base_url: https://old-api.example.com
    disable: true
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Ejemplos

### Especificación Mínima

```yaml
specs:
  - domain: dadjokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

### Especificación con Autenticación

```yaml
specs:
  - domain: binance
    llm_title: Binance Market Data API
    base_url: https://api.binance.com
    auth:
      type: hmac
      config:
        api_key: $(BINANCE_API_KEY)
        secret_key: $(BINANCE_SECRET_KEY)
    collections:
      - llm_title: Market Data
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/binance.yaml
```

### Especificación con Múltiples Colecciones

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
      - llm_title: Marine
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
```

### Especificación con Instrucción LLM y Etiquetas

```yaml
specs:
  - domain: rickandmorty
    llm_title: Rick and Morty API
    llm_instruction: "Use esta API para obtener información sobre personajes, episodios y ubicaciones del programa Rick and Morty."
    base_url: https://rickandmortyapi.com/api
    tags: ["entertainment", "public"]
    collections:
      - llm_title: Characters
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/rick-and-morty.json
```

## Relacionados

- [Configuración de Especificaciones (config)](../configuration/spec-settings.md) — referencia YAML completa
- [Cascada de Configuración](../configuration/cascade.md) — cómo las configuraciones se anulan entre sí
- [Resumen de Autenticación](../auth/overview.md) — los 9 métodos de autenticación
- [Cliente HTTP](../configuration/http-client.md) — configuración del cliente HTTP
- [Colecciones](./collections) — archivos de especificación dentro de una especificación
