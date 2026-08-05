# Intégration avec VS Code

## Via .vscode/mcp.json

1. Installez l'extension MCP pour VS Code (par exemple, MCP Client par org.mcp ou similaire).
2. Créez `.vscode/mcp.json` à la racine du projet :

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

> // "${{workspaceFolder}}" sera passé comme chemin d'espace de travail

3. Rechargez la fenêtre VS Code (Ctrl+Maj+P → "Reload Window").
4. Utilisez l'assistant IA — il connaîtra désormais vos API.

## Alternative : Via les paramètres VS Code

Vous pouvez également configurer dans `.vscode/settings.json` :

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

## Utilisation

Après la configuration, l'assistant IA de VS Code peut travailler avec vos API via swag2mcp.

## Autres

Vous ne voyez pas votre client ? Toutes les intégrations MCP suivent le même modèle :
- Définissez la commande sur `swag2mcp` avec l'argument `mcp`
- Ajoutez éventuellement un chemin d'espace de travail : `mcp /chemin/vers/espace-de-travail`
- Consultez la documentation de votre client pour l'emplacement et le format exacts du fichier de configuration

La plupart des clients MCP prennent en charge le transport stdio, et certains prennent en charge HTTP (SSE / Streamable HTTP).
