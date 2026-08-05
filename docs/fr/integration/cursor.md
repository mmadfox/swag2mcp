# Intégration avec Cursor

## stdio

### Via Settings UI

1. Ouvrir les paramètres de Cursor (Cmd+, / Ctrl+,)
2. Go to **Serveurs MCP**
3. Click **Ajouter un nouveau serveur**
4. Fill in:
   - **Nom:** `swag2mcp`
   - **Type:** `command`
   - **Commande:** `swag2mcp mcp`
5. Click **Enregistrer**

### Via le fichier de configuration

In `~/.cursor/mcp.json`:

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

After connecting, Cursor AI Agent can:

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
