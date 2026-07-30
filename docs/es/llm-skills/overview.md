# Skills for LLM

Los skills son archivos Markdown que enseñan a su agente LLM a trabajar con swag2mcp de manera más efectiva. Se cargan como parte del prompt del sistema del agente y le dan al LLM instrucciones precisas para formatear respuestas y entender comandos CLI.

## Skills disponibles

| Skill | Descripción | Descargar |
|-------|-------------|-----------|
| **swag2mcp-format** | Formatea las respuestas de las herramientas MCP en tablas Markdown compactas y legibles | [SKILL.md](https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md) |
| **swag2mcp-cli** | Referencia CLI completa — el LLM conoce cada comando, bandera y opción de configuración | [SKILL.md](https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md) |

## Por qué son importantes los skills

Sin un skill de formato, el LLM decide por sí mismo cómo mostrar los resultados de las herramientas — a menudo verboso e inconsistente. El skill de formato garantiza un estilo limpio y uniforme: tablas compactas para listas, encabezados en línea para detalles y esquemas concisos.

El skill CLI permite al LLM responder preguntas "cómo hago..." sobre comandos de swag2mcp con precisión, sin adivinar.

## Instalación a través del agente LLM

Copie esta solicitud en su IDE con IA (OpenCode, Cursor, Claude Desktop, VS Code, etc.):

```
Agrega los skills de swag2mcp a mi proyecto:

1. Crea el directorio .agents/skills/swag2mcp-format/ y agrega el skill desde https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md
2. Crea el directorio .agents/skills/swag2mcp-cli/ y agrega el skill desde https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md
```

El agente descargará ambos archivos de skills y los colocará en los directorios correctos.

## Instalación manual

Si su cliente LLM no admite la configuración basada en agente, descargue los archivos manualmente:

```bash
mkdir -p .agents/skills/swag2mcp-format
mkdir -p .agents/skills/swag2mcp-cli

curl -o .agents/skills/swag2mcp-format/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md
curl -o .agents/skills/swag2mcp-cli/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md
```

## Configurar su cliente LLM

Cada cliente LLM e IDE tiene su propia forma de instalar skills. El siguiente ejemplo es para **OpenCode** — consulte la documentación de su cliente para el método correcto.

```json
{
  "skills": [
    {
      "name": "swag2mcp-format",
      "sourceURL": "https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md"
    },
    {
      "name": "swag2mcp-cli",
      "sourceURL": "https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md"
    }
  ]
}
```

## Reinicio requerido

**Después de agregar los skills, reinicie su cliente LLM o IDE.** Algunas herramientas cargan los skills solo al inicio. Si los skills no parecen tener efecto, intente:

- **OpenCode**: Reinicie la aplicación o ejecute el comando opencode nuevamente
- **Cursor**: Cierre y vuelva a abrir la ventana (`Cmd+Shift+W` / `Ctrl+Shift+W`)
- **Claude Desktop**: Salga y vuelva a iniciar la aplicación
- **VS Code**: Recargue la ventana (`Ctrl+Shift+P` → "Developer: Reload Window")
