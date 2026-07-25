# Outils de découverte

Les outils de découverte permettent au LLM de naviguer dans la hiérarchie des spécifications : trouver toutes les spécifications, explorer une spécification pour voir ses collections et explorer les balises d'une collection. Commencez par `spec_list` pour voir les API disponibles, puis utilisez les ID pour approfondir.

---

## spec_list

### Objectif

Lister toutes les spécifications API enregistrées dans l'espace de travail. C'est le point de départ de toute session — le LLM l'appelle en premier pour découvrir les API disponibles.

### Quand l'utiliser

- Au début d'une session pour voir les API configurées
- Après avoir ajouté ou supprimé des spécifications pour actualiser la liste
- Lorsque vous avez besoin d'un ID de spécification pour d'autres outils

### Fonctionnement

Renvoie une liste de toutes les spécifications avec leur ID unique et leur nom de domaine. Aucun paramètre nécessaire.

### Paramètres

Aucun.

### Réponse

```json
{
  "specs": [
    {
      "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
      "domain": "meteo"
    },
    {
      "id": "b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7",
      "domain": "dadjoke"
    }
  ]
}
```

| Champ | Type | Description |
|-------|------|-------------|
| `id` | string | Hachage MD5 de 32 caractères, identifiant unique de la spécification |
| `domain` | string | Nom de domaine de la spécification (par exemple « meteo », « dadjoke ») |

### Nuances

- Renvoie uniquement `id` et `domain` — pour les détails complets (collections, balises), utilisez `spec_by_id`
- Tous les ID sont des chaînes hexadécimales MD5 de 32 caractères (`^[0-9a-f]{32}$`)
- Si aucune spécification n'est configurée, renvoie un tableau vide

---

## spec_by_id

### Objectif

Obtenir des informations détaillées sur une spécification spécifique : son domaine, toutes ses collections et leurs statistiques (nombre de balises, nombre de méthodes).

### Quand l'utiliser

- Après `spec_list` pour voir les collections d'une spécification
- Lorsque vous avez besoin des ID de collection pour une navigation ultérieure

### Fonctionnement

Prend un ID de spécification et renvoie les métadonnées de la spécification ainsi que toutes ses collections avec leurs décomptes.

### Paramètres

| Paramètre | Type | Obligatoire | Description |
|-----------|------|----------|-------------|
| `id` | string | Oui | Hachage MD5 de 32 caractères de la spécification |

### Réponse

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collections": [
    {
      "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
      "title": "Weather Forecast",
      "llmTitle": "Forecast API",
      "countTags": 3,
      "countMethods": 12
    }
  ]
}
```

| Champ | Type | Description |
|-------|------|-------------|
| `spec.id` | string | Identifiant de la spécification |
| `spec.domain` | string | Nom de domaine de la spécification |
| `collections[].id` | string | Identifiant de la collection |
| `collections[].title` | string | Titre lisible |
| `collections[].llmTitle` | string | Titre adapté au LLM (optionnel) |
| `collections[].countTags` | int | Nombre de balises dans la collection |
| `collections[].countMethods` | int | Nombre de méthodes HTTP dans la collection |

### Nuances

- Renvoie une erreur `not_found` si l'ID de spécification n'existe pas
- L'`id` doit être une chaîne hexadécimale MD5 valide de 32 caractères

---

## collection_by_spec

### Objectif

Lister toutes les collections d'une spécification spécifique. Similaire à `spec_by_id` mais renvoie uniquement la liste des collections sans les métadonnées supplémentaires de la spécification.

### Quand l'utiliser

- Lorsque vous avez déjà l'ID de spécification et que vous avez simplement besoin de la liste des collections
- Comme alternative plus légère à `spec_by_id`

### Paramètres

| Paramètre | Type | Obligatoire | Description |
|-----------|------|----------|-------------|
| `specId` | string | Oui | Hachage MD5 de 32 caractères de la spécification |

### Réponse

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collections": [
    {
      "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
      "title": "Weather Forecast",
      "llmTitle": "Forecast API",
      "countTags": 3,
      "countMethods": 12
    }
  ]
}
```

### Nuances

- Renvoie `not_found` si la spécification n'existe pas
- Mêmes données que `spec_by_id` mais sans l'enveloppe supplémentaire de la spécification

---

## collection_by_id

### Objectif

Obtenir des informations détaillées sur une collection spécifique : ses métadonnées, la spécification parente et toutes les balises de la collection.

### Quand l'utiliser

- Après `collection_by_spec` pour voir les balises d'une collection
- Lorsque vous avez besoin des ID de balise pour `tag_by_id` ou `endpoint_by_tag`

### Paramètres

| Paramètre | Type | Obligatoire | Description |
|-----------|------|----------|-------------|
| `id` | string | Oui | Hachage MD5 de 32 caractères de la collection |

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
  "tags": [
    {
      "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "title": "forecast",
      "countMethods": 5
    },
    {
      "id": "e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7",
      "title": "current",
      "countMethods": 7
    }
  ]
}
```

| Champ | Type | Description |
|-------|------|-------------|
| `spec` | object | Spécification parente (id, domain) |
| `collection` | object | Métadonnées de la collection (id, title, countMethods) |
| `tags[]` | array | Liste des balises avec id, title, countMethods |

### Nuances

- Renvoie `not_found` si l'ID de collection n'existe pas
- Les balises sont renvoyées avec leurs ID — utilisez `endpoint_by_tag(tagId)` pour voir les points de terminaison réels

---

## tag_by_spec

### Objectif

Lister toutes les balises d'une spécification entière, couvrant toutes les collections. Utile pour obtenir une vue d'ensemble de toutes les balises disponibles.

### Quand l'utiliser

- Lorsque vous souhaitez voir toutes les balises d'une spécification sans explorer chaque collection
- Lorsque vous ne savez pas quelle collection contient la balise dont vous avez besoin

### Paramètres

| Paramètre | Type | Obligatoire | Description |
|-----------|------|----------|-------------|
| `specId` | string | Oui | Hachage MD5 de 32 caractères de la spécification |

### Réponse

```json
{
  "tags": [
    {
      "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "title": "forecast",
      "countMethods": 5
    },
    {
      "id": "e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7",
      "title": "current",
      "countMethods": 7
    }
  ]
}
```

### Nuances

- Renvoie `not_found` si la spécification n'existe pas
- Les balises sont agrégées à partir de toutes les collections de la spécification

---

## tag_by_collection

### Objectif

Lister toutes les balises d'une collection spécifique. Contrairement à `tag_by_spec`, cela renvoie également les métadonnées de la spécification parente et de la collection.

### Quand l'utiliser

- Après `collection_by_id` pour confirmer la liste des balises
- Lorsque vous avez besoin du contexte complet (spécification + collection + balises)

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
  "tags": [
    {
      "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "title": "forecast",
      "countMethods": 5
    }
  ]
}
```

### Nuances

- Renvoie `not_found` si la collection n'existe pas
- Mêmes données de balises que `tag_by_spec` mais limitées à une seule collection

---

## tag_by_id

### Objectif

Obtenir des informations sur une seule balise : son ID, son titre et le nombre de méthodes qu'elle contient. Cela vous renseigne sur la balise elle-même — pour voir les points de terminaison réels, utilisez `endpoint_by_tag`.

### Quand l'utiliser

- Lorsque vous avez un ID de balise et que vous souhaitez confirmer son nom et sa taille
- Avant d'appeler `endpoint_by_tag` pour comprendre combien de points de terminaison attendre

### Paramètres

| Paramètre | Type | Obligatoire | Description |
|-----------|------|----------|-------------|
| `id` | string | Oui | Hachage MD5 de 32 caractères de la balise |

### Réponse

```json
{
  "tag": {
    "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
    "title": "forecast",
    "countMethods": 5
  }
}
```

| Champ | Type | Description |
|-------|------|-------------|
| `tag.id` | string | Identifiant de la balise |
| `tag.title` | string | Nom de balise lisible |
| `tag.countMethods` | int | Nombre de méthodes HTTP dans cette balise |

### Nuances

- Renvoie `not_found` si la balise n'existe pas
- Cet outil renvoie uniquement les métadonnées de la balise — utilisez `endpoint_by_tag` pour obtenir la liste réelle des points de terminaison
