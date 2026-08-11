# Tout client MCP

swag2mcp est un **serveur MCP** (Model Context Protocol). Cela signifie qu'il fonctionne avec **tout client MCP** — pas seulement ceux listés dans cette section. Si votre éditeur, IDE ou agent prend en charge le protocole MCP, vous pouvez y connecter swag2mcp.

## Modèle universel

Chaque client MCP utilise la même configuration de base. Ajoutez swag2mcp comme serveur MCP avec :

- **Commande :** `swag2mcp`
- **Arguments :** `mcp` (plus un chemin d'espace de travail facultatif : `mcp /path/to/workspace`)

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

L'emplacement exact et le format du fichier de configuration (JSON, TOML, paramètres GUI) varient selon le client — **consultez la documentation MCP de votre client** pour savoir où le placer.

## Transports

- **stdio** — fonctionne partout ; la plupart des clients MCP le prennent en charge
- **HTTP (SSE / Streamable HTTP)** — pris en charge par les clients avec une option de transport HTTP

Voir la référence de la commande [`mcp`](/fr/cli/mcp) pour les indicateurs de transport.

## Intégrations testées

| Client | Guide |
|--------|-------|
| OpenCode | [OpenCode](/fr/integration/opencode) |
| Cursor | [Cursor](/fr/integration/cursor) |
| Claude Desktop | [Claude Desktop](/fr/integration/claude) |
| VS Code | [VS Code](/fr/integration/vscode) |
| Crush | [Crush](/fr/integration/crush) |

> Si votre client n'est pas dans la liste, cela ne signifie **pas** qu'il n'est pas pris en charge. Tant qu'il parle le protocole MCP, utilisez le modèle universel ci-dessus et suivez le manuel de votre client pour l'emplacement de la configuration.
