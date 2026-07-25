# Recherche en texte intégral

## Aperçu

swag2mcp inclut un moteur de recherche en texte intégral intégré (bluge) qui indexe tous les points d'accès de toutes les specs. Le LLM peut rechercher des points d'accès par méthode, chemin, résumé ou étiquette — même sans connaître l'ID du point d'accès.

## Comment fonctionne l'indexation

Lorsqu'une spec est ajoutée ou mise à jour, chaque point d'accès est indexé. Les champs suivants sont consultables :

| Champ | Description | Exemple |
|-------|-------------|---------|
| `method` | Méthode HTTP | `GET`, `POST`, `PUT` |
| `path` | Chemin du point d'accès API | `/api/v1/users/{id}` |
| `summary` | Résumé OpenAPI | « Trouver un animal par ID » |
| `tag` | Catégorie du point d'accès | « animaux », « utilisateurs » |
| `_all` | Tous les champs combinés | méthode + chemin + étiquette + résumé |

L'index est reconstruit à chaque démarrage du serveur MCP. Il est stocké en mémoire pour des recherches rapides.

## Syntaxe de requête

La recherche prend en charge une syntaxe de requête riche pour un filtrage précis :

| Exemple | Description |
|---------|-------------|
| `animal` | Recherche textuelle simple dans tous les champs |
| `method:GET` | Trouver tous les points d'accès GET |
| `tag:animaux` | Trouver les points d'accès dans l'étiquette « animaux » |
| `path:"/api/v1/users"` | Correspondance exacte de chemin |
| `+method:POST +tag:animal` | Doit correspondre aux deux conditions |
| `-method:DELETE` | Exclure les méthodes DELETE |
| `create~` | Recherche floue (tolérante aux fautes de frappe) |
| `cr*` | Recherche par joker |
| `"trouver animal"` | Recherche par phrase |
| `+summary:animal -method:DELETE` | Inclure « animal » dans le résumé, exclure DELETE |

### Recherche par champ spécifique

Vous pouvez rechercher dans des champs spécifiques en utilisant la syntaxe `champ:valeur` :

```
method:GET
tag:animaux
path:"/animal/findByStatus"
summary:"trouver animal par statut"
```

### Opérateurs booléens

- `+` — le terme doit correspondre (ET)
- `-` — le terme ne doit pas correspondre (NON)
- Espace entre les termes — OU (n'importe quel terme peut correspondre)

### Recherche floue et par joker

- `terme~` — recherche floue (correspond à des mots similaires, gère les fautes de frappe)
- `te*` — joker (correspond à n'importe quels caractères)
- `te?t` — joker pour un seul caractère

## Exemples

```
# Trouver toutes les requêtes GET
method:GET

# Trouver les requêtes POST dans l'étiquette animal
+method:POST +tag:animal

# Trouver des points d'accès par chemin exact
path:"/animal/findByStatus"

# Trouver par description
"trouver animal par statut"

# Trouver tout sauf DELETE
+summary:animal -method:DELETE

# Recherche floue pour « create » (gère les fautes de frappe)
create~
```

## Outil MCP

L'outil MCP `search` expose le moteur de recherche au LLM :

```
→ search(query: "trouver animal par statut", limit: 5)
← GET /animal/findByStatus — Trouve les animaux par statut
   GET /animal/{animalId} — Trouver un animal par ID
```

### Paramètres

| Paramètre | Requis | Description |
|-----------|--------|-------------|
| `query` | Oui | Requête de recherche (prend en charge la syntaxe structurée) |
| `limit` | Oui | Nombre maximum de résultats (1-50) |

## Notes importantes

- **L'index est en mémoire** — il est reconstruit à chaque démarrage du serveur MCP. Il n'y a pas de fichier d'index persistant.
- **Tous les champs sont en minuscules** — les recherches sont insensibles à la casse
- **La limite est plafonnée à 50** — vous ne pouvez pas demander plus de 50 résultats
- **Une syntaxe de requête invalide** retourne un message d'erreur utile avec des exemples
- **Le champ `_all`** combine méthode, chemin, étiquette et résumé pour les recherches textuelles simples
