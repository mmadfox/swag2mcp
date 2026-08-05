# Integración con VS Code

## Via .vscode/mcp.json

1. Instale la extensión MCP para VS Code (por ejemplo, MCP Client de org.mcp o similar).
2. Cree `.vscode/mcp.json` en la raíz del proyecto:

```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "${workspaceFolder}"]
    }
  }
}
```

> // "${{workspaceFolder}}" se pasará como ruta del espacio de trabajo

3. Recargue la ventana de VS Code (Ctrl+Mayús+P → "Reload Window").
4. Use el asistente de IA — ahora conocerá sus APIs.

## Alternativa: Mediante configuración de VS Code

También puede configurar en `.vscode/settings.json`:

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

Después de la configuración, el asistente de IA de VS Code puede trabajar con sus APIs a través de swag2mcp.

## Otros

¿No ve su cliente? Todas las integraciones MCP siguen el mismo patrón:
- Establezca el comando a `swag2mcp` con el argumento `mcp`
- Opcionalmente agregue una ruta de espacio de trabajo: `mcp /path/to/workspace`
- Consulte la documentación de su cliente para la ubicación y formato exactos del archivo de configuración

La mayoría de los clientes MCP admiten el transporte stdio, y algunos admiten HTTP (SSE / HTTP Streamable).
