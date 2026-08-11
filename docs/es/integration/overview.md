# Cualquier cliente MCP

swag2mcp es un **servidor MCP** (Model Context Protocol). Esto significa que funciona con **cualquier cliente MCP** — no solo con los listados en esta sección. Si su editor, IDE o agente admite el protocolo MCP, puede conectarle swag2mcp.

## Patrón universal

Cada cliente MCP usa la misma configuración básica. Agregue swag2mcp como servidor MCP con:

- **Comando:** `swag2mcp`
- **Argumentos:** `mcp` (más una ruta de espacio de trabajo opcional: `mcp /path/to/workspace`)

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

La ubicación exacta y el formato del archivo de configuración (JSON, TOML, ajustes GUI) varían según el cliente — **consulte la documentación MCP de su cliente** sobre dónde colocarlo.

## Transportes

- **stdio** — funciona en todas partes; la mayoría de los clientes MCP lo admiten
- **HTTP (SSE / Streamable HTTP)** — admitido por clientes con opción de transporte HTTP

Consulte la referencia del comando [`mcp`](/es/cli/mcp) para las banderas de transporte.

## Integraciones probadas

| Cliente | Guía |
|---------|------|
| OpenCode | [OpenCode](/es/integration/opencode) |
| Cursor | [Cursor](/es/integration/cursor) |
| Claude Desktop | [Claude Desktop](/es/integration/claude) |
| VS Code | [VS Code](/es/integration/vscode) |
| Crush | [Crush](/es/integration/crush) |

> Si su cliente no está en la lista, eso **no** significa que no sea compatible. Mientras hable el protocolo MCP, use el patrón universal anterior y siga el manual de su cliente sobre la ubicación de la configuración.
