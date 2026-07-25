# Colecciones

Una colección es un único archivo OpenAPI/Swagger/Postman que describe una API específica. Apunta a una `location` (URL o ruta de archivo local) y pertenece a una especificación (dominio).

Una especificación puede tener múltiples colecciones — por ejemplo, la especificación "meteo" podría tener colecciones "Pronóstico", "Calidad del Aire" y "Marino", cada una apuntando a un archivo de especificación diferente.

## Campos de la Colección

| Campo | Clave YAML | Requerido | Descripción |
|-------|------------|-----------|-------------|
| [Título LLM](#instrucción-llm) | `llm_title` | ❌ | Nombre mostrado de la colección para el LLM (máx. 120 caracteres). Se auto-rellena desde el documento de especificación si no se establece |
| [Instrucción LLM](#instrucción-llm) | `llm_instruction` | ❌ | Sugerencia corta para el LLM (máx. 360 caracteres). Se auto-rellena desde el documento de especificación si no se establece |
| Título | `title` | ❌ | Anulación del título original de la especificación (se auto-rellena desde el documento analizado) |
| [Ubicación](#ubicación--cómo-se-resuelven-los-archivos-de-especificación) | `location` | ✅ | URL o ruta al archivo de especificación (5–250 caracteres) |
| [Deshabilitar](#deshabilitar) | `disable` | ❌ | Omitir esta colección durante la carga |
| [Cliente HTTP](#anulación-del-cliente-http) | `http_client` | ❌ | Configuración HTTP por colección (encabezados, cookies) |
| [URL Base](#anulación-de-la-url-base) | `base_url` | ❌ | Anular la URL base de la especificación para esta colección |
| [Servidor Simulado](#servidor-simulado) | `base_mock_url` | ❌ | Dirección del servidor simulado en formato `host:port`. Requerido cuando `mock_enabled: true` |

## Ubicación — Cómo se Resuelven los Archivos de Especificación

El campo `location` le indica a swag2mcp dónde encontrar el archivo OpenAPI/Swagger/Postman. Admite varios tipos de origen:

| Origen | Ejemplo | Descripción |
|--------|---------|-------------|
| **URL remota** | `https://raw.githubusercontent.com/.../spec.yaml` | Descargado y almacenado en caché |
| **Archivo local (absoluto)** | `/home/user/my-api.yaml` | Leído del sistema de archivos, almacenado en caché |
| **Archivo local (relativo)** | `./my-api.yaml` | Resuelto a ruta absoluta, almacenado en caché |
| **Archivo local del espacio de trabajo** | `specs/my-api.yaml` | Almacenado en `~/.swag2mcp/specs/`, usado directamente (no en caché) |
| **URI file://** | `file:///home/user/spec.yaml` | Convertido a ruta local, almacenado en caché |

swag2mcp detecta automáticamente el tipo de origen:

- `https://` o `http://` → URL remota (en caché)
- `file://` → archivo local (convertido a ruta del sistema de archivos)
- Todo lo demás → archivo local (con expansión de `~` para el directorio personal)

### URL remotas

Cuando usa una URL remota, swag2mcp descarga el archivo y lo almacena en caché localmente. La caché se reutiliza en inicios posteriores para evitar descargas repetidas.

### Archivos locales

Los archivos locales se leen directamente del sistema de archivos. Si el archivo está fuera del directorio `specs/` del espacio de trabajo, se copia a la caché para mantener la consistencia.

### Archivos locales del espacio de trabajo

El directorio `specs/` dentro del espacio de trabajo (`~/.swag2mcp/specs/`) es el lugar recomendado para archivos de especificación locales. Los archivos almacenados aquí se usan directamente sin almacenamiento en caché. Use una ruta relativa que comience con `specs/` para referenciarlos.

> **Nota:** `specs/` es solo un nombre de directorio (como `cache/` o `responses/`), no el concepto de "especificación". Almacena los archivos OpenAPI/Swagger/Postman reales a los que apuntan las colecciones.

```bash
# Importar un archivo de especificación al espacio de trabajo
swag2mcp import https://example.com/api.yaml myspec

# Después de la importación, la ubicación se convierte en:
# specs/myspec.yaml
```

## Sistema de Caché

swag2mcp almacena en caché los archivos de especificación remotos para evitar descargarlos en cada inicio.

### Cómo funciona

1. Cuando se carga una colección con una URL remota, swag2mcp verifica la caché
2. Si existe una entrada de caché válida (no expirada), se usa directamente
3. Si no, el archivo se descarga, se analiza y se almacena en la caché

### Estructura de la caché

```
~/.swag2mcp/
  cache/
    {sha256_hash}.spec    # Contenido del archivo de especificación en caché
    {sha256_hash}.meta    # Metadatos de la caché (JSON)
```

Cada archivo en caché tiene un archivo de metadatos que contiene:

```json
{
  "source": "https://example.com/api.yaml",
  "source_type": "url",
  "cached_at": "2024-01-01T00:00:00Z",
  "mod_time": "2024-01-01T00:00:00Z",
  "ttl_sec": 3600
}
```

### TTL de la caché

Cada archivo en caché recibe un **TTL aleatorio** entre 1 hora y 48 horas. Esto evita que todos los archivos en caché expiren al mismo tiempo (problema de la avalancha).

### Clave de caché

La clave de caché es un hash SHA-256 de la cadena de ubicación sin procesar (primeros 16 bytes = 32 caracteres hexadecimales).

### Gestión de la caché

```bash
# Limpiar caché y respuestas, volver a descargar todos los archivos de especificación
swag2mcp update

# Limpiar solo caché y respuestas
swag2mcp clean
```

- `swag2mcp update` — valida la configuración, limpia `cache/` y `responses/`, luego vuelve a almacenar en caché todas las ubicaciones de colecciones
- `swag2mcp clean` — elimina todo el contenido de `cache/` y `responses/`, más los scripts de autenticación huérfanos
- Las respuestas antiguas se limpian automáticamente después de 48 horas al iniciar el servidor MCP

## Validación

Cada colección se valida cuando se carga la configuración. La validación se ejecuta en cada inicio de `swag2mcp mcp`. Si falla, el servidor MCP no se iniciará — en algunos IDEs esto significa que el servidor simplemente no se conectará, y el LLM recibe un mensaje de error claro explicando qué corregir.

| Verificación | Regla |
|--------------|-------|
| **Ubicación** | Requerida, 5–250 caracteres |
| **Accesibilidad de la ubicación** | Debe ser una URL accesible o un archivo existente |
| **Validez de la ubicación** | Debe ser un archivo OpenAPI 3.x, Swagger 2.0 o Postman válido |
| **Título LLM** | Máx. 120 caracteres, letras/dígitos/puntuación básica |
| **Instrucción LLM** | Máx. 360 caracteres, mismo conjunto de caracteres que el título |
| **URL Base** | Debe ser una URL válida si se establece |
| **URL Base Simulada** | Debe ser `host:port` o `host:port/path` donde host es `localhost`, `127.0.0.1` o `0.0.0.0` |
| **Simulado requerido** | Si `mock_enabled: true`, cada colección debe tener `base_mock_url` |
| **Puertos simulados duplicados** | No dos colecciones pueden compartir el mismo puerto simulado |

Para diagnosticar problemas antes de iniciar el servidor, use el comando [`validate`](../cli/validate.md):

```bash
# Validar espacio de trabajo predeterminado (~/.swag2mcp)
swag2mcp validate

# Validar un espacio de trabajo de proyecto personalizado
swag2mcp validate ./my-project
```

## Agregar Colecciones

### Mediante configuración YAML

Edite `~/.swag2mcp/swag2mcp.yaml` directamente:

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

Después de editar, reinicie el servidor MCP (`swag2mcp mcp`) para que los cambios surtan efecto.

### Mediante CLI

```bash
# Modo interactivo
swag2mcp add collection

# No interactivo con YAML
swag2mcp add collection --yaml 'spec_domain: meteo
llm_title: Forecast
location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml'

# Redirigir desde stdin
cat collection.yaml | swag2mcp add collection --yaml -

# Mostrar ejemplo YAML
swag2mcp add collection --example
```

### Mediante Importación

```bash
# Importar un archivo de especificación al espacio de trabajo
swag2mcp import https://example.com/api.yaml
```

## Instrucción LLM

Las colecciones pueden tener su propia `llm_instruction` (hasta 360 caracteres) para una guía más específica. Esto se inyecta en el prompt del sistema de swag2mcp junto con la instrucción a nivel de especificación.

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        llm_instruction: "Use esta colección para el clima actual y pronósticos diarios."
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        llm_instruction: "Use esta colección para el índice de calidad del aire y datos de contaminación."
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
```

Si `llm_title` no se establece, se rellena automáticamente desde el campo `title` del documento de especificación. Si `llm_instruction` no se establece, se rellena desde el campo `description` del documento de especificación.

## Deshabilitar

Establezca `disable: true` para omitir una colección. No se cargará, indexará ni estará disponible para el LLM.

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
        disable: true
```

## Anulación de la URL Base

Cada colección puede anular la `base_url` de la especificación. Esto es útil cuando diferentes colecciones dentro de la misma especificación usan diferentes endpoints de API.

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
      - llm_title: Marine
        base_url: https://marine-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
```

## Anulación del Cliente HTTP

Las colecciones pueden anular la configuración HTTP (encabezados, cookies) de los niveles de especificación y global.

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        http_client:
          headers:
            X-API-Version: "2"
          cookies:
            - name: session
              value: abc123
```

La configuración en cascada: global → especificación → colección. Consulte [Cascada de Configuración](../configuration/cascade.md) para más detalles.

## Servidor Simulado

Cuando `mock_enabled: true` se establece a nivel de configuración, cada colección debe tener `base_mock_url` establecido. Esto le indica a swag2mcp dónde se está ejecutando el servidor simulado para esta colección.

```yaml
mock_enabled: true
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        base_mock_url: localhost:8080
```

Consulte [Servidor Simulado](../advanced/mock-server.md) para más detalles.

## Ejemplos

### Colección Mínima

```yaml
specs:
  - domain: dadjokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

### Colección Completa con Todos los Campos

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        llm_instruction: "Use para el clima actual y pronósticos diarios."
        title: "Custom Title"
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        disable: false
        base_url: https://forecast-api.open-meteo.com
        base_mock_url: localhost:8080
        http_client:
          headers:
            X-Custom: value
```

### Múltiples Colecciones por Especificación

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
      - llm_title: Marine
        base_url: https://marine-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
```

### Archivo Local en el Espacio de Trabajo (Directorio specs/)

```yaml
specs:
  - domain: myapi
    llm_title: My Internal API
    base_url: https://api.mycompany.com
    collections:
      - llm_title: Users
        location: specs/users.openapi.json
      - llm_title: Orders
        location: specs/orders.openapi.json
```

### Colección Deshabilitada

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
        disable: true
```

## Relacionados

- [Configuración de Colecciones (config)](../configuration/collection-settings.md) — referencia YAML completa
- [Cascada de Configuración](../configuration/cascade.md) — cómo las configuraciones se anulan entre sí
- [Especificaciones](./specs) — contenedores lógicos para colecciones
- [Cliente HTTP](../configuration/http-client.md) — configuración del cliente HTTP
- [Servidor Simulado](../advanced/mock-server.md) — configuración del servidor simulado
- [CLI: validate](../cli/validate.md) — referencia del comando validate
- [CLI: update](../cli/update.md) — referencia del comando update
- [CLI: clean](../cli/clean.md) — referencia del comando clean
