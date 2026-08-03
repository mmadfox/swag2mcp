# swag2mcp-cli

El skill **swag2mcp-cli** le da a su LLM la referencia completa de la CLI de swag2mcp — cada comando, bandera, argumento y opción de configuración. Con este skill, el LLM puede responder preguntas "cómo hago..." con precisión, sin adivinar.

## Lo que cubre

Los 13 comandos CLI:

| Comando | Propósito |
|---------|-----------|
| `init` | Inicializar un espacio de trabajo y configuración |
| `add` | Agregar una especificación o colección |
| `delete` | Eliminar una especificación o colección |
| `ls` | Listar especificaciones configuradas |
| `run` | Iniciar el explorador API TUI |
| `validate` | Validar el archivo de configuración |
| `clean` | Limpiar datos en caché |
| `update` | Actualizar caché desde la configuración |
| `mcp` | Iniciar el servidor MCP |
| `version` | Mostrar información de versión |
| `info` | Mostrar información de ejecución |
| `import` | Importar archivos de especificación o restaurar espacio de trabajo desde copia de seguridad |
| `export` | Exportar el espacio de trabajo a un archivo ZIP |

Más todas las banderas, estructura del archivo de configuración, métodos de autenticación y opciones avanzadas.

## Enlace directo

<https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md>

## Instalación a través del agente LLM

Copie esta solicitud en su IDE con IA:

```
Crea el directorio .agents/skills/swag2mcp-cli/ y agrega el skill desde
https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md
```

## Instalación manual

```bash
mkdir -p .agents/skills/swag2mcp-cli
curl -o .agents/skills/swag2mcp-cli/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md
```

## Reinicio requerido

Después de agregar el skill, reinicie su cliente LLM o IDE (ver [Visión general](overview.md#reinicio-requerido)).
