# ls

## Objectif

Lister toutes les **specs** configurées et leurs **collections** dans un format lisible par l'humain. C'est le moyen principal d'inspecter les API disponibles dans votre espace de travail.

## Quand l'utiliser

- Vous voulez voir quelles API sont configurées
- Vous devez trouver un ID de spec ou de collection
- Vous voulez vérifier combien de points d'accès chaque collection a
- Vous voulez filtrer les specs par étiquettes

## Syntaxe

```bash
swag2mcp ls [chemin] [drapeaux]
```

## Arguments

| Argument | Position | Requis | Description |
|----------|----------|--------|-------------|
| `chemin` | 1 | Non | Répertoire de l'espace de travail. S'il est omis, résolution via les règles de résolution de chemin. |

## Drapeaux

| Drapeau | Raccourci | Type | Défaut | Description |
|---------|-----------|------|--------|-------------|
| `--tags` | `-t` | `string` | `""` | Filtrer les specs par étiquettes (séparées par des virgules) |

## Comment cela fonctionne

### Lister toutes les specs

Affiche chaque spec avec son domaine, ses collections et le nombre de points d'accès :

```bash
swag2mcp ls
```

Exemple de sortie :

```
Spécifications :
  dadjoke (https://icanhazdadjoke.com)
    blagues (3 points d'accès)
  meteo (https://meteo.swagger.io/v2)
    previsions (5 points d'accès)
    actuel (8 points d'accès)
  binance (https://api.binance.com)
    donnees-marche (12 points d'accès)
```

### Filtrer par étiquettes

Affiche uniquement les specs qui ont les étiquettes spécifiées :

```bash
swag2mcp ls --tags=public
swag2mcp ls --tags=public,interne
```

## Vérification post-commande

Utilisez `ls` après `add`, `delete`, `update` ou `import` pour confirmer que l'état de l'espace de travail correspond à vos attentes.

## Nuances

- **Auto-initialisation :** Si aucun fichier de configuration n'existe, `ls` exécute automatiquement l'assistant d'initialisation d'abord.
- **Filtrage par étiquettes :** Les étiquettes sont séparées par des virgules. Les specs correspondant à **n'importe laquelle** des étiquettes spécifiées sont affichées (logique OU).
- **Format de sortie :** La sortie est en texte brut, pas en JSON. Pour une sortie lisible par machine, utilisez `info`.
