# Étiquettes

Une étiquette (tag) est une catégorie qui regroupe des points d'accès connexes au sein d'une collection. Les étiquettes peuvent exister ou non — toutes les collections n'en ont pas, et une collection peut avoir n'importe quel nombre d'étiquettes.

Les étiquettes proviennent du fichier OpenAPI/Swagger/Postman lui-même. Il n'y a **aucun paramètre de configuration YAML** pour les étiquettes — vous ne pouvez pas créer, renommer ou supprimer des étiquettes dans `swag2mcp.yaml`. La seule façon de modifier les étiquettes est de modifier le fichier de spécification original.

## Hiérarchie

```
Spec (domaine, par ex. « meteo »)
  └── Collection (fichier de spec, par ex. forecast.yml)
        └── Étiquette « meteo »
              └── GET /forecast
              └── GET /forecast/hourly
        └── Étiquette « alertes »
              └── GET /alerts
```

## Comment les étiquettes sont créées

Les étiquettes sont extraites du document de spécification lors de l'analyse :

**OpenAPI 3.x / Swagger 2.0** — la liste `tags` de chaque opération devient des étiquettes :

```yaml
paths:
  /animal:
    get:
      tags: ["animaux"]
      summary: "Trouver un animal par ID"
    post:
      tags: ["animaux"]
      summary: "Ajouter un nouvel animal"
  /animal/{animalId}/uploadImage:
    post:
      tags: ["images_animal"]
      summary: "Télécharge une image"
```

**Postman** — chaque dossier de premier niveau devient une étiquette. Les dossiers imbriqués utilisent le nom du dernier dossier.

Si un point d'accès n'a pas d'étiquette, il est placé sous une étiquette `"default"`.

## Objectif

Les étiquettes aident le LLM à trouver des groupes de points d'accès connexes. Au lieu de rechercher dans tous les points d'accès d'une collection, le LLM peut d'abord trouver la bonne étiquette, puis lister uniquement les points d'accès qu'elle contient.

## Outils MCP pour les étiquettes

| Outil | Description |
|-------|-------------|
| `tag_by_spec` | Toutes les étiquettes dans une spec entière |
| `tag_by_collection` | Étiquettes dans une collection spécifique |
| `tag_by_id` | Détails de l'étiquette (titre, nombre de méthodes) |
| `endpoint_by_tag` | Points d'accès regroupés sous une étiquette |

## Exemple

```
Requête : "Afficher toutes les étiquettes dans la collection animaux"
→ tag_by_collection(collectionId: "...")
→ Résultat : animaux (5 méthodes), images_animal (1 méthode)
```

## Limitations

- Les étiquettes sont en lecture seule du point de vue de la configuration. Pour ajouter, renommer ou supprimer des étiquettes, modifiez le fichier OpenAPI/Swagger/Postman original et exécutez `swag2mcp update`.
- Les étiquettes ne peuvent pas être filtrées ou désactivées par collection dans la configuration YAML.
