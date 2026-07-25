# Intégration avec Cursor

## stdio

Dans les paramètres de Cursor, ajoutez le serveur MCP :

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

Après la connexion, l'agent IA de Cursor peut :

- Explorer vos API
- Trouver des points de terminaison pertinents
- Appeler des API et afficher les résultats
- Aider à déboguer les requêtes

## Autres

Vous ne voyez pas votre client ? Toutes les intégrations MCP suivent le même modèle :
- Définissez la commande sur `swag2mcp` avec l'argument `mcp`
- Ajoutez éventuellement un chemin d'espace de travail : `mcp /chemin/vers/espace-de-travail`
- Consultez la documentation de votre client pour l'emplacement et le format exacts du fichier de configuration

La plupart des clients MCP prennent en charge le transport stdio, et certains prennent en charge HTTP (SSE / Streamable HTTP).
