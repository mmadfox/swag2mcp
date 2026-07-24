# info

## Objectif

Afficher un résumé complet de l'environnement d'exécution swag2mcp au format **JSON**. Cela inclut la version, le chemin de l'espace de travail, le résumé des specs, les paramètres du client HTTP, la configuration du transport MCP, les méthodes d'authentification et l'état du mode mock.

## Quand l'utiliser

- Vous voulez un aperçu lisible par machine de l'espace de travail
- Vous devez vérifier la configuration d'exécution pour le débogage
- Vous voulez voir combien de specs et de points d'accès sont actifs
- Vous devez vérifier les paramètres du client HTTP ou du transport MCP

## Syntaxe

```bash
swag2mcp info [chemin]
```

## Arguments

| Argument | Position | Requis | Description |
|----------|----------|--------|-------------|
| `chemin` | 1 | Non | Répertoire de l'espace de travail. S'il est omis, résolution via les règles de résolution de chemin. |

## Drapeaux

Aucun.

## Comment cela fonctionne

```bash
swag2mcp info
swag2mcp info ./mon-espace-travail
```

## Sortie

La sortie est un objet JSON avec la structure suivante :

```json
{
  "version": "v1.2.0",
  "workspace": "~/.swag2mcp",
  "uptime": "2h 15m",
  "specs": {
    "total": 4,
    "active": 3,
    "disabled": 1,
    "collections": 6,
    "endpoints": 42
  },
  "http_client": {
    "timeout": "30s",
    "max_response_size": "2 Ko",
    "proxy": "aucun",
    "follow_redirects": true,
    "max_redirects": 10,
    "randomize": false
  },
  "mcp": {
    "transport": "stdio",
    "addr": ":8080",
    "path": "/mcp"
  },
  "auth_methods": ["bearer", "api-key"],
  "mock_enabled": false
}
```

## Vérification post-commande

Utilisez `info` pour confirmer que l'espace de travail a été chargé correctement et que toutes les specs sont actives avant de démarrer le serveur MCP.

## Nuances

- **Auto-initialisation :** Si aucun fichier de configuration n'existe, `info` exécute automatiquement l'assistant d'initialisation d'abord.
- **JSON uniquement :** La sortie est toujours en JSON. Pour une sortie lisible par l'humain, utilisez `ls`.
- **`max_response_size` :** Affiché dans un format lisible par l'humain (par exemple, `"1 Ko"`, `"2 Mo"`).
- **Pas d'index de texte intégral :** `info` désactive l'indexation en texte intégral car il n'a besoin que de la configuration et des métadonnées des specs.
