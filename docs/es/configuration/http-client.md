# Cliente HTTP

swag2mcp utiliza un cliente HTTP configurable para todas las llamadas a la API. Estas configuraciones se definen globalmente y pueden anularse a nivel de especificación y colección.

## Configuración

```yaml
http_client:
  timeout: 30s
  max_response_size: 1048576
  user_agent: "swag2mcp-global/1.0"
  follow_redirects: true
  max_redirects: 10
  random: false
  proxy:
    url: ""
    username: ""
    password: ""
    bypass: []
  headers:
    "Accept": "application/json"
  cookies:
    - name: "session"
      value: "abc123"
      domain: ".example.com"
      path: "/"
```

## Tiempo de Espera

Controla cuánto tiempo espera swag2mcp una respuesta de la API antes de rendirse.

- **Tipo:** duración (formato Go: `30s`, `60s`, `2m`)
- **Valor predeterminado:** `30s`
- **Rango:** 1 segundo a 5 minutos
- **Efecto:** Si la API no responde dentro de este tiempo, la solicitud falla con un error de tiempo de espera.
- **Cuándo aumentar:** APIs lentas, cargas útiles grandes, redes no confiables.
- **Cuándo disminuir:** APIs internas, verificaciones de salud, escenarios de fallo rápido.

```yaml
http_client:
  timeout: 60s
```

## Tamaño Máximo de Respuesta

Limita cuán grande puede ser una respuesta antes de que swag2mcp la guarde en disco en lugar de devolverla en línea al LLM.

- **Tipo:** `int` (bytes)
- **Valor predeterminado:** `1048576` (1 MB)
- **Rango:** 256 a 10,485,760 bytes (10 MB)
- **Efecto:** Cuando una respuesta excede este límite, se guarda en `{workspace}/responses/` como un archivo JSON. El LLM recibe una referencia de archivo y puede explorarla con las herramientas `response_outline`, `response_compress` y `response_slice`.
- **Cuándo aumentar:** APIs que devuelven conjuntos de datos grandes (informes, registros, análisis).
- **Cuándo disminuir:** Ventana de contexto LLM limitada, o cuando prefiere acceso basado en archivos para todas las respuestas.

```yaml
http_client:
  max_response_size: 4194304  # 4 MB
```

## Agente de Usuario

El encabezado `User-Agent` enviado con cada solicitud. Algunas APIs requieren un agente de usuario específico o bloquean agentes de usuario conocidos de bots.

- **Tipo:** `string`
- **Valor predeterminado:** `"swag2mcp-global/1.0"`
- **Efecto:** Identifica su aplicación ante el servidor de la API.
- **Cuándo cambiar:** La API requiere un agente de usuario específico, o desea identificar su aplicación para análisis.

```yaml
http_client:
  user_agent: "MyApp/1.0"
```

## Seguir Redirecciones

Controla si swag2mcp sigue automáticamente las redirecciones HTTP (códigos de estado 3xx).

- **Tipo:** `bool`
- **Valor predeterminado:** `true`
- **Efecto:** Cuando es `true`, swag2mcp sigue las redirecciones hasta `max_redirects` veces. Cuando es `false`, la respuesta de redirección se devuelve tal cual.
- **Cuándo deshabilitar:** APIs que redirigen en un bucle, endpoints sensibles a la seguridad donde desea inspeccionar los destinos de redirección manualmente.

```yaml
http_client:
  follow_redirects: false
```

## Redirecciones Máximas

Limita cuántas redirecciones sigue swag2mcp antes de detenerse.

- **Tipo:** `int`
- **Valor predeterminado:** `10`
- **Rango:** 0 a 50
- **Efecto:** Si la API redirige más veces que este límite, la solicitud falla.
- **Cuándo cambiar:** APIs con cadenas de redirección largas, o reducir para fallar más rápido en bucles de redirección.

```yaml
http_client:
  max_redirects: 5
```

## Aleatorizador

Agrega encabezados aleatorios similares a los de un navegador a cada solicitud para evitar la identificación y el bloqueo.

- **Tipo:** `bool`
- **Valor predeterminado:** `false`
- **Efecto:** Cuando es `true`, swag2mcp genera encabezados aleatorios para cada solicitud: `User-Agent` (de un conjunto de cadenas de navegador reales), `Accept`, `Accept-Language`, `Accept-Encoding`, `Cache-Control`. Esto anula la configuración de `user_agent`.
- **Cuándo habilitar:** APIs que bloquean solicitudes basadas en User-Agent o patrones de encabezados, escenarios de scraping.

```yaml
http_client:
  random: true
```

## Proxy

Un servidor proxy actúa como intermediario entre swag2mcp y la API de destino. Todo el tráfico HTTP se enruta a través de él.

**Cuándo podría necesitar un proxy:**
- **Red corporativa** — todo el tráfico saliente debe pasar por un proxy de la empresa
- **Restricciones geográficas** — algunas APIs están bloqueadas por región, un proxy en la región correcta lo evita
- **IP estática** — APIs que requieren lista blanca de IP
- **Anonimato** — ocultar la IP de origen de la API de destino

### URL del Proxy

- **Tipo:** `string`
- **Valor predeterminado:** `""` (sin proxy)
- **Esquemas admitidos:** `http`, `https`, `socks5`, `socks5h`
- **Admite `$(VAR)`:** ✅ resuelto en tiempo de ejecución

| Esquema | Descripción | Caso de uso |
|---------|-------------|-------------|
| `http` | Proxy HTTP para tráfico HTTP | Proxies corporativos, proxy básico |
| `https` | Proxy HTTPS (túnel CONNECT) | Proxies corporativos seguros |
| `socks5` | Proxy SOCKS5 (DNS resuelto localmente) | Propósito general, cualquier protocolo |
| `socks5h` | Proxy SOCKS5 (DNS resuelto en el proxy) | Cuando el proxy tiene mejor resolución DNS |

### Autenticación del Proxy

Si el proxy requiere autenticación, proporcione `username` y `password`:

- **Admite `$(VAR)`:** ✅ resuelto en tiempo de ejecución para los tres campos (`url`, `username`, `password`)

```yaml
http_client:
  proxy:
    url: "http://proxy.example.com:8080"
    username: "proxyuser"
    password: "$(PROXY_PASSWORD)"
```

### Exclusión del Proxy

Una lista de dominios que **no** deben pasar por el proxy. Útil para servicios internos, localhost o APIs que solo son accesibles directamente.

```yaml
http_client:
  proxy:
    url: "http://proxy.example.com:8080"
    bypass:
      - "localhost"
      - "127.0.0.1"
      - "*.internal.company.com"
      - "api.local"
```

La exclusión admite patrones comodín (`*.example.com` coincide con cualquier subdominio).

## Encabezados

Encabezados HTTP personalizados agregados a cada solicitud. Los encabezados se fusionan entre niveles de cascada:

```
Encabezados globales → Encabezados de especificación (fusionados) → Encabezados de colección (fusionados)
```

Los encabezados de colección anulan los encabezados de especificación, que anulan los encabezados globales para la misma clave.

```yaml
http_client:
  headers:
    "Accept": "application/json"
    "Accept-Language": "en-US"
```

Los valores de los encabezados admiten la resolución de `$(ENV_VAR)`.

## Cookies

Cookies enviadas con cada solicitud. Las cookies se fusionan entre niveles de cascada (el nivel inferior anula el global para el mismo nombre de cookie).

```yaml
http_client:
  cookies:
    - name: "session"
      value: "abc123"
      domain: ".example.com"
      path: "/"
      secure: false
      http_only: false
```

### Campos de Cookie

| Campo | Requerido | Descripción |
|-------|-----------|-------------|
| `name` | Sí | Nombre de la cookie |
| `value` | Sí | Valor de la cookie (admite resolución de `$(ENV_VAR)`) |
| `domain` | No | Ámbito del dominio (por ejemplo, `.example.com`) |
| `path` | No | Ámbito de la ruta (por ejemplo, `/`) |
| `secure` | No | Solo enviar sobre HTTPS |
| `http_only` | No | No accesible mediante JavaScript |

## Encabezados Personalizados a Nivel de Especificación

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    http_client:
      headers:
        "Accept": "application/json"
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Cookies a Nivel de Especificación

```yaml
specs:
  - domain: example
    llm_title: Example API
    base_url: https://api.example.com
    http_client:
      cookies:
        - name: "session"
          value: "abc123"
        - name: "csrf"
          value: "$(CSRF_TOKEN)"
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Cascada

Las configuraciones del cliente HTTP se transmiten en cascada de global a especificación a colección. Todas las configuraciones pueden anularse en cada nivel:

```
Global (http_client)
    ↓ anula (todas las configuraciones)
Especificación (specs[].http_client)
    ↓ anula (todas las configuraciones)
Colección (specs[].collections[].http_client)
```

**Todas las configuraciones del cliente HTTP** (tiempo de espera, proxy, agente de usuario, redirecciones, tamaño de respuesta, aleatorizador, encabezados, cookies) pueden anularse tanto a nivel de especificación como de colección.

Consulte [Cascada de Configuración](./cascade) para más detalles.
