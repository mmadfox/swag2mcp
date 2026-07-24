# Integración con Claude Desktop

## stdio

En `claude_desktop_config.json`:

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

## Espacio de Trabajo Personalizado

```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "/path/to/workspace"]
    }
  }
}
```

## Uso

Después de reiniciar Claude Desktop, puede:

- "Muéstrame la lista de todas las APIs"
- "Encuentra el endpoint para crear un pedido"
- "Llama a la API del clima para Moscú"

## Otros

¿No ve su cliente? Todas las integraciones MCP siguen el mismo patrón:
- Establezca el comando a `swag2mcp` con el argumento `mcp`
- Opcionalmente agregue una ruta de espacio de trabajo: `mcp /path/to/workspace`
- Consulte la documentación de su cliente para la ubicación y formato exactos del archivo de configuración

La mayoría de los clientes MCP admiten el transporte stdio, y algunos admiten HTTP (SSE / HTTP Streamable).
