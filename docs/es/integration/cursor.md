# Integración con Cursor

## stdio

En la configuración de Cursor, agregue el servidor MCP:

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

## Uso

Después de conectarse, el Agente IA de Cursor puede:

- Explorar sus APIs
- Encontrar endpoints relevantes
- Llamar APIs y mostrar resultados
- Ayudar a depurar solicitudes

## Otros

¿No ve su cliente? Todas las integraciones MCP siguen el mismo patrón:
- Establezca el comando a `swag2mcp` con el argumento `mcp`
- Opcionalmente agregue una ruta de espacio de trabajo: `mcp /path/to/workspace`
- Consulte la documentación de su cliente para la ubicación y formato exactos del archivo de configuración

La mayoría de los clientes MCP admiten el transporte stdio, y algunos admiten HTTP (SSE / HTTP Streamable).
