# Intégration avec VS Code

## Via les paramètres VS Code

Dans `.vscode/settings.json` :

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

## Via une extension

Installez l'extension MCP pour VS Code et ajoutez :

```json
{
  "mcp.servers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp"]
    }
  }
}
```

## Utilisation

Après la configuration, l'assistant IA de VS Code peut travailler avec vos API via swag2mcp.

## Autres

Vous ne voyez pas votre client ? Toutes les intégrations MCP suivent le même modèle :
- Définissez la commande sur `swag2mcp` avec l'argument `mcp`
- Ajoutez éventuellement un chemin d'espace de travail : `mcp /chemin/vers/espace-de-travail`
- Consultez la documentation de votre client pour l'emplacement et le format exacts du fichier de configuration

La plupart des clients MCP prennent en charge le transport stdio, et certains prennent en charge HTTP (SSE / Streamable HTTP).
