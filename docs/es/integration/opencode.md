# Integración con OpenCode

## stdio

```json
{
  "mcp": {
    "swag2mcp": {
      "type": "local",
      "command": ["swag2mcp", "mcp"],
      "enabled": true
    }
  }
}
```

## HTTP

```json
{
  "mcp": {
    "swag2mcp": {
      "type": "local",
      "command": ["swag2mcp", "mcp", "--transport", "sse", "--http-addr", "127.0.0.1:8080"],
      "enabled": true
    }
  }
}
```

## Consultas de Ejemplo

Una vez conectado, puede preguntar:

- "¿Qué APIs tienes?"
- "Muestra todos los endpoints en petstore"
- "Encuentra una API para crear un usuario"
- "Llama a GET /pet/1 y muestra el resultado"

## Otros

¿No ve su cliente? Todas las integraciones MCP siguen el mismo patrón:
- Establezca el comando a `swag2mcp` con el argumento `mcp`
- Opcionalmente agregue una ruta de espacio de trabajo: `mcp /path/to/workspace`
- Consulte la documentación de su cliente para la ubicación y formato exactos del archivo de configuración

La mayoría de los clientes MCP admiten el transporte stdio, y algunos admiten HTTP (SSE / HTTP Streamable).
