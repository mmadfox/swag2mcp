# Intégration avec Claude Desktop

## stdio

Dans `claude_desktop_config.json` :

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

## Espace de travail personnalisé

```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "/chemin/vers/espace-de-travail"]
    }
  }
}
```

## Utilisation

Après avoir redémarré Claude Desktop, vous pouvez :

- « Affichez-moi la liste de toutes les API »
- « Trouvez le point de terminaison pour créer une commande »
- « Appelez l'API météo pour Paris »

## Autres

Vous ne voyez pas votre client ? Toutes les intégrations MCP suivent le même modèle :
- Définissez la commande sur `swag2mcp` avec l'argument `mcp`
- Ajoutez éventuellement un chemin d'espace de travail : `mcp /chemin/vers/espace-de-travail`
- Consultez la documentation de votre client pour l'emplacement et le format exacts du fichier de configuration

La plupart des clients MCP prennent en charge le transport stdio, et certains prennent en charge HTTP (SSE / Streamable HTTP).
