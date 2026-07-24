# Solución de Problemas

## Problemas de Instalación

### swag2mcp: command not found

El binario no está en su PATH.

```bash
# Verifique si Go está instalado
go version

# Encuentre dónde instala Go los binarios
go env GOPATH
# Generalmente ~/go o ~/go/bin

# Agregue al PATH (agregue esto a ~/.zshrc o ~/.bashrc)
export PATH=$PATH:$(go env GOPATH)/bin

# O use la ruta completa
~/go/bin/swag2mcp --version
```

Si descargó un binario desde GitHub Releases, asegúrese de que esté en un directorio que esté en su PATH:

```bash
# Mover a /usr/local/bin (macOS/Linux)
sudo mv swag2mcp /usr/local/bin/
```

### permiso denegado

El binario no tiene permisos de ejecución.

```bash
# Para go install (corregir propiedad)
sudo chown -R $(whoami) $(go env GOPATH)

# Para binario descargado
chmod +x /path/to/swag2mcp
```

### Versión de Go demasiado antigua

swag2mcp requiere Go 1.23+.

```bash
go version
# Si la versión < 1.23, actualice Go:
# https://go.dev/dl/
```

### Servidor simulado no encontrado

El servidor simulado es un binario separado. Instálelo explícitamente:

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

## Problemas de Configuración

### Archivo de configuración no encontrado

swag2mcp no puede encontrar `swag2mcp.yaml`.

```bash
# Cree una nueva configuración
swag2mcp init

# O especifique la ruta explícitamente
swag2mcp mcp /path/to/workspace
swag2mcp ls /path/to/workspace
```

**Causa común:** Ejecutó `swag2mcp mcp` desde un directorio aleatorio y buscó `~/.swag2mcp/` en lugar del espacio de trabajo de su proyecto. Siempre pase la ruta explícitamente.

### Espacio de trabajo incorrecto cargado

swag2mcp cargó un espacio de trabajo diferente al esperado.

**Orden de resolución:** `[path]` explícito → directorio actual (`./`) → `~/.swag2mcp/`. Si ejecuta `swag2mcp mcp` sin una ruta desde un directorio que no tiene `swag2mcp.yaml`, recurre a `~/.swag2mcp/`.

**Solución:** Siempre pase la ruta del espacio de trabajo: `swag2mcp mcp /path/to/your/workspace`

### Error de análisis YAML

El archivo de configuración tiene sintaxis YAML inválida.

```bash
# Valide la configuración
swag2mcp validate

# Errores comunes:
# - Tabuladores en lugar de espacios (YAML requiere espacios)
# - Sangría faltante para campos anidados
# - Cadenas sin comillas con caracteres especiales (: # & {)
```

**Consejo:** Use un linter YAML o un editor con soporte YAML para detectar errores de sintaxis.

### La validación falla: "no specifications defined"

El archivo de configuración existe pero no tiene especificaciones.

```bash
# Agregue una especificación
swag2mcp add spec

# O edite swag2mcp.yaml y agregue al menos una especificación
```

### La validación falla: "duplicate domain"

Dos especificaciones tienen el mismo valor de `domain`. Los dominios deben ser únicos.

```bash
# Liste las especificaciones actuales
swag2mcp ls

# Verifique si hay dominios duplicados en swag2mcp.yaml
```

### La validación falla: "invalid spec location"

La URL o ruta de archivo de `location` no es accesible o no es un archivo de especificación válido.

```bash
# Verifique si la URL es accesible
curl -I https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml

# Verifique si el archivo local existe
ls -la ./specs/my-api.yaml

# Verifique que el archivo sea OpenAPI/Swagger/Postman válido
# (no solo cualquier página JSON o HTML)
```

**Causa común:** El campo `location` apunta al endpoint de la API en sí (por ejemplo, `https://api.example.com/v1/users`) en lugar de la URL del archivo de especificación. La ubicación debe apuntar a un archivo OpenAPI/Swagger/Postman.

## Problemas del Servidor MCP

### Puerto ya en uso

Otro proceso está usando el puerto.

```bash
# Encuentre el proceso
lsof -i :8080

# Mátelo
kill <PID>

# O use un puerto diferente
swag2mcp mcp --transport sse --http-addr :9090
```

### Conexión rechazada

El servidor MCP no está en ejecución o no es accesible.

```bash
# Asegúrese de que el servidor esté en ejecución
swag2mcp mcp --transport sse --http-addr 127.0.0.1:8080

# En otra terminal, verifique el endpoint de salud
curl http://127.0.0.1:8080/health

# Si usa una ruta personalizada
curl http://127.0.0.1:8080/custom-path/health
```

### Las herramientas MCP no aparecen en el cliente LLM

El cliente LLM no puede ver ninguna herramienta.

```bash
# Verifique que las especificaciones estén cargadas
swag2mcp ls

# Verifique que las especificaciones no estén deshabilitadas
swag2mcp validate

# Verifique los registros del servidor
swag2mcp mcp --logfile /tmp/swag2mcp.log
cat /tmp/swag2mcp.log

# Verifique que la ruta del espacio de trabajo en su configuración del IDE sea correcta
# (debe ser una ruta absoluta)
```

**Causas comunes:**
- Ruta de espacio de trabajo incorrecta en la configuración del IDE
- Todas las especificaciones tienen `disable: true`
- Las especificaciones están filtradas por `--tags`
- El archivo de configuración no existe en la ruta especificada

### El protocolo de enlace MCP falla (transporte HTTP)

Para los transportes SSE y Streamable HTTP, el protocolo MCP requiere inicialización antes de que funcionen las llamadas a herramientas.

```
Paso 1: POST /mcp → {"method":"initialize", ...}
Paso 2: POST /mcp → {"method":"notifications/initialized"}
Paso 3: POST /mcp → {"method":"tools/list", ...}  ← ahora funciona
```

Asegúrese de que su cliente LLM complete el protocolo de enlace antes de llamar a las herramientas.

### La verificación de salud devuelve 404

La ruta del endpoint de salud puede diferir de la ruta MCP.

```bash
# Endpoint de salud predeterminado
curl http://127.0.0.1:8080/health

# Si cambió la ruta MCP, la salud sigue en /health
# (no se ve afectada por --http-path)
```

### Herramienta auth no disponible

La herramienta MCP `auth` no aparece.

La herramienta `auth` está **deshabilitada por defecto** (`--disable-llm-auth=true`). Esto es intencional para seguridad en producción.

```bash
# Habilite la herramienta auth
swag2mcp mcp --disable-llm-auth=false
```

## Problemas de Autenticación

### 401 No Autorizado

La API rechazó la solicitud debido a credenciales faltantes o inválidas.

```bash
# Verifique que la autenticación esté configurada
swag2mcp info

# Valide la configuración
swag2mcp validate

# Verifique que las variables de entorno estén establecidas
echo $MY_TOKEN

# Verifique que el token no haya expirado (los tokens bearer son estáticos)
```

**Causas comunes:**
- El token falta o está vacío
- La variable de entorno no está establecida
- El token ha expirado (los tokens bearer no se renuevan automáticamente)
- Tipo de autenticación incorrecto configurado

### 403 Prohibido

La API rechazó la solicitud debido a permisos insuficientes.

- El token puede no tener los alcances requeridos
- La clave de API puede no tener acceso a este recurso
- Consulte la documentación de la API para conocer los permisos requeridos

### Endpoint de token OAuth2 inalcanzable

swag2mcp no puede alcanzar la URL del token OAuth2.

```bash
# Verifique el token_url en su configuración
# Verifique que la URL sea correcta y accesible
curl -X POST https://auth.example.com/oauth/token \
  -d "grant_type=client_credentials" \
  -d "client_id=test" \
  -d "client_secret=test"

# Verifique la conectividad de red
# Verifique la configuración del proxy si está detrás de un proxy corporativo
```

### La autenticación Digest falla

swag2mcp no puede completar el protocolo de enlace de autenticación Digest.

- El servidor debe devolver un encabezado `WWW-Authenticate: Digest ...` con una respuesta 401
- El desafío se almacena en caché durante 5 minutos — si el servidor cambia su nonce, espere a que la caché expire
- Verifique que el nombre de usuario y la contraseña sean correctos

### Discrepancia de firma HMAC

La API rechazó la solicitud firmada con HMAC.

- Verifique que `api_key` y `secret_key` sean correctos
- Verifique que la API use firma HMAC-SHA256 estilo Binance
- Algunos exchanges usan métodos de firma diferentes — la autenticación HMAC es específicamente para APIs compatibles con Binance

### La autenticación de script falla

El script de autenticación externo falló.

```bash
# Verifique que el script exista
ls -la ~/.swag2mcp/auth_scripts/my-domain.sh

# Ejecute el script manualmente para probar
sh ~/.swag2mcp/auth_scripts/my-domain.sh

# Verifique el formato de salida del script (debe ser JSON: {"token": "...", "expires_in": 3600})
# Verifique que el script se complete dentro de 30 segundos
# Verifique que el script tenga permisos de ejecución
chmod +x ~/.swag2mcp/auth_scripts/my-domain.sh
```

## Problemas de Búsqueda

### Sin resultados de búsqueda

La búsqueda no devolvió endpoints.

```bash
# Verifique que las especificaciones estén cargadas
swag2mcp ls

# Verifique que las especificaciones no estén deshabilitadas
swag2mcp validate

# Pruebe una consulta más simple
# Pruebe buscar por método: method:GET
# Pruebe buscar por etiqueta: tag:pets

# El índice se reconstruye en cada inicio del servidor MCP
# Si acaba de agregar una especificación, reinicie el servidor
```

### La búsqueda devuelve resultados irrelevantes

La consulta es demasiado amplia o ambigua.

- Use filtros de campo para acotar: `method:GET +tag:pets`
- Use frases exactas: `"find pet by status"`
- Use el parámetro `limit` para obtener resultados más enfocados

## Problemas de Llamadas a la API

### invoke devuelve un error

La llamada a la API falló.

```bash
# Verifique el mensaje de error — incluye el código de estado HTTP
# Errores 4xx: verifique parámetros, autenticación o permisos
# Errores 5xx: el servidor de la API tiene un problema

# Siempre inspeccione el endpoint antes de invocar
inspect(endpointId: "...")

# Verifique que todos los parámetros requeridos estén proporcionados
# Verifique los tipos de parámetros (cadena, número, booleano)
```

### Error de límite de velocidad

El LLM llamó al mismo endpoint demasiado rápido.

Cada endpoint tiene un período de enfriamiento de 10 segundos. Espere antes de llamar de nuevo, o deshabilite el limitador de velocidad:

```yaml
disable_ratelimiter: true
```

### Respuesta demasiado grande (se devolvió fileRef)

La respuesta excedió `max_response_size`.

Esto es normal. Use las herramientas de respuesta para explorar los datos:

```
1. response_outline(path) → entender la estructura
2. response_compress(path, mode: "first_of_array") → obtener una muestra
3. response_slice(path, jsonPath: "data.0") → obtener datos específicos
```

O aumente el límite:

```yaml
http_client:
  max_response_size: 4194304  # 4 MB
```

### Respuestas lentas de la API

La API está tardando demasiado en responder.

```yaml
http_client:
  timeout: 120s  # Aumentar desde el valor predeterminado de 30s
```

## Problemas del Espacio de Trabajo

### swag2mcp init falla: "directory is not empty"

El directorio de destino ya tiene archivos.

```bash
# Use --force para sobrescribir
swag2mcp init --force

# O use un directorio diferente
swag2mcp init ./new-workspace
```

### swag2mcp update falla

Uno o más archivos de especificación no pudieron descargarse.

```bash
# Verifique el mensaje de error para saber qué URL falló
# Verifique que la URL sea accesible
curl -I <url-fallida>

# Verifique la conectividad de red
# Verifique la configuración del proxy
```

### Export no crea ZIP

El argumento `[output]` debe ser una ruta de archivo que termine en `.zip`, no un directorio.

```bash
# Correcto
swag2mcp export /path/to/workspace /path/to/backup.zip

# Incorrecto (no se creará ningún ZIP)
swag2mcp export /path/to/workspace /some/directory
```

### Import falla: "not a valid swag2mcp backup"

El archivo ZIP no fue creado por `swag2mcp export`.

Solo los archivos ZIP creados por `swag2mcp export` pueden importarse. El archivo tiene una estructura interna específica (`swag2mcp.yaml`, `specs/`, `auth_scripts/`).

## Problemas de la TUI

### La TUI no se renderiza correctamente

La terminal es demasiado pequeña o no admite las funciones requeridas.

- Tamaño mínimo de terminal: 80×24 caracteres
- La TUI usa Bubbletea y funciona en la mayoría de terminales modernas
- Intente redimensionar su ventana de terminal
- Pruebe con un emulador de terminal diferente

### La TUI muestra "no specs found"

El espacio de trabajo no tiene especificaciones configuradas.

```bash
# Verifique las especificaciones
swag2mcp ls

# Agregue una especificación
swag2mcp add spec
```

## Problemas del Servidor Simulado

### El servidor simulado no inicia

```bash
# Verifique que mock_enabled: true esté en la configuración
# Verifique que cada colección tenga base_mock_url establecido
# Verifique que los puertos no estén en uso
lsof -i :9090

# Verifique los registros del servidor simulado
swag2mcp-mock mockserver
```

### El servidor simulado devuelve respuestas vacías

El archivo de especificación puede no tener esquemas de respuesta definidos.

- El servidor simulado genera datos a partir de esquemas de respuesta
- Si no se encuentra ningún esquema, devuelve `{}`
- Verifique que su especificación OpenAPI tenga `responses` con `schema` definido

## Problemas de Red

### Falló la conexión del proxy

swag2mcp no puede conectarse a través del proxy configurado.

```bash
# Verifique el formato de la URL del proxy (debe incluir esquema: http://, https://, socks5://)
# Verifique las credenciales del proxy
# Verifique la lista de exclusión — el destino puede estar en la lista de exclusión
# Pruebe el proxy con curl
curl -x http://proxy.company.com:8080 https://api.example.com
```

### Errores TLS/SSL

Falló la verificación del certificado.

- Si usa un certificado autofirmado para el servidor MCP, el cliente debe confiar en él
- Para el servidor simulado con `--tls`, se genera un certificado autofirmado automáticamente
- Para llamadas a la API, swag2mcp usa el almacén de certificados del sistema

## Otros Problemas

### Alto uso de disco

Los directorios de caché y respuestas pueden crecer con el tiempo.

```bash
# Limpiar todo
swag2mcp clean

# Las respuestas antiguas (>48h) se limpian automáticamente al iniciar el servidor MCP
# Los archivos de caché expiran aleatoriamente entre 1 y 48 horas
```

### "command not found" después de go install

El directorio de `go install` no está en su PATH.

```bash
# Encuentre dónde instala Go los binarios
go env GOPATH
# Agregue al PATH
export PATH=$PATH:$(go env GOPATH)/bin
```

### El LLM no usa las herramientas correctamente

El LLM puede necesitar mejores instrucciones o una habilidad de formato.

- Use `llm_instruction` en la configuración de su especificación para describir lo que hace la API
- Considere usar la [habilidad swag2mcp-format](https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md) para un formato de salida consistente
- La calidad de las respuestas del LLM depende del modelo y de las instrucciones que recibe

### ¿Cómo reporto un error?

Abra un issue en [GitHub](https://github.com/mmadfox/swag2mcp/issues) con:
- La versión de swag2mcp (`swag2mcp --version`)
- Su sistema operativo y arquitectura
- El comando exacto que ejecutó
- El mensaje de error completo
- Su archivo de configuración (con los secretos eliminados)
