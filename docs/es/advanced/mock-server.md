# Servidor Simulado

## Descripción General

El servidor simulado genera respuestas de API falsas basadas en sus esquemas OpenAPI. Le permite probar su integración de API sin realizar llamadas HTTP reales. Esto es útil para desarrollo, pruebas de agentes LLM y demostraciones.

El servidor simulado es un **binario separado** — `swag2mcp-mock`. No está incluido en el binario principal `swag2mcp` y debe instalarse por separado.

## Instalación

```bash
# Opción 1: Descargar desde GitHub Releases
# Busque swag2mcp-mock_<version>_<os>_<arch>.tar.gz

# Opción 2: Instalar con Go
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

## Configuración

Habilite el servidor simulado en su configuración:

```yaml
mock_enabled: true
mock_auth:
  oauth2_port: 9090
  digest_port: 9091
  hmac_port: 9092

specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
        base_mock_url: "127.0.0.1:9090"
```

## Parámetros

### mock_enabled

- **Tipo:** `bool`
- **Valor predeterminado:** `false`
- **Efecto:** Cuando es `true`, cada colección activa debe tener `base_mock_url` establecido. El servidor simulado inicia servidores HTTP para cada colección.

### mock_auth

Puertos para servidores de autenticación simulados. Estos simulan endpoints de autenticación OAuth2, Digest y HMAC para que pueda probar APIs autenticadas sin credenciales reales.

| Campo | Valor predeterminado | Descripción |
|-------|---------------------|-------------|
| `oauth2_port` | `9090` | Puerto para el servidor de token OAuth2 simulado |
| `digest_port` | `9091` | Puerto para el servidor de autenticación Digest simulado |
| `hmac_port` | `9092` | Puerto para el servidor de autenticación HMAC simulado |

### base_mock_url (por colección)

- **Tipo:** `string`
- **Requerido:** Sí (cuando `mock_enabled: true`)
- **Formato:** `host:port` (por ejemplo, `localhost:8080`, `127.0.0.1:9000`)
- **Efecto:** Cada colección obtiene su propio servidor HTTP en esta dirección. El servidor responde a todos los endpoints definidos en la especificación con datos generados aleatoriamente.

## Iniciar el servidor simulado

```bash
# Iniciar con configuración predeterminada
swag2mcp-mock mockserver

# Iniciar con TLS
swag2mcp-mock mockserver --tls

# Iniciar con certificado TLS personalizado
swag2mcp-mock mockserver --tls --tls-cert cert.pem --tls-key key.pem
```

### Banderas TLS

| Bandera | Descripción |
|---------|-------------|
| `--tls` | Habilitar TLS con un certificado autofirmado |
| `--tls-cert` | Ruta al archivo de certificado TLS |
| `--tls-key` | Ruta al archivo de clave TLS |

Si `--tls` está establecido sin `--tls-cert` y `--tls-key`, se genera un certificado autofirmado automáticamente para `localhost`.

## Qué hace el servidor simulado

Cuando inicia el servidor simulado, este:

1. **Analiza todos los archivos de especificación** — lee la especificación OpenAPI/Swagger de cada colección
2. **Registra controladores** — crea un controlador HTTP para cada ruta y método definido en la especificación
3. **Genera datos falsos** — responde con datos generados aleatoriamente que coinciden con el esquema de respuesta (tipos, formatos y estructura correctos)
4. **Inicia servidores de autenticación** — simula endpoints de autenticación OAuth2, Digest y HMAC para pruebas

### Probando el simulacro

```bash
# En una terminal:
swag2mcp-mock mockserver

# En otra terminal:
curl http://localhost:8080/pets
# → [{"id":1,"name":"Pet_name","status":"available"}]
```

## Cómo se generan los datos falsos

El servidor simulado genera datos falsos realistas basados en el esquema OpenAPI:

- **Cadenas** — palabras aleatorias, oraciones o valores específicos de formato (email, URL, UUID, fecha, teléfono, etc.)
- **Números** — enteros y flotantes aleatorios dentro del rango especificado
- **Booleanos** — verdadero/falso aleatorio
- **Arreglos** — de 1 a 3 elementos aleatorios
- **Objetos** — todas las propiedades rellenadas con valores aleatorios
- **Enumeraciones** — valor aleatorio de la lista de enumeración
- **Campos anulables** — a veces devuelve `null` (~10% de probabilidad)

## Casos de uso

- **Desarrollo** — pruebe su integración sin acceso real a la API
- **Pruebas de agentes LLM** — verifique que el LLM pueda descubrir, inspeccionar e invocar endpoints
- **Demostraciones** — muestre swag2mcp funcionando sin configurar APIs reales
- **Pruebas de carga** — pruebe el servidor MCP bajo carga sin golpear APIs reales

## Notas importantes

- **Binario separado** — `swag2mcp-mock` no está incluido en el binario principal `swag2mcp`. Instálelo por separado.
- **Cada colección obtiene su propio puerto** — configure `base_mock_url` por colección
- **Los servidores de autenticación simulados son globales** — los servidores OAuth2, Digest y HMAC se ejecutan en los puertos configurados independientemente de cuántas colecciones tenga
- **Los fallos de análisis de especificaciones no son fatales** — si la especificación de una colección no puede analizarse, se omite con una advertencia
- **TLS autofirmado** — cuando usa `--tls` sin certificados, se genera un certificado autofirmado solo para localhost
