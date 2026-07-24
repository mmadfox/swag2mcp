# Inicio Rápido

Ponga swag2mcp en marcha en 2 minutos.

## 1. Inicializar

### Directorio personal (recomendado)

Configuración única para todo su sistema. La configuración se almacena en su carpeta personal.

::: code-group

```bash [macOS / Linux]
swag2mcp init
# Crea ~/.swag2mcp/swag2mcp.yaml
```

```powershell [Windows]
swag2mcp.exe init
# Crea %USERPROFILE%\.swag2mcp\swag2mcp.yaml
```

:::

### Directorio del proyecto

Para un espacio de trabajo aislado dentro de su proyecto.

::: code-group

```bash [macOS / Linux]
mkdir -p ./swag2mcp && swag2mcp init ./swag2mcp
```

```powershell [Windows]
mkdir ./swag2mcp; swag2mcp.exe init ./swag2mcp
```

:::

### Desde ZIP

Si tiene un espacio de trabajo ya preparado (por ejemplo, de un colega):

```bash
swag2mcp import --from-zip workspace.zip
```

## 2. Instalar Habilidades del Agente (recomendado)

Instale las habilidades de swag2mcp para enseñar a su agente de IA todos los comandos, banderas, formato de configuración y ejemplos del mundo real.

Pídale a su agente:

```bash
"Agregue la habilidad swag2mcp-cli desde https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-cli/SKILL.md"
"Agregue la habilidad swag2mcp-format desde https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md"
```

> Algunos IDEs requieren un reinicio después de agregar habilidades.

## 3. Configuración del Cliente LLM / IDE

Configure su IDE para conectarse a swag2mcp. El IDE iniciará el servidor MCP automáticamente cuando sea necesario.

::: code-group

```json [OpenCode]
{
  "mcp": {
    "swag2mcp": {
      "type": "local",
      "command": ["swag2mcp", "mcp"],
      "enabled": true
    }
  }
}
```

```json [Claude Desktop]
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp"]
    }
  }
}
```

```json [Crush]
{
  "mcp": {
    "swag2mcp": {
      "type": "stdio",
      "command": "swag2mcp",
      "args": ["mcp"]
    }
  }
}
```

:::

Para otros IDEs (Cursor, VS Code, JetBrains) consulte la [guía de Integración](../integration/opencode.md).

> Si inicializó el espacio de trabajo en una ruta personalizada (por ejemplo, `./swag2mcp`), use la ruta completa en el comando:
> `"command": ["swag2mcp", "mcp", "/absolute/path/to/swag2mcp"]`

> **Después de cualquier cambio de configuración, reinicie el servidor MCP** para que los cambios surtan efecto.

## 4. Iniciar Servidor MCP

### stdio (predeterminado) — para IDE local

Nada que configurar. Su IDE inicia swag2mcp automáticamente mediante la configuración anterior.

```bash
swag2mcp mcp
```

### SSE / HTTP Streamable — para acceso remoto

```bash
swag2mcp mcp --transport sse --http-addr :8080
```

O configure en `swag2mcp.yaml`:

```yaml
mcp:
  transport: sse
  addr: ":8080"
  path: "/mcp"
```

Consulte la [referencia del Servidor MCP](../configuration/mcp-server.md) para todas las banderas.

### Filtrar especificaciones por etiquetas

```bash
swag2mcp mcp --tags weather,public
```

Solo las especificaciones con etiquetas coincidentes estarán disponibles para el LLM.

### Verificar que funciona

Después de conectarse, pregúntele a su agente LLM:

```bash
"¿Qué herramientas MCP soportas?"
```

Si el agente lista las herramientas de swag2mcp (`spec_list`, `search`, `invoke`, etc.) — todo está funcionando.

### Consultas de ejemplo para probar

| Pregunte a su agente | Qué sucede |
|----------------------|------------|
| "¿Qué tiempo hace en Nueva York?" | `invoke` — llama a la API de pronóstico de Open-Meteo |
| "¿Cuál es el precio actual de BTC?" | `invoke` — llama a la API de ticker de Binance |
| "Cuéntame un chiste de papá" | `invoke` — llama a la API de icanhazdadjoke |
| "Muéstrame a Pikachu" | `invoke` — llama a la API de PokéAPI por nombre |
| "¿Quién es Rick Sanchez?" | `invoke` — llama a la API de personajes de Rick and Morty |
| "¿Cuál es la calidad del aire en Pekín?" | `invoke` — llama a la API de calidad del aire de Open-Meteo |
| "¿Qué tan altas son las olas cerca de Portugal?" | `invoke` — llama a la API marina de Open-Meteo |
| "Busca chistes sobre perros" | `invoke` — llama al endpoint de búsqueda de dadjoke |
| "Lista todos los Pokémon" | `invoke` — llama al endpoint de lista de PokéAPI |
| "¿Cuál es la elevación del Monte Everest?" | `invoke` — llama a la API de elevación de Open-Meteo |

## 5. ¿Qué Sigue?

- [Conceptos](../concepts/overview.md) — entienda la arquitectura
- [Configuración](../configuration/config-file.md) — personalice los ajustes
- [Comandos CLI](../cli/overview.md) — referencia completa de comandos
