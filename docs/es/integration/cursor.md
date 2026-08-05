# Integración con Cursor

## stdio

### Via Settings UI

1. Abrir configuración de Cursor (Cmd+, / Ctrl+,)
2. Go to **Servidores MCP**
3. Click **Añadir nuevo servidor**
4. Fill in:
   - **Nombre:** `swag2mcp`
   - **Tipo:** `command`
   - **Comando:** `swag2mcp mcp`
5. Click **Guardar**

### Mediante archivo de configuración

In `~/.cursor/mcp.json`:

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

After connecting, Cursor AI Agent can:

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
