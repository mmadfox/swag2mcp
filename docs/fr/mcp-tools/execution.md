# Outils d'exécution

Les outils d'exécution sont le cœur de swag2mcp : **search** trouve les points de terminaison lorsque vous n'avez pas d'ID, **inspect** révèle le contrat OpenAPI complet et **invoke** effectue l'appel API réel. Utilisez-les toujours dans cet ordre : search → inspect → invoke.

---

## search

### Objectif

Le seul outil pour trouver des points de terminaison lorsque vous n'avez pas d'ID de point de terminaison. Effectue une recherche en texte intégral dans tous les points de terminaison de toutes les spécifications en utilisant le moteur de recherche bluge.

### Quand l'utiliser

- Lorsque vous ne connaissez pas l'ID du point de terminaison
- Lorsque vous souhaitez trouver des points de terminaison par mots-clés, méthode, balise ou chemin
- Lorsque vous devez découvrir quels points de terminaison existent pour une fonctionnalité spécifique

### Fonctionnement

Recherche dans l'index en texte intégral de toutes les spécifications. Prend en charge les requêtes structurées avec des filtres de champ, des opérateurs booléens, la recherche floue, les caractères génériques, etc.

### Paramètres

| Paramètre | Type | Obligatoire | Description |
|-----------|------|----------|-------------|
| `query` | string | Oui | Requête de recherche (prend en charge la syntaxe structurée) |
| `limit` | int | Oui | Nombre maximal de résultats à renvoyer (1-50) |

### Syntaxe de requête

| Exemple | Description |
|---------|-------------|
| `pet` | Recherche textuelle simple dans tous les champs |
| `method:GET` | Filtrer par méthode HTTP |
| `tag:pet` | Filtrer par nom de balise |
| `path:"/api/v1/users"` | Recherche de chemin exact |
| `+method:POST +tag:pet` | Doit correspondre aux deux conditions |
| `-method:DELETE` | Exclure les méthodes DELETE |
| `create~` | Recherche floue (tolérante aux fautes de frappe) |
| `path:/api/v1/*` | Recherche par chemin avec caractère générique |
| `/pattern/` | Recherche par expression régulière |
| `term^3` | Augmenter la pertinence d'un terme |

**Champs consultables :** `method` (mot-clé), `tag` (mot-clé), `path` (texte), `summary` (texte), `_all` (champ de texte par défaut).

**Non pris en charge :** parenthèses pour le regroupement, opérateurs explicites `AND`/`OR`, regroupement de champs.

### Réponse

```json
{
  "endpoints": [
    {
      "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
      "tagId": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "tagName": "forecast",
      "collectionId": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
      "collectionTitle": "Weather Forecast",
      "specId": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
      "specDomain": "meteo",
      "method": "GET",
      "path": "/v1/forecast",
      "summary": "Obtenir les prévisions météorologiques pour un lieu"
    }
  ]
}
```

Chaque résultat inclut la filiation complète (spécification → collection → balise) afin que le LLM puisse naviguer vers les points de terminaison connexes.

### Nuances

- `limit` doit être compris entre 1 et 50 (renvoie `validation_failed` sinon)
- `query` est obligatoire (renvoie `validation_failed` si vide)
- Les résultats sont renvoyés par ordre de pertinence (meilleure correspondance en premier)
- Utilisez des filtres de champ (`method:GET`, `tag:pet`) pour affiner les résultats
- Pour une correspondance exacte de chemin, utilisez des guillemets : `path:"/v1/forecast"`

---

## inspect

### Objectif

Récupérer l'objet complet de l'opération OpenAPI pour un point de terminaison : tous les paramètres, le schéma du corps de la requête, les schémas de réponse, l'URL de base et l'URL complète. C'est l'outil à appeler **avant** `invoke` pour comprendre le contrat du point de terminaison.

### Quand l'utiliser

- Toujours avant `invoke` — vous avez besoin du contrat complet pour effectuer un appel correct
- Lorsque vous devez expliquer les détails techniques d'une API à l'utilisateur
- Lorsque vous devez connaître les paramètres obligatoires, la structure du corps de la requête ou le format de réponse

### Fonctionnement

Recherche le point de terminaison dans l'index et renvoie l'objet complet de l'opération OpenAPI avec tous les schémas résolus.

### Paramètres

| Paramètre | Type | Obligatoire | Description |
|-----------|------|----------|-------------|
| `endpointId` | string | Oui | Hachage MD5 de 32 caractères du point de terminaison |

### Réponse

```json
{
  "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
  "tagId": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
  "collectionId": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
  "specId": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
  "specDomain": "meteo",
  "method": "POST",
  "path": "/pet",
  "baseUrl": "https://meteo.swagger.io/v2",
  "fullUrl": "https://meteo.swagger.io/v2/pet",
  "operation": {
    "id": "addPet",
    "tags": ["pet"],
    "summary": "Ajouter un nouvel animal",
    "description": "Ajouter un nouvel animal à la boutique",
    "deprecated": false,
    "parameters": [
      {
        "name": "petId",
        "in": "path",
        "description": "ID de l'animal",
        "required": true,
        "schema": {
          "type": "integer",
          "format": "int64"
        }
      }
    ],
    "requestBody": {
      "description": "Objet animal à ajouter",
      "required": true,
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "name": { "type": "string" },
              "status": { "type": "string", "enum": ["available", "pending", "sold"] }
            },
            "required": ["name"]
          }
        }
      }
    },
    "responses": {
      "200": {
        "description": "Opération réussie",
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/Pet"
            }
          }
        }
      },
      "405": {
        "description": "Entrée invalide"
      }
    }
  }
}
```

| Champ | Type | Description |
|-------|------|-------------|
| `baseUrl` | string | URL de base de l'API (depuis la configuration) |
| `fullUrl` | string | URL complète du point de terminaison (base + chemin) |
| `operation.parameters[]` | array | Paramètres avec nom, emplacement (path/query/header/cookie), description, indicateur obligatoire et schéma |
| `operation.requestBody` | object | Corps de la requête avec type de contenu et schéma |
| `operation.responses` | map | Codes de réponse avec descriptions et schémas |
| `operation.deprecated` | bool | Indique si le point de terminaison est obsolète |

### Nuances

- Renvoie `not_found` si le point de terminaison n'existe pas
- C'est le **seul** outil qui renvoie l'opération OpenAPI complète — `endpoint_by_id` renvoie uniquement un résumé
- Appelez toujours `inspect` avant `invoke` pour comprendre les paramètres obligatoires et la structure du corps
- L'objet `operation` inclut les références `$ref` qui sont résolues en leurs définitions de schéma complètes

---

## invoke

### Objectif

Exécuter un véritable appel API vers un point de terminaison. C'est le seul outil qui effectue de véritables requêtes HTTP. L'authentification est appliquée automatiquement — vous n'avez pas besoin d'appeler `auth` d'abord.

### Quand l'utiliser

- Uniquement après avoir appelé `inspect` pour comprendre le contrat du point de terminaison
- Uniquement avec confirmation explicite de l'utilisateur pour les opérations destructrices (POST, PUT, PATCH, DELETE)
- Lorsque l'utilisateur demande d'appeler une API et que vous avez tous les paramètres obligatoires

### Fonctionnement

1. Recherche le point de terminaison dans l'index
2. Substitue les paramètres de chemin dans l'URL
3. Ajoute les paramètres de requête
4. Ajoute les en-têtes et cookies
5. Sérialise le corps de la requête en JSON
6. Obtient et applique automatiquement l'authentification (jeton, en-têtes, paramètres de requête)
7. Effectue la requête HTTP
8. Renvoie la réponse ou l'enregistre dans un fichier si elle est trop volumineuse

### Paramètres

| Paramètre | Type | Obligatoire | Description |
|-----------|------|----------|-------------|
| `endpointId` | string | Oui | Hachage MD5 de 32 caractères du point de terminaison |
| `parameters` | object | Non | Paramètres de chemin, de requête et d'en-tête sous forme de paires clé-valeur |
| `requestBody` | object | Non | Corps de la requête pour les requêtes POST/PUT/PATCH |
| `headers` | object | Non | En-têtes HTTP supplémentaires à envoyer |
| `cookies` | object | Non | Cookies HTTP supplémentaires à envoyer |

### Réponse (en ligne)

```json
{
  "statusCode": 200,
  "headers": {
    "content-type": "application/json"
  },
  "body": {
    "id": 1,
    "name": "Rex",
    "status": "available"
  }
}
```

### Réponse (référence de fichier — lorsque le corps dépasse la limite de taille)

```json
{
  "statusCode": 200,
  "headers": {
    "content-type": "application/json"
  },
  "fileRef": {
    "path": "/Users/utilisateur/.swag2mcp/responses/response_a1b2c3d4.json",
    "size": 1572864,
    "sizeHint": "1.5 Mo",
    "maxSizeHint": "2 Ko",
    "message": "La réponse dépasse la limite de 2 Ko et a été enregistrée sur le disque.",
    "openCmd": "open /Users/utilisateur/.swag2mcp/responses/response_a1b2c3d4.json"
  }
}
```

| Champ | Type | Description |
|-------|------|-------------|
| `statusCode` | int | Code d'état de la réponse HTTP |
| `headers` | object | En-têtes de la réponse HTTP |
| `body` | any | Corps de la réponse (présent lorsqu'il est dans la limite de taille) |
| `fileRef` | object | Référence de fichier (présente lorsque le corps dépasse la limite de taille) |

### Travailler avec de grandes réponses

Lorsque `invoke` renvoie une `fileRef`, utilisez les outils de réponse pour explorer les données :

1. **`response_outline(path)`** — obtenir le résumé structurel (clés, types, longueurs de tableaux)
2. **`response_compress(path, mode)`** — compresser les données pour les adapter en ligne
3. **`response_slice(path, jsonPath)`** — extraire un fragment spécifique

### Nuances

- **L'authentification est automatique :** L'outil `invoke` obtient et applique automatiquement l'authentification à partir de la configuration d'authentification de la spécification. Vous n'avez **pas** besoin d'appeler `auth` d'abord.
- **Limitation de débit :** Chaque point de terminaison a un délai de refroidissement de 10 secondes. Un deuxième appel au même point de terminaison dans les 10 secondes est silencieusement bloqué (renvoie une erreur `rate_limit`).
- **Limite de taille de réponse :** La valeur par défaut est de 2 Ko (configurable via `max_response_size`). Si la réponse dépasse cette limite, elle est enregistrée dans `{workspace}/responses/` et une `FileReference` est renvoyée au lieu du `body` en ligne.
- **Gestion des paramètres :** Les paramètres de chemin sont substitués dans l'URL. Les paramètres de requête sont ajoutés. Les paramètres de la requête remplacent les valeurs par défaut de la spécification de l'opération.
- **Corps de la requête :** Pour POST/PUT/PATCH, le corps est sérialisé en JSON. `Content-Type` est automatiquement défini sur `application/json`.
- **Gestion des erreurs :** Les erreurs HTTP (non-2xx) sont renvoyées sous forme d'`invoke_error` avec le code d'état et le corps de la réponse dans l'indice.
- **Opérations destructrices :** N'invoquez jamais POST/PUT/PATCH/DELETE sans confirmation explicite de l'utilisateur.
