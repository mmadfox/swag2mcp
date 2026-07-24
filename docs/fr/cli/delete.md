# delete

## Objectif

Supprimer une **spec** (service API) ou une **collection** (fichier de spécification) de la configuration. C'est l'inverse de `add`.

## Quand l'utiliser

- Une API n'est plus nécessaire
- Vous voulez supprimer un fichier de spécification spécifique d'une spec
- Vous nettoyez votre espace de travail

## Syntaxe

```bash
swag2mcp delete spec [chemin]
swag2mcp delete collection [chemin]
```

## Arguments

| Argument | Position | Requis | Description |
|----------|----------|--------|-------------|
| `chemin` | 1 | Non | Répertoire de l'espace de travail. S'il est omis, résolution via les règles de résolution de chemin. |

## Drapeaux

Aucun. Les deux sous-commandes sont purement interactives.

## Comment cela fonctionne

### Supprimer une spec

Vous invite à sélectionner une spec dans une liste, puis demande confirmation avant de supprimer.

```bash
swag2mcp delete spec
```

### Supprimer une collection

Vous invite à sélectionner une spec, puis une collection dans cette spec, puis demande confirmation.

```bash
swag2mcp delete collection
```

## Trouver les IDs

Les invites interactives affichent des noms lisibles par l'humain, pas des IDs. Si vous avez besoin d'IDs pour référence :

```bash
# Lister toutes les specs avec leurs IDs
swag2mcp ls

# Lister les collections pour une spec spécifique
swag2mcp ls --tags
```

## Vérification post-commande

```bash
swag2mcp ls [chemin]
# La spec ou collection supprimée ne devrait plus apparaître
```

## Nuances

- **TTY requis :** Les deux commandes nécessitent un terminal interactif. Elles ne fonctionneront **pas** dans les pipelines CI/CD, les tâches cron ou les scripts non interactifs.
- **Pas de `--force` ou `--yes` :** Il n'y a aucun moyen d'ignorer l'invite de confirmation. C'est intentionnel pour éviter les suppressions accidentelles.
- **Auto-initialisation :** Si aucun fichier de configuration n'existe, `delete` exécute automatiquement l'assistant d'initialisation d'abord.
- **Pas de mode YAML :** Contrairement à `add`, il n'y a pas de drapeau `--yaml`. La suppression est toujours interactive.
