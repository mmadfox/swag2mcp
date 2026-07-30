# swag2mcp-format

El skill **swag2mcp-format** enseña a su LLM a mostrar las respuestas de las herramientas MCP de swag2mcp en un formato Markdown limpio, compacto y legible. Sin este skill, el LLM decide por sí mismo cómo formatear las respuestas — a menudo de manera verbosa e inconsistente.

## Lo que cubre

Todas las herramientas MCP de swag2mcp:

- `spec_list`, `spec_by_id` — resumen y detalles de especificaciones
- `collection_by_spec`, `collection_by_id` — colecciones con etiquetas
- `tag_by_spec`, `tag_by_collection`, `tag_by_id` — listados de etiquetas
- `endpoint_by_spec`, `endpoint_by_collection`, `endpoint_by_tag`, `endpoint_by_id` — listas de endpoints
- `search` — resultados de búsqueda
- `inspect` — detalles completos de la operación con esquemas compactos
- `invoke` — resultados de llamadas API
- `auth` — información de autenticación
- `info` — información de ejecución

## Enlace directo

<https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md>

## Instalación a través del agente LLM

Copie esta solicitud en su IDE con IA:

```
Crea el directorio .agents/skills/swag2mcp-format/ y agrega el skill desde
https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md
```

## Instalación manual

```bash
mkdir -p .agents/skills/swag2mcp-format
curl -o .agents/skills/swag2mcp-format/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md
```

## Reinicio requerido

Después de agregar el skill, reinicie su cliente LLM o IDE (ver [Visión general](overview.md#reinicio-requerido)).
