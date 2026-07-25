# Integración con Crush

## stdio

```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp"]
    }
  }
}
```

## HTTP

```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "--transport", "sse", "--http-addr", "127.0.0.1:8080"]
    }
  }
}
```

## Otros

¿No ve su cliente? Todas las integraciones MCP siguen el mismo patrón:
- Establezca el comando a `swag2mcp` con el argumento `mcp`
- Opcionalmente agregue una ruta de espacio de trabajo: `mcp /path/to/workspace`
- Consulte la documentación de su cliente para la ubicación y formato exactos del archivo de configuración

La mayoría de los clientes MCP admiten el transporte stdio, y algunos admiten HTTP (SSE / HTTP Streamable).
