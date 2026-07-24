# Outils MCP

## Aperçu

swag2mcp fournit **19 outils MCP** qui donnent à un agent LLM un accès complet à vos API via le Model Context Protocol. Ces outils couvrent l'ensemble du flux de travail : découvrir les API disponibles, naviguer dans la hiérarchie des spécifications, rechercher et inspecter des points de terminaison, exécuter des appels API et travailler avec de grandes réponses.

### Ce que les outils résolvent

- **Découverte** — le LLM peut trouver des spécifications, des collections et des balises sans connaître les ID à l'avance
- **Navigation** — descendre de la spécification → collection → balise → point de terminaison dans une hiérarchie structurée
- **Recherche** — recherche en texte intégral dans tous les points de terminaison lorsque vous n'avez pas d'ID
- **Inspection** — obtenir l'objet complet de l'opération OpenAPI avant d'effectuer un appel
- **Exécution** — invoquer de véritables appels API avec authentification automatique
- **Gestion des grandes réponses** — structurer, compresser et découper les réponses trop volumineuses pour être incluses en ligne

### Lecture seule ou modifiable

| Type | Nombre | Outils |
|------|-------|-------|
| **Lecture seule** | 17 | Tous les outils de découverte, point de terminaison, recherche, inspection, info et réponse |
| **Modifiable** | 2 | `invoke` (effectue de véritables appels HTTP), `auth` (récupère les jetons) |

Les outils en lecture seule sont marqués avec `ReadOnlyHint=true` et `IdempotentHint=true` dans le protocole MCP, signalant au LLM qu'ils peuvent être appelés sans effets secondaires.

### Gestion des erreurs

Tous les outils renvoient des erreurs sous forme d'objets `LLMError` structurés avec un code lisible par machine et un message lisible qui explique ce qui n'a pas fonctionné et quoi faire ensuite :

| Code d'erreur | Signification |
|------------|---------|
| `validation_failed` | Entrée invalide (mauvais format d'ID, champs obligatoires manquants) |
| `not_found` | Entité non trouvée dans l'index ou l'espace de travail |
| `rate_limit` | Deuxième appel `invoke` dans les 10 secondes sur le même point de terminaison |
| `invoke_error` | Échec de l'appel HTTP, échec du téléchargement |
| `auth_error` | Échec de la récupération du jeton d'authentification |
| `config_error` | Échec du chargement ou de l'enregistrement du fichier de configuration |
| `parse_error` | Échec de l'analyse du fichier de spécification |

## Catégories

| Catégorie | Outils | Description |
|----------|-------|-------------|
| **Découverte** | `spec_list`, `spec_by_id`, `collection_by_spec`, `collection_by_id`, `tag_by_spec`, `tag_by_collection`, `tag_by_id` | Naviguer dans la hiérarchie des spécifications : trouver des spécifications, des collections et des balises |
| **Points de terminaison** | `endpoint_by_spec`, `endpoint_by_collection`, `endpoint_by_tag`, `endpoint_by_id` | Voir les points de terminaison à différents niveaux de la hiérarchie |
| **Exécution** | `search`, `inspect`, `invoke` | Rechercher, inspecter le contrat complet et appeler des API |
| **Utilitaires** | `auth`, `info`, `response_outline`, `response_compress`, `response_slice` | Jetons d'authentification, informations d'exécution et gestion des grandes réponses |
| **Compétences** | [Guide de formatage](/mcp-tools/skills) | Personnaliser l'affichage des réponses des outils |

## Liste complète

| Outil | Description |
|------|-------------|
| `spec_list` | Lister toutes les spécifications API dans l'espace de travail |
| `spec_by_id` | Obtenir des informations détaillées sur une spécification avec ses collections |
| `collection_by_spec` | Lister les collections d'une spécification |
| `collection_by_id` | Obtenir les détails d'une collection avec ses balises |
| `tag_by_spec` | Lister toutes les balises d'une spécification |
| `tag_by_collection` | Lister les balises d'une collection |
| `tag_by_id` | Obtenir les détails d'une balise (ID, titre, nombre de méthodes) |
| `endpoint_by_spec` | Lister tous les points de terminaison d'une spécification |
| `endpoint_by_collection` | Lister les points de terminaison d'une collection |
| `endpoint_by_tag` | Lister les points de terminaison d'une balise |
| `endpoint_by_id` | Résumé rapide d'un point de terminaison (méthode, chemin, résumé) |
| `search` | Recherche en texte intégral dans tous les points de terminaison |
| `inspect` | Détails complets de l'opération OpenAPI (paramètres, schémas) |
| `invoke` | Exécuter un véritable appel API |
| `auth` | Obtenir le jeton d'authentification ou les en-têtes pour une spécification |
| `info` | Informations d'exécution (version, spécifications, configuration) |
| `response_outline` | Résumé structurel d'un fichier de réponse volumineux |
| `response_compress` | Compresser une grande réponse pour l'adapter en ligne |
| `response_slice` | Extraire un fragment d'une grande réponse |

## Hiérarchie de navigation

```
spec_list
  └── spec_by_id(id)
        └── collection_by_spec(specId)
              └── collection_by_id(id)
                    └── tag_by_collection(collectionId)
                          └── tag_by_id(id)
                                └── endpoint_by_tag(tagId)
                                      └── endpoint_by_id(id)
                                            └── inspect(endpointId)
                                                  └── invoke(endpointId)
```

Lorsque vous n'avez pas d'ID, utilisez `search` pour trouver des points de terminaison par requête. Lorsque `invoke` renvoie une `fileRef` (réponse trop volumineuse), utilisez `response_outline` → `response_compress` ou `response_slice` pour explorer les données.
