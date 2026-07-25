# Outils de points de terminaison

Les outils de points de terminaison permettent au LLM de visualiser les points de terminaison API à différents niveaux de la hiérarchie : tous les points de terminaison d'une spécification, d'une collection, d'une balise ou un résumé d'un seul point de terminaison. Utilisez-les pour découvrir les opérations disponibles avant d'inspecter ou d'invoquer.

---

## endpoint_by_spec

### Objectif

Lister tous les points de terminaison d'une spécification entière, couvrant toutes les collections et balises. Renvoie la vue la plus complète — chaque point de terminaison de la spécification avec son contexte complet (balise, collection, spécification).

### Quand l'utiliser

- Lorsque vous souhaitez voir tous les points de terminaison disponibles dans une spécification
- Lorsque vous ne savez pas quelle collection ou balise contient le point de terminaison dont vous avez besoin
- Après `spec_by_id` pour obtenir la liste complète des points de terminaison

### Paramètres

| Paramètre | Type | Obligatoire | Description |
|-----------|------|----------|-------------|
| `specId` | string | Oui | Hachage MD5 de 32 caractères de la spécification |

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

| Champ | Type | Description |
|-------|------|-------------|
| `id` | string | Identifiant du point de terminaison |
| `tagId` | string | Identifiant de la balise parente |
| `tagName` | string | Nom de balise lisible |
| `collectionId` | string | Identifiant de la collection parente |
| `collectionTitle` | string | Titre de collection lisible |
| `specId` | string | Identifiant de la spécification parente |
| `specDomain` | string | Nom de domaine de la spécification |
| `method` | string | Méthode HTTP (GET, POST, PUT, DELETE, etc.) |
| `path` | string | Chemin API (par exemple /v1/forecast) |
| `summary` | string | Résumé lisible de ce que fait le point de terminaison |

### Nuances

- Renvoie `not_found` si la spécification n'existe pas
- Chaque point de terminaison inclut sa filiation complète (spécification → collection → balise) pour le contexte
- Pour un résumé rapide d'un seul point de terminaison, utilisez `endpoint_by_id`

---

## endpoint_by_collection

### Objectif

Lister tous les points de terminaison d'une collection spécifique, indépendamment de leur balise. Renvoie les points de terminaison groupés par collection avec les métadonnées de la spécification et de la collection.

### Quand l'utiliser

- Après `collection_by_id` pour voir tous les points de terminaison d'une collection
- Lorsque vous souhaitez explorer la surface API complète d'une collection

### Paramètres

| Paramètre | Type | Obligatoire | Description |
|-----------|------|----------|-------------|
| `collectionId` | string | Oui | Hachage MD5 de 32 caractères de la collection |

### Réponse

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collection": {
    "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
    "title": "Weather Forecast",
    "countMethods": 12
  },
  "endpoints": [
    {
      "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
      "tagId": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "tagName": "forecast",
      "method": "GET",
      "path": "/v1/forecast",
      "summary": "Obtenir les prévisions météorologiques pour un lieu"
    }
  ]
}
```

### Nuances

- Renvoie `not_found` si la collection n'existe pas
- Inclut les métadonnées de la spécification et de la collection pour le contexte
- Les points de terminaison de toutes les balises de la collection sont renvoyés ensemble

---

## endpoint_by_tag

### Objectif

Lister tous les points de terminaison regroupés sous une balise spécifique. C'est la vue la plus ciblée — les points de terminaison d'une balise dans une collection.

### Quand l'utiliser

- Après `tag_by_id` pour voir les points de terminaison réels d'une balise
- Lorsque vous connaissez la balise et souhaitez voir ses opérations

### Paramètres

| Paramètre | Type | Obligatoire | Description |
|-----------|------|----------|-------------|
| `tagId` | string | Oui | Hachage MD5 de 32 caractères de la balise |

### Réponse

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collection": {
    "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
    "title": "Weather Forecast",
    "countMethods": 12
  },
  "tag": {
    "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
    "title": "forecast",
    "countMethods": 5
  },
  "endpoints": [
    {
      "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
      "method": "GET",
      "path": "/v1/forecast",
      "summary": "Obtenir les prévisions météorologiques pour un lieu"
    }
  ]
}
```

### Nuances

- Renvoie `not_found` si la balise n'existe pas
- Inclut le contexte complet : spécification, collection et métadonnées de la balise
- Les points de terminaison sont limités à une seule balise dans une seule collection

---

## endpoint_by_id

### Objectif

Obtenir un résumé rapide d'un seul point de terminaison : méthode, chemin, résumé et état d'obsolescence. C'est un outil léger — pour l'objet complet de l'opération OpenAPI (paramètres, corps de la requête, schémas de réponse), utilisez `inspect`.

### Quand l'utiliser

- Lorsque vous avez un ID de point de terminaison et souhaitez un rappel rapide de ce qu'il fait
- Avant de décider d'appeler `inspect` pour obtenir les détails complets
- Lorsque vous devez confirmer la méthode et le chemin avant d'invoquer

### Paramètres

| Paramètre | Type | Obligatoire | Description |
|-----------|------|----------|-------------|
| `id` | string | Oui | Hachage MD5 de 32 caractères du point de terminaison |

### Réponse

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collection": {
    "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
    "title": "Weather Forecast",
    "countMethods": 12
  },
  "tag": {
    "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
    "title": "forecast",
    "countMethods": 5
  },
  "endpoint": {
    "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
    "method": "GET",
    "path": "/v1/forecast",
    "summary": "Obtenir les prévisions météorologiques pour un lieu"
  }
}
```

| Champ | Type | Description |
|-------|------|-------------|
| `endpoint.id` | string | Identifiant du point de terminaison |
| `endpoint.method` | string | Méthode HTTP |
| `endpoint.path` | string | Chemin API |
| `endpoint.summary` | string | Résumé lisible |

### Nuances

- Renvoie `not_found` si le point de terminaison n'existe pas
- Ceci est un **résumé rapide** — il ne renvoie pas les paramètres, le corps de la requête ou les schémas de réponse
- Pour les détails techniques complets (nécessaires avant `invoke`), utilisez `inspect`
