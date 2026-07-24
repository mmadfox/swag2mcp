# Intégration avec OpenCode

## stdio

```json
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

## HTTP

```json
{
  "mcp": {
    "swag2mcp": {
      "type": "local",
      "command": ["swag2mcp", "mcp", "--transport", "sse", "--http-addr", "127.0.0.1:8080"],
      "enabled": true
    }
  }
}
```

## Exemples de requêtes

Une fois connecté, vous pouvez demander :

- « Quelles API avez-vous ? »
- « Affichez tous les points de terminaison de petstore »
- « Trouvez une API pour créer un utilisateur »
- « Appelez GET /pet/1 et affichez le résultat »

## Autres

Vous ne voyez pas votre client ? Toutes les intégrations MCP suivent le même modèle :
- Définissez la commande sur `swag2mcp` avec l'argument `mcp`
- Ajoutez éventuellement un chemin d'espace de travail : `mcp /chemin/vers/espace-de-travail`
- Consultez la documentation de votre client pour l'emplacement et le format exacts du fichier de configuration

La plupart des clients MCP prennent en charge le transport stdio, et certains prennent en charge HTTP (SSE / Streamable HTTP).
