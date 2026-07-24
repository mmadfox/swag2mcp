# Outils utilitaires

Les outils utilitaires fournissent des fonctionnalités de support : récupération des jetons d'authentification, obtention d'informations d'exécution et travail avec de grandes réponses API qui ne tiennent pas en ligne.

---

## auth

### Objectif

Récupérer un jeton d'authentification, des en-têtes ou des paramètres de requête pour une spécification spécifique. Cela donne au LLM un accès aux identifiants qui peuvent être utilisés en dehors de swag2mcp (par exemple, pour générer une commande curl).

### Quand l'utiliser

- Uniquement lorsque l'utilisateur demande explicitement le jeton brut ou les identifiants
- Lors de la génération d'une commande curl ou d'un extrait de code nécessitant une authentification
- Lorsque l'utilisateur souhaite voir quelle méthode d'authentification est configurée

### Quand NE PAS l'utiliser

- **N'appelez pas** `auth` avant `inspect` ou `invoke` — `invoke` obtient et applique automatiquement l'authentification
- **N'appelez pas** `auth` simplement pour vérifier si l'authentification est configurée — utilisez `info` à la place

### Fonctionnement

Recherche la configuration d'authentification de la spécification et exécute le flux d'authentification (échange de jetons, exécution de script, etc.) pour obtenir les identifiants actuels.

### Paramètres

| Paramètre | Type | Obligatoire | Description |
|-----------|------|----------|-------------|
| `specId` | string | Oui | Hachage MD5 de 32 caractères de la spécification |

### Réponse

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "headers": {
    "Authorization": "Bearer eyJhbGciOiJIUzI1NiIs...",
    "X-API-Key": "ma-clé-api"
  },
  "queryParams": {
    "api_key": "ma-clé-api"
  }
}
```

| Champ | Type | Description |
|-------|------|-------------|
| `token` | string | Valeur brute du jeton (jeton bearer, clé API, etc.) |
| `headers` | object | En-têtes HTTP à inclure dans les requêtes |
| `queryParams` | object | Paramètres de requête à inclure dans les requêtes |

### Nuances

- **Désactivé par défaut en production :** L'indicateur `--disable-llm-auth` (valeur par défaut : `true`) supprime complètement l'outil `auth` de la liste des outils MCP. Le LLM ne peut pas voir ni demander de jetons. Définissez `--disable-llm-auth=false` pour l'activer pour le débogage ou les jetons de courte durée.
- **`invoke` gère l'authentification automatiquement :** Vous n'avez pas besoin d'appeler `auth` avant `invoke`. Le service d'invocation obtient et applique automatiquement l'authentification correcte.
- **Prend en charge 9 méthodes d'authentification :** `none`, `basic`, `bearer`, `digest`, `hmac`, `oauth2-cc` (identifiants client), `oauth2-pwd` (mot de passe), `api-key`, `script`.
- Renvoie `auth_error` si la méthode d'authentification échoue (par exemple, point de terminaison de jeton OAuth2 inaccessible, échec d'exécution de script).

---

## info

### Objectif

Renvoie un résumé complet de l'environnement d'exécution de swag2mcp : version, chemin de l'espace de travail, spécifications actives, paramètres du client HTTP, configuration du transport MCP, méthodes d'authentification et état du mode de simulation.

### Quand l'utiliser

- Lorsque l'utilisateur pose des questions sur la configuration du système
- Lorsque vous devez vérifier les paramètres d'exécution (délai d'attente, limite de taille de réponse, transport)
- Lorsque vous devez savoir quelles méthodes d'authentification sont disponibles
- Lors du dépannage de problèmes de configuration

### Fonctionnement

Renvoie un instantané pré-calculé de l'état d'exécution. Aucun paramètre nécessaire.

### Paramètres

Aucun.

### Réponse

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
    "follow_redirects": true,
    "max_redirects": 10,
    "randomize": false,
    "proxy": null,
    "headers": {},
    "cookies": []
  },
  "mcp": {
    "transport": "stdio",
    "addr": ":8080",
    "path": "/mcp",
    "auth_enabled": false
  },
  "auth": {
    "methods": ["bearer", "api-key"]
  },
  "mock": {
    "enabled": false
  }
}
```

| Champ | Type | Description |
|-------|------|-------------|
| `version` | string | Version de swag2mcp |
| `workspace` | string | Chemin du répertoire de l'espace de travail |
| `uptime` | string | Temps de fonctionnement du serveur (lisible) |
| `specs` | object | Résumé des spécifications : total, actives, désactivées, collections, points de terminaison |
| `http_client` | object | Configuration du client HTTP |
| `http_client.max_response_size` | string | Taille maximale de réponse au format lisible (par exemple « 2 Ko ») |
| `mcp` | object | Configuration du serveur MCP |
| `auth` | object | Méthodes d'authentification disponibles |
| `mock` | object | État du serveur de simulation |

### Nuances

- `max_response_size` est affiché au format lisible (par exemple, « 1 Ko », « 2 Mo »)
- `uptime` est calculé à partir de l'heure de démarrage du serveur
- Les données sont un instantané pris au moment du démarrage — elles reflètent l'état au moment où le serveur MCP a démarré

---

## response_outline

### Objectif

Obtenir un résumé structurel de haut niveau d'un fichier de réponse JSON volumineux qui a été enregistré sur le disque par `invoke`. Il renvoie la forme des données — clés, types, longueurs de tableaux et indices de navigation — sans renvoyer les valeurs réelles.

### Quand l'utiliser

- Immédiatement après que `invoke` a renvoyé une `fileRef` (réponse trop volumineuse pour être en ligne)
- C'est la **première étape obligatoire** du flux de travail pour les grandes réponses

### Fonctionnement

Lit le fichier de réponse enregistré et analyse sa structure : type de premier niveau, clés, longueurs de tableaux, profondeur d'imbrication et indices de compression.

### Paramètres

| Paramètre | Type | Obligatoire | Description |
|-----------|------|----------|-------------|
| `path` | string | Oui | Chemin absolu depuis `fileRef.path` |
| `maxDepth` | int | Non | Profondeur de récursion maximale (par défaut : 3) |
| `maxArrayItems` | int | Non | Nombre d'éléments de tableau à inspecter (par défaut : 5) |

### Réponse

```json
{
  "outline": {
    "type": "object",
    "size": 1572864,
    "lineCount": 12500,
    "depth": 3,
    "structure": {
      "type": "object",
      "keys": ["data", "meta", "error"],
      "data": {
        "type": "array",
        "length": 500,
        "items": {
          "type": "object",
          "keys": ["id", "name", "status", "createdAt"]
        }
      }
    },
    "schemaHint": "objet avec 3 clés : data (tableau[500]), meta (objet), error (null)",
    "keys": ["data", "meta", "error"],
    "itemCount": 500,
    "itemType": "object",
    "compressionHints": [
      "response_compress(path, 'first_of_array', 'data')",
      "response_compress(path, 'sample_array', 'data', arrayHead=3, arrayTail=2)",
      "response_compress(path, 'keys_only', 'data')",
      "response_compress(path, 'select_keys', 'data', selectKeys=[id, name])"
    ],
    "navigationHints": {
      "paths": ["data", "meta", "error"],
      "arrays": [
        {"path": "data", "length": 500}
      ]
    }
  }
}
```

| Champ | Type | Description |
|-------|------|-------------|
| `type` | string | Type de premier niveau : « object » ou « array » |
| `size` | int | Taille du fichier en octets |
| `lineCount` | int | Nombre de lignes dans le fichier |
| `depth` | int | Profondeur d'imbrication maximale inspectée |
| `structure` | object | Structure récursive avec clés, types, longueurs de tableaux |
| `schemaHint` | string | Résumé en une ligne de la forme de premier niveau |
| `keys` | array | Clés de premier niveau (pour les objets) |
| `itemCount` | int | Longueur du tableau (pour les tableaux) |
| `compressionHints` | array | Appels `response_compress` suggérés avec paramètres |
| `navigationHints` | object | Chemins de premier niveau et tableaux avec longueurs |

### Nuances

- Renvoie `validation_failed` si le chemin est invalide ou ne se trouve pas dans le répertoire des réponses
- Renvoie `not_found` si le fichier n'existe pas
- Renvoie `validation_failed` si le fichier n'est pas un JSON valide
- Le champ `compressionHints` fournit des suggestions prêtes à l'emploi pour les appels `response_compress`

---

## response_compress

### Objectif

Réduire une valeur JSON dans un fichier de réponse enregistré afin qu'elle tienne dans la limite de taille de réponse et puisse être renvoyée au LLM en ligne. Plusieurs modes de compression vous permettent de choisir le bon compromis entre taille et information.

### Quand l'utiliser

- Après `response_outline` pour comprendre la structure
- Lorsque vous devez obtenir des données d'une grande réponse en ligne
- Lorsque `response_slice` est trop étroit et que vous avez besoin d'une vue plus large

### Fonctionnement

Lit le fichier de réponse enregistré, navigue jusqu'au chemin JSON spécifié, applique le mode de compression et renvoie le résultat compressé. Si le résultat dépasse encore la limite de taille, il est enregistré dans un nouveau fichier.

### Paramètres

| Paramètre | Type | Obligatoire | Description |
|-----------|------|----------|-------------|
| `path` | string | Oui | Chemin absolu depuis `fileRef.path` |
| `jsonPath` | string | Non | Chemin vers la valeur à compresser (par exemple `data` ou `data.0`) |
| `mode` | string | Oui | Mode de compression (voir le tableau ci-dessous) |
| `arrayHead` | int | Non | Éléments de début à conserver en mode `sample_array` (par défaut : 3) |
| `arrayTail` | int | Non | Éléments de fin à conserver en mode `sample_array` (par défaut : 2) |
| `stringLen` | int | Non | Longueur maximale des chaînes en mode `truncate_strings` (par défaut : 80) |
| `selectKeys` | array | Non | Clés à conserver en mode `select_keys` |

### Modes de compression

| Mode | Description | Meilleur pour |
|------|-------------|----------|
| `first_of_array` | Conserver uniquement le premier élément d'un tableau | Lorsque tous les éléments ont la même structure |
| `sample_array` | Conserver le début et la fin d'un tableau | Lorsque vous devez voir la plage de valeurs |
| `truncate_strings` | Raccourcir chaque chaîne à `stringLen` caractères | Lorsque les chaînes sont très longues mais que la structure est importante |
| `keys_only` | Remplacer les valeurs d'objet par leurs noms de type | Lorsque vous avez uniquement besoin de la structure |
| `select_keys` | Conserver uniquement les clés spécifiées dans chaque objet | Lorsque vous avez besoin de champs spécifiques de nombreux objets |

### Réponse

```json
{
  "body": [
    { "id": 1, "name": "Rex", "status": "available" },
    { "id": 2, "name": "Max", "status": "pending" }
  ],
  "hint": "Tableau compressé de 500 à 2 éléments en utilisant le mode first_of_array"
}
```

| Champ | Type | Description |
|-------|------|-------------|
| `body` | any | Valeur JSON compressée (présente lorsqu'elle est dans la limite de taille) |
| `fileRef` | object | Référence de fichier (présente si encore trop volumineuse) |
| `hint` | string | Explication de ce qui a été compressé |

### Nuances

- Si le résultat compressé dépasse encore `max_response_size`, il est enregistré dans un nouveau fichier et une `FileReference` est renvoyée
- Valeurs par défaut : `arrayHead=3`, `arrayTail=2`, `stringLen=80`
- Renvoie `validation_failed` pour un chemin invalide, un JSONPath invalide ou un fichier non JSON
- Renvoie `not_found` si le fichier n'existe pas ou si le JSONPath ne correspond pas

---

## response_slice

### Objectif

Extraire un fragment spécifique d'un fichier de réponse JSON enregistré par chemin JSON logique ou par plage de lignes. Contrairement à `response_compress`, cela vous donne les données brutes et non modifiées.

### Quand l'utiliser

- Lorsque vous avez besoin d'un élément ou d'une valeur spécifique d'une grande réponse
- Lorsque `response_compress` ne vous donne pas assez de détails
- Lorsque vous souhaitez naviguer dans une réponse étape par étape

### Fonctionnement

Lit le fichier de réponse enregistré et extrait un fragment par chemin JSON (par exemple, `data.3.name`) ou par plage de lignes (par exemple, `120-240`). Renvoie des indices de navigation pour parcourir les tableaux et les objets.

### Paramètres

| Paramètre | Type | Obligatoire | Description |
|-----------|------|----------|-------------|
| `path` | string | Oui | Chemin absolu depuis `fileRef.path` |
| `jsonPath` | string | Non | Chemin logique vers la valeur (par exemple `data.3.name`) |
| `line` | int | Non | Numéro de ligne (base 1) pour centrer le fragment |
| `range` | string | Non | Plage de lignes au format `début-fin` (par exemple `120-240`) |
| `around` | int | Non | Lignes à inclure autour de `line` (par défaut : 20) |

### Réponse

```json
{
  "slice": {
    "lines": [120, 130],
    "fragment": "{\n  \"id\": 1,\n  \"name\": \"Rex\"\n}",
    "value": {
      "id": 1,
      "name": "Rex"
    },
    "jsonPath": "data.0",
    "context": "object",
    "isComplete": true,
    "nextLine": 131,
    "prevLine": 119,
    "nextPath": "data.1",
    "prevPath": null
  }
}
```

| Champ | Type | Description |
|-------|------|-------------|
| `lines` | array | Plage de lignes (base 1) [début, fin] |
| `fragment` | string | Texte JSON brut (lorsqu'il est assez petit) |
| `value` | any | Valeur JSON extraite |
| `jsonPath` | string | Le chemin JSON utilisé |
| `context` | string | « object », « array » ou « value » |
| `isComplete` | bool | Vrai lorsque la valeur est un fragment JSON valide |
| `nextLine` | int | Ligne suivante suggérée pour la navigation par lignes |
| `prevLine` | int | Ligne précédente suggérée |
| `nextPath` | string | Chemin JSON suivant suggéré pour la navigation dans les tableaux |
| `prevPath` | string | Chemin JSON précédent suggéré |

### Nuances

- **Préférez `jsonPath` aux numéros de ligne** — les chemins JSON sont stables et descriptifs, les numéros de ligne changent si le fichier est régénéré
- Si le fragment extrait dépasse `max_response_size`, il est enregistré dans un nouveau fichier et une `FileReference` est renvoyée
- La valeur par défaut de `around` est de 20 lignes
- La réponse inclut `nextPath`/`prevPath` pour parcourir les tableaux et `nextLine`/`prevLine` pour la navigation par lignes
- Renvoie `validation_failed` pour un chemin invalide, un JSONPath invalide, une ligne/plage invalide ou un fichier non JSON
- Renvoie `not_found` si le fichier n'existe pas ou si le JSONPath ne correspond pas
