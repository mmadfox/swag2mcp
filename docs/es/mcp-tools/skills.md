# Habilidades

## Personalizar el Formato de Salida

Cada herramienta MCP de swag2mcp devuelve datos JSON estructurados. Cómo se **presentan** estos datos al usuario depende de la habilidad de formato del LLM — y usted puede controlarlo completamente.

### La habilidad de formato predeterminada

swag2mcp incluye una habilidad de formato incorporada que define markdown compacto y legible por humanos para cada respuesta de herramienta:

[swag2mcp-format SKILL.md](https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md)

Esta habilidad cubre las 19 herramientas MCP con:
- Tablas compactas para listas (especificaciones, colecciones, etiquetas, endpoints)
- Encabezados en línea para vistas de detalle
- Representación compacta de esquemas para `inspect`
- Estilo consistente en todas las respuestas

### Por qué las habilidades son importantes

Los mismos datos pueden presentarse de formas radicalmente diferentes dependiendo de la habilidad:

| Estilo | Ejemplo de salida |
|-------|-------------------|
| **Tablas compactas** (predeterminado) | `GET /pet/{petId}` — Find pet by ID |
| **Verboso** | `Method: GET, Path: /pet/{petId}, Summary: Find pet by ID, Deprecated: false` |
| **Mínimo** | `GET /pet/{petId}` |
| **Técnico** | `GET /pet/{petId} → 200: Pet object, 404: Not found` |
| **Personalizado** | Cualquier formato que pueda describir |

### Crear su propia habilidad

Puede escribir su propia habilidad de formato describiendo el formato de salida exacto que desea. La habilidad es un archivo markdown con reglas de formato para cada herramienta. Aquí hay algunas ideas:

- **Salida JSON** — devolver JSON sin procesar para consumo por máquina
- **Estilo CSV** — datos tabulares para importación a hojas de cálculo
- **Amigable para diagramas** — diagramas Mermaid o ASCII de la estructura de la API
- **Mínimo** — solo método y ruta, nada más
- **Estilo documentación** — descripciones completas, ejemplos y notas

### El único límite es el modelo

La calidad de la salida formateada depende completamente de la capacidad del LLM para seguir sus reglas de formato. Una habilidad bien escrita con ejemplos claros produce una salida consistente y confiable. Una habilidad vaga produce resultados inconsistentes.

Puede:
- Usar la habilidad predeterminada tal cual
- Bifurcarla y ajustar el formato a su gusto
- Escribir la suya propia desde cero
- Cambiar entre habilidades dependiendo de la tarea

### Cómo usar una habilidad

Las habilidades son cargadas por el cliente LLM (OpenCode, Cursor, Claude Desktop, etc.) como parte de su prompt del sistema o configuración del agente. Consulte la documentación de su cliente para saber cómo adjuntar un archivo de habilidad.

Para OpenCode, las habilidades se configuran en `opencode.json`:

```json
{
  "skills": [
    {
      "name": "swag2mcp-format",
      "sourceURL": "https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md"
    }
  ]
}
```
