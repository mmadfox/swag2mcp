# Points d'accès

Un point d'accès est une méthode HTTP + chemin spécifique qui peut être invoquée (par exemple, `GET /api/users/{id}`). Les points d'accès sont les opérations API réelles que le LLM découvre, inspecte et appelle.

## Structure

Chaque point d'accès contient :

- **Méthode HTTP** : GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS
- **Chemin** : `/api/v1/users/{id}`
- **Résumé** : une courte description de ce que fait le point d'accès — très utile pour que le LLM comprenne son objectif en un coup d'œil
- **Description** : une explication détaillée du comportement, des paramètres et des cas d'utilisation du point d'accès
- **Paramètres** : chemin, requête, en-tête, cookie
- **Corps de la requête** : pour POST/PUT/PATCH
- **Réponses** : codes d'état et schémas de réponse

Les champs `summary` et `description` proviennent du fichier OpenAPI/Swagger/Postman. Ils sont le principal moyen pour le LLM de comprendre ce que fait un point d'accès. Des résumés bien rédigés rendent la découverte des points d'accès beaucoup plus efficace.

## Outils MCP pour les points d'accès

| Outil | Description |
|-------|-------------|
| `endpoint_by_spec` | Tous les points d'accès dans une spec |
| `endpoint_by_collection` | Points d'accès dans une collection |
| `endpoint_by_tag` | Points d'accès dans une étiquette |
| `endpoint_by_id` | Résumé rapide d'un point d'accès |
| `inspect` | Détails complets du point d'accès (schémas, paramètres) |
| `invoke` | Appeler le point d'accès |
| `search` | Rechercher des points d'accès par texte |

## Points d'accès dépréciés

Les points d'accès marqués comme `deprecated` dans la spec sont affichés avec un avis lors de l'inspection.

## Configuration

Les points d'accès sont **en lecture seule** du point de vue de swag2mcp. Il n'y a pas de paramètres de configuration YAML pour les points d'accès — vous ne pouvez pas ajouter, supprimer, renommer ou modifier des points d'accès dans `swag2mcp.yaml`.

Pour modifier les points d'accès (en ajouter de nouveaux, mettre à jour les résumés, modifier les paramètres, marquer comme dépréciés), modifiez le fichier OpenAPI/Swagger/Postman original et exécutez `swag2mcp update` pour ré-analyser et ré-indexer.

## Exemple

```
Requête : "Afficher les détails pour GET /animal/{animalId}"
→ inspect(endpointId: "abc123...")
→ Résultat :
  GET /animal/{animalId}
  Résumé : Trouver un animal par ID
  Description : Retourne un seul animal par son ID
  Paramètres :
    - animalId (chemin, entier, requis)
  Réponses :
    - 200 : Objet Animal
    - 400 : Erreur
    - 404 : Non trouvé
```
