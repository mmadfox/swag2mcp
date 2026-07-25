# Flujo de Trabajo CLI

Esta página muestra ejemplos reales de uso de swag2mcp desde la terminal — desde la inicialización hasta las operaciones diarias.

## Inicio rápido

```bash
# 1. Inicializar un espacio de trabajo
mkdir -p .swag2mcp && swag2mcp init ./.swag2mcp

# 2. Listar sus especificaciones
swag2mcp ls
```

## Agregar una especificación con YAML

### Especificación simple (API pública)

```bash
swag2mcp add spec --yaml - <<EOF
domain: meteo
llm_title: Open-Meteo Weather API
base_url: https://api.open-meteo.com
collections:
  - llm_title: Weather Forecast
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
EOF
```

### Especificación con autenticación (token bearer desde env)

```bash
swag2mcp add spec --yaml - <<EOF
domain: my-api
llm_title: My Protected API
base_url: https://api.example.com/v1
auth:
  type: bearer
  config:
    token: \$(MY_TOKEN)
collections:
  - llm_title: Users
    location: https://raw.githubusercontent.com/my-org/my-api/main/users.yaml
EOF
```

### Especificación con múltiples colecciones

```bash
swag2mcp add spec --yaml - <<EOF
domain: meteo
llm_title: Open-Meteo APIs
base_url: https://api.open-meteo.com
collections:
  - llm_title: Forecast
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
  - llm_title: Air Quality
    base_url: https://air-quality-api.open-meteo.com
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
EOF
```

## Agregar una colección a una especificación existente

```bash
swag2mcp add collection --yaml - <<EOF
spec_domain: meteo
llm_title: Marine Weather
location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
EOF
```

## Listar especificaciones

```bash
$ swag2mcp ls
Specifications:
  dadjoke (https://icanhazdadjoke.com)
    jokes (3 endpoints)
  meteo (https://api.open-meteo.com)
    forecast (5 endpoints)
    air-quality (8 endpoints)
    marine (4 endpoints)
```

### Filtrar por etiquetas

```bash
swag2mcp ls --tags=public
```

## Ver información de ejecución

```bash
$ swag2mcp info
{
  "version": "v1.2.0",
  "workspace": "/home/user/.swag2mcp",
  "specs": {
    "total": 2,
    "active": 2,
    "disabled": 0,
    "collections": 4,
    "endpoints": 20
  },
  "http_client": {
    "timeout": "30s",
    "max_response_size": "2 KB",
    "follow_redirects": true
  },
  "mcp": {
    "transport": "stdio"
  },
  "auth": {
    "methods": ["bearer"]
  }
}
```

## Validar configuración

```bash
$ swag2mcp validate
✅ Configuration is valid.
✓ Spec dadjoke: OK
✓ Spec meteo: OK
```

## Iniciar el servidor MCP

### stdio (para integración con IDE)

```bash
swag2mcp mcp
```

### HTTP (para acceso remoto)

```bash
swag2mcp mcp --transport sse --http-addr :8080
```

### Con filtro de etiquetas

```bash
swag2mcp mcp --tags=public
```

## Actualizar especificaciones

Actualizar todos los archivos de especificación en caché:

```bash
swag2mcp update
```

## Limpiar caché

```bash
swag2mcp clean
```

## Exportación e importación

### Copia de seguridad de su espacio de trabajo

```bash
swag2mcp export --output ~/backups/swag2mcp-2026-07-24.zip
```

### Restaurar en otra máquina

```bash
# En la máquina nueva
swag2mcp import --from-zip swag2mcp-2026-07-24.zip
```

## Explorador TUI interactivo

```bash
swag2mcp run
```

Abre una interfaz de terminal de pantalla completa para buscar, navegar e invocar APIs.

## Servidor simulado

```bash
# Instalar el binario simulado
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest

# Iniciar servidores simulados
swag2mcp-mock mockserver
```
