# Gestion de la taille des réponses

## Aperçu

Les réponses API peuvent être très volumineuses — parfois trop grandes pour tenir dans la fenêtre de contexte du LLM. swag2mcp gère automatiquement les tailles de réponse en sauvegardant les réponses trop grandes sur le disque et en fournissant des outils pour les explorer.

## Comment cela fonctionne

1. **Vous appelez `invoke`** — swag2mcp effectue la requête API
2. **Si la réponse est petite** (dans la limite) — elle est retournée en ligne au LLM
3. **Si la réponse est trop volumineuse** (dépasse la limite) — elle est sauvegardée dans `{espace-travail}/responses/` sous forme de fichier JSON. Le LLM reçoit une référence de fichier au lieu de la réponse complète

### Exemple : petite réponse (en ligne)

```json
{
  "statusCode": 200,
  "body": {
    "id": 1,
    "name": "Rex",
    "status": "available"
  }
}
```

### Exemple : grande réponse (référence de fichier)

```json
{
  "statusCode": 200,
  "fileRef": {
    "path": "/Users/utilisateur/.swag2mcp/responses/response_a1b2c3d4.json",
    "size": 1572864,
    "sizeHint": "1.5 Mo",
    "maxSizeHint": "2 Ko",
    "message": "La réponse dépasse la limite de 2 Ko et a été sauvegardée sur le disque.",
    "openCmd": "open /Users/utilisateur/.swag2mcp/responses/response_a1b2c3d4.json"
  }
}
```

## Configuration

```yaml
http_client:
  max_response_size: 1048576  # 1 Mo en octets
```

### max_response_size

- **Type :** `int` (octets)
- **Défaut :** `1048576` (1 Mo)
- **Plage :** 256 à 10 485 760 octets (10 Mo)
- **Effet :** Les réponses plus grandes que cela sont sauvegardées sur le disque au lieu d'être retournées en ligne
- **Quand augmenter :** API qui retournent de grands ensembles de données (rapports, journaux, analyses)
- **Quand diminuer :** Fenêtre de contexte LLM limitée, ou lorsque vous préférez un accès basé sur les fichiers

## Travailler avec de grandes réponses

Lorsque `invoke` retourne un `fileRef`, utilisez ces trois outils pour explorer les données :

### 1. response_outline — comprendre la structure

Obtenez un résumé structurel de la réponse : clés, types, longueurs de tableaux et indices de navigation.

```json
→ response_outline(path: "/chemin/vers/fichier.json")
← {
    "type": "object",
    "size": 1572864,
    "keys": ["data", "meta"],
    "itemCount": 500,
    "compressionHints": [
      "response_compress(path, 'first_of_array', 'data')",
      "response_compress(path, 'sample_array', 'data', arrayHead=3, arrayTail=2)"
    ]
  }
```

### 2. response_compress — obtenir une version plus petite

Compressez les données pour qu'elles tiennent en ligne. Plusieurs modes de compression vous permettent de choisir le bon compromis.

| Mode | Description | Meilleur pour |
|------|-------------|---------------|
| `first_of_array` | Conserver uniquement le premier élément d'un tableau | Quand tous les éléments ont la même structure |
| `sample_array` | Conserver le début (3) et la fin (2) d'un tableau | Quand vous devez voir la plage de valeurs |
| `truncate_strings` | Raccourcir chaque chaîne à N caractères | Quand les chaînes sont très longues |
| `keys_only` | Remplacer les valeurs par leurs noms de type | Quand vous avez seulement besoin de la structure |
| `select_keys` | Conserver uniquement les clés spécifiées | Quand vous avez besoin de champs spécifiques |

```json
→ response_compress(path: "/chemin/vers/fichier.json", mode: "first_of_array", jsonPath: "data")
← {
    "body": [{ "id": 1, "name": "Rex" }],
    "hint": "Tableau compressé de 500 à 1 élément en utilisant le mode first_of_array"
  }
```

### 3. response_slice — extraire un fragment spécifique

Obtenez un élément ou une valeur spécifique par chemin JSON ou plage de lignes.

```json
→ response_slice(path: "/chemin/vers/fichier.json", jsonPath: "data.0")
← {
    "slice": {
      "value": { "id": 1, "name": "Rex" },
      "jsonPath": "data.0",
      "nextPath": "data.1",
      "prevPath": null
    }
  }
```

## Flux de travail complet

```
1. invoke(point d'accès) → fileRef (la réponse fait 1,5 Mo)
2. response_outline(chemin) → structure : { data: Array(500) }
3. response_compress(chemin, mode: "first_of_array", jsonPath: "data") → premier élément
4. response_slice(chemin, jsonPath: "data.0") → détails complets du premier élément
5. response_slice(chemin, jsonPath: "data.1") → deuxième élément
```

## Nettoyage automatique

Lorsque le serveur MCP démarre (`swag2mcp mcp`), les fichiers de réponse de plus de 48 heures sont automatiquement supprimés. Vous pouvez également les nettoyer manuellement :

```bash
swag2mcp clean
```

## Notes importantes

- **La limite est en octets** — `1048576` = 1 Mo, `2097152` = 2 Mo, etc.
- **Les références de fichier incluent une commande d'ouverture** — sur macOS c'est `open`, sur Linux c'est `xdg-open`
- **Les fichiers de réponse sont nommés avec des suffixes aléatoires** — pas de conflits entre les appels simultanés
- **Le répertoire des réponses est créé automatiquement** — aucune configuration manuelle nécessaire
